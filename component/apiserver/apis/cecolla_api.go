package apis

import (
	"fmt"

	common "github.com/hootrhino/rhilex/component/apiserver/common"
	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/apiserver/server"
	"github.com/hootrhino/rhilex/component/apiserver/service"
	"github.com/hootrhino/rhilex/component/intercache"
	"github.com/hootrhino/rhilex/component/interdb"
	"github.com/hootrhino/rhilex/component/luaruntime"
	"gorm.io/gorm"

	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"

	"encoding/json"

	"github.com/gin-gonic/gin"
)

func InitCecollaRoute() {
	cecollaApi := server.RouteGroup(server.ContextUrl("/cecollas"))
	{
		cecollaApi.POST(("/create"), server.AddRoute(CreateCecolla))
		cecollaApi.PUT(("/update"), server.AddRoute(UpdateCecolla))
		cecollaApi.PUT(("/updateAction"), server.AddRoute(UpdateCecollaAction))
		cecollaApi.DELETE(("/del"), server.AddRoute(DeleteCecolla))
		cecollaApi.GET(("/detail"), server.AddRoute(CecollaDetail))
		cecollaApi.GET("/group", server.AddRoute(ListCecollaByGroup))
		cecollaApi.GET("/listByGroup", server.AddRoute(ListCecollaByGroup))
		cecollaApi.GET("/list", server.AddRoute(ListCecolla))
		cecollaApi.PUT("/restart", server.AddRoute(RestartCecolla))
		cecollaApi.GET("/cecollaErrMsg", server.AddRoute(GetCecollaErrorMsg))
		cecollaApi.GET("/cecollaSchema", server.AddRoute(GetCecollaSchema))
	}
}

type CecollaVo struct {
	UUID        string                 `json:"uuid"`
	Gid         string                 `json:"gid"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Action      string                 `json:"action"`
	State       int                    `json:"state"`
	ErrMsg      string                 `json:"errMsg"`
	Config      map[string]interface{} `json:"config"`
	Description string                 `json:"description"`
}

/*
*
* 列表先读数据库，然后读内存，合并状态后输出
*
 */
func CecollaDetail(c *gin.Context, ruleEngine typex.Rhilex) {
	uuid, _ := c.GetQuery("uuid")
	mCecolla, err := service.GetMCecollaWithUUID(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400EmptyObj(err))
		return
	}
	CecollaVo := CecollaVo{}
	CecollaVo.UUID = mCecolla.UUID
	CecollaVo.Name = mCecolla.Name
	CecollaVo.Type = mCecolla.Type
	CecollaVo.Action = mCecolla.Action
	CecollaVo.Description = mCecolla.Description
	CecollaVo.Config = mCecolla.GetConfig()
	Slot := intercache.GetSlot("__DefaultRuleEngine")
	if Slot != nil {
		CacheValue, ok := Slot[mCecolla.UUID]
		if ok {
			CecollaVo.ErrMsg = CacheValue.ErrMsg
		}
	}
	//
	cecolla := ruleEngine.GetCecolla(mCecolla.UUID)
	if cecolla == nil {
		CecollaVo.State = int(typex.DEV_STOP)
	} else {
		CecollaVo.State = int(cecolla.Cecolla.Status())
	}
	Group := service.GetResourceGroup(mCecolla.UUID)
	CecollaVo.Gid = Group.UUID
	c.JSON(common.HTTP_OK, common.OkWithData(CecollaVo))
}

/*
*
* 新版本的Dashboard设备不分组列表
*
 */
func ListCecolla(c *gin.Context, ruleEngine typex.Rhilex) {
	pager, err := service.ReadPageRequest(c)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	count, MCecollas := service.PageCecolla(pager.Current, pager.Size)
	cecollas := []CecollaVo{}
	for _, mCecolla := range MCecollas {
		CecollaVo := CecollaVo{}
		CecollaVo.UUID = mCecolla.UUID
		CecollaVo.Name = mCecolla.Name
		CecollaVo.Type = mCecolla.Type
		CecollaVo.Action = mCecolla.Action
		CecollaVo.Description = mCecolla.Description
		CecollaVo.Config = mCecolla.GetConfig()
		//
		cecolla := ruleEngine.GetCecolla(mCecolla.UUID)
		if cecolla == nil {
			CecollaVo.State = int(typex.DEV_STOP)
		} else {
			CecollaVo.State = int(cecolla.Cecolla.Status())
		}
		Group := service.GetResourceGroup(mCecolla.UUID)
		CecollaVo.Gid = Group.UUID

		cecollas = append(cecollas, CecollaVo)
	}

	Result := service.WrapPageResult(*pager, cecollas, count)
	c.JSON(common.HTTP_OK, common.OkWithData(Result))
}

/*
*
* 分组查看
*
 */
func ListCecollaByGroup(c *gin.Context, ruleEngine typex.Rhilex) {
	pager, err := service.ReadPageRequest(c)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	Gid, _ := c.GetQuery("gid")
	count, MCecollas := service.PageCecollaByGroup(pager.Current, pager.Size, Gid)
	cecollas := []CecollaVo{}
	for _, mCecolla := range MCecollas {
		CecollaVo := CecollaVo{}
		CecollaVo.UUID = mCecolla.UUID
		CecollaVo.Name = mCecolla.Name
		CecollaVo.Type = mCecolla.Type
		CecollaVo.Action = mCecolla.Action
		CecollaVo.Description = mCecolla.Description
		CecollaVo.Config = mCecolla.GetConfig()
		//
		cecolla := ruleEngine.GetCecolla(mCecolla.UUID)
		if cecolla == nil {
			CecollaVo.State = int(typex.DEV_STOP)
		} else {
			CecollaVo.State = int(cecolla.Cecolla.Status())
		}
		Group := service.GetResourceGroup(mCecolla.UUID)
		CecollaVo.Gid = Group.UUID

		cecollas = append(cecollas, CecollaVo)
	}

	Result := service.WrapPageResult(*pager, cecollas, count)
	c.JSON(common.HTTP_OK, common.OkWithData(Result))
}

// 重启
func RestartCecolla(c *gin.Context, ruleEngine typex.Rhilex) {
	uuid, _ := c.GetQuery("uuid")
	err := ruleEngine.RestartCecolla(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())

}

// 删除设备
func DeleteCecolla(c *gin.Context, ruleEngine typex.Rhilex) {
	uuid, _ := c.GetQuery("uuid")
	txErr := interdb.DB().Transaction(func(tx *gorm.DB) error {
		Group := service.GetResourceGroup(uuid)
		err3 := service.DeleteCecolla(uuid)
		if err3 != nil {
			return err3
		}
		// 解除关联
		err2 := interdb.DB().Where("gid=? and rid =?", Group.UUID, uuid).
			Delete(&model.MGenericGroupRelation{}).Error
		if err2 != nil {
			return err2
		}
		ruleEngine.RemoveCecolla(uuid)
		return nil
	})
	if txErr != nil {
		c.JSON(common.HTTP_OK, common.Error400(txErr))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

// 创建设备
func CreateCecolla(c *gin.Context, ruleEngine typex.Rhilex) {
	form := CecollaVo{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	configJson, err := json.Marshal(form.Config)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if err := ruleEngine.CheckCecollaType(typex.CecollaType(form.Type)); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if service.CheckCecollaNameDuplicate(form.Name) {
		c.JSON(common.HTTP_OK, common.Error("Cecolla Name Duplicated"))
		return
	}
	template :=
		`--
-- Go https://www.hootrhino.com for more tutorials
--
-- ID: %s
-- NAME = "%s"
-- DESCRIPTION = "%s"
--
-- Action Main
--
function Main(Payload)
	Debug("== Payload ==" .. Payload)
end
`
	newUUID := utils.CecUuid()
	MCecolla := model.MCecolla{
		UUID:        newUUID,
		Type:        form.Type,
		Name:        form.Name,
		Description: form.Description,
		Config:      string(configJson),
	}
	MCecolla.Action = fmt.Sprintf(template, MCecolla.UUID, MCecolla.Name, MCecolla.Description)
	if err := service.InsertCecolla(&MCecolla); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 给一个分组
	if err := service.BindResource(form.Gid, MCecolla.UUID); err != nil {
		c.JSON(common.HTTP_OK, common.Errorf("Group not found:%s", form.Gid))
		return
	}
	if err := server.LoadNewestCecolla(newUUID, ruleEngine); err != nil {
		c.JSON(common.HTTP_OK, common.OkWithMsg(err.Error()))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())

}

/**
 * 更新物模型
 *
 */
func UpdateCecollaAction(c *gin.Context, ruleEngine typex.Rhilex) {
	type Form struct {
		UUID   string `json:"uuid"`
		Action string `json:"action"`
	}
	form := Form{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if err := luaruntime.ValidateCecollaletSyntax([]byte(form.Action)); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	mCecolla, err := service.GetMCecollaWithUUID(form.UUID)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	mCecolla.Action = form.Action
	if err := service.UpdateCecolla(form.UUID, mCecolla); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())

}

// 更新设备
func UpdateCecolla(c *gin.Context, ruleEngine typex.Rhilex) {

	form := CecollaVo{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	configJson, err := json.Marshal(form.Config)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if ok, r := utils.IsValidNameLength(form.Name); !ok {
		c.JSON(common.HTTP_OK, common.Error(r))
		return
	}
	//
	// 取消绑定分组,删除原来旧的分组
	txErr := service.ReBindResource(func(tx *gorm.DB) error {
		MCecolla := model.MCecolla{
			UUID:        form.UUID,
			Type:        form.Type,
			Name:        form.Name,
			Action:      "",
			Description: form.Description,
			Config:      string(configJson),
		}
		return tx.Model(MCecolla).
			Where("uuid=?", form.UUID).
			Updates(&MCecolla).Error
	}, form.UUID, form.Gid)
	if txErr != nil {
		c.JSON(common.HTTP_OK, common.Error400(txErr))
		return
	}
	if err := server.LoadNewestCecolla(form.UUID, ruleEngine); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* 获取设备挂了的异常信息
* __DefaultRuleEngine：用于RHILEX内部存储一些KV键值对
 */
func GetCecollaErrorMsg(c *gin.Context, ruleEngine typex.Rhilex) {
	uuid, _ := c.GetQuery("uuid")
	Slot := intercache.GetSlot("__DefaultRuleEngine")
	if Slot != nil {
		CacheValue, ok := Slot[uuid]
		if ok {
			c.JSON(common.HTTP_OK, common.OkWithData(CacheValue.ErrMsg))
			return
		}
	}
	c.JSON(common.HTTP_OK, common.OkWithData("--"))
}

/**
 * 获取物模型
 *
 */
func GetCecollaSchema(c *gin.Context, ruleEngine typex.Rhilex) {
	uuid, _ := c.GetQuery("uuid")
	cecolla := ruleEngine.GetCecolla(uuid)
	if cecolla == nil {
		c.JSON(common.HTTP_OK, common.OkWithMsg("Cecolla Not Exists:"+uuid))
		return
	}
	response, err := cecolla.Cecolla.OnCtrl([]byte("GetSchema"), []byte{})
	if err != nil {
		c.JSON(common.HTTP_OK, common.OkWithMsg(err.Error()))
		return
	}
	switch T := response.(type) {
	case map[string]any:
		c.JSON(common.HTTP_OK, common.OkWithData(T))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithMsg("Empty Schema"))
}
