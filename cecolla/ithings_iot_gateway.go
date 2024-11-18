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
	// 网关类型的设备，可通过与云端的数据通信，对其下的子设备进行绑定与解绑操作。实现此类功能需利用如下两个 Topic：
	// 数据上行 Topic（用于发布）：$gateway/up/topo/${productid}/${devicename}
	// 数据下行 Topic（用于订阅）：$gateway/down/topo/${productid}/${devicename}
	_ithings_topology_up   = "$gateway/up/topo/%s/%s"   //数据上行 Topic（用于发布）
	_ithings_topology_down = "$gateway/down/topo/%s/%s" //数据下行 Topic（用于订阅）
	// 网关类型的设备，可通过与云端的数据通信，代理其下的子设备进行上线与下线操作。此类功能所用到的 Topic 与网关子设备拓扑管理的 Topic 一致：
	// 数据上行 Topic（用于发布）：$gateway/up/status/${productid}/${devicename}
	// 数据下行 Topic（用于订阅）：$gateway/down/status/${productid}/${devicename}
	_ithings_gateway_status_up   = "$gateway/up/status/%s/%s"   //数据上行 Topic（用于发布）
	_ithings_gateway_status_down = "$gateway/down/status/%s/%s" //数据下行 Topic（用于订阅）
)

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
	status                 typex.CecollaState
	mainConfig             IThingsGatewayMainConfig
	client                 mqtt.Client
	authInfo               IThingsMQTTAuthInfo
	propertyUpTopic        string
	propertyDownTopic      string
	actionUpTopic          string
	actionDownTopic        string
	gatewayTopicUp         string
	gatewayTopicDown       string
	topologyTopicUp        string
	topologyTopicDown      string
	gatewayStatusTopicUp   string
	gatewayStatusTopicDown string
	// 子设备
	IThingsSubDevices []ithings.IThingsSubDevice
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
	hd.IThingsSubDevices = make([]ithings.IThingsSubDevice, 0)
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
	// 拓扑关系
	hd.topologyTopicUp = fmt.Sprintf(_ithings_topology_up,
		hd.mainConfig.ProductId, hd.mainConfig.DeviceName)
	hd.topologyTopicDown = fmt.Sprintf(_ithings_topology_down,
		hd.mainConfig.ProductId, hd.mainConfig.DeviceName)
	// 子设备上下线
	hd.gatewayStatusTopicUp = fmt.Sprintf(_ithings_gateway_status_up,
		hd.mainConfig.ProductId, hd.mainConfig.DeviceName)
	hd.gatewayStatusTopicDown = fmt.Sprintf(_ithings_gateway_status_down,
		hd.mainConfig.ProductId, hd.mainConfig.DeviceName)
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
			glogger.GLogger.Info("Connect IThings with Gateway Mode")
			token0 := hd.client.Subscribe(hd.gatewayTopicDown, 1, func(c mqtt.Client, msg mqtt.Message) {
				glogger.Debug("IThingsGateway gateway Down: ", hd.gatewayTopicDown, string(msg.Payload()))
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
				}
			})
			if token0.Error() == nil {
				// 获取物模型
				hd.client.Publish(hd.gatewayTopicUp, 1, false,
					fmt.Sprintf(`{"method":"getSchema","msgToken":"%s","payload":{"productID":"%s"}}`,
						uuid.NewString(), hd.mainConfig.SubProduct))
			}

			// 获取拓扑
			token1 := hd.client.Subscribe(hd.topologyTopicDown, 1, func(c mqtt.Client, msg mqtt.Message) {
				glogger.Debug("IThingsGateway topology Down: ", hd.gatewayTopicDown, string(msg.Payload()))
				response := ithings.IthingsTopologyResponse{}
				errUnmarshal := json.Unmarshal(msg.Payload(), &response)
				if errUnmarshal != nil {
					glogger.GLogger.Error(errUnmarshal)
					return
				}
				hd.IThingsSubDevices = response.Payload.Devices
			})
			if token1.Error() == nil {
				hd.client.Publish(hd.topologyTopicUp, 1, false,
					fmt.Sprintf(`{"method": "getTopo","msgToken": "%s"}`, uuid.NewString()))
			}

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

type SubDeviceParam struct {
	Timestamp int64  `json:"timestamp"`
	ProductId string `json:"productID"`
	DeviceId  string `json:"deviceID"`
	Param     string `json:"param"`
	Value     any    `json:"value"`
}

/**
 * 外部参数
 *
 */
func (hd *IThingsGateway) OnCtrl(cmd []byte, b []byte) (any, error) {
	Cmd := string(cmd)
	// 子设备上线
	if Cmd == "SubDeviceSetOnline" {
		subDeviceParam := SubDeviceParam{}
		if errUnmarshal := json.Unmarshal(b, &subDeviceParam); errUnmarshal != nil {
			return nil, errUnmarshal
		}
		payload := `{"method":"online","msgToken":"%s","payload":{"devices":[{"productID":"%s","deviceName":"%s"}]}}`
		token := hd.client.Publish(hd.gatewayStatusTopicUp, 1, false,
			fmt.Sprintf(payload, uuid.NewString(), subDeviceParam.ProductId, subDeviceParam.DeviceId))
		glogger.Debugf("SubDevice SetOnline: %s %s", hd.gatewayStatusTopicUp,
			fmt.Sprintf(payload, uuid.NewString(), subDeviceParam.ProductId, subDeviceParam.DeviceId))

		if token.Error() != nil {
			glogger.Error(token.Error())
			return nil, token.Error()
		}
	}
	// 来自点位表的数据同步，批量上报
	if Cmd == "PackReportSubDeviceParams" {
		subDeviceParam := SubDeviceParam{}
		if errUnmarshal := json.Unmarshal(b, &subDeviceParam); errUnmarshal != nil {
			return nil, errUnmarshal
		}
		packReport := ithings.NewIthingsPackReport(subDeviceParam.Timestamp,
			subDeviceParam.ProductId, subDeviceParam.DeviceId,
			subDeviceParam.Param, subDeviceParam.Value)
		glogger.Debugf("PackReport SubDevice Params: %s %s", hd.propertyUpTopic, packReport.String())
		token := hd.client.Publish(hd.propertyUpTopic, 1, false, packReport.String())
		if token.Error() != nil {
			glogger.Error(token.Error())
			return nil, token.Error()
		}
		return nil, nil
	}
	// 返回物模型
	if Cmd == "GetSchema" {
		glogger.Debug("GetSchema")
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
		IthingsPropertyReport := ithings.IthingsPropertyReport{
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
		IthingsGetPropertyReply := ithings.IthingsGetPropertyReply{
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
