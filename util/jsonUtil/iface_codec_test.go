package jsonUtil

import (
	"fmt"
	jsonIter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type testA struct {
	Id int `json:"id,omitempty"`
}

func (this *testA) getType() string { return "a" }
func (this *testA) val() int        { return this.Id }

type testB struct {
	Id int `json:"id,omitempty"`
}

func (this *testB) getType() string { return "b" }
func (this *testB) val() int        { return this.Id }

type testInterface interface {
	getType() string
	val() int
}

// 接口类型无法被反序列化
func Test_1(t *testing.T) {
	{
		var v1 testInterface
		v1 = &testA{Id: 1}
		str := MustMarshalToString(v1)
		assert.Equal(t, `{"id":1}`, str)
		t.Logf("json str: %v", str)

		var v2 testInterface
		err := UnmarshalFromString(str, &v2)
		assert.Nil(t, v2)
		assert.Error(t, err)
		t.Logf("unmarshal error: %v", err)
	}

	{
		type tmp struct {
			Obj testInterface `json:"obj,omitempty"`
		}

		v1 := &tmp{Obj: &testA{Id: 1}}
		str := MustMarshalToString(v1)
		assert.Equal(t, `{"obj":{"id":1}}`, str)
		t.Logf("json str: %v", str)

		var v2 *tmp
		err := UnmarshalFromString(str, &v2)
		assert.NotNil(t, v2)
		assert.Nil(t, v2.Obj)
		assert.Error(t, err)
		t.Logf("unmarshal error: %v", err)
	}
}

// 使用 结构体 代替接口类型，配合工厂方法，绕开 反序列化接口 的需求。
func Test_2(t *testing.T) {
	a1 := &testA{Id: 1}

	{
		v1 := &testInterfaceDesc{Type: a1.getType(), Data: MustMarshal(a1)}
		str := MustMarshalToString(v1)
		t.Logf("json str: %v", str)

		var v2 *testInterfaceDesc
		err := UnmarshalFromString(str, &v2)
		assert.NoError(t, err)

		a2, err := newTestInterface(v2.Type, v2.Data)
		assert.NoError(t, err)
		assert.Equal(t, a1.Id, a2.(*testA).Id)
		t.Logf("a2: [%v] %+v", reflect.TypeOf(a2), a2)
	}

	{
		type tmp struct {
			Obj *testInterfaceDesc `json:"obj,omitempty"`
		}

		v1 := &tmp{Obj: newTestInterfaceDesc(a1)}
		str := MustMarshalToString(v1)
		t.Logf("json str: %v", str)

		var v2 *tmp
		err := UnmarshalFromString(str, &v2)
		assert.NoError(t, err)

		a2, err := v2.Obj.GetTestInterface()
		assert.NoError(t, err)
		assert.Equal(t, a1.Id, a2.(*testA).Id)
		t.Logf("a2: [%v] %+v", reflect.TypeOf(a2), a2)

		// hidden danger ?????
		n := 0
		for i := 0; i < 1000; i++ {
			a3, _ := v2.Obj.GetTestInterface()
			n += a3.val()
		}
		assert.Equal(t, 1000, n)
	}

	// 每个需要反序列化的接口，都需要定义一个【额外的结构】体用来序列化和反序列化；
	// 在所有需要用到该接口的地方，都需要通过【额外的转换代码】来实现  接口->结构体  或  结构体->接口  的转换。
	// 结构体->接口 转换后的对象如果需要多次使用，则调用方需要通过额外的代码存储转换结果而不是每次都重新转换，调用方必须充分了解这里的【代码隐患】。
}

type testInterfaceDesc struct {
	Type string              `json:"type,omitempty"`
	Data jsonIter.RawMessage `json:"data,omitempty"`
}

func newTestInterfaceDesc(a testInterface) *testInterfaceDesc {
	if !reflect2.IsNil(a) {
		return &testInterfaceDesc{Type: a.getType(), Data: MustMarshal(a)}
	}
	return nil
}

func (this *testInterfaceDesc) GetTestInterface() (testInterface, error) {
	if obj, err := newTestInterface(this.Type, this.Data); err != nil {
		return nil, err
	} else {
		return obj.(testInterface), nil
	}
}

// 通过自定义 JSON 序列化与反序列化方法，实现接口字段反序列化
func Test_3_1(t *testing.T) {
	a1 := &testA{Id: 1}

	v1 := &TestStruct1{Id: 123, Name: "abc", Obj: a1}
	str := MustMarshalToString(v1)

	// 能成功的反序列化接口类型的字段
	var v2 *TestStruct1
	err := UnmarshalFromString(str, &v2)
	assert.NoError(t, err)
	assert.Equal(t, a1.Id, v2.Obj.(*testA).Id)
	t.Logf("a2: [%v] %+v", reflect.TypeOf(v2.Obj), v2.Obj)
}

func Test_3_2(t *testing.T) {
	a1 := &testA{Id: 1}

	v1 := &TestStruct1{Id: 123, Name: "abc", Obj: a1}
	str := MustMarshalToString(v1)
	t.Logf("json str: %v", str)

	// 【额外的内部结构体】
}

type TestStruct1 struct {
	Id    int         `json:"id,omitempty"`
	Name  string      `json:"name,omitempty"`
	Other interface{} `json:"other,omitempty"`

	Obj testInterface `json:"obj,omitempty"`
}

func (this *TestStruct1) MarshalJSON() ([]byte, error) {
	wrap := &testStructJsonWrap{
		Id:    this.Id,
		Name:  this.Name,
		Other: this.Other,
	}
	if this.Obj != nil {
		wrap.ObjType = this.Obj.getType()
		wrap.ObjData = MustMarshal(this.Obj)
	}
	return Marshal(wrap)
}

func (this *TestStruct1) UnmarshalJSON(bytes []byte) error {
	var wrap *testStructJsonWrap
	if err := Unmarshal(bytes, &wrap); err != nil {
		return err
	}
	*this = TestStruct1{
		Id:    wrap.Id,
		Name:  wrap.Name,
		Other: wrap.Other,
	}
	if wrap.ObjType != "" {
		if v, err := newTestInterface(wrap.ObjType, wrap.ObjData); err != nil {
			return err
		} else {
			this.Obj = v.(testInterface)
		}
	}
	return nil
}

type testStructJsonWrap struct {
	Id    int         `json:"id,omitempty"`
	Name  string      `json:"name,omitempty"`
	Other interface{} `json:"other,omitempty"`

	ObjType string              `json:"obj_type,omitempty"`
	ObjData jsonIter.RawMessage `json:"obj_data,omitempty"`
}

func Test_4(t *testing.T) {
	type TestStruct2 struct {
		Id    int           `json:"id,omitempty"`
		Name  string        `json:"name,omitempty"`
		Other interface{}   `json:"other,omitempty"`
		Obj   testInterface `json:"obj,omitempty"`
	}

	registerCodec()

	a1 := &testA{Id: 1}

	v1 := &TestStruct2{Id: 123, Name: "abc", Obj: a1}
	str := MustMarshalToString(v1)

	// 能成功的反序列化接口类型的字段
	var v2 *TestStruct2
	err := UnmarshalFromString(str, &v2)
	assert.NoError(t, err)
	assert.Equal(t, a1.Id, v2.Obj.(*testA).Id)
	t.Logf("a2: [%v] %+v", reflect.TypeOf(v2.Obj), v2.Obj)

	t.Logf("json str: %v", str)
}

func registerCodec() {
	RegisterInterfaceCodec(InterfaceCodec{
		Type: reflect.TypeOf(new(testInterface)),
		Encode: func(obj interface{}) (typ string, data interface{}) {
			return obj.(testInterface).getType(), obj
		},
		Decode: newTestInterface,
	})
}

func newTestInterface(typ string, data interface{}) (obj interface{}, err error) {
	switch typ {
	case "a":
		obj = &testA{}
	case "b":
		obj = &testB{}
	default:
		return nil, fmt.Errorf("unknown type '%s'", typ)
	}
	if err = UnmarshalFromObject(data, obj); err != nil {
		return nil, err
	}
	return obj, nil
}
