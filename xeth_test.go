// Copyright © 2018-2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xeth

import (
	"flag"
	"fmt"
	"sync"
	"testing"
)

var Continue = flag.Bool("test.continue", false,
	"continue after ifinfo and fib dumps unil SIGINT")

func Test(t *testing.T) {
	flag.Parse()

	err := Init()
	if err != nil {
		t.Fatal(err)
	}

	DumpIfInfo()
	for buf := range RxCh {
		if Class(buf) == ClassBreak {
			break
		}
		msg := Parse(buf)
		fmt.Println(msg)
		Pool(msg)
	}

	Range(func(xid Xid) bool {
		fmt.Println("xid", uint32(xid), xid)
		return true
	})

	DumpFib()
	for buf := range RxCh {
		if Class(buf) == ClassBreak {
			break
		}
		msg := Parse(buf)
		fmt.Println(msg)
		Pool(msg)
	}

	if *Continue {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			for buf := range RxCh {
				msg := Parse(buf)
				fmt.Println(msg)
				Pool(msg)
			}
		}()
		wg.Wait()
	}

	if err := Close(); err != nil {
		if IsSignal(err) {
			t.Log(err)
		} else {
			t.Error(err)
		}
	}
	if val := Cloned.Count(); val != 0 {
		t.Log("cloned", val)
	}
	if val := Parsed.Count(); val != 0 {
		t.Log("parsed", val)
	}
	if val := Dropped.Count(); val != 0 {
		t.Log("dropped", val)
	}
	if val := Sent.Count(); val != 0 {
		t.Log("sent", val)
	}
}
