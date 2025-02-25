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
	"github.com/vapourismo/knx-go/knx"
)

type KNXGatewayCommonConfig struct {
}

type KNXGatewayConfig struct {
}

type KNXGatewayMainConfig struct {
	CommonConfig     KNXGatewayCommonConfig `json:"commonConfig" validate:"required"`
	KNXGatewayConfig KNXGatewayConfig       `json:"knxConfig"`
}

type KNXGateway struct {
	typex.XStatus
	status      typex.DeviceState
	mainConfig  KNXGatewayMainConfig
	GroupTunnel knx.GroupTunnel
}

func NewKNXGateway(e typex.Rhilex) typex.XDevice {
	hd := new(KNXGateway)
	hd.RuleEngine = e
	return hd
}

func (hd *KNXGateway) Init(devId string, configMap map[string]any) error {
	hd.PointId = devId
	if err := utils.BindSourceConfig(configMap, &hd.mainConfig); err != nil {
		glogger.GLogger.Error(err)
		return err
	}

	return nil
}

func (hd *KNXGateway) Start(cctx typex.CCTX) error {
	hd.Ctx = cctx.Ctx
	hd.CancelCTX = cctx.CancelCTX
	GroupTunnel, err := knx.NewGroupTunnel("127.0.0.1:3671", knx.DefaultTunnelConfig)
	if err != nil {
		return err
	}
	hd.GroupTunnel = GroupTunnel
	hd.status = typex.DEV_UP
	return nil
}

func (hd *KNXGateway) Status() typex.DeviceState {
	return hd.status
}

func (hd *KNXGateway) Stop() {
	hd.status = typex.DEV_DOWN
	hd.CancelCTX()
}

func (hd *KNXGateway) Details() *typex.Device {
	return hd.RuleEngine.GetDevice(hd.PointId)
}

func (hd *KNXGateway) SetState(status typex.DeviceState) {
	hd.status = status
}

func (hd *KNXGateway) OnDCACall(UUID string, Command string, Args any) typex.DCAResult {
	return typex.DCAResult{}
}

func (hd *KNXGateway) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return []byte{}, nil
}
