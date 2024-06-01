package device

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/intercache"
	"github.com/hootrhino/rhilex/component/interdb"

	"github.com/hootrhino/gobacnet"
	"github.com/hootrhino/gobacnet/apdus"
	"github.com/hootrhino/gobacnet/btypes"

	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

type bacnetCommonConfig struct {
	Frequency int `json:"frequency" title:"采集间隔，单位毫秒"`
}
type bacnetConfig struct {
	Mode string `json:"mode" title:"bacnet运行模式"`

	Ip     string `json:"ip" title:"bacnet设备ip(仅type=SINGLE生效)"`
	Port   int    `json:"port" title:"bacnet端口，通常是47808(仅type=SINGLE生效)"`
	IsMstp int    `json:"isMstp" title:"是否为mstp设备，若是则子网号必须填写(仅type=SINGLE时生效)"`
	Subnet int    `json:"subnet" title:"子网号(仅type=SINGLE 且 isMstp=1 时生效)"`

	LocalIp    string `json:"LocalIp" title:"本地ip地址(仅type=BROADCAST时有效)"`
	SubnetCIDR int    `json:"subnetCidr" title:"子网掩码长度(仅type=BROADCAST时有效)"`
	LocalPort  int    `json:"localPort" title:"本地监听端口，填0表示默认47808(有的模拟器必须本地监听47808才能正常交互)"`
	//
	DeviceId  uint32 `json:"deviceId"`
	VendorId  uint16 `json:"vendorId"`
	NetWorkId uint16 `json:"netWorkId"`
}

type bacnetDataPoint struct {
	UUID           string            `json:"uuid"`
	Tag            string            `json:"tag" validate:"required" title:"数据Tag"`
	BacnetDeviceId uint32            `json:"bacnetDeviceId" title:"bacnet设备id(若isMstp=1，则deviceId应该必填；若是纯bacnetip设备，则填1即可)"`
	ObjectType     btypes.ObjectType `json:"objectType" title:"object类型"`
	ObjectId       uint32            `json:"objectId" title:"object的id"`

	property btypes.PropertyData
}
type BacnetMainConfig struct {
	BacnetConfig bacnetConfig       `json:"bacnetConfig" validate:"required"`
	CommonConfig bacnetCommonConfig `json:"commonConfig" validate:"required"`
}
type GenericBacnetIpDevice struct {
	typex.XStatus
	bacnetClient bacnet.Client
	status       typex.DeviceState
	RuleEngine   typex.Rhilex
	mainConfig   BacnetMainConfig
	// 点位表
	SubDeviceDataPoints []bacnetDataPoint           // 读到子设备的点位
	SelfPropertyData    map[uint32][2]btypes.Object // 自己的点位
	remoteDeviceMap     map[uint32]btypes.Device
}

func NewGenericBacnetIpDevice(e typex.Rhilex) typex.XDevice {
	g := new(GenericBacnetIpDevice)
	g.RuleEngine = e
	g.mainConfig = BacnetMainConfig{
		CommonConfig: bacnetCommonConfig{Frequency: 1000},
		BacnetConfig: bacnetConfig{
			Mode:       "BROADCAST",
			LocalIp:    "127.0.0.1",
			SubnetCIDR: 24,
			LocalPort:  47808,
			DeviceId:   2580,
			VendorId:   2580,
		},
	}
	g.SubDeviceDataPoints = make([]bacnetDataPoint, 0)
	g.status = typex.DEV_DOWN
	return g
}

func (dev *GenericBacnetIpDevice) Init(devId string, configMap map[string]interface{}) error {
	dev.PointId = devId
	// 先给个空的
	dev.remoteDeviceMap = make(map[uint32]btypes.Device)

	intercache.RegisterSlot(devId)
	err := utils.BindSourceConfig(configMap, &dev.mainConfig)
	if err != nil {
		return err
	}
	var dataPoints []model.MBacnetDataPoint
	err = interdb.DB().Table("m_bacnet_data_points").
		Where("device_uuid=?", devId).Find(&dataPoints).Error

	for _, mDataPoint := range dataPoints {
		dataPoint := bacnetDataPoint{
			UUID:           mDataPoint.UUID,
			Tag:            mDataPoint.Tag,
			BacnetDeviceId: mDataPoint.BacnetDeviceId,
			ObjectType:     getObjectTypeByNumber(mDataPoint.ObjectType),
			ObjectId:       mDataPoint.ObjectId,
		}
		// Cache Value
		intercache.SetValue(dev.PointId, mDataPoint.UUID, intercache.CacheValue{
			UUID:          mDataPoint.UUID,
			Status:        0,
			LastFetchTime: uint64(time.Now().UnixMilli()),
			Value:         "",
			ErrMsg:        "Loading",
		})
		dev.SubDeviceDataPoints = append(dev.SubDeviceDataPoints, dataPoint)
	}
	if err != nil {
		glogger.GLogger.Error(err)
		return err
	}

	return nil
}

func getObjectTypeByNumber(strType string) btypes.ObjectType {
	switch strType {
	case "AI":
		return btypes.AnalogInput
	case "AO":
		return btypes.AnalogOutput
	case "AV":
		return btypes.AnalogValue
	case "BI":
		return btypes.BinaryInput
	case "BO":
		return btypes.BinaryOutput
	case "BV":
		return btypes.BinaryValue
	case "MI":
		return btypes.MultiStateInput
	case "MO":
		return btypes.MultiStateOutput
	case "MV":
		return btypes.MultiStateValue
	}
	return btypes.AnalogInput
}

func (dev *GenericBacnetIpDevice) Start(cctx typex.CCTX) error {
	dev.CancelCTX = cctx.CancelCTX
	dev.Ctx = cctx.Ctx
	PropertyData := map[uint32][2]btypes.Object{}
	// 将nodeConfig对应的配置信息
	for idx, BacnetDataPoint := range dev.SubDeviceDataPoints {
		SubPropertyData := btypes.PropertyData{
			Object: btypes.Object{
				ID: btypes.ObjectID{
					Type:     BacnetDataPoint.ObjectType,
					Instance: btypes.ObjectInstance(BacnetDataPoint.ObjectId),
				},
				Properties: []btypes.Property{
					{
						Type:       btypes.PropPresentValue, // Present value
						ArrayIndex: btypes.ArrayAll,
					},
				},
			},
		}
		dev.SubDeviceDataPoints[idx].property = SubPropertyData
		// 配置自身的点位
		PropertyData[uint32(BacnetDataPoint.ObjectId)] = apdus.NewAIPropertyWithRequiredFields(BacnetDataPoint.Tag,
			uint32(BacnetDataPoint.ObjectId), float32(0.00), "")
	}
	// 广播模式监听
	if dev.mainConfig.BacnetConfig.Mode == "BROADCAST" {
		// 创建一个bacnet ip的本地网络
		// PropertyData:= map[uint32][2]btypes.Object{
		// 	1: apdus.NewAIPropertyWithRequiredFields("temp", 1, float32(3.14), "empty"),
		// 	2: apdus.NewAIPropertyWithRequiredFields("humi", 2, float32(77.67), "empty"),
		// 	3: apdus.NewAIPropertyWithRequiredFields("pres", 3, float32(101.11), "empty"),
		// },
		client, err := bacnet.NewClient(&bacnet.ClientBuilder{
			Ip:           dev.mainConfig.BacnetConfig.LocalIp,
			Port:         dev.mainConfig.BacnetConfig.LocalPort,
			SubnetCIDR:   dev.mainConfig.BacnetConfig.SubnetCIDR,
			DeviceId:     dev.mainConfig.BacnetConfig.DeviceId,         // 参数化
			VendorId:     uint32(dev.mainConfig.BacnetConfig.VendorId), // 参数化
			NetWorkId:    dev.mainConfig.BacnetConfig.NetWorkId,        // 参数化
			PropertyData: PropertyData,                                 // 点位表, 需要更新为动态
		})
		if err != nil {
			return err
		}

		dev.bacnetClient = client
		client.SetLogger(glogger.GLogger.Logger)
		go dev.bacnetClient.ClientRun()
		go func(ctx context.Context) {
			// 定时刷新device列表 后续可以优化下逻辑
			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					/// 迁移配置到前端控制
					devices, err := client.WhoIs(&bacnet.WhoIsOpts{
						Low:             -1,
						High:            -1,
						GlobalBroadcast: true,
					})
					if err != nil {
						glogger.GLogger.Error(err)
						continue
					}
					if len(devices) > 0 {
						deviceMap := make(map[uint32]btypes.Device)
						for i := range devices {
							deviceMap[uint32(devices[i].DeviceID)] = devices[i]
						}
						dev.remoteDeviceMap = nil
						dev.remoteDeviceMap = deviceMap
					}

				case <-ctx.Done():
					return
				}
			}
		}(dev.Ctx)

	}

	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Duration(dev.mainConfig.CommonConfig.Frequency) * time.Millisecond)
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			default:
			}

			read, err2 := dev.ReadProperty()
			if err2 != nil {
				glogger.GLogger.Error(err2)
			} else {
				dev.RuleEngine.WorkDevice(dev.Details(), string(read))
			}
			<-ticker.C
		}
	}(dev.Ctx)

	dev.status = typex.DEV_UP
	return nil
}

func (dev *GenericBacnetIpDevice) OnRead(cmd []byte, data []byte) (int, error) {
	read, err := dev.ReadProperty()
	if err != nil {
		return 0, err
	}
	len := copy(data, read)
	return len, nil
}

type ReturnValue struct {
	Tag              string      `json:"tag"`
	DeviceId         uint32      `json:"deviceId"`
	PropertyType     string      `json:"propertyType"`
	PropertyInstance uint32      `json:"propertyInstance"`
	Value            interface{} `json:"value"`
}

/*
*
* 局域网广播
*
 */
func (dev *GenericBacnetIpDevice) ReadProperty() ([]byte, error) {
	retMap := map[string]ReturnValue{}
	for _, SubDeviceDataPoint := range dev.SubDeviceDataPoints {
		var bacnetDeviceId uint32
		if dev.mainConfig.BacnetConfig.Mode == "SINGLE" {
			bacnetDeviceId = 1
		} else {
			bacnetDeviceId = SubDeviceDataPoint.BacnetDeviceId
		}
		if device, ok := dev.remoteDeviceMap[bacnetDeviceId]; ok {
			property, err := dev.bacnetClient.ReadProperty(device, SubDeviceDataPoint.property)
			if err != nil {
				glogger.GLogger.Errorf("bacnet Client Read Property failed. tag = %v, err=%v", SubDeviceDataPoint.Tag, err)
				intercache.SetValue(dev.PointId, SubDeviceDataPoint.UUID, intercache.CacheValue{
					UUID:          SubDeviceDataPoint.UUID,
					Status:        0,
					LastFetchTime: uint64(time.Now().UnixMilli()),
					Value:         "",
					ErrMsg:        err.Error(),
				})
				dev.bacnetClient.GetBacnetIPServer().
					UpdateAIPropertyValue(uint32(SubDeviceDataPoint.ObjectId), float32(0))
				continue
			}
			ReturnValue := ReturnValue{
				Tag:              SubDeviceDataPoint.Tag,
				DeviceId:         bacnetDeviceId,
				PropertyType:     property.Object.ID.Type.String(),
				PropertyInstance: uint32(property.Object.ID.Instance),
			}
			if len(property.Object.Properties) > 0 {
				ReturnValue.Value = property.Object.Properties[0].Data
			} else {
				ReturnValue.Value = uint32(0)
			}
			retMap[SubDeviceDataPoint.Tag] = ReturnValue
			dev.bacnetClient.GetBacnetIPServer().
				UpdateAIPropertyValue(uint32(SubDeviceDataPoint.ObjectId), ReturnValue.Value)

			intercache.SetValue(dev.PointId, SubDeviceDataPoint.UUID, intercache.CacheValue{
				UUID:          SubDeviceDataPoint.UUID,
				Status:        1,
				LastFetchTime: uint64(time.Now().UnixMilli()),
				Value:         fmt.Sprintf("%v", ReturnValue.Value),
				ErrMsg:        "",
			})
		}
	}
	bytes, _ := json.Marshal(retMap)
	glogger.GLogger.Debug(string(bytes))
	return bytes, nil
}

func (dev *GenericBacnetIpDevice) OnWrite(cmd []byte, data []byte) (int, error) {
	//TODO implement me
	return 0, errors.New("not Support")
}

func (dev *GenericBacnetIpDevice) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return nil, errors.New("not Support")
}

func (dev *GenericBacnetIpDevice) Status() typex.DeviceState {
	return dev.status
}

func (dev *GenericBacnetIpDevice) Stop() {
	dev.status = typex.DEV_DOWN
	if dev.CancelCTX != nil {
		dev.CancelCTX()
	}
	if dev.bacnetClient != nil {
		dev.bacnetClient.Close()
	}
	intercache.UnRegisterSlot(dev.PointId)
}

func (dev *GenericBacnetIpDevice) Details() *typex.Device {
	return dev.RuleEngine.GetDevice(dev.PointId)
}

func (dev *GenericBacnetIpDevice) SetState(state typex.DeviceState) {
	dev.status = state
}

func (dev *GenericBacnetIpDevice) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}
