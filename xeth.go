/* XETH driver sideband control.
 *
 * Copyright(c) 2018 Platina Systems, Inc.
 *
 * This program is free software; you can redistribute it and/or modify it
 * under the terms and conditions of the GNU General Public License,
 * version 2, as published by the Free Software Foundation.
 *
 * This program is distributed in the hope it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
 * FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for
 * more details.
 *
 * You should have received a copy of the GNU General Public License along with
 * this program; if not, write to the Free Software Foundation, Inc.,
 * 51 Franklin St - Fifth Floor, Boston, MA 02110-1301 USA.
 *
 * The full GNU General Public License is included in this distribution in
 * the file called "COPYING".
 *
 * Contact Information:
 * sw@platina.com
 * Platina Systems, 3180 Del La Cruz Blvd, Santa Clara, CA 95054
 */
package xeth

import (
	"fmt"
	"io"
	"net"
	"os"
	"syscall"
	"time"
)

//go:generate sh -c "go tool cgo -godefs godefs.go > godefed.go"

const netname = "unixpacket"

var (
	Count struct {
		Tx struct {
			Sent, Dropped uint64
		}
	}
	// Receive message channel feed from sock by gorx
	RxCh <-chan []byte

	sb struct {
		name string
		addr *net.UnixAddr
		sock *net.UnixConn

		rxch chan []byte
		txch chan []byte
	}
)

// Connect to @xeth socket and run channel service routines
// driver :: XETH driver name (e.g. "platina-mk1")
func Start(driver string) error {
	var err error
	sb.name = driver
	sb.addr, err = net.ResolveUnixAddr(netname, "@xeth")
	if err != nil {
		return err
	}
	for {
		sb.sock, err = net.DialUnix(netname, nil, sb.addr)
		if err == nil {
			break
		}
		if !isEAGAIN(err) {
			return err
		}
	}
	sb.rxch = make(chan []byte, 4)
	sb.txch = make(chan []byte, 4)
	RxCh = sb.rxch
	go gorx()
	go gotx()

	// load Interface cache
	DumpIfInfo()
	UntilBreak(func(buf []byte) error {
		return nil
	})

	return nil
}

// Close @xeth socket and shutdown service routines
func Stop() {
	const (
		SHUT_RD = iota
		SHUT_WR
		SHUT_RDWR
	)
	if sb.sock == nil {
		return
	}
	close(sb.txch)
	sock := sb.sock
	sb.sock = nil
	if f, err := sock.File(); err == nil {
		syscall.Shutdown(int(f.Fd()), SHUT_RDWR)
	}
	sock.Close()
	Range(func(xid Xid, xeth *Xeth) bool {
		xeth.Uppers.Range(func(k, v interface{}) bool {
			xeth.Uppers.Delete(k.(Xid))
			return true
		})
		xeth.Lowers.Range(func(k, v interface{}) bool {
			xeth.Lowers.Delete(k.(Xid))
			return true
		})
		Delete(xid)
		return true
	})
}

// Return driver name (e.g. "platina-mk1")
func String() string { return sb.name }

// Send carrier state change message
func Carrier(xid Xid, on bool) error {
	buf := Pool.Get(SizeofMsgCarrier)
	defer Pool.Put(buf)
	ToMsgCarrier(buf).Set(xid, on)
	return tx(buf, 0)
}

// Send DumpFib request
func DumpFib() error {
	buf := Pool.Get(SizeofMsgDumpFibInfo)
	defer Pool.Put(buf)
	ToMsgHeader(buf).Set(MsgKindDumpFibInfo)
	return tx(buf, 0)
}

// Send DumpIfInfo request then flush RxCh until break to cache ifinfos
func CacheIfInfo() {
	if err := DumpIfInfo(); err == nil {
		UntilBreak(func(buf []byte) error { return nil })
	}
}

// Send DumpIfinfo request
func DumpIfInfo() error {
	buf := Pool.Get(SizeofMsgDumpIfInfo)
	defer Pool.Put(buf)
	ToMsgHeader(buf).Set(MsgKindDumpIfInfo)
	return tx(buf, 0)
}

// Send stat update message
func SetLinkStat(xid Xid, stat LinkStat, count uint64) error {
	buf := Pool.Get(SizeofMsgStat)
	defer Pool.Put(buf)
	ToMsgLinkStat(buf).Set(xid, stat, count)
	return tx(buf, 10*time.Millisecond)
}

func SetEthtoolStat(xid Xid, stat EthtoolStat, count uint64) error {
	buf := Pool.Get(SizeofMsgStat)
	defer Pool.Put(buf)
	ToMsgEthtoolStat(buf).Set(xid, stat, count)
	return tx(buf, 10*time.Millisecond)
}

// Send speed change message
func Speed(xid Xid, mbps Mbps) error {
	buf := Pool.Get(SizeofMsgSpeed)
	defer Pool.Put(buf)
	ToMsgSpeed(buf).Set(xid, mbps)
	return tx(buf, 0)
}

// Send through leaky bucket
func Tx(buf []byte) {
	msg := Pool.Get(len(buf))
	copy(msg, buf)
	select {
	case sb.txch <- msg:
		Count.Tx.Sent++
	default:
		Count.Tx.Dropped++
		Pool.Put(msg)
	}
}

func UntilBreak(f func([]byte) error) error {
	for buf := range RxCh {
		if KindOf(buf) == MsgKindBreak {
			Pool.Put(buf)
			break
		}
		err := f(buf)
		Pool.Put(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func UntilSig(sig <-chan os.Signal, f func([]byte) error) error {
	for {
		select {
		case <-sig:
			return nil
		case buf, ok := <-RxCh:
			if !ok {
				return nil
			}
			err := f(buf)
			Pool.Put(buf)
			if err != nil {
				return err
			}
		}
	}
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

func gorx() {
	const minrxto = 10 * time.Millisecond
	const maxrxto = 320 * time.Millisecond
	rxto := minrxto
	rxbuf := Pool.Get(PageSize)
	defer Pool.Put(rxbuf)
	rxoob := Pool.Get(PageSize)
	defer Pool.Put(rxoob)
	defer close(sb.rxch)
	for sb.sock != nil {
		err := sb.sock.SetReadDeadline(time.Now().Add(rxto))
		if err != nil {
			fmt.Fprintln(os.Stderr, "xeth set rx deadline", err)
			break
		}
		n, noob, flags, addr, err :=
			sb.sock.ReadMsgUnix(rxbuf, rxoob)
		_ = noob
		_ = flags
		_ = addr
		if n == 0 || isTimeout(err) {
			if rxto < maxrxto {
				rxto *= 2
			}
		} else if err == nil {
			rxto = minrxto
			kind := KindOf(rxbuf[:n])
			if err = kind.validate(rxbuf[:n]); err != nil {
				fmt.Fprintln(os.Stderr, "xeth rx", err)
				break
			}
			kind.update(rxbuf[:n])
			msg := Pool.Get(n)
			copy(msg, rxbuf[:n])
			sb.rxch <- msg
		} else {
			e, ok := err.(*os.SyscallError)
			if !ok || e.Err.Error() != "EOF" {
				fmt.Fprintln(os.Stderr, "xeth rx", err)
			}
			break
		}
	}
}

func gotx() {
	for msg := range sb.txch {
		tx(msg, 10*time.Millisecond)
		Pool.Put(msg)
	}
}

func tx(buf []byte, timeout time.Duration) error {
	var oob []byte
	var dl time.Time
	if sb.sock == nil {
		return io.EOF
	}
	if timeout != time.Duration(0) {
		dl = time.Now().Add(timeout)
	}
	err := sb.sock.SetWriteDeadline(dl)
	if err != nil {
		return err
	}
	_, _, err = sb.sock.WriteMsgUnix(buf, oob, nil)
	return err
}
