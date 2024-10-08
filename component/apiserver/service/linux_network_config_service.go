package service

import (
	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/interdb"
)

/*
*
  - 默认静态IP
    114DNS:
    IPv4: 114.114.114.114, 114.114.115.115
    IPv6: 2400:3200::1, 2400:3200:baba::1
    阿里云DNS:
    IPv4: 223.5.5.5, 223.6.6.6
    腾讯DNS:
    IPv4: 119.29.29.29, 119.28.28.28
    百度DNS:
    IPv4: 180.76.76.76
    DNSPod DNS (也称为Dnspod Public DNS):
    IPv4: 119.29.29.29, 182.254.116.116
*/

func GetEthConfig(Interface string) (model.MNetworkConfig, error) {
	MNetworkConfig := model.MNetworkConfig{}
	err := interdb.DB().
		Where("interface=?", Interface).
		Find(&MNetworkConfig).Error
	return MNetworkConfig, err
}

func UpdateEthConfig(MNetworkConfig model.MNetworkConfig) error {
	Model := model.MNetworkConfig{}
	return interdb.DB().
		Model(Model).
		Where("interface=?", MNetworkConfig.Interface).
		Updates(MNetworkConfig).Error
}

/*
*
* 初始化网卡配置参数
*
 */
func InitNetWorkConfig() error {
	dhcp0 := true
	dhcp1 := false
	eth0 := model.MNetworkConfig{
		Type:      "ETH",
		Interface: "eth0",
		Address:   "192.168.1.100",
		Netmask:   "255.255.255.0",
		Gateway:   "192.168.1.1",
		DNS: model.StringList{
			"8.8.8.8",
			"114.114.114.114",
		},
		DHCPEnabled: &dhcp0,
	}
	eth1 := model.MNetworkConfig{
		Type:      "ETH",
		Interface: "eth1",
		Address:   "192.168.64.100",
		Netmask:   "255.255.255.0",
		Gateway:   "192.168.64.1",
		DNS: model.StringList{
			"8.8.8.8",
			"114.114.114.114",
		},
		DHCPEnabled: &dhcp1,
	}
	var err error
	err = interdb.DB().Where("interface=? and id=1", "eth0").FirstOrCreate(&eth0).Error
	if err != nil {
		return err
	}
	err = interdb.DB().Where("interface=? and id=2", "eth1").FirstOrCreate(&eth1).Error
	if err != nil {
		return err
	}
	return nil
}

/*
*
* 匹配: /etc/network/interfaces
*
 */
func GetAllNetConfig() ([]model.MNetworkConfig, error) {
	// 查出前两个网卡的配置
	ethCfg := []model.MNetworkConfig{}
	err := interdb.DB().
		Where("interface=? or interface=?", "eth0", "eth1").
		Find(&ethCfg).Error
	return ethCfg, err
}
