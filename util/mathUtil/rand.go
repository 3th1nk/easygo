package mathUtil

import (
	"math"
	"math/rand"
)

func Float(x float64, decimals ...int) float64 {
	return Float2(0, x, decimals...)
}

func Float2(x, y float64, decimals ...int) float64 {
	var theDecimal int
	if len(decimals) != 0 && decimals[0] > 0 {
		theDecimal = decimals[0]
	} else {
		theDecimal = 3
	}
	n := math.Pow10(theDecimal)
	xn, yn := x*n, y*n
	return xn + float64(rand.Int63n(int64(math.Floor(yn-xn))))/n
}
