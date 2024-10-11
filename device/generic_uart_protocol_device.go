// Copyright (C) 2023 wwhai
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
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package device

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hootrhino/rhilex/common"

	serial "github.com/hootrhino/goserial"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

// 读出来的字节缓冲默认大小
const __DEFAULT_BUFFER_SIZE = 1024

type _CPDCommonConfig struct {
	Mode      string `json:"mode" validate:"required"`      // 传输协议
	RetryTime int    `json:"retryTime" validate:"required"` // 几次以后重启,0 表示不重启
}

/*
*
* 自定义协议
*
 */
type _GenericUartProtocolConfig struct {
	CommonConfig _CPDCommonConfig  `json:"commonConfig" validate:"required"`
	HostConfig   common.HostConfig `json:"hostConfig"`
	UartConfig   common.UartConfig `json:"uartConfig"`
}
type GenericUartProtocolDevice struct {
	typex.XStatus
	status     typex.DeviceState
	RuleEngine typex.Rhilex
	serialPort serial.Port // 串口
	mainConfig _GenericUartProtocolConfig
	errorCount int // 记录最大容错数，默认5次，出错超过5此就重启

}

func NewGenericUartProtocolDevice(e typex.Rhilex) typex.XDevice {
	mdev := new(GenericUartProtocolDevice)
	mdev.RuleEngine = e
	mdev.mainConfig = _GenericUartProtocolConfig{
		CommonConfig: _CPDCommonConfig{},
		HostConfig: common.HostConfig{
			Host:    "127.0.0.1",
			Port:    502,
			Timeout: 3000,
		},
		UartConfig: common.UartConfig{
			Timeout:  3000,
			Uart:     "/dev/ttyS1",
			BaudRate: 9600,
			DataBits: 8,
			Parity:   "N",
			StopBits: 1,
		},
	}
	return mdev

}

// 初始化
func (mdev *GenericUartProtocolDevice) Init(devId string, configMap map[string]interface{}) error {
	mdev.PointId = devId
	if err := utils.BindSourceConfig(configMap, &mdev.mainConfig); err != nil {
		return err
	}
	if !utils.SContains([]string{`UART`},
		mdev.mainConfig.CommonConfig.Mode) {
		return errors.New("option only 'UART'")
	}
	return nil
}

// 启动
func (mdev *GenericUartProtocolDevice) Start(cctx typex.CCTX) error {
	mdev.Ctx = cctx.Ctx
	mdev.CancelCTX = cctx.CancelCTX
	mdev.errorCount = 0
	mdev.status = typex.DEV_DOWN

	// 现阶段暂时只支持RS485串口, 以后有需求再支持TCP、UDP
	if mdev.mainConfig.CommonConfig.Mode == "UART" {

		config := serial.Config{
			Address:  mdev.mainConfig.UartConfig.Uart,
			BaudRate: mdev.mainConfig.UartConfig.BaudRate,
			DataBits: mdev.mainConfig.UartConfig.DataBits,
			Parity:   mdev.mainConfig.UartConfig.Parity,
			StopBits: mdev.mainConfig.UartConfig.StopBits,
			Timeout:  time.Duration(mdev.mainConfig.UartConfig.Timeout) * time.Millisecond,
		}
		serialPort, err := serial.Open(&config)
		if err != nil {
			glogger.GLogger.Error("serialPort start failed:", err)
			return err
		}
		mdev.serialPort = serialPort
		mdev.status = typex.DEV_UP
		return nil
	}
	return fmt.Errorf("unsupported Mode:%s", mdev.mainConfig.CommonConfig.Mode)
}

/*
*
* 数据读出来，对数据结构有要求, 其中Key必须是个数字或者数字字符串, 例如 1 or "1"
*
 */
func (mdev *GenericUartProtocolDevice) OnRead(cmd []byte, data []byte) (int, error) {
	return 0, errors.New("unknown read command:" + string(cmd))

}

/*
*
* 写进来的数据格式 参考@Protocol
*
 */

// 把数据写入设备
func (mdev *GenericUartProtocolDevice) OnWrite(cmd []byte, data []byte) (int, error) {
	return 0, errors.New("unknown write command:" + string(cmd))
}

/*
*
* 外部指令交互, 常用来实现自定义协议等
*
 */
func (mdev *GenericUartProtocolDevice) OnCtrl(cmd []byte, _ []byte) ([]byte, error) {
	glogger.GLogger.Debug("Time slice SliceRequest:", string(cmd))
	return mdev.ctrl(cmd)
}

// 设备当前状态
func (mdev *GenericUartProtocolDevice) Status() typex.DeviceState {
	if mdev.errorCount >= mdev.mainConfig.CommonConfig.RetryTime {
		mdev.CancelCTX()
		mdev.status = typex.DEV_DOWN
	}
	return mdev.status
}

// 停止设备
func (mdev *GenericUartProtocolDevice) Stop() {
	mdev.status = typex.DEV_DOWN
	if mdev.CancelCTX != nil {
		mdev.CancelCTX()
	}
	if mdev.mainConfig.CommonConfig.Mode == "UART" {
		if mdev.serialPort != nil {
			mdev.serialPort.Close()
		}
	}
}

// 真实设备
func (mdev *GenericUartProtocolDevice) Details() *typex.Device {
	return mdev.RuleEngine.GetDevice(mdev.PointId)
}

// 状态
func (mdev *GenericUartProtocolDevice) SetState(status typex.DeviceState) {
	mdev.status = status
}

/*
*
* 设备服务调用
*
 */
func (mdev *GenericUartProtocolDevice) OnDCACall(_ string, Command string,
	Args interface{}) typex.DCAResult {

	return typex.DCAResult{}
}

// --------------------------------------------------------------------------------------------------
// 内部函数
// --------------------------------------------------------------------------------------------------
func (mdev *GenericUartProtocolDevice) ctrl(args []byte) ([]byte, error) {
	hexs, err1 := hex.DecodeString(string(args))
	if err1 != nil {
		glogger.GLogger.Error(err1)
		return nil, err1
	}
	glogger.GLogger.Debug("Custom Protocol Device Request:", hexs)
	result := [__DEFAULT_BUFFER_SIZE]byte{}
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(mdev.mainConfig.HostConfig.Timeout)*time.Millisecond)
	count := 0
	var errSliceRequest error = nil
	if mdev.mainConfig.CommonConfig.Mode == "UART" {
		count, errSliceRequest = utils.SliceRequest(ctx, mdev.serialPort,
			hexs, result[:], false,
			time.Duration(30)*time.Millisecond /*30ms wait*/)
	}
	cancel()
	if errSliceRequest != nil {
		glogger.GLogger.Error("Custom Protocol Device Request error: ", errSliceRequest)
		mdev.errorCount++
		return nil, errSliceRequest
	}
	dataMap := map[string]string{}
	dataMap["in"] = string(args)
	out := hex.EncodeToString(result[:count])
	glogger.GLogger.Debug("Custom Protocol Device Response:", out)
	dataMap["out"] = out
	bytes, _ := json.Marshal(dataMap)
	return []byte(bytes), nil
}
