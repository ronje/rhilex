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
	"encoding/json"
	"time"

	"github.com/expr-lang/expr/vm"
)

type RhilexModel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

// 告警日志
type MAlarmLog struct {
	RhilexModel
	UUID      string `gorm:"not null"` // UUID
	RuleId    string `gorm:"not null"` // 规则ID
	Source    string `gorm:"not null"` // 告警源，某个设备
	EventType string `gorm:"not null"` // 告警标识符
	Ts        uint64 `gorm:"not null"` // 时间戳
	Summary   string `gorm:"not null"` // 概览
	Info      string `gorm:"not null"` // 内容
}

// 告警规则：
// 举例：某个传感器达到30度的时候告警，但是这个冷却时间会有5分钟，如果不做限制会持续不断的生成大量日志
// 在Interval时间内，发生了Threshold条告警才真实告警。
type MAlarmRule struct {
	RhilexModel
	UUID        string `gorm:"uniqueIndex"` // UUID
	Name        string `gorm:"not null"`    // 名称
	ExprDefine  string `gorm:"not null"`
	Interval    uint64 `gorm:"not null"` // 执行周期
	Threshold   uint64 `gorm:"not null"` // 单次触发的日志数量阈值
	HandleId    string // 事件处理器，目前是北向ID
	Description string // 描述
}

func (m *MAlarmRule) GetExprDefine() []ExprDefine {
	outputs := []ExprDefine{}
	json.Unmarshal([]byte(m.ExprDefine), &outputs)
	return outputs
}

func (m *MAlarmRule) SetExprDefine(outputs []ExprDefine) {
	b, _ := json.Marshal(outputs)
	m.ExprDefine = string(b)
}

// 输出规则
type ExprDefine struct {
	Expr      string `json:"expr"`
	EventType string `json:"eventType"`
	program   *vm.Program
}
