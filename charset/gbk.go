package charset

import (
	"github.com/3th1nk/easygo/util/logs"
	"golang.org/x/text/encoding/simplifiedchinese"
)

func IsGB2312(s string) bool {
	for i := 0; i < len(s); {
		c0 := s[i]
		switch {
		default:
			return false

		case c0 < 0x80:
			i++

		case (0xA1 <= c0 && c0 <= 0xA9) || (0xB0 <= c0 && c0 <= 0xF7):
			// 	双字节
			//	0xA1-0xA9  0xA1-0xFE
			//	0xB0-0xF7  0xA1-0xFE
			if i+1 >= len(s) {
				// 数据不完整
				return false
			}

			c1 := s[i+1]
			if c1 < 0xA1 || 0xFE < c1 {
				return false
			}
			i += 2
		}
	}
	return true
}

func IsGBK(s string) bool {
	for i := 0; i < len(s); {
		c0 := s[i]
		switch {
		default:
			return false

		case c0 <= 0x80:
			// Microsoft's Code Page 936 extends GBK 1.0 to encode the euro sign U+20AC as 0x80
			i++
			continue

		case c0 < 0xFF:
			// 双字节 0x81–0xFE 0x40–0xFE(0x7F除外)
			if i+1 >= len(s) {
				// 数据不完整
				return false
			}

			c1 := s[i+1]
			if c1 < 0x40 || c1 == 0x7F || 0xFE < c1 {
				return false
			}
			i += 2
		}
	}
	return true
}

func IsGB18030(s string) bool {
	i := 0
loop:
	for i < len(s) {
		c0 := s[i]
		switch {
		default:
			return false

		case c0 <= 0x80:
			// Microsoft's Code Page 936 extends GBK 1.0 to encode the euro sign U+20AC as 0x80
			i++
			continue

		case c0 < 0xFF:
			if i+1 >= len(s) {
				// 数据不完整
				return false
			}

			c1 := s[i+1]
			switch {
			default:
				return false

			case 0x40 <= c1 && c1 != 0x7F && c1 < 0xFF:
				// 双字节 0x81–0xFE 0x40–0xFE(0x7F除外)
				i += 2
				goto loop

			case 0x30 <= c1 && c1 < 0x40:
				// 四字节 0x81-0xFE  0x30-0x39  0x81-0xFE  0x30-0x39
				if i+3 >= len(s) {
					// 数据不完整
					return false
				}

				c2 := s[i+2]
				if c2 < 0x81 || 0xFF <= c2 {
					return false
				}

				c3 := s[i+3]
				if c3 < 0x30 || 0x3A <= c3 {
					return false
				}

				i += 4
				goto loop
			}
		}
	}
	return true
}

func ToGBK(s string) string {
	if IsUTF8(s) {
		s1, err := simplifiedchinese.GBK.NewEncoder().String(s)
		if err == nil {
			return s1
		}
		logs.Default.Error(err.Error())
	}
	return s
}

func ToGB18030(s string) string {
	if IsUTF8(s) {
		s1, err := simplifiedchinese.GB18030.NewEncoder().String(s)
		if err == nil {
			return s1
		}
		logs.Default.Error(err.Error())
	}
	return s
}
