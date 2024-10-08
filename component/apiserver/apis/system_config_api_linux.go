package apis

import (
	"net"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rhilex/archsupport/haas506"
	"github.com/hootrhino/rhilex/archsupport/rhilexg1"
	"github.com/hootrhino/rhilex/archsupport/rhilexpro1"
	common "github.com/hootrhino/rhilex/component/apiserver/common"
	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/apiserver/service"
	"github.com/hootrhino/rhilex/ossupport"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

/*
*
* 设置音量
*
 */
func SetVolume(c *gin.Context, ruleEngine typex.Rhilex) {
	if runtime.GOOS != "linux" {
		c.JSON(common.HTTP_OK, common.Error("OS Not Support:"+runtime.GOOS))
		return
	}
	type Form struct {
		Volume int `json:"volume"`
	}
	DtoCfg := Form{}
	if err0 := c.ShouldBindJSON(&DtoCfg); err0 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err0))
		return
	}
	v, err := ossupport.SetVolume(DtoCfg.Volume)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(v))

}

/*
*
* 获取音量的值
*
 */
func GetVolume(c *gin.Context, ruleEngine typex.Rhilex) {
	if runtime.GOOS != "linux" {
		c.JSON(common.HTTP_OK, common.Error("OS Not Support:"+runtime.GOOS))
		return
	}
	v, err := ossupport.GetVolume()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	if v == "" {
		c.JSON(common.HTTP_OK, common.Error("Volume get failed, please check system"))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(map[string]string{
		"volume": v,
	}))
}

/*
*
* WIFI
*
 */
func GetWifi(c *gin.Context, ruleEngine typex.Rhilex) {
	iface, _ := c.GetQuery("iface")
	MWifiConfig, err := service.GetWlanConfig(iface)
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	Cfg := ossupport.WlanConfig{
		Wlan0: ossupport.WLANInterface{
			Interface: MWifiConfig.Interface,
			SSID:      MWifiConfig.SSID,
			Password:  MWifiConfig.Password,
			Security:  MWifiConfig.Security,
		},
	}
	c.JSON(common.HTTP_OK, common.OkWithData(Cfg))

}

/*
*
*
*通过nmcli配置WIFI
 */
func SetWifi(c *gin.Context, ruleEngine typex.Rhilex) {
	if runtime.GOOS != "linux" {
		c.JSON(common.HTTP_OK, common.Error("OS Not Support:"+runtime.GOOS))
		return
	}
	type Form struct {
		Interface string `json:"interface"`
		SSID      string `json:"ssid"`
		Password  string `json:"password"`
		Security  string `json:"security"` // wpa2-psk wpa3-psk
	}

	DtoCfg := Form{}
	if err0 := c.ShouldBindJSON(&DtoCfg); err0 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err0))
		return
	}
	if !utils.SContains([]string{"wpa2-psk", "wpa3-psk"}, DtoCfg.Security) {
		c.JSON(common.HTTP_OK,
			common.Error(("Only support 2 valid security algorithm:wpa2-psk,wpa3-psk")))
		return
	}
	if !utils.SContains([]string{"wlan0"}, DtoCfg.Interface) {
		c.JSON(common.HTTP_OK, common.Error(("Only support wlan0")))
		return
	}

	MNetCfg := model.MNetworkConfig{
		Type:      "WIFI",
		Interface: DtoCfg.Interface,
		SSID:      DtoCfg.SSID,
		Password:  DtoCfg.Password,
		Security:  DtoCfg.Security,
	}
	if err := service.UpdateWlanConfig(MNetCfg); err != nil {
		if err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
	}
	if typex.DefaultVersionInfo.Product == "RHILEXPRO1" {
		errSetWifi := rhilexpro1.SetWifi(MNetCfg.Interface, MNetCfg.SSID, MNetCfg.Password, 3*time.Second)
		if errSetWifi != nil {
			c.JSON(common.HTTP_OK, common.Error400(errSetWifi))
			return
		}
		goto END
	}
	if typex.DefaultVersionInfo.Product == "HAAS506LD1" {
		errSetWifi := haas506.SetWifi(MNetCfg.Interface, MNetCfg.SSID, MNetCfg.Password, 3*time.Second)
		if errSetWifi != nil {
			c.JSON(common.HTTP_OK, common.Error400(errSetWifi))
			return
		}
		goto END
	}
	if typex.DefaultVersionInfo.Product == "RHILEXG1" {
		errSetWifi := rhilexg1.SetWifi(MNetCfg.Interface, MNetCfg.SSID, MNetCfg.Password, 3*time.Second)
		if errSetWifi != nil {
			c.JSON(common.HTTP_OK, common.Error400(errSetWifi))
			return
		}
		goto END
	}
END:
	c.JSON(common.HTTP_OK, common.Error("Unsupported Product:"+typex.DefaultVersionInfo.Product))

}

/*
*
* 生成最新的ETC配置
*
 */
func ApplyNewestEtcEthConfig() error {
	return nil

}

/*
*
* 时区设置
*
 */
type timeVo struct {
	SysTime     string `json:"sysTime"`
	SysTimeZone string `json:"sysTimeZone"`
	EnableNtp   bool   `json:"enableNtp"`
}

/*
*
  - 设置时间、时区
  - sudo date -s "2023-08-07 15:30:00"
    获取时间: date "+%Y-%m-%d %H:%M:%S" -> 2023-08-07 15:30:00
*/
func SetSystemTime(c *gin.Context, ruleEngine typex.Rhilex) {
	if runtime.GOOS != "linux" {
		c.JSON(common.HTTP_OK, common.Error("OS Not Support:"+runtime.GOOS))
		return
	}
	DtoCfg := timeVo{}
	if err0 := c.ShouldBindJSON(&DtoCfg); err0 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err0))
		return
	}
	if validTimeZone(DtoCfg.SysTimeZone) {
		c.JSON(common.HTTP_OK, common.Error("Invalid TimeZone:"+DtoCfg.SysTimeZone))
		return
	}

	err1 := ossupport.SetSystemTime(DtoCfg.SysTime)
	if err1 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err1))
		return
	}
	err2 := ossupport.SetTimeZone(DtoCfg.SysTimeZone)
	if err2 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err2))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())

}

/*
*
* 获取系统时间
*
 */
func GetSystemTime(c *gin.Context, ruleEngine typex.Rhilex) {
	if runtime.GOOS != "linux" {
		c.JSON(common.HTTP_OK, common.Error("OS Not Support:"+runtime.GOOS))
		return
	}
	SysTime, err := ossupport.GetSystemTime()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	SysTimeZone, err := ossupport.GetTimeZone()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}

	c.JSON(common.HTTP_OK, common.OkWithData(timeVo{
		EnableNtp:   true,
		SysTime:     SysTime,
		SysTimeZone: SysTimeZone.CurrentTimezone,
	}))
}

/*
*
* 设置静态网络IP等, 当前只支持Linux 其他的没测试暂时不予支持

	{
	  "name": "eth0",
	  "interface": "eth0",
	  "address": "192.168.1.100",
	  "netmask": "255.255.255.0",
	  "gateway": "192.168.1.1",
	  "dns": ["8.8.8.8", "8.8.4.4"],
	  "dhcp_enabled": false
	}
*/
func isValidIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil
}

var (
	once     sync.Once
	timeZone *regexp.Regexp
)

func validTimeZone(timezone string) bool {
	once.Do(func() {
		regexPattern := `^[A-Za-z]+/[A-Za-z_]+$`
		timeZone = regexp.MustCompile(regexPattern)
	})

	return timeZone.MatchString(timezone)
}

/*
*
* 展示网络配置信息
*
 */
func GetEthNetwork(c *gin.Context, ruleEngine typex.Rhilex) {
	c.JSON(common.HTTP_OK, common.OkWithData(map[string]ossupport.EtcNetworkConfig{}))

}

/*
*
* 获取当前网络情况
*
 */
func GetCurrentNetConnection(c *gin.Context, ruleEngine typex.Rhilex) {
	nmcliOutput, err := ossupport.GetCurrentNetConnection()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(nmcliOutput))
}

/*
*
* 设置两个网口
*
 */
func SetEthNetwork(c *gin.Context, ruleEngine typex.Rhilex) {
	if runtime.GOOS != "linux" {
		c.JSON(common.HTTP_OK, common.Error("OS Not Support:"+runtime.GOOS))
		return
	}
	type Form struct {
		Interface   string   `json:"interface"` // eth1 eth0
		Address     string   `json:"address"`
		Netmask     string   `json:"netmask"`
		Gateway     string   `json:"gateway"`
		DNS         []string `json:"dns"`
		DHCPEnabled bool     `json:"dhcp_enabled"`
	}
	DtoCfg := Form{}
	if err0 := c.ShouldBindJSON(&DtoCfg); err0 != nil {
		c.JSON(common.HTTP_OK, common.Error400(err0))
		return
	}

	if !isValidIP(DtoCfg.Address) {
		c.JSON(common.HTTP_OK,
			common.Error(("Invalid IP:" + DtoCfg.Address)))
		return
	}
	if !isValidIP(DtoCfg.Gateway) {
		c.JSON(common.HTTP_OK,
			common.Error(("Invalid Gateway IP:" + DtoCfg.Address)))
		return
	}
	if !isValidSubnetMask(DtoCfg.Netmask) {
		c.JSON(common.HTTP_OK,
			common.Error(("Invalid SubnetMask:" + DtoCfg.Address)))
		return
	}
	for _, dns := range DtoCfg.DNS {
		if !isValidIP(dns) {
			c.JSON(common.HTTP_OK,
				common.Error(("Invalid DNS IP:" + DtoCfg.Address)))
			return
		}
	}

	MNetCfg := model.MNetworkConfig{
		Type:        "ETH",
		Interface:   DtoCfg.Interface,
		Address:     DtoCfg.Address,
		Netmask:     DtoCfg.Netmask,
		Gateway:     DtoCfg.Gateway,
		DNS:         DtoCfg.DNS,
		DHCPEnabled: &DtoCfg.DHCPEnabled,
	}
	if err := service.UpdateEthConfig(MNetCfg); err != nil {
		if err != nil {
			c.JSON(common.HTTP_OK, common.Error400(err))
			return
		}
	}
	if typex.DefaultVersionInfo.Product == "RHILEXPRO1" {
		config := []rhilexpro1.NetworkInterfaceConfig{
			rhilexpro1.NetworkInterfaceConfig{
				Interface:   MNetCfg.Interface,
				Address:     MNetCfg.Address,
				Netmask:     MNetCfg.Netmask,
				Gateway:     MNetCfg.Gateway,
				DHCPEnabled: *MNetCfg.DHCPEnabled,
			},
		}
		errSetWifi := rhilexpro1.SetEthernet(config)
		if errSetWifi != nil {
			c.JSON(common.HTTP_OK, common.Error400(errSetWifi))
			return
		}
		goto END
	}
	if typex.DefaultVersionInfo.Product == "HAAS506LD1" {
		config := []haas506.NetworkInterfaceConfig{
			haas506.NetworkInterfaceConfig{
				Interface:   MNetCfg.Interface,
				Address:     MNetCfg.Address,
				Netmask:     MNetCfg.Netmask,
				Gateway:     MNetCfg.Gateway,
				DHCPEnabled: *MNetCfg.DHCPEnabled,
			},
		}
		errSetWifi := haas506.SetEthernet(config)
		if errSetWifi != nil {
			c.JSON(common.HTTP_OK, common.Error400(errSetWifi))
			return
		}
		goto END
	}
	if typex.DefaultVersionInfo.Product == "RHILEXG1" {
		config := []rhilexg1.NetworkInterfaceConfig{
			rhilexg1.NetworkInterfaceConfig{
				Interface:   MNetCfg.Interface,
				Address:     MNetCfg.Address,
				Netmask:     MNetCfg.Netmask,
				Gateway:     MNetCfg.Gateway,
				DHCPEnabled: *MNetCfg.DHCPEnabled,
			},
		}
		errSetWifi := rhilexg1.SetEthernet(config)
		if errSetWifi != nil {
			c.JSON(common.HTTP_OK, common.Error400(errSetWifi))
			return
		}
		goto END
	}
END:
	c.JSON(common.HTTP_OK, common.Ok())

}

/*
*
* 更新时间
*
 */
func UpdateTimeByNtp(c *gin.Context, ruleEngine typex.Rhilex) {
	if err := ossupport.UpdateTimeByNtp(); err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.Ok())

}

func isValidSubnetMask(mask string) bool {
	// 分割子网掩码为4个整数
	parts := strings.Split(mask, ".")
	if len(parts) != 4 {
		return false
	}

	// 将每个部分转换为整数
	var octets [4]int
	for i, part := range parts {
		octet, err := strconv.Atoi(part)
		if err != nil || octet < 0 || octet > 255 {
			return false
		}
		octets[i] = octet
	}

	// 判断是否为有效的子网掩码
	var bits int
	for _, octet := range octets {
		bits += bitsInByte(octet)
	}

	return bits >= 1 && bits <= 32
}

func bitsInByte(b int) int {
	count := 0
	for b > 0 {
		count += b & 1
		b >>= 1
	}
	return count
}
