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
	"fmt"
	"sync"

	core "github.com/hootrhino/rhilex/config"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
)

var DefaultPluginTypeManager *PluginTypeManager

type PluginTypeManager struct {
	e        typex.Rhilex
	registry sync.Map
}

func InitPluginTypeManager(e typex.Rhilex) {
	DefaultPluginTypeManager = &PluginTypeManager{
		e:        e,
		registry: sync.Map{},
	}
}

func (rm *PluginTypeManager) All() []typex.XPlugin {
	data := make([]typex.XPlugin, 0)
	rm.registry.Range(func(key, value any) bool {
		data = append(data, value.(typex.XPlugin))
		return true
	})
	return data
}
func (rm *PluginTypeManager) Count() int {
	count := int(0)
	rm.registry.Range(func(key, value any) bool {
		count++
		return true
	})
	return count
}
func (rm *PluginTypeManager) Find(name string) typex.XPlugin {
	p, ok := rm.registry.Load(name)
	if ok {
		return p.(typex.XPlugin)
	}
	return nil
}
func (rm *PluginTypeManager) LoadPlugin(sectionK string, p typex.XPlugin) error {
	section := utils.GetINISection(core.GlobalConfig.IniPath, sectionK)
	if err := p.Init(section); err != nil {
		return err
	}
	_, ok := rm.registry.Load(p.PluginMetaInfo().UUID)
	if ok {
		return fmt.Errorf("plugin already installed:" + p.PluginMetaInfo().Name)
	}
	if err := p.Start(rm.e); err != nil {
		return err
	}
	if p.PluginMetaInfo().UUID != "LicenseManager" {
		rm.registry.Store(p.PluginMetaInfo().UUID, p)
		glogger.GLogger.Infof("Plugin start successfully:[%v]", p.PluginMetaInfo().Name)
	}
	return nil

}
func (rm *PluginTypeManager) Stop() {
	rm.registry.Range(func(key, value any) bool {
		plugin := value.(typex.XPlugin)
		glogger.GLogger.Infof("Stop plugin:(%s)", plugin.PluginMetaInfo().Name)
		plugin.Stop()
		glogger.GLogger.Infof("Stop plugin:(%s) Successfully", plugin.PluginMetaInfo().Name)
		return true
	})
}
