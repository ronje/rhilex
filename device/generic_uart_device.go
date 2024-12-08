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
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"strings"
	"sync"
	"time"

	serial "github.com/hootrhino/goserial"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/resconfig"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

type GenericUartCommonConfig struct {
	AutoRequest *bool `json:"autoRequest" validate:"required"`
}
type GenericUartRwConfig struct {
	Tag        string `json:"tag" validate:"required"`
	TimeSlice  uint64 `json:"timeSlice" validate:"required"`
	ReadFormat string `json:"readFormat" validate:"required" myself:"RAW,HEX,UTF8"` // 读取格式, "RAW"|"HEX"|"UTF8"
}
type GenericUartMainConfig struct {
	CommonConfig  GenericUartCommonConfig `json:"commonConfig" validate:"required"`
	RwConfig      GenericUartRwConfig     `json:"rwConfig" validate:"required"`
	UartConfig    resconfig.UartConfig    `json:"uartConfig"`
	CecollaConfig resconfig.CecollaConfig `json:"cecollaConfig"`
	AlarmConfig   resconfig.AlarmConfig   `json:"alarmConfig"`
}

type GenericUartDevice struct {
	typex.XStatus
	serialPort serial.Port
	status     typex.DeviceState
	RuleEngine typex.Rhilex
	mainConfig GenericUartMainConfig
	locker     sync.Locker
}

/*
*
* 通用串口透传，纯粹的串口读取网关
*
 */
func NewGenericUartDevice(e typex.Rhilex) typex.XDevice {
	uart := new(GenericUartDevice)
	uart.locker = &sync.Mutex{}
	uart.mainConfig = GenericUartMainConfig{
		CommonConfig: GenericUartCommonConfig{
			AutoRequest: func() *bool {
				b := true
				return &b
			}(),
		},
		RwConfig: GenericUartRwConfig{
			TimeSlice:  50,
			ReadFormat: "HEX",
			Tag:        "uart",
		},
		UartConfig: resconfig.UartConfig{
			Timeout:  3000,
			Uart:     "/dev/ttyS1",
			BaudRate: 9600,
			DataBits: 8,
			Parity:   "N",
			StopBits: 1,
		},
		CecollaConfig: resconfig.CecollaConfig{
			Enable: func() *bool {
				b := false
				return &b
			}(),
			EnableCreateSchema: func() *bool {
				b := true
				return &b
			}(),
		},
		AlarmConfig: resconfig.AlarmConfig{
			Enable: func() *bool {
				b := false
				return &b
			}(),
		},
	}
	uart.RuleEngine = e
	return uart
}

//  初始化
func (uart *GenericUartDevice) Init(devId string, configMap map[string]interface{}) error {
	uart.PointId = devId

	if err := utils.BindSourceConfig(configMap, &uart.mainConfig); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	if err := uart.mainConfig.UartConfig.Validate(); err != nil {
		return nil
	}
	if uart.mainConfig.RwConfig.TimeSlice < 30 {
		errA := fmt.Errorf("TimeSlice Must Great than 30, but current is: %v",
			uart.mainConfig.RwConfig.TimeSlice)
		glogger.GLogger.Error(errA)
		return errA
	}
	ReadFormatTypes := []string{"HEX", "RAW", "UTF8"}
	if !slices.Contains(ReadFormatTypes, uart.mainConfig.RwConfig.ReadFormat) {
		errA := fmt.Errorf("ReadFormat Only Support Type: %v", ReadFormatTypes)
		glogger.GLogger.Error(errA)
		return errA
	}
	return nil
}

// 启动
func (uart *GenericUartDevice) Start(cctx typex.CCTX) error {
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
	if !*uart.mainConfig.CommonConfig.AutoRequest {
		uart.status = typex.DEV_UP
		return nil
	}
	go func(ctx context.Context) {
		result := [2048]byte{}
		sliceTimer := time.NewTimer(time.Duration(uart.mainConfig.RwConfig.TimeSlice) * time.Millisecond)
		sliceTimer.Stop()
		peerCount := 0
		for {
			select {
			case <-ctx.Done():
				return
			case <-sliceTimer.C:
				mapV := map[string]interface{}{
					"tag": uart.mainConfig.RwConfig.Tag,
				}
				switch uart.mainConfig.RwConfig.ReadFormat {
				case "HEX":
					mapV["value"] = hex.EncodeToString(result[:peerCount])
				case "RAW":
					Value := []uint32{} // JSON会把[]Uint8识别为二进制，然后转换成Base64
					for i := 0; i < peerCount; i++ {
						Value = append(Value, uint32(result[i]))
					}
					mapV["value"] = Value
				case "UTF8":
					mapV["value"] = string(result[:peerCount])
				default:
					mapV["value"] = ""
					glogger.GLogger.Error("Not supported type:", uart.mainConfig.RwConfig.ReadFormat)
				}
				glogger.GLogger.Debug("Serial Port Read: ", result[:peerCount])
				bytes, _ := json.Marshal(mapV)
				uart.RuleEngine.WorkDevice(uart.Details(), string(bytes))
				for i := 0; i < peerCount; i++ {
					result[i] = 0 // 清空
				}
				peerCount = 0 // re-init index
			default:
				n, errR := io.ReadAtLeast(uart.serialPort, result[peerCount:], 1)
				if errR != nil {
					if !strings.Contains(errR.Error(), "timeout") {
						glogger.GLogger.Error(errR)
					}
				}
				if n != 0 {
					peerCount += n
					sliceTimer.Reset(time.Duration(uart.mainConfig.RwConfig.TimeSlice) * time.Millisecond)
				}
			}
		}
	}(uart.Ctx)
	uart.status = typex.DEV_UP
	return nil
}

// 从设备里面读数据出来:
func (uart *GenericUartDevice) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	result := [2048]byte{}
	if string(cmd) == "HEX" {
		hexs, err1 := hex.DecodeString(string(cmd))
		if err1 != nil {
			glogger.GLogger.Error(err1)
			return nil, err1
		}
		n, errSliceRequest := utils.SliceRequest(uart.Ctx, uart.serialPort,
			hexs, result[:], false, time.Duration(uart.mainConfig.RwConfig.TimeSlice)*time.Millisecond)
		if errSliceRequest != nil {
			return []byte{}, errSliceRequest
		}
		return result[:n], nil
	}
	if string(cmd) == "STRING" {
		// s := "t1.txt=\"RHILEX\"\xFF\xFF\xFF"
		n, err := uart.serialPort.Write(args)
		if err != nil {
			return nil, err
		}
		return result[:n], nil

	}
	return []byte{}, fmt.Errorf("unsupported cmd, must one of : STRING|HEX")
}

// 设备当前状态
func (uart *GenericUartDevice) Status() typex.DeviceState {
	if uart.serialPort == nil {
		uart.status = typex.DEV_DOWN
	}
	return uart.status
}

// 停止设备
func (uart *GenericUartDevice) Stop() {
	uart.status = typex.DEV_DOWN
	if uart.CancelCTX != nil {
		uart.CancelCTX()
	}
	if uart.serialPort != nil {
		uart.serialPort.Close()
	}

}

func (uart *GenericUartDevice) Details() *typex.Device {
	return uart.RuleEngine.GetDevice(uart.PointId)
}

func (uart *GenericUartDevice) SetState(status typex.DeviceState) {
	uart.status = status
}

func (uart *GenericUartDevice) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}
