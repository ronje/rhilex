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

package haas506

import "github.com/hootrhino/rhilex/archsupport"

func GetSysDevTree() archsupport.DeviceTree {
	return archsupport.DeviceTree{
		Network: []archsupport.DeviceNode{
			{Name: "eth0", Type: archsupport.ETHNET, Status: 1},
			{Name: "eth1", Type: archsupport.ETHNET, Status: 1},
		},
		Wlan: []archsupport.DeviceNode{
			{Name: "wlan0", Type: archsupport.WLAN, Status: 1},
		},
		MNet4g: []archsupport.DeviceNode{
			{Name: "usb2", Type: archsupport.NM4G, Status: 1},
		},
		MNet5g: []archsupport.DeviceNode{
			{Name: "usb1", Type: archsupport.NM5G, Status: 1},
		},
		CanBus: []archsupport.DeviceNode{
			{Name: "can1", Type: archsupport.CAN, Status: 1},
			{Name: "can2", Type: archsupport.CAN, Status: 1},
		},
	}
}
