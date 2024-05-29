package influxdb

import (
	"strings"
	"unicode"
)

// FormatMeasurement 统一measurement命名规范，避免特殊字符
func FormatMeasurement(s string) string {
	return strings.Map(func(r rune) rune {
		// 统一小写
		r = unicode.ToLower(r)
		// 替换 A-Z a-z 0-9 _-. 之外的字符
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' || r == '.' {
			return r
		}
		return '_'
	}, s)
}

// EscapeTagValue 写入数据时，需要将tag字段值中不合规的字符进行转义
//	行协议：https://docs.influxdata.com/influxdb/v1/write_protocols/line_protocol_reference
func EscapeTagValue(s string) string {
	// 统一去除前后空格
	s = strings.TrimSpace(s)
	// 尾部反斜杠需要特殊处理
	s = EscapeTagTailBackslash(s)

	// 空格、逗号、等号需要转义
	return strings.NewReplacer(
		" ", "\\ ",
		",", "\\,",
		"=", "\\=",
	).Replace(s)
}

// EscapeFieldValue 写入数据时，需要将field值中不合规的字符进行转义
func EscapeFieldValue(s string) string {
	// 统一去除前后空格
	s = strings.TrimSpace(s)

	// 单引号、反斜杠需要转义
	return strings.NewReplacer(
		`"`, `\"`,
		`\`, `\\`,
	).Replace(s)
}

// EscapeCondValue 查询数据时，需要将条件值中部分字符进行转义
//  由于tag和field写入时转义稍有差别，所以查询条件也需要区分，通常都是用tag字段作为条件，如果是field字段，需要额外指定 isFieldVal=true
//	https://docs.influxdata.com/influxdb/v1/query_language/data_exploration/#string-literals
func EscapeCondValue(s string, isFieldVal ...bool) string {
	s = strings.TrimSpace(s)
	// field字段尾部反斜杠写入时不需要特殊处理，所以查询条件中也不需要处理
	if len(isFieldVal) == 0 || !isFieldVal[0] {
		s = EscapeTagTailBackslash(s)
	}

	// 单引号、反斜杠需要转义
	return strings.NewReplacer(
		`'`, `\'`,
		`\`, `\\`,
	).Replace(s)
}

// EscapeTagTailBackslash 写入数据时，标签字段尾部的反斜杠会对后面的空格进行转义，导致行协议无法正确解析，所以需要额外处理，查询时也需要对齐
func EscapeTagTailBackslash(s string) string {
	if strings.HasSuffix(s, `\`) {
		return s + ` `
	}
	return s
}

// UnescapeQueryResultValue 还原查询结果中tag、field值的转义
//	写入时(RawWrite除外)tag、field值转义时都移除了前后空格，唯一特殊情况是以反斜杠结尾的tag字段，特殊追加了一个空格，所以查询结果中需要去掉，以保持一致
func UnescapeQueryResultValue(s string) string {
	if strings.HasSuffix(s, `\ `) {
		s = s[:len(s)-1]
	}
	return s
}

// Quote 给字符串加上双引号
func Quote(s string) string {
	return `"` + s + `"`
}

// QuoteIfNeed 如有必要，给字段的Key加上双引号，否则返回原字符串
//	https://docs.influxdata.com/influxdb/v1/query_language/explore-data/#quoting
func QuoteIfNeed(s string) string {
	if len(s) == 0 {
		return s
	}
	// 数字开头
	if s[0] >= '0' && s[0] <= '9' {
		return Quote(s)
	}
	// 关键字
	if isKeyword(s) {
		return Quote(s)
	}
	// 包含 A-z,0-9,_ 之外的字符
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return Quote(s)
		}
	}
	return s
}

// SingleQuote 给查询条件值加上单引号
//	https://docs.influxdata.com/influxdb/v1/troubleshooting/frequently-asked-questions/#when-should-i-single-quote-and-when-should-i-double-quote-in-queries
func SingleQuote(s string) string {
	return `'` + s + `'`
}

// trimParentheses 移除左右的小括号
func trimParentheses(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")") {
		s = s[1 : len(s)-1]
	}
	return s
}
