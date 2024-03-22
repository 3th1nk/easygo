package cyclelist

import "sync/atomic"

func New(size int) *CycleList {
	return &CycleList{cap: int64(size), data: make([]interface{}, size)}
}

// 循环链表
type CycleList struct {
	cap  int64
	end  int64
	data []interface{}
}

func (this *CycleList) Size() int {
	return int(this.cap)
}

func (this *CycleList) Add(a interface{}) {
	idx := atomic.AddInt64(&this.end, 1) % this.cap
	this.data[idx] = a
}

func (this *CycleList) Walk(f func(a interface{})) {
	stop := this.end % this.cap
	pos := (stop + 1) % this.cap
	for {
		if v := this.data[pos]; v != nil {
			f(v)
		}
		if pos == stop {
			break
		}
		pos = (pos + 1) % this.cap
	}
}

func (this *CycleList) ReverseWalk(f func(a interface{})) {
	pos := this.end % this.cap
	stop := (pos + 1) % this.cap
	for {
		v := this.data[pos]
		if v == nil {
			break
		}

		f(v)

		if pos == stop {
			break
		} else if pos = pos - 1; pos == -1 {
			pos = this.cap - 1
		}
	}
}
