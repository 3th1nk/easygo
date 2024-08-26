package influxdb

import (
	"fmt"
	"github.com/3th1nk/easygo/util/logs"
	"github.com/3th1nk/easygo/util/runtimeUtil"
	"github.com/panjf2000/ants/v2"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultFlushSize     = 1000
	defaultFlushInterval = 5 * time.Second

	stateNotRun   = 0
	stateRunning  = 1
	stateStopping = 2
)

type Client struct {
	addr            string // influxdb地址, 如：http://127.0.0.1:8086
	username        string // 用户名
	password        string // 密码 或 token
	mu              sync.RWMutex
	bucketGroups    map[string]*bucketGroup // key: db+rp  value: *bucketGroup
	groupSize       int                     // 每组桶的数量，默认16
	queryEpoch      string                  // 查询返回的时间格式, 默认s
	writePrecision  string                  // 写入数据的时间精度, 默认s
	writeSortTagKey bool                    // 写入数据时是否对tag key排序, 默认开启
	flushInterval   time.Duration           // 异步写入数据的时间间隔, 默认5s
	flushSize       int                     // 异步写入数据的行数，InfluxDB没有硬性限制单次写入行数，官方建议单次写入5000行
	debugger        *debugger               // 调试器, 默认不开启
	writePoolSize   int                     //
	writePool       *ants.Pool              // 写协程池
	state           int32                   // (异步协程)运行状态 0:未运行 1:运行中 2:正在停止
	stopSignal      chan struct{}           // 信号
	logger          logs.Logger
}

func NewClient(addr string, opts ...Option) *Client {
	c := &Client{
		addr:            strings.TrimSuffix(addr, "/"),
		bucketGroups:    make(map[string]*bucketGroup, 8),
		groupSize:       16,
		queryEpoch:      "s",
		writePrecision:  "s",
		writeSortTagKey: true,
		flushInterval:   defaultFlushInterval,
		flushSize:       defaultFlushSize,
		writePoolSize:   30,
		stopSignal:      make(chan struct{}, 1),
		logger:          logs.Default,
	}
	for _, opt := range opts {
		opt(c)
	}

	pool, err := ants.NewPool(c.writePoolSize, ants.WithLogger(c))
	if err != nil {
		panic(fmt.Sprintf("[InfluxDB] 创建写协程池失败: %v", err))
	}
	c.writePool = pool

	c.startAsyncWrite()
	return c
}

// Close 关闭客户端，释放资源
//	！！！务必在程序退出时调用，否则可能会导致数据丢失！！！
func (this *Client) Close() {
	this.stopAsyncWrite()
	close(this.stopSignal)
	this.writePool.Release()
}

func (this *Client) startAsyncWrite() {
	if atomic.LoadInt32(&this.state) == stateNotRun {
		go this.asyncWriter()
	}
}

func (this *Client) stopAsyncWrite() {
	atomic.StoreInt32(&this.state, stateStopping)
	this.stopSignal <- struct{}{}

	// 等待存量的写入任务完成
	for atomic.LoadInt32(&this.state) != stateNotRun {
		this.logger.Debug("[InfluxDB] 等待存量数据写入完成...")
		time.Sleep(time.Millisecond * 100)
	}
}

// Printf ants.Logger实现
func (this *Client) Printf(format string, args ...interface{}) {
	if logs.IsErrorEnable(this.logger) {
		this.logger.Error(format, args...)
	}
}

// Flush 强制刷新写入，立即将所有缓存的数据写入InfluxDB，一般无需手动调用
func (this *Client) Flush() {
	defer runtimeUtil.Recover()

	this.mu.RLock()
	if len(this.bucketGroups) == 0 {
		this.mu.RUnlock()
		return
	}

	groups := make([]*bucketGroup, 0, len(this.bucketGroups))
	for _, group := range this.bucketGroups {
		groups = append(groups, group)
	}
	this.mu.RUnlock()

	for _, group := range groups {
		url := this.buildWriteUrl(group.Db(), group.Rp())
		group.Range(func(bck *bucket) {
			if bck.Len() == 0 {
				return
			}

			if err := this.writePool.Submit(func() {
				lines := bck.Pop()
				if len(lines) == 0 {
					return
				}

				_ = this.doBatchWrite(url, lines)
			}); err != nil {
				this.logger.Warn("[InfluxDB] 获取写协程异常")
			}
		})
	}
}

// asyncWriter 异步写入
func (this *Client) asyncWriter() {
	atomic.StoreInt32(&this.state, stateRunning)
	defer func() {
		atomic.StoreInt32(&this.state, stateNotRun)
	}()

	ticker := time.NewTicker(this.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-this.stopSignal:
			this.logger.Debug("[InfluxDB] 存量数据开始写入")
			this.Flush()
			this.logger.Debug("[InfluxDB] 存量数据写入完成")
			return

		case <-ticker.C:
			this.logger.Debug("[InfluxDB] 定时批量写入数据")
			this.Flush()
		}
	}
}

// doBatchWrite 批量写入
func (this *Client) doBatchWrite(writeUrl string, lines []string) error {
	if len(lines) == 0 {
		return nil
	}

	for i := 0; i < len(lines); i += this.flushSize {
		end := i + this.flushSize
		if end > len(lines) {
			end = len(lines)
		}

		data := strings.Join(lines[i:end], "\n")
		if this.debugger != nil {
			this.logger.Debug("[InfluxDB] 写入数据 url=%v, line_count=%d, lines=\n%s", writeUrl, len(lines[i:end]), data)
		}
		resBody, err := doRequest(http.MethodPost, writeUrl, "", data, nil)
		if err != nil {
			this.countWrite(false)
			this.logger.Error("[InfluxDB] url=%v, err=%v, resp=%v", writeUrl, err.Error(), string(resBody))
			this.logger.Error(err.Error())
			return err
		}
		this.countWrite(true)
	}
	return nil
}
