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

package globalinit

import (
	"log"

	"github.com/hootrhino/rhilex/typex"
)

var __DefaultGlobalInitManager *GlobalInitManager

type InitFunc struct {
	Name string
	Func func()
}

// GlobalInitManager 用于管理注册的函数
type GlobalInitManager struct {
	r         typex.Rhilex
	functions []InitFunc // 存储注册的函数
}

func NewGlobalInitManager(r typex.Rhilex) *GlobalInitManager {
	return &GlobalInitManager{r: r, functions: []InitFunc{}}
}

// Register 函数用于注册一个新的函数
func (fm *GlobalInitManager) Register(fn InitFunc) {
	fm.functions = append(fm.functions, fn)
}

// Start 执行所有已注册的函数
func (fm *GlobalInitManager) CallAllFunc() {
	for _, fn := range fm.functions {
		log.Println("Executing init function:", fn.Name)
		fn.Func()
	}
}
func InitGlobalInitManager(r typex.Rhilex) {
	__DefaultGlobalInitManager = NewGlobalInitManager(r)
}

func Register(fn InitFunc) {
	__DefaultGlobalInitManager.Register(fn)
}
