package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIfInt(t *testing.T) {
	assert.Equal(t, 1, IfInt(true, 1, 2))
}
