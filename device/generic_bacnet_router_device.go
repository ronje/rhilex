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
	"net"
	"time"

	bacnet "github.com/hootrhino/gobacnet"
	"github.com/hootrhino/gobacnet/apdus"
	"github.com/hootrhino/gobacnet/btypes"
	"github.com/hootrhino/rhilex/component/apiserver/model"
	"github.com/hootrhino/rhilex/component/intercache"
	"github.com/hootrhino/rhilex/component/interdb"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

type BacnetRouterConfig struct {
	Mode        string `json:"mode" validate:"required"` // IP/MSTP
	LocalPort   int    `json:"localPort" validate:"required"`
	NetworkCidr string `json:"networkCidr" validate:"required"`
	DeviceId    uint32 `json:"deviceId" validate:"required"`
	VendorId    uint32 `json:"vendorId" validate:"required"`
	NetWorkId   uint16 `json:"netWorkId" validate:"required"`
}

type BacnetRouterMainConfig struct {
	BacnetRouterConfig BacnetRouterConfig `json:"bacnetRouterConfig" validate:"required"`
}

type BacnetRouter struct {
	typex.XStatus
	status           typex.DeviceState
	mainConfig       BacnetRouterMainConfig
	bacnetClient     bacnet.Client
	selfPropertyData map[uint32][2]btypes.Object
}

type BacnetRouterDataPoint struct {
	UUID       string            `json:"uuid"`
	Tag        string            `json:"tag" validate:"required" title:"数据Tag"`
	ObjectId   uint32            `json:"id" title:"object的id"`
	ObjectType btypes.ObjectType `json:"type" title:"object类型"`
}

func NewBacnetRouter(e typex.Rhilex) typex.XDevice {
	br := new(BacnetRouter)
	br.RuleEngine = e
	br.mainConfig = BacnetRouterMainConfig{
		BacnetRouterConfig: BacnetRouterConfig{
			Mode:      "BROADCAST",
			LocalPort: 47808,
			DeviceId:  2580,
			VendorId:  2580,
			NetWorkId: 2580,
		},
	}
	br.status = typex.DEV_DOWN
	return br
}

func (br *BacnetRouter) Init(devId string, configMap map[string]interface{}) error {
	br.PointId = devId

	intercache.RegisterSlot(devId)
	err := utils.BindSourceConfig(configMap, &br.mainConfig)
	if err != nil {
		return err
	}
	var MBacnetRouterDataPoints []model.MBacnetRouterDataPoint
	err = interdb.DB().Table("m_bacnet_router_data_points").
		Where("device_uuid=?", devId).Find(&MBacnetRouterDataPoints).Error
	if err != nil {
		glogger.GLogger.Error(err)
		return err
	}

	br.selfPropertyData = make(map[uint32][2]btypes.Object)
	for _, mDataPoint := range MBacnetRouterDataPoints {
		// Map Model to Point
		// Cache Value
		intercache.SetValue(br.PointId, mDataPoint.UUID, intercache.CacheValue{
			UUID:          mDataPoint.UUID,
			Status:        0,
			LastFetchTime: uint64(time.Now().UnixMilli()),
			Value:         "",
			ErrMsg:        "-/-",
		})
		br.selfPropertyData[mDataPoint.ObjectId] = apdus.NewAIPropertyWithRequiredFields(mDataPoint.Tag,
			mDataPoint.ObjectId, float32(0), "-/-")
	}

	return nil
}

func (br *BacnetRouter) Start(cctx typex.CCTX) error {
	br.Ctx = cctx.Ctx
	br.CancelCTX = cctx.CancelCTX
	// 创建一个bacnet ip的本地网络
	IP, IPNet, errParseCIDR := net.ParseCIDR(br.mainConfig.BacnetRouterConfig.NetworkCidr)
	if errParseCIDR != nil {
		glogger.GLogger.Error(errParseCIDR)
		return errParseCIDR
	}
	MaskSize, _ := IPNet.Mask.Size()
	client, err := bacnet.NewClient(&bacnet.ClientBuilder{
		Ip:           IP.String(),
		SubnetCIDR:   MaskSize,
		Port:         br.mainConfig.BacnetRouterConfig.LocalPort,
		DeviceId:     br.mainConfig.BacnetRouterConfig.DeviceId,  // RHILEX 自身的ID
		VendorId:     br.mainConfig.BacnetRouterConfig.VendorId,  // RHILEX 自身的厂家
		NetWorkId:    br.mainConfig.BacnetRouterConfig.NetWorkId, // RHILEX 自身的网络号
		PropertyData: br.selfPropertyData,                        // RHILEX 点位表
	})
	if err != nil {
		return err
	}
	go client.StartPoll(br.Ctx)
	br.bacnetClient = client
	client.SetLogger(glogger.GLogger.Logger)
	br.status = typex.DEV_UP
	return nil
}

func (br *BacnetRouter) Status() typex.DeviceState {
	return typex.DEV_UP
}

func (br *BacnetRouter) Stop() {
	br.status = typex.DEV_DOWN
	if br.CancelCTX != nil {
		br.CancelCTX()
	}
	if br.bacnetClient != nil {
		br.bacnetClient.ClientClose(false)
		br.bacnetClient.Close()
	}
	intercache.UnRegisterSlot(br.PointId)
}

func (br *BacnetRouter) Details() *typex.Device {
	return br.RuleEngine.GetDevice(br.PointId)
}

func (br *BacnetRouter) SetState(status typex.DeviceState) {
	br.status = status
}

func (br *BacnetRouter) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}

func (br *BacnetRouter) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return []byte{}, nil
}

func (br *BacnetRouter) OnRead(cmd []byte, data []byte) (int, error) {

	return 0, nil
}

func (br *BacnetRouter) OnWrite(cmd []byte, b []byte) (int, error) {
	return 0, nil
}
