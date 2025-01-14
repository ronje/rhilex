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

// 获取Camera列表
func AllCamera() []model.MCamera {
	m := []model.MCamera{}
	interdb.InterDb().Find(&m)
	return m

}
func GetCameraWithUUID(uuid string) (*model.MCamera, error) {
	m := model.MCamera{}
	if err := interdb.InterDb().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return nil, err
	} else {
		return &m, nil
	}
}

// 删除Camera
func DeleteCamera(uuid string) error {
	result := interdb.InterDb().Where("uuid = ?", uuid).Delete(&model.MCamera{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete camera with uuid %s: %v", uuid, result.Error)
	}
	return nil
}

// 创建Camera
func InsertCamera(Camera *model.MCamera) error {
	result := interdb.InterDb().Create(Camera)
	if result.Error != nil {
		return fmt.Errorf("failed to insert camera: %v", result.Error)
	}
	return nil
}

// 更新Camera
func UpdateCamera(Camera *model.MCamera) error {
	result := interdb.InterDb().Model(model.MCamera{}).Where("uuid = ?", Camera.UUID).Updates(*Camera)
	if result.Error != nil {
		return fmt.Errorf("failed to update camera with uuid %s: %v", Camera.UUID, result.Error)
	}
	return nil
}

// 获取Camera列表
func PageCamera(current, size int) (int64, []model.MCamera, error) {
	var count int64
	var MCameras []model.MCamera

	// 查询总数
	if err := interdb.InterDb().Model(&model.MCamera{}).Count(&count).Error; err != nil {
		return 0, nil, fmt.Errorf("failed to count cameras: %v", err)
	}

	// 查询分页数据
	offset := (current - 1) * size
	if err := interdb.InterDb().
		Order("created_at DESC").
		Limit(size).
		Offset(offset).
		Find(&MCameras).Error; err != nil {
		return 0, nil, fmt.Errorf("failed to get cameras: %v", err)
	}

	return count, MCameras, nil
}
