# UDP 接入文档：Rhilex UDP 服务

## 1. 服务概述
Rhilex UDP 服务提供了一个基于 UDP 协议的服务，用于接收客户端发送的 UDP 数据包。服务启动后，会监听指定的地址和端口，当接收到客户端的数据包时，会对数据进行处理，并返回响应信息。

## 2. 配置信息

### 2.1 配置结构体
服务的配置信息通过 `RHILEXUdpConfig` 结构体进行管理，包含以下几个关键配置项：

- **Host**：
    - **类型**：字符串
    - **说明**：服务监听的主机地址，用于指定服务运行所在的主机地址，客户端将向该地址发送 UDP 数据包。此配置项为必填项。
    - **示例**："0.0.0.0"

- **Port**：
    - **类型**：整数
    - **说明**：服务监听的端口号，客户端需要将 UDP 数据包发送到该端口。此配置项为必填项。
    - **示例**：6200

- **MaxDataLength**：
    - **类型**：整数
    - **说明**：最大数据包长度，用于限制服务接收的 UDP 数据包的最大长度。
    - **示例**：1024

### 2.2 默认配置
如果在初始化服务时未提供特定的配置信息，将使用以下默认配置：
- **Host**："0.0.0.0"
- **Port**：6200
- **MaxDataLength**：1024
以下是接入 Rhilex UDP 服务并以 JSON 格式发送数据的简要说明：

## 接入步骤

1. **确定服务地址和端口**：
   明确 Rhilex UDP 服务监听的主机地址（`Host`）和端口号（`Port`），默认配置下，`Host` 为 `"0.0.0.0"`，`Port` 为 `6200`。

2. **构造 JSON 数据**：
   准备要发送的数据，并将其转换为 JSON 格式的字符串。

3. **发送 UDP 数据包**：
   将 JSON 格式的字符串转换为字节数组，使用编程语言提供的 UDP 套接字功能，将数据包发送到指定的服务地址和端口。

### 示例代码（Python）

```python
import socket
import json

# 服务端地址和端口
server_host = "服务端Host"
server_port = 服务端Port

# 构造 JSON 数据
data = {
    "key1": "value1",
    "key2": "value2"
}
json_data = json.dumps(data).encode('utf-8')

# 创建 UDP 套接字
sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

# 发送数据
sock.sendto(json_data, (server_host, server_port))

# 接收响应
data, addr = sock.recvfrom(1024)
print(f"Received from server: {data.decode()}")

# 关闭套接字
sock.close()
```

### 示例代码（Go）

```go
package main

import (
    "encoding/json"
    "fmt"
    "net"
)

func main() {
    // 服务端地址和端口
    serverAddr, err := net.ResolveUDPAddr("udp", "服务端Host:服务端Port")
    if err != nil {
        fmt.Println("ResolveUDPAddr error:", err)
        return
    }

    localAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:0")
    if err != nil {
        fmt.Println("ResolveUDPAddr error:", err)
        return
    }

    conn, err := net.DialUDP("udp", localAddr, serverAddr)
    if err != nil {
        fmt.Println("DialUDP error:", err)
        return
    }
    defer conn.Close()

    // 构造 JSON 数据
    data := map[string]string{
        "key1": "value1",
        "key2": "value2",
    }
    jsonData, err := json.Marshal(data)
    if err != nil {
        fmt.Println("JSON encoding error:", err)
        return
    }

    // 发送数据
    _, err = conn.Write(jsonData)
    if err != nil {
        fmt.Println("Write error:", err)
        return
    }

    buffer := make([]byte, 1024)
    n, _, err := conn.ReadFromUDP(buffer)
    if err != nil {
        fmt.Println("ReadFromUDP error:", err)
        return
    }
    fmt.Println("Received from server:", string(buffer[:n]))
}
```

### 注意事项
- 发送的 JSON 数据转换为字节数组后的长度不要超过服务端配置的 `MaxDataLength`（默认值为 1024），以免数据丢失或处理异常。
- 若服务端配置发生变化，需相应调整客户端代码中的服务地址和端口。