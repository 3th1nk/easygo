package _gen

import (
	"github.com/3th1nk/easygo/util"
	"github.com/3th1nk/easygo/util/strUtil"
	"testing"
)

func TestGenUnion(t *testing.T) {
	util.Println("\n\n\n\n")
	genUnion("int")
	genUnion("int32")
	genUnion("int64")
	util.Println("\n\n\n\n")
}

func genUnion(typ string) {
	util.Println(strUtil.Replace(`
// 取两个 []{type} 的并集
func Union{type2}(a ...[]{type}) []{type} {
	// 过滤空切片
	b := make([][]{type}, 0, len(a))
	// b 中所有切片的长度之和
	var bn int
	for _, v := range a {
		if n := len(v); n != 0 {
			b = append(b, v)
			bn += n
		}
	}

	if n := len(b); n == 0 {
		return nil
	} else if n == 1 {
		return b[0]
	}

	c := make([]{type}, len(b[0]), bn)
	copy(c, b[0])
	for _, arr := range b[1:] {
		for _, n := range arr {
			if -1 == IndexOf{type2}(c, n) {
				c = append(c, n)
			}
		}
	}
	return c
}
`, map[string]string{
		"{type}":  typ,
		"{type2}": strUtil.UcFirst(typ),
	}))
}
