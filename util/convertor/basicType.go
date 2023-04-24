package convertor

import (
	"fmt"
	"github.com/3th1nk/easygo/util/regexpUtil"
	jsonIter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"reflect"
	"strconv"
	"strings"
)

// BasicType 简化版的反射类型
type BasicType int

const (
	BasicType_Others  BasicType = -1
	BasicType_Invalid BasicType = 0
	BasicType_Nil     BasicType = 1
	BasicType_Bool    BasicType = 2
	BasicType_Int     BasicType = 3
	BasicType_Uint    BasicType = 4
	BasicType_Float   BasicType = 5
	BasicType_String  BasicType = 6
	BasicType_Slice   BasicType = 7
	BasicType_Map     BasicType = 8
	BasicType_Struct  BasicType = 9
)

func GetBasicType0(a interface{}, parseStrValue ...bool) BasicType {
	basicType, _, _ := GetBasicType(a, parseStrValue...)
	return basicType
}

// 获取变量的基本类型
//   parseStrValue: 如果 a 是字符串，是否继续解析字符串内容的值类型（使用了 GetStrValueType 方法）。例：
//                  GetBasicType("123", true) 将返回 BasicType_Int、GetBasicType("true", true) 将返回 BasicType_Bool。
func GetBasicType(a interface{}, parseStrValue ...bool) (BasicType, reflect.Type, reflect.Value) {
	refType := reflect.TypeOf(a)
	refValue := reflect.ValueOf(a)
	if reflect2.IsNil(a) {
		return BasicType_Nil, refType, refValue
	}

	kind := refType.Kind()
	if kind == reflect.Ptr {
		refType = refType.Elem()
		refValue = refValue.Elem()
	}
	switch reflect.TypeOf(a).Kind() {
	case reflect.Invalid:
		return BasicType_Invalid, refType, refValue
	case reflect.Bool:
		return BasicType_Bool, refType, refValue
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return BasicType_Int, refType, refValue
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return BasicType_Uint, refType, refValue
	case reflect.Float32, reflect.Float64:
		return BasicType_Float, refType, refValue
	case reflect.String:
		if len(parseStrValue) != 0 && parseStrValue[0] {
			a, b, c, _ := GetStrValueType(refValue.String())
			return a, b, c
		}
		return BasicType_String, refType, refValue
	case reflect.Slice, reflect.Array:
		return BasicType_Slice, refType, refValue
	case reflect.Map:
		return BasicType_Map, refType, refValue
	case reflect.Struct:
		return BasicType_Struct, refType, refValue
	default:
		return BasicType_Others, refType, refValue
	}
}

func ReflectTypeToBasicType(refType reflect.Type) BasicType {
	switch refType.Kind() {
	case reflect.Invalid:
		return BasicType_Invalid
	case reflect.Bool:
		return BasicType_Bool
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return BasicType_Int
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return BasicType_Uint
	case reflect.Float32, reflect.Float64:
		return BasicType_Float
	case reflect.String:
		return BasicType_String
	case reflect.Slice, reflect.Array:
		return BasicType_Slice
	case reflect.Map:
		return BasicType_Map
	case reflect.Struct:
		return BasicType_Struct
	default:
		return BasicType_Others
	}
}

// 获取字符串值的格式。
//   BasicType: 返回字符串可能的基础类型。可能是以下值：
//      String:  普通字符串
//      Float:   带有小数点的数字
//      Int:     负整数、或小于等于 MaxInt64 的非负整数
//      Uint:    非负且大于 MaxInt64 的整数
//      Bool:    "bool" 或者 "false"
//      Nil:     "null"
//      Slice:   JSON 数组
//      Map:     JSON 对象
func GetStrValueType(str string) (BasicType, reflect.Type, reflect.Value, interface{}) {
	if str = strings.TrimSpace(str); str == "" {
		return BasicType_String, reflect.TypeOf(str), reflect.ValueOf(str), str
	}

	// float|int|uint
	if pos := strings.Index(str, "."); pos != -1 {
		if regexpUtil.IsFloat(str) {
			if val, err := strconv.ParseFloat(str, 10); reflect2.IsNil(err) {
				return BasicType_Float, reflect.TypeOf(val), reflect.ValueOf(val), val
			}
		}
	} else {
		if regexpUtil.IsInt(str) {
			if val, err := strconv.ParseInt(str, 10, 64); reflect2.IsNil(err) {
				return BasicType_Int, reflect.TypeOf(val), reflect.ValueOf(val), val
			}
			if val, err := strconv.ParseUint(str, 10, 64); reflect2.IsNil(err) {
				return BasicType_Uint, reflect.TypeOf(val), reflect.ValueOf(val), val
			}
		}
	}

	// bool
	if val, err := strconv.ParseBool(strings.ToLower(strings.TrimSpace(str))); reflect2.IsNil(err) {
		return BasicType_Bool, reflect.TypeOf(val), reflect.ValueOf(val), val
	}

	// slice|map
	switch obj := jsonIter.Get([]byte(str)); obj.ValueType() {
	case jsonIter.StringValue:
		val := obj.ToString()
		return BasicType_String, reflect.TypeOf(val), reflect.ValueOf(val), val
	case jsonIter.NilValue:
		// jsonIter 是以`n`开头的表示为nil类型，存在以`n`开头的字符串，所以这里再对比一下是否为null
		if str == "null" {
			return BasicType_Nil, reflect.TypeOf(nil), reflect.ValueOf(nil), nil
		}
	case jsonIter.ArrayValue:
		// jsonIter 是以`[`开头的表示为map，存在以`[`开头的字符串，所以这里再反序列化一次，验证是否为对象
		if err := jsonApi.UnmarshalFromString(str, &[]interface{}{}); err == nil {
			val := obj.GetInterface()
			return BasicType_Slice, reflect.TypeOf(val), reflect.ValueOf(val), val
		}
	case jsonIter.ObjectValue:
		// jsonIter 是以`{`开头的表示为map，存在以`{`开头的字符串，所以这里再反序列化一次，验证是否为对象
		if err := jsonApi.UnmarshalFromString(str, &map[string]interface{}{}); err == nil {
			val := obj.GetInterface()
			return BasicType_Map, reflect.TypeOf(val), reflect.ValueOf(val), val
		}
	}

	return BasicType_String, reflect.TypeOf(str), reflect.ValueOf(str), str
}

// 判断变量是否为空（nil、0、""、长度为 0 的 Slice 或 Map）
func IsEmpty(a interface{}) bool {
	return IsEmpty2(a)
}

// 判断变量是否为空（nil、0、""、长度为 0 的 Slice 或 Map）
//   falseAsEmpty: 是否把 false 当作空值。默认为 true。
func IsEmpty2(a interface{}, falseAsEmpty ...bool) bool {
	if reflect2.IsNil(a) {
		return true
	}
	switch t, _, v := GetBasicType(a); t {
	case BasicType_Bool:
		return (len(falseAsEmpty) == 0 || falseAsEmpty[0]) && !v.Bool()
	case BasicType_Int:
		return v.Int() == 0
	case BasicType_Uint:
		return v.Uint() == 0
	case BasicType_Float:
		return v.Float() == 0
	case BasicType_String, BasicType_Slice, BasicType_Map:
		return v.Len() == 0
	}
	return false
}

func (this BasicType) GetValue(val reflect.Value) interface{} {
	switch this {
	case BasicType_Nil:
		return nil
	case BasicType_Bool:
		return val.Bool()
	case BasicType_Int:
		return val.Int()
	case BasicType_Uint:
		return val.Uint()
	case BasicType_Float:
		return val.Float()
	case BasicType_String:
		return val.String()
	default:
		return val.Interface()
	}
}

func (this BasicType) String() string {
	switch this {
	case BasicType_Invalid:
		return "invalid"
	case BasicType_Nil:
		return "nil"
	case BasicType_Bool:
		return "bool"
	case BasicType_Int:
		return "int"
	case BasicType_Uint:
		return "uint"
	case BasicType_Float:
		return "float"
	case BasicType_String:
		return "string"
	case BasicType_Slice:
		return "slice"
	case BasicType_Map:
		return "map"
	case BasicType_Struct:
		return "struct"
	case BasicType_Others:
		return "others"
	default:
		return fmt.Sprintf("unknown(%d)", int(this))
	}
}
