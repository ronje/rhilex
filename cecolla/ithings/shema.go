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

type (
	// Model 物模型协议-数据模板定义
	ModelSimple struct {
		Properties PropertiesSimple `json:"properties,omitempty"` //属性
		Events     EventsSimple     `json:"events,omitempty"`     //事件
		Actions    ActionsSimple    `json:"actions,omitempty"`    //行为
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
		Type    DataType          `json:"type"`              //参数类型:bool int string struct float timestamp array enum
		Mapping map[string]string `json:"mapping,omitempty"` //枚举及bool类型:bool enum
	}
)

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
