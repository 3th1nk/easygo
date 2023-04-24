package arrUtil

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestStrToInt(t *testing.T) {
	arr, err := StrToInt(strings.Split("1,2,3", ","))
	assert.NoError(t, err)
	assert.Equal(t, 3, len(arr))
	assert.Equal(t, 1, arr[0])
	assert.Equal(t, 2, arr[1])
	assert.Equal(t, 3, arr[2])
}
