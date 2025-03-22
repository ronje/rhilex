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

package resconfig

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

/**
 * 验证器接口
 *
 */
type ConfigValidator interface {
	Validate() error
}

// ConvertAndValidateJSONToStruct 接收JSON数据和结构体指针，将JSON映射到结构体并进行验证
func ValidateConfig(jsonData []byte, structPtr any) error {
	// 反射获取结构体类型
	val := reflect.ValueOf(structPtr)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return errors.New("structPtr must be a non-nil pointer to a struct")
	}

	// 将JSON数据映射到结构体
	if err := json.Unmarshal(jsonData, structPtr); err != nil {
		return err
	}

	// 获取结构体的值
	val = val.Elem()

	// 检查结构体是否实现了Validator接口
	if val.Type().Implements(reflect.TypeOf((*ConfigValidator)(nil)).Elem()) {
		validator := val.Addr().Interface().(ConfigValidator)
		return validator.Validate()
	}
	return fmt.Errorf("not implement ConfigValidator")
}
