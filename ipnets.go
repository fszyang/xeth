// Copyright © 2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

package xeth

import (
	"net"
	"sync"
	"unsafe"
)

type DevAddIPNet struct {
	Xid
	*net.IPNet
}

type DevDelIPNet struct {
	Xid
	Prefix string
}

var (
	ipnets sync.Map

	poolIPNet = sync.Pool{
		New: func() interface{} {
			return &net.IPNet{
				IP:   make([]byte, net.IPv6len, net.IPv6len),
				Mask: make([]byte, net.IPv6len, net.IPv6len),
			}
		},
	}
)

func (xid Xid) IPNets() (l []*net.IPNet) {
	if v, ok := ipnets.Load(xid); ok {
		l = v.([]*net.IPNet)
	}
	return
}

func (xid Xid) deleteIPNets() {
	for _, entry := range xid.IPNets() {
		entry.IP = entry.IP[:cap(entry.IP)]
		entry.Mask = entry.Mask[:cap(entry.Mask)]
		poolIPNet.Put(entry)
	}
	ipnets.Delete(xid)
}

func (xid Xid) addIP(addr, mask uint32) *DevAddIPNet {
	ip := net.IP(make([]byte, net.IPv4len, net.IPv4len))
	*(*uint32)(unsafe.Pointer(&ip[0])) = addr
	l := xid.IPNets()
	for _, entry := range l {
		if ip.Equal(entry.IP) {
			return &DevAddIPNet{xid, entry}
		}
	}
	clone := poolIPNet.Get().(*net.IPNet)
	*(*uint32)(unsafe.Pointer(&clone.IP[0])) = addr
	*(*uint32)(unsafe.Pointer(&clone.Mask[0])) = mask
	clone.IP = clone.IP[:net.IPv4len]
	clone.Mask = clone.Mask[:net.IPv4len]
	ipnets.Store(xid, append(l, clone))
	return &DevAddIPNet{xid, clone}
}

func (xid Xid) delIP(addr, mask uint32) *DevDelIPNet {
	ip := net.IP(make([]byte, net.IPv4len, net.IPv4len))
	*(*uint32)(unsafe.Pointer(&ip[0])) = addr
	l := xid.IPNets()
	for i, entry := range l {
		if ip.Equal(entry.IP) {
			prefix := entry.String()
			n := len(l) - 1
			copy(l[i:], l[i+1:])
			ipnets.Store(xid, l[:n])
			entry.IP = entry.IP[:cap(entry.IP)]
			entry.Mask = entry.Mask[:cap(entry.Mask)]
			poolIPNet.Put(entry)
			return &DevDelIPNet{xid, prefix}
		}
	}
	return nil
}

func (xid Xid) addIP6(addr []byte, len int) *DevAddIPNet {
	ip := net.IP(addr)
	l := xid.IPNets()
	for _, entry := range l {
		if ip.Equal(entry.IP) {
			return &DevAddIPNet{xid, entry}
		}
	}
	clone := poolIPNet.Get().(*net.IPNet)
	copy(clone.IP, ip)
	copy(clone.Mask, net.CIDRMask(len, net.IPv6len*8))
	ipnets.Store(xid, append(l, clone))
	return &DevAddIPNet{xid, clone}
}

func (xid Xid) delIP6(addr []byte) *DevDelIPNet {
	ip := net.IP(addr)
	l := xid.IPNets()
	for i, entry := range l {
		if ip.Equal(entry.IP) {
			prefix := entry.String()
			n := len(l) - 1
			copy(l[i:], l[i+1:])
			ipnets.Store(xid, l[:n])
			entry.IP = entry.IP[:cap(entry.IP)]
			entry.Mask = entry.Mask[:cap(entry.Mask)]
			poolIPNet.Put(entry)
			return &DevDelIPNet{xid, prefix}
		}
	}
	return nil
}
