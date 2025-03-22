<!--
 Copyright (C) 2025 wwhai

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU Affero General Public License as
 published by the Free Software Foundation, either version 3 of the
 License, or (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU Affero General Public License for more details.

 You should have received a copy of the GNU Affero General Public License
 along with this program.  If not, see <https://www.gnu.org/licenses/>.
-->

# RHILEX规则引擎设计

## 1. 概述
RHILEX规则引擎是一个基于Go语言开发的强大规则处理系统，它利用Lua脚本的灵活性来定义规则表达式。该引擎的核心机制在于执行一系列Lua函数，这些函数按照特定的逻辑顺序依次处理输入数据，从而实现复杂的规则处理流程。

## 2. 规则定义
在RHILEX规则引擎中，规则是通过一个Lua函数列表来定义的，这个列表被命名为`Actions`。以下是一个典型的`Actions`列表示例：

```lua
Actions = {
    function (in_data)
        new_data = handle(in_data)
        return new_data, bool
    end,
    function (in_data)
        new_data = handle(in_data)
        return new_data, bool
    end,
    function (in_data)
        new_data = handle(in_data)
        return new_data, bool
    end
}
```

### 2.1 函数结构
`Actions`列表中的每个函数都具有相同的结构：
- **输入参数**：每个函数都接受一个输入参数`in_data`，该参数是前一个函数的输出数据（如果有的话），或者是最初的输入数据（对于第一个函数）。
- **处理逻辑**：在函数内部，通过调用`handle`函数对输入数据进行处理，生成新的数据`new_data`。`handle`函数是一个自定义的处理函数，具体的处理逻辑根据业务需求而定。
- **返回值**：每个函数都返回两个值，第一个值是处理后的数据`new_data`，第二个值是一个布尔值`bool`。

## 3. 规则执行机制
规则引擎会按照顺序依次执行`Actions`列表中的每个函数。具体的执行流程如下：

### 3.1 初始输入
首先，将初始的输入数据传递给`Actions`列表中的第一个函数作为`in_data`参数。

### 3.2 函数执行
引擎会依次执行`Actions`列表中的每个函数，对于每个函数：
- 调用该函数，并将当前的输入数据`in_data`传递给它。
- 函数执行完毕后，返回两个值：`new_data`和`bool`。

### 3.3 数据传递逻辑
根据函数返回的布尔值`bool`，决定是否将当前函数的输出数据`new_data`传递给下一个函数：
- 如果`bool`为`true`，则将当前函数的输出数据`new_data`作为下一个函数的输入参数`in_data`。
- 如果`bool`为`false`，则停止执行后续的函数，整个规则处理流程结束。

### 3.4 循环执行
如果当前函数的`bool`值为`true`，则继续执行下一个函数，直到`Actions`列表中的所有函数都执行完毕，或者遇到某个函数的`bool`值为`false`为止。

## 4. 代码示例
以下是一个简单的Go代码示例，展示了如何使用RHILEX规则引擎执行规则：

```go
package main

import (
    "github.com/yuin/gopher-lua"
)

// 执行规则引擎
func executeRules(actions []*lua.LFunction, input lua.LValue) (lua.LValue, bool) {
    currentInput := input
    for _, action := range actions {
        // 调用Lua函数
        result := action.Call([]lua.LValue{currentInput})
        if len(result) != 2 {
            // 处理错误情况
            return nil, false
        }
        newData := result[0]
        continueFlag := lua.LVAsBool(result[1])
        if!continueFlag {
            return newData, false
        }
        currentInput = newData
    }
    return currentInput, true
}

func main() {
    L := lua.NewState()
    defer L.Close()

    // 加载Lua脚本
    script := `
        function handle(data)
            return data * 2
        end

        Actions = {
            function (in_data)
                new_data = handle(in_data)
                return new_data, true
            end,
            function (in_data)
                new_data = handle(in_data)
                return new_data, true
            end,
            function (in_data)
                new_data = handle(in_data)
                return new_data, false
            end
        }
    `
    if err := L.DoString(script); err != nil {
        panic(err)
    }

    // 获取Actions列表
    actionsTable := L.GetGlobal("Actions")
    if actionsTable.Type() != lua.LTTable {
        panic("Actions is not a table")
    }

    var actions []*lua.LFunction
    actionsTable.ForEach(func(key, value lua.LValue) {
        if value.Type() == lua.LTFunction {
            actions = append(actions, value.(*lua.LFunction))
        }
    })

    // 初始输入数据
    input := lua.LNumber(10)

    // 执行规则引擎
    result, _ := executeRules(actions, input)
    println(result.String())
}
```

## 5. 总结
RHILEX规则引擎通过Lua脚本的灵活性和Go语言的高效性，提供了一种强大的规则处理机制。通过定义一系列的Lua函数，并根据函数的返回值来决定数据的传递逻辑，实现了复杂的规则处理流程。这种机制可以广泛应用于各种需要根据规则进行数据处理的场景，如业务规则引擎、数据验证、工作流管理等。