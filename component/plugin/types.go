// Copyright (C) 2024 wwhai
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

package plugin

import (
	"github.com/hootrhino/rhilex/typex"
	"gopkg.in/ini.v1"
)

// 插件的服务参数
type ServiceArg struct {
	UUID string      `json:"uuid"` // 插件UUID, Rhilex用来查找插件的
	Name string      `json:"name"` // 服务名, 在服务中响应识别
	Args interface{} `json:"args"` // 服务参数
}
type ServiceResult struct {
	Out interface{} `json:"out"`
}

type XPluginMetaInfo struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

type XPlugin interface {
	Init(*ini.Section) error
	Start(typex.Rhilex) error
	Service(ServiceArg) ServiceResult
	Stop() error
	PluginMetaInfo() XPluginMetaInfo
}
