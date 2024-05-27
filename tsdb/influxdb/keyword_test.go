package influxdb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_usingFunction(t *testing.T) {
	for k, v := range map[string]bool{
		"*":                   false,
		"count":               false,
		"count(\"a\"":         false,
		"count(\"a\")":        true,
		"count(\"a\") as num": true,
	} {
		assert.Equal(t, v, usingFunction(k))
	}
}
