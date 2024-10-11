package service

import (
	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/interdb"
)

/*
*
* 配置WIFI
*
 */
func UpdateWlanConfig(MNetworkConfig model.MNetworkConfig) error {
	Model := model.MNetworkConfig{Interface: MNetworkConfig.Interface}
	return interdb.DB().
		Model(Model).
		Where("interface=? and type=\"WIFI\"", MNetworkConfig.Interface).
		FirstOrCreate(&MNetworkConfig).Error
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
