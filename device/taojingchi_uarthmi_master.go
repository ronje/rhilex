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

package device

import (
	"encoding/json"
	"errors"
	"strings"

	"fmt"

	"sync"
	"time"

	serial "github.com/hootrhino/goserial"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/resconfig"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

type TaoJingChiHmiMainConfig struct {
	UartConfig resconfig.UartConfig `json:"uartConfig"`
}

type TaoJingChiHmiDevice struct {
	typex.XStatus
	serialPort serial.Port
	status     typex.SourceState
	RuleEngine typex.Rhilex
	mainConfig TaoJingChiHmiMainConfig
	locker     sync.Locker
}

/*
*
* 通用串口透传，纯粹的串口读取网关
*
 */
func NewTaoJingChiHmiDevice(e typex.Rhilex) typex.XDevice {
	uart := new(TaoJingChiHmiDevice)
	uart.locker = &sync.Mutex{}
	uart.mainConfig = TaoJingChiHmiMainConfig{
		UartConfig: resconfig.UartConfig{
			Uart:     "/dev/ttyS1",
			BaudRate: 9600,
			Timeout:  30,
			DataBits: 8,
			Parity:   "N",
			StopBits: 1,
		},
	}
	uart.RuleEngine = e
	return uart
}

//  初始化
func (uart *TaoJingChiHmiDevice) Init(devId string, configMap map[string]any) error {
	uart.PointId = devId

	if err := utils.BindSourceConfig(configMap, &uart.mainConfig); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	if err := uart.mainConfig.UartConfig.Validate(); err != nil {
		return nil
	}
	return nil
}

// 启动
func (uart *TaoJingChiHmiDevice) Start(cctx typex.CCTX) error {
	uart.Ctx = cctx.Ctx
	uart.CancelCTX = cctx.CancelCTX

	config := serial.Config{
		Address:  uart.mainConfig.UartConfig.Uart,
		BaudRate: uart.mainConfig.UartConfig.BaudRate,
		DataBits: uart.mainConfig.UartConfig.DataBits,
		Parity:   uart.mainConfig.UartConfig.Parity,
		StopBits: uart.mainConfig.UartConfig.StopBits,
		Timeout:  time.Duration(uart.mainConfig.UartConfig.Timeout) * time.Millisecond,
	}
	serialPort, errOpen := serial.Open(&config)
	if errOpen != nil {
		glogger.GLogger.Error("serial port start failed err:", errOpen, ", config:", config)
		return errOpen
	}
	uart.serialPort = serialPort

	go func(uart *TaoJingChiHmiDevice) {
		buffer := [1024]byte{}
		for {
			select {
			case <-uart.Ctx.Done():
				return
			default:
			}
			data, err := readUntilEndMarker(uart.serialPort, buffer[:])
			if err != nil {
				glogger.Error(err)
				uart.status = typex.SOURCE_DOWN
				return
			}
			uart.RuleEngine.WorkDevice(uart.Details(), string(data))
		}
	}(uart)
	uart.status = typex.SOURCE_UP
	return nil
}

// 读取数据直到遇到 \xFF\xFF\xFF 结尾
func readUntilEndMarker(serialPort serial.Port, buffer []byte) ([]byte, error) {
	if len(buffer) < 3 {
		return nil, errors.New("buffer is too small to hold end marker")
	}
	var data []byte
	// 读取数据直到遇到结束标记 \xFF\xFF\xFF
	for {
		readBytes, err := serialPort.Read(buffer)
		if err != nil {
			if strings.Contains(err.Error(), "timeout") {
				continue
			}
			return nil, err
		}
		if readBytes >= 3 &&
			buffer[readBytes-3] == 0xFF &&
			buffer[readBytes-2] == 0xFF &&
			buffer[readBytes-1] == 0xFF {
			return append(data, buffer[:readBytes-3]...), nil
		}
		data = append(data, buffer[:readBytes]...)
		if readBytes > 2 && buffer[readBytes-2] == 0xFF && buffer[readBytes-1] == 0xFF {
			copy(buffer, buffer[readBytes-2:])
		}
	}
}

// 从设备里面读数据出来:
// t1.txt="OK"\xff\xff\xff
func (uart *TaoJingChiHmiDevice) OnCtrl(cmd []byte, args []byte) ([]byte, error) {

	if string(cmd) == "WriteToHmi" {
		cmds := []string{}
		if errUnmarshal := json.Unmarshal(args, &cmds); errUnmarshal != nil {
			return nil, errUnmarshal
		}
		for _, ctrl := range cmds {
			// 陶晶池指令默认需要加上包尾 \xFF\xFF\xFF
			_, err := uart.serialPort.Write([]byte(ctrl + "\xFF\xFF\xFF"))
			if err != nil {
				return nil, err
			}
		}
		return []byte{}, nil
	}
	return []byte{}, fmt.Errorf("unsupported cmd")
}

// 设备当前状态
func (uart *TaoJingChiHmiDevice) Status() typex.SourceState {
	if uart.serialPort == nil {
		uart.status = typex.SOURCE_DOWN
	}
	return uart.status
}

// 停止设备
func (uart *TaoJingChiHmiDevice) Stop() {
	uart.status = typex.SOURCE_DOWN
	if uart.CancelCTX != nil {
		uart.CancelCTX()
	}
	if uart.serialPort != nil {
		uart.serialPort.Close()
	}

}

func (uart *TaoJingChiHmiDevice) Details() *typex.Device {
	return uart.RuleEngine.GetDevice(uart.PointId)
}

func (uart *TaoJingChiHmiDevice) SetState(status typex.SourceState) {
	uart.status = status
}

func (uart *TaoJingChiHmiDevice) OnDCACall(UUID string, Command string, Args any) typex.DCAResult {
	return typex.DCAResult{}
}
