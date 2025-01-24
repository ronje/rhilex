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
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/hootrhino/rhilex/component/eventbus"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

type TransceiverForwarderConfig struct {
	// 要监听哪个外设的数据?这个参数就是外部通信模块的NAME
	ComName string `json:"comName" validate:"required"`
}
type TransceiverForwarder struct {
	typex.XStatus
	mainConfig TransceiverForwarderConfig
	subscriber eventbus.Subscriber
}

func NewTransceiverForwarder(r typex.Rhilex) typex.XSource {
	u := TransceiverForwarder{}
	u.mainConfig = TransceiverForwarderConfig{}
	u.subscriber = eventbus.Subscriber{
		Callback: func(Topic string, Event eventbus.EventMessage) {
			if u.mainConfig.ComName != "" {
				handleEvent(&u, Event)
			}
		},
	}
	u.RuleEngine = r
	return &u
}

func handleEvent(u *TransceiverForwarder, Event eventbus.EventMessage) {
	switch T := Event.Payload.(type) {
	case []byte:
		comData := RuleData{
			ComName: u.mainConfig.ComName,
			Data:    hex.EncodeToString(T),
		}
		processData(u, comData)
	case string:
		comData := RuleData{
			ComName: u.mainConfig.ComName,
			Data:    T,
		}
		processData(u, comData)
	default:
		glogger.GLogger.Error(fmt.Errorf("unsupported data type:%v", T))
	}
}

func processData(u *TransceiverForwarder, comData RuleData) {
	work, err := u.RuleEngine.WorkInEnd(u.RuleEngine.GetInEnd(u.PointId),
		comData.String())
	if !work {
		glogger.GLogger.Error(err)
	}
}

func (u *TransceiverForwarder) Init(inEndId string, configMap map[string]interface{}) error {
	u.PointId = inEndId
	if err := utils.BindSourceConfig(configMap, &u.mainConfig); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	return nil
}

func (u *TransceiverForwarder) Start(cctx typex.CCTX) error {
	u.Ctx = cctx.Ctx
	u.CancelCTX = cctx.CancelCTX
	//lineS := "event.transceiver.data." + tc.mainConfig.Name
	eventbus.Subscribe("event.transceiver.data."+u.mainConfig.ComName, &u.subscriber)
	return nil

}

func (u *TransceiverForwarder) Status() typex.SourceState {
	return typex.SOURCE_UP
}

func (u *TransceiverForwarder) Stop() {
	eventbus.UnSubscribe(u.mainConfig.ComName, &u.subscriber)
	if u.CancelCTX != nil {
		u.CancelCTX()
	}
}

func (u *TransceiverForwarder) Details() *typex.InEnd {
	return u.RuleEngine.GetInEnd(u.PointId)
}

/*
*
* 从内部总线拿数据
* internotify.Insert(...)
 */
type RuleData struct {
	ComName string `json:"comName"`
	Data    string `json:"data"`
}

func (O RuleData) String() string {
	bytes, _ := json.Marshal(O)
	return string(bytes)
}
