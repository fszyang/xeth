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

import "unsafe"

type LinkStat int
type MsgLinkStat MsgStat

const (
	LinkStatRxPackets LinkStat = iota
	LinkStatTxPackets
	LinkStatRxBytes
	LinkStatTxBytes
	LinkStatRxErrors
	LinkStatTxErrors
	LinkStatRxDropped
	LinkStatTxDropped
	LinkStatMulticast
	LinkStatCollisions
	LinkStatRxLengthErrors
	LinkStatRxOverErrors
	LinkStatRxCrcErrors
	LinkStatRxFrameErrors
	LinkStatRxFifoErrors
	LinkStatRxMissedErrors
	LinkStatTxAbortedErrors
	LinkStatTxCarrierErrors
	LinkStatTxFifoErrors
	LinkStatTxHeartbeatErrors
	LinkStatTxWindowErrors
	LinkStatRxCompressed
	LinkStatTxCompressed
	LinkStatRxNohandler
)

func (stat LinkStat) String() string {
	s, found := map[LinkStat]string{
		LinkStatRxPackets:         "rx-packets",
		LinkStatTxPackets:         "tx-packets",
		LinkStatRxBytes:           "rx-bytes",
		LinkStatTxBytes:           "tx-bytes",
		LinkStatRxErrors:          "rx-errors",
		LinkStatTxErrors:          "tx-errors",
		LinkStatRxDropped:         "rx-dropped",
		LinkStatTxDropped:         "tx-dropped",
		LinkStatMulticast:         "multicast",
		LinkStatCollisions:        "collisions",
		LinkStatRxLengthErrors:    "rx-length-errors",
		LinkStatRxOverErrors:      "rx-over-errors",
		LinkStatRxCrcErrors:       "rx-crc-errors",
		LinkStatRxFrameErrors:     "rx-frame-errors",
		LinkStatRxFifoErrors:      "rx-fifo-errors",
		LinkStatRxMissedErrors:    "rx-missed-errors",
		LinkStatTxAbortedErrors:   "tx-aborted-errors",
		LinkStatTxCarrierErrors:   "tx-carrier-errors",
		LinkStatTxFifoErrors:      "tx-fifo-errors",
		LinkStatTxHeartbeatErrors: "tx-heartbeat-errors",
		LinkStatTxWindowErrors:    "tx-window-errors",
		LinkStatRxCompressed:      "rx-compressed",
		LinkStatTxCompressed:      "tx-compressed",
		LinkStatRxNohandler:       "rx-nohandler",
	}[stat]
	if found {
		return s
	}
	return "invalid-link-stat"
}

func ToMsgLinkStat(buf []byte) *MsgLinkStat {
	return (*MsgLinkStat)(unsafe.Pointer(&buf[0]))
}

func (msg *MsgLinkStat) Set(xid Xid, stat LinkStat, count uint64) {
	msg.Header.Set(MsgKindLinkStat)
	msg.Xid = uint32(xid)
	msg.Index = uint32(stat)
	msg.Count = count
}
