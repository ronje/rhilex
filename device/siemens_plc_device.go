package device

import (
	"context"

	"encoding/json"
	"errors"

	"sync"
	"time"

	"github.com/hootrhino/rhilex/component/intercache"

	"github.com/jinzhu/copier"

	"github.com/hootrhino/rhilex/common"
	"github.com/hootrhino/rhilex/component/interdb"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
	"github.com/robinson/gos7"
)

// 点位表
type __SiemensDataPoint struct {
	UUID            string  `json:"uuid"`
	DeviceUUID      string  `json:"device_uuid"`
	SiemensAddress  string  `json:"siemensAddress"` // 西门子的地址字符串
	Tag             string  `json:"tag"`
	Alias           string  `json:"alias"`
	Frequency       int64   `json:"frequency"`
	Status          int     `json:"status"`        // 运行时数据
	LastFetchTime   uint64  `json:"lastFetchTime"` // 运行时数据
	Value           string  `json:"value"`         // 运行时数据
	AddressType     string  `json:"-"`             // // 西门子解析后的地址信息: 寄存器类型: DB I Q
	DataBlockType   string  `json:"-"`             // // 西门子解析后的地址信息: 数据类型: INT UINT ....
	DataBlockOrder  string  `json:"-"`             //  西门子解析后的地址信息: 数据类型: INT UINT ....
	Weight          float64 `json:"-"`             // 权重
	DataBlockNumber int     `json:"-"`             // // 西门子解析后的地址信息: 数据块号: 100...
	ElementNumber   int     `json:"-"`             // // 西门子解析后的地址信息: 元素号:1000...
	DataSize        int     `json:"-"`             // // 西门子解析后的地址信息: 位号,0-8，只针对I、Q
	BitNumber       int     `json:"-"`             // // 西门子解析后的地址信息: 位号,0-8，只针对I、Q
}

// https://cloudvpn.beijerelectronics.com/hc/en-us/articles/4406049761169-Siemens-S7
type S1200CommonConfig struct {
	Timeout      int  `json:"timeout" validate:"required"`      // 5s
	IdleTimeout  int  `json:"idleTimeout" validate:"required"`  // 5s
	AutoRequest  bool `json:"autoRequest" validate:"required"`  // false
	BatchRequest bool `json:"batchRequest" validate:"required"` // 批量采集
}

type S1200Config struct {
	Host  string `json:"host" validate:"required"`  // 127.0.0.1:502
	Model string `json:"model" validate:"required"` // s7-200 s7-1500
	Rack  int    `json:"rack" validate:"required"`  // 0
	Slot  int    `json:"slot" validate:"required"`  // 1
}
type S1200MainConfig struct {
	CommonConfig  S1200CommonConfig    `json:"commonConfig" validate:"required"` // 通用配置
	S1200Config   S1200Config          `json:"s1200Config" validate:"required"`  // 通用配置
	CecollaConfig common.CecollaConfig `json:"cecollaConfig"`
}

// https://www.ad.siemens.com.cn/productportal/prods/s7-1200_plc_easy_plus/07-Program/02-basic/01-Data_Type/01-basic.html
type SIEMENS_PLC struct {
	typex.XStatus
	status              typex.DeviceState
	RuleEngine          typex.Rhilex
	mainConfig          S1200MainConfig
	client              gos7.Client
	handler             *gos7.TCPClientHandler
	locker              sync.Mutex
	__SiemensDataPoints map[string]*__SiemensDataPoint
}

/*
*
* 西门子 S1200 系列 PLC
*
 */
func NewSIEMENS_PLC(e typex.Rhilex) typex.XDevice {
	s1200 := new(SIEMENS_PLC)
	s1200.RuleEngine = e
	s1200.locker = sync.Mutex{}
	s1200.mainConfig = S1200MainConfig{
		CommonConfig: S1200CommonConfig{
			Timeout:      1000,
			IdleTimeout:  3000,
			AutoRequest:  false,
			BatchRequest: false,
		},
		S1200Config: S1200Config{
			Host: "127.0.0.1:1500",
			Rack: 0,
			Slot: 1,
		},
	}
	s1200.__SiemensDataPoints = map[string]*__SiemensDataPoint{}
	return s1200
}

// 初始化
func (s1200 *SIEMENS_PLC) Init(devId string, configMap map[string]interface{}) error {
	s1200.PointId = devId
	intercache.RegisterSlot(s1200.PointId)
	if err := utils.BindSourceConfig(configMap, &s1200.mainConfig); err != nil {
		glogger.GLogger.Error(err)
		return err
	}
	// 合并数据库里面的点位表
	// TODO 这里需要优化一下，而不是直接查表这种形式，应该从物模型组件来加载
	// DataSchema = schema.load(uuid)
	// DataSchema.update(k, v)
	var list []__SiemensDataPoint
	errDb := interdb.DB().Table("m_siemens_data_points").
		Where("device_uuid=?", devId).Find(&list).Error
	if errDb != nil {
		return errDb
	}
	// 开始解析地址表
	for _, SiemensDataPoint := range list {
		// 频率不能太快
		if SiemensDataPoint.Frequency < 1 {
			return errors.New("'frequency' must grate than 50 millisecond")
		}
		//
		AddressInfo, err1 := utils.ParseSiemensDB(SiemensDataPoint.SiemensAddress)
		if err1 != nil {
			return err1
		}
		SiemensDataPoint.DataBlockNumber = AddressInfo.DataBlockNumber
		SiemensDataPoint.ElementNumber = AddressInfo.ElementNumber
		SiemensDataPoint.AddressType = AddressInfo.AddressType
		SiemensDataPoint.BitNumber = AddressInfo.BitNumber
		SiemensDataPoint.DataSize = AddressInfo.DataBlockSize

		// 提前缓冲
		NewSiemensDataPoint := __SiemensDataPoint{}
		copier.Copy(&NewSiemensDataPoint, &SiemensDataPoint)
		s1200.__SiemensDataPoints[SiemensDataPoint.UUID] = &NewSiemensDataPoint
		intercache.SetValue(s1200.PointId, SiemensDataPoint.UUID, intercache.CacheValue{
			UUID:          SiemensDataPoint.UUID,
			Status:        0,
			LastFetchTime: 0,
			Value:         "0",
			ErrMsg:        "",
		})
	}

	return nil
}

// 启动
func (s1200 *SIEMENS_PLC) Start(cctx typex.CCTX) error {
	s1200.Ctx = cctx.Ctx
	s1200.CancelCTX = cctx.CancelCTX
	//
	s1200.handler = gos7.NewTCPClientHandler(
		s1200.mainConfig.S1200Config.Host, // 127.0.0.1:1500
		s1200.mainConfig.S1200Config.Rack, // 0
		s1200.mainConfig.S1200Config.Slot) // 1
	s1200.handler.Timeout = time.Duration(
		s1200.mainConfig.CommonConfig.Timeout) * time.Millisecond
	s1200.handler.IdleTimeout = time.Duration(
		s1200.mainConfig.CommonConfig.IdleTimeout) * time.Millisecond
	if err := s1200.handler.Connect(); err != nil {
		return err
	} else {
		s1200.status = typex.DEV_UP
	}

	s1200.client = gos7.NewClient(s1200.handler)
	if !s1200.mainConfig.CommonConfig.AutoRequest {
		s1200.status = typex.DEV_UP
		return nil
	}
	go func(ctx context.Context) {
		for {
			select {
			case <-s1200.Ctx.Done():
				return
			case <-time.After(4 * time.Millisecond):
				// Continue loop
			}
			s1200.locker.Lock()
			ReadPLCRegisterValues := s1200.Read()
			s1200.locker.Unlock()
			if len(ReadPLCRegisterValues) < 1 {
				time.Sleep(50 * time.Second)
				continue
			}
			if !s1200.mainConfig.CommonConfig.BatchRequest {
				if len(ReadPLCRegisterValues) > 0 {
					if bytes, err := json.Marshal(ReadPLCRegisterValues); err != nil {
						glogger.GLogger.Error(err)
					} else {
						s1200.RuleEngine.WorkDevice(s1200.Details(), string(bytes))
					}
				}
			}
		}
	}(cctx.Ctx)
	return nil
}


// 设备当前状态
func (s1200 *SIEMENS_PLC) Status() typex.DeviceState {
	if s1200.client == nil {
		return typex.DEV_DOWN
	}
	return s1200.status

}

// 停止设备
func (s1200 *SIEMENS_PLC) Stop() {
	s1200.status = typex.DEV_DOWN
	if s1200.CancelCTX != nil {
		s1200.CancelCTX()
	}
	if s1200.handler != nil {
		s1200.handler.Close()
	}
	intercache.UnRegisterSlot(s1200.PointId)
}

// 真实设备
func (s1200 *SIEMENS_PLC) Details() *typex.Device {
	return s1200.RuleEngine.GetDevice(s1200.PointId)
}

// 状态
func (s1200 *SIEMENS_PLC) SetState(status typex.DeviceState) {
	s1200.status = status
}

func (s1200 *SIEMENS_PLC) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}
func (s1200 *SIEMENS_PLC) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return []byte{}, nil
}
func (s1200 *SIEMENS_PLC) Write(cmd []byte, data []byte) (int, error) {
	return 0, nil
}

// 字节格式:[dbNumber1, start1, size1, dbNumber2, start2, size2]
// 读: db --> dbNumber, start, size, buffer[]
var rData = [common.T_2KB]byte{} // 一次最大接受2KB数据
// SIEMENS_PLC: 当读多字节寄存器的时候，需要考虑UTF8
var __siemensReadResult = [256]byte{0}

type ReadPLCRegisterValue struct {
	Tag           string `json:"tag"`
	Alias         string `json:"alias"`
	LastFetchTime uint64 `json:"lastFetchTime"`
	Value         string `json:"value"`
}

func (s1200 *SIEMENS_PLC) Read() []ReadPLCRegisterValue {
	values := []ReadPLCRegisterValue{}
	for uuid, db := range s1200.__SiemensDataPoints {
		//DB 4字节
		if db.AddressType == "DB" {
			lastTimes := uint64(time.Now().UnixMilli())

			// 00.00.00.01 | 00.00.00.02 | 00.00.00.03 | 00.00.00.04
			// 根据类型解析长度
			if err := s1200.client.AGReadDB(db.DataBlockNumber,
				db.ElementNumber, db.DataSize, rData[:]); err != nil {
				glogger.GLogger.Error(err)
				intercache.SetValue(s1200.PointId, uuid, intercache.CacheValue{
					UUID:          uuid,
					Status:        1,
					LastFetchTime: lastTimes,
					Value:         "",
					ErrMsg:        err.Error(),
				})
				continue
			}
			// ValidData := [4]byte{} // 固定4字节，以后有8自己的时候再支持
			copy(__siemensReadResult[:], rData[:db.DataSize])
			Value := utils.ParseModbusValue(db.DataSize, db.DataBlockType, db.DataBlockOrder,
				float32(db.Weight), __siemensReadResult)
			PlcReadReg := ReadPLCRegisterValue{
				Tag:           db.Tag,
				Alias:         db.Alias,
				Value:         Value,
				LastFetchTime: lastTimes,
			}
			values = append(values, PlcReadReg)
			intercache.SetValue(s1200.PointId, uuid, intercache.CacheValue{
				UUID:          uuid,
				Status:        0,
				LastFetchTime: lastTimes,
				Value:         Value,
				ErrMsg:        "",
			})
			if !s1200.mainConfig.CommonConfig.BatchRequest {
				if bytes, errMarshal := json.Marshal(PlcReadReg); errMarshal != nil {
					glogger.GLogger.Error(errMarshal)
				} else {
					s1200.RuleEngine.WorkDevice(s1200.Details(), string(bytes))
				}
			}
		}
		if db.Frequency < 10 {
			db.Frequency = 100 // 不能太快
		}
		time.Sleep(time.Duration(db.Frequency) * time.Millisecond)
	}
	return values
}
