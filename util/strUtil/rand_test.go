package strUtil

import (
	"github.com/3th1nk/easygo/util"
	"github.com/3th1nk/easygo/util/mathUtil"
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRand(t *testing.T) {
	for i, n := 0, 5; i < n; i++ {
		util.Println(Rand(8))
	}
}

// 启动 worker 个协程、每个协程执行 loop 次 Rand 函数，统计性能
func TestRandPerf(t *testing.T) {
	length, worker, loop := 16, 100, 10000
	wg := sync.WaitGroup{}
	wg.Add(worker)
	start := time.Now()
	for i := 0; i < worker; i++ {
		go func() {
			defer wg.Done()
			for j, n := 0, loop; j < n; j++ {
				Rand(length)
			}
		}()
	}
	wg.Wait()
	took := time.Since(start)
	util.PrintTimeLn("total=%v, worker=%v, took=%v, %v/op", loop*worker, worker, took, took/time.Duration(loop))
}

// 启动 worker 个协程、每个协程执行 loop 次 Rand 函数，检测冲突数量。期望冲突数量为 0。
func TestRandConcurrent(t *testing.T) {
	length, worker, loop, conflict, rand := 12, 100, 10000, int32(0), int32(0)
	dict := sync.Map{}
	wg := sync.WaitGroup{}
	wg.Add(worker)
	start := time.Now()
	for i := 0; i < worker; i++ {
		go func() {
			defer wg.Done()
			for j, n := 0, loop; j < n; j++ {
				str := Rand(length)
				if _, loaded := dict.LoadOrStore(str, 1); loaded {
					util.PrintTimeLn("conflict: %v", str)
					atomic.AddInt32(&conflict, 1)
				} else {
					atomic.AddInt32(&rand, 1)
				}
			}
		}()
	}
	wg.Wait()
	took := time.Since(start)
	util.PrintTimeLn("total=%v, worker=%v, conflict=%v, rand=%v, took=%v, %v/op", loop*worker, worker, conflict, rand, took, took/time.Duration(loop))
	assert.Equal(t, 0, int(conflict))
}

// 测试 Rand 平均无冲突次数。
func TestRandConflict(t *testing.T) {
	for _, v := range [][]int{
		{6, 10, 1000000, 1},
		{8, 10, 1000000, 1},
	} {
		arr := testRandConflict(v[0], v[1], v[2], v[3])
		if len(arr) != 0 {
			avg := mathUtil.SumInt(arr) / len(arr)
			max := mathUtil.MaxInt(arr...)
			min := mathUtil.MinInt(arr...)
			t.Logf("rand(len=%v, cnt=%v*%v, tolerate=%v): avg=%v, max=%v, min=%v, [%v]", v[0], v[1], v[2], v[3], avg, max, min, JoinInt(arr, ","))
		} else {
			t.Logf("rand(len=%v, cnt=%v*%v, tolerate=%v): no conflict", v[0], v[1], v[2], v[3])
		}
	}
}

func testRandConflict(length, loop, count, tolerate int) (arr []int) {
	arr = make([]int, 0, loop)
	str, loaded := "", true
	for i := 0; i < loop; i++ {
		dict, size, conflict := sync.Map{}, 0, 0
		for j := 0; j < count; j++ {
			str = Rand(length)
			if _, loaded = dict.LoadOrStore(str, 1); loaded {
				if conflict = conflict + 1; conflict >= tolerate {
					arr = append(arr, size)
					break
				}
			}
			size++
		}
	}
	return
}
