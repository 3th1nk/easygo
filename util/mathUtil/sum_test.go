package mathUtil

import "testing"

func TestSumInt64(t *testing.T) {
	if n := SumInt64([]int64{1, 2, 3, 4, 5}); n != 15 {
		t.Errorf("assert faild: %v", n)
	}
	if n := SumInt64([]int64{1, 2}); n != 3 {
		t.Errorf("assert faild: %v", n)
	}
	if n := SumInt64([]int64{1}); n != 1 {
		t.Errorf("assert faild: %v", n)
	}
}
