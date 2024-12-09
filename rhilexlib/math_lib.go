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
)

/*
*
*随机数0-N
*
 */
func MathRandomInt(L *lua.LState) int {
	v := L.ToNumber(2)
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	r := rng.Intn(int(v)) + 1
	L.Push(lua.LNumber(r))
	return 1
}

func MathAbs(L *lua.LState) int {
	arg := L.CheckNumber(1)
	L.Push(lua.LNumber(math.Abs(float64(arg))))
	return 1
}

func MathAcos(L *lua.LState) int {
	arg := L.CheckNumber(1)
	L.Push(lua.LNumber(math.Acos(float64(arg))))
	return 1
}

func MathAsin(L *lua.LState) int {
	arg := L.CheckNumber(1)
	L.Push(lua.LNumber(math.Asin(float64(arg))))
	return 1
}

func MathAtan(L *lua.LState) int {
	arg := L.CheckNumber(1)
	L.Push(lua.LNumber(math.Atan(float64(arg))))
	return 1
}

func MathAtan2(L *lua.LState) int {
	y := L.CheckNumber(1)
	x := L.CheckNumber(2)
	L.Push(lua.LNumber(math.Atan2(float64(y), float64(x))))
	return 1
}

func MathCeil(L *lua.LState) int {
	arg := L.CheckNumber(1)
	L.Push(lua.LNumber(math.Ceil(float64(arg))))
	return 1
}

func MathCos(L *lua.LState) int {
	arg := L.CheckNumber(1)
	L.Push(lua.LNumber(math.Cos(float64(arg))))
	return 1
}

func MathExp(L *lua.LState) int {
	arg := L.CheckNumber(1)
	L.Push(lua.LNumber(math.Exp(float64(arg))))
	return 1
}

func MathFloor(L *lua.LState) int {
	arg := L.CheckNumber(1)
	L.Push(lua.LNumber(math.Floor(float64(arg))))
	return 1
}

func MathLog(L *lua.LState) int {
	arg := L.CheckNumber(1)
	L.Push(lua.LNumber(math.Log(float64(arg))))
	return 1
}

func MathMax(L *lua.LState) int {
	if L.GetTop() == 0 {
		L.Push(lua.LNumber(math.Inf(1)))
		return 1
	}
	max := L.CheckNumber(1)
	for i := 2; i <= L.GetTop(); i++ {
		arg := L.CheckNumber(i)
		if arg > max {
			max = arg
		}
	}
	L.Push(lua.LNumber(max))
	return 1
}

func MathMin(L *lua.LState) int {
	if L.GetTop() == 0 {
		L.Push(lua.LNumber(math.Inf(-1)))
		return 1
	}
	min := L.CheckNumber(1)
	for i := 2; i <= L.GetTop(); i++ {
		arg := L.CheckNumber(i)
		if arg < min {
			min = arg
		}
	}
	L.Push(lua.LNumber(min))
	return 1
}

func MathPow(L *lua.LState) int {
	base := L.CheckNumber(1)
	exp := L.CheckNumber(2)
	L.Push(lua.LNumber(math.Pow(float64(base), float64(exp))))
	return 1
}

func MathSin(L *lua.LState) int {
	arg := L.CheckNumber(1)
	L.Push(lua.LNumber(math.Sin(float64(arg))))
	return 1
}

func MathSqrt(L *lua.LState) int {
	arg := L.CheckNumber(1)
	L.Push(lua.LNumber(math.Sqrt(float64(arg))))
	return 1
}

func MathTan(L *lua.LState) int {
	arg := L.CheckNumber(1)
	L.Push(lua.LNumber(math.Tan(float64(arg))))
	return 1
}

func MathRound(L *lua.LState) int {
	arg := L.CheckNumber(1)
	L.Push(lua.LNumber(math.Round(float64(arg))))
	return 1
}
