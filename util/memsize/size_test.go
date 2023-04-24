package memsize

import (
	"github.com/3th1nk/easygo/util"
	"github.com/3th1nk/easygo/util/mathUtil"
	"github.com/3th1nk/easygo/util/strUtil"
	"math/rand"
	"testing"
	"time"
)

func TestSize(t *testing.T) {
	arrB := make([]*tmpB, 500)
	for i := range arrB {
		arrB[i] = &tmpB{
			Id:   rand.Int(),
			Name: strUtil.Rand(20 + rand.Intn(20)),
		}
		if n := rand.Intn(1000); n <= i {
			arrB[i].A = arrB[n]
		}
	}

	arrA := make([]*tmpA, 10000)
	for i := range arrA {
		arrA[i] = &tmpA{
			Id:     rand.Int(),
			Name:   strUtil.Rand(10 + rand.Intn(40)),
			IntArr: newIntArr(rand.Intn(10)),
			StrArr: newStrArr(rand.Intn(10), 10, 50),
			Dict0:  map[int]int{},
			Dict1:  map[string]int{},
			Dict2:  map[string]*tmpB{},
		}
		for j := 0; j < 100; j++ {
			arrA[i].Dict0[rand.Int()] = rand.Int()
		}
		for j := 0; j < 100; j++ {
			arrA[i].Dict1[strUtil.Rand(8)] = rand.Int()
		}
		for j := 0; j < 5; j++ {
			arrA[i].Dict2[strUtil.Rand(8)] = arrB[rand.Intn(len(arrB))]
		}
		if rand.Int()%3 == 0 {
			arrA[i].Obj = &tmpB{
				Id:   rand.Int(),
				Name: strUtil.Rand(20 + rand.Intn(20)),
				A: &tmpB{
					Id:   rand.Int(),
					Name: strUtil.Rand(20 + rand.Intn(20)),
				},
				B: []*tmpB{},
			}
		} else {
			arrA[i].Obj = arrB[rand.Intn(len(arrB))]
		}
		arrA[i].ObjArr = make([]*tmpB, mathUtil.MaxInt(0, rand.Intn(10)-5))
		for j := range arrA[i].ObjArr {
			arrA[i].ObjArr[j] = arrB[rand.Intn(len(arrB))]
		}
	}

	{
		start := time.Now()
		size := 0
		for _, val := range arrA {
			size += Size(val)
		}
		util.Println("memsize=%v, took=%v", size, time.Since(start).Round(time.Millisecond))
	}
}

type tmpA struct {
	Id     int
	Name   string
	IntArr []int
	StrArr []string
	Dict0  map[int]int
	Dict1  map[string]int
	Dict2  map[string]*tmpB
	Obj    *tmpB
	ObjArr []*tmpB
}

type tmpB struct {
	Id   int
	Name string
	A    *tmpB
	B    []*tmpB
}

func newIntArr(n int) (arr []int) {
	arr = make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rand.Int()
	}
	return
}

func newStrArr(n int, minLen, maxLen int) (arr []string) {
	arr = make([]string, n)
	for i := 0; i < n; i++ {
		arr[i] = strUtil.Rand(minLen + rand.Intn(maxLen-minLen))
	}
	return
}
