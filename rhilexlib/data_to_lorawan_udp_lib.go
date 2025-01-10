package rhilexlib

import (
	"errors"

	lua "github.com/hootrhino/gopher-lua"
	"github.com/hootrhino/rhilex/component/interqueue"
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"
)

/*
*
* Data To LoraWan UDP 1700
*
 */
func DataToSemtechUdp(rx typex.Rhilex, uuid string) func(*lua.LState) int {
	return func(l *lua.LState) int {
		id := l.ToString(2)
		data := l.ToString(3)
		err := handleLoraWanUDPFormat(rx, id, data)
		if err != nil {
			l.Push(lua.LString(err.Error()))
			return 1
		}
		l.Push(lua.LNil)
		return 1
	}
}

func handleLoraWanUDPFormat(e typex.Rhilex, uuid string, incoming string) error {
	outEnd := e.GetOutEnd(uuid)
	if outEnd != nil {
		return interqueue.PushOutQueue(outEnd, incoming)
	}
	msg := "target not found:" + uuid
	glogger.GLogger.Error(msg)
	return errors.New(msg)

}
