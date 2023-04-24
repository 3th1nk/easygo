package internal

import (
	"github.com/axgle/mahonia"
)

func Translate(s, srcEncode, dstEncode string) string {
	srcDecoder := mahonia.NewDecoder(srcEncode)
	dstDecoder := mahonia.NewDecoder(dstEncode)
	src := srcDecoder.ConvertString(s)
	_, dst, _ := dstDecoder.Translate([]byte(src), true)
	return string(dst)
}
