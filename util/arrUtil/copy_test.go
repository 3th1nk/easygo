package arrUtil

import (
	"github.com/3th1nk/easygo/util/convertor"
	"testing"
)

// 验证 copy 函数是浅拷贝
func TestCopy(t *testing.T) {
	type abc struct {
		Id   int
		Name string
	}

	arr1 := []*abc{
		{Id: 1, Name: "a"},
		{Id: 2, Name: "b"},
		{Id: 3, Name: "c"},
		{Id: 4, Name: "d"},
		{Id: 5, Name: "e"},
	}
	arr2 := make([]*abc, len(arr1))
	copy(arr2, arr1)

	for _, item := range arr2 {
		item.Name = item.Name + item.Name + item.Name
	}

	str1 := convertor.ToStringNoError(arr1)
	str2 := convertor.ToStringNoError(arr2)
	if str1 != str2 {
		t.Errorf("assert faild: \n  %v\n  %v", str1, str2)
	}
}
