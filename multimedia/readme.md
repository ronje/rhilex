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

### 概述
`MultimediaResourceManager` 是一个用于管理多媒体资源的管理器，它封装了 `GatewayResourceManager` 并提供了一系列方法来初始化、加载、停止、重启多媒体资源，以及获取资源列表、详情和状态等功能。

### 类定义
```go
type MultimediaResourceManager struct {
    RuleEngine        typex.Rhilex
    MultimediaStreams *xmanager.GatewayResourceManager
}
```
- `RuleEngine`：类型为 `typex.Rhilex`，用于处理相关规则的引擎。
- `MultimediaStreams`：类型为 `*xmanager.GatewayResourceManager`，用于管理多媒体资源。

### 方法

#### 1. `InitMultimediaRuntime`
```go
func InitMultimediaRuntime(re typex.Rhilex)
```
- **功能描述**：初始化多媒体运行时，创建 `MultimediaResourceManager` 实例并注册 `__MultimediaBinding` 插槽。
- **参数**：
  - `re`：类型为 `typex.Rhilex`，用于初始化 `MultimediaResourceManager` 的规则引擎。
- **返回值**：无
- **示例**：
```go
var re typex.Rhilex // 假设已经正确初始化
multimedia.InitMultimediaRuntime(re)
```

#### 2. `StopMultimediaRuntime`
```go
func StopMultimediaRuntime()
```
- **功能描述**：停止所有多媒体资源，并取消注册 `__MultimediaBinding` 插槽。如果 `MultimediaResourceManager` 未初始化，则直接返回。
- **参数**：无
- **返回值**：无
- **示例**：
```go
multimedia.StopMultimediaRuntime()
```

#### 3. `LoadMultimediaResource`
```go
func LoadMultimediaResource(uuid string, name string, resourceType string, configMap map[string]interface{}, description string) error
```
- **功能描述**：加载多媒体资源。如果 `MultimediaResourceManager` 未初始化，则返回错误信息。
- **参数**：
  - `uuid`：字符串类型，资源的唯一标识符。
  - `name`：字符串类型，资源的名称。
  - `resourceType`：字符串类型，资源的类型，如 "RTSP"、"RTMP" 等。
  - `configMap`：`map[string]interface{}` 类型，资源的配置信息。
  - `description`：字符串类型，资源的描述信息。
- **返回值**：如果加载失败，返回错误信息；否则返回 `nil`。
- **示例**：
```go
err := multimedia.LoadMultimediaResource("uuid-1", "name", "RTSP", map[string]interface{}{}, "description")
if err != nil {
    panic(err)
}
```

#### 4. `RestartMultimediaResource`
```go
func RestartMultimediaResource(uuid string) error
```
- **功能描述**：重启指定 UUID 的多媒体资源。如果 `MultimediaResourceManager` 未初始化，则返回错误信息。
- **参数**：
  - `uuid`：字符串类型，资源的唯一标识符。
- **返回值**：如果重启失败，返回错误信息；否则返回 `nil`。
- **示例**：
```go
err := multimedia.RestartMultimediaResource("uuid-1")
if err != nil {
    panic(err)
}
```

#### 5. `StopMultimediaResource`
```go
func StopMultimediaResource(uuid string) error
```
- **功能描述**：停止指定 UUID 的多媒体资源。如果 `MultimediaResourceManager` 未初始化，则返回错误信息。
- **参数**：
  - `uuid`：字符串类型，资源的唯一标识符。
- **返回值**：如果停止失败，返回错误信息；否则返回 `nil`。
- **示例**：
```go
err := multimedia.StopMultimediaResource("uuid-1")
if err != nil {
    panic(err)
}
```

#### 6. `GetMultimediaResourceList`
```go
func GetMultimediaResourceList() []*xmanager.GatewayResourceWorker
```
- **功能描述**：获取所有多媒体资源的列表。如果 `MultimediaResourceManager` 未初始化，则返回 `nil`。
- **参数**：无
- **返回值**：返回多媒体资源的列表，类型为 `[]*xmanager.GatewayResourceWorker`。
- **示例**：
```go
resources := multimedia.GetMultimediaResourceList()
fmt.Println(resources)
```

#### 7. `GetMultimediaResourceDetails`
```go
func GetMultimediaResourceDetails(uuid string) (*xmanager.GatewayResourceWorker, error)
```
- **功能描述**：获取指定 UUID 的多媒体资源的详细信息。如果 `MultimediaResourceManager` 未初始化，则返回错误信息。
- **参数**：
  - `uuid`：字符串类型，资源的唯一标识符。
- **返回值**：如果资源存在，返回资源的详细信息，类型为 `*xmanager.GatewayResourceWorker`；否则返回错误信息。
- **示例**：
```go
details, err := multimedia.GetMultimediaResourceDetails("uuid-1")
if err != nil {
    panic(err)
}
fmt.Println(details)
```

#### 8. `GetMultimediaResourceStatus`
```go
func GetMultimediaResourceStatus(uuid string) (xmanager.GatewayResourceState, error)
```
- **功能描述**：获取指定 UUID 的多媒体资源的当前状态。如果 `MultimediaResourceManager` 未初始化，则返回错误信息。
- **参数**：
  - `uuid`：字符串类型，资源的唯一标识符。
- **返回值**：如果资源存在，返回资源的状态，类型为 `xmanager.GatewayResourceState`；否则返回错误信息。
- **示例**：
```go
status, err := multimedia.GetMultimediaResourceStatus("uuid-1")
if err != nil {
    panic(err)
}
fmt.Println(status)
```

#### 9. `StartMultimediaResourceMonitoring`
```go
func StartMultimediaResourceMonitoring()
```
- **功能描述**：开始监控所有多媒体资源的状态。如果 `MultimediaResourceManager` 未初始化，则直接返回。
- **参数**：无
- **返回值**：无
- **示例**：
```go
multimedia.StartMultimediaResourceMonitoring()
```

#### 10. `RegisterMultimediaResourceType`
```go
func RegisterMultimediaResourceType(resourceType string, factory func(map[string]interface{}) (xmanager.GatewayResource, error))
```
- **功能描述**：注册多媒体资源类型和对应的工厂函数。如果 `MultimediaResourceManager` 未初始化，则直接返回。
- **参数**：
  - `resourceType`：字符串类型，资源的类型。
  - `factory`：函数类型，创建资源实例的工厂函数。
- **返回值**：无
- **示例**：
```go
factory := func(configMap map[string]interface{}) (xmanager.GatewayResource, error) {
    // 实现工厂函数
    return nil, nil
}
multimedia.RegisterMultimediaResourceType("RTSP", factory)
```