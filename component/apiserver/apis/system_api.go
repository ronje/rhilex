package apis

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"

	common "github.com/hootrhino/rhilex/component/apiserver/common"
	"github.com/hootrhino/rhilex/component/apiserver/server"
	"github.com/hootrhino/rhilex/component/apiserver/service"
	"github.com/hootrhino/rhilex/component/applet"
	"github.com/hootrhino/rhilex/component/intermetric"
	"github.com/hootrhino/rhilex/component/security"
	core "github.com/hootrhino/rhilex/config"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/ossupport"
	"github.com/hootrhino/rhilex/periphery"
	"github.com/hootrhino/rhilex/periphery/en6400"
	"github.com/hootrhino/rhilex/periphery/haas506"
	"github.com/hootrhino/rhilex/periphery/rhilexg1"
	"github.com/hootrhino/rhilex/periphery/rhilexpro1"
	"github.com/hootrhino/rhilex/registry"
	"github.com/hootrhino/rhilex/utils"

	"github.com/gin-gonic/gin"
	"github.com/hootrhino/rhilex/typex"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
)

func InitSystemRoute() {
	osApi := server.RouteGroup(server.ContextUrl("/os"))
	{
		osApi.GET(("/netInterfaces"), server.AddRoute(GetNetInterfaces))
		osApi.GET(("/osRelease"), server.AddRoute(CatOsRelease))
		osApi.GET(("/system"), server.AddRoute(System))
		osApi.GET(("/uarts"), server.AddRoute(GetUartList))
		osApi.GET(("/startedAt"), server.AddRoute(StartedAt))
		osApi.GET(("/getVideos"), server.AddRoute(GetVideos))
		osApi.GET(("/getGpuInfo"), server.AddRoute(GetGpuInfo))
		osApi.GET(("/sysConfig"), server.AddRoute(GetSysConfig))
		osApi.POST(("/resetInterMetric"), server.AddRoute(ResetInterMetric))
		osApi.GET(("/getSecurityLicense"), server.AddRoute(GetSecurityLicense))
	}
	systemApi := server.RouteGroup(server.ContextUrl("/"))
	{
		systemApi.GET(("/ping"), server.AddRoute(Ping))
	}
	server.DefaultApiServer.Route().
		GET(server.ContextUrl("statistics"), server.AddRoute(Statistics))
	server.DefaultApiServer.Route().
		POST(server.ContextUrl("login"), server.AddRoute(Login))
	server.DefaultApiServer.Route().
		POST(server.ContextUrl("validateRule"), server.AddRoute(ValidateLuaSyntax))
}

// 启动时间
var __StartedAt = time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")

/*
*
* 健康检查接口, 一般用来监视是否工作
*
 */
func Ping(c *gin.Context, ruleEngine typex.Rhilex) {
	c.JSON(common.HTTP_OK, common.OkWithData("PONG"))
}

/**
 * 前端证书
 *
 */
//
type SecurityLiscnse struct {
	PublicKey string `json:"public_key"`
}

func GetSecurityLicense(c *gin.Context, ruleEngine typex.Rhilex) {
	if !security.CheckLicense() {
		security.InitSecurityLicense()
	}
	_, pub, err := security.ReadLocalKeys()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
	} else {
		c.JSON(common.HTTP_OK, common.OkWithData(SecurityLiscnse{
			PublicKey: pub,
		}))
	}
}

// Get all plugins
func Plugins(c *gin.Context, ruleEngine typex.Rhilex) {
	data := []interface{}{}
	plugins := registry.DefaultPluginRegistry.All()
	for _, plugin := range plugins {
		data = append(data, plugin.PluginMetaInfo())
	}
	c.JSON(common.HTTP_OK, common.OkWithData(data))
}

/*
*
* 本地网卡
*
 */
func GetNetInterfaces(c *gin.Context, ruleEngine typex.Rhilex) {
	interfaces, err := ossupport.GetAvailableInterfaces()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
	} else {
		c.JSON(common.HTTP_OK, common.OkWithData(interfaces))
	}
}

// 计算资源数据
func source_count(e typex.Rhilex) map[string]int {
	return map[string]int{
		"inends":  len(e.AllInEnds()),
		"outends": len(e.AllOutEnds()),
		"rules":   len(e.AllRules()),
		"devices": len(e.AllDevices()),
		"goods":   0,
		"plugins": registry.DefaultPluginRegistry.Count(),
		"apps":    applet.AppCount(),
	}
}

/*
*
* 获取系统指标, Go 自带这个不准, 后期版本需要更换跨平台实现
*
 */
func System(c *gin.Context, ruleEngine typex.Rhilex) {
	cpuPercent, _ := cpu.Percent(time.Duration(1)*time.Second, true)
	diskInfo, _ := disk.Usage("/")
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	// var m runtime.MemStats
	// runtime.ReadMemStats(&m)
	// ip, err0 := utils.HostNameI()
	memPercent, _ := service.GetMemPercent()
	hardWareInfo := map[string]interface{}{
		"version":    typex.MainVersion,
		"diskInfo":   calculateDiskInfo(diskInfo),
		"memPercent": memPercent,
		"cpuPercent": calculateCpuPercent(cpuPercent),
		"osArch":     ruleEngine.Version().Arch,
		"osDist":     ruleEngine.Version().Dist,
		"product":    typex.DefaultVersionInfo.Product,
		"startedAt":  __StartedAt,
		"osUpTime": func() string {
			result, err := ossupport.GetUptime()
			if err != nil {
				return "0 days 0 hours 0 minutes"
			}
			return result
		}(),
	}
	c.JSON(common.HTTP_OK, common.OkWithData(gin.H{
		"hardWareInfo": hardWareInfo,
		"statistic":    intermetric.GetMetric(),
		"sourceCount":  source_count(ruleEngine),
	}))
}

/*
*
* SnapshotDump
*
 */
func SnapshotDump(c *gin.Context, ruleEngine typex.Rhilex) {
	c.Header("Content-Type", "text/plain")
	c.Header("Content-Disposition",
		fmt.Sprintf("attachment;filename=SnapshotDump_%v.json", time.Now().UnixMilli()))
	c.Writer.Write([]byte(ruleEngine.SnapshotDump()))
	c.Writer.Flush()
}

// Get statistics data
func Statistics(c *gin.Context, ruleEngine typex.Rhilex) {
	c.JSON(common.HTTP_OK, common.OkWithData(intermetric.GetMetric()))
}

// Get statistics data
func SourceCount(c *gin.Context, ruleEngine typex.Rhilex) {
	c.JSON(common.HTTP_OK, common.OkWithData(map[string]int{
		"inends":  len(ruleEngine.AllInEnds()),
		"outends": len(ruleEngine.AllOutEnds()),
		"rules":   len(ruleEngine.AllRules()),
		"devices": len(ruleEngine.AllDevices()),
		"plugins": registry.DefaultPluginRegistry.Count(),
	}))
}

/*
*
* 获取本地的串口列表
*
 */
func GetUartList(c *gin.Context, ruleEngine typex.Rhilex) {
	c.JSON(common.HTTP_OK, common.OkWithData(service.GetOsPort()))
}

/*
*
* 计算开机时间
*
 */
func StartedAt(c *gin.Context, ruleEngine typex.Rhilex) {
	c.JSON(common.HTTP_OK, common.OkWithData(__StartedAt))
}

func calculateDiskInfo(diskInfo *disk.UsageStat) float64 {
	value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", diskInfo.UsedPercent), 64)
	return value

}

// 计算CPU平均使用率
func calculateCpuPercent(cpus []float64) float64 {
	var acc float64 = 0
	for _, v := range cpus {
		acc += v
	}
	value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", acc/float64(len(cpus))), 64)
	return value
}

func CatOsRelease(c *gin.Context, ruleEngine typex.Rhilex) {
	r, err := utils.CatOsRelease()
	if err != nil {
		c.JSON(common.HTTP_OK, common.Error400(err))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(r))
}

/*
*
* 重置度量值
*
 */
func ResetInterMetric(c *gin.Context, ruleEngine typex.Rhilex) {
	intermetric.Reset()
	c.JSON(common.HTTP_OK, common.Ok())
}

/*
*
* 获取视频接口
*
 */
func GetVideos(c *gin.Context, ruleEngine typex.Rhilex) {
	if runtime.GOOS == "windows" {
		L, _ := ossupport.GetWindowsVideos()
		c.JSON(common.HTTP_OK, common.OkWithData(L))
	} else {
		L, _ := ossupport.GetUnixVideos()
		c.JSON(common.HTTP_OK, common.OkWithData(L))
	}
}

/*
*
* 获取GPU信息
*
 */
func GetGpuInfo(c *gin.Context, ruleEngine typex.Rhilex) {
	GpuInfos, err := periphery.GetGpuInfoWithNvidiaSmi()
	if err != nil {
		glogger.GLogger.Error(err)
	}
	c.JSON(common.HTTP_OK, common.OkWithData(GpuInfos))
}

/*
*
* 获取控制面板上支持的设备
*
 */
func GetDeviceCtrlTree(c *gin.Context, ruleEngine typex.Rhilex) {
	env := os.Getenv("ARCHSUPPORT")
	if env == "RHILEXG1" {
		c.JSON(common.HTTP_OK, common.OkWithData(rhilexg1.GetSysDevTree()))
		return
	}
	if env == "RHILEXPRO1" {
		c.JSON(common.HTTP_OK, common.OkWithData(rhilexpro1.GetSysDevTree()))
		return
	}
	if env == "EN6400" {
		c.JSON(common.HTTP_OK, common.OkWithData(en6400.GetSysDevTree()))
		return
	}

	if env == "HAAS506LD1" {
		c.JSON(common.HTTP_OK, common.OkWithData(haas506.GetSysDevTree()))
		return
	}
	c.JSON(common.HTTP_OK, common.OkWithData(periphery.DefaultDeviceTree()))
}

/**
 * 系统配置
 *
 */
func GetSysConfig(c *gin.Context, ruleEngine typex.Rhilex) {
	c.JSON(common.HTTP_OK, common.OkWithData(core.GlobalConfig))
}
