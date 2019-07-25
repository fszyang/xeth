// Copyright © 2018-2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xeth

type EthtoolFlagBits uint32

type DevEthtoolFlags struct {
	Xid
	EthtoolFlagBits
}

func (xid Xid) RxEthtoolFlags(flags uint32) *DevEthtoolFlags {
	m := xid.Map()
	bits := EthtoolFlagBits(flags)
	if flags == 0 {
		m.Delete(EthtoolFlagsAttr)
	} else {
		m.Store(EthtoolFlagsAttr, bits)
	}
	return &DevEthtoolFlags{xid, bits}
}

func (bits EthtoolFlagBits) Test(bit uint) bool {
	mask := uint32(1 << bit)
	return (uint32(bits) & mask) == mask
}
