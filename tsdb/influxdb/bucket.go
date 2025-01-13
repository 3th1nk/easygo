package influxdb

import (
	"sync"
)

const (
	bucketSize = 10000 // 桶容量
)

func makeBucketGroupName(db, rp string) string {
	if rp == "" {
		return db
	}
	return db + "." + rp
}

// bucket 桶
//	写入相同db相同rp的数据会被聚合到同一个bucket中批量发送
type bucket struct {
	mu        sync.RWMutex
	size      int      // 桶容量
	flushSize int      // influxdb单次写入量上限
	lines     []string // 行数据
}

func newBucket(size, flushSize int) *bucket {
	return &bucket{
		size:      size,
		flushSize: flushSize,
		lines:     make([]string, 0, size),
	}
}

// findMaxMultiple 从一个大数值中取给定数值的最大倍数
func findMaxMultiple(bigNum, divisor int) int {
	if divisor == 0 {
		return 0 // 防止除以0的错误
	}
	return (bigNum / divisor) * divisor
}

// Push 添加数据，如果满了，则触发onFull回调
func (this *bucket) Push(lines []string, onFull func(lines []string) error) error {
	if len(lines) == 0 {
		return nil
	}

	this.mu.Lock()
	total := len(this.lines) + len(lines)
	if total < this.size {
		this.lines = append(this.lines, lines...)
		this.mu.Unlock()
		return nil
	}

	// 取单次写入量上限的最大整数倍，减少实际的批量写入次数
	popNum := findMaxMultiple(total, this.flushSize)

	all := make([]string, popNum)
	if popNum < this.size {
		copy(all, this.lines[:popNum])
		this.lines = append(this.lines[popNum:], lines...)
	} else {
		copy(all, this.lines)
		offset := popNum - len(this.lines)
		copy(all[len(this.lines):], lines[:offset])
		this.lines = append(this.lines[:0], lines[offset:]...)
	}
	this.mu.Unlock()

	return onFull(all)
}

func (this *bucket) Len() int {
	this.mu.RLock()
	defer this.mu.RUnlock()
	return len(this.lines)
}

func (this *bucket) Pop(n ...int) []string {
	this.mu.Lock()
	defer this.mu.Unlock()

	if len(n) == 0 {
		n = append(n, len(this.lines))
	}
	if n[0] > len(this.lines) {
		n[0] = len(this.lines)
	}
	if n[0] == 0 {
		return nil
	}

	arr := make([]string, n[0])
	copy(arr, this.lines[:n[0]])
	this.lines = this.lines[n[0]:]
	return arr
}

type bucketGroup struct {
	db      string
	rp      string // 为空时，表示默认rp
	size    int
	buckets []*bucket
}

func newBucketGroup(db, rp string, size, flushSize int) *bucketGroup {
	buckets := make([]*bucket, size)
	for i := range buckets {
		buckets[i] = newBucket(bucketSize, flushSize)
	}

	return &bucketGroup{
		db:      db,
		rp:      rp,
		size:    size,
		buckets: buckets,
	}
}

func (this *bucketGroup) Size() int {
	return this.size
}

func (this *bucketGroup) Db() string {
	return this.db
}

func (this *bucketGroup) Rp() string {
	return this.rp
}

func (this *bucketGroup) Get(idx int) *bucket {
	if idx < 0 || idx >= this.size {
		return nil
	}

	return this.buckets[idx]
}

func (this *bucketGroup) Push(idx int, lines []string, onFull func(lines []string) error) error {
	if bck := this.Get(idx); bck != nil {
		return bck.Push(lines, onFull)
	}
	return nil
}

func (this *bucketGroup) Range(fn func(bck *bucket)) {
	for _, bck := range this.buckets {
		fn(bck)
	}
}
