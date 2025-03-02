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

type TemplateDeviceCommonConfig struct {
}

type TemplateDeviceConfig struct {
}

type TemplateDeviceMainConfig struct {
	CommonConfig         TemplateDeviceCommonConfig `json:"commonConfig" validate:"required"`
	TemplateDeviceConfig TemplateDeviceConfig       `json:"templateDeviceConfig"`
}

type TemplateDevice struct {
	typex.XStatus
	status     typex.SourceState
	mainConfig TemplateDeviceMainConfig
}

func NewTemplateDevice(e typex.Rhilex) typex.XDevice {
	hd := new(TemplateDevice)
	hd.RuleEngine = e
	return hd
}

func (hd *TemplateDevice) Init(devId string, configMap map[string]any) error {
	hd.PointId = devId
	if err := utils.BindSourceConfig(configMap, &hd.mainConfig); err != nil {
		glogger.GLogger.Error(err)
		return err
	}

	return nil
}

func (hd *TemplateDevice) Start(cctx typex.CCTX) error {
	hd.Ctx = cctx.Ctx
	hd.CancelCTX = cctx.CancelCTX

	hd.status = typex.SOURCE_UP
	return nil
}

func (hd *TemplateDevice) Status() typex.SourceState {
	return hd.Status
}

func (hd *TemplateDevice) Stop() {
	hd.status = typex.SOURCE_DOWN
	hd.CancelCTX()
}

func (hd *TemplateDevice) Details() *typex.Device {
	return hd.RuleEngine.GetDevice(hd.PointId)
}

func (hd *TemplateDevice) SetState(status typex.SourceState) {
	hd.status = status
}

func (hd *TemplateDevice) OnDCACall(UUID string, Command string, Args any) typex.DCAResult {
	return typex.DCAResult{}
}

func (hd *TemplateDevice) OnCtrl(cmd []byte, args []byte) ([]byte, error) {
	return []byte{}, nil
}
