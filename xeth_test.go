// Copyright © 2018-2019 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xeth

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
)

var Continue = flag.Bool("test.continue", false,
	"continue after ifinfo and fib dumps unil SIGINT")

func Test(t *testing.T) {
	var wg sync.WaitGroup
	defer wg.Wait()
	stopch := make(chan struct{})
	defer close(stopch)

	flag.Parse()

	task, err := Start(&wg, stopch)
	if err != nil {
		t.Fatal(err)
	}

	task.DumpIfInfo()
	for buf := range task.RxCh {
		if Class(buf) == ClassBreak {
			break
		}
		// Load the attribute cache through Parse
		Pool(Parse(buf))
	}

	Range(func(xid Xid) bool {
		attrs := xid.Attrs()
		fmt.Print(attrs.IfInfoName(),
			", xid ", uint32(xid),
			", ifindex ", attrs.IfInfoIfIndex(),
			", netns ", attrs.IfInfoNetNs(),
			", ipnets ", attrs.IPNets(),
			"\n")
		return true
	})

	task.DumpFib()
	for buf := range task.RxCh {
		if Class(buf) == ClassBreak {
			break
		}
		msg := Parse(buf)
		fmt.Println(msg)
		Pool(msg)
	}

	if *Continue {
		sigch := make(chan os.Signal, 1)
		signal.Notify(sigch,
			syscall.SIGTERM,
			syscall.SIGINT,
			syscall.SIGHUP,
			syscall.SIGQUIT)
		go func() {
			wg.Add(1)
			defer wg.Done()
			for buf := range task.RxCh {
				msg := Parse(buf)
				fmt.Println(msg)
				Pool(msg)
			}
		}()
		<-sigch
		signal.Stop(sigch)
	}
}
