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
	"os"

	"github.com/gin-gonic/gin"
	common "github.com/hootrhino/rhilex/component/apiserver/common"
	"github.com/hootrhino/rhilex/periphery/haas506"
	"github.com/hootrhino/rhilex/typex"
)

/**
 * 获取基本信息
 *
 */

type NM4gInfoVo struct {
	Up    bool   `json:"up"`
	CSQ   int32  `json:"csq"`
	ICCID string `json:"iccid"`
	IMEL  string `json:"imel"`
	COPS  string `json:"cops"`
}

/**
 * 开启4G
 *
 */
func Turnon4g(c *gin.Context, ruleEngine typex.Rhilex) {
	env := os.Getenv("ARCHSUPPORT")
	if env == "HAAS506LD1" {
		err := haas506.ML307RTurnOn4G()
		if err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

/**
 * 关闭4g
 *
 */
func Turnoff4g(c *gin.Context, ruleEngine typex.Rhilex) {
	env := os.Getenv("ARCHSUPPORT")
	if env == "HAAS506LD1" {
		err := haas506.ML307RTurnOff4G()
		if err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

func Get4GBaseInfo(c *gin.Context, ruleEngine typex.Rhilex) {
	Info := NM4gInfoVo{
		CSQ:   0,
		ICCID: "UNKNOWN",
		IMEL:  "UNKNOWN",
		COPS:  "UNKNOWN",
	}
	env := os.Getenv("ARCHSUPPORT")
	if env == "HAAS506LD1" {
		BaseInfo := haas506.Get4GBaseInfo()
		Info.ICCID = BaseInfo.ICCID
		Info.IMEL = BaseInfo.IMEL
		Info.CSQ = BaseInfo.CSQ
		Info.COPS = BaseInfo.COPS
		Info.Up = BaseInfo.Up
		c.JSON(common.HTTP_OK, common.OkWithData(Info))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(Info))
}
