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

import "sync"

type Matcher func(*Xeth, ...interface{}) bool

var xeths sync.Map

func Of(xid Xid) *Xeth {
	if v, found := xeths.Load(xid); found {
		return v.(*Xeth)
	}
	return nil
}

// Return list of Xeth's matching the given qualifier which may be,
//	IfIndex
//	string or IfName
//	Matcher(xeth, args...) == true
func Matching(qualifier interface{}, args ...interface{}) (list []*Xeth) {
	xeths.Range(func(k, v interface{}) bool {
		xeth := v.(*Xeth)
		switch t := qualifier.(type) {
		case IfIndex:
			if xeth.IfInfo.IfIndex == t {
				list = append(list, xeth)
			}
		case string:
			if xeth.IfInfo.Name == t {
				list = append(list, xeth)
			}
		case IfName:
			if xeth.IfInfo.Name == t.String() {
				list = append(list, xeth)
			}
		case Matcher:
			if t(xeth, args...) {
				list = append(list, xeth)
			}
		default:
			return false
		}
		return true
	})
	return
}

func Delete(xid Xid) {
	xeths.Delete(xid)
}

func Range(f func(xid Xid, xeth *Xeth) bool) {
	xeths.Range(func(k, v interface{}) bool {
		return f(k.(Xid), v.(*Xeth))
	})
}

func Update(xid Xid, args ...interface{}) *Xeth {
	xeth := Of(xid)
	if xeth == nil {
		xeth = new(Xeth)
		xeth.Xid = xid
		xeths.Store(xid, xeth)
	}
	xeth.cache(args...)
	return xeth
}
