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

package test

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
)

// ---

// ### **正确的数据包实例**
// 这些数据包符合协议的设计，可以成功解析并通过校验。

// #### **正确的包1**（数据包长度正常）
// ```text
// 数据包： 0xAA 0x55 0x03 0x01 0x02 0x03 0x06
// 解释：
//   - Header: 0xAA55 (2字节)
//   - Length: 0x03 (3字节)
//   - Data: 0x01 0x02 0x03 (3字节)
//   - Checksum: 0x06 (校验和为 0x01 + 0x02 + 0x03 = 0x06)
// ```

// **解析过程：**
// - `Header` 确认成功（0xAA55）。
// - `Length` = 3，读取后续的 3 字节数据。
// - 数据：`0x01 0x02 0x03`，校验和计算结果为 `0x06`。
// - 校验和正确，数据包有效，成功解析。

// ---

// #### **正确的包2**（最大长度数据包）
// ```text
// 数据包： 0xAA 0x55 0xFF 0x01 0x02 0x03 ... 0xFF 0xFE
// 解释：
//   - Header: 0xAA55 (2字节)
//   - Length: 0xFF (255字节)
//   - Data: 0x01 0x02 0x03 ... 0xFE (255字节)
//   - Checksum: 0xFE (假设为正确的校验和)
// ```

// **解析过程：**
// - `Header` 确认成功（0xAA55）。
// - `Length` = 255，读取后续的 255 字节数据。
// - 校验和正确，数据包有效，成功解析。

// ---

// ### **极端错误的数据包实例**
// 这些数据包违反协议规则，可能是恶意攻击或传输错误引起的。我们需要确保解析过程能够正确识别并丢弃这些无效数据。

// #### **错误包1**（伪造的 Length 字段导致越界）
// ```text
// 数据包： 0xAA 0x55 0x0A 0x01 0x02 0x03 0xAA 0x55 0x02 0x04 0x05 0x09
// 解释：
//   - Header: 0xAA55 (2字节)
//   - Length: 0x0A (10字节)
//   - Data: 0x01 0x02 0x03 (只有3字节数据)
//   - Checksum: 0x09 (假设校验和)
// ```

// **解析过程：**
// - `Header` 确认成功（0xAA55）。
// - `Length` = 10，但实际数据只有 3 字节，这会导致越界。
// - 系统会检测到数据不足 10 字节，报错 `incomplete packet data`，丢弃该数据包。

// **防御措施**：
// - 长度字段 `Length` 应该被验证是否符合缓冲区的实际可用空间。
// - 如果数据不足，则丢弃数据包并重新同步。

// ---

// #### **错误包2**（伪造的 `Header` 字段）
// ```text
// 数据包： 0xFF 0xFF 0x05 0x01 0x02 0x03 0x06
// 解释：
//   - Header: 0xFFFF (无效头部，协议应为0xAA55)
//   - Length: 0x05 (5字节)
//   - Data: 0x01 0x02 0x03 0x04 0x05 (5字节)
//   - Checksum: 0x06 (假设校验和)
// ```

// **解析过程：**
// - `Header` 为 0xFFFF，不是预期的 0xAA55，因此识别为无效包。
// - 系统会抛出 `no valid header found` 错误，丢弃该数据包并重新同步。

// **防御措施**：
// - 数据包头部需要严格验证，确保匹配预定义的 `Header` 值。
// - 无效的 `Header` 会导致数据包被丢弃。

// ---

// #### **错误包3**（Length 超过最大长度）
// ```text
// 数据包： 0xAA 0x55 0xFF 0x01 0x02 0x03 ... 0xFF 0xFE
// 解释：
//   - Header: 0xAA55 (2字节)
//   - Length: 0xFF (255字节)
//   - Data: 0x01 0x02 0x03 ... 0xFE (255字节)
//   - Checksum: 0xFE (假设为正确的校验和)
// ```

// **解析过程：**
// - `Length` 为 0xFF，表示数据包应该有 255 字节。
// - 如果协议规定最大数据包长度为 128 字节，则该数据包超出最大长度范围，解析时会报错。

// **防御措施**：
// - 在解析时，应该对 `Length` 字段进行最大长度检查。如果超过协议规定的最大长度（例如 128 字节），丢弃该数据包。

// ---

// #### **错误包4**（数据包中没有 `Checksum` 校验字段）
// ```text
// 数据包： 0xAA 0x55 0x03 0x01 0x02 0x03
// 解释：
//   - Header: 0xAA55 (2字节)
//   - Length: 0x03 (3字节)
//   - Data: 0x01 0x02 0x03 (3字节)
//   - Checksum 缺失（没有校验和字段）
// ```

// **解析过程：**
// - `Header` 确认成功（0xAA55）。
// - `Length` = 3，读取 3 字节数据，但缺少 `Checksum` 字段，系统会报错 `invalid checksum`。
// - 数据包无效，丢弃。

// **防御措施**：
// - 必须在数据包中包含 `Checksum` 校验字段，缺失校验字段的包会被丢弃。
// - 校验和验证是数据完整性的核心步骤。

// ---

// #### **错误包5**（伪造的 `Data` 内容）
// ```text
// 数据包： 0xAA 0x55 0x05 0x01 0x02 0x03 0xFF 0xFF 0xFF 0xFF
// 解释：
//   - Header: 0xAA55 (2字节)
//   - Length: 0x05 (5字节)
//   - Data: 0x01 0x02 0x03 0xFF 0xFF (数据中有异常值，可能是恶意篡改)
//   - Checksum: 0xFF (伪造的校验和)
// ```

// **解析过程：**
// - `Header` 确认成功（0xAA55）。
// - `Length` = 5，读取 5 字节数据，内容为 `0x01 0x02 0x03 0xFF 0xFF`，数据明显有误。
// - 校验和为伪造的 `0xFF`，计算出的校验和与实际值不符，系统会报错 `invalid checksum`，丢弃该数据包。

// **防御措施**：
// - 校验和字段不能被伪造，数据内容如果不符合校验规则，则丢弃数据包。

// ---

// ### **总结**

// 以下是数据包的正确与错误示例：

// | **数据包类型**     | **描述**                              | **解析结果** |
// |-------------------|-------------------------------------|-------------|
// | **正确包1**       | 合法的固定长度数据包                    | 解析成功       |
// | **正确包2**       | 最大长度数据包                         | 解析成功       |
// | **错误包1**       | 伪造的 `Length` 字段导致数据不足         | `incomplete packet data` |
// | **错误包2**       | 无效的 `Header` 字段                    | `no valid header found`  |
// | **错误包3**       | `Length` 超过最大长度                   | `length exceeds maximum allowed size` |
// | **错误包4**       | 缺少 `Checksum` 校验字段               | `invalid checksum` |
// | **错误包5**       | 伪造的 `Data` 和校验和                 | `invalid checksum` |

// 数据包规则
const (
	Header       = 0xAA55 // 数据包头部标识
	MaxPacketLen = 256    // 数据包最大允许长度
)

// ParsePacket 解析数据包
func ParsePacket(buffer []byte) ([]byte, []byte, error) {
	if len(buffer) < 4 { // Header (2 bytes) + Length (1 byte) + Checksum (1 byte)
		return nil, buffer, errors.New("incomplete packet")
	}

	// 查找Header
	startIndex := bytes.Index(buffer, []byte{byte(Header >> 8), byte(Header & 0xFF)})
	if startIndex == -1 {
		return nil, nil, errors.New("no valid header found")
	}

	// 确保缓冲区中有足够的字节来读取Length字段
	if len(buffer[startIndex:]) < 3 {
		return nil, buffer[startIndex:], errors.New("incomplete packet header")
	}

	// 检查Length字段
	length := int(buffer[startIndex+2])
	if length > MaxPacketLen {
		return nil, buffer[startIndex+3:], errors.New("length exceeds maximum allowed size")
	}

	// 检查缓冲区中是否包含完整数据包
	totalPacketLen := 3 + length + 1 // Header(2) + Length(1) + Data(length) + Checksum(1)
	if len(buffer[startIndex:]) < totalPacketLen {
		return nil, buffer[startIndex:], errors.New("incomplete packet data")
	}

	// 校验Checksum
	data := buffer[startIndex+3 : startIndex+3+length]
	checksum := buffer[startIndex+3+length]
	if calculateChecksum(data) != checksum {
		return nil, buffer[startIndex+3+length:], errors.New("invalid checksum")
	}

	// 返回完整包和剩余缓冲区
	return buffer[startIndex : startIndex+totalPacketLen], buffer[startIndex+totalPacketLen:], nil
}

// calculateChecksum 计算校验和
func calculateChecksum(data []byte) byte {
	var sum byte
	for _, b := range data {
		sum += b
	}
	return sum
}

// 实现单元测试函数 Test_ParsePacket
// go test -timeout 30s -run ^TestParsePacket$ github.com/hootrhino/rhilex/test -v -count=1
func TestParsePacket(t *testing.T) {
	// 测试用例集
	tests := []struct {
		name          string
		buffer        []byte
		expectedError string
	}{
		{
			name:          "Valid packet (Normal case)",
			buffer:        []byte{0xAA, 0x55, 0x03, 0x01, 0x02, 0x03, 0x06}, // 正常数据包
			expectedError: "",
		},
		{
			name:          "Valid packet (Max Length)",
			buffer:        append([]byte{0xAA, 0x55, 0xFF}, bytes.Repeat([]byte{0x01}, 255)...), // 长度为 255 字节的数据包
			expectedError: "",
		},
		{
			name:          "Malformed packet (Invalid Header)",
			buffer:        []byte{0xFF, 0xFF, 0x03, 0x01, 0x02, 0x03, 0x06}, // 错误的 Header 字段
			expectedError: "no valid header found",
		},
		{
			name:          "Malformed packet (Length exceeds available buffer)",
			buffer:        []byte{0xAA, 0x55, 0x0A, 0x01, 0x02, 0x03}, // Length 为 10，数据只有 3 字节
			expectedError: "incomplete packet data",
		},
		{
			name:          "Malformed packet (Checksum mismatch)",
			buffer:        []byte{0xAA, 0x55, 0x03, 0x01, 0x02, 0x03, 0xFF}, // 错误的 Checksum
			expectedError: "invalid checksum",
		},
		{
			name:          "Malformed packet (Length exceeds maximum allowed size)",
			buffer:        append([]byte{0xAA, 0x55, 0xFF}, bytes.Repeat([]byte{0x01}, 255)...), // 超过最大长度（假设最大为 128 字节）
			expectedError: "length exceeds maximum allowed size",
		},
		{
			name:          "Malformed packet (Missing Checksum)",
			buffer:        []byte{0xAA, 0x55, 0x03, 0x01, 0x02, 0x03}, // 缺少 Checksum 字段
			expectedError: "invalid checksum",
		},
	}

	// 遍历所有测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packet, remaining, err := ParsePacket(tt.buffer)

			// 如果预期没有错误，检查是否有解析错误
			if tt.expectedError == "" {
				if err != nil {
					t.Errorf("Expected no error, but got %v", err)
				}
				if len(packet) == 0 {
					t.Error("Expected valid packet, but got nil or empty packet")
				}
				if len(remaining) == 0 {
					t.Error("Expected non-empty remaining buffer")
				}
			} else {
				// 如果预期有错误，检查错误信息
				if err == nil {
					t.Errorf("Expected error %v, but got no error", tt.expectedError)
				} else if err.Error() != tt.expectedError {
					t.Errorf("Expected error %v, but got %v", tt.expectedError, err)
				}
			}
		})
	}
}

// 测试 Checksum 计算函数
func TestCalculateChecksum(t *testing.T) {
	tests := []struct {
		data             []byte
		expectedChecksum byte
	}{
		{
			data:             []byte{0x01, 0x02, 0x03},
			expectedChecksum: 0x06, // 0x01 + 0x02 + 0x03 = 0x06
		},
		{
			data:             []byte{0xFF, 0xFF, 0xFF},
			expectedChecksum: 0xFD, // 0xFF + 0xFF + 0xFF = 0xFD
		},
		{
			data:             []byte{0x00, 0x00, 0x00},
			expectedChecksum: 0x00, // 0x00 + 0x00 + 0x00 = 0x00
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Checksum for %v", tt.data), func(t *testing.T) {
			result := calculateChecksum(tt.data)
			if result != tt.expectedChecksum {
				t.Errorf("Expected checksum %v, but got %v", tt.expectedChecksum, result)
			}
		})
	}
}
