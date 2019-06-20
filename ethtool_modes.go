// Copyright © 2018-2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xeth

import "sync"

type EthtoolLinkModeBits []uint8

type DevSupportedLinkModes Xid
type DevAdvertisingLinkModes Xid
type DevLPAdvertisingLinkModes Xid

var (
	supportedLinkModes     sync.Map
	advertisingLinkModes   sync.Map
	lpadvertisingLinkModes sync.Map
)

func (xid Xid) SupportedLinkModes() (modes EthtoolLinkModeBits) {
	if v, ok := supportedLinkModes.Load(xid); ok {
		modes = v.(EthtoolLinkModeBits)
	}
	return
}

func (xid Xid) AdvertisingLinkModes() (modes EthtoolLinkModeBits) {
	if v, ok := advertisingLinkModes.Load(xid); ok {
		modes = v.(EthtoolLinkModeBits)
	}
	return
}

func (xid Xid) LPAdvertisingLinkModes() (modes EthtoolLinkModeBits) {
	if v, ok := lpadvertisingLinkModes.Load(xid); ok {
		modes = v.(EthtoolLinkModeBits)
	}
	return
}

func (xid Xid) deleteSupportedLinkModes() {
	supportedLinkModes.Delete(xid)
}

func (xid Xid) deleteAdvertisingLinkModes() {
	advertisingLinkModes.Delete(xid)
}

func (xid Xid) deleteLPAdvertisingLinkModes() {
	lpadvertisingLinkModes.Delete(xid)
}

func (xid Xid) supportedLinkModes(modes []uint8) DevSupportedLinkModes {
	supported := xid.SupportedLinkModes()
	if supported == nil || len(supported) != len(modes) {
		supported = make(EthtoolLinkModeBits, len(modes))
	}
	copy(supported, modes)
	supportedLinkModes.Store(xid, supported)
	return DevSupportedLinkModes(xid)
}

func (xid Xid) advertisingLinkModes(modes []uint8) DevAdvertisingLinkModes {
	advertising := xid.AdvertisingLinkModes()
	if advertising == nil || len(advertising) != len(modes) {
		advertising = make(EthtoolLinkModeBits, len(modes))
	}
	copy(advertising, modes)
	advertisingLinkModes.Store(xid, advertising)
	return DevAdvertisingLinkModes(xid)
}

func (xid Xid) lpadvertisingLinkModes(modes []uint8) DevLPAdvertisingLinkModes {
	lpadvertising := xid.LPAdvertisingLinkModes()
	if lpadvertising == nil || len(lpadvertising) != len(modes) {
		lpadvertising = make(EthtoolLinkModeBits, len(modes))
	}
	copy(lpadvertising, modes)
	lpadvertisingLinkModes.Store(xid, lpadvertising)
	return DevLPAdvertisingLinkModes(xid)
}

func (bits EthtoolLinkModeBits) Test(bit uint) bool {
	if bit < uint(len(bits)*8) {
		mask := uint8(bit) & (8 - 1)
		return (bits[bit/8] & mask) == mask
	}
	return false
}
