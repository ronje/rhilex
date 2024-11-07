package service

import (
	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/interdb"
)

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
