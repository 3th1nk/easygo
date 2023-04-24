package mapUtil

import (
	"github.com/3th1nk/easygo/util/convertor"
	"github.com/3th1nk/easygo/util/jsonUtil"
	"github.com/modern-go/reflect2"
	"sort"
	"strings"
)

type StringObjectMap map[string]interface{}

func NewStringObjectMap(a ...map[string]interface{}) StringObjectMap {
	dict := StringObjectMap{}
	dict.SetMulti(a...)
	return dict
}

func ToStringObjectMapSlice(src []map[string]interface{}) (dest []StringObjectMap) {
	if src == nil {
		return nil
	}
	dest = make([]StringObjectMap, len(src))
	for i, v := range src {
		dest[i] = v
	}
	return
}

func (this StringObjectMap) Copy() StringObjectMap {
	return NewStringObjectMap(this)
}

func (this StringObjectMap) ToDB() (bytes []byte, err error) {
	if len(this) != 0 {
		return jsonUtil.Marshal(this)
	}
	return []byte{}, nil
}

func (this *StringObjectMap) FromDB(bytes []byte) error {
	return jsonUtil.Unmarshal(bytes, this)
}

func (this StringObjectMap) Delete(key ...string) {
	if len(this) != 0 {
		for _, k := range key {
			delete(this, k)
		}
	}
}

func (this StringObjectMap) Set(key string, val interface{}) {
	this[key] = val
}

func (this StringObjectMap) SetMulti(a ...map[string]interface{}) {
	for _, item := range a {
		for key, val := range item {
			this[key] = val
		}
	}
}

func (this StringObjectMap) Keys() (arr []string) {
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

func (this StringObjectMap) Values() (arr []interface{}) {
	if this != nil {
		arr = make([]interface{}, len(this))
		n := 0
		for _, val := range this {
			arr[n] = val
			n++
		}
	}
	return
}

func (this StringObjectMap) ValuesSortedByKey() (arr []interface{}) {
	if this != nil {
		keys := this.Keys()
		sort.Strings(keys)
		arr = make([]interface{}, len(this))
		n := 0
		for _, k := range keys {
			arr[n] = this.MustGet(k)
			n++
		}
	}
	return
}

func (this StringObjectMap) Contains(key string) bool {
	if len(this) != 0 {
		_, ok := this[key]
		return ok
	}
	return false
}

func (this StringObjectMap) ContainsAny(keys ...string) bool {
	if len(this) != 0 {
		for _, s := range keys {
			if _, ok := this[s]; ok {
				return true
			}
		}
	}
	return false
}

func (this StringObjectMap) ContainsAll(keys []string) bool {
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

func (this StringObjectMap) Get(key string, defaultVal ...interface{}) (val interface{}, found bool) {
	if len(this) != 0 {
		if val, found = this[key]; found {
			return
		}
	}
	if len(defaultVal) != 0 {
		return defaultVal[0], false
	}
	return nil, false
}

func (this StringObjectMap) MustGet(key string, defaultVal ...interface{}) interface{} {
	v, _ := this.Get(key, defaultVal...)
	return v
}

func (this StringObjectMap) GetString(key string, defaultVal ...string) (val string, found bool, err error) {
	if len(this) != 0 {
		if v, ok := this[key]; ok {
			val, err := convertor.ToString(v)
			if err != nil && len(defaultVal) != 0 {
				return defaultVal[0], true, err
			}
			return val, true, err
		}
	}
	if len(defaultVal) != 0 {
		return defaultVal[0], false, nil
	}
	return "", false, nil
}

func (this StringObjectMap) MustGetString(key string, defaultVal ...string) string {
	v, _, _ := this.GetString(key, defaultVal...)
	return v
}

func (this StringObjectMap) GetInt(key string, defaultVal ...int) (val int, found bool, err error) {
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

func (this StringObjectMap) MustGetInt(key string, defaultVal ...int) int {
	v, _, _ := this.GetInt(key, defaultVal...)
	return v
}

func (this StringObjectMap) GetInt64(key string, defaultVal ...int64) (val int64, found bool, err error) {
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

func (this StringObjectMap) MustGetInt64(key string, defaultVal ...int64) int64 {
	val, _, _ := this.GetInt64(key, defaultVal...)
	return val
}

func (this StringObjectMap) GetFloat(key string, defaultVal ...float64) (val float64, found bool, err error) {
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

func (this StringObjectMap) MustGetFloat(key string, defaultVal ...float64) float64 {
	V, _, _ := this.GetFloat(key, defaultVal...)
	return V
}

func (this StringObjectMap) GetBool(key string, defaultVal ...bool) (val bool, found bool, err error) {
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

func (this StringObjectMap) MustGetBool(key string, defaultVal ...bool) bool {
	V, _, _ := this.GetBool(key, defaultVal...)
	return V
}

func (this StringObjectMap) GetStringMap(key string) (val StringMap, found bool, err error) {
	if len(this) != 0 {
		var v interface{}
		v, found = this[key]
		if v != nil {
			if val, _ = v.(map[string]string); val != nil {
				return
			} else if m, _ := v.(map[string]interface{}); m != nil {
				val = StringObjectMap(m).ToStringMap()
				return
			} else {
				err = jsonUtil.UnmarshalFromObject(v, &val)
				return
			}
		}
	}
	return
}

func (this StringObjectMap) MustGetStringMap(key string) (val StringMap) {
	v, _, _ := this.GetStringMap(key)
	return v
}

func (this StringObjectMap) GetStringObjectMap(key string) (val StringObjectMap, found bool, err error) {
	if len(this) != 0 {
		var v interface{}
		v, found = this[key]
		if v != nil {
			if val, _ = v.(map[string]interface{}); val != nil {
				return
			} else if m, _ := v.(map[string]string); m != nil {
				val = StringMap(m).ToStringObjectMap()
				return
			} else {
				err = jsonUtil.UnmarshalFromObject(v, &val)
				return
			}
		}
	}
	return
}

func (this StringObjectMap) MustGetStringObjectMap(key string) (val StringObjectMap) {
	v, _, _ := this.GetStringObjectMap(key)
	return v
}

func (this StringObjectMap) GetToObject(key string, obj interface{}) (found bool, err error) {
	if len(this) != 0 {
		if v, ok := this[key]; ok && v != nil {
			str, err := convertor.ToString(v)
			if err != nil {
				return true, err
			} else if str != "" {
				err = jsonUtil.UnmarshalFromString(str, obj)
				return true, err
			}
		}
	}
	return
}

func (this StringObjectMap) Mapping(f func(key string, val interface{}) (newKey string, newVal interface{})) StringObjectMap {
	if this != nil {
		result := make(StringObjectMap, len(this))
		for k, v := range this {
			k2, v2 := f(k, v)
			result[k2] = v2
		}
		return result
	}
	return nil
}

func (this StringObjectMap) Subset(keys ...string) StringObjectMap {
	return this.SubsetF(func(key string, _ interface{}) bool {
		for _, item := range keys {
			if item == key {
				return true
			}
		}
		return false
	})
}

func (this StringObjectMap) SubsetF(f func(key string, val interface{}) bool) StringObjectMap {
	if this != nil {
		result := make(StringObjectMap, len(this))
		for key, val := range this {
			if f(key, val) {
				result[key] = val
			}
		}
		return result
	}
	return nil
}

func (this StringObjectMap) TrimNilValues() StringObjectMap {
	return this.SubsetF(func(_ string, val interface{}) bool {
		return !reflect2.IsNil(val)
	})
}

// 移除值为空的元素。
//   falseAsEmpty: 是否把 false 当作空值。默认为 true。
func (this StringObjectMap) TrimEmptyValues(falseAsEmpty ...bool) StringObjectMap {
	return this.SubsetF(func(key string, val interface{}) bool {
		return !convertor.IsEmpty2(val, falseAsEmpty...)
	})
}

func (this StringObjectMap) ToLowerKeyMap() StringObjectMap {
	if this != nil {
		a := make(StringObjectMap, len(this))
		for k, v := range this {
			a[strings.ToLower(k)] = v
		}
		return a
	}
	return nil
}

func (this StringObjectMap) ToStringMap(f ...func(key string, val interface{}) (newKey, newVal string)) StringMap {
	if this != nil {
		var theF func(key string, val interface{}) (newKey, newVal string)
		if len(f) != 0 && f[0] != nil {
			theF = f[0]
		} else {
			theF = func(key string, val interface{}) (newKey, newVal string) {
				return key, convertor.ToStringNoError(val)
			}
		}

		result := make(StringMap, len(this))
		for k, v := range this {
			k2, v2 := theF(k, v)
			result[k2] = v2
		}
		return result
	}
	return nil
}
