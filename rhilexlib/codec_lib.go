// Copyright (C) 2024 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package rhilexlib

import (
	"reflect"

	"github.com/hootrhino/rhilex/typex"

	lua "github.com/hootrhino/gopher-lua"
)

/*
*
* GRPC 解码
*
 */
func Request(rx typex.Rhilex) func(*lua.LState) int {
	return request(rx)
}
func RPCDecode(rx typex.Rhilex) func(*lua.LState) int {
	return request(rx)
}

/*
*
* GRPC 编码
*
 */
func RPCEncode(rx typex.Rhilex) func(*lua.LState) int {
	return request(rx)
}
func request(rx typex.Rhilex) func(*lua.LState) int {
	return func(l *lua.LState) int {
		id := l.ToString(2)               // UUID
		data := l.ToString(3)             // Data
		target := rx.GetOutEnd(id).Target // Codec Target
		// 两个返回值
		// () -> data ,err
		if target.Details().Type != typex.GRPC_CODEC_TARGET {
			l.Push(lua.LNil)                                             // Data
			l.Push(lua.LString("Only support 'GRPC_CODEC_TARGET' type")) // Error
			return 2
		}
		r, err := target.To(data)
		if err != nil {
			l.Push(lua.LNil)                 // Data
			l.Push(lua.LString(err.Error())) // Error
			return 2
		}
		switch t := r.(type) {
		case string:
			l.Push(lua.LString(t)) // Data
			l.Push(lua.LNil)       // Error
			return 2
		case []uint8:
			l.Push(lua.LString(t)) // Data
			l.Push(lua.LNil)       // Error
			return 2
		}
		l.Push(lua.LNil)                                                                        // Data
		l.Push(lua.LString("result must string, but current is:" + reflect.TypeOf(r).String())) // Error
		return 2

	}
}
