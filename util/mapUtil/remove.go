package mapUtil

import (
	"github.com/3th1nk/easygo/util/arrUtil"
	"github.com/3th1nk/easygo/util/convertor"
	"strings"
)

func RemoveMapObjectDuplicate(arr []StringObjectMap, keys ...string) []StringObjectMap {
	if len(arr) == 1 {
		return arr
	}
	n := len(keys)

	val := arrUtil.RemoveDuplicate(arr, func(val interface{}) string {
		obj := val.(StringObjectMap)
		if n == 0 {
			return convertor.ToStringNoError(val)
		} else if n == 1 {
			return convertor.ToStringNoError(obj[keys[0]])
		} else {
			sli := make([]string, n)
			for i, key := range keys {
				sli[i] = convertor.ToStringNoError(obj[key])
			}
			return strings.Join(sli, "*")
		}
	})
	if val != nil {
		return val.([]StringObjectMap)
	}
	return arr
}

func RemoveMapDuplicate(arr []map[string]interface{}, keys ...string) []map[string]interface{} {
	if len(arr) == 1 {
		return arr
	}
	n := len(keys)

	val := arrUtil.RemoveDuplicate(arr, func(val interface{}) string {
		obj := val.(map[string]interface{})

		if n == 0 {
			return convertor.ToStringNoError(val)
		} else if n == 1 {
			return convertor.ToStringNoError(obj[keys[0]])
		} else {
			sli := make([]string, n)
			for i, key := range keys {
				sli[i] = convertor.ToStringNoError(obj[key])
			}
			return strings.Join(sli, "*")
		}
	})
	if val != nil {
		return val.([]map[string]interface{})
	}
	return arr
}

func RemoveMapStringDuplicate(arr []StringMap, keys ...string) []StringMap {
	if len(arr) == 1 {
		return arr
	}
	n := len(keys)

	val := arrUtil.RemoveDuplicate(arr, func(val interface{}) string {
		obj := val.(StringMap)
		if n == 0 {
			return convertor.ToStringNoError(val)
		} else if n == 1 {
			return convertor.ToStringNoError(obj[keys[0]])
		} else {
			sli := make([]string, n)
			for i, key := range keys {
				sli[i] = convertor.ToStringNoError(obj[key])
			}
			return strings.Join(sli, "*")
		}
	})
	if val != nil {
		return val.([]StringMap)
	}
	return arr
}
