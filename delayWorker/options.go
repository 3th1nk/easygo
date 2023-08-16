package delayWorker

import (
	"github.com/3th1nk/easygo/util/logs"
)

type Options struct {
	QueueSize int
	DelaySec  int
	DelayCnt  int
	Parallel  bool
	Logger    logs.Logger

	// 调试模式
	debug   bool
	counter *counter
}

func DefaultOptions() *Options {
	return &Options{
		QueueSize: 1000,
		DelaySec:  5,
		DelayCnt:  100,
		Parallel:  false,
		Logger:    logs.Default,
	}
}

func (opt *Options) ensure() {
	if opt.DelaySec <= 0 {
		opt.DelaySec = 5
	}
	if opt.DelayCnt <= 0 {
		opt.DelayCnt = 100
	}
	if opt.QueueSize <= 0 {
		opt.QueueSize = 1000
	}
	if opt.Logger == nil {
		opt.Logger = logs.Default
	}
}

func (opt *Options) statistics() string {
	if opt.debug && opt.counter != nil {
		return opt.counter.String()
	}
	return ""
}

func (opt *Options) incReceive(n int) {
	if opt.debug {
		if opt.counter == nil {
			opt.counter = &counter{}
		}
		opt.counter.IncReceive(n)
	}
}

func (opt *Options) incDrop(n int) {
	if opt.debug {
		if opt.counter == nil {
			opt.counter = &counter{}
		}
		opt.counter.IncDrop(n)
	}
}

func (opt *Options) incDone(n int) {
	if opt.debug {
		if opt.counter == nil {
			opt.counter = &counter{}
		}
		opt.counter.IncDone(n)
	}
}

func (this *Worker) Debug() *Worker {
	this.opt.debug = true
	return this
}

func (this *Worker) WithDelaySec(sec int) *Worker {
	this.opt.DelaySec = sec
	return this
}

func (this *Worker) WithDelayCnt(cnt int) *Worker {
	this.opt.DelayCnt = cnt
	return this
}

func (this *Worker) WithQueueSize(size int) *Worker {
	this.opt.QueueSize = size
	return this
}

func (this *Worker) WithParallel(enable bool) *Worker {
	this.opt.Parallel = enable
	return this
}

func (this *Worker) WithLogger(logger logs.Logger) *Worker {
	this.opt.Logger = logger
	return this
}
