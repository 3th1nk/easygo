package arrUtil

import (
	"github.com/3th1nk/easygo/util/convertor"
	"strconv"
	"testing"
)

func TestRemoveAt(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}

	for i, n := 0, len(arr); i < n; i++ {
		val := RemoveIntAt(arr, i)
		t.Logf("val: %v, i=%v", val, i)

		if len(val) != len(arr)-1 {
			t.Errorf("assert faild")
			return
		}

		if i != 0 {
			if val[i-1] != arr[i-1] {
				t.Errorf("assert faild")
				return
			}
		}

		if i != n-1 {
			if val[i] != arr[i+1] {
				t.Errorf("assert faild")
				return
			}
		}
	}
}

func TestRemoveInt(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 1, 2, 3, 4, 5}

	for i := 0; i < 10; i++ {
		val := RemoveInt(arr, i)
		t.Logf("remove(%v):  %v", i, convertor.ToStringNoError(val))

		val = RemoveIntN(arr, 1, i)
		t.Logf("removeN(%v): %v", i, convertor.ToStringNoError(val))
	}
}

func TestRemoveString(t *testing.T) {
	arr := []string{"1", "2", "3", "4", "5", "1", "2", "3", "4", "5"}

	for i := 0; i < 10; i++ {
		val := RemoveString(arr, strconv.Itoa(i))
		t.Logf("remove(%v):  %v", i, convertor.ToStringNoError(val))

		val = RemoveStringN(arr, 1, strconv.Itoa(i))
		t.Logf("removeN(%v): %v", i, convertor.ToStringNoError(val))
	}
}

func TestRemoveDuplicate(t *testing.T) {
	type a struct {
		Name string
		Val  int
	}
	arr := []a{
		{Name: "1", Val: 1},
		{Name: "1", Val: 1},
		{Name: "2", Val: 2},
		{Name: "1", Val: 1},
		{Name: "3", Val: 3},
		{Name: "2", Val: 2},
	}

	arr2 := RemoveDuplicate(arr, func(val interface{}) string {
		return val.(a).Name
	})

	t.Log(convertor.ToStringNoError(arr2))
}
