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

/*
*
* RhilexG1 硬件接口相关管理
* 警告：此处会随着硬件不同而不兼容，要移植的时候记得统一一下目标硬件的端口
*
 */
package uartctrl

import (
	"encoding/json"
	"fmt"
	"runtime"
	"sync"

	"github.com/hootrhino/rhilex/typex"
)

var __UartsManager *UartsManager

type UartsManager struct {
	Interfaces sync.Map
	rhilex     typex.Rhilex
}

func InitUartsManager(rhilex typex.Rhilex) *UartsManager {
	__UartsManager = &UartsManager{
		Interfaces: sync.Map{},
		rhilex:     rhilex,
	}
	return __UartsManager
}

/*
*
* 这里记录一些RHILEXG1网关的硬件接口信息,同时记录串口是否被占用等
*
 */
type UartConfig struct {
	Timeout  int    `json:"timeout"`
	Uart     string `json:"uart"`
	BaudRate int    `json:"baudRate"`
	DataBits int    `json:"dataBits"`
	Parity   string `json:"parity"`
	StopBits int    `json:"stopBits"`
}
type UartOccupy struct {
	UUID string `json:"uuid"` // UUID
	Type string `json:"type"` // DEVICE, OS,... Other......
	Name string `json:"name"` // 占用的设备名称
}

func (O UartOccupy) String() string {
	return fmt.Sprintf("Occupied By: (%s,%s), Type is %s", O.UUID, O.Name, O.Type)
}

type SystemUart struct {
	UUID        string      `json:"uuid"`        // 接口名称
	Name        string      `json:"name"`        // 接口名称
	Alias       string      `json:"alias"`       // 别名
	Busy        bool        `json:"busy"`        // 运行时数据，是否被占
	OccupyBy    UartOccupy  `json:"occupyBy"`    // 运行时数据，被谁占用了 UUID
	Type        string      `json:"type"`        // 接口类型, UART(串口),USB(USB),FD(通用文件句柄)
	Description string      `json:"description"` // 额外备注
	Config      interface{} `json:"config"`      // 配置, 串口配置、或者网卡、USB等
}

func (v SystemUart) String() string {
	b, _ := json.Marshal(v)
	return string(b)
}

/*
*
* 加载配置到运行时, 需要刷新与配置相关的所有设备
*
 */
func SetUart(Port SystemUart) {
	__UartsManager.Interfaces.Store(Port.Name, Port)
	refreshUart(Port.Name)
}
func RemovePort(PortName string) {
	__UartsManager.Interfaces.Delete(PortName)
	refreshUart(PortName)
}
func RefreshPort(Port SystemUart) {
	__UartsManager.Interfaces.Store(Port.Name, Port)
	refreshUart(Port.Name)
}

/*
*
* 刷新所有关联的设备, 也就是 OccupyBy=UUID 需要重载
*
 */
func refreshUart(name string) {
	Object, ok := __UartsManager.Interfaces.Load(name)
	if !ok {
		return
	}
	Port := Object.(SystemUart)
	if Port.Busy {
		if Port.OccupyBy.Type == "DEVICE" {
			UUID := Port.OccupyBy.UUID
			if Device := __UartsManager.rhilex.GetDevice(UUID); Device != nil {
				// 拉闸 DEV_DOWN 以后就重启了, 然后就会拉取最新的配置
				Device.Device.SetState(typex.DEV_DOWN)
			}
		}
	}

}

/*
*
* 获取一个运行时状态
*
 */
func GetUart(name string) (SystemUart, error) {
	if Object, ok := __UartsManager.Interfaces.Load(name); ok {
		return Object.(SystemUart), nil
	}
	return SystemUart{}, fmt.Errorf("interface not exists:%s", name)
}

/*
*
* 所有的接口
*
 */
func AllUart() []SystemUart {
	result := []SystemUart{}
	__UartsManager.Interfaces.Range(func(key, Object any) bool {
		// 如果不是被rhilex占用；则需要检查是否被操作系统进程占用了
		Port := Object.(SystemUart)
		if Port.OccupyBy.Type != "DEVICE" {
			if err := CheckSerialBusy(Port.Name); err != nil {
				SetInterfaceBusy(Port.Name, UartOccupy{
					UUID: runtime.GOOS,
					Type: "OS",
					Name: runtime.GOOS,
				})
			} else {
				FreeInterfaceBusy(Port.Name)
			}
		}
		result = append(result, Port)
		return true
	})

	return result
}

/*
*
* 忙碌
*
 */
func SetInterfaceBusy(name string, OccupyBy UartOccupy) {
	if Object, ok := __UartsManager.Interfaces.Load(name); ok {
		Port := Object.(SystemUart)
		Port.Busy = true
		Port.OccupyBy = OccupyBy
		__UartsManager.Interfaces.Store(name, Port)
	}
}

/*
*
* 释放
*
 */
func FreeInterfaceBusy(name string) {
	if Object, ok := __UartsManager.Interfaces.Load(name); ok {
		Port := Object.(SystemUart)
		if Port.OccupyBy.Type == "DEVICE" {
			Port.Busy = false
			Port.OccupyBy = UartOccupy{
				"-", "-", "-",
			}
			__UartsManager.Interfaces.Store(name, Port)
		}
	}
}
