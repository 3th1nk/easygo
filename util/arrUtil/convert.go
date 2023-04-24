package arrUtil

import (
	"github.com/3th1nk/easygo/util/convertor"
	"github.com/modern-go/reflect2"
	"reflect"
)

// ------------------------------------------------------------------------------ ToInterface

// 将数组或切片转换为指定的切片类型。
//
// 参数：
//   a: 要转换的数组。
//   dstElemType: 目标元素的类型。
//   f: 转换时的自定义回调函数。
func ToType(slice interface{}, dstElemType reflect.Type, f ...func(a interface{}) interface{}) interface{} {
	if reflect2.IsNil(slice) {
		return nil
	}

	inVal := reflect.ValueOf(slice)
	if kind := inVal.Kind(); kind != reflect.Array && kind != reflect.Slice {
		panic("参数 slice 必须是切片类型")
	}

	var theF func(a interface{}) interface{}
	if len(f) != 0 {
		theF = f[0]
	}

	n := inVal.Len()
	outVal := reflect.MakeSlice(reflect.SliceOf(dstElemType), n, n)
	for i := 0; i < n; i++ {
		v := inVal.Index(i).Interface()
		if theF != nil {
			v = theF(v)
		}
		outVal.Index(i).Set(reflect.ValueOf(v))
	}
	return outVal.Interface()
}

// 将数组转换为 interface 数组
// 参数:
//   slice: 要转换的源，必须是切片类型
func ToInterface(slice interface{}, f ...func(a interface{}) interface{}) []interface{} {
	return ToSlice(slice, f...)
}

// 将数组转换为 interface 数组
// 参数:
//   slice: 要转换的源，必须是切片类型
func ToSlice(slice interface{}, f ...func(a interface{}) interface{}) []interface{} {
	if reflect2.IsNil(slice) {
		return nil
	}

	reflectVal := reflect.ValueOf(slice)
	if kind := reflectVal.Kind(); kind != reflect.Array && kind != reflect.Slice {
		panic("参数 slice 必须是切片类型")
	}

	var theF func(a interface{}) interface{}
	if len(f) != 0 {
		theF = f[0]
	}

	n := reflectVal.Len()
	arr := make([]interface{}, n)
	for i := 0; i < n; i++ {
		v := reflectVal.Index(i).Interface()
		if theF != nil {
			v = theF(v)
		}
		arr[i] = v
	}
	return arr
}

// ------------------------------------------------------------------------------ ToStr

// 将数组转换为 string 数组
//   slice: 要转换的源，必须是切片类型
func ToStr(slice interface{}, f ...func(i int) (val string, err error)) ([]string, error) {
	if len(f) == 0 || f[0] == nil {
		return ToStrIf(slice, nil)
	} else {
		return ToStrIf(slice, func(i int) (val string, ok bool, err error) {
			ok = true
			val, err = f[0](i)
			return
		})
	}
}

// 将数组转换为 string 数组
//   slice: 要转换的源，必须是切片类型
//   f: 过滤函数，返回值中的 ok 表示是否需要加入到结果数组中。
func ToStrIf(slice interface{}, f func(i int) (val string, ok bool, err error)) ([]string, error) {
	if reflect2.IsNil(slice) {
		return nil, nil
	}

	reflectVal := reflect.ValueOf(slice)
	if kind := reflectVal.Kind(); kind != reflect.Array && kind != reflect.Slice {
		panic("参数 slice 必须是切片类型")
	}

	if f == nil {
		f = func(i int) (val string, ok bool, err error) {
			ok = true
			val, err = convertor.ToString(reflectVal.Index(i).Interface())
			return
		}
	}

	if val, err := doSliceTo(slice, reflect.TypeOf([]string{}), func(i int) (val interface{}, ok bool, err error) {
		return f(i)
	}); !reflect2.IsNil(err) {
		return nil, err
	} else {
		return val.([]string), nil
	}
}

func ToStrNoError(slice interface{}, f ...func(i int) string) []string {
	if len(f) == 0 || f[0] == nil {
		val, _ := ToStrIf(slice, nil)
		return val
	} else {
		val, _ := ToStrIf(slice, func(i int) (val string, ok bool, err error) {
			return f[0](i), true, nil
		})
		return val
	}
}

func ToStrIfNoError(slice interface{}, f func(i int) (val string, ok bool)) []string {
	if f == nil {
		val, _ := ToStrIf(slice, nil)
		return val
	} else {
		val, _ := ToStrIf(slice, func(i int) (val string, ok bool, err error) {
			err = nil
			val, ok = f(i)
			return
		})
		return val
	}
}

// ------------------------------------------------------------------------------ ToInt

// 将数组转换为 int 数组
//   slice: 要转换的源，必须是切片类型
func ToInt(slice interface{}, f ...func(i int) (val int, err error)) ([]int, error) {
	if len(f) == 0 || f[0] == nil {
		return ToIntIf(slice, nil)
	} else {
		return ToIntIf(slice, func(i int) (val int, ok bool, err error) {
			ok = true
			val, err = f[0](i)
			return
		})
	}
}

// 将数组转换为 int 数组
//   slice: 要转换的源，必须是切片类型
//   f: 过滤函数，返回值中的 ok 表示是否需要加入到结果数组中。
func ToIntIf(slice interface{}, f func(i int) (val int, ok bool, err error)) ([]int, error) {
	if reflect2.IsNil(slice) {
		return nil, nil
	}

	reflectVal := reflect.ValueOf(slice)
	if kind := reflectVal.Kind(); kind != reflect.Array && kind != reflect.Slice {
		panic("参数 slice 必须是切片类型")
	}

	if f == nil {
		f = func(i int) (int, bool, error) {
			if val, err := convertor.ToInt(reflectVal.Index(i).Interface()); !reflect2.IsNil(err) {
				return 0, false, err
			} else {
				return val, true, nil
			}
		}
	}

	if val, err := doSliceTo(slice, reflect.TypeOf([]int{}), func(i int) (val interface{}, ok bool, err error) {
		return f(i)
	}); !reflect2.IsNil(err) {
		return nil, err
	} else {
		return val.([]int), nil
	}
}

func ToIntNoError(slice interface{}, f ...func(i int) int) []int {
	if len(f) == 0 || f[0] == nil {
		val, _ := ToIntIf(slice, nil)
		return val
	} else {
		val, _ := ToIntIf(slice, func(i int) (val int, ok bool, err error) {
			return f[0](i), true, nil
		})
		return val
	}
}

func ToIntIfNoError(slice interface{}, f func(i int) (val int, ok bool)) []int {
	if f == nil {
		val, _ := ToIntIf(slice, nil)
		return val
	} else {
		val, _ := ToIntIf(slice, func(i int) (val int, ok bool, err error) {
			err = nil
			val, ok = f(i)
			return
		})
		return val
	}
}

// ------------------------------------------------------------------------------ ToInt64

// 将数组转换为 int64 数组
//   slice: 要转换的源，必须是切片类型
func ToInt64(slice interface{}, f ...func(i int) (val int64, err error)) ([]int64, error) {
	if len(f) == 0 || f[0] == nil {
		return ToInt64If(slice, nil)
	} else {
		return ToInt64If(slice, func(i int) (val int64, ok bool, err error) {
			ok = true
			val, err = f[0](i)
			return
		})
	}
}

// 将数组转换为 int64 数组
//   slice: 要转换的源，必须是切片类型
//   f: 过滤函数，返回值中的 ok 表示是否需要加入到结果数组中。
func ToInt64If(slice interface{}, f func(i int) (val int64, ok bool, err error)) ([]int64, error) {
	if reflect2.IsNil(slice) {
		return nil, nil
	}

	reflectVal := reflect.ValueOf(slice)
	if kind := reflectVal.Kind(); kind != reflect.Array && kind != reflect.Slice {
		panic("参数 slice 必须是切片类型")
	}

	if f == nil {
		f = func(i int) (int64, bool, error) {
			if val, err := convertor.ToInt64(reflectVal.Index(i).Interface()); !reflect2.IsNil(err) {
				return 0, false, err
			} else {
				return val, true, nil
			}
		}
	}

	if val, err := doSliceTo(slice, reflect.TypeOf([]int64{}), func(i int) (val interface{}, ok bool, err error) {
		return f(i)
	}); !reflect2.IsNil(err) {
		return nil, err
	} else {
		return val.([]int64), nil
	}
}

func ToInt64NoError(slice interface{}, f ...func(i int) int64) []int64 {
	if len(f) == 0 || f[0] == nil {
		val, _ := ToInt64If(slice, nil)
		return val
	} else {
		val, _ := ToInt64If(slice, func(i int) (val int64, ok bool, err error) {
			return f[0](i), true, nil
		})
		return val
	}
}

func ToInt64IfNoError(slice interface{}, f func(i int) (val int64, ok bool)) []int64 {
	if f == nil {
		val, _ := ToInt64If(slice, nil)
		return val
	} else {
		val, _ := ToInt64If(slice, func(i int) (val int64, ok bool, err error) {
			err = nil
			val, ok = f(i)
			return
		})
		return val
	}
}

// ------------------------------------------------------------------------------ ToInt32

// 将数组转换为 int32 数组
//   slice: 要转换的源，必须是切片类型
func ToInt32(slice interface{}, f ...func(i int) (val int32, err error)) ([]int32, error) {
	if len(f) == 0 || f[0] == nil {
		return ToInt32If(slice, nil)
	} else {
		return ToInt32If(slice, func(i int) (val int32, ok bool, err error) {
			ok = true
			val, err = f[0](i)
			return
		})
	}
}

// 将数组转换为 int32 数组
//   slice: 要转换的源，必须是切片类型
//   f: 过滤函数，返回值中的 ok 表示是否需要加入到结果数组中。
func ToInt32If(slice interface{}, f func(i int) (val int32, ok bool, err error)) ([]int32, error) {
	if reflect2.IsNil(slice) {
		return nil, nil
	}

	reflectVal := reflect.ValueOf(slice)
	if kind := reflectVal.Kind(); kind != reflect.Array && kind != reflect.Slice {
		panic("参数 slice 必须是切片类型")
	}

	if f == nil {
		f = func(i int) (int32, bool, error) {
			if val, err := convertor.ToInt32(reflectVal.Index(i).Interface()); !reflect2.IsNil(err) {
				return 0, false, err
			} else {
				return val, true, nil
			}
		}
	}

	if val, err := doSliceTo(slice, reflect.TypeOf([]int32{}), func(i int) (val interface{}, ok bool, err error) {
		return f(i)
	}); !reflect2.IsNil(err) {
		return nil, err
	} else {
		return val.([]int32), nil
	}
}

func ToInt32NoError(slice interface{}, f ...func(i int) int32) []int32 {
	if len(f) == 0 || f[0] == nil {
		val, _ := ToInt32If(slice, nil)
		return val
	} else {
		val, _ := ToInt32If(slice, func(i int) (val int32, ok bool, err error) {
			return f[0](i), true, nil
		})
		return val
	}
}

func ToInt32IfNoError(slice interface{}, f func(i int) (val int32, ok bool)) []int32 {
	if f == nil {
		val, _ := ToInt32If(slice, nil)
		return val
	} else {
		val, _ := ToInt32If(slice, func(i int) (val int32, ok bool, err error) {
			err = nil
			val, ok = f(i)
			return
		})
		return val
	}
}

// ------------------------------------------------------------------------------ ToFloat

// 将数组转换为 float64 数组
//   slice: 要转换的源，必须是切片类型
func ToFloat(slice interface{}, f ...func(i int) (val float64, err error)) ([]float64, error) {
	if len(f) == 0 || f[0] == nil {
		return ToFloatIf(slice, nil)
	} else {
		return ToFloatIf(slice, func(i int) (val float64, ok bool, err error) {
			ok = true
			val, err = f[0](i)
			return
		})
	}
}

// 将数组转换为 float64 数组
//   slice: 要转换的源，必须是切片类型
//   f: 过滤函数，返回值中的 ok 表示是否需要加入到结果数组中。
func ToFloatIf(slice interface{}, f func(i int) (val float64, ok bool, err error)) ([]float64, error) {
	if reflect2.IsNil(slice) {
		return nil, nil
	}

	reflectVal := reflect.ValueOf(slice)
	if kind := reflectVal.Kind(); kind != reflect.Array && kind != reflect.Slice {
		panic("参数 slice 必须是切片类型")
	}

	if f == nil {
		f = func(i int) (float64, bool, error) {
			if val, err := convertor.ToFloat(reflectVal.Index(i).Interface()); !reflect2.IsNil(err) {
				return 0, false, err
			} else {
				return val, true, nil
			}
		}
	}

	if val, err := doSliceTo(slice, reflect.TypeOf([]float64{}), func(i int) (val interface{}, ok bool, err error) {
		return f(i)
	}); !reflect2.IsNil(err) {
		return nil, err
	} else {
		return val.([]float64), nil
	}
}

func ToFloatNoError(slice interface{}, f ...func(i int) float64) []float64 {
	if len(f) == 0 || f[0] == nil {
		val, _ := ToFloatIf(slice, nil)
		return val
	} else {
		val, _ := ToFloatIf(slice, func(i int) (val float64, ok bool, err error) {
			return f[0](i), true, nil
		})
		return val
	}
}

func ToFloatIfNoError(slice interface{}, f func(i int) (val float64, ok bool)) []float64 {
	if f == nil {
		val, _ := ToFloatIf(slice, nil)
		return val
	} else {
		val, _ := ToFloatIf(slice, func(i int) (val float64, ok bool, err error) {
			err = nil
			val, ok = f(i)
			return
		})
		return val
	}
}

// ------------------------------------------------------------------------------ ToBool

// 将数组转换为 int 数组
//   slice: 要转换的源，必须是切片类型
func ToBool(slice interface{}, f ...func(i int) (val bool, err error)) ([]bool, error) {
	if len(f) == 0 || f[0] == nil {
		return ToBoolIf(slice, nil)
	} else {
		return ToBoolIf(slice, func(i int) (val bool, ok bool, err error) {
			ok = true
			val, err = f[0](i)
			return
		})
	}
}

// 将数组转换为 int 数组
//   slice: 要转换的源，必须是切片类型
//   f: 过滤函数，返回值中的 ok 表示是否需要加入到结果数组中。
func ToBoolIf(slice interface{}, f func(i int) (val bool, ok bool, err error)) ([]bool, error) {
	if reflect2.IsNil(slice) {
		return nil, nil
	}

	reflectVal := reflect.ValueOf(slice)
	if kind := reflectVal.Kind(); kind != reflect.Array && kind != reflect.Slice {
		panic("参数 slice 必须是切片类型")
	}

	if f == nil {
		f = func(i int) (bool, bool, error) {
			if val, err := convertor.ToBool(reflectVal.Index(i).Interface()); !reflect2.IsNil(err) {
				return false, false, err
			} else {
				return val, true, nil
			}
		}
	}

	if val, err := doSliceTo(slice, reflect.TypeOf([]bool{}), func(i int) (val interface{}, ok bool, err error) {
		return f(i)
	}); !reflect2.IsNil(err) {
		return nil, err
	} else {
		return val.([]bool), nil
	}
}

func ToBoolNoError(slice interface{}, f ...func(i int) bool) []bool {
	if len(f) == 0 || f[0] == nil {
		val, _ := ToBoolIf(slice, nil)
		return val
	} else {
		val, _ := ToBoolIf(slice, func(i int) (val bool, ok bool, err error) {
			return f[0](i), true, nil
		})
		return val
	}
}

func ToBoolIfNoError(slice interface{}, f func(i int) (val bool, ok bool)) []bool {
	if f == nil {
		val, _ := ToBoolIf(slice, nil)
		return val
	} else {
		val, _ := ToBoolIf(slice, func(i int) (val bool, ok bool, err error) {
			err = nil
			val, ok = f(i)
			return
		})
		return val
	}
}

// ------------------------------------------------------------------------------ slice
func doSliceTo(slice interface{}, elemType reflect.Type, f func(i int) (val interface{}, ok bool, err error)) (interface{}, error) {
	srcVal := reflect.ValueOf(slice)
	n := srcVal.Len()
	arr := reflect.MakeSlice(elemType, 0, n)
	for i := 0; i < n; i++ {
		if val, ok, err := f(i); !reflect2.IsNil(err) {
			return nil, err
		} else if ok {
			arr = reflect.Append(arr, reflect.ValueOf(val))
		}
	}
	return arr.Interface(), nil
}
