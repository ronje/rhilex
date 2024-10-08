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

	"github.com/gin-gonic/gin"
	common "github.com/hootrhino/rhilex/component/apiserver/common"
	"github.com/hootrhino/rhilex/component/apiserver/server"
	"github.com/hootrhino/rhilex/typex"
)

func InitSysMenuPermissionRoute() {
	route := server.RouteGroup(server.ContextUrl("/menu"))
	route.GET("/main", server.AddRoute(GetSysMenus))
	route.GET("/distConfig", server.AddRoute(GetDistConfigMenus))
}

type SysMenuPermissionVo struct {
	Id     int32  `json:"id"`
	Key    string `json:"key"`
	Access bool   `json:"access"`
}

func GetSysMenus(c *gin.Context, ruleEngine typex.Rhilex) {
	allMenu := []SysMenuPermissionVo{
		{Id: 0, Key: "dashboard", Access: true},
		{Id: 1, Key: "device", Access: true},
		{Id: 2, Key: "schema", Access: true},
		{Id: 3, Key: "repository", Access: true},
		{Id: 4, Key: "inend", Access: true},
		{Id: 5, Key: "outend", Access: true},
		{Id: 6, Key: "app", Access: true},
		{Id: 7, Key: "plugin", Access: true},
		{Id: 8, Key: "module", Access: true},
		{Id: 9, Key: "system", Access: true},
	}
	c.JSON(common.HTTP_OK, common.OkWithData(allMenu))
}

func GetDistConfigMenus(c *gin.Context, ruleEngine typex.Rhilex) {
	allMenu := []SysMenuPermissionVo{
		{Id: 0, Key: "resource", Access: true},
		{Id: 2, Key: "port", Access: true},
		{Id: 9, Key: "user", Access: true},
		{Id: 8, Key: "backup", Access: true},
	}
	linuxMenu := []SysMenuPermissionVo{
		{Id: 1, Key: "netStatus", Access: true},
		{Id: 3, Key: "network", Access: true},
		{Id: 4, Key: "wifi", Access: true}, // 依赖nmcli
		{Id: 5, Key: "time", Access: true},
		{Id: 6, Key: "firmware", Access: true},
		{Id: 7, Key: "reboot", Access: true},
	}
	if runtime.GOOS == "linux" {
		allMenu = append(allMenu, linuxMenu...)
	}
	c.JSON(common.HTTP_OK, common.OkWithData(allMenu))
}
