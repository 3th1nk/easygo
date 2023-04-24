package mapUtil

import (
	"github.com/3th1nk/easygo/util/jsonUtil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringObjectMap_Set0(t *testing.T) {
	{
		obj := &struct {
			Data StringObjectMap `json:"data,omitempty"`
		}{}

		assert.Equal(t, `{}`, jsonUtil.MustMarshalToString(obj))

		assert.Panics(t, func() {
			obj.Data.Set("id", 1)
		})
	}
	{
		obj := &struct {
			Data StringObjectMap `json:"data,omitempty"`
		}{Data: StringObjectMap{}}

		assert.Equal(t, `{}`, jsonUtil.MustMarshalToString(obj))

		obj.Data.Set("id", 1)
		assert.Equal(t, 1, obj.Data.MustGetInt("id"))

		obj.Data.Set("name", "abc")
		assert.Equal(t, "abc", obj.Data.MustGetString("name"))

		assert.Equal(t, `{"data":{"id":1,"name":"abc"}}`, jsonUtil.MustMarshalToString(obj))
	}
}

func TestStringObjectMap_Set1(t *testing.T) {
	{
		var a StringObjectMap
		assert.Panics(t, func() {
			a.Set("id", 1)
		})
	}

	{
		a := make(StringObjectMap)
		a.Set("id", 1)
		if val := a.MustGetInt("id"); val != 1 {
			t.Errorf("assert faild: %v", val)
		}
	}

}

func TestStringObjectMap(t *testing.T) {
	a := make(StringObjectMap)
	a.SetMulti(map[string]interface{}{
		"id":   "1",
		"name": "abc",
		"dict": map[string]interface{}{
			"abc": 123,
			"bcd": "aaa",
		},
	})
	if v := a.MustGetInt("id"); v != 1 {
		t.Errorf("assert faild: %v", v)
	}
	if v, found := a.Get("asdbsadf"); v != nil || found {
		t.Errorf("assert faild: %v", v)
	}
	var dict StringObjectMap
	if found, err := a.GetToObject("dict", &dict); dict == nil || !found || err != nil {
		t.Errorf("assert faild: %v", jsonUtil.MustMarshalToString(dict))
	} else if v := StringObjectMap(dict).MustGetInt("abc"); v != 123 {
		t.Errorf("assert faild: %v, %v", jsonUtil.MustMarshalToString(v), jsonUtil.MustMarshalToString(dict))
	} else if v := dict["bcd"]; v != "aaa" {
		t.Errorf("assert faild: %v, %v", jsonUtil.MustMarshalToString(v), jsonUtil.MustMarshalToString(dict))
	}

	b := make(StringObjectMap)
	jsonUtil.UnmarshalFromString(`{"id": "1", "name": "abc"}`, &b)
	if v := b.MustGetInt("id"); v != 1 {
		t.Errorf("assert faild: %v", v)
	}
	if v := b.MustGetInt("IDID"); v != 0 {
		t.Errorf("assert faild: %v", v)
	}

	b1 := make(StringObjectMap)
	b1.FromDB([]byte(`{"id": "1", "name": "abc"}`))
	if v := b1.MustGetInt("id"); v != 1 {
		t.Errorf("assert faild: %v", v)
	}
	if v := b1.MustGetInt("IDID"); v != 0 {
		t.Errorf("assert faild: %v", v)
	}

	c := map[string]interface{}(b)
	if v := c["id"]; v != "1" {
		t.Errorf("assert faild: %v", v)
	}

	d := StringObjectMap(c)
	if v := d["id"]; v != "1" {
		t.Errorf("assert faild: %v", v)
	}
	if v, _ := d.Get("ID"); v != nil {
		t.Errorf("assert faild: %v", v)
	}

	e := CaseInsensitiveStringObjectMap(b)
	if v := e["id"]; v != "1" {
		t.Errorf("assert faild: %v", v)
	}
	if v, _ := e.Get("ID"); v != "1" {
		t.Errorf("assert faild: %v", v)
	}
}

func TestStringObjectMap_TrimEmptyValues(t *testing.T) {
	a := StringObjectMap{
		"obj":     nil,
		"arr_nil": []string(nil),
		"map_nil": map[string]string(nil),
		"str":     "",
		"int":     0,
		"int8":    int8(0),
		"float":   float64(0),
		"arr1":    []interface{}{},
		"arr2":    []*StringObjectMap{},
		"map1":    map[string]interface{}{},
		"map2":    map[string]int{},
	}
	b := a.TrimEmptyValues()
	assert.Equal(t, 0, len(b))
	assert.Equal(t, "{}", jsonUtil.MustMarshalToString(b))
	bytes, _ := b.ToDB()
	assert.Equal(t, "", string(bytes))
}
