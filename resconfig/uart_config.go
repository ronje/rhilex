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

package resconfig

import (
	"errors"

	"github.com/hootrhino/rhilex/component/uartctrl"
)

type UartConfig struct {
	// 固定写法，表示串口最小一个包耗时，一般50毫秒足够
	Uart     string `json:"uart"`
	Timeout  int    `json:"timeout"`
	BaudRate int    `json:"baudRate"`
	DataBits int    `json:"dataBits"`
	Parity   string `json:"parity"`
	StopBits int    `json:"stopBits"`
}

func (uc *UartConfig) Validate() error {
	// 检查Uart是否为非空字符串
	if uc.Uart == "" {
		return errors.New("uart cannot be empty")
	}
	// 检查BaudRate是否在标准波特率范围内
	commonBaudRates := []int{110, 300, 600, 1200, 2400, 4800,
		9600, 14400, 19200, 38400, 57600, 115200, 128000, 256000}
	validBaudRate := false
	for _, br := range commonBaudRates {
		if uc.BaudRate == br {
			validBaudRate = true
			break
		}
	}
	if !validBaudRate {
		return errors.New("baudRate is not a standard value")
	}
	// 检查DataBits是否合理
	if uc.DataBits != 5 && uc.DataBits != 6 && uc.DataBits != 7 && uc.DataBits != 8 {
		return errors.New("dataBits must be 5, 6, 7, or 8")
	}
	// 检查Parity是否合理
	allowedParities := []string{"N", "O", "E", "M", "S"}
	validParity := false
	for _, p := range allowedParities {
		if uc.Parity == p {
			validParity = true
			break
		}
	}
	if !validParity {
		return errors.New("parity must be N, O, E, M, or S")
	}
	// 检查StopBits是否合理
	if uc.StopBits != 1 && uc.StopBits != 2 {
		return errors.New("stopBits must be 1 or 2")
	}
	// CheckSerialBusy
	if err := uartctrl.CheckSerialBusy(uc.Uart); err != nil {
		return err
	}
	return nil
}
