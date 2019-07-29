// Copyright © 2018-2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xeth

import (
	"net"
	"sync"
	"unsafe"

	"github.com/platinasystems/xeth/internal"
)

type FibEntryEvent uint8

type RtScope uint8
type RtTable uint32
type Rtn uint8
type RtnhFlags uint32

type FibEntry struct {
	net.IPNet
	NHs []*NH
	NetNs
	RtTable
	FibEntryEvent
	Rtn
	Tos uint8
}

type NH struct {
	net.IP
	Xid
	Weight int32
	RtnhFlags
	RtScope
}

var poolFibEntry = sync.Pool{
	New: func() interface{} {
		return &FibEntry{
			IPNet: net.IPNet{
				IP:   make([]byte, net.IPv6len, net.IPv6len),
				Mask: make([]byte, net.IPv6len, net.IPv6len),
			},
		}
	},
}
var poolNH = sync.Pool{
	New: func() interface{} {
		return &NH{
			IP: make([]byte, net.IPv6len, net.IPv6len),
		}
	},
}

func (nh *NH) Pool() {
	nh.IP = nh.IP[:cap(nh.IP)]
	poolNH.Put(nh)
}

func (fe *FibEntry) Pool() {
	for _, nh := range fe.NHs {
		nh.Pool()
	}
	fe.NHs = fe.NHs[:0]
	fe.IPNet.IP = fe.IPNet.IP[:cap(fe.IPNet.IP)]
	fe.IPNet.Mask = fe.IPNet.Mask[:cap(fe.IPNet.Mask)]
	poolFibEntry.Put(fe)
}

func fib4(msg *internal.MsgFibEntry) *FibEntry {
	fe := poolFibEntry.Get().(*FibEntry)
	fe.NetNs = NetNs(msg.Net)
	*(*uint32)(unsafe.Pointer(&fe.IPNet.IP[0])) = msg.Address
	*(*uint32)(unsafe.Pointer(&fe.IPNet.Mask[0])) = msg.Mask
	fe.IPNet.IP = fe.IPNet.IP[:net.IPv4len]
	fe.IPNet.Mask = fe.IPNet.Mask[:net.IPv4len]
	fe.FibEntryEvent = FibEntryEvent(msg.Event)
	fe.Rtn = Rtn(msg.Type)
	fe.RtTable = RtTable(msg.Table)
	fe.Tos = msg.Tos
	for _, nh := range msg.NextHops() {
		xid := fe.NetNs.Xid(nh.Ifindex)
		if xid == 0 {
			continue
		}
		fenh := poolNH.Get().(*NH)
		*(*uint32)(unsafe.Pointer(&fenh.IP[0])) = nh.Gw
		fenh.IP = fenh.IP[:net.IPv4len]
		fenh.Xid = xid
		fenh.Weight = nh.Weight
		fenh.RtnhFlags = RtnhFlags(nh.Flags)
		fenh.RtScope = RtScope(nh.Scope)
		fe.NHs = append(fe.NHs, fenh)
	}
	return fe
}

func fib6(msg *internal.MsgFib6Entry) *FibEntry {
	netns := NetNs(msg.Net)
	fe := poolFibEntry.Get().(*FibEntry)
	fe.NetNs = netns
	copy(fe.IPNet.IP, msg.Address[:])
	fe.IPNet.Mask = net.CIDRMask(int(msg.Length), net.IPv6len*8)
	fe.FibEntryEvent = FibEntryEvent(msg.Event)
	fe.Rtn = Rtn(msg.Type)
	fe.RtTable = RtTable(msg.Table)
	nhxid := netns.Xid(msg.Nh.Ifindex)
	if nhxid == 0 {
		return fe
	}
	nh := poolNH.Get().(*NH)
	copy(nh.IP, msg.Nh.Gw[:])
	nh.Xid = nhxid
	nh.Weight = msg.Nh.Weight
	nh.RtnhFlags = RtnhFlags(msg.Nh.Flags)
	fe.NHs = append(fe.NHs, nh)
	for _, sibling := range msg.Siblings() {
		sibxid := netns.Xid(sibling.Ifindex)
		if sibxid == 0 {
			continue
		}
		nh = poolNH.Get().(*NH)
		copy(nh.IP, sibling.Gw[:])
		nh.Xid = sibxid
		nh.Weight = sibling.Weight
		nh.RtnhFlags = RtnhFlags(sibling.Flags)
		fe.NHs = append(fe.NHs, nh)
	}
	return fe
}
