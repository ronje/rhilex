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

package model

// 告警规则：
// 举例：某个传感器达到30度的时候告警，但是这个冷却时间会有5分钟，如果不做限制会持续不断的生成大量日志
// 在Interval时间内，发生了Threshold条告警才真实告警。
type MAlarmRule struct {
	RhilexModel
	UUID        string `gorm:"uniqueIndex"` // UUID
	Name        string `gorm:"not null"`    // 名称
	Expr        string `gorm:"not null"`    // 表达式
	Interval    uint64 `gorm:"not null"`    // 执行周期
	Threshold   uint64 `gorm:"not null"`    // 单次触发的日志数量阈值
	Description string // 描述
}
