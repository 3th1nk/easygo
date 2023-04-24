package timeRoundedCounter

import (
	"github.com/3th1nk/easygo/util/arrUtil"
	"github.com/3th1nk/easygo/util/mathUtil"
	"sort"
	"sync"
	"time"
)

func NewProperty(name string, round time.Duration, capacity int) *PropertyCounter {
	c := &PropertyCounter{autoSortProperty: true}
	c.base = doNew(name, round, capacity, func() CounterItem {
		return &propertyCountItem{p: c, val: make(map[string]int, 2)}
	})
	registerNamedCounter(c)
	return c
}

type PropertyCounter struct {
	base             *CustomCounter
	properties       []string
	autoSortProperty bool
	mu               sync.RWMutex
}

func (this *PropertyCounter) Name() string {
	return this.base.name
}

func (this *PropertyCounter) Round() time.Duration {
	return this.base.round
}

func (this *PropertyCounter) ItemProperties() []string {
	return this.properties
}

func (this *PropertyCounter) SetPropertyNames(autoSort bool, property ...string) *PropertyCounter {
	this.autoSortProperty = autoSort
	this.properties = property
	return this
}

func (this *PropertyCounter) Add(property string, n int, t ...time.Time) {
	this.addF(property, func(v int) int { return v + n }, t...)
}

func (this *PropertyCounter) Set(property string, n int, t ...time.Time) {
	this.addF(property, func(_ int) int { return n }, t...)
}

func (this *PropertyCounter) Max(property string, n int, t ...time.Time) {
	this.addF(property, func(v int) int { return mathUtil.MaxInt(v, n) }, t...)
}

func (this *PropertyCounter) Min(property string, n int, t ...time.Time) {
	this.addF(property, func(v int) int { return mathUtil.MinInt(v, n) }, t...)
}

func (this *PropertyCounter) addF(property string, f func(v int) int, t ...time.Time) {
	this.base.Add(func(a CounterItem) {
		x := a.(*propertyCountItem)

		this.mu.Lock()
		defer this.mu.Unlock()

		v, ok := x.val[property]
		x.val[property] = f(v)
		if !ok && !arrUtil.ContainsString(this.properties, property) {
			this.properties = append(this.properties, property)
			if this.autoSortProperty {
				sort.Strings(this.properties)
			}
		}
	}, t...)
}

func (this *PropertyCounter) GetAll(rtrimZero ...bool) []CounterItem {
	return this.base.GetAll(rtrimZero...)
}

type propertyCountItem struct {
	p   *PropertyCounter
	val map[string]int
}

func (this *propertyCountItem) IsZero() bool {
	this.p.mu.RLock()
	defer this.p.mu.RUnlock()

	for _, v := range this.val {
		if v != 0 {
			return false
		}
	}
	return true
}

func (this *propertyCountItem) Reset() {
	this.p.mu.RLock()
	defer this.p.mu.RUnlock()

	for k := range this.val {
		this.val[k] = 0
	}
}

func (this *propertyCountItem) MRTValues() []int {
	arr := make([]int, len(this.p.properties))
	for i, property := range this.p.properties {
		arr[i] = this.val[property]
	}
	return arr
}

func (this *propertyCountItem) Values() []interface{} {
	src := this.MRTValues()
	dst := make([]interface{}, len(src))
	for i, v := range src {
		dst[i] = v
	}
	return dst
}
