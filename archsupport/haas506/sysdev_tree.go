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

import (
	"log"
	"os"

	"github.com/hootrhino/rhilex/archsupport"
	"github.com/hootrhino/rhilex/ossupport"
)

func init() {
	env := os.Getenv("ARCHSUPPORT")
	if env == "HAAS506LD1" {
		InitDevTree()
	}
}
func GetSysDevTree() archsupport.DeviceTree {
	return archsupport.DeviceTree{
		Network: []archsupport.DeviceNode{
			{Name: "eth1", Type: archsupport.ETHNET, Status: 1},
			{Name: "eth2", Type: archsupport.ETHNET, Status: 1},
		},
		Wlan: []archsupport.DeviceNode{
			{Name: "wlan0", Type: archsupport.WLAN, Status: 1},
		},
		MNet4g: []archsupport.DeviceNode{
			{Name: "eth0", Type: archsupport.NM4G, Status: 1},
		},
		MNet5g: []archsupport.DeviceNode{},
		CanBus: []archsupport.DeviceNode{
			{Name: "can1", Type: archsupport.CAN, Status: 1},
			{Name: "can2", Type: archsupport.CAN, Status: 1},
		},
	}
}

// 初始化一些硬件配置
func InitDevTree() {
	{
		log.Println("Ip Link Set Iface eth0 up")
		err := ossupport.IpLinkSetIface("eth0", "up")
		if err != nil {
			log.Println("Ip Link Set Iface eth0 up error:", err)
		}
	}
	{
		log.Println("Ip Link Set Iface eth1 up")
		err := ossupport.IpLinkSetIface("eth1", "up")
		if err != nil {
			log.Println("Ip Link Set Iface eth1 up error:", err)
		}
	}
	{
		log.Println("Ip Link Set Iface eth2 up")
		err := ossupport.IpLinkSetIface("eth2", "up")
		if err != nil {
			log.Println("Ip Link Set Iface eth2 up error:", err)
		}
	}
	{
		log.Println("Ip Link Set Iface wlan0 up")
		err := ossupport.IpLinkSetIface("wlan0", "up")
		if err != nil {
			log.Println("Ip Link Set Iface wlan0 up error:", err)
		}
	}
	{
		log.Println("Ip Link Set Iface can1 up")
		err := ossupport.IpLinkSetIface("can1", "up")
		if err != nil {
			log.Println("Ip Link Set Iface can1 up error:", err)
		}
	}
	{
		log.Println("Ip Link Set Iface can2 up")
		err := ossupport.IpLinkSetIface("can2", "up")
		if err != nil {
			log.Println("Ip Link Set Iface can2 up error:", err)
		}
	}
}
