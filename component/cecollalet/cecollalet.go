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

package cecollalet

import (
	"context"
	"log"
	"runtime"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rhilex/typex"
)

// lua 虚拟机的参数
const p_VM_Registry_Size int = 1024 * 1024    // 默认堆栈大小
const p_VM_Registry_MaxSize int = 1024 * 1024 // 默认最大堆栈
const p_VM_Registry_GrowStep int = 32         // 默认CPU消耗
type CecollaletState int

/*
*
* 轻量级应用
*
 */
type Cecollalet struct {
	UUID            string             `json:"uuid"`            // 名称
	Name            string             `json:"name"`            // 名称
	Version         string             `json:"version"`         // 版本号
	Description     string             `json:"description"`     // 版本号
	AutoStart       bool               `json:"autoStart"`       // 自动启动
	CecollaletState CecollaletState    `json:"cecollaletState"` // 状态: 1 运行中, 0 停止
	KilledBy        string             `json:"-"`               // 被谁杀死的: RHILEX|EXCEPT|NORMAL|""
	luaMainFunc     *lua.LFunction     `json:"-"`               // Main
	vm              *lua.LState        `json:"-"`               // lua 环境
	ctx             context.Context    `json:"-"`               // context
	cancel          context.CancelFunc `json:"-"`               // Cancel
}

func NewCecollalet(uuid, Name, Version string) *Cecollalet {
	cecollalet := new(Cecollalet)
	cecollalet.Name = Name
	cecollalet.UUID = uuid
	cecollalet.Version = Version
	cecollalet.KilledBy = "NORMAL"
	cecollalet.vm = lua.NewState(lua.Options{
		RegistrySize:     p_VM_Registry_Size,
		RegistryMaxSize:  p_VM_Registry_MaxSize,
		RegistryGrowStep: p_VM_Registry_GrowStep,
	})
	return cecollalet
}

func (cecollalet *Cecollalet) SetCnC(ctx context.Context, cancel context.CancelFunc) {
	cecollalet.ctx = ctx
	cecollalet.cancel = cancel
	cecollalet.vm.SetContext(cecollalet.ctx)
}
func (cecollalet *Cecollalet) SetMainFunc(f *lua.LFunction) {
	cecollalet.luaMainFunc = f
}
func (cecollalet *Cecollalet) GetMainFunc() *lua.LFunction {
	return cecollalet.luaMainFunc
}

func (cecollalet *Cecollalet) VM() *lua.LState {
	return cecollalet.vm
}

/*
*
* 源码bug，没有等字节码执行结束就直接给释放stack了，问题处在state.go:1391, 已经给作者提了issue，
* 如果1个月内不解决，准备自己fork一个过来维护.
* Issue: https://github.com/hootrhino/gopher-lua/discussions/430
 */
func (cecollalet *Cecollalet) Stop() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("[gopher-lua] cecollalet Stop:", cecollalet.UUID, ", with recover error: ", err)
		}
	}()
	cecollalet.CecollaletState = 0
	if cecollalet.cancel != nil {
		cecollalet.cancel()
	}
}

/*
*
* 清理内存
*
 */
func (cecollalet *Cecollalet) Remove() {
	cecollalet.Stop()
	runtime.GC()
}

/*
*
* cecollalet Stack 管理器
*
 */
type XCecollalet interface {
	GetRhilex() typex.Rhilex
	ListCecollalet() []*Cecollalet
	LoadCecollalet(cecollalet *Cecollalet) error
	GetCecollalet(uuid string) *Cecollalet
	RemoveCecollalet(uuid string) error
	UpdateCecollalet(cecollalet Cecollalet) error
	StartCecollalet(uuid string) error
	StopCecollalet(uuid string) error
	Stop()
}
