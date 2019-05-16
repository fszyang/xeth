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
	"unsafe"
)

const MsgKindInvalid = 0xff

type MsgKind int

func ToMsgHeader(buf []byte) *MsgHeader {
	return (*MsgHeader)(unsafe.Pointer(&buf[0]))
}

func MsgKindOf(buf []byte) MsgKind {
	var kind MsgKind = MsgKindInvalid
	h := ToMsgHeader(buf)
	if len(buf) >= SizeofMsg &&
		h.Z64 == 0 &&
		h.Z32 == 0 &&
		h.Z16 == 0 &&
		h.Version == MsgVersion {
		kind = MsgKind(h.Kind)
	}
	return kind
}

func (kind MsgKind) String() string {
	s, found := map[MsgKind]string{
		MsgKindBreak:           "break",
		MsgKindLinkStat:        "link-stat",
		MsgKindEthtoolStat:     "ethtool-stat",
		MsgKindEthtoolFlags:    "ethtool-flags",
		MsgKindEthtoolSettings: "ethtool-settings",
		MsgKindDumpIfInfo:      "dump-ifinfo",
		MsgKindCarrier:         "carrier",
		MsgKindSpeed:           "speed",
		MsgKindIfInfo:          "ifinfo",
		MsgKindIfa:             "ifa",
		MsgKindDumpFibInfo:     "dump-fibinfo",
		MsgKindFibEntry:        "fib-entry",
		MsgKindIfDel:           "ifdel",
		MsgKindNeighUpdate:     "neigh-update",
		MsgKindIfVid:           "ifvid",
		MsgKindChangeUpperXid:  "change-upper-xid",
	}[kind]
	if found {
		return s
	}
	return "invalid-msg-kind"
}

func (kind MsgKind) update(buf []byte) {
	switch kind {
	case MsgKindChangeUpperXid:
		msg := ToMsgChangeUpperXid(buf)
		xid := Xid(msg.Upper)
		Update(xid, msg)
	case MsgKindIfa:
		msg := ToMsgIfa(buf)
		xid := Xid(msg.Xid)
		Update(xid, msg)
	case MsgKindIfInfo:
		msg := ToMsgIfinfo(buf)
		xid := Xid(msg.Xid)
		switch msg.Reason {
		case IfInfoReasonNew:
			Update(xid, msg)
		case IfInfoReasonDel:
			Delete(xid)
		case IfInfoReasonUp:
			Update(xid, net.Flags(msg.Flags))
		case IfInfoReasonDown:
			Update(xid, net.Flags(msg.Flags))
		case IfInfoReasonDump:
			Update(xid, msg)
		case IfInfoReasonReg:
			if ifentry := Of(xid); ifentry != nil {
				ifentry.Netns = Netns(msg.Net)
			} else {
				Update(xid, msg)
			}
		case IfInfoReasonUnreg:
			if ifentry := Of(xid); ifentry != nil {
				ifentry.Netns = DefaultNetns
			}
		}
	case MsgKindEthtoolFlags:
		msg := ToMsgEthtoolFlags(buf)
		xid := Xid(msg.Xid)
		Update(xid, msg)
	case MsgKindEthtoolSettings:
		msg := ToMsgEthtoolSettings(buf)
		xid := Xid(msg.Xid)
		Update(xid, msg)
	}
}

func (kind MsgKind) validate(buf []byte) error {
	n, found := map[MsgKind]int{
		MsgKindChangeUpperXid:  SizeofMsgChangeUpperXid,
		MsgKindEthtoolFlags:    SizeofMsgEthtoolFlags,
		MsgKindEthtoolSettings: SizeofMsgEthtoolSettings,
		MsgKindIfa:             SizeofMsgIfa,
		MsgKindIfInfo:          SizeofMsgIfInfo,
		MsgKindNeighUpdate:     SizeofMsgNeighUpdate,
	}[kind]
	if (kind == MsgKindFibEntry && len(buf) < SizeofMsgFibEntry) ||
		(found && len(buf) != n) {
		return fmt.Errorf("%s: mismatched length of %d", kind, n)
	} else if !found && kind != MsgKindFibEntry {
		return fmt.Errorf("unsupported msg kind %d", kind)
	}
	return nil
}

func ToMsgCarrier(buf []byte) *MsgCarrier {
	return (*MsgCarrier)(unsafe.Pointer(&buf[0]))
}

func ToMsgChangeUpperXid(buf []byte) *MsgChangeUpperXid {
	return (*MsgChangeUpperXid)(unsafe.Pointer(&buf[0]))
}

func ToMsgEthtoolFlags(buf []byte) *MsgEthtoolFlags {
	return (*MsgEthtoolFlags)(unsafe.Pointer(&buf[0]))
}

func ToMsgEthtoolSettings(buf []byte) *MsgEthtoolSettings {
	return (*MsgEthtoolSettings)(unsafe.Pointer(&buf[0]))
}

func ToMsgIfa(buf []byte) *MsgIfa {
	return (*MsgIfa)(unsafe.Pointer(&buf[0]))
}

func ToMsgIfinfo(buf []byte) *MsgIfinfo {
	return (*MsgIfinfo)(unsafe.Pointer(&buf[0]))
}

func ToMsgNeighUpdate(buf []byte) *MsgNeighUpdate {
	return (*MsgNeighUpdate)(unsafe.Pointer(&buf[0]))
}

func ToMsgSpeed(buf []byte) *MsgSpeed {
	return (*MsgSpeed)(unsafe.Pointer(&buf[0]))
}

func ToMsgStat(buf []byte) *MsgStat {
	return (*MsgStat)(unsafe.Pointer(&buf[0]))
}

func KindOf(buf []byte) MsgKind {
	return MsgKind(ToMsgHeader(buf).Kind)
}

func (h *MsgHeader) Set(kind MsgKind) {
	h.Z64 = 0
	h.Z32 = 0
	h.Z16 = 0
	h.Version = MsgVersion
	h.Kind = uint8(kind)
}
