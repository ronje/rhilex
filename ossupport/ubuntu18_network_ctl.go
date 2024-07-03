// Copyright (C) 2023 wwhai
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
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package ossupport

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
)

//--------------------------------------------------------------------------------------
// 注意: 这些设置主要是针对System Ubuntu18.04 的，有可能在不同的发行版有不同的指令，不一定通用
// ！！！！ Warning: MUST RUN WITH SUDO or ROOT USER  ！！！！
//--------------------------------------------------------------------------------------

/*
*
* Ubuntu: 刷新DNS，
*
 */
func ReloadDNS() error {
	cmd := exec.Command("systemctl", "restart", "NetworkManager")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error executing nmcli: %s", err.Error()+":"+string(output))
	}
	return nil
}

/*
*
rer@revb-h3:~$ nmcli device status

	DEVICE           TYPE      STATE         CONNECTION
	usb0             ethernet  connected     Wired connection 1
	wlx0cc6551c5026  wifi      connected     AABBCC
	eth1             ethernet  connected     eth1
	eth0             ethernet  disconnected  --
	lo               loopback  unmanaged     --

*
*/
type DeviceStatus struct {
	DEVICE     string `json:"device"`
	TYPE       string `json:"type"`
	STATE      string `json:"state"`
	CONNECTION string `json:"connection"`
}

func GetCurrentNetConnection() ([]DeviceStatus, error) {
	nmcliCmd := exec.Command("nmcli", "device", "status")
	output, err := nmcliCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Error executing nmcli: %s", err.Error()+":"+string(output))
	}
	nmcliOutputStr := string(output)
	deviceStatuses := parseNmcliOutput(nmcliOutputStr)
	return deviceStatuses, nil
}
func parseNmcliOutput(output string) []DeviceStatus {
	var deviceStatuses []DeviceStatus

	// 按行分割输出
	lines := strings.Split(output, "\n")

	// 如果没有输出行，返回空切片
	if len(lines) == 0 {
		return deviceStatuses
	}

	// 获取列名
	headers := strings.Fields(lines[0])

	// 遍历剩余的行，每行是一个设备状态
	for _, line := range lines[1:] {
		fields := strings.Fields(line)

		// 如果列数不匹配列名数，跳过该行
		if len(fields) != len(headers) {
			continue
		}

		// 创建一个新的设备状态结构体，并填充数据
		var status DeviceStatus
		for i, header := range headers {
			switch header {
			case "DEVICE":
				status.DEVICE = fields[i]
			case "TYPE":
				status.TYPE = fields[i]
			case "STATE":
				status.STATE = fields[i]
			case "CONNECTION":
				status.CONNECTION = fields[i]
			}
		}

		// 将设备状态添加到切片
		deviceStatuses = append(deviceStatuses, status)
	}

	return deviceStatuses
}

type NetInterfaceInfo struct {
	Name string `json:"name"`
	Mac  string `json:"mac"`
	Addr string `json:"addr"`
}

/*
*
* 获取网卡
*
 */
func GetAvailableInterfaces() ([]NetInterfaceInfo, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	netInterfaces := make([]NetInterfaceInfo, 0, len(interfaces))
	for _, inter := range interfaces {
		info := NetInterfaceInfo{
			Name: inter.Name,
			Mac:  inter.HardwareAddr.String(),
		}
		addrs, err := inter.Addrs()
		if err != nil {
			continue
		}
		for i := range addrs {
			addr := addrs[i].String()
			cidr, _, _ := net.ParseCIDR(addr)
			if cidr == nil {
				continue
			}
			if cidr.To4() != nil {
				info.Addr = addr
				break
			}
		}
		netInterfaces = append(netInterfaces, info)
	}
	return netInterfaces, nil
}
