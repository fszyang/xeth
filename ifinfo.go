/* Copyright(c) 2018 Platina Systems, Inc.
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
)

type IfIndex int32
type IfLinkIndex int32

type IfInfo struct {
	Name string
	Xid
	DevKind
	IfIndex
	IfLinkIndex
	Netns
	IfInfoReason
	net.Flags
	net.HardwareAddr
}

func (ifinfo *IfInfo) String() string { return fmt.Sprint(ifinfo) }

func (ifinfo *IfInfo) Format(f fmt.State, c rune) {
	fmt.Fprint(f, ifinfo.IfIndex, ": ", ifinfo.Name)
	if ifinfo.IfLinkIndex > 0 {
		links := Matching(ifinfo.IfLinkIndex)
		if len(links) == 1 {
			fmt.Fprint(f, "@", links[0].IfInfo.Name)
		} else {
			fmt.Fprint(f, "@", ifinfo.IfLinkIndex)
		}
	}
	fmt.Fprint(f, ": xid ", ifinfo.Xid)
	if ifinfo.Flags != 0 {
		fmt.Fprint(f, " <", ifinfo.Flags, ">")
	}
	fmt.Fprint(f, " reason ", ifinfo.IfInfoReason)
	if ifinfo.Netns != DefaultNetns {
		fmt.Fprint(f, " netns ", ifinfo.Netns)
	}
	fmt.Fprint(f, "\n    link/", ifinfo.DevKind)
	fmt.Fprint(f, " ", ifinfo.HardwareAddr)
}
