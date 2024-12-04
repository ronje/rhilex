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
func GetMInEnd(uuid string) (*model.MInEnd, error) {
	m := new(model.MInEnd)
	if err := interdb.DB().Table("m_in_ends").Where("uuid=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}
func GetMInEndWithUUID(uuid string) (*model.MInEnd, error) {
	m := new(model.MInEnd)
	if err := interdb.DB().Table("m_in_ends").Where("uuid=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func InsertMInEnd(i *model.MInEnd) error {
	return interdb.DB().Table("m_in_ends").Create(i).Error
}

func DeleteMInEnd(uuid string) error {
	return interdb.DB().Where("uuid=?", uuid).Delete(&model.MInEnd{}).Error
}

func UpdateMInEnd(uuid string, i *model.MInEnd) error {
	m := model.MInEnd{}
	if err := interdb.DB().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return err
	} else {
		interdb.DB().Model(m).Updates(*i)
		return nil
	}
}
