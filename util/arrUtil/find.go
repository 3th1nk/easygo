package arrUtil

import (
	"github.com/3th1nk/easygo/util"
	"github.com/modern-go/reflect2"
	"reflect"
)

// 从数组中查找符合条件的元素，并返回包含这些元素的新数组
//   slice: 要查找的数组
//   match: 匹配函数
//   n: 最多返回的数量。<0 表示全部
func Find(slice interface{}, match func(i int) bool, limit ...int) interface{} {
	if v := doFind(slice, match, util.IfEmptyIntSlice(limit, -1)); v != nil {
		return v.Interface()
	}
	return nil
}

// 从数组中查找符合条件的元素，并返回包含这些元素的新数组
//   slice: 要查找的数组
//   match: 匹配函数
//   n: 最多返回的数量。<0 表示全部
func FindN(slice interface{}, match func(i int) bool, limit int, offset ...int) interface{} {
	if v := doFind(slice, match, limit, offset...); v != nil {
		return v.Interface()
	}
	return nil
}

func doFind(slice interface{}, match func(i int) bool, limit int, offset ...int) *reflect.Value {
	if reflect2.IsNil(slice) {
		return nil
	}

	reflectVal := reflect.ValueOf(slice)
	if kind := reflectVal.Kind(); kind != reflect.Array && kind != reflect.Slice {
		panic("param 'slice' must be slice")
	}

	if limit < 0 {
		limit = reflectVal.Len()
	}
	if limit == 0 {
		return nil
	}

	result, cnt, skip, theOffset := reflect.MakeSlice(reflectVal.Type(), 0, limit), 0, 0, util.IfEmptyIntSlice(offset, 0)
	for i, n := 0, reflectVal.Len(); i < n; i++ {
		if match == nil || match(i) {
			if skip < theOffset {
				skip++
			} else if result, cnt = reflect.Append(result, reflectVal.Index(i)), cnt+1; cnt == limit {
				break
			}
		}
	}

	return &result
}

// 返回数组中查找符合条件的元素的个数
//   slice: 要查找的数组
//   match: 匹配函数
func Count(slice interface{}, match func(i int) bool) int {
	if v := doFind(slice, match, -1); v != nil {
		return v.Len()
	}
	return 0
}

// 从数组中查找第一个符合条件的元素，并返回该元素
//   a: 要查找的数组
//   match: 匹配函数
func First(slice interface{}, match func(i int) bool) (interface{}, bool) {
	if v := doFind(slice, match, 1); v != nil && v.Len() != 0 {
		return v.Index(0).Interface(), true
	}
	return nil, false
}

// 从数组中查找第一个符合条件的元素，并返回该元素
//   a: 要查找的数组
//   match: 匹配函数
func MustFirst(slice interface{}, match func(i int) bool) interface{} {
	v, _ := First(slice, match)
	return v
}

func FindInt(a []int, match func(i int, v int) bool, limit ...int) []int {
	if a == nil {
		return nil
	}

	theLimit := -1
	if len(limit) != 0 && limit[0] >= 0 {
		theLimit = limit[0]
	}

	arr, j := make([]int, len(a)), 0
	for i, v := range a {
		if theLimit == j {
			break
		} else if match == nil || match(i, v) {
			arr[j] = v
			j++
		}
	}
	return arr[:j]
}

func FirstInt(a []int, match ...func(i int, v int) bool) (v int, ok bool) {
	var f func(i int, v int) bool
	if len(match) != 0 {
		f = match[0]
	}
	arr := FindInt(a, f, 1)
	if len(arr) != 0 {
		return arr[0], true
	}
	return
}

func MustFirstInt(a []int, match ...func(i int, v int) bool) (v int) {
	v, _ = FirstInt(a, match...)
	return
}

func FindInt64(a []int64, match func(i int, v int64) bool, limit ...int) []int64 {
	if a == nil {
		return nil
	}

	theLimit := -1
	if len(limit) != 0 && limit[0] >= 0 {
		theLimit = limit[0]
	}

	arr, j := make([]int64, len(a)), 0
	for i, v := range a {
		if theLimit == j {
			break
		} else if match == nil || match(i, v) {
			arr[j] = v
			j++
		}
	}
	return arr[:j]
}

func FirstInt64(a []int64, match ...func(i int, v int64) bool) (v int64, ok bool) {
	var f func(i int, v int64) bool
	if len(match) != 0 {
		f = match[0]
	}
	arr := FindInt64(a, f, 1)
	if len(arr) != 0 {
		return arr[0], true
	}
	return
}

func MustFirstInt64(a []int64, match ...func(i int, v int64) bool) (v int64) {
	v, _ = FirstInt64(a, match...)
	return
}

func FindString(a []string, match func(i int, s string) bool, limit ...int) []string {
	if a == nil {
		return nil
	}

	arr, j := make([]string, len(a)), 0
	for i, s := range a {
		if len(limit) > 0 && limit[0] <= j {
			break
		} else if match == nil || match(i, s) {
			arr[j] = s
			j++
		}
	}
	return arr[:j]
}

func FirstString(a []string, match ...func(i int, s string) bool) (s string, ok bool) {
	var f func(i int, s string) bool
	if len(match) != 0 {
		f = match[0]
	}
	arr := FindString(a, f, 1)
	if len(arr) != 0 {
		return arr[0], true
	}
	return
}

func MustFirstString(a []string, match ...func(i int, s string) bool) (s string) {
	s, _ = FirstString(a, match...)
	return
}
