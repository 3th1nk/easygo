package timeUtil

// ------------------------------------------------------------------------------
// 对计时器的一个扩展：
//    1、允许修改计时器的间隔。这样如果某个计时器的间隔是基于配置文件的，可以在配置文件变更后直接修改间隔。
//    2、允许手动触发一次 Tick、允许指定第一次触发 Tick 的延期时间（默认的计时器是一个计时周期）。有时候需要在创建计时器的时候立即或者在很短时间内就先触发一次，比如用于更新缓存的计时器。
// ------------------------------------------------------------------------------

import (
	"github.com/3th1nk/easygo/util/osUtil"
	"github.com/3th1nk/easygo/util/runtimeUtil"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

func NewTicker(interval, after time.Duration, f func(t time.Time)) *Ticker {
	obj := &Ticker{du: interval, f: f, last: time.Now().Add(after - interval)}
	activeTicker.Store(obj, obj)
	obj.Restart(after)
	return obj
}

var (
	activeTicker sync.Map
)

func init() {
	osUtil.OnSignalExitNoLog(func(_ syscall.Signal) {
		activeTicker.Range(func(_, val interface{}) bool {
			val.(*Ticker).Stop(0)
			return true
		})
	})
}

type Ticker struct {
	mu      sync.RWMutex
	du      time.Duration
	ticker  *innerTicker
	running int64
	last    time.Time
	f       func(t time.Time)
}

// 获取 Ticker 间隔
func (this *Ticker) GetDuration() time.Duration {
	return this.du
}

// 设置 Ticker 间隔
func (this *Ticker) SetDuration(d time.Duration) {
	if d != this.du {
		this.du = d
		this.Restart(this.last.Add(d).Sub(time.Now()))
	}
}

// 手动触发
func (this *Ticker) Trigger() {
	this.doTrigger(time.Now())
}

func (this *Ticker) doTrigger(t time.Time) {
	this.last = t
	atomic.AddInt64(&this.running, 1)
	defer func() {
		atomic.AddInt64(&this.running, -1)
		runtimeUtil.HandleRecover("timeUtil.Ticker panic", recover())
	}()
	this.f(t)
}

// 停止
//
// 参数：
//   d: 最长等待时间。
func (this *Ticker) Stop(d time.Duration) bool {
	this.mu.Lock()
	defer this.mu.Unlock()

	if this.ticker != nil {
		this.ticker.stop()
		this.ticker = nil

		activeTicker.Delete(this)
	}

	if this.running == 0 {
		return true
	} else if d <= 0 {
		return false
	} else {
		expired := false
		time.AfterFunc(d, func() { expired = true })
		for !expired && this.running != 0 {
			time.Sleep(10 * time.Millisecond)
		}
		return this.running == 0
	}
}

func (this *Ticker) Restart(after time.Duration) *Ticker {
	this.mu.Lock()
	defer this.mu.Unlock()

	if this.ticker != nil {
		this.ticker.stop()
	}
	this.ticker = startTicker(this.du, after, this.doTrigger)

	return this
}

type innerTicker struct {
	t        *time.Ticker
	stopChan chan bool
	stopped  bool
}

func startTicker(d, after time.Duration, f func(time.Time)) *innerTicker {
	obj := &innerTicker{stopChan: make(chan bool, 1)}
	time.AfterFunc(after, func() {
		if obj.stopped {
			return
		}
		obj.t = time.NewTicker(d)
		go func() {
			defer obj.t.Stop()
			for {
				select {
				case t := <-obj.t.C:
					f(t)
				case <-obj.stopChan:
					return
				}
			}
		}()
		f(time.Now())
	})
	return obj
}

func (this *innerTicker) stop() {
	this.stopped = true
	this.stopChan <- true
}
