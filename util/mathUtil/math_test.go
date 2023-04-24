package mathUtil

import (
	"fmt"
	"github.com/3th1nk/easygo/util"
	"math"
	"testing"
)

func TestModf(t *testing.T) {
	for _, v := range []float64{1, 1.0001, 1.000000000001, -1, -1.0001, -1.000000000001} {
		i, f := math.Modf(v)
		util.Println("%16s: i=%v, f=%v", fmt.Sprintf("%v", v), i, fmt.Sprintf("%.20f", f))
	}
}
