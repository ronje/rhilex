// Copyright (C) 2025 wwhai
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
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"
	common "github.com/hootrhino/rhilex/component/apiserver/common"
	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/apiserver/server"
	"github.com/hootrhino/rhilex/component/apiserver/service"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

func InitMultiMediaRoute() {
	Apis := server.RouteGroup(server.ContextUrl("/multimedia"))
	{
		Apis.POST(("/camera/create"), server.AddRoute(CreateCamera))
		Apis.PUT(("/camera/update"), server.AddRoute(UpdateCamera))
		Apis.GET(("/camera/detail"), server.AddRoute(CameraDetail))
		Apis.GET("/camera/list", server.AddRoute(ListCamera))
		Apis.DELETE(("/camera/del"), server.AddRoute(DeleteCamera))
	}
}

// RTSP推拉流设置参数
type CameraVo struct {
	UUID string `json:"uuid" validate:"required"`
	// 名称
	Name string `json:"name" validate:"required"`
	// 设备类型
	Type string `json:"type" validate:"required"`
	// 拉流地址
	StreamUrl string `json:"streamUrl" validate:"required"`
	// 是否开启推流
	EnablePush *bool `json:"enablePush"`
	// 推流地址
	PushUrl string `json:"pushUrl"`
	// 是否开启AI模型处理
	EnableAi *bool `json:"enableAi"`
	// AI模型选择
	AiModel string `json:"aiModel"`
}

func CreateCamera(c *gin.Context, ruleEngine typex.Rhilex) {
	EnablePush := false
	EnableAi := false
	cameraVo := CameraVo{
		EnablePush: &EnablePush,
		EnableAi:   &EnableAi,
	}
	if err := c.ShouldBindJSON(&cameraVo); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 验证参数
	if !utils.SContains([]string{"RTSP", "RTMP"}, cameraVo.Type) {
		c.JSON(common.HTTP_OK, common.Error400(fmt.Errorf("Invalid camera type [%s]", cameraVo.Type)))
		return
	}
	// 验证其他参数
	if _, err := url.ParseRequestURI(cameraVo.StreamUrl); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if *cameraVo.EnablePush {
		if _, err := url.ParseRequestURI(cameraVo.PushUrl); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
	}
	if *cameraVo.EnableAi {
		if !utils.SContains([]string{"YOLOV8", "FACENET"}, cameraVo.AiModel) {
			c.JSON(common.HTTP_OK, common.Error("Only Support one of YOLOV8 or FACENET"))
			return
		}
	}
	if errSave := service.InsertCamera(&model.MCamera{
		UUID:       utils.CameraUuid(),
		Name:       cameraVo.Name,
		Type:       cameraVo.Type,
		StreamUrl:  cameraVo.StreamUrl,
		EnablePush: cameraVo.EnablePush,
		PushUrl:    cameraVo.PushUrl,
		EnableAi:   cameraVo.EnableAi,
		AiModel:    cameraVo.AiModel,
	}); errSave != nil {
		c.JSON(common.HTTP_OK, common.Error400(errSave))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

func UpdateCamera(c *gin.Context, ruleEngine typex.Rhilex) {
	cameraVo := CameraVo{}
	if err := c.ShouldBindJSON(&cameraVo); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 验证参数
	if !utils.SContains([]string{"RTSP", "RTMP"}, cameraVo.Type) {
		c.JSON(common.HTTP_OK, common.Error400(fmt.Errorf("Invalid camera type [%s]", cameraVo.Type)))
		return
	}
	// 验证其他参数
	if _, err := url.ParseRequestURI(cameraVo.StreamUrl); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if *cameraVo.EnablePush {
		if _, err := url.ParseRequestURI(cameraVo.PushUrl); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
	}
	if *cameraVo.EnableAi {
		if !utils.SContains([]string{"YOLOV8", "FACENET"}, cameraVo.AiModel) {
			c.JSON(common.HTTP_OK, common.Error("Only Support one of YOLOV8 or FACENET"))
			return
		}
	}
	// 保存到数据库
	if errSave := service.UpdateCamera(&model.MCamera{
		UUID:       cameraVo.UUID,
		Name:       cameraVo.Name,
		Type:       cameraVo.Type,
		StreamUrl:  cameraVo.StreamUrl,
		EnablePush: cameraVo.EnablePush,
		PushUrl:    cameraVo.PushUrl,
		EnableAi:   cameraVo.EnableAi,
		AiModel:    cameraVo.AiModel,
	}); errSave != nil {
		c.JSON(common.HTTP_OK, common.Error400(errSave))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}
func CameraDetail(c *gin.Context, ruleEngine typex.Rhilex) {
	uuid, ok := c.GetQuery("uuid")
	if !ok {
		c.JSON(common.HTTP_OK, common.Error400(fmt.Errorf("missing uuid")))
		return
	}
	Model, err := service.GetCameraWithUUID(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 返回结果
	c.JSON(common.HTTP_OK, common.OkWithData(CameraVo{
		UUID:       Model.UUID,
		Name:       Model.Name,
		Type:       Model.Type,
		StreamUrl:  Model.StreamUrl,
		EnablePush: Model.EnablePush,
		PushUrl:    Model.PushUrl,
		EnableAi:   Model.EnableAi,
		AiModel:    Model.AiModel,
	}))
}
func ListCamera(c *gin.Context, ruleEngine typex.Rhilex) {
	pager, err := service.ReadPageRequest(c)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	count, MCameras, err := service.PageCamera(pager.Current, pager.Size)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	cameras := []CameraVo{}
	for _, MCamera := range MCameras {
		cameras = append(cameras, CameraVo{
			UUID:       MCamera.UUID,
			Name:       MCamera.Name,
			Type:       MCamera.Type,
			StreamUrl:  MCamera.StreamUrl,
			EnablePush: MCamera.EnablePush,
			PushUrl:    MCamera.PushUrl,
			EnableAi:   MCamera.EnableAi,
			AiModel:    MCamera.AiModel,
		})
	}
	Result := service.WrapPageResult(*pager, cameras, count)
	c.JSON(common.HTTP_OK, common.OkWithData(Result))
}

// 删除Camera
func DeleteCamera(c *gin.Context, ruleEngine typex.Rhilex) {
	uuid, ok := c.GetQuery("uuid")
	if !ok {
		c.JSON(common.HTTP_OK, common.Error400(fmt.Errorf("missing uuid")))
		return
	}
	if err := service.DeleteCamera(uuid); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(fmt.Errorf("delete failed: %v", err)))
		return
	}

	c.JSON(common.HTTP_OK, common.Ok())
}
