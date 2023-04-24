package timeRoundedCounter

import "time"

type MRTCounterItem interface {
	CounterItem

	MRTValues() []*MRTCountValue
}

type MRTCountValue struct {
	Count int
	MRT   time.Duration
}

type mrtDurValue struct {
	cnt int
	dur time.Duration
}

func (this *mrtDurValue) add(n int, d time.Duration) {
	this.cnt += n
	this.dur += d
}

func (this *mrtDurValue) toCountValue() *MRTCountValue {
	if this.cnt == 0 {
		return &MRTCountValue{}
	}

	val := &MRTCountValue{Count: this.cnt, MRT: this.dur / time.Duration(this.cnt)}
	if val.MRT < time.Millisecond {
		val.MRT = val.MRT.Round(time.Microsecond)
	} else if val.MRT < time.Second {
		val.MRT = val.MRT.Round(time.Millisecond)
	} else {
		val.MRT = val.MRT.Round(time.Second)
	}
	return val
}
