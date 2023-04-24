package mathUtil

func SumFloat(a []float64) (n float64) {
	for _, v := range a {
		n += v
	}
	return n
}

func SumInt64(a []int64) (n int64) {
	for _, v := range a {
		n += v
	}
	return n
}

func SumInt32(a []int32) (n int32) {
	for _, v := range a {
		n += v
	}
	return n
}

func SumInt(a []int) (n int) {
	for _, v := range a {
		n += v
	}
	return n
}
