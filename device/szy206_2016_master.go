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
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"

	serial "github.com/hootrhino/goserial"
	"github.com/hootrhino/rhilex/component/intercache"
	"github.com/hootrhino/rhilex/component/interdb"
	"github.com/hootrhino/rhilex/device/szy2062016"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/resconfig"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

// SZY206_2016_Data 用于封装M-Bus采集的数据
type SZY206_2016_DataPoint struct {
	UUID      string `json:"uuid"`
	MeterId   string `json:"meterId"`
	MeterType string `json:"meterType"`
	Tag       string `json:"tag"`
	Alias     string `json:"alias"`
	Frequency int64  `json:"frequency"`
}
type SZY206_2016_MasterGatewayCommonConfig struct {
	Mode         string `json:"mode" validate:"required"` // UART | TCP
	AutoRequest  *bool  `json:"autoRequest" validate:"required"`
	BatchRequest *bool  `json:"batchRequest" validate:"required"`
}

type SZY206_2016_MasterGatewayMainConfig struct {
	CommonConfig  SZY206_2016_MasterGatewayCommonConfig `json:"commonConfig" validate:"required"`
	HostConfig    resconfig.HostConfig                  `json:"hostConfig"`
	UartConfig    resconfig.UartConfig                  `json:"uartConfig"`
	CecollaConfig resconfig.CecollaConfig               `json:"cecollaConfig"`
	AlarmConfig   resconfig.AlarmConfig                 `json:"alarmConfig"`
}

/**
 *
 * SZY206_2016_
 */

type SZY206_2016_MasterGateway struct {
	typex.XStatus
	status      typex.DeviceState
	mainConfig  SZY206_2016_MasterGatewayMainConfig
	DataPoints  map[string]SZY206_2016_DataPoint
	uartHandler *szy2062016.SZY206ClientHandler
	tcpHandler  *szy2062016.SZY206ClientHandler
}

func NewSZY206_2016_MasterGateway(e typex.Rhilex) typex.XDevice {
	gw := new(SZY206_2016_MasterGateway)
	gw.RuleEngine = e
	gw.mainConfig = SZY206_2016_MasterGatewayMainConfig{
		CommonConfig: SZY206_2016_MasterGatewayCommonConfig{
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
			Host: "127.0.0.1",
			Port: 10065,
		},
		UartConfig: resconfig.UartConfig{
			Timeout:  3000,
			Uart:     "/dev/ttyS1",
			BaudRate: 2400,
			DataBits: 8,
			Parity:   "E",
			StopBits: 1,
		},
	}
	gw.DataPoints = map[string]SZY206_2016_DataPoint{}
	return gw
}

func (gw *SZY206_2016_MasterGateway) Init(devId string, configMap map[string]interface{}) error {
	gw.PointId = devId
	intercache.RegisterSlot(gw.PointId)

	if err := utils.BindSourceConfig(configMap, &gw.mainConfig); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	if !utils.SContains([]string{"UART", "TCP"}, gw.mainConfig.CommonConfig.Mode) {
		return errors.New("unsupported mode, only can be one of 'TCP' or 'UART'")
	}
	if err := gw.mainConfig.UartConfig.Validate(); err != nil {
		return nil
	}
	var DLT645_ModbusPointList []SZY206_2016_DataPoint
	PointLoadErr := interdb.DB().Table("m_szy2062016_data_points").
		Where("device_uuid=?", devId).Find(&DLT645_ModbusPointList).Error
	if PointLoadErr != nil {
		return PointLoadErr
	}
	LastFetchTime := uint64(time.Now().UnixMilli())
	for _, SZY206_2016_Point := range DLT645_ModbusPointList {
		if SZY206_2016_Point.Frequency < 1 {
			return errors.New("'frequency' must grate than 50 millisecond")
		}
		gw.DataPoints[SZY206_2016_Point.UUID] = SZY206_2016_DataPoint{
			UUID:      SZY206_2016_Point.UUID,
			MeterId:   SZY206_2016_Point.MeterId,
			MeterType: SZY206_2016_Point.MeterType,
			Tag:       SZY206_2016_Point.Tag,
			Alias:     SZY206_2016_Point.Alias,
			Frequency: SZY206_2016_Point.Frequency,
		}
		intercache.SetValue(gw.PointId, SZY206_2016_Point.UUID, intercache.CacheValue{
			UUID:          SZY206_2016_Point.UUID,
			Status:        0,
			LastFetchTime: LastFetchTime,
			Value:         "0",
			ErrMsg:        "--",
		})
	}
	return nil
}

func (gw *SZY206_2016_MasterGateway) Start(cctx typex.CCTX) error {
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
		gw.uartHandler = szy2062016.NewSZY206ClientHandler(serialPort)
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
		gw.tcpHandler = szy2062016.NewSZY206ClientHandler(tcpconn)
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

/**
 * 读到的数据
 *
 */
type SZY2062016ReadData struct {
	Tag     string `json:"tag"`
	MeterId string `json:"meterId"`
	Value   int64  `json:"value"`
}

func (gw *SZY206_2016_MasterGateway) work(handler *szy2062016.SZY206ClientHandler) {
	for {
		select {
		case <-gw.Ctx.Done():
			return
		default:
		}
		SZY2062016ReadDataList := []SZY2062016ReadData{}
		for _, DataPoint := range gw.DataPoints {
			lastTimes := uint64(time.Now().UnixMilli())
			NewValue := intercache.CacheValue{
				UUID:          DataPoint.UUID,
				LastFetchTime: lastTimes,
				Value:         "0",
			}
			MeterSn, err1 := utils.HexStringToBytes(DataPoint.MeterId)
			if err1 != nil {
				glogger.GLogger.Error(err1)
				NewValue.Status = 0
				NewValue.ErrMsg = err1.Error()
				intercache.SetValue(gw.PointId, DataPoint.UUID, NewValue)
				continue
			}
			if len(MeterSn) != 6 {
				glogger.GLogger.Error("invalid MeterId:", DataPoint.MeterId)
				NewValue.Status = 0
				NewValue.ErrMsg = string("invalid MeterId:" + DataPoint.MeterId)
				intercache.SetValue(gw.PointId, DataPoint.UUID, NewValue)
				continue
			}
			Address := utils.ByteReverse(MeterSn)
			frame := szy2062016.SZY206Frame0x00{
				Start1:     szy2062016.CTRL_CODE_FRAME_START,
				DataLength: 0x07,
				Start2:     szy2062016.CTRL_CODE_FRAME_START,
				CtrlCode:   szy2062016.SetControlCode(1),
				Address:    [5]byte{Address[0], Address[1], Address[2], Address[3], Address[4]},
				DataArea:   []byte{},
				End:        szy2062016.CTRL_CODE_FRAME_END,
			}
			Bytes, err2 := frame.Encode()
			if err2 != nil {
				glogger.GLogger.Error(err2)
				NewValue.Status = 0
				NewValue.ErrMsg = err2.Error()
				intercache.SetValue(gw.PointId, DataPoint.UUID, NewValue)
				continue
			}
			Resp, errRequest := handler.Request(Bytes)
			if errRequest != nil {
				glogger.GLogger.Error(errRequest)
				NewValue.Status = 0
				NewValue.ErrMsg = errRequest.Error()
				intercache.SetValue(gw.PointId, DataPoint.UUID, NewValue)
				continue
			}
			DLT645Frame0x11, errDecode := handler.DecodeSZY206Frame0x00(Resp)
			if errDecode != nil {
				glogger.GLogger.Error(errDecode)
				NewValue.Status = 0
				NewValue.ErrMsg = errDecode.Error()
				intercache.SetValue(gw.PointId, DataPoint.UUID, NewValue)
				continue
			}
			glogger.GLogger.Debug(DLT645Frame0x11.String())

			Value, errValue := DLT645Frame0x11.GetData()
			if errValue != nil {
				glogger.GLogger.Error(errValue)
				NewValue.Status = 0
				NewValue.ErrMsg = errValue.Error()
				intercache.SetValue(gw.PointId, DataPoint.UUID, NewValue)
				continue
			}
			NewValue.Value = Value
			NewValue.Status = 1
			NewValue.ErrMsg = ""
			intercache.SetValue(gw.PointId, DataPoint.UUID, NewValue)
			if !*gw.mainConfig.CommonConfig.BatchRequest {
				if bytes, err := json.Marshal(SZY2062016ReadData{
					MeterId: DataPoint.MeterId,
					Value:   Value,
				}); err != nil {
					glogger.GLogger.Error(err)
				} else {
					glogger.GLogger.Debug(string(bytes))
					gw.RuleEngine.WorkDevice(gw.Details(), string(bytes))
				}
			} else {
				SZY2062016ReadDataList = append(SZY2062016ReadDataList, SZY2062016ReadData{
					MeterId: DataPoint.MeterId,
					Tag:     DataPoint.Tag,
					Value:   Value,
				})
			}
			time.Sleep(time.Duration(DataPoint.Frequency) * time.Millisecond)
		}
		if *gw.mainConfig.CommonConfig.BatchRequest {
			if len(SZY2062016ReadDataList) > 0 {
				if bytes, err := json.Marshal(SZY2062016ReadDataList); err != nil {
					glogger.GLogger.Error(err)
				} else {
					glogger.GLogger.Debug(string(bytes))
					gw.RuleEngine.WorkDevice(gw.Details(), string(bytes))
				}
			}
		}
	}
}
func (gw *SZY206_2016_MasterGateway) Status() typex.DeviceState {
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

func (gw *SZY206_2016_MasterGateway) Stop() {
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

func (gw *SZY206_2016_MasterGateway) Details() *typex.Device {
	return gw.RuleEngine.GetDevice(gw.PointId)
}

func (gw *SZY206_2016_MasterGateway) SetState(status typex.DeviceState) {
	gw.status = status
}

func (gw *SZY206_2016_MasterGateway) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}

func (gw *SZY206_2016_MasterGateway) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return []byte{}, nil
}
