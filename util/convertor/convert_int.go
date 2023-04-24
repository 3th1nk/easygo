package convertor

import (
	"fmt"
	"github.com/modern-go/reflect2"
	"strings"
)

func ToInt64(a interface{}) (int64, error) {
	if reflect2.IsNil(a) {
		return 0, nil
	}

	switch t := a.(type) {
	case string:
		return StrToInt64(t)
	case *string:
		return StrToInt64(*t)
	case bool:
		return int64(BoolToInt(t)), nil
	case *bool:
		return int64(BoolToInt(*t)), nil
	case float64:
		return int64(t), nil
	case *float64:
		return int64(*t), nil
	case float32:
		return int64(t), nil
	case *float32:
		return int64(*t), nil
	case int:
		return int64(t), nil
	case *int:
		return int64(*t), nil
	case int64:
		return t, nil
	case *int64:
		return *t, nil
	case int32:
		return int64(t), nil
	case *int32:
		return int64(*t), nil
	case int16:
		return int64(t), nil
	case *int16:
		return int64(*t), nil
	case int8:
		return int64(t), nil
	case *int8:
		return int64(*t), nil
	case uint:
		return int64(t), nil
	case *uint:
		return int64(*t), nil
	case uint64:
		return int64(t), nil
	case *uint64:
		return int64(*t), nil
	case uint32:
		return int64(t), nil
	case *uint32:
		return int64(*t), nil
	case uint16:
		return int64(t), nil
	case *uint16:
		return int64(*t), nil
	case uint8:
		return int64(t), nil
	case *uint8:
		return int64(*t), nil
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
		return reflectValue.Int(), nil
	case BasicType_Uint:
		return int64(reflectValue.Uint()), nil
	case BasicType_Float:
		return int64(reflectValue.Float()), nil
	case BasicType_String:
		return StrToInt64(a.(string))
	default:
		return 0, fmt.Errorf("can't convert %s(%v) to int", strings.Trim(reflectType.PkgPath()+"."+reflectType.Name(), "."), ToStringNoError(a))
	}
}

func ToUint64(a interface{}) (uint64, error) {
	if reflect2.IsNil(a) {
		return 0, nil
	}

	switch t := a.(type) {
	case string:
		return StrToUint64(t)
	case *string:
		return StrToUint64(*t)
	case bool:
		return uint64(BoolToInt(t)), nil
	case *bool:
		return uint64(BoolToInt(*t)), nil
	case float64:
		return uint64(t), nil
	case *float64:
		return uint64(*t), nil
	case float32:
		return uint64(t), nil
	case *float32:
		return uint64(*t), nil
	case int:
		return uint64(t), nil
	case *int:
		return uint64(*t), nil
	case int64:
		return uint64(t), nil
	case *int64:
		return uint64(*t), nil
	case int32:
		return uint64(t), nil
	case *int32:
		return uint64(*t), nil
	case int16:
		return uint64(t), nil
	case *int16:
		return uint64(*t), nil
	case int8:
		return uint64(t), nil
	case *int8:
		return uint64(*t), nil
	case uint:
		return uint64(t), nil
	case *uint:
		return uint64(*t), nil
	case uint64:
		return t, nil
	case *uint64:
		return *t, nil
	case uint32:
		return uint64(t), nil
	case *uint32:
		return uint64(*t), nil
	case uint16:
		return uint64(t), nil
	case *uint16:
		return uint64(*t), nil
	case uint8:
		return uint64(t), nil
	case *uint8:
		return uint64(*t), nil
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
		return uint64(reflectValue.Int()), nil
	case BasicType_Uint:
		return reflectValue.Uint(), nil
	case BasicType_Float:
		return uint64(reflectValue.Float()), nil
	case BasicType_String:
		return StrToUint64(a.(string))
	default:
		return 0, fmt.Errorf("can't convert %s(%v) to uint", strings.Trim(reflectType.PkgPath()+"."+reflectType.Name(), "."), ToStringNoError(a))
	}
}

func ToInt64NoError(a interface{}) int64 {
	v, _ := ToInt64(a)
	return v
}

func ToUint64NoError(a interface{}) uint64 {
	v, _ := ToUint64(a)
	return v
}

func ToInt(a interface{}) (int, error) {
	v, err := ToInt64(a)
	return int(v), err
}

func ToIntNoError(a interface{}) int {
	v, _ := ToInt64(a)
	return int(v)
}

func ToUint(a interface{}) (uint, error) {
	v, err := ToUint64(a)
	return uint(v), err
}

func ToUintNoError(a interface{}) uint {
	v, _ := ToUint64(a)
	return uint(v)
}

func ToInt32(a interface{}) (int32, error) {
	v, err := ToInt64(a)
	return int32(v), err
}

func ToInt32NoError(a interface{}) int32 {
	v, _ := ToInt64(a)
	return int32(v)
}

func ToUint32(a interface{}) (uint32, error) {
	v, err := ToUint64(a)
	return uint32(v), err
}

func ToUint32NoError(a interface{}) uint32 {
	v, _ := ToUint64(a)
	return uint32(v)
}
