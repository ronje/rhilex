package service

import (
	"errors"

	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/interdb"
	"gorm.io/gorm"
)

func GetEthConfig(Interface string) (model.MNetworkConfig, error) {
	MNetworkConfig := model.MNetworkConfig{}
	result := interdb.InterDb().
		Where("interface=? AND type =\"ETHNET\"", Interface).
		Find(&MNetworkConfig)
	return MNetworkConfig, result.Error
}

func UpdateEthConfig(MNetworkConfig model.MNetworkConfig) error {
	return interdb.InterDb().Transaction(func(tx *gorm.DB) error {
		var existingConfig model.MNetworkConfig
		if err := tx.Where("interface = ? AND type = ?", MNetworkConfig.Interface, "ETHNET").
			First(&existingConfig).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return tx.Create(&MNetworkConfig).Error
			}
			return err
		}
		return tx.Model(&existingConfig).Updates(MNetworkConfig).Error
	})
}

func AllEthConfig() ([]model.MNetworkConfig, error) {
	MNetworkConfig := []model.MNetworkConfig{}
	err := interdb.InterDb().
		Where("type =\"ETHNET\"").
		Find(&MNetworkConfig).Error
	return MNetworkConfig, err
}
