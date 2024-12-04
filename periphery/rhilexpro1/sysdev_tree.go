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

package rhilexpro1

import "github.com/hootrhino/rhilex/periphery"

func GetSysDevTree() periphery.DeviceTree {
	return periphery.DeviceTree{
		Network: []periphery.DeviceNode{
			{Name: "eth0", Type: "ethernet", Status: 1},
			{Name: "eth1", Type: "ethernet", Status: 1},
		},
		Wlan: []periphery.DeviceNode{
			{Name: "wlan0", Type: "wlan", Status: 1},
		},
		MNet4g: []periphery.DeviceNode{},
		MNet5g: []periphery.DeviceNode{},
		CanBus: []periphery.DeviceNode{},
	}
}
