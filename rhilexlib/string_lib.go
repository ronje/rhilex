package rhilexlib

import (
	"strings"

	lua "github.com/hootrhino/gopher-lua"
)

/*
*
* Table 转成 String, {1,2,3,4,5} -> "12345"
*
 */
func T2Str(l *lua.LState) int {
	table := l.ToTable(1)
	args := []string{}
	table.ForEach(func(l1, value lua.LValue) {
		args = append(args, value.String())
	})
	r := strings.Join(args, "")
	l.Push(lua.LString(r))
	return 1
}

// {255,255,255} -> "\0xFF\0xFF\0xFF"
func Bin2Str(l *lua.LState) int {
	table := l.ToTable(1)
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

// StringToUpper 将字符串转换为大写
func StringToUpper(L *lua.LState) int {
	s := L.ToString(1)
	L.Push(lua.LString(strings.ToUpper(s)))
	return 1
}

// StringToLower 将字符串转换为小写
func StringToLower(L *lua.LState) int {
	s := L.ToString(1)
	L.Push(lua.LString(strings.ToLower(s)))
	return 1
}

// StringTrim 去除字符串两端的空白字符
func StringTrim(L *lua.LState) int {
	s := L.ToString(1)
	L.Push(lua.LString(strings.TrimSpace(s)))
	return 1
}

// StringTrimLeft 去除字符串左侧的空白字符
func StringTrimLeft(L *lua.LState) int {
	s := L.ToString(1)
	L.Push(lua.LString(strings.TrimLeft(s, " ")))
	return 1
}

// StringTrimRight 去除字符串右侧的空白字符
func StringTrimRight(L *lua.LState) int {
	s := L.ToString(1)
	L.Push(lua.LString(strings.TrimRight(s, " ")))
	return 1
}

// StringReplace 替换字符串中的子串
func StringReplace(L *lua.LState) int {
	s := L.ToString(1)
	old := L.ToString(2)
	new := L.ToString(3)
	L.Push(lua.LString(strings.Replace(s, old, new, -1)))
	return 1
}

// StringRepeat 重复字符串多次
func StringRepeat(L *lua.LState) int {
	s := L.ToString(1)
	count := L.ToInt(2)
	L.Push(lua.LString(strings.Repeat(s, count)))
	return 1
}

// StringContains 检查字符串是否包含子串
func StringContains(L *lua.LState) int {
	s := L.ToString(1)
	substr := L.ToString(2)
	L.Push(lua.LBool(strings.Contains(s, substr)))
	return 1
}

// StringHasPrefix 检查字符串是否以指定前缀开始
func StringHasPrefix(L *lua.LState) int {
	s := L.ToString(1)
	prefix := L.ToString(2)
	L.Push(lua.LBool(strings.HasPrefix(s, prefix)))
	return 1
}

// StringHasSuffix 检查字符串是否以指定后缀结束
func StringHasSuffix(L *lua.LState) int {
	s := L.ToString(1)
	suffix := L.ToString(2)
	L.Push(lua.LBool(strings.HasSuffix(s, suffix)))
	return 1
}
