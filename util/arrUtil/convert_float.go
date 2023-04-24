package arrUtil

import "strconv"

func FloatToObj(a []float64) []interface{} {
	arr := make([]interface{}, len(a))
	for i, v := range a {
		arr[i] = v
	}
	return arr
}

func FloatToStr(a []float64) []string {
	arr := make([]string, len(a))
	for i, v := range a {
		arr[i] = strconv.FormatFloat(v, 'f', -1, 64)
	}
	return arr
}

func FloatToInt(a []float64) []int {
	arr := make([]int, len(a))
	for i, v := range a {
		arr[i] = int(v)
	}
	return arr
}

func FloatToInt64(a []float64) []int64 {
	arr := make([]int64, len(a))
	for i, v := range a {
		arr[i] = int64(v)
	}
	return arr
}

func IntToFloat(a []int) []float64 {
	arr := make([]float64, len(a))
	for i, v := range a {
		arr[i] = float64(v)
	}
	return arr
}

func Int64ToFloat(a []int64) []float64 {
	arr := make([]float64, len(a))
	for i, v := range a {
		arr[i] = float64(v)
	}
	return arr
}
