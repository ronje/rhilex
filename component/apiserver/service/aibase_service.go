// Copyright (C) 2025 wwhai
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

// 获取AiBase列表
func AllAiBase() []model.MAiBase {
	m := []model.MAiBase{}
	interdb.InterDb().Find(&m)
	return m

}
func GetAiBaseWithUUID(uuid string) (*model.MAiBase, error) {
	m := model.MAiBase{}
	if err := interdb.InterDb().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return nil, err
	} else {
		return &m, nil
	}
}

// 删除AiBase
func DeleteAiBase(uuid string) error {
	return interdb.InterDb().Where("uuid=?", uuid).Delete(&model.MAiBase{}).Error
}

// 创建AiBase
func InsertAiBase(AiBase *model.MAiBase) error {
	return interdb.InterDb().Create(AiBase).Error
}

// 更新AiBase
func UpdateAiBase(AiBase *model.MAiBase) error {
	m := model.MAiBase{}
	if err := interdb.InterDb().Where("uuid=?", AiBase.UUID).First(&m).Error; err != nil {
		return err
	} else {
		interdb.InterDb().Model(m).Updates(*AiBase)
		return nil
	}
}
