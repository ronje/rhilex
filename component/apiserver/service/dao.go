package service

import (
	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/interdb"
)

// -----------------------------------------------------------------------------------
func GetMRule(uuid string) (*model.MRule, error) {
	m := new(model.MRule)
	return m, interdb.DB().Where("uuid=?", uuid).First(m).Error
}
func GetAllMRule() ([]model.MRule, error) {
	m := []model.MRule{}
	return m, interdb.DB().Find(&m).Error
}

func GetMRuleWithUUID(uuid string) (*model.MRule, error) {
	m := new(model.MRule)
	return m, interdb.DB().Where("uuid=?", uuid).First(m).Error
}

func InsertMRule(r *model.MRule) error {
	return interdb.DB().Table("m_rules").Create(r).Error
}

func DeleteMRule(uuid string) error {
	return interdb.DB().Table("m_rules").Where("uuid=?", uuid).Delete(&model.MRule{}).Error
}

func UpdateMRule(uuid string, r *model.MRule) error {
	return interdb.DB().Model(r).Where("uuid=?", uuid).Updates(*r).Error
}

// -----------------------------------------------------------------------------------
func GetMInEnd(uuid string) (*model.MInEnd, error) {
	m := new(model.MInEnd)
	if err := interdb.DB().Table("m_in_ends").Where("uuid=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}
func GetMInEndWithUUID(uuid string) (*model.MInEnd, error) {
	m := new(model.MInEnd)
	if err := interdb.DB().Table("m_in_ends").Where("uuid=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func InsertMInEnd(i *model.MInEnd) error {
	return interdb.DB().Table("m_in_ends").Create(i).Error
}

func DeleteMInEnd(uuid string) error {
	return interdb.DB().Where("uuid=?", uuid).Delete(&model.MInEnd{}).Error
}

func UpdateMInEnd(uuid string, i *model.MInEnd) error {
	m := model.MInEnd{}
	if err := interdb.DB().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return err
	} else {
		interdb.DB().Model(m).Updates(*i)
		return nil
	}
}

// -----------------------------------------------------------------------------------
func GetMOutEnd(id string) (*model.MOutEnd, error) {
	m := new(model.MOutEnd)
	if err := interdb.DB().First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}
func GetMOutEndWithUUID(uuid string) (*model.MOutEnd, error) {
	m := new(model.MOutEnd)
	if err := interdb.DB().Where("uuid=?", uuid).First(m).Error; err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func InsertMOutEnd(o *model.MOutEnd) error {
	return interdb.DB().Table("m_out_ends").Create(o).Error
}

func DeleteMOutEnd(uuid string) error {
	return interdb.DB().Where("uuid=?", uuid).Delete(&model.MOutEnd{}).Error
}

func UpdateMOutEnd(uuid string, o *model.MOutEnd) error {
	m := model.MOutEnd{}
	if err := interdb.DB().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return err
	} else {
		interdb.DB().Model(m).Updates(*o)
		return nil
	}
}

func AllDevices() []model.MDevice {
	devices := []model.MDevice{}
	interdb.DB().Find(&devices)
	return devices
}

// -------------------------------------------------------------------------------------

// 获取设备列表
func GetMDeviceWithUUID(uuid string) (*model.MDevice, error) {
	m := new(model.MDevice)
	return m, interdb.DB().Where("uuid=?", uuid).First(m).Error
}

// 检查名称是否重复
func CheckDeviceCount(T string) int64 {
	Count := int64(0)
	interdb.DB().Model(model.MDevice{}).Where("type=?", T).Count(&Count)
	return Count
}

// 检查名称是否重复
func CheckNameDuplicate(name string) bool {
	Count := int64(0)
	interdb.DB().Model(model.MDevice{}).Where("name=?", name).Count(&Count)
	return Count > 0
}

// 删除设备
func DeleteDevice(uuid string) error {
	return interdb.DB().Where("uuid=?", uuid).Delete(&model.MDevice{}).Error
}

// 创建设备
func InsertDevice(o *model.MDevice) error {
	return interdb.DB().Table("m_devices").Create(o).Error
}

// 更新设备信息
func UpdateDevice(uuid string, o *model.MDevice) error {
	m := model.MDevice{}
	if err := interdb.DB().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return err
	} else {
		interdb.DB().Model(m).Updates(*o)
		return nil
	}
}

// -------------------------------------------------------------------------------------
// Goods
// -------------------------------------------------------------------------------------

// 获取Goods列表
func AllGoods() []model.MGoods {
	m := []model.MGoods{}
	interdb.DB().Find(&m)
	return m

}
func GetGoodsWithUUID(uuid string) (*model.MGoods, error) {
	m := model.MGoods{}
	if err := interdb.DB().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return nil, err
	} else {
		return &m, nil
	}
}

// 删除Goods
func DeleteGoods(uuid string) error {
	return interdb.DB().Where("uuid=?", uuid).Delete(&model.MGoods{}).Error
}

// 创建Goods
func InsertGoods(goods *model.MGoods) error {
	return interdb.DB().Table("m_goods").Create(goods).Error
}

// 更新Goods
func UpdateGoods(goods model.MGoods) error {
	return interdb.DB().Model(goods).
		Where("uuid=?", goods.UUID).Updates(&goods).Error
}

// -------------------------------------------------------------------------------------
// App Dao
// -------------------------------------------------------------------------------------

// 获取App列表
func AllApp() []model.MApplet {
	m := []model.MApplet{}
	interdb.DB().Find(&m)
	return m

}
func GetMAppWithUUID(uuid string) (*model.MApplet, error) {
	m := model.MApplet{}
	if err := interdb.DB().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return nil, err
	} else {
		return &m, nil
	}
}

// 删除App
func DeleteApp(uuid string) error {
	return interdb.DB().Where("uuid=?", uuid).Delete(&model.MApplet{}).Error
}

// 创建App
func InsertApp(app *model.MApplet) error {
	return interdb.DB().Create(app).Error
}

// 更新App
func UpdateApp(app *model.MApplet) error {
	m := model.MApplet{}
	if err := interdb.DB().Where("uuid=?", app.UUID).First(&m).Error; err != nil {
		return err
	} else {
		interdb.DB().Model(m).Where("uuid=?", app.UUID).Updates(*app)
		return nil
	}
}

// 获取AiBase列表
func AllAiBase() []model.MAiBase {
	m := []model.MAiBase{}
	interdb.DB().Find(&m)
	return m

}
func GetAiBaseWithUUID(uuid string) (*model.MAiBase, error) {
	m := model.MAiBase{}
	if err := interdb.DB().Where("uuid=?", uuid).First(&m).Error; err != nil {
		return nil, err
	} else {
		return &m, nil
	}
}

// 删除AiBase
func DeleteAiBase(uuid string) error {
	return interdb.DB().Where("uuid=?", uuid).Delete(&model.MAiBase{}).Error
}

// 创建AiBase
func InsertAiBase(AiBase *model.MAiBase) error {
	return interdb.DB().Create(AiBase).Error
}

// 更新AiBase
func UpdateAiBase(AiBase *model.MAiBase) error {
	m := model.MAiBase{}
	if err := interdb.DB().Where("uuid=?", AiBase.UUID).First(&m).Error; err != nil {
		return err
	} else {
		interdb.DB().Model(m).Updates(*AiBase)
		return nil
	}
}
