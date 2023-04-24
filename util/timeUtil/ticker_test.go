package timeUtil

import (
	"github.com/3th1nk/easygo/util"
	"github.com/stretchr/testify/assert"
	"math"
	"sync"
	"testing"
	"time"
)

func TestUtilTicker(t *testing.T) {
	start := time.Now()
	util.Println("start: %v", start.Format("04:05.999"))

	var arr []time.Time

	d1 := 300 * time.Millisecond
	s1 := time.Second
	after := 0 * time.Millisecond
	ticker := NewTicker(d1, after, func(t time.Time) {
		arr = append(arr, t)
		util.Println(" tick: %v", t.Format("04:05.999"))
	})
	time.Sleep(s1)
	n1 := int(math.Floor(float64(s1+d1-after) / float64(d1)))
	util.Println("n1=%v", n1)
	assert.Equal(t, n1, len(arr))

	d2 := 500 * time.Millisecond
	s2 := 3 * time.Second
	ticker.SetDuration(d2)
	time.Sleep(s2)
	n2 := n1 + int(math.Floor(float64(s1+s2+d1-after-time.Duration(n1)*d1)/float64(d2)))
	util.Println("n2=%v", n2)
	assert.Equal(t, n2, len(arr))

	assert.Equal(t, 0, int((arr[0].Sub(start) - after).Milliseconds()))

	if !ticker.Stop(time.Second) {
		t.Error("stop error")
	}
	time.Sleep(time.Second)
	if !ticker.Stop(time.Second) {
		t.Error("stop error")
	}
}

// 验证： 如果 timer 间隔设置为 1s、每次触发执行一个 3s 的逻辑，则实际两次触发的间隔是 1s 还是 4s
//   即：验证 timer 的间隔是指 ‘下一次开始 - 上一次开始’ 还是 ‘下一次开始 - 上一次结束’
// 验证结果：
//   间隔为 Max(间隔, 执行耗时)
func TestTimeTicker_Interval(t *testing.T) {
	interval, span, maxLoop := time.Second, 3*time.Second, 3
	wait := sync.WaitGroup{}
	wait.Add(2)
	go func() {
		start := time.Now()
		loop := 0
		ticker := time.NewTicker(interval)
		for range ticker.C {
			time.Sleep(span)
			loop++
			if loop == maxLoop {
				break
			}
		}
		t.Logf("time.Ticker: interval=%v, span=%v, loop=%v, timespan=%v", interval, span, maxLoop, time.Now().Sub(start))
		wait.Done()
	}()
	go func() {
		start := time.Now()
		loop := 0
		NewTicker(interval, 0, func(t time.Time) {
			time.Sleep(span)
			loop++
		})
		for loop < maxLoop {
			time.Sleep(20 * time.Millisecond)
		}
		t.Logf("timeUtil.Ticker: interval=%v, span=%v, loop=%v, timespan=%v", interval, span, maxLoop, time.Now().Sub(start))
		wait.Done()
	}()
	wait.Wait()
}

func TestTimeAfterFunc_1(t *testing.T) {
	stop := false
	time.AfterFunc(2*time.Second, func() {
		tick := time.NewTicker(time.Second)
		for range tick.C {
			util.PrintTimeLn("tick")
		}
		if stop {
			tick.Stop()
		}
	})
	time.AfterFunc(5*time.Second, func() {
		stop = true
	})
	for !stop {
		time.Sleep(200 * time.Millisecond)
		util.PrintTimeLn("    wait...")
	}
}

func TestTimeAfterFunc_2(t *testing.T) {
	wg := sync.WaitGroup{}

	for i, n := 0, 100; i < n; i++ {
		start := time.Now()
		var after time.Time
		wg.Add(1)
		time.AfterFunc(0, func() {
			after = time.Now()
			wg.Done()
		})
		wg.Wait()
		t.Logf("d=%v", after.Sub(start))
	}
}
