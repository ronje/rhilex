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

package szy2062016

const (
	// CRC多项式X^7 + X^6 + X^5 + X^2 + 1的二进制表示为0b1110001
	polynomial = 0xE1 // 0b11100001
)

// 计算CRC校验码
func crc(data []byte) byte {
	// 初始值
	crc := byte(0)

	for _, b := range data {
		crc ^= b                 // XOR当前字节
		for i := 0; i < 8; i++ { // 处理8位
			if (crc & 0x80) != 0 { // 如果最高位为1
				crc = (crc << 1) ^ polynomial // 左移并XOR多项式
			} else {
				crc <<= 1 // 仅左移
			}
		}
	}
	return crc
}
