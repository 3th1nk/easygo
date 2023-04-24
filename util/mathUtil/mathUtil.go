package mathUtil

func IsZeroFloat6(f float64) bool {
	return f >= -0.000001 && f < 0.000001
}
