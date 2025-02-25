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

package engine

import (
	"encoding/json"
	"fmt"
	"runtime"

	lua "github.com/hootrhino/gopher-lua"

	"github.com/hootrhino/rhilex/component/intermetric"
	"github.com/hootrhino/rhilex/component/interqueue"
	"github.com/hootrhino/rhilex/component/luaexecutor"
	"github.com/hootrhino/rhilex/component/orderedmap"
	core "github.com/hootrhino/rhilex/config"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/registry"
	"github.com/hootrhino/rhilex/typex"
	"github.com/hootrhino/rhilex/utils"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/sirupsen/logrus"
)

// 规则引擎
type RuleEngine struct {
	Rules   *orderedmap.OrderedMap[string, *typex.Rule]   `json:"rules"`
	InEnds  *orderedmap.OrderedMap[string, *typex.InEnd]  `json:"inends"`
	OutEnds *orderedmap.OrderedMap[string, *typex.OutEnd] `json:"outends"`
	Devices *orderedmap.OrderedMap[string, *typex.Device] `json:"devices"`
	Config  *typex.RhilexConfig                           `json:"config"`
}

func NewRuleEngine(config typex.RhilexConfig) typex.Rhilex {
	return &RuleEngine{
		Rules:   orderedmap.NewOrderedMap[string, *typex.Rule](),
		InEnds:  orderedmap.NewOrderedMap[string, *typex.InEnd](),
		OutEnds: orderedmap.NewOrderedMap[string, *typex.OutEnd](),
		Devices: orderedmap.NewOrderedMap[string, *typex.Device](),
		Config:  &config,
	}
}

/*
*
* Engine Start
*
 */
func (e *RuleEngine) Start() *typex.RhilexConfig {
	return e.Config
}
func (e *RuleEngine) Version() typex.VersionInfo {
	return typex.DefaultVersionInfo
}

func (e *RuleEngine) GetConfig() *typex.RhilexConfig {
	return e.Config
}

// Stop
func (e *RuleEngine) Stop() {
	glogger.GLogger.Info("Ready to stop RHILEX")
	// 资源 TODO: 后期重构设备资源等，独立资源管理器。
	for _, inEnd := range e.InEnds.Values() {
		if inEnd.Source != nil {
			glogger.GLogger.Infof("Stop InEnd:(%s,%s)", inEnd.Name, inEnd.UUID)
			e.GetInEnd(inEnd.UUID).State = typex.SOURCE_STOP
			inEnd.Source.Stop()
		}
		glogger.GLogger.Infof("Stop InEnd:(%s,%s) Successfully", inEnd.Name, inEnd.UUID)
	}
	for _, outEnd := range e.OutEnds.Values() {
		if outEnd.Target != nil {
			glogger.GLogger.Infof("Stop NewTarget:(%s,%s)", outEnd.Name, outEnd.UUID)
			e.GetOutEnd(outEnd.UUID).State = typex.SOURCE_STOP
			outEnd.Target.Stop()
			glogger.GLogger.Infof("Stop NewTarget:(%s,%s) Successfully", outEnd.Name, outEnd.UUID)
		}
	}
	for _, device := range e.Devices.Values() {
		glogger.GLogger.Infof("Stop Device:(%s)", device.Name)
		e.GetDevice(device.UUID).State = typex.DEV_STOP
		device.Device.Stop()
		glogger.GLogger.Infof("Stop Device:(%s) Successfully", device.Name)
	}
}

// 核心功能: Work, 主要就是推流进队列
func (e *RuleEngine) WorkInEnd(in *typex.InEnd, data string) (bool, error) {
	if err := interqueue.PushInQueue(in, data); err != nil {
		return false, err
	}
	return true, nil
}

// 核心功能: Work, 主要就是推流进队列
func (e *RuleEngine) WorkDevice(Device *typex.Device, data string) (bool, error) {
	if err := interqueue.PushDeviceQueue(Device, data); err != nil {
		return false, err
	}
	return true, nil
}

// RunSourceCallbacks 执行针对资源端的规则脚本
func (e *RuleEngine) RunSourceCallbacks(in *typex.InEnd, callbackArgs string) {
	// 执行来自资源的脚本
	for _, rule := range in.BindRules {
		if rule.Status == typex.RULE_RUNNING {
			if err := e.executeRule(&rule, callbackArgs); err != nil {
				return
			}
		}
	}
}

// RunDeviceCallbacks 执行针对设备端的规则脚本
func (e *RuleEngine) RunDeviceCallbacks(Device *typex.Device, callbackArgs string) {
	for _, rule := range Device.BindRules {
		if err := e.executeRule(&rule, callbackArgs); err != nil {
			return
		}
	}
}

// executeRule 执行规则脚本的公共逻辑
func (e *RuleEngine) executeRule(rule *typex.Rule, callbackArgs string) error {
	// 执行规则动作
	_, errA := luaexecutor.ExecuteActions(rule, lua.LString(callbackArgs))
	if errA != nil {
		e.handleError(rule, errA)
		return errA
	}

	// 执行成功处理
	_, errS := luaexecutor.ExecuteSuccess(rule.LuaVM)
	if errS != nil {
		glogger.GLogger.Error(errS)
		return errS
	}

	return nil
}

// handleError 处理规则执行错误
func (e *RuleEngine) handleError(rule *typex.Rule, err error) {
	Debugger, Ok := rule.LuaVM.GetStack(1)
	if !Ok {
		return
	}

	LValue, _ := rule.LuaVM.GetInfo("f", Debugger, lua.LNil)
	rule.LuaVM.GetInfo("l", Debugger, lua.LNil)
	rule.LuaVM.GetInfo("S", Debugger, lua.LNil)
	rule.LuaVM.GetInfo("u", Debugger, lua.LNil)
	rule.LuaVM.GetInfo("n", Debugger, lua.LNil)

	LFunction := LValue.(*lua.LFunction)
	LastCall := lua.DbgCall{
		Name: "_main",
	}
	if len(LFunction.Proto.DbgCalls) > 0 {
		LastCall = LFunction.Proto.DbgCalls[0]
	}

	glogger.GLogger.WithFields(logrus.Fields{
		"topic": "rule/log/" + rule.UUID,
	}).Warnf("Function Name: [%s],"+
		"What: [%s], Source Line: [%d],"+
		" Last Call: [%s], Error message: %s",
		Debugger.Name, Debugger.What, Debugger.CurrentLine,
		LastCall.Name, err.Error(),
	)
}

func (e *RuleEngine) GetInEnd(uuid string) *typex.InEnd {
	v, ok := (e.InEnds).Get(uuid)
	if ok {
		return v
	}
	return nil
}

func (e *RuleEngine) SaveInEnd(in *typex.InEnd) {
	e.InEnds.Set(in.UUID, in)
}

func (e *RuleEngine) RemoveInEnd(uuid string) {
	if inEnd := e.GetInEnd(uuid); inEnd != nil {
		if inEnd.Source != nil {
			glogger.GLogger.Infof("InEnd [%s, %s] ready to stop", uuid, inEnd.Name)
			inEnd.Source.Stop()
			glogger.GLogger.Infof("InEnd [%s, %s] stopped", uuid, inEnd.Name)
			e.InEnds.Delete(uuid)
			glogger.GLogger.Infof("InEnd [%s, %s] has been deleted", uuid, inEnd.Name)
			inEnd = nil
		}
	}
}

func (e *RuleEngine) AllInEnds() []*typex.InEnd {
	return e.InEnds.Values()
}

func (e *RuleEngine) GetOutEnd(id string) *typex.OutEnd {
	v, ok := e.OutEnds.Get(id)
	if ok {
		return v
	}
	return nil
}

func (e *RuleEngine) SaveOutEnd(out *typex.OutEnd) {
	e.OutEnds.Set(out.UUID, out)
}

func (e *RuleEngine) RemoveOutEnd(uuid string) {
	if outEnd := e.GetOutEnd(uuid); outEnd != nil {
		if outEnd.Target != nil {
			glogger.GLogger.Infof("OutEnd [%s, %s] ready to stop", uuid, outEnd.Name)
			outEnd.Target.Stop()
			glogger.GLogger.Infof("OutEnd [%s, %s] stopped", uuid, outEnd.Name)
			e.OutEnds.Delete(uuid)
			glogger.GLogger.Infof("OutEnd [%s, %s] has been deleted", uuid, outEnd.Name)
			outEnd = nil
		}
	}
}

func (e *RuleEngine) AllOutEnds() []*typex.OutEnd {
	return e.OutEnds.Values()
}

// -----------------------------------------------------------------
// 获取运行时快照
// -----------------------------------------------------------------
func (e *RuleEngine) SnapshotDump() string {
	parts, _ := disk.Partitions(true)
	diskInfo, _ := disk.Usage(parts[0].Mountpoint)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	system := map[string]any{
		"version":  typex.MainVersion,
		"diskInfo": int(diskInfo.UsedPercent),
		"system":   utils.BToMb(m.Sys),
		"alloc":    utils.BToMb(m.Alloc),
		"total":    utils.BToMb(m.TotalAlloc),
		"osArch":   runtime.GOOS + "-" + runtime.GOARCH,
	}
	data := map[string]any{
		"rules":      e.Rules.Values(),
		"inends":     e.InEnds.Values(),
		"outends":    e.OutEnds.Values(),
		"devices":    e.Devices.Values(),
		"statistics": intermetric.GetMetric(),
		"system":     system,
		"config":     core.GlobalConfig,
	}
	b, err := json.Marshal(data)
	if err != nil {
		glogger.GLogger.Error(err)
		return err.Error()
	}
	return string(b)
}

// 重启源
func (e *RuleEngine) RestartInEnd(uuid string) error {
	if InEnd, ok := e.InEnds.Get(uuid); ok {
		if InEnd.Source != nil {
			InEnd.Source.Details().State = typex.SOURCE_DOWN // Down 以后会被自动拉起来
		}
		return nil
	}
	return fmt.Errorf("inend not exists:%s", uuid)
}

// 重启目标
func (e *RuleEngine) RestartOutEnd(uuid string) error {
	if OutEnd, ok := e.OutEnds.Get(uuid); ok {
		if OutEnd.Target != nil {
			OutEnd.Target.Details().State = typex.SOURCE_DOWN // Down 以后会被自动拉起来
		}
		return nil
	}
	return fmt.Errorf("outEnd not exists:%s", uuid)

}

// 重启设备
func (e *RuleEngine) RestartDevice(uuid string) error {
	if Device, ok := e.Devices.Get(uuid); ok {
		if Device.Device != nil {
			Device.Device.SetState(typex.DEV_DOWN) // Down 以后会被自动拉起来
		}
		return nil
	}
	return fmt.Errorf("device not exists:%s", uuid)
}

func (e *RuleEngine) CheckSourceType(Type typex.InEndType) error {
	keys := registry.DefaultSourceRegistry.AllKeys()
	if utils.SContains(keys, string(Type)) {
		return nil
	}
	return fmt.Errorf("Source Type Not Support:%s", Type)
}

// 0.7.0
// 更新设备的运行时状态
func (e *RuleEngine) SetDeviceStatus(uuid string, DeviceState typex.DeviceState) {
	Device := e.GetDevice(uuid)
	if Device != nil {
		Device.State = DeviceState
	}
}
func (e *RuleEngine) SetSourceStatus(uuid string, SourceState typex.SourceState) {
	Source := e.GetInEnd(uuid)
	if Source != nil {
		Source.State = SourceState
	}
}
func (e *RuleEngine) SetTargetStatus(uuid string, SourceState typex.SourceState) {
	Outend := e.GetOutEnd(uuid)
	if Outend != nil {
		Outend.State = SourceState
	}
}
func (e *RuleEngine) CheckDeviceType(Type typex.DeviceType) error {
	keys := registry.DefaultDeviceRegistry.AllKeys()
	if utils.SContains(keys, string(Type)) {
		return nil
	}
	return fmt.Errorf("Device Type Not Support:%s", Type)
}

func (e *RuleEngine) CheckTargetType(Type typex.TargetType) error {
	keys := registry.DefaultTargetRegistry.AllKeys()
	if utils.SContains(keys, string(Type)) {
		return nil
	}
	return fmt.Errorf("Target Type Not Support:%s", Type)
}
