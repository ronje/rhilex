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
	"github.com/hootrhino/rhilex/component/orderedmap"
	"github.com/hootrhino/rhilex/target"
	"github.com/hootrhino/rhilex/typex"
)

var DefaultTargetTypeManager *TargetTypeManager

type TargetTypeManager struct {
	e        typex.Rhilex
	registry *orderedmap.OrderedMap[typex.TargetType, *typex.XConfig]
}

func InitTargetTypeManager(e typex.Rhilex) {
	DefaultTargetTypeManager = &TargetTypeManager{
		e:        e,
		registry: orderedmap.NewOrderedMap[typex.TargetType, *typex.XConfig](),
	}
	LoadAllTargetType(e)
}

func LoadAllTargetType(e typex.Rhilex) {

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

func (rm *TargetTypeManager) Register(name typex.TargetType, f *typex.XConfig) {
	f.Type = string(name)
	rm.registry.Set(name, f)
}

func (rm *TargetTypeManager) Find(name typex.TargetType) *typex.XConfig {
	p, ok := rm.registry.Get(name)
	if ok {
		return p
	}
	return nil
}
func (rm *TargetTypeManager) All() []*typex.XConfig {
	return rm.registry.Values()
}

/**
 * 获取所有类型
 *
 */
func (rm *TargetTypeManager) AllKeys() []string {
	data := []string{}
	for _, k := range rm.registry.Keys() {
		data = append(data, k.String())
	}
	return data
}

func (rm *TargetTypeManager) Stop() {
}
