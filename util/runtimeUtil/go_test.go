package runtimeUtil

import (
	"github.com/3th1nk/easygo/util"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGoWait(t *testing.T) {
	wait := []time.Duration{
		10 * time.Millisecond,
		20 * time.Millisecond,
		30 * time.Millisecond,
		40 * time.Millisecond,
		50 * time.Millisecond,
	}
	waitF := make([]func(d time.Duration) bool, len(wait))
	for i, v := range wait {
		v1 := v
		waitF[i] = func(d time.Duration) bool {
			time.Sleep(v1)
			return v1 <= d
		}
	}

	test := func(d time.Duration, expectAllDone bool, maxTook ...time.Duration) {
		maxTook = append(maxTook, d+20*time.Millisecond)
		start := time.Now()
		allDone := GoWait(d, len(waitF), func(i int, d time.Duration) (done bool) {
			return waitF[i](d)
		})
		took := time.Now().Sub(start)
		assert.Equal(t, expectAllDone, allDone)
		assert.LessOrEqual(t, took.Milliseconds(), maxTook[0].Milliseconds())
		util.PrintTimeLn("wait(%vms), allDone=%v, took=%v", d.Milliseconds(), allDone, took)
	}

	test(10*time.Millisecond, false)
	test(20*time.Millisecond, false)
	test(30*time.Millisecond, false)
	test(40*time.Millisecond, false)
	test(50*time.Millisecond, false)
	test(70*time.Millisecond, true, 70*time.Millisecond)
	test(80*time.Millisecond, true, 70*time.Millisecond)
	test(90*time.Millisecond, true, 70*time.Millisecond)
}

func TestGoWaitN(t *testing.T) {
	arr := make([]int, 200)
	GoWait(time.Second, len(arr), func(i int, d time.Duration) (done bool) {
		arr[i] = i
		return true
	})
	for i, v := range arr {
		assert.Equal(t, i, v)
	}
}
