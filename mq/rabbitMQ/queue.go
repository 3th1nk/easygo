/*
1、在消费函数处理结束后再发送Ack；
2、支持设置消费协程数量，并发消费；
3、支持停止消费协程、并等待所有已经获取到的消息处理完毕，以便实现进程优雅退出*/
package rabbitMQ

import (
	"context"
	"fmt"
	"github.com/3th1nk/easygo/counter/jobCounter"
	"github.com/3th1nk/easygo/util"
	"github.com/3th1nk/easygo/util/arrUtil"
	"github.com/3th1nk/easygo/util/jsonUtil"
	"github.com/3th1nk/easygo/util/logs"
	"github.com/3th1nk/easygo/util/osUtil"
	"github.com/3th1nk/easygo/util/runtimeUtil"
	"github.com/3th1nk/easygo/util/strUtil"
	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"github.com/panjf2000/ants/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

func New(name string, pool *ChannelPool, opt ...Options) *Queue {
	obj := &Queue{name: name, channelPool: pool}
	pool.queue = append(pool.queue, obj)
	if len(opt) == 0 {
		obj.opt = DefaultOptions
	} else {
		obj.opt = opt[0]
	}
	if !obj.opt.NoAutoDeclare {
		if obj.opt.QueueDeclareOptions == nil {
			obj.opt.QueueDeclareOptions = &DefaultQueueDeclareOptions
		}
		updateQueueDeclareOptions(obj.opt.QueueDeclareOptions, DefaultOptions.QueueDeclareOptions, &DefaultQueueDeclareOptions)
	}
	if obj.name == "" {
		obj.name = strUtil.Rand(6)
	}
	if !obj.opt.NoJobCounter {
		obj.jobCounter = jobCounter.New("rabbitMQ - " + name)
	}

	if obj.opt.ProducerSize > 0 {
		var err error
		if obj.opt.ProducerSize == 1 {
			obj.producer, err = obj.NewProducer()
		} else {
			obj.producer, err = obj.NewGroupProducer(obj.opt.ProducerSize)
		}
		if err != nil {
			runtimeUtil.PanicIfError(err)
		}
	}

	allQueue = append(allQueue, obj)

	osUtil.OnSignalExit(func(_ syscall.Signal) (stopped bool) {
		return obj.Stop(15 * time.Second)
	}, fmt.Sprintf("stop rabbitMQ %v", obj.name))

	if logger := obj.GetLogger(); logs.IsDebugEnable(logger) {
		logger.Debug("[mq.%s] new rabbitMQ: %v", obj.name, jsonUtil.MustMarshalToString(obj.opt))
	}

	for _, f := range onNewQueue {
		f(obj)
	}

	return obj
}

func OnNewQueue(f func(*Queue)) {
	onNewQueue = append(onNewQueue, f)
}

func AllQueue() []*Queue {
	return allQueue
}

var (
	DefaultOptions = Options{
		QueueDeclareOptions: &DefaultQueueDeclareOptions,
		ConsumeConcurrent:   1000,
		ProducerSize:        10,
	}

	DefaultQueueDeclareOptions = QueueDeclareOptions{
		Capacity:  100000,      //
		ExpireSec: 86400,       // 3600*24
		Overflow:  "drop-head", //
	}

	allQueue    = make([]*Queue, 0, 2) //
	declaredMap map[string]bool        //
	declareLock sync.RWMutex           //

	onNewQueue []func(*Queue)
)

type Options struct {
	NoAutoDeclare       bool                 `json:"no_auto_declare,omitempty"`       // 指示不自动声明队列。
	QueueDeclareOptions *QueueDeclareOptions `json:"queue_declare_options,omitempty"` // 声明队列时的默认参数。
	ConsumeConcurrent   int                  `json:"consume_concurrent,omitempty"`    // 消费者默认并发数
	NoJobCounter        bool                 `json:"no_job_counter,omitempty"`        // 是否不使用 JobCounter。
	NoProduceLog        bool                 `json:"no_produce_log,omitempty"`        // 当发送消息时是否记录日志
	NoConsumeLog        bool                 `json:"no_consume_log,omitempty"`        // 当接收到消息时是否记录日志
	NoWaitOnStop        bool                 `json:"no_wait_on_stop,omitempty"`       // 调用 Stop 方法时，是否等待当前正则执行的任务处理结束。
	Logger              logs.Logger          `json:"-"`                               //

	ProducerSize int `json:"producer_size,omitempty"` // 生产者数量
}

type QueueDeclareOptions struct {
	Capacity  int    `json:"capacity,omitempty"`   // 100000
	ExpireSec int    `json:"expire_sec,omitempty"` // x-expires，秒
	Overflow  string `json:"overflow,omitempty"`   // x-overflow
}

func updateQueueDeclareOptions(opt *QueueDeclareOptions, source ...*QueueDeclareOptions) {
	if opt == nil {
		return
	}
	source, _ = arrUtil.Find(source, func(i int) bool { return source[i] != nil }).([]*QueueDeclareOptions)

	if opt.Capacity <= 0 {
		for _, v := range source {
			if v.Capacity > 0 {
				opt.Capacity = v.Capacity
				break
			}
		}
	}
	if opt.ExpireSec <= 0 {
		for _, v := range source {
			if v.ExpireSec > 0 {
				opt.ExpireSec = v.ExpireSec
				break
			}
		}
	}
	if opt.Overflow == "" {
		for _, v := range source {
			if v.Overflow != "" {
				opt.Overflow = v.Overflow
				break
			}
		}
	}
}

// ------------------------------------------------------------------------------

type Queue struct {
	name        string              //
	opt         Options             //
	channelPool *ChannelPool        //
	consumer    int32               // 当前正在运行的 consumer 数量
	jobCounter  *jobCounter.Counter //

	producer IProducer

	stop      []func(wait time.Duration) bool
	onProduce []func(e *ProduceEventArgs)
	onConsume []func(e *ConsumeEventArgs)
}

func (this *Queue) Name() string {
	return this.name
}

func (this *Queue) QueueDeclareOptions() *QueueDeclareOptions {
	return this.opt.QueueDeclareOptions
}

func (this *Queue) GetLogger() logs.Logger {
	if !reflect2.IsNil(this.opt.Logger) {
		return this.opt.Logger
	}
	return logs.Default
}

func (this *Queue) SetLogger(logger logs.Logger) {
	this.opt.Logger = logger
}

// ------------------------------------------------------------------------------ celeryImpl methods

type ProduceEventArgs struct {
	Time     time.Time
	Exchange string
	Queue    string
	Data     string
	Error    error
	Took     time.Duration
}

// 设置 Produce 事件回调函数
func (this *Queue) OnProduce(f func(e *ProduceEventArgs)) {
	this.onProduce = append(this.onProduce, f)
}

type ConsumeEventArgs struct {
	Time  time.Time
	Queue string
	Data  string
	Error error
	Took  time.Duration
}

// 设置 Produce 事件回调函数
func (this *Queue) OnConsume(f func(e *ConsumeEventArgs)) {
	this.onConsume = append(this.onConsume, f)
}

func (this *Queue) fireConsume(e *ConsumeEventArgs) {
	for _, f := range this.onConsume {
		func() {
			defer runtimeUtil.HandleRecover("")
			f(e)
		}()
	}
}

// 获取当前正在使用的参数
func (this *Queue) GetOptions() Options {
	return this.opt
}

func (this *Queue) GetPool() *ChannelPool {
	return this.channelPool
}

// 声明一个队列
//   queue: 队列名称
//   opt: 队列声明参数。默认使用 GetOptions().QueueDeclareOptions
func (this *Queue) DeclareQueue(queue string, opt ...QueueDeclareOptions) (err error) {
	var theOpt *QueueDeclareOptions
	if len(opt) == 0 {
		if theOpt = this.opt.QueueDeclareOptions; theOpt == nil {
			if theOpt = DefaultOptions.QueueDeclareOptions; theOpt == nil {
				return fmt.Errorf("missing 'QueueDeclareOptions'")
			}
		}
	} else {
		theOpt = &opt[0]
	}
	updateQueueDeclareOptions(theOpt, this.opt.QueueDeclareOptions, DefaultOptions.QueueDeclareOptions, &DefaultQueueDeclareOptions)
	optStr := jsonUtil.MustMarshalToString(theOpt)

	logger := this.GetLogger()
	err = this.channelPool.do(fmt.Sprintf("declareQueue(%v, %v)", queue, optStr), 0, 0, func(ch *amqp.Channel) error {
		return this.doDeclareQueue(ch, queue, theOpt)
	}, false)
	if err != nil {
		logger.Error("[mq.%s] declare queue error: %v, queue=%v, opt=%v", this.name, err, queue, optStr)
	} else if logs.IsDebugEnable(logger) {
		logger.Debug("[mq.%s] declare queue: %v, opt=%v", this.name, queue, optStr)
	}
	return err
}

func (this *Queue) doDeclareQueue(ch *amqp.Channel, queue string, opt *QueueDeclareOptions) (err error) {
	_, err = ch.QueueDeclare(queue, true, false, false, false, map[string]interface{}{
		"x-expires":    opt.ExpireSec * 1000,
		"x-max-length": opt.Capacity,
		"x-overflow":   opt.Overflow,
	})
	if err != nil {
		if str := err.Error(); strings.Contains(str, "inequivalent arg") && strings.Contains(str, "but current is") {
			_ = ch.Close()
			if logger := this.GetLogger(); logs.IsWarnEnable(logger) {
				logger.Warn("[mq.%s] declare queue args error: %v", this.name, err)
			}
			err = nil
		} else {
			return err
		}
	}
	this.setDeclared("", queue)
	return nil
}

// 声明一个交换机
//   exchange: 交换机
//   kind: 交换机类型，默认 fanout
func (this *Queue) DeclareExchange(exchange string, kind ...string) (err error) {
	theKind := append(kind, "fanout")[0]
	logger := this.GetLogger()
	err = this.channelPool.do(fmt.Sprintf("declareExchange(%v, %v)", exchange, theKind), 1, 0, func(ch *amqp.Channel) error {
		return this.doDeclareExchange(ch, exchange, theKind)
	}, false)
	if err != nil {
		logger.Error("[mq.%s] declare exchange error: %v, exchange=%v, kind=%v", this.name, err, exchange, theKind)
	} else if logs.IsDebugEnable(logger) {
		logger.Debug("[mq.%s] declare exchange: %v, kind=%v", this.name, exchange, theKind)
	}
	return err
}

func (this *Queue) doDeclareExchange(ch *amqp.Channel, exchange, kind string) (err error) {
	err = ch.ExchangeDeclare(exchange, kind, true, false, false, false, map[string]interface{}{})
	if err != nil {
		return err
	}
	this.setDeclared(exchange, "")
	return nil
}

// 将队列绑定到交换机上
func (this *Queue) BindQueue(queue, exchange string) (err error) {
	logger := this.GetLogger()
	err = this.channelPool.do(fmt.Sprintf("bindQueue(%v, %v)", queue, exchange), 1, 0, func(ch *amqp.Channel) error {
		return ch.QueueBind(queue, "", exchange, false, nil)
	}, false)
	if err != nil {
		logger.Error("[mq.%s] bind queue error: %v, exchange=%v, queue=%v", this.name, err, exchange, queue)
	} else if logs.IsDebugEnable(logger) {
		logger.Debug("[mq.%s] bind queue: exchange=%v, queue=%v", this.name, exchange, queue)
	}
	return err
}

// 取消队列绑定
func (this *Queue) UnbindQueue(queue, exchange string) (err error) {
	logger := this.GetLogger()
	err = this.channelPool.do(fmt.Sprintf("unbindQueue(%v, %v)", queue, exchange), 1, 0, func(ch *amqp.Channel) error {
		return ch.QueueUnbind(queue, "", exchange, nil)
	}, false)
	if err != nil {
		logger.Error("[mq.%s] unbind queue error: %v, exchange=%v, queue=%v", this.name, err, exchange, queue)
	} else if logs.IsDebugEnable(logger) {
		logger.Debug("[mq.%s] unbind queue: exchange=%v, queue=%v", this.name, exchange, queue)
	}
	return err
}

// 从服务端移除一个队列
func (this *Queue) RemoveQueue(queue string) (err error) {
	logger := this.GetLogger()
	err = this.channelPool.do(fmt.Sprintf("removeQueue(%v)", queue), 1, 0, func(ch *amqp.Channel) error {
		_, err = ch.QueueDelete(queue, false, false, false)
		return err
	}, false)
	if err != nil {
		logger.Error("[mq.%s] remove queue error: %v, queue=%v", this.name, err, queue)
	} else if logs.IsDebugEnable(logger) {
		_, key := this.getDeclaredMapKey("", queue)
		declareLock.Lock()
		delete(declaredMap, key)
		declareLock.Unlock()

		logger.Debug("[mq.%s] remove queue: %v", this.name, queue)
	}
	return err
}

// 从服务端移除一个交换机
func (this *Queue) RemoveExchange(exchange string) (err error) {
	logger := this.GetLogger()
	err = this.channelPool.do(fmt.Sprintf("removeExchange(%v)", exchange), 1, 0, func(ch *amqp.Channel) error {
		return ch.ExchangeDelete(exchange, false, false)
	}, false)
	if err != nil {
		logger.Error("[mq.%s] remove exchange error: %v, exchange=%v", this.name, err, exchange)
	} else if logs.IsDebugEnable(logger) {
		key, _ := this.getDeclaredMapKey(exchange, "")
		declareLock.Lock()
		delete(declaredMap, key)
		declareLock.Unlock()

		logger.Debug("[mq.%s] remove exchange: %v", this.name, exchange)
	}
	return err
}

func (this *Queue) checkAndDeclare(exchange, queue string) error {
	exchangeKey, queueKey := this.getDeclaredMapKey(exchange, queue)

	declareLock.RLock()
	declareExchange := exchangeKey != "" && !declaredMap[exchangeKey]
	declareQueue := queueKey != "" && !declaredMap[queueKey]
	declareLock.RUnlock()

	if !declareExchange && !declareQueue {
		return nil
	}

	logger := this.GetLogger()
	if declareExchange {
		if err := this.channelPool.do(fmt.Sprintf("declareExchange(%s)", exchange), 1, 0, func(ch *amqp.Channel) error {
			return this.doDeclareExchange(ch, exchange, "fanout")
		}, false); err != nil {
			logger.Error("[mq.%s] auto declare exchange error: %v, exchange=%v", this.name, err, exchange)
			return err
		} else if logs.IsDebugEnable(logger) {
			logger.Debug("[mq.%s] auto declare exchange: %v", this.name, exchange)
		}
	}
	if declareQueue {
		if err := this.channelPool.do(fmt.Sprintf("declareQueue(%s)", exchange), 1, 0, func(ch *amqp.Channel) error {
			return this.doDeclareQueue(ch, queue, this.opt.QueueDeclareOptions)
		}, false); err != nil {
			logger.Error("[mq.%s] auto declare queue error: %v, queue=%v, opt=%v", this.name, err, queue, jsonUtil.MustMarshalToString(&this.opt.QueueDeclareOptions))
			return err
		} else if logs.IsDebugEnable(logger) {
			logger.Debug("[mq.%s] auto declare queue: %v, opt=%v", this.name, queue, jsonUtil.MustMarshalToString(&this.opt.QueueDeclareOptions))
		}
	}

	return nil
}

func (this *Queue) getDeclaredMapKey(exchange, queue string) (exchangeKey, queueKey string) {
	if exchange != "" {
		exchangeKey = "exchange:" + exchange
	}
	if queue != "" {
		queueKey = "queue:" + queue
	}
	return
}

func (this *Queue) setDeclared(exchange, queue string) {
	exchangeKey, queueKey := this.getDeclaredMapKey(exchange, queue)

	declareLock.Lock()
	defer declareLock.Unlock()

	if declaredMap == nil {
		declaredMap = make(map[string]bool, 64)
	}
	if exchangeKey != "" {
		declaredMap[exchangeKey] = true
	}
	if queueKey != "" {
		declaredMap[queueKey] = true
	}
}

// 检测 MQ 是否可以访问。如果不可用则返回对应的错误信息。
func (this *Queue) Ping() error {
	return this.channelPool.Ping()
}

// 发送一条消息。
// 发送结束后将触发 OnProduce 事件。
func (this *Queue) Produce(exchange, queue, data string) (err error) {
	if !this.opt.NoAutoDeclare {
		if err = this.checkAndDeclare(exchange, queue); err != nil {
			return err
		}
	}

	// produce event
	if len(this.onProduce) != 0 {
		e := &ProduceEventArgs{Time: time.Now(), Exchange: exchange, Queue: queue, Data: data}
		defer func() {
			e.Error, e.Took = err, time.Since(e.Time)
			for _, f := range this.onProduce {
				func() {
					defer runtimeUtil.Recover().Handle()
					f(e)
				}()
			}
		}()
	}

	return this.send(exchange, queue, data)
}

func (this *Queue) GetDefaultProducer() IProducer {
	if this.producer != nil {
		return this.producer
	}

	p, err := this.NewProducer()
	runtimeUtil.PanicIfError(err)

	return p
}

func (this *Queue) send(exchange, queue string, data string) error {
	if this.producer != nil {
		return this.producer.Send(exchange, queue, data)
	} else {
		err := this.channelPool.do(fmt.Sprintf("procude(%v%v, %v)", exchange, queue, jsoniter.Get([]byte(data), "sid").GetInterface()), 0, 0, func(ch *amqp.Channel) error {
			return ch.PublishWithContext(context.Background(), exchange, queue, false, false, amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				Body:         []byte(data),
			})
		}, false)
		if err != nil {
			this.GetLogger().Error("[mq.%s] produce error: %v, exchange=%v, queue=%v, msg=%v", this.name, err, exchange, queue, data)
		} else if !this.opt.NoProduceLog && logs.IsDebugEnable(this.GetLogger()) {
			if exchange != "" {
				this.GetLogger().Debug(`[mq.%s] produce: exchange=%v, msg=%v`, this.name, exchange, data)
			} else {
				this.GetLogger().Debug(`[mq.%s] produce: queue=%v, msg=%v`, this.name, queue, data)
			}
		}
		return err
	}
}

type ConsumeOptions struct {
	// 交换机，队列绑定的交换机，如果队列不存在，重建队列+绑定交换机
	Exchange string
	// handler 超时时间
	Timeout time.Duration
	// With a prefetch count greater than zero, the server will deliver that many
	// messages to consumers before acknowledgments are received.
	Qos int
	// 当创建新的 Channel 时触发。
	OnChannel func(ch *amqp.Channel)
}

// 启动指定队列的消费者。
//   queue: 队列名
//   concurrent: 并发数（消费者协程数）
//   handler: 消息处理函数
// 当从队列中接收到消息并处理结束后会触发 Consume 事件。
func (this *Queue) Consume(queue string, concurrent int, handler func(data string) (ack bool), opt ...ConsumeOptions) (stop func(wait time.Duration) (stopped bool), err error) {
	var theOpt ConsumeOptions
	if len(opt) != 0 {
		theOpt = opt[0]
	}

	if !this.opt.NoAutoDeclare {
		err = this.DeclareAndBindExchangeQueue(theOpt.Exchange, queue)
		if err != nil {
			return
		}
	}

	atomic.AddInt32(&this.consumer, 1)
	if concurrent == 0 {
		concurrent = DefaultOptions.ConsumeConcurrent
	}

	activeJob, stopRequired, stopChan := int32(0), int32(0), make(chan bool, 2)
	stop = func(wait time.Duration) bool {
		if atomic.CompareAndSwapInt32(&stopRequired, 0, 1) {
			stopChan <- true
		}

		// 等待
		if !this.opt.NoWaitOnStop && activeJob > 0 && wait > 0 {
			expired := false
			time.AfterFunc(wait, func() {
				expired = true
				if activeJob > 0 {
					this.GetLogger().Error("[mq.%s] stop consume timeout: %v, activeJob=%v, wait=%v", this.name, queue, activeJob, wait)
				}
			})
			for !expired && activeJob > 0 {
				time.Sleep(10 * time.Millisecond)
			}
		}

		return this.opt.NoWaitOnStop || activeJob <= 0
	}
	this.stop = append(this.stop, stop)

	go func() {
		jobPool, _ := ants.NewPool(concurrent, ants.WithExpiryDuration(time.Minute)) // 协程超时时间

		// 退出时 running 自减。如果自减结果为 0 ，表明当前 consumer 已经停止，同时 this.consumer 自减。
		defer func() {
			runtimeUtil.HandleRecover("", recover())
			atomic.AddInt32(&this.consumer, -1)
			jobPool.Release()

			if logger := this.GetLogger(); logs.IsWarnEnable(logger) {
				logger.Warn("queue=%v exchange=%v consume exit stopSignal=%v", queue, theOpt.Exchange, stopRequired)
			}
		}()

		for stopRequired == 0 {
			consumeError := this.channelPool.do(fmt.Sprintf("consume(%s)", queue), 0, time.Second, func(ch *amqp.Channel) (err error) {
				if theOpt.OnChannel != nil {
					func() {
						defer runtimeUtil.Recover().Handle()
						theOpt.OnChannel(ch)
					}()
				}

				if theOpt.Qos != 0 {
					if err = ch.Qos(theOpt.Qos, 0, false); err != nil {
						if logger := this.GetLogger(); logs.IsErrorEnable(logger) {
							logger.Error("queue=%v set PrefetchCount=%v error: %v", queue, theOpt.Qos, err)
						}
						// 设置出错了，把这个参数置为 0，这样下次重试时就不会再次尝试设置
						theOpt.Qos = 0
						return err
					}
				}

				delivery, err := ch.Consume(queue, "", false, false, false, false, nil)
				if err != nil {
					return err
				}

				if logger := this.GetLogger(); logs.IsInfoEnable(logger) {
					logger.Info("[mq.%s] start consumer: queue=%v, concurrent=%v, qos=%v", this.name, queue, concurrent, theOpt.Qos)
				}

				for {
					select {
					case <-stopChan:
						return nil
					case msg, ok := <-delivery:
						if !ok {
							return fmt.Errorf("delivery failed")
						}
						msgBody := string(msg.Body)

						if !this.opt.NoConsumeLog {
							if logger := this.GetLogger(); logs.IsDebugEnable(logger) {
								logger.Debug("[mq.%s] consume: queue=%v, msg=%v", this.name, queue, msgBody)
							}
						}

						var jobId string
						if this.jobCounter != nil {
							jobId = fmt.Sprintf("handle(%v, %v) [%v]", queue, util.ShortStr(msgBody), strUtil.Rand(8))
							this.jobCounter.Add(jobId)
						}

						start := time.Now()
						_ = jobPool.Submit(func() {
							ack, handleError := true, error(nil)
							func() {
								defer runtimeUtil.Recover().Err(&handleError).Handle()
								ack = handler(msgBody)
							}()

							// 如果消费者panic了，重试的时候还是会panic，这里直接确认
							if handleError != nil {
								ack = true
								this.GetLogger().Error("[mq.%s] handle msg error: %v, queue=%v, msg=%v", this.name, handleError, queue, util.ShortStr(msgBody))
							}

							if ack {
								if err := msg.Ack(false); err != nil {
									this.GetLogger().Error("[mq.%s] ack msg error: %v, queue=%v, msg=%v", this.name, err, queue, util.ShortStr(msgBody))
								}
							} else {
								if err := msg.Reject(true); err != nil {
									this.GetLogger().Error("[mq.%s] reject msg error: %v, queue=%v, msg=%v", this.name, err, queue, util.ShortStr(msgBody))
								} else {
									this.GetLogger().Warn("[mq.%s] reject msg: queue=%v, msg=%v", this.name, queue, util.ShortStr(msgBody))
								}
							}

							if len(this.onConsume) > 0 {
								this.fireConsume(&ConsumeEventArgs{
									Time:  start,
									Queue: queue,
									Data:  msgBody,
									Error: handleError,
									Took:  time.Since(start),
								})
							}

							atomic.AddInt32(&activeJob, -1)
							if jobId != "" {
								this.jobCounter.Done(jobId)
							}
						})
					}
				}
			}, true)
			if consumeError != nil {
				if strings.Contains(consumeError.Error(), "NOT_FOUND - no queue") {
					consumeError = this.DeclareAndBindExchangeQueue(theOpt.Exchange, queue)
					this.GetLogger().Error("[mq.%s] consume error: auto-declare exchange=%v, queue=%v, error=%v", this.name, theOpt.Exchange, queue, consumeError)
				} else {
					this.GetLogger().Error("[mq.%s] consume error: queue=%v, error=%v", this.name, queue, consumeError)
				}

				// fire event
				if len(this.onConsume) != 0 {
					this.fireConsume(&ConsumeEventArgs{
						Time:  time.Now(),
						Queue: queue,
						Error: consumeError,
					})
				}

				// 队列消费出错表明队列目前不可用，需要等待一段较长的时间
				timeout := time.Now().Add(15 * time.Second)
				for stopRequired == 0 && time.Now().Before(timeout) {
					time.Sleep(200 * time.Millisecond)
				}
			}
		}
	}()

	return
}

func (this *Queue) onHandleFinish(msg *amqp.Delivery, ack bool, args *ConsumeEventArgs) {
	defer func() {
		runtimeUtil.HandleRecover("", recover())
		// fire event
		this.fireConsume(args)
	}()

	logger := this.GetLogger()
	shortData := util.ShortStr(args.Data)

	if args.Error != nil {
		logger.Error("[mq.%s] handle msg error: %v, queue=%v, msg=%v", this.name, args.Error, args.Queue, shortData)
	}

	if ack {
		if err := msg.Ack(false); err != nil {
			logger.Error("[mq.%s] ack msg error: %v, queue=%v, msg=%v", this.name, err, args.Queue, shortData)
		}
	} else {
		if err := msg.Reject(true); err != nil {
			logger.Error("[mq.%s] reject msg error: %v, queue=%v, msg=%v", this.name, err, args.Queue, shortData)
		} else if logs.IsWarnEnable(logger) {
			logger.Warn("[mq.%s] reject msg: queue=%v, msg=%v", this.name, args.Queue, shortData)
		}
	}
}

func (this *Queue) Do(f func(ch *amqp.Channel) error) error {
	var jobName string
	_, file, line, ok := runtime.Caller(1)
	if ok {
		_, file := filepath.Split(file)
		jobName = fmt.Sprintf("Do(%p), %v:%v", f, file, line)
	} else {
		jobName = fmt.Sprintf("Do(%p)", f)
	}
	return this.channelPool.do(jobName, 1, 0, func(ch *amqp.Channel) error {
		return f(ch)
	}, true) // 由于不确定回调函数中有没有更改 Channel 的参数，保险起见自动关闭。
}

// 停止所有的 Consumer（如果有启动的话）、并关闭连接。
func (this *Queue) Stop(wait time.Duration) (stopped bool) {
	stopped = runtimeUtil.GoWait(wait, len(this.stop), func(i int, d time.Duration) (done bool) {
		return this.stop[i](d)
	})

	// 等待一小段时间
	time.Sleep(20 * time.Millisecond)

	return stopped
}

func (this *Queue) DeclareAndBindExchangeQueue(exchange, queue string) error {
	if queue != "" {
		if err := this.DeclareQueue(queue); err != nil {
			return err
		}
	}

	if exchange != "" {
		if err := this.DeclareExchange(exchange); err != nil {
			return err
		}
		return this.BindQueue(queue, exchange)
	}

	return nil
}
