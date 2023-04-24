package arrUtil

import (
	"fmt"
	"github.com/3th1nk/easygo/util/strUtil"
	jsonIter "github.com/json-iterator/go"
	"reflect"
	"strings"
	"unsafe"
)

// CommaSepInt 以逗号分隔的整数
type CommaSepInt []int

func (this CommaSepInt) String() string {
	return strUtil.JoinInt(this, ",")
}

func (this *CommaSepInt) FromString(str string) error {
	if str == "" {
		*this = []int{}
	} else {
		*this, _ = strUtil.SplitToInt(str, ",", false)
	}
	return nil
}

func (this *CommaSepInt) FromDB(bytes []byte) error {
	return this.FromString(string(bytes))
}

func (this *CommaSepInt) ToDB() ([]byte, error) {
	return []byte(this.String()), nil
}

// 以逗号分隔的 string
//
// 注意：
//   每个字符串元素中不能再包含逗号，否则会导致序列化与反序列化结果不一致
type CommaSepString []string

func (this CommaSepString) String() string {
	return strings.Join(this, ",")
}

func (this *CommaSepString) FromString(str string) error {
	if str == "" {
		*this = []string{}
	} else {
		*this = strings.Split(str, ",")
	}
	return nil
}

func (this *CommaSepString) FromDB(bytes []byte) error {
	return this.FromString(string(bytes))
}

func (this *CommaSepString) ToDB() ([]byte, error) {
	return []byte(this.String()), nil
}

func init() {
	commaSepIntType := reflect.TypeOf(CommaSepInt{}).String()
	jsonIter.RegisterTypeEncoderFunc(commaSepIntType, func(ptr unsafe.Pointer, stream *jsonIter.Stream) {
		stream.WriteString((*CommaSepInt)(ptr).String())
	}, func(ptr unsafe.Pointer) bool {
		return len(*(*CommaSepInt)(ptr)) == 0
	})
	jsonIter.RegisterTypeDecoderFunc(commaSepIntType, func(ptr unsafe.Pointer, iter *jsonIter.Iterator) {
		any := iter.ReadAny()
		switch any.ValueType() {
		case jsonIter.InvalidValue, jsonIter.NilValue:
			*((*CommaSepInt)(ptr)) = nil
		case jsonIter.StringValue, jsonIter.NumberValue:
			(*CommaSepInt)(ptr).FromString(any.ToString())
		case jsonIter.ArrayValue:
			*(*CommaSepInt)(ptr), _ = ToInt(any.GetInterface())
		default:
			iter.ReportError("", fmt.Sprintf("cannot convert '%v' to CommaSepInt", any.GetInterface()))
		}
	})

	commaSepStringType := reflect.TypeOf(CommaSepString{}).String()
	jsonIter.RegisterTypeEncoderFunc(commaSepStringType, func(ptr unsafe.Pointer, stream *jsonIter.Stream) {
		stream.WriteString((*CommaSepString)(ptr).String())
	}, func(ptr unsafe.Pointer) bool {
		return len(*(*CommaSepString)(ptr)) == 0
	})
	jsonIter.RegisterTypeDecoderFunc(commaSepStringType, func(ptr unsafe.Pointer, iter *jsonIter.Iterator) {
		any := iter.ReadAny()
		switch any.ValueType() {
		case jsonIter.InvalidValue, jsonIter.NilValue:
			*((*CommaSepString)(ptr)) = nil
		case jsonIter.StringValue, jsonIter.NumberValue:
			(*CommaSepString)(ptr).FromString(any.ToString())
		case jsonIter.ArrayValue:
			*(*CommaSepString)(ptr), _ = ToStr(any.GetInterface())
		default:
			iter.ReportError("", fmt.Sprintf("cannot convert '%v' to CommaSepString", any.GetInterface()))
		}
	})
}
