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
	"encoding/json"
	"errors"
	"fmt"

	"sync"
	"time"

	"github.com/expr-lang/expr"
	"github.com/hootrhino/rhilex/component/orderedmap"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

var __DefaultAlarmCenter *AlarmCenter

type QueueData struct {
	RuleId    string
	Source    string
	Expr      string
	EventType string
	HandleId  string
	Data      string
	In        map[string]any
}
type AlarmCenter struct {
	e         typex.Rhilex
	registry  *orderedmap.OrderedMap[string, *AlarmRule]
	caches    []MAlarmLog
	QueueData chan QueueData
}

func InitAlarmCenter(e typex.Rhilex) {
	__DefaultAlarmCenter = &AlarmCenter{
		e:         e,
		registry:  orderedmap.NewOrderedMap[string, *AlarmRule](),
		caches:    []MAlarmLog{},
		QueueData: make(chan QueueData, 1024),
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
	go func() {
		for {
			select {
			case <-context.Background().Done():
				return
			case qData := <-__DefaultAlarmCenter.QueueData:
				__DefaultAlarmCenter.caches = append(__DefaultAlarmCenter.caches, MAlarmLog{
					UUID:      utils.AlarmLogUuid(),
					Ts:        uint64(time.Now().UnixMilli()),
					RuleId:    qData.RuleId,
					Source:    qData.Source,
					EventType: qData.EventType,
					Summary:   fmt.Sprintf("EventType:[ %s ] | Source:[ %s ]", qData.EventType, qData.Source),
					Info:      qData.Data,
				})
				Target := __DefaultAlarmCenter.e.GetOutEnd(qData.HandleId)
				if Target != nil {
					if Target.Target != nil {
						// 直接把这个EventType输出到对面
						Target.Target.To(qData.EventType)
					}
				}
			}
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
func LoadAlarmRule(uuid string, alarmRule AlarmRule) error {
	ExprDefines := []ExprDefine{}
	for _, exprDefine := range alarmRule.ExprDefines {
		program, err := expr.Compile(exprDefine.Expr, expr.AsBool())
		if err != nil {
			return err
		}
		ExprDefines = append(ExprDefines, ExprDefine{
			Expr:      exprDefine.Expr,
			EventType: exprDefine.EventType,
			program:   program,
		})
	}
	NewAlarm := NewAlarmRule(alarmRule.Threshold, alarmRule.Interval, alarmRule.HandleId, ExprDefines)
	__DefaultAlarmCenter.registry.Set(uuid, NewAlarm)
	return nil
}
func ReLoadAlarmRule(uuid string, alarmRule AlarmRule) error {
	RemoveExpr(uuid)
	return LoadAlarmRule(uuid, alarmRule)
}

// Run Expr
func RunExpr(ruleId, Source string, in map[string]any) (bool, error) {
	AlarmRule, ok := __DefaultAlarmCenter.registry.Get(ruleId)
	if ok {
		group := sync.WaitGroup{}
		group.Add(len(AlarmRule.ExprDefines))
		for _, ExprDefine := range AlarmRule.ExprDefines {
			go func() {
				defer group.Done()
				output, err := expr.Run(ExprDefine.program, in)
				if err != nil {
					return
				}
				switch T := output.(type) {
				case bool:
					if T {
						if AlarmRule.AddLog() {
							__DefaultAlarmCenter.QueueData <- QueueData{
								In:        in,
								RuleId:    ruleId,
								Source:    Source,
								Expr:      ExprDefine.Expr,
								EventType: ExprDefine.EventType,
								HandleId:  AlarmRule.HandleId,
								Data:      BeautifulMapPrint(in),
							}
						}
					}
				}
			}()
		}
		group.Wait()
		return true, nil
	}
	return false, errors.New("AlarmRule not exists in registry:" + ruleId)
}

// Remove Expr
func RemoveExpr(uuid string) {
	__DefaultAlarmCenter.registry.Delete(uuid)
}

// 输入数据检查规则
func Input(ruleId, Source string, in map[string]any) (bool, error) {
	glogger.GLogger.Debugf("AlarmCenter Input RuleId:%s Source:%s Data:%s", ruleId, Source, BeautifulMapPrint(in))
	return RunExpr(ruleId, Source, in)
}

// BeautifulMapPrint 打印一个map，具有漂亮的格式，并返回格式化的字符串。
func BeautifulMapPrint(m map[string]interface{}) string {
	b, _ := json.Marshal(m)
	return string(b)
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
	if err != nil {
		return false, err
	}
	T, ok := output.(bool)
	if !ok {
		return false, fmt.Errorf("%s", "invalid expr result")
	}
	return T, err
}
