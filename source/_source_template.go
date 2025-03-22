// Copyright (C) 2025 wwhai
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
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package source

import (
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

type MySourceConfig struct {
}

// MySource 实现 XSource 接口的具体结构体
type MySource struct {
	typex.XStatus
	mainConfig MySourceConfig
}

// Init 实现 Init 方法
func (ms *MySource) Init(inEndId string, configMap map[string]any) error {
	ms.PointId = inEndId
	if err := utils.BindSourceConfig(configMap, &ms.mainConfig); err != nil {
		glogger.GLogger.Errorf("Failed to bind source config: %v", err)
		return err
	}
	return nil
}

// Start 实现 Start 方法
func (ms *MySource) Start(cctx typex.CCTX) error {

	return nil
}

// Status 实现 Status 方法
func (ms *MySource) Status() typex.SourceState {
	return ms.SourceState
}

// Details 实现 Details 方法
func (ms *MySource) Details() *typex.InEnd {
	return ms.RuleEngine.GetInEnd(ms.PointId)
}

// Stop 实现 Stop 方法
func (ms *MySource) Stop() {
	ms.SourceState = typex.SOURCE_DOWN
	if ms.CancelCTX != nil {
		ms.CancelCTX()
	}
}
