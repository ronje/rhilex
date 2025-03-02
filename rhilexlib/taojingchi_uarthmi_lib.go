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
	"encoding/json"

	"github.com/hootrhino/rhilex/typex"

	lua "github.com/hootrhino/gopher-lua"
)

/*
*

	local err = tjchmi:WriteToHmi("$uuid", {
		"t0.txt=\"Hello-1\"",
		"t1.txt=\"Hello-2\"",
		"t2.txt=\"Hello-3\""
	})

*
*/
func TJCWriteToHmi(rx typex.Rhilex) func(*lua.LState) int {
	return func(l *lua.LState) int {
		devUUID := l.ToString(2)
		data := l.ToTable(3)
		Device := rx.GetDevice(devUUID)
		args := []string{}
		data.ForEach(func(l1, l2 lua.LValue) {
			args = append(args, l2.String())
		})
		if Device != nil {
			if Device.Device.Status() == typex.SOURCE_UP {
				bytes, errMarshal := json.Marshal(args)
				if errMarshal != nil {
					l.Push(lua.LString(errMarshal.Error()))
					return 1
				}
				_, err := Device.Device.OnCtrl([]byte("WriteToHmi"), bytes)
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
