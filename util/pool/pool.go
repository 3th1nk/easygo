package pool

import (
	"github.com/3th1nk/easygo/util/mathUtil"
	"github.com/3th1nk/easygo/util/timeUtil"
	"sync"
	"time"
)

// 创建缓存池。
//
// 参数：
//    capacity:
//       最大创建的对象数。一旦超过该阈值，在 Get 时就会等待其他协程归还（Put）对象。
//    ttl:
//       对象在被放回 Pool 之后，过多久时间会被自动清除（从 Pool 中移除、并调用 del 方法释放对象）。
//       如果 TTL ≤ 0，则对象在调用 Pool.Put 方法时，会立即调用 del 方法进行释放、并且不会放回资源池，相当于
//       不使用缓存池。
//    new:
//       对象的创建方法。
//    del:
//       对象的释放方法。
//       如果 del 方法为 nil、且对象实现了 Close() 方法，则默认会使用 Close 方法释放；否则需在 del 内部自行
//       释放资源。
func New(capacity int, ttl time.Duration, new func() interface{}, del func(x interface{})) *Pool {
	if new == nil {
		new = func() interface{} { return nil }
	}
	pool := &Pool{new: new, del: del}
	pool.SetCapacity(capacity).SetTTL(ttl)
	return pool
}

// 资源对象的缓存池
//   用于代替 sync.Pool 以便在复用的同时又能够及时释放资源。
//
//   网络连接池就是其中的一种典型场景：
//     为了在高并发情况下充分复用网络连接，调用方使用完毕后不会立即关闭网络连接，而是暂时放在资源池中，以便下次使
//     用时直接取出复用。
//     但如果业务在某个时刻对资源使用非常高、但在接下来的相当长时间内使用情况又非常低，这时候高峰期缓存的资源对象
//     就会在内存中迟迟得不到利用，此时反而降低了内存使用率。
//
//   Pool 通过定时器对 “放回到资源池中超过指定时间的对象” 进行释放，这样既能在高峰期及时复用、又能在低峰时及时释
//   放资源。
type Pool struct {
	capacity int
	size     int
	ttl      time.Duration      // 对象归还给缓存池之后多久被回收
	new      func() interface{} // 当需要创建新的对象时触发
	del      func(interface{})  // 当对象（因 TTL 过期或其他情况）需要被销毁时触发
	idle     []*wrapper         // 已归还给缓存池但尚未清除的对象
	ticker   *timeUtil.Ticker   // 回收计时器
	mu       sync.RWMutex       //
}

type wrapper struct {
	obj    interface{}
	expire time.Time
}

func (this *Pool) SetCapacity(capacity int) *Pool {
	if capacity <= 0 {
		this.capacity = -1
	} else {
		this.capacity = capacity
	}
	return this
}

func (this *Pool) SetTTL(ttl time.Duration) *Pool {
	this.ttl = ttl
	if ttl <= 0 {
		if this.ticker != nil {
			this.ticker.Stop(0)
		}
	} else {
		if this.ticker == nil {
			this.ticker = timeUtil.NewTicker(ttl, ttl, this.doCheckExpire)
		} else {
			this.ticker.SetDuration(ttl)
		}
	}
	return this
}

func (this *Pool) Capacity() int {
	return this.capacity
}

// 已经从对缓存池中申请出去的对象个数
func (this *Pool) Size() int {
	return this.size
}

func (this *Pool) TTL() time.Duration {
	return this.ttl
}

// 可复用的空闲对象个数
func (this *Pool) Idle() int {
	return len(this.idle)
}

// 获取一个对象。
//   如果当前对象池中有空闲对象则直接返回，否则尝试构造新的对象。
//
// 参数：
//   wait: 等待时间。如果设置了 capacity、并且当前已经达到容量上限，则会等待其他协程归还（Put）对象。
//      wait = 0: 不等待，立即返回。
//      wait < 0: 一直等待，直到获得空闲对象。
func (this *Pool) Get(wait time.Duration) interface{} {
	v, _ := this.TryGet(wait)
	return v
}

// 从缓冲池中获取一个对象。
//   如果当前对象池中有空闲对象则直接返回，否则尝试构造新的对象。
//
// 参数：
//   wait: 等待时间。如果设置了 capacity、并且当前已经达到容量上限，则会等待其他协程归还（Put）对象。
//      wait = 0: 不等待，立即返回。
//      wait < 0: 一直等待，直到获得空闲对象。
func (this *Pool) TryGet(wait time.Duration) (interface{}, bool) {
	if v, ok := this.getNoWait(); ok {
		return v, true
	}

	if wait != 0 {
		timeout := false
		if wait > 0 {
			time.AfterFunc(wait, func() { timeout = true })
		}
		sleep := time.Duration(mathUtil.MinMaxInt64(int64(wait)/4, int64(time.Millisecond), int64(100*time.Millisecond)))
		for !timeout {
			if v, ok := this.getNoWait(); ok {
				return v, true
			}
			time.Sleep(sleep)
		}
	}

	return nil, false
}

func (this *Pool) getNoWait() (interface{}, bool) {
	this.mu.Lock()
	defer this.mu.Unlock()

	// 如果 idle 非空，取出空闲对象返回
	if len(this.idle) != 0 {
		v := this.idle[0].obj
		this.idle = this.idle[1:]
		this.size++
		return v, true
	}

	if this.capacity == -1 || this.size < this.capacity {
		this.size++
		return this.new(), true
	}

	return nil, false
}

// 归还对象
func (this *Pool) Put(x interface{}) {
	this.mu.Lock()
	defer this.mu.Unlock()

	if this.ttl > 0 && len(this.idle) <= 2048 {
		this.idle = append(this.idle, &wrapper{obj: x, expire: time.Now().Add(this.ttl)})
	} else {
		this.doDel(x)
	}
	this.size--
}

// 关闭缓冲池
func (this *Pool) Close() {
	if this.ticker != nil {
		this.ticker.Stop(0)
	}

	this.mu.Lock()
	defer this.mu.Unlock()

	for _, v := range this.idle {
		this.doDel(v.obj)
		v.obj = nil
	}
	this.idle = nil
}

func (this *Pool) CheckExpire() {
	this.doCheckExpire(time.Now())
}

func (this *Pool) doCheckExpire(t time.Time) {
	if len(this.idle) == 0 {
		return
	}

	this.mu.Lock()
	defer this.mu.Unlock()

	arr, cnt := make([]*wrapper, len(this.idle)), 0
	for _, v := range this.idle {
		if v.expire.After(t) {
			arr[cnt], cnt = v, cnt+1
		} else {
			this.doDel(v.obj)
			v.obj = nil
		}
	}
	if cnt != len(arr) {
		arr = arr[:cnt]
	}
	this.idle = arr
}

func (this *Pool) doDel(obj interface{}) {
	if this.del != nil {
		this.del(obj)
	} else if v, _ := obj.(interface{ Close() }); v != nil {
		v.Close()
	}
}
