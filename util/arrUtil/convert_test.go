package arrUtil

import (
	"fmt"
	"github.com/modern-go/reflect2"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestStrToInt64(t *testing.T) {
	arr, err := ToInt64([]string{"1", "2", "3"})
	if !reflect2.IsNil(err) {
		t.Error(err)
		return
	}
	if len(arr) != 3 || arr[0] != 1 || arr[1] != 2 || arr[2] != 3 {
		t.Errorf("assert faild: %v", arr)
	}
}

func TestSliceToStr1(t *testing.T) {
	in := []int{1, 2, 3}
	arr, err := ToStr(in)
	if !reflect2.IsNil(err) {
		t.Error(err)
		return
	}
	if len(arr) != 3 || arr[0] != "1" || arr[1] != "2" || arr[2] != "3" {
		t.Errorf("assert faild: %v", arr)
	}
}

func TestSliceToStr2(t *testing.T) {
	in := []int{1, 2, 3}
	arr := ToStrNoError(in, func(i int) string {
		return fmt.Sprintf("%v%v", in[i], in[i])
	})
	if len(arr) != 3 || arr[0] != "11" || arr[1] != "22" || arr[2] != "33" {
		t.Errorf("assert faild: %v", arr)
	}
}

func TestSliceToInt(t *testing.T) {
	in := []string{"1", "2", "3"}
	out, err := ToInt(in, nil)
	if !reflect2.IsNil(err) {
		t.Error(err)
		return
	}
	if out[0] != 1 || out[1] != 2 || out[2] != 3 {
		t.Errorf("assert faild: %v", out)
	}
}

func TestToInterface(t *testing.T) {
	{
		arr := ToInterface([]int{1, 2, 3})
		assert.Equal(t, 3, len(arr))
		assert.Equal(t, 1, arr[0])
		assert.Equal(t, 2, arr[1])
		assert.Equal(t, 3, arr[2])
	}

	{
		arr := ToInterface([]string{"a", "b", "c"})
		assert.Equal(t, 3, len(arr))
		assert.Equal(t, "a", arr[0])
		assert.Equal(t, "b", arr[1])
		assert.Equal(t, "c", arr[2])
	}

	type tmp struct {
		Id   int
		Name string
	}

	{
		arr := ToInterface([]*tmp{
			{Id: 1, Name: "a"},
			{Id: 2, Name: "b"},
			{Id: 3, Name: "c"},
		}, func(a interface{}) interface{} {
			return a.(*tmp).Name
		})
		assert.Equal(t, 3, len(arr))
		assert.Equal(t, "a", arr[0])
		assert.Equal(t, "b", arr[1])
		assert.Equal(t, "c", arr[2])
	}

	{
		arr := ToInterface([]tmp{
			{Id: 1, Name: "a"},
			{Id: 2, Name: "b"},
			{Id: 3, Name: "c"},
		}, func(a interface{}) interface{} {
			return a.(tmp).Name
		})
		assert.Equal(t, 3, len(arr))
		assert.Equal(t, "a", arr[0])
		assert.Equal(t, "b", arr[1])
		assert.Equal(t, "c", arr[2])
	}
}

func TestToSliceType(t *testing.T) {
	arr := make([]interface{}, 4)
	for i := range arr {
		arr[i] = int64(i)
	}

	intArr, _ := ToType(arr, reflect.TypeOf(int64(0))).([]int64)
	assert.Equal(t, len(arr), len(intArr))
	for i := range intArr {
		assert.Equal(t, arr[i], intArr[i])
	}
}

func ExampleToType() {
	arr := make([]interface{}, 4)
	for i := range arr {
		arr[i] = int64(i)
	}

	// 下面代码将会 Panic，因为无法从 []interface{} 直接转换为 []int64，即使切片内的元素类型相同。
	// intArr := arr.([]int64)

	// 通过 ToType(dstType=typeof(int64)) 可以实现转换。
	intArr, _ := ToType(arr, reflect.TypeOf(int64(0))).([]int64)
	if len(arr) != len(intArr) {
		panic("length not equal")
	}
	for i := range intArr {
		if arr[i] != intArr[i] {
			panic(fmt.Sprintf("val[%v] not equal", i))
		}
	}
}
