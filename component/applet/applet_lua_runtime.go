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

package applet

import (
	"fmt"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rhilex/rhilexlib"
	"github.com/hootrhino/rhilex/typex"
)

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

/*
*
  - 分组加入函数
*/
func AddAppLibToGroup(app *Application, rx typex.Rhilex,
	ModuleName string, funcs map[string]func(l *lua.LState) int) {
	var table *lua.LTable
	if ModuleName == "_G" {
		table = app.vm.G.Global
	} else {
		table = app.vm.NewTable()
	}
	app.vm.SetGlobal(ModuleName, table)
	for funcName, f := range funcs {
		table.RawSetString(funcName, app.vm.NewClosure(f))
	}
	app.vm.Push(table)
}

func LoadAppLibGroup(app *Application, e typex.Rhilex) {
	{
		Funcs := map[string]func(l *lua.LState) int{
			"ToHttp":       rhilexlib.DataToHttp(e, app.UUID),
			"ToMqtt":       rhilexlib.DataToMqtt(e, app.UUID),
			"ToUdp":        rhilexlib.DataToUdp(e, app.UUID),
			"ToTcp":        rhilexlib.DataToTcp(e, app.UUID),
			"ToTdEngine":   rhilexlib.DataToTdEngine(e, app.UUID),
			"ToMongoDB":    rhilexlib.DataToMongoDB(e, app.UUID),
			"ToSemtechUdp": rhilexlib.DataToSemtechUdp(e, app.UUID),
			"ToUart":       rhilexlib.DataToUart(e, app.UUID),
			"ToGreptimeDB": rhilexlib.DataToGreptimeDB(e),
		}
		AddAppLibToGroup(app, e, "data", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"Debug": rhilexlib.DebugAPP(e, app.UUID),
			"Throw": rhilexlib.Throw(e, app.UUID),
		}
		AddAppLibToGroup(app, e, "_G", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"VSet":             rhilexlib.StoreSet(e, app.UUID),
			"VSetWithDuration": rhilexlib.StoreSetWithDuration(e, app.UUID),
			"VGet":             rhilexlib.StoreGet(e, app.UUID),
			"VDel":             rhilexlib.StoreDelete(e, app.UUID),
		}
		AddAppLibToGroup(app, e, "kv", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"Time":       rhilexlib.Time(e, app.UUID),
			"TimeMs":     rhilexlib.TimeMs(e, app.UUID),
			"TsUnix":     rhilexlib.TsUnix(e, app.UUID),
			"TsUnixNano": rhilexlib.TsUnixNano(e, app.UUID),
			"NtpTime":    rhilexlib.NtpTime(e, app.UUID),
			"Sleep":      rhilexlib.Sleep(e, app.UUID),
		}
		AddAppLibToGroup(app, e, "time", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"HToN":         rhilexlib.HToN(e, app.UUID),
			"HsubToN":      rhilexlib.HsubToN(e, app.UUID),
			"MatchHex":     rhilexlib.MatchHex(e, app.UUID),
			"MatchUInt":    rhilexlib.MatchUInt(e, app.UUID),
			"Bytes2Hexs":   rhilexlib.Bytes2Hexs(e, app.UUID),
			"Hexs2Bytes":   rhilexlib.Hexs2Bytes(e, app.UUID),
			"ABCD":         rhilexlib.ABCD(e, app.UUID),
			"DCBA":         rhilexlib.DCBA(e, app.UUID),
			"BADC":         rhilexlib.BADC(e, app.UUID),
			"CDAB":         rhilexlib.CDAB(e, app.UUID),
			"TwoBytesHOrL": rhilexlib.TwoBytesHOrL(e, app.UUID),
			"Int16HOrL":    rhilexlib.Int16HOrL(e, app.UUID),
		}
		AddAppLibToGroup(app, e, "hex", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"MB":            rhilexlib.MatchBinary(e, app.UUID),
			"MBHex":         rhilexlib.MatchBinaryHex(e, app.UUID),
			"B2BS":          rhilexlib.ByteToBitString(e, app.UUID),
			"Bit":           rhilexlib.GetABitOnByte(e, app.UUID),
			"B2I64":         rhilexlib.ByteToInt64(e, app.UUID),
			"B64S2B":        rhilexlib.B64S2B(e, app.UUID),
			"BS2B":          rhilexlib.BitStringToBytes(e, app.UUID),
			"Bin2F32":       rhilexlib.BinToFloat32(e, app.UUID),
			"Bin2F64":       rhilexlib.BinToFloat64(e, app.UUID),
			"Bin2F32Big":    rhilexlib.BinToFloat32(e, app.UUID),
			"Bin2F64Big":    rhilexlib.BinToFloat64(e, app.UUID),
			"Bin2F32Little": rhilexlib.BinToFloat32Little(e, app.UUID),
			"Bin2F64Little": rhilexlib.BinToFloat64Little(e, app.UUID),
		}
		AddAppLibToGroup(app, e, "binary", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"T2J": rhilexlib.JSONE(e, app.UUID),
			"J2T": rhilexlib.JSOND(e, app.UUID),
		}
		AddAppLibToGroup(app, e, "json", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"ReadDevice":  rhilexlib.ReadDevice(e, app.UUID),
			"WriteDevice": rhilexlib.WriteDevice(e, app.UUID),
			"CtrlDevice":  rhilexlib.CtrlDevice(e, app.UUID),
			"ReadSource":  rhilexlib.ReadSource(e, app.UUID),
			"WriteSource": rhilexlib.WriteSource(e, app.UUID),
		}
		AddAppLibToGroup(app, e, "device", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"T2Str":   rhilexlib.T2Str(e, app.UUID),
			"Bin2Str": rhilexlib.Bin2Str(e, app.UUID),
		}
		AddAppLibToGroup(app, e, "string", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"F5":        rhilexlib.F5(e, app.UUID),
			"F6":        rhilexlib.F6(e, app.UUID),
			"F15":       rhilexlib.F15(e, app.UUID),
			"F16":       rhilexlib.F16(e, app.UUID),
			"ParseByte": rhilexlib.ParseModbusByte(e, app.UUID),
		}
		AddAppLibToGroup(app, e, "modbus", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"DO1Set":  rhilexlib.RHILEXG1_DO1Set(e, app.UUID),
			"DO1Get":  rhilexlib.RHILEXG1_DO1Get(e, app.UUID),
			"DO2Set":  rhilexlib.RHILEXG1_DO2Set(e, app.UUID),
			"DO2Get":  rhilexlib.RHILEXG1_DO2Get(e, app.UUID),
			"DI1Get":  rhilexlib.RHILEXG1_DI1Get(e, app.UUID),
			"DI2Get":  rhilexlib.RHILEXG1_DI2Get(e, app.UUID),
			"DI3Get":  rhilexlib.RHILEXG1_DI3Get(e, app.UUID),
			"Led1On":  rhilexlib.Led1On(e, app.UUID),
			"Led1Off": rhilexlib.Led1Off(e, app.UUID),
		}
		AddAppLibToGroup(app, e, "rhilexg1", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"XOR":   rhilexlib.XOR(e, app.UUID),
			"CRC16": rhilexlib.CRC16(e, app.UUID),
		}
		AddAppLibToGroup(app, e, "misc", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"GPIOGet": rhilexlib.RASPI4_GPIOGet(e, app.UUID),
			"GPIOSet": rhilexlib.RASPI4_GPIOSet(e, app.UUID),
		}
		AddAppLibToGroup(app, e, "raspi4b", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"GPIOGet": rhilexlib.WKYWS1608_GPIOGet(e, app.UUID),
			"GPIOSet": rhilexlib.WKYWS1608_GPIOSet(e, app.UUID),
		}
		AddAppLibToGroup(app, e, "ws1608", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"TFloat":    rhilexlib.TruncateFloat(e, app.UUID),
			"RandomInt": rhilexlib.RandomInt(e, app.UUID),
		}
		AddAppLibToGroup(app, e, "math", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"PlayMusic": rhilexlib.PlayMusic(e, app.UUID),
		}
		AddAppLibToGroup(app, e, "audio", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"Request": rhilexlib.Request(e),
		}
		AddAppLibToGroup(app, e, "rpc", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"Execute": rhilexlib.JqSelect(e, app.UUID),
		}
		AddAppLibToGroup(app, e, "jq", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"Ping": rhilexlib.PingIp(e, app.UUID),
		}
		AddAppLibToGroup(app, e, "network", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"Get":  rhilexlib.HttpGet(e, app.UUID),
			"Post": rhilexlib.HttpPost(e, app.UUID),
		}
		AddAppLibToGroup(app, e, "http", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"LedOn":  rhilexlib.EN6400_LedOn(e, app.UUID),
			"LedOff": rhilexlib.EN6400_LedOff(e, app.UUID),
		}
		AddAppLibToGroup(app, e, "en6400", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"Save":       rhilexlib.InsertToDataCenterTable(e, app.UUID),
			"List":       rhilexlib.QueryDataCenterList(e, app.UUID),
			"Last":       rhilexlib.QueryDataCenterLast(e, app.UUID),
			"UpdateLast": rhilexlib.UpdateDataCenterLast(e, app.UUID),
		}
		AddAppLibToGroup(app, e, "rds", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"Ctrl": rhilexlib.CtrlComRF(e),
		}
		AddAppLibToGroup(app, e, "rfcom", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"ParseDOxygen": rhilexlib.ApureParseOxygen(e),
		}
		AddAppLibToGroup(app, e, "apure", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"F5": rhilexlib.SlaverF5(e),
			"F6": rhilexlib.SlaverF6(e),
		}
		AddAppLibToGroup(app, e, "modbus_slaver", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"ActionReplySuccess":   rhilexlib.IthingsActionReplySuccess(e),
			"ActionReplyFailure":   rhilexlib.IthingsActionReplyFailure(e),
			"PropertyReplySuccess": rhilexlib.IthingsPropertyReplySuccess(e),
			"PropertyReplyFailure": rhilexlib.IthingsPropertyReplyFailure(e),
			"PropertyReport":       rhilexlib.IthingsPropertyReport(e),
			"GetPropertyReply":     rhilexlib.IthingsGetPropertyReply(e),
		}
		AddAppLibToGroup(app, e, "ithings", Funcs)
	}
	{
		Funcs := map[string]func(l *lua.LState) int{
			"ActionReplySuccess":   rhilexlib.IthingsActionReplySuccess(e),
			"ActionReplyFailure":   rhilexlib.IthingsActionReplyFailure(e),
			"PropertyReplySuccess": rhilexlib.IthingsPropertyReplySuccess(e),
			"PropertyReplyFailure": rhilexlib.IthingsPropertyReplyFailure(e),
			"PropertyReport":       rhilexlib.IthingsPropertyReport(e),
			"GetPropertyReply":     rhilexlib.IthingsGetPropertyReply(e),
		}
		AddAppLibToGroup(app, e, "tciothub", Funcs)
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
		AddAppLibToGroup(app, e, "haas506", Funcs)
	}
}
