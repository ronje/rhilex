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
	"os"
	"os/signal"
	"strings"
	"syscall"

	plugins "github.com/hootrhino/rhilex/plugin"
	"github.com/hootrhino/rhilex/plugin/discover"
	wdog "github.com/hootrhino/rhilex/plugin/generic_watchdog"
	modbusscrc "github.com/hootrhino/rhilex/plugin/modbus_crc_tools"
	modbusscanner "github.com/hootrhino/rhilex/plugin/modbus_scanner"
	ngrokc "github.com/hootrhino/rhilex/plugin/ngrokc"
	telemetry "github.com/hootrhino/rhilex/plugin/telemetry"
	ttyterminal "github.com/hootrhino/rhilex/plugin/ttyd_terminal"
	usbmonitor "github.com/hootrhino/rhilex/plugin/usb_monitor"
	ini "gopkg.in/ini.v1"

	apiServer "github.com/hootrhino/rhilex/component/apiserver"
	"github.com/hootrhino/rhilex/component/globalinit"
	"github.com/hootrhino/rhilex/component/performance"
	core "github.com/hootrhino/rhilex/config"
	glogger "github.com/hootrhino/rhilex/glogger"
	icmpsender "github.com/hootrhino/rhilex/plugin/icmp_sender"
	license_manager "github.com/hootrhino/rhilex/plugin/license_manager"
	typex "github.com/hootrhino/rhilex/typex"
)

func RunRhilex(iniPath string) {
	mainConfig := core.InitGlobalConfig(iniPath)
	glogger.StartGLogger(glogger.LogConfig{
		AppID:         mainConfig.AppId,
		LogLevel:      mainConfig.LogLevel,
		EnableConsole: mainConfig.EnableConsole,
		DebugMode:     mainConfig.DebugMode,
		LogMaxSize:    mainConfig.LogMaxSize,
		LogMaxBackups: mainConfig.LogMaxBackups,
		LogMaxAge:     mainConfig.LogMaxAge,
		LogCompress:   mainConfig.LogCompress,
	})
	globalinit.InitGlobalInitManager()
	glogger.StartNewRealTimeLogger(mainConfig.LogLevel)
	performance.SetDebugMode(mainConfig.EnablePProf)
	performance.SetGomaxProcs(mainConfig.GomaxProcs)
	//
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTERM)
	engine := InitRuleEngine(mainConfig)
	engine.Start()
	license_manager := license_manager.NewLicenseManager(engine)
	if err := plugins.LoadPlugin("plugin.license_manager", license_manager); err != nil {
		glogger.GLogger.Error(err)
		return
	}
	apiServer := apiServer.NewHttpApiServer(engine)
	if err := plugins.LoadPlugin("plugin.http_server", apiServer); err != nil {
		glogger.GLogger.Error(err)
		return
	}
	// Load Plugin
	loadOtherPlugin()
	s := <-c
	glogger.GLogger.Warn("RHILEX Receive Stop Signal: ", s)
	engine.Stop()
}

// loadPlugin 根据Ini配置信息，加载插件
func loadOtherPlugin() {
	cfg, _ := ini.ShadowLoad(core.GlobalConfig.IniPath)
	sections := cfg.ChildSections("plugin")
	for _, section := range sections {
		name := strings.TrimPrefix(section.Name(), "plugin.")
		enable, err := section.GetKey("enable")
		if err != nil {
			continue
		}
		if !enable.MustBool(false) {
			glogger.GLogger.Warnf("Plugin is disable:%s", name)
			continue
		}
		var plugin typex.XPlugin
		if name == "usbmonitor" {
			plugin = usbmonitor.NewUsbMonitor()
		}
		if name == "icmpsender" {
			plugin = icmpsender.NewICMPSender()
		}
		if name == "modbus_scanner" {
			plugin = modbusscanner.NewModbusScanner()
		}
		if name == "ttyd" {
			plugin = ttyterminal.NewWebTTYPlugin()
		}
		if name == "modbus_crc_tools" {
			plugin = modbusscrc.NewModbusCrcCalculator()
		}
		if name == "soft_wdog" {
			plugin = wdog.NewGenericWatchDog()
		}
		if name == "ngrokc" {
			plugin = ngrokc.NewNgrokClient()
		}
		if name == "telemetry" {
			plugin = telemetry.NewTelemetry()
		}
		if name == "discover" {
			plugin = discover.NewDiscoverPlugin()
		}
		if plugin != nil {
			if err := plugins.LoadPlugin(section.Name(), plugin); err != nil {
				glogger.GLogger.Error(err)
			}
		}
	}
}
