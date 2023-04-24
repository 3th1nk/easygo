package timeRoundedCounter

import (
	"fmt"
	"time"
)

func NewStepMRT(name string, round time.Duration, capacity int, step ...time.Duration) *StepMRTCounter {
	c := &StepMRTCounter{step: step}
	c.base = doNew(name, round, capacity, func() CounterItem {
		val := make([]*mrtDurValue, len(step)+1)
		for i := range val {
			val[i] = &mrtDurValue{}
		}
		return &stepMrtCountItem{p: c, val: val}
	})
	registerNamedCounter(c)
	return c
}

type StepMRTCounter struct {
	base       *CustomCounter
	step       []time.Duration
	stepFormat func(step time.Duration) string
}

func (this *StepMRTCounter) Name() string {
	return this.base.name
}

func (this *StepMRTCounter) Round() time.Duration {
	return this.base.round
}

func (this *StepMRTCounter) Steps() []time.Duration {
	return this.step
}

func (this *StepMRTCounter) SetStepFormat(f func(step time.Duration) string) *StepMRTCounter {
	this.stepFormat = f
	return this
}

func (this *StepMRTCounter) ItemProperties() []string {
	stepCount := len(this.step)
	if stepCount == 0 {
		return []string{"MRT"}
	}

	stepFormat := this.stepFormat
	if stepFormat == nil {
		stepFormat = func(step time.Duration) string {
			if step < time.Millisecond {
				return fmt.Sprintf("%dus", step.Microseconds())
			} else if step < time.Second {
				return fmt.Sprintf("%dms", step.Milliseconds())
			}
			return fmt.Sprintf("%.2fs", step.Seconds())
		}
	}

	arr, maxStep := make([]string, stepCount+2), ""
	for i, step := range this.step {
		if step == 0 {
			arr[i] = "0"
		} else {
			str := stepFormat(step)
			arr[i] = "≤ " + str
			maxStep = str
		}
	}
	arr[stepCount] = "> " + maxStep
	arr[stepCount+1] = "Total"
	return arr
}

func (this *StepMRTCounter) Add(n int, d time.Duration, t ...time.Time) {
	this.base.Add(func(a CounterItem) {
		x := a.(*stepMrtCountItem)
		for i, step := range this.step {
			if d <= step {
				x.val[i].add(n, d)
				return
			}
		}
		x.val[len(this.step)].add(n, d)
	}, t...)
}

func (this *StepMRTCounter) GetAll(rtrimZero ...bool) []CounterItem {
	return this.base.GetAll(rtrimZero...)
}

func (this *StepMRTCounter) GetAllMRT(rtrimZero ...bool) []MRTCounterItem {
	src := this.base.GetAll(rtrimZero...)
	dst := make([]MRTCounterItem, len(src))
	for i, v := range src {
		dst[i] = v.(MRTCounterItem)
	}
	return dst
}

type stepMrtCountItem struct {
	p   *StepMRTCounter
	val []*mrtDurValue
}

func (this *stepMrtCountItem) IsZero() bool {
	for _, v := range this.val {
		if v.cnt != 0 {
			return false
		}
	}
	return true
}

func (this *stepMrtCountItem) Reset() {
	for _, v := range this.val {
		v.cnt, v.dur = 0, 0
	}
}

func (this *stepMrtCountItem) MRTValues() []*MRTCountValue {
	if stepCount := len(this.p.step); stepCount == 0 {
		return []*MRTCountValue{this.val[0].toCountValue()}
	} else {
		valItems, total := make([]*MRTCountValue, stepCount+2), mrtDurValue{}
		for i, v := range this.val {
			valItems[i] = v.toCountValue()
			if v.dur > 0 {
				total.add(v.cnt, v.dur)
			}
		}
		valItems[stepCount+1] = total.toCountValue()
		return valItems
	}
}

func (this *stepMrtCountItem) Values() []interface{} {
	src := this.MRTValues()
	dst := make([]interface{}, len(src))
	for i, v := range src {
		dst[i] = v
	}
	return dst
}
