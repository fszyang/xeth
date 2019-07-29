// Copyright © 2018-2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xeth

import (
	"net"
	"sync"
	"syscall"

	"github.com/platinasystems/xeth/internal"
)

type Neighbor struct {
	NetNs
	Xid
	net.IP
	net.HardwareAddr
}

var poolNeighbor = sync.Pool{
	New: func() interface{} {
		return &Neighbor{
			IP: make([]byte, net.IPv6len, net.IPv6len),
			HardwareAddr: make([]byte,
				internal.SizeofEthAddr,
				internal.SizeofEthAddr),
		}
	},
}

func (note *Neighbor) Pool() {
	note.IP = note.IP[:cap(note.IP)]
	poolNeighbor.Put(note)
}

func neighbor(msg *internal.MsgNeighUpdate) *Neighbor {
	note := poolNeighbor.Get().(*Neighbor)
	netns := NetNs(msg.Net)
	note.NetNs = netns
	note.Xid = netns.Xid(msg.Ifindex)
	if msg.Family == syscall.AF_INET {
		copy(note.IP, msg.Dst[:net.IPv4len])
		note.IP = note.IP[:net.IPv4len]
	} else {
		copy(note.IP, msg.Dst[:int(msg.Len)])
		note.IP = note.IP[:cap(note.IP)]
	}
	copy(note.HardwareAddr, msg.Lladdr[:])
	return note
}
