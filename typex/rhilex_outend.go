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

import (
	"github.com/hootrhino/rhilex/utils"
)

type OutEnd struct {
	UUID        string      `json:"uuid"`
	State       SourceState `json:"state"`
	Type        TargetType  `json:"type"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	//
	Config map[string]interface{} `json:"config"`
	Target XTarget                `json:"-"`
}

func NewOutEnd(t TargetType,
	n string,
	d string,
	c map[string]interface{}) *OutEnd {
	return &OutEnd{
		UUID:        utils.MakeUUID("OUTEND"),
		Type:        TargetType(t),
		State:       SOURCE_DOWN,
		Name:        n,
		Description: d,
		Config:      c,
	}
}

func (out *OutEnd) GetConfig(k string) interface{} {
	return (out.Config)[k]
}
