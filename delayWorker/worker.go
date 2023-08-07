package delayWorker

import (
	"github.com/3th1nk/easygo/util/strUtil"
	"github.com/panjf2000/ants/v2"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type Worker struct {
	name    string
	receive chan interface{}
	buf     []interface{}
	handler func(jobs []interface{})
	running int32
	mu      sync.Mutex
	pool    *ants.Pool
	opt     *Options
}

func New(name string, handler func(jobs []interface{})) *Worker {
	if handler == nil {
		panic("handler is nil")
	}
	if name == "" {
		name = strUtil.Rand(6)
	}

	return &Worker{
		name:    name,
		handler: handler,
		opt:     &Options{},
	}
}

func (this *Worker) init() error {
	this.opt.ensure()

	if this.opt.Parallel {
		if runtime.NumCPU() > 1 {
			var err error
			this.pool, err = ants.NewPool(runtime.NumCPU())
			if err != nil {
				return err
			}
		} else {
			this.opt.Parallel = false
			this.opt.Logger.Warn("delayWorker[%s] parallel is disabled because of cpu core is 1", this.name)
		}
	}

	this.receive = make(chan interface{}, this.opt.QueueSize)
	this.buf = make([]interface{}, 0, this.opt.DelayCnt)
	return nil
}

func (this *Worker) Close() {
	if this.Running() {
		this.setRunning(0)
		this.opt.Logger.Debug("delayWorker[%s] close", this.name)

		this.do()
		for job := range this.receive {
			this.buf = append(this.buf, job)
			if len(this.buf) >= this.opt.DelayCnt {
				this.do()
			}
		}
		this.do()

		close(this.receive)
		if this.pool != nil {
			this.pool.Release()
		}
	}
}

func (this *Worker) Running() bool {
	return atomic.LoadInt32(&this.running) > 0
}

func (this *Worker) setRunning(val int) {
	atomic.AddInt32(&this.running, int32(val))
}

func (this *Worker) Push(job interface{}) {
	if this.Running() {
		this.receive <- job
	} else {
		this.opt.Logger.Warn("delayWorker[%s] is not running, job is dropped", this.name)
	}
}

func (this *Worker) do() {
	if len(this.buf) == 0 {
		return
	}
	this.mu.Lock()
	if len(this.buf) == 0 {
		this.mu.Unlock()
		return
	}
	this.opt.Logger.Debug("delayWorker[%s] do %d jobs", this.name, len(this.buf))

	if this.opt.Parallel {
		copyBuf := make([]interface{}, len(this.buf))
		copy(copyBuf, this.buf)
		this.buf = this.buf[:0]
		this.mu.Unlock()
		_ = this.pool.Submit(func() {
			this.handler(copyBuf)
		})
	} else {
		this.handler(this.buf)
		this.buf = this.buf[:0]
		this.mu.Unlock()
	}
}

func (this *Worker) Run() error {
	if this.Running() {
		return nil
	}
	if err := this.init(); err != nil {
		return err
	}
	this.setRunning(1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				this.opt.Logger.Fatal("delayWorker[%s] panic:%s", this.name, r)
			}
			this.setRunning(0)
		}()

		var timer *time.Timer
		for this.Running() {
			job, ok := <-this.receive
			if !ok {
				break
			}

			if len(this.buf) == 0 {
				timer = time.AfterFunc(time.Second*time.Duration(this.opt.DelaySec), func() {
					this.opt.Logger.Debug("delayWorker[%s] timer", this.name)
					this.do()
				})
			}

			this.mu.Lock()
			this.buf = append(this.buf, job)
			this.mu.Unlock()

			if len(this.buf) >= this.opt.DelayCnt {
				timer.Stop()
				this.opt.Logger.Debug("delayWorker[%s] buf is full", this.name)
				this.do()
			}
		}
	}()

	return nil
}
