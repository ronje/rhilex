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
	"errors"
	"fmt"
	"io"
)

type DLT645ClientHandler struct {
	DataLinkLayer DLT6452007DataLinkLayer
	Transporter   dlt645SerialTransporter
}

// NewRTUClientHandler allocates and initializes a RTUClientHandler.
func NewDLT645ClientHandler(Transporter io.ReadWriteCloser) *DLT645ClientHandler {
	handler := &DLT645ClientHandler{
		DataLinkLayer: DLT6452007DataLinkLayer{},
		Transporter:   dlt645SerialTransporter{port: Transporter},
	}
	return handler
}

/**
 * 打包帧
 *
 */
func (handler *DLT645ClientHandler) DecodeDLT645Frame0x11(data []byte) (DLT645Frame0x11, error) {
	// 6 (Address) + 1 (CtrlCode) + 1 (DataLength) + 4 (DataType) + 1 (Checksum) + 0-N (minimum DataArea) + 1 (End)
	frame := DLT645Frame0x11{}
	if len(data) < 15 {
		return frame, fmt.Errorf("data too short to be a valid DLT645 frame")
	}
	frame.Start = data[0]
	copy(frame.Address[:], data[1:7])
	frame.CtrlCode = data[7]
	frame.DataLength = data[8]
	copy(frame.DataType[:], data[9:13])
	if 13+int(frame.DataLength)-4 > len(data) {
		return frame, fmt.Errorf("data too large to be a valid DLT645 frame")
	}
	frame.DataArea = data[13 : 13+int(frame.DataLength)-4]
	frame.CheckSum = data[len(data)-2]
	frame.End = data[len(data)-1]

	if frame.Start != 0x68 || frame.End != 0x16 {
		return frame, fmt.Errorf("invalid start or end byte")
	}

	if int(frame.DataLength-4) != len(frame.DataArea) {
		return frame, fmt.Errorf("data length mismatch")
	}
	CheckCrcErr := handler.DataLinkLayer.CheckCrc(data[0:13+int(frame.DataLength)-4], frame.CheckSum)
	if CheckCrcErr != nil {
		return frame, CheckCrcErr
	}

	return frame, nil
}

// 解包DLT645协议帧
func (handler *DLT645ClientHandler) DecodeDLT645Frame0x11Response(data []byte) (DLT645Frame0x11, error) {
	frame := DLT645Frame0x11{
		Start:      data[0],
		Address:    data[2:8],
		CtrlCode:   data[8],
		DataLength: data[9],
		DataType:   [4]byte{data[10], data[11], data[12], data[13]},
	}

	if len(data) < 16 { // 至少包含起始符、地址域、控制码、数据长度、校验和和结束符
		return frame, errors.New("invalid frame length")
	}

	if frame.DataLength > 0 {
		if 14+frame.DataLength-4 > byte(len(data)) {
			return frame, errors.New("invalid frame length")
		}
		frame.DataArea = data[14 : 14+frame.DataLength-4]
	}

	frame.CheckSum = data[len(data)-2]
	frame.End = data[len(data)-1]

	CheckCrcErr := handler.DataLinkLayer.CheckCrc(data[0:len(data)-2], frame.CheckSum)
	if CheckCrcErr != nil {
		return frame, CheckCrcErr
	}
	return frame, nil
}
