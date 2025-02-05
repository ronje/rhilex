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

package transceiver

import (
	"context"
	"fmt"
	"sync"
	"time"

	serial "github.com/hootrhino/goserial"
	"github.com/hootrhino/rhilex/component/eventbus"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

// DefaultTransceiver 默认收发器结构体
type DefaultTransceiver struct {
	rhilex     typex.Rhilex
	locker     sync.Mutex
	serialPort serial.Port
	cancel     context.CancelFunc
	ctx        context.Context
	mainConfig TransceiverConfig
	Event      string
	EventType  string
}

// NewTransceiver 创建一个新的收发器实例
func NewTransceiver(rhilex typex.Rhilex) Transceiver {
	return &DefaultTransceiver{
		rhilex: rhilex,
		locker: sync.Mutex{},
		mainConfig: TransceiverConfig{
			Name:              "default_transceiver",
			Address:           "COM1",
			BaudRate:          9600,
			DataBits:          8,
			Parity:            "N",
			StopBits:          1,
			IOTimeout:         0, // IOTimeout * time.Millisecond
			TransportProtocol: 1,
		},
	}
}

// Start 启动收发器
func (tc *DefaultTransceiver) Start(Config TransceiverConfig) error {
	tc.mainConfig = Config
	tc.EventType = "transceiver.up.data"
	tc.Event = fmt.Sprintf("transceiver.up.data.%s", tc.mainConfig.Name)
	glogger.GLogger.Info("Transceiver Init:", tc.mainConfig.Name)
	ctx, cancel := context.WithCancel(context.Background())
	tc.ctx = ctx
	tc.cancel = cancel
	config := serial.Config{
		Address:  tc.mainConfig.Address,
		BaudRate: tc.mainConfig.BaudRate,
		DataBits: tc.mainConfig.DataBits,
		Parity:   tc.mainConfig.Parity,
		StopBits: tc.mainConfig.StopBits,
		Timeout:  time.Duration(tc.mainConfig.IOTimeout) * time.Millisecond,
	}

	serialPort, errOpen := serial.Open(&config)
	if errOpen != nil {
		glogger.GLogger.Error("serial port start failed", "err", errOpen, "config", config)
		return errOpen
	}
	tc.serialPort = serialPort
	// 开启协程来监听串口的数据
	go func() {
		for {
			select {
			case <-tc.ctx.Done():
				return
			default:
			}
			ctx1, cancel1 := context.WithTimeout(context.Background(), config.Timeout)
			n, buffer := utils.ReadInLeastTimeout(ctx1, serialPort, config.Timeout)
			cancel1()
			if n > 0 {
				glogger.GLogger.Debug("Transceiver.Receive:", string(buffer[:n]))
				eventbus.Publish(tc.Event, eventbus.EventMessage{
					Topic:   tc.Event,
					From:    "transceiver",
					Type:    tc.EventType,
					Event:   tc.Event,
					Ts:      uint64(time.Now().UnixMilli()),
					Payload: buffer[:n],
				})
			}
		}
	}()
	glogger.GLogger.Info("Transceiver Started:", tc.mainConfig.Name)
	return nil
}

// Ctrl 控制收发器
func (tc *DefaultTransceiver) Ctrl(topic, args []byte, timeout time.Duration) ([]byte, error) {
	return []byte("CMD SEND SUCCESS"), nil
}

// Info 获取收发器信息
func (tc *DefaultTransceiver) Info() CommunicatorInfo {
	return CommunicatorInfo{
		Name:     tc.mainConfig.Name,
		Model:    "RHILEX TRANSCEIVER",
		Type:     LORA,
		Vendor:   "RHILEX-TECH",
		Mac:      "00:00:00:00:00:00",
		Firmware: "v0.0.1",
	}
}

// Status 获取收发器状态
func (tc *DefaultTransceiver) Status() TransceiverStatus {
	if tc.serialPort == nil {
		return TransceiverStatus{
			Code:  TC_DOWN,
			Error: fmt.Errorf("serial port is not opened"),
		}
	}
	return TransceiverStatus{
		Code:  TC_UP,
		Error: nil,
	}
}

// Stop 停止收发器
func (tc *DefaultTransceiver) Stop() {
	tc.cancel()
	if tc.serialPort != nil {
		tc.serialPort.Close()
	}
	glogger.GLogger.Info("Transceiver Stopped:", tc.mainConfig.Name)
}
