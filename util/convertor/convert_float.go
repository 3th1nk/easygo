package convertor

import (
	"fmt"
	"github.com/modern-go/reflect2"
	"strconv"
	"strings"
)

func ToFloat(a interface{}) (float64, error) {
	if reflect2.IsNil(a) {
		return 0, nil
	}

	switch t := a.(type) {
	case string:
		return StrToFloat(t)
	case *string:
		return StrToFloat(*t)
	case bool:
		return BoolToFloat(t), nil
	case *bool:
		return BoolToFloat(*t), nil
	case float64:
		return t, nil
	case *float64:
		return *t, nil
	case float32:
		return float64(t), nil
	case *float32:
		return float64(*t), nil
	case int:
		return float64(t), nil
	case *int:
		return float64(*t), nil
	case int64:
		return float64(t), nil
	case *int64:
		return float64(*t), nil
	case int32:
		return float64(t), nil
	case *int32:
		return float64(*t), nil
	case int16:
		return float64(t), nil
	case *int16:
		return float64(*t), nil
	case int8:
		return float64(t), nil
	case *int8:
		return float64(*t), nil
	case uint:
		return float64(t), nil
	case *uint:
		return float64(*t), nil
	case uint64:
		return float64(t), nil
	case *uint64:
		return float64(*t), nil
	case uint32:
		return float64(t), nil
	case *uint32:
		return float64(*t), nil
	case uint16:
		return float64(t), nil
	case *uint16:
		return float64(*t), nil
	case uint8:
		return float64(t), nil
	case *uint8:
		return float64(*t), nil
	}

	t, reflectType, reflectValue := GetBasicType(a)
	switch t {
	case BasicType_Bool:
		if reflectValue.Bool() {
			return 1, nil
		} else {
			return 0, nil
		}
	case BasicType_Int:
		return float64(reflectValue.Int()), nil
	case BasicType_Uint:
		return float64(reflectValue.Uint()), nil
	case BasicType_Float:
		return reflectValue.Float(), nil
	case BasicType_String:
		return StrToFloat(a.(string))
	default:
		return 0, fmt.Errorf("can't convert %s(%v) to float", strings.Trim(reflectType.PkgPath()+"."+reflectType.Name(), "."), ToStringNoError(a))
	}
}

func ToFloatNoError(a interface{}) float64 {
	v, _ := ToFloat(a)
	return v
}

func ToFloat32(a interface{}) (float32, error) {
	v, err := ToFloat(a)
	return float32(v), err
}

func ToFloat32NoError(a interface{}) float32 {
	v, _ := ToFloat(a)
	return float32(v)
}

func StrToFloat(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	} else if v, err := strconv.ParseFloat(s, 64); err == nil {
		return v, nil
	}
	return 0, fmt.Errorf("can't convert string(%s) to float", s)
}

func StrToFloatNoError(s string) float64 {
	v, _ := StrToFloat(s)
	return v
}

func BoolToFloat(v bool) float64 {
	if v {
		return 1
	}
	return 0
}
