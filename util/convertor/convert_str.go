package convertor

import (
	"fmt"
	"github.com/3th1nk/easygo/util/timeUtil"
	"github.com/modern-go/reflect2"
	"strconv"
	"time"
)

type ConvertToString interface {
	ToString() (str string, err error)
}

func ToString(a interface{}, api ...JsonAPI) (string, error) {
	if reflect2.IsNil(a) {
		return "", nil
	}

	switch t := a.(type) {
	case []byte:
		return string(t), nil
	case *[]byte:
		return string(*t), nil
	case byte:
		return string([]byte{t}), nil
	case *byte:
		return string([]byte{*t}), nil
	case string:
		return t, nil
	case *string:
		return *t, nil
	case bool:
		return strconv.FormatBool(t), nil
	case *bool:
		return strconv.FormatBool(*t), nil
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64), nil
	case *float64:
		return strconv.FormatFloat(*t, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(t), 'f', -1, 64), nil
	case *float32:
		return strconv.FormatFloat(float64(*t), 'f', -1, 64), nil
	case int:
		return strconv.FormatInt(int64(t), 10), nil
	case *int:
		return strconv.FormatInt(int64(*t), 10), nil
	case int64:
		return strconv.FormatInt(t, 10), nil
	case *int64:
		return strconv.FormatInt(*t, 10), nil
	case int32:
		return strconv.FormatInt(int64(t), 10), nil
	case *int32:
		return strconv.FormatInt(int64(*t), 10), nil
	case int16:
		return strconv.FormatInt(int64(t), 10), nil
	case *int16:
		return strconv.FormatInt(int64(*t), 10), nil
	case int8:
		return strconv.FormatInt(int64(t), 10), nil
	case *int8:
		return strconv.FormatInt(int64(*t), 10), nil
	case uint:
		return strconv.FormatInt(int64(t), 10), nil
	case *uint:
		return strconv.FormatInt(int64(*t), 10), nil
	case uint64:
		return strconv.FormatInt(int64(t), 10), nil
	case *uint64:
		return strconv.FormatInt(int64(*t), 10), nil
	case uint32:
		return strconv.FormatInt(int64(t), 10), nil
	case *uint32:
		return strconv.FormatInt(int64(*t), 10), nil
	case uint16:
		return strconv.FormatInt(int64(t), 10), nil
	case *uint16:
		return strconv.FormatInt(int64(*t), 10), nil
	case time.Time:
		return t.Format("2006-01-02 15:04:05"), nil
	case *time.Time:
		return t.Format("2006-01-02 15:04:05"), nil
	case timeUtil.JsonTime:
		return t.Format("2006-01-02 15:04:05"), nil
	case *timeUtil.JsonTime:
		return t.Format("2006-01-02 15:04:05"), nil
	}

	if v, ok := a.(error); ok {
		return v.Error(), nil
	} else if v, ok := a.(ConvertToString); ok {
		return v.ToString()
	} else if i, ok := a.(fmt.Stringer); ok {
		return i.String(), nil
	}

	switch t, _, reflectValue := GetBasicType(a); t {
	case BasicType_Bool:
		return strconv.FormatBool(reflectValue.Bool()), nil
	case BasicType_Int:
		return strconv.FormatInt(reflectValue.Int(), 10), nil
	case BasicType_Uint:
		return strconv.FormatUint(reflectValue.Uint(), 10), nil
	case BasicType_Float:
		return strconv.FormatFloat(reflectValue.Float(), 'f', -1, 64), nil
	case BasicType_String:
		return reflectValue.String(), nil
	default:
		var theApi JsonAPI
		if len(api) != 0 && api[0] != nil {
			theApi = api[0]
		} else {
			theApi = jsonApi
		}
		if str, err := theApi.MarshalToString(a); err != nil {
			return "", err
		} else if str != "nil" && str != "[]" && str != "{}" {
			return str, nil
		}
		return "", nil
	}
}

func ToStringNoError(a interface{}, api ...JsonAPI) string {
	v, _ := ToString(a, api...)
	return v
}

func StrToInt64(s string) (int64, error) {
	if s == "" {
		return 0, nil
	} else if v, err := strconv.ParseInt(s, 10, 64); reflect2.IsNil(err) {
		return v, nil
	} else if v, err := strconv.ParseFloat(s, 64); reflect2.IsNil(err) {
		return int64(v), nil
	} else if v, err := strconv.ParseBool(s); reflect2.IsNil(err) {
		if v {
			return 1, nil
		}
		return 0, nil
	}
	return 0, fmt.Errorf("can't convert string(%s) to int", s)
}

func StrToInt64NoError(s string) int64 {
	v, _ := StrToInt64(s)
	return v
}

func StrToUint64(s string) (uint64, error) {
	if s == "" {
		return 0, nil
	} else if v, err := strconv.ParseUint(s, 10, 64); reflect2.IsNil(err) {
		return v, nil
	} else if v, err := strconv.ParseFloat(s, 64); reflect2.IsNil(err) {
		return uint64(v), nil
	} else if v, err := strconv.ParseBool(s); reflect2.IsNil(err) {
		if v {
			return 1, nil
		}
		return 0, nil
	}
	return 0, fmt.Errorf("can't convert string(%s) to uint", s)
}

func StrToUint64NoError(s string) uint64 {
	v, _ := StrToUint64(s)
	return v
}

func StrToInt(s string) (int, error) {
	v, err := StrToInt64(s)
	return int(v), err
}

func StrToIntNoError(s string) int {
	v, _ := StrToInt64(s)
	return int(v)
}

func StrToUint(s string) (uint, error) {
	v, err := StrToInt64(s)
	return uint(v), err
}

func StrToUintNoError(s string) uint {
	v, _ := StrToUint64(s)
	return uint(v)
}

func StrToInt32(s string) (int32, error) {
	v, err := StrToInt64(s)
	return int32(v), err
}

func StrToInt32NoError(s string) int32 {
	v, _ := StrToInt64(s)
	return int32(v)
}

func StrToUint32(s string) (uint32, error) {
	v, err := StrToInt64(s)
	return uint32(v), err
}

func StrToUint32NoError(s string) uint32 {
	v, _ := StrToUint64(s)
	return uint32(v)
}
