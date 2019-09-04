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
	attrs := MayMakeLinkAttrs(xid)
	if len(attrs.IfInfoName()) > 0 {
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
	attrs := LinkAttrs(xid)
	flags := attrs.IfInfoFlags()
	flags |= net.FlagUp
	attrs.IfInfoFlags(flags)
	return DevUp(xid)
}

func (xid Xid) RxDown() DevDown {
	attrs := LinkAttrs(xid)
	flags := attrs.IfInfoFlags()
	flags &^= net.FlagUp
	attrs.IfInfoFlags(flags)
	return DevDown(xid)
}

func (xid Xid) RxReg(netns NetNs) *DevReg {
	xidattrs := LinkAttrs(xid)
	ifindex := xidattrs.IfInfoIfIndex()
	if netns != DefaultNetNs {
		DefaultNetNs.Xid(ifindex, 0)
		netns.Xid(ifindex, xid)
		xidattrs.IfInfoNetNs(netns)
	} else {
		DefaultNetNs.Xid(ifindex, xid)
	}
	return &DevReg{xid, netns}
}

func (xid Xid) RxUnreg() DevUnreg {
	xidattrs := LinkAttrs(xid)
	ifindex := xidattrs.IfInfoIfIndex()
	oldns := xidattrs.IfInfoNetNs()
	oldns.Xid(ifindex, 0)
	DefaultNetNs.Xid(ifindex, xid)
	xidattrs.IfInfoNetNs(DefaultNetNs)
	return DevUnreg(xid)
}
