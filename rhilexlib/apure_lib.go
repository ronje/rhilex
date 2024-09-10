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
	"encoding/binary"
	"encoding/hex"
	"math"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rhilex/typex"
)

/*
*
* 上海阔思科技的非主流传感器解析支持, 小端模式
* local Value = apure:ParseDOxygen("0001") -> 12.56
*
 */
func ApureParseOxygen(rx typex.Rhilex) func(*lua.LState) int {
	return func(L *lua.LState) int {
		hexValue := L.ToString(2)
		Byte, err := hex.DecodeString(hexValue)
		if err != nil {
			L.Push(lua.LNumber(0))
			return 1
		}
		if len(Byte) != 4 {
			L.Push(lua.LNumber(0))
			return 1
		}
		// 00 01 02 03
		b1 := [2]byte{Byte[1], Byte[0]}
		b2 := [2]byte{Byte[3], Byte[2]}
		Value := binary.LittleEndian.Uint16(b1[:])
		Round := binary.LittleEndian.Uint16(b2[:])
		if Round == 1 {
			L.Push(lua.LNumber(ApureTruncateFloat(float64(Value)*0.1, 1)))
			return 1
		}
		if Round == 2 {
			L.Push(lua.LNumber(ApureTruncateFloat(float64(Value)*0.1, 2)))
			return 1
		}
		if Round == 3 {
			L.Push(lua.LNumber(ApureTruncateFloat(float64(Value)*0.1, 3)))
			return 1
		}
		if Round == 4 {
			L.Push(lua.LNumber(ApureTruncateFloat(float64(Value)*0.1, 4)))
			return 1
		}
		L.Push(lua.LNumber(0))
		return 1
	}
}

// ApureTruncateFloat 截断一个浮点数到指定的小数位数，不进行四舍五入
// number: 需要截断的浮点数
// decimalPlaces: 要保留的小数位数
func ApureTruncateFloat(number float64, decimalPlaces int) float64 {
	shift := math.Pow(10, float64(decimalPlaces))
	return math.Trunc(number*shift) / shift
}
