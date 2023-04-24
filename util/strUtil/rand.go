package strUtil

import (
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// 生成随机字符串（默认使用小写字母和数字）
//
// 在一般场景下，可以使用随机字符串代替 UUID（UUID 太长了）：
//   UUID（Universally Unique Identifier）中的 Universally 就表名了他对重复率的目标：每秒产生10亿笔UUID，100年后产生一次重复的机率是50%。
//   实际场景中，我们只需要使 “重复机率小到可以忽略不计” 即可，即可认为其是 “事实唯一” 的。
//
// 注意：
//   随机函数并不保证唯一，根据随机长度不同，平均无冲突次数不同。可以根据业务场景选择适当的随机长度。
//   平均无冲突次数 是指：平均执行多少次 Rand 方法，会产生一次相同的随即结果。
//   由于随机性，平均无冲突次数 只是一个参考值，并不是绝对值。
//
// 重复机率计算：
//   随机字符串默认由 0-9、a-z 组成（可通过重载方法更改 Bucket）
//   重复机率：每一位有 36 种可能，n 位 对应 36^n 种可能，平均无冲突次数为 36^(n/2)。
//     6 位：平均无冲突次数为 36^3 = 46656
//     8 位：平均无冲突次数为 36^4 = 1679616 = 170万
//     10 位：平均无冲突次数为 36^5 = 60466176 = 6000万
//     12 位：平均无冲突次数为 36^6 = 2176782336 = 21亿
//     14 位：平均无冲突次数为 36^7 = 78364164096 = 780亿
//     16 位：平均无冲突次数为 36^8 = 2821109907456 = 2.8万亿
func Rand(length int) string {
	return doRand(length, []byte("0123456789abcdefghijklmnopqrstuvwxyz"), 36)
}

// 生成由小写字母组成的随机字符串。
//
// 注意：
//   随机函数并不保证唯一，请参考 Rand 方法的备注。
//
// 重复机率：每一位有 26 种可能，n 位 对应 26^n 种可能，平均无冲突次数为 26^(n/2)。
//    6 位：平均无冲突次数为 26^3 = 17576
//    8 位：平均无冲突次数为 26^4 = 456976 = 45万
//    10 位：平均无冲突次数为 26^5 = 11881376 = 1200万
//    12 位：平均无冲突次数为 26^6 = 308915776 = 3亿
//    14 位：平均无冲突次数为 26^7 = 8031810176 = 80亿
//    16 位：平均无冲突次数为 26^8 = 208827064576 = 2000亿
func RandL(length int) string {
	return doRand(length, []byte("abcdefghijklmnopqrstuvwxyz"), 26)
}

// 生成由大写字母组成的随机字符串。
// 重复机率：每一位有 26 种可能，n 位 对应 26^n 种可能，平均无冲突次数为 26^(n/2)。
//    6 位：平均无冲突次数为 26^3 = 17576
//    8 位：平均无冲突次数为 26^4 = 456976 = 45万
//    10 位：平均无冲突次数为 26^5 = 11881376 = 1200万
//    12 位：平均无冲突次数为 26^6 = 308915776 = 3亿
//    14 位：平均无冲突次数为 26^7 = 8031810176 = 80亿
//    16 位：平均无冲突次数为 26^8 = 208827064576 = 2000亿
func RandU(length int) string {
	return doRand(length, []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ"), 26)
}

// 生成由小写字母和数字组成的随机字符串。
// 重复机率：每一位有 36 种可能，n 位 对应 36^n 种可能，平均无冲突次数为 36^(n/2)。
//    6 位：平均无冲突次数为 36^3 = 46656
//    8 位：平均无冲突次数为 36^4 = 1679616 = 170万
//    10 位：平均无冲突次数为 36^5 = 60466176 = 6000万
//    12 位：平均无冲突次数为 36^6 = 2176782336 = 21亿
//    14 位：平均无冲突次数为 36^7 = 78364164096 = 780亿
//    16 位：平均无冲突次数为 36^8 = 2821109907456 = 2.8万亿
func RandLN(length int) string {
	return doRand(length, []byte("0123456789abcdefghijklmnopqrstuvwxyz"), 36)
}

// 生成由大写字母和数字组成的随机字符串。
// 重复机率：每一位有 36 种可能，n 位 对应 36^n 种可能，平均无冲突次数为 36^(n/2)。
//    6 位：平均无冲突次数为 36^3 = 46656
//    8 位：平均无冲突次数为 36^4 = 1679616 = 170万
//    10 位：平均无冲突次数为 36^5 = 60466176 = 6000万
//    12 位：平均无冲突次数为 36^6 = 2176782336 = 21亿
//    14 位：平均无冲突次数为 36^7 = 78364164096 = 780亿
//    16 位：平均无冲突次数为 36^8 = 2821109907456 = 2.8万亿
func RandUN(length int) string {
	return doRand(length, []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"), 36)
}

// 生成由大写字母、小写字母和数字组成的随机字符串。
// 重复机率：每一位有 62 种可能，n 位 对应 62^n 种可能，平均无冲突次数为 62^(n/2)。
//    6 位：平均无冲突次数为 62^3 = 238328 = 23万
//    8 位：平均无冲突次数为 62^4 = 14776336 = 1400万
//    10 位：平均无冲突次数为 62^5 = 916132832 = 9亿
//    12 位：平均无冲突次数为 62^6 = 56800235584 = 560亿
//    14 位：平均无冲突次数为 62^7 = 3521614606208 = 3.5万亿
func RandULN(length int) string {
	return doRand(length, []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"), 62)
}

// 根据指定的字符集生成随机字符串。
// 重复机率：每一位有 len(bucket) 种可能，长度 n 就有 len(bucket)^n 种可能。
func RandB(length int, bucket string) string {
	return doRand(length, []byte(bucket), len(bucket))
}

func doRand(length int, bucket []byte, size int) string {
	idx := atomic.AddInt64(&rdIdx, 1) & rdMask
	rd := rdList[idx]
	buf := make([]byte, length)
	rd.Lock()
	for i := length - 1; i >= 0; i-- {
		v := rd.Intn(size)
		buf[i] = bucket[v]
	}
	rd.Unlock()
	return string(buf)
}

const (
	rdCnt  = 2 << 4
	rdMask = rdCnt - 1
)

var (
	rdIdx, rdList = int64(0), func() (arr []*lockedRand) {
		ts := time.Now().UnixNano()
		base := rand.New(rand.NewSource(ts))
		arr = make([]*lockedRand, rdCnt)
		arr[0] = &lockedRand{Rand: base}
		for i := 1; i < rdCnt; i++ {
			arr[i] = &lockedRand{Rand: rand.New(rand.NewSource(ts + base.Int63n(math.MaxInt32)))}
		}
		return
	}()
)

type lockedRand struct {
	sync.Mutex
	*rand.Rand
}
