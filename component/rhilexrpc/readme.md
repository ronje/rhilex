## 文档
支持 GRPC 用户自定义协议接入。
### 协议定义
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