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
func GetMRule(uuid string) (*model.MRule, error) {
	m := new(model.MRule)
	return m, interdb.InterDb().Where("uuid=?", uuid).First(m).Error
}
func GetAllMRule() ([]model.MRule, error) {
	m := []model.MRule{}
	return m, interdb.InterDb().Find(&m).Error
}

func GetMRuleWithUUID(uuid string) (*model.MRule, error) {
	m := new(model.MRule)
	return m, interdb.InterDb().Where("uuid=?", uuid).First(m).Error
}

func InsertMRule(r *model.MRule) error {
	return interdb.InterDb().Table("m_rules").Create(r).Error
}

func DeleteMRule(uuid string) error {
	return interdb.InterDb().Table("m_rules").Where("uuid=?", uuid).Delete(&model.MRule{}).Error
}

func UpdateMRule(uuid string, r *model.MRule) error {
	return interdb.InterDb().Model(r).Where("uuid=?", uuid).Updates(*r).Error
}

// -----------------------------------------------------------------------------------
func AllMRules() []model.MRule {
	rules := []model.MRule{}
	interdb.InterDb().Table("m_rules").Find(&rules)
	return rules
}

func AllMInEnd() []model.MInEnd {
	inends := []model.MInEnd{}
	interdb.InterDb().Table("m_in_ends").Find(&inends)
	return inends
}

func AllMOutEnd() []model.MOutEnd {
	outends := []model.MOutEnd{}
	interdb.InterDb().Table("m_out_ends").Find(&outends)
	return outends
}

func AllMUser() []model.MUser {
	users := []model.MUser{}
	interdb.InterDb().Find(&users)
	return users
}
