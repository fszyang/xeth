// Copyright © 2018-2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xeth

import (
	"fmt"
	"net"
	"sync"
)

type Attr uint8

type Attrs sync.Map

const (
	EthtoolAutoNegAttr Attr = iota
	EthtoolDevPortAttr
	EthtoolDuplexAttr
	EthtoolFlagsAttr
	EthtoolSpeedAttr
	IPNetsAttr
	IfInfoNameAttr
	IfInfoIfIndexAttr
	IfInfoNetNsAttr
	IfInfoFlagsAttr
	IfInfoDevKindAttr
	IfInfoHardwareAddrAttr
	LinkModesAdvertisingAttr
	LinkModesLPAdvertisingAttr
	LinkModesSupportedAttr
	LinkUpAttr
	LowersAttr
	StatNamesAttr
	StatsAttr
	UppersAttr
)

var XidAttrMaps, InvalidXidAttrMap sync.Map

func Range(f func(xid Xid) bool) {
	XidAttrMaps.Range(func(k, v interface{}) bool {
		return f(k.(Xid))
	})
}

func (xid Xid) Valid() bool {
	_, ok := XidAttrMaps.Load(xid)
	return ok
}

func (xid Xid) Map() (m *sync.Map) {
	if v, ok := XidAttrMaps.Load(xid); ok {
		m = v.(*sync.Map)
	} else if true {
		panic(fmt.Errorf("xid %d hasn't been mapped", uint32(xid)))
	} else {
		m = &InvalidXidAttrMap
	}
	return
}

func (xid Xid) RxDelete() (note DevDel) {
	defer XidAttrMaps.Delete(xid)
	note = DevDel(xid)
	v, ok := XidAttrMaps.Load(xid)
	if !ok {
		return
	}
	m := v.(*sync.Map)
	attrs := (*Attrs)(m)
	for _, entry := range attrs.IPNets() {
		entry.IP = entry.IP[:cap(entry.IP)]
		entry.Mask = entry.Mask[:cap(entry.Mask)]
		poolIPNet.Put(entry)
	}
	m.Range(func(key, value interface{}) bool {
		defer m.Delete(key)
		if method, found := value.(pooler); found {
			method.Pool()
		}
		return true
	})
	return
}

func (xid Xid) Attrs() *Attrs {
	return (*Attrs)(xid.Map())
}

func (attrs *Attrs) Map() *sync.Map {
	return (*sync.Map)(attrs)
}

func (attrs *Attrs) EthtoolAutoNeg() (an AutoNeg) {
	m := attrs.Map()
	if v, ok := m.Load(EthtoolAutoNegAttr); ok {
		an = v.(AutoNeg)
	}
	return
}

func (attrs *Attrs) EthtoolDuplex() (duplex Duplex) {
	m := attrs.Map()
	if v, ok := m.Load(EthtoolDuplexAttr); ok {
		duplex = v.(Duplex)
	}
	return
}

func (attrs *Attrs) EthtoolDevPort() (devport DevPort) {
	m := attrs.Map()
	if v, ok := m.Load(EthtoolDevPortAttr); ok {
		devport = v.(DevPort)
	}
	return
}

func (attrs *Attrs) EthtoolFlags(set ...EthtoolFlagBits) (bits EthtoolFlagBits) {
	m := attrs.Map()
	if len(set) > 0 {
		bits = set[0]
		m.Store(EthtoolFlagsAttr, bits)
	} else if v, ok := m.Load(EthtoolFlagsAttr); ok {
		bits = v.(EthtoolFlagBits)
	}
	return
}

func (attrs *Attrs) EthtoolSpeed() (mbps uint32) {
	m := attrs.Map()
	if v, ok := m.Load(EthtoolSpeedAttr); ok {
		mbps = v.(uint32)
	}
	return
}

func (attrs *Attrs) IfInfoName() (name string) {
	m := attrs.Map()
	if v, ok := m.Load(IfInfoNameAttr); ok {
		name = v.(string)
	}
	return
}

func (attrs *Attrs) IfInfoIfIndex() (ifindex int32) {
	m := attrs.Map()
	if v, ok := m.Load(IfInfoIfIndexAttr); ok {
		ifindex = v.(int32)
	}
	return
}

func (attrs *Attrs) IfInfoNetNs() (netns NetNs) {
	m := attrs.Map()
	if v, ok := m.Load(IfInfoNetNsAttr); ok {
		netns = v.(NetNs)
	}
	return
}

func (attrs *Attrs) IfInfoFlags() (flags net.Flags) {
	m := attrs.Map()
	if v, ok := m.Load(IfInfoFlagsAttr); ok {
		flags = v.(net.Flags)
	}
	return
}

func (attrs *Attrs) IfInfoDevKind() (devkind DevKind) {
	m := attrs.Map()
	if v, ok := m.Load(IfInfoDevKindAttr); ok {
		devkind = v.(DevKind)
	}
	return
}

func (attrs *Attrs) IfInfoHardwareAddr() (ha net.HardwareAddr) {
	m := attrs.Map()
	if v, ok := m.Load(IfInfoHardwareAddrAttr); ok {
		ha = v.(net.HardwareAddr)
	}
	return
}

func (attrs *Attrs) IPNets(set ...[]*net.IPNet) (l []*net.IPNet) {
	m := attrs.Map()
	if len(set) > 0 {
		l = set[0]
		m.Store(IPNetsAttr, l)
	} else if v, ok := m.Load(IPNetsAttr); ok {
		l = v.([]*net.IPNet)
	}
	return
}

func (attrs *Attrs) IsAdminUp() bool {
	return attrs.IfInfoFlags()&net.FlagUp == net.FlagUp
}

func (attrs *Attrs) IsAutoNeg() bool {
	return attrs.EthtoolAutoNeg() == AUTONEG_ENABLE
}

func (attrs *Attrs) IsBridge() bool {
	return attrs.IfInfoDevKind() == DevKindBridge
}

func (attrs *Attrs) IsLag() bool {
	return attrs.IfInfoDevKind() == DevKindLag
}

func (attrs *Attrs) IsPort() bool {
	return attrs.IfInfoDevKind() == DevKindPort
}

func (attrs *Attrs) IsVlan() bool {
	return attrs.IfInfoDevKind() == DevKindVlan
}

func (attrs *Attrs) LinkModesSupported(set ...EthtoolLinkModeBits) (modes EthtoolLinkModeBits) {
	m := attrs.Map()
	if len(set) > 0 {
		modes = set[0]
		m.Store(LinkModesSupportedAttr, modes)
	} else if v, ok := m.Load(LinkModesSupportedAttr); ok {
		modes = v.(EthtoolLinkModeBits)
	}
	return
}

func (attrs *Attrs) LinkModesAdvertising(set ...EthtoolLinkModeBits) (modes EthtoolLinkModeBits) {
	m := attrs.Map()
	if len(set) > 0 {
		modes = set[0]
		m.Store(LinkModesAdvertisingAttr, modes)
	} else if v, ok := m.Load(LinkModesAdvertisingAttr); ok {
		modes = v.(EthtoolLinkModeBits)
	}
	return
}

func (attrs *Attrs) LinkModesLPAdvertising(set ...EthtoolLinkModeBits) (modes EthtoolLinkModeBits) {
	m := attrs.Map()
	if len(set) > 0 {
		modes = set[0]
		m.Store(LinkModesLPAdvertisingAttr, modes)
	} else if v, ok := m.Load(LinkModesLPAdvertisingAttr); ok {
		modes = v.(EthtoolLinkModeBits)
	}
	return
}

func (attrs *Attrs) LinkUp(set ...bool) (up bool) {
	m := attrs.Map()
	if len(set) > 0 {
		up = set[0]
		m.Store(LinkUpAttr, up)
	} else if v, ok := m.Load(LinkUpAttr); ok {
		up = v.(bool)
	}
	return
}

func (attrs *Attrs) Lowers(set ...[]Xid) (xids []Xid) {
	m := attrs.Map()
	if len(set) > 0 {
		xids = set[0]
		m.Store(LowersAttr, xids)
	} else if v, ok := m.Load(LowersAttr); ok {
		xids = v.([]Xid)
	}
	return
}

func (attrs *Attrs) Uppers(set ...[]Xid) (xids []Xid) {
	m := attrs.Map()
	if len(set) > 0 {
		xids = set[0]
		m.Store(UppersAttr, xids)
	} else if v, ok := m.Load(UppersAttr); ok {
		xids = v.([]Xid)
	}
	return
}

func (attrs *Attrs) Stats(set ...[]uint64) (stats []uint64) {
	m := attrs.Map()
	if len(set) > 0 {
		stats = set[0]
		m.Store(StatsAttr, stats)
	} else if v, ok := m.Load(StatsAttr); ok {
		stats = v.([]uint64)
	}
	return
}

func (attrs *Attrs) StatNames(set ...[]string) (names []string) {
	m := attrs.Map()
	if len(set) > 0 {
		names = set[0]
		m.Store(StatNamesAttr, names)
	} else if v, ok := m.Load(StatNamesAttr); ok {
		names = v.([]string)
	}
	return
}
