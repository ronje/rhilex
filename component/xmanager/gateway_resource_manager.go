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

// gateway_resource_worker.go
package xmanager

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hootrhino/rhilex/component/orderedmap"
)

// GatewayResourceManager 通用资源管理器
type GatewayResourceManager struct {
	resources *orderedmap.OrderedMap[string, *GatewayResourceWorker]
	types     map[string]func(map[string]interface{}) (GatewayResource, error)
	mu        sync.RWMutex
}

// NewGatewayResourceManager 创建新的资源管理器
func NewGatewayResourceManager() *GatewayResourceManager {
	return &GatewayResourceManager{
		resources: orderedmap.NewOrderedMap[string, *GatewayResourceWorker](),
		types:     make(map[string]func(map[string]interface{}) (GatewayResource, error)),
	}
}

// RegisterType 注册资源类型和其对应的 worker 实现
func (m *GatewayResourceManager) RegisterType(resourceType string, factory func(map[string]interface{}) (GatewayResource, error)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.types[resourceType] = factory
}

// LoadResource 加载资源
func (m *GatewayResourceManager) LoadResource(uuid string, name string, resourceType string, configMap map[string]interface{}, description string) error {
	m.mu.RLock()
	factory, exists := m.types[resourceType]
	m.mu.RUnlock()
	if !exists {
		return fmt.Errorf("unsupported resource type: %s", resourceType)
	}

	worker, err := factory(configMap)
	if err != nil {
		return err
	}

	err = worker.Init(uuid, configMap)
	if err != nil {
		return err
	}

	grw := &GatewayResourceWorker{
		Worker:      worker,
		UUID:        uuid,
		Name:        name,
		Type:        resourceType,
		Config:      configMap,
		Description: description,
	}

	m.resources.Set(uuid, grw)
	return nil
}

// RestartResource 重启资源
func (m *GatewayResourceManager) RestartResource(uuid string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	worker, exists := m.resources.Get(uuid)
	if !exists {
		return fmt.Errorf("resource not found: %s", uuid)
	}
	worker.Worker.Stop()
	ctx := context.Background()
	return worker.Worker.Start(ctx)
}

// StopResource 停止资源
func (m *GatewayResourceManager) StopResource(uuid string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	worker, exists := m.resources.Get(uuid)
	if !exists {
		return fmt.Errorf("resource not found: %s", uuid)
	}
	worker.Worker.Stop()
	m.resources.Delete(uuid)
	return nil
}

// GetResourceList 获取资源列表
func (m *GatewayResourceManager) GetResourceList() []*GatewayResourceWorker {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.resources.Values()
}

// GetResourceDetails 获取资源详情
func (m *GatewayResourceManager) GetResourceDetails(uuid string) (*GatewayResourceWorker, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	worker, exists := m.resources.Get(uuid)
	if !exists {
		return nil, fmt.Errorf("resource not found: %s", uuid)
	}
	return worker, nil
}

// GetResourceStatus 获取资源状态
func (m *GatewayResourceManager) GetResourceStatus(uuid string) (GatewayResourceState, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	worker, exists := m.resources.Get(uuid)
	if !exists {
		return MEDIA_DOWN, fmt.Errorf("resource not found: %s", uuid)
	}
	return worker.Worker.Status(), nil
}

// StartMonitoring 开始资源监控
// StartMonitoring 开始资源监控
func (m *GatewayResourceManager) StartMonitoring() {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			m.mu.RLock()
			uuids := m.resources.Keys()
			workers := make([]*GatewayResourceWorker, len(uuids))
			for i, uuid := range uuids {
				worker, _ := m.resources.Get(uuid)
				workers[i] = worker
			}
			m.mu.RUnlock()

			for _, worker := range workers {
				status := worker.Worker.Status()

				switch status {
				case MEDIA_DOWN:
					m.RestartResource(worker.UUID)
				case MEDIA_STOP, MEDIA_DISABLE:
					m.StopResource(worker.UUID)
				case MEDIA_PENDING:
					continue
				}
			}
		}
	}()
}
