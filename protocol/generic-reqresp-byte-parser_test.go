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
	"testing"

	"github.com/stretchr/testify/assert"
)

// 单元测试
func TestGenericByteParser(t *testing.T) {
	// 定义包头和包尾
	edger := PacketEdger{
		Head: [2]byte{0xAF, 0x00},
		Tail: [2]byte{0xFA, 0x00},
	}

	// 创建数据校验器
	checker := &SimpleChecker{}

	// 创建 GenericByteParser 实例
	parser := NewGenericByteParser(checker, edger)

	// 测试数据：正确的有效数据包
	validPacket := []byte{
		0xAF, 0x00, // Header
		0x03,             // Length (3 bytes data)
		0x01, 0x02, 0x03, // Data
		0x00,       // Checksum (0x01 ^ 0x02 ^ 0x03)
		0xFA, 0x00, // Tail
	}

	t.Run("Valid Packet", func(t *testing.T) {
		data, err := parser.ParseBytes(validPacket)
		assert.NoError(t, err)
		assert.Equal(t, []byte{0x01, 0x02, 0x03}, data)
	})

	// 测试数据：包头错误
	invalidHeader := []byte{
		0x00, 0x00, // 错误的包头
		0x03, // Length
		0x01, 0x02, 0x03,
		0x00,       // Checksum
		0xFA, 0x00, // Tail
	}

	t.Run("Invalid Header", func(t *testing.T) {
		data, err := parser.ParseBytes(invalidHeader)
		assert.Error(t, err)
		assert.Nil(t, data)
		assert.Equal(t, "no valid header found", err.Error())
	})

	// 测试数据：包尾错误
	invalidTail := []byte{
		0xAF, 0x00, // Header
		0x03, // Length
		0x01, 0x02, 0x03,
		0x00,       // Checksum
		0x00, 0x00, // 错误的包尾
	}

	t.Run("Invalid Tail", func(t *testing.T) {
		data, err := parser.ParseBytes(invalidTail)
		assert.Error(t, err)
		assert.Nil(t, data)
		assert.Equal(t, "no valid tail found", err.Error())
	})

	// 测试数据：数据长度不匹配
	lengthMismatch := []byte{
		0xAF, 0x00, // Header
		0x03, // Length (3 bytes)
		0x01, 0x02, 0x03,
		0x00,       // Checksum (expecting 0x00)
		0xFA, 0x00, // Tail
	}

	t.Run("Length Mismatch", func(t *testing.T) {
		data, err := parser.ParseBytes(lengthMismatch)
		assert.Error(t, err)
		assert.Nil(t, data)
		assert.Equal(t, "data length mismatch", err.Error())
	})

	// 测试数据：校验和错误
	invalidChecksum := []byte{
		0xAF, 0x00, // Header
		0x03, // Length
		0x01, 0x02, 0x03,
		0x01,       // 错误的校验和（正确应为0x00）
		0xFA, 0x00, // Tail
	}

	t.Run("Invalid Checksum", func(t *testing.T) {
		data, err := parser.ParseBytes(invalidChecksum)
		assert.Error(t, err)
		assert.Nil(t, data)
		assert.Equal(t, "data validation failed: data is empty", err.Error())
	})

	// 测试数据：空数据
	t.Run("Empty Data", func(t *testing.T) {
		data, err := parser.ParseBytes([]byte{})
		assert.Error(t, err)
		assert.Nil(t, data)
		assert.Equal(t, "no valid header found", err.Error())
	})

	// 测试数据：包头和包尾不匹配
	t.Run("Header Tail Mismatch", func(t *testing.T) {
		mismatchedPacket := []byte{
			0xAF, 0x00, // Header
			0x03, // Length
			0x01, 0x02, 0x03,
			0x00,       // Checksum
			0xFB, 0x00, // 错误的包尾
		}
		data, err := parser.ParseBytes(mismatchedPacket)
		assert.Error(t, err)
		assert.Nil(t, data)
		assert.Equal(t, "no valid tail found", err.Error())
	})
}
