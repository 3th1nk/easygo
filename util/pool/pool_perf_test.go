package pool

import (
	"fmt"
	"github.com/3th1nk/easygo/internal/_test"
	"github.com/3th1nk/easygo/util"
	"github.com/panjf2000/ants/v2"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestPool_Get_Perf(t *testing.T) {
	capacity := []int{-1, 10000}
	ttl := []time.Duration{0, time.Minute}
	wait := []time.Duration{0, 100 * time.Millisecond, time.Second}
	for _, c := range capacity {
		for _, tt := range ttl {
			for _, w := range wait {
				pool := New(c, tt, func() interface{} { return 1 }, nil)
				result := _test.PerfIf(func(i int) bool {
					return nil != pool.Get(w)
				}, _test.PerfOptions{Name: fmt.Sprintf("cap=%v, ttl=%v, wait=%v", pool.Capacity(), pool.TTL(), w)})
				if pool.Capacity() > 0 {
					assert.True(t, int(result.Ok) <= pool.Capacity()+10, fmt.Sprintf("%v / %v", result.Ok, pool.Capacity()))
				}
			}
		}
	}
}

func TestAntsPool_Perf(t *testing.T) {
	test := func(loop int, d time.Duration) {
		wg := sync.WaitGroup{}

		{
			wg.Add(loop)
			start := time.Now()
			for i := 0; i < loop; i++ {
				go func() {
					time.Sleep(d)
					wg.Done()
				}()
			}
			wg.Wait()
			util.Println("go(%v, %v): took=%v", loop, d, time.Since(start).Round(time.Millisecond))
		}

		{
			pool, _ := ants.NewPool(1000)
			wg.Add(loop)
			start := time.Now()
			for i := 0; i < loop; i++ {
				pool.Submit(func() {
					time.Sleep(d)
					wg.Done()
				})
			}
			wg.Wait()
			util.Println("go(%v, %v): took=%v", loop, d, time.Since(start).Round(time.Millisecond))
		}
	}

	test(1000000, 200*time.Millisecond)
}

// go=1,  ok=2665154,  fail=0,  avg=2665154/s,  mrt=375ns
// go=50,  ok=1509606,  fail=0,  avg=1509606/s,  mrt=662ns
// go=1000,  ok=12616744,  fail=0,  avg=1261674/s,  mrt=792ns
// go=2000,  ok=12187689,  fail=0,  avg=1218769/s,  mrt=820ns
// go=10000,  ok=11159349,  fail=0,  avg=1115935/s,  mrt=896ns
func TestPool_Perf(t *testing.T) {
	pool := New(2000, 5*time.Second, func() interface{} { return 1 }, nil)
	_test.PerfIf(func(i int) bool {
		if v := pool.Get(-1); v != nil {
			pool.Put(v)
			return true
		}
		return false
	}, _test.PerfOptions{
		Name:      "pool",
		Goroutine: 1000,
	})
}
