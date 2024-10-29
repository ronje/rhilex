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
	"github.com/hootrhino/rhilex/source"
	"github.com/hootrhino/rhilex/typex"
)

var DefaultSourceTypeManager *SourceTypeManager

type SourceTypeManager struct {
	e        typex.Rhilex
	registry *orderedmap.OrderedMap[typex.InEndType, *typex.XConfig]
}

func InitSourceTypeManager(e typex.Rhilex) {
	DefaultSourceTypeManager = &SourceTypeManager{
		e:        e,
		registry: orderedmap.NewOrderedMap[typex.InEndType, *typex.XConfig](),
	}
	LoadAllSourceType(e)
}

func LoadAllSourceType(e typex.Rhilex) {
	DefaultSourceTypeManager.Register(typex.CUSTOM_PROTOCOL_SERVER,
		&typex.XConfig{
			Engine:    e,
			NewSource: source.NewCustomProtocol,
		},
	)
	DefaultSourceTypeManager.Register(typex.COMTC_EVENT_FORWARDER,
		&typex.XConfig{
			Engine:    e,
			NewSource: source.NewTransceiverForwarder,
		},
	)
	DefaultSourceTypeManager.Register(typex.HTTP_SERVER,
		&typex.XConfig{
			Engine:    e,
			NewSource: source.NewHttpInEndSource,
		},
	)
	DefaultSourceTypeManager.Register(typex.COAP_SERVER,
		&typex.XConfig{
			Engine:    e,
			NewSource: source.NewCoAPInEndSource,
		},
	)
	DefaultSourceTypeManager.Register(typex.GRPC_SERVER,
		&typex.XConfig{
			Engine:    e,
			NewSource: source.NewGrpcInEndSource,
		},
	)

	DefaultSourceTypeManager.Register(typex.UDP_SERVER,
		&typex.XConfig{
			Engine:    e,
			NewSource: source.NewUdpInEndSource,
		},
	)
	DefaultSourceTypeManager.Register(typex.TCP_SERVER,
		&typex.XConfig{
			Engine:    e,
			NewSource: source.NewTcpSource,
		},
	)
	DefaultSourceTypeManager.Register(typex.INTERNAL_EVENT,
		&typex.XConfig{
			Engine:    e,
			NewSource: source.NewInternalEventSource,
		},
	)
	DefaultSourceTypeManager.Register(typex.GENERIC_MQTT_SERVER,
		&typex.XConfig{
			Engine:    e,
			NewSource: source.NewGenericMqttSource,
		},
	)
	DefaultSourceTypeManager.Register(typex.GENERIC_MQTT_SERVER,
		&typex.XConfig{
			Engine:    e,
			NewSource: source.NewMqttServer,
		},
	)
}

func (rm *SourceTypeManager) Register(name typex.InEndType, f *typex.XConfig) {
	f.Type = string(name)
	rm.registry.Set(name, f)
}

func (rm *SourceTypeManager) Find(name typex.InEndType) *typex.XConfig {
	p, ok := rm.registry.Get(name)
	if ok {
		return p
	}
	return nil
}
func (rm *SourceTypeManager) All() []*typex.XConfig {
	return rm.registry.Values()
}

/**
 * 获取所有类型
 *
 */
func (rm *SourceTypeManager) AllKeys() []string {
	data := []string{}
	for _, k := range rm.registry.Keys() {
		data = append(data, k.String())
	}
	return data
}

func (rm *SourceTypeManager) Stop() {
}
