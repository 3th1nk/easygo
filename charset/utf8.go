package charset

import (
	"github.com/3th1nk/easygo/util/logs"
	"golang.org/x/text/encoding/simplifiedchinese"
	"strings"
)

func notUTF8Tail(ch byte) bool {
	return ch < 0x80 || ch > 0xBF
}

func IsUTF8(s string) bool {
	for i := 0; i < len(s); {
		c0 := s[i]
		switch {
		default:
			return false

		case c0 < 0x80:
			i++
			continue

		case 0xC2 <= c0 && c0 < 0xDF:
			if i+1 >= len(s) {
				// 数据不完整
				return false
			}

			// UTF8-2
			// 0xC2-0xDF  0x80-0xBF
			c1 := s[i+1]
			if c1 < 0x80 || 0xBF < c1 {
				return false
			}
			i += 2
			continue

		case 0xE0 <= c0 && c0 < 0xFF:
			if i+2 >= len(s) {
				// 数据不完整
				return false
			}

			// UTF8-3
			//	0xE0 	   0xA0-0xBF  	0x80-0xBF
			//	0xE1-0xEC  0x80-0xBF 	0x80-0xBF
			//	0xED 	   0x80-0x9F 	0x80-0xBF
			//	0xEE-0xEF  0x80-0xBF 	0x80-0xBF
			c1, c2 := s[i+1], s[i+2]
			if (c0 == 0xE0 && (c1 < 0xA0 || c1 > 0xBF || notUTF8Tail(c2))) ||
				(0xE1 <= c0 && c0 <= 0xEC && (notUTF8Tail(c1) || notUTF8Tail(c2))) ||
				(c0 == 0xED && (c1 < 0x80 || c1 > 0x9F || notUTF8Tail(c2))) ||
				(0xEE <= c0 && c0 < 0xFF && (notUTF8Tail(c1) || notUTF8Tail(c2))) {
				return false
			}
			i += 3
			continue

		case 0xF0 <= c0 && c0 < 0xF5:
			if i+3 >= len(s) {
				// 数据不完整
				return false
			}

			//	UTF8-4
			//	0xF0 	   0x90-0xBF 	0x80-0xBF 	0x80-0xBF
			//	0xF1-0xF3  0x80-0xBF 	0x80-0xBF 	0x80-0xBF
			//	0xF4 	   0x80-0x8F 	0x80-0xBF 	0x80-0xBF
			c1, c2, c3 := s[i+1], s[i+2], s[i+3]
			if (c0 == 0xF0 && (c1 < 0x90 || 0xBF < c1 || notUTF8Tail(c2) || notUTF8Tail(c3))) ||
				(0xF1 <= c0 && c0 <= 0xF3 && (notUTF8Tail(c1) || notUTF8Tail(c2) || notUTF8Tail(c3))) ||
				(c0 == 0xF4 && (c1 < 0x80 || 0x8F < c1 || notUTF8Tail(c2) || notUTF8Tail(c3))) {
				return false
			}
			i += 4
			continue
		}
	}
	return true
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
	if IsGB18030(s) {
		s1, err := simplifiedchinese.GB18030.NewDecoder().String(s)
		if err == nil {
			return s1
		}
		logs.Default.Error(err.Error())
	}
	return s
}
