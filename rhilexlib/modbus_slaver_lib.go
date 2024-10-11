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
	"fmt"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rhilex/typex"
)

/*
*
* 写入单个线圈 modbus_slaver:F5("${UUID}", 1, 0)
*
 */
func SlaverF5(rx typex.Rhilex) func(*lua.LState) int {
	return func(l *lua.LState) int {
		uuid := l.ToString(2)
		Device := rx.GetDevice(uuid)
		if Device == nil {
			l.Push(lua.LString("Device not exists"))
			return 1
		}
		if Device.Device == nil {
			l.Push(lua.LString("Device not exists"))
			return 1
		}
		addr := l.ToNumber(3)
		if addr > 65535 {
			l.Push(lua.LString("Invalid address"))
			return 1
		}
		value := l.ToNumber(4)
		if value > 1 {
			l.Push(lua.LString("Invalid value"))
			return 1
		}
		if value == 0 {
			_, err := Device.Device.OnCtrl([]byte("CTRL_F5"),
				[]byte(fmt.Sprintf("%d,%d", addr, 0)))
			if err != nil {
				l.Push(lua.LString("Invalid value"))
				return 1
			}

		}
		if value == 1 {
			_, err := Device.Device.OnCtrl([]byte("CTRL_F5"),
				[]byte(fmt.Sprintf("%d,%d", addr, 1)))
			if err != nil {
				l.Push(lua.LString("Invalid value"))
				return 1
			}
		}
		l.Push(lua.LNil)
		return 1
	}
}

/*
*
* 写入保持寄存器 modbus_slaver:F6("${UUID}", 1, 0xABCD)
*
 */
func SlaverF6(rx typex.Rhilex) func(*lua.LState) int {
	return func(l *lua.LState) int {
		uuid := l.ToString(2)
		Device := rx.GetDevice(uuid)
		if Device == nil {
			l.Push(lua.LString("Device not exists"))
			return 1
		}
		if Device.Device == nil {
			l.Push(lua.LString("Device not exists"))
			return 1
		}
		addr := l.ToNumber(3)
		if addr > 0xFFFF {
			l.Push(lua.LString("Invalid address"))
			return 1
		}
		value := l.ToNumber(4)
		if value > 0xFFFF {
			l.Push(lua.LString("Invalid value"))
			return 1
		}
		_, err := Device.Device.OnCtrl([]byte("CTRL_F6"),
			[]byte(fmt.Sprintf("%d,%d", addr, value)))
		if err != nil {
			l.Push(lua.LString("Invalid value"))
			return 1
		}
		l.Push(lua.LNil)
		return 1
	}
}
