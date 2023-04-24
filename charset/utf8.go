package charset

import "github.com/3th1nk/easygo/charset/internal"

/*
	UTF-8编码规则：如果只有一个字节则取值为0x00-0x7F，其余字节按长度进行以下拓展：

	UTF-8由4种编码方式实现，即UTF8-1 / UTF8-2 / UTF8-3 / UTF8-4，编码范围如下
					第一字节		第二字节		第三字节	第四字节
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

func IsUTF8(s string) bool {
	var nBytes int
	for _, chr := range s {
		// 跳过ASCII字符
		if nBytes == 0 && chr < 0x80 {
			continue
		}

		// 多字节字符, 计算字节数
		if nBytes == 0 {
			if chr >= 0xFC && chr <= 0xFD {
				nBytes = 6
			} else if chr >= 0xF8 {
				nBytes = 5
			} else if chr >= 0xF0 {
				nBytes = 4
			} else if chr >= 0xE0 {
				nBytes = 3
			} else if chr >= 0xC0 {
				nBytes = 2
			} else {
				return false
			}
		} else {
			if (chr & 0xC0) != 0x80 {
				return false
			}
		}
		nBytes--
	}

	if nBytes != 0 {
		return false
	}
	return true
}

func ToUTF8(s string) string {
	if IsGB18030(s) {
		return internal.Translate(s, "GB18030", "UTF-8")
	}
	return s
}
