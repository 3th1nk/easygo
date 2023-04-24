package nilType

import (
	"github.com/3th1nk/easygo/util/jsonUtil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBool_Batch(t *testing.T) {
	// []: Bool, Valid, String, MarshalToString
	cases := [][]interface{}{
		{NewBool(nil), false, "", `"null"`},
		{NewBool((*struct{})(nil)), false, "", `"null"`},
		{NewBool(true), true, "true", "true"},
		{NewBool(false), true, "false", "false"},
	}

	for i, arr := range cases {
		valid, str, jsonStr := arr[1].(bool), arr[2].(string), arr[3].(string)

		v1 := arr[0].(Bool)
		if expect, actual := valid, v1.IsNotNil(); expect != actual {
			t.Errorf("[%d] assert faild: expect %v, but %v", i, expect, actual)
		}
		if expect, actual := str, v1.String(); expect != actual {
			t.Errorf("[%d] assert faild: expect %v, but %v", i, expect, actual)
		}
		if expect, actual := jsonStr, jsonUtil.MustMarshalToString(v1); expect != actual {
			t.Errorf("[%d] assert faild: expect %v, but %v", i, expect, actual)
		}

		v2 := &v1
		if expect, actual := valid, v2.IsNotNil(); expect != actual {
			t.Errorf("[%d] assert faild: expect %v, but %v", i, expect, actual)
		}
		if expect, actual := str, v2.String(); expect != actual {
			t.Errorf("[%d] assert faild: expect %v, but %v", i, expect, actual)
		}
		if expect, actual := jsonStr, jsonUtil.MustMarshalToString(v2); expect != actual {
			t.Errorf("[%d] assert faild: expect %v, but %v", i, expect, actual)
		}

		v3 := Bool{}
		_ = jsonUtil.UnmarshalFromString(jsonStr, &v3)
		if actual, expect := v1.val, v3.val; actual != expect {
			t.Errorf("[%d] assert faild: expect %v, but %v, jsonStr=%v", i, actual, expect, jsonStr)
		}
		if actual, expect := v1.ok, v3.ok; actual != expect {
			t.Errorf("[%d] assert faild: expect %v, but %v, jsonStr=%v", i, actual, expect, jsonStr)
		}

		v4 := &Bool{}
		_ = jsonUtil.UnmarshalFromString(jsonStr, v4)
		if actual, expect := v1.val, v4.val; actual != expect {
			t.Errorf("[%d] assert faild: expect %v, but %v, jsonStr=%v", i, actual, expect, jsonStr)
		}
		if actual, expect := v1.ok, v4.ok; actual != expect {
			t.Errorf("[%d] assert faild: expect %v, but %v, jsonStr=%v", i, actual, expect, jsonStr)
		}
	}
}

func TestBool_MarshalJSON(t *testing.T) {
	v := NewBool(false)
	str := jsonUtil.MustMarshalToString(v)
	assert.Equal(t, str, "false")
}

func TestBool_UnmarshalJSON(t *testing.T) {
	jsonStr := "true"
	v := NewBool()
	err := jsonUtil.UnmarshalFromString(jsonStr, &v)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(v)
}

func TestBool_InStruct(t *testing.T) {
	type structA struct {
		Val Bool `json:"val,omitempty"`
	}
	type structB struct {
		Val *Bool `json:"val,omitempty"`
	}

	{
		v := structB{}
		jsonUtil.UnmarshalFromString(`{"val":""}`, &v)
		if assert.NotNil(t, v.Val) {
			assert.Equal(t, v.Val.IsNotNil(), false)
		}
	}

	{
		v := structA{}
		assert.Equal(t, v.Val.IsNotNil(), false)
		assert.Equal(t, jsonUtil.MustMarshalToString(v), `{}`)
	}

	{
		v := structA{Val: NewBool(false)}
		assert.Equal(t, v.Val.IsNotNil(), true)
		assert.Equal(t, jsonUtil.MustMarshalToString(v), `{"val":false}`)
	}

	{
		v := structA{Val: NewBool(true)}
		assert.Equal(t, v.Val.IsNotNil(), true)
		assert.Equal(t, jsonUtil.MustMarshalToString(v), `{"val":true}`)
	}

	{
		v := structA{}
		jsonUtil.UnmarshalFromString(`{}`, &v)
		assert.Equal(t, v.Val.IsNotNil(), false)
	}

	{
		v := structA{}
		jsonUtil.UnmarshalFromString(`{"val":null}`, &v)
		assert.Equal(t, v.Val.IsNotNil(), false)
	}

	{
		v := structA{}
		jsonUtil.UnmarshalFromString(`{"val":false}`, &v)
		assert.Equal(t, v.Val.IsNotNil(), true)
		assert.Equal(t, v.Val.BoolValue(), false)
	}

	{
		v := structA{}
		jsonUtil.UnmarshalFromString(`{"val":true}`, &v)
		assert.Equal(t, v.Val.IsNotNil(), true)
		assert.Equal(t, v.Val.BoolValue(), true)
	}

	{
		v := structB{}
		assert.Equal(t, v.Val.IsNotNil(), false)
		assert.Equal(t, jsonUtil.MustMarshalToString(v), `{}`)
	}

	{
		v := structB{Val: nil}
		assert.Equal(t, v.Val.IsNotNil(), false)
		assert.Equal(t, jsonUtil.MustMarshalToString(v), `{}`)
	}

	{
		v := structB{Val: &Bool{ok: true, val: false}}
		assert.Equal(t, v.Val.IsNotNil(), true)
		assert.Equal(t, jsonUtil.MustMarshalToString(v), `{"val":false}`)
	}

	{
		v := structB{Val: &Bool{ok: true, val: true}}
		assert.Equal(t, v.Val.IsNotNil(), true)
		assert.Equal(t, jsonUtil.MustMarshalToString(v), `{"val":true}`)
	}

	{
		v := structB{}
		jsonUtil.UnmarshalFromString(`{}`, &v)
		assert.Nil(t, v.Val)
	}

	{
		v := structB{}
		jsonUtil.UnmarshalFromString(`{"val":null}`, &v)
		assert.Nil(t, v.Val)
		assert.Equal(t, v.Val.IsNotNil(), false)
	}

	{
		v := structB{}
		jsonUtil.UnmarshalFromString(`{"val":""}`, &v)
		if assert.NotNil(t, v.Val) {
			assert.Equal(t, v.Val.IsNotNil(), false)
		}
	}

	{
		v := structB{}
		jsonUtil.UnmarshalFromString(`{"val":false}`, &v)
		if assert.NotNil(t, v.Val) {
			assert.Equal(t, v.Val.IsNotNil(), true)
			assert.Equal(t, v.Val.BoolValue(), false)
		}
	}

	{
		v := structB{}
		jsonUtil.UnmarshalFromString(`{"val":true}`, &v)
		if assert.NotNil(t, v.Val) {
			assert.Equal(t, v.Val.IsNotNil(), true)
			assert.Equal(t, v.Val.BoolValue(), true)
		}
	}
}
