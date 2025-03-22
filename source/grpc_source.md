<!--
 Copyright (C) 2024 wwhai

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

# GRPC Server


# RhilexRpc gRPC 服务文档

## 一、服务概述
RhilexRpc 是一个基于 gRPC 协议的服务，主要用于处理外部应用的请求。它通过开启一个 gRPC 服务器，为客户端提供远程调用接口，实现分布式系统间的通信和数据交互。

## 二、服务配置

### 2.1 配置项说明
服务的配置信息通过 `GrpcConfig` 结构体进行管理，包含以下几个关键配置项：

- **Host**：
    - **类型**：字符串
    - **说明**：gRPC 服务器监听的主机地址，是服务运行所在的主机地址，用于客户端连接到服务器。此配置项为必填项。
    - **示例**："127.0.0.1"

- **Port**：
    - **类型**：整数
    - **说明**：gRPC 服务器监听的端口号，客户端通过该端口与服务器进行通信。端口号的取值范围必须在 1 到 65535 之间，且为必填项。
    - **示例**：2583

- **Type**：
    - **类型**：字符串
    - **说明**：服务类型，用于标识该服务的特定类型或用途，可根据实际需求进行配置，为可选项。
    - **示例**："example_type"

- **CacheOfflineData**：
    - **类型**：布尔指针
    - **说明**：用于指示是否缓存离线数据，根据业务需求决定是否启用离线数据缓存功能，为可选项。
    - **示例**：true 或 false

### 2.2 默认配置
如果在初始化服务时未提供特定的配置信息，将使用以下默认配置：
- **Host**："127.0.0.1"
- **Port**：2583

## 三、协议格式
```proto
syntax = "proto3";
option go_package = "./;rhilexrpc";
option java_multiple_files = false;
option java_package = "rhilexrpc";
option java_outer_classname = "RhilexRpc";

package rhilexrpc;

service RhilexRpc {
  rpc Request (RpcRequest) returns (RpcResponse) {}
}

message RpcRequest {
  bytes value = 1;
}

message RpcResponse {
  int32 code = 1;
  string message = 2;
  bytes data = 3;
}
```