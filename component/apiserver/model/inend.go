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

package model

import "encoding/json"

type MInEnd struct {
	RhilexModel
	// UUID for origin source ID
	UUID        string     `gorm:"uniqueIndex"`
	Type        string     `gorm:"not null"`
	Name        string     `gorm:"not null"`
	BindRules   StringList `json:"bindRules"` // 与之关联的规则表["A","B","C"]
	Description string
	Config      string
	XDataModels string
}

func (md MInEnd) GetConfig() map[string]interface{} {
	result := make(map[string]interface{})
	err := json.Unmarshal([]byte(md.Config), &result)
	if err != nil {
		return map[string]interface{}{}
	}
	return result
}
