package convertor

import (
	"fmt"
	"github.com/modern-go/reflect2"
	"strconv"
	"strings"
)

func ToBool(a interface{}) (bool, error) {
	if reflect2.IsNil(a) {
		return false, nil
	}

	switch t := a.(type) {
	case string:
		return StrToBool(t)
	case *string:
		return StrToBool(*t)
	case bool:
		return t, nil
	case *bool:
		return *t, nil
	case float64:
		return t != 0, nil
	case *float64:
		return *t != 0, nil
	case float32:
		return t != 0, nil
	case *float32:
		return *t != 0, nil
	case int:
		return t != 0, nil
	case *int:
		return *t != 0, nil
	case int64:
		return t != 0, nil
	case *int64:
		return *t != 0, nil
	case int32:
		return t != 0, nil
	case *int32:
		return *t != 0, nil
	case int16:
		return t != 0, nil
	case *int16:
		return *t != 0, nil
	case int8:
		return t != 0, nil
	case *int8:
		return *t != 0, nil
	case uint:
		return t != 0, nil
	case *uint:
		return *t != 0, nil
	case uint64:
		return t != 0, nil
	case *uint64:
		return *t != 0, nil
	case uint32:
		return t != 0, nil
	case *uint32:
		return *t != 0, nil
	case uint16:
		return t != 0, nil
	case *uint16:
		return *t != 0, nil
	case uint8:
		return t != 0, nil
	case *uint8:
		return *t != 0, nil
	}

	t, reflectType, reflectValue := GetBasicType(a)
	switch t {
	case BasicType_Bool:
		return reflectValue.Bool(), nil
	case BasicType_Int:
		return reflectValue.Int() != 0, nil
	case BasicType_Uint:
		return reflectValue.Uint() != 0, nil
	case BasicType_Float:
		return reflectValue.Float() != 0, nil
	case BasicType_String:
		return StrToBool(a.(string))
	default:
		return false, fmt.Errorf("can't convert %s(%v) to bool", strings.Trim(reflectType.PkgPath()+"."+reflectType.Name(), "."), ToStringNoError(a))
	}
}

func ToBoolNoError(a interface{}) bool {
	v, _ := ToBool(a)
	return v
}

func StrToBool(s string) (bool, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	if s == "" {
		return false, nil
	} else if v, err := strconv.ParseBool(s); err == nil {
		return v, nil
	} else if v, err := strconv.ParseFloat(s, 64); err == nil {
		return v != 0, nil
	}
	return false, fmt.Errorf("can't convert string(%s) to bool", s)
}

func StrToBoolNoError(s string) bool {
	v, _ := StrToBool(s)
	return v
}

func BoolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
