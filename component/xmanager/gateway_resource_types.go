// Copyright (C) 2025 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// gateway_resource_types.go
package xmanager

import (
	"context"
	"fmt"
)

// GatewayResourceState 资源状态类型
type GatewayResourceState int

// to string
func (s GatewayResourceState) String() string {
	switch s {
	case MEDIA_DOWN:
		return "DOWN"
	case MEDIA_UP:
		return "UP"
	case MEDIA_PAUSE:
		return "PAUSE"
	case MEDIA_STOP:
		return "STOP"
	case MEDIA_PENDING:
		return "PENDING"
	case MEDIA_DISABLE:
		return "DISABLE"
	default:
		return "UNKNOWN"
	}
}

const (
	// 故障
	MEDIA_DOWN GatewayResourceState = 0
	// 启用
	MEDIA_UP GatewayResourceState = 1
	// 暂停
	MEDIA_PAUSE GatewayResourceState = 2
	// 停止
	MEDIA_STOP GatewayResourceState = 3
	// 准备
	MEDIA_PENDING GatewayResourceState = 4
	// 禁用
	MEDIA_DISABLE GatewayResourceState = 5
)

// 资源服务参数
type ResourceServiceArg struct {
	UUID string
	Args []any
}

// To string
func (s *ResourceServiceArg) String() string {
	return fmt.Sprintf(" ResourceServiceArg UUID: %s, Args: %v", s.UUID, s.Args)
}

// 资源服务
type ResourceServiceRequest struct {
	Name   string               // 服务名称
	Method string               // 服务方法
	Args   []ResourceServiceArg // 服务参数
}

// ResourceServiceReturn 资源服务返回
type ResourceServiceResponse struct {
	Type   string
	Result any
	Error  error
}

// to string
func (s *ResourceServiceResponse) String() string {
	return fmt.Sprintf("ResourceServiceResponse Type: %s, Result: %v, Error: %v", s.Type, s.Result, s.Error)
}

// 资源服务
type ResourceService struct {
	Name        string                  // 服务名称
	Description string                  // 服务描述
	Method      string                  // 服务方法
	Args        []ResourceServiceArg    // 服务参数
	Response    ResourceServiceResponse // 服务返回
}

// to string
func (s *ResourceService) String() string {
	return fmt.Sprintf("ResourceService Name: %s, Description: %s, Method: %s, Args: %v, Response: %v",
		s.Name, s.Description, s.Method, s.Args, s.Response)
}

// GatewayResource 多媒体资源工作接口
type GatewayResource interface {
	Init(uuid string, configMap map[string]any) error
	Start(context.Context) error
	Status() GatewayResourceState
	Services() []ResourceService
	OnService(request ResourceServiceRequest) (ResourceServiceResponse, error)
	Details() *GatewayResourceWorker
	Stop()
}
