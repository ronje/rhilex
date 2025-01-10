<!--
 Copyright (C) 2023 wwhai

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU Affero General Public License as
 published by the Free Software Foundation, either version 3 of the
 License, or (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU Affero General Public License for more details.

 You should have received a copy of the GNU Affero General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.
-->

## 内部组件
# RHILEX 系统组件开发文档
## 概述
RHILEX系统组件是构建在RHILEX框架之上的可插拔模块，它们为系统提供了可扩展性和灵活性。每个组件都有其独特的功能和职责，通过标准的接口与系统交互。本文档描述了如何定义和使用RHILEX系统组件。
## 包引用
```go
package component
```
## 导入依赖
```go
import (
	"github.com/hootrhino/rhilex/typex"
)
```
## 结构体定义
### XComponentMetaInfo
```go
type XComponentMetaInfo struct {
	UUID    string `json:"uuid"`    // UUID，组件的唯一标识符
	Name    string `json:"name"`    // 组件名，组件的名称
	Version string `json:"version"` // 版本，组件的版本号
}
```
### CallArgs
```go
type CallArgs struct {
	ComponentName string // 组件名称，调用时指定的组件名
	ServiceName   string // 服务名称，调用时指定的服务名
}
```
### CallResult
```go
type CallResult struct {
	Code   int    // 状态码，表示调用结果的状态
	Result any    // 结果，调用返回的结果数据
}
```
### ServiceSpec
```go
type ServiceSpec struct {
	CallArgs   CallArgs   // 调用参数，服务调用的入参
	CallResult CallResult // 调用结果，服务调用的出参
}
```
### XComponent
```go
type XComponent interface {
	Init(cfg map[string]any) error     // 初始化组件配置
	Start(rhilex typex.Rhilex) error   // 启动组件
	Call(CallArgs) (CallResult, error) // 调用组件接口
	Services() map[string]ServiceSpec  // 获取组件提供的服务列表
	MetaInfo() XComponentMetaInfo      // 获取组件的元信息
	Stop() error                       // 停止组件
}
```
## 接口说明
### Init
```go
Init(cfg map[string]any) error
```
- **描述**: 初始化组件配置。
- **参数**:
  - `cfg map[string]any`: 组件配置信息。
- **返回值**:
  - `error`: 如果初始化过程中发生错误，则返回错误信息。
### Start
```go
Start(rhilex typex.Rhilex) error
```
- **描述**: 启动组件。
- **参数**:
  - `rhilex typex.Rhilex`: RHILEX框架的实例。
- **返回值**:
  - `error`: 如果启动过程中发生错误，则返回错误信息。
### Call
```go
Call(CallArgs) (CallResult, error)
```
- **描述**: 调用组件接口。
- **参数**:
  - `CallArgs`: 调用参数，包含组件名称和服务名称。
- **返回值**:
  - `CallResult`: 调用结果，包含状态码和结果数据。
  - `error`: 如果调用过程中发生错误，则返回错误信息。
### Services
```go
Services() map[string]ServiceSpec
```
- **描述**: 获取组件提供的服务列表。
- **返回值**:
  - `map[string]ServiceSpec`: 包含组件所有服务的映射表。
### MetaInfo
```go
MetaInfo() XComponentMetaInfo
```
- **描述**: 获取组件的元信息。
- **返回值**:
  - `XComponentMetaInfo`: 组件的元信息，包括UUID、名称和版本。
### Stop
```go
Stop() error
```
- **描述**: 停止组件。
- **返回值**:
  - `error`: 如果停止过程中发生错误，则返回错误信息。
## 开发指南
1. **定义组件**: 创建一个新的结构体，实现`XComponent`接口。
2. **初始化配置**: 在`Init`方法中解析并验证配置参数。
3. **启动组件**: 在`Start`方法中启动组件所需的所有资源和服务。
4. **实现服务**: 定义组件提供的`ServiceSpec`，并在`Call`方法中处理服务调用。
5. **获取元信息**: 在`MetaInfo`方法中返回组件的元信息。
6. **停止组件**: 在`Stop`方法中释放组件占用的资源。
确保遵循上述指南和接口定义来开发RHILEX系统组件，以保证组件的兼容性和稳定性。
