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

package cjt1882004

import (
	"errors"
	"fmt"
	"time"

	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
	"github.com/sirupsen/logrus"
)

type CJT188ClientHandler struct {
	logger        *logrus.Logger
	DataLinkLayer *CJT1882004DataLinkLayer
	Transporter   *CJT188SerialTransporter
}

// NewRTUClientHandler allocates and initializes a RTUClientHandler.
func NewCJT188ClientHandler(Transporter typex.GenericRWC) *CJT188ClientHandler {
	handler := &CJT188ClientHandler{
		DataLinkLayer: &CJT1882004DataLinkLayer{},
		Transporter: NewCJT188SerialTransporter(TransporterConfig{
			ReadTimeout:  100 * time.Millisecond,
			WriteTimeout: 100 * time.Millisecond,
			Port:         Transporter,
		}),
	}
	return handler
}

/**
 * 关闭
 *
 */
func (handler *CJT188ClientHandler) Close() error {
	if handler.Transporter.port == nil {
		return fmt.Errorf("invalid Transporter")
	}
	return handler.Transporter.port.Close()
}

func (handler *CJT188ClientHandler) SetLogger(logger *logrus.Logger) {
	handler.logger = logger
}

/**
 * 发送请求
 *
 */
func (handler *CJT188ClientHandler) Request(data []byte) ([]byte, error) {
	data = append([]byte{0xFE, 0xFE, 0xFE, 0xFE}, data...)
	request := ""
	for _, b := range data {
		request += fmt.Sprintf("0x%02x ", b)
	}
	handler.logger.Debug("CJT188ClientHandler.Request:", request)
	r, e := handler.Transporter.SendFrame(data)
	handler.logger.Debug("Transporter.Received:", utils.ByteDumpHexString(r))
	return r, e
}

// 解包CJT188协议帧
func (handler *CJT188ClientHandler) DecodeCJT188Frame0x01(data []byte) (CJT188Frame0x01, error) {
	handler.logger.Debug("CJT188ClientHandler.DecodeCJT188Frame0x11:", data)
	frame := CJT188Frame0x01{
		Start:        data[0],
		MeterType:    data[1],
		Address:      [7]byte{data[2], data[3], data[4], data[5], data[6], data[7], data[8]},
		CtrlCode:     data[9],
		DataLength:   data[10],
		DataType:     [2]byte{data[11], data[12]},
		SerialNumber: data[13],
	}
	if len(data) < 15 {
		return frame, errors.New("invalid frame length")
	}

	if int(data[10]) > len(data) {
		return frame, errors.New("invalid frame length")
	}

	if frame.DataLength > 0 {
		frame.DataArea = data[14 : 14+data[10]-3]
	}
	frame.CheckSum = data[len(data)-2]
	frame.End = data[len(data)-1]
	if frame.Start != 0x68 || frame.End != 0x16 {
		return frame, fmt.Errorf("invalid start or end byte")
	}
	CheckCrcErr := handler.DataLinkLayer.CheckCrc(data[0:len(data)-2], frame.CheckSum)
	if CheckCrcErr != nil {
		return frame, CheckCrcErr
	}
	return frame, nil
}

/**
 * 打包帧
 *
 */
func (handler *CJT188ClientHandler) DecodeCJT188Frame0x01Response(data []byte) (CJT188Frame0x01, error) {
	return handler.DecodeCJT188Frame0x01(data)
}
