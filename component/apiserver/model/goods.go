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

//
// 外挂
//

type MGoods struct {
	RhilexModel
	UUID        string `gorm:"uniqueIndex"`
	LocalPath   string `gorm:"not null"`
	GoodsType   string `gorm:"not null"` // LOCAL, EXTERNAL
	ExecuteType string `gorm:"not null"` // exe,elf,js,py....
	AutoStart   bool  `gorm:"not null"`
	NetAddr     string `gorm:"not null"`
	Args        string `gorm:"not null"`
	Description string `gorm:"not null"`
}
