package mapUtil

import (
	"github.com/3th1nk/easygo/util/convertor"
	"github.com/3th1nk/easygo/util/jsonUtil"
	"strings"
)

type CaseInsensitiveStringMap map[string]string

func NewCaseInsensitiveStringMap(a ...map[string]string) CaseInsensitiveStringMap {
	dict := CaseInsensitiveStringMap{}
	dict.SetMulti(a...)
	return dict
}

func ToCaseInsensitiveStringMapSlice(src []map[string]string) (dest []CaseInsensitiveStringMap) {
	if src == nil {
		return nil
	}
	dest = make([]CaseInsensitiveStringMap, len(src))
	for i, v := range src {
		dest[i] = v
	}
	return
}

func (this CaseInsensitiveStringMap) Copy() CaseInsensitiveStringMap {
	return NewCaseInsensitiveStringMap(this)
}

func (this CaseInsensitiveStringMap) ToDB() (bytes []byte, err error) {
	if len(this) != 0 {
		return jsonUtil.Marshal(this)
	}
	return []byte{}, nil
}

func (this *CaseInsensitiveStringMap) FromDB(bytes []byte) error {
	return jsonUtil.Unmarshal(bytes, this)
}

func (this CaseInsensitiveStringMap) Delete(key ...string) {
	if len(this) != 0 {
		for _, k := range key {
			for s := range this {
				if strings.EqualFold(k, s) {
					delete(this, s)
					break
				}
			}
		}
	}
}

func (this CaseInsensitiveStringMap) Set(key, val string) {
	if _, found, realKey := this.tryGet(key); found {
		this[realKey] = val
	} else {
		this[key] = val
	}
}

func (this CaseInsensitiveStringMap) SetMulti(a ...map[string]string) {
	for _, item := range a {
		for key, val := range item {
			if _, found, realKey := this.tryGet(key); found {
				this[realKey] = val
			} else {
				this[key] = val
			}
		}
	}
}

func (this CaseInsensitiveStringMap) Keys() (arr []string) {
	if this != nil {
		arr = make([]string, 0, len(this))
		n := 0
		for key, _ := range this {
			arr[n] = key
			n++
		}
	}
	return
}

func (this CaseInsensitiveStringMap) Values() (arr []string) {
	if this != nil {
		arr = make([]string, 0, len(this))
		idx := 0
		for _, val := range this {
			arr[idx] = val
			idx++
		}
	}
	return
}

func (this CaseInsensitiveStringMap) tryGet(key string) (v string, found bool, realKey string) {
	if this != nil {
		if v, found := this[key]; found {
			return v, true, key
		}
		for s, val := range this {
			if strings.EqualFold(s, key) {
				return val, true, s
			}
		}
	}
	return
}

func (this CaseInsensitiveStringMap) Contains(key string) bool {
	_, ok, _ := this.tryGet(key)
	return ok
}

func (this CaseInsensitiveStringMap) ContainsAny(keys ...string) bool {
	if len(this) != 0 {
		for _, s := range keys {
			if _, ok, _ := this.tryGet(s); ok {
				return true
			}
		}
	}
	return false
}

func (this CaseInsensitiveStringMap) ContainsAll(keys []string) bool {
	if len(this) != 0 {
		for _, s := range keys {
			if _, ok, _ := this.tryGet(s); !ok {
				return false
			}
		}
		return true
	}
	return false
}

func (this CaseInsensitiveStringMap) Get(key string, defaultVal ...string) (val string, found bool) {
	if val, found, _ = this.tryGet(key); found {
		return
	} else if len(defaultVal) != 0 {
		return defaultVal[0], false
	}
	return "", false
}

func (this CaseInsensitiveStringMap) MustGet(key string, defaultVal ...string) string {
	v, _ := this.Get(key, defaultVal...)
	return v
}

func (this CaseInsensitiveStringMap) GetInt(key string, defaultVal ...int) (val int, found bool, err error) {
	if v, found, _ := this.tryGet(key); found {
		val, err := convertor.ToInt(v)
		if err != nil && len(defaultVal) != 0 {
			return defaultVal[0], true, err
		}
		return val, true, err
	} else if len(defaultVal) != 0 {
		return defaultVal[0], false, nil
	}
	return 0, false, nil
}

func (this CaseInsensitiveStringMap) MustGetInt(key string, defaultVal ...int) int {
	v, _, _ := this.GetInt(key, defaultVal...)
	return v
}

func (this CaseInsensitiveStringMap) GetInt64(key string, defaultVal ...int64) (val int64, found bool, err error) {
	if v, ok, _ := this.tryGet(key); ok {
		val, err := convertor.ToInt64(v)
		if err != nil && len(defaultVal) != 0 {
			return defaultVal[0], true, err
		}
		return val, true, err
	} else if len(defaultVal) != 0 {
		return defaultVal[0], false, nil
	}
	return 0, false, nil
}

func (this CaseInsensitiveStringMap) MustGetInt64(key string, defaultVal ...int64) int64 {
	val, _, _ := this.GetInt64(key, defaultVal...)
	return val
}

func (this CaseInsensitiveStringMap) GetFloat(key string, defaultVal ...float64) (val float64, found bool, err error) {
	if v, ok, _ := this.tryGet(key); ok {
		val, err := convertor.ToFloat(v)
		if err != nil && len(defaultVal) != 0 {
			return defaultVal[0], true, err
		}
		return val, true, err
	} else if len(defaultVal) != 0 {
		return defaultVal[0], false, nil
	}
	return 0, false, nil
}

func (this CaseInsensitiveStringMap) MustGetFloat(key string, defaultVal ...float64) float64 {
	v, _, _ := this.GetFloat(key, defaultVal...)
	return v
}

func (this CaseInsensitiveStringMap) GetBool(key string, defaultVal ...bool) (val bool, found bool, err error) {
	if v, ok, _ := this.tryGet(key); ok {
		val, err := convertor.ToBool(v)
		if err != nil && len(defaultVal) != 0 {
			return defaultVal[0], true, err
		}
		return val, true, err
	} else if len(defaultVal) != 0 {
		return defaultVal[0], false, nil
	}
	return false, false, nil
}

func (this CaseInsensitiveStringMap) MustGetBool(key string, defaultVal ...bool) bool {
	v, _, _ := this.GetBool(key, defaultVal...)
	return v
}

func (this CaseInsensitiveStringMap) GetToObject(key string, obj interface{}) (found bool, err error) {
	if v, ok, _ := this.tryGet(key); ok && v != "" {
		return true, jsonUtil.UnmarshalFromString(v, obj)
	}
	return
}

func (this CaseInsensitiveStringMap) Mapping(f func(key, val string) (newKey, newVal string)) CaseInsensitiveStringMap {
	if this != nil {
		result := make(CaseInsensitiveStringMap, len(this))
		for k, v := range this {
			k2, v2 := f(k, v)
			result.Set(k2, v2)
		}
		return result
	}
	return nil
}

func (this CaseInsensitiveStringMap) Subset(keys ...string) CaseInsensitiveStringMap {
	return this.SubsetF(func(key, _ string) bool {
		for _, item := range keys {
			if strings.EqualFold(key, item) {
				return true
			}
		}
		return false
	})
}

func (this CaseInsensitiveStringMap) SubsetF(f func(key, val string) bool) CaseInsensitiveStringMap {
	if this != nil {
		result := make(CaseInsensitiveStringMap, len(this))
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
func (this CaseInsensitiveStringMap) TrimEmptyValues() CaseInsensitiveStringMap {
	return this.SubsetF(func(_, val string) bool { return val != "" })
}

func (this CaseInsensitiveStringMap) ToLowerKeyMap() StringMap {
	if this != nil {
		a := make(StringMap, len(this))
		for k, v := range this {
			a[strings.ToLower(k)] = v
		}
		return a
	}
	return nil
}

func (this CaseInsensitiveStringMap) ToStringObjectMap(f ...func(key, val string) (newKey string, newVal interface{})) StringObjectMap {
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
