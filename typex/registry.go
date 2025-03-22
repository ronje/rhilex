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

package typex

import (
	"github.com/hootrhino/rhilex/component/orderedmap"
)

// KeyType 定义了 Registry 中键的类型约束
type KeyType interface {
	string | DeviceType | InEndType | TargetType
}

// ValueType 定义了 Registry 中值的类型约束
type ValueType interface {
	*XConfig | *XPlugin
}

// Registry 是一个泛型注册表，用于存储键值对
type Registry[K KeyType, V ValueType] struct {
	registry *orderedmap.OrderedMap[K, V]
}

// NewRegistry 创建一个新的 Registry 实例
func NewRegistry[K KeyType, V ValueType]() *Registry[K, V] {
	return &Registry[K, V]{
		registry: orderedmap.NewOrderedMap[K, V](),
	}
}

// Register 向注册表中添加一个键值对
func (r *Registry[K, V]) Register(key K, val V) {
	r.registry.Set(key, val)
}

// Get 从注册表中获取指定键的值，如果键不存在则返回 nil 和 false
func (r *Registry[K, V]) Get(key K) (V, bool) {
	return r.registry.Get(key)
}

// All 返回注册表中所有的值
func (r *Registry[K, V]) All() []V {
	return r.registry.Values()
}

// Remove 从注册表中删除指定键的键值对
func (r *Registry[K, V]) Remove(key K) {
	r.registry.Delete(key)
}

// Count 返回注册表中键值对的数量
func (r *Registry[K, V]) Count() int {
	return r.registry.Size()
}

// Keys 返回注册表中所有的键
func (r *Registry[K, V]) Keys() []K {
	return r.registry.Keys()
}

// Values 返回注册表中所有的值
func (r *Registry[K, V]) Values() []V {
	return r.registry.Values()
}

// Find 查找注册表中指定键的值，如果键存在则返回值和 true，否则返回 nil 和 false
func (r *Registry[K, V]) Find(key K) (V, bool) {
	return r.registry.Get(key)
}
