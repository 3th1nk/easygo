package arrUtil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// 测试 在截取切片时，切片的 cap 不会改动
func TestSliceCap(t *testing.T) {
	a := make([]int, 0, 16)
	a = append(a, 1, 2, 3, 4, 5)

	b := a[:len(a)-1]
	assert.Equal(t, 4, len(b))
	assert.Equal(t, 16, cap(b))

	c := b[1:]
	assert.Equal(t, 3, len(c))
	assert.Equal(t, 15, cap(c))
}

func TestSlice(t *testing.T) {
	rawArr := []int{1, 2, 3, 4, 5}

	if arr := rawArr[0:0]; len(arr) != 0 {
		t.Error("assert faild")
	}

	if arr := rawArr[1:1]; len(arr) != 0 {
		t.Error("assert faild")
	}

	if arr := rawArr[4:4]; len(arr) != 0 {
		t.Error("assert faild")
	}

	if arr := rawArr[5:5]; len(arr) != 0 {
		t.Error("assert faild")
	}

	if arr := rawArr[1:2]; len(arr) != 1 || arr[0] != 2 {
		t.Error("assert faild")
	}

	if arr := rawArr[2:4]; len(arr) != 2 || arr[0] != 3 || arr[1] != 4 {
		t.Error("assert faild")
	}

	if arr := append(rawArr[0:0], 1); len(arr) != 1 || arr[0] != 1 {
		t.Error("assert faild")
	}

	if arr := append(rawArr[4:4], 1); len(arr) != 1 || arr[0] != 1 {
		t.Error("assert faild")
	}

	if arr := append(rawArr[0:0], rawArr[4:4]...); len(arr) != 0 {
		t.Error("assert faild")
	}
}

func TestSliceModify(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for i := len(arr) - 1; i >= 0; i-- {
		if arr[i]%3 == 0 {
			arr = append(arr[:i], arr[i+1:]...)
		}
	}
	if len(arr) != 6 {
		t.Error("assert faild")
	}
	if arr[0] != 1 || arr[1] != 2 || arr[2] != 4 || arr[3] != 5 || arr[4] != 7 || arr[5] != 8 {
		t.Error("assert faild")
	}
}
