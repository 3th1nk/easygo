package arrUtil

import (
	"github.com/3th1nk/easygo/util/strUtil"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestUnionString(t *testing.T) {
	assert.Equal(t, "", strings.Join(UnionString(nil, nil), ","))
	assert.Equal(t, "", strings.Join(UnionString([]string{}, nil), ","))
	assert.Equal(t, "", strings.Join(UnionString(nil, []string{}), ","))
	assert.Equal(t, "", strings.Join(UnionString([]string{}, []string{}), ","))

	assert.Equal(t, "a,b", strings.Join(UnionString([]string{"a", "b"}, nil), ","))
	assert.Equal(t, "a,b", strings.Join(UnionString([]string{"a", "b"}, []string{"a"}), ","))
	assert.Equal(t, "a,b", strings.Join(UnionString([]string{"a", "b"}, []string{"b"}), ","))
	assert.Equal(t, "a,b", strings.Join(UnionString([]string{"a", "b"}, []string{"b", "a"}), ","))

	assert.Equal(t, "a,b", strings.Join(UnionString(nil, []string{"a", "b"}), ","))
	assert.Equal(t, "a,b", strings.Join(UnionString([]string{"a"}, []string{"a", "b"}), ","))
	assert.Equal(t, "b,a", strings.Join(UnionString([]string{"b"}, []string{"a", "b"}), ","))
	assert.Equal(t, "b,a", strings.Join(UnionString([]string{"b", "a"}, []string{"a", "b"}), ","))
}

func TestUnionInt(t *testing.T) {
	assert.Equal(t, "1,2,3", strUtil.JoinInt(UnionInt([]int{1, 2, 3}, nil), ","))
	assert.Equal(t, "1,2,3", strUtil.JoinInt(UnionInt([]int{1, 2, 3}, []int{1, 3}), ","))
	assert.Equal(t, "1,2,3", strUtil.JoinInt(UnionInt([]int{1, 2, 3}, []int{3, 1}), ","))
	assert.Equal(t, "1,2,3", strUtil.JoinInt(UnionInt([]int{1, 2, 3}, []int{3}, []int{3}, []int{2}), ","))

	assert.Equal(t, "1,2,3,4", strUtil.JoinInt(UnionInt([]int{1, 2, 3}, []int{3, 1, 4}), ","))
	assert.Equal(t, "1,2,3,4,5,6", strUtil.JoinInt(UnionInt([]int{1, 2, 3}, []int{3, 1, 4}, []int{1, 1, 5}, []int{1, 1, 5, 5, 3, 2, 6}), ","))
}
