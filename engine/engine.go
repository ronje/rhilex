// Copyright (C) 2023 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package engine

import (
	"encoding/json"
	"fmt"
	"runtime"
	"sync"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rhilex/component/aibase"
	"github.com/hootrhino/rhilex/component/appstack"
	"github.com/hootrhino/rhilex/component/hwportmanager"
	"github.com/hootrhino/rhilex/component/interkv"
	"github.com/hootrhino/rhilex/component/rhilexmanager"
	"github.com/hootrhino/rhilex/component/ruleengine"
	transceiver "github.com/hootrhino/rhilex/component/transceivercom/transceiver"
	core "github.com/hootrhino/rhilex/config"

	intercache "github.com/hootrhino/rhilex/component/intercache"

	"github.com/hootrhino/rhilex/component/shellymanager"
	supervisor "github.com/hootrhino/rhilex/component/supervisor"

	datacenter "github.com/hootrhino/rhilex/component/datacenter"
	"github.com/hootrhino/rhilex/component/interdb"
	"github.com/hootrhino/rhilex/component/intermetric"
	"github.com/hootrhino/rhilex/component/internotify"
	"github.com/hootrhino/rhilex/component/interqueue"
	"github.com/hootrhino/rhilex/component/trailer"
	"github.com/hootrhino/rhilex/device"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/source"
	"github.com/hootrhino/rhilex/target"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/sirupsen/logrus"
)

/*
*
* 全局默认引擎，未来主要留给外部使用
*
 */
var __DefaultRuleEngine *RuleEngine

const __DEFAULT_DB_PATH string = "./rhilex.db"

// 规则引擎
type RuleEngine struct {
	locker            sync.Mutex
	Rules             *sync.Map                        `json:"rules"`
	Plugins           *sync.Map                        `json:"plugins"`
	InEnds            *sync.Map                        `json:"inends"`
	OutEnds           *sync.Map                        `json:"outends"`
	Devices           *sync.Map                        `json:"devices"`
	Config            *typex.RhilexConfig              `json:"config"`
	DeviceTypeManager *rhilexmanager.DeviceTypeManager `json:"-"` // 待迁移组件: component/rhilexmanager
	SourceTypeManager *rhilexmanager.SourceTypeManager `json:"-"` // 待迁移组件: component/rhilexmanager
	TargetTypeManager *rhilexmanager.TargetTypeManager `json:"-"` // 待迁移组件: component/rhilexmanager
}

func InitRuleEngine(config typex.RhilexConfig) typex.Rhilex {
	__DefaultRuleEngine = &RuleEngine{
		locker:            sync.Mutex{},
		DeviceTypeManager: rhilexmanager.NewDeviceTypeManager(),
		SourceTypeManager: rhilexmanager.NewSourceTypeManager(),
		TargetTypeManager: rhilexmanager.NewTargetTypeManager(),
		Plugins:           &sync.Map{},
		Rules:             &sync.Map{},
		InEnds:            &sync.Map{},
		OutEnds:           &sync.Map{},
		Devices:           &sync.Map{},
		Config:            &config,
	}

	// Internal DB
	interdb.Init(__DefaultRuleEngine, __DEFAULT_DB_PATH)
	// Internal kv Store
	interkv.InitInterKVStore(core.GlobalConfig.MaxKvStoreSize)
	// Shelly Device Registry
	shellymanager.InitShellyDeviceRegistry(__DefaultRuleEngine)
	// SuperVisor Admin
	supervisor.InitResourceSuperVisorAdmin(__DefaultRuleEngine)
	// Init Global Value Registry
	intercache.InitGlobalValueRegistry(__DefaultRuleEngine)
	// Internal Bus
	internotify.InitInternalEventBus(__DefaultRuleEngine, core.GlobalConfig.MaxQueueSize)
	// Load hardware Port Manager
	hwportmanager.InitHwPortsManager(__DefaultRuleEngine)
	// Internal Metric
	intermetric.InitInternalMetric(__DefaultRuleEngine)
	// trailer
	trailer.InitTrailerRuntime(__DefaultRuleEngine)
	// lua appstack manager
	appstack.InitAppStack(__DefaultRuleEngine)
	// current only support Internal ai
	aibase.InitAlgorithmRuntime(__DefaultRuleEngine)
	// Data center: future version maybe support
	datacenter.InitDataCenter(__DefaultRuleEngine)
	// Internal Queue
	interqueue.InitYQueue(__DefaultRuleEngine, core.GlobalConfig.MaxQueueSize)
	// Internal BUS
	interqueue.StartYQueue()
	// Init Transceiver Communicator Manager
	transceiver.InitTransceiverCommunicatorManager(__DefaultRuleEngine)
	return __DefaultRuleEngine
}

/*
*
* Engine Start
*
 */
func (e *RuleEngine) Start() *typex.RhilexConfig {
	// Resource Manager
	e.InitDeviceTypeManager()
	e.InitSourceTypeManager()
	e.InitTargetTypeManager()
	intercache.RegisterSlot("__DefaultRuleEngine")
	return e.Config
}

func (e *RuleEngine) GetPlugins() *sync.Map {
	return e.Plugins
}
func (e *RuleEngine) AllPlugins() *sync.Map {
	return e.Plugins
}

func (e *RuleEngine) Version() typex.VersionInfo {
	return typex.DefaultVersionInfo
}

func (e *RuleEngine) GetConfig() *typex.RhilexConfig {
	return e.Config
}

// Stop
func (e *RuleEngine) Stop() {
	glogger.GLogger.Info("[*] Ready to stop rhilex")
	// 所有的APP停了
	appstack.Stop()
	// 外挂停了
	trailer.Stop()
	// 资源
	e.InEnds.Range(func(key, value interface{}) bool {
		inEnd := value.(*typex.InEnd)
		if inEnd.Source != nil {
			glogger.GLogger.Infof("Stop InEnd:(%s,%s)", inEnd.Name, inEnd.UUID)
			e.GetInEnd(inEnd.UUID).State = typex.SOURCE_STOP
			inEnd.Source.Stop()
		}
		glogger.GLogger.Infof("Stop InEnd:(%s,%s) Successfully", inEnd.Name, inEnd.UUID)
		return true
	})
	// 停止所有外部资源
	e.OutEnds.Range(func(key, value interface{}) bool {
		outEnd := value.(*typex.OutEnd)
		if outEnd.Target != nil {
			glogger.GLogger.Infof("Stop NewTarget:(%s,%s)", outEnd.Name, outEnd.UUID)
			e.GetOutEnd(outEnd.UUID).State = typex.SOURCE_STOP
			outEnd.Target.Stop()
			glogger.GLogger.Infof("Stop NewTarget:(%s,%s) Successfully", outEnd.Name, outEnd.UUID)
		}
		return true
	})
	// 停止所有插件
	e.Plugins.Range(func(key, value interface{}) bool {
		plugin := value.(typex.XPlugin)
		glogger.GLogger.Infof("Stop plugin:(%s)", plugin.PluginMetaInfo().Name)
		plugin.Stop()
		glogger.GLogger.Infof("Stop plugin:(%s) Successfully", plugin.PluginMetaInfo().Name)
		return true
	})
	// 停止所有设备
	e.Devices.Range(func(key, value interface{}) bool {
		Device := value.(*typex.Device)
		glogger.GLogger.Infof("Stop Device:(%s)", Device.Name)
		e.GetDevice(Device.UUID).State = typex.DEV_STOP
		Device.Device.Stop()
		glogger.GLogger.Infof("Stop Device:(%s) Successfully", Device.Name)
		return true
	})
	// Flush Shelly Device Cache
	glogger.GLogger.Info("Flush Shelly Device Cache")
	shellymanager.Flush()
	// Internal Cache
	glogger.GLogger.Info("Flush Internal Cache")
	intercache.Flush()
	// AI Runtime
	glogger.GLogger.Info("Stop AI Runtime")
	aibase.Stop()
	// Stop transceiver
	glogger.GLogger.Info("Stop transceiver")
	transceiver.Stop()
	// UnRegister __DefaultRuleEngine
	intercache.UnRegisterSlot("__DefaultRuleEngine")
	// UnRegister __DeviceConfigMap
	intercache.UnRegisterSlot("__DeviceConfigMap")
	// END
	glogger.GLogger.Info("[√] Stop rhilex successfully")
	if err := glogger.Close(); err != nil {
		fmt.Println("Close logger error: ", err)
	}
}

// 核心功能: Work, 主要就是推流进队列
func (e *RuleEngine) WorkInEnd(in *typex.InEnd, data string) (bool, error) {
	if err := interqueue.DefaultYQueue.PushInQueue(in, data); err != nil {
		return false, err
	}
	return true, nil
}

// 核心功能: Work, 主要就是推流进队列
func (e *RuleEngine) WorkDevice(Device *typex.Device, data string) (bool, error) {
	if err := interqueue.DefaultYQueue.PushDeviceQueue(Device, data); err != nil {
		return false, err
	}
	return true, nil
}

/*
*
* 执行针对资源端的规则脚本
*
 */
func (e *RuleEngine) RunSourceCallbacks(in *typex.InEnd, callbackArgs string) {
	// 执行来自资源的脚本
	for _, rule := range in.BindRules {
		if rule.Status == typex.RULE_RUNNING {
			_, errA := ruleengine.ExecuteActions(&rule, lua.LString(callbackArgs))
			if errA != nil {
				Debugger, Ok := rule.LuaVM.GetStack(1)
				if Ok {
					LValue, _ := rule.LuaVM.GetInfo("f", Debugger, lua.LNil)
					rule.LuaVM.GetInfo("l", Debugger, lua.LNil)
					rule.LuaVM.GetInfo("S", Debugger, lua.LNil)
					rule.LuaVM.GetInfo("u", Debugger, lua.LNil)
					rule.LuaVM.GetInfo("n", Debugger, lua.LNil)
					LFunction := LValue.(*lua.LFunction)
					LastCall := lua.DbgCall{
						Name: "_main", Pc: 0,
					}
					if len(LFunction.Proto.DbgCalls) > 0 {
						LastCall = LFunction.Proto.DbgCalls[0]
					}
					glogger.GLogger.WithFields(logrus.Fields{
						"topic": "rule/log/" + rule.UUID,
					}).Warnf("Function Name: [%s],"+
						"What: [%s], Source Line: [%d],"+
						" Last Call: [%s], Error message: %s",
						Debugger.Name, Debugger.What, Debugger.CurrentLine,
						LastCall.Name, errA.Error(),
					)
				}
				// _, err0 := ruleengine.ExecuteFailed(rule.LuaVM, lua.LString(errA.Error()))
				// if err0 != nil {
				// 	glogger.GLogger.Error(err0)
				// }
			} else {
				_, errS := ruleengine.ExecuteSuccess(rule.LuaVM)
				if errS != nil {
					glogger.GLogger.Error(errS)
					return // lua 是规则链，有短路原则，中途出错会中断
				}
			}
		}
	}
}

/*
*
* 执行针对设备端的规则脚本
*
 */
func (e *RuleEngine) RunDeviceCallbacks(Device *typex.Device, callbackArgs string) {
	for _, rule := range Device.BindRules {
		_, errA := ruleengine.ExecuteActions(&rule, lua.LString(callbackArgs))
		if errA != nil {
			Debugger, Ok := rule.LuaVM.GetStack(1)
			if Ok {
				LValue, _ := rule.LuaVM.GetInfo("f", Debugger, lua.LNil)
				rule.LuaVM.GetInfo("l", Debugger, lua.LNil)
				rule.LuaVM.GetInfo("S", Debugger, lua.LNil)
				rule.LuaVM.GetInfo("u", Debugger, lua.LNil)
				rule.LuaVM.GetInfo("n", Debugger, lua.LNil)
				LFunction := LValue.(*lua.LFunction)
				LastCall := lua.DbgCall{
					Name: "_main", Pc: 0,
				}
				if len(LFunction.Proto.DbgCalls) > 0 {
					LastCall = LFunction.Proto.DbgCalls[0]
				}
				glogger.GLogger.WithFields(logrus.Fields{
					"topic": "rule/log/" + rule.UUID,
				}).Warnf("Function Name: [%s],"+
					"What: [%s], Source Line: [%d],"+
					" Last Call: [%s], Error message: %s",
					Debugger.Name, Debugger.What, Debugger.CurrentLine,
					LastCall.Name, errA.Error(),
				)
			}
			// _, err1 := ruleengine.ExecuteFailed(rule.LuaVM, lua.LString(errA.Error()))
			// if err1 != nil {
			// 	glogger.GLogger.Error(err1)
			// }
		} else {
			_, err2 := ruleengine.ExecuteSuccess(rule.LuaVM)
			if err2 != nil {
				glogger.GLogger.WithFields(logrus.Fields{
					"topic": "rule/log/" + rule.UUID,
				}).Info("RunLuaCallbacks error:", err2)
				return
			}
		}
	}
}

func (e *RuleEngine) GetInEnd(uuid string) *typex.InEnd {
	v, ok := (e.InEnds).Load(uuid)
	if ok {
		return v.(*typex.InEnd)
	}
	return nil
}

func (e *RuleEngine) SaveInEnd(in *typex.InEnd) {
	e.InEnds.Store(in.UUID, in)
}

func (e *RuleEngine) RemoveInEnd(uuid string) {
	if inEnd := e.GetInEnd(uuid); inEnd != nil {
		if inEnd.Source != nil {
			glogger.GLogger.Infof("InEnd [%s, %s] ready to stop", uuid, inEnd.Name)
			inEnd.Source.Stop()
			glogger.GLogger.Infof("InEnd [%s, %s] stopped", uuid, inEnd.Name)
			e.InEnds.Delete(uuid)
			glogger.GLogger.Infof("InEnd [%s, %s] has been deleted", uuid, inEnd.Name)
			inEnd = nil
		}
	}
}

func (e *RuleEngine) AllInEnds() *sync.Map {
	return e.InEnds
}

func (e *RuleEngine) GetOutEnd(id string) *typex.OutEnd {
	v, ok := e.OutEnds.Load(id)
	if ok {
		return v.(*typex.OutEnd)
	} else {
		return nil
	}

}

func (e *RuleEngine) SaveOutEnd(out *typex.OutEnd) {
	e.OutEnds.Store(out.UUID, out)

}

func (e *RuleEngine) RemoveOutEnd(uuid string) {
	if outEnd := e.GetOutEnd(uuid); outEnd != nil {
		if outEnd.Target != nil {
			glogger.GLogger.Infof("OutEnd [%s, %s] ready to stop", uuid, outEnd.Name)
			outEnd.Target.Stop()
			glogger.GLogger.Infof("OutEnd [%s, %s] stopped", uuid, outEnd.Name)
			e.OutEnds.Delete(uuid)
			glogger.GLogger.Infof("OutEnd [%s, %s] has been deleted", uuid, outEnd.Name)
			outEnd = nil
		}
	}
}

func (e *RuleEngine) AllOutEnds() *sync.Map {
	return e.OutEnds
}

// -----------------------------------------------------------------
// 获取运行时快照
// -----------------------------------------------------------------
func (e *RuleEngine) SnapshotDump() string {
	inends := []interface{}{}
	rules := []interface{}{}
	plugins := []interface{}{}
	outends := []interface{}{}
	devices := []interface{}{}
	e.AllInEnds().Range(func(key, value interface{}) bool {
		inends = append(inends, value)
		return true
	})
	e.AllRules().Range(func(key, value interface{}) bool {
		rules = append(rules, value)
		return true
	})
	e.AllPlugins().Range(func(key, value interface{}) bool {
		plugins = append(plugins, (value.(typex.XPlugin)).PluginMetaInfo())
		return true
	})
	e.AllOutEnds().Range(func(key, value interface{}) bool {
		outends = append(outends, value)
		return true
	})
	e.AllDevices().Range(func(key, value interface{}) bool {
		devices = append(devices, value)
		return true
	})

	parts, _ := disk.Partitions(true)
	diskInfo, _ := disk.Usage(parts[0].Mountpoint)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	system := map[string]interface{}{
		"version":  typex.MainVersion,
		"diskInfo": int(diskInfo.UsedPercent),
		"system":   utils.BToMb(m.Sys),
		"alloc":    utils.BToMb(m.Alloc),
		"total":    utils.BToMb(m.TotalAlloc),
		"osArch":   runtime.GOOS + "-" + runtime.GOARCH,
	}
	data := map[string]interface{}{
		"rules":      rules,
		"plugins":    plugins,
		"inends":     inends,
		"outends":    outends,
		"devices":    devices,
		"statistics": intermetric.GetMetric(),
		"system":     system,
		"config":     core.GlobalConfig,
	}
	b, err := json.Marshal(data)
	if err != nil {
		glogger.GLogger.Error(err)
		return err.Error()
	}
	return string(b)
}

// 重启源
func (e *RuleEngine) RestartInEnd(uuid string) error {
	if Value, ok := e.InEnds.Load(uuid); ok {
		InEnd := Value.(*typex.InEnd)
		if InEnd.Source != nil {
			InEnd.Source.Details().State = typex.SOURCE_DOWN // Down 以后会被自动拉起来
		}
		return nil
	}
	return fmt.Errorf("inend not exists:%s", uuid)
}

// 重启目标
func (e *RuleEngine) RestartOutEnd(uuid string) error {
	if Value, ok := e.OutEnds.Load(uuid); ok {
		OutEnd := Value.(*typex.OutEnd)
		if OutEnd.Target != nil {
			OutEnd.Target.Details().State = typex.SOURCE_DOWN // Down 以后会被自动拉起来
		}
		return nil
	}
	return fmt.Errorf("outEnd not exists:%s", uuid)

}

// 重启设备
func (e *RuleEngine) RestartDevice(uuid string) error {
	if Value, ok := e.Devices.Load(uuid); ok {
		Device := Value.(*typex.Device)
		if Device.Device != nil {
			Device.Device.SetState(typex.DEV_DOWN) // Down 以后会被自动拉起来
		}
		return nil
	}
	return fmt.Errorf("device not exists:%s", uuid)
}

/*
*
* 初始化设备管理器
*
 */

func (e *RuleEngine) InitDeviceTypeManager() error {
	e.DeviceTypeManager.Register(typex.KNX_GATEWAY,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewKNXGateway,
		},
	)
	e.DeviceTypeManager.Register(typex.LORA_WAN_GATEWAY,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewLoraGateway,
		},
	)
	e.DeviceTypeManager.Register(typex.TENCENT_IOTHUB_GATEWAY,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewTencentIoTGateway,
		},
	)
	e.DeviceTypeManager.Register(typex.SMART_HOME_CONTROLLER,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewShellyGen1ProxyGateway,
		},
	)
	e.DeviceTypeManager.Register(typex.GENERIC_HTTP_DEVICE,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewGenericHttpDevice,
		},
	)
	e.DeviceTypeManager.Register(typex.GENERIC_CAMERA,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewVideoCamera,
		},
	)
	e.DeviceTypeManager.Register(typex.SIEMENS_PLC,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewSIEMENS_PLC,
		},
	)
	e.DeviceTypeManager.Register(typex.GENERIC_MODBUS_MASTER,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewGenericModbusMaster,
		},
	)
	e.DeviceTypeManager.Register(typex.GENERIC_MODBUS_SLAVER,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewGenericModbusSlaver,
		},
	)
	e.DeviceTypeManager.Register(typex.GENERIC_UART_RW,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewGenericUartDevice,
		},
	)
	e.DeviceTypeManager.Register(typex.GENERIC_SNMP,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewGenericSnmpDevice,
		},
	)
	e.DeviceTypeManager.Register(typex.GENERIC_UART_PROTOCOL,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewGenericUartProtocolDevice,
		},
	)
	e.DeviceTypeManager.Register(typex.GENERIC_OPCUA,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewGenericOpcuaDevice,
		},
	)
	e.DeviceTypeManager.Register(typex.GENERIC_CAMERA,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewVideoCamera,
		},
	)
	e.DeviceTypeManager.Register(typex.GENERIC_AIS_RECEIVER,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewAISDeviceMaster,
		},
	)
	e.DeviceTypeManager.Register(typex.GENERIC_BACNET_IP,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewGenericBacnetIpDevice,
		},
	)
	e.DeviceTypeManager.Register(typex.BACNET_ROUTER_GW,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewBacnetRouter,
		},
	)
	e.DeviceTypeManager.Register(typex.HNC8,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewHNC8_CNC,
		},
	)
	e.DeviceTypeManager.Register(typex.KDN,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewKDN_CNC,
		},
	)
	return nil
}

/*
*
* 初始化输入资源管理器
*
 */
func (e *RuleEngine) InitSourceTypeManager() error {
	e.SourceTypeManager.Register(typex.COMTC_EVENT_FORWARDER,
		&typex.XConfig{
			Engine:    e,
			NewSource: source.NewTransceiverForwarder,
		},
	)
	e.SourceTypeManager.Register(typex.MQTT,
		&typex.XConfig{
			Engine:    e,
			NewSource: source.NewMqttInEndSource,
		},
	)
	e.SourceTypeManager.Register(typex.HTTP,
		&typex.XConfig{
			Engine:    e,
			NewSource: source.NewHttpInEndSource,
		},
	)
	e.SourceTypeManager.Register(typex.COAP,
		&typex.XConfig{
			Engine:    e,
			NewSource: source.NewCoAPInEndSource,
		},
	)
	e.SourceTypeManager.Register(typex.GRPC,
		&typex.XConfig{
			Engine:    e,
			NewSource: source.NewGrpcInEndSource,
		},
	)
	e.SourceTypeManager.Register(typex.NATS_SERVER,
		&typex.XConfig{
			Engine:    e,
			NewSource: source.NewNatsSource,
		},
	)
	e.SourceTypeManager.Register(typex.UDP_SERVER,
		&typex.XConfig{
			Engine:    e,
			NewSource: source.NewUdpInEndSource,
		},
	)
	e.SourceTypeManager.Register(typex.TCP_SERVER,
		&typex.XConfig{
			Engine:    e,
			NewSource: source.NewTcpSource,
		},
	)
	e.SourceTypeManager.Register(typex.GENERIC_IOT_HUB,
		&typex.XConfig{
			Engine:    e,
			NewSource: source.NewIoTHubSource,
		},
	)
	e.SourceTypeManager.Register(typex.INTERNAL_EVENT,
		&typex.XConfig{
			Engine:    e,
			NewSource: source.NewInternalEventSource,
		},
	)
	e.SourceTypeManager.Register(typex.GENERIC_MQTT,
		&typex.XConfig{
			Engine:    e,
			NewSource: source.NewGenericMqttSource,
		},
	)
	e.SourceTypeManager.Register(typex.GENERIC_MQTT_SERVER,
		&typex.XConfig{
			Engine:    e,
			NewSource: source.NewMqttServer,
		},
	)
	return nil
}

/*
*
* 初始化输出资源管理器
*
 */
func (e *RuleEngine) InitTargetTypeManager() error {

	e.TargetTypeManager.Register(typex.SEMTECH_UDP_FORWARDER,
		&typex.XConfig{
			Engine:    e,
			NewTarget: target.NewSemtechUdpForwarder,
		},
	)
	e.TargetTypeManager.Register(typex.GENERIC_UART_TARGET,
		&typex.XConfig{
			Engine:    e,
			NewTarget: target.NewGenericUart,
		},
	)
	e.TargetTypeManager.Register(typex.MONGO_SINGLE,
		&typex.XConfig{
			Engine:    e,
			NewTarget: target.NewMongoTarget,
		},
	)
	e.TargetTypeManager.Register(typex.MQTT_TARGET,
		&typex.XConfig{
			Engine:    e,
			NewTarget: target.NewMqttTarget,
		},
	)
	e.TargetTypeManager.Register(typex.NATS_TARGET,
		&typex.XConfig{
			Engine:    e,
			NewTarget: target.NewNatsTarget,
		},
	)
	e.TargetTypeManager.Register(typex.HTTP_TARGET,
		&typex.XConfig{
			Engine:    e,
			NewTarget: target.NewHTTPTarget,
		},
	)
	e.TargetTypeManager.Register(typex.TDENGINE_TARGET,
		&typex.XConfig{
			Engine:    e,
			NewTarget: target.NewTdEngineTarget,
		},
	)
	e.TargetTypeManager.Register(typex.GRPC_CODEC_TARGET,
		&typex.XConfig{
			Engine:    e,
			NewTarget: target.NewCodecTarget,
		},
	)
	e.TargetTypeManager.Register(typex.UDP_TARGET,
		&typex.XConfig{
			Engine:    e,
			NewTarget: target.NewUdpTarget,
		},
	)
	e.TargetTypeManager.Register(typex.TCP_TRANSPORT,
		&typex.XConfig{
			Engine:    e,
			NewTarget: target.NewTTcpTarget,
		},
	)
	return nil
}

/*
*-----------------------------------------------------------------
* 0.6.8 New Api: 将注册权交给设备
*-----------------------------------------------------------------
 */
func RegisterNewDevice(Type typex.DeviceType, Cfg *typex.XConfig) error {
	Cfg.Engine = __DefaultRuleEngine
	__DefaultRuleEngine.DeviceTypeManager.Register(Type, Cfg)
	return nil
}
func RegisterNewSource(Type typex.InEndType, Cfg *typex.XConfig) error {
	Cfg.Engine = __DefaultRuleEngine
	__DefaultRuleEngine.SourceTypeManager.Register(Type, Cfg)
	return nil
}
func RegisterNewTarget(Type typex.TargetType, Cfg *typex.XConfig) error {
	Cfg.Engine = __DefaultRuleEngine
	__DefaultRuleEngine.TargetTypeManager.Register(Type, Cfg)
	return nil
}
