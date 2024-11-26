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

package dmodbus

const (
	READ_COIL                        = 1  //  Read Coil
	READ_DISCRETE_INPUT              = 2  //  Read Discrete Input
	READ_HOLDING_REGISTERS           = 3  //  Read Holding Registers
	READ_INPUT_REGISTERS             = 4  //  Read Input Registers
	WRITE_SINGLE_COIL                = 5  //  Write Single Coil
	WRITE_SINGLE_HOLDING_REGISTER    = 6  //  Write Single Holding Register
	WRITE_MULTIPLE_COILS             = 15 //  Write Multiple Coils
	WRITE_MULTIPLE_HOLDING_REGISTERS = 16 //  Write Multiple Holding Registers
)

/*
*
* 采集到的数据
*
 */
type ModbusRegister struct {
	UUID      string  `json:"UUID"`
	Tag       string  `json:"tag" validate:"required" title:"数据Tag"`         // 数据Tag
	Alias     string  `json:"alias" validate:"required" title:"别名"`          // 别名
	Function  int     `json:"function" validate:"required" title:"Modbus功能"` // Function
	SlaverId  byte    `json:"slaverId" validate:"required" title:"从机ID"`     // 从机ID
	Address   uint16  `json:"address" validate:"required" title:"地址"`        // Address
	Frequency int64   `json:"frequency" validate:"required" title:"采集频率"`    // 间隔
	Quantity  uint16  `json:"quantity" validate:"required" title:"数量"`       // Quantity
	DataType  string  `json:"dataType"`                                      // 运行时数据
	DataOrder string  `json:"dataOrder"`                                     // 运行时数据
	Weight    float64 `json:"weight"`
	Value     string  `json:"value,omitempty"` // 运行时数据. Type, Order不同值也不同
}

type RegisterList []*ModbusRegister

func (r RegisterList) Len() int {
	return len(r)
}

func (r RegisterList) Less(i, j int) bool {
	if r[i].SlaverId == r[j].SlaverId {
		if r[i].Function == r[j].Function {
			if r[i].Frequency == r[j].Frequency {
				return r[i].Address < r[j].Address
			} else {
				return r[i].Frequency < r[j].Frequency
			}
		} else {
			return r[i].Function < r[j].Function
		}
	} else {
		return r[i].SlaverId < r[j].SlaverId
	}
}

func (r RegisterList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
