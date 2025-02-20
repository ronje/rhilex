# 如何编写导出函数

在Go语言中，有时候我们需要编写导出函数，这些函数可以被Lua脚本调用。这里介绍一种常见的导出函数的编写方式，函数的导出形式为
```go
func Name(rx typex.Rhilex) func(*lua.LState) int {
    return func(L *lua.LState) int {
        ///...
        return 0
        }
    }
```

## 编写导出函数

### 函数结构
导出函数通常是一个高阶函数，它接受一些参数（这里是 `rx typex.Rhilex`），然后返回一个可以被Lua调用的函数。这个返回的函数接受一个 `*lua.LState` 类型的参数，返回一个 `int` 类型的值，表示返回给Lua的返回值数量。

### 获取参数
在返回的函数内部，我们可以通过 `L.ToString`、`L.ToNumber`、`L.ToBoolean` 等方法来获取Lua传递过来的参数。例如，`L.ToString(2)` 表示获取Lua调用时的第二个参数并将其转换为字符串。参数的索引从1开始。

### 返回参数
在返回的函数内部，我们可以使用 `L.Push` 方法将返回值压入Lua的栈中，然后返回栈中返回值的数量。例如，`L.Push(lua.LNumber(123))` 表示将一个Lua数字类型的值 `123` 压入栈中，然后 `return 1` 表示返回值的数量为1。

### 示例代码
以下是一个简单的示例，假设我们要编写一个导出函数，用于计算两个数的和：

```go
package main

import (
	"github.com/hootrhino/gopher-lua"
)

// Add 导出函数，用于计算两个数的和
func Add(rx interface{}) func(*lua.LState) int {
	return func(L *lua.LState) int {
		// 获取Lua传递过来的第一个参数并转换为数字
		a := L.ToNumber(1)
		// 获取Lua传递过来的第二个参数并转换为数字
		b := L.ToNumber(2)

		// 计算两个数的和
		result := a + b

		// 将结果压入Lua的栈中
		L.Push(lua.LNumber(result))

		// 返回栈中返回值的数量
		return 1
	}
}
```

## Lua调用示例
假设我们在Go代码中已经注册了 `Add` 函数，那么在Lua中可以这样调用：

```lua
-- 假设Go代码中已经注册了Add函数
local addFunc = Add(nil) -- 调用Go的导出函数
local result = addFunc(1, 2) -- 调用返回的函数，传递两个参数
print(result) -- 输出结果 3
```