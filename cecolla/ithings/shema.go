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
	"time"

	"github.com/google/uuid"
)

type (
	// Schema 物模型协议-数据模板定义
	SchemaSimple struct {
		Properties PropertiesSimple `json:"properties"` //属性
		Events     EventsSimple     `json:"events"`     //事件
		Actions    ActionsSimple    `json:"actions"`    //行为
	}

	/*事件*/
	EventSimple struct {
		Identifier string       `json:"id"`     //标识符 (统一)
		Name       string       `json:"name"`   //功能名称
		Type       EventType    `json:"type"`   //事件类型: 1:信息:info  2:告警alert  3:故障:fault
		Params     ParamSimples `json:"params"` //事件参数
	}
	EventsSimple []EventSimple

	ParamSimple struct {
		Identifier string `json:"id"`   //参数标识符
		Name       string `json:"name"` //参数名称
		Define            //参数定义
	}
	ParamSimples []ParamSimple
	/*行为*/
	ActionSimple struct {
		Identifier string       `json:"id"`     //标识符 (统一)
		Name       string       `json:"name"`   //功能名称
		Dir        ActionDir    `json:"dir"`    //调用方向
		Input      ParamSimples `json:"input"`  //调用参数
		Output     ParamSimples `json:"output"` //返回参数
	}
	ActionsSimple []ActionSimple

	/*属性*/
	PropertySimple struct {
		Identifier string       `json:"id"`   //标识符 (统一)
		Name       string       `json:"name"` //功能名称
		Mode       PropertyMode `json:"mode"` //读写类型:rw(可读可写) r(只读)
		Define                  //数据定义
	}
	PropertiesSimple []PropertySimple
	/*数据类型定义*/
	Define struct {
		Type    DataType          `json:"type"`    //参数类型:bool int string struct float timestamp array enum
		Mapping map[string]string `json:"mapping"` //枚举及bool类型:bool enum
	}
)

func (O SchemaSimple) String() string {
	bytes, _ := json.Marshal(O)
	return string(bytes)
}

// 数据类型
type DataType string

const (
	DataTypeBool      DataType = "bool"
	DataTypeInt       DataType = "int"
	DataTypeString    DataType = "string"
	DataTypeStruct    DataType = "struct"
	DataTypeFloat     DataType = "float"
	DataTypeTimestamp DataType = "timestamp"
	DataTypeArray     DataType = "array"
	DataTypeEnum      DataType = "enum"
)

// 属性读写类型: r(只读) rw(可读可写)
type PropertyMode string

const (
	PropertyModeR  PropertyMode = "r"
	PropertyModeRW PropertyMode = "rw"
)

// 事件类型: 信息:info  告警alert  故障:fault
type EventType = string

const (
	EventTypeInfo  EventType = "info"
	EventTypeAlert EventType = "alert"
	EventTypeFault EventType = "fault"
)

// 行为的执行方向
type ActionDir string

const (
	ActionDirUp   ActionDir = "up"   //向上调用
	ActionDirDown ActionDir = "down" //向下调用
)

type IthingsGetPropertyReply struct {
	Method    string                 `json:"method"`
	Timestamp int64                  `json:"timestamp"`
	MsgToken  string                 `json:"msgToken"`
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

type IthingsResponse struct {
	Method   string         `json:"method"`
	MsgToken string         `json:"msgToken"`
	Code     int            `json:"code"`
	Payload  IthingsPayload `json:"payload"`
}

type IthingsPayload struct {
	ProductId string       `json:"productId"`
	Schema    SchemaSimple `json:"schema"`
}

func (O IthingsResponse) String() string {
	bytes, _ := json.Marshal(O)
	return string(bytes)
}

type IthingsTopologyResponse struct {
	Method   string                  `json:"method"`
	MsgToken string                  `json:"msgToken"`
	Code     int                     `json:"code"`
	Payload  IthingsSubDevicePayload `json:"payload"`
}

func (O IthingsTopologyResponse) String() string {
	bytes, _ := json.Marshal(O)
	return string(bytes)
}

type IthingsSubDevicePayload struct {
	Devices []IThingsSubDevice `json:"devices"`
}
type IThingsSubDevice struct {
	ProductID  string `json:"productID"`
	DeviceName string `json:"deviceName"`
	// Signature    string `json:"signature"`
	// Random       string `json:"random"`
	// Timestamp    string `json:"timestamp"`
	// SignMethod   string `json:"signMethod"`
	// DeviceSecret string `json:"deviceSecret"`
}

/**
 * 批量上报子设备
 *
 */
type IthingsPackReport struct {
	Method     string              `json:"method"`
	MsgToken   string              `json:"msgToken"`
	Timestamp  int64               `json:"timestamp"`
	Properties []IthingsProperties `json:"properties,omitempty"`
	SubDevices []IthingsSubDevices `json:"subDevices,omitempty"`
}

func NewIthingsPackReport(Timestamp int64, ProductID, DeviceName string, Param string, Value any) IthingsPackReport {
	return IthingsPackReport{
		Method:    "packReport",
		MsgToken:  uuid.NewString(),
		Timestamp: time.Now().UnixMilli(),
		SubDevices: []IthingsSubDevices{
			{
				ProductID:  ProductID,
				DeviceName: DeviceName,
				Properties: []IthingsProperties{
					{
						Timestamp: Timestamp,
						Params:    map[string]any{Param: Value},
					},
				},
			},
		},
	}
}
func (O IthingsPackReport) String() string {
	bytes, _ := json.Marshal(O)
	return string(bytes)
}

type IthingsSubDevices struct {
	ProductID  string              `json:"productID"`
	DeviceName string              `json:"deviceName"`
	Properties []IthingsProperties `json:"properties"`
}

type IthingsProperties struct {
	Timestamp int64          `json:"timestamp,omitempty"`
	Params    map[string]any `json:"params"`
}

/**
 * 创建物模型
 *
 */
type IthingsCreateSchema struct {
	Method     string                         `json:"method"`
	MsgToken   string                         `json:"msgToken"`
	Timestamp  int64                          `json:"timestamp"`
	Properties []IthingsCreateSchemaPropertie `json:"properties"`
}

func (O IthingsCreateSchema) String() string {
	bytes, _ := json.Marshal(O)
	return string(bytes)
}

type IthingsCreateSchemaPropertie struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	ProductId  string `json:"productID,omitempty"`
	DeviceName string `json:"deviceName,omitempty"`
}
