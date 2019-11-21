// Copyright © 2018-2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xeth

type EthtoolLinkModeBits uint64

type DevLinkModesSupported Xid
type DevLinkModesAdvertising Xid
type DevLinkModesLPAdvertising Xid

func (xid Xid) RxSupported(modes uint64) DevLinkModesSupported {
	LinkOf(xid).LinkModesSupported(EthtoolLinkModeBits(modes))
	return DevLinkModesSupported(xid)
}

func (xid Xid) RxAdvertising(modes uint64) DevLinkModesAdvertising {
	LinkOf(xid).LinkModesAdvertising(EthtoolLinkModeBits(modes))
	return DevLinkModesAdvertising(xid)
}

func (xid Xid) RxLPAdvertising(modes uint64) DevLinkModesLPAdvertising {
	LinkOf(xid).LinkModesLPAdvertising(EthtoolLinkModeBits(modes))
	return DevLinkModesLPAdvertising(xid)
}

func (bits EthtoolLinkModeBits) Test(bit uint) bool {
	if bit < 64 {
		mask := EthtoolLinkModeBits(1 << bit)
		return (bits & mask) == mask
	}
	return false
}
