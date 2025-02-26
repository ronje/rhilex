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
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package rhilexlib

import (
	lua "github.com/hootrhino/gopher-lua"
	en6400 "github.com/hootrhino/rhilex/periphery/en6400"
	"github.com/hootrhino/rhilex/typex"
)

/*
*
  - 打开LED
    local err = en6400:LedOn()
    if err == nil then
    print("LED off")
    end
*/
func EN6400_LedOn(rx typex.Rhilex, uuid string) func(*lua.LState) int {
	return func(l *lua.LState) int {
		err := en6400.EN6400_GPIO231Set(int(1))
		if err != nil {
			l.Push(lua.LString(err.Error()))
		} else {
			l.Push(lua.LNil)
		}
		return 1
	}

}

/*
*
  - 关闭LED
    local err = en6400:LedOff()
    if err == nil then
    print("LED off")
    end

*
*/
func EN6400_LedOff(rx typex.Rhilex, uuid string) func(*lua.LState) int {
	return func(l *lua.LState) int {
		err := en6400.EN6400_GPIO231Set(int(0))
		if err != nil {
			l.Push(lua.LString(err.Error()))
		} else {
			l.Push(lua.LNil)
		}
		return 1
	}

}

/*
*
  - 加速度计
    local acc, err = en6400.GetAccelerator()
    if err == nil then
    print(acc.x, acc.y, acc.z)
    end

*
*/
func EN6400_GetAccelerator(rx typex.Rhilex, uuid string) func(*lua.LState) int {
	return func(l *lua.LState) int {
		Acceleration, err := en6400.ReadAcceleration()
		if err != nil {
			l.Push(lua.LNil)
			l.Push(lua.LString(err.Error()))
		} else {
			t := lua.LTable{}
			t.RawSetString("x", lua.LNumber(Acceleration.X))
			t.RawSetString("y", lua.LNumber(Acceleration.Y))
			t.RawSetString("z", lua.LNumber(Acceleration.Z))
			l.Push(&t)
			l.Push(lua.LNil)
		}
		return 1
	}
}
