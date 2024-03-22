package rabbitMQ

import (
	"crypto/tls"
	"fmt"
	"github.com/3th1nk/easygo/dataStruct/cyclelist"
	"github.com/3th1nk/easygo/util/arrUtil"
	"github.com/3th1nk/easygo/util/convertor"
	"github.com/3th1nk/easygo/util/pool"
	"github.com/3th1nk/easygo/util/runtimeUtil"
	"github.com/3th1nk/easygo/util/strUtil"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/singleflight"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"
)

func NewChannelPool(url string, opt ...ChannelPoolOptions) (*ChannelPool, error) {
	obj := &ChannelPool{url: strUtil.Split(url, ",", true, strings.TrimSpace), sg: &singleflight.Group{}}
	if len(opt) == 0 {
		obj.opt = DefaultPoolOptions
	} else {
		obj.opt = opt[0]
	}
	if obj.opt.Capacity <= 0 {
		obj.opt.Capacity = DefaultPoolOptions.Capacity
	}
	if obj.opt.TTL <= 0 {
		obj.opt.TTL = DefaultPoolOptions.TTL
	}
	obj.pool = pool.New(obj.opt.Capacity, obj.opt.TTL, func() interface{} {
		return &chanWrap{pool: obj}
	}, func(x interface{}) {
		x.(*chanWrap).close()
	})
	if obj.opt.LastJobSize > 0 {
		obj.lastJobList = cyclelist.New(obj.opt.LastJobSize)
	}

	// 测试连接（必须在 pool 初始化之后）
	if err := obj.Ping(); err != nil {
		return nil, err
	}

	poolLock.Lock()
	poolList = append(poolList, obj)
	poolLock.Unlock()

	return obj, nil
}

func AllChannelPool() []*ChannelPool {
	return poolList
}

var (
	poolList []*ChannelPool
	poolLock sync.Mutex
)

type ChannelPool struct {
	opt         ChannelPoolOptions                                     //
	url         []string                                               // 连接字符串列表
	pool        *pool.Pool                                             //
	lastJobList *cyclelist.CycleList                                   //
	conn        *amqp.Connection                                       // 当前活动连接
	connUrl     string                                                 // 当前活动连接所使用的连接字符串
	connLock    sync.Mutex                                             //
	queue       []*Queue                                               //
	sg          *singleflight.Group                                    // 用于控制同一时间点只有一个协程发起重连
	onGet       []func(job string, err error)                          //
	onPut       []func(job string, time time.Time, took time.Duration) //
}

type ChannelPoolOptions struct {
	TLSConfig   *tls.Config            // 当使用 ExternalAuth 时的 TLS 配置。为空表示不使用 ExternalAuth
	Capacity    int                    // amqp.Channel 缓存池容量上限。
	TTL         time.Duration          // amqp.Channel 在归还给缓存池之后多久被回收。默认 5 秒。
	Properties  map[string]interface{} //
	LastJobSize int                    //

	OnConnect    func(host string, err error)
	OnDisconnect func(host string, err error)
}

var DefaultPoolOptions = ChannelPoolOptions{
	Capacity: 2000,
	TTL:      5 * time.Second,
}

type JobInfo struct {
	Job  string        `json:"job,omitempty"`
	Time time.Time     `json:"time,omitempty"`
	Took time.Duration `json:"took,omitempty"`
}

func (this *ChannelPool) OnGet(f func(job string, err error)) {
	this.onGet = append(this.onGet, f)
}

func (this *ChannelPool) OnPut(f func(job string, time time.Time, took time.Duration)) {
	this.onPut = append(this.onPut, f)
}

// 获取当前正在使用的连接字符串
func (this *ChannelPool) ConnectedUrl() string {
	return this.connUrl
}

// LocalAddr 获取本地连接信息
func (this *ChannelPool) LocalAddr() net.Addr {
	if this.conn != nil {
		return this.conn.LocalAddr()
	}
	return nil
}

func (this *ChannelPool) Capacity() int {
	return this.opt.Capacity
}

func (this *ChannelPool) SetCapacity(n int) {
	this.opt.Capacity = n
	this.pool.SetCapacity(n)
}

func (this *ChannelPool) TTL() time.Duration {
	return this.opt.TTL
}

func (this *ChannelPool) SetTTL(d time.Duration) {
	this.opt.TTL = d
	this.pool.SetTTL(d)
}

// 已经从对缓存池中申请出去的 amqp.Channel 个数
func (this *ChannelPool) Size() int {
	return this.pool.Size()
}

// 可复用的空闲 amqp.Channel 个数
func (this *ChannelPool) Idle() int {
	return this.pool.Idle()
}

// 遍历当前活动对象
func (this *ChannelPool) LastJob() []*JobInfo {
	if this.lastJobList != nil {
		arr, cnt := make([]*JobInfo, this.lastJobList.Size()), 0
		this.lastJobList.ReverseWalk(func(a interface{}) {
			arr[cnt], cnt = a.(*JobInfo), cnt+1
		})
		if cnt != len(arr) {
			return arr[:cnt]
		}
		return arr
	}
	return nil
}

// 关闭连接
func (this *ChannelPool) Close(d time.Duration) (closed bool) {
	queueList := this.queue
	closed = runtimeUtil.GoWait(d, len(queueList), func(i int, d time.Duration) (done bool) {
		return queueList[i].Stop(d)
	})
	this.queue = nil

	this.pool.Close()

	if this.conn != nil {
		_ = this.conn.Close()
		this.conn, this.connUrl = nil, ""
	}

	return
}

func (this *ChannelPool) Ping() error {
	return this.do("ping", 0, 0, func(ch *amqp.Channel) error {
		return nil
	}, false)
}

func (this *ChannelPool) Do(job string, f func(ch *amqp.Channel)) error {
	ch, err := this.getChannel(job)
	if err != nil {
		return err
	}
	defer this.putChannel(ch)
	f(ch.channel)
	return nil
}

// 获取一个 Channel，并执行自定义函数。
//   autoClose: 函数执行完毕之后是否自动关闭 Channel，默认为 false 以复用 Channel。
//     但： Consume 函数总是应当自动关闭，因为 Consume 函数中可能对 Channel 做了其他参数初始化
func (this *ChannelPool) do(jobId string, retry int, interval time.Duration, f func(ch *amqp.Channel) error, autoClose bool) (jobError error) {
	ch, err := this.getChannel(jobId)
	if err != nil {
		return fmt.Errorf("get channel error: %v", err)
	}

	defer func() {
		if autoClose || jobError != nil {
			ch.close()
		}
		this.putChannel(ch)
	}()

	enableReconnect, retried, hasPanic := true, 0, false
	for retried <= retry {
		// 如果当前是在重试，则通过 open 方法获取新的 Channel，在新的 Channel 上重试
		if retried != 0 {
			time.Sleep(interval)

			if err = ch.close().open(); err != nil {
				return fmt.Errorf("reopen channel error: %v", err)
			}
		}

	reconnect:
		func() {
			defer runtimeUtil.Recover().OnPanic(func(p interface{}) {
				hasPanic, err = true, convertor.ToError(p)
			}).Handle()
			err = f(ch.channel)
		}()
		if err == nil {
			return nil
		} else if hasPanic {
			return fmt.Errorf("panic: %v", err)
		} else if err == amqp.ErrClosed || err == amqp.ErrChannelMax {
			// 连接不可用，关闭 Channel
			ch.close()
			if enableReconnect && nil == ch.open() {
				enableReconnect = false
				goto reconnect
			}
		}
		retried++
	}
	return err
}

func (this *ChannelPool) getChannel(job string) (ch *chanWrap, err error) {
	if len(this.onGet) != 0 {
		func() {
			defer runtimeUtil.Recover().Handle()
			for _, f := range this.onGet {
				f(job, err)
			}
		}()
	}

	ch, _ = this.pool.Get(-1).(*chanWrap)
	if ch == nil {
		return nil, &ChannelPoolBusyError{size: this.pool.Size()}
	}

	if err = ch.open(); err != nil {
		this.pool.Put(ch)
		return nil, &ChannelOpenError{cause: err}
	}

	ch.job = job
	if this.lastJobList != nil && job != "" {
		ch.time = time.Now()
	}

	return
}

func (this *ChannelPool) putChannel(ch *chanWrap) {
	if this.lastJobList != nil && ch.job != "" {
		ch.took = time.Since(ch.time)
		this.lastJobList.Add(&JobInfo{Job: ch.job, Time: ch.time, Took: ch.took})
	}

	if len(this.onPut) != 0 {
		for _, f := range this.onPut {
			f(ch.job, ch.time, ch.took)
		}
	}

	this.pool.Put(ch)
}

func (this *ChannelPool) reconnect() error {
	if this.conn != nil {
		shadowConn, shadowUrl := this.conn, this.connUrl
		this.conn, this.connUrl = nil, ""
		time.AfterFunc(10*time.Second, func() {
			err := shadowConn.Close()
			if this.opt.OnDisconnect != nil {
				_, _, _, host, _ := ParseUrl(shadowUrl)
				this.opt.OnDisconnect(host, err)
			}
		})
	}

	urlList, dialError := this.url, error(nil)
	for len(urlList) != 0 {
		idx := rand.Intn(len(urlList))
		url := AutoEncodeUrl(urlList[idx])
		cfg := amqp.Config{}
		if this.opt.TLSConfig != nil {
			cfg.Heartbeat = 10 * time.Second
			cfg.TLSClientConfig = this.opt.TLSConfig
			cfg.SASL = []amqp.Authentication{&amqp.ExternalAuth{}}
		}
		if len(this.opt.Properties) != 0 {
			cfg.Properties = this.opt.Properties
		}
		conn, err := amqp.DialConfig(url, cfg)
		if this.opt.OnConnect != nil {
			_, _, _, host, _ := ParseUrl(url)
			this.opt.OnConnect(host, err)
		}
		if err == nil {
			this.conn, this.connUrl = conn, url
			dialError = nil
			break
		} else {
			if conn != nil {
				_ = conn.Close()
			}
			urlList = arrUtil.RemoveStringAt(urlList, idx)
			dialError = err
		}
	}
	return dialError
}

type chanWrap struct {
	pool    *ChannelPool
	conn    *amqp.Connection
	channel *amqp.Channel
	job     string
	time    time.Time
	took    time.Duration
}

func (this *chanWrap) close() *chanWrap {
	if this.channel != nil {
		// TODO amqp channel close 不清楚什么原因会卡住，上层使用互斥锁导致死锁。这里加个超时控制
		channel := this.channel
		runtimeUtil.GoWait(time.Second*3, 1, func(_ int, _ time.Duration) (done bool) {
			_ = channel.Close()
			return true
		})
		this.channel = nil
	}
	this.conn = nil
	return this
}

func (this *chanWrap) open() (err error) {
	this.pool.connLock.Lock()
	defer this.pool.connLock.Unlock()

	if this.conn == nil || this.conn.IsClosed() {
		this.close()
	} else if this.channel != nil && !this.channel.IsClosed() {
		return nil
	}

	if this.pool.conn == nil || this.pool.conn.IsClosed() {
		// 如果当前 channel 不为 nil、但 conn 已经关闭，则释放当前 channel
		this.close()
		// 确保创建连接
		if err = this.pool.reconnect(); err != nil {
			return
		}
	}

	// 尝试获取 Channel
	this.channel, err = this.pool.conn.Channel()
	if err == amqp.ErrClosed {
		// 如果连接已关闭或 Channel 数量超限，就重连后再次尝试获取 Channel
		if err = this.pool.reconnect(); err == nil {
			this.channel, err = this.pool.conn.Channel()
		}
	}
	if err == nil {
		this.conn = this.pool.conn
	}

	return err
}
