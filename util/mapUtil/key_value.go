package mapUtil

import (
	"github.com/3th1nk/easygo/util/sortUtil"
	"github.com/modern-go/reflect2"
	"reflect"
	"sort"
)

func Keys(m interface{}) (keys interface{}) {
	if reflect2.IsNil(m) {
		return nil
	}

	reflectVal := reflect.ValueOf(m)
	if kind := reflectVal.Kind(); kind != reflect.Map {
		panic("参数必须是 Map 类型")
	}

	n := reflectVal.Len()
	result := reflect.MakeSlice(reflect.SliceOf(reflectVal.Type().Key()), n, n)
	for iter, idx := reflectVal.MapRange(), 0; iter.Next(); {
		result.Index(idx).Set(iter.Key())
		idx++
	}
	return result.Interface()
}

func StringKeys(m interface{}) (keys []string) {
	if reflect2.IsNil(m) {
		return nil
	}
	return Keys(m).([]string)
}

func SortedStringKeys(m interface{}) (keys []string) {
	keys = StringKeys(m)
	sort.Strings(keys)
	return
}

func IntKeys(m interface{}) (keys []int) {
	if reflect2.IsNil(m) {
		return nil
	}
	return Keys(m).([]int)
}

func SortedIntKeys(m interface{}) (keys []int) {
	keys = IntKeys(m)
	sort.Ints(keys)
	return
}

func Int64Keys(m interface{}) (keys []int64) {
	if reflect2.IsNil(m) {
		return nil
	}
	return Keys(m).([]int64)
}

func SortedInt64Keys(m interface{}) (keys []int64) {
	keys = Int64Keys(m)
	sortUtil.Int64s(keys)
	return
}

func Values(m interface{}) (values interface{}) {
	if reflect2.IsNil(m) {
		return nil
	}

	reflectVal := reflect.ValueOf(m)
	if kind := reflectVal.Kind(); kind != reflect.Map {
		panic("参数必须是 Map 类型")
	}

	n := reflectVal.Len()
	result := reflect.MakeSlice(reflect.SliceOf(reflectVal.Type().Elem()), n, n)
	for iter, idx := reflectVal.MapRange(), 0; iter.Next(); {
		result.Index(idx).Set(iter.Value())
		idx++
	}
	return result.Interface()
}

func StringValues(m interface{}) (keys []string) {
	if reflect2.IsNil(m) {
		return nil
	}
	return Values(m).([]string)
}

func SortedStringValues(m interface{}) (keys []string) {
	keys = StringValues(m)
	sort.Strings(keys)
	return
}

func IntValues(m interface{}) (keys []int) {
	if reflect2.IsNil(m) {
		return nil
	}
	return Values(m).([]int)
}

func SortedIntValues(m interface{}) (keys []int) {
	keys = IntValues(m)
	sort.Ints(keys)
	return
}
