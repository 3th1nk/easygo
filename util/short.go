package util

import "fmt"

var (
	DefaultShortStrMax = 1023
)

// ShortStr 长字符串截断、只保留前面的一部分，其余部分用 “...(%d more)” 代替。
func ShortStr(str string, max ...int) string {
	var theMax int
	if len(max) != 0 {
		theMax = max[0]
		if n := theMax % 3; n != 0 {
			theMax += 3 - n
		}
	} else {
		theMax = DefaultShortStrMax
	}
	n := len(str)
	if theMax <= 0 {
		return fmt.Sprintf("str(%d)", n)
	} else if n > theMax {
		return fmt.Sprintf("%s...(%d chars)", str[:theMax], n)
	}
	return str
}
