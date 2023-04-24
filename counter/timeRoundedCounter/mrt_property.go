package timeRoundedCounter

import (
	"github.com/3th1nk/easygo/util/arrUtil"
	"sort"
	"sync"
	"time"
)

func NewPropertyMRT(name string, round time.Duration, capacity int) *PropertyMRTCounter {
	c := &PropertyMRTCounter{autoSortProperty: true}
	c.base = doNew(name, round, capacity, func() CounterItem {
		return &PropertyMrtCountItem{p: c, val: make(map[string]*mrtDurValue, 2)}
	})
	registerNamedCounter(c)
	return c
}

type PropertyMRTCounter struct {
	base             *CustomCounter
	properties       []string
	autoSortProperty bool
	mu               sync.RWMutex
}

func (this *PropertyMRTCounter) Name() string {
	return this.base.name
}

func (this *PropertyMRTCounter) Round() time.Duration {
	return this.base.round
}

func (this *PropertyMRTCounter) ItemProperties() []string {
	return this.properties
}

func (this *PropertyMRTCounter) SetPropertyNames(autoSort bool, property ...string) *PropertyMRTCounter {
	this.autoSortProperty = autoSort
	this.properties = property
	return this
}

func (this *PropertyMRTCounter) Add(property string, n int, d time.Duration, t ...time.Time) {
	this.base.Add(func(v CounterItem) {
		v.(*PropertyMrtCountItem).Add(property, n, d)
	}, t...)
}

func (this *PropertyMRTCounter) Set(property string, n int, d time.Duration, t ...time.Time) {
	this.base.Add(func(v CounterItem) {
		v.(*PropertyMrtCountItem).Set(property, n, d)
	}, t...)
}

func (this *PropertyMRTCounter) AddF(f func(a *PropertyMrtCountItem)) {
	this.base.Add(func(v CounterItem) {
		f(v.(*PropertyMrtCountItem))
	})
}

func (this *PropertyMRTCounter) GetAll(rtrimZero ...bool) []CounterItem {
	return this.base.GetAll(rtrimZero...)
}

func (this *PropertyMRTCounter) GetAllMRT(rtrimZero ...bool) []MRTCounterItem {
	src := this.base.GetAll(rtrimZero...)
	dst := make([]MRTCounterItem, len(src))
	for i, v := range src {
		dst[i] = v.(MRTCounterItem)
	}
	return dst
}

type PropertyMrtCountItem struct {
	p   *PropertyMRTCounter
	val map[string]*mrtDurValue
}

func (this *PropertyMrtCountItem) Add(property string, n int, d time.Duration) {
	this.addF(property, func(v *mrtDurValue) {
		v.cnt, v.dur = v.cnt+n, v.dur+d
	})
}

func (this *PropertyMrtCountItem) Set(property string, n int, d time.Duration) {
	this.addF(property, func(v *mrtDurValue) {
		v.cnt, v.dur = n, d
	})
}

func (this *PropertyMrtCountItem) addF(property string, f func(v *mrtDurValue)) {
	this.p.mu.Lock()
	defer this.p.mu.Unlock()

	v := this.val[property]
	if v == nil {
		v = &mrtDurValue{}
		this.val[property] = v

		if !arrUtil.ContainsString(this.p.properties, property) {
			this.p.properties = append(this.p.properties, property)
			if this.p.autoSortProperty {
				sort.Strings(this.p.properties)
			}
		}
	}
	f(v)
}

func (this *PropertyMrtCountItem) IsZero() bool {
	this.p.mu.RLock()
	defer this.p.mu.RUnlock()

	for _, v := range this.val {
		if v.cnt != 0 {
			return false
		}
	}
	return true
}

func (this *PropertyMrtCountItem) Reset() {
	this.p.mu.RLock()
	defer this.p.mu.RUnlock()

	for _, v := range this.val {
		v.cnt, v.dur = 0, 0
	}
}

func (this *PropertyMrtCountItem) Properties() []string {
	properties := make([]string, len(this.p.properties))
	copy(properties, this.p.properties)
	return properties
}

func (this *PropertyMrtCountItem) MRTValues() []*MRTCountValue {
	arr := make([]*MRTCountValue, len(this.p.properties))
	for i, property := range this.p.properties {
		if v := this.val[property]; v != nil {
			arr[i] = v.toCountValue()
		} else {
			arr[i] = &MRTCountValue{}
		}
	}
	return arr
}

func (this *PropertyMrtCountItem) Values() []interface{} {
	src := this.MRTValues()
	dst := make([]interface{}, len(src))
	for i, v := range src {
		dst[i] = v
	}
	return dst
}
