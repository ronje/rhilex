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
	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/interdb"
)

// -------------------------------------------------------------------------------------
// AlarmRule Dao
// -------------------------------------------------------------------------------------

func GetMAlarmRuleWithUUID(uuid string) (*model.MAlarmRule, error) {
	m := model.MAlarmRule{}
	return &m, interdb.InterDb().Where("uuid=?", uuid).First(&m).Error
}

// 删除AlarmRule
func DeleteAlarmRule(uuid string) error {
	return interdb.InterDb().Where("uuid=?", uuid).Delete(&model.MAlarmRule{}).Error
}

// 创建AlarmRule
func InsertAlarmRule(AlarmRule *model.MAlarmRule) error {
	return interdb.InterDb().Create(AlarmRule).Error
}

// 更新AlarmRule
func UpdateAlarmRule(AlarmRule *model.MAlarmRule) error {
	return interdb.InterDb().Model(&model.MAlarmRule{}).
		Where("uuid=?", AlarmRule.UUID).Updates(*AlarmRule).Error
}

// 分页
func PageAlarmRule(current, size int) (int64, []model.MAlarmRule) {
	sql := `SELECT * FROM m_alarm_rules ORDER BY created_at DESC limit ? offset ?;`
	MAlarmRules := []model.MAlarmRule{}
	offset := (current - 1) * size
	interdb.InterDb().Raw(sql, size, offset).Find(&MAlarmRules)
	var count int64
	interdb.InterDb().Model(&model.MAlarmRule{}).Count(&count)
	return count, MAlarmRules
}
