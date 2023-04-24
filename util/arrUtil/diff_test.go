package arrUtil

import (
	"github.com/3th1nk/easygo/util/strUtil"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestDiffInt(t *testing.T) {
	a, _ := strUtil.SplitToInt("1,2,3", ",", true)
	b, _ := strUtil.SplitToInt("1,2,4", ",", true)
	matches, added, removed := DiffInt(a, b)
	assert.Equal(t, "1,2", strUtil.JoinInt(matches, ","))
	assert.Equal(t, "4", strUtil.JoinInt(added, ","))
	assert.Equal(t, "3", strUtil.JoinInt(removed, ","))
}

func TestDiffString(t *testing.T) {
	a := strings.Split("a,b,c", ",")
	b := strings.Split("a,c,d", ",")
	matches, added, removed := DiffString(a, b)
	assert.Equal(t, "a,c", strings.Join(matches, ","))
	assert.Equal(t, "d", strings.Join(added, ","))
	assert.Equal(t, "b", strings.Join(removed, ","))
}
