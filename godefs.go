// +build ignore

package xeth

/*
#include <stdint.h>
#include <linux/types.h>
#include <errno.h>
#include "xeth.h"
*/
import "C"

const (
	VlanVidMask = C.XETH_VLAN_VID_MASK
	VlanNVid    = C.XETH_VLAN_N_VID

	MsgVersion = C.XETH_MSG_VERSION

	IflaUnspec = C.XETH_IFLA_UNSPEC
	IflaXid    = C.XETH_IFLA_XID
	IflaVid    = C.XETH_IFLA_VID
	IflaKind   = C.XETH_IFLA_KIND

	DevKindUnspec = C.XETH_DEV_KIND_UNSPEC
	DevKindPort   = C.XETH_DEV_KIND_PORT
	DevKindVlan   = C.XETH_DEV_KIND_VLAN
	DevKindBridge = C.XETH_DEV_KIND_BRIDGE
	DevKindLag    = C.XETH_DEV_KIND_LAG

	MsgKindBreak           = C.XETH_MSG_KIND_BREAK
	MsgKindLinkStat        = C.XETH_MSG_KIND_LINK_STAT
	MsgKindEthtoolStat     = C.XETH_MSG_KIND_ETHTOOL_STAT
	MsgKindEthtoolFlags    = C.XETH_MSG_KIND_ETHTOOL_FLAGS
	MsgKindEthtoolSettings = C.XETH_MSG_KIND_ETHTOOL_SETTINGS
	MsgKindDumpIfInfo      = C.XETH_MSG_KIND_DUMP_IFINFO
	MsgKindCarrier         = C.XETH_MSG_KIND_CARRIER
	MsgKindSpeed           = C.XETH_MSG_KIND_SPEED
	MsgKindIfInfo          = C.XETH_MSG_KIND_IFINFO
	MsgKindIfa             = C.XETH_MSG_KIND_IFA
	MsgKindDumpFibInfo     = C.XETH_MSG_KIND_DUMP_FIBINFO
	MsgKindFibEntry        = C.XETH_MSG_KIND_FIBENTRY
	MsgKindIfDel           = C.XETH_MSG_KIND_IFDEL
	MsgKindNeighUpdate     = C.XETH_MSG_KIND_NEIGH_UPDATE
	MsgKindIfVid           = C.XETH_MSG_KIND_IFVID
	MsgKindChangeUpperXid  = C.XETH_MSG_KIND_CHANGE_UPPER_XID

	CarrierOff = C.XETH_CARRIER_OFF
	CarrierOn  = C.XETH_CARRIER_ON

	IfInfoReasonNew   = C.XETH_IFINFO_REASON_NEW
	IfInfoReasonDel   = C.XETH_IFINFO_REASON_DEL
	IfInfoReasonUp    = C.XETH_IFINFO_REASON_UP
	IfInfoReasonDown  = C.XETH_IFINFO_REASON_DOWN
	IfInfoReasonDump  = C.XETH_IFINFO_REASON_DUMP
	IfInfoReasonReg   = C.XETH_IFINFO_REASON_REG
	IfInfoReasonUnreg = C.XETH_IFINFO_REASON_UNREG

	SizeofIfName             = C.XETH_IFNAMSIZ
	SizeofEthAddr            = C.XETH_ALEN
	SizeofJumboFrame         = C.XETH_SIZEOF_JUMBO_FRAME
	SizeofMsg                = C.sizeof_struct_xeth_msg
	SizeofMsgBreak           = C.sizeof_struct_xeth_msg_break
	SizeofMsgCarrier         = C.sizeof_struct_xeth_msg_carrier
	SizeofMsgChangeUpperXid  = C.sizeof_struct_xeth_msg_change_upper_xid
	SizeofMsgDumpFibInfo     = C.sizeof_struct_xeth_msg_dump_fibinfo
	SizeofMsgDumpIfInfo      = C.sizeof_struct_xeth_msg_dump_ifinfo
	SizeofMsgEthtoolFlags    = C.sizeof_struct_xeth_msg_ethtool_flags
	SizeofMsgEthtoolSettings = C.sizeof_struct_xeth_msg_ethtool_settings
	SizeofMsgIfa             = C.sizeof_struct_xeth_msg_ifa
	SizeofMsgIfInfo          = C.sizeof_struct_xeth_msg_ifinfo
	SizeofNextHop            = C.sizeof_struct_xeth_next_hop
	SizeofMsgFibEntry        = C.sizeof_struct_xeth_msg_fibentry
	SizeofMsgNeighUpdate     = C.sizeof_struct_xeth_msg_neigh_update
	SizeofMsgSpeed           = C.sizeof_struct_xeth_msg_speed
	SizeofMsgStat            = C.sizeof_struct_xeth_msg_stat
)

type MsgHeader C.struct_xeth_msg_header

type MsgBreak C.struct_xeth_msg_break

type MsgCarrier C.struct_xeth_msg_carrier

type MsgChangeUpperXid C.struct_xeth_msg_change_upper_xid

type MsgEthtoolFlags C.struct_xeth_msg_ethtool_flags

type MsgEthtoolSettings C.struct_xeth_msg_ethtool_settings

type NextHop C.struct_xeth_next_hop

type MsgFibentry C.struct_xeth_msg_fibentry

type MsgIfa C.struct_xeth_msg_ifa

type MsgIfinfo C.struct_xeth_msg_ifinfo

type MsgNeighUpdate C.struct_xeth_msg_neigh_update

type MsgSpeed C.struct_xeth_msg_speed

type MsgStat C.struct_xeth_msg_stat
