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
	"errors"
	"fmt"
)

// AppLayerFrame 应用层帧结构体
type AppLayerFrame struct {
	Length  uint16
	Payload []byte
	Crc16   uint16
}

// ToString 转换为字符串
func (f *AppLayerFrame) ToString() string {
	return fmt.Sprintf("Length: %d, Payload: %x, Crc16: %x", f.Length, f.Payload, f.Crc16)
}

// DecodeAppLayerFrame 解码数据字节数组为 AppLayerFrame 结构
func DecodeAppLayerFrame(data []byte) (AppLayerFrame, error) {
	if len(data) < 4 {
		return AppLayerFrame{}, errors.New("data too short to decode")
	}
	frame := AppLayerFrame{}
	frame.Length = binary.BigEndian.Uint16(data[:2])

	if len(data) < int(frame.Length)+4 {
		return AppLayerFrame{}, errors.New("data length mismatch")
	}
	frame.Payload = data[2 : 2+frame.Length]
	frame.Crc16 = binary.BigEndian.Uint16(data[2+frame.Length:])

	return frame, nil
}

// Encode 将 AppLayerFrame 编码为字节数组
func (f *AppLayerFrame) Encode() ([]byte, error) {
	if len(f.Payload) != int(f.Length) {
		return nil, errors.New("payload length does not match the specified Length")
	}
	data := make([]byte, 2+len(f.Payload)+2)
	binary.BigEndian.PutUint16(data[:2], f.Length)
	copy(data[2:2+len(f.Payload)], f.Payload)
	binary.BigEndian.PutUint16(data[len(data)-2:], f.Crc16)
	return data, nil
}
