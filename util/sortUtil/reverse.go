package sortUtil

import (
	"reflect"
)

func ReverseInt(slice []int) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		tmp := slice[i]
		slice[i], slice[j] = slice[j], tmp
	}
}

func ReverseInt64(slice []int64) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		tmp := slice[i]
		slice[i], slice[j] = slice[j], tmp
	}
}

func ReverseFloat(slice []float64) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		tmp := slice[i]
		slice[i], slice[j] = slice[j], tmp
	}
}

func ReverseString(slice []string) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		tmp := slice[i]
		slice[i], slice[j] = slice[j], tmp
	}
}

func Reverse(slice interface{}) {
	reflectVal := reflect.ValueOf(slice)
	if kind := reflectVal.Kind(); kind != reflect.Array && kind != reflect.Slice {
		panic("参数 slice 必须是切片类型")
	}

	for i, j := 0, reflectVal.Len()-1; i < j; i, j = i+1, j-1 {
		v1, v2 := reflectVal.Index(i), reflectVal.Index(j)
		tmp := reflect.ValueOf(v1.Interface())
		v1.Set(v2)
		v2.Set(tmp)
	}
}
