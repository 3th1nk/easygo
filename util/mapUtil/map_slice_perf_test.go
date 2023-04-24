package mapUtil

import (
	"fmt"
	"github.com/3th1nk/easygo/internal/_test"
	"github.com/3th1nk/easygo/util"
	"github.com/3th1nk/easygo/util/arrUtil"
	"testing"
)

// 测试从 map 和 slice 中查询的性能差异
func TestMapSlicePerf_Query(t *testing.T) {
	for _, size := range []int{10, 100, 1000, 10000} {
		doTestMapSliceQuery_Perf(t, "Query Exists", size, func(key string, m map[interface{}]interface{}) bool {
			_ = m[key]
			return true
		}, func(key string, arr []interface{}) bool {
			for _, v := range arr {
				if v == key {
					break
				}
			}
			return true
		})
	}

	for _, size := range []int{10, 100, 1000, 10000} {
		doTestMapSliceQuery_Perf(t, "Query NotExists", size, func(key string, m map[interface{}]interface{}) bool {
			_ = m[key+key]
			return true
		}, func(key string, arr []interface{}) bool {
			for _, v := range arr {
				if v == key+key {
					break
				}
			}
			return true
		})
	}
}

// 测试 构建 map 和 slice 并从中查询的性能差异
func TestMapSlicePerf_MakeAndQuery(t *testing.T) {
	for _, size := range []int{10, 100, 1000, 10000} {
		keys := makeKeys(size)
		doTestMapSliceQuery_Perf(t, "MakeAndQuery Exists", size, func(key string, m map[interface{}]interface{}) bool {
			m = makeMap(keys)
			_ = m[key]
			return true
		}, func(key string, arr []interface{}) bool {
			arr = makeSlice(keys)
			for _, v := range arr {
				if v == key {
					break
				}
			}
			return true
		})
	}

	for _, size := range []int{10, 100, 1000, 10000} {
		keys := makeKeys(size)
		doTestMapSliceQuery_Perf(t, "MakeAndQuery NotExists", size, func(key string, m map[interface{}]interface{}) bool {
			m = makeMap(keys)
			_ = m[key]
			return true
		}, func(key string, arr []interface{}) bool {
			arr = makeSlice(keys)
			for _, v := range arr {
				if v == key+key {
					break
				}
			}
			return true
		})
	}
}

func doTestMapSliceQuery_Perf(t *testing.T, name string, size int, f1 func(key string, m map[interface{}]interface{}) bool, f2 func(key string, arr []interface{}) bool) {
	util.Println("========================= %v(%v):", name, size)

	keys := makeKeys(size)
	dict := makeMap(keys)
	arr := makeSlice(keys)

	a := _test.Perf(func(i int) {
		for _, key := range keys {
			if !f1(key, dict) {
				return
			}
		}
	}, _test.PerfOptions{NoHead: true})

	b := _test.Perf(func(i int) {
		for _, key := range keys {
			if !f2(key, arr) {
				return
			}
		}
	}, _test.PerfOptions{NoHead: true})

	util.Println("    map 的性能是 slice 的 %.2f%%\n", 100*a.QPS/b.QPS)
}

func makeKeys(size int) []string {
	keys := make([]string, size)
	for i := range keys {
		keys[i] = fmt.Sprintf("test_key_%v", i)
	}
	return keys
}

func makeMap(keys []string) map[interface{}]interface{} {
	m := make(map[interface{}]interface{}, len(keys))
	for i, k := range keys {
		m[k] = i
	}
	return m
}

func makeSlice(keys []string) []interface{} {
	return arrUtil.StrToObj(keys)
}
