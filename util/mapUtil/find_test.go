package mapUtil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testStruct struct {
	Str string
}

func TestFind1(t *testing.T) {
	a := map[string]*testStruct{
		"a": {Str: "a"},
		"b": {Str: "b"},
		"c": {Str: "c"},
		"d": {Str: "d"},
	}
	b := Find(a, func(key, val interface{}) bool {
		return key.(string) != "b"
	}).(map[string]*testStruct)
	assert.Equal(t, 3, len(b))
	assert.NotNil(t, b["a"])
	assert.NotNil(t, b["c"])
	assert.NotNil(t, b["d"])
}

func TestFind2(t *testing.T) {
	a := map[*testStruct]string{
		{Str: "a"}: "a",
		{Str: "b"}: "b",
		{Str: "c"}: "c",
		{Str: "d"}: "d",
	}
	b := Find(a, func(key, val interface{}) bool {
		return val.(string) != "b"
	}).(map[*testStruct]string)
	assert.Equal(t, 3, len(b))
	for _, val := range b {
		assert.NotEqual(t, "b", val)
	}
}

func TestFind3(t *testing.T) {
	a := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
	}
	b := Find(a, func(key, val interface{}) bool {
		return key.(string) != "b"
	}).(map[string]int)
	assert.Equal(t, 3, len(b))
	assert.NotZero(t, b["a"])
	assert.NotZero(t, b["c"])
	assert.NotZero(t, b["d"])
}

func TestFind4(t *testing.T) {
	a := map[int]string{
		1: "a",
		2: "b",
		3: "c",
		4: "d",
	}
	b := Find(a, func(key, val interface{}) bool {
		return val.(string) != "b"
	}).(map[int]string)
	assert.Equal(t, 3, len(b))
	assert.NotZero(t, b[1])
	assert.NotZero(t, b[3])
	assert.NotZero(t, b[4])
}
