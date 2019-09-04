// Copyright © 2018-2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xeth

import "github.com/platinasystems/xeth/internal"

type AutoNeg uint8
type Duplex uint8
type DevPort uint8

type DevEthtoolSettings Xid

func (xid Xid) RxEthtoolSettings(msg *internal.MsgEthtoolSettings) DevEthtoolSettings {
	attrs := LinkAttrs(xid)
	attrs.EthtoolSpeed(msg.Speed)
	attrs.EthtoolAutoNeg(AutoNeg(msg.Autoneg))
	attrs.EthtoolDuplex(Duplex(msg.Duplex))
	attrs.EthtoolDevPort(DevPort(msg.Port))
	return DevEthtoolSettings(xid)
}
