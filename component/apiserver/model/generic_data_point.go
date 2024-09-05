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

import "gopkg.in/square/go-jose.v2/json"

/**
 *
 * 通用点位
 */

type MDataPoint struct {
	RhilexModel
	UUID       string `gorm:"uniqueIndex"`
	DeviceUuid string `gorm:"not null"`
	Tag        string `gorm:"not null"`
	Alias      string `gorm:"not null"`
	Frequency  int    `gorm:"not null"`
	Config     string `gorm:"not null"`
}

func (mdp MDataPoint) GetConfig() map[string]interface{} {
	result := make(map[string]interface{})
	err := json.Unmarshal([]byte(mdp.Config), &result)
	if err != nil {
		return map[string]interface{}{}
	}
	return result
}

func (md MDevice) GetConfig() map[string]interface{} {
	result := make(map[string]interface{})
	err := json.Unmarshal([]byte(md.Config), &result)
	if err != nil {
		return map[string]interface{}{}
	}
	return result
}
