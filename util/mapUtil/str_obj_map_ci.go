package mapUtil

import (
	"github.com/3th1nk/easygo/util/convertor"
	"github.com/3th1nk/easygo/util/jsonUtil"
	"github.com/modern-go/reflect2"
	"strings"
)

type CaseInsensitiveStringObjectMap map[string]interface{}

func NewCaseInsensitiveStringObjectMap(a ...map[string]interface{}) CaseInsensitiveStringObjectMap {
	dict := CaseInsensitiveStringObjectMap{}
	dict.SetMulti(a...)
	return dict
}

func ToCaseInsensitiveStringObjectMapSlice(src []map[string]interface{}) (dest []CaseInsensitiveStringObjectMap) {
	if src == nil {
		return nil
	}
	dest = make([]CaseInsensitiveStringObjectMap, len(src))
	for i, v := range src {
		dest[i] = v
	}
	return
}

func (this CaseInsensitiveStringObjectMap) Copy() CaseInsensitiveStringObjectMap {
	return NewCaseInsensitiveStringObjectMap(this)
}

func (this CaseInsensitiveStringObjectMap) ToDB() (bytes []byte, err error) {
	if len(this) != 0 {
		return jsonUtil.Marshal(this)
	}
	return []byte{}, nil
}

func (this *CaseInsensitiveStringObjectMap) FromDB(bytes []byte) error {
	return jsonUtil.Unmarshal(bytes, this)
}

func (this CaseInsensitiveStringObjectMap) Delete(key ...string) {
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

func (this CaseInsensitiveStringObjectMap) Set(key string, val interface{}) {
	if _, found, realKey := this.tryGet(key); found {
		this[realKey] = val
	} else {
		this[key] = val
	}
}

func (this CaseInsensitiveStringObjectMap) SetMulti(a ...map[string]interface{}) {
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

func (this CaseInsensitiveStringObjectMap) Keys() (arr []string) {
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

func (this CaseInsensitiveStringObjectMap) Values() (arr []interface{}) {
	if this != nil {
		arr = make([]interface{}, 0, len(this))
		n := 0
		for _, val := range this {
			arr[n] = val
			n++
		}
	}
	return
}

func (this CaseInsensitiveStringObjectMap) tryGet(key string) (v interface{}, found bool, realKey string) {
	if len(this) != 0 {
		if v, found := this[key]; found {
			return v, true, key
		}
	}
	for s, val := range this {
		if strings.EqualFold(s, key) {
			return val, true, s
		}
	}
	return
}

func (this CaseInsensitiveStringObjectMap) Contains(key string) bool {
	_, ok, _ := this.tryGet(key)
	return ok
}

func (this CaseInsensitiveStringObjectMap) ContainsAny(keys ...string) bool {
	if len(this) != 0 {
		for _, s := range keys {
			if _, ok, _ := this.tryGet(s); ok {
				return true
			}
		}
	}
	return false
}

func (this CaseInsensitiveStringObjectMap) ContainsAll(keys []string) bool {
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

func (this CaseInsensitiveStringObjectMap) Get(key string, defaultVal ...interface{}) (val interface{}, found bool) {
	if val, found, _ = this.tryGet(key); found {
		return
	} else if len(defaultVal) != 0 {
		return defaultVal[0], false
	}
	return nil, false
}

func (this CaseInsensitiveStringObjectMap) MustGet(key string, defaultVal ...interface{}) interface{} {
	v, _ := this.Get(key, defaultVal...)
	return v
}

func (this CaseInsensitiveStringObjectMap) GetString(key string, defaultVal ...string) (val string, found bool, err error) {
	if v, ok, _ := this.tryGet(key); ok {
		val, err := convertor.ToString(v)
		if err != nil && len(defaultVal) != 0 {
			return defaultVal[0], true, err
		}
		return val, true, err
	} else if len(defaultVal) != 0 {
		return defaultVal[0], false, nil
	}
	return "", false, nil
}

func (this CaseInsensitiveStringObjectMap) MustGetString(key string, defaultVal ...string) string {
	v, _, _ := this.GetString(key, defaultVal...)
	return v
}

func (this CaseInsensitiveStringObjectMap) GetInt(key string, defaultVal ...int) (val int, found bool, err error) {
	if v, ok, _ := this.tryGet(key); ok {
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

func (this CaseInsensitiveStringObjectMap) MustGetInt(key string, defaultVal ...int) int {
	v, _, _ := this.GetInt(key, defaultVal...)
	return v
}

func (this CaseInsensitiveStringObjectMap) GetInt64(key string, defaultVal ...int64) (val int64, found bool, err error) {
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

func (this CaseInsensitiveStringObjectMap) MustGetInt64(key string, defaultVal ...int64) int64 {
	val, _, _ := this.GetInt64(key, defaultVal...)
	return val
}

func (this CaseInsensitiveStringObjectMap) GetFloat(key string, defaultVal ...float64) (val float64, found bool, err error) {
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

func (this CaseInsensitiveStringObjectMap) MustGetFloat(key string, defaultVal ...float64) float64 {
	v, _, _ := this.GetFloat(key, defaultVal...)
	return v
}

func (this CaseInsensitiveStringObjectMap) GetBool(key string, defaultVal ...bool) (val bool, found bool, err error) {
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

func (this CaseInsensitiveStringObjectMap) MustGetBool(key string, defaultVal ...bool) bool {
	v, _, _ := this.GetBool(key, defaultVal...)
	return v
}

func (this CaseInsensitiveStringObjectMap) GetToObject(key string, obj interface{}) (found bool, err error) {
	if v, ok, _ := this.tryGet(key); ok && v != nil {
		str, err := convertor.ToString(v)
		if err != nil {
			return true, err
		} else if str != "" {
			err = jsonUtil.UnmarshalFromString(str, obj)
			return true, err
		}
	}
	return
}

func (this CaseInsensitiveStringObjectMap) Mapping(f func(key string, val interface{}) (newKey string, newVal interface{})) CaseInsensitiveStringObjectMap {
	if len(this) != 0 {
		result := make(CaseInsensitiveStringObjectMap, len(this))
		for k, v := range this {
			k2, v2 := f(k, v)
			result.Set(k2, v2)
		}
		return result
	}
	return nil
}

func (this CaseInsensitiveStringObjectMap) Subset(keys ...string) CaseInsensitiveStringObjectMap {
	return this.SubsetF(func(key string, _ interface{}) bool {
		for _, item := range keys {
			if strings.EqualFold(key, item) {
				return true
			}
		}
		return false
	})
}

func (this CaseInsensitiveStringObjectMap) SubsetF(f func(key string, val interface{}) bool) CaseInsensitiveStringObjectMap {
	if len(this) != 0 {
		result := make(CaseInsensitiveStringObjectMap, len(this))
		for key, val := range this {
			if f(key, val) {
				result[key] = val
			}
		}
		return result
	}
	return nil
}

func (this CaseInsensitiveStringObjectMap) TrimNilValues() CaseInsensitiveStringObjectMap {
	return this.SubsetF(func(_ string, val interface{}) bool {
		return !reflect2.IsNil(val)
	})
}

// 移除值为空的元素。
//   falseAsEmpty: 是否把 false 当作空值。默认为 true。
func (this CaseInsensitiveStringObjectMap) TrimEmptyValues(falseAsEmpty ...bool) CaseInsensitiveStringObjectMap {
	return this.SubsetF(func(key string, val interface{}) bool {
		return !convertor.IsEmpty2(val, falseAsEmpty...)
	})
}

func (this CaseInsensitiveStringObjectMap) ToLowerKeyMap() StringObjectMap {
	if this != nil {
		a := make(StringObjectMap, len(this))
		for k, v := range this {
			a[strings.ToLower(k)] = v
		}
		return a
	}
	return nil
}

func (this CaseInsensitiveStringObjectMap) ToStringMap(f ...func(key string, val interface{}) (newKey, newVal string)) StringMap {
	if len(this) != 0 {
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
