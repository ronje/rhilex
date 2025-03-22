# RedisTarget

## 1. 概述
`RedisTarget` 是一个用于将数据存储到 Redis 的组件，它提供了初始化、启动、停止、状态检查以及数据存储等功能。该组件通过实现一系列方法，使得用户可以方便地将数据以哈希表（HMSET）的形式存储到 Redis 中，并且可以实时检查 Redis 连接的状态。

## 2. 结构体说明

### 2.1 `RedisTargetConfig`
```go
type RedisTargetConfig struct {
    Address  string `json:"address"`
    Password string `json:"password"`
    DB       int    `json:"db"`
}
```
- **功能**：用于存储 `RedisTarget` 的配置信息。
- **字段说明**：
  - `Address`：Redis 服务器的地址，格式为 `host:port`。
  - `Password`：连接 Redis 服务器所需的密码。
  - `DB`：要使用的 Redis 数据库编号。

### 2.2 `RedisTarget`
```go
type RedisTarget struct {
    typex.XStatus
    mainConfig RedisTargetConfig
    status     typex.SourceState
    client     *redis.Client
}
```
- **功能**：实现了将数据存储到 Redis 的目标组件。
- **字段说明**：
  - `typex.XStatus`：嵌入的状态相关结构体。
  - `mainConfig`：`RedisTargetConfig` 类型，存储 Redis 连接的配置信息。
  - `status`：`typex.SourceState` 类型，表示当前 Redis 连接的状态。
  - `client`：`*redis.Client` 类型，Redis 客户端实例，用于与 Redis 服务器进行通信。

### 2.3 `RedisData`
```go
type RedisData struct {
    Key   string
    Data  map[string]any
}
```
- **功能**：用于封装要存储到 Redis 的数据和键。
- **字段说明**：
  - `Key`：要存储数据的 Redis 键名。
  - `Data`：要存储的数据，以 `map[string]any` 的形式表示。

## 3. 方法说明

### 3.1 `NewRedisTarget`
```go
func NewRedisTarget(e typex.Rhilex) typex.XTarget
```
- **功能**：创建一个新的 `RedisTarget` 实例。
- **参数**：
  - `e`：`typex.Rhilex` 类型，可能是一个规则引擎实例。
- **返回值**：
  - `typex.XTarget` 类型，返回创建的 `RedisTarget` 实例。

### 3.2 `Init`
```go
func (rt *RedisTarget) Init(outEndId string, configMap map[string]any) error
```
- **功能**：初始化 `RedisTarget` 组件。
- **参数**：
  - `outEndId`：输出端点的 ID。
  - `configMap`：包含 Redis 配置信息的映射。
- **返回值**：
  - `error` 类型，如果初始化过程中出现错误，返回相应的错误信息；否则返回 `nil`。

### 3.3 `Start`
```go
func (rt *RedisTarget) Start(cctx typex.CCTX) error
```
- **功能**：启动 `RedisTarget` 组件。
- **参数**：
  - `cctx`：`typex.CCTX` 类型，包含上下文信息和取消函数。
- **返回值**：
  - `error` 类型，如果启动过程中出现错误，返回相应的错误信息；否则返回 `nil`。

### 3.4 `Status`
```go
func (rt *RedisTarget) Status() typex.SourceState
```
- **功能**：获取 `RedisTarget` 的当前状态，通过向 Redis 发送 `PING` 命令来检查连接状态。
- **返回值**：
  - `typex.SourceState` 类型，返回 Redis 连接的状态。

### 3.5 `To`
```go
func (rt *RedisTarget) To(data any) (any, error)
```
- **功能**：将数据存储到 Redis 中，使用 `HMSET` 命令。
- **参数**：
  - `data`：`any` 类型，需要是 `RedisData` 结构体的实例，包含要存储的数据和键。
- **返回值**：
  - `any` 类型，目前返回 `nil`。
  - `error` 类型，如果存储过程中出现错误，返回相应的错误信息；否则返回 `nil`。

### 3.6 `Stop`
```go
func (rt *RedisTarget) Stop()
```
- **功能**：停止 `RedisTarget` 组件，关闭 Redis 客户端连接。

### 3.7 `Details`
```go
func (rt *RedisTarget) Details() *typex.OutEnd
```
- **功能**：获取 `RedisTarget` 关联的输出端点的详细信息。
- **返回值**：
  - `*typex.OutEnd` 类型，返回输出端点的详细信息。

## 4. 使用示例
```go
package main

import (
    "context"
    "fmt"

    "github.com/hootrhino/rhilex/typex"
)

func main() {
    // 模拟Rhilex实例
    var mockRhilex typex.Rhilex

    // 创建RedisTarget实例
    redisTarget := NewRedisTarget(mockRhilex)

    // 初始化配置
    configMap := map[string]any{
        "address":  "localhost:6379",
        "password": "",
        "db":       0,
    }

    // 初始化RedisTarget
    err := redisTarget.Init("outEndId", configMap)
    if err != nil {
        panic(err)
    }

    // 启动RedisTarget
    ctx, cancel := context.WithCancel(context.Background())
    cctx := typex.CCTX{
        Ctx:        ctx,
        CancelCTX:  cancel,
    }
    err = redisTarget.Start(cctx)
    if err != nil {
        panic(err)
    }

    // 检查Redis连接状态
    status := redisTarget.Status()
    fmt.Println("Redis connection status:", status)

    // 准备数据
    data := map[string]any{
        "field1": "value1",
        "field2": "value2",
    }

    redisData := RedisData{
        Key:  "test_key",
        Data: data,
    }

    // 存储数据到Redis
    _, err = redisTarget.To(redisData)
    if err != nil {
        panic(err)
    }

    // 停止RedisTarget
    redisTarget.Stop()
}
```

## 5. 注意事项
- 在使用 `To` 方法存储数据时，传入的 `data` 参数必须是 `RedisData` 结构体的实例，否则会返回错误。
- 在调用 `Init` 方法初始化 `RedisTarget` 时，确保 `configMap` 中包含正确的 Redis 配置信息，否则可能会导致初始化失败。
- 在调用 `Stop` 方法停止 `RedisTarget` 时，会关闭 Redis 客户端连接，之后不能再使用该实例进行数据存储操作。