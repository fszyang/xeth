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

func (xid Xid) RxIfInfo(msg *internal.MsgIfInfo) (note interface{}) {
	var m *sync.Map
	if v, ok := XidAttrMaps.Load(xid); ok {
		m = v.(*sync.Map)
	} else {
		m = new(sync.Map)
		XidAttrMaps.Store(xid, m)
	}
	if _, ok := m.Load(IfInfoNameAttr); ok {
		note = DevDump(xid)
	} else {
		note = DevNew(xid)
		name := make([]byte, internal.SizeofIfName)
		for i, c := range msg.Ifname[:] {
			if c == 0 {
				name = name[:i]
				break
			} else {
				name[i] = byte(c)
			}
		}
		m.Store(IfInfoNameAttr, string(name))
		m.Store(IfInfoDevKindAttr, DevKind(msg.Kind))
		ha := make(net.HardwareAddr, internal.SizeofEthAddr)
		copy(ha, msg.Addr[:])
		m.Store(IfInfoHardwareAddrAttr, ha)
	}
	m.Store(IfInfoIfIndexAttr, msg.Ifindex)
	m.Store(IfInfoNetNsAttr, NetNs(msg.Net))
	m.Store(IfInfoFlagsAttr, net.Flags(msg.Flags))
	return note
}

func (xid Xid) RxUp() DevUp {
	attrs := xid.Attrs()
	flags := attrs.IfInfoFlags()
	flags |= net.FlagUp
	attrs.Map().Store(IfInfoFlagsAttr, flags)
	return DevUp(xid)
}

func (xid Xid) RxDown() DevDown {
	attrs := xid.Attrs()
	flags := attrs.IfInfoFlags()
	flags &^= net.FlagUp
	attrs.Map().Store(IfInfoFlagsAttr, flags)
	return DevDown(xid)
}

func (xid Xid) RxReg(netns NetNs) *DevReg {
	xid.Map().Store(IfInfoNetNsAttr, netns)
	return &DevReg{xid, netns}
}

func (xid Xid) RxUnreg() DevUnreg {
	xid.Map().Store(IfInfoNetNsAttr, DefaultNetNs)
	return DevUnreg(xid)
}
