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

package internotify

import (
	"time"

	"github.com/hootrhino/rhilex/typex"
	"gorm.io/gorm"
)

/*
*
* 监控chan
*
 */
type MInternalNotify struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UUID      string    `gorm:"not null"` // UUID
	Type      string    `gorm:"not null"` // INFO | ERROR | WARNING
	Status    int       `gorm:"not null"` // 1 未读 2 已读
	Event     string    `gorm:"not null"` // 字符串
	Ts        uint64    `gorm:"not null"` // 时间戳
	Summary   string    `gorm:"not null"` // 概览，为了节省流量，在消息列表只显示这个字段，Info值为“”
	Info      string    `gorm:"not null"` // 消息内容，是个文本，详情显示
}

/*
*
* Sqlite 数据持久层
*
 */
type SqliteDAO struct {
	engine typex.Rhilex
	name   string   // 框架可以根据名称来选择不同的数据库驱动,为以后扩展准备
	db     *gorm.DB // Sqlite 驱动
}
