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

package typex

import "github.com/hootrhino/rhilex/utils"

type Cecolla struct {
	UUID        string                 `json:"uuid"`
	Name        string                 `json:"name"`
	Type        CecollaType            `json:"type"`
	Description string                 `json:"description"`
	State       CecollaState           `json:"state"`
	Config      map[string]interface{} `json:"config"`
	Action      string                 `json:"action"`
	Cecolla     XCecolla               `json:"-"`
}

func NewCecolla(t CecollaType, name string,
	description string, config map[string]interface{}) *Cecolla {
	return &Cecolla{
		UUID:        utils.CecUuid(),
		Name:        name,
		Type:        t,
		State:       CEC_DOWN,
		Description: description,
		Config:      config,
	}
}
