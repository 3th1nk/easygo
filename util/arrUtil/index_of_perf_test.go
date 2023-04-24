package arrUtil

import (
	"github.com/3th1nk/easygo/util/strUtil"
	"math/rand"
	"testing"
	"time"
)

func TestIndexOfInt_Pref(t *testing.T) {
	perfTestInt(t, 100, 100)
	perfTestInt(t, 10000, 100)
}

func TestIndexOfString_Pref(t *testing.T) {
	perfTestString(t, 100, 100)
	perfTestString(t, 1000, 100)
	perfTestString(t, 10000, 100)
}

func TestIndexOfObject_Pref(t *testing.T) {
	perfTestObj(t, 100, 100)
	perfTestObj(t, 10000, 100)
}

func perfTestInt(t *testing.T, arrSize, findSize int) {
	pt := &indexOfPerfTest{arrSize: arrSize, findSize: findSize, elemType: "int"}
	// arr: 被查找的数组
	// find1: 要查找的数字全部都在 arr 里面
	// find2: 要查找的数字有 50% 在 arr 里面
	// find3: 要查找的数字全都不在 arr 里面
	arr, find1, find2, find3 := make([]int, pt.arrSize), make([]int, pt.findSize), make([]int, pt.findSize), make([]int, pt.findSize)
	for i := 0; i < pt.arrSize; i++ {
		arr[i] = rand.Intn(pt.arrSize * 5)
	}
	for i := 0; i < pt.findSize; i++ {
		find1[i] = arr[rand.Intn(pt.arrSize)]
		if i%2 == 0 {
			find2[i] = arr[rand.Intn(pt.arrSize)]
		} else {
			find2[i] = pt.arrSize + rand.Intn(pt.arrSize)
		}
		find3[i] = pt.arrSize + rand.Intn(pt.arrSize)
	}
	pt.run(t, 100, func(i int) {
		IndexOfInt(arr, find1[i])
	})
	pt.run(t, 50, func(i int) {
		IndexOfInt(arr, find2[i])
	})
	pt.run(t, 0, func(i int) {
		IndexOfInt(arr, find3[i])
	})
}

func perfTestString(t *testing.T, arrSize, findSize int) {
	pt := &indexOfPerfTest{arrSize: arrSize, findSize: findSize, elemType: "string"}
	// arr: 被查找的数组
	// find1: 要查找的数字全部都在 arr 里面
	// find2: 要查找的数字有 50% 在 arr 里面
	// find3: 要查找的数字全都不在 arr 里面
	arr, find1, find2, find3 := make([]string, pt.arrSize), make([]string, pt.findSize), make([]string, pt.findSize), make([]string, pt.findSize)
	for i := 0; i < pt.arrSize; i++ {
		arr[i] = strUtil.Rand(64)
	}
	for i := 0; i < pt.findSize; i++ {
		find1[i] = arr[rand.Intn(pt.arrSize)]
		if i%2 == 0 {
			find2[i] = arr[rand.Intn(pt.arrSize)]
		} else {
			find2[i] = strUtil.Rand(64)
		}
		find3[i] = strUtil.Rand(64)
	}
	pt.run(t, 100, func(i int) {
		IndexOfString(arr, find1[i])
	})
	pt.run(t, 50, func(i int) {
		IndexOfString(arr, find2[i])
	})
	pt.run(t, 0, func(i int) {
		IndexOfString(arr, find3[i])
	})
}

func perfTestObj(t *testing.T, arrSize, findSize int) {
	type tmp struct {
		Id   int
		Name string
	}
	pt := &indexOfPerfTest{arrSize: arrSize, findSize: findSize, elemType: "struct"}
	// arr: 被查找的数组
	// find1: 要查找的数字全部都在 arr 里面
	// find2: 要查找的数字有 50% 在 arr 里面
	// find3: 要查找的数字全都不在 arr 里面
	arr, find1, find2, find3 := make([]*tmp, pt.arrSize), make([]*tmp, pt.findSize), make([]*tmp, pt.findSize), make([]*tmp, pt.findSize)
	for i := 0; i < pt.arrSize; i++ {
		arr[i] = &tmp{
			Id:   rand.Intn(pt.arrSize * 5),
			Name: strUtil.Rand(12),
		}
	}
	for i := 0; i < pt.findSize; i++ {
		find1[i] = arr[rand.Intn(pt.arrSize)]
		if i%2 == 0 {
			find2[i] = &tmp{
				Id:   arr[rand.Intn(pt.arrSize)].Id,
				Name: arr[rand.Intn(pt.arrSize)].Name,
			}
		} else {
			find2[i] = &tmp{
				Id:   pt.arrSize + rand.Intn(pt.arrSize),
				Name: strUtil.Rand(12),
			}
		}
		find3[i] = &tmp{
			Id:   pt.arrSize + rand.Intn(pt.arrSize),
			Name: strUtil.Rand(12),
		}
	}
	pt.run(t, 100, func(j int) {
		IndexOf(arr, func(i int) bool { return arr[i].Name == find1[j].Name })
	})
	pt.run(t, 50, func(j int) {
		IndexOf(arr, func(i int) bool { return arr[i].Name == find2[j].Name })
	})
	pt.run(t, 0, func(j int) {
		IndexOf(arr, func(i int) bool { return arr[i].Name == find3[j].Name })
	})
}

type indexOfPerfTest struct {
	arrSize  int
	findSize int
	elemType string
}

// 'IndexOf' 函数的性能测试：
// 打印在长度为 ArrSize 的切片中查找，每秒钟能够执行的次数。
// matchPercent 表示用于查找的字符串包含在切片中（返回结果不等于-1）的比例。
//     100 表示要查找的元素存在于切片中，平均遍历一半就能找到并返回其索引下标；
//     0 表示要查找的元素不存在与切片中，将完整的遍历整个切片、最终返回 -1。
func (this *indexOfPerfTest) run(t *testing.T, matchPercent int, f func(i int)) {
	d, stop, c, i := time.Second, false, 0, 0
	time.AfterFunc(d, func() {
		stop = true
	})
	for !stop {
		f(i)
		if i = i + 1; i == this.findSize {
			c, i = c+1, 0
		}
	}
	n := c*this.findSize + i
	t.Logf("在 %v[%v] 中查找(命中率 %d%%):  %0.f/sec, %v/loop", this.elemType, this.arrSize, matchPercent, float64(n)/d.Seconds(), time.Duration(float64(d)/float64(n)))
}
