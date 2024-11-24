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

package luaruntime

import (
	"errors"
	"fmt"

	lua "github.com/hootrhino/gopher-lua"
	core "github.com/hootrhino/rhilex/config"
	"github.com/hootrhino/rhilex/rhilexlib"
	"github.com/hootrhino/rhilex/typex"
)

const (
	SUCCESS_KEY string = "Success"
	FAILED_KEY  string = "Failed"
	ACTIONS_KEY string = "Actions"
)

// VerifyLuaSyntax Verify Lua Syntax
func VerifyLuaSyntax(r *typex.Rule) error {
	tempVm := lua.NewState(lua.Options{
		SkipOpenLibs:     true,
		RegistrySize:     0,
		RegistryMaxSize:  0,
		RegistryGrowStep: 0,
	})

	if err := tempVm.DoString(r.Success); err != nil {
		return err
	}
	if tempVm.GetGlobal(SUCCESS_KEY).Type() != lua.LTFunction {
		return errors.New("'Success' callback function missed")
	}

	if err := tempVm.DoString(r.Failed); err != nil {
		return err
	}
	if tempVm.GetGlobal(FAILED_KEY).Type() != lua.LTFunction {
		return errors.New("'Failed' callback function missed")
	}
	if err := tempVm.DoString(r.Actions); err != nil {
		return err
	}
	//
	// validate lua syntax
	//
	actionsTable := tempVm.GetGlobal(ACTIONS_KEY)
	if actionsTable != nil && actionsTable.Type() == lua.LTTable {
		valid := true
		actionsTable.(*lua.LTable).ForEach(func(idx, f lua.LValue) {
			//
			// golang function in lua is '*lua.LFunction' type
			//
			if !(f.Type() == lua.LTFunction) {
				valid = false
			}
		})
		if !valid {
			return errors.New("Invalid function type")
		}
	} else {
		return errors.New("'Actions' must be a functions table")
	}
	// 释放语法验证阶段的临时虚拟机
	tempVm.Close()
	tempVm = nil
	// 交给规则脚本
	r.LuaVM.DoString(r.Success)
	r.LuaVM.DoString(r.Actions)
	r.LuaVM.DoString(r.Failed)
	return nil
}

// 临时校验语法
func ValidateLuaSyntax(bytes []byte) error {
	// 把虚拟机参数全部设置为0是为了防止误操作产生垃圾数据
	tempVm := lua.NewState(lua.Options{
		SkipOpenLibs:     true,
		RegistrySize:     0,
		RegistryMaxSize:  0,
		RegistryGrowStep: 0,
	})
	if err := tempVm.DoString(string(bytes)); err != nil {
		return err
	}
	// 检查函数入口
	AppMain := tempVm.GetGlobal("Main")
	if AppMain == nil {
		return fmt.Errorf("'Main' field not exists")
	}
	if AppMain.Type() != lua.LTFunction {
		return fmt.Errorf("'Main' must be function(arg)")
	}
	tempVm.Close()
	tempVm = nil
	return nil
}

// 临时校验语法
func ValidateCecollaletSyntax(bytes []byte) error {
	return ValidateLuaSyntax(bytes)
}

// 临时校验语法
func ValidateAppletSyntax(bytes []byte) error {
	return ValidateLuaSyntax(bytes)
}

/*
*
  - 分组加入函数
*/
func AddRuleLibToGroup(rx typex.Rhilex, LState *lua.LState,
	ModuleName string, funcs map[string]func(*lua.LState) int) {
	var table *lua.LTable
	if ModuleName == "_G" {
		table = LState.G.Global
	} else {
		table = LState.NewTable()
	}
	LState.SetGlobal(ModuleName, table)
	for funcName, f := range funcs {
		table.RawSetString(funcName, LState.NewClosure(f))
	}
	LState.Push(table)
}

func LoadRuleLibGroup(e typex.Rhilex, scope, uuid string, LState *lua.LState) {
	{
		Funcs := map[string]func(l *lua.LState) int{
			"ToHttp":       rhilexlib.DataToHttp(e, uuid),
			"ToMqtt":       rhilexlib.DataToMqtt(e, uuid),
			"ToUdp":        rhilexlib.DataToUdp(e, uuid),
			"ToTcp":        rhilexlib.DataToTcp(e, uuid),
			"ToTdEngine":   rhilexlib.DataToTdEngine(e, uuid),
			"ToMongoDB":    rhilexlib.DataToMongoDB(e, uuid),
			"ToSemtechUdp": rhilexlib.DataToSemtechUdp(e, uuid),
			"ToUart":       rhilexlib.DataToUart(e, uuid),
			"ToGreptimeDB": rhilexlib.DataToGreptimeDB(e),
		}
		AddRuleLibToGroup(e, LState, "data", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{}
		Funcs["Throw"] = rhilexlib.Throw(e, uuid)
		if scope == "RULE" {
			Funcs["Debug"] = rhilexlib.DebugRule(e, uuid)
		}
		if scope == "CECOLLA" {
			Funcs["Debug"] = rhilexlib.DebugCecolla(e, uuid)
		}
		if scope == "APPLET" {
			Funcs["Debug"] = rhilexlib.DebugAPP(e, uuid)
		}
		AddRuleLibToGroup(e, LState, "_G", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"VSet":             rhilexlib.StoreSet(e, uuid),
			"VSetWithDuration": rhilexlib.StoreSetWithDuration(e, uuid),
			"VGet":             rhilexlib.StoreGet(e, uuid),
			"VDel":             rhilexlib.StoreDelete(e, uuid),
		}
		AddRuleLibToGroup(e, LState, "kv", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"Time":       rhilexlib.Time(e, uuid),
			"TimeMs":     rhilexlib.TimeMs(e, uuid),
			"TsUnix":     rhilexlib.TsUnix(e, uuid),
			"TsUnixNano": rhilexlib.TsUnixNano(e, uuid),
			"NtpTime":    rhilexlib.NtpTime(e, uuid),
			"Sleep":      rhilexlib.Sleep(e, uuid),
		}
		AddRuleLibToGroup(e, LState, "time", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"HToN":         rhilexlib.HToN(e, uuid),
			"HsubToN":      rhilexlib.HsubToN(e, uuid),
			"MatchHex":     rhilexlib.MatchHex(e, uuid),
			"MatchUInt":    rhilexlib.MatchUInt(e, uuid),
			"Bytes2Hexs":   rhilexlib.Bytes2Hexs(e, uuid),
			"Hexs2Bytes":   rhilexlib.Hexs2Bytes(e, uuid),
			"ABCD":         rhilexlib.ABCD(e, uuid),
			"DCBA":         rhilexlib.DCBA(e, uuid),
			"BADC":         rhilexlib.BADC(e, uuid),
			"CDAB":         rhilexlib.CDAB(e, uuid),
			"TwoBytesHOrL": rhilexlib.TwoBytesHOrL(e, uuid),
			"Int16HOrL":    rhilexlib.Int16HOrL(e, uuid),
		}
		AddRuleLibToGroup(e, LState, "hex", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"MB":            rhilexlib.MatchBinary(e, uuid),
			"MBHex":         rhilexlib.MatchBinaryHex(e, uuid),
			"B2BS":          rhilexlib.ByteToBitString(e, uuid),
			"Bit":           rhilexlib.GetABitOnByte(e, uuid),
			"B2I64":         rhilexlib.ByteToInt64(e, uuid),
			"B64S2B":        rhilexlib.B64S2B(e, uuid),
			"BS2B":          rhilexlib.BitStringToBytes(e, uuid),
			"Bin2F32":       rhilexlib.BinToFloat32(e, uuid),
			"Bin2F64":       rhilexlib.BinToFloat64(e, uuid),
			"Bin2F32Big":    rhilexlib.BinToFloat32(e, uuid),
			"Bin2F64Big":    rhilexlib.BinToFloat64(e, uuid),
			"Bin2F32Little": rhilexlib.BinToFloat32Little(e, uuid),
			"Bin2F64Little": rhilexlib.BinToFloat64Little(e, uuid),
		}
		AddRuleLibToGroup(e, LState, "binary", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"T2J": rhilexlib.JSONE(e, uuid),
			"J2T": rhilexlib.JSOND(e, uuid),
		}
		AddRuleLibToGroup(e, LState, "json", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"ReadDevice":  rhilexlib.ReadDevice(e, uuid),
			"WriteDevice": rhilexlib.WriteDevice(e, uuid),
			"CtrlDevice":  rhilexlib.CtrlDevice(e, uuid),
		}
		AddRuleLibToGroup(e, LState, "device", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"T2Str":   rhilexlib.T2Str(e, uuid),
			"Bin2Str": rhilexlib.Bin2Str(e, uuid),
		}
		AddRuleLibToGroup(e, LState, "string", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"F5":         rhilexlib.F5(e, uuid),
			"F6":         rhilexlib.F6(e, uuid),
			"F15":        rhilexlib.F15(e, uuid),
			"F16":        rhilexlib.F16(e, uuid),
			"WritePoint": rhilexlib.WriteToModbusSheetRegisterWithTag(e),
			"ParseByte":  rhilexlib.ParseModbusByte(e, uuid),
		}
		AddRuleLibToGroup(e, LState, "modbus", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"DO1Set":  rhilexlib.RHILEXG1_DO1Set(e, uuid),
			"DO1Get":  rhilexlib.RHILEXG1_DO1Get(e, uuid),
			"DO2Set":  rhilexlib.RHILEXG1_DO2Set(e, uuid),
			"DO2Get":  rhilexlib.RHILEXG1_DO2Get(e, uuid),
			"DI1Get":  rhilexlib.RHILEXG1_DI1Get(e, uuid),
			"DI2Get":  rhilexlib.RHILEXG1_DI2Get(e, uuid),
			"DI3Get":  rhilexlib.RHILEXG1_DI3Get(e, uuid),
			"Led1On":  rhilexlib.Led1On(e, uuid),
			"Led1Off": rhilexlib.Led1Off(e, uuid),
		}
		AddRuleLibToGroup(e, LState, "rhilexg1", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"XOR":   rhilexlib.XOR(e, uuid),
			"CRC16": rhilexlib.CRC16(e, uuid),
		}
		AddRuleLibToGroup(e, LState, "misc", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"TFloat":    rhilexlib.TruncateFloat(e, uuid),
			"RandomInt": rhilexlib.RandomInt(e, uuid),
		}
		AddRuleLibToGroup(e, LState, "math", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"PlayMusic": rhilexlib.PlayMusic(e, uuid),
		}
		AddRuleLibToGroup(e, LState, "audio", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"Request": rhilexlib.Request(e),
		}
		AddRuleLibToGroup(e, LState, "rpc", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"Execute": rhilexlib.JqSelect(e, uuid),
		}
		AddRuleLibToGroup(e, LState, "jq", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"Ping": rhilexlib.PingIp(e, uuid),
		}
		AddRuleLibToGroup(e, LState, "network", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"Get":  rhilexlib.HttpGet(e, uuid),
			"Post": rhilexlib.HttpPost(e, uuid),
		}
		AddRuleLibToGroup(e, LState, "http", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"Save":       rhilexlib.InsertToDataCenterTable(e, uuid),
			"List":       rhilexlib.QueryDataCenterList(e, uuid),
			"Last":       rhilexlib.QueryDataCenterLast(e, uuid),
			"UpdateLast": rhilexlib.UpdateDataCenterLast(e, uuid),
		}
		AddRuleLibToGroup(e, LState, "rds", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"ParseDOxygen": rhilexlib.ApureParseOxygen(e),
		}
		AddRuleLibToGroup(e, LState, "apure", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"F5": rhilexlib.SlaverF5(e),
			"F6": rhilexlib.SlaverF6(e),
		}
		AddRuleLibToGroup(e, LState, "modbus_slaver", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"CtrlReplySuccess":        rhilexlib.IthingsCtrlReplySuccess(e),
			"CtrlReplyFailure":        rhilexlib.IthingsCtrlReplyFailure(e),
			"ActionReplySuccess":      rhilexlib.IthingsActionReplySuccess(e),
			"ActionReplyFailure":      rhilexlib.IthingsActionReplyFailure(e),
			"PropertyReplySuccess":    rhilexlib.IthingsPropertyReplySuccess(e),
			"PropertyReplyFailure":    rhilexlib.IthingsPropertyReplyFailure(e),
			"PropertyReport":          rhilexlib.IthingsPropertyReport(e),
			"GetProperties":           rhilexlib.IthingsGetProperties(e),
			"GetPropertyReplySuccess": rhilexlib.IthingsGetPropertyReplySuccess(e),
		}
		AddRuleLibToGroup(e, LState, "ithings", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"CtrlReplySuccess":     rhilexlib.TencentIothubCtrlReplySuccess(e),
			"CtrlReplyFailure":     rhilexlib.TencentIothubCtrlReplyFailure(e),
			"ActionReplySuccess":   rhilexlib.TencentIothubActionReplySuccess(e),
			"ActionReplyFailure":   rhilexlib.TencentIothubActionReplyFailure(e),
			"PropertyReplySuccess": rhilexlib.TencentIothubPropertyReplySuccess(e),
			"PropertyReplyFailure": rhilexlib.TencentIothubPropertyReplyFailure(e),
			"PropertyReport":       rhilexlib.TencentIothubPropertyReport(e),
			"GetPropertyReply":     rhilexlib.TencentIothubGetPropertyReply(e),
		}
		AddRuleLibToGroup(e, LState, "tciothub", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			// LED On
			"Led2On": rhilexlib.HAAS506_Led2On(e),
			"Led3On": rhilexlib.HAAS506_Led3On(e),
			"Led4On": rhilexlib.HAAS506_Led4On(e),
			"Led5On": rhilexlib.HAAS506_Led5On(e),
			// Led Off
			"Led2Off": rhilexlib.HAAS506_Led2Off(e),
			"Led3Off": rhilexlib.HAAS506_Led2Off(e),
			"Led4Off": rhilexlib.HAAS506_Led2Off(e),
			"Led5Off": rhilexlib.HAAS506_Led2Off(e),
			// DO On
			"DO1On": rhilexlib.HAAS506_DO1_On(e),
			"DO2On": rhilexlib.HAAS506_Do2_On(e),
			"DO3On": rhilexlib.HAAS506_Do3_On(e),
			"DO4On": rhilexlib.HAAS506_Do4_On(e),
			// DO Off
			"DO1Off": rhilexlib.HAAS506_DO1_Off(e),
			"DO2Off": rhilexlib.HAAS506_Do2_Off(e),
			"DO3Off": rhilexlib.HAAS506_Do3_Off(e),
			"DO4Off": rhilexlib.HAAS506_Do4_Off(e),
			// AI
			"GetAI1": rhilexlib.HAAS506_AI1_Get(e),
			"GetAI2": rhilexlib.HAAS506_AI2_Get(e),
			"GetAI3": rhilexlib.HAAS506_AI3_Get(e),
			"GetAI4": rhilexlib.HAAS506_AI4_Get(e),
			"GetAI5": rhilexlib.HAAS506_AI5_Get(e),
		}
		AddRuleLibToGroup(e, LState, "haas506ld1", Funcs)
	}
}

/*
*
* 加载外部扩展库
*
 */
func LoadExtLuaLib(e typex.Rhilex, LState *lua.LState) error {
	for _, s := range core.GlobalConfig.ExtLibs {
		err := LState.DoFile(s)
		if err != nil {
			return err
		}
	}
	return nil
}
