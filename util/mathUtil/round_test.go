package mathUtil

import (
	"math"
	"strconv"
	"testing"
)

func TestRound(t *testing.T) {
	f := 1234567.7654321
	s := strconv.FormatFloat(f, 'f', -1, 64)
	t.Log(s)
	for n := -7; n <= 7; n++ {
		t.Logf("round(%v): %v", n, strconv.FormatFloat(Round(f, n), 'f', -1, 64))
	}
}

func TestPow(t *testing.T) {
	for _, n := range []int{-2, -1, 0, 1, 2} {
		t.Logf("Pow(10, %v): %v", n, math.Pow10(n))
	}
}

func TestMathRound(t *testing.T) {
	for _, v := range []float64{
		-0.7,
		-0.5,
		-0.2,
		0,
		0.2,
		0.5,
		0.7,
	} {
		n := math.Round(v)
		t.Logf("%v: %v", v, n)
	}
}
