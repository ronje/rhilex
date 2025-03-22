// Copyright (C) 2024 wwhai
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

package orderedmap

import "sync"

type OrderedMap[K comparable, V any] struct {
	mu       sync.RWMutex
	keys     []K
	values   map[K]V
	keyIndex map[K]int
}

func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		keys:     []K{},
		values:   make(map[K]V),
		keyIndex: make(map[K]int),
	}
}

func (om *OrderedMap[K, V]) Set(key K, value V) {
	om.mu.Lock()
	defer om.mu.Unlock()

	if _, exists := om.values[key]; !exists {
		om.keyIndex[key] = len(om.keys)
		om.keys = append(om.keys, key)
	}
	om.values[key] = value
}

func (om *OrderedMap[K, V]) Get(key K) (V, bool) {
	om.mu.RLock()
	defer om.mu.RUnlock()

	value, exists := om.values[key]
	return value, exists
}

func (om *OrderedMap[K, V]) Delete(key K) {
	om.mu.Lock()
	defer om.mu.Unlock()

	if index, exists := om.keyIndex[key]; exists {
		delete(om.values, key)
		delete(om.keyIndex, key)

		// Remove the key from the keys slice
		copy(om.keys[index:], om.keys[index+1:])
		om.keys = om.keys[:len(om.keys)-1]

		// Update the keyIndex map for the keys after the deleted key
		for i := index; i < len(om.keys); i++ {
			om.keyIndex[om.keys[i]] = i
		}
	}
}

func (om *OrderedMap[K, V]) Keys() []K {
	om.mu.RLock()
	defer om.mu.RUnlock()

	return append([]K(nil), om.keys...)
}

func (om *OrderedMap[K, V]) Values() []V {
	om.mu.RLock()
	defer om.mu.RUnlock()

	values := make([]V, 0, len(om.values))
	for _, key := range om.keys {
		values = append(values, om.values[key])
	}
	return values
}

func (om *OrderedMap[K, V]) Size() int {
	om.mu.RLock()
	defer om.mu.RUnlock()

	return len(om.values)
}
