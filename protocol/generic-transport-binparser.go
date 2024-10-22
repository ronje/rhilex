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

type ParsedData map[string]interface{}

// "ID:32:int:BE; Name:40:string:BE; Age:16:int:LE"

// ParseBinary function
func ParseBinary(expr string, data []byte) (ParsedData, error) {
	parsedData := ParsedData{}
	cursor := 0

	// Split the expression; each expression is of the form Key:Length:Type:Endian
	fields := strings.Split(expr, ";")
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}

		// Parse Key, Length, Type, Endian from the expression
		parts := strings.Split(field, ":")
		if len(parts) != 4 {
			return nil, fmt.Errorf("expression format error: %s", field)
		}

		key := parts[0]
		lengthStr := parts[1] // Length is now in bits
		dataType := parts[2]
		endian := parts[3]

		// Convert Length (bits) to byte length and ensure byte alignment
		lengthBits, err := strconv.Atoi(lengthStr)
		if err != nil {
			return nil, fmt.Errorf("invalid length: %s", lengthStr)
		}
		lengthBytes := (lengthBits + 7) / 8 // Align to bytes

		// Check if there is enough data to parse
		if cursor+lengthBytes > len(data) {
			return nil, fmt.Errorf("data length insufficient to parse %s", key)
		}

		// Set byte order based on Endian
		var order binary.ByteOrder
		if endian == "BE" {
			order = binary.BigEndian
		} else if endian == "LE" {
			order = binary.LittleEndian
		} else {
			return nil, fmt.Errorf("unsupported byte order: %s", endian)
		}

		// Parse data based on Type
		switch dataType {
		case "int":
			if lengthBits == 8 {
				parsedData[key] = int(data[cursor])
			} else if lengthBits == 16 {
				parsedData[key] = int(order.Uint16(data[cursor : cursor+lengthBytes]))
			} else if lengthBits == 32 {
				parsedData[key] = int(order.Uint32(data[cursor : cursor+lengthBytes]))
			} else {
				return nil, fmt.Errorf("unsupported int length: %d bits", lengthBits)
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
				return nil, fmt.Errorf("unsupported float length: %d bits", lengthBits)
			}
		default:
			return nil, fmt.Errorf("unsupported data type: %s", dataType)
		}

		// Move the cursor
		cursor += lengthBytes
	}

	return parsedData, nil
}

func (parsedData ParsedData) String() string {
	jsonData, _ := json.Marshal(parsedData)
	return string(jsonData)
}
