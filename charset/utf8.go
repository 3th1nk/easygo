package charset

import (
	"github.com/3th1nk/easygo/util/logs"
	"golang.org/x/text/encoding/simplifiedchinese"
	"strings"
	"unicode/utf8"
)

func IsUTF8(s string) bool {
	return utf8.ValidString(s)
}

// IsUTF8NotStrict 兼容非标准的UTF-8实现
func IsUTF8NotStrict(s string) bool {
	// 跳过 字节序标记
	if len(s) > 2 && s[0] == 0xEF && s[1] == 0xBB && s[2] == 0xBF {
		s = s[3:]
	}

	// 替换 非标准空字符(0xC0 0x80)
	nullBytes, stdNullBytes := []byte{0xC0, 0x80}, []byte{0x00}
	s = strings.ReplaceAll(s, string(nullBytes), string(stdNullBytes))

	return IsUTF8(s)
}

func ToUTF8(s string) string {
	if IsUTF8(s) {
		return s
	}
	if IsGB18030(s) {
		s1, err := simplifiedchinese.GB18030.NewDecoder().String(s)
		if err == nil {
			return s1
		}
		logs.Default.Error(err.Error())
	}
	return s
}
