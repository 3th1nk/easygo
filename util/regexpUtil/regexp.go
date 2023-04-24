package regexpUtil

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	PatternUint = `\d+`
	PatternInt  = `(?:-)?\d+`

	patternIpv4     = `\d{1,3}.\d{1,3}.\d{1,3}.\d{1,3}`
	patternIpv4Addr = `\d{1,3}.\d{1,3}.\d{1,3}.\d{1,3}(?::\d+)?`
)

var (
	regexIntFind   = regexp.MustCompile(PatternInt)
	regexIntMatch  = regexp.MustCompile(`^` + PatternInt + `$`)
	regexUintFind  = regexp.MustCompile(PatternUint)
	regexUintMatch = regexp.MustCompile(`^` + PatternUint + `$`)

	regexIpv4Find     = regexp.MustCompile(patternIpv4)
	regexIpv4Match    = regexp.MustCompile(`^` + patternIpv4 + `$`)
	regexIpv4AddrFind = regexp.MustCompile(patternIpv4Addr)
)

func IsInt(s string) bool {
	return regexIntMatch.MatchString(strings.TrimSpace(s))
}

func FindInt(s string) []string {
	return regexIntFind.FindAllString(s, -1)
}

func IsUint(s string) bool {
	return regexUintMatch.MatchString(strings.TrimSpace(s))
}

func FindUint(s string) []string {
	return regexUintFind.FindAllString(s, -1)
}

func IsFloat(s string) bool {
	if s = strings.TrimSpace(s); s == "" {
		return false
	}

	hasDot := false
	for i, c := range s {
		if c == '-' {
			if i != 0 {
				return false
			}
		} else if c == '.' {
			if hasDot {
				return false
			} else {
				hasDot = true
			}
		} else if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func IsIPv4(s string) bool {
	s = strings.TrimSpace(s)
	return regexIpv4Match.MatchString(s) && verifyIPv4Str(s)
}

func FindIPv4(s string) []string {
	matches := regexIpv4Find.FindAllStringSubmatch(s, -1)
	arr := make([]string, 0, len(matches))
	for _, v := range matches {
		if verifyIPv4Str(v[0]) {
			arr = append(arr, v[0])
		}
	}
	return arr
}

func FindIPv4Addr(s string) []string {
	matches := regexIpv4AddrFind.FindAllString(s, -1)
	arr := make([]string, 0, len(matches))
	for _, str := range matches {
		if pos := strings.Index(str, ":"); pos != -1 {
			if verifyIPv4Str(str[:pos]) {
				arr = append(arr, str)
			}
		} else {
			if verifyIPv4Str(str) {
				arr = append(arr, str)
			}
		}
	}
	return arr
}

func FindFirstIPv4Addr(s string) string {
	if arr := FindIPv4Addr(s); len(arr) != 0 {
		return arr[0]
	}
	return ""
}

func verifyIPv4Str(s string) bool {
	arr := strings.Split(s, ".")
	if len(arr) != 4 {
		return false
	}
	for _, s := range arr {
		if !IsUint(s) {
			return false
		}
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return false
		} else if n > 255 {
			return false
		}
	}
	return true
}
