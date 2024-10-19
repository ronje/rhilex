package httpserver

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hootrhino/rhilex/component/crontask"
	dataschema "github.com/hootrhino/rhilex/component/dataschema"
	"github.com/hootrhino/rhilex/component/internotify"
	"github.com/shirou/gopsutil/cpu"

	"github.com/hootrhino/rhilex/component/apiserver/apis"
	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/apiserver/server"
	"github.com/hootrhino/rhilex/component/apiserver/service"
	"github.com/hootrhino/rhilex/component/applet"
	"github.com/hootrhino/rhilex/component/interdb"
	"github.com/hootrhino/rhilex/component/trailer"

	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"

	"gopkg.in/ini.v1"

	_ "github.com/mattn/go-sqlite3"
)

type _serverConfig struct {
	DbPath string `ini:"dbpath"`
	Port   int    `ini:"port"`
}
type ApiServerPlugin struct {
	uuid       string
	ruleEngine typex.Rhilex
	mainConfig _serverConfig
}

func NewHttpApiServer(ruleEngine typex.Rhilex) *ApiServerPlugin {
	return &ApiServerPlugin{
		uuid:       "HTTP-API-SERVER",
		mainConfig: _serverConfig{Port: 2580},
		ruleEngine: ruleEngine,
	}
}

/*
*
* 初始化RHILEX, 初始化数据到运行时
*
 */
func initRhilex(engine typex.Rhilex) {
	go GetCpuUsage()
	//
	// Load inend from sqlite
	//
	for _, minEnd := range service.AllMInEnd() {
		if err := server.LoadNewestInEnd(minEnd.UUID, engine); err != nil {
			glogger.GLogger.Error("InEnd load failed:", err)
		}
	}

	//
	// Load out from sqlite
	//
	for _, mOutEnd := range service.AllMOutEnd() {
		if err := server.LoadNewestOutEnd(mOutEnd.UUID, engine); err != nil {
			glogger.GLogger.Error("OutEnd load failed:", err)
		}
	}
	// 加载设备
	for _, mDevice := range service.AllDevices() {
		glogger.GLogger.Debug("LoadNewestDevice mDevice.BindRules: ", mDevice.BindRules.String())
		if err := server.LoadNewestDevice(mDevice.UUID, engine); err != nil {
			glogger.GLogger.Error("Device load failed:", err)
		}

	}
	// 加载外挂
	for _, mGoods := range service.AllGoods() {
		newGoods := trailer.GoodsInfo{
			UUID:        mGoods.UUID,
			AutoStart:   mGoods.AutoStart,
			LocalPath:   mGoods.LocalPath,
			NetAddr:     mGoods.NetAddr,
			Args:        mGoods.Args,
			ExecuteType: mGoods.ExecuteType,
			Description: mGoods.Description,
		}
		if err := trailer.StartProcess(newGoods); err != nil {
			glogger.GLogger.Error("Goods load failed:", err)
		}
	}
	//
	// APP stack
	//
	for _, mApp := range service.AllApp() {
		app := applet.NewApplication(
			mApp.UUID,
			mApp.Name,
			mApp.Version,
		)
		if err := applet.LoadApp(app, mApp.LuaSource); err != nil {
			glogger.GLogger.Error(err)
			continue
		}
		if *mApp.AutoStart {
			glogger.GLogger.Debug("App autoStart allowed:", app.UUID)
			if err1 := applet.StartApp(app.UUID); err1 != nil {
				glogger.GLogger.Error("App autoStart failed:", err1)
			}
		}
	}

}

func (hs *ApiServerPlugin) Init(config *ini.Section) error {
	if err := utils.InIMapToStruct(config, &hs.mainConfig); err != nil {
		return err
	}
	server.StartRhilexApiServer(hs.ruleEngine, hs.mainConfig.Port)
	interdb.DB().Exec("VACUUM;")
	interdb.RegisterModel(
		&model.MInEnd{},
		&model.MOutEnd{},
		&model.MRule{},
		&model.MUser{},
		&model.MDevice{},
		&model.MGoods{},
		&model.MApplet{},
		&model.MGenericGroup{},
		&model.MGenericGroupRelation{},
		&model.MNetworkConfig{},
		&model.MIotSchema{},
		&model.MIotProperty{},
		&model.MIpRoute{},
		&model.MUart{},
		&model.MInternalNotify{},
		&model.MUserLuaTemplate{},
		&model.MModbusDataPoint{},
		&model.MSiemensDataPoint{},
		&model.MSnmpOid{},
		&model.MCjt1882004DataPoint{},
		&model.MDlt6452007DataPoint{},
		&model.MSzy2062016DataPoint{},
		&model.MUserProtocolDataPoint{},
		&model.MBacnetDataPoint{},
		&model.MBacnetRouterDataPoint{},
		&model.MMBusDataPoint{},
		&model.MDataPoint{},
		&model.MCronRebootConfig{},
	)
	// 初始化所有预制参数
	server.DefaultApiServer.InitializeProduct()
	server.DefaultApiServer.InitializeGenericOSData()
	server.DefaultApiServer.InitializeWindowsData()
	server.DefaultApiServer.InitializeUnixData()
	// InitDataSchemaCache
	dataschema.InitDataSchemaCache(hs.ruleEngine)
	// Cron Reboot Executor
	crontask.InitCronRebootExecutor(hs.ruleEngine)
	initRhilex(hs.ruleEngine)
	return nil
}

/*
*
* 加载路由
*
 */
func (hs *ApiServerPlugin) LoadRoute() {
	// User
	apis.InitUserRoute()
	// In End
	apis.InitInEndRoute()
	// Rules
	apis.InitRulesRoute()
	// Out End
	apis.InitOutEndRoute()
	// System API
	apis.InitSystemRoute()
	// backup
	apis.InitBackupRoute()
	// 设备管理
	apis.InitDeviceRoute()
	// Modbus Slaver
	apis.InitModbusSlaverRoute()
	// S1200 点位表
	apis.InitSiemensS7Route()
	// applet
	apis.InitAppletRoute()
	// plugins
	apis.InitPluginsRoute()
	// 分组管理
	apis.InitGroupRoute()
	// 用户LUA代码段管理
	apis.InitUserLuaRoute()
	// System Permission
	apis.InitSysMenuPermissionRoute()
	// System Settings
	apis.LoadSystemSettingsAPI()
	// Modbus
	apis.InitModbusRoute()
	// Mbus
	apis.InitMBusRoute()
	// DLT645
	apis.InitDlt6452007Route()
	// USer Protocol
	apis.InitUserProtocolRoute()
	// Init Internal Notify Route
	apis.InitInternalNotifyRoute()
	// Snmp Route
	apis.InitSnmpRoute()
	// Bacnet Route
	apis.InitBacnetIpRoute()
	// Bacnet Router
	apis.InitBacnetRouterRoute()
	// Data Schema
	apis.InitDataSchemaApi()
	// Data Center
	apis.InitDataCenterApi()
	// Transceiver
	apis.InitTransceiverRoute()
	// Data Point Route
	apis.InitDataPointRoute()
	// Mqtt Server
	apis.InitMqttSourceServerRoute()
	// Cron Reboot
	apis.InitCronRebootRoute()
}

// ApiServerPlugin Start
func (hs *ApiServerPlugin) Start(r typex.Rhilex) error {
	hs.ruleEngine = r
	hs.LoadRoute()
	glogger.GLogger.Infof("Http server started on :%v", hs.mainConfig.Port)
	return nil
}

func (hs *ApiServerPlugin) Stop() error {
	dataschema.FlushDataSchemaCache()
	return nil
}

func (hs *ApiServerPlugin) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		UUID:        hs.uuid,
		Name:        "RHILEX HTTP RESTFul Api Server",
		Version:     "v1.0.0",
		Description: "RHILEX HTTP RESTFul Api Server",
	}
}

/*
*
* 服务调用接口
*
 */
func (*ApiServerPlugin) Service(arg typex.ServiceArg) typex.ServiceResult {
	return typex.ServiceResult{Out: "ApiServerPlugin"}
}
func GetCpuUsage() {
	for {
		select {
		case <-typex.GCTX.Done():
			{
				return
			}
		default:
			{
			}
		}
		cpuPercent, _ := cpu.Percent(time.Duration(10)*time.Second, true)
		V := calculateCpuPercent(cpuPercent)
		// TODO 这个比例需要通过参数适配
		if V > 90 {
			internotify.Push(internotify.BaseEvent{
				Type:    `WARNING`,
				Event:   `system.cpu.load`,
				Ts:      uint64(time.Now().UnixMilli()),
				Summary: "High CPU Usage",
				Info:    fmt.Sprintf("High CPU Usage: %.2f%%, please maintain the device", V),
			})
		}
	}

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
