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
	"errors"
	"time"

	"github.com/hootrhino/rhilex/component/intercache"
	"github.com/hootrhino/rhilex/component/interdb"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/resconfig"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

// MBusData 用于封装M-Bus采集的数据
type MBusDataPoint struct {
	UUID         string `json:"uuid"`
	DeviceID     string `json:"deviceId"`
	SlaverId     string `json:"slaverId"`
	Type         string `json:"Type"`
	Manufacturer string `json:"Manufacturer"`
	Tag          string `json:"tag"`
	Alias        string `json:"alias"`
	Frequency    int64  `json:"frequency"`
	DataLength   int    `json:"dataLength"`
}
type MBusEn13433MasterGatewayCommonConfig struct {
	Mode         string `json:"mode" validate:"required"`
	AutoRequest  *bool  `json:"autoRequest" validate:"required"`
	BatchRequest *bool  `json:"batchRequest" validate:"required"`
}

type MBusConfig struct {
	HostConfig resconfig.HostConfig `json:"hostConfig"`
}

type MBusEn13433MasterGatewayMainConfig struct {
	CommonConfig  MBusEn13433MasterGatewayCommonConfig `json:"commonConfig" validate:"required"`
	MBusConfig    MBusConfig                           `json:"MBusConfig"`
	UartConfig    resconfig.UartConfig                 `json:"uartConfig"`
	CecollaConfig resconfig.CecollaConfig              `json:"cecollaConfig"`
	AlarmConfig   resconfig.AlarmConfig                `json:"alarmConfig"`
}

/**
 *
 * Mbus
 */

type MBusEn13433MasterGateway struct {
	typex.XStatus
	status         typex.DeviceState
	mainConfig     MBusEn13433MasterGatewayMainConfig
	MBusDataPoints map[string]MBusDataPoint
}

func NewMBusEn13433MasterGateway(e typex.Rhilex) typex.XDevice {
	gw := new(MBusEn13433MasterGateway)
	gw.RuleEngine = e
	gw.mainConfig = MBusEn13433MasterGatewayMainConfig{
		CommonConfig: MBusEn13433MasterGatewayCommonConfig{
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
		MBusConfig: MBusConfig{
			HostConfig: resconfig.HostConfig{
				Host: "127.0.0.1",
				Port: 10065,
			},
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
	gw.MBusDataPoints = map[string]MBusDataPoint{}
	return gw
}

func (gw *MBusEn13433MasterGateway) Init(devId string, configMap map[string]interface{}) error {
	gw.PointId = devId
	intercache.RegisterSlot(gw.PointId)

	if err := utils.BindSourceConfig(configMap, &gw.mainConfig); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	if !utils.SContains([]string{"UART", "TCP"}, gw.mainConfig.CommonConfig.Mode) {
		return errors.New("unsupported mode, only can be one of 'TCP' or 'UART'")
	}
	var ModbusPointList []MBusDataPoint
	PointLoadErr := interdb.InterDb().Table("m_mbus_data_points").
		Where("device_uuid=?", devId).Find(&ModbusPointList).Error
	if PointLoadErr != nil {
		return PointLoadErr
	}
	LastFetchTime := uint64(time.Now().UnixMilli())
	for _, MbusPoint := range ModbusPointList {
		if MbusPoint.Frequency < 1 {
			return errors.New("'frequency' must grate than 50 millisecond")
		}
		gw.MBusDataPoints[MbusPoint.UUID] = MBusDataPoint{
			UUID:         MbusPoint.UUID,
			DeviceID:     MbusPoint.DeviceID,
			SlaverId:     MbusPoint.SlaverId,
			Type:         MbusPoint.Type,
			Manufacturer: MbusPoint.Manufacturer,
			Tag:          MbusPoint.Tag,
			Alias:        MbusPoint.Alias,
			DataLength:   MbusPoint.DataLength,
		}
		intercache.SetValue(gw.PointId, MbusPoint.UUID, intercache.CacheValue{
			UUID:          MbusPoint.UUID,
			Status:        0,
			LastFetchTime: LastFetchTime,
			Value:         "",
			ErrMsg:        "--",
		})
	}
	return nil
}

func (gw *MBusEn13433MasterGateway) Start(cctx typex.CCTX) error {
	gw.Ctx = cctx.Ctx
	gw.CancelCTX = cctx.CancelCTX

	gw.status = typex.DEV_UP
	return nil
}

func (gw *MBusEn13433MasterGateway) Status() typex.DeviceState {
	return gw.status
}

func (gw *MBusEn13433MasterGateway) Stop() {
	gw.status = typex.DEV_DOWN
	if gw.CancelCTX != nil {
		gw.CancelCTX()
	}
	intercache.UnRegisterSlot(gw.PointId)
}

func (gw *MBusEn13433MasterGateway) Details() *typex.Device {
	return gw.RuleEngine.GetDevice(gw.PointId)
}

func (gw *MBusEn13433MasterGateway) SetState(status typex.DeviceState) {
	gw.status = status
}

func (gw *MBusEn13433MasterGateway) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}

func (gw *MBusEn13433MasterGateway) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return []byte{}, nil
}
