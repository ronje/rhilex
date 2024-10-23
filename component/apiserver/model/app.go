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

/*
*
* LUA应用
*
 */
type MApp struct {
	RhilexModel
	UUID        string `gorm:"uniqueIndex"` // 名称
	Name        string `gorm:"not null"`    // 名称
	Version     string `gorm:"not null"`    // 版本号
	AutoStart   *bool  `gorm:"not null"`    // 允许启动
	LuaSource   string `gorm:"not null"`    // LuaSource
	Description string `gorm:"not null"`    // 文件路径, 是相对于main的apps目录
}
