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

func GetMUser(username string) (*model.MUser, error) {
	m := new(model.MUser)
	return m, interdb.DB().Where("username=?", username).First(m).Error
}
func Login(username, pwd string) (*model.MUser, error) {
	m := new(model.MUser)
	return m, interdb.DB().
		Where("username=? AND password=?", username, pwd).
		First(m).Error
}

func InsertMUser(o *model.MUser) error {
	return interdb.DB().Model(o).Create(o).Error
}
func InitMUser(o *model.MUser) error {
	return interdb.DB().Model(o).FirstOrCreate(o).Error
}
func ClearAllUser() error {
	return interdb.DB().Model(&model.MUser{}).Exec(`DELETE FROM m_users`).Error
}

func UpdateMUser(oldName string, o *model.MUser) error {
	return interdb.DB().Model(o).
		Where("username=?", oldName).
		Updates(*o).Error
}
