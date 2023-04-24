package mapUtil

import (
	"github.com/3th1nk/easygo/util/convertor"
	"github.com/3th1nk/easygo/util/jsonUtil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringMap(t *testing.T) {
	a := NewCaseInsensitiveStringMap(map[string]string{
		"id":   "1",
		"name": "abc",
	})
	if v := a.MustGetInt("id"); v != 1 {
		t.Errorf("assert faild: %v", v)
	}
	if v, _ := a.Get("asdbsadf"); v != "" {
		t.Errorf("assert faild: %v", v)
	}

	b := make(StringMap)
	jsonUtil.UnmarshalFromString(`{"id": "1", "name": "abc"}`, &b)
	if v := b.MustGetInt("id"); v != 1 {
		t.Errorf("assert faild: %v", v)
	}
	if v := b.MustGetInt("IDID"); v != 0 {
		t.Errorf("assert faild: %v", v)
	}

	c := map[string]string(b)
	if v := c["id"]; v != "1" {
		t.Errorf("assert faild: %v", v)
	}

	d := StringMap(c)
	if v := d["id"]; v != "1" {
		t.Errorf("assert faild: %v", v)
	}
	if v, _ := d.Get("id"); v != "1" {
		t.Errorf("assert faild: %v", v)
	}

	e := CaseInsensitiveStringMap(b)
	if v := e["id"]; v != "1" {
		t.Errorf("assert faild: %v", v)
	}
	if v, _ := e.Get("ID"); v != "1" {
		t.Errorf("assert faild: %v", v)
	}
}

func TestStringObjectMap_ToStringMap(t *testing.T) {
	{
		a := map[string]string(nil)
		b, err := convertor.ToStringMap(a)
		assert.NoError(t, err)
		assert.Nil(t, b)
	}
	{
		a := map[string]string{"id": "123"}
		b, err := convertor.ToStringMap(a)
		if assert.NoError(t, err) && assert.NotEqual(t, 0, len(b)) {
			assert.Equal(t, "123", b["id"])
		}
		b, err = convertor.ToStringMap(&a)
		if assert.NoError(t, err) && assert.NotEqual(t, 0, len(b)) {
			assert.Equal(t, "123", b["id"])
		}
	}
	{
		a := StringMap{"id": "123"}
		b, err := convertor.ToStringMap(a)
		if assert.NoError(t, err) && assert.NotEqual(t, 0, len(b)) {
			assert.Equal(t, "123", b["id"])
		}
		b, err = convertor.ToStringMap(&a)
		if assert.NoError(t, err) && assert.NotEqual(t, 0, len(b)) {
			assert.Equal(t, "123", b["id"])
		}
	}
	{
		a := CaseInsensitiveStringMap{"id": "123"}
		b, err := convertor.ToStringMap(a)
		if assert.NoError(t, err) && assert.NotEqual(t, 0, len(b)) {
			assert.Equal(t, "123", b["id"])
		}
		b, err = convertor.ToStringMap(&a)
		if assert.NoError(t, err) && assert.NotEqual(t, 0, len(b)) {
			assert.Equal(t, "123", b["id"])
		}
	}
	{
		type A string
		type B string
		a := map[A]B{"id": "111"}
		b, err := convertor.ToStringMap(a)
		if assert.NoError(t, err) && assert.NotEqual(t, 0, len(b)) {
			assert.Equal(t, "111", b["id"])
		}
		b, err = convertor.ToStringMap(&a)
		if assert.NoError(t, err) && assert.NotEqual(t, 0, len(b)) {
			assert.Equal(t, "111", b["id"])
		}
	}
}
