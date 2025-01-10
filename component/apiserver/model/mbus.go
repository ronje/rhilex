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

// 点位
type MMBusDataPoint struct {
	RhilexModel
	UUID         string
	DeviceUuid   string
	SlaverId     string
	Type         string
	Manufacturer string
	Tag          string
	Alias        string
	DataLength   *uint64
	Frequency    *uint64  `gorm:"default:50"`
	Weight       *Decimal `gorm:"column:weight;default:1"`
}

func (MMBusDataPoint) TableName() string {
	return "m_mbus_data_points"
}
