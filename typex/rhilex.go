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
	"context"
)

// Rhilex interface
type Rhilex interface {
	//
	// 启动规则引擎
	//
	Start() *RhilexConfig

	//
	// 执行任务
	//
	WorkInEnd(*InEnd, string) (bool, error)
	WorkDevice(*Device, string) (bool, error)
	//
	// 获取配置
	//
	GetConfig() *RhilexConfig
	//
	// 加载输入
	//
	LoadInEndWithCtx(in *InEnd, ctx context.Context, cancelCTX context.CancelFunc) error
	//
	// 获取输入
	//
	GetInEnd(string) *InEnd
	//
	// 保存输入
	//
	SaveInEnd(*InEnd)
	//
	// 删除输入
	//
	RemoveInEnd(string)
	//
	// 所有输入列表
	//
	AllInEnds() []*InEnd
	//
	// 加载输出
	//
	LoadOutEndWithCtx(in *OutEnd, ctx context.Context, cancelCTX context.CancelFunc) error
	//
	// 所有输出
	//
	AllOutEnds() []*OutEnd
	//
	// 获取输出
	//
	GetOutEnd(string) *OutEnd
	//
	// 保存输出
	//
	SaveOutEnd(*OutEnd)
	//
	// 删除输出
	//
	RemoveOutEnd(string)
	//
	// 加载规则
	//
	LoadRule(*Rule) error
	//
	// 所有规则列表
	//
	AllRules() []*Rule
	//
	// 获取规则
	//
	GetRule(id string) *Rule
	//
	// 删除规则
	//
	RemoveRule(uuid string)
	//
	// 获取版本
	//
	Version() VersionInfo

	//
	// 停止规则引擎
	//
	Stop()
	//
	// Snapshot Dump
	//
	SnapshotDump() string
	//
	// 加载设备
	//
	LoadDeviceWithCtx(in *Device, ctx context.Context, cancelCTX context.CancelFunc) error
	//
	// 获取设备
	//
	GetDevice(string) *Device

	//
	// 保存设备
	//
	SaveDevice(*Device)
	//
	//
	//
	AllDevices() []*Device
	//
	// 删除设备
	//
	RemoveDevice(string)
	//
	// 重启源
	//
	RestartInEnd(uuid string) error
	//
	// 重启目标
	//
	RestartOutEnd(uuid string) error

	//
	// 重启设备
	//
	RestartDevice(uuid string) error
	//
	SetDeviceStatus(uuid string, s DeviceState)
	//
	SetSourceStatus(uuid string, s SourceState)
	// 检查类型是否支持
	CheckSourceType(Type InEndType) error
	CheckDeviceType(Type DeviceType) error
	CheckTargetType(Type TargetType) error
	// 云边协同
	CheckCecollaType(Type CecollaType) error
	GetCecolla(string) *Cecolla
	SaveCecolla(*Cecolla)
	AllCecollas() []*Cecolla
	RestartCecolla(uuid string) error
	RemoveCecolla(uuid string)
	LoadCecollaWithCtx(cecolla *Cecolla, ctx context.Context, cancelCTX context.CancelFunc) error
}
