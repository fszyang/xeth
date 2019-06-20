// Copyright © 2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package xeth

import "sync"

type Xid uint32

type DevJoin struct{ lower, upper Xid }
type DevQuit struct{ lower, upper Xid }

var uppers, lowers sync.Map

var xids struct {
	sync.Mutex
	list []Xid
}

func Range(f func(xid Xid) bool) {
	xids.Lock()
	defer xids.Unlock()
	for _, xid := range xids.list {
		if !f(xid) {
			break
		}
	}
}

func (xid Xid) Uppers() (xids []Xid) {
	if v, ok := uppers.Load(xid); ok {
		xids = v.([]Xid)
	}
	return
}

func (xid Xid) Lowers() (xids []Xid) {
	if v, ok := lowers.Load(xid); ok {
		xids = v.([]Xid)
	}
	return
}

func (lower Xid) Join(upper Xid) *DevJoin {
	uppers.Store(lower, upper.addTo(lower.Uppers()))
	lowers.Store(upper, lower.addTo(upper.Lowers()))
	return &DevJoin{lower, upper}
}

func (lower Xid) Quit(upper Xid) *DevQuit {
	uppers.Store(lower, upper.delFrom(lower.Uppers()))
	lowers.Store(upper, lower.delFrom(upper.Lowers()))
	return &DevQuit{lower, upper}
}

func (xid Xid) deleteUppers() {
	uppers.Delete(xid)
}

func (xid Xid) deleteLowers() {
	lowers.Delete(xid)
}

func (xid Xid) addTo(xids []Xid) []Xid {
	for _, entry := range xids {
		if entry == xid {
			return xids
		}
	}
	return append(xids, xid)
}

func (xid Xid) addToXids() {
	xids.Lock()
	defer xids.Unlock()
	xids.list = xid.addTo(xids.list)
}

func (xid Xid) delFrom(xids []Xid) []Xid {
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

func (xid Xid) delFromXids() {
	xids.Lock()
	defer xids.Unlock()
	xids.list = xid.delFrom(xids.list)
}
