// Copyright (C) 2025 wwhai
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

// RTSP推拉流设置参数
type MCamera struct {
	RhilexModel
	UUID       string `gorm:"uniqueIndex"` // UUID
	Name       string `gorm:"not null"`    // 名称
	Type       string `gorm:"not null"`    // 设备类型
	StreamUrl  string `gorm:"not null"`    // 拉流地址
	EnablePush *bool  `gorm:"not null"`    // 是否开启推流
	PushUrl    string `gorm:"not null"`    // 推流地址
	EnableAi   *bool  `gorm:"not null"`    // 开启AI？
	AiModel    string `gorm:"not null"`    // AI模型: YOLOV8 | FACENET| .....
}
