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

# protocol 文档

## 1. 包概述
`protocol` 包实现了一个通用的请求 - 响应协议处理模块，包含了从应用层到传输层的一系列功能实现，用于处理数据的编码、解码、发送和接收，同时提供了主从模式的协议处理。该包遵循 GNU Affero General Public License 协议。

## 2. 主要结构体及功能

### 2.1 `GenericAppLayer`
- **功能**：应用层的实现，负责处理应用帧的请求、写入、读取操作，同时统计错误包的数量。
- **字段**：
  - `errTxCount`：发送错误包计数器。
  - `errRxCount`：接收错误包计数器。
  - `transport`：指向 `Transport` 结构体的指针，用于底层数据传输。
- **方法**：
  - `NewGenericAppLayerAppLayer(config ExchangeConfig) *GenericAppLayer`：根据配置信息创建一个新的 `GenericAppLayer` 实例。
  - `Request(appFrame *ApplicationFrame) (*ApplicationFrame, error)`：发送请求帧并等待响应帧，返回响应帧或错误信息。
  - `Write(appFrame *ApplicationFrame) error`：对应用帧进行编码并通过传输层发送。
  - `Read() (*ApplicationFrame, error)`：从传输层读取数据并解码为应用帧。
  - `Status() error`：获取传输层的状态。
  - `Close() error`：关闭传输层连接。

### 2.2 `Transport`
- **功能**：传输层的实现，负责数据的读写操作，包括设置读写超时、添加包头包尾、数据解析等。
- **字段**：
  - `buffer`：数据缓冲区。
  - `writer`：数据写入器。
  - `reader`：数据读取器。
  - `readTimeout`：读取超时时间。
  - `writeTimeout`：写入超时时间。
  - `parser`：字节解析器，用于解析数据包。
  - `port`：通用端口，实现了 `io.ReadWriteCloser` 接口。
  - `logger`：日志记录器。
- **方法**：
  - `NewTransport(config ExchangeConfig) *Transport`：根据配置信息创建一个新的 `Transport` 实例。
  - `Write(data []byte) error`：将数据添加包头包尾后写入端口。
  - `Read() ([]byte, error)`：从端口读取数据并进行解析。
  - `Status() error`：检查端口状态。
  - `Close() error`：关闭端口。

### 2.3 `GenericProtocolMaster`
- **功能**：主模式的协议处理，通过 `GenericProtocolHandler` 进行请求和停止操作。
- **字段**：
  - `handler`：指向 `GenericProtocolHandler` 结构体的指针。
- **方法**：
  - `NewGenericProtocolMaster(config ExchangeConfig) *GenericProtocolMaster`：根据配置信息创建一个新的 `GenericProtocolMaster` 实例。
  - `Request(appFrame *ApplicationFrame) (*ApplicationFrame, error)`：发送请求并获取响应。
  - `Stop()`：停止主模式的协议处理。

### 2.4 `GenericProtocolSlaver`
- **功能**：从模式的协议处理，通过 `GenericProtocolHandler` 进行数据读取，并提供循环处理和停止操作。
- **字段**：
  - `handler`：指向 `GenericProtocolHandler` 结构体的指针。
  - `ctx`：上下文。
  - `cancel`：取消函数，用于停止操作。
- **方法**：
  - `NewGenericProtocolSlaver(ctx context.Context, cancel context.CancelFunc, config ExchangeConfig) *GenericProtocolSlaver`：根据上下文、取消函数和配置信息创建一个新的 `GenericProtocolSlaver` 实例。
  - `StartLoop(callback func(*ApplicationFrame, error))`：启动循环处理，不断读取数据并通过回调函数处理。
  - `Stop()`：停止从模式的协议处理。

### 2.5 `GenericProtocolHandler`
- **功能**：协议处理的核心，封装了 `GenericAppLayer` 的操作。
- **字段**：
  - `appLayer`：指向 `GenericAppLayer` 结构体的指针。
- **方法**：
  - `NewGenericProtocolHandler(config ExchangeConfig) *GenericProtocolHandler`：根据配置信息创建一个新的 `GenericProtocolHandler` 实例。
  - `Request(appFrame *ApplicationFrame) (*ApplicationFrame, error)`：发送请求并获取响应。
  - `Write(appFrame *ApplicationFrame) error`：写入应用帧。
  - `Read() (*ApplicationFrame, error)`：读取应用帧。
  - `Status() error`：获取状态。
  - `Close() error`：关闭连接。

## 3. 配置结构体

### 3.1 `ExchangeConfig`
- **功能**：存储协议处理所需的配置信息。
- **字段**：
  - `Port`：通用端口，实现了 `io.ReadWriteCloser` 接口。
  - `ReadTimeout`：读取超时时间。
  - `WriteTimeout`：写入超时时间。
  - `PacketEdger`：包头包尾信息。
  - `Logger`：日志记录器。
- **方法**：
  - `NewExchangeConfig() ExchangeConfig`：创建一个默认配置的 `ExchangeConfig` 实例。

### 3.2 `PacketEdger`
- **功能**：定义数据包的包头和包尾。
- **字段**：
  - `Head`：包头，长度为 2 字节的数组。
  - `Tail`：包尾，长度为 2 字节的数组。

## 4. 接口

### 4.1 `DataChecker`
- **功能**：数据检查器接口，用于检查数据的有效性。
- **方法**：
  - `CheckData(data []byte) error`：检查数据是否有效，返回错误信息。

### 4.2 `GenericPort`
- **功能**：通用端口接口，扩展了 `io.ReadWriteCloser` 接口，添加了设置读写超时的方法。
- **方法**：
  - `SetReadDeadline(t time.Time) error`：设置读取超时时间。
  - `SetWriteDeadline(t time.Time) error`：设置写入超时时间。

## 5. 测试
- **`TestGenericProtocolMaster`**：对 `GenericProtocolMaster` 进行单元测试，创建一个 `GenericProtocolMaster` 实例，发送请求并检查响应。

## 6. 使用示例
以下是一个简单的使用示例，展示如何创建一个 `GenericProtocolMaster` 实例并发送请求：

```go
package main

import (
    "log"
    "testing"

    "github.com/sirupsen/logrus"
    "your_package_path/protocol"
)

func main() {
    // 初始化日志
    Logger := logrus.StandardLogger()
    Logger.SetLevel(logrus.DebugLevel)

    // 配置信息
    config := protocol.ExchangeConfig{
        Port:         protocol.NewGenericReadWriteCloser(),
        ReadTimeout:  5000,
        WriteTimeout: 5000,
        PacketEdger: protocol.PacketEdger{
            Head: [2]byte{0xAB, 0xAB},
            Tail: [2]byte{0xBA, 0xBA},
        },
        Logger: Logger,
    }

    // 创建主模式协议处理实例
    master := protocol.NewGenericProtocolMaster(config)

    // 创建请求帧
    request := protocol.NewApplicationFrame([]byte{0xFF, 0xFF, 0xFF, 0xFF})
    log.Println("Request:", request.ToString())

    // 发送请求并获取响应
    response, err := master.Request(request)
    if err != nil {
        log.Fatal(err)
    } else {
        log.Println("Response:", response.ToString())
    }

    // 停止协议处理
    master.Stop()
}
```

## 7. 注意事项
- 该包依赖于 `github.com/sirupsen/logrus` 进行日志记录，使用时需要确保该依赖已经正确安装。
- 在使用 `GenericProtocolSlaver` 的 `StartLoop` 方法时，需要注意回调函数的实现，避免出现死循环或资源泄漏。
- 错误处理方面，各个方法可能会返回不同的错误信息，调用者需要根据具体情况进行处理。