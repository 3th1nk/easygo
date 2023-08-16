package delayWorker

import (
	"fmt"
	"sync/atomic"
)

type counter struct {
	receive uint64
	drop    uint64
	done    uint64
}

func (c *counter) Receive() uint64 {
	return atomic.LoadUint64(&c.receive)
}

func (c *counter) Drop() uint64 {
	return atomic.LoadUint64(&c.drop)
}

func (c *counter) Done() uint64 {
	return atomic.LoadUint64(&c.done)
}

func (c *counter) IncReceive(n int) {
	atomic.AddUint64(&c.receive, uint64(n))
}

func (c *counter) IncDrop(n int) {
	atomic.AddUint64(&c.drop, uint64(n))
}

func (c *counter) IncDone(n int) {
	atomic.AddUint64(&c.done, uint64(n))
}

func (c *counter) String() string {
	return fmt.Sprintf("receive:%d drop:%d done:%d", c.Receive(), c.Drop(), c.Done())
}
