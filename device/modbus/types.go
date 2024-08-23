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

package modbus

/*
*
* 点位表
*
 */
type ModbusPoint struct {
	UUID      string  `json:"uuid,omitempty"` // 当UUID为空时新建
	Tag       string  `json:"tag"`
	Alias     string  `json:"alias"`
	Function  int     `json:"function"`
	SlaverId  byte    `json:"slaverId"`
	Address   uint16  `json:"address"`
	Frequency int64   `json:"frequency"`
	Quantity  uint16  `json:"quantity"`
	Value     string  `json:"value,omitempty"` // 运行时数据
	DataType  string  `json:"dataType"`        // 运行时数据
	DataOrder string  `json:"dataOrder"`       // 运行时数据
	Weight    float64 `json:"weight"`          // 权重
}
