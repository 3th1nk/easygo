package charset

import (
	"github.com/3th1nk/easygo/util/logs"
	"golang.org/x/text/encoding/simplifiedchinese"
	"strings"
)

/*
	UTF-8编码规则：如果只有一个字节则取值为0x00-0x7F，其余字节按长度进行以下拓展：

	UTF-8由4种编码方式实现，即UTF8-1 / UTF8-2 / UTF8-3 / UTF8-4，编码范围如下
					第一字节		第二字节		第三字节	    第四字节
	UTF8-1			0x00-0x7F

	UTF8-2 			0xC2-0xDF  0x80-0xBF

	UTF8-3 			0xE0 	   0xA0-0xBF  	0x80-0xBF
					0xE1-0xEC  0x80-0xBF 	0x80-0xBF
					0xED 	   0x80-0x9F 	0x80-0xBF
					0xEE-0xEF  0x80-0xBF 	0x80-0xBF

	UTF8-4 			0xF0 	   0x90-0xBF 	0x80-0xBF 	0x80-0xBF
					0xF1-0xF3  0x80-0xBF 	0x80-0xBF 	0x80-0xBF
					0xF4 	   0x80-0x8F 	0x80-0xBF 	0x80-0xBF

	注：每种编码可能有多个编码范围，每个编码范围间，以空格作为每个字节的分隔符。
	例如UTF8-3的第一个编码，其第一个字节取值必须为0xE0，第二个字节范围为0xA0-0xBF，第三个字节为0x80-0xBF。
*/

func utf8TailError(ch byte) bool {
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
			if (c0 == 0xE0 && (c1 < 0xA0 || c1 > 0xBF || utf8TailError(c2))) ||
				(0xE1 <= c0 && c0 <= 0xEC && (utf8TailError(c1) || utf8TailError(c2))) ||
				(c0 == 0xED && (c1 < 0x80 || c1 > 0x9F || utf8TailError(c2))) ||
				(0xEE <= c0 && c0 < 0xFF && (utf8TailError(c1) || utf8TailError(c2))) {
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
			if (c0 == 0xF0 && (c1 < 0x90 || 0xBF < c1 || utf8TailError(c2) || utf8TailError(c3))) ||
				(0xF1 <= c0 && c0 <= 0xF3 && (utf8TailError(c1) || utf8TailError(c2) || utf8TailError(c3))) ||
				(c0 == 0xF4 && (c1 < 0x80 || 0x8F < c1 || utf8TailError(c2) || utf8TailError(c3))) {
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
