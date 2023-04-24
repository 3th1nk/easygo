package timeRoundedCounter

import (
	"fmt"
	"github.com/3th1nk/easygo/util/mathUtil"
	"sync"
	"time"
)

func NewCustom(name string, round time.Duration, capacity int, newItem func() CounterItem) *CustomCounter {
	c := doNew(name, round, capacity, newItem)
	registerNamedCounter(c)
	return c
}

func doNew(name string, round time.Duration, capacity int, newItem func() CounterItem) *CustomCounter {
	if round <= 0 {
		panic(fmt.Errorf("invalid param 'round'"))
	}
	if capacity <= 0 {
		panic(fmt.Errorf("invalid param 'capacity'"))
	}
	if newItem == nil {
		panic(fmt.Errorf("missing param 'newItem'"))
	}

	data := make([]CounterItem, capacity)
	for i := 0; i < capacity; i++ {
		v := newItem()
		v.Reset()
		data[i] = v
	}
	return &CustomCounter{name: name, round: round, data: data, newItem: newItem}
}

type CustomCounter struct {
	name    string
	round   time.Duration
	data    []CounterItem
	newItem func() CounterItem
	maxNS   int64
	maxPos  int
	lock    sync.Mutex
}

func (this *CustomCounter) Name() string { return this.name }

func (this *CustomCounter) Round() time.Duration { return this.round }

func (this *CustomCounter) Capacity() int { return len(this.data) }

// 设置计数器的值
//
// Perf:
//   [perf-1s]  total=14877334,  avg=14877334/s,  mrt=67ns
func (this *CustomCounter) Add(f func(v CounterItem), t ...time.Time) {
	var item CounterItem
	if len(t) != 0 {
		item = this.itemToWrite(t[0])
	} else {
		item = this.itemToWrite(time.Now())
	}
	if item != nil {
		f(item)
	}
}

// Perf:
//   [perf-1s] total=1720411,  avg=1720411/s,  mrt=581ns
func (this *CustomCounter) GetAll(rtrimZero ...bool) []CounterItem {
	return this.GetN(time.Now(), 0, -1, rtrimZero...)
}

func (this *CustomCounter) GetN(t time.Time, offset, limit int, rtrimZero ...bool) []CounterItem {
	if this.maxNS == 0 {
		return nil
	}

	if offset < 0 {
		offset = 0
	}
	capacity := len(this.data)
	if limit < 0 {
		limit = capacity - offset
	} else if n := capacity - offset; limit > n {
		limit = n
	}
	if limit <= 0 {
		return nil
	}

	arr, cnt := make([]CounterItem, limit), 0
	pos := (this.maxPos - offset) % capacity
	if pos < 0 {
		pos += capacity
	}
	overflow := int((t.Truncate(this.round).UnixNano() - this.maxNS) / int64(this.round))
	if overflow > 0 {
		cnt = mathUtil.MinInt(overflow, capacity)
		for i := 0; i < cnt; i++ {
			v := this.newItem()
			v.Reset()
			arr[i] = v
		}
		this.reverseWalk(pos, limit-overflow, func(a CounterItem) {
			arr[cnt], cnt = a, cnt+1
		})
	} else {
		this.reverseWalk(pos, limit, func(a CounterItem) {
			arr[cnt], cnt = a, cnt+1
		})
	}
	if cnt != len(arr) {
		arr = arr[:cnt]
	}

	if len(rtrimZero) == 0 || rtrimZero[0] {
		found := false
		for i := limit - 1; i >= 0; i-- {
			if !arr[i].IsZero() {
				arr, found = arr[:i+1], true
				break
			}
		}
		if !found {
			arr = nil
		}
	}

	return arr
}

func (this *CustomCounter) itemToWrite(t time.Time) (item CounterItem) {
	this.lock.Lock()
	defer this.lock.Unlock()

	capacity := len(this.data)
	ns := t.Truncate(this.round).UnixNano()
	if this.maxNS == 0 {
		this.maxNS = ns
		return this.data[0]
	} else {
		step := int((this.maxNS - ns) / int64(this.round))
		if step >= capacity {
			return nil
		} else {
			pos := (this.maxPos - step) % capacity
			if pos >= capacity {
				pos -= capacity
			} else if pos < 0 {
				pos += capacity
			}
			if step < 0 {
				this.reverseWalk(pos, -step, func(a CounterItem) {
					a.Reset()
				})
				this.maxNS, this.maxPos = ns, pos
			}
			return this.data[pos]
		}
	}
}

func (this *CustomCounter) reverseWalk(pos, limit int, f func(a CounterItem)) {
	for i := 0; i < limit; i++ {
		f(this.data[pos])
		if pos = pos - 1; pos == -1 {
			pos = len(this.data) - 1
		}
	}
}
