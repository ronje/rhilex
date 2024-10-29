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

package userprotocol

import (
	"fmt"
	"time"

	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
	"github.com/sirupsen/logrus"
)

type UserProtocolClientHandler struct {
	logger        *logrus.Logger
	DataLinkLayer *UserProtocolDataLinkLayer
	Transporter   *UserProtocolTransporter
}

// NewRTUClientHandler allocates and initializes a RTUClientHandler.
func NewUserProtocolClientHandler(Transporter typex.GenericRWC) *UserProtocolClientHandler {
	handler := &UserProtocolClientHandler{
		DataLinkLayer: &UserProtocolDataLinkLayer{},
		Transporter: NewUserProtocolTransporter(TransporterConfig{
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
func (handler *UserProtocolClientHandler) Close() error {
	if handler.Transporter.port == nil {
		return fmt.Errorf("invalid Transporter")
	}
	return handler.Transporter.port.Close()
}

func (handler *UserProtocolClientHandler) SetLogger(logger *logrus.Logger) {
	handler.logger = logger
}

/**
 * 发送请求
 *
 */
func (handler *UserProtocolClientHandler) Request(data []byte) ([]byte, error) {
	request := ""
	for _, b := range data {
		request += fmt.Sprintf("0x%02x ", b)
	}
	handler.logger.Debug("UserProtocolClientHandler.Request:", request)
	r, e := handler.Transporter.SendFrame(data)
	handler.logger.Debug("Transporter.Received:", utils.ByteDumpHexString(r))
	return r, e
}
