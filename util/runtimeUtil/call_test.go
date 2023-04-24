package runtimeUtil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCallFunc(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			// defer 后直接跟 testHandler，testHandler 内的 recover 是可以捕获 panic 的。
			// 故不应该再走到这里
			t.Errorf("这里不应该捕获到 panic 的")
		}
	}()

	p := CallFunc(func() {
		panic("trigger panic")
	})
	assert.Equal(t, "trigger panic", p)
}

func TestCallErrFunc(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			// defer 后直接跟 testHandler，testHandler 内的 recover 是可以捕获 panic 的。
			// 故不应该再走到这里
			t.Errorf("这里不应该捕获到 panic 的")
		}
	}()

	_, p := CallFuncE(func() error {
		panic("trigger panic")
	})
	assert.Equal(t, "trigger panic", p)
}
