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

package rhilexg1

import "github.com/hootrhino/rhilex/periphery"

func GetSysDevTree() periphery.DeviceTree {
	wlanList, _ := getWlanList()
	Wlans := []periphery.DeviceNode{}
	for _, wlan := range wlanList {
		Wlans = append(Wlans, periphery.DeviceNode{Name: wlan.Name, Type: periphery.WLAN, Status: 1})
	}
	ethList, _ := getEthList()
	eths := []periphery.DeviceNode{}
	for _, eth := range ethList {
		eths = append(eths, periphery.DeviceNode{Name: eth.Name, Type: periphery.ETHNET, Status: 1})
	}
	return periphery.DeviceTree{
		Network: eths,
		Wlan:    Wlans,
		MNet4g: []periphery.DeviceNode{
			{Name: "usb0", Type: periphery.NM4G, Status: 1},
		},
		MNet5g: []periphery.DeviceNode{},
		CanBus: []periphery.DeviceNode{},
	}
}
