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
	transceiver "github.com/hootrhino/rhilex/component/transceivercom/transceiver"
	"time"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rhilex/typex"
)

/*
*
* 向系统的通信模组发送数据
*
 */
func CtrlComRF(rx typex.Rhilex) func(*lua.LState) int {
	return func(l *lua.LState) int {
		Name := l.ToString(2)
		Topic := l.ToString(3)
		Args := l.ToString(4)
		Result, Err := transceiver.Ctrl(Name, []byte(Topic), []byte(Args), 300*time.Millisecond)
		if Err != nil {
			l.Push(lua.LNil)
			l.Push(lua.LString(Err.Error()))
		} else {
			l.Push(lua.LString(Result))
			l.Push(lua.LNil)
		}
		return 2
	}
}
