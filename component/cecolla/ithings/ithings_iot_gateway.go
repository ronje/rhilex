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

package ithings

import (
	"encoding/json"
	"fmt"
	"strings"

	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rhilex/component/cecollalet"
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
	_ithings_PropertyUpTopic   = "$thing/up/property/%v/%v"
	_ithings_PropertyDownTopic = "$thing/down/property/%v/%v"
	// 动作
	_ithings_ActionDownTopic = "$thing/down/action/%v/%v"
	_ithings_ActionUpTopic   = "$thing/up/action/%v/%v"
	// 设备从云端接收最新消息使用的 Topic：
	//     请求 Topic： $gateway/up/thing/{ProductID}/{devicename}
	//     响应 Topic： $gateway/down/thing/{ProductID}/{devicename}
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
	Cecollalet             *cecollalet.Cecollalet
	Action                 string
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
	IThingsSubDevices []IThingsSubDevice
	// 自己的物模型
	GatewaySchema *SchemaSimple
	// 子设备的物模型
	SubDeviceSchema *SchemaSimple
	// 缓存数据值的地方，设备report过来以后，保存在此处
	DevicePropertiesSlot map[string]any
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
	hd.IThingsSubDevices = []IThingsSubDevice{}
	hd.DevicePropertiesSlot = map[string]any{}
	return hd
}

type IThingsMQTTAuthInfo struct {
	ClientID string `json:"clientid"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

func GenerateIThingsMQTTAuthInfo(productID, DeviceName, secret string) (IThingsMQTTAuthInfo, error) {
	c, u, p := GenSecretDeviceInfo("hmacsha256", productID, DeviceName, secret)
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
	hd.Cecollalet = cecollalet.NewCecollalet(hd.PointId, fmt.Sprintf("Action: %s", devId), "v1.0.0")
	// FIXME 这种形式来获取配置其实是不合理的，不过暂时先这么搞，后期重新设计架构
	Detail := hd.RuleEngine.GetCecolla(hd.PointId)
	if Detail != nil {
		return cecollalet.LoadCecollalet(hd.Cecollalet, Detail.Action)
	}
	return fmt.Errorf("Get Cecolla Failed:%s", hd.PointId)
}

// 启动
func (hd *IThingsGateway) Start(cctx typex.CCTX) error {
	hd.Ctx = cctx.Ctx
	hd.CancelCTX = cctx.CancelCTX
	// 自身属性
	hd.propertyDownTopic = fmt.Sprintf(_ithings_PropertyDownTopic,
		hd.mainConfig.ProductId, hd.mainConfig.DeviceName)
	hd.propertyUpTopic = fmt.Sprintf(_ithings_PropertyUpTopic,
		hd.mainConfig.ProductId, hd.mainConfig.DeviceName)
	// 自身动作
	hd.actionDownTopic = fmt.Sprintf(_ithings_ActionDownTopic,
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
		glogger.Infof("IThings Connected Success")
		// 属性下发
		if token := hd.client.Subscribe(hd.propertyDownTopic, 1,
			hd.OnGatewayPropertyReceived); token.Error() != nil {
			glogger.Error(token.Error())
		}
		// 动作下发
		if token := hd.client.Subscribe(hd.actionDownTopic, 1,
			hd.OnGatewayActionReceived); token.Error() != nil {
			glogger.Error(token.Error())
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
				response := IthingsResponse{}
				errUnmarshal := json.Unmarshal(msg.Payload(), &response)
				if errUnmarshal != nil {
					glogger.GLogger.Error(errUnmarshal)
					return
				}
				// DataType = "bool"
				// DataType = "int"
				// DataType = "string"
				// DataType = "struct"
				// DataType = "float"
				// 初始化本地属性槽
				if response.Payload.ProductId == hd.mainConfig.ProductId {
					hd.GatewaySchema = &response.Payload.Schema
					for _, Property := range hd.GatewaySchema.Properties {
						if Property.Type == "bool" {
							hd.DevicePropertiesSlot[fmt.Sprintf("%s:%s:%s",
								hd.mainConfig.ProductId, hd.mainConfig.DeviceName, Property.Identifier)] = false
						}
						if Property.Type == "int" {
							hd.DevicePropertiesSlot[fmt.Sprintf("%s:%s:%s",
								hd.mainConfig.ProductId, hd.mainConfig.DeviceName, Property.Identifier)] = 0
						}
						if Property.Type == "string" {
							hd.DevicePropertiesSlot[fmt.Sprintf("%s:%s:%s",
								hd.mainConfig.ProductId, hd.mainConfig.DeviceName, Property.Identifier)] = ""
						}
						if Property.Type == "float" {
							hd.DevicePropertiesSlot[fmt.Sprintf("%s:%s:%s",
								hd.mainConfig.ProductId, hd.mainConfig.DeviceName, Property.Identifier)] = 0
						}
						if Property.Type == "struct" {
							hd.DevicePropertiesSlot[fmt.Sprintf("%s:%s:%s",
								hd.mainConfig.ProductId, hd.mainConfig.DeviceName, Property.Identifier)] = nil
						}

					}
					glogger.Debug("Get Gateway Schema Success:", hd.GatewaySchema.String())
				}
				// 初始化子设备的属性槽
				if response.Payload.ProductId == hd.mainConfig.SubProduct {
					hd.SubDeviceSchema = &response.Payload.Schema
					for _, Property := range hd.SubDeviceSchema.Properties {
						if Property.Type == "bool" {
							hd.DevicePropertiesSlot[fmt.Sprintf("%s:%s:%s",
								hd.mainConfig.ProductId, hd.mainConfig.DeviceName, Property.Identifier)] = false
						}
						if Property.Type == "int" {
							hd.DevicePropertiesSlot[fmt.Sprintf("%s:%s:%s",
								hd.mainConfig.ProductId, hd.mainConfig.DeviceName, Property.Identifier)] = 0
						}
						if Property.Type == "string" {
							hd.DevicePropertiesSlot[fmt.Sprintf("%s:%s:%s",
								hd.mainConfig.ProductId, hd.mainConfig.DeviceName, Property.Identifier)] = ""
						}
						if Property.Type == "float" {
							hd.DevicePropertiesSlot[fmt.Sprintf("%s:%s:%s",
								hd.mainConfig.ProductId, hd.mainConfig.DeviceName, Property.Identifier)] = 0
						}
						if Property.Type == "struct" {
							hd.DevicePropertiesSlot[fmt.Sprintf("%s:%s:%s",
								hd.mainConfig.ProductId, hd.mainConfig.DeviceName, Property.Identifier)] = nil
						}
					}
					glogger.Debug("Get SubDevice Schema Success:", hd.SubDeviceSchema.String())
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
				response := IthingsTopologyResponse{}
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
		if !hd.client.IsConnectionOpen() || !hd.client.IsConnected() {
			return typex.CEC_DOWN
		}
	}
	return hd.status
}

// 停止设备
func (hd *IThingsGateway) Stop() {
	intercache.UnRegisterSlot(hd.PointId)
	if hd.Cecollalet != nil {
		cecollalet.StopCecollalet(hd.PointId)
		cecollalet.RemoveCecollalet(hd.PointId)
	}
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
	Timestamp  int64  `json:"timestamp"`
	ProductId  string `json:"productID"`
	DeviceName string `json:"deviceName"`
	Param      string `json:"param"`
	Value      any    `json:"value"`
}

// 子设备订阅
type SubDeviceTopic struct {
	ProductId  string `json:"productID"`
	DeviceName string `json:"deviceName"`
}

/**
 * 获取属性
 *
 */
//
type GetPropertiesCmd struct {
	Token       string   `json:"token"`
	ProductId   string   `json:"productID"`
	DeviceName  string   `json:"deviceName"`
	Identifiers []string `json:"identifiers"`
}

/**
 * 外部参数
 *
 */
func (hd *IThingsGateway) OnCtrl(cmd []byte, b []byte) (any, error) {
	if hd.client == nil {
		return nil, fmt.Errorf("invalid mqtt connection")
	}
	Cmd := string(cmd)
	// 获取某个产品某个设备的某个属性,用于子设备数据获取
	if Cmd == "GetProperties" {
		ctrlCmd := GetPropertiesCmd{}
		if errUnmarshal := json.Unmarshal(b, &ctrlCmd); errUnmarshal != nil {
			return nil, errUnmarshal
		}
		values := map[string]any{}
		for _, param := range ctrlCmd.Identifiers {
			value := hd.DevicePropertiesSlot[fmt.Sprintf("%s:%s:%s", ctrlCmd.ProductId, ctrlCmd.DeviceName, param)]
			values[param] = value
		}
		if bytes, errMarshal := json.Marshal(values); errMarshal != nil {
			return nil, errMarshal
		} else {
			return bytes, nil
		}
	}
	// 新建子设备点位表的物模型
	if Cmd == "CreateSubDeviceSchema" {
		Properties := []IthingsCreateSchemaPropertie{}
		if errUnmarshal := json.Unmarshal(b, &Properties); errUnmarshal != nil {
			return nil, errUnmarshal
		}
		for _, Property := range Properties {
			createSchema := IthingsCreateSchema{
				Method:    "createSchema",
				MsgToken:  uuid.NewString(),
				Timestamp: time.Now().UnixMilli(),
				Properties: []IthingsCreateSchemaPropertie{{
					Id:   Property.Id,
					Name: Property.Name,
					Type: Property.Type,
				}},
			}
			glogger.Debug("Create SubDevice Schema:", createSchema.String())
			token := hd.client.Publish(fmt.Sprintf(_ithings_gateway_up,
				Property.ProductId, Property.DeviceName), 1, false, createSchema.String())
			if token.Error() != nil {
				glogger.Error(token.Error())
				return nil, token.Error()
			}
			// 在物模型创建好以后初始化值
			hd.DevicePropertiesSlot[fmt.Sprintf("%s:%s:%s",
				hd.mainConfig.ProductId, hd.mainConfig.DeviceName, Property.Name)] = new(int)
		}

		// 再次刷新获取物模型
		hd.client.Publish(hd.gatewayTopicUp, 1, false,
			fmt.Sprintf(`{"method":"getSchema","msgToken":"%s","payload":{"productID":"%s"}}`,
				uuid.NewString(), hd.mainConfig.SubProduct))
	}
	// 子设备上线
	if Cmd == "SubDeviceSetOnline" {
		subDeviceParam := SubDeviceParam{}
		if errUnmarshal := json.Unmarshal(b, &subDeviceParam); errUnmarshal != nil {
			return nil, errUnmarshal
		}
		payload := `{"method":"online","msgToken":"%s","payload":{"devices":[{"productID":"%s","deviceName":"%s"}]}}`
		token := hd.client.Publish(hd.gatewayStatusTopicUp, 1, false,
			fmt.Sprintf(payload, uuid.NewString(), subDeviceParam.ProductId, subDeviceParam.DeviceName))
		glogger.Debugf("SubDevice SetOnline: %s %s", hd.gatewayStatusTopicUp,
			fmt.Sprintf(payload, uuid.NewString(), subDeviceParam.ProductId, subDeviceParam.DeviceName))

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
		hd.DevicePropertiesSlot[fmt.Sprintf("%s:%s:%s",
			subDeviceParam.ProductId, subDeviceParam.DeviceName, subDeviceParam.Param)] = subDeviceParam.Value
		packReport := NewIthingsPackReport(subDeviceParam.Timestamp,
			subDeviceParam.ProductId, subDeviceParam.DeviceName,
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
	// 上报属性
	if Cmd == "GetPropertyReplySuccess" {
		ctrlCmd := GetPropertiesCmd{}
		if errUnmarshal := json.Unmarshal(b, &ctrlCmd); errUnmarshal != nil {
			return nil, errUnmarshal
		}
		params := map[string]any{}
		for _, param := range ctrlCmd.Identifiers {
			value := hd.DevicePropertiesSlot[fmt.Sprintf("%s:%s:%s", ctrlCmd.ProductId, ctrlCmd.DeviceName, param)]
			params[param] = value
		}
		IthingsGetPropertyReply := IthingsGetPropertyReply{
			Method:    "getReportReply",
			MsgToken:  ctrlCmd.Token,
			Timestamp: time.Now().UnixMilli(),
			Code:      200,
			Data:      params,
			Msg:       "success",
		}
		glogger.Debug("IthingsGetPropertyReply:", IthingsGetPropertyReply.String())
		hd.client.Publish(fmt.Sprintf(_ithings_PropertyUpTopic, ctrlCmd.ProductId, ctrlCmd.DeviceName),
			1, false, IthingsGetPropertyReply.String())
		goto END
	}
END:
	return nil, nil
}

/**
 * 执行本地回调
 *
 */
type IthingsDownMsg struct {
	Method   string `json:"method"`
	MsgToken string `json:"msgToken"`
}

/**
 * 处理网关的物模型消息
 *
 */
func (hd *IThingsGateway) OnGatewayPropertyReceived(c mqtt.Client, msg mqtt.Message) {
	glogger.Debug("IThingsGateway.OnGatewayPropertyReceived == ", string(msg.Payload()))
	downMsg := IthingsDownMsg{}
	if errUnmarshal := json.Unmarshal(msg.Payload(), &downMsg); errUnmarshal != nil {
		glogger.Error(errUnmarshal)
		return
	}
	if downMsg.Method == "control" || downMsg.Method == "getReport" {
		// "$thing/down/action/%v/%v"
		fields := strings.Split(msg.Topic(), "/")
		Env := lua.LTable{}
		if len(fields) == 5 {
			Env.RawSetString("Product", lua.LString(fields[3]))
			Env.RawSetString("Device", lua.LString(fields[4]))
			Env.RawSetString("Payload", lua.LString(msg.Payload()))
		}
		err := cecollalet.StartCecollalet(hd.PointId, &Env)
		if err != nil {
			glogger.Error(err)
		}
	}
}
func (hd *IThingsGateway) OnGatewayActionReceived(c mqtt.Client, msg mqtt.Message) {
	glogger.Debug("IThingsGateway.OnGatewayActionReceived == ", string(msg.Payload()))
	downMsg := IthingsDownMsg{}
	if errUnmarshal := json.Unmarshal(msg.Payload(), &downMsg); errUnmarshal != nil {
		return
	}
	if downMsg.Method == "action" {
		// "$thing/down/action/%v/%v"
		fields := strings.Split(msg.Topic(), "/")
		Env := lua.LTable{}
		if len(fields) == 5 {
			Env.RawSetString("Product", lua.LString(fields[3]))
			Env.RawSetString("Device", lua.LString(fields[4]))
			Env.RawSetString("Payload", lua.LString(msg.Payload()))
		}
		err := cecollalet.StartCecollalet(hd.PointId, &Env)
		if err != nil {
			glogger.Error(err)
		}
	}

}

/**
 * 子设备的代理消息
 *
 */
func (hd *IThingsGateway) OnSubdevicePropertyReceived(c mqtt.Client, msg mqtt.Message) {
	glogger.Debug("IThingsGateway.OnSubdevicePropertyReceived == ", msg)

}
func (hd *IThingsGateway) OnSubdeviceActionReceived(c mqtt.Client, msg mqtt.Message) {
	glogger.Debug("IThingsGateway.OnSubdeviceActionReceived == ", msg)

}
