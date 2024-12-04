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
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/hootrhino/rhilex/resconfig"
	"github.com/hootrhino/rhilex/component/intercache"
	"github.com/hootrhino/rhilex/component/interdb"
	userproto "github.com/hootrhino/rhilex/device/useprotocol"

	serial "github.com/hootrhino/goserial"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

type UserProtocolCommonConfig struct {
	Mode         string `json:"mode" validate:"required"` // UART | TCP
	AutoRequest  *bool  `json:"autoRequest" validate:"required"`
	BatchRequest *bool  `json:"batchRequest" validate:"required"`
}
type GenericUserProtocolDataPoint struct {
	UUID      string `json:"uuid"`
	Command   string `json:"command"`
	Tag       string `json:"tag"`
	Alias     string `json:"alias"`
	Frequency int64  `json:"frequency"`
}

/*
*
* 自定义协议
*
 */
type GenericUserProtocolConfig struct {
	CommonConfig UserProtocolCommonConfig `json:"commonConfig" validate:"required"`
	HostConfig   resconfig.HostConfig        `json:"hostConfig"`
	UartConfig   resconfig.UartConfig        `json:"uartConfig"`
}
type GenericUserProtocolDevice struct {
	typex.XStatus
	status      typex.DeviceState
	RuleEngine  typex.Rhilex
	uartHandler *userproto.UserProtocolClientHandler
	tcpHandler  *userproto.UserProtocolClientHandler
	mainConfig  GenericUserProtocolConfig
	DataPoints  map[string]GenericUserProtocolDataPoint
}

func NewGenericUserProtocolDevice(e typex.Rhilex) typex.XDevice {
	gw := new(GenericUserProtocolDevice)
	gw.RuleEngine = e
	gw.mainConfig = GenericUserProtocolConfig{
		CommonConfig: UserProtocolCommonConfig{
			Mode: "UART",
			AutoRequest: func() *bool {
				b := false
				return &b
			}(),
			BatchRequest: func() *bool {
				b := false
				return &b
			}(),
		},
		HostConfig: resconfig.HostConfig{
			Host:    "127.0.0.1",
			Port:    502,
			Timeout: 3000,
		},
		UartConfig: resconfig.UartConfig{
			Timeout:  3000,
			Uart:     "COM1",
			BaudRate: 9600,
			DataBits: 8,
			Parity:   "N",
			StopBits: 1,
		},
	}
	gw.DataPoints = map[string]GenericUserProtocolDataPoint{}

	return gw
}

// 初始化
func (gw *GenericUserProtocolDevice) Init(devId string, configMap map[string]interface{}) error {
	gw.PointId = devId
	intercache.RegisterSlot(gw.PointId)

	if err := utils.BindSourceConfig(configMap, &gw.mainConfig); err != nil {
		return err
	}
	if !utils.SContains([]string{"UART", "TCP"}, gw.mainConfig.CommonConfig.Mode) {
		return errors.New("unsupported mode, only can be one of 'TCP' or 'UART'")
	}
	if err := gw.mainConfig.UartConfig.Validate(); err != nil {
		return nil
	}
	var DataPoints []GenericUserProtocolDataPoint
	PointLoadErr := interdb.DB().Table("m_user_protocol_data_points").
		Where("device_uuid=?", devId).Find(&DataPoints).Error
	if PointLoadErr != nil {
		return PointLoadErr
	}
	LastFetchTime := uint64(time.Now().UnixMilli())
	for _, DataPoint := range DataPoints {
		if DataPoint.Frequency < 1 {
			return errors.New("'frequency' must grate than 50 millisecond")
		}
		gw.DataPoints[DataPoint.UUID] = GenericUserProtocolDataPoint{
			UUID:      DataPoint.UUID,
			Command:   DataPoint.Command,
			Tag:       DataPoint.Tag,
			Alias:     DataPoint.Alias,
			Frequency: DataPoint.Frequency,
		}
		intercache.SetValue(gw.PointId, DataPoint.UUID, intercache.CacheValue{
			UUID:          DataPoint.UUID,
			Status:        0,
			LastFetchTime: LastFetchTime,
			Value:         "0",
			ErrMsg:        "--",
		})
	}
	return nil
}

// 启动
func (gw *GenericUserProtocolDevice) Start(cctx typex.CCTX) error {
	gw.Ctx = cctx.Ctx
	gw.CancelCTX = cctx.CancelCTX
	if gw.mainConfig.CommonConfig.Mode == "UART" {
		config := serial.Config{
			Address:  gw.mainConfig.UartConfig.Uart,
			BaudRate: gw.mainConfig.UartConfig.BaudRate,
			DataBits: gw.mainConfig.UartConfig.DataBits,
			Parity:   gw.mainConfig.UartConfig.Parity,
			StopBits: gw.mainConfig.UartConfig.StopBits,
			Timeout:  time.Duration(gw.mainConfig.UartConfig.Timeout) * time.Millisecond,
		}
		serialPort, errOpen := serial.Open(&config)
		if errOpen != nil {
			glogger.GLogger.Error("serial port start failed err:", errOpen, ", config:", config)
			return errOpen
		}
		gw.uartHandler = userproto.NewUserProtocolClientHandler(serialPort)
		gw.uartHandler.SetLogger(glogger.Logrus)
		if *gw.mainConfig.CommonConfig.AutoRequest {
			go gw.work(gw.uartHandler)
		}
		goto END

	}
	if gw.mainConfig.CommonConfig.Mode == "TCP" {
		tcpconn, errDial := net.Dial("tcp",
			fmt.Sprintf("%s:%d", gw.mainConfig.HostConfig.Host,
				gw.mainConfig.HostConfig.Port))
		if errDial != nil {
			return errDial
		}
		gw.tcpHandler = userproto.NewUserProtocolClientHandler(tcpconn)
		gw.uartHandler.SetLogger(glogger.Logrus)
		if *gw.mainConfig.CommonConfig.AutoRequest {
			go gw.work(gw.tcpHandler)
		}
		goto END
	}
END:
	gw.status = typex.DEV_UP
	return nil
}

type UserProtocolReadData struct {
	Tag     string `json:"tag"`
	Command string `json:"command"`
	Value   string `json:"value"`
}

func (gw *GenericUserProtocolDevice) work(handler *userproto.UserProtocolClientHandler) {
	for {
		select {
		case <-gw.Ctx.Done():
			return
		default:
		}
		UserProtocolReadDataList := []UserProtocolReadData{}
		for _, DataPoint := range gw.DataPoints {
			lastTimes := uint64(time.Now().UnixMilli())
			NewValue := intercache.CacheValue{
				UUID:          DataPoint.UUID,
				LastFetchTime: lastTimes,
				Value:         "0",
			}
			CommandBytes, err1 := utils.HexStringToBytes(DataPoint.Command)
			if err1 != nil {
				glogger.GLogger.Error(err1)
				NewValue.Status = 0
				NewValue.ErrMsg = err1.Error()
				intercache.SetValue(gw.PointId, DataPoint.UUID, NewValue)
				continue
			}
			RespBytes, errRequest := handler.Request(CommandBytes)
			if errRequest != nil {
				glogger.GLogger.Error(errRequest)
				NewValue.Status = 0
				NewValue.ErrMsg = errRequest.Error()
				intercache.SetValue(gw.PointId, DataPoint.UUID, NewValue)
				continue
			}
			Value := hex.EncodeToString(RespBytes)
			NewValue.Status = 1
			NewValue.ErrMsg = ""
			intercache.SetValue(gw.PointId, DataPoint.UUID, NewValue)
			if !*gw.mainConfig.CommonConfig.BatchRequest {
				if bytes, err := json.Marshal(UserProtocolReadData{
					Tag:     DataPoint.Tag,
					Command: DataPoint.Command,
					Value:   Value,
				}); err != nil {
					glogger.GLogger.Error(err)
				} else {
					glogger.GLogger.Debug(string(bytes))
					gw.RuleEngine.WorkDevice(gw.Details(), string(bytes))
				}
			} else {
				UserProtocolReadDataList = append(UserProtocolReadDataList, UserProtocolReadData{
					Tag:     DataPoint.Tag,
					Command: DataPoint.Command,
					Value:   Value,
				})
			}
			time.Sleep(time.Duration(DataPoint.Frequency) * time.Millisecond)
		}
		if *gw.mainConfig.CommonConfig.BatchRequest {
			if len(UserProtocolReadDataList) > 0 {
				if bytes, err := json.Marshal(UserProtocolReadDataList); err != nil {
					glogger.GLogger.Error(err)
				} else {
					glogger.GLogger.Debug(string(bytes))
					gw.RuleEngine.WorkDevice(gw.Details(), string(bytes))
				}
			}
		}
	}
}

/*
*
* 外部指令交互, 常用来实现自定义协议等
*
 */
func (gw *GenericUserProtocolDevice) OnCtrl(cmd []byte, _ []byte) ([]byte, error) {
	return nil, errors.New("unknown write command:" + string(cmd))
}

// 设备当前状态
func (gw *GenericUserProtocolDevice) Status() typex.DeviceState {
	if gw.mainConfig.CommonConfig.Mode == "UART" {
		if gw.uartHandler == nil {
			return typex.DEV_DOWN
		}
	}
	if gw.mainConfig.CommonConfig.Mode == "TCP" {
		if gw.tcpHandler == nil {
			return typex.DEV_DOWN
		}
	}
	return gw.status
}

// 停止设备
func (gw *GenericUserProtocolDevice) Stop() {
	gw.status = typex.DEV_DOWN
	if gw.CancelCTX != nil {
		gw.CancelCTX()
	}
	if gw.mainConfig.CommonConfig.Mode == "UART" {
		if gw.uartHandler != nil {
			gw.uartHandler.Close()
		}
	}
	if gw.mainConfig.CommonConfig.Mode == "TCP" {
		if gw.tcpHandler != nil {
			gw.tcpHandler.Close()
		}
	}
	intercache.UnRegisterSlot(gw.PointId)
}

// 真实设备
func (gw *GenericUserProtocolDevice) Details() *typex.Device {
	return gw.RuleEngine.GetDevice(gw.PointId)
}

// 状态
func (gw *GenericUserProtocolDevice) SetState(status typex.DeviceState) {
	gw.status = status
}

/*
*
* 设备服务调用
*
 */
func (gw *GenericUserProtocolDevice) OnDCACall(_ string, Command string,
	Args interface{}) typex.DCAResult {

	return typex.DCAResult{}
}
