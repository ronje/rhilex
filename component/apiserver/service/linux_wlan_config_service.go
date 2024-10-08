package service

import (
	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/interdb"
)

/*
*
* 配置WIFI Wlan0
*
 */
func UpdateWlanConfig(MNetworkConfig model.MNetworkConfig) error {
	Model := model.MNetworkConfig{Interface: MNetworkConfig.Interface}
	return interdb.DB().
		Model(Model).
		Where("interface=? and type=\"WIFI\"", MNetworkConfig.Interface).
		Updates(MNetworkConfig).Error
}

/*
*
* 获取Wlan0的配置信息
*
 */
func GetWlanConfig(Interface string) (model.MNetworkConfig, error) {
	MWifiConfig := model.MNetworkConfig{}
	err := interdb.DB().
		Where("interface=? and type=\"WIFI\"", Interface).
		Find(&MWifiConfig).Error
	return MWifiConfig, err
}

/*
*
* 初始化网卡配置参数
*
 */
func InitWlanConfig() error {
	// 默认给DHCP
	wlan0 := model.MNetworkConfig{
		Type:      "WIFI",
		Interface: "wlan0",
		SSID:      "example.wifi",
		Password:  "123456",
		Security:  "wpa2-psk",
		Address:   "192.168.10.1",
		Netmask:   "25.255.255.0",
		Gateway:   "192.168.1.1",
		DNS: []string{
			"8.8.8.8",
		},
		DHCPEnabled: new(bool),
	}
	err := interdb.DB().
		Where("interface=? and type=\"WIFI\"", "wlan0").
		FirstOrCreate(&wlan0).Error
	if err != nil {
		return err
	}
	return nil
}
