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

package alarmcenter

import (
	"errors"
	"fmt"
	"time"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/hootrhino/rhilex/component/orderedmap"
	"github.com/hootrhino/rhilex/typex"
)

var __DefaultAlarmCenter *AlarmCenter

type AlarmCenter struct {
	e        typex.Rhilex
	registry *orderedmap.OrderedMap[string, *AlarmRule]
}

func InitAlarmCenter(e typex.Rhilex) {
	__DefaultAlarmCenter = &AlarmCenter{
		e:        e,
		registry: orderedmap.NewOrderedMap[string, *AlarmRule](),
	}

}

// Stop
func StopAlarmCenter() {
	for _, v := range __DefaultAlarmCenter.registry.Keys() {
		__DefaultAlarmCenter.registry.Delete(v)
	}
}

// Load Expr
func LoadExpr(uuid, exprs string, Threshold uint64, Interval time.Duration) (*vm.Program, error) {
	Program, err := expr.Compile(exprs, expr.AsBool())
	if err != nil {
		return nil, err
	}
	__DefaultAlarmCenter.registry.Set(uuid, NewAlarmRule(Threshold, Interval, Program))
	return Program, nil
}

// Run Expr
func RunExpr(uuid string, in map[string]any) (bool, error) {
	AlarmRule, ok := __DefaultAlarmCenter.registry.Get(uuid)
	if ok {
		output, err := expr.Run(AlarmRule.program, in)
		switch T := output.(type) {
		case bool:
			if T {
				// TODO 触发规则
				if AlarmRule.AddLog() {
					fmt.Println("====== AlarmRule.AddLog() ========== ", in)
				}
			}
		}
		return false, err
	}
	return false, errors.New("Invalid Expr vm")
}

// Remove Expr
func RemoveExpr(uuid string) {
	__DefaultAlarmCenter.registry.Delete(uuid)
}

// 输入数据检查规则
func Input(uuid string, in map[string]any) (bool, error) {
	return RunExpr(uuid, in)
}

// 测试
func VerifyExpr(exprs string) (bool, error) {
	_, err := expr.Compile(exprs, expr.AsBool())
	if err != nil {
		return false, err
	}
	return true, nil
}

// 测试
func TestRunExpr(exprs string, in map[string]any) (bool, error) {
	Program, err := expr.Compile(exprs, expr.AsBool())
	if err != nil {
		return false, err
	}
	output, err := expr.Run(Program, in)
	return output.(bool), err
}
