// Copyright (C) 2025 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the tera3  of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package source

import (
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

// 3轴加速度传感器，IIC协议
type Axis3AccSensorConfig struct {
	Address uint16 `json:"address"`
}

// Axis3AccSensor 实现 XSource 接口的具体结构体
type Axis3AccSensor struct {
	typex.XStatus
	mainConfig Axis3AccSensorConfig
}

func NewAxis3AccSensor(e typex.Rhilex) typex.XSource {
	c := Axis3AccSensor{
		mainConfig: Axis3AccSensorConfig{Address: 0x18},
	}
	c.RuleEngine = e
	return &c
}

// Init 实现 Init 方法
func (a3 *Axis3AccSensor) Init(inEndId string, configMap map[string]any) error {
	a3.PointId = inEndId
	if err := utils.BindSourceConfig(configMap, &a3.mainConfig); err != nil {
		glogger.GLogger.Errorf("Failed to bind source config: %v", err)
		return err
	}
	return nil
}

// Start 实现 Start 方法
func (a3 *Axis3AccSensor) Start(cctx typex.CCTX) error {

	return nil
}

// Status 实现 Status 方法
func (a3 *Axis3AccSensor) Status() typex.SourceState {
	return a3.SourceState
}

// Details 实现 Details 方法
func (a3 *Axis3AccSensor) Details() *typex.InEnd {
	return a3.RuleEngine.GetInEnd(a3.PointId)
}

// Stop 实现 Stop 方法
func (a3 *Axis3AccSensor) Stop() {
	a3.SourceState = typex.SOURCE_DOWN
	if a3.CancelCTX != nil {
		a3.CancelCTX()
	}
}
