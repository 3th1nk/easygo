package runtimeUtil

import (
	"sync"
	"sync/atomic"
	"time"
)

func Go(n int, f func(i int)) {
	if n <= 0 {
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer func() {
				wg.Done()
				HandleRecover("", recover())
			}()
			f(i)
		}(i)
	}
	wg.Wait()
}

// GoWait 并发超时执行，如果超时时间未执行完成，提前返回
func GoWait(wait time.Duration, n int, f func(i int, wait time.Duration) (done bool)) (allDone bool) {
	if n == 0 || f == nil {
		return true
	}

	undone := int32(n)
	doneChan := make(chan bool)
	for i := 0; i < n; i++ {
		go func(i int) {
			if f(i, wait) && atomic.AddInt32(&undone, -1) == 0 {
				close(doneChan)
			}
		}(i)
	}

	if wait > 0 {
		select {
		case <-doneChan:
			return true
		case <-time.After(wait):
			return false
		}
	}

	return undone == 0
}
