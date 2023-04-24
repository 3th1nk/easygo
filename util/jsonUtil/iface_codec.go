package jsonUtil

import (
	"fmt"
	jsonIter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"reflect"
	"strings"
	"unsafe"
)

// 注册接口类型的序列化与反序列化逻辑。
//  默认情况下无法将 JSON 字符串反序列化到接口对象中，因为持续不知道对应的具体的类型。
//  InterfaceCodec 通过在序列化时通过额外的字段把具体的类型信息也序列化到 JSON 字符串中，这样在反序列化时候就能够通过类型名称创建不同类型的对象。
//
// Example:
//  type testA struct {
//     AAA int `json:"aaa,omitempty"`
//  }
//  type testB struct {
//     BBB int `json:"bbb,omitempty"`
//  }
//  func (this *testA) getType() string { return "a" }
//  func (this *testB) getType() string { return "b" }
//
//  type testStruct struct {
//     Obj testInterface `json:"obj,omitempty"`
//  }
//  type testInterface interface {
//     getType() string
//  }
//
//  a := &testStruct{
//     Obj: &testA{ AAA: 111},
//  }
//  util.Println(MustMarshalToString(a)) // {"obj":{"type":"a","data":{"aaa":111}}}
//
//  b := &testStruct{
//     Obj: &testB{BBB: 222},
//  }
//  util.Println(MustMarshalToString(b)) // {"obj":{"type":"b","data":{"bbb":222}}}
func RegisterInterfaceCodec(i InterfaceCodec) {
	if i.Type == nil {
		panic(fmt.Errorf("type must not be nil"))
	}
	if i.Encode == nil {
		panic(fmt.Errorf("encode must not be nil"))
	}
	if i.Decode == nil {
		panic(fmt.Errorf("decode must not be nil"))
	}

	t := i.Type
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Interface {
		panic(fmt.Errorf("type '%v' is not interface, but %v", i.Type, i.Type.Kind()))
	}
	i.Type = t

	if i.TypeField == "" {
		i.TypeField = "type"
	}
	if i.DataField == "" {
		i.DataField = "data"
	}

	typeStr, codec := i.Type.String(), &codecWrap{codec: &i}
	jsonIter.RegisterTypeEncoder(typeStr, codec)
	jsonIter.RegisterTypeDecoder(typeStr, codec)
}

// 接口类型的编解码选项。
//   Type:      接口类型
//   TypeField: 序列化时用来存储类型名称的字段
//   DataField: 序列化时用来存储数据的字段
//   Encode:    序列化方法，指定接口对象、返回具体的类型名称和数据对象。
//   Decode:    反序列化方法，通过类型名称和数据对象返回接口对象。
type InterfaceCodec struct {
	Type      reflect.Type                                            // 接口类型
	TypeField string                                                  // 序列化时用来存储类型名称的字段
	DataField string                                                  // 序列化时用来存储数据的字段
	Encode    func(obj interface{}) (typ string, data interface{})    // 序列化方法，指定接口对象、返回具体的类型名称和数据对象。
	Decode    func(typ string, data interface{}) (interface{}, error) // 反序列化方法，通过类型名称和数据对象返回接口对象。
}

type codecWrap struct {
	codec *InterfaceCodec
}

func (this *codecWrap) IsEmpty(ptr unsafe.Pointer) bool {
	return reflect2.IsNil(reflect.NewAt(this.codec.Type, ptr).Elem().Interface())
}

func (this *codecWrap) Encode(ptr unsafe.Pointer, stream *jsonIter.Stream) {
	obj := reflect.NewAt(this.codec.Type, ptr).Elem().Interface()
	if reflect2.IsNil(obj) {
		stream.WriteNil()
		return
	}

	typ, data := this.codec.Encode(obj)
	if data != nil {
		if v := reflect.ValueOf(data); v.Type() == this.codec.Type {
			data = v.Elem().Interface()
		}
	}

	stream.WriteObjectStart()
	stream.WriteObjectField(this.codec.TypeField)
	stream.WriteString(typ)
	if !reflect2.IsNil(data) {
		// 保留写入 data 前的缓存
		bufBeforeData := stream.Buffer()

		stream.WriteMore()
		stream.WriteObjectField(this.codec.DataField)
		posBeforeDataVal := stream.Buffered()
		stream.WriteVal(data)
		// 判断写入的 data 是否是 null、空字符串、空数组、空对象，如果是则回退到 bufBeforeData 丢弃 data 部分
		dataStr := string(stream.Buffer()[posBeforeDataVal:])
		if n := len(dataStr); n >= 2 {
			if c := dataStr[0]; c == '"' || c == '{' || c == '[' {
				dataStr = strings.TrimSpace(dataStr[1 : n-1])
				if dataStr == "" {
					stream.SetBuffer(bufBeforeData)
				} else {
					lines, empty := strings.Split(dataStr, "\n"), true
					for _, s := range lines {
						if s = strings.TrimSpace(s); s != "" {
							empty = false
							break
						}
					}
					if empty {
						stream.SetBuffer(bufBeforeData)
					}
				}
			}
		}
	}
	stream.WriteObjectEnd()
}

func (this *codecWrap) Decode(ptr unsafe.Pointer, iter *jsonIter.Iterator) {
	any := iter.ReadAny()
	switch any.ValueType() {
	case jsonIter.NilValue:
		return
	case jsonIter.ObjectValue:
		typ := any.Get(this.codec.TypeField).ToString()
		data := any.Get(this.codec.DataField).GetInterface()
		obj, err := this.codec.Decode(typ, data)
		if obj != nil {
			reflect.NewAt(this.codec.Type, ptr).Elem().Set(reflect.ValueOf(obj))
		}
		if err != nil {
			iter.Error = fmt.Errorf("decode %v error: %v, data=%v(%v)", typ, err, reflect.TypeOf(data), MustMarshalToString(data))
		}
	default:
		val := any.GetInterface()
		iter.Error = fmt.Errorf("cannot decode %v(%v) to %v", reflect.TypeOf(val), MustMarshalToString(val), this.codec.Type.String())
	}
}
