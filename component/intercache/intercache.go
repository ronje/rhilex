// Copyright (C) 2023 wwhai
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
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package intercache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hootrhino/rhilex/typex"
)

type CacheValue struct {
	UUID           string
	Status         int // 1 正常；0 错误，填充 ErrMsg
	ErrMsg         string
	LastFetchTime  uint64
	ExpirationTime time.Time // 缓存过期时间
	Value          interface{}
}

var __DefaultValueCache *GlobalValueRegistry

func RegisterSlot(Slot string) {
	__DefaultValueCache.RegisterSlot(Slot)
}
func UnRegisterSlot(Slot string) {
	__DefaultValueCache.UnRegisterSlot(Slot)
}
func GetSlot(Slot string) map[string]CacheValue {
	return __DefaultValueCache.GetSlot(Slot)
}
func SetValue(Slot, K string, V CacheValue) {
	__DefaultValueCache.SetValue(Slot, K, V)
}
func GetValue(Slot, K string) CacheValue {
	return __DefaultValueCache.GetValue(Slot, K)
}
func DeleteValue(Slot, K string) {
	__DefaultValueCache.DeleteValue(Slot, K)
}

func Size() uint64 {
	return __DefaultValueCache.Size()
}
func Flush() {
	__DefaultValueCache.Flush()
}

type GlobalValueRegistry struct {
	Slots      map[string]map[string]CacheValue
	slotKeys   []string
	ruleEngine typex.Rhilex
	locker     sync.RWMutex
}

func InitGlobalValueRegistry(ruleEngine typex.Rhilex) *GlobalValueRegistry {
	__DefaultValueCache = &GlobalValueRegistry{
		ruleEngine: ruleEngine,
		Slots:      map[string]map[string]CacheValue{},
		locker:     sync.RWMutex{},
		slotKeys:   []string{},
	}
	__DefaultValueCache.startTimeoutChecker()
	return __DefaultValueCache
}
func (M *GlobalValueRegistry) RegisterSlot(Slot string) {
	M.locker.Lock()
	defer M.locker.Unlock()
	M.Slots[Slot] = map[string]CacheValue{}
}
func (M *GlobalValueRegistry) GetSlot(Slot string) map[string]CacheValue {
	M.locker.Lock()
	defer M.locker.Unlock()
	if S, ok := M.Slots[Slot]; ok {
		return S
	}
	return nil
}
func (M *GlobalValueRegistry) SetValue(Slot, K string, V CacheValue) {
	M.locker.Lock()
	defer M.locker.Unlock()
	if S, ok := M.Slots[Slot]; ok {
		S[K] = V
		M.Slots[Slot] = S
	}
	found := false
	for i, k := range M.slotKeys {
		if k == K {
			M.slotKeys[i] = K
			found = true
			break
		}
	}
	if !found {
		M.slotKeys = append(M.slotKeys, K)
	}
}
func (M *GlobalValueRegistry) GetValue(Slot, K string) CacheValue {
	M.locker.Lock()
	defer M.locker.Unlock()
	if S, ok := M.Slots[Slot]; ok {
		return S[K]
	}
	return CacheValue{}
}
func (M *GlobalValueRegistry) DeleteValue(Slot, K string) {
	M.locker.Lock()
	defer M.locker.Unlock()
	if S, ok := M.Slots[Slot]; ok {
		delete(S, Slot)
	}
	for i, k := range M.slotKeys {
		if k == K {
			M.slotKeys = append(M.slotKeys[:i], M.slotKeys[i+1:]...)
			break
		}
	}
}
func (M *GlobalValueRegistry) List() []map[string]CacheValue {
	CacheValues := []map[string]CacheValue{}
	M.locker.RLock()
	defer M.locker.RUnlock()
	for _, k := range M.slotKeys {
		V := M.Slots[k]
		if V != nil {
			CacheValues = append(CacheValues, V)
		}
	}
	return CacheValues
}
func (M *GlobalValueRegistry) UnRegisterSlot(SlotName string) {
	M.locker.Lock()
	defer M.locker.Unlock()
	Slot := M.Slots[SlotName]
	for ks := range Slot {
		delete(Slot, ks)
	}
	delete(M.Slots, SlotName)
	for i := 0; i < len(M.slotKeys); i++ {
		M.slotKeys = append(M.slotKeys[:i], M.slotKeys[i+1:]...)
	}
}
func (M *GlobalValueRegistry) Size() uint64 {
	return uint64(len(M.Slots))
}

func (M *GlobalValueRegistry) ClearSlot() {
	M.locker.Lock()
	defer M.locker.Unlock()
	for _, sk := range M.slotKeys {
		for ks := range M.Slots[sk] {
			delete(M.Slots[sk], ks)
		}
		delete(M.Slots, sk)
	}
	for i := 0; i < len(M.slotKeys); i++ {
		M.slotKeys = append(M.slotKeys[:i], M.slotKeys[i+1:]...)
	}
}

func (M *GlobalValueRegistry) Flush() {
	for slotName, slot := range M.Slots {
		for k := range slot {
			delete(slot, k)
		}
		delete(M.Slots, slotName)
	}
}

// 设置带超时的缓存值
// 设置带超时的缓存值
func (M *GlobalValueRegistry) SetWithTimeout(Slot, K string, V CacheValue, timeoutMs int) error {
	M.locker.Lock()
	defer M.locker.Unlock()

	if timeoutMs > 0 {
		// 计算过期时间
		expirationTime := time.Now().Add(time.Duration(timeoutMs) * time.Millisecond)
		V.ExpirationTime = expirationTime
	} else {
		// 如果超时时间为0，不设置过期时间
		V.ExpirationTime = time.Time{} // 零值，表示没有过期时间
	}

	// 存储值
	if S, ok := M.Slots[Slot]; ok {
		S[K] = V
		M.Slots[Slot] = S
	} else {
		return fmt.Errorf("slot not found: %s", Slot)
	}

	return nil
}

// 启动全局超时检查器
func (M *GlobalValueRegistry) startTimeoutChecker() {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-context.Background().Done():
				return
			case <-ticker.C:
				M.cleanupExpiredItems()
			}
		}
	}()
}

// 清理过期的条目
func (M *GlobalValueRegistry) cleanupExpiredItems() {
	M.locker.Lock()
	defer M.locker.Unlock()

	// 遍历所有槽位，检查每个键值对的过期时间
	for slot, items := range M.Slots {
		for key, value := range items {
			// 只有设置了过期时间，并且过期时间已到，才执行删除操作
			if !value.ExpirationTime.IsZero() {
				// 如果当前时间已经超过了过期时间，则删除该条目
				if time.Now().After(value.ExpirationTime) {
					delete(M.Slots[slot], key)
				}
			}
		}
	}
}
