package arrUtil

import (
	"github.com/3th1nk/easygo/util/strUtil"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_Distinct(t *testing.T) {
	assert.Equal(t, "1,0,2", strUtil.Join(DistinctInt([]int{1, 0, 1, 2}), ","))

	assert.Equal(t, "3,1,2,5,4,0", strUtil.Join(DistinctInt([]int{3, 1, 2, 5, 3, 4, 0, 3, 2, 1}), ","))

	assert.Equal(t, "3,1,2,5,4", strUtil.Join(DistinctString(strings.Split("3,1,2,5,3,4,2,1", ",")), ","))
	assert.Equal(t, "a,b,d,,c,A", strUtil.Join(DistinctString(strings.Split("a,b,a,b,d,,a,c,,A", ",")), ","))
	assert.Equal(t, "a,B,,d,c", strUtil.Join(DistinctString(strings.Split("a,B,,a,,b,,d,a,,c,A", ","), true), ","))

	assert.Equal(t, "1,3,5", strUtil.Join(DistinctSortedInt([]int{1, 1, 1, 1, 3, 3, 5, 5, 5, 5}), ","))

	assert.Equal(t, "a,A,,b", strUtil.Join(DistinctSortedString([]string{"a", "A", "", "b"}, false), ","))
	assert.Equal(t, "a,b", strUtil.Join(DistinctSortedString([]string{"a", "A", "b"}, true), ","))

	assert.Equal(t, "a,A,B,b,c,d", strUtil.Join(DistinctSortedString(strings.Split("a,a,A,B,b,b,c,d,d,d", ","), false), ","))
	assert.Equal(t, "a,B,c,d", strUtil.Join(DistinctSortedString(strings.Split("a,a,A,B,b,b,c,d,d,d", ","), true), ","))
}
