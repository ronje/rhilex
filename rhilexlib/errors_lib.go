package rhilexlib

import (
	"github.com/hootrhino/rhilex/glogger"
	"github.com/hootrhino/rhilex/typex"

	lua "github.com/hootrhino/gopher-lua"
)

func Throw(rx typex.Rhilex, uuid string) func(*lua.LState) int {
	return func(l *lua.LState) int {
		glogger.GLogger.Error(l.ToString(1))
		return 0
	}
}
