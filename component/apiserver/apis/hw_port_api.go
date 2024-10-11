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
	"encoding/json"

	"github.com/gin-gonic/gin"
	common "github.com/hootrhino/rhilex/component/apiserver/common"
	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/apiserver/server"
	"github.com/hootrhino/rhilex/component/apiserver/service"
	"github.com/hootrhino/rhilex/component/uartctrl"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

func InitHwIfaceRoute() {
	// 硬件接口API
	HwIFaceApi := server.DefaultApiServer.GetGroup(server.ContextUrl("/hwiface"))
	{
		HwIFaceApi.GET("/detail", server.AddRoute(GetUartDetail))
		HwIFaceApi.GET("/list", server.AddRoute(AllUarts))
		HwIFaceApi.POST("/update", server.AddRoute(UpdateUartConfig))
		HwIFaceApi.GET("/refresh", server.AddRoute(RefreshPortList))
	}
}

type UartVo struct {
	UUID        string       `json:"uuid"`
	Name        string       `json:"name"`   // 接口名称
	Type        string       `json:"type"`   // 接口类型, UART(串口),USB(USB),FD(通用文件句柄)
	Alias       string       `json:"alias"`  // 别名
	Config      any          `json:"config"` // 配置
	Busy        bool         `json:"busy"`   // 运行时数据，是否被占
	OccupyBy    UartOccupyVo `json:"occupyBy"`
	Description string       `json:"description"` // 额外备注

}
type UartOccupyVo struct {
	UUID string `json:"uuid"` // UUID
	Type string `json:"type"` // DEVICE, Other......
	Name string `json:"name"` // DEVICE, Other......
}
type UartConfigVo struct {
	Timeout  int    `json:"timeout"`
	Uart     string `json:"uart"`
	BaudRate int    `json:"baudRate"`
	DataBits int    `json:"dataBits"`
	Parity   string `json:"parity"`
	StopBits int    `json:"stopBits"`
}

func (u UartConfigVo) JsonString() string {
	if bytes, err := json.Marshal(u); err != nil {
		return "{}"
	} else {
		return string(bytes)
	}
}

/*
*
* 针对刚插入硬件的情况，需要及时刷新
*
 */
func RefreshPortList(c *gin.Context, ruleEngine typex.Rhilex) {
	if err := service.ReScanUartConfig(); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* 硬件接口
*
 */
func AllUarts(c *gin.Context, ruleEngine typex.Rhilex) {
	UartVos := []UartVo{}
	for _, port := range uartctrl.AllUart() {
		UartVos = append(UartVos, UartVo{
			UUID:  port.UUID,
			Name:  port.Name,
			Type:  port.Type,
			Alias: port.Alias,
			Busy:  port.Busy,
			OccupyBy: UartOccupyVo{
				UUID: port.OccupyBy.UUID,
				Type: port.OccupyBy.Type,
				Name: port.OccupyBy.Name,
			},
			Description: port.Description,
		})
	}
	c.JSON(common.HTTP_OK, common.OkWithData(UartVos))
}

/*
*
* 更新接口参数
*
 */
func UpdateUartConfig(c *gin.Context, ruleEngine typex.Rhilex) {
	type Form struct {
		UUID        string       `json:"uuid"`
		Config      UartConfigVo `json:"config"` // 配置, 串口配置、或者网卡、USB等
		Description string       `json:"description"`
	}
	form := Form{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if err := service.UpdateUartConfig(model.MUart{
		UUID:        form.UUID,
		Config:      form.Config.JsonString(),
		Description: form.Description,
	}); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	MUart, err1 := service.GetUartConfig(form.UUID)
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	HwIPort := uartctrl.SystemUart{
		UUID:        MUart.UUID,
		Name:        MUart.Name,
		Type:        MUart.Type,
		Alias:       MUart.Alias,
		Description: MUart.Description,
	}
	// 串口类
	if MUart.Type == "UART" {
		config := uartctrl.UartConfig{}
		utils.BindSourceConfig(MUart.GetConfig(), &config)
		HwIPort.Config = config
	}
	if MUart.Type == "FD" {
		HwIPort.Config = nil
	}
	// 刷新接口参数
	uartctrl.RefreshPort(HwIPort)
	c.JSON(common.HTTP_OK, common.Ok())

}

/*
*
* 获取详情
*
 */
func GetUartDetail(c *gin.Context, ruleEngine typex.Rhilex) {
	uuid, _ := c.GetQuery("uuid")
	Port, err1 := uartctrl.GetUart(uuid)
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(UartVo{
		UUID:        Port.UUID,
		Name:        Port.Name,
		Type:        Port.Type,
		Alias:       Port.Alias,
		Config:      Port.Config,
		Description: Port.Description,
		Busy:        Port.Busy,
		OccupyBy: UartOccupyVo{
			Port.OccupyBy.UUID, Port.OccupyBy.Type, Port.Name,
		},
	}))

}
