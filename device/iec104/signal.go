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

package iec104
//Signal 104信号
type Signal struct {
	TypeID  uint    `json:"type_id"` //类型id，1:单点遥信，9:单点遥测
	Address uint32  `json:"address"` //地址
	Value   float64 `json:"value"`   //值
	Quality byte    `json:"quality"` //品质描述
	Ts      float64 `json:"ts"`      //毫秒时间戳
}
