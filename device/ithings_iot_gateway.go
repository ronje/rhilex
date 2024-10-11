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
	"github.com/google/uuid"
	"github.com/hootrhino/rhilex/device/ithings"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

// 设备属性上行请求 Topic： $thing/up/property/{ProductID}/{DeviceName}
// 设备属性下行响应 Topic： $thing/down/property/{ProductID}/{DeviceName}
// 设备响应行为执行结果或设备请求服务端行为 Topic： $thing/up/action/{ProductID}/{DeviceName}
// 应用调用设备行为或服务端响应设备请求执行结果 Topic： $thing/down/action/{ProductID}/{DeviceName}
const (
	// 属性
	_ithings_PropertyUpTopic    = "$thing/up/property/%v/%v"
	_ithings_PropertyTopic      = "$thing/down/property/%v/%v"
	_ithings_PropertyReplyTopic = "$thing/down/property/%v/%v"
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
}
type IThingsGatewayMainConfig struct {
	CommonConfig IThingsGatewayConfig `json:"ithingsConfig" validate:"required"` // 通用配置
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
	status             typex.DeviceState
	mainConfig         IThingsGatewayMainConfig
	client             mqtt.Client
	authInfo           IThingsMQTTAuthInfo
	propertyUpTopic    string
	propertyDownTopic  string
	propertyReplyTopic string
	actionUpTopic      string
	actionDownTopic    string
	topologyTopicUp    string
	topologyTopicDown  string
	//
	IThingsSubDevices []IThingsSubDevice
}

func NewIThingsGateway(e typex.Rhilex) typex.XDevice {
	hd := new(IThingsGateway)
	hd.RuleEngine = e
	hd.mainConfig = IThingsGatewayMainConfig{
		CommonConfig: IThingsGatewayConfig{
			ServerEndpoint: "tcp://127.0.0.1:1883",
			Mode:           "DEVICE",
			ProductId:      "",
			DeviceName:     "",
			DevicePsk:      "",
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
	hd.propertyReplyTopic = fmt.Sprintf(_ithings_PropertyReplyTopic,
		hd.mainConfig.CommonConfig.ProductId, hd.mainConfig.CommonConfig.DeviceName)
	// 动作
	hd.actionDownTopic = fmt.Sprintf(_ithings_ActionTopic,
		hd.mainConfig.CommonConfig.ProductId, hd.mainConfig.CommonConfig.DeviceName)
	hd.actionUpTopic = fmt.Sprintf(_ithings_ActionUpTopic,
		hd.mainConfig.CommonConfig.ProductId, hd.mainConfig.CommonConfig.DeviceName)
	// 子设备
	hd.topologyTopicUp = fmt.Sprintf(_ithings_operation_up,
		hd.mainConfig.CommonConfig.ProductId, hd.mainConfig.CommonConfig.DeviceName)
	hd.topologyTopicDown = fmt.Sprintf(_ithings_operation_down,
		hd.mainConfig.CommonConfig.ProductId, hd.mainConfig.CommonConfig.DeviceName)

	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		glogger.GLogger.Infof("IThings Connected Success")
		// 属性下发
		if token := hd.client.Subscribe(hd.propertyDownTopic, 1, func(c mqtt.Client, msg mqtt.Message) {
			glogger.GLogger.Debug("Property Down, Topic: [", msg.Topic(), "] Payload: ", string(msg.Payload()))
			hd.RuleEngine.WorkDevice(hd.Details(), string(msg.Payload()))
		}); token.Error() != nil {
			glogger.GLogger.Error(token.Error())
		}
		// 动作下发
		if token := hd.client.Subscribe(hd.actionDownTopic, 1, func(c mqtt.Client, msg mqtt.Message) {
			glogger.GLogger.Debug("Action Down, Topic: [", msg.Topic(), "] Payload: ", string(msg.Payload()))
			hd.RuleEngine.WorkDevice(hd.Details(), string(msg.Payload()))
		}); token.Error() != nil {
			glogger.GLogger.Error(token.Error())
		}
		// 网关模式:
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
				glogger.GLogger.Info("Connect IThings with Gateway Mode")
				hd.client.Publish(hd.topologyTopicUp, 1, false, `{"type": "describe_sub_devices"}`)
			}
		}
	}

	var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		glogger.GLogger.Warnf("IThings Disconnect: %v, %v try to reconnect", err, hd.status)
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

/**
 * Lua输入进来的指令
 *
 */
type IThingsInputMsg struct {
	Type string      `json:"type"`
	Cmd  interface{} `json:"cmd"`
}

func (hd *IThingsGateway) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return []byte{}, nil
}

func (hd *IThingsGateway) OnRead(cmd []byte, data []byte) (int, error) {

	return 0, nil
}

// ActionReplySuccess
// ActionReplyFailure
// PropertyReplySuccess
// PropertyReplyFailure
// LUA 调用接口
func (hd *IThingsGateway) OnWrite(cmd []byte, b []byte) (int, error) {
	Cmd := string(cmd)
	ActionResp := `{"method": "actionReply","msgToken": "%s","code": 200,"msg":"success"}`
	PropertyResp := `{"method": "reportReply","msgToken": "%s","code": 200,"msg":"success"}`
	if Cmd == "ActionReplySuccess" {
		Token := string(b)
		msg := fmt.Sprintf(ActionResp, Token)
		hd.client.Publish(hd.actionUpTopic, 1, false, msg)
		goto END
	}
	if Cmd == "ActionReplyFailure" {
		Token := string(b)
		msg := fmt.Sprintf(ActionResp, Token)
		hd.client.Publish(hd.actionUpTopic, 1, false, msg)
		goto END
	}
	if Cmd == "PropertyReplySuccess" {
		Token := string(b)
		msg := fmt.Sprintf(PropertyResp, Token)
		hd.client.Publish(hd.propertyUpTopic, 1, false, msg)
		goto END
	}
	if Cmd == "PropertyReplyFailure" {
		Token := string(b)
		msg := fmt.Sprintf(PropertyResp, Token)
		hd.client.Publish(hd.propertyUpTopic, 1, false, msg)
		goto END
	}
	if Cmd == "PropertyReport" {
		params := map[string]interface{}{}
		if errUnmarshal := json.Unmarshal(b, &params); errUnmarshal != nil {
			return 0, errUnmarshal
		}
		IthingsPropertyReport := IthingsPropertyReport{
			Method:    "report",
			MsgToken:  uuid.NewString(),
			Timestamp: time.Now().UnixMilli(),
			Params:    params,
		}
		hd.client.Publish(hd.propertyUpTopic, 1, false, IthingsPropertyReport.String())
		goto END
	}
	if Cmd == "GetPropertyReply" {
		params := map[string]interface{}{}
		if errUnmarshal := json.Unmarshal(b, &params); errUnmarshal != nil {
			return 0, errUnmarshal
		}
		IthingsGetPropertyReply := IthingsGetPropertyReply{
			Method:    "getReportReply",
			Type:      "report",
			MsgToken:  uuid.NewString(),
			Timestamp: time.Now().UnixMilli(),
			Code:      0,
			Data:      params,
			Msg:       "success",
		}
		hd.client.Publish(hd.propertyReplyTopic, 1, false, IthingsGetPropertyReply.String())
		goto END
	}
END:
	return 0, nil
}

type IthingsGetPropertyReply struct {
	Method    string                 `json:"method"`
	Timestamp int64                  `json:"timestamp"`
	MsgToken  string                 `json:"msgToken"`
	Type      string                 `json:"type"`
	Code      int                    `json:"code"`
	Data      map[string]interface{} `json:"data"`
	Msg       string                 `json:"msg"`
}

func (O IthingsGetPropertyReply) String() string {
	bytes, _ := json.Marshal(O)
	return string(bytes)
}

type IthingsPropertyReport struct {
	Method    string                 `json:"method"`
	MsgToken  string                 `json:"msgToken"`
	Timestamp int64                  `json:"timestamp"`
	Params    map[string]interface{} `json:"params"`
}

func (O IthingsPropertyReport) String() string {
	bytes, _ := json.Marshal(O)
	return string(bytes)
}

type IThingsSubDeviceMessage struct {
	Method  string                         `json:"method"`
	Payload IThingsSubDeviceMessagePayload `json:"payload"`
}
type IThingsSubDeviceMessagePayload struct {
	Devices []IThingsSubDevice `json:"devices"`
}

// IThingsHubMQTTAuthInfo 腾讯云 Iot Hub MQTT 认证信息
type IThingsHubMQTTAuthInfo struct {
	ClientID string `json:"clientid"`
	UserName string `json:"username"`
	Password string `json:"password"`
}
