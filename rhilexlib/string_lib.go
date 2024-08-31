package rhilexlib

import (
	"strings"

	"github.com/hootrhino/rhilex/typex"

	lua "github.com/hootrhino/gopher-lua"
)

/*
*
* Table 转成 String, {1,2,3,4,5} -> "12345"
*
 */
func T2Str(rx typex.Rhilex, uuid string) func(*lua.LState) int {
	return func(l *lua.LState) int {
		table := l.ToTable(2)
		args := []string{}
		table.ForEach(func(l1, value lua.LValue) {
			args = append(args, value.String())
		})
		r := strings.Join(args, "")
		l.Push(lua.LString(r))
		return 1
	}
}

// {255,255,255} -> "\0xFF\0xFF\0xFF"
func Bin2Str(rx typex.Rhilex, uuid string) func(*lua.LState) int {
	return func(l *lua.LState) int {
		table := l.ToTable(2)
		args := []byte{}
		table.ForEach(func(l1, value lua.LValue) {
			switch value.Type() {
			case lua.LTNumber:
				if lua.LVAsNumber(value) >= 0 && lua.LVAsNumber(value) <= 255 {
					args = append(args, byte(lua.LVAsNumber(value)))
				}
			default:
				return
			}
		})
		l.Push(lua.LString(string(args)))
		return 1
	}
}
