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
	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/apiserver/server"
	"github.com/hootrhino/rhilex/component/apiserver/service"
	"github.com/hootrhino/rhilex/component/interdb"
	"github.com/hootrhino/rhilex/typex"
)

func InitInternalNotifyRoute() {
	// 站内公告
	internalNotifyApi := server.DefaultApiServer.GetGroup(server.ContextUrl("/notify"))
	{
		internalNotifyApi.PUT("/clear", server.AddRoute(ClearInternalNotifies))
		internalNotifyApi.PUT("/read", server.AddRoute(ReadInternalNotifies))
		internalNotifyApi.GET("/pageList", server.AddRoute(PageInternalNotifies))
	}
}

/*
*
* 内部事件
*
 */
type InternalNotifyVo struct {
	UUID    string `json:"uuid"`           // UUID
	Type    string `json:"type"`           // INFO | ERROR | WARNING
	Status  int    `json:"status"`         // 1 未读 2 已读
	Event   string `json:"event"`          // 字符串
	Ts      uint64 `json:"ts"`             // 时间戳
	Summary string `json:"summary"`        // 概览，为了节省流量，在消息列表只显示这个字段，Info值为“”
	Info    string `json:"info,omitempty"` // 消息内容，是个文本，详情显示
}

/*
*
* 站内消息
*
 */
func InternalNotifies(c *gin.Context, ruleEngine typex.Rhilex) {
	data := []InternalNotifyVo{}
	models := service.AllInternalNotifies()
	for _, model := range models {
		data = append(data, InternalNotifyVo{
			UUID:    model.UUID,
			Type:    model.Type,
			Event:   model.Event,
			Ts:      model.Ts,
			Summary: model.Summary,
			Info:    model.Info,
			Status:  model.Status,
		})

	}
	c.JSON(common.HTTP_OK, common.OkWithData(data))
}

/*
*
* 支持分页
*
 */
func PageInternalNotifies(c *gin.Context, ruleEngine typex.Rhilex) {
	pager, err := service.ReadPageRequest(c)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if pager.Size > 100 {
		c.JSON(common.HTTP_OK, common.Error("Query size too large, Must less than 100"))
		return
	}
	DbTx := interdb.DB().Scopes(service.Paginate(*pager))
	records := []InternalNotifyVo{}
	result := DbTx.Model(model.MInternalNotify{}).Order("id DESC").Scan(&records)
	if result.Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(result.Error))
		return
	}
	var count int64
	err1 := interdb.DB().Model(&model.MInternalNotify{}).Count(&count).Error
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	Result := service.WrapPageResult(*pager, records, count)
	c.JSON(common.HTTP_OK, common.OkWithData(Result))
}

/*
*
* 清空
*
 */
func ClearInternalNotifies(c *gin.Context, ruleEngine typex.Rhilex) {
	if err := service.ClearInternalNotifies(); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* 阅读
*
 */
func ReadInternalNotifies(c *gin.Context, ruleEngine typex.Rhilex) {
	uuid, _ := c.GetQuery("uuid")
	if err := service.ReadInternalNotifies(uuid); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}
