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
	"context"
	"errors"
	"time"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/hootrhino/rhilex/component/orderedmap"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

var __DefaultAlarmCenter *AlarmCenter

type AlarmCenter struct {
	e        typex.Rhilex
	registry *orderedmap.OrderedMap[string, *AlarmRule]
	caches   []MAlarmLog
}

func InitAlarmCenter(e typex.Rhilex) {
	__DefaultAlarmCenter = &AlarmCenter{
		e:        e,
		registry: orderedmap.NewOrderedMap[string, *AlarmRule](),
		caches:   []MAlarmLog{},
	}
	InitAlarmDb(e)
	go func() {
		for {
			select {
			case <-context.Background().Done():
				return
			default:
			}
			batchSize := len(__DefaultAlarmCenter.caches)
			if batchSize >= 0 {
				tx := AlarmDb().CreateInBatches(__DefaultAlarmCenter.caches, batchSize)
				if tx.Error != nil {
					glogger.GLogger.Error(tx.Error)
				}
				__DefaultAlarmCenter.caches = []MAlarmLog{}
			}
			time.Sleep(1 * time.Second)
		}
	}()
}

// Stop
func StopAlarmCenter() {
	for _, v := range __DefaultAlarmCenter.registry.Keys() {
		__DefaultAlarmCenter.registry.Delete(v)
	}

}

// Load Expr
func LoadExpr(uuid, exprs string, Threshold uint64, Interval time.Duration, Type string) (*vm.Program, error) {
	Program, err := expr.Compile(exprs, expr.AsBool())
	if err != nil {
		return nil, err
	}
	__DefaultAlarmCenter.registry.Set(uuid, NewAlarmRule(Threshold, Interval, Type, Program))
	return Program, nil
}

// Run Expr
func RunExpr(ruleId, Source string, in map[string]any) (bool, error) {

	AlarmRule, ok := __DefaultAlarmCenter.registry.Get(ruleId)
	if ok {
		output, err := expr.Run(AlarmRule.program, in)
		switch T := output.(type) {
		case bool:
			if T {
				if AlarmRule.AddLog() {
					__DefaultAlarmCenter.caches = append(__DefaultAlarmCenter.caches, MAlarmLog{
						UUID:      utils.AlarmLogUuid(),
						Ts:        uint64(time.Now().UnixMilli()),
						RuleId:    ruleId,
						Source:    Source,
						EventType: AlarmRule.EventType,
						Summary:   "WARNING",
						Info:      AlarmRule.program.Source().String(),
					})
					Target := __DefaultAlarmCenter.e.GetOutEnd(AlarmRule.HandleId)
					if Target != nil {
						if Target.Target != nil {
							// 直接把这个EventType输出到对面
							Target.Target.To(AlarmRule.EventType)
						}
					}
				}
			}
		}
		return false, err
	}
	return false, errors.New("AlarmRule not exists in registry:" + ruleId)
}

// Remove Expr
func RemoveExpr(uuid string) {
	__DefaultAlarmCenter.registry.Delete(uuid)
}

// 输入数据检查规则
func Input(ruleId, Source string, in map[string]any) (bool, error) {
	return RunExpr(ruleId, Source, in)
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
