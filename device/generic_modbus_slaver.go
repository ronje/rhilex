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
	"fmt"
	"time"

	mbserver "github.com/hootrhino/gomodbus-server"
	serial "github.com/hootrhino/goserial"
	"github.com/hootrhino/rhilex/common"
	"github.com/hootrhino/rhilex/component/hwportmanager"
	"github.com/hootrhino/rhilex/component/intercache"
	"github.com/hootrhino/rhilex/component/interdb"
	modbus_device "github.com/hootrhino/rhilex/device/modbus"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

type ModbusSlaverCommonConfig struct {
	Mode string `json:"mode" validate:"required"` // UART | TCP
}
type ModbusSlaverConfig struct {
	CommonConfig ModbusSlaverCommonConfig `json:"commonConfig" validate:"required"`
	HostConfig   common.HostConfig        `json:"hostConfig"`
	PortUuid     string                   `json:"portUuid"`
}

type ModbusSlaver struct {
	typex.XStatus
	status       typex.DeviceState
	mainConfig   ModbusSlaverConfig
	hwPortConfig hwportmanager.UartConfig
	registers    map[string]*common.RegisterRW
	server       *mbserver.Server
}

func NewGenericModbusSlaver(e typex.Rhilex) typex.XDevice {
	mdev := new(ModbusSlaver)
	mdev.RuleEngine = e
	mdev.mainConfig = ModbusSlaverConfig{
		CommonConfig: ModbusSlaverCommonConfig{Mode: "TCP"},
		PortUuid:     "/dev/ttyS0",
		HostConfig:   common.HostConfig{Host: "0.0.0.0", Port: 1502, Timeout: 3000},
	}
	mdev.registers = map[string]*common.RegisterRW{}
	mdev.status = typex.DEV_DOWN

	return mdev
}

func (mdev *ModbusSlaver) Init(devId string, configMap map[string]interface{}) error {
	mdev.PointId = devId
	intercache.RegisterSlot(mdev.PointId)
	if err := utils.BindSourceConfig(configMap, &mdev.mainConfig); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	if !utils.SContains([]string{"UART", "TCP"}, mdev.mainConfig.CommonConfig.Mode) {
		return fmt.Errorf("unsupported mode, only can be one of 'TCP' or 'UART'")
	}
	// 合并数据库里面的点位表
	var ModbusPointList []modbus_device.ModbusPoint
	modbusPointLoadErr := interdb.DB().Table("m_modbus_data_points").
		Where("device_uuid=?", devId).Find(&ModbusPointList).Error
	if modbusPointLoadErr != nil {
		return modbusPointLoadErr
	}
	for _, ModbusPoint := range ModbusPointList {
		mdev.registers[ModbusPoint.UUID] = &common.RegisterRW{
			UUID:      ModbusPoint.UUID,
			Tag:       ModbusPoint.Tag,
			Alias:     ModbusPoint.Alias,
			Function:  ModbusPoint.Function,
			SlaverId:  ModbusPoint.SlaverId,
			Address:   ModbusPoint.Address,
			Quantity:  ModbusPoint.Quantity,
			Frequency: ModbusPoint.Frequency,
			DataType:  ModbusPoint.DataType,
			DataOrder: ModbusPoint.DataOrder,
			Weight:    ModbusPoint.Weight,
		}
		LastFetchTime := uint64(time.Now().UnixMilli())
		intercache.SetValue(mdev.PointId, ModbusPoint.UUID, intercache.CacheValue{
			UUID:          ModbusPoint.UUID,
			LastFetchTime: LastFetchTime,
		})
	}
	if mdev.mainConfig.CommonConfig.Mode == "UART" {
		hwPort, err := hwportmanager.GetHwPort(mdev.mainConfig.PortUuid)
		if err != nil {
			return err
		}
		if hwPort.Busy {
			return fmt.Errorf("UART is busying now, Occupied By:%s", hwPort.OccupyBy)
		}
		switch tCfg := hwPort.Config.(type) {
		case hwportmanager.UartConfig:
			mdev.hwPortConfig = tCfg
		default:
			return fmt.Errorf("Invalid config:%s", hwPort.Config)
		}
	}
	return nil
}

func (mdev *ModbusSlaver) Start(cctx typex.CCTX) error {
	mdev.Ctx = cctx.Ctx
	mdev.CancelCTX = cctx.CancelCTX
	mdev.server = mbserver.NewServerWithContext(mdev.Ctx)
	mdev.server.SetLogger(glogger.Logrus)
	// 点位, 需要和数据库关联起来
	mdev.server.HoldingRegisters = []uint16{
		0x01, 0x02, 0x03, 0x04, 0x05,
		0x21, 0x22, 0x23, 0x24, 0x25,
	}
	mdev.server.InputRegisters = []uint16{
		0x01, 0x01, 0x01, 0x01, 0x01,
		0x01, 0x01, 0x01, 0x01, 0x01,
	}
	mdev.server.DiscreteInputs = []byte{
		0x01, 0x02, 0x03, 0x04, 0x05,
		0x21, 0x22, 0x23, 0x24, 0x25,
	}
	mdev.server.Coils = []byte{
		0x01, 0x01, 0x01, 0x01, 0x01,
		0x01, 0x01, 0x01, 0x01, 0x01,
	}
	mdev.server.SetOnRequest(func(s *mbserver.Server, frame mbserver.Framer) {
		glogger.GLogger.Debug("Received Modbus Request:", frame.GetFunction(), frame.GetData())
	})
	if mdev.mainConfig.CommonConfig.Mode == "UART" {
		hwPort, err := hwportmanager.GetHwPort(mdev.mainConfig.PortUuid)
		if err != nil {
			return err
		}
		if hwPort.Busy {
			return fmt.Errorf("UART is busying now, Occupied By:%s", hwPort.OccupyBy)
		}
		err1 := mdev.server.ListenRTU(&serial.Config{
			Address:  mdev.hwPortConfig.Uart,
			BaudRate: mdev.hwPortConfig.BaudRate,
			DataBits: mdev.hwPortConfig.DataBits,
			Parity:   mdev.hwPortConfig.Parity,
			StopBits: mdev.hwPortConfig.StopBits,
			Timeout:  time.Duration(mdev.hwPortConfig.Timeout) * (time.Millisecond),
		})
		if err1 != nil {
			return err1
		}

	}
	if mdev.mainConfig.CommonConfig.Mode == "TCP" {
		err2 := mdev.server.ListenTCP(fmt.Sprintf("%s:%d",
			mdev.mainConfig.HostConfig.Host, mdev.mainConfig.HostConfig.Port))
		if err2 != nil {
			return err2
		}
	}
	mdev.status = typex.DEV_UP
	return nil
}

func (mdev *ModbusSlaver) Status() typex.DeviceState {
	return typex.DEV_UP
}

func (mdev *ModbusSlaver) Stop() {
	mdev.status = typex.DEV_DOWN
	if mdev.CancelCTX != nil {
		mdev.CancelCTX()
	}
	if mdev.server != nil {
		mdev.server.Close()
	}
	if mdev.mainConfig.CommonConfig.Mode == "UART" {
		hwportmanager.FreeInterfaceBusy(mdev.mainConfig.PortUuid)
	}
	intercache.UnRegisterSlot(mdev.PointId) // 卸载点位表
}

func (mdev *ModbusSlaver) Details() *typex.Device {
	return mdev.RuleEngine.GetDevice(mdev.PointId)
}

func (mdev *ModbusSlaver) SetState(status typex.DeviceState) {
	mdev.status = status
}

func (mdev *ModbusSlaver) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}

func (mdev *ModbusSlaver) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return []byte{}, nil
}

func (mdev *ModbusSlaver) OnRead(cmd []byte, data []byte) (int, error) {
	return 0, nil
}

func (mdev *ModbusSlaver) OnWrite(cmd []byte, b []byte) (int, error) {
	return 0, nil
}
