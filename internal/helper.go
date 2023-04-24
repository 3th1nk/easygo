package internal

import "fmt"

func ShortStr(s string, max ...int) string {
	var theMax int
	if len(max) != 0 {
		theMax = max[0]
	} else {
		theMax = 256
	}
	n := len(s)
	if n > theMax {
		return fmt.Sprintf("%s...(%d more)", s[:theMax], n-theMax)
	}
	return s
}
