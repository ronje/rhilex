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

package rhilexmanager

import (
	ithings "github.com/hootrhino/rhilex/component/cecolla/ithings"
	tencent "github.com/hootrhino/rhilex/component/cecolla/tencent"
	"github.com/hootrhino/rhilex/component/orderedmap"
	"github.com/hootrhino/rhilex/typex"
)

var DefaultCecollaTypeManager *CecollaTypeManager

type CecollaTypeManager struct {
	e        typex.Rhilex
	registry *orderedmap.OrderedMap[typex.CecollaType, *typex.XConfig]
}

func InitCecollaTypeManager(e typex.Rhilex) {
	DefaultCecollaTypeManager = &CecollaTypeManager{
		e:        e,
		registry: orderedmap.NewOrderedMap[typex.CecollaType, *typex.XConfig](),
	}
	LoadAllCecType(e)
}

func LoadAllCecType(e typex.Rhilex) {

	DefaultCecollaTypeManager.Register(typex.TENCENT_IOTHUB_CEC,
		&typex.XConfig{
			Engine:     e,
			NewCecolla: tencent.NewTencentIoTGateway,
		},
	)
	DefaultCecollaTypeManager.Register(typex.ITHINGS_IOTHUB_CEC,
		&typex.XConfig{
			Engine:     e,
			NewCecolla: ithings.NewIThingsGateway,
		},
	)
}
func (rm *CecollaTypeManager) Register(name typex.CecollaType, f *typex.XConfig) {
	f.Type = string(name)
	rm.registry.Set(name, f)
}

func (rm *CecollaTypeManager) Find(name typex.CecollaType) *typex.XConfig {
	if xcfg, ok := rm.registry.Get(name); ok {
		return xcfg
	}
	return nil
}
func (rm *CecollaTypeManager) All() []*typex.XConfig {
	return rm.registry.Values()
}

/**
 * 获取所有类型
 *
 */
func (rm *CecollaTypeManager) AllKeys() []string {
	data := []string{}
	for _, k := range rm.registry.Keys() {
		data = append(data, k.String())
	}
	return data
}

func (rm *CecollaTypeManager) Stop() {
}
