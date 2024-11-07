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
	"fmt"
	"log"

	"github.com/hootrhino/rhilex/archsupport/en6400"
	"github.com/hootrhino/rhilex/archsupport/haas506"
	"github.com/hootrhino/rhilex/archsupport/rhilexg1"
	"github.com/hootrhino/rhilex/archsupport/rhilexpro1"
)

var __DefaultGlobalInitManager *GlobalInitManager

type InitFunc struct {
	Name string
	Func func() error
}

// GlobalInitManager 用于管理注册的函数
type GlobalInitManager struct {
	functions []InitFunc // 存储注册的函数
}

func NewGlobalInitManager() *GlobalInitManager {
	return &GlobalInitManager{functions: []InitFunc{}}
}

// Register 函数用于注册一个新的函数
func (fm *GlobalInitManager) Register(fn InitFunc) {
	fm.functions = append(fm.functions, fn)
}

// Start 执行所有已注册的函数
func (fm *GlobalInitManager) CallAllFunc() {
	for _, fn := range fm.functions {
		log.Println("Executing init function:", fn.Name)
		if err := fn.Func(); err != nil {
			panic(fmt.Sprintf("Executing init function:%s, error:%s", fn.Name, err))
		}
	}
}
func InitGlobalInitManager() {
	__DefaultGlobalInitManager = NewGlobalInitManager()
	RegisterGlobalInit()
	__DefaultGlobalInitManager.CallAllFunc()
}

func Register(fn InitFunc) {
	if __DefaultGlobalInitManager == nil {
		panic("Need Init GlobalInitManager")
	}
	__DefaultGlobalInitManager.Register(fn)
}

/**
 * 注册全局初始化函数
 *
 */
func RegisterGlobalInit() {
	__DefaultGlobalInitManager.Register(InitFunc{
		Name: "en6400",
		Func: en6400.Init_EN6400,
	})
	__DefaultGlobalInitManager.Register(InitFunc{
		Name: "haas506",
		Func: haas506.Init_HAAS506LD1,
	})
	__DefaultGlobalInitManager.Register(InitFunc{
		Name: "rhilexg1",
		Func: rhilexg1.Init_RHILEXG1,
	})
	__DefaultGlobalInitManager.Register(InitFunc{
		Name: "rhilexpro1",
		Func: rhilexpro1.Init_RHILEXPRO1,
	})

}
