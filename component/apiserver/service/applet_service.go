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
// App Dao
// -------------------------------------------------------------------------------------

// 获取App列表
func AllApp() []model.MApplet {
	m := []model.MApplet{}
	interdb.DB().Find(&m)
	return m

}
func GetMAppWithUUID(uuid string) (*model.MApplet, error) {
	m := model.MApplet{}
	if err := interdb.DB().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return nil, err
	} else {
		return &m, nil
	}
}

// 删除App
func DeleteApp(uuid string) error {
	return interdb.DB().Where("uuid=?", uuid).Delete(&model.MApplet{}).Error
}

// 创建App
func InsertApp(app *model.MApplet) error {
	return interdb.DB().Create(app).Error
}

// 更新App
func UpdateApp(app *model.MApplet) error {
	m := model.MApplet{}
	if err := interdb.DB().Where("uuid=?", app.UUID).First(&m).Error; err != nil {
		return err
	} else {
		interdb.DB().Model(m).Where("uuid=?", app.UUID).Updates(*app)
		return nil
	}
}
