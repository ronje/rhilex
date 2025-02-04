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

package plugin

import (
	"fmt"

	"github.com/hootrhino/rhilex/component/orderedmap"
	core "github.com/hootrhino/rhilex/config"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

var __DefaultPluginRegistry *PluginRegistry

type PluginRegistry struct {
	e        typex.Rhilex
	registry *orderedmap.OrderedMap[string, typex.XPlugin]
}

func InitPluginRegistry(e typex.Rhilex) {
	__DefaultPluginRegistry = &PluginRegistry{
		e:        e,
		registry: orderedmap.NewOrderedMap[string, typex.XPlugin](),
	}
}

func (rm *PluginRegistry) All() []typex.XPlugin {
	return rm.registry.Values()
}

func (rm *PluginRegistry) Count() int {
	return rm.registry.Size()
}

func (rm *PluginRegistry) Find(name string) typex.XPlugin {
	p, ok := rm.registry.Get(name)
	if ok {
		return p
	}
	return nil
}
func (rm *PluginRegistry) LoadPlugin(sectionK string, p typex.XPlugin) error {
	section := utils.GetINISection(core.GlobalConfig.IniPath, sectionK)
	if err := p.Init(section); err != nil {
		return err
	}
	_, ok := rm.registry.Get(p.PluginMetaInfo().UUID)
	if ok {
		return fmt.Errorf("plugin already installed:%s", p.PluginMetaInfo().Name)
	}
	if err := p.Start(rm.e); err != nil {
		return err
	}
	if p.PluginMetaInfo().UUID != "LicenseManager" {
		rm.registry.Set(p.PluginMetaInfo().UUID, p)
		glogger.GLogger.Infof("Plugin start successfully:[%v]", p.PluginMetaInfo().Name)
	}
	return nil

}

func Stop() {
	for _, plugin := range __DefaultPluginRegistry.registry.Values() {
		glogger.GLogger.Infof("Stop plugin:(%s)", plugin.PluginMetaInfo().Name)
		plugin.Stop()
		glogger.GLogger.Infof("Stop plugin:(%s) Successfully", plugin.PluginMetaInfo().Name)
	}
}

// __DefaultPluginRegistry 作为包级别的变量，在程序启动时被初始化，对外提供封装后的接口
func LoadPlugin(sectionK string, p typex.XPlugin) error {
	return __DefaultPluginRegistry.LoadPlugin(sectionK, p)
}
func All() []typex.XPlugin {
	return __DefaultPluginRegistry.All()
}
func Count() int {
	return __DefaultPluginRegistry.Count()
}
func Find(name string) typex.XPlugin {
	return __DefaultPluginRegistry.Find(name)
}
