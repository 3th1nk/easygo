package jsonUtil

import (
	"fmt"
	"github.com/3th1nk/easygo/crypto/aes"
	jsonIter "github.com/json-iterator/go"
	"unsafe"
)

// 注册加密字段。
//   加密字段会使用 AES 加密算法在 JSON 序列化和反序列化时自动加解密。
func RegisterEncryptField(typ string, encryptKey string, encryptField ...string) {
	for _, s := range encryptField {
		jsonIter.RegisterFieldEncoderFunc(typ, s, func(ptr unsafe.Pointer, stream *jsonIter.Stream) {
			str, err := aes.Encrypt(*(*string)(ptr), encryptKey, true)
			if err != nil {
				panic(err)
			}
			stream.WriteString(str)
		}, func(ptr unsafe.Pointer) bool {
			return *(*string)(ptr) == ""
		})
		jsonIter.RegisterFieldDecoderFunc(typ, s, func(ptr unsafe.Pointer, iter *jsonIter.Iterator) {
			any := iter.ReadAny()
			switch any.ValueType() {
			case jsonIter.InvalidValue, jsonIter.NilValue:
				*((*string)(ptr)) = ""
			case jsonIter.StringValue:
				*(*string)(ptr), _ = aes.Decrypt(any.ToString(), encryptKey)
			default:
				iter.ReportError("", fmt.Sprintf("cannot convert '%v' to CommaSepString", any.GetInterface()))
			}
		})
	}
}
