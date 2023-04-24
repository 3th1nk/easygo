package jsonUtil

import (
	jsonIter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"unsafe"
)

func NoOmitemptyApi() *Api {
	return defaultApi.NoOmitemptyApi()
}

// 返回一个忽略字段上 omitempty 标记的 Api
func (this *Api) NoOmitemptyApi() *Api {
	api := this.Clone()
	api.RegisterExtension(&noOmitemptyExtension{})
	return api
}

type noOmitemptyExtension struct {
	jsonIter.DummyExtension
}

func (extension *noOmitemptyExtension) DecorateEncoder(typ reflect2.Type, encoder jsonIter.ValEncoder) jsonIter.ValEncoder {
	return &noOmitemptyEncoder{encoder: encoder}
}

type noOmitemptyEncoder struct {
	encoder jsonIter.ValEncoder
}

func (this *noOmitemptyEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

func (this *noOmitemptyEncoder) Encode(ptr unsafe.Pointer, stream *jsonIter.Stream) {
	this.encoder.Encode(ptr, stream)
}
