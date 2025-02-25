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

package utils

import "reflect"

// IsArrayAndGetTypeList 检查给定的 any 是否为数组，并返回其元素类型列表。
func IsArrayAndGetValueList(i any) ([]any, bool) {
	if i == nil {
		return []any{}, false
	}
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Array && v.Kind() != reflect.Slice {
		return []any{i}, false
	}

	values := make([]any, v.Len())
	for i := 0; i < v.Len(); i++ {
		values[i] = v.Index(i).Interface()
	}
	return values, true
}
