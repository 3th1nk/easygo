package arrUtil

func ShiftString(a []string) (out []string, val string, ok bool) {
	if len(a) != 0 {
		return a[1:], a[0], true
	}
	out = a
	return
}

func MustShiftString(a []string, defaultVal ...string) (out []string, val string) {
	out, val, ok := ShiftString(a)
	if !ok && len(defaultVal) != 0 {
		val = defaultVal[0]
	}
	return
}

func ShiftInt(a []int) (out []int, val int, ok bool) {
	if len(a) != 0 {
		return a[1:], a[0], true
	}
	out = a
	return
}

func MustShiftInt(a []int, defaultVal ...int) (out []int, val int) {
	out, val, ok := ShiftInt(a)
	if !ok && len(defaultVal) != 0 {
		val = defaultVal[0]
	}
	return
}

func ShiftInt64(a []int64) (out []int64, val int64, ok bool) {
	if len(a) != 0 {
		return a[1:], a[0], true
	}
	out = a
	return
}

func MustShiftInt64(a []int64, defaultVal ...int64) (out []int64, val int64) {
	out, val, ok := ShiftInt64(a)
	if !ok && len(defaultVal) != 0 {
		val = defaultVal[0]
	}
	return
}

func ShiftFloat(a []float64) (out []float64, val float64, ok bool) {
	if len(a) != 0 {
		return a[1:], a[0], true
	}
	out = a
	return
}

func MustShiftFloat(a []float64, defaultVal ...float64) (out []float64, val float64) {
	out, val, ok := ShiftFloat(a)
	if !ok && len(defaultVal) != 0 {
		val = defaultVal[0]
	}
	return
}

func ShiftBool(a []bool) (out []bool, val bool, ok bool) {
	if len(a) != 0 {
		return a[1:], a[0], true
	}
	out = a
	return
}

func MustShiftBool(a []bool, defaultVal ...bool) (out []bool, val bool) {
	out, val, ok := ShiftBool(a)
	if !ok && len(defaultVal) != 0 {
		val = defaultVal[0]
	}
	return
}
