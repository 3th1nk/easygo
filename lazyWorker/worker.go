package lazyWorker

import (
	"github.com/3th1nk/easygo/util/strUtil"
	"github.com/panjf2000/ants/v2"
	"runtime"
	"sync/atomic"
	"time"
)

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
		state:   int32(StateCreated),
		opt:     opt,
	}
}

func (this *Worker) init() error {
	if this.inited {
		return nil
	}
	this.inited = true

	this.opt.ensure()
	if !this.opt.Serialization {
		if runtime.NumCPU() > 1 {
			var err error
			this.pool, err = ants.NewPool(runtime.NumCPU())
			if err != nil {
				this.inited = false
				return err
			}
		} else {
			this.opt.Serialization = true
			this.opt.Logger.Warn("lazyWorker[%s] concurrency is disabled because of cpu core is 1", this.name)
		}
	}

	this.receive = make(chan interface{}, this.opt.QueueSize)
	this.buf = make([]interface{}, 0, this.opt.LazySize)
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

func (this *Worker) getState() State {
	return State(atomic.LoadInt32(&this.state))
}

func (this *Worker) setState(state State) {
	atomic.StoreInt32(&this.state, int32(state))
}

func (this *Worker) cmpAndSwapState(old, new State) bool {
	return atomic.CompareAndSwapInt32(&this.state, int32(old), int32(new))
}

func (this *Worker) Stop() {
	if this.cmpAndSwapState(StateRunning, StateStopping) {
		this.opt.Logger.Debug("lazyWorker[%s] is stopping", this.name)

		// wait for all jobs to be processed
		this.stop <- true
		for this.getState() == StateStopping {
			this.opt.Logger.Debug("lazyWorker[%s] is waiting stopped, current state is:%s", this.name, this.getState())
			time.Sleep(time.Millisecond * 100)
		}
	} else if this.cmpAndSwapState(StateCreated, StateStopped) {
		// do nothing
	}

	if s := this.opt.statistics(); s != "" {
		this.opt.Logger.Debug("lazyWorker[%s] is stopped, %s", this.name, s)
	} else {
		this.opt.Logger.Debug("lazyWorker[%s] is stopped", this.name)
	}
}

func (this *Worker) Push(job interface{}) {
	if this.getState() == StateRunning {
		this.receive <- job
		this.opt.incReceive(1)
	} else {
		this.opt.Logger.Warn("lazyWorker[%s] is not running, job is dropped", this.name)
		this.opt.incDrop(1)
	}
}

func (this *Worker) do() {
	if len(this.buf) == 0 {
		return
	}

	//st := time.Now()
	defer func() {
		if r := recover(); r != nil {
			this.opt.Logger.Error("lazyWorker[%s] do panic:%s", this.name, r)
		}
		//this.opt.Logger.Debug("lazyWorker[%s] do %d jobs, cost:%s", this.name, len(this.buf), time.Since(st))
		this.opt.incDone(len(this.buf))
		// !!!clear buf even if panic occurred
		this.buf = this.buf[:0]
	}()

	if !this.opt.Serialization {
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
	if !this.cmpAndSwapState(StateCreated, StateRunning) &&
		!this.cmpAndSwapState(StateStopped, StateRunning) {
		this.opt.Logger.Error("lazyWorker[%s] is in state %s, cannot be started", this.name, this.getState())
		return
	}

	if err := this.init(); err != nil {
		this.opt.Logger.Error("lazyWorker[%s] init error:%s", this.name, err)
		this.setState(StateStopped)
		return
	}
	this.opt.Logger.Debug("lazyWorker[%s] is running", this.name)

	go func() {
		var ticker = time.NewTicker(this.opt.LazyInterval)
		defer func() {
			if r := recover(); r != nil {
				this.opt.Logger.Fatal("lazyWorker[%s] panic:%s", this.name, r)
			}
			ticker.Stop()

			// processed all jobs in receive and buf before exit
			processing := true
			for processing {
				select {
				case job := <-this.receive:
					this.buf = append(this.buf, job)
					if len(this.buf) >= this.opt.LazySize {
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

			this.setState(StateStopped)
		}()

		for {
			select {
			case <-this.stop:
				this.opt.Logger.Debug("lazyWorker[%s] receive stop signal", this.name)
				return

			case <-ticker.C:
				//this.opt.Logger.Debug("lazyWorker[%s] ticker", this.name)
				this.do()

			case job := <-this.receive:
				this.buf = append(this.buf, job)
				if len(this.buf) >= this.opt.LazySize {
					//this.opt.Logger.Debug("lazyWorker[%s] buf is full", this.name)
					this.do()
				}
			}
		}
	}()
}
