package cyclelist

import (
	"github.com/3th1nk/easygo/util/strUtil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCycleList_Add_1(t *testing.T) {
	cl := New(10)
	for i := 1; i <= 5; i++ {
		cl.Add(i)
	}

	{
		arr1 := make([]int, 0, cl.Size())
		cl.Walk(func(a interface{}) {
			arr1 = append(arr1, a.(int))
		})
		assert.Equal(t, "1,2,3,4,5", strUtil.JoinInt(arr1, ","))

		arr2 := make([]int, 0, cl.Size())
		cl.ReverseWalk(func(a interface{}) {
			arr2 = append(arr2, a.(int))
		})
		for i, v := range arr1 {
			assert.Equal(t, v, arr2[len(arr2)-1-i])
		}
	}

	for i := 6; i <= 11; i++ {
		cl.Add(i)
	}
	{
		arr1 := make([]int, 0, cl.Size())
		cl.Walk(func(a interface{}) {
			arr1 = append(arr1, a.(int))
		})
		assert.Equal(t, "2,3,4,5,6,7,8,9,10,11", strUtil.JoinInt(arr1, ","))

		arr2 := make([]int, 0, cl.Size())
		cl.ReverseWalk(func(a interface{}) {
			arr2 = append(arr2, a.(int))
		})
		for i, v := range arr1 {
			assert.Equal(t, v, arr2[len(arr2)-1-i])
		}
	}

	for i := 12; i <= 18; i++ {
		cl.Add(i)
	}
	{
		arr1 := make([]int, 0, cl.Size())
		cl.Walk(func(a interface{}) {
			arr1 = append(arr1, a.(int))
		})
		assert.Equal(t, "9,10,11,12,13,14,15,16,17,18", strUtil.JoinInt(arr1, ","))

		arr2 := make([]int, 0, cl.Size())
		cl.ReverseWalk(func(a interface{}) {
			arr2 = append(arr2, a.(int))
		})
		for i, v := range arr1 {
			assert.Equal(t, v, arr2[len(arr2)-1-i])
		}
	}

	for i := 11; i <= 20; i++ {
		cl.Add(i)
	}
	{
		arr1 := make([]int, 0, cl.Size())
		cl.Walk(func(a interface{}) {
			arr1 = append(arr1, a.(int))
		})
		assert.Equal(t, "11,12,13,14,15,16,17,18,19,20", strUtil.JoinInt(arr1, ","))

		arr2 := make([]int, 0, cl.Size())
		cl.ReverseWalk(func(a interface{}) {
			arr2 = append(arr2, a.(int))
		})
		for i, v := range arr1 {
			assert.Equal(t, v, arr2[len(arr2)-1-i])
		}
	}
}

func TestCycleList_Add_2(t *testing.T) {
	cl := New(10)
	for i := 1; i <= 20; i++ {
		cl.Add(i)
	}
}
