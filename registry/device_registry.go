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

package registry

import (
	"github.com/hootrhino/rhilex/component/orderedmap"
	"github.com/hootrhino/rhilex/device"
	"github.com/hootrhino/rhilex/typex"
)

var DefaultDeviceRegistry *DeviceRegistry

type DeviceRegistry struct {
	e        typex.Rhilex
	registry *orderedmap.OrderedMap[typex.DeviceType, *typex.XConfig]
}

func InitDeviceRegistry(e typex.Rhilex) {
	DefaultDeviceRegistry = &DeviceRegistry{
		e:        e,
		registry: orderedmap.NewOrderedMap[typex.DeviceType, *typex.XConfig](),
	}
	LoadAllDeviceType(e)
}

func LoadAllDeviceType(e typex.Rhilex) {
	DefaultDeviceRegistry.Register(typex.GENERIC_NEMA_GNS_PROTOCOL,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewNemaGpsMasterDevice,
		},
	)
	DefaultDeviceRegistry.Register(typex.TAOJINGCHI_UARTHMI_MASTER,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewTaoJingChiHmiDevice,
		},
	)
	DefaultDeviceRegistry.Register(typex.SZY2062016_MASTER,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewSZY206_2016_MasterGateway,
		},
	)
	DefaultDeviceRegistry.Register(typex.CJT1882004_MASTER,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewCJT188_2004_MasterGateway,
		},
	)
	DefaultDeviceRegistry.Register(typex.DLT6452007_MASTER,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewDLT645_2007_MasterGateway,
		},
	)
	DefaultDeviceRegistry.Register(typex.KNX_GATEWAY,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewKNXGateway,
		},
	)
	DefaultDeviceRegistry.Register(typex.LORA_WAN_GATEWAY,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewLoraGateway,
		},
	)
	DefaultDeviceRegistry.Register(typex.GENERIC_HTTP_DEVICE,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewGenericHttpDevice,
		},
	)
	DefaultDeviceRegistry.Register(typex.GENERIC_CAMERA,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewVideoCamera,
		},
	)
	DefaultDeviceRegistry.Register(typex.SIEMENS_PLC,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewSIEMENS_PLC,
		},
	)
	DefaultDeviceRegistry.Register(typex.GENERIC_MODBUS_MASTER,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewGenericModbusMaster,
		},
	)
	DefaultDeviceRegistry.Register(typex.GENERIC_MODBUS_SLAVER,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewGenericModbusSlaver,
		},
	)
	DefaultDeviceRegistry.Register(typex.GENERIC_UART_RW,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewGenericUartDevice,
		},
	)
	DefaultDeviceRegistry.Register(typex.GENERIC_SNMP,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewGenericSnmpDevice,
		},
	)
	DefaultDeviceRegistry.Register(typex.GENERIC_USER_PROTOCOL,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewGenericUserProtocolDevice,
		},
	)
	DefaultDeviceRegistry.Register(typex.GENERIC_CAMERA,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewVideoCamera,
		},
	)
	DefaultDeviceRegistry.Register(typex.GENERIC_AIS_RECEIVER,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewAISDeviceMaster,
		},
	)
	DefaultDeviceRegistry.Register(typex.GENERIC_BACNET_IP,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewGenericBacnetIpDevice,
		},
	)
	DefaultDeviceRegistry.Register(typex.BACNET_ROUTER_GW,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewBacnetRouter,
		},
	)
	DefaultDeviceRegistry.Register(typex.GENERIC_MBUS_EN13433_MASTER,
		&typex.XConfig{
			Engine:    e,
			NewDevice: device.NewMBusEn13433MasterGateway,
		},
	)
}
func (rm *DeviceRegistry) Register(name typex.DeviceType, f *typex.XConfig) {
	f.Type = string(name)
	rm.registry.Set(name, f)
}

func (rm *DeviceRegistry) Find(name typex.DeviceType) *typex.XConfig {
	if xcfg, ok := rm.registry.Get(name); ok {
		return xcfg
	}
	return nil
}
func (rm *DeviceRegistry) All() []*typex.XConfig {
	return rm.registry.Values()
}

/**
 * 获取所有类型
 *
 */
func (rm *DeviceRegistry) AllKeys() []string {
	data := []string{}
	for _, k := range rm.registry.Keys() {
		data = append(data, k.String())
	}
	return data
}

func (rm *DeviceRegistry) Stop() {
}
