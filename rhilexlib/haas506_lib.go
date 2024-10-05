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
// LED2
func HAAS506_Led2On(rx typex.Rhilex) func(*lua.LState) int {
	return func(l *lua.LState) int {
		err := archsupport.HAAS506_LEDSet(int(2), 1)
		if err != nil {
			l.Push(lua.LString(err.Error()))
		} else {
			l.Push(lua.LNil)
		}
		return 1
	}

}

// LED3
func HAAS506_Led3On(rx typex.Rhilex) func(*lua.LState) int {
	return func(l *lua.LState) int {
		err := archsupport.HAAS506_LEDSet(int(3), 1)
		if err != nil {
			l.Push(lua.LString(err.Error()))
		} else {
			l.Push(lua.LNil)
		}
		return 1
	}

}

// LED4
func HAAS506_Led4On(rx typex.Rhilex) func(*lua.LState) int {
	return func(l *lua.LState) int {
		err := archsupport.HAAS506_LEDSet(int(4), 1)
		if err != nil {
			l.Push(lua.LString(err.Error()))
		} else {
			l.Push(lua.LNil)
		}
		return 1
	}

}

// LED5
func HAAS506_Led5On(rx typex.Rhilex) func(*lua.LState) int {
	return func(l *lua.LState) int {
		err := archsupport.HAAS506_LEDSet(int(5), 1)
		if err != nil {
			l.Push(lua.LString(err.Error()))
		} else {
			l.Push(lua.LNil)
		}
		return 1
	}

}

// Off
// LED2
func HAAS506_Led2Off(rx typex.Rhilex) func(*lua.LState) int {
	return func(l *lua.LState) int {
		err := archsupport.HAAS506_LEDSet(int(2), 0)
		if err != nil {
			l.Push(lua.LString(err.Error()))
		} else {
			l.Push(lua.LNil)
		}
		return 1
	}

}

// LED3
func HAAS506_Led3Off(rx typex.Rhilex) func(*lua.LState) int {
	return func(l *lua.LState) int {
		err := archsupport.HAAS506_LEDSet(int(3), 0)
		if err != nil {
			l.Push(lua.LString(err.Error()))
		} else {
			l.Push(lua.LNil)
		}
		return 1
	}

}

// LED4
func HAAS506_Led4Off(rx typex.Rhilex) func(*lua.LState) int {
	return func(l *lua.LState) int {
		err := archsupport.HAAS506_LEDSet(int(4), 0)
		if err != nil {
			l.Push(lua.LString(err.Error()))
		} else {
			l.Push(lua.LNil)
		}
		return 1
	}

}

// LED5
func HAAS506_Led5Off(rx typex.Rhilex) func(*lua.LState) int {
	return func(l *lua.LState) int {
		err := archsupport.HAAS506_LEDSet(int(5), 0)
		if err != nil {
			l.Push(lua.LString(err.Error()))
		} else {
			l.Push(lua.LNil)
		}
		return 1
	}

}
