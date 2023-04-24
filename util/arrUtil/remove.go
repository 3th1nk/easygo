package arrUtil

import (
	"fmt"
	"github.com/modern-go/reflect2"
	"reflect"
)

func RemoveStringAt(a []string, i ...int) []string {
	if reflect2.IsNil(a) {
		return a
	}
	return doRemoveAt(reflect.ValueOf(a), i...).([]string)
}

func RemoveIntAt(a []int, i ...int) []int {
	if reflect2.IsNil(a) {
		return a
	}
	return doRemoveAt(reflect.ValueOf(a), i...).([]int)
}

func RemoveAt(a interface{}, i ...int) interface{} {
	if reflect2.IsNil(a) {
		return a
	}
	reflectVal := reflect.ValueOf(a)
	if kind := reflectVal.Kind(); kind != reflect.Array && kind != reflect.Slice {
		panic("参数 a 必须是切片类型")
	}
	return doRemoveAt(reflectVal, i...)
}

func doRemoveAt(reflectVal reflect.Value, i ...int) interface{} {
	n := reflectVal.Len()
	result := reflect.MakeSlice(reflectVal.Type(), 0, n)
	for _, i := range i {
		if i < 0 || i >= n {
			panic(fmt.Sprintf("index out of range: %v(len=%v)", i, n))
		}
		result = reflect.AppendSlice(result, reflectVal.Slice(0, i))
		result = reflect.AppendSlice(result, reflectVal.Slice(i+1, n))
	}
	return result.Interface()
}

func RemoveString(a []string, val ...string) []string {
	return RemoveStringN(a, -1, val...)
}

func RemoveStringN(a []string, n int, val ...string) []string {
	if reflect2.IsNil(a) {
		return a
	}
	return doRemove(reflect.ValueOf(a), n, func(i int) bool {
		for _, vv := range val {
			if vv == a[i] {
				return true
			}
		}
		return false
	}).([]string)
}

func RemoveStringIf(a []string, f func(i int) bool) []string {
	return RemoveStringIfN(a, -1, f)
}

func RemoveStringIfN(a []string, n int, f func(i int) bool) []string {
	if reflect2.IsNil(a) {
		return nil
	}
	return RemoveIfN(a, n, f).([]string)
}

func RemoveInt(a []int, val ...int) []int {
	return RemoveIntN(a, -1, val...)
}

func RemoveIntN(a []int, n int, val ...int) []int {
	if reflect2.IsNil(a) {
		return a
	}
	return doRemove(reflect.ValueOf(a), n, func(i int) bool {
		for _, vv := range val {
			if vv == a[i] {
				return true
			}
		}
		return false
	}).([]int)
}

func RemoveIntIf(a []int, f func(i int) bool) []int {
	return RemoveIntIfN(a, -1, f)
}

func RemoveIntIfN(a []int, n int, f func(i int) bool) []int {
	if reflect2.IsNil(a) {
		return nil
	}
	return RemoveIfN(a, n, f).([]int)
}

func Remove(a interface{}, v interface{}) interface{} {
	return RemoveN(a, -1, v)
}

func RemoveN(a interface{}, n int, v interface{}) interface{} {
	if reflect2.IsNil(a) {
		return a
	}
	reflectVal := reflect.ValueOf(a)
	if kind := reflectVal.Kind(); kind != reflect.Array && kind != reflect.Slice {
		panic("参数 a 必须是切片类型")
	}
	return doRemove(reflectVal, n, func(i int) bool {
		return v == reflectVal.Index(i).Interface()
	})
}

func RemoveIf(a interface{}, f func(i int) bool) interface{} {
	return RemoveIfN(a, -1, f)
}

func RemoveIfN(a interface{}, n int, f func(i int) bool) interface{} {
	if reflect2.IsNil(a) {
		return a
	}
	reflectVal := reflect.ValueOf(a)
	if kind := reflectVal.Kind(); kind != reflect.Array && kind != reflect.Slice {
		panic("参数 a 必须是切片类型")
	}
	return doRemove(reflectVal, n, f)
}

func doRemove(reflectVal reflect.Value, n int, f func(i int) bool) interface{} {
	if n == 0 {
		return reflectVal.Interface()
	}

	count, removedCount := reflectVal.Len(), 0
	result := reflect.MakeSlice(reflectVal.Type(), 0, count)
	for i := 0; i < count; i++ {
		if removedCount == n || f == nil || !f(i) {
			result = reflect.Append(result, reflectVal.Index(i))
		} else if n > 0 {
			removedCount++
		}
	}

	return result.Interface()
}

func RemoveDuplicate(arr interface{}, f func(val interface{}) string) interface{} {
	if v := doRemoveDuplicate(arr, f); v != nil {
		return v.Interface()
	}
	return nil
}

func doRemoveDuplicate(arr interface{}, f func(val interface{}) string) *reflect.Value {
	if reflect2.IsNil(arr) || f == nil {
		return nil
	}

	reflectVal := reflect.ValueOf(arr)
	if kind := reflectVal.Kind(); kind != reflect.Array && kind != reflect.Slice {
		panic("param 'slice' must be slice")
	}

	n := reflectVal.Len()

	result := reflect.MakeSlice(reflectVal.Type(), 0, n)
	exists := make(map[string]bool, n)
	for i, n := 0, reflectVal.Len(); i < n; i++ {
		v := reflectVal.Index(i)

		if key := f(v.Interface()); !exists[key] {
			result = reflect.Append(result, reflectVal.Index(i))
			exists[key] = true
		}
	}

	return &result
}
