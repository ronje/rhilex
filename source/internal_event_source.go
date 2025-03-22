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

package source

import (
	"encoding/json"

	"github.com/hootrhino/rhilex/component/eventbus"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

/*
*
* 用来将内部消息总线的事件推到资源脚本
*
 */
type __InternalEventSourceConfig struct {
	// - ALL: 全部事件
	// - SOURCE: 南向事件
	// - DEVICE: 设备事件
	// - TARGET: 北向事件
	// - SYSTEM: 系统事件
	// - HARDWARE: 硬件事件
	Type string `json:"type"`
}
type InternalEventSource struct {
	typex.XStatus
	mainConfig __InternalEventSourceConfig
	subscriber eventbus.Subscriber
}

func NewInternalEventSource(r typex.Rhilex) typex.XSource {
	u := InternalEventSource{}
	u.mainConfig = __InternalEventSourceConfig{
		Type: "ALL",
	}
	u.subscriber = eventbus.Subscriber{
		Callback: func(Topic string, Event eventbus.EventMessage) {
			bytes, _ := json.Marshal(event{
				Type:  Event.Type,
				Event: Event.Event,
				Ts:    Event.Ts,
				Info:  Event.Payload,
			})
			u.RuleEngine.WorkInEnd(u.RuleEngine.GetInEnd(u.PointId), string(bytes))
		},
	}
	u.RuleEngine = r
	return &u
}

func (u *InternalEventSource) Init(inEndId string, configMap map[string]any) error {
	u.PointId = inEndId
	if err := utils.BindSourceConfig(configMap, &u.mainConfig); err != nil {
		return err
	}
	return nil
}

func (u *InternalEventSource) Start(cctx typex.CCTX) error {
	u.Ctx = cctx.Ctx
	u.CancelCTX = cctx.CancelCTX
	eventbus.Subscribe("*", &u.subscriber)
	return nil

}

func (u *InternalEventSource) Status() typex.SourceState {
	return typex.SOURCE_UP
}

func (u *InternalEventSource) Stop() {
	eventbus.UnSubscribe(u.PointId, &u.subscriber)
	if u.CancelCTX != nil {
		u.CancelCTX()
	}
}

func (u *InternalEventSource) Details() *typex.InEnd {
	return u.RuleEngine.GetInEnd(u.PointId)
}

// 来自外面的数据
func (*InternalEventSource) DownStream([]byte) (int, error) {
	return 0, nil
}

// 上行数据
func (*InternalEventSource) UpStream([]byte) (int, error) {
	return 0, nil
}

type event struct {
	// - ALL: 全部
	// - SOURCE: 南向事件
	// - DEVICE: 设备事件
	// - TARGET: 北向事件
	// - SYSTEM: 系统内部事件
	// - HARDWARE: 硬件事件
	Type  string `json:"type"`
	Event string `json:"event"`
	Ts    uint64 `json:"ts"`
	Info  any    `json:"info"`
}
