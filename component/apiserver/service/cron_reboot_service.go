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

/**
 * 获取配置
 *
 */

func GetCronRebootConfig() (*model.MCronRebootConfig, error) {
	m := new(model.MCronRebootConfig)
	m.ID = 1
	m.Enable = new(bool)
	m.CronExpr = "0 0 0 0 0"
	return m, interdb.DB().Model(m).First(m).Error
}

/**
 * 更新
 *
 */
func UpdateMCronRebootConfig(o *model.MCronRebootConfig) error {
	return interdb.DB().Model(o).Where("id=?", 1).Updates(o).Error
}
