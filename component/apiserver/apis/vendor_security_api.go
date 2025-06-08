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
	"github.com/gin-gonic/gin"
	common "github.com/hootrhino/rhilex/component/apiserver/common"
	"github.com/hootrhino/rhilex/typex"
)

/*
*
* 获取一机一密
*
 */

func GetVendorKey(c *gin.Context, ruleEngine typex.Rhilex) {
	type LocalLicense struct {
		Type              string `json:"type"` // FREETRIAL | COMMERCIAL
		DeviceID          string `json:"device_id"`
		AuthorizeAdmin    string `json:"authorize_admin"`
		AuthorizePassword string `json:"authorize_password"`
		BeginAuthorize    int64  `json:"begin_authorize"`
		EndAuthorize      int64  `json:"end_authorize"`
		Iface             string `json:"iface"`
		MAC               string `json:"mac"`
		License           string `json:"license"`
	}
	c.JSON(common.HTTP_OK, common.OkWithData(LocalLicense{
		Type:              "COMMERCIAL",
		DeviceID:          "RHILEX-0001",
		AuthorizeAdmin:    "admin",
		AuthorizePassword: "password",
		BeginAuthorize:    0,
		EndAuthorize:      0,
		Iface:             "eth0",
		MAC:               "00:00:00:00:00:00",
		License:           "RHILEX-0001",
	}))
}
