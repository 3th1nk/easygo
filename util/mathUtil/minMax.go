package mathUtil

import "math"

const (
	MAX_Int = int(^uint(0) >> 1)
	MIN_Int = -MAX_Int
)

func MaxInt(a ...int) (n int) {
	if len(a) == 0 {
		return
	}
	n = MIN_Int
	for _, v := range a {
		if t := v; t > n {
			n = t
		}
	}
	return
}

func MaxInt64(a ...int64) (n int64) {
	if len(a) == 0 {
		return
	}
	n = int64(math.MinInt64)
	for _, v := range a {
		if v > n {
			n = v
		}
	}
	return
}

func MaxInt32(a ...int32) (n int32) {
	if len(a) == 0 {
		return
	}
	n = int32(math.MinInt32)
	for _, v := range a {
		if v > n {
			n = v
		}
	}
	return
}

func MinInt(a ...int) (n int) {
	if len(a) == 0 {
		return 0
	}
	n = MAX_Int
	for _, v := range a {
		if t := v; t < n {
			n = t
		}
	}
	return n
}

func MinInt64(a ...int64) (n int64) {
	if len(a) == 0 {
		return
	}
	n = int64(math.MaxInt64)
	for _, v := range a {
		if v < n {
			n = v
		}
	}
	return
}

func MinInt32(a ...int32) (n int32) {
	if len(a) == 0 {
		return
	}
	n = int32(math.MaxInt32)
	for _, v := range a {
		if v < n {
			n = v
		}
	}
	return
}

// 将 a 的值限定在 [min, max] 区间内：如果 a < min 则返回 min；如果 a > max 则返回 max；否则返回 a 本身。
func MinMaxInt(a, min, max int) (n int) {
	if a < min {
		return min
	} else if a > max {
		return max
	} else {
		return a
	}
}

// 将 a 的值限定在 [min, max] 区间内：如果 a < min 则返回 min；如果 a > max 则返回 max；否则返回 a 本身。
func MinMaxInt64(a, min, max int64) (n int64) {
	if a < min {
		return min
	} else if a > max {
		return max
	} else {
		return a
	}
}

// 将 a 的值限定在 [min, max] 区间内：如果 a < min 则返回 min；如果 a > max 则返回 max；否则返回 a 本身。
func MinMaxInt32(a, min, max int32) (n int32) {
	if a < min {
		return min
	} else if a > max {
		return max
	} else {
		return a
	}
}
