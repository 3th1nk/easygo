package arrUtil

import "strconv"

func IntToObj(a []int) []interface{} {
	arr := make([]interface{}, len(a))
	for i, v := range a {
		arr[i] = v
	}
	return arr
}

func IntToStr(a []int) []string {
	arr := make([]string, len(a))
	for i, v := range a {
		arr[i] = strconv.FormatInt(int64(v), 10)
	}
	return arr
}

func IntTo64(a []int) []int64 {
	arr := make([]int64, len(a))
	for i, v := range a {
		arr[i] = int64(v)
	}
	return arr
}

func Int64ToObj(a []int64) []interface{} {
	arr := make([]interface{}, len(a))
	for i, v := range a {
		arr[i] = v
	}
	return arr
}

func Int64ToStr(a []int64) []string {
	arr := make([]string, len(a))
	for i, v := range a {
		arr[i] = strconv.FormatInt(v, 10)
	}
	return arr
}

func Int64ToInt(a []int64) []int {
	arr := make([]int, len(a))
	for i, v := range a {
		arr[i] = int(v)
	}
	return arr
}

func IntToMap(arr []int) map[int]bool {
	m := make(map[int]bool, len(arr))
	for _, val := range arr {
		m[val] = true
	}
	return m
}

func Int64ToMap(arr []int64) map[int64]bool {
	m := make(map[int64]bool, len(arr))
	for _, val := range arr {
		m[val] = true
	}
	return m
}
