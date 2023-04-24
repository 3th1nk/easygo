package arrUtil

import (
	"fmt"
	"strings"
	"testing"
)

func TestIndexOfSortedInt(t *testing.T) {
	arr := []int{1, 3, 5, 7, 9, 11}
	for _, innerArr := range [][]int{
		{0, -1},
		{1, 0},
		{2, -1},
		{3, 1},
		{4, -1},
		{5, 2},
		{6, -1},
		{7, 3},
		{8, -1},
		{9, 4},
		{10, -1},
		{11, 5},
		{12, -1},
	} {
		val, expect := IndexOfSortedInt(arr, innerArr[0]), innerArr[1]
		if val != expect {
			t.Error(fmt.Sprintf("assert faild: search(%v) assert %v, but %v", innerArr[0], expect, val))
		} else {
			t.Logf("IndexOfSortedInt(%v): %v", innerArr[0], val)
		}

		val, expect = IndexOfInt(arr, innerArr[0]), innerArr[1]
		if val != expect {
			t.Error(fmt.Sprintf("assert faild: search(%v) assert %v, but %v", innerArr[0], expect, val))
		} else {
			t.Logf("IndexOfInt(%v): %v", innerArr[0], val)
		}

		val, expect = IndexOf(arr, func(i int) bool {
			return arr[i] == innerArr[0]
		}), innerArr[1]
		if val != expect {
			t.Error(fmt.Sprintf("assert faild: search(%v) assert %v, but %v", innerArr[0], expect, val))
		} else {
			t.Logf("IndexOf(%v): %v", innerArr[0], val)
		}
	}
}

func TestIndexOfAnyString(t *testing.T) {
	a := strings.Split("a,b,c,d,e", ",")

	for s, n := range map[string]int{
		"a,b,c":    0,
		"b,c,a":    0,
		"b,c,e":    1,
		"e,b,c":    1,
		"aa,bb,cc": -1,
	} {
		find := strings.Split(s, ",")
		idx := IndexOfAnyString(a, find)
		if idx != n {
			t.Errorf("assert faild: expect %v, but %v", n, idx)
		}
	}
}
