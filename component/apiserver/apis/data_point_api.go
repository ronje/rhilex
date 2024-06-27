package apis

import (
	"github.com/gin-gonic/gin"
	common "github.com/hootrhino/rhilex/component/apiserver/common"
	"github.com/hootrhino/rhilex/component/apiserver/dto"
	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/apiserver/server"
	"github.com/hootrhino/rhilex/component/apiserver/service"
	"github.com/hootrhino/rhilex/component/intercache"
	"github.com/hootrhino/rhilex/component/interdb"
	"github.com/hootrhino/rhilex/typex"
)

func InitDataPointRoute() {
	route := server.RouteGroup(server.ContextUrl("/datapoint"))
	route.POST(("/sheetImport"), server.AddRoute(DataPointSheetImport))
	route.GET(("/sheetExport"), server.AddRoute(DataPointSheetExport))
	route.GET(("/list"), server.AddRoute(DataPointSheetPageList))
	route.POST(("/update"), server.AddRoute(DataPointSheetCreateOrUpdate))
	route.DELETE(("/delIds"), server.AddRoute(DataPointSheetDeleteByUUIDs))
	route.DELETE(("/delAll"), server.AddRoute(DataPointSheetDeleteAll))
}

func DataPointSheetImport(c *gin.Context, ruleEngine typex.Rhilex) {
	// device_uuid
	// file: points
	// import framework
	// TODO validator
}

func DataPointSheetExport(c *gin.Context, ruleEngine typex.Rhilex) {
	// device_uuid
	// export framework
	// TODO validator
}

func DataPointSheetPageList(c *gin.Context, ruleEngine typex.Rhilex) {
	pager, err := service.ReadPageRequest(c)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	deviceUuid, _ := c.GetQuery("device_uuid")
	db := interdb.DB()
	tx := db.Scopes(service.Paginate(*pager))
	var count int64
	err = interdb.DB().Model(&model.MDataPoint{}).
		Where("device_uuid=?", deviceUuid).Count(&count).Error
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	var records []model.MDataPoint
	result := tx.Order("created_at DESC").Find(&records,
		&model.MDataPoint{DeviceUuid: deviceUuid})
	if result.Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(result.Error))
		return
	}

	var recordsVo []dto.DataPointVO
	Slot := intercache.GetSlot(deviceUuid)
	if Slot != nil {
		for _, record := range records {
			value, ok := Slot[record.UUID]
			pointVo := dto.DataPointVO{
				UUID:       record.UUID,
				DeviceUUID: record.DeviceUuid,
				Tag:        record.Tag,
				Alias:      record.Alias,
				Config:     record.GetConfig(),
				ErrMsg:     value.ErrMsg,
			}
			if ok {
				pointVo.Status = func() uint32 {
					if value.Value == "" {
						return 0
					}
					return 1
				}()
				pointVo.LastFetchTime = value.LastFetchTime
				pointVo.Value = value.Value
				recordsVo = append(recordsVo, pointVo)
			} else {
				recordsVo = append(recordsVo, pointVo)
			}
		}
	}

	Result := service.WrapPageResult(*pager, recordsVo, count)
	c.JSON(common.HTTP_OK, common.OkWithData(Result))
}

func DataPointSheetCreateOrUpdate(c *gin.Context, ruleEngine typex.Rhilex) {
	type Form struct {
		DeviceUUID string                              `json:"device_uuid"`
		Points     []dto.BacnetDataPointCreateOrUpdate `json:"points"`
	}
	form := Form{}
	err := c.ShouldBindJSON(&form)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// TODO 从device_uuid中获取其类型，类型获取对应的validator、import、export方法
	ruleEngine.RestartDevice(form.DeviceUUID)
	c.JSON(common.HTTP_OK, common.Ok())
}

func DataPointSheetDeleteByUUIDs(c *gin.Context, ruleEngine typex.Rhilex) {
	type Form struct {
		UUIDs      []string `json:"uuids"`
		DeviceUUID string   `json:"device_uuid"`
	}
	form := Form{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	err := service.DeleteByUuids(form.DeviceUUID, form.UUIDs)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	ruleEngine.RestartDevice(form.DeviceUUID)
	c.JSON(common.HTTP_OK, common.Ok())
}

func DataPointSheetDeleteAll(c *gin.Context, ruleEngine typex.Rhilex) {
	type Form struct {
		DeviceUUID string `json:"device_uuid"`
	}
	form := Form{}
	if Error := c.ShouldBindJSON(&form); Error != nil {
		c.JSON(common.HTTP_OK, common.Error400(Error))
		return
	}
	err := service.DeleteAllByDeviceUuid(form.DeviceUUID)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	ruleEngine.RestartDevice(form.DeviceUUID)
	c.JSON(common.HTTP_OK, common.Ok())
}
