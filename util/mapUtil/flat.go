package mapUtil

import (
	"github.com/3th1nk/easygo/util"
	"github.com/3th1nk/easygo/util/arrUtil"
	"github.com/3th1nk/easygo/util/convertor"
	"github.com/3th1nk/easygo/util/regexpUtil"
	"github.com/modern-go/reflect2"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

const (
	defaultFlatSep = "."
)

func FlatMap(a interface{}, sep ...string) (StringObjectMap, error) {
	return FlatMapF(a, func(prefix, key string) string { return prefix + util.IfEmptyStringSlice(sep, defaultFlatSep) + key })
}

func MustFlatMap(a interface{}, sep ...string) StringObjectMap {
	val, _ := FlatMapF(a, func(prefix, key string) string { return prefix + util.IfEmptyStringSlice(sep, defaultFlatSep) + key })
	return val
}

func FlatMapF(a interface{}, joinKey func(prefix, key string) string) (StringObjectMap, error) {
	if reflect2.IsNil(a) {
		return nil, nil
	}

	if joinKey == nil {
		joinKey = func(prefix, key string) string { return prefix + defaultFlatSep + key }
	}

	dict := make(map[string]interface{}, 16)
	flat := &flatter{
		new: func() interface{} { return make(map[string]interface{}, 32) },
		set: func(dict interface{}, key string, val interface{}) { dict.(map[string]interface{})[key] = val },
		each: func(dict interface{}, f func(key string, val interface{})) {
			for key, val := range dict.(map[string]interface{}) {
				f(key, val)
			}
		},
	}
	if err := flat.do(dict, "", a, joinKey); err != nil {
		return nil, err
	}
	return dict, nil
}

func MustFlatMapF(a interface{}, joinKey func(prefix, key string) string) StringObjectMap {
	val, _ := FlatMapF(a, joinKey)
	return val
}

func FlatStringMap(a interface{}, sep ...string) (StringMap, error) {
	return FlatStringMapF(a, func(prefix, key string) string { return prefix + util.IfEmptyStringSlice(sep, defaultFlatSep) + key })
}

func MustFlatStringMap(a interface{}, sep ...string) StringMap {
	val, _ := FlatStringMapF(a, func(prefix, key string) string { return prefix + util.IfEmptyStringSlice(sep, defaultFlatSep) + key })
	return val
}

func FlatStringMapF(a interface{}, joinKey func(prefix, key string) string) (StringMap, error) {
	if reflect2.IsNil(a) {
		return nil, nil
	}

	if joinKey == nil {
		joinKey = func(prefix, key string) string { return prefix + defaultFlatSep + key }
	}

	dict := make(map[string]string, 16)
	flat := &flatter{
		new: func() interface{} { return make(map[string]string, 32) },
		set: func(dict interface{}, key string, val interface{}) {
			dict.(map[string]string)[key] = convertor.ToStringNoError(val)
		},
		each: func(dict interface{}, f func(key string, val interface{})) {
			for key, val := range dict.(map[string]interface{}) {
				f(key, val)
			}
		},
	}
	if err := flat.do(dict, "", a, joinKey); err != nil {
		return nil, err
	}
	return dict, nil
}

func MustFlatStringMapF(a interface{}, joinKey func(prefix, key string) string) StringMap {
	val, _ := FlatStringMapF(a, joinKey)
	return val
}

type flatter struct {
	new  func() interface{}
	set  func(dict interface{}, key string, val interface{})
	each func(dict interface{}, f func(key string, val interface{}))
}

func (this *flatter) do(dict interface{}, prefix string, val interface{}, joinKey func(prefix, key string) string) error {
	if reflect2.IsNil(val) {
		this.set(dict, prefix, val)
		return nil
	}

	refValue := reflect.ValueOf(val)
	for refValue.Kind() == reflect.Ptr {
		refValue = refValue.Elem()
	}
	refType := refValue.Type()
	switch refType.Kind() {
	default:
		this.set(dict, prefix, val)
	case reflect.Slice, reflect.Array:
		for i, n := 0, refValue.Len(); i < n; i++ {
			key := strconv.Itoa(i)
			if prefix != "" {
				key = joinKey(prefix, key)
			}
			if err := this.do(dict, key, refValue.Index(i).Interface(), joinKey); err != nil {
				return err
			}
		}
	case reflect.Map:
		for _, k := range refValue.MapKeys() {
			key := convertor.ToStringNoError(k.Interface())
			if prefix != "" {
				key = joinKey(prefix, key)
			}
			if err := this.do(dict, key, refValue.MapIndex(k).Interface(), joinKey); err != nil {
				return err
			}
		}
	case reflect.Struct:
		innerDict := this.new()
		if err := this.do(innerDict, prefix, val, joinKey); err != nil {
			return err
		} else {
			this.each(innerDict, func(key string, val interface{}) {
				if prefix != "" {
					key = joinKey(prefix, key)
				}
				this.set(dict, key, val)
			})
		}
	}
	return nil
}

func UnFlatMap(in map[string]interface{}, sep ...string) (interface{}, error) {
	return doUnFlatMap(in, util.IfEmptyStringSlice(sep, defaultFlatSep), nil)
}

func MustUnFlatMap(in map[string]interface{}, sep ...string) interface{} {
	val, _ := doUnFlatMap(in, util.IfEmptyStringSlice(sep, defaultFlatSep), nil)
	return val
}

func UnFlatMapF(in map[string]interface{}, splitKey func(path string) (prefix, key string)) (interface{}, error) {
	return doUnFlatMap(in, "", splitKey)
}

func MustUnFlatMapF(in map[string]interface{}, splitKey func(path string) (prefix, key string)) interface{} {
	val, _ := doUnFlatMap(in, "", splitKey)
	return val
}

func doUnFlatMap(in map[string]interface{}, sep string, splitKey func(path string) (prefix, key string)) (interface{}, error) {
	if len(in) == 0 {
		return nil, nil
	}

	if splitKey == nil {
		if sep == "" {
			sep = defaultFlatSep
		}
		sepLen := len(sep)
		splitKey = func(path string) (prefix, key string) {
			if pos := strings.LastIndex(path, sep); pos != -1 {
				return path[:pos], path[pos+sepLen:]
			}
			return "", path
		}
	}

	type splitInfo struct {
		path    []string
		key     []string
		idx     []int
		val     interface{}
		isSlice bool
	}
	prefixMap := make(map[string]*splitInfo, len(in))
	for _, path := range StringKeys(in) {
		for path != "" {
			prefix, key := splitKey(path)
			if obj, _ := prefixMap[prefix]; obj == nil {
				prefixMap[prefix] = &splitInfo{path: []string{path}, key: []string{key}}
			} else if !arrUtil.ContainsString(obj.key, key) {
				obj.path = append(obj.path, path)
				obj.key = append(obj.key, key)
			}
			path = prefix
		}
	}

	for _, obj := range prefixMap {
		// 先将 isSlice 置为 true，然后根据后续逻辑验证是否是数组
		obj.isSlice = true

		// 判断 key 是否全部都是数字，如果是就赋值到 idx 中
		obj.idx = make([]int, len(obj.key))
		for i, s := range obj.key {
			if s == "" {
				obj.isSlice = false
				break
			} else if !regexpUtil.IsInt(s) {
				obj.isSlice = false
				break
			} else {
				obj.idx[i], _ = strconv.Atoi(s)
			}
		}

		// 再次判断 idx 是否是 0-n 的数组
		if obj.isSlice {
			arr := make([]int, len(obj.idx))
			copy(arr, obj.idx)
			sort.Ints(arr)
			for i, n := range arr {
				if i != n {
					obj.isSlice = false
					break
				}
			}
		}

		if obj.isSlice {
			obj.val = make([]interface{}, len(obj.path))
		} else {
			obj.val = make(map[string]interface{}, 8)
		}
	}

	var setValue func(obj *splitInfo) interface{}
	setValue = func(obj *splitInfo) interface{} {
		if obj.isSlice {
			arr := obj.val.([]interface{})
			for idx, path := range obj.path {
				if val, _ := prefixMap[path]; val != nil {
					arr[obj.idx[idx]] = setValue(val)
				} else if val, ok := in[path]; ok {
					arr[obj.idx[idx]] = val
				}
			}
		} else {
			dict := obj.val.(map[string]interface{})
			for idx, path := range obj.path {
				if val, _ := prefixMap[path]; val != nil {
					dict[obj.key[idx]] = setValue(val)
				} else if val, ok := in[path]; ok {
					dict[obj.key[idx]] = val
				}
			}
		}
		return obj.val
	}
	return setValue(prefixMap[""]), nil
}
