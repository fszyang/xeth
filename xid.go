// Copyright © 2018-2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xeth

import (
	"fmt"
	"net"
	"sync"
)

type Xid uint32
type XidAttr uint8
type XidAttrs sync.Map

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

var XidAttrMaps sync.Map

func Range(f func(xid Xid) bool) {
	XidAttrMaps.Range(func(k, v interface{}) bool {
		return f(k.(Xid))
	})
}

// Valid() if xid has mapped attributes
func (xid Xid) Valid() bool {
	_, ok := XidAttrMaps.Load(xid)
	return ok
}

// get mapped attrs but panic if unavailable
func (xid Xid) Attrs() (attrs *XidAttrs) {
	if v, ok := XidAttrMaps.Load(xid); ok {
		attrs = (*XidAttrs)(v.(*sync.Map))
	} else if true {
		panic(fmt.Errorf("xid %d hasn't been mapped", uint32(xid)))
	}
	return
}

// make the xid's attrs map if it's not already available
func (xid Xid) attrs() (attrs *XidAttrs) {
	if v, ok := XidAttrMaps.Load(xid); ok {
		attrs = (*XidAttrs)(v.(*sync.Map))
	} else {
		m := new(sync.Map)
		attrs = (*XidAttrs)(m)
		XidAttrMaps.Store(xid, m)
	}
	return
}

func (xid Xid) RxDelete() (note DevDel) {
	defer XidAttrMaps.Delete(xid)
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
	m := attrs.Map()
	m.Range(func(key, value interface{}) bool {
		defer m.Delete(key)
		if method, found := value.(pooler); found {
			method.Pool()
		}
		return true
	})
	return
}

func (attrs *XidAttrs) Map() *sync.Map {
	return (*sync.Map)(attrs)
}

func (attrs *XidAttrs) Delete(attr XidAttr) {
	attrs.Map().Delete(attr)
}

func (attrs *XidAttrs) EthtoolAutoNeg(set ...AutoNeg) (an AutoNeg) {
	m := attrs.Map()
	if len(set) > 0 {
		an = set[0]
		m.Store(EthtoolAutoNegXidAttr, an)
	} else if v, ok := m.Load(EthtoolAutoNegXidAttr); ok {
		an = v.(AutoNeg)
	}
	return
}

func (attrs *XidAttrs) EthtoolDuplex(set ...Duplex) (duplex Duplex) {
	m := attrs.Map()
	if len(set) > 0 {
		duplex = set[0]
		m.Store(EthtoolDuplexXidAttr, duplex)
	} else if v, ok := m.Load(EthtoolDuplexXidAttr); ok {
		duplex = v.(Duplex)
	}
	return
}

func (attrs *XidAttrs) EthtoolDevPort(set ...DevPort) (devport DevPort) {
	m := attrs.Map()
	if len(set) > 0 {
		devport = set[0]
		m.Store(EthtoolDevPortXidAttr, devport)
	} else if v, ok := m.Load(EthtoolDevPortXidAttr); ok {
		devport = v.(DevPort)
	}
	return
}

func (attrs *XidAttrs) EthtoolFlags(set ...EthtoolFlagBits) (bits EthtoolFlagBits) {
	m := attrs.Map()
	if len(set) > 0 {
		bits = set[0]
		m.Store(EthtoolFlagsXidAttr, bits)
	} else if v, ok := m.Load(EthtoolFlagsXidAttr); ok {
		bits = v.(EthtoolFlagBits)
	}
	return
}

func (attrs *XidAttrs) EthtoolSpeed(set ...uint32) (mbps uint32) {
	m := attrs.Map()
	if len(set) > 0 {
		mbps = set[0]
		m.Store(EthtoolSpeedXidAttr, mbps)
	} else if v, ok := m.Load(EthtoolSpeedXidAttr); ok {
		mbps = v.(uint32)
	}
	return
}

func (attrs *XidAttrs) IfInfoName(set ...string) (name string) {
	m := attrs.Map()
	if len(set) > 0 {
		name = set[0]
		m.Store(IfInfoNameXidAttr, name)
	} else if v, ok := m.Load(IfInfoNameXidAttr); ok {
		name = v.(string)
	}
	return
}

func (attrs *XidAttrs) IfInfoIfIndex(set ...int32) (ifindex int32) {
	m := attrs.Map()
	if len(set) > 0 {
		ifindex = set[0]
		m.Store(IfInfoIfIndexXidAttr, ifindex)
	} else if v, ok := m.Load(IfInfoIfIndexXidAttr); ok {
		ifindex = v.(int32)
	}
	return
}

func (attrs *XidAttrs) IfInfoNetNs(set ...NetNs) (netns NetNs) {
	m := attrs.Map()
	if len(set) > 0 {
		netns = set[0]
		m.Store(IfInfoNetNsXidAttr, netns)
	} else if v, ok := m.Load(IfInfoNetNsXidAttr); ok {
		netns = v.(NetNs)
	}
	return
}

func (attrs *XidAttrs) IfInfoFlags(set ...net.Flags) (flags net.Flags) {
	m := attrs.Map()
	if len(set) > 0 {
		flags = set[0]
		m.Store(IfInfoFlagsXidAttr, flags)
	} else if v, ok := m.Load(IfInfoFlagsXidAttr); ok {
		flags = v.(net.Flags)
	}
	return
}

func (attrs *XidAttrs) IfInfoDevKind(set ...DevKind) (devkind DevKind) {
	m := attrs.Map()
	if len(set) > 0 {
		devkind = set[0]
		m.Store(IfInfoDevKindXidAttr, devkind)
	} else if v, ok := m.Load(IfInfoDevKindXidAttr); ok {
		devkind = v.(DevKind)
	}
	return
}

func (attrs *XidAttrs) IfInfoHardwareAddr(set ...net.HardwareAddr) (ha net.HardwareAddr) {
	m := attrs.Map()
	if len(set) > 0 {
		ha = set[0]
		m.Store(IfInfoHardwareAddrXidAttr, ha)
	} else if v, ok := m.Load(IfInfoHardwareAddrXidAttr); ok {
		ha = v.(net.HardwareAddr)
	}
	return
}

func (attrs *XidAttrs) IPNets(set ...[]*net.IPNet) (l []*net.IPNet) {
	m := attrs.Map()
	if len(set) > 0 {
		l = set[0]
		m.Store(IPNetsXidAttr, l)
	} else if v, ok := m.Load(IPNetsXidAttr); ok {
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
	m := attrs.Map()
	if len(set) > 0 {
		modes = set[0]
		m.Store(LinkModesSupportedXidAttr, modes)
	} else if v, ok := m.Load(LinkModesSupportedXidAttr); ok {
		modes = v.(EthtoolLinkModeBits)
	}
	return
}

func (attrs *XidAttrs) LinkModesAdvertising(set ...EthtoolLinkModeBits) (modes EthtoolLinkModeBits) {
	m := attrs.Map()
	if len(set) > 0 {
		modes = set[0]
		m.Store(LinkModesAdvertisingXidAttr, modes)
	} else if v, ok := m.Load(LinkModesAdvertisingXidAttr); ok {
		modes = v.(EthtoolLinkModeBits)
	}
	return
}

func (attrs *XidAttrs) LinkModesLPAdvertising(set ...EthtoolLinkModeBits) (modes EthtoolLinkModeBits) {
	m := attrs.Map()
	if len(set) > 0 {
		modes = set[0]
		m.Store(LinkModesLPAdvertisingXidAttr, modes)
	} else if v, ok := m.Load(LinkModesLPAdvertisingXidAttr); ok {
		modes = v.(EthtoolLinkModeBits)
	}
	return
}

func (attrs *XidAttrs) LinkUp(set ...bool) (up bool) {
	m := attrs.Map()
	if len(set) > 0 {
		up = set[0]
		m.Store(LinkUpXidAttr, up)
	} else if v, ok := m.Load(LinkUpXidAttr); ok {
		up = v.(bool)
	}
	return
}

func (attrs *XidAttrs) Lowers(set ...[]Xid) (xids []Xid) {
	m := attrs.Map()
	if len(set) > 0 {
		xids = set[0]
		m.Store(LowersXidAttr, xids)
	} else if v, ok := m.Load(LowersXidAttr); ok {
		xids = v.([]Xid)
	}
	return
}

func (attrs *XidAttrs) Uppers(set ...[]Xid) (xids []Xid) {
	m := attrs.Map()
	if len(set) > 0 {
		xids = set[0]
		m.Store(UppersXidAttr, xids)
	} else if v, ok := m.Load(UppersXidAttr); ok {
		xids = v.([]Xid)
	}
	return
}

func (attrs *XidAttrs) Stats(set ...[]uint64) (stats []uint64) {
	m := attrs.Map()
	if len(set) > 0 {
		stats = set[0]
		m.Store(StatsXidAttr, stats)
	} else if v, ok := m.Load(StatsXidAttr); ok {
		stats = v.([]uint64)
	}
	return
}

func (attrs *XidAttrs) StatNames(set ...[]string) (names []string) {
	m := attrs.Map()
	if len(set) > 0 {
		names = set[0]
		m.Store(StatNamesXidAttr, names)
	} else if v, ok := m.Load(StatNamesXidAttr); ok {
		names = v.([]string)
	}
	return
}
