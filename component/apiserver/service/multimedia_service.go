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
	"fmt"

	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/interdb"
)

// 获取MultiMedia列表
func AllMultiMedia() []model.MMultiMedia {
	m := []model.MMultiMedia{}
	interdb.InterDb().Find(&m)
	return m

}
func GetMultiMediaWithUUID(uuid string) (*model.MMultiMedia, error) {
	m := model.MMultiMedia{}
	if err := interdb.InterDb().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return nil, err
	} else {
		return &m, nil
	}
}

// 删除MultiMedia
func DeleteMultiMedia(uuid string) error {
	result := interdb.InterDb().Where("uuid = ?", uuid).Delete(&model.MMultiMedia{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete MultiMedia with uuid %s: %v", uuid, result.Error)
	}
	return nil
}

// 创建MultiMedia
func InsertMultiMedia(MultiMedia *model.MMultiMedia) error {
	result := interdb.InterDb().Create(MultiMedia)
	if result.Error != nil {
		return fmt.Errorf("failed to insert MultiMedia: %v", result.Error)
	}
	return nil
}

// 更新MultiMedia
func UpdateMultiMedia(MultiMedia *model.MMultiMedia) error {
	result := interdb.InterDb().Model(model.MMultiMedia{}).Where("uuid = ?", MultiMedia.UUID).Updates(*MultiMedia)
	if result.Error != nil {
		return fmt.Errorf("failed to update MultiMedia with uuid %s: %v", MultiMedia.UUID, result.Error)
	}
	return nil
}

// 获取MultiMedia列表
func PageMultiMedia(current, size int) (int64, []model.MMultiMedia, error) {
	var count int64
	var MMultiMedias []model.MMultiMedia

	// 查询总数
	if err := interdb.InterDb().Model(&model.MMultiMedia{}).Count(&count).Error; err != nil {
		return 0, nil, fmt.Errorf("failed to count MultiMedias: %v", err)
	}

	// 查询分页数据
	offset := (current - 1) * size
	if err := interdb.InterDb().
		Order("created_at DESC").
		Limit(size).
		Offset(offset).
		Find(&MMultiMedias).Error; err != nil {
		return 0, nil, fmt.Errorf("failed to get MultiMedias: %v", err)
	}

	return count, MMultiMedias, nil
}
