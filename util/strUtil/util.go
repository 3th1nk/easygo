package strUtil

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"github.com/modern-go/reflect2"
	"net/url"
	"strings"
)

// ------------------------------------------------------------------------------ encode & decode
func UrlDecode(s string) string {
	v, _ := url.QueryUnescape(s)
	return v
}

func Sha1(s string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(s)))
}

func Md5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// ------------------------------------------------------------------------------
func Replace(str string, replace map[string]string) string {
	for key, val := range replace {
		str = strings.Replace(str, key, val, -1)
	}
	return str
}

func For(in []string, f func(int)) {
	if reflect2.IsNil(in) || f == nil {
		return
	}
	for i, n := 0, len(in); i < n; i++ {
		f(i)
	}
}

func Map(in []string, f func(i int, s string) string) []string {
	if reflect2.IsNil(in) || f == nil {
		return in
	}

	out := make([]string, len(in))
	for i, s := range in {
		out[i] = f(i, s)
	}
	return out
}
