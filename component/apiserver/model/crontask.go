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

import "time"

/**
 * 定时任务
 */
type MCronTask struct {
	RhilexModel
	UUID      string    `gorm:"uniqueIndex,not null; default:''" json:"uuid"`
	Name      string    `gorm:"not null;" json:"name"`
	CronExpr  string    `gorm:"not null" json:"cronExpr"` // quartz cron expr
	Enable    *bool     `json:"enable"`                   // 是否启用定时任务
	TaskType  string    `json:"taskType"`                 // CRON_TASK_TYPE，目前只有LINUX_SHELL
	Command   string    `json:"command"`                  // 根据TaskType而定，TaskType=LINUX_SHELL时Command=/bin/bash
	Args      *string   `json:"args"`                     // "-param1 -param2 -param3"
	IsRoot    *bool     `json:"isRoot"`                   // 是否使用root用户运行，目前不使用，默认和rhilex用户一致
	WorkDir   string    `json:"workDir"`                  // 目前不使用，默认工作路径和网关工作路径保持一致
	Env       string    `json:"env"`                      // ["A=e1", "B=e2", "C=e3"]
	Script    string    `json:"script"`                   // 脚本内容，base64编码
	UpdatedAt time.Time `json:"updatedAt"`
}

/**
 * 任务结果
 */
type MCronResult struct {
	RhilexModel
	TaskUuid  string    `gorm:"not null; default:''" json:"taskUuid,omitempty"`
	Status    string    `json:"status"`             // CRON_RESULT_STATUS
	ExitCode  string    `json:"exitCode,omitempty"` // 0-success other-failed
	LogPath   string    `json:"logPath,omitempty"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}
