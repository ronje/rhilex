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

package dlt6452007

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// ReadDataRequest represents a DLT645 read data request.
type ReadDataRequest struct {
	Address [6]byte // 地址域
	DataID  uint32  // 数据标识
}

// ReadDataResponse represents a DLT645 read data response.
type ReadDataResponse struct {
	Address [6]byte // 地址域
	DataID  uint32  // 数据标识
	Status  byte    // 状态域
	Value   []byte  // 数据值
}

// PackReadDataRequest constructs a byte slice from a ReadDataRequest.
func PackReadDataRequest(req ReadDataRequest) []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteByte(0x68) // 帧起始符
	buffer.Write(req.Address[:])
	buffer.WriteByte(0x11) // 读取数据的控制码
	buffer.WriteByte(0x04) // 数据长度
	var dataIDBytes [4]byte
	binary.BigEndian.PutUint32(dataIDBytes[:], req.DataID)
	buffer.Write(dataIDBytes[:])
	checksum := calculateChecksum(buffer.Bytes()[1:])
	buffer.WriteByte(checksum)
	buffer.WriteByte(0x16) // 帧结束符
	return buffer.Bytes()
}

// UnpackReadDataResponse parses a byte slice into a ReadDataResponse.
func UnpackReadDataResponse(packet []byte) (*ReadDataResponse, error) {
	if len(packet) < 12 { // 最小帧长度
		return nil, fmt.Errorf("invalid packet length")
	}

	if packet[0] != 0x68 || packet[len(packet)-1] != 0x16 { // 检查帧起始符和结束符
		return nil, fmt.Errorf("invalid frame header or end")
	}

	var resp ReadDataResponse
	copy(resp.Address[:], packet[1:7])
	resp.Status = packet[9]
	resp.DataID = binary.BigEndian.Uint32(packet[10:14])
	resp.Value = packet[14 : len(packet)-2] // 排除校验和和帧结束符

	// 校验和验证
	calculatedChecksum := calculateChecksum(packet[1 : len(packet)-2])
	if packet[len(packet)-2] != calculatedChecksum {
		return nil, fmt.Errorf("invalid checksum")
	}

	return &resp, nil
}

// calculateChecksum computes the checksum for the given byte slice.
func calculateChecksum(data []byte) byte {
	var checksum byte
	for _, b := range data {
		checksum += b
	}
	return checksum
}
