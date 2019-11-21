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
	l := mayMakeLinkOf(xid)
	if len(l.IfInfoName()) > 0 {
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
		l.IfInfoName(string(name))
		l.IfInfoDevKind(DevKind(msg.Kind))
		ha := make(net.HardwareAddr, internal.SizeofEthAddr)
		copy(ha, msg.Addr[:])
		l.IfInfoHardwareAddr(ha)
	}
	l.IfInfoIfIndex(msg.Ifindex)
	l.IfInfoNetNs(NetNs(msg.Net))
	l.IfInfoFlags(net.Flags(msg.Flags))
	return note
}

func (xid Xid) RxUp() (up DevUp) {
	l := expectLinkOf(xid, "admin-up")
	if l == nil {
		return
	}
	flags := l.IfInfoFlags()
	flags |= net.FlagUp
	l.IfInfoFlags(flags)
	up = DevUp(xid)
	return
}

func (xid Xid) RxDown() (down DevDown) {
	l := expectLinkOf(xid, "admin-down")
	if l == nil {
		return
	}
	flags := l.IfInfoFlags()
	flags &^= net.FlagUp
	l.IfInfoFlags(flags)
	down = DevDown(xid)
	return
}

func (xid Xid) RxReg(netns NetNs) (reg *DevReg) {
	l := expectLinkOf(xid, "netns-reg")
	if l == nil {
		return
	}
	ifindex := l.IfInfoIfIndex()
	if netns != DefaultNetNs {
		DefaultNetNs.Xid(ifindex, 0)
		netns.Xid(ifindex, xid)
		l.IfInfoNetNs(netns)
	} else {
		DefaultNetNs.Xid(ifindex, xid)
	}
	reg = &DevReg{xid, netns}
	return
}

func (xid Xid) RxUnreg() (unreg DevUnreg) {
	l := expectLinkOf(xid, "netns-reg")
	if l == nil {
		return
	}
	ifindex := l.IfInfoIfIndex()
	oldns := l.IfInfoNetNs()
	oldns.Xid(ifindex, 0)
	DefaultNetNs.Xid(ifindex, xid)
	l.IfInfoNetNs(DefaultNetNs)
	unreg = DevUnreg(xid)
	return
}
