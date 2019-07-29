// Copyright © 2018-2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xeth

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"syscall"
)

type NetNs uint64
type NetNsAttr uint8
type NetNsAttrs sync.Map

const DefaultNetNs NetNs = 1

const (
	PathNetNsAttr NetNsAttr = iota
	XidOfIfIndexNetNsAttr
)

var NetNsAttrMaps sync.Map

func NetNsRange(f func(ns NetNs) bool) {
	NetNsAttrMaps.Range(func(k, v interface{}) bool {
		return f(k.(NetNs))
	})
}

func (ns NetNs) Valid() bool {
	_, ok := NetNsAttrMaps.Load(ns)
	return ok
}

// get mapped attrs but panic if unavailable
func (ns NetNs) Attrs() (attrs *NetNsAttrs) {
	if v, ok := NetNsAttrMaps.Load(ns); ok {
		attrs = (*NetNsAttrs)(v.(*sync.Map))
	} else if true {
		panic(fmt.Errorf("netns %d hasn't been mapped", uint64(ns)))
	}
	return
}

func (ns NetNs) String() string {
	return ns.Base()
}

func (ns NetNs) Base() string {
	return filepath.Base(ns.Path())
}

func (ns NetNs) Path() string {
	if !ns.Valid() {
		return ns.path()
	}
	return ns.Attrs().Path()
}

func (ns NetNs) path() string {
	if ns == DefaultNetNs {
		return "default"
	}
	var nspath string
	filepath.Walk("/run",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if len(nspath) > 0 {
				return filepath.SkipDir
			}
			stat := info.Sys().(*syscall.Stat_t)
			if stat.Ino == uint64(ns) {
				nspath = filepath.Join(path, info.Name())
				return filepath.SkipDir
			}
			return nil
		})
	if len(nspath) == 0 {
		nspath = fmt.Sprintf("not-found(%#x)", uint64(ns))
	}
	return nspath
}

func (ns NetNs) XidOfIfIndexMap() (m *sync.Map) {
	attrs := ns.attrs()
	if v, ok := attrs.Map().Load(XidOfIfIndexNetNsAttr); ok {
		m = v.(*sync.Map)
	} else {
		m = new(sync.Map)
		attrs.Map().Store(XidOfIfIndexNetNsAttr, m)
	}
	return
}

func (attrs *NetNsAttrs) String() string {
	return attrs.Base()
}

func (attrs *NetNsAttrs) Base() string {
	return filepath.Base(attrs.Path())
}

func (attrs *NetNsAttrs) Path() (path string) {
	if v, ok := attrs.Map().Load(PathNetNsAttr); ok {
		path = v.(string)
	}
	return
}

func (attrs *NetNsAttrs) Map() *sync.Map {
	return (*sync.Map)(attrs)
}

// make the xid's attrs map and path if it's not already available
func (ns NetNs) attrs() *NetNsAttrs {
	if v, ok := NetNsAttrMaps.Load(ns); ok {
		return (*NetNsAttrs)(v.(*sync.Map))
	}
	m := new(sync.Map)
	NetNsAttrMaps.Store(ns, m)
	m.Store(PathNetNsAttr, ns.path())
	return (*NetNsAttrs)(m)
}
