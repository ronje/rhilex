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
	"math"
	"math/rand"
	"time"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rhilex/typex"
)

/*``
*
* XOR 校验
*
 */
func RandomInt(rx typex.Rhilex, uuid string) func(L *lua.LState) int {
	return func(L *lua.LState) int {
		v := L.ToNumber(2)
		L.Push(lua.LNumber(randomInt(int(v))))
		return 1
	}
}

/*
*
*随机数0-N
*
 */
func randomInt(v int) int {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	return rng.Intn(v) + 1
}

func truncateFloat(number float64, decimalPlaces int) float64 {
	scale := math.Pow(10, float64(decimalPlaces))
	result := math.Floor(number*scale) / scale
	return result
}

/*
*
* 取小数位 applib:Float(number, decimalPlaces) -> float
*
 */
func TruncateFloat(rx typex.Rhilex, uuid string) func(L *lua.LState) int {
	return func(l *lua.LState) int {
		number := l.ToNumber(2)
		decimalPlaces := l.ToInt(3)
		l.Push(lua.LNumber(truncateFloat(float64(number), decimalPlaces)))
		return 1
	}
}
