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
type XethLinkAttr uint8

const (
	XethLinkAttrEthtoolAutoNeg XethLinkAttr = iota
	XethLinkAttrEthtoolDevPort
	XethLinkAttrEthtoolDuplex
	XethLinkAttrEthtoolFlags
	XethLinkAttrEthtoolSpeed
	XethLinkAttrIPNets
	XethLinkAttrIfInfoName
	XethLinkAttrIfInfoIfIndex
	XethLinkAttrIfInfoNetNs
	XethLinkAttrIfInfoFlags
	XethLinkAttrIfInfoDevKind
	XethLinkAttrIfInfoHardwareAddr
	XethLinkAttrLinkModesAdvertising
	XethLinkAttrLinkModesLPAdvertising
	XethLinkAttrLinkModesSupported
	XethLinkAttrLinkUp
	XethLinkAttrLowers
	XethLinkAttrStatNames
	XethLinkAttrStats
	XethLinkAttrUppers
	XethLinkAttrXid
)

type XethLinkAttrs struct {
	*sync.Map
}

var LinkAttrMaps sync.Map

// get attrs map but panic if unavailable
func LinkAttrMap(xid Xid) (m *sync.Map) {
	if v, ok := LinkAttrMaps.Load(xid); ok {
		m = v.(*sync.Map)
	} else if true {
		panic(fmt.Errorf("xid (%d, %d) hasn't been mapped",
			uint32(xid/VlanNVid), uint32(xid&VlanVidMask)))
	}
	return
}

func LinkAttrs(xid Xid) XethLinkAttrs {
	return XethLinkAttrs{LinkAttrMap(xid)}
}

func LinkRange(f func(xid Xid, m *sync.Map) bool) {
	LinkAttrMaps.Range(func(k, v interface{}) bool {
		return f(k.(Xid), v.(*sync.Map))
	})
}

func ListXids() (xids Xids) {
	// scan docker containers to cache their name space attributes
	DockerScan()
	LinkRange(func(xid Xid, m *sync.Map) bool {
		xids = append(xids, xid)
		return true
	})
	return
}

// make the xid attrs map if it's not already available
func MayMakeLinkAttrMap(xid Xid) (m *sync.Map) {
	if v, ok := LinkAttrMaps.Load(xid); ok {
		m = v.(*sync.Map)
	} else {
		m = new(sync.Map)
		LinkAttrMaps.Store(xid, m)
		m.Store(XethLinkAttrXid, xid)
	}
	return
}

func MayMakeLinkAttrs(xid Xid) XethLinkAttrs {
	return XethLinkAttrs{MayMakeLinkAttrMap(xid)}
}

func RxDelete(xid Xid) (note DevDel) {
	defer LinkAttrMaps.Delete(xid)
	note = DevDel(xid)
	if !Valid(xid) {
		return
	}
	attrs := LinkAttrs(xid)
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

// Valid() if xid has mapped attributes
func Valid(xid Xid) bool {
	_, ok := LinkAttrMaps.Load(xid)
	return ok
}

func (attrs XethLinkAttrs) EthtoolAutoNeg(set ...AutoNeg) (an AutoNeg) {
	if len(set) > 0 {
		an = set[0]
		attrs.Store(XethLinkAttrEthtoolAutoNeg, an)
	} else if v, ok := attrs.Load(XethLinkAttrEthtoolAutoNeg); ok {
		an = v.(AutoNeg)
	}
	return
}

func (attrs XethLinkAttrs) EthtoolDuplex(set ...Duplex) (duplex Duplex) {
	if len(set) > 0 {
		duplex = set[0]
		attrs.Store(XethLinkAttrEthtoolDuplex, duplex)
	} else if v, ok := attrs.Load(XethLinkAttrEthtoolDuplex); ok {
		duplex = v.(Duplex)
	}
	return
}

func (attrs XethLinkAttrs) EthtoolDevPort(set ...DevPort) (devport DevPort) {
	if len(set) > 0 {
		devport = set[0]
		attrs.Store(XethLinkAttrEthtoolDevPort, devport)
	} else if v, ok := attrs.Load(XethLinkAttrEthtoolDevPort); ok {
		devport = v.(DevPort)
	}
	return
}

func (attrs XethLinkAttrs) EthtoolFlags(set ...EthtoolFlagBits) (bits EthtoolFlagBits) {
	if len(set) > 0 {
		bits = set[0]
		attrs.Store(XethLinkAttrEthtoolFlags, bits)
	} else if v, ok := attrs.Load(XethLinkAttrEthtoolFlags); ok {
		bits = v.(EthtoolFlagBits)
	}
	return
}

func (attrs XethLinkAttrs) EthtoolSpeed(set ...uint32) (mbps uint32) {
	if len(set) > 0 {
		mbps = set[0]
		attrs.Store(XethLinkAttrEthtoolSpeed, mbps)
	} else if v, ok := attrs.Load(XethLinkAttrEthtoolSpeed); ok {
		mbps = v.(uint32)
	}
	return
}

func (attrs XethLinkAttrs) IfInfoName(set ...string) (name string) {
	if len(set) > 0 {
		name = set[0]
		attrs.Store(XethLinkAttrIfInfoName, name)
	} else if v, ok := attrs.Load(XethLinkAttrIfInfoName); ok {
		name = v.(string)
	}
	return
}

func (attrs XethLinkAttrs) IfInfoIfIndex(set ...int32) (ifindex int32) {
	if len(set) > 0 {
		ifindex = set[0]
		attrs.Store(XethLinkAttrIfInfoIfIndex, ifindex)
	} else if v, ok := attrs.Load(XethLinkAttrIfInfoIfIndex); ok {
		ifindex = v.(int32)
	}
	return
}

func (attrs XethLinkAttrs) IfInfoNetNs(set ...NetNs) (netns NetNs) {
	if len(set) > 0 {
		netns = set[0]
		attrs.Store(XethLinkAttrIfInfoNetNs, netns)
	} else if v, ok := attrs.Load(XethLinkAttrIfInfoNetNs); ok {
		netns = v.(NetNs)
	}
	return
}

func (attrs XethLinkAttrs) IfInfoFlags(set ...net.Flags) (flags net.Flags) {
	if len(set) > 0 {
		flags = set[0]
		attrs.Store(XethLinkAttrIfInfoFlags, flags)
	} else if v, ok := attrs.Load(XethLinkAttrIfInfoFlags); ok {
		flags = v.(net.Flags)
	}
	return
}

func (attrs XethLinkAttrs) IfInfoDevKind(set ...DevKind) (devkind DevKind) {
	if len(set) > 0 {
		devkind = set[0]
		attrs.Store(XethLinkAttrIfInfoDevKind, devkind)
	} else if v, ok := attrs.Load(XethLinkAttrIfInfoDevKind); ok {
		devkind = v.(DevKind)
	}
	return
}

func (attrs XethLinkAttrs) IfInfoHardwareAddr(set ...net.HardwareAddr) (ha net.HardwareAddr) {
	if len(set) > 0 {
		ha = set[0]
		attrs.Store(XethLinkAttrIfInfoHardwareAddr, ha)
	} else if v, ok := attrs.Load(XethLinkAttrIfInfoHardwareAddr); ok {
		ha = v.(net.HardwareAddr)
	}
	return
}

func (attrs XethLinkAttrs) IPNets(set ...[]*net.IPNet) (l []*net.IPNet) {
	if len(set) > 0 {
		l = set[0]
		attrs.Store(XethLinkAttrIPNets, l)
	} else if v, ok := attrs.Load(XethLinkAttrIPNets); ok {
		l = v.([]*net.IPNet)
	}
	return
}

func (attrs XethLinkAttrs) IsAdminUp() bool {
	return attrs.IfInfoFlags()&net.FlagUp == net.FlagUp
}

func (attrs XethLinkAttrs) IsAutoNeg() bool {
	return attrs.EthtoolAutoNeg() == AUTONEG_ENABLE
}

func (attrs XethLinkAttrs) IsBridge() bool {
	return attrs.IfInfoDevKind() == DevKindBridge
}

func (attrs XethLinkAttrs) IsLag() bool {
	return attrs.IfInfoDevKind() == DevKindLag
}

func (attrs XethLinkAttrs) IsPort() bool {
	return attrs.IfInfoDevKind() == DevKindPort
}

func (attrs XethLinkAttrs) IsVlan() bool {
	return attrs.IfInfoDevKind() == DevKindVlan
}

func (attrs XethLinkAttrs) LinkModesSupported(set ...EthtoolLinkModeBits) (modes EthtoolLinkModeBits) {
	if len(set) > 0 {
		modes = set[0]
		attrs.Store(XethLinkAttrLinkModesSupported, modes)
	} else if v, ok := attrs.Load(XethLinkAttrLinkModesSupported); ok {
		modes = v.(EthtoolLinkModeBits)
	}
	return
}

func (attrs XethLinkAttrs) LinkModesAdvertising(set ...EthtoolLinkModeBits) (modes EthtoolLinkModeBits) {
	if len(set) > 0 {
		modes = set[0]
		attrs.Store(XethLinkAttrLinkModesAdvertising, modes)
	} else if v, ok := attrs.Load(XethLinkAttrLinkModesAdvertising); ok {
		modes = v.(EthtoolLinkModeBits)
	}
	return
}

func (attrs XethLinkAttrs) LinkModesLPAdvertising(set ...EthtoolLinkModeBits) (modes EthtoolLinkModeBits) {
	if len(set) > 0 {
		modes = set[0]
		attrs.Store(XethLinkAttrLinkModesLPAdvertising, modes)
	} else if v, ok := attrs.Load(XethLinkAttrLinkModesLPAdvertising); ok {
		modes = v.(EthtoolLinkModeBits)
	}
	return
}

func (attrs XethLinkAttrs) LinkUp(set ...bool) (up bool) {
	if len(set) > 0 {
		up = set[0]
		attrs.Store(XethLinkAttrLinkUp, up)
	} else if v, ok := attrs.Load(XethLinkAttrLinkUp); ok {
		up = v.(bool)
	}
	return
}

func (attrs XethLinkAttrs) Lowers(set ...[]Xid) (xids []Xid) {
	if len(set) > 0 {
		xids = set[0]
		attrs.Store(XethLinkAttrLowers, xids)
	} else if v, ok := attrs.Load(XethLinkAttrLowers); ok {
		xids = v.([]Xid)
	}
	return
}

func (attrs XethLinkAttrs) Uppers(set ...[]Xid) (xids []Xid) {
	if len(set) > 0 {
		xids = set[0]
		attrs.Store(XethLinkAttrUppers, xids)
	} else if v, ok := attrs.Load(XethLinkAttrUppers); ok {
		xids = v.([]Xid)
	}
	return
}

func (attrs XethLinkAttrs) Stats(set ...[]uint64) (stats []uint64) {
	if len(set) > 0 {
		stats = set[0]
		attrs.Store(XethLinkAttrStats, stats)
	} else if v, ok := attrs.Load(XethLinkAttrStats); ok {
		stats = v.([]uint64)
	}
	return
}

func (attrs XethLinkAttrs) StatNames(set ...[]string) (names []string) {
	if len(set) > 0 {
		names = set[0]
		attrs.Store(XethLinkAttrStatNames, names)
	} else if v, ok := attrs.Load(XethLinkAttrStatNames); ok {
		names = v.([]string)
	}
	return
}

func (attrs XethLinkAttrs) String() string {
	return attrs.IfInfoName()
}

func (xids Xids) Cut(i int) Xids {
	copy(xids[i:], xids[i+1:])
	return xids[:len(xids)-1]
}

func (xids Xids) FilterContainer(re *regexp.Regexp) Xids {
	for i := 0; i < len(xids); {
		ns := LinkAttrs(xids[i]).IfInfoNetNs()
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
		if re.MatchString(LinkAttrs(xids[i]).IfInfoName()) {
			i += 1
		} else {
			xids = xids.Cut(i)
		}
	}
	return xids
}

func (xids Xids) FilterNetNs(re *regexp.Regexp) Xids {
	for i := 0; i < len(xids); {
		ns := LinkAttrs(xids[i]).IfInfoNetNs()
		if re.MatchString(ns.String()) {
			i += 1
		} else {
			xids = xids.Cut(i)
		}
	}
	return xids
}
