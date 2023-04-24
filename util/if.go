package util

func If(condition bool, valIfTrue, valIfFalse interface{}) interface{} {
	if condition {
		return valIfTrue
	} else {
		return valIfFalse
	}
}

func IfUser(condition bool, valIfTrue, valIfFalse func() interface{}) interface{} {
	if condition {
		return valIfTrue()
	} else {
		return valIfFalse()
	}
}

func IfEmptySlice(val []interface{}, valIfEmpty interface{}) interface{} {
	if len(val) != 0 {
		return val[0]
	}
	return valIfEmpty
}

func IfString(condition bool, valIfTrue, valIfFalse string) string {
	if condition {
		return valIfTrue
	} else {
		return valIfFalse
	}
}

func IfEmptyString(val, valIfEmpty string) string {
	if val != "" {
		return val
	}
	return valIfEmpty
}

func IfEmptyStringSlice(val []string, valIfEmpty string) string {
	if len(val) != 0 {
		return val[0]
	}
	return valIfEmpty
}

func IfNotEmptyString(val, valIfNotEmpty string) string {
	if val != "" {
		return valIfNotEmpty
	}
	return ""
}

func IfUserString(condition bool, valIfTrue, valIfFalse func() string) string {
	if condition {
		return valIfTrue()
	} else {
		return valIfFalse()
	}
}

func IfUserEmptyString(val string, valIfEmpty func() string) string {
	if val == "" {
		return valIfEmpty()
	}
	return val
}

func IfInt(condition bool, valIfTrue, valIfFalse int) int {
	if condition {
		return valIfTrue
	} else {
		return valIfFalse
	}
}

func IfEmptyInt(val, valIfEmpty int) int {
	if val != 0 {
		return val
	}
	return valIfEmpty
}

func IfEmptyIntSlice(val []int, valIfEmpty int) int {
	if len(val) != 0 {
		return val[0]
	}
	return valIfEmpty
}

func IfNotEmptyInt(val, valIfNotEmpty int) int {
	if val != 0 {
		return valIfNotEmpty
	}
	return 0
}

func IfUserInt(condition bool, valIfTrue, valIfFalse func() int) int {
	if condition {
		return valIfTrue()
	} else {
		return valIfFalse()
	}
}

func IfUserEmptyInt(val int, valIfEmpty func() int) int {
	if val == 0 {
		return valIfEmpty()
	}
	return val
}

func IfInt64(condition bool, valIfTrue, valIfFalse int64) int64 {
	if condition {
		return valIfTrue
	} else {
		return valIfFalse
	}
}

func IfEmptyInt64(val, valIfEmpty int64) int64 {
	if val == 0 {
		return valIfEmpty
	}
	return val
}

func IfEmptyInt64Slice(val []int64, valIfEmpty int64) int64 {
	if len(val) != 0 {
		return val[0]
	}
	return valIfEmpty
}

func IfNotEmptyInt64(val, valIfNotEmpty int64) int64 {
	if val != 0 {
		return valIfNotEmpty
	}
	return 0
}

func IfUserInt64(condition bool, valIfTrue, valIfFalse func() int64) int64 {
	if condition {
		return valIfTrue()
	} else {
		return valIfFalse()
	}
}

func IfUserEmptyInt64(val int64, valIfEmpty func() int64) int64 {
	if val == 0 {
		return valIfEmpty()
	}
	return val
}

func IfInt32(condition bool, valIfTrue, valIfFalse int32) int32 {
	if condition {
		return valIfTrue
	} else {
		return valIfFalse
	}
}

func IfEmptyInt32(val, valIfEmpty int32) int32 {
	if val == 0 {
		return valIfEmpty
	}
	return val
}

func IfEmptyInt32Slice(val []int32, valIfEmpty int32) int32 {
	if len(val) != 0 {
		return val[0]
	}
	return valIfEmpty
}

func IfNotEmptyInt32(val, valIfNotEmpty int32) int32 {
	if val != 0 {
		return valIfNotEmpty
	}
	return 0
}

func IfUserInt32(condition bool, valIfTrue, valIfFalse func() int32) int32 {
	if condition {
		return valIfTrue()
	} else {
		return valIfFalse()
	}
}

func IfUserEmptyInt32(val int32, valIfEmpty func() int32) int32 {
	if val == 0 {
		return valIfEmpty()
	}
	return val
}

func IfFloat(condition bool, valIfTrue, valIfFalse float64) float64 {
	if condition {
		return valIfTrue
	} else {
		return valIfFalse
	}
}

func IfUserFloat(condition bool, valIfTrue, valIfFalse func() float64) float64 {
	if condition {
		return valIfTrue()
	} else {
		return valIfFalse()
	}
}

func IfEmptyBoolSlice(val []bool, valIfEmpty bool) bool {
	if len(val) != 0 {
		return val[0]
	}
	return valIfEmpty
}
