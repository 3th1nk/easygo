package strUtil

import (
	"strconv"
)

func ToIntArr(a []string) (val []int, err error) {
	n := len(a)
	val = make([]int, n)
	for i, s := range a {
		if v, e := strconv.ParseInt(s, 10, 64); e == nil {
			val[i] = int(v)
		} else if err == nil {
			err = e
		}
	}
	return
}

func ToInt32Arr(a []string) (val []int32, err error) {
	n := len(a)
	val = make([]int32, n)
	for i, s := range a {
		if v, e := strconv.ParseInt(s, 10, 64); e == nil {
			val[i] = int32(v)
		} else if err == nil {
			err = e
		}
	}
	return
}

func ToInt64Arr(a []string) (val []int64, err error) {
	n := len(a)
	val = make([]int64, n)
	for i, s := range a {
		if v, e := strconv.ParseInt(s, 10, 64); e == nil {
			val[i] = v
		} else if err == nil {
			err = e
		}
	}
	return
}
