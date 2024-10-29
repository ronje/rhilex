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

package luaexecutor

import (
	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/sirupsen/logrus"
)

/*
*
* 执行针对资源端的规则脚本
*
 */
func RunSourceCallbacks(in *typex.InEnd, callbackArgs string) {
	// 执行来自资源的脚本
	for _, rule := range in.BindRules {
		if rule.Status == typex.RULE_RUNNING {
			_, errA := ExecuteActions(&rule, lua.LString(callbackArgs))
			if errA != nil {
				Debugger, Ok := rule.LuaVM.GetStack(1)
				if Ok {
					LValue, _ := rule.LuaVM.GetInfo("f", Debugger, lua.LNil)
					rule.LuaVM.GetInfo("l", Debugger, lua.LNil)
					rule.LuaVM.GetInfo("S", Debugger, lua.LNil)
					rule.LuaVM.GetInfo("u", Debugger, lua.LNil)
					rule.LuaVM.GetInfo("n", Debugger, lua.LNil)
					LFunction := LValue.(*lua.LFunction)
					LastCall := lua.DbgCall{
						Name: "_main",
					}
					if len(LFunction.Proto.DbgCalls) > 0 {
						LastCall = LFunction.Proto.DbgCalls[0]
					}
					glogger.GLogger.WithFields(logrus.Fields{
						"topic": "rule/log/" + rule.UUID,
					}).Warnf("Function Name: [%s],"+
						"What: [%s], Source Line: [%d],"+
						" Last Call: [%s], Error message: %s",
						Debugger.Name, Debugger.What, Debugger.CurrentLine,
						LastCall.Name, errA.Error(),
					)
				}
			} else {
				_, errS := ExecuteSuccess(rule.LuaVM)
				if errS != nil {
					glogger.GLogger.Error(errS)
					return // lua 是规则链，有短路原则，中途出错会中断
				}
			}
		}
	}
}

/*
*
* 执行针对设备端的规则脚本
*
 */
func RunDeviceCallbacks(Device *typex.Device, callbackArgs string) {
	for _, rule := range Device.BindRules {
		_, errA := ExecuteActions(&rule, lua.LString(callbackArgs))
		if errA != nil {
			Debugger, Ok := rule.LuaVM.GetStack(1)
			if Ok {
				LValue, _ := rule.LuaVM.GetInfo("f", Debugger, lua.LNil)
				rule.LuaVM.GetInfo("l", Debugger, lua.LNil)
				rule.LuaVM.GetInfo("S", Debugger, lua.LNil)
				rule.LuaVM.GetInfo("u", Debugger, lua.LNil)
				rule.LuaVM.GetInfo("n", Debugger, lua.LNil)
				LFunction := LValue.(*lua.LFunction)
				LastCall := lua.DbgCall{
					Name: "_main",
				}
				if len(LFunction.Proto.DbgCalls) > 0 {
					LastCall = LFunction.Proto.DbgCalls[0]
				}
				glogger.GLogger.WithFields(logrus.Fields{
					"topic": "rule/log/" + rule.UUID,
				}).Warnf("Function Name: [%s],"+
					"What: [%s], Source Line: [%d],"+
					" Last Call: [%s], Error message: %s",
					Debugger.Name, Debugger.What, Debugger.CurrentLine,
					LastCall.Name, errA.Error(),
				)
			}
		} else {
			_, err2 := ExecuteSuccess(rule.LuaVM)
			if err2 != nil {
				glogger.GLogger.WithFields(logrus.Fields{
					"topic": "rule/log/" + rule.UUID,
				}).Info("RunLuaCallbacks error:", err2)
				return
			}
		}
	}
}
