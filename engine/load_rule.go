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

package engine

import (
	"github.com/hootrhino/rhilex/component/luaruntime"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
)

// LoadRule: 每个规则都绑定了资源(FromSource)或者设备(FromDevice)
// 使用MAP来记录RULE的绑定关系, KEY是UUID, Value是规则
func (e *RuleEngine) LoadRule(r *typex.Rule) error {
	// 前置语法验证
	if err := luaruntime.VerifyLuaSyntax(r); err != nil {
		return err
	}
	// 前置自定义库校验
	if err := luaruntime.LoadExtLuaLib(e, r.LuaVM); err != nil {
		return err
	}
	e.SaveRule(r)
	//--------------------------------------------------------------
	// Load LoadBuildInLuaLib
	//--------------------------------------------------------------
	luaruntime.LoadRuleLibGroup(e, "RULE", r.UUID, r.LuaVM)
	glogger.GLogger.Infof("Rule [%s, %s] load successfully", r.UUID, r.Name)
	// 查找输入定义的资源是否存在
	if in := e.GetInEnd(r.FromSource); in != nil {
		(in.BindRules)[r.UUID] = *r
		return nil
	}
	// 查找输入定义的资源是否存在
	if Device := e.GetDevice(r.FromDevice); Device != nil {
		// 绑定资源和规则，建立关联关系
		(Device.BindRules)[r.UUID] = *r
	}
	return nil

}

// GetRule a rule
func (e *RuleEngine) GetRule(id string) *typex.Rule {
	v, ok := (e.Rules).Get(id)
	if ok {
		return v
	}
	return nil

}

func (e *RuleEngine) SaveRule(r *typex.Rule) {
	e.Rules.Set(r.UUID, r)
}

// RemoveRule and inend--rule bindings
func (e *RuleEngine) RemoveRule(ruleId string) {
	if rule := e.GetRule(ruleId); rule != nil {
		for _, inEnd := range e.AllInEnds() {
			for _, r := range inEnd.BindRules {
				if rule.UUID == r.UUID {
					delete(inEnd.BindRules, ruleId)
				}
			}
		}
		for _, device := range e.AllDevices() {
			for _, r := range device.BindRules {
				glogger.GLogger.Debugf("Unlink rule:[%s, %s]", rule.UUID, rule.Name)
				if rule.UUID == r.UUID {
					delete(device.BindRules, ruleId)
				}
			}
		}
		e.Rules.Delete(ruleId)
		glogger.GLogger.Infof("Rule [%s, %s] has been deleted", ruleId, rule.Name)
	}
}

func (e *RuleEngine) AllRules() []*typex.Rule {
	return e.Rules.Values()
}
