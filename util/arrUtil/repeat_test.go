package arrUtil

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Repeat_Int(t *testing.T) {
	n := 100

	{
		arr, ok := Repeat(n, 10).([]int)
		if !assert.True(t, ok) {
			return
		}
		for _, item := range arr {
			assert.Equal(t, item, n)
		}
	}

	{
		arr, ok := Repeat(n, 10, func(i int, a interface{}) interface{} {
			return n + i
		}).([]int)
		if !assert.True(t, ok) {
			return
		}
		for idx, item := range arr {
			assert.Equal(t, item, n+idx)
		}
	}
}

func Test_Repeat(t *testing.T) {
	type type1 struct {
		n int
	}
	obj := &type1{n: 100}
	arr, ok := Repeat(obj, 10, func(i int, a interface{}) interface{} {
		a.(*type1).n = obj.n + i
		return a
	}).([]*type1)
	if !assert.True(t, ok) {
		return
	}
	assert.Equal(t, arr[0].n, obj.n)
	assert.Equal(t, arr[1].n, obj.n+1)
	assert.Equal(t, arr[2].n, obj.n+2)
	assert.Equal(t, arr[3].n, obj.n+3)
}

func TestRepeatInt(t *testing.T) {
	n := 100

	{
		arr := RepeatInt(100, 10)
		for _, item := range arr {
			assert.Equal(t, item, n)
		}
	}

	{
		arr := RepeatInt(n, 10, func(i int, a int) int {
			return a + i
		})
		assert.Equal(t, arr[0], n+0)
		assert.Equal(t, arr[1], n+1)
	}
}

func TestRepeatString(t *testing.T) {
	s := "abc"

	{
		arr := RepeatString(s, 10)
		for _, item := range arr {
			assert.Equal(t, item, s)
		}
	}

	{
		arr := RepeatString(s, 10, func(i int, s string) string {
			return fmt.Sprintf("%s.%d", s, i)
		})
		assert.Equal(t, arr[0], s+".0")
		assert.Equal(t, arr[1], s+".1")
	}
}
