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
	"encoding/json"
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
		Apis.POST(("/create"), server.AddRoute(CreateMultiMedia))
		Apis.PUT(("/update"), server.AddRoute(UpdateMultiMedia))
		Apis.GET(("/detail"), server.AddRoute(MultiMediaDetail))
		Apis.GET("/list", server.AddRoute(ListMultiMedia))
		Apis.DELETE(("/del"), server.AddRoute(DeleteMultiMedia))
	}
}

type MultimediaConfig struct {
	StreamUrl  string `json:"streamUrl" validate:"required"`
	EnablePush *bool  `json:"enablePush"`
	PushUrl    string `json:"pushUrl"`
	EnableAi   *bool  `json:"enableAi"`
	AiModel    string `json:"aiModel"`
}

func (cfg MultimediaConfig) Validate() error {
	if cfg.StreamUrl == "" {
		return fmt.Errorf("StreamUrl is required")
	}
	if cfg.PushUrl == "" {
		return fmt.Errorf("PushUrl is required")
	}
	return nil
}
func (ms MultimediaConfig) FromString(s string) {
	json.Unmarshal([]byte(s), &ms)
}

func (ms MultimediaConfig) JsonString() string {
	jsonStr, _ := json.Marshal(ms)
	return string(jsonStr)
}

// RTSP推拉流设置参数
type MultiMediaVo struct {
	UUID        string           `json:"uuid"` // 如果空串就是新建, 非空就是更新
	Type        string           `json:"type" binding:"required"`
	Name        string           `json:"name" binding:"required"`
	Description string           `json:"description"`
	Config      MultimediaConfig `json:"config" binding:"required"`
}

func CreateMultiMedia(c *gin.Context, ruleEngine typex.Rhilex) {
	MultiMediaVo := MultiMediaVo{}
	if err := c.ShouldBindJSON(&MultiMediaVo); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 验证参数
	if !utils.SContains([]string{"RTSP", "RTMP"}, MultiMediaVo.Type) {
		c.JSON(common.HTTP_OK, common.Error400(fmt.Errorf("Invalid MultiMedia type [%s]", MultiMediaVo.Type)))
		return
	}
	// 验证其他参数
	if _, err := url.ParseRequestURI(MultiMediaVo.Config.StreamUrl); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if *MultiMediaVo.Config.EnablePush {
		if _, err := url.ParseRequestURI(MultiMediaVo.Config.PushUrl); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
	}
	if *MultiMediaVo.Config.EnablePush {
		if !utils.SContains([]string{"YOLOV8", "FACENET"}, MultiMediaVo.Config.AiModel) {
			c.JSON(common.HTTP_OK, common.Error("Only Support one of YOLOV8 or FACENET"))
			return
		}
	}
	configJson, err := json.Marshal(MultiMediaVo.Config)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if errSave := service.InsertMultiMedia(&model.MMultiMedia{
		UUID:        utils.MultimediaUuid(),
		Name:        MultiMediaVo.Name,
		Type:        MultiMediaVo.Type,
		Config:      string(configJson),
		Description: MultiMediaVo.Description,
	}); errSave != nil {
		c.JSON(common.HTTP_OK, common.Error400(errSave))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

func UpdateMultiMedia(c *gin.Context, ruleEngine typex.Rhilex) {
	MultiMediaVo := MultiMediaVo{}
	if err := c.ShouldBindJSON(&MultiMediaVo); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 验证参数
	if !utils.SContains([]string{"RTSP", "RTMP"}, MultiMediaVo.Type) {
		c.JSON(common.HTTP_OK, common.Error400(fmt.Errorf("Invalid MultiMedia type [%s]", MultiMediaVo.Type)))
		return
	}
	// 验证其他参数
	if _, err := url.ParseRequestURI(MultiMediaVo.Config.StreamUrl); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if *MultiMediaVo.Config.EnablePush {
		if _, err := url.ParseRequestURI(MultiMediaVo.Config.PushUrl); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
	}
	if *MultiMediaVo.Config.EnablePush {
		if !utils.SContains([]string{"YOLOV8", "FACENET"}, MultiMediaVo.Config.AiModel) {
			c.JSON(common.HTTP_OK, common.Error("Only Support one of YOLOV8 or FACENET"))
			return
		}
	}
	// 保存到数据库
	if errSave := service.UpdateMultiMedia(&model.MMultiMedia{
		UUID:        MultiMediaVo.UUID,
		Name:        MultiMediaVo.Name,
		Type:        MultiMediaVo.Type,
		Config:      MultiMediaVo.Config.JsonString(),
		Description: MultiMediaVo.Description,
	}); errSave != nil {
		c.JSON(common.HTTP_OK, common.Error400(errSave))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}
func MultiMediaDetail(c *gin.Context, ruleEngine typex.Rhilex) {
	uuid, ok := c.GetQuery("uuid")
	if !ok {
		c.JSON(common.HTTP_OK, common.Error400(fmt.Errorf("missing uuid")))
		return
	}
	Model, err := service.GetMultiMediaWithUUID(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 返回结果
	Config := MultimediaConfig{}
	Config.FromString(Model.Config)
	c.JSON(common.HTTP_OK, common.OkWithData(MultiMediaVo{
		UUID:        Model.UUID,
		Name:        Model.Name,
		Type:        Model.Type,
		Config:      Config,
		Description: Model.Description,
	}))
}
func ListMultiMedia(c *gin.Context, ruleEngine typex.Rhilex) {
	pager, err := service.ReadPageRequest(c)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	count, MMultiMedias, err := service.PageMultiMedia(pager.Current, pager.Size)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	MultiMedias := []MultiMediaVo{}
	for _, MMultiMedia := range MMultiMedias {
		// 返回结果
		Config := MultimediaConfig{}
		Config.FromString(MMultiMedia.Config)
		MultiMedias = append(MultiMedias, MultiMediaVo{
			UUID:        MMultiMedia.UUID,
			Name:        MMultiMedia.Name,
			Type:        MMultiMedia.Type,
			Config:      Config,
			Description: MMultiMedia.Description,
		})
	}
	Result := service.WrapPageResult(*pager, MultiMedias, count)
	c.JSON(common.HTTP_OK, common.OkWithData(Result))
}

// 删除MultiMedia
func DeleteMultiMedia(c *gin.Context, ruleEngine typex.Rhilex) {
	uuid, ok := c.GetQuery("uuid")
	if !ok {
		c.JSON(common.HTTP_OK, common.Error400(fmt.Errorf("missing uuid")))
		return
	}
	if err := service.DeleteMultiMedia(uuid); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(fmt.Errorf("delete failed: %v", err)))
		return
	}

	c.JSON(common.HTTP_OK, common.Ok())
}
