package mathUtil

import (
	"math"
	"strconv"
)

// 保留浮点数小数点后 N 位。
//   备注: 当 decimals 为负数时忽略，直接 返回原始值；
func Round(f float64, decimals int) float64 {
	if decimals < 0 {
		return f
	} else if decimals == 0 {
		return math.Round(f)
	} else {
		str := strconv.FormatFloat(f, 'f', decimals, 64)
		v, _ := strconv.ParseFloat(str, 64)
		return v
	}
}

// 保留浮点数小数点后 6 位。
func Round6(f float64) float64 {
	return Round(f, 6)
}
