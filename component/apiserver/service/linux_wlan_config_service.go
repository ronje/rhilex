package service

import (
	"errors"

	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/interdb"
	"gorm.io/gorm"
)

/*
*
* 配置WIFI
*
 */
func UpdateWlanConfig(MNetworkConfig model.MNetworkConfig) error {
	// 使用事务来确保数据的一致性
	return interdb.DB().Transaction(func(tx *gorm.DB) error {
		// 检查记录是否存在
		var existingConfig model.MNetworkConfig
		if err := tx.Where("interface = ? AND type = ?", MNetworkConfig.Interface, "WIFI").
			First(&existingConfig).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return tx.Create(&MNetworkConfig).Error
			}
			return err
		}
		return tx.Model(&existingConfig).Updates(MNetworkConfig).Error
	})
}

/*
*
* 获取Wlan0的配置信息
*
 */
func GetWlanConfig(Interface string) (model.MNetworkConfig, error) {
	MWifiConfig := model.MNetworkConfig{}
	return MWifiConfig, interdb.DB().
		Where("interface=? and type=\"WIFI\"", Interface).
		Find(&MWifiConfig).Error
}

func AllWlanConfig() ([]model.MNetworkConfig, error) {
	MNetworkConfig := []model.MNetworkConfig{}
	err := interdb.DB().
		Where("type =\"WIFI\"").
		Find(&MNetworkConfig).Error
	return MNetworkConfig, err
}
