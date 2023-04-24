package strUtil

import (
	"fmt"
	jsonIter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	jsonApi = jsonIter.Config{
		EscapeHTML:              true,
		MarshalFloatWith6Digits: true,
		SortMapKeys:             true,
	}.Froze()
)

func JoinInt(a []int, sep string, f ...func(i int, v int) string) string {
	var theF func(i int, v int) string
	if len(f) != 0 && f[0] != nil {
		theF = f[0]
	} else {
		theF = func(i int, v int) string {
			return strconv.FormatInt(int64(v), 10)
		}
	}

	arr := make([]string, len(a))
	for i, s := range a {
		arr[i] = theF(i, s)
	}
	return strings.Join(arr, sep)
}

func JoinInt64(a []int64, sep string, f ...func(i int, v int64) string) string {
	var theF func(i int, v int64) string
	if len(f) != 0 && f[0] != nil {
		theF = f[0]
	} else {
		theF = func(i int, v int64) string {
			return strconv.FormatInt(v, 10)
		}
	}

	arr := make([]string, len(a))
	for i, s := range a {
		arr[i] = theF(i, s)
	}
	return strings.Join(arr, sep)
}

func JoinInt32(a []int32, sep string, f ...func(i int, v int32) string) string {
	var theF func(i int, v int32) string
	if len(f) != 0 && f[0] != nil {
		theF = f[0]
	} else {
		theF = func(i int, v int32) string {
			return strconv.FormatInt(int64(v), 10)
		}
	}

	arr := make([]string, len(a))
	for i, s := range a {
		arr[i] = theF(i, s)
	}
	return strings.Join(arr, sep)
}

func JoinStr(a []string, sep string, f func(i int, v string) string) string {
	arr := make([]string, len(a))
	for i, s := range a {
		arr[i] = f(i, s)
	}
	return strings.Join(arr, sep)
}

func Join(slice interface{}, sep string, f ...func(i int) string) string {
	if len(f) == 0 || f[0] == nil {
		return JoinIf(slice, sep, nil)
	} else {
		return JoinIf(slice, sep, func(i int) (s string, b bool) {
			return f[0](i), true
		})
	}
}

func JoinIf(slice interface{}, sep string, f func(i int) (string, bool)) string {
	if reflect2.IsNil(slice) {
		return ""
	}

	reflectVal := reflect.ValueOf(slice)
	if kind := reflectVal.Kind(); kind != reflect.Array && kind != reflect.Slice {
		panic("参数 alice 必须是切片类型")
	}

	if f == nil {
		f = func(i int) (string, bool) {
			a := reflectVal.Index(i).Interface()
			if reflect2.IsNil(a) {
				return "", true
			}
			switch t := a.(type) {
			case []byte:
				return string(t), true
			case string:
				return t, true
			case bool:
				return strconv.FormatBool(t), true
			case time.Time:
				return t.Format("2006-01-02 15:04:05"), true
			default:
				if i, ok := a.(fmt.Stringer); ok {
					return i.String(), true
				} else {
					if str, err := jsonApi.MarshalToString(a); !reflect2.IsNil(err) {
						return "", true
					} else if str != "nil" && str != "[]" && str != "{}" {
						return str, true
					}
					return "", true
				}
			}
		}
	}

	valLen := reflectVal.Len()
	arr := make([]string, 0, valLen)
	for i := 0; i < valLen; i++ {
		if str, ok := f(i); ok {
			arr = append(arr, str)
		}
	}
	return strings.Join(arr, sep)
}
