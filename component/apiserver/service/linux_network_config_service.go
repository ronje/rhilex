package service

import (
	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/interdb"
)

func GetEthConfig(Interface string) (model.MNetworkConfig, error) {
	MNetworkConfig := model.MNetworkConfig{}
	err := interdb.DB().
		Where("interface=? AND type =\"ETHNET\"", Interface).
		Find(&MNetworkConfig).Error
	return MNetworkConfig, err
}

func UpdateEthConfig(MNetworkConfig model.MNetworkConfig) error {
	return interdb.DB().
		Model(&MNetworkConfig).
		Where("interface=? AND type =\"ETHNET\"", MNetworkConfig.Interface).
		FirstOrCreate(&MNetworkConfig).Error
}

func AllEthConfig() ([]model.MNetworkConfig, error) {
	MNetworkConfig := []model.MNetworkConfig{}
	err := interdb.DB().
		Where("type =\"ETHNET\"").
		Find(&MNetworkConfig).Error
	return MNetworkConfig, err
}
