/* Copyright(c) 2018-2019 Platina Systems, Inc.
 *
 * This program is free software; you can redistribute it and/or modify it
 * under the terms and conditions of the GNU General Public License,
 * version 2, as published by the Free Software Foundation.
 *
 * This program is distributed in the hope it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
 * FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for
 * more details.
 *
 * You should have received a copy of the GNU General Public License along with
 * this program; if not, write to the Free Software Foundation, Inc.,
 * 51 Franklin St - Fifth Floor, Boston, MA 02110-1301 USA.
 *
 * The full GNU General Public License is included in this distribution in
 * the file called "COPYING".
 *
 * Contact Information:
 * sw@platina.com
 * Platina Systems, 3180 Del La Cruz Blvd, Santa Clara, CA 95054
 */

package xeth

import (
	"fmt"
	"net"
	"sync"
)

type Xeth struct {
	IfInfo
	EthtoolPrivFlags
	EthtoolSettings
	IPNets []*net.IPNet
	Uppers Associates
	Lowers Associates
	Attr   sync.Map
}

func (xeth *Xeth) String() string { return fmt.Sprint(xeth) }

func (xeth *Xeth) Format(f fmt.State, c rune) {
	fmt.Fprint(f, &xeth.IfInfo)
	if xeth.EthtoolPrivFlags != 0 {
		fmt.Fprint(f, " <", xeth.EthtoolPrivFlags, ">")
	}
	if xeth.EthtoolSettings.Speed != 0 {
		fmt.Fprint(f, " speed ", xeth.EthtoolSettings.Speed)
		fmt.Fprint(f, " autoneg ", xeth.EthtoolSettings.Autoneg)
	}
	if xeth.Uppers.NotEmpty() {
		fmt.Fprint(f, " uppers [", &xeth.Uppers, "]")
	}
	if xeth.Lowers.NotEmpty() {
		fmt.Fprint(f, " lowers [", &xeth.Lowers, "]")
	}
	for _, ipnet := range xeth.IPNets {
		fmt.Fprint(f, "\n    ")
		if ipnet.IP.To4() != nil {
			fmt.Fprint(f, "inet ", ipnet)
		} else {
			fmt.Fprint(f, "inet6 ", ipnet)
		}
	}
}

func (xeth *Xeth) cache(args ...interface{}) {
	for _, v := range args {
		switch t := v.(type) {
		case string:
			xeth.dub(t)
		case *MsgChangeUpperXid:
			var upper, lower *Xeth
			upperXid := Xid(t.Upper)
			lowerXid := Xid(t.Lower)
			if xeth.Xid == upperXid {
				upper = xeth
				lower = Of(lowerXid)
			} else {
				upper = Of(upperXid)
				lower = xeth
			}
			if t.Linking > 0 {
				lower.Uppers.Store(upperXid, upper)
				upper.Lowers.Store(lowerXid, lower)
			} else {
				lower.Uppers.Delete(upperXid)
				upper.Lowers.Delete(lowerXid)
			}
		case *MsgIfinfo:
			xeth.dub((*IfName)(&t.Ifname).String())
			xeth.Xid = Xid(t.Xid)
			xeth.DevKind = DevKind(t.Kind)
			xeth.IfIndex = IfIndex(t.Ifindex)
			xeth.Netns = Netns(t.Net)
			xeth.HardwareAddr =
				make(net.HardwareAddr, len(t.Addr[:]))
			copy(xeth.HardwareAddr, t.Addr[:])
			xeth.IfInfo.Flags = net.Flags(t.Flags)
			xeth.IfInfoReason = IfInfoReason(t.Reason)
		case IfInfoReason:
			xeth.IfInfoReason = t
		case net.HardwareAddr:
			copy(xeth.HardwareAddr, t)
		case DevKind:
			xeth.DevKind = t
		case net.Flags:
			xeth.Flags = t
		case Netns:
			xeth.Netns = t
		case *MsgIfa:
			switch t.Event {
			case IFA_ADD:
				xeth.IPNets = append(xeth.IPNets, t.IPNet())
			case IFA_DEL:
				ipnet := t.IPNet()
				n := len(xeth.IPNets)
				for i, x := range xeth.IPNets {
					if x.IP.Equal(ipnet.IP) {
						copy(xeth.IPNets[i:],
							xeth.IPNets[i+1:])
						xeth.IPNets[n-1] = nil
						xeth.IPNets = xeth.IPNets[:n-1]
						break
					}
				}
			}
		case *MsgEthtoolFlags:
			xeth.EthtoolPrivFlags.cache(t)
		case EthtoolPrivFlags:
			xeth.EthtoolPrivFlags.cache(t)
		case *MsgEthtoolSettings:
			xeth.EthtoolSettings.cache(t)
		case Mbps:
			xeth.EthtoolSettings.cache(t)
		case Duplex:
			xeth.EthtoolSettings.cache(t)
		case DevPort:
			xeth.EthtoolSettings.cache(t)
		case Autoneg:
			xeth.EthtoolSettings.cache(t)
		}
	}
}

func (xeth *Xeth) dub(name string) {
	if xeth.Name != name {
		xeth.Name = name
	}
}
