// Copyright © 2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package xeth

type DevJoin struct{ Lower, Upper Xid }
type DevQuit struct{ Lower, Upper Xid }

func (lower Xid) join(upper Xid) *DevJoin {
	lowerAttrs := LinkAttrsOf(lower)
	upperAttrs := LinkAttrsOf(upper)
	lowerAttrs.Uppers(upper.List(lowerAttrs.Uppers()))
	upperAttrs.Lowers(lower.List(upperAttrs.Lowers()))
	return &DevJoin{lower, upper}
}

func (lower Xid) quit(upper Xid) *DevQuit {
	lowerAttrs := LinkAttrsOf(lower)
	upperAttrs := LinkAttrsOf(upper)
	lowerAttrs.Uppers(upper.Delist(lowerAttrs.Uppers()))
	upperAttrs.Lowers(lower.Delist(upperAttrs.Lowers()))
	return &DevQuit{lower, upper}
}

func (xid Xid) List(xids []Xid) []Xid {
	for _, entry := range xids {
		if entry == xid {
			return xids
		}
	}
	return append(xids, xid)
}

func (xid Xid) Delist(xids []Xid) []Xid {
	for i, entry := range xids {
		if entry == xid {
			n := len(xids) - 1
			copy(xids[i:], xids[i+1:])
			xids = xids[:n]
			break
		}
	}
	return xids
}
