package mapUtil

import (
	"fmt"
	"github.com/3th1nk/easygo/util/mathUtil"
	"github.com/3th1nk/easygo/util/sortUtil"
	"github.com/modern-go/reflect2"
	"reflect"
)

// 求两个或多个 map 的差集：
//   如果 key 在 src 中存在、但 dest 中不存在，则结果集中的值为该类型的默认值；
//   如果 key 在 dest 中存在、但 src 中不存在、或者值与 src 中不相等，则结果集中的值为 dest 中的值。
// src 和 dest 的类型必须相同。
func Diff(src interface{}, dest ...interface{}) interface{} {
	return diff(src, dest, nil)
}

// 求两个或多个 map 的差集：
//   如果 key 在 src 中存在、但 dest 中不存在，则结果集中的值为该类型的默认值；
//   如果 key 在 dest 中存在、但 src 中不存在、或者值与 src 中不相等，则结果集中的值为 dest 中的值。
func DiffStr(src map[string]string, dest ...map[string]string) map[string]string {
	args := make([]interface{}, len(dest))
	for i, item := range dest {
		args[i] = item
	}
	return diff(src, args, nil).(map[string]string)
}

// 求两个或多个 map 的差集：
//   如果 key 在 src 中存在、但 dest 中不存在，则结果集中的值为该类型的默认值；
//   如果 key 在 dest 中存在、但 src 中不存在、或者值与 src 中不相等，则结果集中的值为 dest 中的值。
func DiffStr2(src map[string]string, dest map[string]string) (matches, changed, added, removed map[string]string) {
	if len(src) == 0 {
		return nil, nil, dest, nil
	} else if len(dest) == 0 {
		return nil, nil, nil, src
	}

	maxLen := mathUtil.MaxInt(len(src), len(dest))
	matches, changed, added, removed = make(map[string]string, maxLen), make(map[string]string, maxLen), make(map[string]string, maxLen), make(map[string]string, maxLen)
	for k, v := range src {
		if v2, ok := dest[k]; !ok {
			removed[k] = v
		} else if v != v2 {
			changed[k] = v2
		} else {
			matches[k] = v
		}
	}
	for k, v := range dest {
		if _, ok := src[k]; !ok {
			added[k] = v
		}
	}
	return
}

// 求两个或多个 map 的差集：
//   如果 key 在 src 中存在、但 dest 中不存在，则结果集中的值为该类型的默认值；
//   如果 key 在 dest 中存在、但 src 中不存在、或者值与 src 中不相等，则结果集中的值为 dest 中的值。
func DiffStrObject(src map[string]interface{}, dest ...map[string]interface{}) map[string]interface{} {
	args := make([]interface{}, len(dest))
	for i, item := range dest {
		args[i] = item
	}
	return diff(src, args, nil).(map[string]interface{})
}

func DiffStrObject2(src map[string]interface{}, dest map[string]interface{}, f ...func(a, b interface{}) bool) (matches, changed, added, removed map[string]interface{}) {
	a, b, c, d := Diff2(src, dest, f...)
	matches, _ = a.(map[string]interface{})
	changed, _ = b.(map[string]interface{})
	added, _ = c.(map[string]interface{})
	removed, _ = d.(map[string]interface{})
	return
}

// 求两个或多个 map 的并集
func Union(src interface{}, dest ...interface{}) interface{} {
	args := make([]interface{}, 0, len(dest)+1)
	if !reflect2.IsNil(src) {
		args = append(args, src)
	}
	args = append(args, dest...)
	return diff(nil, args, nil)
}

func UnionStr(src map[string]string, dest ...map[string]string) map[string]string {
	args := make([]interface{}, len(dest))
	for i, item := range dest {
		args[i] = item
	}
	return Union(src, args...).(map[string]string)
}

func UnionStrObject(src map[string]interface{}, dest ...map[string]interface{}) map[string]interface{} {
	args := make([]interface{}, len(dest))
	for i, item := range dest {
		args[i] = item
	}
	return Union(src, args...).(map[string]interface{})
}

func diff(src interface{}, dest []interface{}, f func(a, b interface{}) bool) interface{} {
	// 确认参数类型相同，并取出所有的 map
	types := make([]reflect.Type, 0, len(dest)+1)
	if !reflect2.IsNil(src) {
		types = append(types, reflect.TypeOf(src))
	}
	for _, v := range dest {
		if !reflect2.IsNil(v) {
			types = append(types, reflect.TypeOf(v))
		}
	}
	if len(types) == 0 {
		return nil
	}
	for i, n := 1, len(types); i < n; i++ {
		if types[0] != types[i] {
			panic(fmt.Sprintf("参数类型不同 (%v, %v)", types[0], types[i]))
		}
	}

	isPtrType := types[0].Kind() == reflect.Ptr

	// 获取参数对应的 map 对象（处理指向 map 的指针），并判断类型
	maps := make([]reflect.Value, 0, len(dest))
	if !reflect2.IsNil(src) {
		if isPtrType {
			maps = append(maps, reflect.ValueOf(src).Elem())
		} else {
			maps = append(maps, reflect.ValueOf(src))
		}
	}
	for _, v := range dest {
		if !reflect2.IsNil(v) {
			if isPtrType {
				maps = append(maps, reflect.ValueOf(v).Elem())
			} else {
				maps = append(maps, reflect.ValueOf(v))
			}
		}
	}
	if maps[0].Kind() != reflect.Map {
		panic(fmt.Sprintf("参数不是 map 类型 (%v)", types[0]))
	}

	var resultMap reflect.Value
	if reflect2.IsNil(src) {
		// src 为 nil，直接返回 dest 的合并结果
		resultMap = mergeMaps(maps)
	} else if len(maps) == 1 {
		// src 不为 nil、dest 为空，返回空 map
		resultMap = reflect.MakeMap(maps[0].Type())
	} else {
		resultMap = reflect.MakeMap(maps[0].Type())
		srcMap := maps[0]
		destMap := mergeMaps(maps[1:])
		// 在 src 中存在、但是在 dest 中不存在或者值不相等的元素
		for _, key := range srcMap.MapKeys() {
			destVal := destMap.MapIndex(key)
			if destVal.IsValid() {
				if f != nil {
					if !f(srcMap.MapIndex(key).Interface(), destVal.Interface()) {
						resultMap.SetMapIndex(key, destVal)
					}
				} else {
					if !reflect.DeepEqual(srcMap.MapIndex(key).Interface(), destVal.Interface()) {
						resultMap.SetMapIndex(key, destVal)
					}
				}
			} else {
				resultMap.SetMapIndex(key, reflect.New(srcMap.MapIndex(key).Type()).Elem())
			}
		}
		// 在 dest 中存在、但是在 src 中不存在的元素
		for _, key := range destMap.MapKeys() {
			if !srcMap.MapIndex(key).IsValid() {
				resultMap.SetMapIndex(key, destMap.MapIndex(key))
			}
		}
	}

	if isPtrType {
		result := reflect.New(types[0]).Elem()
		resultVal := reflect.New(resultMap.Type())
		resultVal.Elem().Set(resultMap)
		result.Set(resultVal)
		return result.Interface()
	}
	return resultMap.Interface()
}

func Diff2(src interface{}, dest interface{}, f ...func(a, b interface{}) bool) (matches, changed, added, removed interface{}) {
	if src == 0 {
		return nil, nil, dest, nil
	} else if dest == nil {
		return nil, nil, nil, src
	}

	srcMap, destMap := reflect.ValueOf(src), reflect.ValueOf(dest)
	typeSrc, typeDest := srcMap.Type(), destMap.Type()
	if typeSrc != typeDest {
		panic(fmt.Errorf("param type dismatch: (%v, %v)", typeSrc, typeDest))
	}
	if typeSrc.Kind() != reflect.Map {
		panic(fmt.Errorf("param must be map"))
	}

	matchesMap, changedMap, addedMap, removedMap := reflect.MakeMap(typeSrc), reflect.MakeMap(typeSrc), reflect.MakeMap(typeSrc), reflect.MakeMap(typeSrc)

	// 在 src 中存在、但是在 dest 中不存在或者值不相等的元素
	for _, key := range srcMap.MapKeys() {
		srcVal := srcMap.MapIndex(key)
		destVal := destMap.MapIndex(key)
		if destVal.IsValid() {
			var eq bool
			if len(f) != 0 && f[0] != nil {
				eq = f[0](srcVal.Interface(), destVal.Interface())
			} else {
				eq = reflect.DeepEqual(srcVal.Interface(), destVal.Interface())
			}
			if eq {
				matchesMap.SetMapIndex(key, srcVal)
			} else {
				changedMap.SetMapIndex(key, destVal)
			}
		} else {
			removedMap.SetMapIndex(key, srcVal)
		}
	}
	// 在 dest 中存在、但是在 src 中不存在的元素
	for _, key := range destMap.MapKeys() {
		if !srcMap.MapIndex(key).IsValid() {
			addedMap.SetMapIndex(key, destMap.MapIndex(key))
		}
	}

	return matchesMap.Interface(), changedMap.Interface(), addedMap.Interface(), removedMap.Interface()
}

func mergeMaps(src []reflect.Value) reflect.Value {
	dest := reflect.MakeMap(src[0].Type())
	for _, item := range src {
		for _, key := range item.MapKeys() {
			dest.SetMapIndex(key, item.MapIndex(key))
		}
	}
	return dest
}

// LeftDiffKeyInt64 提取Key差异：左边有，右边没有
func LeftDiffKeyInt64(left map[int64]bool, right map[int64]bool, sorted ...bool) []int64 {

	keys := make([]int64, 0, len(left))

	for k := range left {
		if right == nil || !right[k] {
			keys = append(keys, k)
		}
	}

	if len(sorted) > 0 && sorted[0] {
		sortUtil.Int64s(keys)
	}

	return keys
}
