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

/**
 * 动作成功
 *
 */
func IthingsActionReplySuccess(rx typex.Rhilex) func(*lua.LState) int {
	return func(stateStack *lua.LState) int {
		uuid := stateStack.ToString(2)
		token := stateStack.ToString(3)
		Device := rx.GetDevice(uuid)
		if Device != nil {
			if Device.Device != nil {
				_, err := Device.Device.OnWrite([]byte("ActionReplySuccess"), []byte(token))
				if err != nil {
					stateStack.Push(lua.LString(err.Error()))
					return 1
				}
			}
		}
		stateStack.Push(lua.LNil)
		return 1
	}
}

/**
 * 动作失败
 *
 */
func IthingsActionReplyFailure(rx typex.Rhilex) func(*lua.LState) int {
	return func(stateStack *lua.LState) int {
		uuid := stateStack.ToString(2)
		token := stateStack.ToString(3)
		Device := rx.GetDevice(uuid)
		if Device != nil {
			if Device.Device != nil {
				_, err := Device.Device.OnWrite([]byte("ActionReplyFailure"), []byte(token))
				if err != nil {
					stateStack.Push(lua.LString(err.Error()))
					return 1
				}
			}
		}
		stateStack.Push(lua.LNil)
		return 1
	}
}

/**
 * 属性成功
 *
 */
func IthingsPropertyReplySuccess(rx typex.Rhilex) func(*lua.LState) int {
	return func(stateStack *lua.LState) int {
		uuid := stateStack.ToString(2)
		token := stateStack.ToString(3)
		Device := rx.GetDevice(uuid)
		if Device != nil {
			if Device.Device != nil {
				_, err := Device.Device.OnWrite([]byte("PropertyReplySuccess"), []byte(token))
				if err != nil {
					stateStack.Push(lua.LString(err.Error()))
					return 1
				}
			}
		}
		stateStack.Push(lua.LNil)
		return 1
	}
}

/**
 * 属性失败
 *
 */
func IthingsPropertyReplyFailure(rx typex.Rhilex) func(*lua.LState) int {
	return func(stateStack *lua.LState) int {
		uuid := stateStack.ToString(2)
		token := stateStack.ToString(3)
		Device := rx.GetDevice(uuid)
		if Device != nil {
			if Device.Device != nil {
				_, err := Device.Device.OnWrite([]byte("PropertyReplyFailure"), []byte(token))
				if err != nil {
					stateStack.Push(lua.LString(err.Error()))
					return 1
				}
			}
		}
		stateStack.Push(lua.LNil)
		return 1
	}
}

/**
 * 上传属性
 *
 */
func IthingsPropertyReport(rx typex.Rhilex) func(*lua.LState) int {
	return func(stateStack *lua.LState) int {
		uuid := stateStack.ToString(2)
		paramsTable := stateStack.ToTable(3)
		params := map[string]interface{}{}
		paramsTable.ForEach(func(k, v lua.LValue) {
			params[k.String()] = v
		})
		Device := rx.GetDevice(uuid)
		if Device != nil {
			if Device.Device != nil {
				bytes, errMarshal := json.Marshal(params)
				if errMarshal != nil {
					stateStack.Push(lua.LString(errMarshal.Error()))
					return 1
				}
				_, err := Device.Device.OnWrite([]byte("PropertyReport"), []byte(bytes))
				if err != nil {
					stateStack.Push(lua.LString(err.Error()))
					return 1
				}
			}
		}
		stateStack.Push(lua.LNil)
		return 1
	}
}
