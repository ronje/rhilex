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

# 内网端口穿透

## 一、概述
本系统是一个使用 Go 语言实现的内网端口穿透工具，具备客户端认证功能。它允许将内网中的服务暴露到公网上，使得外部用户可以通过访问服务端来访问内网服务。系统由服务端和客户端两部分组成，服务端负责管理客户端连接并进行数据转发，客户端负责连接服务端并将本地服务信息传递给服务端。

## 二、功能特性
1. **内网端口穿透**：实现内网服务端口到公网服务端的转发，让外部用户能够访问内网服务。
2. **客户端管理**：服务端可以管理连接的客户端，记录客户端信息，支持添加、移除客户端操作。
3. **认证功能**：客户端连接服务端时需要进行认证，采用简单的用户名和密码认证方式，确保系统安全性。

## 三、代码结构

### 服务端（`server.go`）
1. **`Client` 结构体**：表示一个连接的客户端，包含客户端连接、目标地址和存活状态。
2. **`ClientManager` 结构体**：管理所有连接的客户端，使用 `map` 存储客户端信息，并使用 `sync.Mutex` 进行并发控制。提供 `AddClient`、`RemoveClient` 和 `GetClient` 方法来管理客户端。
3. **认证逻辑**：
    - `AuthRequest` 和 `AuthResponse` 结构体用于处理认证请求和响应。
    - `authenticate` 函数负责读取客户端的认证请求，验证用户名和密码，并发送认证响应。
4. **数据转发**：`handleClient` 函数负责在客户端和目标地址之间进行数据转发。
5. **主函数**：监听服务端端口，接受客户端连接，进行认证，连接目标地址，添加客户端到管理列表，并启动数据转发协程。

### 客户端（`client.go`）
1. **认证逻辑**：
    - `AuthRequest` 和 `AuthResponse` 结构体用于处理认证请求和响应。
    - `authenticate` 函数负责构造认证请求并发送给服务端，读取服务端的认证响应，根据响应结果判断认证是否成功。
2. **数据转发**：`handleServer` 函数负责在服务端和本地服务之间进行数据转发。
3. **主函数**：连接到服务端，进行认证，发送本地地址到服务端，连接本地服务并启动数据转发。

## 四、使用方法

### 1. 编译代码
在服务端和客户端所在的机器上分别执行以下命令进行代码编译：
```sh
go build server.go
go build client.go
```

### 2. 运行服务端
在服务端机器上执行编译后的服务端程序：
```sh
./server
```
服务端将开始监听 `:8080` 端口，等待客户端连接。

### 3. 运行客户端
在客户端机器上执行编译后的客户端程序：
```sh
./client
```
客户端将连接到服务端，进行认证，认证成功后将本地服务信息发送给服务端，并开始数据转发。

### 4. 注意事项
- 请将 `client.go` 代码中的 `server_ip` 替换为实际的服务端 IP 地址。
- 目前代码中的用户名和密码是硬编码的，在实际应用中可以从配置文件或数据库中读取。

## 五、安全性考虑
- 本系统采用的是简单的用户名和密码认证方式，安全性较低。在实际应用中，建议使用更安全的认证机制，如 SSL/TLS 认证、OAuth 等。
- 服务端和客户端之间的数据传输没有进行加密，建议在生产环境中使用加密通道，如 SSH 隧道或 VPN 等。

## 六、扩展建议
1. **日志记录**：可以添加更详细的日志记录功能，方便调试和监控系统运行状态。
2. **配置管理**：将一些常量（如服务端端口、用户名、密码等）提取到配置文件中，方便修改和管理。
3. **错误处理**：进一步完善错误处理机制，提高系统的稳定性和可靠性。
4. **并发控制**：在高并发场景下，可以考虑使用更高效的并发控制机制，如 goroutine 池等。