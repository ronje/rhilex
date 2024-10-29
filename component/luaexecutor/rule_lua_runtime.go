// Copyright (C) 2023 wwhai
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

package luaexecutor

import (
	"errors"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rhilex/component/interpipeline"
	"github.com/hootrhino/rhilex/typex"
)

const (
	SUCCESS_KEY string = "Success"
	FAILED_KEY  string = "Failed"
	ACTIONS_KEY string = "Actions"
)

// LUA Callback : Success
func ExecuteSuccess(vm *lua.LState) (interface{}, error) {
	return interpipeline.Execute(vm, SUCCESS_KEY)
}

// LUA Callback : Failed

func ExecuteFailed(vm *lua.LState, arg lua.LValue) (interface{}, error) {
	return interpipeline.Execute(vm, FAILED_KEY, arg)
}

/*
*
* Execute Lua Callback
*
 */
func ExecuteActions(rule *typex.Rule, arg lua.LValue) (lua.LValue, error) {
	// 原始 lua 数据结构
	luaOriginTable := rule.LuaVM.GetGlobal(ACTIONS_KEY)
	if luaOriginTable != nil && luaOriginTable.Type() == lua.LTTable {
		// 断言成包含回调的 table
		switch funcsTable := luaOriginTable.(type) {
		case *lua.LTable:
			{
				funcs := make(map[string]*lua.LFunction, funcsTable.Len())
				var err error = nil
				funcsTable.ForEach(func(idx, f lua.LValue) {
					if f.Type() == lua.LTFunction {
						funcs[idx.String()] = f.(*lua.LFunction)
					} else {
						err = errors.New(f.String() + " not a lua function")
						return
					}
				})
				if err != nil {
					return nil, err
				}
				// Rule may stop
				if rule.Status != typex.RULE_STOP {
					return interpipeline.RunPipline(rule.LuaVM, funcs, arg)
				}
				return lua.LNil, nil
			}
		default:
			{
				return nil, errors.New("'Actions' is not functions type Table")
			}
		}
	}
	return nil, errors.New("'Actions' not a lua table or not exist")

}
