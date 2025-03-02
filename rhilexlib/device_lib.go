package rhilexlib

import (
	"encoding/hex"

	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"

	lua "github.com/hootrhino/gopher-lua"
)

/*
*
* return Hex String
*
 */
func CtrlDevice(rx typex.Rhilex, uuid string) func(*lua.LState) int {
	return func(l *lua.LState) int {
		// write(uuid,cmd,data)
		devUUID := l.ToString(2)
		cmd := l.ToString(3)
		data := l.ToString(4)
		Device := rx.GetDevice(devUUID)
		if Device != nil {
			if Device.Device.Status() == typex.SOURCE_UP {
				result, err := Device.Device.OnCtrl([]byte(cmd), []byte(data))
				//
				CtrlResponse := hex.EncodeToString(result)
				if err != nil {
					glogger.GLogger.Error(err)
					l.Push(lua.LNil)
					l.Push(lua.LString(err.Error()))
					return 2
				} else {
					l.Push(lua.LString(CtrlResponse))
					l.Push(lua.LNil)
					return 2
				}
			} else {
				l.Push(lua.LNil)
				l.Push(lua.LString("device down:" + devUUID))
				return 2
			}

		}
		l.Push(lua.LNil)
		l.Push(lua.LString("device not exists:" + devUUID))
		return 2
	}
}
