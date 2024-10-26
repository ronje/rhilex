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

package rhilexmanager

import (
	"github.com/hootrhino/rhilex/target"
	"github.com/hootrhino/rhilex/typex"
)

var DefaultTargetTypeManager *TargetTypeManager

type TargetTypeManager struct {
	registry map[typex.TargetType]*typex.XConfig
}

func InitTargetTypeManager(e typex.Rhilex) {
	DefaultTargetTypeManager = &TargetTypeManager{
		registry: map[typex.TargetType]*typex.XConfig{},
	}

	DefaultTargetTypeManager.Register(typex.SEMTECH_UDP_FORWARDER,
		&typex.XConfig{
			Engine:    e,
			NewTarget: target.NewSemtechUdpForwarder,
		},
	)
	DefaultTargetTypeManager.Register(typex.GENERIC_UART_TARGET,
		&typex.XConfig{
			Engine:    e,
			NewTarget: target.NewGenericUart,
		},
	)
	DefaultTargetTypeManager.Register(typex.MONGO_SINGLE,
		&typex.XConfig{
			Engine:    e,
			NewTarget: target.NewMongoTarget,
		},
	)
	DefaultTargetTypeManager.Register(typex.MQTT_TARGET,
		&typex.XConfig{
			Engine:    e,
			NewTarget: target.NewMqttTarget,
		},
	)
	DefaultTargetTypeManager.Register(typex.HTTP_TARGET,
		&typex.XConfig{
			Engine:    e,
			NewTarget: target.NewHTTPTarget,
		},
	)
	DefaultTargetTypeManager.Register(typex.TDENGINE_TARGET,
		&typex.XConfig{
			Engine:    e,
			NewTarget: target.NewTdEngineTarget,
		},
	)
	DefaultTargetTypeManager.Register(typex.RHILEX_GRPC_TARGET,
		&typex.XConfig{
			Engine:    e,
			NewTarget: target.NewRhilexRpcTarget,
		},
	)
	DefaultTargetTypeManager.Register(typex.UDP_TARGET,
		&typex.XConfig{
			Engine:    e,
			NewTarget: target.NewUUdpTarget,
		},
	)
	DefaultTargetTypeManager.Register(typex.TCP_TRANSPORT,
		&typex.XConfig{
			Engine:    e,
			NewTarget: target.NewTTcpTarget,
		},
	)
	DefaultTargetTypeManager.Register(typex.GREPTIME_DATABASE,
		&typex.XConfig{
			Engine:    e,
			NewTarget: target.NewGrepTimeDbTarget,
		},
	)
}
func NewTargetTypeManager() *TargetTypeManager {
	return &TargetTypeManager{
		registry: map[typex.TargetType]*typex.XConfig{},
	}

}
func (rm *TargetTypeManager) Register(name typex.TargetType, f *typex.XConfig) {
	rm.registry[name] = f
}

func (rm *TargetTypeManager) Find(name typex.TargetType) *typex.XConfig {

	return rm.registry[name]
}
func (rm *TargetTypeManager) All() []*typex.XConfig {
	data := make([]*typex.XConfig, 0)
	for _, v := range rm.registry {
		data = append(data, v)
	}
	return data
}

/**
 * 获取所有类型
 *
 */
func (rm *TargetTypeManager) AllKeys() []string {
	data := []string{}
	for k := range rm.registry {
		data = append(data, k.String())
	}
	return data
}
