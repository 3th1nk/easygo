package nilType

import (
	"database/sql/driver"
	"fmt"
	"github.com/3th1nk/easygo/util/convertor"
	"github.com/3th1nk/easygo/util/jsonUtil"
	jsonIter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"reflect"
	"unsafe"
)

func NewString(a ...interface{}) (val String) {
	if len(a) != 0 {
		_ = val.SetValue(a[0])
	}
	return
}

func init() {
	boolType := reflect.TypeOf(String{})
	jsonIter.RegisterTypeEncoderFunc(boolType.String(), func(ptr unsafe.Pointer, stream *jsonIter.Stream) {
		v := (*String)(ptr)
		if !reflect2.IsNil(v) && v.ok {
			_, _ = stream.Write(jsonUtil.MustMarshal(v.val))
		} else {
			stream.WriteString("null")
		}
	}, func(pointer unsafe.Pointer) bool {
		return !(*String)(pointer).ok
	})
	jsonIter.RegisterTypeDecoderFunc(boolType.String(), func(ptr unsafe.Pointer, iter *jsonIter.Iterator) {
		err := (*String)(ptr).SetValue(iter.Read())
		if err != nil {
			iter.ReportError("NilStringDecoder", fmt.Sprintf("convert error: %v", err))
		}
	})
}

type String struct {
	ok  bool
	val string
}

func (this String) StringValue() string {
	if !reflect2.IsNil(this) && this.ok {
		return this.val
	}
	return ""
}

func (this *String) IsNotNil() bool {
	if !reflect2.IsNil(this) {
		return this.ok
	}
	return false
}

func (this *String) InterfaceValue() interface{} {
	if !reflect2.IsNil(this) && this.ok {
		return this.val
	}
	return nil
}

func (this *String) String() string {
	if !reflect2.IsNil(this) && this.ok {
		return convertor.ToStringNoError(this.val)
	}
	return ""
}

func (this *String) SetValue(val interface{}) (err error) {
	if !reflect2.IsNil(val) {
		this.val, err = convertor.ToString(val)
		this.ok = err == nil
		return err
	} else {
		this.ok = false
		return nil
	}
}

func (this *String) Scan(val interface{}) (err error) {
	return this.SetValue(val)
}

func (this String) Value() (driver.Value, error) {
	if !reflect2.IsNil(this) {
		if this.ok {
			return this.val, nil
		}
	}
	return nil, nil
}

func (this String) MarshalJSON() ([]byte, error) {
	if !reflect2.IsNil(this) && this.ok {
		return jsonUtil.Marshal(this.val)
	} else {
		return []byte("null"), nil
	}
}

func (this *String) UnmarshalJSON(bytes []byte) error {
	if !reflect2.IsNil(this) {
		if len(bytes) != 0 {
			obj := jsonUtil.Get(bytes)
			switch obj.ValueType() {
			case jsonIter.InvalidValue:
				return obj.LastError()
			default:
				return this.Scan(obj.GetInterface())
			}
		} else {
			this.ok = false
		}
	}
	return nil
}
