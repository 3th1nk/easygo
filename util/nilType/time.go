package nilType

import (
	"database/sql/driver"
	"fmt"
	"github.com/3th1nk/easygo/util/convertor"
	"github.com/3th1nk/easygo/util/jsonUtil"
	"github.com/3th1nk/easygo/util/timeUtil"
	jsonIter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"reflect"
	"time"
	"unsafe"
)

func NewTime(a ...interface{}) (val Time) {
	if len(a) != 0 {
		_ = val.SetValue(a[0])
	}
	return
}

func init() {
	boolType := reflect.TypeOf(Time{})
	jsonIter.RegisterTypeEncoderFunc(boolType.String(), func(ptr unsafe.Pointer, stream *jsonIter.Stream) {
		v := (*Time)(ptr)
		if !reflect2.IsNil(v) && v.ok {
			_, _ = stream.Write(jsonUtil.MustMarshal(v.val))
		} else {
			stream.WriteString("null")
		}
	}, func(pointer unsafe.Pointer) bool {
		return !(*Time)(pointer).ok
	})
	jsonIter.RegisterTypeDecoderFunc(boolType.String(), func(ptr unsafe.Pointer, iter *jsonIter.Iterator) {
		err := (*Time)(ptr).SetValue(iter.Read())
		if err != nil {
			iter.ReportError("NilTimeDecoder", fmt.Sprintf("convert error: %v", err))
		}
	})
}

type Time struct {
	ok  bool
	val time.Time
}

func (this Time) TimeValue() time.Time {
	if !reflect2.IsNil(this) && this.ok {
		return this.val
	}
	return time.Time{}
}

func (this *Time) IsNotNil() bool {
	if !reflect2.IsNil(this) {
		return this.ok
	}
	return false
}

func (this *Time) InterfaceValue() interface{} {
	if !reflect2.IsNil(this) && this.ok {
		return this.val
	}
	return nil
}

func (this *Time) String() string {
	if !reflect2.IsNil(this) && this.ok {
		return convertor.ToStringNoError(this.val)
	}
	return ""
}

func (this *Time) SetValue(val interface{}) (err error) {
	if !reflect2.IsNil(val) {
		this.val, err = timeUtil.Parse(convertor.ToStringNoError(val))
		this.ok = err == nil
		return err
	} else {
		this.ok = false
		return nil
	}
}

func (this *Time) Scan(val interface{}) (err error) {
	return this.SetValue(val)
}

func (this Time) Value() (driver.Value, error) {
	if !reflect2.IsNil(this) {
		if this.ok {
			return this.val, nil
		}
	}
	return nil, nil
}

func (this Time) MarshalJSON() ([]byte, error) {
	if !reflect2.IsNil(this) && this.ok {
		return jsonUtil.Marshal(this.val)
	} else {
		return []byte("null"), nil
	}
}

func (this *Time) UnmarshalJSON(bytes []byte) error {
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
