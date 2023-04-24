package mapUtil

import (
	"github.com/modern-go/reflect2"
	"reflect"
)

// 从字典中查找符合条件的元素，并返回包含这些元素的新字典
//   slice: 要查找的数组
//   match: 匹配函数
//   count: 最多返回的数量，<0 表示全部
func Find(dict interface{}, match func(key, val interface{}) bool) interface{} {
	if v := doFind(dict, match); v != nil {
		return v.Interface()
	}
	return nil
}

func doFind(dict interface{}, match func(key, val interface{}) bool) *reflect.Value {
	if reflect2.IsNil(dict) {
		return nil
	}

	reflectVal := reflect.ValueOf(dict)
	if kind := reflectVal.Kind(); kind != reflect.Map {
		panic("参数 dict 必须是字典类型")
	}

	maxLen := reflectVal.Len()
	result := reflect.MakeMapWithSize(reflectVal.Type(), maxLen)

	keys := reflectVal.MapKeys()
	for _, key := range keys {
		val := reflectVal.MapIndex(key)
		if match == nil || match(key.Interface(), val.Interface()) {
			result.SetMapIndex(key, val)
		}
	}

	return &result
}
