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
	"time"

	"github.com/gin-gonic/gin"
	common "github.com/hootrhino/rhilex/component/apiserver/common"
	"github.com/hootrhino/rhilex/component/apiserver/server"
	transceiver "github.com/hootrhino/rhilex/transceiver"
	"github.com/hootrhino/rhilex/typex"
)

/*
*
* 初始化路由
*
 */
func InitTransceiverRoute() {
	route := server.RouteGroup(server.ContextUrl("/transceiver"))
	route.POST("/ctrl", server.AddRoute(TransceiverCtrl))
	route.GET("/list", server.AddRoute(TransceiverList))
	route.GET("/detail", server.AddRoute(TransceiverDetail))
}

type TransceiverInfoVo struct {
	Name     string `json:"name"`
	Model    string `json:"model"`
	Type     int    `json:"type"`
	Vendor   string `json:"vendor"`
	Mac      string `json:"mac"`
	Firmware string `json:"firmware"`
	Status   int    `json:"status"`
	ErrMsg   string `json:"errMsg"`
}

/*
*
* 通讯模块列表
*
 */
func TransceiverList(c *gin.Context, ruleEngine typex.Rhilex) {
	TransceiverInfos := []TransceiverInfoVo{}
	for _, Info := range transceiver.List() {
		TransceiverIn := TransceiverInfoVo{
			Name:     Info.Name,
			Model:    Info.Model,
			Type:     int(Info.Type),
			Vendor:   Info.Vendor,
			Mac:      Info.Mac,
			Firmware: Info.Firmware,
		}
		transceivercommunicator := transceiver.GetCommunicator(Info.Name)
		if transceivercommunicator != nil {
			Status := transceivercommunicator.Status()
			TransceiverIn.Status = int(Status.Code)
			if Status.Error != nil {
				TransceiverIn.ErrMsg = Status.Error.Error()
			}
		}
		TransceiverInfos = append(TransceiverInfos, TransceiverIn)
	}
	c.JSON(common.HTTP_OK, common.OkWithData(TransceiverInfos))
}

type TransceiverCtrlCmd struct {
	Name   string `json:"name"`
	Topic  string `json:"topic"`
	Args   string `json:"args"`
	Result any    `json:"result,omitempty"`
}

/*
*
* 控制指令
*
 */
func TransceiverCtrl(c *gin.Context, ruleEngine typex.Rhilex) {
	transceiverCtrlCmdVo := TransceiverCtrlCmd{}
	if err := c.ShouldBindJSON(&transceiverCtrlCmdVo); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	Result, Err := transceiver.Ctrl(
		transceiverCtrlCmdVo.Name,
		[]byte(transceiverCtrlCmdVo.Topic),
		[]byte(transceiverCtrlCmdVo.Args),
		300*time.Millisecond)
	if Err != nil {
		c.JSON(common.HTTP_OK, common.Error400(Err))
		return
	}
	transceiverCtrlCmdVo.Result = string(Result)
	c.JSON(common.HTTP_OK, common.OkWithData(transceiverCtrlCmdVo))
}

/*
*
* 详情数据，比如常见的指令表，资源等
*
 */
func TransceiverDetail(c *gin.Context, ruleEngine typex.Rhilex) {
	Name, _ := c.GetQuery("name")
	transceivercommunicator := transceiver.GetCommunicator(Name)
	if transceivercommunicator != nil {
		Info := transceivercommunicator.Info()
		TransceiverIn := TransceiverInfoVo{
			Name:     Info.Name,
			Model:    Info.Model,
			Type:     int(Info.Type),
			Vendor:   Info.Vendor,
			Mac:      Info.Mac,
			Firmware: Info.Firmware,
		}
		Status := transceivercommunicator.Status()
		TransceiverIn.Status = int(Status.Code)
		if Status.Error != nil {
			TransceiverIn.ErrMsg = Status.Error.Error()
		}
		c.JSON(common.HTTP_OK, common.OkWithData(TransceiverIn))
		return
	}
	c.JSON(common.HTTP_OK, common.Error("Transceiver not exists:"+Name))
}
