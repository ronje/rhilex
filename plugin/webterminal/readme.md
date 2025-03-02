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

# WebTerminal

## 一、简介
WebTerminal 是一个基于 Go 语言开发的 Web 终端应用，通过结合伪终端（`pty`）技术和 WebSocket 通信，实现了在 Web 浏览器中模拟终端操作的功能。它允许用户在网页上输入命令，并将命令执行结果实时反馈到浏览器中，就像在本地终端中操作一样。

## 二、技术架构
1. **后端**：
    - 使用 Go 语言编写，借助 `os/exec` 和 `pty` 包创建并管理伪终端，将终端的输入输出重定向到 WebSocket 连接。
    - 采用 `gorilla/websocket` 库处理 WebSocket 通信，实现与前端的数据传输。
    - 利用 `http` 包搭建 HTTP 服务器，监听指定端口，处理 WebSocket 连接请求。
2. **前端**：
    - 基于 `xterm.js` 库实现终端界面的渲染和交互，模拟真实终端的显示效果和输入行为。
    - 通过 WebSocket 与后端进行通信，将用户输入的命令发送到后端，并接收后端返回的命令执行结果进行显示。

## 三、主要功能
1. **终端模拟**：在 Web 页面中呈现一个类似于传统终端的界面，支持命令输入和结果显示。
2. **命令执行**：用户在终端界面输入的命令会被发送到后端，后端通过伪终端执行命令，并将执行结果实时返回给前端显示。
3. **WebSocket 通信**：使用 WebSocket 协议实现前后端的双向通信，确保数据的实时传输和交互的流畅性。
4. **并发控制**：通过 `busy` 标志和 `sync.Mutex` 实现对终端使用状态的并发控制，避免多个用户同时使用时产生冲突。
5. **Ping 机制**：后端定时向 WebSocket 连接发送 Ping 消息，以检测连接的有效性，确保连接的稳定性。
6. **优雅关闭**：支持通过上下文（`context`）实现服务器的优雅关闭，在关闭时正确处理资源释放和正在进行的连接。

## 四、代码结构
1. **结构体定义**：
`WebTerminal` 结构体是核心数据结构，包含了与终端相关的资源（如伪终端文件 `terminalPty`）、HTTP 服务器实例 `httpServer`、WebSocket 升级器 `upgrader` 以及并发控制相关的字段 `busy` 和 `mu`。
```go
type WebTerminal struct {
    terminalPty *os.File
    httpServer  *http.Server
    upgrader    websocket.Upgrader
    busy        bool
    mu          sync.Mutex
    ctx         context.Context
    cancel      context.CancelFunc
    wg          sync.WaitGroup
}
```
2. **方法实现**：
    - **初始化方法 `Init`**：目前为空实现，可用于读取配置文件等初始化操作。
    - **启动方法 `Start`**：创建伪终端，配置 WebSocket 升级器，启动 HTTP 服务器并监听指定端口，处理 WebSocket 连接请求。
    - **停止方法 `Stop`**：取消上下文，等待所有 goroutine 完成，关闭伪终端文件，优雅关闭 HTTP 服务器。
    - **重启方法 `Restart`**：先调用 `Stop` 方法停止服务，然后重新创建上下文并调用 `Start` 方法启动服务。
    - **插件元信息方法 `PluginMetaInfo`**：返回插件的元信息，包括 UUID、名称、版本和描述。
    - **服务调用方法 `Service`**：目前为空实现，可根据实际需求扩展为处理特定服务请求的逻辑。
    - **终端处理方法 `handleTerminal`**：处理 WebSocket 连接的建立，实现终端输入输出的双向数据传输，包括发送 Ping 消息、将伪终端输出重定向到 WebSocket、将 WebSocket 输入发送到伪终端等功能。

## 五、使用说明
1. **启动服务**：运行包含 `WebTerminal` 代码的 Go 程序，服务器将监听指定端口（默认为 `:2579`）。
2. **前端连接**：在前端页面中，使用 `xterm.js` 初始化终端界面，并通过 WebSocket 连接到后端服务器的 `/ws` 路径。
3. **操作终端**：在终端界面中输入命令，命令执行结果将实时显示在终端中。

## 六、注意事项
1. 确保在使用过程中正确处理资源的获取和释放，避免资源泄漏。
2. 由于使用了并发控制，要注意在多用户或多连接情况下的线程安全问题。
3. 在处理 WebSocket 连接时，要考虑网络异常等情况，确保连接的稳定性和可靠性。
4. 对于命令执行结果的处理，要注意处理可能的错误和异常情况，避免出现未预期的行为。