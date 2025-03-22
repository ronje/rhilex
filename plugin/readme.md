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

# 插件开发指南
## 一、引言
`DemoPlugin` 模板提供了插件开发的基本结构和必要接口，开发者只需根据实际功能需求，在模板基础上进行扩展和实现即可。

## 二、开发环境准备
1. 确保已经安装了 Go 语言开发环境，且版本满足项目要求。
2. 确认已引入 `github.com/hootrhino/rhilex/typex` 和 `gopkg.in/ini.v1` 这两个依赖包。若未引入，可使用以下命令进行获取：
```bash
go get github.com/hootrhino/rhilex/typex
go get gopkg.in/ini.v1
```

## 三、插件结构剖析
1. **结构体定义**：
```go
type DemoPlugin struct {
}
```
`DemoPlugin` 结构体是插件的核心载体，开发者可在此结构体中添加字段，用于存储插件运行时所需的各种数据和状态信息，如数据库连接对象、配置参数等。

2. **实例创建函数**：
```go
func NewDemoPlugin() *DemoPlugin {
	return &DemoPlugin{}
}
```
`NewDemoPlugin` 函数用于创建插件实例，开发者一般无需修改此函数的基本结构，若插件有特殊的初始化逻辑，可在函数内部添加相应代码。

3. **初始化方法 `Init`**：
```go
func (dm *DemoPlugin) Init(config *ini.Section) error {
	return nil
}
```
`Init` 方法在插件初始化阶段被调用。开发者应在此方法中实现从 `config` 参数（即 INI 配置文件的一个 section）中读取插件所需的配置信息，并进行相应的初始化操作，如初始化数据库连接、设置内部参数等。若初始化过程中出现错误，需返回具体的错误信息；若成功，则返回 `nil`。

4. **启动方法 `Start`**：
```go
func (dm *DemoPlugin) Start(typex.Rhilex) error {
	return nil
}
```
`Start` 方法用于启动插件的核心功能。开发者需根据插件的实际功能，在此方法中编写启动相关服务、开启监听等逻辑代码。方法参数 `typex.Rhilex` 提供了框架相关的上下文信息，可按需使用。若启动过程中出现错误，返回对应的错误信息；启动成功则返回 `nil`。

5. **停止方法 `Stop`**：
```go
func (dm *DemoPlugin) Stop() error {
	return nil
}
```
`Stop` 方法在插件停止时被调用。开发者应在此方法中实现释放插件占用的资源，如关闭数据库连接、停止正在运行的服务等操作。若停止过程中出现错误，返回错误信息；停止成功则返回 `nil`。

6. **元信息方法 `PluginMetaInfo`**：
```go
func (dm *DemoPlugin) PluginMetaInfo() typex.XPluginMetaInfo {
	return typex.XPluginMetaInfo{
		UUID:        "DEMO01",
		Name:        "DemoPlugin",
		Version:     "v0.0.1",
		Description: "DEMO01",
	}
}
```
`PluginMetaInfo` 方法用于返回插件的元信息。开发者需根据插件的实际情况，修改 `UUID`（唯一标识符）、`Name`（名称）、`Version`（版本）和 `Description`（描述）字段的值，以便准确标识和描述插件。

7. **服务调用方法 `Service`**：
```go
func (dm *DemoPlugin) Service(arg typex.ServiceArg) typex.ServiceResult {
	return typex.ServiceResult{}
}
```
`Service` 方法是插件的服务调用接口。当外部对插件发起服务请求时，会调用此方法。开发者需根据具体的服务逻辑，解析 `arg` 参数（传入的服务调用参数），执行相应的操作，并将结果填充到 `typex.ServiceResult` 结构体中返回。

## 四、开发步骤
1. **创建插件结构体**：
复制 `DemoPlugin` 结构体，并将其重命名为符合插件功能的名称，如 `MyNewPlugin`。根据插件需求，在结构体中添加必要的字段。
```go
type MyNewPlugin struct {
    // 添加自定义字段，如数据库连接对象等
    dbConn interface{}
}
```
2. **实现实例创建函数**：
复制 `NewDemoPlugin` 函数，并将其重命名为 `NewMyNewPlugin`。在函数内部，根据插件结构体的字段，进行必要的初始化操作。
```go
func NewMyNewPlugin() *MyNewPlugin {
    return &MyNewPlugin{
        // 初始化自定义字段
        dbConn: nil,
    }
}
```
3. **编写初始化逻辑**：
在 `Init` 方法中，从 `config` 参数读取配置信息，根据配置进行插件的初始化工作，如建立数据库连接、读取文件路径等。
```go
func (mn *MyNewPlugin) Init(config *ini.Section) error {
    // 读取数据库连接配置
    dbConnStr := config.Key("db_connection_string").String()
    // 尝试建立数据库连接
    var err error
    mn.dbConn, err = establishDBConnection(dbConnStr)
    if err != nil {
        return err
    }
    return nil
}
```
4. **实现启动逻辑**：
在 `Start` 方法中，编写启动插件核心功能的代码，如启动服务、开启定时任务等。
```go
func (mn *MyNewPlugin) Start(rh typex.Rhilex) error {
    // 启动基于数据库连接的服务
    err := startMyService(mn.dbConn)
    if err != nil {
        return err
    }
    return nil
}
```
5. **编写停止逻辑**：
在 `Stop` 方法中，实现释放资源的逻辑，如关闭数据库连接、停止服务等。
```go
func (mn *MyNewPlugin) Stop() error {
    // 关闭数据库连接
    err := closeDBConnection(mn.dbConn)
    if err != nil {
        return err
    }
    return nil
}
```
6. **设置插件元信息**：
在 `PluginMetaInfo` 方法中，修改 `UUID`、`Name`、`Version` 和 `Description` 字段，使其符合插件的实际情况。
```go
func (mn *MyNewPlugin) PluginMetaInfo() typex.XPluginMetaInfo {
    return typex.XPluginMetaInfo{
        UUID:        "MY_PLUGIN_001",
        Name:        "My New Plugin",
        Version:     "v1.0.0",
        Description: "This is a new plugin for specific functionality",
    }
}
```
7. **实现服务调用逻辑**：
在 `Service` 方法中，根据传入的 `arg` 参数，执行相应的业务逻辑，并将结果填充到 `typex.ServiceResult` 结构体中返回。
```go
func (mn *MyNewPlugin) Service(arg typex.ServiceArg) typex.ServiceResult {
    // 解析服务调用参数
    reqData := parseServiceArg(arg)
    // 执行具体的业务逻辑
    resultData, err := performBusinessLogic(reqData, mn.dbConn)
    if err != nil {
        // 处理错误情况
        return typex.ServiceResult{
            Error: err.Error(),
        }
    }
    // 返回正常结果
    return typex.ServiceResult{
        Data: resultData,
    }
}
```