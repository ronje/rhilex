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

package target

import (
	"errors"
	"fmt"
	"time"

	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

//
/*
*
* 单向的MQTT客户端，不支持subscribe，订阅了不生效
*
 */
type MqttTargetConfig struct {
	Host     string `json:"host" validate:"required" title:"服务地址"`
	Port     int    `json:"port" validate:"required" title:"服务端口"`
	ClientId string `json:"clientId" validate:"required" title:"客户端ID"`
	Username string `json:"username" validate:"required" title:"连接账户"`
	Password string `json:"password" validate:"required" title:"连接密码"`
	PubTopic string `json:"pubTopic" title:"上报TOPIC" info:"上报TOPIC"` // 上报数据的 Topic
	SubTopic string `json:"subTopic" title:"订阅TOPIC" info:"订阅TOPIC"` // 上报数据的 Topic
}
type mqttOutEndTarget struct {
	typex.XStatus
	client     mqtt.Client
	mainConfig MqttTargetConfig
	status     typex.SourceState
}

func NewMqttTarget(e typex.Rhilex) typex.XTarget {
	m := new(mqttOutEndTarget)
	m.RuleEngine = e
	m.mainConfig = MqttTargetConfig{
		Host: "127.0.0.1",
		Port: 1883,
	}
	m.status = typex.SOURCE_DOWN
	return m
}

func (mq *mqttOutEndTarget) Init(outEndId string, configMap map[string]interface{}) error {
	mq.PointId = outEndId
	if err := utils.BindSourceConfig(configMap, &mq.mainConfig); err != nil {
		return err
	}
	return nil
}
func (mq *mqttOutEndTarget) Start(cctx typex.CCTX) error {
	mq.Ctx = cctx.Ctx
	mq.CancelCTX = cctx.CancelCTX

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%v", mq.mainConfig.Host, mq.mainConfig.Port))
	opts.SetClientID(mq.mainConfig.ClientId)
	opts.SetUsername(mq.mainConfig.Username)
	opts.SetPassword(mq.mainConfig.Password)
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		glogger.GLogger.Infof("Mqtt Connected Success")
	})
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		glogger.GLogger.Warn("Mqtt Connect lost:", err)
	})
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(false)    //不需要自动重连, 交给RHILEX管理
	opts.SetMaxReconnectInterval(0) // 不需要自动重连, 交给RHILEX管理
	opts.SetConnectTimeout(5 * time.Second)
	opts.SetPingTimeout(5 * time.Second)
	mq.client = mqtt.NewClient(opts)
	token := mq.client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	mq.status = typex.SOURCE_UP
	return nil
}

func (mq *mqttOutEndTarget) Stop() {
	mq.status = typex.SOURCE_DOWN
	if mq.CancelCTX != nil {
		mq.CancelCTX()
	}
	if mq.client != nil {
		mq.client.Disconnect(10)
	}
}

func (mq *mqttOutEndTarget) Status() typex.SourceState {
	if mq.client != nil {
		if mq.client.IsConnected() && mq.client.IsConnectionOpen() {
			return typex.SOURCE_UP
		}
		return typex.SOURCE_DOWN
	}
	return mq.status
}

func (mq *mqttOutEndTarget) Details() *typex.OutEnd {
	return mq.RuleEngine.GetOutEnd(mq.PointId)
}

func (mq *mqttOutEndTarget) To(data interface{}) (interface{}, error) {
	if mq.client != nil {
		switch s := data.(type) {
		case string:
			glogger.GLogger.Debug("Target publish:", mq.mainConfig.PubTopic, 1, false, data)
			token := mq.client.Publish(mq.mainConfig.PubTopic, 1, false, s)
			return nil, token.Error()
		default:
			return nil, errors.New("Invalid mqtt data type")
		}
	}
	return nil, errors.New("mqtt client is nil")
}
