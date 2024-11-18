package apis

import (
	"fmt"

	common "github.com/hootrhino/rhilex/component/apiserver/common"
	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/apiserver/server"
	"github.com/hootrhino/rhilex/component/apiserver/service"
	"github.com/hootrhino/rhilex/component/intercache"
	"github.com/hootrhino/rhilex/component/interdb"
	"gorm.io/gorm"

	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"

	"encoding/json"

	"github.com/gin-gonic/gin"
)

func InitDeviceRoute() {
	deviceApi := server.RouteGroup(server.ContextUrl("/devices"))
	{
		deviceApi.POST(("/create"), server.AddRoute(CreateDevice))
		deviceApi.PUT(("/update"), server.AddRoute(UpdateDevice))
		deviceApi.DELETE(("/del"), server.AddRoute(DeleteDevice))
		deviceApi.GET(("/detail"), server.AddRoute(DeviceDetail))
		deviceApi.GET("/group", server.AddRoute(ListDeviceGroup))
		deviceApi.GET("/listByGroup", server.AddRoute(ListDeviceByGroup))
		deviceApi.GET("/list", server.AddRoute(ListDevice))
		deviceApi.PUT("/restart", server.AddRoute(RestartDevice))
		deviceApi.GET("/deviceErrMsg", server.AddRoute(GetDeviceErrorMsg))
	}
}

type DeviceVo struct {
	UUID        string                 `json:"uuid"`
	Gid         string                 `json:"gid"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
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
func DeviceDetail(c *gin.Context, ruleEngine typex.Rhilex) {
	uuid, _ := c.GetQuery("uuid")
	mdev, err := service.GetMDeviceWithUUID(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400EmptyObj(err))
		return
	}
	DeviceVo := DeviceVo{}
	DeviceVo.UUID = mdev.UUID
	DeviceVo.Name = mdev.Name
	DeviceVo.Type = mdev.Type
	DeviceVo.Description = mdev.Description
	DeviceVo.Config = mdev.GetConfig()
	Slot := intercache.GetSlot("__DefaultRuleEngine")
	if Slot != nil {
		CacheValue, ok := Slot[mdev.UUID]
		if ok {
			DeviceVo.ErrMsg = CacheValue.ErrMsg
		}
	}
	//
	device := ruleEngine.GetDevice(mdev.UUID)
	if device == nil {
		DeviceVo.State = int(typex.DEV_STOP)
	} else {
		DeviceVo.State = int(device.Device.Status())
	}
	Group := service.GetResourceGroup(mdev.UUID)
	DeviceVo.Gid = Group.UUID
	c.JSON(common.HTTP_OK, common.OkWithData(DeviceVo))
}

/*
*
* 新版本的Dashboard设备不分组列表
*
 */
func ListDevice(c *gin.Context, ruleEngine typex.Rhilex) {
	pager, err := service.ReadPageRequest(c)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	count, MDevices := service.PageDevice(pager.Current, pager.Size)
	devices := []DeviceVo{}
	for _, mdev := range MDevices {
		DeviceVo := DeviceVo{}
		DeviceVo.UUID = mdev.UUID
		DeviceVo.Name = mdev.Name
		DeviceVo.Type = mdev.Type
		DeviceVo.Description = mdev.Description
		DeviceVo.Config = mdev.GetConfig()
		//
		device := ruleEngine.GetDevice(mdev.UUID)
		if device == nil {
			DeviceVo.State = int(typex.DEV_STOP)
		} else {
			DeviceVo.State = int(device.Device.Status())
		}
		Group := service.GetResourceGroup(mdev.UUID)
		DeviceVo.Gid = Group.UUID

		devices = append(devices, DeviceVo)
	}

	Result := service.WrapPageResult(*pager, devices, count)
	c.JSON(common.HTTP_OK, common.OkWithData(Result))
}

/*
*
* 分组查看
*
 */
func ListDeviceByGroup(c *gin.Context, ruleEngine typex.Rhilex) {
	pager, err := service.ReadPageRequest(c)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	Gid, _ := c.GetQuery("uuid")
	count, MDevices := service.PageDeviceByGroup(pager.Current, pager.Size, Gid)
	devices := []DeviceVo{}
	for _, mdev := range MDevices {
		DeviceVo := DeviceVo{}
		DeviceVo.UUID = mdev.UUID
		DeviceVo.Name = mdev.Name
		DeviceVo.Type = mdev.Type
		DeviceVo.Description = mdev.Description
		DeviceVo.Config = mdev.GetConfig()
		//
		device := ruleEngine.GetDevice(mdev.UUID)
		if device == nil {
			DeviceVo.State = int(typex.DEV_STOP)
		} else {
			DeviceVo.State = int(device.Device.Status())
		}
		Group := service.GetResourceGroup(mdev.UUID)
		DeviceVo.Gid = Group.UUID

		devices = append(devices, DeviceVo)
	}

	Result := service.WrapPageResult(*pager, devices, count)
	c.JSON(common.HTTP_OK, common.OkWithData(Result))
}

// 重启
func RestartDevice(c *gin.Context, ruleEngine typex.Rhilex) {
	uuid, _ := c.GetQuery("uuid")
	err := ruleEngine.RestartDevice(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())

}

// 删除设备
func DeleteDevice(c *gin.Context, ruleEngine typex.Rhilex) {
	uuid, _ := c.GetQuery("uuid")
	Mdev, err := service.GetMDeviceWithUUID(uuid)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 检查是否有规则被绑定了
	for _, ruleId := range Mdev.BindRules {
		if ruleId != "" {
			_, err0 := service.GetMRuleWithUUID(ruleId)
			if err0 != nil {
				c.JSON(common.HTTP_OK, common.Error400(err0))
				return
			}
			c.JSON(common.HTTP_OK, common.Error("Device already have rule binding:"+Mdev.BindRules.String()))
			return
		}

	}

	// GENERIC_MODBUS_MASTER需要同步删除点位表记录
	if Mdev.Type == typex.GENERIC_MODBUS_MASTER.String() {
		if err := service.DeleteAllModbusPointByDevice(uuid); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		goto NEXT
	}
	// 西门子需要同步删除点位表记录
	if Mdev.Type == typex.SIEMENS_PLC.String() {
		if err := service.DeleteAllSiemensPointByDevice(uuid); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		goto NEXT
	}
	// SNMP需要同步删除点位表记录
	if Mdev.Type == typex.GENERIC_SNMP.String() {
		if err := service.DeleteAllSnmpOidByDevice(uuid); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		goto NEXT
	}
	// BacnetIP需要同步删除点位表记录
	if Mdev.Type == typex.GENERIC_BACNET_IP.String() {
		if err := service.DeleteAllBacnetDataPointByDeviceUuid(uuid); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		goto NEXT
	}
	// BACNET GW需要同步删除点位表记录
	if Mdev.Type == typex.BACNET_ROUTER_GW.String() {
		if err := service.DeleteAllBacnetRouterPointByDeviceUuid(uuid); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		goto NEXT
	}
	// DLT6452007_MASTER需要同步删除点位表记录
	if Mdev.Type == (typex.DLT6452007_MASTER).String() {
		if err := service.DeleteAllMDlt6452007ByDevice(uuid); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		goto NEXT
	}
	// CJT1882004_MASTER需要同步删除点位表记录
	if Mdev.Type == (typex.CJT1882004_MASTER).String() {
		if err := service.DeleteAllMCjt1882004ByDevice(uuid); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		goto NEXT
	}
	// GENERIC_USER_PROTOCOL需要同步删除点位表记录
	if Mdev.Type == (typex.GENERIC_USER_PROTOCOL).String() {
		if err := service.DeleteAllMUserProtocolByDevice(uuid); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		goto NEXT
	}
	// SZY2062016_MASTER
	if Mdev.Type == (typex.SZY2062016_MASTER).String() {
		if err := service.DeleteAllMUserProtocolByDevice(uuid); err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
		goto NEXT
	}

NEXT:
	old := ruleEngine.GetDevice(uuid)
	if old != nil {
		if old.Device.Status() == typex.DEV_UP {
			old.Device.Stop()
		}
	}
	// 事务
	txErr := interdb.DB().Transaction(func(tx *gorm.DB) error {
		Group := service.GetResourceGroup(uuid)
		err3 := service.DeleteDevice(uuid)
		if err3 != nil {
			return err3
		}
		// 解除关联
		err2 := interdb.DB().Where("gid=? and rid =?", Group.UUID, uuid).
			Delete(&model.MGenericGroupRelation{}).Error
		if err2 != nil {
			return err2
		}
		ruleEngine.RemoveDevice(uuid)
		return nil
	})
	if txErr != nil {
		c.JSON(common.HTTP_OK, common.Error400(txErr))
		return
	}
	// 删除云边协同
	intercache.DeleteValue("__CecollaBinding", uuid)
	c.JSON(common.HTTP_OK, common.Ok())
}

type CecollaConfig struct {
	Enable    bool   `json:"enable"`    // 是否开启
	CecollaId string `json:"cecollaId"` // Cecolla UUID
}
type DeviceConfig struct {
	CecollaConfig CecollaConfig `json:"cecollaConfig"`
}

// 创建设备
func CreateDevice(c *gin.Context, ruleEngine typex.Rhilex) {
	form := DeviceVo{}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	configJson, err := json.Marshal(form.Config)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	deviceConfig := DeviceConfig{}
	json.Unmarshal(configJson, &deviceConfig)
	if deviceConfig.CecollaConfig.Enable {
		value := intercache.GetValue("__CecollaBinding", deviceConfig.CecollaConfig.CecollaId)
		if value.Value != nil {
			c.JSON(common.HTTP_OK, common.Error400(fmt.Errorf("Cecolla already bind to device:%s", value.Value)))
			return
		}
	}
	if err := ruleEngine.CheckDeviceType(typex.DeviceType(form.Type)); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if service.CheckDeviceNameDuplicate(form.Name) {
		c.JSON(common.HTTP_OK, common.Error("Device Name Duplicated"))
		return
	}
	// 检查个人版的创建权限: 以下三种情况，以及2个数量
	// - GENERIC_UART_RW
	// - GENERIC_MODBUS_MASTER
	// - GENERIC_MODBUS_SLAVER
	if typex.License.Type == "FREETRIAL" {
		if !utils.SContains([]string{
			"GENERIC_UART_RW",
			"GENERIC_MODBUS_MASTER",
			"GENERIC_MODBUS_SLAVER",
		}, form.Type) {
			c.JSON(common.HTTP_OK,
				common.Error("FREETRIAL Version Only Allow Create Limited Device Type"))
			return
		}
		if count := service.CheckDeviceCount(form.Name); count > (2) {
			c.JSON(common.HTTP_OK,
				common.Error("FREETRIAL Version Only Allow Create 2 Input Devices"))
			return
		}
	}

	if ok, r := utils.IsValidNameLength(form.Name); !ok {
		c.JSON(common.HTTP_OK, common.Error(r))
		return
	}
	isSingle := false
	// 内部通知单例模式
	if form.Type == typex.INTERNAL_EVENT.String() {
		for _, device := range ruleEngine.AllDevices() {
			if device.Type.String() == form.Type {
				isSingle = true
			}
		}
	}
	if isSingle {
		msg := fmt.Errorf("The %s is singleton Device, can not create multiple", form.Name)
		c.JSON(common.HTTP_OK, common.Error400(msg))
		return
	}

	newUUID := utils.DeviceUuid()
	MDevice := model.MDevice{
		UUID:        newUUID,
		Type:        form.Type,
		Name:        form.Name,
		Description: form.Description,
		Config:      string(configJson),
		BindRules:   []string{},
	}
	if err := service.InsertDevice(&MDevice); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	// 给一个分组
	if err := service.BindResource(form.Gid, MDevice.UUID); err != nil {
		c.JSON(common.HTTP_OK, common.Errorf("Group not found:%s", form.Gid))
		return
	}
	if err := server.LoadNewestDevice(newUUID, ruleEngine); err != nil {
		c.JSON(common.HTTP_OK, common.OkWithMsg(err.Error()))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())

}

// 更新设备
func UpdateDevice(c *gin.Context, ruleEngine typex.Rhilex) {

	form := DeviceVo{}
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
	deviceConfig := DeviceConfig{}
	json.Unmarshal(configJson, &deviceConfig)
	if deviceConfig.CecollaConfig.Enable {
		value := intercache.GetValue("__CecollaBinding", deviceConfig.CecollaConfig.CecollaId)
		if value.Value != form.UUID {
			c.JSON(common.HTTP_OK, common.Error400(fmt.Errorf("Cecolla already bind to device:%s", value.Value)))
			return
		}
	}
	// 检查个人版的创建权限: 以下三种情况，以及2个数量
	// - GENERIC_UART_RW
	// - GENERIC_MODBUS_MASTER
	// - GENERIC_MODBUS_SLAVER
	if typex.License.Type == "FREETRIAL" {
		if !utils.SContains([]string{
			"GENERIC_UART_RW",
			"GENERIC_MODBUS_MASTER",
			"GENERIC_MODBUS_SLAVER",
		}, form.Type) {
			c.JSON(common.HTTP_OK,
				common.Error("FREETRIAL Version Only Allow Create Limited Device Type"))
			return
		}
		if count := service.CheckDeviceCount(form.Name); count > (2) {
			c.JSON(common.HTTP_OK,
				common.Error("FREETRIAL Version Only Allow Create 2 Input Devices"))
			return
		}
	}
	//
	// 取消绑定分组,删除原来旧的分组
	txErr := service.ReBindResource(func(tx *gorm.DB) error {
		MDevice := model.MDevice{
			UUID:        form.UUID,
			Type:        form.Type,
			Name:        form.Name,
			Description: form.Description,
			Config:      string(configJson),
		}
		return tx.Model(MDevice).
			Where("uuid=?", form.UUID).
			Updates(&MDevice).Error
	}, form.UUID, form.Gid)
	if txErr != nil {
		c.JSON(common.HTTP_OK, common.Error400(txErr))
		return
	}
	if err := server.LoadNewestDevice(form.UUID, ruleEngine); err != nil {
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
func GetDeviceErrorMsg(c *gin.Context, ruleEngine typex.Rhilex) {
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
