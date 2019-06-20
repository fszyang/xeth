// Copyright © 2018-2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xeth

import (
	"net"
	"sync"

	"github.com/platinasystems/xeth/internal"
)

type DevKind uint8
type NetNs uint64

type IfInfo struct {
	Name    string
	IfIndex int32
	NetNs
	net.Flags
	DevKind
	net.HardwareAddr
}

type DevNew Xid
type DevDel Xid
type DevUp Xid
type DevDown Xid
type DevDump Xid
type DevUnreg Xid
type DevReg struct {
	Xid
	NetNs
}

var (
	ifinfos sync.Map

	poolIfInfo = sync.Pool{
		New: func() interface{} {
			return &IfInfo{
				HardwareAddr: make([]byte,
					internal.SizeofEthAddr,
					internal.SizeofEthAddr),
			}
		},
	}
)

func (ifinfo *IfInfo) Pool() {
	poolIfInfo.Put(ifinfo)
}

func (xid Xid) IfInfo() (ifinfo *IfInfo) {
	if v, ok := ifinfos.Load(xid); ok {
		ifinfo = v.(*IfInfo)
	}
	return
}

func (xid Xid) deleteIfInfo() {
	if ifinfo := xid.IfInfo(); ifinfo != nil {
		poolIfInfo.Put(ifinfo)
	}
	ifinfos.Delete(xid)
}

func (xid Xid) del() DevDel {
	xid.delFromXids()
	xid.deleteIfInfo()
	xid.deleteIPNets()
	xid.deleteUppers()
	xid.deleteLowers()
	xid.deleteEthtoolFlags()
	xid.deleteEthtoolSettings()
	xid.deleteSupportedLinkModes()
	xid.deleteAttrs()
	return DevDel(xid)
}

func (xid Xid) ifinfo(msg *internal.MsgIfInfo) (note interface{}) {
	ifinfo := xid.IfInfo()
	if ifinfo == nil {
		ifinfo = poolIfInfo.Get().(*IfInfo)
		xid.addToXids()
		note = DevNew(xid)
	} else {
		note = DevDump(xid)
	}
	ifinfo.Name = ""
	for i, c := range msg.Ifname[:] {
		if c == 0 {
			ifinfo.Name = string(msg.Ifname[:i])
		}
	}
	if len(ifinfo.Name) == 0 {
		ifinfo.Name = string(msg.Ifname[:])
	}
	ifinfo.IfIndex = msg.Ifindex
	ifinfo.NetNs = NetNs(msg.Net)
	ifinfo.Flags = net.Flags(msg.Flags)
	ifinfo.DevKind = DevKind(msg.Kind)
	copy(ifinfo.HardwareAddr, msg.Addr[:])
	ifinfos.Store(xid, ifinfo)
	return note
}

func (xid Xid) up() DevUp {
	if ifinfo := xid.IfInfo(); ifinfo != nil {
		ifinfo.Flags |= net.FlagUp
	}
	return DevUp(xid)
}

func (xid Xid) down() DevDown {
	if ifinfo := xid.IfInfo(); ifinfo != nil {
		ifinfo.Flags &^= net.FlagUp
	}
	return DevDown(xid)
}

func (xid Xid) reg(netns NetNs) *DevReg {
	if ifinfo := xid.IfInfo(); ifinfo != nil {
		ifinfo.NetNs = netns
	}
	return &DevReg{xid, netns}
}

func (xid Xid) unreg() DevUnreg {
	if ifinfo := xid.IfInfo(); ifinfo != nil {
		ifinfo.NetNs = 1
	}
	return DevUnreg(xid)
}
