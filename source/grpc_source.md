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
## 协议格式
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