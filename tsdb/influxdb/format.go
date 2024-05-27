package influxdb

import (
	"strings"
	"unicode"
)

// ToLowerAndEscape 将字符串转为小写，并将特殊字符(下划线、连接符、点号除外)替换为下划线，可用于格式化measurement名称
func ToLowerAndEscape(s string) string {
	s = strings.ToLower(s)
	s = strings.Replace(s, "-", "_", -1)
	s = strings.Replace(s, "/", "_", -1)
	s = strings.Replace(s, " ", "_", -1)
	return s
}

// EscapeTagValue 将tag值中不合规的字符进行转义
func EscapeTagValue(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Replace(s, `\`, `\\`, -1)
	s = strings.Replace(s, "\n", "\\n", -1)
	for _, k := range []string{" ", ",", "="} {
		s = strings.Replace(s, k, `\\`+k, -1)
	}
	if strings.HasSuffix(s, `\`) {
		// influxdb不允许value以斜杠结尾，加一个空格
		s = s + ` `
	}
	return s
}

// EscapeFieldValue 将field值中不合规的字符进行转义
func EscapeFieldValue(value string) string {
	value = strings.TrimSpace(value)
	value = strings.Replace(value, "\n", "", -1)
	value = strings.Replace(value, `"`, ``, -1)
	if strings.HasSuffix(value, `\`) {
		// influxdb不允许value以斜杠结尾，加一个空格
		value = value + `\ `
	}
	return value
}

// Quote 给字符串加上双引号
func Quote(s string) string {
	return `"` + s + `"`
}

// QuoteIfNeed 如有必要，给字符串加上双引号，否则返回原字符串
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

// SingleQuote 给字符串加上单引号
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
