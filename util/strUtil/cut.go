package strUtil

import "unicode/utf8"

// CutHeadIfOverflow 当字符串字节长度超过maxLen时从头开始按字符集截取字符直到等于或小于
// example：
// str := "12你好"
// CutHeadIfOverflow(str,6)>>> 你好
// CutHeadIfOverflow(str,5)>>> 好
func CutHeadIfOverflow(str string, maxLen uint64) string {
	// 为什么不用递归? golang貌似没有对尾递归进行优化。
	strRune := []rune(str)
	for uint64(utf8.RuneCountInString(str)) > maxLen {
		strRune = strRune[1:]
		str = string(strRune)
	}
	return str
}

// CutTailIfOverflow 当字符串字节长度超过maxLen时从尾部开始按字符集截取字符直到等于或小于
// example：
// str := "12你好"
// CutHeadIfOverflow(str,6)>>> 12你
// CutHeadIfOverflow(str,4)>>> 12
func CutTailIfOverflow(str string, maxLen uint64) string {
	strRune := []rune(str)
	for uint64(utf8.RuneCountInString(str)) > maxLen {
		strRune = strRune[:len(strRune)-1]
		str = string(strRune)
	}
	return str
}
