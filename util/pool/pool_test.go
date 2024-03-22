package pool

import (
	"github.com/3th1nk/easygo/util"
	"github.com/3th1nk/easygo/util/timeUtil"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestPool_GetCap_1(t *testing.T) {
	const capacity = 10
	pool := New(capacity, 0, func() interface{} { return rand.Intn(999999) }, nil)
	for i := 0; i < capacity+10; i++ {
		v := pool.Get(0)
		if i < capacity {
			assert.NotNil(t, v)
		} else {
			assert.Nil(t, v)
		}
	}
	assert.Equal(t, 0, pool.Idle())
	t.Logf("Size: %d, Idle: %d", pool.Size(), pool.Idle())
}

func TestPool_GetCap_2(t *testing.T) {
	const capacity = 10
	const wait = 100 * time.Millisecond
	pool := New(capacity, 0, func() interface{} { return rand.Intn(999999) }, nil)
	for i, n := 0, capacity+10; i < n; i++ {
		start := time.Now()
		v := pool.Get(wait)
		ts := time.Since(start)
		if i < capacity {
			assert.NotNil(t, v)
			assert.GreaterOrEqual(t, time.Millisecond, ts, ts)
		} else {
			assert.Nil(t, v)
			assert.LessOrEqual(t, wait, ts, ts)
		}
	}
	assert.Equal(t, 0, pool.Idle())
	t.Logf("Size: %d, Idle: %d", pool.Size(), pool.Idle())
}

func TestPool_TryGet(t *testing.T) {
	const capacity = 10
	pool := New(capacity, 0, func() interface{} { return rand.Intn(999999) }, nil)
	for i, n := 0, capacity+10; i < n; i++ {
		start := time.Now()
		v, ok := pool.TryGet(0)
		ts := time.Since(start)
		assert.GreaterOrEqual(t, time.Millisecond, ts, ts)
		if i < capacity {
			assert.NotNil(t, v)
			assert.Equal(t, true, ok)
		} else {
			assert.Nil(t, v)
			assert.Equal(t, false, ok)
		}
	}
	assert.Equal(t, 0, pool.Idle())
	t.Logf("Size: %d, Idle: %d", pool.Size(), pool.Idle())
}

func TestPool_PutNoCap(t *testing.T) {
	pool := New(-1, 0, func() interface{} { return rand.Intn(999999) }, nil)
	for i := 0; i < 10000; i++ {
		pool.Put(rand.Intn(999999))
	}
	assert.Equal(t, 0, pool.Idle())
	t.Logf("Size: %d, Idle: %d", pool.Size(), pool.Idle())
}

func TestPool_TTL(t *testing.T) {
	const ttl = 100 * time.Millisecond
	const count = 10

	got, expire := 0, 0
	pool := New(-1, ttl, func() interface{} { return rand.Intn(999999) }, func(x interface{}) { expire++ })
	for i, n := 0, count*3; i < n; i++ {
		v := pool.Get(0)
		assert.NotNil(t, v)
		got++
	}
	t.Logf("Got: %d, Size: %d, Idle: %d", got, pool.Size(), pool.Idle())

	for i := 0; i < count; i++ {
		pool.Put(rand.Intn(999999))
		time.Sleep(ttl / count)
	}
	assert.Equal(t, count, expire+pool.Idle())
	t.Logf("Put: %d, Size: %d, Idle: %d, Expire: %d", count, pool.Size(), pool.Idle(), expire)
}

func TestPool_TTL_Expire(t *testing.T) {
	const ttl = 100 * time.Millisecond
	const count = 10

	expire := 0
	pool := New(-1, ttl, func() interface{} { return rand.Intn(999999) }, func(x interface{}) { expire++ })

	for i := 0; i < count; i++ {
		pool.Put(rand.Intn(999999))
		time.Sleep(ttl / count)
	}
	assert.Equal(t, count, expire+pool.Idle())
	t.Logf("Put: %d, Size: %d, Idle: %d, Expire: %d", count, pool.Size(), pool.Idle(), expire)

	time.Sleep(ttl)
	pool.CheckExpire()
	assert.Equal(t, 0, pool.Idle())
	t.Logf("Put: %d, Size: %d, Idle: %d, Expire: %d", count, pool.Size(), pool.Idle(), expire)
}

func TestPoolConcurrence(t *testing.T) {
	const ttl = 100 * time.Millisecond
	pool := New(10000, ttl, func() interface{} { return rand.Intn(999999) }, nil)

	const worker, loop = 100, 5000
	var got, put, done int32
	var maxIdle = 0

	ticker := timeUtil.NewTicker(500*time.Millisecond, 500*time.Millisecond, func(t time.Time) {
		util.PrintTimeLn("got: %d, put: %d, done: %d, idle: %d, maxIdle: %d", got, put, done, pool.Idle(), maxIdle)
	})

	wg := sync.WaitGroup{}
	wg.Add(worker + 1)
	for i := 0; i < worker; i++ {
		go func() {
			workerGot := 0
			for workerGot != loop {
				if v, ok := pool.TryGet(0 * time.Millisecond); ok {
					workerGot++
					atomic.AddInt32(&got, 1)
					time.AfterFunc(time.Duration(10+rand.Int63n(50))*time.Millisecond, func() {
						pool.Put(v)
						if n := pool.Idle(); n > maxIdle {
							maxIdle = n
						}
						if atomic.AddInt32(&put, 1) == worker*loop {
							wg.Done()
						}
					})
				}
			}
			atomic.AddInt32(&done, 1)
			wg.Done()
		}()
	}
	wg.Wait()
	ticker.Trigger()

	time.Sleep(ttl)
	ticker.Trigger()

	ticker.Stop(0)
}
