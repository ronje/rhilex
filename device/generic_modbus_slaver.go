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
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	mbserver "github.com/hootrhino/gomodbus-server"
	serial "github.com/hootrhino/goserial"
	"github.com/hootrhino/rhilex/common"
	"github.com/hootrhino/rhilex/component/hwportmanager"
	"github.com/hootrhino/rhilex/component/intercache"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

type ModbusSlaverCommonConfig struct {
	Mode         string `json:"mode" validate:"required"` // UART | TCP
	MaxRegisters int    `json:"maxRegisters" validate:"required"`
	SlaverId     int16  `json:"slaverId" validate:"required"`
}
type ModbusSlaverConfig struct {
	CommonConfig ModbusSlaverCommonConfig `json:"commonConfig" validate:"required"`
	HostConfig   common.HostConfig        `json:"hostConfig"`
	PortUuid     string                   `json:"portUuid"`
}

type ModbusSlaver struct {
	typex.XStatus
	status           typex.DeviceState
	mainConfig       ModbusSlaverConfig
	hwPortConfig     hwportmanager.UartConfig
	registers        map[string]*common.RegisterRW
	server           *mbserver.Server
	HoldingRegisters []uint16 // [5] = WriteSingleCoil
	InputRegisters   []uint16 // [6] = WriteHoldingRegister
	DiscreteInputs   []byte   // [15] = WriteMultipleCoils
	Coils            []byte   // [16] = WriteHoldingRegisters
}

func NewGenericModbusSlaver(e typex.Rhilex) typex.XDevice {
	mdev := new(ModbusSlaver)
	mdev.RuleEngine = e
	mdev.mainConfig = ModbusSlaverConfig{
		CommonConfig: ModbusSlaverCommonConfig{Mode: "TCP", MaxRegisters: 64, SlaverId: 1},
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
	mdev.HoldingRegisters = make([]uint16, mdev.mainConfig.CommonConfig.MaxRegisters)
	mdev.InputRegisters = make([]uint16, mdev.mainConfig.CommonConfig.MaxRegisters)
	mdev.DiscreteInputs = make([]byte, mdev.mainConfig.CommonConfig.MaxRegisters)
	mdev.Coils = make([]byte, mdev.mainConfig.CommonConfig.MaxRegisters)
	for i := 0; i < mdev.mainConfig.CommonConfig.MaxRegisters; i++ {
		CoilUUID := fmt.Sprintf("%s_Coils:%d", mdev.PointId, i)
		HoldingRegisterUUID := fmt.Sprintf("%s_HoldingRegisters:%d", mdev.PointId, i)
		InputRegisterUUID := fmt.Sprintf("%s_InputRegisters:%d", mdev.PointId, i)
		DiscreteInputUUID := fmt.Sprintf("%s_DiscreteInputs:%d", mdev.PointId, i)
		//
		LastFetchTime := uint64(time.Now().UnixMilli())
		intercache.SetValue(mdev.PointId, HoldingRegisterUUID, intercache.CacheValue{
			UUID:          HoldingRegisterUUID,
			LastFetchTime: LastFetchTime,
			Value:         "0",
		})
		intercache.SetValue(mdev.PointId, InputRegisterUUID, intercache.CacheValue{
			UUID:          InputRegisterUUID,
			LastFetchTime: LastFetchTime,
			Value:         "0",
		})
		intercache.SetValue(mdev.PointId, DiscreteInputUUID, intercache.CacheValue{
			UUID:          DiscreteInputUUID,
			LastFetchTime: LastFetchTime,
			Value:         "0",
		})
		intercache.SetValue(mdev.PointId, CoilUUID, intercache.CacheValue{
			UUID:          CoilUUID,
			LastFetchTime: LastFetchTime,
			Value:         "0",
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

type modbusSlaverValue struct {
	Register      int16  `json:"register"`
	SlaverId      byte   `json:"slaverId"`
	LastFetchTime uint64 `json:"lastFetchTime"`
	Value         string `json:"value"`
}

func (O modbusSlaverValue) String() string {
	bytes, _ := json.Marshal(O)
	return string(bytes)
}
func uint16ToBytes(value uint16) []byte {
	return []byte{byte(value >> 8), byte(value)}
}
func (mdev *ModbusSlaver) Start(cctx typex.CCTX) error {
	mdev.Ctx = cctx.Ctx
	mdev.CancelCTX = cctx.CancelCTX
	mdev.server = mbserver.NewServerWithContext(mdev.Ctx)
	mdev.server.SetLogger(glogger.Logrus)
	// 点位, 需要和数据库关联起来
	mdev.server.InputRegisters = mdev.InputRegisters
	mdev.server.DiscreteInputs = mdev.DiscreteInputs
	mdev.server.Coils = mdev.Coils
	mdev.server.HoldingRegisters = mdev.HoldingRegisters
	// [5] = WriteSingleCoil
	// [6] = WriteHoldingRegister
	// [15] = WriteMultipleCoils
	// [16] = WriteHoldingRegisters
	mdev.server.SetOnRequest(func(s *mbserver.Server, frame mbserver.Framer) {
		FunCode := frame.GetFunction()
		register, numRegs, endRegister := getRegisterAddressAndNumber(frame)
		glogger.GLogger.Debug("Modbus OnRequest: ", register, numRegs, endRegister)
		if register > mdev.mainConfig.CommonConfig.MaxRegisters {
			glogger.GLogger.Error("exceed MaxRegisters:", register, numRegs, endRegister)
			return
		}
		if FunCode == 5 { // 更新线圈
			LastFetchTime := uint64(time.Now().UnixMilli())
			LastValue := mdev.Coils[register]
			UUID := fmt.Sprintf("%s_Coils:%d", mdev.PointId, register)
			CacheValue := intercache.CacheValue{
				UUID:          UUID,
				LastFetchTime: LastFetchTime,
				Value:         "0", // 更新后的值
			}
			if LastValue == 0xFF {
				CacheValue.Value = "1"
			}
			intercache.SetValue(mdev.PointId, UUID, CacheValue)
			if bytes, errMarshal := json.Marshal(modbusSlaverValue{
				Register:      int16(register),
				SlaverId:      byte(mdev.mainConfig.CommonConfig.SlaverId),
				LastFetchTime: LastFetchTime,
				Value: func() string {
					if CacheValue.Value == "1" {
						return "1"
					}
					return "0"
				}(),
			}); errMarshal != nil {
				glogger.GLogger.Error(errMarshal)
			} else {
				mdev.RuleEngine.WorkDevice(mdev.Details(), string(bytes))
			}
		}
		if FunCode == 6 { // 更新HoldingRegisters
			LastFetchTime := uint64(time.Now().UnixMilli())
			LastValue := mdev.HoldingRegisters[register]
			UUID := fmt.Sprintf("%s_HoldingRegisters:%d", mdev.PointId, register)
			intercache.SetValue(mdev.PointId, UUID, intercache.CacheValue{
				UUID:          UUID,
				LastFetchTime: LastFetchTime,
				Value:         hex.EncodeToString(uint16ToBytes(LastValue)),
			})
			if bytes, errMarshal := json.Marshal(modbusSlaverValue{
				Register:      int16(register),
				SlaverId:      byte(mdev.mainConfig.CommonConfig.SlaverId),
				LastFetchTime: LastFetchTime,
				Value:         hex.EncodeToString(uint16ToBytes(LastValue)),
			}); errMarshal != nil {
				glogger.GLogger.Error(errMarshal)
			} else {
				mdev.RuleEngine.WorkDevice(mdev.Details(), string(bytes))
			}
		}
		// 15 16暂时不支持
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

/*
*
* 提取Modbus报文数据
*
 */
func getRegisterAddressAndNumber(frame mbserver.Framer) (int, int, int) {
	data := frame.GetData()
	register := int(binary.BigEndian.Uint16(data[0:2]))
	numRegs := int(binary.BigEndian.Uint16(data[2:4]))
	endRegister := register + numRegs
	return register, numRegs, endRegister
}

/*
*
* 提取Modbus报文数据
*
 */
func getRegisterAddressAndValue(frame mbserver.Framer) (int, uint16) {
	data := frame.GetData()
	register := int(binary.BigEndian.Uint16(data[0:2]))
	value := binary.BigEndian.Uint16(data[2:4])
	return register, value
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
