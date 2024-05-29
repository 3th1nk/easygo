package influxdb

import "strings"

var (
	// keywordTable 关键字表
	keywordTable = map[string]interface{}{
		"ALL":           nil,
		"ALTER":         nil,
		"ANY":           nil,
		"AS":            nil,
		"ASC":           nil,
		"BEGIN":         nil,
		"BY":            nil,
		"CREATE":        nil,
		"CONTINUOUS":    nil,
		"DATABASE":      nil,
		"DATABASES":     nil,
		"DEFAULT":       nil,
		"DELETE":        nil,
		"DESC":          nil,
		"DESTINATIONS":  nil,
		"DIAGNOSTICS":   nil,
		"DISTINCT":      nil,
		"DROP":          nil,
		"DURATION":      nil,
		"END":           nil,
		"EVERY":         nil,
		"EXPLAIN":       nil,
		"FIELD":         nil,
		"FOR":           nil,
		"FROM":          nil,
		"GRANT":         nil,
		"GRANTS":        nil,
		"GROUP":         nil,
		"GROUPS":        nil,
		"IN":            nil,
		"INF":           nil,
		"INSERT":        nil,
		"INTO":          nil,
		"KEY":           nil,
		"KEYS":          nil,
		"KILL":          nil,
		"LIMIT":         nil,
		"SHOW":          nil,
		"MEASUREMENT":   nil,
		"MEASUREMENTS":  nil,
		"NAME":          nil,
		"OFFSET":        nil,
		"ON":            nil,
		"ORDER":         nil,
		"PASSWORD":      nil,
		"POLICY":        nil,
		"POLICIES":      nil,
		"PRIVILEGES":    nil,
		"QUERIES":       nil,
		"QUERY":         nil,
		"READ":          nil,
		"REPLICATION":   nil,
		"RESAMPLE":      nil,
		"RETENTION":     nil,
		"REVOKE":        nil,
		"SELECT":        nil,
		"SERIES":        nil,
		"SET":           nil,
		"SHARD":         nil,
		"SHARDS":        nil,
		"SLIMIT":        nil,
		"SOFFSET":       nil,
		"STATS":         nil,
		"SUBSCRIPTION":  nil,
		"SUBSCRIPTIONS": nil,
		"TAG":           nil,
		"TO":            nil,
		"USER":          nil,
		"USERS":         nil,
		"VALUES":        nil,
		"WHERE":         nil,
		"WITH":          nil,
		"WRITE":         nil,
	}

	// functionTable 函数表
	//	https://docs.influxdata.com/influxdb/v1/query_language/functions
	functionTable = map[string]interface{}{
		"COUNT":                             nil,
		"DISTINCT":                          nil,
		"INTEGRAL":                          nil,
		"MEAN":                              nil,
		"MEDIAN":                            nil,
		"MODE":                              nil,
		"SPREAD":                            nil,
		"STDDEV":                            nil,
		"SUM":                               nil,
		"BOTTOM":                            nil,
		"FIRST":                             nil,
		"LAST":                              nil,
		"MAX":                               nil,
		"MIN":                               nil,
		"PERCENTILE":                        nil,
		"SAMPLE":                            nil,
		"TOP":                               nil,
		"ABS":                               nil,
		"ACOS":                              nil,
		"ASIN":                              nil,
		"ATAN":                              nil,
		"ATAN2":                             nil,
		"CEIL":                              nil,
		"COS":                               nil,
		"CUMULATIVE_SUM":                    nil,
		"DERIVATIVE":                        nil,
		"DIFFERENCE":                        nil,
		"ELAPSED":                           nil,
		"EXP":                               nil,
		"FLOOR":                             nil,
		"HISTOGRAM":                         nil,
		"LN":                                nil,
		"LOG":                               nil,
		"LOG2":                              nil,
		"LOG10":                             nil,
		"MOVING_AVERAGE":                    nil,
		"NON_NEGATIVE_DERIVATIVE":           nil,
		"NON_NEGATIVE_DIFFERENCE":           nil,
		"POW":                               nil,
		"ROUND":                             nil,
		"SIN":                               nil,
		"SQRT":                              nil,
		"TAN":                               nil,
		"HOLT_WINTERS":                      nil,
		"CHANDE_MOMENTUM_OSCILLATOR":        nil,
		"EXPONENTIAL_MOVING_AVERAGE":        nil,
		"DOUBLE_EXPONENTIAL_MOVING_AVERAGE": nil,
		"KAUFMANS_EFFICIENCY_RATIO":         nil,
		"KAUFMANS_ADAPTIVE_MOVING_AVERAGE":  nil,
		"TRIPLE_EXPONENTIAL_MOVING_AVERAGE": nil,
		"TRIPLE_EXPONENTIAL_DERIVATIVE":     nil,
		"RELATIVE_STRENGTH_INDEX":           nil,
	}
)

// isKeyword 判断是否为关键字
func isKeyword(s string) bool {
	_, ok := keywordTable[strings.ToUpper(s)]
	return ok
}

// usingFunction 判断是否使用了函数
func usingFunction(s string) bool {
	if len(s) == 0 {
		return false
	}

	// 找到第一个左括号
	idx := strings.Index(s, "(")
	if idx < 2 { // 函数名至少2个字符
		return false
	}

	// 判断左括号前是否为函数名
	if _, ok := functionTable[strings.ToUpper(s[:idx])]; !ok {
		return false
	}

	// 需要有一个成对的右括号，但不一定是最后一个字符，因为有可能用了别名
	//	！！！这里只是一个简单判断！！！
	if !strings.Contains(s[idx+1:], ")") {
		return false
	}

	return true
}
