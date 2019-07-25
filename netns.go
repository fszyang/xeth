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

const DefaultNetNs NetNs = 1

var nscache struct {
	sync.Mutex
	path map[NetNs]string
}

func (ns NetNs) Path() string {
	nscache.Lock()
	defer nscache.Unlock()
	if nscache.path == nil {
		nscache.path = map[NetNs]string{
			DefaultNetNs: "default",
		}
	}
	nspath, ok := nscache.path[ns]
	if ok {
		return nspath
	}
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
	nscache.path[ns] = nspath
	return nspath
}

func (ns NetNs) Base() string   { return filepath.Base(ns.Path()) }
func (ns NetNs) String() string { return ns.Base() }
