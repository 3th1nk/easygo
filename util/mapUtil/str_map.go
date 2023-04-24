package mapUtil

import (
	"github.com/3th1nk/easygo/util/convertor"
	"github.com/3th1nk/easygo/util/jsonUtil"
	"sort"
	"strings"
)

type StringMap map[string]string

func NewStringMap(a ...map[string]string) StringMap {
	dict := StringMap{}
	dict.SetMulti(a...)
	return dict
}

func ToStringMapSlice(src []map[string]string) (dest []StringMap) {
	if src == nil {
		return nil
	}
	dest = make([]StringMap, len(src))
	for i, v := range src {
		dest[i] = v
	}
	return
}

func (this StringMap) Copy() StringMap {
	return NewStringMap(this)
}

func (this StringMap) ToDB() (bytes []byte, err error) {
	if len(this) != 0 {
		return jsonUtil.Marshal(this)
	}
	return []byte{}, nil
}

func (this *StringMap) FromDB(bytes []byte) error {
	return jsonUtil.Unmarshal(bytes, this)
}

func (this StringMap) Delete(key ...string) {
	if len(this) != 0 {
		for _, k := range key {
			delete(this, k)
		}
	}
}

func (this StringMap) Set(key, val string) {
	this[key] = val
}

func (this StringMap) SetMulti(a ...map[string]string) {
	for _, item := range a {
		for key, val := range item {
			this[key] = val
		}
	}
}

func (this StringMap) Keys() (arr []string) {
	if this != nil {
		arr = make([]string, len(this))
		n := 0
		for key, _ := range this {
			arr[n] = key
			n++
		}
	}
	return
}

func (this StringMap) Values() (arr []string) {
	if this != nil {
		arr = make([]string, len(this))
		n := 0
		for _, val := range this {
			arr[n] = val
			n++
		}
	}
	return
}

func (this StringMap) ValuesSortedByKey() (arr []string) {
	if this != nil {
		keys := this.Keys()
		sort.Strings(keys)
		arr = make([]string, len(this))
		n := 0
		for _, k := range keys {
			arr[n] = this.MustGet(k)
			n++
		}
	}
	return
}

func (this StringMap) Contains(key string) bool {
	if len(this) != 0 {
		_, found := this[key]
		return found
	}
	return false
}

func (this StringMap) ContainsAny(keys ...string) bool {
	if len(this) != 0 {
		for _, s := range keys {
			if _, ok := this[s]; ok {
				return true
			}
		}
	}
	return false
}

func (this StringMap) ContainsAll(keys []string) bool {
	if len(this) != 0 {
		for _, s := range keys {
			if _, ok := this[s]; !ok {
				return false
			}
		}
		return true
	}
	return false
}

func (this StringMap) Get(key string, defaultVal ...string) (val string, found bool) {
	if len(this) != 0 {
		if val, found = this[key]; found {
			return
		}
	}
	if len(defaultVal) != 0 {
		return defaultVal[0], false
	}
	return "", false
}

func (this StringMap) MustGet(key string, defaultVal ...string) string {
	v, _ := this.Get(key, defaultVal...)
	return v
}

func (this StringMap) GetInt(key string, defaultVal ...int) (val int, found bool, err error) {
	if len(this) != 0 {
		if v, ok := this[key]; ok {
			val, err := convertor.ToInt(v)
			if err != nil && len(defaultVal) != 0 {
				return defaultVal[0], true, err
			}
			return val, true, err
		}
	}
	if len(defaultVal) != 0 {
		return defaultVal[0], false, nil
	}
	return 0, false, nil
}

func (this StringMap) MustGetInt(key string, defaultVal ...int) int {
	val, _, _ := this.GetInt(key, defaultVal...)
	return val
}

func (this StringMap) GetInt64(key string, defaultVal ...int64) (val int64, found bool, err error) {
	if len(this) != 0 {
		if v, ok := this[key]; ok {
			val, err := convertor.ToInt64(v)
			if err != nil && len(defaultVal) != 0 {
				return defaultVal[0], true, err
			}
			return val, true, err
		}
	}
	if len(defaultVal) != 0 {
		return defaultVal[0], false, nil
	}
	return 0, false, nil
}

func (this StringMap) MustGetInt64(key string, defaultVal ...int64) int64 {
	val, _, _ := this.GetInt64(key, defaultVal...)
	return val
}

func (this StringMap) GetFloat(key string, defaultVal ...float64) (val float64, found bool, err error) {
	if len(this) != 0 {
		if v, ok := this[key]; ok {
			val, err := convertor.ToFloat(v)
			if err != nil && len(defaultVal) != 0 {
				return defaultVal[0], true, err
			}
			return val, true, err
		}
	}
	if len(defaultVal) != 0 {
		return defaultVal[0], false, nil
	}
	return 0, false, nil
}

func (this StringMap) MustGetFloat(key string, defaultVal ...float64) float64 {
	v, _, _ := this.GetFloat(key, defaultVal...)
	return v
}

func (this StringMap) GetBool(key string, defaultVal ...bool) (val bool, found bool, err error) {
	if len(this) != 0 {
		if v, ok := this[key]; ok {
			val, err := convertor.ToBool(v)
			if err != nil && len(defaultVal) != 0 {
				return defaultVal[0], true, err
			}
			return val, true, err
		}
	}
	if len(defaultVal) != 0 {
		return defaultVal[0], false, nil
	}
	return false, false, nil
}

func (this StringMap) MustGetBool(key string, defaultVal ...bool) bool {
	v, _, _ := this.GetBool(key, defaultVal...)
	return v
}

func (this StringMap) GetStringMap(key string) (val StringMap, found bool, err error) {
	if len(this) != 0 {
		var str string
		str, found = this[key]
		if str != "" {
			err = jsonUtil.UnmarshalFromString(str, &val)
		}
	}
	return
}

func (this StringMap) MustGetStringMap(key string) (val StringMap) {
	v, _, _ := this.GetStringMap(key)
	return v
}

func (this StringMap) GetStringObjectMap(key string) (val StringObjectMap, found bool, err error) {
	if len(this) != 0 {
		var str string
		str, found = this[key]
		if str != "" {
			err = jsonUtil.UnmarshalFromString(str, &val)
		}
	}
	return
}

func (this StringMap) MustGetStringObjectMap(key string) (val StringObjectMap) {
	v, _, _ := this.GetStringObjectMap(key)
	return v
}

func (this StringMap) GetToObject(key string, obj interface{}) (found bool, err error) {
	if len(this) != 0 {
		if v, ok := this[key]; ok && v != "" {
			return true, jsonUtil.UnmarshalFromString(v, obj)
		}
	}
	return
}

func (this StringMap) Mapping(f func(key, val string) (newKey, newVal string)) StringMap {
	if this != nil {
		result := make(StringMap, len(this))
		for k, v := range this {
			k2, v2 := f(k, v)
			result[k2] = v2
		}
		return result
	}
	return nil
}

func (this StringMap) Subset(keys ...string) StringMap {
	return this.SubsetF(func(key, _ string) bool {
		for _, item := range keys {
			if item == key {
				return true
			}
		}
		return false
	})
}

func (this StringMap) SubsetF(f func(key, val string) bool) StringMap {
	if this != nil {
		result := make(StringMap, len(this))
		for key, val := range this {
			if f(key, val) {
				result[key] = val
			}
		}
		return result
	}
	return nil
}

// 移除值为空的元素。
func (this StringMap) TrimEmptyValues() StringMap {
	return this.SubsetF(func(_, val string) bool { return val != "" })
}

func (this StringMap) ToLowerKeyMap() StringMap {
	if this != nil {
		a := make(StringMap, len(this))
		for k, v := range this {
			a[strings.ToLower(k)] = v
		}
		return a
	}
	return nil
}

func (this StringMap) ToStringObjectMap(f ...func(key, val string) (newKey string, newVal interface{})) StringObjectMap {
	if this != nil {
		var theF func(key, val string) (newKey string, newVal interface{})
		if len(f) != 0 && f[0] != nil {
			theF = f[0]
		}

		result := make(StringObjectMap, len(this))
		for k, v := range this {
			if theF == nil {
				result[k] = v
			} else {
				k2, v2 := theF(k, v)
				result[k2] = v2
			}
		}
		return result
	}
	return nil
}
