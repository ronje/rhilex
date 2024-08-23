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

package dataschema

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/hootrhino/rhilex/utils"
)

type IoTPropertyType string

// "INTEGER", "BOOL", "FLOAT", "STRING", "GEO"
const (
	// 目前边缘侧暂时只支持常见类型
	IoTPropertyTypeString  IoTPropertyType = "STRING"
	IoTPropertyTypeInteger IoTPropertyType = "INTEGER"
	IoTPropertyTypeFloat   IoTPropertyType = "FLOAT"
	IoTPropertyTypeBool    IoTPropertyType = "BOOL"
	IoTPropertyTypeGeo     IoTPropertyType = "GEO"
)

// string
type IoTPropertyString string

// int
type IoTPropertyInteger int64

// float
type IoTPropertyFloat float32

// bool
type IoTPropertyBool bool

// 地理坐标系统
type IoTPropertyGeo string

/*
* 物模型,边缘端目前暂时只支持属性
*
 */
type IoTSchema struct {
	IoTProperties map[string]IoTProperty `json:"iotProperties"`
}

// 规则
type IoTPropertyRule struct {
	DefaultValue any       `json:"defaultValue"` // 默认值: 0 false ''
	Max          int       `json:"max"`          // int|float|string: 最大值
	Min          int       `json:"min"`          // int|float|string: 最小值
	TrueLabel    string    `json:"trueLabel"`    // bool: 真值label
	FalseLabel   string    `json:"falseLabel"`   // bool: 假值label
	Round        int       `json:"round"`        // float: 小数点位
	validator    Validator `json:"-"`
}

// 物模型属性
type IoTProperty struct {
	UUID        string          `json:"uuid"`            // Cache uuid
	Label       string          `json:"label"`           // UI显示的那个文本
	Name        string          `json:"name"`            // 变量关联名
	Description string          `json:"description"`     // 额外信息
	Type        IoTPropertyType `json:"type"`            // 类型, 只能是上面几种
	Rw          string          `json:"rw"`              // R读 W写 RW读写
	Unit        string          `json:"unit"`            // 单位 例如：摄氏度、米、牛等等
	Value       interface{}     `json:"value,omitempty"` // Value 是运行时值, 前端不用填写
	Rule        IoTPropertyRule `json:"-"`
}

// 验证语法
// "INTEGER", "BOOL", "FLOAT", "STRING", "GEO"
func (I *IoTProperty) CovertAndValidate() error {
	switch I.Type {
	case "INTEGER":
		I.StringValue()
	case "BOOL":
		I.BoolValue()
	case "FLOAT":
		I.FloatValue()
	case "STRING":
		I.StringValue()
	case "GEO":
		I.GeoValue()
	default:
		return fmt.Errorf("Unsupported type:%v", I.Type)
	}
	return I.Rule.validator.Validate(I.Value)
}
func (I *IoTProperty) StringValue() string {
	if I == nil {
		return ""
	}
	switch I.Type {
	case IoTPropertyTypeString:
		{
			I.Rule.validator = StringRule{
				MinLength:    I.Rule.Min,
				MaxLength:    I.Rule.Max,
				DefaultValue: "",
			}
			return I.Value.(string)
		}
	}
	return ""
}
func (I *IoTProperty) IntValue() int {
	if I == nil {
		return 0
	}
	switch I.Type {
	case IoTPropertyTypeInteger:
		{
			I.Rule.validator = IntegerRule{
				Min:          I.Rule.Min,
				Max:          I.Rule.Max,
				DefaultValue: 0,
			}
			return I.Value.(int)
		}
	}
	return 0
}
func (I *IoTProperty) FloatValue() float64 {
	if I == nil {
		return 0
	}
	switch I.Type {
	case IoTPropertyTypeFloat:
		{
			I.Rule.validator = FloatRule{
				Min:          I.Rule.Min,
				Max:          I.Rule.Max,
				DefaultValue: 0.00,
			}
			return I.Value.(float64)
		}
	}
	return 0
}
func (I *IoTProperty) BoolValue() bool {
	if I == nil {
		return false
	}
	switch I.Type {
	case IoTPropertyTypeBool:
		{
			I.Rule.validator = BoolRule{
				TrueLabel:    I.Rule.TrueLabel,
				FalseLabel:   I.Rule.FalseLabel,
				DefaultValue: false,
			}
			return I.Value.(bool)
		}
	}
	return false
}
func (I *IoTProperty) GeoValue() IoTPropertyGeo {
	if I == nil {
		return ""
	}
	switch I.Type {
	case IoTPropertyTypeGeo:
		{
			I.Rule.validator = GeoRule{
				DefaultValue: "0,0",
			}
			return I.Value.(IoTPropertyGeo)
		}
	}
	return "0,0"
}

/*
*
* 验证物模型本身是否合法, 包含了 IoTPropertyType，Rule 的对应关系
*
 */
func (I *IoTProperty) HoldValidator() error {
	switch I.Type {
	case "INTEGER":
		Rule := IntegerRule{
			Min: I.Rule.Min,
			Max: I.Rule.Max,
		}
		switch T := I.Rule.DefaultValue.(type) {
		case int:
			Rule.DefaultValue = T
		case int32:
			Rule.DefaultValue = int(T)
		case int64:
			Rule.DefaultValue = int(T)
		default:
			Rule.DefaultValue = 0
		}
		I.Rule.validator = Rule
		return nil
	case "BOOL":
		Rule := BoolRule{
			TrueLabel:  I.Rule.TrueLabel,
			FalseLabel: I.Rule.FalseLabel,
		}
		switch T := I.Rule.DefaultValue.(type) {
		case bool:
			Rule.DefaultValue = T
		default:
			Rule.DefaultValue = false
		}
		I.Rule.validator = Rule
		return nil

	case "FLOAT":
		Rule := FloatRule{
			Min:   I.Rule.Min,
			Max:   I.Rule.Max,
			Round: I.Rule.Round,
		}
		switch T := I.Rule.DefaultValue.(type) {
		case float32:
			Rule.DefaultValue = float64(T)
		case float64:
			Rule.DefaultValue = float64(T)
		default:
			Rule.DefaultValue = 0
		}
		I.Rule.validator = Rule
		return nil

	case "STRING":
		Rule := StringRule{
			MinLength: I.Rule.Min,
			MaxLength: I.Rule.Max,
		}
		switch T := I.Rule.DefaultValue.(type) {
		case string:
			Rule.DefaultValue = T
		default:
			Rule.DefaultValue = ""
		}
		I.Rule.validator = Rule
		return nil

	case "GEO":
		Rule := GeoRule{DefaultValue: "0,0"}
		I.Rule.validator = Rule
		return nil
	}
	return fmt.Errorf("Unsupported Validator type:%v", I.Type)
}

/*
*
* 验证规则
*
 */
func (I IoTProperty) ValidateRule() error {
	return I.Rule.validator.Validate(I.Value)
}

/*
*
* 物模型规则 : String|Float|Int|Bool
*
 */
type Validator interface {
	Validate(Value interface{}) error
}

/*
*
* 字符串规则
*
 */
type StringRule struct {
	MinLength    int    `json:"minLength"`
	MaxLength    int    `json:"maxLength"`
	DefaultValue string `json:"defaultValue"`
}

func (S StringRule) Validate(Value interface{}) error {
	switch SV := Value.(type) {
	case string:
		L := len(SV)
		if L >= S.MaxLength {
			return fmt.Errorf("Value exceed Max Length:%v", L)
		} else {
			return nil
		}
	}
	return fmt.Errorf("Invalid String type: %v, Expect UTF8 string", Value)
}

/*
*
* 整数规则
*
 */
type IntegerRule struct {
	DefaultValue int `json:"defaultValue"`
	Max          int `json:"max"`
	Min          int `json:"min"`
}

func (V IntegerRule) Validate(Value interface{}) error {
	switch T := Value.(type) {
	case int32:
		if T < int32(V.Max) && T > int32(V.Min) {
			return nil
		}
	case int64:
		if T < int64(V.Max) && T > int64(V.Min) {
			return nil
		}
	}
	return fmt.Errorf("Invalid Int type:%v", Value)
}

/*
*
* 浮点数规则
*
 */
type FloatRule struct {
	DefaultValue float64 `json:"defaultValue"`
	Max          int     `json:"max"`
	Min          int     `json:"min"`
	Round        int     `json:"round"`
}

func (V FloatRule) Validate(Value interface{}) error {
	switch T := Value.(type) {
	case float32:
		if T < float32(V.Max) && T > float32(V.Min) {
			return nil
		}
	case float64:
		if T < float64(V.Max) && T > float64(V.Min) {
			return nil
		}
	}
	return fmt.Errorf("Invalid Float type:%v", Value)
}

/*
*
* 布尔规则
*
 */
type BoolRule struct {
	DefaultValue bool   `json:"defaultValue"`
	TrueLabel    string `json:"trueLabel"`  // UI界面对应的label
	FalseLabel   string `json:"falseLabel"` // UI界面对应的label
}

func (V BoolRule) Validate(Value interface{}) error {
	switch Value.(type) {
	case bool:
		return nil
	}
	return fmt.Errorf("Invalid Bool type:%v", Value)
}

/*
*
* 地理坐标规则
*
 */
type GeoRule struct {
	DefaultValue IoTPropertyGeo `json:"defaultValue"`
}

func (V GeoRule) Validate(Value interface{}) error {
	switch T := Value.(type) {
	case string:
		if isValidGEO(T) {
			return nil
		}
	}
	return fmt.Errorf("Invalid Coordinate type:%v", Value)
}

// isValidGEO 验证字符串是否是有效的地理坐标
func isValidGEO(coord string) bool {
	regexPattern := `^(\-?\d+(\.\d+)?),\s*(\-?\d+(\.\d+)?)$`
	matched, _ := regexp.MatchString(regexPattern, coord)
	if !matched {
		return false
	}
	parts := strings.Split(coord, ",")
	if len(parts) != 2 {
		return false
	}
	latitude, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil || latitude < -90 || latitude > 90 {
		return false
	}
	longitude, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil || longitude < -180 || longitude > 180 {
		return false
	}
	return true
}

/*
*
  - Check Type:"INTEGER", "BOOL", "FLOAT", "STRING", "GEO"
*/
func CheckPropertyType(s string) error {
	Types := []string{"INTEGER", "BOOL", "FLOAT", "STRING", "GEO"}
	if !slices.Contains(Types, s) {
		return fmt.Errorf("Invalid Property Type, Must one of:%v", Types)
	}
	return nil
}

/*
*
* 验证R\W类型
*
 */
func ValidateRw(s string) error {
	if !utils.SContains([]string{"R", "W", "RW"}, s) {
		return fmt.Errorf("RW Value Only Support 'R' or 'W' or 'RW'")
	}
	return nil
}
