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
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

//	{
//		DevEUI: "01:02:03:04:05:06:07:08",
//		AppKey: "0123456789ABCDEF0123456789ABCDE",
//		FrequencyBand: "CN470",
//		JoinType: "OTAA",
//		Channels: []string{"01", "02", "03"},
//		DataRates: []int{1, 2, 3},
//		TxPower: 14,
//		AntennaPolarization: "Linear",
//		GatewayEUI: "11:12:13:14:15:16:17:18",
//		ServerIP: "192.168.1.1",
//		ServerPort: "1234",
//	}
type LoraGatewayConfig struct {
	DevEUI              string   // 设备欧盟标识符
	AppKey              string   // 应用程序密钥
	FrequencyBand       string   // 频段
	JoinType            string   // 入网方式 (OTAA 或 ABP)
	Channels            []string // 信道列表
	DataRates           []int    // 数据速率列表
	TxPower             int      // 传输功率
	AntennaPolarization string   // 天线极化
	GatewayEUI          string   // 网关EUI
	ServerIP            string   // 服务器IP地址
	ServerPort          string   // 服务器端口
}

type LoraGatewayMainConfig struct {
	LoraGatewayConfig LoraGatewayConfig `json:"loraGatewayConfig"`
}

type LoraGateway struct {
	typex.XStatus
	status     typex.DeviceState
	mainConfig LoraGatewayMainConfig
}

func NewLoraGateway(e typex.Rhilex) typex.XDevice {
	hd := new(LoraGateway)
	hd.RuleEngine = e
	return hd
}

func (hd *LoraGateway) Init(devId string, configMap map[string]interface{}) error {
	hd.PointId = devId
	if err := utils.BindSourceConfig(configMap, &hd.mainConfig); err != nil {
		glogger.GLogger.Error(err)
		return err
	}

	return nil
}

func (hd *LoraGateway) Start(cctx typex.CCTX) error {
	hd.Ctx = cctx.Ctx
	hd.CancelCTX = cctx.CancelCTX

	hd.status = typex.DEV_UP
	return nil
}

func (hd *LoraGateway) Status() typex.DeviceState {
	return hd.status
}

func (hd *LoraGateway) Stop() {
	hd.status = typex.DEV_DOWN
	hd.CancelCTX()
}

func (hd *LoraGateway) Details() *typex.Device {
	return hd.RuleEngine.GetDevice(hd.PointId)
}

func (hd *LoraGateway) SetState(status typex.DeviceState) {
	hd.status = status
}

func (hd *LoraGateway) OnDCACall(UUID string, Command string, Args interface{}) typex.DCAResult {
	return typex.DCAResult{}
}

func (hd *LoraGateway) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return []byte{}, nil
}

func (hd *LoraGateway) OnRead(cmd []byte, data []byte) (int, error) {

	return 0, nil
}

func (hd *LoraGateway) OnWrite(cmd []byte, b []byte) (int, error) {
	return 0, nil
}
