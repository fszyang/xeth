// Copyright © 2018-2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This package provides a sideband control interface to an XETH driver.
// Usage,
//	err := xeth.Init()
//	defer xeth.Close()
//	xeth.DumpIfInfo()
//	for buf := range xeth.RxCh {
//		if xeth.Class(buf) == xeth.ClassBreak {
//			break
//		}
//		msg := xeth.Parse(buf)
//		...
//		xeth.Pool(msg)
//	}
//	...
//	xeth.DumpFib()
//	for buf := range xeth.RxCh {
//		if xeth.Class(buf) == xeth.ClassBreak {
//			break
//		}
//		msg := xeth.Parse(buf)
//		...
//		xeth.Pool(msg)
//	}
//	...
//	var wg sync.WaitGroup
//	wg.Add(1)
//	go func() {
//		defer wg.Done()
//		for buf := range xeth.RxCh {
//			msg := xeth.Parse(buf)
//			...
//			xeth.Pool(msg)
//		}
//	}()
//	...
//	wg.Wait()
package xeth

import (
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"

	"github.com/platinasystems/xeth/internal"
)

//go:generate sh -c "go tool cgo -godefs godefs.go > godefed.go"
//go:generate sh -c "go tool cgo -godefs internal/godefs.go > internal/godefed.go"

const netname = "unixpacket"

const (
	ClassUnknown = iota
	ClassBreak
	ClassInterface
	ClassAddress
	ClassFib
	ClassNeighbor
)

type Counter uint64
type SigErr struct{ os.Signal }

type Break struct{}

type Buffer interface{ buffer }

type pooler interface {
	Pool()
}

var (
	Cloned  Counter       // cloned received messages
	Parsed  Counter       // messages parsed by user
	Dropped Counter       // messages that overflowed transmit channel
	Sent    Counter       // messages and exception frames sent to driver
	RxCh    <-chan Buffer // cloned msgs received from driver
)

var (
	sock *net.UnixConn

	loch chan buffer // low priority, leaky-bucket tx channel
	hich chan buffer // high priority, unbuffered, no-drop tx channel
	rxch chan Buffer // deep channel of cloned msgs received from driver

	sigch chan os.Signal

	wg sync.WaitGroup

	sigErr error // signal that terminated service routines
	rxErr  error // error that stopped the rx service
	txErr  error // error that stopped the tx service
)

// Connect @xeth socket and run channel service routines.
func Init() error {
	addr, err := net.ResolveUnixAddr(netname, "@xeth")
	if err != nil {
		return err
	}

	for {
		sock, err = net.DialUnix(netname, nil, addr)
		if err == nil {
			break
		}
		if !isEAGAIN(err) {
			return err
		}
	}

	loch = make(chan buffer, 4)
	hich = make(chan buffer)
	rxch = make(chan Buffer, 1024)
	sigch = make(chan os.Signal, 1)

	RxCh = rxch

	wg.Add(2)
	go gorx()
	go gotx()

	return nil
}

func Class(buf Buffer) int {
	switch kind(buf) {
	case internal.MsgKindBreak:
		return ClassBreak
	case internal.MsgKindChangeUpperXid,
		internal.MsgKindEthtoolFlags,
		internal.MsgKindEthtoolLinkModesSupported,
		internal.MsgKindEthtoolLinkModesAdvertising,
		internal.MsgKindEthtoolLinkModesLPAdvertising,
		internal.MsgKindEthtoolSettings,
		internal.MsgKindIfInfo:
		return ClassInterface
	case internal.MsgKindIfa, internal.MsgKindIfa6:
		return ClassAddress
	case internal.MsgKindFibEntry, internal.MsgKindFib6Entry:
		return ClassFib
	case internal.MsgKindNeighUpdate:
		return ClassNeighbor
	default:
		return ClassUnknown
	}
}

// Close socket and shutdown service routines
func Close() error {
	if sock == nil {
		return os.ErrClosed
	}
	sigch <- syscall.SIGKILL
	close(loch)
	close(hich)
	wg.Wait()
	t := sock
	sock = nil
	f, err := t.File()
	if err == nil {
		syscall.Shutdown(int(f.Fd()), SHUT_RDWR)
	}
	t.Close()
	if err == nil {
		if err = rxErr; err == nil {
			if err = txErr; err == nil {
				err = sigErr
			}
		}
	}
	return err
}

// request fib dump
func DumpFib() {
	buf := newBuffer(internal.SizeofMsgDumpFibInfo)
	msg := (*internal.MsgHeader)(buf.pointer())
	msg.Set(internal.MsgKindDumpFibInfo)
	hich <- buf
}

// request ifinfo dump
func DumpIfInfo() {
	buf := newBuffer(internal.SizeofMsgDumpIfInfo)
	msg := (*internal.MsgHeader)(buf.pointer())
	msg.Set(internal.MsgKindDumpIfInfo)
	hich <- buf
}

// Send an exception frame to driver through leaky-bucket channel.
func ExceptionFrame(b []byte) {
	queue(cloneBuffer(b))
}

// true if err is type *SigErr
func IsSignal(err error) bool {
	_, match := err.(*SigErr)
	return match
}

// parse driver message and cache ifinfo in xid maps.
func Parse(buf Buffer) interface{} {
	defer Parsed.inc()
	defer buf.pool()
	switch kind(buf) {
	case internal.MsgKindBreak:
		return Break{}
	case internal.MsgKindChangeUpperXid:
		msg := (*internal.MsgChangeUpperXid)(buf.pointer())
		lower := Xid(msg.Lower)
		upper := Xid(msg.Upper)
		if msg.Linking != 0 {
			return lower.Join(upper)
		} else {
			return lower.Quit(upper)
		}
	case internal.MsgKindEthtoolFlags:
		msg := (*internal.MsgEthtoolFlags)(buf.pointer())
		xid := Xid(msg.Xid)
		return xid.ethtoolFlags(msg.Flags)
	case internal.MsgKindEthtoolLinkModesSupported:
		msg := (*internal.MsgEthtoolLinkModes)(buf.pointer())
		xid := Xid(msg.Xid)
		return xid.supportedLinkModes(msg.Modes())
	case internal.MsgKindEthtoolLinkModesAdvertising:
		msg := (*internal.MsgEthtoolLinkModes)(buf.pointer())
		xid := Xid(msg.Xid)
		return xid.advertisingLinkModes(msg.Modes())
	case internal.MsgKindEthtoolLinkModesLPAdvertising:
		msg := (*internal.MsgEthtoolLinkModes)(buf.pointer())
		xid := Xid(msg.Xid)
		return xid.lpadvertisingLinkModes(msg.Modes())
	case internal.MsgKindEthtoolSettings:
		msg := (*internal.MsgEthtoolSettings)(buf.pointer())
		xid := Xid(msg.Xid)
		return xid.ethtoolSettings(msg)
	case internal.MsgKindFibEntry:
		msg := (*internal.MsgFibEntry)(buf.pointer())
		return fib4(msg)
	case internal.MsgKindFib6Entry:
		msg := (*internal.MsgFib6Entry)(buf.pointer())
		return fib6(msg)
	case internal.MsgKindIfa:
		msg := (*internal.MsgIfa)(buf.pointer())
		xid := Xid(msg.Xid)
		if msg.Event == internal.IFA_ADD {
			return xid.addIP(msg.Address, msg.Mask)
		} else {
			return xid.delIP(msg.Address, msg.Mask)
		}
	case internal.MsgKindIfa6:
		msg := (*internal.MsgIfa6)(buf.pointer())
		xid := Xid(msg.Xid)
		if msg.Event == internal.IFA_ADD {
			addr := []byte(msg.Address[:])
			length := int(msg.Length)
			return xid.addIP6(addr, length)
		} else {
			return xid.delIP6(msg.Address[:])
		}
	case internal.MsgKindIfInfo:
		msg := (*internal.MsgIfInfo)(buf.pointer())
		xid := Xid(msg.Xid)
		switch msg.Reason {
		case internal.IfInfoReasonNew:
			return xid.ifinfo(msg)
		case internal.IfInfoReasonDump:
			return xid.ifinfo(msg)
		case internal.IfInfoReasonDel:
			return xid.del()
		case internal.IfInfoReasonUp:
			return xid.up()
		case internal.IfInfoReasonDown:
			return xid.down()
		case internal.IfInfoReasonReg:
			netns := NetNs(msg.Net)
			return xid.reg(netns)
		case internal.IfInfoReasonUnreg:
			return xid.unreg()
		}
	case internal.MsgKindNeighUpdate:
		msg := (*internal.MsgNeighUpdate)(buf.pointer())
		return neighbor(msg)
	}
	return nil
}

func Pool(msg interface{}) {
	if method, found := msg.(pooler); found {
		method.Pool()
	}
}

// Send carrier change to driver through hi-priority channel.
func SetCarrier(xid Xid, on bool) {
	buf := newBuffer(internal.SizeofMsgCarrier)
	msg := (*internal.MsgCarrier)(buf.pointer())
	msg.Header.Set(internal.MsgKindCarrier)
	msg.Xid = uint32(xid)
	if on {
		msg.Flag = internal.CarrierOn
	} else {
		msg.Flag = internal.CarrierOff
	}
	hich <- buf
}

// Send ethtool stat change to driver through leaky-bucket channel.
func SetEthtoolStat(xid Xid, stat uint32, n uint64) {
	setStat(internal.MsgKindEthtoolStat, xid, stat, n)
}

// Send link stat change to driver through leaky-bucket channel.
func SetLinkStat(xid Xid, stat uint32, n uint64) {
	setStat(internal.MsgKindLinkStat, xid, stat, n)
}

// Send speed change to driver through hi-priority channel.
func SetSpeed(xid Xid, mbps uint32) {
	buf := newBuffer(internal.SizeofMsgSpeed)
	msg := (*internal.MsgSpeed)(buf.pointer())
	msg.Header.Set(internal.MsgKindSpeed)
	msg.Xid = uint32(xid)
	msg.Mbps = mbps
	hich <- buf
}

func gorx() {
	const minrxto = 10 * time.Millisecond
	const maxrxto = 320 * time.Millisecond

	rxto := minrxto
	rxbuf := make([]byte, PageSize, PageSize)
	rxoob := make([]byte, PageSize, PageSize)
	ptr := unsafe.Pointer(&rxbuf[0])
	h := (*internal.MsgHeader)(ptr)

	signal.Notify(sigch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP,
		syscall.SIGQUIT)

	defer func() {
		rxbuf = rxbuf[:0]
		rxoob = rxoob[:0]
		signal.Stop(sigch)
		close(rxch)
		wg.Done()
	}()

	for !signaled() {
		rxErr = sock.SetReadDeadline(time.Now().Add(rxto))
		if rxErr != nil {
			break
		}
		n, noob, flags, addr, err := sock.ReadMsgUnix(rxbuf, rxoob)
		_ = noob
		_ = flags
		_ = addr
		if signaled() {
			break
		}
		if n == 0 || isTimeout(err) {
			if rxto < maxrxto {
				rxto *= 2
			}
		} else if err != nil {
			e, ok := err.(*os.SyscallError)
			if !ok || e.Err.Error() != "EOF" {
				rxErr = err
			}
			break
		} else if rxErr = h.Validate(rxbuf[:n]); rxErr != nil {
			break
		} else {
			rxto = minrxto
			rxch <- cloneBuffer(rxbuf[:n])
			Cloned.inc()
		}
	}
}

func gotx() {
	defer wg.Done()
	for txErr == nil {
		select {
		case buf, ok := <-hich:
			if ok {
				txErr = tx(buf, 0)
			} else {
				return
			}
		case buf, ok := <-loch:
			if ok {
				txErr = tx(buf, 10*time.Millisecond)
			} else {
				return
			}
		}
		if txErr == nil {
			Sent.inc()
		}
	}
}

func kind(buf buffer) uint8 {
	return (*internal.MsgHeader)(buf.pointer()).Kind
}

// Send through low-priority, leaky-bucket.
func queue(buf buffer) {
	select {
	case loch <- buf:
	default:
		buf.pool()
		Dropped.inc()
	}
}

func setStat(kind uint8, xid Xid, stat uint32, n uint64) {
	buf := newBuffer(internal.SizeofMsgStat)
	msg := (*internal.MsgStat)(buf.pointer())
	msg.Header.Set(kind)
	msg.Xid = uint32(xid)
	msg.Index = stat
	msg.Count = n
	queue(buf)
}

func signaled() bool {
	select {
	case sig := <-sigch:
		if sig != syscall.SIGKILL {
			sigErr = &SigErr{sig}
		}
		return true
	default:
		return false
	}
}

func tx(buf buffer, timeout time.Duration) error {
	var oob []byte
	var dl time.Time
	defer buf.pool()
	if sock == nil {
		return io.EOF
	}
	if timeout != time.Duration(0) {
		dl = time.Now().Add(timeout)
	}
	err := sock.SetWriteDeadline(dl)
	if err != nil {
		return err
	}
	_, _, err = sock.WriteMsgUnix(buf.bytes(), oob, nil)
	return err
}

func isEAGAIN(err error) bool {
	if err != nil {
		if operr, ok := err.(*net.OpError); ok {
			if oserr, ok := operr.Err.(*os.SyscallError); ok {
				if oserr.Err == syscall.EAGAIN {
					return true
				}
			}
		}
	}
	return false
}

func isTimeout(err error) bool {
	if err != nil {
		if op, ok := err.(*net.OpError); ok {
			return op.Timeout()
		}
	}
	return false
}

func (count *Counter) Count() uint64 {
	return atomic.LoadUint64((*uint64)(count))
}

func (count *Counter) Reset() {
	atomic.StoreUint64((*uint64)(count), 0)
}

func (count *Counter) inc() {
	atomic.AddUint64((*uint64)(count), 1)
}

func (sig *SigErr) Error() string {
	return sig.String()
}
