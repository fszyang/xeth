// +build ignore

package xeth

/*
#include <stdint.h>
#include <linux/types.h>
#include <errno.h>
#include "internal/xeth.h"
*/
import "C"

const (
	VlanVidMask = C.XETH_VLAN_VID_MASK
	VlanNVid    = C.XETH_VLAN_N_VID
)

const (
	DevKindUnspec = C.XETH_DEV_KIND_UNSPEC
	DevKindPort   = C.XETH_DEV_KIND_PORT
	DevKindVlan   = C.XETH_DEV_KIND_VLAN
	DevKindBridge = C.XETH_DEV_KIND_BRIDGE
	DevKindLag    = C.XETH_DEV_KIND_LAG
)

const (
	SHUT_RD = iota
	SHUT_WR
	SHUT_RDWR
)

const (
	FIB_EVENT_ENTRY_REPLACE = iota
	FIB_EVENT_ENTRY_APPEND
	FIB_EVENT_ENTRY_ADD
	FIB_EVENT_ENTRY_DEL
	FIB_EVENT_RULE_ADD
	FIB_EVENT_RULE_DEL
	FIB_EVENT_NH_ADD
	FIB_EVENT_NH_DEL
)

const (
	RTN_UNSPEC = iota
	RTN_UNICAST
	RTN_LOCAL
	RTN_BROADCAST
	RTN_ANYCAST
	RTN_MULTICAST
	RTN_BLACKHOLE
	RTN_UNREACHABLE
	RTN_PROHIBIT
	RTN_THROW
	RTN_NAT
	RTN_XRESOLVE
	__RTN_MAX
)

const RTN_MAX = __RTN_MAX - 1

const (
	RTNH_F_DEAD       = 1  // Nexthop is dead (used by multipath)
	RTNH_F_PERVASIVE  = 2  // Do recursive gateway lookup
	RTNH_F_ONLINK     = 4  // Gateway is forced on link
	RTNH_F_OFFLOAD    = 8  // offloaded route
	RTNH_F_LINKDOWN   = 16 // carrier-down on nexthop
	RTNH_F_UNRESOLVED = 32 // The entry is unresolved (ipmr)
)

const (
	RT_SCOPE_UNIVERSE = 0
	RT_SCOPE_SITE     = 200
	RT_SCOPE_LINK     = 253
	RT_SCOPE_HOST     = 254
	RT_SCOPE_NOWHERE  = 255
)

const (
	RT_TABLE_UNSPEC = 0
	// User defined values
	RT_TABLE_COMPAT  = 252
	RT_TABLE_DEFAULT = 253
	RT_TABLE_MAIN    = 254
	RT_TABLE_LOCAL   = 255
	RT_TABLE_MAX     = 0xFFFFFFFF
)

const (
	LinkStatRxPackets = iota
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

var IndexofLinkStat = map[string]int{
	"rx-packets":          LinkStatRxPackets,
	"tx-packets":          LinkStatTxPackets,
	"rx-bytes":            LinkStatRxBytes,
	"tx-bytes":            LinkStatTxBytes,
	"rx-errors":           LinkStatRxErrors,
	"tx-errors":           LinkStatTxErrors,
	"rx-dropped":          LinkStatRxDropped,
	"tx-dropped":          LinkStatTxDropped,
	"multicast":           LinkStatMulticast,
	"collisions":          LinkStatCollisions,
	"rx-length-errors":    LinkStatRxLengthErrors,
	"rx-over-errors":      LinkStatRxOverErrors,
	"rx-crc-errors":       LinkStatRxCrcErrors,
	"rx-frame-errors":     LinkStatRxFrameErrors,
	"rx-fifo-errors":      LinkStatRxFifoErrors,
	"rx-missed-errors":    LinkStatRxMissedErrors,
	"tx-aborted-errors":   LinkStatTxAbortedErrors,
	"tx-carrier-errors":   LinkStatTxCarrierErrors,
	"tx-fifo-errors":      LinkStatTxFifoErrors,
	"tx-heartbeat-errors": LinkStatTxHeartbeatErrors,
	"tx-window-errors":    LinkStatTxWindowErrors,
	"rx-compressed":       LinkStatRxCompressed,
	"tx-compressed":       LinkStatTxCompressed,
	"rx-nohandler":        LinkStatRxNohandler,
}

const (
	AUTONEG_DISABLE = 0
	AUTONEG_ENABLE  = 1
	AUTONEG_UNKNOWN = 0xff
)

const (
	DUPLEX_HALF    = 0
	DUPLEX_FULL    = 1
	DUPLEX_UNKNOWN = 0xff
)

const (
	PORT_TP    = 0
	PORT_AUI   = 1
	PORT_MII   = 2
	PORT_FIBRE = 3
	PORT_BNC   = 4
	PORT_DA    = 5
	PORT_NONE  = 0xef
	PORT_OTHER = 0xff
)

const (
	ETH_TP_MDI_INVALID = 0 // status: unknown; control: unsupported
	ETH_TP_MDI         = 1 // status: MDI;     control: force MDI
	ETH_TP_MDI_X       = 2 // status: MDI-X;   control: force MDI-X
	ETH_TP_MDI_AUTO    = 3 //                  control: auto-select
)

const (
	ETHTOOL_LINK_MODE_10baseT_Half = iota
	ETHTOOL_LINK_MODE_10baseT_Full
	ETHTOOL_LINK_MODE_100baseT_Half
	ETHTOOL_LINK_MODE_100baseT_Full
	ETHTOOL_LINK_MODE_1000baseT_Half
	ETHTOOL_LINK_MODE_1000baseT_Full
	ETHTOOL_LINK_MODE_Autoneg
	ETHTOOL_LINK_MODE_TP
	ETHTOOL_LINK_MODE_AUI
	ETHTOOL_LINK_MODE_MII
	ETHTOOL_LINK_MODE_FIBRE
	ETHTOOL_LINK_MODE_BNC
	ETHTOOL_LINK_MODE_10000baseT_Full
	ETHTOOL_LINK_MODE_Pause
	ETHTOOL_LINK_MODE_Asym_Pause
	ETHTOOL_LINK_MODE_2500baseX_Full
	ETHTOOL_LINK_MODE_Backplane
	ETHTOOL_LINK_MODE_1000baseKX_Full
	ETHTOOL_LINK_MODE_10000baseKX4_Full
	ETHTOOL_LINK_MODE_10000baseKR_Full
	ETHTOOL_LINK_MODE_10000baseR_FEC
	ETHTOOL_LINK_MODE_20000baseMLD2_Full
	ETHTOOL_LINK_MODE_20000baseKR2_Full
	ETHTOOL_LINK_MODE_40000baseKR4_Full
	ETHTOOL_LINK_MODE_40000baseCR4_Full
	ETHTOOL_LINK_MODE_40000baseSR4_Full
	ETHTOOL_LINK_MODE_40000baseLR4_Full
	ETHTOOL_LINK_MODE_56000baseKR4_Full
	ETHTOOL_LINK_MODE_56000baseCR4_Full
	ETHTOOL_LINK_MODE_56000baseSR4_Full
	ETHTOOL_LINK_MODE_56000baseLR4_Full
	ETHTOOL_LINK_MODE_25000baseCR_Full
	ETHTOOL_LINK_MODE_25000baseKR_Full
	ETHTOOL_LINK_MODE_25000baseSR_Full
	ETHTOOL_LINK_MODE_50000baseCR2_Full
	ETHTOOL_LINK_MODE_50000baseKR2_Full
	ETHTOOL_LINK_MODE_100000baseKR4_Full
	ETHTOOL_LINK_MODE_100000baseSR4_Full
	ETHTOOL_LINK_MODE_100000baseCR4_Full
	ETHTOOL_LINK_MODE_100000baseLR4_ER4_Full
	ETHTOOL_LINK_MODE_50000baseSR2_Full
	ETHTOOL_LINK_MODE_1000baseX_Full
	ETHTOOL_LINK_MODE_10000baseCR_Full
	ETHTOOL_LINK_MODE_10000baseSR_Full
	ETHTOOL_LINK_MODE_10000baseLR_Full
	ETHTOOL_LINK_MODE_10000baseLRM_Full
	ETHTOOL_LINK_MODE_10000baseER_Full
	ETHTOOL_LINK_MODE_2500baseT_Full
	ETHTOOL_LINK_MODE_5000baseT_Full
	ETHTOOL_LINK_MODE_NBITS
)