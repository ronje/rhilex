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

package typex

import (
	"fmt"

	"gopkg.in/ini.v1"
)

type PluginType int

func (s PluginType) String() string {
	return fmt.Sprintf("%d", s)
}

//
// 插件开发步骤：
// 1 在ini配置文件中增加配置项
// 2 实现插件接口: XPlugin
// 3 LoadPlugin(sectionK string, p typex.XPlugin)
//

// 插件的服务参数
type ServiceArg struct {
	UUID string `json:"uuid"` // 插件UUID, Rhilex用来查找插件的
	Name string `json:"name"` // 服务名, 在服务中响应识别
	Args any    `json:"args"` // 服务参数
}
type ServiceResult struct {
	Out any `json:"out"`
}

/*
*
* 插件的元信息结构体
*   注意：插件信息这里uuid，name有些是固定写死的，比较特殊，不要轻易改变已有的，否则会导致接口失效
*        只要是已有的尽量不要改这个UUID。
*
 */
type XPluginMetaInfo struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

/*
*
* 插件: 用来增强RHILEX的外部功能，本色不属于RHILEX
*
 */
type XPlugin interface {
	Init(*ini.Section) error          // 参数为外部配置
	Start(Rhilex) error               // 启动插件
	Service(ServiceArg) ServiceResult // 对外提供一些服务
	Stop() error                      // 停止插件
	PluginMetaInfo() XPluginMetaInfo  // 插件的元信息
}
