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
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rhilex/archsupport/haas506"
	"github.com/hootrhino/rhilex/archsupport/rhilexg1"
	"github.com/hootrhino/rhilex/archsupport/rhilexpro1"
	common "github.com/hootrhino/rhilex/component/apiserver/common"
	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/apiserver/service"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

/*
*
* WIFI
*
 */
type WifiConfigVo struct {
	Interface string `json:"interface"` // eth1 eth0
	SSID      string `json:"ssid"`
	Password  string `json:"password"`
	Security  string `json:"security"` // wpa2-psk wpa3-psk
}

/**
 * 获取WIFI配置
 *
 */
func GetWifi(c *gin.Context, ruleEngine typex.Rhilex) {
	iface, _ := c.GetQuery("iface")
	MWifiConfig, err := service.GetWlanConfig(iface)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	vo := WifiConfigVo{
		Interface: MWifiConfig.Interface,
		SSID:      MWifiConfig.SSID,
		Password:  MWifiConfig.Password,
		Security:  MWifiConfig.Security,
	}
	c.JSON(common.HTTP_OK, common.OkWithData(vo))

}

/*
*
*
*配置WIFI
 */
func SetWifi(c *gin.Context, ruleEngine typex.Rhilex) {
	if runtime.GOOS != "linux" {
		c.JSON(common.HTTP_OK, common.Error("OS Not Support:"+runtime.GOOS))
		return
	}
	DtoCfg := WifiConfigVo{}
	if err0 := c.ShouldBindJSON(&DtoCfg); err0 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err0))
		return
	}
	if !utils.SContains([]string{"wpa2-psk", "wpa3-psk"}, DtoCfg.Security) {
		c.JSON(common.HTTP_OK,
			common.Error(("Only support 2 valid security algorithm:wpa2-psk,wpa3-psk")))
		return
	}

	MNetCfg := model.MNetworkConfig{
		Type:      "WIFI",
		Interface: DtoCfg.Interface,
		SSID:      DtoCfg.SSID,
		Password:  DtoCfg.Password,
		Security:  DtoCfg.Security,
	}
	if err := service.UpdateWlanConfig(MNetCfg); err != nil {
		if err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
	}
	if typex.DefaultVersionInfo.Product == "RHILEXPRO1" {
		errSetWifi := rhilexpro1.SetWifi(MNetCfg.Interface, MNetCfg.SSID, MNetCfg.Password, 3*time.Second)
		if errSetWifi != nil {
			c.JSON(common.HTTP_OK, common.Error400(errSetWifi))
			return
		}
		goto END
	}
	if typex.DefaultVersionInfo.Product == "HAAS506LD1" {
		errSetWifi := haas506.SetWifi(MNetCfg.Interface, MNetCfg.SSID, MNetCfg.Password, 3*time.Second)
		if errSetWifi != nil {
			c.JSON(common.HTTP_OK, common.Error400(errSetWifi))
			return
		}
		goto END
	}
	if typex.DefaultVersionInfo.Product == "RHILEXG1" {
		errSetWifi := rhilexg1.SetWifi(MNetCfg.Interface, MNetCfg.SSID, MNetCfg.Password, 3*time.Second)
		if errSetWifi != nil {
			c.JSON(common.HTTP_OK, common.Error400(errSetWifi))
			return
		}
		goto END
	}
END:
	c.JSON(common.HTTP_OK, common.Error("Unsupported Product:"+typex.DefaultVersionInfo.Product))

}
