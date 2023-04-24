package timeRoundedCounter

import (
	"fmt"
	"github.com/3th1nk/easygo/internal/_test"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestCounter_GetAll(t *testing.T) {
	const format = "15:04:05"
	cl := newTestCounter(time.Minute, 10)
	now := time.Now().Truncate(time.Minute)

	// 第一次增加数据， maxNS 被置为当前时间
	cl.Add(nil, now.Add(-5*time.Minute))
	assert.Equal(t, now.Add(-5*time.Minute).Format(format), timeUtil.FromNS(cl.base.maxNS).Format(format))
	assert.Equal(t, 0, cl.base.maxPos)
	assert.Equal(t, "0,0,0,0,0,1", toString(cl.base.GetAll()))

	// 增加了更早时间的数据， maxNS 不变
	cl.Add(nil, now.Add(-1*time.Minute))
	assert.Equal(t, now.Add(-1*time.Minute).Format(format), timeUtil.FromNS(cl.base.maxNS).Format(format))
	assert.Equal(t, 4, cl.base.maxPos)
	assert.Equal(t, "0,1,0,0,0,1", toString(cl.base.GetAll()))

	// 增加了更早时间的数据， maxNS 不变
	cl.Add(nil, now)
	assert.Equal(t, now.Format(format), timeUtil.FromNS(cl.base.maxNS).Format(format))
	assert.Equal(t, 5, cl.base.maxPos)
	assert.Equal(t, "1,1,0,0,0,1", toString(cl.base.GetAll()))
}

func TestCounter_MaxNS(t *testing.T) {
	const format = "15:04:05"
	cl := newTestCounter(time.Minute, 10)
	now := timeUtil.ParseNoErr("2022-01-01 12:00:00").Truncate(cl.base.round)

	// 第一次增加数据， maxNS 被置为当前时间
	cl.Add(nil, now)
	assert.Equal(t, now.Format(format), timeUtil.FromNS(cl.base.maxNS).Format(format))
	assert.Equal(t, 0, cl.base.maxPos)
	assert.Equal(t, "1", toString(cl.base.GetN(now, 0, -1)))

	// 增加了更早时间的数据， maxNS 不变
	cl.Add(nil, now.Add(-1*time.Minute))
	assert.Equal(t, now.Format(format), timeUtil.FromNS(cl.base.maxNS).Format(format))
	assert.Equal(t, 0, cl.base.maxPos)
	assert.Equal(t, "1,1", toString(cl.base.GetN(now, 0, -1)))

	// 增加了更晚时间的数据，maxNS 改变
	cl.Add(nil, now.Add(1*time.Minute))
	assert.Equal(t, now.Add(1*cl.base.round).Format(format), timeUtil.FromNS(cl.base.maxNS).Format(format))
	assert.Equal(t, 1, cl.base.maxPos)
	assert.Equal(t, "1,1,1", toString(cl.base.GetN(now, 0, -1)))

	// 改变了现有数据， maxNS 不变
	cl.Add(nil, now.Add(-1*time.Minute))
	assert.Equal(t, now.Add(1*cl.base.round).Format(format), timeUtil.FromNS(cl.base.maxNS).Format(format))
	assert.Equal(t, 1, cl.base.maxPos)
	assert.Equal(t, "1,1,2", toString(cl.base.GetN(now, 0, -1)))

	// 增加了更晚时间的数据，maxNS 改变
	cl.Add(nil, now.Add(1*time.Minute))
	assert.Equal(t, now.Add(1*cl.base.round).Format(format), timeUtil.FromNS(cl.base.maxNS).Format(format))
	assert.Equal(t, 1, cl.base.maxPos)
	assert.Equal(t, "2,1,2", toString(cl.base.GetN(now, 0, -1)))

	// 增加了更晚时间的数据，maxNS 改变
	cl.Add(nil, now.Add(8*time.Minute))
	assert.Equal(t, now.Add(8*cl.base.round).Format(format), timeUtil.FromNS(cl.base.maxNS).Format(format))
	assert.Equal(t, 8, cl.base.maxPos)
	assert.Equal(t, "1,0,0,0,0,0,0,2,1,2", toString(cl.base.GetN(now, 0, -1)))

	// 新增的数据跨越了尾部，maxNS 改变、尾部数据被改写
	cl.Add(nil, now.Add(9*time.Minute))
	assert.Equal(t, now.Add(9*cl.base.round).Format(format), timeUtil.FromNS(cl.base.maxNS).Format(format))
	assert.Equal(t, 9, cl.base.maxPos)
	assert.Equal(t, "1,1,0,0,0,0,0,0,2,1", toString(cl.base.GetN(now, 0, -1)))

	// 新增的数据跨越了尾部，maxNS 改变、尾部数据被改写
	cl.Add(nil, now.Add(10*time.Minute))
	assert.Equal(t, now.Add(10*cl.base.round).Format(format), timeUtil.FromNS(cl.base.maxNS).Format(format))
	assert.Equal(t, 0, cl.base.maxPos)
	assert.Equal(t, "1,1,1,0,0,0,0,0,0,2", toString(cl.base.GetN(now, 0, -1)))

	// 新增数据跨越了整个区间，maxNS 改变、所有数据被重置
	cl.Add(nil, now.Add(20*time.Minute))
	assert.Equal(t, now.Add(20*cl.base.round).Format(format), timeUtil.FromNS(cl.base.maxNS).Format(format))
	assert.Equal(t, 0, cl.base.maxPos)
	assert.Equal(t, "1,0,0,0,0,0,0,0,0,0", toString(cl.base.GetN(now, 0, -1, false)))

	// 新增的数据早于最早的时间，请求会被忽略，不会改动任何数据
	cl.Add(nil, now.Add(10*time.Minute))
	assert.Equal(t, now.Add(20*cl.base.round).Format(format), timeUtil.FromNS(cl.base.maxNS).Format(format))
	assert.Equal(t, 0, cl.base.maxPos)
	assert.Equal(t, "1,0,0,0,0,0,0,0,0,0", toString(cl.base.GetN(now, 0, -1, false)))

	// 新增的数据早于最早的时间，请求会被忽略，不会改动任何数据
	cl.Add(nil, now.Add(8*time.Minute))
	assert.Equal(t, now.Add(20*cl.base.round).Format(format), timeUtil.FromNS(cl.base.maxNS).Format(format))
	assert.Equal(t, 0, cl.base.maxPos)
	assert.Equal(t, "1,0,0,0,0,0,0,0,0,0", toString(cl.base.GetN(now, 0, -1, false)))
}

// Perf:
//   [perf-1s]  total=14877334,  avg=14877334/s,  mrt=67ns
func TestCounter_Add_Perf(t *testing.T) {
	cl := newTestCounter(time.Minute, 10)
	ts := time.Now()
	cl.Add(nil, ts)
	_test.Perf(func(i int) {
		cl.Add(nil, ts)
		ts = ts.Add(time.Millisecond)
	})
}

func TestCounter_GetN(t *testing.T) {
	cl := newTestCounter(time.Minute, 10)
	now := time.Now()

	for i := -20; i < -10; i++ {
		cl.Add(&testCountItem{val: i}, now.Add(time.Duration(i)*time.Minute))
	}
	assert.Equal(t, "-11,-12,-13,-14,-15,-16,-17,-18,-19,-20", toString(cl.base.GetN(timeUtil.FromNS(cl.base.maxNS), 0, -1)))
	assert.Equal(t, "", toString(cl.base.GetAll()))

	for i := -5; i <= 0; i++ {
		cl.Add(&testCountItem{val: i}, now.Add(time.Duration(i)*time.Minute))
	}
	assert.Equal(t, "0,-1,-2,-3,-4,-5", toString(cl.base.GetAll()))
	assert.Equal(t, "0,-1,-2,-3,-4,-5,0,0,0,0", toString(cl.base.GetAll(false)))
}

func TestCounter_GetAll_Perf(t *testing.T) {
	cl := newTestCounter(time.Minute, 60)
	now := time.Now()

	for i := -5; i <= 0; i++ {
		cl.Add(&testCountItem{val: i}, now.Add(time.Duration(i)*time.Minute))
	}
	assert.Equal(t, "0,-1,-2,-3,-4,-5", toString(cl.base.GetAll()))

	_test.Perf(func(i int) {
		cl.base.GetAll()
	})

	_test.Perf(func(i int) {
		cl.base.GetAll(false)
	})
}

func newTestCounter(round time.Duration, capacity int) *testCounter {
	c := &testCounter{}
	c.base = doNew("", round, capacity, func() CounterItem {
		return &testCountItem{}
	})
	return c
}

func toString(arr []CounterItem) string {
	dst := make([]string, len(arr))
	for i, v := range arr {
		dst[i] = fmt.Sprintf("%v", v)
	}
	return strings.Join(dst, ",")
}

type testCounter struct {
	base *CustomCounter
}

func (this *testCounter) Name() string {
	return this.base.name
}

func (this *testCounter) Round() time.Duration {
	return this.base.round
}

func (this *testCounter) ItemProperties() []string {
	return []string{"Val"}
}

func (this *testCounter) Add(a *testCountItem, t ...time.Time) {
	this.base.Add(func(v CounterItem) {
		if a == nil {
			v.(*testCountItem).val++
		} else {
			v.(*testCountItem).val += a.val
		}
	}, t...)
}

type testCountItem struct {
	val int
}

func (this *testCountItem) String() string {
	return strconv.Itoa(this.val)
}

func (this *testCountItem) IsZero() bool {
	return this.val == 0
}

func (this *testCountItem) Reset() {
	this.val = 0
}
