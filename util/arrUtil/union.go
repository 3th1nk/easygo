package arrUtil

// 取两个 []int 的并集
func UnionInt(a ...[]int) []int {
	// 过滤空切片
	b := make([][]int, 0, len(a))
	// b 中所有切片的长度之和
	var bn int
	for _, v := range a {
		if n := len(v); n != 0 {
			b = append(b, v)
			bn += n
		}
	}

	if n := len(b); n == 0 {
		return nil
	} else if n == 1 {
		return b[0]
	}

	c := make([]int, len(b[0]), bn)
	copy(c, b[0])
	for _, arr := range b[1:] {
		for _, n := range arr {
			if -1 == IndexOfInt(c, n) {
				c = append(c, n)
			}
		}
	}
	return c
}

// 取两个 []int32 的并集
func UnionInt32(a ...[]int32) []int32 {
	// 过滤空切片
	b := make([][]int32, 0, len(a))
	// b 中所有切片的长度之和
	var bn int
	for _, v := range a {
		if n := len(v); n != 0 {
			b = append(b, v)
			bn += n
		}
	}

	if n := len(b); n == 0 {
		return nil
	} else if n == 1 {
		return b[0]
	}

	c := make([]int32, len(b[0]), bn)
	copy(c, b[0])
	for _, arr := range b[1:] {
		for _, n := range arr {
			if -1 == IndexOfInt32(c, n) {
				c = append(c, n)
			}
		}
	}
	return c
}

// 取两个 []int64 的并集
func UnionInt64(a ...[]int64) []int64 {
	// 过滤空切片
	b := make([][]int64, 0, len(a))
	// b 中所有切片的长度之和
	var bn int
	for _, v := range a {
		if n := len(v); n != 0 {
			b = append(b, v)
			bn += n
		}
	}

	if n := len(b); n == 0 {
		return nil
	} else if n == 1 {
		return b[0]
	}

	c := make([]int64, len(b[0]), bn)
	copy(c, b[0])
	for _, arr := range b[1:] {
		for _, n := range arr {
			if -1 == IndexOfInt64(c, n) {
				c = append(c, n)
			}
		}
	}
	return c
}

// 取两个数组的并集。
func UnionString(a, b []string, ignoreCase ...bool) []string {
	if a == nil {
		return b
	} else if b == nil {
		return a
	}
	if len(a) == 0 {
		return b
	} else if len(b) == 0 {
		return a
	}

	c := make([]string, len(a), len(a)+len(b))
	copy(c, a)
	for _, s := range b {
		if -1 == IndexOfString(a, s, ignoreCase...) {
			c = append(c, s)
		}
	}
	return c
}
