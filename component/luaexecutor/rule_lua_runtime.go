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
	"fmt"

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
// ExecuteSuccess 执行成功回调，调用 interpipeline.Execute 函数并传递 SUCCESS_KEY
// vm 是 Lua 虚拟机的状态
// 返回 interpipeline.Execute 函数的执行结果和错误信息
func ExecuteSuccess(vm *lua.LState) (interface{}, error) {
	// 调用 interpipeline.Execute 函数并传递 SUCCESS_KEY
	result, err := interpipeline.Execute(vm, SUCCESS_KEY)
	if err != nil {
		// 如果执行过程中出现错误，记录错误信息并返回
		// 这里可以根据实际需求添加更详细的日志记录，如使用日志库
		return nil, fmt.Errorf("Execute Success Error: %w", err)
	}
	return result, nil
}

// LUA Callback : Failed
// ExecuteFailed 执行失败回调，调用 interpipeline.Execute 函数并传递 FAILED_KEY 和额外的参数
// vm 是 Lua 虚拟机的状态
// arg 是传递给 interpipeline.Execute 函数的额外参数
// 返回 interpipeline.Execute 函数的执行结果和错误信息
func ExecuteFailed(vm *lua.LState, arg lua.LValue) (interface{}, error) {
	// 调用 interpipeline.Execute 函数并传递 FAILED_KEY 和额外的参数
	result, err := interpipeline.Execute(vm, FAILED_KEY, arg)
	if err != nil {
		// 如果执行过程中出现错误，记录错误信息并返回
		// 这里可以根据实际需求添加更详细的日志记录，如使用日志库
		return nil, fmt.Errorf("Execute Failed Error: %w", err)
	}
	return result, nil
}

/*
*
* Execute Lua Callback
*
 */
func ExecuteActions(rule *typex.Rule, arg lua.LValue) (lua.LValue, error) {
	// 原始 lua 数据结构
	luaOriginTable := rule.LuaVM.GetGlobal(ACTIONS_KEY)
	// 检查 'Actions' 是否存在且为 Lua 表
	if luaOriginTable == nil || luaOriginTable.Type() != lua.LTTable {
		return nil, errors.New("'Actions' not a lua table or not exist")
	}
	// 断言成包含回调的 table
	funcsTable, ok := luaOriginTable.(*lua.LTable)
	if !ok {
		return nil, errors.New("'Actions' is not functions type Table")
	}

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
