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

package USBMonitorPlugin

import (
	"errors"

	"github.com/hootrhino/rhilex/typex"

	"gopkg.in/ini.v1"
)

/*
*
* USB 热插拔监控器, 方便观察USB状态, 本插件只支持Linux！！！
*
 */
type USBMonitorPlugin struct {
	uuid string
}

func NewUSBMonitorPlugin() typex.XPlugin {
	return &USBMonitorPlugin{
		uuid: "USB-MONITOR",
	}
}
func (usbm *USBMonitorPlugin) Init(_ *ini.Section) error {
	return nil

}

func (usbm *USBMonitorPlugin) Start(_ typex.Rhilex) error {
	return errors.New("USB monitor plugin not support")
}

func (usbm *USBMonitorPlugin) Stop() error {
	return nil
}

func (usbm *USBMonitorPlugin) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		UUID:        usbm.uuid,
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
func (usbm *USBMonitorPlugin) Service(arg typex.ServiceArg) typex.ServiceResult {
	return typex.ServiceResult{}
}
