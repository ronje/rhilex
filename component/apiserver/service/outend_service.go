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

// -----------------------------------------------------------------------------------
func GetMOutEnd(id string) (*model.MOutEnd, error) {
	m := new(model.MOutEnd)
	if err := interdb.InterDb().First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}
func GetMOutEndWithUUID(uuid string) (*model.MOutEnd, error) {
	m := new(model.MOutEnd)
	if err := interdb.InterDb().Where("uuid=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func InsertMOutEnd(o *model.MOutEnd) error {
	return interdb.InterDb().Table("m_out_ends").Create(o).Error
}

func DeleteMOutEnd(uuid string) error {
	return interdb.InterDb().Where("uuid=?", uuid).Delete(&model.MOutEnd{}).Error
}

func UpdateMOutEnd(uuid string, o *model.MOutEnd) error {
	m := model.MOutEnd{}
	if err := interdb.InterDb().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return err
	} else {
		interdb.InterDb().Model(m).Updates(*o)
		return nil
	}
}
