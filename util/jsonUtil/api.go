// ------------------------------------------------------------------------------
// 对 jsonIter 的进一步封装
//   1、默认注册弱类型解释器，支持将 string("123") 解析为 number(123)、将 string("true")|number(!0) 解析为 bool(true) 等。
//   2、time.Time 默认序列化为 yyyy-MM-dd HH:MM:SS 格式
//   3、序列化 Map 时默认对 Key 进行排序
//   4、简化接口方法
// ------------------------------------------------------------------------------
package jsonUtil

import (
	"bytes"
	jsonIter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
)

func NewApi(config *jsonIter.Config) *Api {
	return &Api{API: config.Froze(), cfg: config}
}

func DefaultApi() *Api {
	return defaultApi
}

func SetDefaultApi(api *Api) {
	defaultApi = api
}

func SortMapKeysApi(sort ...bool) *Api {
	if len(sort) == 0 || sort[0] {
		return sortMapKeysApi
	}
	return unSortMapKeysApi
}

func (this *Api) IndentApi(indentAndPrefix ...string) *Api {
	api := this.Clone()
	if len(indentAndPrefix) == 0 {
		api.indent, api.prefix = "    ", ""
	} else if len(indentAndPrefix) == 1 {
		api.indent, api.prefix = indentAndPrefix[0], ""
	} else {
		api.indent, api.prefix = indentAndPrefix[0], indentAndPrefix[1]
	}
	return api
}

var (
	unSortMapKeysApi = NewApi(&jsonIter.Config{
		EscapeHTML:              false,
		MarshalFloatWith6Digits: true,
	})

	sortMapKeysApi = NewApi(&jsonIter.Config{
		EscapeHTML:              false,
		MarshalFloatWith6Digits: true,
		SortMapKeys:             true,
	})

	defaultApi = sortMapKeysApi
)

// ------------------------------------------------------------------------------ Api
type Api struct {
	jsonIter.API
	indent string
	prefix string
	cfg    *jsonIter.Config
}

func (this *Api) Clone() *Api { return NewApi(this.cfg) }

func (this *Api) Get(data []byte, path ...interface{}) jsonIter.Any {
	return this.API.Get(data, path...)
}

func (this *Api) GetString(str string, path ...interface{}) jsonIter.Any {
	return this.API.Get([]byte(str), path...)
}

func (this *Api) Unmarshal(data []byte, v interface{}) error {
	if len(data) != 0 {
		return this.API.Unmarshal(data, v)
	}
	return nil
}

func (this *Api) UnmarshalFromString(str string, v interface{}) error {
	if len(str) != 0 {
		return this.API.UnmarshalFromString(str, v)
	}
	return nil
}

func (this *Api) UnmarshalFromObject(src, dest interface{}) (err error) {
	return this.doUnmarshalFromObject(src, dest, nil)
}

func (this *Api) doUnmarshalFromObject(src, dest interface{}, f func(srcData []byte, dest interface{}) error) (err error) {
	if reflect2.IsNil(src) {
		return nil
	}

	var data []byte
	switch t := src.(type) {
	case []byte:
		if len(t) == 0 {
			return nil
		}
		data = t
	case string:
		if len(t) == 0 {
			return nil
		}
		data = []byte(t)
	default:
		if data, err = this.Marshal(src); err != nil {
			return err
		}
	}

	if f != nil {
		if err = f(data, dest); err != nil {
			return err
		}
	}

	return this.Unmarshal(data, dest)
}

func (this *Api) Marshal(v interface{}) ([]byte, error) {
	return this.doEncode(v, this.indent, this.prefix)
}

func (this *Api) MustMarshal(v interface{}) []byte {
	rtn, _ := this.Marshal(v)
	return rtn
}

func (this *Api) MarshalIndent(v interface{}, indentAndPrefix ...string) ([]byte, error) {
	var indent, prefix string
	if len(indentAndPrefix) == 0 {
		indent, prefix = "    ", ""
	} else if len(indentAndPrefix) == 1 {
		indent, prefix = indentAndPrefix[0], ""
	} else {
		indent, prefix = indentAndPrefix[0], indentAndPrefix[1]
	}
	return this.doEncode(v, indent, prefix)
}

func (this *Api) doEncode(v interface{}, indent, prefix string) ([]byte, error) {
	if indent != "" || prefix != "" {
		buf, err := this.API.MarshalIndent(v, "", indent)
		if err != nil {
			return nil, err
		}
		if prefix != "" {
			buf = bytes.ReplaceAll(buf, []byte{'\n'}, append([]byte{'\n'}, []byte(prefix)...))
			buf = append([]byte(prefix), buf...)
		}
		return buf, nil
	}
	return this.API.Marshal(v)
}

func (this *Api) MustMarshalIndent(v interface{}, indentAndPrefix ...string) []byte {
	rtn, _ := this.MarshalIndent(v, indentAndPrefix...)
	return rtn
}

func (this *Api) MarshalToString(v interface{}) (string, error) {
	if b, err := this.API.MarshalToString(v); err != nil {
		return "", err
	} else {
		return string(b), nil
	}
}

func (this *Api) MustMarshalToString(v interface{}) string {
	rtn, _ := this.MarshalToString(v)
	return rtn
}

func (this *Api) MarshalToStringIndent(v interface{}, indentAndPrefix ...string) (string, error) {
	if val, err := this.MarshalIndent(v, indentAndPrefix...); err != nil {
		return "", err
	} else {
		return string(val), nil
	}
}

func (this *Api) MustMarshalToStringIndent(v interface{}, indentAndPrefix ...string) string {
	if val, err := this.MarshalIndent(v, indentAndPrefix...); err != nil {
		return ""
	} else {
		return string(val)
	}
}
