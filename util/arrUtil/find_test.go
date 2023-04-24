package arrUtil

import (
	"fmt"
	"github.com/modern-go/reflect2"
	"testing"
)

func TestFindInterface(t *testing.T) {
	type abc struct{ Id int }
	arr := []*abc{
		{Id: 1},
		{Id: 3},
		{Id: 5},
		{Id: 7},
	}

	resultArr, ok := Find(arr, func(i int) bool {
		return arr[i].Id == 3
	}, -1).([]*abc)
	if !ok || len(resultArr) != 1 || resultArr[0].Id != 3 {
		t.Error(fmt.Sprintf("assert faild"))
	}

	resultArr, ok = Find(arr, func(i int) bool {
		return arr[i].Id > 1
	}).([]*abc)
	if !ok || len(resultArr) != 3 || resultArr[0].Id != 3 || resultArr[1].Id != 5 || resultArr[2].Id != 7 {
		t.Error(fmt.Sprintf("assert faild"))
	}

	resultArr, ok = Find(arr, func(i int) bool {
		return arr[i].Id == 123456576
	}, -1).([]*abc)
	if !ok || len(resultArr) != 0 {
		t.Error(fmt.Sprintf("assert faild"))
	}

	resultVal, _ := First(arr, func(i int) bool {
		return arr[i].Id > 1
	})
	if reflect2.IsNil(resultVal) || resultVal.(*abc).Id != 3 {
		t.Error(fmt.Sprintf("assert faild"))
	}

	resultVal, _ = First(arr, func(i int) bool {
		return arr[i].Id < 1
	})
	if !reflect2.IsNil(resultVal) {
		t.Error(fmt.Sprintf("assert faild"))
	}
}

func TestFindInt(t *testing.T) {
	arr := []int{1, 3, 5, 7}

	resultArr, ok := Find(arr, func(i int) bool {
		return arr[i] == 3
	}).([]int)
	if !ok || len(resultArr) != 1 || resultArr[0] != 3 {
		t.Error(fmt.Sprintf("assert faild"))
	}

	resultArr, ok = Find(arr, func(i int) bool {
		return arr[i] > 1
	}).([]int)
	if !ok || len(resultArr) != 3 || resultArr[0] != 3 || resultArr[1] != 5 || resultArr[2] != 7 {
		t.Error(fmt.Sprintf("assert faild"))
	}

	resultArr, ok = FindN(arr, func(i int) bool {
		return arr[i] > 1
	}, 5, 1).([]int)
	if !ok || len(resultArr) != 2 || resultArr[0] != 5 || resultArr[1] != 7 {
		t.Error(fmt.Sprintf("assert faild"))
	}

	resultArr, ok = Find(arr, func(i int) bool {
		return arr[i] == 123456576
	}).([]int)
	if !ok || len(resultArr) != 0 {
		t.Error(fmt.Sprintf("assert faild"))
	}

	resultVal, _ := First(arr, func(i int) bool {
		return arr[i] > 1
	})
	if resultVal != 3 {
		t.Error(fmt.Sprintf("assert faild"))
	}

	resultVal, _ = First(arr, func(i int) bool {
		return arr[i] < 1
	})
	if !reflect2.IsNil(resultVal) {
		t.Error(fmt.Sprintf("assert faild"))
	}
}

func TestFindString(t *testing.T) {
	arr := []string{"1", "3", "5", "7"}

	resultArr, ok := Find(arr, func(i int) bool {
		return arr[i] == "3"
	}).([]string)
	if !ok || len(resultArr) != 1 || resultArr[0] != "3" {
		t.Error(fmt.Sprintf("assert faild"))
	}

	resultArr, ok = Find(arr, func(i int) bool {
		return arr[i] > "1"
	}).([]string)
	if !ok || len(resultArr) != 3 || resultArr[0] != "3" || resultArr[1] != "5" || resultArr[2] != "7" {
		t.Error(fmt.Sprintf("assert faild"))
	}

	resultArr, ok = Find(arr, func(i int) bool {
		return arr[i] == "123456576"
	}).([]string)
	if !ok || len(resultArr) != 0 {
		t.Error(fmt.Sprintf("assert faild"))
	}

	resultVal, _ := First(arr, func(i int) bool {
		return arr[i] > "1"
	})
	if resultVal != "3" {
		t.Error(fmt.Sprintf("assert faild"))
	}

	resultVal, _ = First(arr, func(i int) bool {
		return arr[i] < "1"
	})
	if !reflect2.IsNil(resultVal) {
		t.Error(fmt.Sprintf("assert faild"))
	}
}

func TestFind(t *testing.T) {
	arr := make([]string, 0, 8)
	arr = append(arr, "a", "b", "c")

	result := Find(arr, func(i int) bool {
		return arr[i] != "a"
	}).([]string)
	if n := len(result); n != 2 {
		t.Errorf("assert faild: %v", n)
	} else {
		if result[0] != "b" {
			t.Errorf("assert faild: %v", result[0])
		}
		if result[1] != "c" {
			t.Errorf("assert faild: %v", result[1])
		}
	}

	result = Find(arr, func(i int) bool {
		return arr[i] == "not exist"
	}).([]string)
	if n := len(result); n != 0 {
		t.Errorf("assert faild: %v", n)
	}
}
