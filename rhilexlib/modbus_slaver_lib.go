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
* 写入单个线圈
*    modbus_slaver:F5("${UUID}", 1, 0)
*
 */
func SlaverF5(rx typex.Rhilex) func(*lua.LState) int {
	return func(l *lua.LState) int {
		uuid := l.ToString(2)
		addr := l.ToString(3)
		value := l.ToString(4)
		fmt.Println(uuid, addr, value)
		return 1
	}
}

/*
*
* 写入保持寄存器
*
 */
func SlaverF6(rx typex.Rhilex) func(*lua.LState) int {
	return func(l *lua.LState) int {
		return 1
	}
}
