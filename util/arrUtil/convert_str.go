package arrUtil

import "github.com/3th1nk/easygo/util/convertor"

func StrToObj(a []string) []interface{} {
	arr := make([]interface{}, len(a))
	for i, v := range a {
		arr[i] = v
	}
	return arr
}

func StrToInt(a []string) ([]int, error) {
	arr, err := make([]int, len(a)), error(nil)
	for i, s := range a {
		if arr[i], err = convertor.StrToInt(s); err != nil {
			return nil, err
		}
	}
	return arr, nil
}

func StrToIntNoErr(a []string) []int {
	arr, _ := StrToInt(a)
	return arr
}

func StrToInt64(a []string) ([]int64, error) {
	arr, err := make([]int64, len(a)), error(nil)
	for i, s := range a {
		if arr[i], err = convertor.StrToInt64(s); err != nil {
			return nil, err
		}
	}
	return arr, nil
}

func StrToInt64NoErr(a []string) []int64 {
	arr, _ := StrToInt64(a)
	return arr
}

func StrToFloat(a []string, ignoreErr ...bool) ([]float64, error) {
	arr := make([]float64, 0, len(a))
	for _, s := range a {
		f, err := convertor.StrToFloat(s)
		if err == nil {
			arr = append(arr, f)
		} else if len(ignoreErr) == 0 || !ignoreErr[0] {
			return nil, err
		}
	}
	return arr, nil
}

func StrToFloatNoErr(a []string, ignoreErr ...bool) []float64 {
	arr, _ := StrToFloat(a, ignoreErr...)
	return arr
}

func StrToMap(arr []string) map[string]bool {
	m := make(map[string]bool, len(arr))
	for _, val := range arr {
		m[val] = true
	}
	return m
}
