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

package protocol

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// 定义一个泛型结构体来存储解析后的数据
type ParsedData map[string]interface{}

// ParseBinary 函数
func ParseBinary(expr string, data []byte) (ParsedData, error) {
	parsedData := ParsedData{}
	cursor := 0

	// 分割表达式, 每个表达式形如 Key:Length:Type:Endian
	fields := strings.Split(expr, ";")
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}

		// 解析表达式中的 Key, Length, Type, Endian
		parts := strings.Split(field, ":")
		if len(parts) != 4 {
			return nil, fmt.Errorf("表达式格式错误: %s", field)
		}

		key := parts[0]
		lengthStr := parts[1] // 长度现在是位数
		dataType := parts[2]
		endian := parts[3]

		// 将 Length（位）转换为位长度，并确保字节对齐
		lengthBits, err := strconv.Atoi(lengthStr)
		if err != nil {
			return nil, fmt.Errorf("无效的长度: %s", lengthStr)
		}
		lengthBytes := (lengthBits + 7) / 8 // 按字节对齐

		// 检查是否有足够的数据可供解析
		if cursor+lengthBytes > len(data) {
			return nil, fmt.Errorf("数据长度不足以解析 %s", key)
		}

		// 根据 Endian 设置字节序
		var order binary.ByteOrder
		if endian == "BE" {
			order = binary.BigEndian
		} else if endian == "LE" {
			order = binary.LittleEndian
		} else {
			return nil, fmt.Errorf("不支持的字节序: %s", endian)
		}

		// 根据 Type 解析数据
		switch dataType {
		case "int":
			if lengthBits == 8 {
				parsedData[key] = int(data[cursor])
			} else if lengthBits == 16 {
				parsedData[key] = int(order.Uint16(data[cursor : cursor+lengthBytes]))
			} else if lengthBits == 32 {
				parsedData[key] = int(order.Uint32(data[cursor : cursor+lengthBytes]))
			} else {
				return nil, fmt.Errorf("不支持的 int 长度: %d bits", lengthBits)
			}
		case "string":
			parsedData[key] = string(data[cursor : cursor+lengthBytes])
		case "float":
			if lengthBits == 32 {
				bits := order.Uint32(data[cursor : cursor+lengthBytes])
				parsedData[key] = float32(bits)
			} else if lengthBits == 64 {
				bits := order.Uint64(data[cursor : cursor+lengthBytes])
				parsedData[key] = float64(bits)
			} else {
				return nil, fmt.Errorf("不支持的 float 长度: %d bits", lengthBits)
			}
		default:
			return nil, fmt.Errorf("不支持的数据类型: %s", dataType)
		}

		// 移动光标
		cursor += lengthBytes
	}

	return parsedData, nil
}

func (parsedData ParsedData) String() string {
	jsonData, _ := json.Marshal(parsedData)
	return string(jsonData)
}
