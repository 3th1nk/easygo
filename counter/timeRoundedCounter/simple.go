package timeRoundedCounter

import (
	"time"
)

func NewSimple(name string, round time.Duration, capacity int) *SimpleCounter {
	c := &SimpleCounter{}
	c.base = doNew(name, round, capacity, func() CounterItem { return &simpleCountItem{} })
	registerNamedCounter(c)
	return c
}

func (this *SimpleCounter) ItemProperties() []string {
	return []string{"Count"}
}

func (this *SimpleCounter) Add(n int, t ...time.Time) {
	this.base.Add(func(a CounterItem) {
		a.(*simpleCountItem).cnt += n
	}, t...)
}

func (this *SimpleCounter) GetAll(rtrimZero ...bool) []CounterItem {
	return this.base.GetAll(rtrimZero...)
}

type SimpleCounter struct {
	base *CustomCounter
}

func (this *SimpleCounter) Name() string {
	return this.base.name
}

func (this *SimpleCounter) Round() time.Duration {
	return this.base.round
}

type simpleCountItem struct {
	cnt int
}

func (this *simpleCountItem) IsZero() bool {
	return this.cnt != 0
}

func (this *simpleCountItem) Reset() {
	this.cnt = 0
}

func (this *simpleCountItem) Values() []interface{} {
	return []interface{}{this.cnt}
}
