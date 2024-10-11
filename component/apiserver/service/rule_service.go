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
func AllMRules() []model.MRule {
	rules := []model.MRule{}
	interdb.DB().Table("m_rules").Find(&rules)
	return rules
}

func AllMInEnd() []model.MInEnd {
	inends := []model.MInEnd{}
	interdb.DB().Table("m_in_ends").Find(&inends)
	return inends
}

func AllMOutEnd() []model.MOutEnd {
	outends := []model.MOutEnd{}
	interdb.DB().Table("m_out_ends").Find(&outends)
	return outends
}

func AllMUser() []model.MUser {
	users := []model.MUser{}
	interdb.DB().Find(&users)
	return users
}
