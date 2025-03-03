// Copyright (C) 2025 wwhai
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
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package usbmonitor

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"

	"gopkg.in/ini.v1"
)

/*
*
* USB 热插拔监控器, 方便观察USB状态, 本插件只支持Linux！！！
*
 */
type USBMonitorPlugin struct {
	uuid    string
	monitor *USBMonitor
}

func NewUSBMonitorPlugin() typex.XPlugin {
	return &USBMonitorPlugin{
		uuid:    "USB-MONITOR",
		monitor: nil,
	}
}
func (usbMonitor *USBMonitorPlugin) Init(_ *ini.Section) error {
	return nil
}

func (usbMonitor *USBMonitorPlugin) Start(_ typex.Rhilex) error {
	// 不支持windows
	if strings.Contains(runtime.GOOS, "windows") {
		return fmt.Errorf("USB Monitor not support windows")
	}
	monitor, err := NewUSBMonitor()
	if err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	usbMonitor.monitor = monitor
	usbMonitor.monitor.callback = func(event, device string) {
		if strings.Contains(device, "ttyUSB") {
			if event == "create" {
				glogger.GLogger.Info("USB device connected: ", device)
			}
			if event == "remove" {
				glogger.GLogger.Info("USB device disconnected: ", device)
			}

		}
	}
	usbMonitor.monitor.Start()
	glogger.GLogger.Info("USB Monitor Started")
	return nil
}

func (usbMonitor *USBMonitorPlugin) Stop() error {
	if usbMonitor.monitor != nil {
		usbMonitor.monitor.Stop()
	}
	glogger.GLogger.Info("USB Monitor Stopped")
	return nil
}

func (usbMonitor *USBMonitorPlugin) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		UUID:        usbMonitor.uuid,
		Name:        "USB Monitor",
		Version:     "v0.0.1",
		Description: "USB Hot Plugin Monitor",
	}
}

/*
*
* 服务调用接口
*
 */
func (usbMonitor *USBMonitorPlugin) Service(arg typex.ServiceArg) typex.ServiceResult {
	return typex.ServiceResult{}
}
