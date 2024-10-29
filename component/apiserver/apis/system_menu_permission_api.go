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

type SysMenu struct {
	Id       int32          `json:"id"`
	Group    string         `json:"group"`
	Key      string         `json:"key"`
	Access   bool           `json:"access"`
	Children []SysMenuChild `json:"children"`
}

type SysMenuChild struct {
	Id     int32  `json:"id"`
	Key    string `json:"key"`
	Access bool   `json:"access"`
}

func GetSysMenus(c *gin.Context, ruleEngine typex.Rhilex) {
	allMenu := []SysMenu{
		{Group: "root", Children: []SysMenuChild{}, Id: 0, Key: "dashboard", Access: true},
		{Group: "root", Children: []SysMenuChild{}, Id: 1, Key: "device", Access: true},
		{Group: "root", Children: []SysMenuChild{}, Id: 2, Key: "schema", Access: true},
		{Group: "root", Children: []SysMenuChild{}, Id: 3, Key: "repository", Access: true},
		{Group: "root", Children: []SysMenuChild{}, Id: 4, Key: "inend", Access: true},
		{Group: "root", Children: []SysMenuChild{}, Id: 5, Key: "outend", Access: true},
		{Group: "root", Children: []SysMenuChild{}, Id: 6, Key: "app", Access: true},
		{Group: "root", Children: []SysMenuChild{}, Id: 7, Key: "plugin", Access: true},
		{Group: "root", Children: []SysMenuChild{}, Id: 8, Key: "module", Access: true},
		{Group: "root", Children: []SysMenuChild{}, Id: 9, Key: "system", Access: true},
	}
	c.JSON(common.HTTP_OK, common.OkWithData(allMenu))
}

// 系统资源

// 网络配置
//     网络状态
//     网卡设置
//     WIFI设置
//     4G设置
//     5G设置
//     CAN设置

// 时间设置
//     系统时间
//     定时重启

// 系统版本
//     固件设置
//     数据备份

// 用户设置
func GetDistConfigMenus(c *gin.Context, ruleEngine typex.Rhilex) {
	allMenu := []SysMenu{
		{Id: 0, Group: "resource", Key: "resource", Access: true, Children: []SysMenuChild{}},
		{Id: 1, Group: "nerworking", Key: "nerworking", Access: true,
			Children: []SysMenuChild{
				{Id: 1, Key: "netStatus", Access: true},
				{Id: 2, Key: "network", Access: true},
				{Id: 3, Key: "wifi", Access: true},
				{Id: 4, Key: "net4g", Access: true},
				{Id: 5, Key: "net5g", Access: false},
				{Id: 6, Key: "can", Access: false},
			},
		},
		{Id: 2, Group: "datetime", Key: "datetime", Access: true,
			Children: []SysMenuChild{
				{Id: 1, Key: "time", Access: true},
				{Id: 2, Key: "reboot", Access: true},
			},
		},
		{Id: 3, Group: "sysver", Key: "sysver", Access: true,
			Children: []SysMenuChild{
				{Id: 1, Key: "firmware", Access: true},
				{Id: 2, Key: "backup", Access: true},
			},
		},
		{Id: 4, Group: "user", Key: "user", Access: true, Children: []SysMenuChild{}},
	}
	c.JSON(common.HTTP_OK, common.OkWithData(allMenu))
}
