// Copyright (C) 2024 wwhai
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
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package device

import (
	"encoding/json"
	"fmt"

	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/hootrhino/rhilex/device/ithings"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

const (
	// 属性
	_ithings_PropertyTopic   = "$thing/down/property/%v/%v"
	_ithings_PropertyUpTopic = "$thing/up/property/%v/%v"
	// 动作
	_ithings_ActionTopic   = "$thing/down/action/%v/%v"
	_ithings_ActionUpTopic = "$thing/up/action/%v/%v"
	// 子设备拓扑
	_ithings_operation_up   = "$gateway/operation/%s/%s"        //数据上行 Topic（用于发布）
	_ithings_operation_down = "$gateway/operation/result/%s/%s" //数据下行 Topic（用于订阅）
)

type IThingsGatewayConfig struct {
	ServerEndpoint string `json:"serverEndpoint" validate:"required"` //服务地址
	Mode           string `json:"mode" validate:"required"`           //模式: DEVICE|GATEWAY
	ProductId      string `json:"productId" validate:"required"`      //产品名
	DeviceName     string `json:"deviceName" validate:"required"`     //设备名
	DevicePsk      string `json:"devicePsk" validate:"required"`      //秘钥
	ClientId       string `json:"clientId" validate:"required"`       //客户端ID
}
type IThingsGatewayMainConfig struct {
	CommonConfig IThingsGatewayConfig `json:"tencentConfig" validate:"required"` // 通用配置
}
type IThingsSubDevice struct {
	ProductID    string `json:"productID"`
	DeviceName   string `json:"deviceName"`
	Signature    string `json:"signature"`
	Random       string `json:"random"`
	Timestamp    string `json:"timestamp"`
	SignMethod   string `json:"signMethod"`
	DeviceSecret string `json:"deviceSecret"`
}

// 腾讯云物联网平台网关
type IThingsGateway struct {
	typex.XStatus
	status            typex.DeviceState
	mainConfig        IThingsGatewayMainConfig
	client            mqtt.Client
	authInfo          IThingsMQTTAuthInfo
	propertyUpTopic   string
	propertyDownTopic string
	actionUpTopic     string
	actionDownTopic   string
	topologyTopicUp   string
	topologyTopicDown string
	//
	IThingsSubDevices []IThingsSubDevice
}

func NewIThingsGateway(e typex.Rhilex) typex.XDevice {
	hd := new(IThingsGateway)
	hd.RuleEngine = e
	hd.mainConfig = IThingsGatewayMainConfig{
		CommonConfig: IThingsGatewayConfig{
			Mode:       "DEVICE",
			ProductId:  "",
			DeviceName: "",
			DevicePsk:  "",
			ClientId:   "",
		},
	}
	hd.IThingsSubDevices = make([]IThingsSubDevice, 0)
	return hd
}

type IThingsMQTTAuthInfo struct {
	ClientID string `json:"clientid"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

func GenerateIThingsMQTTAuthInfo(productID, deviceName, secret string) (IThingsMQTTAuthInfo, error) {
	c, u, p := ithings.GenSecretDeviceInfo("hmacsha256", productID, deviceName, secret)
	return IThingsMQTTAuthInfo{
		ClientID: c,
		UserName: u,
		Password: p,
	}, nil
}

//  初始化
func (hd *IThingsGateway) Init(devId string, configMap map[string]interface{}) error {
	hd.PointId = devId
	if err := utils.BindSourceConfig(configMap, &hd.mainConfig); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	authInfo, err := GenerateIThingsMQTTAuthInfo(hd.mainConfig.CommonConfig.ProductId,
		hd.mainConfig.CommonConfig.DeviceName, hd.mainConfig.CommonConfig.DevicePsk)
	if err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	hd.authInfo = authInfo
	return nil
}

// 启动
func (hd *IThingsGateway) Start(cctx typex.CCTX) error {
	hd.Ctx = cctx.Ctx
	hd.CancelCTX = cctx.CancelCTX
	// 属性
	hd.propertyDownTopic = fmt.Sprintf(_ithings_PropertyTopic,
		hd.mainConfig.CommonConfig.ProductId, hd.mainConfig.CommonConfig.DeviceName)
	hd.propertyUpTopic = fmt.Sprintf(_ithings_PropertyUpTopic,
		hd.mainConfig.CommonConfig.ProductId, hd.mainConfig.CommonConfig.DeviceName)
	hd.actionDownTopic = fmt.Sprintf(_ithings_ActionTopic,
		hd.mainConfig.CommonConfig.ProductId, hd.mainConfig.CommonConfig.DeviceName)
	hd.actionUpTopic = fmt.Sprintf(_ithings_ActionUpTopic,
		hd.mainConfig.CommonConfig.ProductId, hd.mainConfig.CommonConfig.DeviceName)
	hd.topologyTopicUp = fmt.Sprintf(_ithings_operation_up,
		hd.mainConfig.CommonConfig.ProductId, hd.mainConfig.CommonConfig.DeviceName)
	hd.topologyTopicDown = fmt.Sprintf(_ithings_operation_down,
		hd.mainConfig.CommonConfig.ProductId, hd.mainConfig.CommonConfig.DeviceName)

	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		glogger.GLogger.Infof("IOTHUB Connected Success")
		// 属性下发
		if err := hd.client.Subscribe(hd.propertyDownTopic, 1, func(c mqtt.Client, msg mqtt.Message) {
			hd.RuleEngine.WorkDevice(hd.Details(), string(msg.Payload()))
		}); err != nil {
			glogger.GLogger.Error(err)
		}
		// 动作下发
		if err := hd.client.Subscribe(hd.actionDownTopic, 1, func(c mqtt.Client, msg mqtt.Message) {
			hd.RuleEngine.WorkDevice(hd.Details(), string(msg.Payload()))
		}); err != nil {
			glogger.GLogger.Error(err)
		}
		// 网关模式
		//    数据上行 Topic（用于发布）：$gateway/operation/${productid}/${devicename}
		//    数据下行 Topic（用于订阅）：$gateway/operation/result/${productid}/${devicename}
		if hd.mainConfig.CommonConfig.Mode == "GATEWAY" {
			token := hd.client.Subscribe(hd.topologyTopicDown, 1, func(c mqtt.Client, msg mqtt.Message) {
				glogger.GLogger.Debug("IThingsGateway topologyTopicDown: ", hd.topologyTopicDown, msg)
				SubDeviceMessage := IThingsSubDeviceMessage{}
				if err := json.Unmarshal(msg.Payload(), &SubDeviceMessage); err != nil {
					return
				}
				hd.IThingsSubDevices = append(hd.IThingsSubDevices, SubDeviceMessage.Payload.Devices...)
			})
			if token.Error() != nil {
				glogger.GLogger.Error(token.Error())
			} else {
				// Get Topology: {"type": "describe_sub_devices"}
				glogger.GLogger.Info("Connect iothub with Gateway Mode")
				hd.client.Publish(hd.topologyTopicUp, 1, false, `{"type": "describe_sub_devices"}`)
			}
		}
	}

	var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		glogger.GLogger.Warnf("IOTHUB Disconnect: %v, %v try to reconnect", err, hd.status)
	}
	opts := mqtt.NewClientOptions()
	opts.AddBroker(hd.mainConfig.CommonConfig.ServerEndpoint)
	opts.SetClientID(hd.authInfo.ClientID)
	opts.SetUsername(hd.authInfo.UserName)
	opts.SetPassword(hd.authInfo.Password)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.SetCleanSession(true)
	opts.SetPingTimeout(30 * time.Second)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetConnectTimeout(5 * time.Second)
	opts.SetAutoReconnect(false)
	opts.SetMaxReconnectInterval(0)
	hd.client = mqtt.NewClient(opts)
	if token := hd.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	hd.status = typex.DEV_UP
	return nil
}

type IthingsMethod string

const (
	__IThings_ONLINE   IthingsMethod = "online"
	__IThings_OFFLINE  IthingsMethod = "offline"
	__IThings_TOPOLOGY IthingsMethod = "describeSubDevices"
)

type IThingsSubDeviceMessage struct {
	Method  Method                         `json:"method"`
	Payload IThingsSubDeviceMessagePayload `json:"payload"`
}
type IThingsSubDeviceMessagePayload struct {
	Devices []IThingsSubDevice `json:"devices"`
}

// 设备当前状态
func (hd *IThingsGateway) Status() typex.DeviceState {
	if hd.client != nil {
		if hd.client.IsConnectionOpen() && hd.client.IsConnected() {
			return typex.DEV_UP
		} else {
			return typex.DEV_DOWN
		}
	}
	return hd.status
}

// 停止设备
func (hd *IThingsGateway) Stop() {
	hd.status = typex.DEV_DOWN
	if hd.CancelCTX != nil {
		hd.CancelCTX()
	}
	if hd.client != nil {
		hd.client.Disconnect(50)
	}
}

// 真实设备
func (hd *IThingsGateway) Details() *typex.Device {
	return hd.RuleEngine.GetDevice(hd.PointId)
}

// 状态
func (hd *IThingsGateway) SetState(status typex.DeviceState) {
	hd.status = status

}

func (hd *IThingsGateway) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}

func (hd *IThingsGateway) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return []byte{}, nil
}

func (hd *IThingsGateway) OnRead(cmd []byte, data []byte) (int, error) {

	return 0, nil
}

func (hd *IThingsGateway) OnWrite(cmd []byte, b []byte) (int, error) {
	return 0, nil
}

func (hd *IThingsGateway) subscribe(topic string) error {
	token := hd.client.Subscribe(topic, 1, func(c mqtt.Client, msg mqtt.Message) {
		glogger.GLogger.Debug("IThingsGateway: ", topic, msg)
		hd.RuleEngine.WorkDevice(hd.Details(), string(msg.Payload()))
	})
	if token.Error() != nil {
		return token.Error()
	} else {
		return nil
	}
}

// IThingsHubMQTTAuthInfo 腾讯云 Iot Hub MQTT 认证信息
type IThingsHubMQTTAuthInfo struct {
	ClientID string `json:"clientid"`
	UserName string `json:"username"`
	Password string `json:"password"`
}
