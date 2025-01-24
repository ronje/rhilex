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
	"errors"
	"fmt"
)

// 定义通用解析器
type GenericByteParser struct {
	edger   PacketEdger
	checker DataChecker
}

// 创建一个新的解析器
func NewGenericByteParser(checker DataChecker, edger PacketEdger) *GenericByteParser {
	return &GenericByteParser{
		edger:   edger,
		checker: checker,
	}
}
func (parser *GenericByteParser) PackBytes(frame *ApplicationFrame) ([]byte, error) {
	b, err := frame.Encode()
	if err != nil {
		return nil, err
	}
	bodyLength := len(b)
	packetLength := 2 + bodyLength + 2
	packet := make([]byte, packetLength)
	copy(packet[:2], parser.edger.Head[:])
	copy(packet[2:], b)
	copy(packet[packetLength-2:], parser.edger.Tail[:])
	return packet, nil
}

// 解析字节流，提取有效数据包
func (parser *GenericByteParser) ParseBytes(b []byte) ([]byte, error) {
	// 查找包头
	startIdx := -1
	for i := 0; i < len(b)-1; i++ {
		if b[i] == parser.edger.Head[0] && b[i+1] == parser.edger.Head[1] {
			startIdx = i
			break
		}
	}

	// 如果没有找到包头
	if startIdx == -1 {
		return nil, errors.New("no valid header found")
	}

	// 查找包尾
	endIdx := -1
	for i := startIdx + 2; i < len(b)-1; i++ {
		if b[i] == parser.edger.Tail[0] && b[i+1] == parser.edger.Tail[1] {
			endIdx = i + 2 // 包尾的位置（包含包尾的字节）
			break
		}
	}

	// 如果没有找到包尾
	if endIdx == -1 {
		return nil, errors.New("no valid tail found")
	}

	// 提取包体数据（包含包头和包尾之间的部分）
	packetData := b[startIdx+2 : endIdx-2]

	// 使用用户自定义的校验器验证数据
	if err := parser.checker.CheckData(packetData); err != nil {
		return nil, fmt.Errorf("data validation failed: %v", err)
	}

	// 返回有效数据
	return packetData, nil
}
