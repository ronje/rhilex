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
	"github.com/hootrhino/rhilex/component/apiserver/service"
	"github.com/hootrhino/rhilex/source"
	"github.com/hootrhino/rhilex/typex"
)

func InitMqttSourceServerRoute() {
	route := server.RouteGroup(server.ContextUrl("/inends"))
	route.GET(("/mqttClients"), server.AddRoute(MqttClients))
	route.DELETE(("/mqttClientsKickOut"), server.AddRoute(KickOutMqttClient))
}

/*
*
* 获取Mqtt客户端
*
 */
func MqttClients(c *gin.Context, ruleEngine typex.Rhilex) {
	uuid, _ := c.GetQuery("uuid")
	pager, err := service.ReadPageRequest(c)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	Inend := ruleEngine.GetInEnd(uuid)
	if Inend != nil {
		switch T := Inend.Source.(type) {
		case *source.MqttServer:
			Clients := T.Clients((pager.Current), (pager.Size))
			Result := service.WrapPageResult(*pager, Clients, int64(Clients.Len()))
			c.JSON(common.HTTP_OK, common.OkWithData(Result))
			return
		}
	}
	c.JSON(common.HTTP_OK, common.Error("Inend not exists"))
}

/*
*
* 踢下线
*
 */
func KickOutMqttClient(c *gin.Context, ruleEngine typex.Rhilex) {
	clientId, _ := c.GetQuery("clientId")
	uuid, _ := c.GetQuery("uuid")
	Inend := ruleEngine.GetInEnd(uuid)
	if Inend != nil {
		switch T := Inend.Source.(type) {
		case *source.MqttServer:
			Client, ok := T.FindClients(clientId)
			if ok {
				Client.Stop(nil)
			}
			c.JSON(common.HTTP_OK, common.Ok())
			return
		}
	}
	c.JSON(common.HTTP_OK, common.Error("Inend not exists:"+clientId))
}
