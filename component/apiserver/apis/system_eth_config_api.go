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

package apis

import (
	"net"
	"runtime"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rhilex/archsupport/haas506"
	"github.com/hootrhino/rhilex/archsupport/rhilexg1"
	"github.com/hootrhino/rhilex/archsupport/rhilexpro1"
	common "github.com/hootrhino/rhilex/component/apiserver/common"
	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/apiserver/service"
	"github.com/hootrhino/rhilex/typex"
)

type NetworkConfigVo struct {
	Interface   string   `json:"interface"`
	Address     string   `json:"address"`
	Netmask     string   `json:"netmask"`
	Gateway     string   `json:"gateway"`
	DNS         []string `json:"dns"`
	DHCPEnabled *bool    `json:"dhcp_enabled"`
}

/*
*
* 展示网络配置信息
*
 */
func GetEthNetwork(c *gin.Context, ruleEngine typex.Rhilex) {
	Interface, _ := c.GetQuery("iface")
	MNetworkConfig, err := service.GetEthConfig(Interface)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(NetworkConfigVo{
		Interface:   MNetworkConfig.Interface,
		Address:     MNetworkConfig.Address,
		Netmask:     MNetworkConfig.Netmask,
		Gateway:     MNetworkConfig.Gateway,
		DNS:         []string{"8.8.8.8", "114.114.114.114"},
		DHCPEnabled: MNetworkConfig.DHCPEnabled,
	}))
}

/**
 * 获取网卡配置表
 *
 */
func AllEthNetwork(c *gin.Context, ruleEngine typex.Rhilex) {
	MNetworkConfigs, err := service.AllEthConfig()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	NetworkConfigVos := []NetworkConfigVo{}
	for _, MNetworkConfig := range MNetworkConfigs {
		NetworkConfigVos = append(NetworkConfigVos, NetworkConfigVo{
			Interface:   MNetworkConfig.Interface,
			Address:     MNetworkConfig.Address,
			Netmask:     MNetworkConfig.Netmask,
			Gateway:     MNetworkConfig.Gateway,
			DNS:         []string{"8.8.8.8", "114.114.114.114"},
			DHCPEnabled: MNetworkConfig.DHCPEnabled,
		})
	}
	c.JSON(common.HTTP_OK, common.OkWithData(NetworkConfigVos))
}

/*
 *
 * 设置两个网口
 *
 */
func SetEthNetwork(c *gin.Context, ruleEngine typex.Rhilex) {
	if runtime.GOOS != "linux" {
		c.JSON(common.HTTP_OK, common.Error("OS Not Support:"+runtime.GOOS))
		return
	}

	DtoCfg := NetworkConfigVo{}
	if err0 := c.ShouldBindJSON(&DtoCfg); err0 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err0))
		return
	}

	if !isValidIP(DtoCfg.Address) {
		c.JSON(common.HTTP_OK,
			common.Error(("Invalid IP:" + DtoCfg.Address)))
		return
	}
	if !isValidIP(DtoCfg.Gateway) {
		c.JSON(common.HTTP_OK,
			common.Error(("Invalid Gateway IP:" + DtoCfg.Address)))
		return
	}
	if !isValidSubnetMask(DtoCfg.Netmask) {
		c.JSON(common.HTTP_OK,
			common.Error(("Invalid SubnetMask:" + DtoCfg.Address)))
		return
	}
	for _, dns := range DtoCfg.DNS {
		if !isValidIP(dns) {
			c.JSON(common.HTTP_OK,
				common.Error(("Invalid DNS IP:" + DtoCfg.Address)))
			return
		}
	}

	MNetCfg := model.MNetworkConfig{
		Type:        "ETHNET",
		Interface:   DtoCfg.Interface,
		Address:     DtoCfg.Address,
		Netmask:     DtoCfg.Netmask,
		Gateway:     DtoCfg.Gateway,
		DNS:         DtoCfg.DNS,
		DHCPEnabled: DtoCfg.DHCPEnabled,
	}
	if err1 := service.UpdateEthConfig(MNetCfg); err1 != nil {
		if err1 != nil {
			c.JSON(common.HTTP_OK, common.Error400(err1))
			return
		}
	}
	if typex.DefaultVersionInfo.Product == "RHILEXPRO1" {
		config := []rhilexpro1.NetworkInterfaceConfig{
			{
				Interface:   MNetCfg.Interface,
				Address:     MNetCfg.Address,
				Netmask:     MNetCfg.Netmask,
				Gateway:     MNetCfg.Gateway,
				DHCPEnabled: *MNetCfg.DHCPEnabled,
			},
		}
		errSetWifi := rhilexpro1.SetEthernet(config)
		if errSetWifi != nil {
			c.JSON(common.HTTP_OK, common.Error400(errSetWifi))
			return
		}
		goto END
	}
	if typex.DefaultVersionInfo.Product == "HAAS506LD1" {
		config := []haas506.NetworkInterfaceConfig{
			{
				Interface:   MNetCfg.Interface,
				Address:     MNetCfg.Address,
				Netmask:     MNetCfg.Netmask,
				Gateway:     MNetCfg.Gateway,
				DHCPEnabled: *MNetCfg.DHCPEnabled,
			},
		}
		errSetWifi := haas506.SetEthernet(config)
		if errSetWifi != nil {
			c.JSON(common.HTTP_OK, common.Error400(errSetWifi))
			return
		}
		goto END
	}
	if typex.DefaultVersionInfo.Product == "RHILEXG1" {
		config := []rhilexg1.NetworkInterfaceConfig{
			{
				Interface:   MNetCfg.Interface,
				Address:     MNetCfg.Address,
				Netmask:     MNetCfg.Netmask,
				Gateway:     MNetCfg.Gateway,
				DHCPEnabled: *MNetCfg.DHCPEnabled,
			},
		}
		errSetWifi := rhilexg1.SetEthernet(config)
		if errSetWifi != nil {
			c.JSON(common.HTTP_OK, common.Error400(errSetWifi))
			return
		}
		goto END
	}
END:
	c.JSON(common.HTTP_OK, common.Ok())

}

func isValidSubnetMask(mask string) bool {
	// 分割子网掩码为4个整数
	parts := strings.Split(mask, ".")
	if len(parts) != 4 {
		return false
	}

	// 将每个部分转换为整数
	var octets [4]int
	for i, part := range parts {
		octet, err := strconv.Atoi(part)
		if err != nil || octet < 0 || octet > 255 {
			return false
		}
		octets[i] = octet
	}

	// 判断是否为有效的子网掩码
	var bits int
	for _, octet := range octets {
		bits += bitsInByte(octet)
	}

	return bits >= 1 && bits <= 32
}

func bitsInByte(b int) int {
	count := 0
	for b > 0 {
		count += b & 1
		b >>= 1
	}
	return count
}

func isValidIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil
}
