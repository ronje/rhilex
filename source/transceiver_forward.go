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
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hootrhino/rhilex/component/internotify"
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
}

func NewTransceiverForwarder(r typex.Rhilex) typex.XSource {
	s := TransceiverForwarder{}
	s.mainConfig = TransceiverForwarderConfig{}
	s.RuleEngine = r
	return &s
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
	u.startInternalEventQueue(u.Ctx)
	return nil

}

func (u *TransceiverForwarder) Status() typex.SourceState {
	return typex.SOURCE_UP
}

func (u *TransceiverForwarder) Stop() {
	if u.CancelCTX != nil {
		u.CancelCTX()
	}
}

func (u *TransceiverForwarder) Details() *typex.InEnd {
	return u.RuleEngine.GetInEnd(u.PointId)
}

func (u *TransceiverForwarder) Test(inEndId string) bool {
	return true
}

func (*TransceiverForwarder) DownStream([]byte) (int, error) {
	return 0, nil
}

func (*TransceiverForwarder) UpStream([]byte) (int, error) {
	return 0, nil
}

/*
*
* 从内部总线拿数据
* internotify.Push(...)
 */
type RuleData struct {
	ComName string `json:"comName"`
	Data    string `json:"data"`
}

func (O RuleData) String() string {
	if bytes, err := json.Marshal(O); err != nil {
		return "{}"
	} else {
		return string(bytes)
	}
}
func (u *TransceiverForwarder) startInternalEventQueue(ctxU context.Context) {
	go func(ctx context.Context) {
		Queue := make(chan internotify.BaseEvent, 64) // 64是个魔法数字，未来需要优化为动态配置
		ID := fmt.Sprintf("TransceiverForwarder:%s", u.PointId)
		Subscriber := internotify.Subscriber{
			Id:      ID,
			Channel: &Queue,
		}
		internotify.AddSubscriber(Subscriber)
		glogger.GLogger.Debugf("Start Transceiver Forwarder:%s", u.PointId)
		defer internotify.RemoveSubscriber(ID)
		for {
			select {
			case <-ctxU.Done():
				return
			case Event := <-*Subscriber.Channel:
				// 过滤不感去兴趣的事件
				// "transceiver.upstream.data.$ComName"
				if !strings.Contains(Event.Event, u.mainConfig.ComName) {
					continue
				}
				glogger.GLogger.Debug(ID, " Received Data:", Event.String())
				switch T := Event.Info.(type) {
				case []byte:
					comData := RuleData{
						ComName: u.mainConfig.ComName,
						Data:    hex.EncodeToString(T),
					}
					work, err := u.RuleEngine.WorkInEnd(u.RuleEngine.GetInEnd(u.PointId),
						comData.String())
					if !work {
						glogger.GLogger.Error(err)
						continue
					}
				case string:
					comData := RuleData{
						ComName: u.mainConfig.ComName,
						Data:    (T),
					}
					work, err := u.RuleEngine.WorkInEnd(u.RuleEngine.GetInEnd(u.PointId),
						comData.String())
					if !work {
						glogger.GLogger.Error(err)
						continue
					}
				default:
					glogger.GLogger.Error(fmt.Errorf("unsupported data type:%v", T))
				}
			}
		}
	}(u.Ctx)
}
