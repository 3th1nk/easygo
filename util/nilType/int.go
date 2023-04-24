package nilType

import (
	"database/sql/driver"
	"fmt"
	"github.com/3th1nk/easygo/util/convertor"
	"github.com/3th1nk/easygo/util/jsonUtil"
	jsonIter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"reflect"
	"strconv"
	"unsafe"
)

func NewInt(a ...interface{}) (val Int) {
	if len(a) != 0 {
		_ = val.SetValue(a[0])
	}
	return
}

func init() {
	boolType := reflect.TypeOf(Int{})
	jsonIter.RegisterTypeEncoderFunc(boolType.String(), func(ptr unsafe.Pointer, stream *jsonIter.Stream) {
		v := (*Int)(ptr)
		if !reflect2.IsNil(v) && v.ok {
			_, _ = stream.Write(jsonUtil.MustMarshal(v.val))
		} else {
			stream.WriteString("null")
		}
	}, func(pointer unsafe.Pointer) bool {
		return !(*Int)(pointer).ok
	})
	jsonIter.RegisterTypeDecoderFunc(boolType.String(), func(ptr unsafe.Pointer, iter *jsonIter.Iterator) {
		err := (*Int)(ptr).SetValue(iter.Read())
		if err != nil {
			iter.ReportError("NilIntDecoder", fmt.Sprintf("convert error: %v", err))
		}
	})
}

type Int struct {
	ok  bool
	val int64
}

func (this Int) IntValue() int {
	if !reflect2.IsNil(this) && this.ok {
		return int(this.val)
	}
	return 0
}

func (this Int) Int64Value() int64 {
	if !reflect2.IsNil(this) && this.ok {
		return this.val
	}
	return 0
}

func (this *Int) IsNotNil() bool {
	if !reflect2.IsNil(this) {
		return this.ok
	}
	return false
}

func (this *Int) InterfaceValue() interface{} {
	if !reflect2.IsNil(this) && this.ok {
		return this.val
	}
	return nil
}

func (this *Int) String() string {
	if !reflect2.IsNil(this) && this.ok {
		return convertor.ToStringNoError(this.val)
	}
	return ""
}

func (this *Int) SetValue(val interface{}) (err error) {
	if !reflect2.IsNil(val) {
		this.val, err = strconv.ParseInt(convertor.ToStringNoError(val), 10, 64)
		this.ok = err == nil
		return err
	} else {
		this.ok = false
		return nil
	}
}

func (this *Int) Scan(val interface{}) (err error) {
	return this.SetValue(val)
}

func (this Int) Value() (driver.Value, error) {
	if !reflect2.IsNil(this) {
		if this.ok {
			return this.val, nil
		}
	}
	return nil, nil
}

func (this Int) MarshalJSON() ([]byte, error) {
	if !reflect2.IsNil(this) && this.ok {
		return jsonUtil.Marshal(this.val)
	} else {
		return []byte("null"), nil
	}
}

func (this *Int) UnmarshalJSON(bytes []byte) error {
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
