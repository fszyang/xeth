// Copyright © 2018-2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xeth

type EthtoolLinkModeBits []uint8

type DevLinkModesSupported Xid
type DevLinkModesAdvertising Xid
type DevLinkModesLPAdvertising Xid

func (xid Xid) RxSupported(modes []uint8) DevLinkModesSupported {
	l := LinkOf(xid)
	supported := l.LinkModesSupported()
	if supported == nil || len(supported) != len(modes) {
		supported = make(EthtoolLinkModeBits, len(modes))
	}
	copy(supported, modes)
	l.LinkModesSupported(supported)
	return DevLinkModesSupported(xid)
}

func (xid Xid) RxAdvertising(modes []uint8) DevLinkModesAdvertising {
	l := LinkOf(xid)
	advertising := l.LinkModesAdvertising()
	if advertising == nil || len(advertising) != len(modes) {
		advertising = make(EthtoolLinkModeBits, len(modes))
	}
	copy(advertising, modes)
	l.LinkModesAdvertising(advertising)
	return DevLinkModesAdvertising(xid)
}

func (xid Xid) RxLPAdvertising(modes []uint8) DevLinkModesLPAdvertising {
	l := LinkOf(xid)
	lpadvertising := l.LinkModesLPAdvertising()
	if lpadvertising == nil || len(lpadvertising) != len(modes) {
		lpadvertising = make(EthtoolLinkModeBits, len(modes))
	}
	copy(lpadvertising, modes)
	l.LinkModesLPAdvertising(lpadvertising)
	return DevLinkModesLPAdvertising(xid)
}

func (bits EthtoolLinkModeBits) Test(bit uint) bool {
	if bit < uint(len(bits)*8) {
		mask := uint8(bit) & (8 - 1)
		return (bits[bit/8] & mask) == mask
	}
	return false
}
