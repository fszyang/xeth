// Copyright © 2018-2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xeth

import (
	"sync"
)

type EthtoolFlagBits uint32
type DevEthtoolFlags struct {
	Xid
	EthtoolFlagBits
}

var ethtoolFlags sync.Map

func (xid Xid) EthtoolFlags() (bits EthtoolFlagBits) {
	if v, ok := ethtoolFlags.Load(xid); ok {
		bits = v.(EthtoolFlagBits)
	}
	return
}

func (xid Xid) deleteEthtoolFlags() {
	ethtoolFlags.Delete(xid)
}

func (xid Xid) ethtoolFlags(flags uint32) *DevEthtoolFlags {
	bits := EthtoolFlagBits(flags)
	ethtoolFlags.Store(xid, bits)
	return &DevEthtoolFlags{xid, bits}
}

func (bits EthtoolFlagBits) Test(bit uint) bool {
	mask := uint32(1 << bit)
	return (uint32(bits) & mask) == mask
}
