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

	"github.com/sirupsen/logrus"
)

type DLT645ClientHandler struct {
	logger        *logrus.Logger
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
 * 关闭
 *
 */
func (handler *DLT645ClientHandler) Close() error {
	if handler.Transporter.port == nil {
		return fmt.Errorf("invalid Transporter")
	}
	return handler.Transporter.port.Close()
}

func (handler *DLT645ClientHandler) SetLogger(logger *logrus.Logger) {
	handler.logger = logger
}

/**
 * 发送请求
 *
 */
func (handler *DLT645ClientHandler) Request(data []byte) ([]byte, error) {
	request := ""
	for _, b := range data {
		request += fmt.Sprintf("0x%02x ", b)
	}
	handler.logger.Debug("DLT645ClientHandler.Request:", request)
	r, e := handler.Transporter.SendFrame(data)
	result := ""
	for _, b := range r {
		result += fmt.Sprintf("0x%02x ", b)
	}
	handler.logger.Debug("Transporter.Received:", result)
	return r, e
}

/**
 * 打包帧
 *
 */
func (handler *DLT645ClientHandler) DecodeDLT645Frame0x11(data []byte) (DLT645Frame0x11, error) {
	handler.logger.Debug("DLT645ClientHandler.DecodeDLT645Frame0x11:", data)
	frame := DLT645Frame0x11{
		Start:      data[0],
		Address:    [6]byte{data[1], data[2], data[3], data[4], data[5], data[6]},
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
	if frame.Start != 0x68 || frame.End != 0x16 {
		return frame, fmt.Errorf("invalid start or end byte")
	}
	if int(frame.DataLength-4) != len(frame.DataArea) {
		return frame, fmt.Errorf("data length mismatch")
	}
	CheckCrcErr := handler.DataLinkLayer.CheckCrc(data[0:len(data)-2], frame.CheckSum)
	if CheckCrcErr != nil {
		return frame, CheckCrcErr
	}
	return frame, nil
}

// 解包DLT645协议帧
func (handler *DLT645ClientHandler) DecodeDLT645Frame0x11Response(data []byte) (DLT645Frame0x11, error) {
	frame := DLT645Frame0x11{
		Start:      data[0],
		Address:    [6]byte{data[1], data[2], data[3], data[4], data[5], data[6]},
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
	if frame.Start != 0x68 || frame.End != 0x16 {
		return frame, fmt.Errorf("invalid start or end byte")
	}
	if int(frame.DataLength-4) != len(frame.DataArea) {
		return frame, fmt.Errorf("data length mismatch")
	}
	CheckCrcErr := handler.DataLinkLayer.CheckCrc(data[0:len(data)-2], frame.CheckSum)
	if CheckCrcErr != nil {
		return frame, CheckCrcErr
	}
	return frame, nil
}
