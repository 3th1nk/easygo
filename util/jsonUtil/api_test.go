package jsonUtil

import (
	"encoding/json"
	"fmt"
	"github.com/3th1nk/easygo/util/convertor"
	"github.com/3th1nk/easygo/util/timeUtil"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUnmarshal(t *testing.T) {
	a := &struct {
		A string  `json:"a"`
		B int     `json:"b"`
		C int32   `json:"c"`
		D int64   `json:"d"`
		E float32 `json:"e"`
		F float64 `json:"f"`
		G bool    `json:"g"`
	}{}
	err := UnmarshalFromString(`{"a":"true","b":true,"c":true,"d":true,"e":true,"f":true,"g":true}`, &a)
	if !assert.NoError(t, err) {
		return
	}

	t.Log(MustMarshalToString(a))
	assert.NotEmpty(t, a.A)
	assert.NotEmpty(t, a.B)
	assert.NotEmpty(t, a.C)
	assert.NotEmpty(t, a.D)
	assert.NotEmpty(t, a.E)
	assert.NotEmpty(t, a.F)
	assert.NotEmpty(t, a.G)
}

func TestUnMarshalToNilPtr(t *testing.T) {
	var a *struct {
		A string  `json:"a"`
		B int     `json:"b"`
		C int32   `json:"c"`
		D int64   `json:"d"`
		E float32 `json:"e"`
		F float64 `json:"f"`
		G bool    `json:"g"`
	}
	err := UnmarshalFromString(`{"a":"true","b":true,"c":true,"d":true,"e":true,"f":true,"g":true}`, &a)
	if !assert.NoError(t, err) {
		return
	}

	t.Log(MustMarshalToString(a))
	assert.NotEmpty(t, a.A)
	assert.NotEmpty(t, a.B)
	assert.NotEmpty(t, a.C)
	assert.NotEmpty(t, a.D)
	assert.NotEmpty(t, a.E)
	assert.NotEmpty(t, a.F)
	assert.NotEmpty(t, a.G)
}

func TestUnmarshalToString(t *testing.T) {
	a := &struct {
		A string `json:"a"`
	}{}

	arr := []interface{}{
		nil,
		true,
		123,
		123.456,
		"abc",
		[]interface{}{},
		[]bool{true, false},
		[]int{1, 2, 3},
		[]float32{1.1, 2.2, 3.3},
		[]string{"a", "b", "c"},
		map[string]interface{}{},
		map[string]interface{}{"a": 1, "b": 2},
	}
	for _, item := range arr {
		a.A = ""
		str := MustMarshalToString(map[string]interface{}{"a": item})
		if assert.NotEqual(t, "", str) {
			err := UnmarshalFromString(str, &a)
			if assert.NoError(t, err) {
				if convertor.IsEmpty(item) {
					assert.Equal(t, "", a.A, MustMarshalToString(item))
				} else {
					assert.NotEqual(t, "", a.A, MustMarshalToString(item))
				}
			}
		}
		t.Logf("%s: %v", a.A, MustMarshalToString(item))
	}
}

func TestMarshalUnmarshal(t *testing.T) {
	type tp struct {
		Time time.Time
	}

	a := &tp{
		Time: time.Now(),
	}
	str := MustMarshalToStringIndent(a)
	t.Log(str)

	b := &tp{}
	err := UnmarshalFromString(str, b)
	if err != nil {
		t.Error(err)
		return
	}

	strA := timeUtil.Format(a.Time)
	strB := timeUtil.Format(b.Time)
	if strA != strB {
		t.Errorf("assert faild: %v, %v", strA, strB)
		return
	}
	t.Log(strA)
}

func TestUnmarshalFromString(t *testing.T) {
	type strucA struct {
		Name string `json:"name"`
	}
	var a *strucA
	str := `{"abc":null}`
	err := UnmarshalFromString(str, &a)
	if err != nil {
		t.Error(err)
	} else if a.Name != "" {
		t.Errorf("assert faild: %v", a.Name)
	}
}

func TestMarshalToStringIndent(t *testing.T) {
	obj := map[string]interface{}{
		"id-1":   1,
		"name-1": "a",
		"map-1": map[string]interface{}{
			"id-2":   2,
			"name-2": "b",
			"map-2": map[string]interface{}{
				"id-n":   3,
				"name-n": "c",
			},
		},
	}
	fmt.Println(SortMapKeysApi().MarshalToStringIndent(obj))
	fmt.Println(SortMapKeysApi(false).MarshalToStringIndent(obj))
}

func TestMarshalTime(t *testing.T) {
	a := struct {
		T time.Time `json:"t"`
	}{T: time.Now()}
	if str := MustMarshalToString(a); str != fmt.Sprintf(`{"t":"%v"}`, timeUtil.Format(a.T)) {
		t.Errorf("assert faild: %v", str)
	} else {
		t.Log(str)
		tmp, _ := json.Marshal(a)
		t.Log(string(tmp))
	}
}

func TestMarshalRegexPattern(t *testing.T) {
	type tmp struct {
		Regex string `json:"regex,omitempty"`
	}

	a := &tmp{Regex: `^[\w-.]+@[\w-.]+\.\w+$`}
	str, err := MarshalToString(a)
	assert.NoError(t, err)

	var b *tmp
	err = UnmarshalFromString(str, &b)
	assert.NoError(t, err)

	assert.Equal(t, a.Regex, b.Regex)
}

func Test_UnmarshalFromObject(t *testing.T) {
	type A struct {
		Data string `json:"data"`
	}
	type B struct {
		Data interface{} `json:"data"`
	}

	{
		var a *A
		str, _ := MarshalToString(&B{Data: 123})
		err := UnmarshalFromObject(str, &a)
		if assert.NoError(t, err) && assert.NotNil(t, a) {
			assert.Equal(t, "123", a.Data)
		}
	}
	{
		var a *A
		str, _ := MarshalToString(&B{Data: "abc"})
		err := UnmarshalFromObject(str, &a)
		if assert.NoError(t, err) && assert.NotNil(t, a) {
			assert.Equal(t, "abc", a.Data)
		}
	}
	{
		var a *A
		str, _ := MarshalToString(&B{Data: map[string]string{"name": "abc"}})
		err := UnmarshalFromObject(str, &a)
		if assert.NoError(t, err) && assert.NotNil(t, a) {
			assert.Equal(t, `{"name":"abc"}`, a.Data)
		}
	}
}

func TestMarshalToStringIndent2(t *testing.T) {
	str, err := MarshalToStringIndent(map[string]interface{}{
		"id":   123,
		"name": "abc",
	}, "    ", "---- ")
	assert.NoError(t, err)
	t.Logf("\n%s", str)
}
