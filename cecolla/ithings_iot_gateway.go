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

package cecolla

import (
	"encoding/json"
	"fmt"

	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/hootrhino/rhilex/cecolla/ithings"
	"github.com/hootrhino/rhilex/component/intercache"
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
	_ithings_PropertyUpTopic = "$thing/up/property/%v/%v"
	_ithings_PropertyTopic   = "$thing/down/property/%v/%v"
	// 动作
	_ithings_ActionTopic   = "$thing/down/action/%v/%v"
	_ithings_ActionUpTopic = "$thing/up/action/%v/%v"
	// 设备从云端接收最新消息使用的 Topic：
	//     请求 Topic： $gateway/up/thing/{ProductID}/{DeviceName}
	//     响应 Topic： $gateway/down/thing/{ProductID}/{DeviceName}
	_ithings_gateway_up   = "$gateway/up/thing/%s/%s"   //数据上行 Topic（用于发布）
	_ithings_gateway_down = "$gateway/down/thing/%s/%s" //数据下行 Topic（用于订阅）
)

type IThingsSubDevice struct {
	ProductID    string `json:"productID"`
	DeviceName   string `json:"deviceName"`
	Signature    string `json:"signature"`
	Random       string `json:"random"`
	Timestamp    string `json:"timestamp"`
	SignMethod   string `json:"signMethod"`
	DeviceSecret string `json:"deviceSecret"`
}

/**
 * 主配置
 *
 */
type IThingsGatewayMainConfig struct {
	ServerEndpoint string `json:"serverEndpoint" validate:"required"` //服务地址，默认"tcp://127.0.0.1:1883"
	Mode           string `json:"mode" validate:"required"`           //工作模式: DEVICE|GATEWAY，默认DEVICE
	SubProduct     string `json:"subProduct" validate:"required"`     //子产品名,GATEWAY模式下有效，必填
	ProductId      string `json:"productId" validate:"required"`      //产品名称
	DeviceName     string `json:"deviceName" validate:"required"`     //设备名称
	DevicePsk      string `json:"devicePsk" validate:"required"`      //连接秘钥
}

// 腾讯云物联网平台网关
type IThingsGateway struct {
	typex.XStatus
	status            typex.CecollaState
	mainConfig        IThingsGatewayMainConfig
	client            mqtt.Client
	authInfo          IThingsMQTTAuthInfo
	propertyUpTopic   string
	propertyDownTopic string
	actionUpTopic     string
	actionDownTopic   string
	gatewayTopicUp    string
	gatewayTopicDown  string
	// 子设备
	IThingsSubDevices []IThingsSubDevice
	// 自己的物模型
	GatewaySchema *ithings.SchemaSimple
	// 子设备的物模型
	SubDeviceSchema *ithings.SchemaSimple
}

func NewIThingsGateway(e typex.Rhilex) typex.XCecolla {
	hd := new(IThingsGateway)
	hd.RuleEngine = e
	hd.mainConfig = IThingsGatewayMainConfig{
		ServerEndpoint: "tcp://127.0.0.1:1883",
		Mode:           "DEVICE",
		SubProduct:     "",
		ProductId:      "",
		DeviceName:     "",
		DevicePsk:      "",
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
	intercache.RegisterSlot(hd.PointId)
	if err := utils.BindSourceConfig(configMap, &hd.mainConfig); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	authInfo, err := GenerateIThingsMQTTAuthInfo(hd.mainConfig.ProductId,
		hd.mainConfig.DeviceName, hd.mainConfig.DevicePsk)
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
	// 自身属性
	hd.propertyDownTopic = fmt.Sprintf(_ithings_PropertyTopic,
		hd.mainConfig.ProductId, hd.mainConfig.DeviceName)
	hd.propertyUpTopic = fmt.Sprintf(_ithings_PropertyUpTopic,
		hd.mainConfig.ProductId, hd.mainConfig.DeviceName)
	// 自身动作
	hd.actionDownTopic = fmt.Sprintf(_ithings_ActionTopic,
		hd.mainConfig.ProductId, hd.mainConfig.DeviceName)
	hd.actionUpTopic = fmt.Sprintf(_ithings_ActionUpTopic,
		hd.mainConfig.ProductId, hd.mainConfig.DeviceName)

	// 网关-子设备
	hd.gatewayTopicUp = fmt.Sprintf(_ithings_gateway_up,
		hd.mainConfig.ProductId, hd.mainConfig.DeviceName)
	hd.gatewayTopicDown = fmt.Sprintf(_ithings_gateway_down,
		hd.mainConfig.ProductId, hd.mainConfig.DeviceName)
	//
	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		glogger.GLogger.Infof("IThings Connected Success")
		// 属性下发
		if token := hd.client.Subscribe(hd.propertyDownTopic, 1, func(c mqtt.Client, msg mqtt.Message) {
			glogger.GLogger.Debug("Property Down, Topic: [", msg.Topic(), "] Payload: ", string(msg.Payload()))
		}); token.Error() != nil {
			glogger.GLogger.Error(token.Error())
		}
		// 动作下发
		if token := hd.client.Subscribe(hd.actionDownTopic, 1, func(c mqtt.Client, msg mqtt.Message) {
			glogger.GLogger.Debug("Action Down, Topic: [", msg.Topic(), "] Payload: ", string(msg.Payload()))
		}); token.Error() != nil {
			glogger.GLogger.Error(token.Error())
		}
		// 设备物模型
		hd.client.Publish(hd.gatewayTopicUp, 1, false,
			fmt.Sprintf(`{"method":"getSchema","msgToken":"%s","payload":{"productID":"%s"}}`,
				uuid.NewString(), hd.mainConfig.ProductId))
		// 获取子物模型
		if hd.mainConfig.Mode == "GATEWAY" {
			token := hd.client.Subscribe(hd.gatewayTopicDown, 1, func(c mqtt.Client, msg mqtt.Message) {
				glogger.Debug("IThingsGateway gateway Topic Down: ", hd.gatewayTopicDown, string(msg.Payload()))
				response := ithings.IthingsResponse{}
				errUnmarshal := json.Unmarshal(msg.Payload(), &response)
				if errUnmarshal != nil {
					glogger.GLogger.Error(errUnmarshal)
					return
				}
				if response.Payload.ProductId == hd.mainConfig.ProductId {
					hd.GatewaySchema = &response.Payload.Schema
					glogger.Debug("Get Gateway Schema Success:", hd.GatewaySchema.String())
				}
				if response.Payload.ProductId == hd.mainConfig.SubProduct {
					hd.SubDeviceSchema = &response.Payload.Schema
					glogger.Debug("Get SubDevice Schema Success:", hd.SubDeviceSchema.String())
				}
			})
			if token.Error() != nil {
				glogger.GLogger.Error(token.Error())
				return
			}
			glogger.GLogger.Info("Connect IThings with Gateway Mode")
			hd.client.Publish(hd.gatewayTopicUp, 1, false,
				fmt.Sprintf(`{"method":"getSchema","msgToken":"%s","payload":{"productID":"%s"}}`,
					uuid.NewString(), hd.mainConfig.SubProduct))
		}
	}

	var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		glogger.GLogger.Warnf("IThings Disconnect: %v, %v try to reconnect", err, hd.status)
	}
	opts := mqtt.NewClientOptions()
	opts.AddBroker(hd.mainConfig.ServerEndpoint)
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
	hd.status = typex.CEC_UP
	return nil
}

// 设备当前状态
func (hd *IThingsGateway) Status() typex.CecollaState {
	if hd.client != nil {
		if hd.client.IsConnectionOpen() && hd.client.IsConnected() {
			return typex.CEC_UP
		} else {
			return typex.CEC_DOWN
		}
	}
	return hd.status
}

// 停止设备
func (hd *IThingsGateway) Stop() {
	intercache.UnRegisterSlot(hd.PointId)
	hd.status = typex.CEC_DOWN
	if hd.CancelCTX != nil {
		hd.CancelCTX()
	}
	if hd.client != nil {
		hd.client.Disconnect(50)
	}
}

// 真实设备
func (hd *IThingsGateway) Details() *typex.Cecolla {
	return hd.RuleEngine.GetCecolla(hd.PointId)
}

// 状态
func (hd *IThingsGateway) SetState(status typex.CecollaState) {
	hd.status = status

}

/**
 * Lua输入进来的指令
 *
 */
type IThingsInputMsg struct {
	Type string      `json:"type"`
	Cmd  interface{} `json:"cmd"`
}

// CtrlReplySuccess
// CtrlReplyFailure
// ActionReplySuccess
// ActionReplyFailure
// PropertyReplySuccess
// PropertyReplyFailure
func (hd *IThingsGateway) OnCtrl(cmd []byte, b []byte) (any, error) {
	Cmd := string(cmd)
	// 返回物模型
	if Cmd == "GetSchema" {
		return map[string]any{
			"gatewaySchema":   hd.GatewaySchema,
			"subDeviceSchema": hd.SubDeviceSchema,
		}, nil
	}
	PropertyResp := `{"method": "reportReply","msgToken": "%s","code": 200,"msg":"success"}`
	CtrlResp := `{"method": "controlReply","msgToken": "%s","code": 200,"msg":"success"}`
	ActionResp := `{"method": "actionReply","msgToken": "%s","code": 200,"msg":"success"}`
	if Cmd == "CtrlReplySuccess" {
		Token := string(b)
		msg := fmt.Sprintf(CtrlResp, Token)
		hd.client.Publish(hd.propertyUpTopic, 1, false, msg)
		goto END
	}
	if Cmd == "CtrlReplyFailure" {
		Token := string(b)
		msg := fmt.Sprintf(CtrlResp, Token)
		hd.client.Publish(hd.propertyUpTopic, 1, false, msg)
		goto END
	}
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
			return nil, errUnmarshal
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
			return nil, errUnmarshal
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
		hd.client.Publish(hd.propertyUpTopic, 1, false, IthingsGetPropertyReply.String())
		goto END
	}
END:
	return nil, nil
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
