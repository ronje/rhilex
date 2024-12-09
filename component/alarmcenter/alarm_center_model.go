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

import "time"

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
