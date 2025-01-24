package interkv

import (
	"fmt"
	"strings"
	"time"

	"github.com/hootrhino/rhilex/glogger"
	cache "github.com/wwhai/tinycache"
)

// 自定义错误类型
var __errMaxStoreSizeReached = fmt.Errorf("max store size reached")

var GlobalStore *RhilexStore

func InitInterKVStore(maxSize int) {
	GlobalStore = NewRhilexStore(maxSize)
}

type RhilexStore struct {
	cache   *cache.Cache
	maxSize int
	index   map[string][]string
}

func NewRhilexStore(maxSize int) *RhilexStore {
	return &RhilexStore{
		cache:   cache.New(time.Duration(maxSize), 0),
		maxSize: maxSize,
		index:   make(map[string][]string),
	}
}

/*
*
* 设置过期时间
*
 */
func (rs *RhilexStore) SetWithDuration(k string, v string, d time.Duration) error {
	if rs.cache.ItemCount()+1 > rs.maxSize {
		glogger.GLogger.Error("Max store size reached:", rs.cache.ItemCount())
		return __errMaxStoreSizeReached
	}
	rs.cache.Set(k, v, d)
	rs.updateIndex(k)
	return nil
}

// 设置值
func (rs *RhilexStore) Set(k string, v string) error {
	if rs.cache.ItemCount()+1 > rs.maxSize {
		glogger.GLogger.Error("Max store size reached:", rs.cache.ItemCount())
		return __errMaxStoreSizeReached
	}
	rs.cache.Set(k, v, -1)
	rs.updateIndex(k)
	return nil
}

// 获取值
func (rs *RhilexStore) Get(k string) string {
	v, ok := rs.cache.Get(k)
	if ok {
		return v.(string)
	} else {
		return ""
	}
}

// 删除键值对
func (rs *RhilexStore) Delete(k string) error {
	rs.cache.Delete(k)
	rs.removeFromIndex(k)
	return nil
}

// 统计数量
func (rs *RhilexStore) Count() int {
	return rs.cache.ItemCount()
}

// 模糊查询匹配
// 支持: *AAA AAA* A*B
func (rs *RhilexStore) FuzzyGet(Key string) any {
	keys := rs.getKeysFromIndex(Key)
	for _, k := range keys {
		if v, ok := rs.cache.Get(k); ok {
			return v
		}
	}
	return ""
}

// 更新索引
func (rs *RhilexStore) updateIndex(k string) {
	patterns := generateIndexPatterns(k)
	for _, pattern := range patterns {
		rs.index[pattern] = append(rs.index[pattern], k)
	}
}

// 从索引中移除键
func (rs *RhilexStore) removeFromIndex(k string) {
	patterns := generateIndexPatterns(k)
	for _, pattern := range patterns {
		keys := rs.index[pattern]
		for i, key := range keys {
			if key == k {
				rs.index[pattern] = append(keys[:i], keys[i+1:]...)
				break
			}
		}
	}
}

// 根据查询模式从索引中获取键列表
func (rs *RhilexStore) getKeysFromIndex(pattern string) []string {
	return rs.index[pattern]
}

// 生成索引模式
func generateIndexPatterns(key string) []string {
	key = strings.ToLower(key)
	patterns := []string{key}
	for i := 1; i <= len(key); i++ {
		patterns = append(patterns, key[:i]+"*")
		patterns = append(patterns, "*"+key[i:])
	}
	return patterns
}
