package strUtil

import (
	"fmt"
	"github.com/3th1nk/easygo/util/mathUtil"
	"github.com/modern-go/reflect2"
	"regexp"
	"strings"
)

func Split(s, sep string, ignoreEmpty bool, format ...func(string) string) []string {
	return SplitN(s, sep, ignoreEmpty, -1, format...)
}

func SplitN(s, sep string, ignoreEmpty bool, limit int, format ...func(string) string) []string {
	if limit == 0 {
		return nil
	}

	var theF func(string) string
	if n := len(format); n != 0 && format[0] != nil {
		theF = format[0]
	} else if ignoreEmpty && n == 0 {
		theF = strings.TrimSpace
	}

	result, sepLen, maxIdx, count := make([]string, 0, mathUtil.MaxInt(limit, 4)), len(sep), limit-1, 0
	for s != "" {
		var str string
		if count == maxIdx {
			str, s = s, ""
		} else if pos := strings.Index(s, sep); pos == -1 {
			str, s = s, ""
		} else {
			str, s = s[:pos], s[pos+sepLen:]
		}
		if theF != nil {
			str = theF(str)
		}
		if str != "" || !ignoreEmpty {
			result = append(result, str)
			count++
		}
	}
	return result
}

func Split2(s, sep string, ignoreEmpty bool, f ...func(string) string) (a string, b string) {
	arr := SplitN(s, sep, ignoreEmpty, 2, f...)
	switch len(arr) {
	case 1:
		a = arr[0]
	case 2:
		a, b = arr[0], arr[1]
	}
	return
}

func Split3(s, sep string, ignoreEmpty bool, f ...func(string) string) (a string, b string, c string) {
	arr := SplitN(s, sep, ignoreEmpty, 3, f...)
	switch len(arr) {
	case 1:
		a = arr[0]
	case 2:
		a, b = arr[0], arr[1]
	case 3:
		a, b, c = arr[0], arr[1], arr[2]
	}
	return
}

func SplitToInt64(s, sep string, ignoreEmpty ...bool) ([]int64, error) {
	return ToInt64Arr(Split(s, sep, append(ignoreEmpty, true)[0]))
}

func SplitToInt64NoError(s, sep string, ignoreEmpty ...bool) []int64 {
	v, _ := ToInt64Arr(Split(s, sep, append(ignoreEmpty, true)[0]))
	return v
}

func SplitToInt(s, sep string, ignoreEmpty ...bool) ([]int, error) {
	return ToIntArr(Split(s, sep, append(ignoreEmpty, true)[0]))
}

func SplitToIntNoError(s, sep string, ignoreEmpty ...bool) []int {
	v, _ := ToIntArr(Split(s, sep, append(ignoreEmpty, true)[0]))
	return v
}

func PregSplitToInt64(s, pattern string) ([]int64, error) {
	reg, err := regexp.Compile(pattern)
	if !reflect2.IsNil(err) {
		return nil, fmt.Errorf("正则表达式语法不正确: %v, pattern=%s", err, pattern)
	}
	return ToInt64Arr(reg.Split(s, -1))
}

func PregSplitToInt64NoError(s, pattern string) []int64 {
	v, _ := PregSplitToInt64(s, pattern)
	return v
}

func PregSplitToInt(s, pattern string) ([]int, error) {
	reg, err := regexp.Compile(pattern)
	if !reflect2.IsNil(err) {
		return nil, fmt.Errorf("正则表达式语法不正确: %v, pattern=%s", err, pattern)
	}
	return ToIntArr(reg.Split(s, -1))
}

func PregSplitToIntNoError(s, pattern string) []int {
	v, _ := PregSplitToInt(s, pattern)
	return v
}
