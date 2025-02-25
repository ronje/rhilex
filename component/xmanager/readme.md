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

# xmanager

## 概述
`xmanager` 包提供了一套用于管理网关资源的工具和接口，包括资源的创建、加载、启动、停止、监控等功能。

## 类型定义

### `GatewayResourceWorker`
用于记录流媒体的元信息。

```go
type GatewayResourceWorker struct {
    Worker      GatewayResource        // 实际的实现接口
    UUID        string                 // 资源唯一标识
    Name        string                 // 资源名称
    Type        string                 // 资源类型
    Config      map[string]any // 资源配置
    Description string                 // 资源描述
}
```

**方法**：
- `String() string`：返回资源的字符串表示形式。
- `GetConfig() map[string]any`：获取资源的配置信息。
- `CheckConfig(config any) error`：检查资源配置是否有效。

### `GatewayResourceManager`
通用资源管理器，用于管理多个 `GatewayResourceWorker` 实例。

```go
type GatewayResourceManager struct {
    resources *orderedmap.OrderedMap[string, *GatewayResourceWorker]
    types     map[string]func(map[string]any) (GatewayResource, error)
    mu        sync.RWMutex
}
```

**方法**：
- `NewGatewayResourceManager() *GatewayResourceManager`：创建新的资源管理器实例。
- `RegisterType(resourceType string, factory func(map[string]any) (GatewayResource, error))`：注册资源类型和其对应的 worker 实现。
- `LoadResource(uuid string, name string, resourceType string, configMap map[string]any, description string) error`：加载资源。
- `RestartResource(uuid string) error`：重启指定 UUID 的资源。
- `StopResource(uuid string) error`：停止并删除指定 UUID 的资源。
- `GetResourceList() []*GatewayResourceWorker`：获取所有资源的列表。
- `GetResourceDetails(uuid string) (*GatewayResourceWorker, error)`：获取指定 UUID 资源的详细信息。
- `GetResourceStatus(uuid string) (GatewayResourceState, error)`：获取指定 UUID 资源的状态。
- `StartMonitoring()`：开始资源监控，定期检查资源状态并执行相应操作。

### `GatewayResourceState`
资源状态类型，定义了资源的各种状态。

```go
type GatewayResourceState int

const (
    MEDIA_DOWN GatewayResourceState = 0
    MEDIA_UP GatewayResourceState = 1
    MEDIA_PAUSE GatewayResourceState = 2
    MEDIA_STOP GatewayResourceState = 3
    MEDIA_PENDING GatewayResourceState = 4
    MEDIA_DISABLE GatewayResourceState = 5
)
```

### `ResourceServiceArg`
资源服务参数。

```go
type ResourceServiceArg struct {
    UUID string
    Args []any
}
```

### `ResourceServiceRequest`
资源服务请求。

```go
type ResourceServiceRequest struct {
    Name   string               // 服务名称
    Method string               // 服务方法
    Args   []ResourceServiceArg // 服务参数
}
```

### `ResourceServiceResponse`
资源服务返回。

```go
type ResourceServiceResponse struct {
    Type   string
    Result any
    Error  error
}
```

### `ResourceService`
资源服务。

```go
type ResourceService struct {
    Name        string                  // 服务名称
    Description string                  // 服务描述
    Method      string                  // 服务方法
    Args        []ResourceServiceArg    // 服务参数
    Response    ResourceServiceResponse // 服务返回
}
```

### `GatewayResource`
多媒体资源工作接口，定义了资源的基本操作方法。

```go
type GatewayResource interface {
    Init(uuid string, configMap map[string]any) error
    Start(context.Context) error
    Status() GatewayResourceState
    Services() []ResourceService
    OnService(request ResourceServiceRequest) (ResourceServiceResponse, error)
    Details() *GatewayResourceWorker
    Stop()
}
```

## 使用示例

### 创建资源管理器
```go
manager := NewGatewayResourceManager()
```

### 注册资源类型
```go
manager.RegisterType("exampleType", func(configMap map[string]any) (GatewayResource, error) {
    // 实现资源的创建逻辑
    return nil, nil
})
```

### 加载资源
```go
configMap := map[string]any{
    "key": "value",
}
err := manager.LoadResource("uuid1", "name1", "exampleType", configMap, "description1")
if err != nil {
    // 处理错误
}
```

### 获取资源列表
```go
resourceList := manager.GetResourceList()
for _, resource := range resourceList {
    fmt.Println(resource.String())
}
```

## 注意事项
- 在使用 `GatewayResourceManager` 的方法时，需要注意并发安全，该管理器已经使用读写锁进行了并发控制。
- 当调用 `LoadResource` 方法时，确保所传入的 `resourceType` 已经通过 `RegisterType` 方法注册。
