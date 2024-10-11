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

type MAiBase struct {
	RhilexModel
	UUID        string `gorm:"uniqueIndex"` // 名称
	Name        string `gorm:"not null"`    // 名称
	Type        string `gorm:"not null"`    // 类型
	IsBuildIn   bool   `gorm:"not null"`    // 是否内建
	Version     string `gorm:"not null"`    // 版本号
	Filepath    string `gorm:"not null"`    // 文件路径, 是相对于main的apps目录
	Description string `gorm:"not null"`
}
