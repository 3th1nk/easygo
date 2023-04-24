package strUtil

import (
	"regexp"
	"strings"
)

// 首字母大写
func UcFirst(s string) string {
	n := len(s)
	if n == 0 {
		return s
	} else if n == 1 {
		return strings.ToUpper(s)
	} else {
		c := s[:1]
		if u := strings.ToUpper(c); u == c {
			return s
		} else {
			return u + s[1:]
		}
	}
}

// 首字母小写
func LcFirst(s string) string {
	n := len(s)
	if n == 0 {
		return s
	} else if n == 1 {
		return strings.ToLower(s)
	} else {
		c := s[:1]
		if u := strings.ToLower(c); u == c {
			return s
		} else {
			return u + s[1:]
		}
	}
}

// 将驼峰命名法转换为首字母大写中划线分割格式（Http Header 格式）
func CamelToUpperKebab(s string) string {
	for _, ss := range camelToKebabPattern.FindAllString(s, -1) {
		s = strings.Replace(s, ss, ss[:1]+"-"+ss[1:], -1)
	}
	return UcFirst(s)
}

// 将驼峰命名法转换为小写下划线分割格式
func CamelToSnake(s string) string {
	return PascalToSnake(s)
}

// 将 Pascal 命名法转换为小写下划线分割格式
func PascalToSnake(s string) string {
	arr := caseSplit(s, pascalSplit, strings.ToLower)
	return strings.Join(arr, "_")
}

func pascalSplit(prevCase, currentCase, nextCase letterCase) (split bool) {
	if prevCase == upperCase {
		return currentCase == upperCase && nextCase == lowerCase
	} else if currentCase == upperCase {
		return true
	}
	return false
}

// 将驼峰命名法转换为小写中划线分割格式
func CamelToLowerKebab(s string) string {
	for _, ss := range camelToKebabPattern.FindAllString(s, -1) {
		s = strings.Replace(s, ss, ss[:1]+"-"+ss[1:], -1)
	}
	return strings.ToLower(s)
}

var camelToKebabPattern = regexp.MustCompile("[a-z][A-Z]")

// 将中划线分割格式转换为驼峰命名法
func KebabToCamel(s string) string {
	arr := strings.Split(s, "-")
	for i, s := range arr {
		if i == 0 {
			arr[i] = LcFirst(s)
		} else {
			arr[i] = UcFirst(s)
		}
	}
	return strings.Join(arr, "")
}

// 将中划线分割格式转换为下划线小写格式
func KebabToSnake(s string) string {
	return strings.ToLower(strings.Replace(s, "-", "_", -1))
}

// 下划线小写格式转换为驼峰命名法
func SnakeToCamel(s string) string {
	if s != "" {
		arr := strings.Split(s, "_")
		return strings.Join(Map(arr, func(i int, s string) string {
			if i == 0 {
				return LcFirst(s)
			} else {
				return UcFirst(s)
			}
		}), "")
	}
	return ""
}

// 下划线小写格式转换为 Pascal 格式
func SnakeToPascal(s string) string {
	if s != "" {
		arr := strings.Split(s, "_")
		return strings.Join(Map(arr, func(i int, s string) string {
			return UcFirst(s)
		}), "")
	}
	return ""
}

// 下划线小写格式转换为首字母大写中划线分割格式（Http Header 格式）
func SnakeToUpperKebab(s string) string {
	if s != "" {
		arr := strings.Split(s, "_")
		return strings.Join(Map(arr, func(i int, s string) string {
			return UcFirst(s)
		}), "-")
	}
	return ""
}

type letterCase int

const (
	otherCase = letterCase(0)
	upperCase = letterCase(1)
	lowerCase = letterCase(2)
)

func caseSplit(str string, split func(prevCase, currentCase, nextCase letterCase) (split bool), format func(s string) string) []string {
	if format == nil {
		format = func(s string) string { return s }
	}

	length := len(str)
	if length == 0 {
		return []string{}
	} else if length == 1 {
		return []string{format(str)}
	}

	getCase := func(r byte) letterCase {
		if r >= 'a' && r <= 'z' {
			return lowerCase
		} else if r >= 'A' && r <= 'Z' {
			return upperCase
		}
		return otherCase
	}

	result, lastSplitPos := make([]string, 0, 4), 0
	prevCase, currentCase := getCase(str[0]), getCase(str[1])
	for i, n := 1, length-1; i < n; i++ {
		nextCase := getCase(str[i+1])
		if split(prevCase, currentCase, nextCase) {
			result, lastSplitPos = append(result, format(str[lastSplitPos:i])), i
		}
		prevCase, currentCase = currentCase, nextCase
	}
	if split(prevCase, currentCase, otherCase) {
		result = append(result, format(str[lastSplitPos:length-1]), format(str[length-1:]))
	} else {
		result = append(result, format(str[lastSplitPos:]))
	}
	return result
}
