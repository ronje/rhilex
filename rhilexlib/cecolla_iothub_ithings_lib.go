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
 * 控制指令
 *
 */
func IthingsCtrlReplySuccess(rx typex.Rhilex) func(*lua.LState) int {
	return func(stateStack *lua.LState) int {
		uuid := stateStack.ToString(2)
		token := stateStack.ToString(3)
		Cecolla := rx.GetCecolla(uuid)
		if Cecolla != nil {
			if Cecolla.Cecolla != nil {
				_, err := Cecolla.Cecolla.OnCtrl([]byte("CtrlReplySuccess"), []byte(token))
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
func IthingsCtrlReplyFailure(rx typex.Rhilex) func(*lua.LState) int {
	return func(stateStack *lua.LState) int {
		uuid := stateStack.ToString(2)
		token := stateStack.ToString(3)
		Cecolla := rx.GetCecolla(uuid)
		if Cecolla != nil {
			if Cecolla.Cecolla != nil {
				_, err := Cecolla.Cecolla.OnCtrl([]byte("CtrlReplyFailure"), []byte(token))
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
 * 动作成功
 *
 */
func IthingsActionReplySuccess(rx typex.Rhilex) func(*lua.LState) int {
	return func(stateStack *lua.LState) int {
		uuid := stateStack.ToString(2)
		token := stateStack.ToString(3)
		Cecolla := rx.GetCecolla(uuid)
		if Cecolla != nil {
			if Cecolla.Cecolla != nil {
				_, err := Cecolla.Cecolla.OnCtrl([]byte("ActionReplySuccess"), []byte(token))
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
		Cecolla := rx.GetCecolla(uuid)
		if Cecolla != nil {
			if Cecolla.Cecolla != nil {
				_, err := Cecolla.Cecolla.OnCtrl([]byte("ActionReplyFailure"), []byte(token))
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
		Cecolla := rx.GetCecolla(uuid)
		if Cecolla != nil {
			if Cecolla.Cecolla != nil {
				_, err := Cecolla.Cecolla.OnCtrl([]byte("PropertyReplySuccess"), []byte(token))
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
		Cecolla := rx.GetCecolla(uuid)
		if Cecolla != nil {
			if Cecolla.Cecolla != nil {
				_, err := Cecolla.Cecolla.OnCtrl([]byte("PropertyReplyFailure"), []byte(token))
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
		productId := stateStack.ToString(3)
		deviceName := stateStack.ToString(4)
		luaTable := stateStack.ToTable(5)
		identifiers := []string{}
		luaTable.ForEach(func(i, v lua.LValue) {
			identifiers = append(identifiers, v.String())
		})
		Cecolla := rx.GetCecolla(uuid)
		if Cecolla != nil {
			if Cecolla.Cecolla != nil {
				getPropertiesCmd := GetPropertiesCmd{
					ProductId:   productId,
					DeviceName:  deviceName,
					Identifiers: identifiers,
				}
				bytes, errMarshal := json.Marshal(getPropertiesCmd)
				if errMarshal != nil {
					stateStack.Push(lua.LString(errMarshal.Error()))
					return 1
				}
				_, err := Cecolla.Cecolla.OnCtrl([]byte("PropertyReport"), []byte(bytes))
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
 * 获取属性回复
 *
 */
func IthingsGetPropertyReplySuccess(rx typex.Rhilex) func(*lua.LState) int {
	return func(stateStack *lua.LState) int {
		uuid := stateStack.ToString(2)
		token := stateStack.ToString(3)
		productId := stateStack.ToString(4)
		deviceName := stateStack.ToString(4)
		luaTable := stateStack.ToTable(6)
		identifiers := []string{}
		luaTable.ForEach(func(i, v lua.LValue) {
			identifiers = append(identifiers, v.String())
		})
		Cecolla := rx.GetCecolla(uuid)
		if Cecolla != nil {
			if Cecolla.Cecolla != nil {
				getPropertiesCmd := GetPropertiesCmd{
					Token:       token,
					ProductId:   productId,
					DeviceName:  deviceName,
					Identifiers: identifiers,
				}
				bytes, errMarshal := json.Marshal(getPropertiesCmd)
				if errMarshal != nil {
					stateStack.Push(lua.LString(errMarshal.Error()))
					return 1
				}
				_, err := Cecolla.Cecolla.OnCtrl([]byte("GetPropertyReplySuccess"), bytes)
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
 * 获取属性
 *
 */
//
type GetPropertiesCmd struct {
	Token       string   `json:"token"`
	ProductId   string   `json:"productID"`
	DeviceName  string   `json:"deviceName"`
	Identifiers []string `json:"identifiers"`
}

func IthingsGetProperties(rx typex.Rhilex) func(*lua.LState) int {
	return func(stateStack *lua.LState) int {
		uuid := stateStack.ToString(2)
		productId := stateStack.ToString(3)
		deviceName := stateStack.ToString(4)
		luaTable := stateStack.ToTable(5)
		identifiers := []string{}
		luaTable.ForEach(func(i, v lua.LValue) {
			identifiers = append(identifiers, v.String())
		})
		Cecolla := rx.GetCecolla(uuid)
		if Cecolla != nil {
			if Cecolla.Cecolla != nil {
				getPropertiesCmd := GetPropertiesCmd{
					ProductId:   productId,
					DeviceName:  deviceName,
					Identifiers: identifiers,
				}
				bytes, errMarshal := json.Marshal(getPropertiesCmd)
				if errMarshal != nil {
					stateStack.Push(lua.LString(errMarshal.Error()))
					return 1
				}
				_, err := Cecolla.Cecolla.OnCtrl([]byte("GetProperties"), []byte(bytes))
				if err != nil {
					stateStack.Push(lua.LString(err.Error()))
					return 1
				}
			}
		}
		stateStack.Push(lua.LNil)
		stateStack.Push(lua.LString("cecolla not exists:" + uuid))
		return 2
	}
}
