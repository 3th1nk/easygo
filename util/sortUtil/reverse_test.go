package sortUtil

import (
	"github.com/3th1nk/easygo/util/strUtil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReverseInt(t *testing.T) {
	var arr []int

	arr = []int{1, 2, 3, 4}
	Reverse(arr)
	assert.Equal(t, "4,3,2,1", strUtil.Join(arr, ","))

	arr = []int{1, 2, 3, 4, 5}
	Reverse(arr)
	assert.Equal(t, "5,4,3,2,1", strUtil.Join(arr, ","))
}

func TestReverseString(t *testing.T) {
	var arr []string

	arr = []string{"1", "2", "3", "4"}
	Reverse(arr)
	assert.Equal(t, "4,3,2,1", strUtil.Join(arr, ","))

	arr = []string{"1", "2", "3", "4", "5"}
	Reverse(arr)
	assert.Equal(t, "5,4,3,2,1", strUtil.Join(arr, ","))
}

func TestReverse(t *testing.T) {
	var arr interface{}

	arr = []int{1, 2, 3, 4}
	Reverse(arr)
	assert.Equal(t, "4,3,2,1", strUtil.Join(arr, ","))

	arr = []int{1, 2, 3, 4, 5}
	Reverse(arr)
	assert.Equal(t, "5,4,3,2,1", strUtil.Join(arr, ","))

	arr = []string{"1", "2", "3", "4"}
	Reverse(arr)
	assert.Equal(t, "4,3,2,1", strUtil.Join(arr, ","))

	arr = []string{"1", "2", "3", "4", "5"}
	Reverse(arr)
	assert.Equal(t, "5,4,3,2,1", strUtil.Join(arr, ","))
}
