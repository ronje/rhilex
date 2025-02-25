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

package tunplugin

import (
	"github.com/hootrhino/rhilex/typex"
	"gopkg.in/ini.v1"
)

type TunPlugin struct {
}

func NewTunPlugin() *TunPlugin {
	return &TunPlugin{}
}

func (dm *TunPlugin) Init(config *ini.Section) error {
	return nil
}

func (dm *TunPlugin) Start(typex.Rhilex) error {
	return nil
}
func (dm *TunPlugin) Stop() error {
	return nil
}

func (dm *TunPlugin) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		UUID:        "MICRO-TUNNING",
		Name:        "Simple Tunnel Plugin",
		Version:     "v0.0.1",
		Description: "Simple Tunnel Plugin Can Forward Network Traffic",
	}
}

func (dm *TunPlugin) Service(arg typex.ServiceArg) typex.ServiceResult {
	return typex.ServiceResult{}
}
