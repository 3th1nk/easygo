package arrUtil

import "strings"

func DistinctInt(a []int) []int {
	if len(a) < 400 {
		return distinctIntWithArr(a)
	}
	return distinctIntWithMap(a)
}

func distinctIntWithArr(a []int) []int {
	b, n := make([]int, len(a)), 0
	for _, v := range a {
		if -1 == IndexOfInt(b[:n], v) {
			b[n], n = v, n+1
		}
	}
	return b[:n]
}

func distinctIntWithMap(a []int) []int {
	arr, m := make([]int, 0, len(a)), make(map[int]bool, len(a))
	for _, i := range a {
		if !m[i] {
			arr = append(arr, i)
			m[i] = true
		}
	}
	return arr
}

func DistinctInt64(a []int64) []int64 {
	if len(a) < 400 {
		return distinctInt64WithArr(a)
	}
	return distinctInt64WithMap(a)
}

func distinctInt64WithArr(a []int64) []int64 {
	b, n := make([]int64, len(a)), 0
	for _, v := range a {
		if -1 == IndexOfInt64(b[:n], v) {
			b[n], n = v, n+1
		}
	}
	return b[:n]
}

func distinctInt64WithMap(a []int64) []int64 {
	arr, m := make([]int64, 0, len(a)), make(map[int64]bool, len(a))
	for _, i := range a {
		if !m[i] {
			arr = append(arr, i)
			m[i] = true
		}
	}
	return arr
}

func DistinctInt32(a []int32) []int32 {
	if len(a) < 400 {
		distinctInt32WithArr(a)
	}
	return distinctInt32WithMap(a)
}

func distinctInt32WithArr(a []int32) []int32 {
	b, n := make([]int32, len(a)), 0
	for _, v := range a {
		if -1 == IndexOfInt32(b[:n], v) {
			b[n], n = v, n+1
		}
	}
	return b[:n]
}

func distinctInt32WithMap(a []int32) []int32 {
	arr, m := make([]int32, 0, len(a)), make(map[int32]bool, len(a))
	for _, i := range a {
		if !m[i] {
			arr = append(arr, i)
			m[i] = true
		}
	}
	return arr
}

func DistinctString(a []string, ignoreCase ...bool) []string {
	if len(a) < 300 {
		return distinctStringWithArr(a, ignoreCase...)
	}
	return distinctStringWithMap(a, ignoreCase...)
}

func distinctStringWithArr(a []string, ignoreCase ...bool) []string {
	b, n := make([]string, len(a)), 0
	for _, v := range a {
		if -1 == IndexOfString(b[:n], v, ignoreCase...) {
			b[n], n = v, n+1
		}
	}
	return b[:n]
}

func distinctStringWithMap(a []string, ignoreCase ...bool) []string {
	arr, m := make([]string, 0, len(a)), make(map[string]bool, len(a))
	if len(ignoreCase) == 0 || !ignoreCase[0] {
		for _, i := range a {
			if !m[i] {
				arr = append(arr, i)
				m[i] = true
			}
		}
		return arr
	}

	str := ""
	for _, i := range a {
		str = strings.ToLower(i)
		if !m[str] {
			arr = append(arr, i)
			m[str] = true
		}
	}
	return arr
}

func DistinctSortedInt(a []int) []int {
	b, n, m := make([]int, len(a)), 0, -1
	for _, v := range a {
		if n == 0 || b[m] != v {
			b[n], n, m = v, n+1, m+1
		}
	}
	return b[:n]
}

func DistinctSortedInt64(a []int64) []int64 {
	b, n, m := make([]int64, len(a)), 0, -1
	for _, v := range a {
		if n == 0 || b[m] != v {
			b[n], n, m = v, n+1, m+1
		}
	}
	return b[:n]
}

func DistinctSortedInt32(a []int32) []int32 {
	b, n, m := make([]int32, len(a)), 0, -1
	for _, v := range a {
		if n == 0 || b[m] != v {
			b[n], n, m = v, n+1, m+1
		}
	}
	return b[:n]
}

func DistinctSortedString(a []string, ignoreCase ...bool) []string {
	b, n, m := make([]string, len(a)), 0, -1
	if len(ignoreCase) != 0 && ignoreCase[0] {
		for _, v := range a {
			if n == 0 || !strings.EqualFold(b[m], v) {
				b[n], n, m = v, n+1, m+1
			}
		}
	} else {
		for _, v := range a {
			if n == 0 || b[m] != v {
				b[n], n, m = v, n+1, m+1
			}
		}
	}
	return b[:n]
}
