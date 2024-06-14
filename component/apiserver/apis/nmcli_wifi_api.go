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
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	common "github.com/hootrhino/rhilex/component/apiserver/common"
	"github.com/hootrhino/rhilex/ossupport"
	"github.com/hootrhino/rhilex/typex"
)

func isWiFiInterface(interfaceName string) bool {
	devicePath := fmt.Sprintf("/sys/class/net/%s/device", interfaceName)
	files, err := os.ReadDir(devicePath)
	if err != nil {
		return false
	}
	for _, file := range files {
		if strings.Contains(file.Name(), "802") {
			return true
		}
	}

	return false
}

/*
*
* 扫描WIFI
*
 */

func ScanWiFiSignalWithNmcli(c *gin.Context, ruleEngine typex.Rhilex) {
	interfaces, err := ossupport.GetAvailableInterfaces()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	SupportWifi := false
	for _, IFace := range interfaces {
		if isWiFiInterface(IFace.Name) {
			SupportWifi = true
			break
		}
	}
	if !SupportWifi {
		c.JSON(common.HTTP_OK, common.Error("Device not support Wifi"))
		return
	}
	wLanList := [][2]string{}
	finished := make(chan bool)
	var errSetting error
	go func() {
		errSetting = ossupport.ScanWlanList()
		if errSetting != nil {
			return
		}
		wLanList, errSetting = ossupport.GetWlanListSignal()
		if errSetting != nil {
			return
		}
		finished <- true
	}()
	select {
	case <-time.After(6 * time.Second):
		errReturn := fmt.Errorf("scan WIFI timeout")
		c.JSON(common.HTTP_OK, common.Error400(errReturn))
		return
	case <-finished:
		if errSetting != nil {
			c.JSON(common.HTTP_OK, common.Error400(errSetting))
		} else {
			c.JSON(common.HTTP_OK, common.OkWithData(wLanList))
		}
	}
}

/*
*
* 扫描WIFI
*
 */
func ScanWIFIWithNmcli(c *gin.Context, ruleEngine typex.Rhilex) {
	interfaces, err := ossupport.GetAvailableInterfaces()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	SupportWifi := false
	for _, IFace := range interfaces {
		if isWiFiInterface(IFace.Name) {
			SupportWifi = true
			break
		}
	}
	if !SupportWifi {
		c.JSON(common.HTTP_OK, common.Error("Device not support Wifi"))
		return
	}
	wLanList, err := ossupport.GetWlanListSignal()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(wLanList))
}

/*
*
* 刷新DNS, 暂且认为是ubuntu
*
 */
func RefreshDNS(c *gin.Context, ruleEngine typex.Rhilex) {
	err := ossupport.ReloadDNS()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}
