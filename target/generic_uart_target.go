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

package target

import (
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	serial "github.com/hootrhino/goserial"

	"github.com/hootrhino/rhilex/component/uartctrl"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

type GenericUartMainConfig struct {
	// 通用配置，包含AllowPing、DataMode、PingPacket、Timeout等
	// 是否允许Ping操作
	AllowPing *bool `json:"allowPing" validate:"required"`
	// 数据模式，RAW_STRING|HEX_STRING
	DataMode string `json:"dataMode" validate:"required"` // RAW_STRING|HEX_STRING
	// Ping请求的数据包内容
	PingPacket string `json:"pingPacket" validate:"required"`
	// 请求超时时间，单位为秒
	Timeout *int `json:"timeout" validate:"required"`
	// 端口UUID，用于识别特定的串口设备
	PortUuid string `json:"portUuid" validate:"required"`
}

type GenericUart struct {
	typex.XStatus
	hwPortConfig uartctrl.UartConfig
	status       typex.SourceState
	locker       sync.Mutex
	serialPort   serial.Port
	mainConfig   GenericUartMainConfig
}

func NewGenericUart(e typex.Rhilex) typex.XTarget {
	mdev := new(GenericUart)
	mdev.RuleEngine = e
	mdev.mainConfig = GenericUartMainConfig{
		PortUuid:   "/dev/ttyS1",
		DataMode:   "RAW_STRING",
		PingPacket: "RHILEX",
		AllowPing: func() *bool {
			b := false
			return &b
		}(),
		Timeout: func() *int {
			b := 3000
			return &b
		}(),
	}
	mdev.locker = sync.Mutex{}
	mdev.status = typex.SOURCE_DOWN
	return mdev
}

func (mdev *GenericUart) Init(outEndId string, configMap map[string]interface{}) error {
	mdev.PointId = outEndId
	if err := utils.BindSourceConfig(configMap, &mdev.mainConfig); err != nil {
		return err
	}
	return nil
}
func (mdev *GenericUart) Start(cctx typex.CCTX) error {
	mdev.Ctx = cctx.Ctx
	mdev.CancelCTX = cctx.CancelCTX
	hwPort, err := uartctrl.GetHwPort(mdev.mainConfig.PortUuid)
	if err != nil {
		return err
	}
	if hwPort.Busy {
		return fmt.Errorf("mdev is busying now, Occupied By:%s", hwPort.OccupyBy)
	}
	switch tCfg := hwPort.Config.(type) {
	case uartctrl.UartConfig:
		{
			mdev.hwPortConfig = tCfg
		}
	default:
		{
			return fmt.Errorf("Invalid config:%s", hwPort.Config)
		}
	}
	config := serial.Config{
		Address:  mdev.hwPortConfig.Uart,
		BaudRate: mdev.hwPortConfig.BaudRate,
		DataBits: mdev.hwPortConfig.DataBits,
		Parity:   mdev.hwPortConfig.Parity,
		StopBits: mdev.hwPortConfig.StopBits,
		Timeout:  time.Duration(mdev.hwPortConfig.Timeout) * time.Millisecond,
	}
	serialPort, err := serial.Open(&config)
	if err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	if *mdev.mainConfig.AllowPing {
		go func(serialPort serial.Port) {
			ticker := time.NewTicker(time.Duration(time.Second * 5))
			defer ticker.Stop()
			for {
				select {
				case <-mdev.Ctx.Done():
					return
				default:
				}
				if mdev.serialPort != nil {
					_, err := mdev.serialPort.Write([]byte(mdev.mainConfig.PingPacket))
					if err != nil {
						glogger.GLogger.Error(err)
					}
				}
				<-ticker.C
			}
		}(serialPort)
	}
	mdev.serialPort = serialPort
	mdev.status = typex.SOURCE_UP
	glogger.GLogger.Info("GenericUart started:", mdev.hwPortConfig.Uart)
	return nil
}

func (mdev *GenericUart) Status() typex.SourceState {
	if mdev.serialPort == nil {
		return typex.SOURCE_DOWN
	}
	_, err := mdev.serialPort.Write([]byte{})
	if err != nil {
		glogger.GLogger.Error(err)
		return typex.SOURCE_DOWN
	}
	return mdev.status
}

/*
*
* 数据写到串口
*
 */
func (mdev *GenericUart) To(data interface{}) (interface{}, error) {
	if mdev.serialPort == nil {
		mdev.status = typex.SOURCE_DOWN
		return 0, fmt.Errorf("serial Port invalid")
	}
	if mdev.mainConfig.DataMode == "RAW_STRING" {
		switch S := data.(type) {
		case string:
			return mdev.serialPort.Write([]byte(S))
		}
	}
	if mdev.mainConfig.DataMode == "HEX_STRING" {
		switch t := data.(type) {
		case string:
			Hex, err := hex.DecodeString(t)
			if err != nil {
				return nil, err
			}
			return mdev.serialPort.Write(Hex)
		}
	}
	return 0, fmt.Errorf("Invalid data:%v", data)
}

func (mdev *GenericUart) Stop() {
	mdev.status = typex.SOURCE_DOWN
	if mdev.CancelCTX != nil {
		mdev.CancelCTX()
	}
	if mdev.serialPort != nil {
		mdev.serialPort.Close()
	}
}
func (mdev *GenericUart) Details() *typex.OutEnd {
	return mdev.RuleEngine.GetOutEnd(mdev.PointId)
}
