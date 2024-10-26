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
	"github.com/hootrhino/rhilex/source"
	"github.com/hootrhino/rhilex/typex"
)

var DefaultSourceTypeManager *SourceTypeManager

type SourceTypeManager struct {
	registry map[typex.InEndType]*typex.XConfig
}

func InitSourceTypeManager(e typex.Rhilex) {
	DefaultSourceTypeManager = &SourceTypeManager{
		registry: map[typex.InEndType]*typex.XConfig{},
	}
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

func NewSourceTypeManager() *SourceTypeManager {
	return &SourceTypeManager{
		registry: map[typex.InEndType]*typex.XConfig{},
	}

}
func (rm *SourceTypeManager) Register(name typex.InEndType, f *typex.XConfig) {
	rm.registry[name] = f
}

func (rm *SourceTypeManager) Find(name typex.InEndType) *typex.XConfig {

	return rm.registry[name]
}
func (rm *SourceTypeManager) All() []*typex.XConfig {
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
func (rm *SourceTypeManager) AllKeys() []string {
	data := []string{}
	for k := range rm.registry {
		data = append(data, k.String())
	}
	return data
}
