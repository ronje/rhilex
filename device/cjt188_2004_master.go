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

	serial "github.com/hootrhino/goserial"
	"github.com/hootrhino/rhilex/resconfig"
	"github.com/hootrhino/rhilex/component/intercache"
	"github.com/hootrhino/rhilex/component/interdb"
	"github.com/hootrhino/rhilex/device/cjt1882004"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

// CJT188_2004_Data 用于封装M-Bus采集的数据
type CJT188_2004_DataPoint struct {
	UUID      string `json:"uuid"`
	MeterId   string `json:"meterId"`
	Tag       string `json:"tag"`
	Alias     string `json:"alias"`
	Frequency int64  `json:"frequency"`
}
type CJT188_2004_MasterGatewayresconfigConfig struct {
	Mode         string `json:"mode" validate:"required"` // UART | TCP
	AutoRequest  *bool  `json:"autoRequest" validate:"required"`
	BatchRequest *bool  `json:"batchRequest" validate:"required"`
}

type CJT188_2004_MasterGatewayMainConfig struct {
	resconfigConfig  CJT188_2004_MasterGatewayresconfigConfig `json:"resconfigConfig" validate:"required"`
	HostConfig    resconfig.HostConfig                     `json:"hostConfig"`
	UartConfig    resconfig.UartConfig                     `json:"uartConfig"`
	CecollaConfig resconfig.CecollaConfig                  `json:"cecollaConfig"`
}

/**
 *
 * CJT188_2004_
 */

type CJT188_2004_MasterGateway struct {
	typex.XStatus
	status      typex.DeviceState
	mainConfig  CJT188_2004_MasterGatewayMainConfig
	DataPoints  map[string]CJT188_2004_DataPoint
	uartHandler *cjt1882004.CJT188ClientHandler
	tcpHandler  *cjt1882004.CJT188ClientHandler
}

func NewCJT188_2004_MasterGateway(e typex.Rhilex) typex.XDevice {
	gw := new(CJT188_2004_MasterGateway)
	gw.RuleEngine = e
	gw.mainConfig = CJT188_2004_MasterGatewayMainConfig{
		resconfigConfig: CJT188_2004_MasterGatewayresconfigConfig{
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
	gw.DataPoints = map[string]CJT188_2004_DataPoint{}
	return gw
}

func (gw *CJT188_2004_MasterGateway) Init(devId string, configMap map[string]interface{}) error {
	gw.PointId = devId
	intercache.RegisterSlot(gw.PointId)

	if err := utils.BindSourceConfig(configMap, &gw.mainConfig); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	if !utils.SContains([]string{"UART", "TCP"}, gw.mainConfig.resconfigConfig.Mode) {
		return errors.New("unsupported mode, only can be one of 'TCP' or 'UART'")
	}
	// CheckSerialBusy
	if err := gw.mainConfig.UartConfig.Validate(); err != nil {
		return nil
	}
	var DLT645_ModbusPointList []CJT188_2004_DataPoint
	PointLoadErr := interdb.DB().Table("m_cjt1882004_data_points").
		Where("device_uuid=?", devId).Find(&DLT645_ModbusPointList).Error
	if PointLoadErr != nil {
		return PointLoadErr
	}
	LastFetchTime := uint64(time.Now().UnixMilli())
	for _, CJT188_2004_Point := range DLT645_ModbusPointList {
		if CJT188_2004_Point.Frequency < 1 {
			return errors.New("'frequency' must grate than 50 millisecond")
		}
		gw.DataPoints[CJT188_2004_Point.UUID] = CJT188_2004_DataPoint{
			UUID:      CJT188_2004_Point.UUID,
			MeterId:   CJT188_2004_Point.MeterId,
			Tag:       CJT188_2004_Point.Tag,
			Alias:     CJT188_2004_Point.Alias,
			Frequency: CJT188_2004_Point.Frequency,
		}
		intercache.SetValue(gw.PointId, CJT188_2004_Point.UUID, intercache.CacheValue{
			UUID:          CJT188_2004_Point.UUID,
			Status:        0,
			LastFetchTime: LastFetchTime,
			Value:         "0",
			ErrMsg:        "--",
		})
	}
	return nil
}

func (gw *CJT188_2004_MasterGateway) Start(cctx typex.CCTX) error {
	gw.Ctx = cctx.Ctx
	gw.CancelCTX = cctx.CancelCTX
	if gw.mainConfig.resconfigConfig.Mode == "UART" {
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
		gw.uartHandler = cjt1882004.NewCJT188ClientHandler(serialPort)
		gw.uartHandler.SetLogger(glogger.Logrus)
		if *gw.mainConfig.resconfigConfig.AutoRequest {
			go gw.work(gw.uartHandler)
		}
		goto END
	}
	if gw.mainConfig.resconfigConfig.Mode == "TCP" {
		tcpconn, errDial := net.Dial("tcp",
			fmt.Sprintf("%s:%d", gw.mainConfig.HostConfig.Host,
				gw.mainConfig.HostConfig.Port))
		if errDial != nil {
			return errDial
		}
		gw.tcpHandler = cjt1882004.NewCJT188ClientHandler(tcpconn)
		gw.uartHandler.SetLogger(glogger.Logrus)
		if *gw.mainConfig.resconfigConfig.AutoRequest {
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
type CJT1882004ReadData struct {
	Tag     string  `json:"tag"`
	MeterId string  `json:"meterId"`
	Value   []int64 `json:"value"`
}

func (gw *CJT188_2004_MasterGateway) work(handler *cjt1882004.CJT188ClientHandler) {
	for {
		select {
		case <-gw.Ctx.Done():
			return
		default:
		}
		cjt1882004ReadDataList := []CJT1882004ReadData{}
		for _, DataPoint := range gw.DataPoints {
			lastTimes := uint64(time.Now().UnixMilli())
			NewValue := intercache.CacheValue{
				UUID:          DataPoint.UUID,
				LastFetchTime: lastTimes,
				Value:         "0",
			}
			MeterSn, err1 := hex.DecodeString(DataPoint.MeterId)
			if err1 != nil {
				glogger.GLogger.Error(err1)
				NewValue.Status = 0
				NewValue.ErrMsg = err1.Error()
				intercache.SetValue(gw.PointId, DataPoint.UUID, NewValue)
				continue
			}
			if len(MeterSn) != 7 {
				glogger.GLogger.Error("invalid MeterId:", DataPoint.MeterId)
				NewValue.Status = 0
				NewValue.ErrMsg = string("invalid MeterId:" + DataPoint.MeterId)
				intercache.SetValue(gw.PointId, DataPoint.UUID, NewValue)
				continue
			}
			frame := cjt1882004.CJT188Frame0x01{
				Start:        cjt1882004.CTRL_CODE_FRAME_START,
				MeterType:    0x10,
				Address:      [7]byte{MeterSn[6], MeterSn[5], MeterSn[4], MeterSn[3], MeterSn[2], MeterSn[1], MeterSn[0]},
				CtrlCode:     cjt1882004.CTRL_CODE_READ_DATA,
				DataLength:   0x03,
				DataType:     [2]byte{0x1F, 0x90},
				DataArea:     []byte{},
				SerialNumber: 0x00,
				End:          cjt1882004.CTRL_CODE_FRAME_END,
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
			DLT645Frame0x11, errDecode := handler.DecodeCJT188Frame0x01Response(Resp)
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
				continue
			}
			NewValue.Value = Value
			NewValue.Status = 1
			NewValue.ErrMsg = ""
			intercache.SetValue(gw.PointId, DataPoint.UUID, NewValue)
			if !*gw.mainConfig.resconfigConfig.BatchRequest {
				if bytes, err := json.Marshal(CJT1882004ReadData{
					MeterId: DataPoint.MeterId,
					Tag:     DataPoint.Tag,
					Value:   Value,
				}); err != nil {
					glogger.GLogger.Error(err)
				} else {
					glogger.GLogger.Debug(string(bytes))
					gw.RuleEngine.WorkDevice(gw.Details(), string(bytes))
				}
			} else {
				cjt1882004ReadDataList = append(cjt1882004ReadDataList, CJT1882004ReadData{
					MeterId: DataPoint.MeterId,
					Tag:     DataPoint.Tag,
					Value:   Value,
				})
			}
			time.Sleep(time.Duration(DataPoint.Frequency) * time.Millisecond)
		}
		if *gw.mainConfig.resconfigConfig.BatchRequest {
			if len(cjt1882004ReadDataList) > 0 {
				if bytes, err := json.Marshal(cjt1882004ReadDataList); err != nil {
					glogger.GLogger.Error(err)
				} else {
					glogger.GLogger.Debug(string(bytes))
					gw.RuleEngine.WorkDevice(gw.Details(), string(bytes))
				}
			}
		}
	}
}
func (gw *CJT188_2004_MasterGateway) Status() typex.DeviceState {
	if gw.mainConfig.resconfigConfig.Mode == "UART" {
		if gw.uartHandler == nil {
			return typex.DEV_DOWN
		}
	}
	if gw.mainConfig.resconfigConfig.Mode == "TCP" {
		if gw.tcpHandler == nil {
			return typex.DEV_DOWN
		}
	}
	return gw.status
}

func (gw *CJT188_2004_MasterGateway) Stop() {
	gw.status = typex.DEV_DOWN
	if gw.CancelCTX != nil {
		gw.CancelCTX()
	}
	if gw.mainConfig.resconfigConfig.Mode == "UART" {
		if gw.uartHandler != nil {
			gw.uartHandler.Close()
		}
	}
	if gw.mainConfig.resconfigConfig.Mode == "TCP" {
		if gw.tcpHandler != nil {
			gw.tcpHandler.Close()
		}
	}
	intercache.UnRegisterSlot(gw.PointId)
}

func (gw *CJT188_2004_MasterGateway) Details() *typex.Device {
	return gw.RuleEngine.GetDevice(gw.PointId)
}

func (gw *CJT188_2004_MasterGateway) SetState(status typex.DeviceState) {
	gw.status = status
}

func (gw *CJT188_2004_MasterGateway) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}

func (gw *CJT188_2004_MasterGateway) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return []byte{}, nil
}
