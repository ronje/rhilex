# 南向资源开发

## 1. 概述
外部南向资源主要用于与外部资源（如 MQTT 等）进行对接，实现双向通信功能。通过定义 `XSource` 接口，规范了与外部资源交互的一系列操作，包括资源的测试、初始化、启动、数据传输等，确保了与不同外部资源对接的一致性和可扩展性。

## 2. 接口定义

### 2.1 `XSource` 接口
`XSource` 接口是外部南向资源的核心接口，代表了一个终端资源（例如实际的 MQTT 客户端），定义了与资源交互所需的一系列方法。

### 2.2 接口方法

#### 2.2.1 `Test` 方法
```go
Test(inEndId string) bool
```
- **功能描述**：用于测试资源是否可用。
- **参数**：
  - `inEndId`：资源的标识符，用于唯一标识一个外部资源。
- **返回值**：
  - `bool`：如果资源可用则返回 `true`，否则返回 `false`。

#### 2.2.2 `Init` 方法
```go
Init(inEndId string, configMap map[string]any) error
```
- **功能描述**：用于初始化资源，传递资源配置信息。
- **参数**：
  - `inEndId`：资源的标识符。
  - `configMap`：资源配置的映射，包含了资源初始化所需的各种配置信息，例如连接地址、端口号、认证信息等。
- **返回值**：
  - `error`：如果初始化成功则返回 `nil`，否则返回相应的错误信息。

#### 2.2.3 `Start` 方法
```go
Start(CCTX context.Context) error
```
- **功能描述**：用于启动资源。
- **参数**：
  - `CCTX`：上下文，用于传递一些上下文信息，例如超时时间、取消信号等，具体作用取决于资源的实现。
- **返回值**：
  - `error`：如果启动成功则返回 `nil`，否则返回相应的错误信息。

#### 2.2.4 `DataModels` 方法
```go
DataModels() []XDataModel
```
- **功能描述**：用于获取资源支持的数据模型列表，这些模型对应于云平台的物模型。
- **参数**：无
- **返回值**：
  - `[]XDataModel`：资源支持的数据模型列表，`XDataModel` 是一个自定义的数据模型类型，包含了数据模型的相关信息。

#### 2.2.5 `Status` 方法
```go
Status() SourceState
```
- **功能描述**：用于获取资源的当前状态。
- **参数**：无
- **返回值**：
  - `SourceState`：资源的当前状态，`SourceState` 是一个自定义的枚举类型，用于表示资源的不同状态，例如 `STARTED`、`STOPPED`、`ERROR` 等。

#### 2.2.6 `Details` 方法
```go
Details() *InEnd
```
- **功能描述**：用于获取资源绑定的详细信息。
- **参数**：无
- **返回值**：
  - `*InEnd`：资源绑定的详细信息，`InEnd` 是一个自定义的结构体类型，包含了资源的详细信息，例如资源的名称、描述、配置信息等。

#### 2.2.7 `Stop` 方法
```go
Stop()
```
- **功能描述**：用于停止资源并释放相关资源。
- **参数**：无
- **返回值**：无

#### 2.2.8 `DownStream` 方法
```go
DownStream([]byte) (int, error)
```
- **功能描述**：用于处理下行数据，即从云平台发送到本地资源的数据。
- **参数**：
  - `[]byte`：从云平台发送过来的字节切片数据。
- **返回值**：
  - `int`：实际处理的数据长度。
  - `error`：如果处理成功则返回 `nil`，否则返回相应的错误信息。

#### 2.2.9 `UpStream` 方法
```go
UpStream([]byte) (int, error)
```
- **功能描述**：用于处理上行数据，即从本地资源发送到云平台的数据。
- **参数**：
  - `[]byte`：从本地资源发送的字节切片数据。
- **返回值**：
  - `int`：实际处理的数据长度。
  - `error`：如果处理成功则返回 `nil`，否则返回相应的错误信息。

## 3. 开发步骤

### 3.1 实现 `XSource` 接口
开发者需要创建一个具体的结构体类型，并实现 `XSource` 接口的所有方法。以下是一个简单的示例：

```go
type MyExternalResource struct {
    // 结构体成员，用于存储资源的相关信息
    config map[string]any
    state  SourceState
}

func (m *MyExternalResource) Test(inEndId string) bool {
    // 实现资源可用性测试逻辑
    // 例如，检查与外部资源的连接是否正常
    return true
}

func (m *MyExternalResource) Init(inEndId string, configMap map[string]any) error {
    // 实现资源初始化逻辑
    // 例如，根据 configMap 中的配置信息初始化连接
    m.config = configMap
    m.state = SourceState_INITIALIZED
    return nil
}

func (m *MyExternalResource) Start(CCTX context.Context) error {
    // 实现资源启动逻辑
    // 例如，建立与外部资源的连接
    m.state = SourceState_STARTED
    return nil
}

func (m *MyExternalResource) DataModels() []XDataModel {
    // 实现获取数据模型列表的逻辑
    // 例如，返回资源支持的数据模型列表
    return []XDataModel{}
}

func (m *MyExternalResource) Status() SourceState {
    // 实现获取资源状态的逻辑
    return m.state
}

func (m *MyExternalResource) Details() *InEnd {
    // 实现获取资源详细信息的逻辑
    // 例如，返回资源的名称、描述等信息
    return &InEnd{}
}

func (m *MyExternalResource) Stop() {
    // 实现资源停止逻辑
    // 例如，关闭与外部资源的连接
    m.state = SourceState_STOPPED
}

func (m *MyExternalResource) DownStream(data []byte) (int, error) {
    // 实现处理下行数据的逻辑
    // 例如，将数据发送到外部资源
    return len(data), nil
}

func (m *MyExternalResource) UpStream(data []byte) (int, error) {
    // 实现处理上行数据的逻辑
    // 例如，将数据发送到云平台
    return len(data), nil
}
```

### 3.2 使用 `XSource` 接口
在应用程序中，可以使用实现了 `XSource` 接口的结构体实例来进行资源的操作。以下是一个简单的示例：

```go
func main() {
    // 创建 MyExternalResource 实例
    resource := &MyExternalResource{}

    // 测试资源可用性
    if resource.Test("resource1") {
        fmt.Println("Resource is available")
    } else {
        fmt.Println("Resource is not available")
    }

    // 初始化资源
    configMap := map[string]any{
        "host": "localhost",
        "port": 1883,
    }
    err := resource.Init("resource1", configMap)
    if err != nil {
        fmt.Println("Failed to initialize resource:", err)
        return
    }

    // 启动资源
    ctx := context.Background()
    err = resource.Start(ctx)
    if err != nil {
        fmt.Println("Failed to start resource:", err)
        return
    }

    // 获取资源状态
    state := resource.Status()
    fmt.Println("Resource status:", state)

    // 处理下行数据
    data := []byte("Hello, external resource!")
    n, err := resource.DownStream(data)
    if err != nil {
        fmt.Println("Failed to send downstream data:", err)
    } else {
        fmt.Println("Sent downstream data:", n, "bytes")
    }

    // 处理上行数据
    data = []byte("Hello, cloud platform!")
    n, err = resource.UpStream(data)
    if err != nil {
        fmt.Println("Failed to send upstream data:", err)
    } else {
        fmt.Println("Sent upstream data:", n, "bytes")
    }

    // 停止资源
    resource.Stop()
}
```

## 4. 注意事项
- 在实现 `XSource` 接口的方法时，需要确保方法的实现逻辑符合接口的功能描述，并且处理好各种可能的错误情况。
- 资源的初始化和启动过程可能会涉及到一些异步操作，需要根据实际情况处理好上下文和超时等问题。
- 在处理下行和上行数据时，需要确保数据的格式和编码符合外部资源和云平台的要求。

## 5. 扩展和优化
- 可以根据实际需求，对 `XDataModel` 和 `InEnd` 等自定义类型进行扩展，添加更多的属性和方法，以满足不同的业务需求。
- 在处理数据传输时，可以考虑添加数据缓存、数据压缩、数据加密等功能，以提高数据传输的效率和安全性。
- 可以添加日志记录和监控功能，方便对资源的使用情况进行跟踪和分析。