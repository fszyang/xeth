// Copyright © 2018-2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xeth

import "sync"

var attrs sync.Map

func (xid Xid) Attrs() *sync.Map {
	if v, ok := attrs.Load(xid); ok {
		return v.(*sync.Map)
	}
	m := new(sync.Map)
	attrs.Store(xid, m)
	return m
}

func (xid Xid) deleteAttrs() {
	if v, ok := attrs.Load(xid); ok {
		defer attrs.Delete(xid)
		m := v.(*sync.Map)
		m.Range(func(key, value interface{}) bool {
			defer m.Delete(key)
			if method, found := value.(pooler); found {
				method.Pool()
			}
			return true
		})
	}
}
