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
* 系统配置参数, 直接以String保存，完了以后再加工成Dto结构体
*
 */
type MNetworkConfig struct {
	RhilexModel
	Type        string     `gorm:"not null"` // 类型: ubuntu16, ubuntu18
	Interface   string     `gorm:"not null"` // eth1 eth0
	Address     string     `gorm:"not null"`
	Netmask     string     `gorm:"not null"`
	Gateway     string     `gorm:"not null"`
	DNS         StringList `gorm:"not null"`
	DHCPEnabled *bool      `gorm:"not null"`
}

/*
*
* 无线网络配置
*
 */
type MWifiConfig struct {
	RhilexModel
	Interface string `gorm:"not null"`
	SSID      string `gorm:"not null"`
	Password  string `gorm:"not null"`
	Security  string `gorm:"not null"` // wpa2-psk wpa3-psk
}
