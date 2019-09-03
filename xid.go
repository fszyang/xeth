// Copyright © 2018-2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xeth

import (
	"fmt"
	"net"
	"regexp"
	"sync"
)

type Xid uint32
type Xids []Xid
type XidAttr uint8
type XidAttrs struct {
	sync.Map
}

const (
	EthtoolAutoNegXidAttr XidAttr = iota
	EthtoolDevPortXidAttr
	EthtoolDuplexXidAttr
	EthtoolFlagsXidAttr
	EthtoolSpeedXidAttr
	IPNetsXidAttr
	IfInfoNameXidAttr
	IfInfoIfIndexXidAttr
	IfInfoNetNsXidAttr
	IfInfoFlagsXidAttr
	IfInfoDevKindXidAttr
	IfInfoHardwareAddrXidAttr
	LinkModesAdvertisingXidAttr
	LinkModesLPAdvertisingXidAttr
	LinkModesSupportedXidAttr
	LinkUpXidAttr
	LowersXidAttr
	StatNamesXidAttr
	StatsXidAttr
	UppersXidAttr
)

var XidsAttrs sync.Map

func Range(f func(xid Xid) bool) {
	XidsAttrs.Range(func(k, v interface{}) bool {
		return f(k.(Xid))
	})
}

func NewXids() (xids Xids) {
	DockerScan()
	Range(func(xid Xid) bool {
		xids = append(xids, xid)
		return true
	})
	return
}

func (xids Xids) Cut(i int) Xids {
	copy(xids[i:], xids[i+1:])
	return xids[:len(xids)-1]
}

func (xids Xids) FilterContainer(re *regexp.Regexp) Xids {
	for i := 0; i < len(xids); {
		ns := xids[i].Attrs().IfInfoNetNs()
		if re.MatchString(ns.ContainerName()) ||
			re.MatchString(ns.ContainerId()) {
			i += 1
		} else {
			xids = xids.Cut(i)
		}
	}
	return xids
}

func (xids Xids) FilterName(re *regexp.Regexp) Xids {
	for i := 0; i < len(xids); {
		if re.MatchString(xids[i].Attrs().IfInfoName()) {
			i += 1
		} else {
			xids = xids.Cut(i)
		}
	}
	return xids
}

func (xids Xids) FilterNetNs(re *regexp.Regexp) Xids {
	for i := 0; i < len(xids); {
		ns := xids[i].Attrs().IfInfoNetNs()
		if re.MatchString(ns.String()) {
			i += 1
		} else {
			xids = xids.Cut(i)
		}
	}
	return xids
}

// Valid() if xid has mapped attributes
func (xid Xid) Valid() bool {
	_, ok := XidsAttrs.Load(xid)
	return ok
}

// get mapped attrs but panic if unavailable
func (xid Xid) Attrs() (attrs *XidAttrs) {
	if v, ok := XidsAttrs.Load(xid); ok {
		attrs = v.(*XidAttrs)
	} else if true {
		panic(fmt.Errorf("xid (%d, %d) hasn't been mapped",
			uint32(xid/VlanNVid), uint32(xid&VlanVidMask)))
	}
	return
}

// make the xid's attrs map if it's not already available
func (xid Xid) attrs() (attrs *XidAttrs) {
	if v, ok := XidsAttrs.Load(xid); ok {
		attrs = v.(*XidAttrs)
	} else {
		attrs = new(XidAttrs)
		XidsAttrs.Store(xid, attrs)
	}
	return
}

func (xid Xid) RxDelete() (note DevDel) {
	defer XidsAttrs.Delete(xid)
	note = DevDel(xid)
	if !xid.Valid() {
		return
	}
	attrs := xid.Attrs()
	for _, entry := range attrs.IPNets() {
		entry.IP = entry.IP[:cap(entry.IP)]
		entry.Mask = entry.Mask[:cap(entry.Mask)]
		poolIPNet.Put(entry)
	}
	attrs.Range(func(key, value interface{}) bool {
		defer attrs.Delete(key)
		if method, found := value.(pooler); found {
			method.Pool()
		}
		return true
	})
	return
}

func (xid Xid) String() string {
	return xid.Attrs().IfInfoName()
}

func (attrs *XidAttrs) EthtoolAutoNeg(set ...AutoNeg) (an AutoNeg) {
	if len(set) > 0 {
		an = set[0]
		attrs.Store(EthtoolAutoNegXidAttr, an)
	} else if v, ok := attrs.Load(EthtoolAutoNegXidAttr); ok {
		an = v.(AutoNeg)
	}
	return
}

func (attrs *XidAttrs) EthtoolDuplex(set ...Duplex) (duplex Duplex) {
	if len(set) > 0 {
		duplex = set[0]
		attrs.Store(EthtoolDuplexXidAttr, duplex)
	} else if v, ok := attrs.Load(EthtoolDuplexXidAttr); ok {
		duplex = v.(Duplex)
	}
	return
}

func (attrs *XidAttrs) EthtoolDevPort(set ...DevPort) (devport DevPort) {
	if len(set) > 0 {
		devport = set[0]
		attrs.Store(EthtoolDevPortXidAttr, devport)
	} else if v, ok := attrs.Load(EthtoolDevPortXidAttr); ok {
		devport = v.(DevPort)
	}
	return
}

func (attrs *XidAttrs) EthtoolFlags(set ...EthtoolFlagBits) (bits EthtoolFlagBits) {
	if len(set) > 0 {
		bits = set[0]
		attrs.Store(EthtoolFlagsXidAttr, bits)
	} else if v, ok := attrs.Load(EthtoolFlagsXidAttr); ok {
		bits = v.(EthtoolFlagBits)
	}
	return
}

func (attrs *XidAttrs) EthtoolSpeed(set ...uint32) (mbps uint32) {
	if len(set) > 0 {
		mbps = set[0]
		attrs.Store(EthtoolSpeedXidAttr, mbps)
	} else if v, ok := attrs.Load(EthtoolSpeedXidAttr); ok {
		mbps = v.(uint32)
	}
	return
}

func (attrs *XidAttrs) IfInfoName(set ...string) (name string) {
	if len(set) > 0 {
		name = set[0]
		attrs.Store(IfInfoNameXidAttr, name)
	} else if v, ok := attrs.Load(IfInfoNameXidAttr); ok {
		name = v.(string)
	}
	return
}

func (attrs *XidAttrs) IfInfoIfIndex(set ...int32) (ifindex int32) {
	if len(set) > 0 {
		ifindex = set[0]
		attrs.Store(IfInfoIfIndexXidAttr, ifindex)
	} else if v, ok := attrs.Load(IfInfoIfIndexXidAttr); ok {
		ifindex = v.(int32)
	}
	return
}

func (attrs *XidAttrs) IfInfoNetNs(set ...NetNs) (netns NetNs) {
	if len(set) > 0 {
		netns = set[0]
		attrs.Store(IfInfoNetNsXidAttr, netns)
	} else if v, ok := attrs.Load(IfInfoNetNsXidAttr); ok {
		netns = v.(NetNs)
	}
	return
}

func (attrs *XidAttrs) IfInfoFlags(set ...net.Flags) (flags net.Flags) {
	if len(set) > 0 {
		flags = set[0]
		attrs.Store(IfInfoFlagsXidAttr, flags)
	} else if v, ok := attrs.Load(IfInfoFlagsXidAttr); ok {
		flags = v.(net.Flags)
	}
	return
}

func (attrs *XidAttrs) IfInfoDevKind(set ...DevKind) (devkind DevKind) {
	if len(set) > 0 {
		devkind = set[0]
		attrs.Store(IfInfoDevKindXidAttr, devkind)
	} else if v, ok := attrs.Load(IfInfoDevKindXidAttr); ok {
		devkind = v.(DevKind)
	}
	return
}

func (attrs *XidAttrs) IfInfoHardwareAddr(set ...net.HardwareAddr) (ha net.HardwareAddr) {
	if len(set) > 0 {
		ha = set[0]
		attrs.Store(IfInfoHardwareAddrXidAttr, ha)
	} else if v, ok := attrs.Load(IfInfoHardwareAddrXidAttr); ok {
		ha = v.(net.HardwareAddr)
	}
	return
}

func (attrs *XidAttrs) IPNets(set ...[]*net.IPNet) (l []*net.IPNet) {
	if len(set) > 0 {
		l = set[0]
		attrs.Store(IPNetsXidAttr, l)
	} else if v, ok := attrs.Load(IPNetsXidAttr); ok {
		l = v.([]*net.IPNet)
	}
	return
}

func (attrs *XidAttrs) IsAdminUp() bool {
	return attrs.IfInfoFlags()&net.FlagUp == net.FlagUp
}

func (attrs *XidAttrs) IsAutoNeg() bool {
	return attrs.EthtoolAutoNeg() == AUTONEG_ENABLE
}

func (attrs *XidAttrs) IsBridge() bool {
	return attrs.IfInfoDevKind() == DevKindBridge
}

func (attrs *XidAttrs) IsLag() bool {
	return attrs.IfInfoDevKind() == DevKindLag
}

func (attrs *XidAttrs) IsPort() bool {
	return attrs.IfInfoDevKind() == DevKindPort
}

func (attrs *XidAttrs) IsVlan() bool {
	return attrs.IfInfoDevKind() == DevKindVlan
}

func (attrs *XidAttrs) LinkModesSupported(set ...EthtoolLinkModeBits) (modes EthtoolLinkModeBits) {
	if len(set) > 0 {
		modes = set[0]
		attrs.Store(LinkModesSupportedXidAttr, modes)
	} else if v, ok := attrs.Load(LinkModesSupportedXidAttr); ok {
		modes = v.(EthtoolLinkModeBits)
	}
	return
}

func (attrs *XidAttrs) LinkModesAdvertising(set ...EthtoolLinkModeBits) (modes EthtoolLinkModeBits) {
	if len(set) > 0 {
		modes = set[0]
		attrs.Store(LinkModesAdvertisingXidAttr, modes)
	} else if v, ok := attrs.Load(LinkModesAdvertisingXidAttr); ok {
		modes = v.(EthtoolLinkModeBits)
	}
	return
}

func (attrs *XidAttrs) LinkModesLPAdvertising(set ...EthtoolLinkModeBits) (modes EthtoolLinkModeBits) {
	if len(set) > 0 {
		modes = set[0]
		attrs.Store(LinkModesLPAdvertisingXidAttr, modes)
	} else if v, ok := attrs.Load(LinkModesLPAdvertisingXidAttr); ok {
		modes = v.(EthtoolLinkModeBits)
	}
	return
}

func (attrs *XidAttrs) LinkUp(set ...bool) (up bool) {
	if len(set) > 0 {
		up = set[0]
		attrs.Store(LinkUpXidAttr, up)
	} else if v, ok := attrs.Load(LinkUpXidAttr); ok {
		up = v.(bool)
	}
	return
}

func (attrs *XidAttrs) Lowers(set ...[]Xid) (xids []Xid) {
	if len(set) > 0 {
		xids = set[0]
		attrs.Store(LowersXidAttr, xids)
	} else if v, ok := attrs.Load(LowersXidAttr); ok {
		xids = v.([]Xid)
	}
	return
}

func (attrs *XidAttrs) Uppers(set ...[]Xid) (xids []Xid) {
	if len(set) > 0 {
		xids = set[0]
		attrs.Store(UppersXidAttr, xids)
	} else if v, ok := attrs.Load(UppersXidAttr); ok {
		xids = v.([]Xid)
	}
	return
}

func (attrs *XidAttrs) Stats(set ...[]uint64) (stats []uint64) {
	if len(set) > 0 {
		stats = set[0]
		attrs.Store(StatsXidAttr, stats)
	} else if v, ok := attrs.Load(StatsXidAttr); ok {
		stats = v.([]uint64)
	}
	return
}

func (attrs *XidAttrs) StatNames(set ...[]string) (names []string) {
	if len(set) > 0 {
		names = set[0]
		attrs.Store(StatNamesXidAttr, names)
	} else if v, ok := attrs.Load(StatNamesXidAttr); ok {
		names = v.([]string)
	}
	return
}
