package jsonUtil

import (
	"fmt"
	"github.com/3th1nk/easygo/util/convertor"
	"github.com/3th1nk/easygo/util/timeUtil"
	jsonIter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
	"math"
	"time"
	"unsafe"
)

// 注册弱类型解析
func registerFuzzy() {
	// jsonIter 自带的弱类型解析，支持将 string("123") 解析为 number(123) 等。
	extra.RegisterFuzzyDecoders()

	// 注册 string 类型的弱类型解析，支持将 array、object 解析为字符串
	jsonIter.RegisterTypeDecoderFunc("string", func(ptr unsafe.Pointer, iter *jsonIter.Iterator) {
		val, err := convertor.ToString(iter.Read())
		if err == nil {
			*((*string)(ptr)) = val
		} else {
			iter.ReportError("fuzzyStringDecoder", fmt.Sprintf("convert error: %v", err))
		}
	})
	// 注册 bool 类型的弱类型解析，支持将 number(!=0|==0)、string("true|false") 解析为 bool(true|false)
	jsonIter.RegisterTypeDecoderFunc("bool", func(ptr unsafe.Pointer, iter *jsonIter.Iterator) {
		val, err := convertor.ToBool(iter.Read())
		if err == nil {
			*((*bool)(ptr)) = val
		} else {
			iter.ReportError("fuzzyBoolDecoder", fmt.Sprintf("convert error: %v", err))
		}
	})

	// 默认将 time.Time 解析为 yyyy-MM-dd HH:MM:SS 格式
	jsonIter.RegisterTypeEncoderFunc("time.Time", func(ptr unsafe.Pointer, stream *jsonIter.Stream) {
		stream.WriteString((*time.Time)(ptr).Format(DefaultTimeFormat))
	}, func(ptr unsafe.Pointer) bool {
		return (*(*time.Time)(ptr)).IsZero()
	})
	jsonIter.RegisterTypeDecoderFunc("time.Time", func(ptr unsafe.Pointer, iter *jsonIter.Iterator) {
		valueType := iter.WhatIsNext()
		switch valueType {
		case jsonIter.NilValue:
			*((*time.Time)(ptr)) = time.Time{}
		case jsonIter.NumberValue:
			f := iter.ReadFloat64()
			if f > math.MaxInt32 {
				f = f / 1000
			}
			s := int64(f)
			ns := int64(f-float64(s)) * int64(time.Second)
			*((*time.Time)(ptr)) = time.Unix(s, ns)
		case jsonIter.StringValue:
			s := iter.ReadString()
			if s == "" {
				*((*time.Time)(ptr)) = time.Time{}
			} else {
				var err error
				*((*time.Time)(ptr)), err = timeUtil.Parse(s)
				if err != nil {
					iter.ReportError("fuzzyTimeDecoder", err.Error())
				}
			}
		default:
			iter.ReportError("fuzzyTimeDecoder", "not time.Time or string")
		}
	})
}
