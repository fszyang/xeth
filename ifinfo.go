// Copyright © 2018-2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xeth

import (
	"net"

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
	attrs := xid.attrs()
	if _, ok := attrs.Map().Load(IfInfoNameXidAttr); ok {
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
		attrs.IfInfoName(string(name))
		attrs.IfInfoDevKind(DevKind(msg.Kind))
		ha := make(net.HardwareAddr, internal.SizeofEthAddr)
		copy(ha, msg.Addr[:])
		attrs.IfInfoHardwareAddr(ha)
	}
	attrs.IfInfoIfIndex(msg.Ifindex)
	attrs.IfInfoNetNs(NetNs(msg.Net))
	attrs.IfInfoFlags(net.Flags(msg.Flags))
	return note
}

func (xid Xid) RxUp() DevUp {
	attrs := xid.Attrs()
	flags := attrs.IfInfoFlags()
	flags |= net.FlagUp
	attrs.Map().Store(IfInfoFlagsXidAttr, flags)
	return DevUp(xid)
}

func (xid Xid) RxDown() DevDown {
	attrs := xid.Attrs()
	flags := attrs.IfInfoFlags()
	flags &^= net.FlagUp
	attrs.Map().Store(IfInfoFlagsXidAttr, flags)
	return DevDown(xid)
}

func (xid Xid) RxReg(netns NetNs) *DevReg {
	xidattrs := xid.Attrs()
	ifindex := xidattrs.IfInfoIfIndex()
	if netns != DefaultNetNs {
		DefaultNetNs.XidOfIfIndexMap().Delete(ifindex)
		netns.XidOfIfIndexMap().Store(ifindex, xid)
		xidattrs.IfInfoNetNs(netns)
	} else {
		DefaultNetNs.XidOfIfIndexMap().Store(ifindex, xid)
	}
	return &DevReg{xid, netns}
}

func (xid Xid) RxUnreg() DevUnreg {
	xidattrs := xid.Attrs()
	ifindex := xidattrs.IfInfoIfIndex()
	oldns := xidattrs.IfInfoNetNs()
	oldns.XidOfIfIndexMap().Delete(ifindex)
	DefaultNetNs.XidOfIfIndexMap().Store(ifindex, xid)
	xidattrs.IfInfoNetNs(DefaultNetNs)
	return DevUnreg(xid)
}
