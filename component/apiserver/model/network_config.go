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
	Type        string // 类型: ETH | WIFI
	Interface   string `gorm:"column:interface;uniqueIndex"`
	Address     string
	Netmask     string
	Gateway     string
	DNS         StringList
	DHCPEnabled bool
	SSID        string `gorm:"column:ssid"`
	Password    string
	Security    string // wpa2-psk wpa3-psk
}
