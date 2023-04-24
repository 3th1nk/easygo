package arrUtil

import (
	"github.com/mohae/deepcopy"
	"reflect"
)

func Repeat(a interface{}, n int, f ...func(i int, a interface{}) interface{}) interface{} {
	var theF func(i int, a interface{}) interface{}
	if len(f) != 0 {
		theF = f[0]
	}

	arr := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(a)), 0, n)
	for i := 0; i < n; i++ {
		v := deepcopy.Copy(a)
		if theF != nil {
			v = theF(i, v)
		}
		arr = reflect.Append(arr, reflect.ValueOf(v))
	}
	return arr.Interface()
}

func RepeatString(a string, n int, f ...func(i int, s string) string) []string {
	var theF func(i int, s string) string
	if len(f) != 0 {
		theF = f[0]
	}

	arr := make([]string, n)
	for i := 0; i < n; i++ {
		var aa string
		if theF == nil {
			aa = a
		} else {
			aa = theF(i, a)
		}
		arr[i] = aa
	}
	return arr
}

func RepeatInt(a int, n int, f ...func(i int, a int) int) []int {
	var theF func(i int, a int) int
	if len(f) != 0 {
		theF = f[0]
	}

	arr := make([]int, n)
	for i := 0; i < n; i++ {
		var aa int
		if theF == nil {
			aa = a
		} else {
			aa = theF(i, a)
		}
		arr[i] = aa
	}
	return arr
}
