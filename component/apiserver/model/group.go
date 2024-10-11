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
* 通用分组
*
 */
type MGenericGroup struct {
	RhilexModel
	UUID   string `gorm:"uniqueIndex"`
	Name   string `gorm:"not null"` // 名称
	Type   string `gorm:"not null"` // 组的类型, DEVICE: 设备分组
	Parent string `gorm:"not null"` // 上级, 如果是0表示根节点
}

/*
*
* 分组表的绑定关系表
*
 */
type MGenericGroupRelation struct {
	RhilexModel
	UUID string `gorm:"uniqueIndex"`
	Gid  string `gorm:"not null"` // 分组ID
	Rid  string `gorm:"not null"` // 被绑定方
}
