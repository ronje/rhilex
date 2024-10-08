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

package apis

import (
	"bufio"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	common "github.com/hootrhino/rhilex/component/apiserver/common"
	"github.com/hootrhino/rhilex/ossupport"
	"github.com/hootrhino/rhilex/typex"
)

/*
*
*nmcli device status

	DEVICE           TYPE      STATE      CONNECTION
	eth0             ethernet  connected  eth0
	usb0             ethernet  connected  usb0
	wlx0cc6551c5026  wifi      connected  iotlab4072
	eth1             ethernet  connected  eth1
	lo               loopback  unmanaged  --
*/
type networkDevice struct {
	// 网卡名称
	Device string `json:"device"`
	// 网卡类型
	// ethernet：以太网
	// wifi：WiFi
	// bridge：桥接设备
	Type string `json:"type"`
	// 网络状态
	// connected：已连接到。
	// disconnected：未连接。
	// unmanaged：系统默认
	// unavailable：网络不可用。
	State string `json:"state"`
	// 网络名称
	Connection string `json:"connection"`
}

func GetNmcliDeviceStatus(c *gin.Context, ruleEngine typex.Rhilex) {

	cmd := exec.Command("nmcli", "device", "status")
	output, err := cmd.Output()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	devices, err := parseNmcliDeviceStatus(string(output))
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(devices))
}

// parseNetworkDevices 解析网络设备信息
func parseNmcliDeviceStatus(input string) ([]networkDevice, error) {
	var devices []networkDevice

	// 将输入按行分割
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "DEVICE") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		device := networkDevice{
			Device:     fields[0],
			Type:       fields[1],
			State:      fields[2],
			Connection: fields[3],
		}
		devices = append(devices, device)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return devices, nil
}

/*
* 网卡详情:
*   nmcli device show eth0
*
 */
func GetNmcliDeviceShow(c *gin.Context, ruleEngine typex.Rhilex) {
	ifaceName, _ := c.GetQuery("iface")
	interfaces, err := ossupport.GetAvailableInterfaces()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	ok := false
	for _, iface := range interfaces {
		if iface.Name == ifaceName {
			ok = true
			break
		}
	}
	if !ok {
		c.JSON(common.HTTP_OK, common.Error("interface not exists"))
		return
	}
	cmd := exec.Command("nmcli", "device", "show", ifaceName)
	nmcliOutput, err := cmd.Output()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	device, err := parseNmcliDeviceShow(string(nmcliOutput))
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(device))
}

// GENERAL.DEVICE:                 eth0
// GENERAL.TYPE:                   ethernet
// GENERAL.HWADDR:                 02:81:5E:DF:D4:81
// GENERAL.MTU:                    1500
// GENERAL.STATE:                  100 (connected)
// GENERAL.CONNECTION:             eth0
// GENERAL.CON-PATH:               /org/freedesktop/NetworkManager/ActiveConnection/2
// WIRED-PROPERTIES.CARRIER:       on
// IP4.ADDRESS[1]:                 192.168.1.185/24
// IP4.GATEWAY:                    192.168.1.1
// IP4.ROUTE[1]:                   dst = 0.0.0.0/0, nh = 192.168.1.1, mt = 101
// IP4.ROUTE[2]:                   dst = 192.168.1.0/24, nh = 0.0.0.0, mt = 101
// IP4.DNS[1]:                     192.168.1.1
// IP6.ADDRESS[1]:                 fe80::9460:7480:61a9:cbd2/64
// IP6.GATEWAY:                    --
// IP6.ROUTE[1]:                   dst = ff00::/8, nh = ::, mt = 256, table=255
// IP6.ROUTE[2]:                   dst = fe80::/64, nh = ::, mt = 256
// IP6.ROUTE[3]:                   dst = fe80::/64, nh = ::, mt = 101

type networkDeviceDetail struct {
	Device      string `json:"device"`
	Type        string `json:"type"`
	HWAddr      string `json:"hwAddr"`
	MTU         int    `json:"mtu"`
	State       string `json:"state"`
	Connection  string `json:"connection"`
	Carrier     string `json:"carrier"`
	IPv4Addr    string `json:"ipv4Addr"`
	IPv4Gateway string `json:"ipv4Gateway"`
	IPv4DNS     string `json:"ipv4Dns"`
	IPv6Addr    string `json:"ipv6Addr"`
	IPv6Gateway string `json:"ipv6Gateway"`
}

// parseNMCLIOutput 解析 nmcli 输出
// nmcli device show
func parseNmcliDeviceShow(output string) (*networkDeviceDetail, error) {
	lines := strings.Split(output, "\n")

	device := &networkDeviceDetail{}

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		switch fields[0] {
		case "GENERAL.DEVICE:":
			device.Device = fields[1]
		case "GENERAL.TYPE:":
			device.Type = fields[1]
		case "GENERAL.HWADDR:":
			device.HWAddr = fields[1]
		case "GENERAL.MTU:":
			device.MTU = parseInt(fields[1])
		case "GENERAL.STATE:":
			device.State = fields[1]
		case "GENERAL.CONNECTION:":
			device.Connection = fields[1]
		case "WIRED-PROPERTIES.CARRIER:":
			device.Carrier = fields[1]
		case "IP4.ADDRESS[1]:":
			device.IPv4Addr = fields[1]
		case "IP4.GATEWAY:":
			device.IPv4Gateway = fields[1]
		case "IP4.DNS[1]:":
			device.IPv4DNS = fields[1]
		case "IP6.ADDRESS[1]:":
			device.IPv6Addr = fields[1]
		case "IP6.GATEWAY:":
			device.IPv6Gateway = fields[1]
		}
	}

	return device, nil
}

// parseInt 将字符串转换为整数，如果失败返回 0
func parseInt(s string) int {
	result, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return result
}
