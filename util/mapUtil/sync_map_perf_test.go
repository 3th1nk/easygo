package mapUtil

import (
	"fmt"
	"github.com/3th1nk/easygo/internal/_test"
	"github.com/3th1nk/easygo/util"
	"sync"
	"testing"
)

// 测试 map+RWMutex 和 sync.map 的 读取 性能差异
func TestSyncMap_Load_Perf(t *testing.T) {
	for _, size := range []int{10, 100, 1000, 10000, 100000} {
		doTestSyncMap_Perf(t, "Load Exists", size, func(key string, m map[interface{}]interface{}, mu *sync.RWMutex) bool {
			mu.RLock()
			defer mu.RUnlock()
			_, _ = m[key]
			return true
		}, func(key string, m *sync.Map) bool {
			_, _ = m.Load(key)
			return true
		})
	}

	for _, size := range []int{10, 100, 1000, 10000, 100000} {
		doTestSyncMap_Perf(t, "Load NotExists", size, func(key string, m map[interface{}]interface{}, mu *sync.RWMutex) bool {
			mu.RLock()
			defer mu.RUnlock()
			_, _ = m[key+key]
			return true
		}, func(key string, m *sync.Map) bool {
			_, _ = m.Load(key + key)
			return true
		})
	}
}

// 测试 map+RWMutex 和 sync.map 的 写入 性能差异
func TestSyncMap_Store_Perf(t *testing.T) {
	for _, size := range []int{10, 100, 1000, 10000, 100000} {
		doTestSyncMap_Perf(t, "Store", size, func(key string, m map[interface{}]interface{}, mu *sync.RWMutex) bool {
			mu.Lock()
			defer mu.Unlock()
			m[key] = key
			return true
		}, func(key string, m *sync.Map) bool {
			m.Store(key, key)
			return true
		})
	}
}

// 测试 map+RWMutex 和 sync.map 的 不存在时写入 性能差异
func TestSyncMap_LoadOrStore_Perf(t *testing.T) {
	for _, size := range []int{10, 100, 1000, 10000, 100000} {
		doTestSyncMap_Perf(t, "LoadOrStore Exists", size, func(key string, m map[interface{}]interface{}, mu *sync.RWMutex) bool {
			mu.Lock()
			defer mu.Unlock()
			if _, ok := m[key]; !ok {
				m[key] = key
			}
			return true
		}, func(key string, m *sync.Map) bool {
			m.LoadOrStore(key, key)
			return true
		})
	}

	for _, size := range []int{10, 100, 1000, 10000, 100000} {
		doTestSyncMap_Perf(t, "LoadOrStore NotExists", size, func(key string, m map[interface{}]interface{}, mu *sync.RWMutex) bool {
			mu.Lock()
			defer mu.Unlock()
			newKey := key + key
			if _, ok := m[newKey]; !ok {
				m[newKey] = newKey
			}
			return true
		}, func(key string, m *sync.Map) bool {
			newKey := key + key
			m.LoadOrStore(newKey, newKey)
			return true
		})
	}
}

// 测试 map+RWMutex 和 sync.map 的 遍历 性能差异
func TestSyncMap_Range_Perf(t *testing.T) {
	for _, size := range []int{10, 100, 1000, 10000, 100000} {
		doTestSyncMap_Perf(t, "Range", size, func(_ string, m map[interface{}]interface{}, mu *sync.RWMutex) bool {
			mu.RLock()
			defer mu.RUnlock()
			n := 0
			for _, _ = range m {
				n++
			}
			return false
		}, func(_ string, m *sync.Map) bool {
			n := 0
			m.Range(func(k, v interface{}) bool {
				n++
				return true
			})
			return false
		})
	}
}

// 测试 map+RWMutex 和 sync.map 的性能
//  size: map 大小
//  f1、f2 是针对统一功能（读、写、遍历等），分别使用 map+RWMutex 和使用 sync.map 时对应的函数。
func doTestSyncMap_Perf(t *testing.T, name string, size int, f1 func(key string, m map[interface{}]interface{}, mu *sync.RWMutex) bool, f2 func(keys string, m *sync.Map) bool) {
	util.Println("========================= %v(%v):", name, size)

	keys := make([]string, size)
	m1 := make(map[interface{}]interface{}, size)
	m2 := &sync.Map{}
	mu := &sync.RWMutex{}
	for i := 0; i < size; i++ {
		k := fmt.Sprintf("test_key_%v", i)
		keys[i] = k
		m1[k] = i
		m2.Store(k, i)
	}

	a := _test.Perf(func(i int) {
		for _, key := range keys {
			if !f1(key, m1, mu) {
				return
			}
		}
	}, _test.PerfOptions{NoHead: true})

	b := _test.Perf(func(i int) {
		for _, key := range keys {
			if !f2(key, m2) {
				return
			}
		}
	}, _test.PerfOptions{NoHead: true})

	util.Println("    map+RWMutex 的性能是 sync.map 的 %.2f%%\n", 100*a.QPS/b.QPS)
}
