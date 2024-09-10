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

	"github.com/hootrhino/rhilex/component/lostcache"
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
	// 离线缓存
	CacheOfflineData *bool `json:"cacheOfflineData" title:"离线缓存"`
}

type GenericUart struct {
	typex.XStatus
	uartConfig uartctrl.UartConfig
	status     typex.SourceState
	locker     sync.Mutex
	serialPort serial.Port
	mainConfig GenericUartMainConfig
}

func NewGenericUart(e typex.Rhilex) typex.XTarget {
	mdev := new(GenericUart)
	mdev.RuleEngine = e
	mdev.mainConfig = GenericUartMainConfig{
		PortUuid:         "/dev/ttyS1",
		DataMode:         "RAW_STRING",
		PingPacket:       "RHILEX",
		AllowPing:        new(bool),
		CacheOfflineData: new(bool),
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
	uartPort, err := uartctrl.GetUart(mdev.mainConfig.PortUuid)
	if err != nil {
		return err
	}
	if uartPort.Busy {
		return fmt.Errorf("mdev is busying now, Occupied By:%s", uartPort.OccupyBy)
	}
	switch tCfg := uartPort.Config.(type) {
	case uartctrl.UartConfig:
		{
			mdev.uartConfig = tCfg
		}
	default:
		{
			return fmt.Errorf("Invalid config:%s", uartPort.Config)
		}
	}
	config := serial.Config{
		Address:  mdev.uartConfig.Uart,
		BaudRate: mdev.uartConfig.BaudRate,
		DataBits: mdev.uartConfig.DataBits,
		Parity:   mdev.uartConfig.Parity,
		StopBits: mdev.uartConfig.StopBits,
		Timeout:  time.Duration(mdev.uartConfig.Timeout) * time.Millisecond,
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
	// 补发数据
	if *mdev.mainConfig.CacheOfflineData {
		if CacheData, err1 := lostcache.GetLostCacheData(mdev.PointId); err1 != nil {
			glogger.GLogger.Error(err1)
		} else {
			for _, data := range CacheData {
				_, errTo := mdev.To(data.Data)
				if errTo == nil {
					lostcache.DeleteLostCacheData(data.ID)
				}
			}
		}
	}
	glogger.GLogger.Info("GenericUart started:", mdev.uartConfig.Uart)
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
		switch S := data.(type) {
		case string:
			_, err := mdev.serialPort.Write([]byte(S))
			if *mdev.mainConfig.CacheOfflineData {
				lostcache.SaveLostCacheData(lostcache.CacheDataDto{
					TargetId: mdev.PointId,
					Data:     S,
				})
			}
			return nil, err
		}
		return 0, fmt.Errorf("serial Port invalid")
	}
	if mdev.mainConfig.DataMode == "RAW_STRING" {
		switch S := data.(type) {
		case string:
			_, err := mdev.serialPort.Write([]byte(S))
			if *mdev.mainConfig.CacheOfflineData {
				lostcache.SaveLostCacheData(lostcache.CacheDataDto{
					TargetId: mdev.PointId,
					Data:     S,
				})
			}
			return nil, err
		}
	}
	if mdev.mainConfig.DataMode == "HEX_STRING" {
		switch S := data.(type) {
		case string:
			Hex, err := hex.DecodeString(S)
			if err != nil {
				if *mdev.mainConfig.CacheOfflineData {
					lostcache.SaveLostCacheData(lostcache.CacheDataDto{
						TargetId: mdev.PointId,
						Data:     S,
					})
				}
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
