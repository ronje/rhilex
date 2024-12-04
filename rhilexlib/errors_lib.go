package rhilexlib

import (
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"

	lua "github.com/hootrhino/gopher-lua"
)

/**
 * 抛出异常
 *
 */
func Throw(rx typex.Rhilex, uuid string) func(*lua.LState) int {
	return func(l *lua.LState) int {
		Debugger, Ok := l.GetStack(1)
		if Ok {
			LValue, _ := l.GetInfo("f", Debugger, lua.LNil)
			l.GetInfo("l", Debugger, lua.LNil)
			l.GetInfo("S", Debugger, lua.LNil)
			l.GetInfo("u", Debugger, lua.LNil)
			l.GetInfo("n", Debugger, lua.LNil)
			LFunction := LValue.(*lua.LFunction)
			LastCall := lua.DbgCall{
				Name: "_main",
			}
			if len(LFunction.Proto.DbgCalls) > 0 {
				LastCall = LFunction.Proto.DbgCalls[0]
			}
			glogger.Errorf("Function Name: [%s],"+
				"What: [%s], Source Line: [%d],"+
				" Last Call: [%s]",
				Debugger.Name, Debugger.What, Debugger.CurrentLine,
				LastCall.Name,
			)
		}
		return 0
	}
}
