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
	"math"

	lua "github.com/hootrhino/gopher-lua"
	haas506 "github.com/hootrhino/rhilex/archsupport/haas506"
	"github.com/hootrhino/rhilex/typex"
)

// AI1
func HAAS506_AI1_Get(rx typex.Rhilex) func(*lua.LState) int {
	return func(l *lua.LState) int {
		V, err := haas506.HAAS506_GPIOGetAI1()
		if err != nil {
			l.Push(lua.LNil)
			l.Push(lua.LString(err.Error()))
		} else {
			l.Push(lua.LNumber(HAAS506_AI_RoundFloat32(float64(V), 2)))
			l.Push(lua.LNil)
		}
		return 2
	}
}

// AI2
func HAAS506_AI2_Get(rx typex.Rhilex) func(*lua.LState) int {
	return func(l *lua.LState) int {
		V, err := haas506.HAAS506_GPIOGetAI2()
		if err != nil {
			l.Push(lua.LNil)
			l.Push(lua.LString(err.Error()))
		} else {
			l.Push(lua.LNumber(HAAS506_AI_RoundFloat32(float64(V), 2)))
			l.Push(lua.LNil)
		}
		return 2
	}
}

// AI3
func HAAS506_AI3_Get(rx typex.Rhilex) func(*lua.LState) int {
	return func(l *lua.LState) int {
		V, err := haas506.HAAS506_GPIOGetAI3()
		if err != nil {
			l.Push(lua.LNil)
			l.Push(lua.LString(err.Error()))
		} else {
			l.Push(lua.LNumber(HAAS506_AI_RoundFloat32(float64(V), 2)))
			l.Push(lua.LNil)
		}
		return 2
	}
}

// AI4
func HAAS506_AI4_Get(rx typex.Rhilex) func(*lua.LState) int {
	return func(l *lua.LState) int {
		V, err := haas506.HAAS506_GPIOGetAI4()
		if err != nil {
			l.Push(lua.LNil)
			l.Push(lua.LString(err.Error()))
		} else {
			l.Push(lua.LNumber(HAAS506_AI_RoundFloat32(float64(V), 2)))
			l.Push(lua.LNil)
		}
		return 2
	}
}

// AI5
func HAAS506_AI5_Get(rx typex.Rhilex) func(*lua.LState) int {
	return func(l *lua.LState) int {
		V, err := haas506.HAAS506_GPIOGetAI5()
		if err != nil {
			l.Push(lua.LNil)
			l.Push(lua.LString(err.Error()))
		} else {
			l.Push(lua.LNumber(HAAS506_AI_RoundFloat32(float64(V), 2)))
			l.Push(lua.LNil)
		}
		return 2
	}
}

// 取2位小数
func HAAS506_AI_RoundFloat32(number float64, decimalPlaces int) float64 {
	scale := math.Pow(10, float64(decimalPlaces))
	result := math.Floor(number/100*scale) / scale
	return result
}
