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

package target

import (
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

type TemplateTargetConfig struct {
}
type TemplateTarget struct {
	typex.XStatus
	mainConfig TemplateTargetConfig
	status     typex.SourceState
}

func NewTemplateTarget(e typex.Rhilex) typex.XTarget {
	ht := new(TemplateTarget)
	ht.RuleEngine = e
	ht.mainConfig = TemplateTargetConfig{}
	ht.status = typex.SOURCE_DOWN
	return ht
}

func (ht *TemplateTarget) Init(outEndId string, configMap map[string]any) error {
	ht.PointId = outEndId

	if err := utils.BindSourceConfig(configMap, &ht.mainConfig); err != nil {
		return err
	}

	return nil

}
func (ht *TemplateTarget) Start(cctx typex.CCTX) error {
	ht.Ctx = cctx.Ctx
	ht.CancelCTX = cctx.CancelCTX
	ht.status = typex.SOURCE_UP
	glogger.GLogger.Info("Template Target started")
	return nil
}

func (ht *TemplateTarget) Status() typex.SourceState {

	return ht.status

}
func (ht *TemplateTarget) To(data any) (any, error) {
	return 0, nil
}

func (ht *TemplateTarget) Stop() {
	ht.status = typex.SOURCE_DOWN
	if ht.CancelCTX != nil {
		ht.CancelCTX()
	}
}
func (ht *TemplateTarget) Details() *typex.OutEnd {
	return ht.RuleEngine.GetOutEnd(ht.PointId)
}
