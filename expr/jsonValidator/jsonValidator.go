package jsonValidator

import (
	"fmt"
	"github.com/3th1nk/easygo/util/arrUtil"
	"github.com/3th1nk/easygo/util/comparer"
	"github.com/3th1nk/easygo/util/convertor"
	"github.com/3th1nk/easygo/util/jsonUtil"
	"github.com/3th1nk/easygo/util/mapUtil"
	jsonIter "github.com/json-iterator/go"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// 创建一个表达式验证器。
func New(data interface{}) *Validator {
	obj := &Validator{}
	if data == nil {
		obj.obj = jsonUtil.Wrap(nil)
	} else {
		switch t := data.(type) {
		case []byte:
			obj.obj = jsonUtil.Get(t)
		case string:
			obj.obj = jsonUtil.Get([]byte(t))
		default:
			bytes, _ := jsonUtil.NoOmitemptyApi().Marshal(data)
			obj.obj = jsonUtil.Get(bytes)
		}
	}
	return obj
}

type Validator struct {
	obj jsonIter.Any
}

// 获取要验证的 Json 对象
func (this *Validator) JsonObj() jsonIter.Any {
	return this.obj
}

// 验证指定路径对应的值是否符合表达式。
//   若指定路径存在、且能够与 val 进行对应的比较操作，则返回比较结果；
//   若路径不存在、路径对应的值与 val 之间无法进行对应的比较（比如 int 之间无法进行 contains 比较等），则返回错误。
func (this *Validator) Validate(path string, operator comparer.Operator, val interface{}, option int) (bool, error) {
	if pathVal, err := this.Value(path); err != nil {
		return false, err
	} else {
		return comparer.Compare(pathVal, val, operator, option)
	}
}

// 获取指定路径对应的值。
//   当路径不存在时返回错误。
func (this *Validator) Value(path string) (val interface{}, err error) {
	arr, isSlice, err := this.getValue(path)
	if isSlice {
		return arr, err
	} else if len(arr) != 0 {
		return arr[0], err
	} else {
		return nil, err
	}
}

func (this *Validator) getValue(path string) (val []interface{}, isSlice bool, err error) {
	isSlice, err = this.getValueToSlice(this.obj, path, &val, "")
	return
}

func (this *Validator) getValueToSlice(obj jsonIter.Any, path string, val *[]interface{}, walkedPath string) (isSlice bool, err error) {
	pos := strings.Index(path, "*")
	if pos == -1 {
		v, vv, err := this.getOneValue(obj, path)
		if err != nil {
			return false, err
		} else if v.ValueType() == jsonIter.ArrayValue {
			for i, n := 0, v.Size(); i < n; i++ {
				*val = append(*val, v.Get(i).GetInterface())
			}
			return true, nil
		} else {
			*val = append(*val, vv)
			return false, nil
		}
	}

	prefix, suffix, arrPath := strings.TrimRight(path[:pos], "."), strings.TrimLeft(path[pos+1:], "."), ""
	var arrObj jsonIter.Any
	if prefix == "" {
		arrObj = obj
		arrPath = walkedPath
	} else {
		arrObj = obj.Get(this.splitPath(prefix)...)
		arrPath = strings.TrimLeft(walkedPath+"."+prefix, ".")
	}
	if valType := arrObj.ValueType(); valType == jsonIter.InvalidValue {
		return false, fmt.Errorf("路径 '%v' 无效", arrPath)
	} else if valType != jsonIter.ArrayValue {
		return false, fmt.Errorf("路径 '%v' 不是数组格式", arrPath)
	}

	for i, n := 0, arrObj.Size(); i < n; i++ {
		arrItemObj := arrObj.Get(i)
		if suffix == "" {
			*val = append(*val, arrItemObj)
		} else {
			if _, err := this.getValueToSlice(arrItemObj, suffix, val, arrPath+strconv.Itoa(i)); err != nil {
				return true, err
			}
		}
	}

	return true, nil
}

func (this *Validator) getOneValue(jsonObj jsonIter.Any, path string) (jsonIter.Any, interface{}, error) {
	pathSegs := make([]interface{}, 0, 16)
	for _, str := range strings.Split(path, ".") {
		if n, e := strconv.ParseInt(str, 10, 64); e == nil {
			pathSegs = append(pathSegs, int(n))
		} else {
			pathSegs = append(pathSegs, str)
		}
	}

	obj := jsonObj.Get(pathSegs...)
	switch obj.ValueType() {
	case jsonIter.InvalidValue:
		return obj, nil, fmt.Errorf("路径 %s 无效", path)
	case jsonIter.StringValue:
		return obj, obj.ToString(), nil
	case jsonIter.NumberValue:
		f := obj.ToFloat64()
		n, frac := math.Modf(f)
		if math.Abs(frac) < 0.000001 {
			return obj, int64(n), nil
		} else {
			return obj, f, nil
		}
	case jsonIter.NilValue:
		return obj, nil, nil
	case jsonIter.BoolValue:
		return obj, obj.ToBool(), nil
	case jsonIter.ArrayValue:
		arr := make([]interface{}, 0)
		obj.ToVal(&arr)
		return obj, arr, nil
	case jsonIter.ObjectValue:
		dict := make(map[string]interface{})
		obj.ToVal(&dict)
		return obj, dict, nil
	default:
		return obj, nil, fmt.Errorf("路径 %s 无效", path)
	}
}

func (this *Validator) splitPath(path string) []interface{} {
	pathSegs := make([]interface{}, 0, 16)
	for _, str := range strings.Split(path, ".") {
		if n, e := strconv.ParseInt(str, 10, 64); e == nil {
			pathSegs = append(pathSegs, int(n))
		} else {
			pathSegs = append(pathSegs, str)
		}
	}
	return pathSegs
}

func (this *Validator) String(path string) (val string, err error) {
	arr, _, err := this.getValue(path)
	if err == nil && len(arr) != 0 {
		val, err = convertor.ToString(arr[0])
	}
	return
}

func (this *Validator) StringSlice(path string) (val []string, err error) {
	arr, _, err := this.getValue(path)
	if err == nil {
		val = arrUtil.ToStrNoError(arr)
	}
	return
}

func (this *Validator) Int(path string) (val int, err error) {
	arr, _, err := this.getValue(path)
	if err == nil && len(arr) != 0 {
		val, err = convertor.ToInt(arr[0])
	}
	return
}

func (this *Validator) IntSlice(path string) (val []int, err error) {
	arr, _, err := this.getValue(path)
	if err == nil {
		val = arrUtil.ToIntNoError(arr)
	}
	return
}

func (this *Validator) Int64(path string) (val int64, err error) {
	arr, _, err := this.getValue(path)
	if err == nil && len(arr) != 0 {
		val, err = convertor.ToInt64(arr[0])
	}
	return
}

func (this *Validator) Int64Slice(path string) (val []int64, err error) {
	arr, _, err := this.getValue(path)
	if err == nil {
		val = arrUtil.ToInt64NoError(arr)
	}
	return
}

func (this *Validator) Float(path string) (val float64, err error) {
	arr, _, err := this.getValue(path)
	if err == nil && len(arr) != 0 {
		val, err = convertor.ToFloat(arr[0])
	}
	return
}

func (this *Validator) FloatSlice(path string) (val []float64, err error) {
	arr, _, err := this.getValue(path)
	if err == nil {
		val = arrUtil.ToFloatNoError(arr)
	}
	return
}

func (this *Validator) Bool(path string) (val bool, err error) {
	arr, _, err := this.getValue(path)
	if err == nil && len(arr) != 0 {
		val, err = convertor.ToBool(arr[0])
	}
	return
}

func (this *Validator) BoolSlice(path string) (val []bool, err error) {
	arr, _, err := this.getValue(path)
	if err == nil {
		val = arrUtil.ToBoolNoError(arr)
	}
	return
}

func (this *Validator) StringObjectMap(path string) (val mapUtil.StringObjectMap, err error) {
	arr, _, err := this.getValue(path)
	if err == nil && len(arr) != 0 {
		if obj, ok := arr[0].(map[string]interface{}); !ok {
			return nil, fmt.Errorf("cannot convert %v('%v') to map", reflect.TypeOf(arr[0]), path)
		} else {
			return obj, nil
		}
	}
	return
}
