// Copyright © 2018-2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xeth

import (
	"sync"

	"github.com/platinasystems/xeth/internal"
)

type AutoNeg uint8
type Duplex uint8
type DevPort uint8

type EthtoolSettings struct {
	Speed uint32
	AutoNeg
	Duplex
	DevPort

	PhyAddress  uint8
	MdioSupport uint8
	Mdix        uint8
	MdixCtrl    uint8
}

type DevEthtoolSettings Xid

var (
	ethtoolSettings sync.Map

	poolEthtoolSettings = sync.Pool{
		New: func() interface{} {
			return new(EthtoolSettings)
		},
	}
)

func (xid Xid) EthtoolSettings() (settings *EthtoolSettings) {
	if v, ok := ethtoolSettings.Load(xid); ok {
		settings = v.(*EthtoolSettings)
	}
	return
}

func (xid Xid) deleteEthtoolSettings() {
	if settings := xid.EthtoolSettings(); settings != nil {
		poolEthtoolSettings.Put(settings)
		ethtoolSettings.Delete(xid)
	}
}

func (xid Xid) ethtoolSettings(msg *internal.MsgEthtoolSettings) DevEthtoolSettings {
	settings := xid.EthtoolSettings()
	if settings == nil {
		settings = poolEthtoolSettings.Get().(*EthtoolSettings)
		ethtoolSettings.Store(xid, settings)
	}
	settings.Speed = msg.Speed
	settings.AutoNeg = AutoNeg(msg.Autoneg)
	settings.Duplex = Duplex(msg.Duplex)
	settings.DevPort = DevPort(msg.Port)
	settings.PhyAddress = msg.Phy_address
	settings.MdioSupport = msg.Mdio_support
	settings.Mdix = msg.Eth_tp_mdix
	settings.MdixCtrl = msg.Eth_tp_mdix_ctrl
	return DevEthtoolSettings(xid)
}
