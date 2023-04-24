package arrUtil

import (
	"github.com/modern-go/reflect2"
	"reflect"
	"sort"
	"strings"
)

// 获取指定值在数组中的索引。-1 表示不在数组中
func IndexOfInt(a []int, find int) int {
	for i, n := range a {
		if n == find {
			return i
		}
	}
	return -1
}

func ContainsInt(a []int, find int) bool {
	return -1 != IndexOfInt(a, find)
}

// 获取指定值在数组中的索引。-1 表示不在数组中
func IndexOfAnyInt(a []int, find []int) int {
	for i, n := range a {
		for _, v := range find {
			if n == v {
				return i
			}
		}
	}
	return -1
}

func ContainsAnyInt(a []int, find []int) bool {
	return -1 != IndexOfAnyInt(a, find)
}

// 获取指定值在数组中的索引。-1 表示不在数组中
func IndexOfInt64(a []int64, find int64) int {
	for i, n := range a {
		if n == find {
			return i
		}
	}
	return -1
}

func ContainsInt64(a []int64, find int64) bool {
	return -1 != IndexOfInt64(a, find)
}

// 获取指定值在数组中的索引。-1 表示不在数组中
func IndexOfAnyInt64(a []int64, find []int64) int {
	for i, n := range a {
		for _, v := range find {
			if n == v {
				return i
			}
		}
	}
	return -1
}

func ContainsAnyInt64(a []int64, find []int64) bool {
	return -1 != IndexOfAnyInt64(a, find)
}

// 获取指定值在数组中的索引。-1 表示不在数组中
func IndexOfInt32(a []int32, find int32) int {
	for i, n := range a {
		if n == find {
			return i
		}
	}
	return -1
}

func ContainsInt32(a []int32, find int32) bool {
	return -1 != IndexOfInt32(a, find)
}

// 获取指定值在数组中的索引。-1 表示不在数组中
func IndexOfAnyInt32(a []int32, find []int32) int {
	for i, n := range a {
		for _, v := range find {
			if n == v {
				return i
			}
		}
	}
	return -1
}

func ContainsAnyInt32(a []int32, find []int32) bool {
	return -1 != IndexOfAnyInt32(a, find)
}

// 获取指定值在数组中的索引。-1 表示不在数组中
func IndexOfString(a []string, find string, ignoreCase ...bool) int {
	if len(ignoreCase) != 0 && ignoreCase[0] {
		for i, s := range a {
			if strings.EqualFold(s, find) {
				return i
			}
		}
	} else {
		for i, s := range a {
			if s == find {
				return i
			}
		}
	}
	return -1
}

func ContainsString(a []string, find string, ignoreCase ...bool) bool {
	return -1 != IndexOfString(a, find, ignoreCase...)
}

// 获取指定值在数组中的索引。-1 表示不在数组中
func IndexOfAnyString(a []string, find []string, ignoreCase ...bool) int {
	if len(ignoreCase) != 0 && ignoreCase[0] {
		for i, n := range a {
			for _, v := range find {
				if strings.EqualFold(n, v) {
					return i
				}
			}
		}
	} else {
		for i, n := range a {
			for _, v := range find {
				if n == v {
					return i
				}
			}
		}
	}
	return -1
}

func ContainsAnyString(a []string, find []string, ignoreCase ...bool) bool {
	return -1 != IndexOfAnyString(a, find, ignoreCase...)
}

// 获取指定值在数组中的索引。-1 表示不在数组中
func IndexOfStringF(a []string, f func(i int, str string) bool) int {
	for i, s := range a {
		if f(i, s) {
			return i
		}
	}
	return -1
}

func ContainsStringF(a []string, f func(i int, str string) bool) bool {
	return -1 != IndexOfStringF(a, f)
}

// 获取指定值在有序数组中的索引。-1 表示不在数组中。数组必须是升序排序的
func IndexOfSortedInt(a []int, find int) int {
	pos := sort.SearchInts(a, find)
	if pos < len(a) && a[pos] == find {
		return pos
	}
	return -1
}

// 获取指定值在有序数组中的索引。-1 表示不在数组中。数组必须是升序排序的
func IndexOfSortedInt64(a []int64, find int64) int {
	count := len(a)
	pos := sort.Search(count, func(i int) bool { return a[i] >= find })
	if pos >= count {
		return -1
	} else if v := a[pos]; v != find {
		return -1
	} else {
		return pos
	}
}

// 获取指定值在有序数组中的索引。-1 表示不在数组中。数组必须是升序排序的
func IndexOfSortedString(a []string, find string) int {
	pos := sort.SearchStrings(a, find)
	if pos < len(a) && a[pos] == find {
		return pos
	}
	return -1
}

func IndexOf(a interface{}, match func(i int) bool) int {
	if reflect2.IsNil(a) {
		return -1
	}

	reflectVal := reflect.ValueOf(a)
	if kind := reflectVal.Kind(); kind != reflect.Array && kind != reflect.Slice {
		panic("参数 a 必须是切片类型")
	}

	for i, n := 0, reflectVal.Len(); i < n; i++ {
		if match == nil || match(i) {
			return i
		}
	}

	return -1
}

func Contains(a interface{}, match func(i int) bool) bool {
	return -1 != IndexOf(a, match)
}

func IndexOfValue(slice interface{}, val interface{}, equal ...func(a, b interface{}) bool) int {
	if reflect2.IsNil(slice) {
		return -1
	}

	reflectVal := reflect.ValueOf(slice)
	if kind := reflectVal.Kind(); kind != reflect.Array && kind != reflect.Slice {
		panic("参数 a 必须是切片类型")
	}

	var theEqual func(a, b interface{}) bool
	if len(equal) != 0 {
		theEqual = equal[0]
	}

	for i, n := 0, reflectVal.Len(); i < n; i++ {
		tmp := reflectVal.Index(i).Interface()
		if theEqual != nil {
			if theEqual(val, tmp) {
				return i
			}
		} else {
			if val == tmp {
				return i
			}
		}
	}
	return -1
}

func ContainsValue(slice interface{}, val interface{}, equal ...func(a, b interface{}) bool) bool {
	return -1 != IndexOfValue(slice, val, equal...)
}
