package strUtil

import (
	"fmt"
	"strings"
)

type Builder struct {
	data     []string
	index    int
	capacity int
}

func NewBuilder(str ...string) *Builder {
	return (&Builder{}).Append(str...)
}

func (this *Builder) Empty() bool { return this.index == 0 }

func (this *Builder) Clear() { this.index = 0 }

func (this *Builder) Append(str ...string) *Builder {
	if n := len(str); n > 0 {
		this.prepareAppend(n)
		for _, s := range str {
			this.data[this.index] = s
			this.index++
		}
	}
	return this
}

func (this *Builder) Appendf(format string, a ...interface{}) *Builder {
	return this.Append(fmt.Sprintf(format, a...))
}

func (this *Builder) Reset() *Builder {
	this.index = 0
	return this
}

func (this *Builder) String() string {
	return strings.Join(this.data, "")
}

func (this *Builder) prepareAppend(n int) {
	if n2 := this.index + n; n2 > this.capacity {
		for n2 > this.capacity {
			if this.capacity == 0 {
				this.capacity = 8
			} else {
				this.capacity *= 2
			}
		}
		if this.index == 0 {
			this.data = make([]string, this.capacity)
		} else {
			newData := make([]string, this.capacity)
			copy(newData, this.data)
			this.data = newData
		}
	}
}
