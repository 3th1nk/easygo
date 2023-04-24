package _test

import (
	"github.com/3th1nk/easygo/util"
	"github.com/3th1nk/easygo/util/runtimeUtil"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

func Perf(f func(i int), opt ...PerfOptions) *PerfResult {
	return PerfIf(func(i int) (ok bool) {
		f(i)
		return true
	}, opt...)
}

// 在指定时间内，循环调用回调函数，输出统计数据。
func PerfIf(f func(i int) (ok bool), opt ...PerfOptions) *PerfResult {
	var theOpt PerfOptions
	if len(opt) != 0 {
		theOpt = opt[0]
	}
	if theOpt.Dur <= 0 {
		theOpt.Dur = time.Second
	}
	if theOpt.Goroutine <= 0 {
		theOpt.Goroutine = 1
	}

	if !theOpt.NoHead {
		title, _, _, _ := runtimeUtil.Caller(2)
		if theOpt.Name != "" {
			title += ", " + theOpt.Name
		}
		util.PrintTimeLn("============================================== %s:", title)
	}

	result := &PerfResult{}
	start, stop := time.Now(), false
	time.AfterFunc(theOpt.Dur, func() { stop = true })

	run := func() {
		loop := 0
		for !stop {
			ok := f(loop)
			loop++
			if ok {
				atomic.AddInt32(&result.Ok, 1)
			} else {
				atomic.AddInt32(&result.Fail, 1)
			}
		}
	}

	if theOpt.Goroutine == 1 {
		run()
	} else {
		wg := sync.WaitGroup{}
		for i := 0; i < theOpt.Goroutine; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				run()
			}()
		}
		wg.Wait()
	}

	//
	if result.Took = time.Since(start); result.Took >= 10*time.Second {
		result.Took = result.Took.Round(time.Second)
	} else if result.Took >= time.Second {
		result.Took = result.Took.Round(10 * time.Millisecond)
	} else if result.Took >= time.Millisecond {
		result.Took = result.Took.Round(10 * time.Microsecond)
	}
	//
	if result.Ok != 0 {
		if result.MRT = result.Took / time.Duration(result.Ok); result.MRT >= time.Second {
			result.MRT = result.MRT.Round(time.Millisecond)
		} else if result.MRT >= time.Millisecond {
			result.MRT = result.MRT.Round(time.Microsecond)
		}
	}
	//
	result.QPS = math.Round(100*float64(result.Ok)/result.Took.Seconds()) / 100

	util.PrintTimeLn("go=%d,  ok=%d,  fail=%d,  avg=%.f/s,  mrt=%v", theOpt.Goroutine, result.Ok, result.Fail, result.QPS, result.MRT)

	return result
}

type PerfOptions struct {
	Name      string
	Dur       time.Duration
	Goroutine int
	NoHead    bool
}

type PerfResult struct {
	Took time.Duration `json:"took,omitempty"` //
	Ok   int32         `json:"ok,omitempty"`   //
	Fail int32         `json:"fail,omitempty"` //
	MRT  time.Duration `json:"mrt,omitempty"`  // Mean Response Time，平均响应时间
	QPS  float64       `json:"qps,omitempty"`  // Request Per Second，平均每秒请求次数
}
