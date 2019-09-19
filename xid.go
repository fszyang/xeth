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

type Linker interface {
	EthtoolAutoNeg(set ...AutoNeg) AutoNeg
	EthtoolDuplex(set ...Duplex) Duplex
	EthtoolDevPort(set ...DevPort) DevPort
	EthtoolFlags(set ...EthtoolFlagBits) EthtoolFlagBits
	EthtoolSpeed(set ...uint32) uint32
	IfInfoName(set ...string) string
	IfInfoIfIndex(set ...int32) int32
	IfInfoNetNs(set ...NetNs) NetNs
	IfInfoFlags(set ...net.Flags) net.Flags
	IfInfoDevKind(set ...DevKind) DevKind
	IfInfoHardwareAddr(set ...net.HardwareAddr) net.HardwareAddr
	IPNets(set ...[]*net.IPNet) []*net.IPNet
	IsAdminUp() bool
	IsAutoNeg() bool
	IsBridge() bool
	IsLag() bool
	IsPort() bool
	IsVlan() bool
	LinkModesSupported(set ...EthtoolLinkModeBits) EthtoolLinkModeBits
	LinkModesAdvertising(set ...EthtoolLinkModeBits) EthtoolLinkModeBits
	LinkModesLPAdvertising(set ...EthtoolLinkModeBits) EthtoolLinkModeBits
	LinkUp(set ...bool) bool
	Lowers(set ...[]Xid) []Xid
	Uppers(set ...[]Xid) []Xid
	Stats(set ...[]uint64) []uint64
	StatNames(set ...[]string) []string
	String() string
	Xid() Xid
}

type LinkAttr uint8

const (
	LinkAttrEthtoolAutoNeg LinkAttr = iota
	LinkAttrEthtoolDevPort
	LinkAttrEthtoolDuplex
	LinkAttrEthtoolFlags
	LinkAttrEthtoolSpeed
	LinkAttrIPNets
	LinkAttrIfInfoName
	LinkAttrIfInfoIfIndex
	LinkAttrIfInfoNetNs
	LinkAttrIfInfoFlags
	LinkAttrIfInfoDevKind
	LinkAttrIfInfoHardwareAddr
	LinkAttrLinkModesAdvertising
	LinkAttrLinkModesLPAdvertising
	LinkAttrLinkModesSupported
	LinkAttrLinkUp
	LinkAttrLowers
	LinkAttrStatNames
	LinkAttrStats
	LinkAttrUppers
	LinkAttrXid
)

type LinkAttrs struct {
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

func LinkAttrsOf(xid Xid) LinkAttrs {
	return LinkAttrs{LinkAttrMap(xid)}
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
		m.Store(LinkAttrXid, xid)
	}
	return
}

func MayMakeLinkAttrs(xid Xid) LinkAttrs {
	return LinkAttrs{MayMakeLinkAttrMap(xid)}
}

func RxDelete(xid Xid) (note DevDel) {
	defer LinkAttrMaps.Delete(xid)
	note = DevDel(xid)
	if !Valid(xid) {
		return
	}
	attrs := LinkAttrsOf(xid)
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

func (attrs LinkAttrs) EthtoolAutoNeg(set ...AutoNeg) (an AutoNeg) {
	if len(set) > 0 {
		an = set[0]
		attrs.Store(LinkAttrEthtoolAutoNeg, an)
	} else if v, ok := attrs.Load(LinkAttrEthtoolAutoNeg); ok {
		an = v.(AutoNeg)
	}
	return
}

func (attrs LinkAttrs) EthtoolDuplex(set ...Duplex) (duplex Duplex) {
	if len(set) > 0 {
		duplex = set[0]
		attrs.Store(LinkAttrEthtoolDuplex, duplex)
	} else if v, ok := attrs.Load(LinkAttrEthtoolDuplex); ok {
		duplex = v.(Duplex)
	}
	return
}

func (attrs LinkAttrs) EthtoolDevPort(set ...DevPort) (devport DevPort) {
	if len(set) > 0 {
		devport = set[0]
		attrs.Store(LinkAttrEthtoolDevPort, devport)
	} else if v, ok := attrs.Load(LinkAttrEthtoolDevPort); ok {
		devport = v.(DevPort)
	}
	return
}

func (attrs LinkAttrs) EthtoolFlags(set ...EthtoolFlagBits) (bits EthtoolFlagBits) {
	if len(set) > 0 {
		bits = set[0]
		attrs.Store(LinkAttrEthtoolFlags, bits)
	} else if v, ok := attrs.Load(LinkAttrEthtoolFlags); ok {
		bits = v.(EthtoolFlagBits)
	}
	return
}

func (attrs LinkAttrs) EthtoolSpeed(set ...uint32) (mbps uint32) {
	if len(set) > 0 {
		mbps = set[0]
		attrs.Store(LinkAttrEthtoolSpeed, mbps)
	} else if v, ok := attrs.Load(LinkAttrEthtoolSpeed); ok {
		mbps = v.(uint32)
	}
	return
}

func (attrs LinkAttrs) IfInfoName(set ...string) (name string) {
	if len(set) > 0 {
		name = set[0]
		attrs.Store(LinkAttrIfInfoName, name)
	} else if v, ok := attrs.Load(LinkAttrIfInfoName); ok {
		name = v.(string)
	}
	return
}

func (attrs LinkAttrs) IfInfoIfIndex(set ...int32) (ifindex int32) {
	if len(set) > 0 {
		ifindex = set[0]
		attrs.Store(LinkAttrIfInfoIfIndex, ifindex)
	} else if v, ok := attrs.Load(LinkAttrIfInfoIfIndex); ok {
		ifindex = v.(int32)
	}
	return
}

func (attrs LinkAttrs) IfInfoNetNs(set ...NetNs) (netns NetNs) {
	if len(set) > 0 {
		netns = set[0]
		attrs.Store(LinkAttrIfInfoNetNs, netns)
	} else if v, ok := attrs.Load(LinkAttrIfInfoNetNs); ok {
		netns = v.(NetNs)
	}
	return
}

func (attrs LinkAttrs) IfInfoFlags(set ...net.Flags) (flags net.Flags) {
	if len(set) > 0 {
		flags = set[0]
		attrs.Store(LinkAttrIfInfoFlags, flags)
	} else if v, ok := attrs.Load(LinkAttrIfInfoFlags); ok {
		flags = v.(net.Flags)
	}
	return
}

func (attrs LinkAttrs) IfInfoDevKind(set ...DevKind) (devkind DevKind) {
	if len(set) > 0 {
		devkind = set[0]
		attrs.Store(LinkAttrIfInfoDevKind, devkind)
	} else if v, ok := attrs.Load(LinkAttrIfInfoDevKind); ok {
		devkind = v.(DevKind)
	}
	return
}

func (attrs LinkAttrs) IfInfoHardwareAddr(set ...net.HardwareAddr) (ha net.HardwareAddr) {
	if len(set) > 0 {
		ha = set[0]
		attrs.Store(LinkAttrIfInfoHardwareAddr, ha)
	} else if v, ok := attrs.Load(LinkAttrIfInfoHardwareAddr); ok {
		ha = v.(net.HardwareAddr)
	}
	return
}

func (attrs LinkAttrs) IPNets(set ...[]*net.IPNet) (l []*net.IPNet) {
	if len(set) > 0 {
		l = set[0]
		attrs.Store(LinkAttrIPNets, l)
	} else if v, ok := attrs.Load(LinkAttrIPNets); ok {
		l = v.([]*net.IPNet)
	}
	return
}

func (attrs LinkAttrs) IsAdminUp() bool {
	return attrs.IfInfoFlags()&net.FlagUp == net.FlagUp
}

func (attrs LinkAttrs) IsAutoNeg() bool {
	return attrs.EthtoolAutoNeg() == AUTONEG_ENABLE
}

func (attrs LinkAttrs) IsBridge() bool {
	return attrs.IfInfoDevKind() == DevKindBridge
}

func (attrs LinkAttrs) IsLag() bool {
	return attrs.IfInfoDevKind() == DevKindLag
}

func (attrs LinkAttrs) IsPort() bool {
	return attrs.IfInfoDevKind() == DevKindPort
}

func (attrs LinkAttrs) IsVlan() bool {
	return attrs.IfInfoDevKind() == DevKindVlan
}

func (attrs LinkAttrs) LinkModesSupported(set ...EthtoolLinkModeBits) (modes EthtoolLinkModeBits) {
	if len(set) > 0 {
		modes = set[0]
		attrs.Store(LinkAttrLinkModesSupported, modes)
	} else if v, ok := attrs.Load(LinkAttrLinkModesSupported); ok {
		modes = v.(EthtoolLinkModeBits)
	}
	return
}

func (attrs LinkAttrs) LinkModesAdvertising(set ...EthtoolLinkModeBits) (modes EthtoolLinkModeBits) {
	if len(set) > 0 {
		modes = set[0]
		attrs.Store(LinkAttrLinkModesAdvertising, modes)
	} else if v, ok := attrs.Load(LinkAttrLinkModesAdvertising); ok {
		modes = v.(EthtoolLinkModeBits)
	}
	return
}

func (attrs LinkAttrs) LinkModesLPAdvertising(set ...EthtoolLinkModeBits) (modes EthtoolLinkModeBits) {
	if len(set) > 0 {
		modes = set[0]
		attrs.Store(LinkAttrLinkModesLPAdvertising, modes)
	} else if v, ok := attrs.Load(LinkAttrLinkModesLPAdvertising); ok {
		modes = v.(EthtoolLinkModeBits)
	}
	return
}

func (attrs LinkAttrs) LinkUp(set ...bool) (up bool) {
	if len(set) > 0 {
		up = set[0]
		attrs.Store(LinkAttrLinkUp, up)
	} else if v, ok := attrs.Load(LinkAttrLinkUp); ok {
		up = v.(bool)
	}
	return
}

func (attrs LinkAttrs) Lowers(set ...[]Xid) (xids []Xid) {
	if len(set) > 0 {
		xids = set[0]
		attrs.Store(LinkAttrLowers, xids)
	} else if v, ok := attrs.Load(LinkAttrLowers); ok {
		xids = v.([]Xid)
	}
	return
}

func (attrs LinkAttrs) Uppers(set ...[]Xid) (xids []Xid) {
	if len(set) > 0 {
		xids = set[0]
		attrs.Store(LinkAttrUppers, xids)
	} else if v, ok := attrs.Load(LinkAttrUppers); ok {
		xids = v.([]Xid)
	}
	return
}

func (attrs LinkAttrs) Stats(set ...[]uint64) (stats []uint64) {
	if len(set) > 0 {
		stats = set[0]
		attrs.Store(LinkAttrStats, stats)
	} else if v, ok := attrs.Load(LinkAttrStats); ok {
		stats = v.([]uint64)
	}
	return
}

func (attrs LinkAttrs) StatNames(set ...[]string) (names []string) {
	if len(set) > 0 {
		names = set[0]
		attrs.Store(LinkAttrStatNames, names)
	} else if v, ok := attrs.Load(LinkAttrStatNames); ok {
		names = v.([]string)
	}
	return
}

func (attrs LinkAttrs) String() string {
	return attrs.IfInfoName()
}

func (attrs LinkAttrs) Xid() (xid Xid) {
	if v, ok := attrs.Load(LinkAttrXid); ok {
		xid = v.(Xid)
	}
	return
}

func (xids Xids) Cut(i int) Xids {
	copy(xids[i:], xids[i+1:])
	return xids[:len(xids)-1]
}

func (xids Xids) FilterContainer(re *regexp.Regexp) Xids {
	for i := 0; i < len(xids); {
		ns := LinkAttrsOf(xids[i]).IfInfoNetNs()
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
		if re.MatchString(LinkAttrsOf(xids[i]).IfInfoName()) {
			i += 1
		} else {
			xids = xids.Cut(i)
		}
	}
	return xids
}

func (xids Xids) FilterNetNs(re *regexp.Regexp) Xids {
	for i := 0; i < len(xids); {
		ns := LinkAttrsOf(xids[i]).IfInfoNetNs()
		if re.MatchString(ns.String()) {
			i += 1
		} else {
			xids = xids.Cut(i)
		}
	}
	return xids
}
