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
	archsupport "github.com/hootrhino/rhilex/archsupport"
	"github.com/hootrhino/rhilex/typex"
)

// On
// DO1
func HAAS506_DO1_On(rx typex.Rhilex) func(*lua.LState) int {
	return func(l *lua.LState) int {
		err := archsupport.HAAS506_GPIOSetDO1(1)
		if err != nil {
			l.Push(lua.LString(err.Error()))
		} else {
			l.Push(lua.LNil)
		}
		return 1
	}
}
func HAAS506_Do2_On(rx typex.Rhilex) func(*lua.LState) int {
	return func(l *lua.LState) int {
		err := archsupport.HAAS506_GPIOSetDO2(1)
		if err != nil {
			l.Push(lua.LString(err.Error()))
		} else {
			l.Push(lua.LNil)
		}
		return 1
	}
}
func HAAS506_Do3_On(rx typex.Rhilex) func(*lua.LState) int {
	return func(l *lua.LState) int {
		err := archsupport.HAAS506_GPIOSetDO3(1)
		if err != nil {
			l.Push(lua.LString(err.Error()))
		} else {
			l.Push(lua.LNil)
		}
		return 1
	}
}
func HAAS506_Do4_On(rx typex.Rhilex) func(*lua.LState) int {
	return func(l *lua.LState) int {
		err := archsupport.HAAS506_GPIOSetDO4(1)
		if err != nil {
			l.Push(lua.LString(err.Error()))
		} else {
			l.Push(lua.LNil)
		}
		return 1
	}
}

// Off
// DO1
func HAAS506_DO1_Off(rx typex.Rhilex) func(*lua.LState) int {
	return func(l *lua.LState) int {
		err := archsupport.HAAS506_GPIOSetDO1(0)
		if err != nil {
			l.Push(lua.LString(err.Error()))
		} else {
			l.Push(lua.LNil)
		}
		return 1
	}
}
func HAAS506_Do2_Off(rx typex.Rhilex) func(*lua.LState) int {
	return func(l *lua.LState) int {
		err := archsupport.HAAS506_GPIOSetDO2(0)
		if err != nil {
			l.Push(lua.LString(err.Error()))
		} else {
			l.Push(lua.LNil)
		}
		return 1
	}
}
func HAAS506_Do3_Off(rx typex.Rhilex) func(*lua.LState) int {
	return func(l *lua.LState) int {
		err := archsupport.HAAS506_GPIOSetDO3(0)
		if err != nil {
			l.Push(lua.LString(err.Error()))
		} else {
			l.Push(lua.LNil)
		}
		return 1
	}
}
func HAAS506_Do4_Off(rx typex.Rhilex) func(*lua.LState) int {
	return func(l *lua.LState) int {
		err := archsupport.HAAS506_GPIOSetDO4(0)
		if err != nil {
			l.Push(lua.LString(err.Error()))
		} else {
			l.Push(lua.LNil)
		}
		return 1
	}
}
