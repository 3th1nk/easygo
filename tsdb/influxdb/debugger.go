package influxdb

import (
	"fmt"
	"sync/atomic"
)

type counter struct {
	total uint64
	fail  uint64
}

func newCounter() *counter {
	return &counter{}
}

func (c *counter) IncTotal() {
	atomic.AddUint64(&c.total, 1)
}

func (c *counter) IncFail() {
	atomic.AddUint64(&c.fail, 1)
}

func (c *counter) Reset() {
	atomic.StoreUint64(&c.total, 0)
	atomic.StoreUint64(&c.fail, 0)
}

func (c *counter) Total() uint64 {
	return atomic.LoadUint64(&c.total)
}

func (c *counter) Fail() uint64 {
	return atomic.LoadUint64(&c.fail)
}

func (c *counter) String() string {
	return fmt.Sprintf("total: %d, fail: %d", c.Total(), c.Fail())
}

type debugger struct {
	write *counter
}

func newDebugger() *debugger {
	return &debugger{
		write: newCounter(),
	}
}

func (d *debugger) CountWrite(success bool) {
	d.write.IncTotal()
	if !success {
		d.write.IncFail()
	}
}

func (this *Client) countWrite(success bool) {
	if this.debugger != nil {
		this.debugger.CountWrite(success)
	}
}

func (this *Client) showWriteCount() {
	if this.debugger != nil {
		this.logger.Info(this.debugger.write.String())
	}
}
