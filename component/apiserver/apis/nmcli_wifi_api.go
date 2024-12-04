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
	"github.com/hootrhino/rhilex/periphery/haas506"
	"github.com/hootrhino/rhilex/periphery/rhilexg1"
	"github.com/hootrhino/rhilex/periphery/rhilexpro1"
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

func ScanWIFIList(c *gin.Context, ruleEngine typex.Rhilex) {
	iface, _ := c.GetQuery("iface")
	if !isWiFiInterface(iface) {
		c.JSON(common.HTTP_OK, common.Error("invalid wlan interface"))
		return
	}
	wLanList := [][2]string{}
	finished := make(chan bool)
	var errSetWifi error
	go func() {
		if typex.DefaultVersionInfo.Product == "RHILEXPRO1" {
			wLanList, errSetWifi = rhilexpro1.ScanWlanList(iface)
			if errSetWifi != nil {
				c.JSON(common.HTTP_OK, common.Error400(errSetWifi))
				return
			}
			goto END
		}
		if typex.DefaultVersionInfo.Product == "HAAS506LD1" {
			wLanList, errSetWifi = haas506.ScanWlanList(iface)
			if errSetWifi != nil {
				c.JSON(common.HTTP_OK, common.Error400(errSetWifi))
				return
			}
			goto END
		}
		if typex.DefaultVersionInfo.Product == "RHILEXG1" {
			wLanList, errSetWifi = rhilexg1.ScanWlanList(iface)
			if errSetWifi != nil {
				c.JSON(common.HTTP_OK, common.Error400(errSetWifi))
				return
			}
			goto END
		}
		if errSetWifi != nil {
			finished <- true
			return
		}
	END:
		finished <- true
	}()
	select {
	case <-time.After(10 * time.Second):
		errReturn := fmt.Errorf("scan WIFI timeout")
		c.JSON(common.HTTP_OK, common.Error400(errReturn))
		return
	case <-finished:
		if errSetWifi != nil {
			c.JSON(common.HTTP_OK, common.Error400(errSetWifi))
		} else {
			c.JSON(common.HTTP_OK, common.OkWithData(wLanList))
		}
	}
}
