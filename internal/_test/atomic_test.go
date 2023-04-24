package _test_test

import (
	"github.com/3th1nk/easygo/internal/_test"
	"sync/atomic"
	"testing"
)

func TestAtomic_Perf(t *testing.T) {
	var n32 int32
	var n64 int64

	_test.Perf(func(i int) {
		atomic.AddInt32(&n32, 1)
	})

	_test.Perf(func(i int) {
		atomic.AddInt64(&n64, 1)
	})

	{
		a := &struct {
			n int64
		}{}
		_test.Perf(func(i int) {
			atomic.AddInt64(&(a.n), 1)
		})
	}

	{
		a := &struct {
			a [1]byte
			n int64
		}{}
		_test.Perf(func(i int) {
			atomic.AddInt64(&(a.n), 1)
		})
	}
}
