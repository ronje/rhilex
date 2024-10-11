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

	"github.com/hootrhino/rhilex/common"
	"github.com/hootrhino/rhilex/component/intercache"
	"github.com/hootrhino/rhilex/component/interdb"
	mbus_device "github.com/hootrhino/rhilex/device/mbus"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

type MBusMasterGatewayCommonConfig struct {
	Mode         string `json:"mode" validate:"required"`
	AutoRequest  *bool  `json:"autoRequest" validate:"required"`
	BatchRequest *bool  `json:"batchRequest" validate:"required"`
}

type MBusConfig struct {
	HostConfig common.HostConfig `json:"hostConfig"`
}

type MBusMasterGatewayMainConfig struct {
	CommonConfig MBusMasterGatewayCommonConfig `json:"commonConfig"`
	MBusConfig   MBusConfig                    `json:"MBusConfig"`
	UartConfig   common.UartConfig             `json:"uartConfig"`
}

/**
 *
 * Mbus
 */

type MBusMasterGateway struct {
	typex.XStatus
	status         typex.DeviceState
	mainConfig     MBusMasterGatewayMainConfig
	MBusDataPoints map[string]mbus_device.MBusDataPoint
}

func NewMBusMasterGateway(e typex.Rhilex) typex.XDevice {
	gw := new(MBusMasterGateway)
	gw.RuleEngine = e
	gw.mainConfig = MBusMasterGatewayMainConfig{
		CommonConfig: MBusMasterGatewayCommonConfig{
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
			HostConfig: common.HostConfig{
				Host: "127.0.0.1",
				Port: 10065,
			},
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
	gw.MBusDataPoints = map[string]mbus_device.MBusDataPoint{}
	return gw
}

func (gw *MBusMasterGateway) Init(devId string, configMap map[string]interface{}) error {
	gw.PointId = devId
	intercache.RegisterSlot(gw.PointId)

	if err := utils.BindSourceConfig(configMap, &gw.mainConfig); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	if !utils.SContains([]string{"UART", "TCP"}, gw.mainConfig.CommonConfig.Mode) {
		return errors.New("unsupported mode, only can be one of 'TCP' or 'UART'")
	}
	var ModbusPointList []mbus_device.MBusDataPoint
	PointLoadErr := interdb.DB().Table("m_mbus_data_points").
		Where("device_uuid=?", devId).Find(&ModbusPointList).Error
	if PointLoadErr != nil {
		return PointLoadErr
	}
	LastFetchTime := uint64(time.Now().UnixMilli())
	for _, MbusPoint := range ModbusPointList {
		if MbusPoint.Frequency < 1 {
			return errors.New("'frequency' must grate than 50 millisecond")
		}
		gw.MBusDataPoints[MbusPoint.UUID] = mbus_device.MBusDataPoint{
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
			ErrMsg:        "Loading",
		})
	}
	return nil
}

func (gw *MBusMasterGateway) Start(cctx typex.CCTX) error {
	gw.Ctx = cctx.Ctx
	gw.CancelCTX = cctx.CancelCTX

	gw.status = typex.DEV_UP
	return nil
}

func (gw *MBusMasterGateway) Status() typex.DeviceState {
	// if gw.mainConfig.CommonConfig.Mode == "TCP" {

	// }
	// if gw.mainConfig.CommonConfig.Mode == "UART" {

	// }
	return typex.DEV_UP
}

func (gw *MBusMasterGateway) Stop() {
	gw.status = typex.DEV_DOWN
	if gw.CancelCTX != nil {
		gw.CancelCTX()
	}
	intercache.RegisterSlot(gw.PointId)
}

func (gw *MBusMasterGateway) Details() *typex.Device {
	return gw.RuleEngine.GetDevice(gw.PointId)
}

func (gw *MBusMasterGateway) SetState(status typex.DeviceState) {
	gw.status = status
}

func (gw *MBusMasterGateway) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}

func (gw *MBusMasterGateway) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return []byte{}, nil
}

func (gw *MBusMasterGateway) OnRead(cmd []byte, data []byte) (int, error) {

	return 0, nil
}

func (gw *MBusMasterGateway) OnWrite(cmd []byte, b []byte) (int, error) {
	return 0, nil
}
