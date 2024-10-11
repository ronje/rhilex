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

import "github.com/hootrhino/rhilex/archsupport"

func GetSysDevTree() archsupport.DeviceTree {
	wlanList, _ := getWlanList()
	Wlans := []archsupport.DeviceNode{}
	for _, wlan := range wlanList {
		Wlans = append(Wlans, archsupport.DeviceNode{Name: wlan.Name, Type: archsupport.WLAN, Status: 1})
	}
	ethList, _ := getEthList()
	eths := []archsupport.DeviceNode{}
	for _, eth := range ethList {
		eths = append(eths, archsupport.DeviceNode{Name: eth.Name, Type: archsupport.ETHNET, Status: 1})
	}
	return archsupport.DeviceTree{
		Network: eths,
		Wlan:    Wlans,
		MNet4g: []archsupport.DeviceNode{
			{Name: "usb0", Type: archsupport.NM4G, Status: 1},
		},
		MNet5g: []archsupport.DeviceNode{},
		CanBus: []archsupport.DeviceNode{},
	}
}
