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
	"time"

	"github.com/hootrhino/rhilex/typex"

	lua "github.com/hootrhino/gopher-lua"
)

/*
*
local err = tjchmi:WriteToHmi("WriteToHmi", "t0.txt=\"Hello\"")
*
*
*/
func TJCWriteToHmi(rx typex.Rhilex) func(*lua.LState) int {
	return func(l *lua.LState) int {
		// write(uuid,cmd,data)
		devUUID := l.ToString(2)
		cmd := l.ToString(3)
		data := l.ToString(4)
		Device := rx.GetDevice(devUUID)
		if Device != nil {
			if Device.Device.Status() == typex.DEV_UP {
				_, err := Device.Device.OnCtrl([]byte(cmd), []byte(data))
				if err != nil {
					l.Push(lua.LString(err.Error()))
					return 1
				}
				l.Push(lua.LNil)
				return 1
			}
			l.Push(lua.LString("device down:" + devUUID))
			return 1
		}
		l.Push(lua.LString("device not exists:" + devUUID))
		return 1
	}
}
