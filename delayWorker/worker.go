package delayWorker

import (
	"github.com/3th1nk/easygo/util/strUtil"
	"github.com/panjf2000/ants/v2"
	"runtime"
	"sync/atomic"
	"time"
)

const (
	StateCreated int32 = iota
	StateRunning
	StateStopping
	StateStopped
)

func StateString(state int32) string {
	switch state {
	case StateCreated:
		return "created"
	case StateRunning:
		return "running"
	case StateStopping:
		return "stopping"
	case StateStopped:
		return "stopped"
	default:
		return "unknown"
	}
}

type Worker struct {
	name    string                   // 名称
	receive chan interface{}         // 接收通道
	buf     []interface{}            // 缓冲区
	handler func(jobs []interface{}) // 处理函数
	stop    chan bool                // 停止信号
	inited  bool                     // 是否已经初始化
	state   int32                    // 状态
	pool    *ants.Pool               // 执行协程池
	opt     *Options                 // 可选配置
}

func New(name string, handler func(jobs []interface{}), opts ...*Options) *Worker {
	if handler == nil {
		panic("handler is nil")
	}
	if name == "" {
		name = strUtil.Rand(6)
	}

	var opt *Options
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	} else {
		opt = DefaultOptions()
	}

	return &Worker{
		name:    name,
		handler: handler,
		stop:    make(chan bool, 1),
		inited:  false,
		state:   StateCreated,
		opt:     opt,
	}
}

func (this *Worker) init() error {
	if this.inited {
		return nil
	}
	this.inited = true

	this.opt.ensure()
	if this.opt.Parallel {
		if runtime.NumCPU() > 1 {
			var err error
			this.pool, err = ants.NewPool(runtime.NumCPU())
			if err != nil {
				this.inited = false
				return err
			}
		} else {
			this.opt.Parallel = false
			this.opt.Logger.Warn("delayWorker[%s] parallel is disabled because of cpu core is 1", this.name)
		}
	}

	this.receive = make(chan interface{}, this.opt.QueueSize)
	this.buf = make([]interface{}, 0, this.opt.DelaySize)
	return nil
}

func (this *Worker) release() {
	if this.inited {
		close(this.receive)
		if this.pool != nil {
			this.pool.Release()
		}
		this.inited = false
	}
}

func (this *Worker) State() int32 {
	return atomic.LoadInt32(&this.state)
}

func (this *Worker) Stop() {
	if atomic.CompareAndSwapInt32(&this.state, StateRunning, StateStopping) {
		this.opt.Logger.Debug("delayWorker[%s] is stopping", this.name)

		// wait for all jobs to be processed
		this.stop <- true
		for atomic.LoadInt32(&this.state) == StateStopping {
			this.opt.Logger.Debug("delayWorker[%s] is waiting stopped, current state is:%s", this.name, StateString(atomic.LoadInt32(&this.state)))
			time.Sleep(time.Millisecond * 100)
		}
	} else if atomic.CompareAndSwapInt32(&this.state, StateCreated, StateStopped) {
		// do nothing
	}

	if s := this.opt.statistics(); s != "" {
		this.opt.Logger.Debug("delayWorker[%s] is stopped, %s", this.name, s)
	} else {
		this.opt.Logger.Debug("delayWorker[%s] is stopped", this.name)
	}
}

func (this *Worker) Push(job interface{}) {
	if atomic.LoadInt32(&this.state) == StateRunning {
		this.receive <- job
		this.opt.incReceive(1)
	} else {
		this.opt.Logger.Warn("delayWorker[%s] is not running, job is dropped", this.name)
		this.opt.incDrop(1)
	}
}

func (this *Worker) do() {
	if len(this.buf) == 0 {
		return
	}

	st := time.Now()
	defer func() {
		if r := recover(); r != nil {
			this.opt.Logger.Error("delayWorker[%s] do panic:%s", this.name, r)
		}
		this.opt.Logger.Debug("delayWorker[%s] do %d jobs, cost:%s", this.name, len(this.buf), time.Since(st))
		this.opt.incDone(len(this.buf))
		// !!!clear buf even if panic occurred
		this.buf = this.buf[:0]
	}()

	if this.opt.Parallel {
		copyBuf := make([]interface{}, len(this.buf))
		copy(copyBuf, this.buf)
		_ = this.pool.Submit(func() {
			this.handler(copyBuf)
		})
	} else {
		this.handler(this.buf)
	}
}

func (this *Worker) Run() {
	if !atomic.CompareAndSwapInt32(&this.state, StateCreated, StateRunning) &&
		!atomic.CompareAndSwapInt32(&this.state, StateStopped, StateRunning) {
		this.opt.Logger.Error("delayWorker[%s] is in state %s, cannot be started", this.name, StateString(atomic.LoadInt32(&this.state)))
		return
	}

	if err := this.init(); err != nil {
		this.opt.Logger.Error("delayWorker[%s] init error:%s", this.name, err)
		atomic.StoreInt32(&this.state, StateStopped)
		return
	}
	this.opt.Logger.Debug("delayWorker[%s] is running", this.name)

	go func() {
		var ticker = time.NewTicker(this.opt.DelayTime)
		defer func() {
			if r := recover(); r != nil {
				this.opt.Logger.Fatal("delayWorker[%s] panic:%s", this.name, r)
			}
			ticker.Stop()

			// processed all jobs in receive and buf before exit
			processing := true
			for processing {
				select {
				case job := <-this.receive:
					this.buf = append(this.buf, job)
					if len(this.buf) >= this.opt.DelaySize {
						this.do()
					}

				default:
					this.do()
					processing = false
					break
				}
			}

			// release resources
			this.release()

			atomic.StoreInt32(&this.state, StateStopped)
		}()

		for {
			select {
			case <-this.stop:
				this.opt.Logger.Debug("delayWorker[%s] receive stop signal", this.name)
				return

			case <-ticker.C:
				this.opt.Logger.Debug("delayWorker[%s] ticker", this.name)
				this.do()

			case job := <-this.receive:
				this.buf = append(this.buf, job)
				if len(this.buf) >= this.opt.DelaySize {
					this.opt.Logger.Debug("delayWorker[%s] buf is full", this.name)
					this.do()
				}
			}
		}
	}()
}
