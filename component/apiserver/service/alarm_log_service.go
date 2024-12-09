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

package service

import (
	"github.com/hootrhino/rhilex/component/alarmcenter"
	"github.com/hootrhino/rhilex/component/apiserver/model"
)

// -------------------------------------------------------------------------------------
// AlarmLog Dao
// -------------------------------------------------------------------------------------

func GetMAlarmLogWithUUID(uuid string) (*model.MAlarmLog, error) {
	m := model.MAlarmLog{}
	return &m, alarmcenter.AlarmDb().Where("uuid=?", uuid).First(&m).Error
}

// 删除AlarmLog
func DeleteAlarmLog(uuid string) error {
	return alarmcenter.AlarmDb().Where("uuid=?", uuid).Delete(&model.MAlarmLog{}).Error
}

// 创建AlarmLog
func InsertAlarmLog(AlarmLog *model.MAlarmLog) error {
	return alarmcenter.AlarmDb().Create(AlarmLog).Error
}

// 更新AlarmLog
func UpdateAlarmLog(AlarmLog *model.MAlarmLog) error {
	return alarmcenter.AlarmDb().Model(&model.MAlarmLog{}).
		Where("uuid=?", AlarmLog.UUID).Updates(*AlarmLog).Error
}

// 分页
func PageAlarmLog(current, size int) (int64, []model.MAlarmLog) {
	sql := `SELECT * FROM m_alarm_logs ORDER BY created_at DESC limit ? offset ?;`
	MAlarmLogs := []model.MAlarmLog{}
	tx := alarmcenter.AlarmDb()
	offset := (current - 1) * size
	tx.Raw(sql, size, offset).Find(&MAlarmLogs)
	var count int64
	tx.Model(&model.MAlarmLog{}).Count(&count)
	return count, MAlarmLogs
}

// 分页
func PageAlarmLogByRuleId(ruleId string, current, size int) (int64, []model.MAlarmLog) {
	sql := `SELECT * FROM m_alarm_logs where rule_id=? ORDER BY created_at DESC limit ? offset ?;`
	MAlarmLogs := []model.MAlarmLog{}
	tx := alarmcenter.AlarmDb()
	offset := (current - 1) * size
	tx.Raw(sql, ruleId, size, offset).Find(&MAlarmLogs)
	var count int64
	tx.Model(&model.MAlarmLog{}).Count(&count)
	return count, MAlarmLogs
}
