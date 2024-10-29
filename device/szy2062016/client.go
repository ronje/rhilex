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

import (
	"errors"
	"fmt"
	"time"

	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
	"github.com/sirupsen/logrus"
)

type SZY206ClientHandler struct {
	logger        *logrus.Logger
	DataLinkLayer *SZY2062016DataLinkLayer
	Transporter   *SZY206SerialTransporter
}

// NewRTUClientHandler allocates and initializes a RTUClientHandler.
func NewSZY206ClientHandler(Transporter typex.GenericRWC) *SZY206ClientHandler {
	handler := &SZY206ClientHandler{
		DataLinkLayer: &SZY2062016DataLinkLayer{},
		Transporter: NewSZY206SerialTransporter(TransporterConfig{
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
func (handler *SZY206ClientHandler) Close() error {
	if handler.Transporter.port == nil {
		return fmt.Errorf("invalid Transporter")
	}
	return handler.Transporter.port.Close()
}

func (handler *SZY206ClientHandler) SetLogger(logger *logrus.Logger) {
	handler.logger = logger
}

/**
 * 发送请求
 *
 */
func (handler *SZY206ClientHandler) Request(data []byte) ([]byte, error) {
	request := ""
	for _, b := range data {
		request += fmt.Sprintf("0x%02x ", b)
	}
	handler.logger.Debug("SZY206ClientHandler.Request:", request)
	r, e := handler.Transporter.SendFrame(data)
	handler.logger.Debug("Transporter.Received:", utils.ByteDumpHexString(r))
	return r, e
}

/**
 * 打包帧
 *
 */
func (handler *SZY206ClientHandler) DecodeSZY206Frame0x00(data []byte) (SZY206Frame0x00, error) {
	handler.logger.Debug("SZY206ClientHandler.DecodeSZY206Frame0x00:", data)
	frame := SZY206Frame0x00{
		Start1:     data[0],
		DataLength: data[1],
		Start2:     data[2],
		CtrlCode:   data[3],
		Address:    [5]byte{data[4], data[5], data[6], data[7], data[8]},
		// 数据
	}
	if len(data) < 11 {
		return frame, errors.New("invalid frame length")
	}

	if frame.DataLength > 0 {
		if frame.DataLength-2 > byte(len(data)) {
			return frame, errors.New("invalid frame length")
		}
		frame.DataArea = data[9 : 9+frame.DataLength-6]
	}
	frame.CheckSum = data[len(data)-2]
	frame.End = data[len(data)-1]
	if frame.Start1 != 0x68 || frame.End != 0x16 {
		return frame, fmt.Errorf("invalid start or end byte")
	}
	if int(frame.DataLength-6) != len(frame.DataArea) {
		return frame, fmt.Errorf("data length mismatch")
	}
	CheckCrcErr := handler.DataLinkLayer.CheckCrc(data[0:len(data)-2], frame.CheckSum)
	if CheckCrcErr != nil {
		return frame, CheckCrcErr
	}
	return frame, nil
}

// 解包SZY206协议帧
func (handler *SZY206ClientHandler) DecodeSZY206Frame0x00Response(data []byte) (SZY206Frame0x00, error) {
	return handler.DecodeSZY206Frame0x00(data)
}
