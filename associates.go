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
	"sync"
)

// The upper or lower entries associated with an interface entry
type Associates struct {
	sync.Map
}

func (associates *Associates) NotEmpty() bool {
	t := false
	associates.Range(func(k, v interface{}) bool {
		t = true
		return false
	})
	return t
}

func (associates *Associates) String() string {
	return fmt.Sprint(associates)
}

func (associates *Associates) Format(f fmt.State, c rune) {
	sep := false
	associates.Range(func(k, v interface{}) bool {
		xeth := v.(*Xeth)
		if sep {
			f.Write([]byte(", "))
		}
		f.Write([]byte(xeth.Name))
		sep = true
		return true
	})
}
