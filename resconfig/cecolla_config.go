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

package resconfig

import "fmt"

/**
 * 云边协同
 *
 */
type CecollaConfig struct {
	Enable             *bool  `json:"enable"`             // 是否开启
	CecollaId          string `json:"cecollaId"`          // Cecolla UUID
	EnableCreateSchema *bool  `json:"enableCreateSchema"` // 是否允许设备创建物模型
}

func (c *CecollaConfig) Validate() error {
	if c.CecollaId == "" {
		return fmt.Errorf("invalid cecollaId")
	}
	return nil
}
