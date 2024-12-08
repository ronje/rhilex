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
	"encoding/json"
	"errors"
	"fmt"
	golog "log"
	"sort"
	"strconv"

	"time"

	"github.com/hootrhino/rhilex/component/alarmcenter"
	intercache "github.com/hootrhino/rhilex/component/intercache"
	"github.com/hootrhino/rhilex/component/interdb"
	"github.com/hootrhino/rhilex/device/dmodbus"
	"github.com/hootrhino/rhilex/device/ithings"
	"github.com/hootrhino/rhilex/resconfig"

	modbus "github.com/hootrhino/gomodbus"
	core "github.com/hootrhino/rhilex/config"

	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

type ModbusPoint struct {
	UUID      string  `json:"uuid,omitempty"` // 当UUID为空时新建
	Tag       string  `json:"tag"`
	Alias     string  `json:"alias"`
	Function  int     `json:"function"`
	SlaverId  byte    `json:"slaverId"`
	Address   uint16  `json:"address"`
	Frequency int64   `json:"frequency"`
	Quantity  uint16  `json:"quantity"`
	Value     string  `json:"value,omitempty"` // 运行时数据
	DataType  string  `json:"dataType"`        // 运行时数据
	DataOrder string  `json:"dataOrder"`       // 运行时数据
	Weight    float64 `json:"weight"`          // 权重
}

// 这是个通用Modbus采集器, 主要用来在通用场景下采集数据，因此需要配合规则引擎来使用
//
//	{
//	    "d1":{
//	        "tag":"d1",
//	        "function":3,
//	        "slaverId":1,
//	        "address":0,
//	        "quantity":2,
//	        "value":"..."
//	    },
//	}
type ModbusMasterCommonConfig struct {
	Mode           string `json:"mode" validate:"required"`
	AutoRequest    *bool  `json:"autoRequest" validate:"required"`
	BatchRequest   *bool  `json:"batchRequest" validate:"required"` // 批量采集
	EnableOptimize *bool  `json:"enableOptimize" validate:"required"`
	MaxRegNum      uint16 `json:"maxRegNum" validate:"required"`
}

// Modbus配置参数
type ModbusMasterConfig struct {
	CommonConfig  ModbusMasterCommonConfig `json:"commonConfig" validate:"required"`
	HostConfig    resconfig.HostConfig     `json:"hostConfig"`
	UartConfig    resconfig.UartConfig     `json:"uartConfig"`
	CecollaConfig resconfig.CecollaConfig  `json:"cecollaConfig"`
	AlarmConfig   resconfig.AlarmConfig    `json:"alarmConfig"`
}

type ModbusMasterGroupedTag struct {
	Function  int    `json:"function"`
	SlaverId  byte   `json:"slaverId"`
	Address   uint16 `json:"address"`
	Frequency int64  `json:"frequency"`
	Quantity  uint16 `json:"quantity"`
	Registers map[string]*dmodbus.ModbusRegister
}

func (g *ModbusMasterGroupedTag) String() string {
	tagIds := make([]string, 0, len(g.Registers))
	for k := range g.Registers {
		tagIds = append(tagIds, k)
	}
	str := fmt.Sprintf("func=%v slaveId=%v address=%v quantity=%v frequency=%v tagIds=%v",
		g.Function, g.SlaverId, g.Address, g.Quantity, g.Frequency, tagIds)
	return str
}

type GenericModbusMaster struct {
	typex.XStatus
	status     typex.DeviceState
	RuleEngine typex.Rhilex
	//
	rtuHandler *modbus.RTUClientHandler
	tcpHandler *modbus.TCPClientHandler
	Client     modbus.Client
	//
	mainConfig       ModbusMasterConfig
	retryTimes       int
	Registers        map[string]*dmodbus.ModbusRegister
	RegisterWithTags map[string]*dmodbus.ModbusRegister
	RegisterGroups   []*ModbusMasterGroupedTag
}

/*
*
* 温湿度传感器
*
 */
func NewGenericModbusMaster(e typex.Rhilex) typex.XDevice {
	mdev := new(GenericModbusMaster)
	mdev.RuleEngine = e
	mdev.mainConfig = ModbusMasterConfig{
		CommonConfig: ModbusMasterCommonConfig{
			EnableOptimize: func() *bool {
				b := false
				return &b
			}(),
			AutoRequest: func() *bool {
				b := false
				return &b
			}(),
			BatchRequest: func() *bool {
				b := false
				return &b
			}(),
			MaxRegNum: 32,
		},
		HostConfig: resconfig.HostConfig{
			Host:    "127.0.0.1",
			Port:    502,
			Timeout: 3000,
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
	mdev.Registers = map[string]*dmodbus.ModbusRegister{}
	mdev.RegisterWithTags = map[string]*dmodbus.ModbusRegister{}
	mdev.RegisterGroups = []*ModbusMasterGroupedTag{}
	mdev.Busy = false
	mdev.status = typex.DEV_DOWN
	return mdev
}

//  初始化
func (mdev *GenericModbusMaster) Init(devId string, configMap map[string]interface{}) error {
	mdev.PointId = devId
	mdev.retryTimes = 0
	intercache.RegisterSlot(mdev.PointId)
	if err := utils.BindSourceConfig(configMap, &mdev.mainConfig); err != nil {
		return err
	}
	if !utils.SContains([]string{"UART", "TCP"}, mdev.mainConfig.CommonConfig.Mode) {
		return errors.New("unsupported mode, only can be one of 'TCP' or 'UART'")
	}

	// 合并数据库里面的点位表
	var ModbusPointList []ModbusPoint
	modbusPointLoadErr := interdb.InterDb().Table("m_modbus_data_points").
		Where("device_uuid=?", devId).Find(&ModbusPointList).Error
	if modbusPointLoadErr != nil {
		return modbusPointLoadErr
	}
	LastFetchTime := uint64(time.Now().UnixMilli())
	subDevicesAlias := map[string]string{}
	for _, ModbusPoint := range ModbusPointList {
		// 频率不能太快
		if ModbusPoint.Frequency < 1 {
			return errors.New("'frequency' must grate than 50 millisecond")
		}
		mdev.RegisterWithTags[ModbusPoint.Tag] = &dmodbus.ModbusRegister{
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
		mdev.Registers[ModbusPoint.UUID] = &dmodbus.ModbusRegister{
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
		intercache.SetValue(mdev.PointId, ModbusPoint.UUID, intercache.CacheValue{
			UUID:          ModbusPoint.UUID,
			Status:        0,
			LastFetchTime: LastFetchTime,
			Value:         "0",
			ErrMsg:        "--",
		})
		subDevicesAlias[ModbusPoint.Alias] = ModbusPoint.Alias
	}

	// 子设备上线,推向云边协同
	for _, Alias := range subDevicesAlias {
		if *mdev.mainConfig.CecollaConfig.Enable {
			cecolla := mdev.RuleEngine.GetCecolla(mdev.mainConfig.CecollaConfig.CecollaId)
			if cecolla != nil {
				ProductId, DeviceName, err := ithings.ParseProductInfo(Alias)
				if err != nil {
					glogger.Error(err)
				} else {
					param := ithings.SubDeviceParam{
						Timestamp:  int64(LastFetchTime),
						ProductId:  ProductId,
						DeviceName: DeviceName,
					}
					_, errOnCtrl := cecolla.Cecolla.OnCtrl([]byte("SubDeviceSetOnline"), []byte(param.String()))
					if errOnCtrl != nil {
						glogger.Error(errOnCtrl)
					}
				}
			}
		}
	}
	// 开启优化
	if *mdev.mainConfig.CommonConfig.EnableOptimize {
		rws := make([]*dmodbus.ModbusRegister, len(mdev.Registers))
		idx := 0
		for _, val := range mdev.Registers {
			rws[idx] = val
			idx++
		}
		mdev.RegisterGroups = mdev.groupTags(rws)
		for i, v := range mdev.RegisterGroups {
			glogger.GLogger.Infof("RegisterGroups%v %v", i, v)
		}
	}
	// 清空可能存在的老数据
	if !*mdev.mainConfig.CecollaConfig.Enable {
		intercache.DeleteValue("__CecollaBinding", mdev.mainConfig.CecollaConfig.CecollaId)
	} else {
		value := intercache.GetValue("__CecollaBinding", mdev.mainConfig.CecollaConfig.CecollaId)
		if value.Value == nil {
			intercache.SetValue("__CecollaBinding",
				mdev.mainConfig.CecollaConfig.CecollaId,
				intercache.CacheValue{Value: mdev.PointId})
		} else {
			switch T := value.Value.(type) {
			case string:
				if T != mdev.PointId {
					glogger.GLogger.Errorf("Cecolla already bind to device:%s", value.Value)
					return fmt.Errorf("Cecolla already bind to device:%s", value.Value)
				}
			}
		}
	}

	return nil
}

/*
*
0、分组，Frequency采集时间需要相同
1、寄存器类型分类
2、tag排序
3、限制单次数据采集数量为32个
4、tag address必须连续
*/
func (mdev *GenericModbusMaster) groupTags(registers []*dmodbus.ModbusRegister) []*ModbusMasterGroupedTag {

	sort.Sort(dmodbus.RegisterList(registers))
	result := make([]*ModbusMasterGroupedTag, 0)
	for i := 0; i < len(registers); {
		start := i
		end := i
		cursor := i
		tagGroup := &ModbusMasterGroupedTag{
			Function:  registers[start].Function,
			SlaverId:  registers[start].SlaverId,
			Address:   registers[start].Address,
			Frequency: registers[start].Frequency,
		}
		result = append(result, tagGroup)
		tagGroup.Registers = make(map[string]*dmodbus.ModbusRegister)

		regMaxAddr := uint16(0)
		for end < len(registers) {
			curReg := registers[cursor]
			evaluateReg := registers[end]
			curRegAddr := curReg.Address + curReg.Quantity - 1
			if curRegAddr > regMaxAddr {
				regMaxAddr = curRegAddr
			}
			if tagGroup.SlaverId != evaluateReg.SlaverId {
				break
			}
			if tagGroup.Function != evaluateReg.Function {
				break
			}
			if tagGroup.Frequency != evaluateReg.Frequency {
				break
			}
			if evaluateReg.Address > regMaxAddr+1 {
				break
			}
			totalQuantity := evaluateReg.Address + evaluateReg.Quantity - tagGroup.Address
			if totalQuantity > mdev.mainConfig.CommonConfig.MaxRegNum {
				// 寄存器数量超过单次最大采集寄存器个数
				break
			}
			tagGroup.Registers[evaluateReg.UUID] = evaluateReg
			tagGroup.Quantity = totalQuantity
			cursor = end
			end++
		}
		i = end
	}
	return result
}

// 启动
func (mdev *GenericModbusMaster) Start(cctx typex.CCTX) error {
	mdev.Ctx = cctx.Ctx
	mdev.CancelCTX = cctx.CancelCTX
	mdev.retryTimes = 0
	if mdev.mainConfig.CommonConfig.Mode == "UART" {
		mdev.rtuHandler = modbus.NewRTUClientHandler(mdev.mainConfig.UartConfig.Uart)
		mdev.rtuHandler.BaudRate = mdev.mainConfig.UartConfig.BaudRate
		mdev.rtuHandler.DataBits = mdev.mainConfig.UartConfig.DataBits
		mdev.rtuHandler.Parity = mdev.mainConfig.UartConfig.Parity
		mdev.rtuHandler.StopBits = mdev.mainConfig.UartConfig.StopBits
		mdev.rtuHandler.Timeout = time.Duration(mdev.mainConfig.UartConfig.Timeout) * time.Millisecond
		if core.GlobalConfig.DebugMode {
			mdev.rtuHandler.Logger = golog.New(glogger.GLogger.Writer(),
				"Modbus RTU Mode: "+mdev.PointId+": ", golog.LstdFlags)
		}

		if err := mdev.rtuHandler.Connect(); err != nil {
			return err
		}
		mdev.Client = modbus.NewClient(mdev.rtuHandler)
	}
	if mdev.mainConfig.CommonConfig.Mode == "TCP" {
		mdev.tcpHandler = modbus.NewTCPClientHandler(
			fmt.Sprintf("%s:%v", mdev.mainConfig.HostConfig.Host, mdev.mainConfig.HostConfig.Port),
		)
		if core.GlobalConfig.DebugMode {
			mdev.tcpHandler.Logger = golog.New(glogger.GLogger.Writer(),
				"Modbus TCP Mode: "+mdev.PointId+": ", golog.LstdFlags)
		}
		if err := mdev.tcpHandler.Connect(); err != nil {
			return err
		}
		mdev.Client = modbus.NewClient(mdev.tcpHandler)
	}
	//---------------------------------------------------------------------------------
	// Start
	//---------------------------------------------------------------------------------
	if *mdev.mainConfig.CommonConfig.AutoRequest {
		mdev.retryTimes = 0
		go func(ctx context.Context) {
			for {
				select {
				case <-time.After(4 * time.Millisecond):
					// Continue loop
				case <-ctx.Done():
					return
				}
				ReadRegisterValues := []ReadRegisterValue{}
				if mdev.mainConfig.CommonConfig.Mode == "UART" {
					ReadRegisterValues = mdev.RTURead()
				}
				if mdev.mainConfig.CommonConfig.Mode == "TCP" {
					ReadRegisterValues = mdev.TCPRead()
				}
				if *mdev.mainConfig.CommonConfig.BatchRequest {
					if len(ReadRegisterValues) > 0 {
						if bytes, errMarshal := json.Marshal(ReadRegisterValues); errMarshal != nil {
							mdev.retryTimes++
							glogger.GLogger.Error(errMarshal)
						} else {
							mdev.RuleEngine.WorkDevice(mdev.Details(), string(bytes))
						}
					}
				}
				// 是否预警
				if *mdev.mainConfig.AlarmConfig.Enable {
					Input := map[string]any{}
					Input["data"] = ReadRegisterValues
					_, err := alarmcenter.Input(mdev.mainConfig.AlarmConfig.AlarmRuleId, Input)
					if err != nil {
						glogger.GLogger.Error(err)
					}
				}
			}

		}(mdev.Ctx)
	}
	if *mdev.mainConfig.CecollaConfig.Enable {
		// 新建设备端物模型
		if *mdev.mainConfig.CecollaConfig.EnableCreateSchema {
			Properties := []ithings.IthingsCreateSchemaPropertie{}
			for _, Register := range mdev.Registers {
				ProductId, DeviceId, err := ithings.ParseProductInfo(Register.Alias)
				if err != nil {
					glogger.Error(err)
				}
				Properties = append(Properties, ithings.IthingsCreateSchemaPropertie{
					Id:        Register.Tag,
					Name:      Register.Alias,
					Type:      ithings.ModbusTypeToSchemaType(Register.DataType),
					ProductId: ProductId,
					DeviceId:  DeviceId,
				})
			}
			cecolla := mdev.RuleEngine.GetCecolla(mdev.mainConfig.CecollaConfig.CecollaId)
			if bytes, errMarshal := json.Marshal(Properties); errMarshal != nil {
				glogger.Error(errMarshal)
			} else {
				if cecolla != nil {
					// CreateSubDeviceSchema
					_, errOnCtrl := cecolla.Cecolla.OnCtrl([]byte("CreateSubDeviceSchema"), bytes)
					if errOnCtrl != nil {
						glogger.Error(errOnCtrl)
					}
				}
			}
		}
	}

	mdev.status = typex.DEV_UP
	return nil
}
func (mdev *GenericModbusMaster) RTURead() []ReadRegisterValue {
	return mdev.modbusRead()
}
func (mdev *GenericModbusMaster) TCPRead() []ReadRegisterValue {
	return mdev.modbusRead()
}

// 设备当前状态
func (mdev *GenericModbusMaster) Status() typex.DeviceState {
	// 容错5次
	if mdev.retryTimes > 5 {
		return typex.DEV_DOWN
	}
	return mdev.status
}

// 停止设备
func (mdev *GenericModbusMaster) Stop() {
	mdev.status = typex.DEV_DOWN
	if mdev.CancelCTX != nil {
		mdev.CancelCTX()
	}
	if mdev.mainConfig.CommonConfig.Mode == "UART" {
		if mdev.rtuHandler != nil {
			mdev.rtuHandler.Close()
		}
	}
	if mdev.mainConfig.CommonConfig.Mode == "TCP" {
		if mdev.tcpHandler != nil {
			mdev.tcpHandler.Close()
		}
	}
	intercache.UnRegisterSlot(mdev.PointId) // 卸载点位表
	intercache.DeleteValue("__CecollaBinding", mdev.mainConfig.CecollaConfig.CecollaId)

}

// 真实设备
func (mdev *GenericModbusMaster) Details() *typex.Device {
	return mdev.RuleEngine.GetDevice(mdev.PointId)
}

// 状态
func (mdev *GenericModbusMaster) SetState(status typex.DeviceState) {
	mdev.status = status
}

func (mdev *GenericModbusMaster) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}

/**
 * 写入Modbus寄存器
 *
 */
// POST -> temp , 0x0001
type CtrlCmd struct {
	PointId string `json:"pointId,omitempty"` // 点位Point Id
	Tag     string `json:"tag,omitempty"`     // 点位Point Id
	Value   string `json:"value"`             // 写的值
}

func (O CtrlCmd) String() string {
	bytes, _ := json.Marshal(O)
	return string(bytes)
}

/**
 * 外部控制指令
 *
 */
func (mdev *GenericModbusMaster) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	if mdev.Client == nil {
		return nil, fmt.Errorf("invalid serial port handle")
	}
	glogger.Debug("GenericModbusMaster.OnCtrl, CMD=", string(cmd), ", Args=", string(args))
	// 根据Tag来写点位
	if string(cmd) == "WriteToSheetRegisterWithTag" {
		ctrlCmd := CtrlCmd{}
		if errUnmarshal := json.Unmarshal(args, &ctrlCmd); errUnmarshal != nil {
			return nil, errUnmarshal
		}
		Register, ok := mdev.RegisterWithTags[ctrlCmd.Tag]
		if ok {
			// 单个线圈
			// 0xFF00：表示线圈 ON（开）。任何非零值通常都表示线圈为 ON，但标准中通常用 0xFF00 来表示。
			// 0x0000：表示线圈 OFF（关）。它是唯一的有效值来表示线圈处于关闭状态。
			if Register.Function == 1 {
				if ctrlCmd.Value == "0" || ctrlCmd.Value == "false" {
					_, errW := mdev.Client.WriteSingleCoil(Register.Address, 0x0000)
					if errW != nil {
						return nil, errW
					}
				}
				if ctrlCmd.Value == "1" || ctrlCmd.Value == "true" {
					_, errW := mdev.Client.WriteSingleCoil(Register.Address, 0xFF00)
					if errW != nil {
						return nil, errW
					}
				}
			}
			// 单个寄存器
			if Register.Function == 3 {
				value, err := StringToUint16(ctrlCmd.Value)
				if err != nil {
					return nil, err
				}
				_, errW := mdev.Client.WriteSingleRegister(Register.Address, value)
				if errW != nil {
					return nil, errW
				}
			}
		}
	}
	// 写指令
	if string(cmd) == "WriteToSheetRegister" {
		ctrlCmd := CtrlCmd{}
		if errUnmarshal := json.Unmarshal(args, &ctrlCmd); errUnmarshal != nil {
			return nil, errUnmarshal
		}
		Register, ok := mdev.Registers[ctrlCmd.PointId]
		if ok {
			// 单个线圈
			// 0xFF00：表示线圈 ON（开）。任何非零值通常都表示线圈为 ON，但标准中通常用 0xFF00 来表示。
			// 0x0000：表示线圈 OFF（关）。它是唯一的有效值来表示线圈处于关闭状态。
			if Register.Function == 1 {
				if ctrlCmd.Value == "0" || ctrlCmd.Value == "false" {
					_, errW := mdev.Client.WriteSingleCoil(Register.Address, 0x0000)
					if errW != nil {
						return nil, errW
					}
				}
				if ctrlCmd.Value == "1" || ctrlCmd.Value == "true" {
					_, errW := mdev.Client.WriteSingleCoil(Register.Address, 0xFF00)
					if errW != nil {
						return nil, errW
					}
				}
			}
			// 单个寄存器
			if Register.Function == 3 {
				value, err := StringToUint16(ctrlCmd.Value)
				if err != nil {
					return nil, err
				}
				_, errW := mdev.Client.WriteSingleRegister(Register.Address, value)
				if errW != nil {
					return nil, errW
				}
			}
		}

	}
	return []byte{}, nil
}

// StringToUint16 将字符串转换为 uint16
func StringToUint16(s string) (uint16, error) {
	value, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(value), nil
}

/*
*
* 返回给Lua的数据结构,经过精简后的寄存器
*
 */
type ReadRegisterValue struct {
	Tag           string `json:"tag"`
	Alias         string `json:"alias"`
	SlaverId      byte   `json:"slaverId"`
	LastFetchTime uint64 `json:"lastFetchTime"`
	Value         string `json:"value"`
}

/*
*
* 串口模式
*
 */
func (mdev *GenericModbusMaster) modbusRead() []ReadRegisterValue {
	if *mdev.mainConfig.CommonConfig.EnableOptimize {
		return mdev.modbusSingleRead()
	} else {
		return mdev.modbusSingleRead()
	}
}

func (mdev *GenericModbusMaster) modbusSingleRead() []ReadRegisterValue {
	var err error
	var results []byte
	ModbusRegisters := []ReadRegisterValue{}
	count := len(mdev.Registers)
	if mdev.Client == nil {
		return ModbusRegisters
	}
	// modbusRead: 当读多字节寄存器的时候，需要考虑UTF8
	// Modbus收到的数据全部放进这个全局缓冲区内
	var __modbusReadResult = [256]byte{0} // 放在栈上提高效率
	for uuid, r := range mdev.Registers {
		if mdev.mainConfig.CommonConfig.Mode == "UART" {
			mdev.rtuHandler.SlaveId = r.SlaverId
		}
		if mdev.mainConfig.CommonConfig.Mode == "TCP" {
			mdev.tcpHandler.SlaveId = 0x01
		}
		// 1 字节
		if r.Function == dmodbus.READ_COIL {
			results, err = mdev.Client.ReadCoils(r.Address, r.Quantity)
			lastTimes := uint64(time.Now().UnixMilli())
			if err != nil {
				count--
				glogger.GLogger.Error(err)
				mdev.retryTimes++
				intercache.SetValue(mdev.PointId, uuid, intercache.CacheValue{
					UUID:          uuid,
					Status:        0,
					Value:         "0",
					LastFetchTime: lastTimes,
					ErrMsg:        err.Error(),
				})
				continue
			}
			// ValidData := [4]byte{0, 0, 0, 0}
			copy(__modbusReadResult[:], results[:])

			AnyValue := utils.ParseRegisterValue(len(results), r.DataType, r.DataOrder, float32(r.Weight), __modbusReadResult)
			Value := utils.CovertAnyType(AnyValue)
			Reg := ReadRegisterValue{
				Tag:           r.Tag,
				SlaverId:      r.SlaverId,
				Alias:         r.Alias,
				Value:         Value,
				LastFetchTime: lastTimes,
			}
			ModbusRegisters = append(ModbusRegisters, Reg)
			intercache.SetValue(mdev.PointId, uuid, intercache.CacheValue{
				UUID:          uuid,
				Status:        1,
				Value:         Value,
				LastFetchTime: lastTimes,
				ErrMsg:        "",
			})
			// 数据推向云边协同
			if *mdev.mainConfig.CecollaConfig.Enable {
				cecolla := mdev.RuleEngine.GetCecolla(mdev.mainConfig.CecollaConfig.CecollaId)
				if cecolla != nil {
					ProductId, DeviceName, err := ithings.ParseProductInfo(r.Alias)
					if err != nil {
						glogger.Error(err)
					} else {
						param := ithings.SubDeviceParam{
							Timestamp:  int64(lastTimes),
							Param:      r.Tag,
							ProductId:  ProductId,
							DeviceName: DeviceName,
							Value:      AnyValue,
						}
						_, OnCtrl := cecolla.Cecolla.OnCtrl([]byte("PackReportSubDeviceParams"), []byte(param.String()))
						if OnCtrl != nil {
							glogger.Error(OnCtrl)
						}
					}
				}
			}
			if !*mdev.mainConfig.CommonConfig.BatchRequest {
				if bytes, errMarshal := json.Marshal(Reg); errMarshal != nil {
					glogger.GLogger.Error(errMarshal)
				} else {
					mdev.RuleEngine.WorkDevice(mdev.Details(), string(bytes))
				}
			}
		}
		// 2 字节
		if r.Function == dmodbus.READ_DISCRETE_INPUT {
			results, err = mdev.Client.ReadDiscreteInputs(r.Address, r.Quantity)
			lastTimes := uint64(time.Now().UnixMilli())
			if err != nil {
				count--
				glogger.GLogger.Error(err)
				mdev.retryTimes++
				intercache.SetValue(mdev.PointId, uuid, intercache.CacheValue{
					UUID:          uuid,
					Status:        0,
					Value:         "0",
					LastFetchTime: lastTimes,
					ErrMsg:        err.Error(),
				})
				continue
			}
			// ValidData := [4]byte{0, 0, 0, 0}
			copy(__modbusReadResult[:], results[:])
			AnyValue := utils.ParseRegisterValue(len(results), r.DataType, r.DataOrder, float32(r.Weight), __modbusReadResult)
			Value := utils.CovertAnyType(AnyValue)
			Reg := ReadRegisterValue{
				Tag:           r.Tag,
				SlaverId:      r.SlaverId,
				Alias:         r.Alias,
				Value:         Value,
				LastFetchTime: lastTimes,
			}
			ModbusRegisters = append(ModbusRegisters, Reg)
			intercache.SetValue(mdev.PointId, uuid, intercache.CacheValue{
				UUID:          uuid,
				Status:        1,
				Value:         Value,
				LastFetchTime: lastTimes,
				ErrMsg:        "",
			})
			// 数据推向云边协同
			if *mdev.mainConfig.CecollaConfig.Enable {
				cecolla := mdev.RuleEngine.GetCecolla(mdev.mainConfig.CecollaConfig.CecollaId)
				if cecolla != nil {
					ProductId, DeviceName, err := ithings.ParseProductInfo(r.Alias)
					if err != nil {
						glogger.Error(err)
					} else {
						param := ithings.SubDeviceParam{
							Timestamp:  int64(lastTimes),
							Param:      r.Tag,
							ProductId:  ProductId,
							DeviceName: DeviceName,
							Value:      AnyValue,
						}
						_, OnCtrl := cecolla.Cecolla.OnCtrl([]byte("PackReportSubDeviceParams"), []byte(param.String()))
						if OnCtrl != nil {
							glogger.Error(OnCtrl)
						}
					}
				}
			}
			if !*mdev.mainConfig.CommonConfig.BatchRequest {
				if bytes, errMarshal := json.Marshal(Reg); errMarshal != nil {
					glogger.GLogger.Error(errMarshal)
				} else {
					mdev.RuleEngine.WorkDevice(mdev.Details(), string(bytes))
				}
			}
		}
		// 2 字节
		//
		if r.Function == dmodbus.READ_HOLDING_REGISTERS {
			results, err = mdev.Client.ReadHoldingRegisters(r.Address, r.Quantity)
			lastTimes := uint64(time.Now().UnixMilli())
			if err != nil {
				count--
				glogger.GLogger.Error(err)
				mdev.retryTimes++
				intercache.SetValue(mdev.PointId, uuid, intercache.CacheValue{
					UUID:          uuid,
					Status:        0,
					Value:         "0",
					LastFetchTime: lastTimes,
					ErrMsg:        err.Error(),
				})
				continue
			}
			// ValidData := [4]byte{0, 0, 0, 0}
			copy(__modbusReadResult[:], results[:])
			AnyValue := utils.ParseRegisterValue(len(results), r.DataType, r.DataOrder, float32(r.Weight), __modbusReadResult)
			Value := utils.CovertAnyType(AnyValue)
			Reg := ReadRegisterValue{
				Tag:           r.Tag,
				SlaverId:      r.SlaverId,
				Alias:         r.Alias,
				Value:         Value,
				LastFetchTime: lastTimes,
			}
			ModbusRegisters = append(ModbusRegisters, Reg)
			intercache.SetValue(mdev.PointId, uuid, intercache.CacheValue{
				UUID:          uuid,
				Status:        1,
				Value:         Value,
				LastFetchTime: lastTimes,
				ErrMsg:        "",
			})
			// 数据推向云边协同
			if *mdev.mainConfig.CecollaConfig.Enable {
				cecolla := mdev.RuleEngine.GetCecolla(mdev.mainConfig.CecollaConfig.CecollaId)
				if cecolla != nil {
					ProductId, DeviceName, err := ithings.ParseProductInfo(r.Alias)
					if err != nil {
						glogger.Error(err)
					} else {
						param := ithings.SubDeviceParam{
							Timestamp:  int64(lastTimes),
							Param:      r.Tag,
							ProductId:  ProductId,
							DeviceName: DeviceName,
							Value:      AnyValue,
						}
						_, errOnCtrl := cecolla.Cecolla.OnCtrl([]byte("PackReportSubDeviceParams"), []byte(param.String()))
						if errOnCtrl != nil {
							glogger.Error(errOnCtrl)
						}
					}
				}
			}
			if !*mdev.mainConfig.CommonConfig.BatchRequest {
				if bytes, errMarshal := json.Marshal(Reg); errMarshal != nil {
					glogger.GLogger.Error(errMarshal)
				} else {
					mdev.RuleEngine.WorkDevice(mdev.Details(), string(bytes))
				}
			}
		}
		// 2 字节
		if r.Function == dmodbus.READ_INPUT_REGISTERS {
			results, err = mdev.Client.ReadInputRegisters(r.Address, r.Quantity)
			lastTimes := uint64(time.Now().UnixMilli())
			if err != nil {
				count--
				glogger.GLogger.Error(err)
				mdev.retryTimes++
				intercache.SetValue(mdev.PointId, uuid, intercache.CacheValue{
					UUID:          uuid,
					Status:        0,
					Value:         "0",
					LastFetchTime: lastTimes,
					ErrMsg:        err.Error(),
				})
				continue
			}
			// ValidData := [4]byte{0, 0, 0, 0}
			copy(__modbusReadResult[:], results[:])
			AnyValue := utils.ParseRegisterValue(len(results), r.DataType, r.DataOrder, float32(r.Weight), __modbusReadResult)
			Value := utils.CovertAnyType(AnyValue)
			Reg := ReadRegisterValue{
				Tag:           r.Tag,
				SlaverId:      r.SlaverId,
				Alias:         r.Alias,
				Value:         Value,
				LastFetchTime: lastTimes,
			}
			ModbusRegisters = append(ModbusRegisters, Reg)
			intercache.SetValue(mdev.PointId, uuid, intercache.CacheValue{
				UUID:          uuid,
				Status:        1,
				Value:         Value,
				LastFetchTime: lastTimes,
				ErrMsg:        "",
			})
			// 数据推向云边协同
			if *mdev.mainConfig.CecollaConfig.Enable {
				cecolla := mdev.RuleEngine.GetCecolla(mdev.mainConfig.CecollaConfig.CecollaId)
				if cecolla != nil {
					ProductId, DeviceName, err := ithings.ParseProductInfo(r.Alias)
					if err != nil {
						glogger.Error(err)
					} else {
						param := ithings.SubDeviceParam{
							Timestamp:  int64(lastTimes),
							Param:      r.Tag,
							ProductId:  ProductId,
							DeviceName: DeviceName,
							Value:      AnyValue,
						}
						_, OnCtrl := cecolla.Cecolla.OnCtrl([]byte("PackReportSubDeviceParams"), []byte(param.String()))
						if OnCtrl != nil {
							glogger.Error(OnCtrl)
						}
					}
				}
			}
			if !*mdev.mainConfig.CommonConfig.BatchRequest {
				if bytes, errMarshal := json.Marshal(Reg); errMarshal != nil {
					glogger.GLogger.Error(errMarshal)
				} else {
					mdev.RuleEngine.WorkDevice(mdev.Details(), string(bytes))
				}
			}
		}
		time.Sleep(time.Duration(r.Frequency) * time.Millisecond)
	}
	return ModbusRegisters
}
