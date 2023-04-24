package convertor

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/modern-go/reflect2"
	"github.com/stretchr/testify/assert"
)

// Person has a Name, Age and Address.
type Person struct {
	Name string
	Age  uint
}

func mustMarshalToString(a interface{}) (s string) {
	s, _ = jsonApi.MarshalToString(a)
	return
}

func TestStringPtrToString(t *testing.T) {
	s := "abc"
	assert.Equal(t, s, ToStringNoError(s))
	assert.Equal(t, s, ToStringNoError(&s))

	type tmp string
	v := tmp(s)
	assert.Equal(t, s, ToStringNoError(v))
}

func TestToString(t *testing.T) {
	cases := []*struct {
		obj interface{}
		str string
	}{
		{obj: nil},
		{obj: (error)(nil)},
		{obj: fmt.Errorf("abc"), str: "abc"},
		{obj: 123, str: "123"},
	}
	for _, c := range cases {
		s, err := ToString(c.obj)
		t.Logf("toStr(%+v): %s", c.obj, s)
		assert.NoError(t, err, c.obj)
		if reflect2.IsNil(c.obj) {
			assert.Equal(t, "", s, c.obj)
		} else if c.str != "" {
			assert.Equal(t, c.str, s, c.obj)
		} else {
			assert.NotEqual(t, "", s, c.obj)
		}
	}
}

func TestToBool(t *testing.T) {
	_, err := ToBool("TrUe")
	if !reflect2.IsNil(err) {
		t.Error(fmt.Sprintf("error occured: err=%v", err))
	}

	for _, obj := range []interface{}{
		1, "1", -1, "-1", "true", "True", "TRUE", "TrUe", "tRuE", "3", 3, "-3", -3, 1.23, "1.23", -1.23, "-1.23",
	} {
		val, err := ToBool(obj)
		if !reflect2.IsNil(err) {
			t.Error(fmt.Sprintf("error occured: err=%v, obj=%s", err, mustMarshalToString(obj)))
		} else {
			if val != true {
				t.Error(fmt.Sprintf("accert faild: val=%v, obj=%s", val, mustMarshalToString(obj)))
			}
		}
	}

	for _, obj := range []interface{}{
		nil, "", 0, "0", "false", "FALSE", "FaLSe", "fAlSe", "-0",
	} {
		val, err := ToBool(obj)
		if !reflect2.IsNil(err) {
			t.Error(fmt.Sprintf("error occured: err=%v, obj=%s", err, mustMarshalToString(obj)))
		} else {
			if val != false {
				t.Error(fmt.Sprintf("accert faild: val=%v, obj=%s", val, mustMarshalToString(obj)))
			}
		}
	}

	for _, obj := range []interface{}{
		"-",
		time.Now(),
		[]int{1}, []interface{}{"a"},
		map[string]interface{}{}, map[string]interface{}{"name": "abc"},
		Person{}, &Person{},
		Person{Name: "ZhangSan"}, &Person{Name: "ZhangSan"},
		Person{Name: "ZhangSan", Age: 24}, &Person{Name: "ZhangSan", Age: 24},
	} {
		val, err := ToBool(obj)
		if reflect2.IsNil(err) {
			t.Error(fmt.Sprintf("accert faild: val=%v, obj=%s", val, mustMarshalToString(obj)))
		}
	}
}

func TestToInt(t *testing.T) {
	for obj, expect := range map[interface{}]int64{
		nil:    0,
		true:   1,
		false:  0,
		0:      0,
		1:      1,
		3.2:    3,
		3.8:    3,
		-1:     -1,
		-3.2:   -3,
		-3.8:   -3,
		"":     0,
		"0":    0,
		"1":    1,
		"3.2":  3,
		"3.8":  3,
		"-0":   0,
		"-1":   -1,
		"-3.2": -3,
		"-3.8": -3,
	} {
		val, err := ToInt64(obj)
		if !reflect2.IsNil(err) {
			t.Error(fmt.Sprintf("error occured: err=%v, obj=%s", err, mustMarshalToString(obj)))
		} else if val != expect {
			t.Error(fmt.Sprintf("accert faild: expect %v, but %v, obj=%s", expect, val, mustMarshalToString(obj)))
		}
	}

	for _, obj := range []interface{}{
		"-",
		time.Now(),
		[]int{1}, []interface{}{"a"},
		map[string]interface{}{}, map[string]interface{}{"name": "abc"},
		Person{}, &Person{},
		Person{Name: "ZhangSan"}, &Person{Name: "ZhangSan"},
		Person{Name: "ZhangSan", Age: 24}, &Person{Name: "ZhangSan", Age: 24},
	} {
		val, err := ToInt64(obj)
		if reflect2.IsNil(err) {
			t.Error(fmt.Sprintf("accert faild: val=%v, item=%s", val, mustMarshalToString(obj)))
		}
	}
}

func TestToUint(t *testing.T) {
	for obj, expect := range map[interface{}]uint64{
		nil:    0,
		true:   1,
		false:  0,
		"":     0,
		0:      0,
		1:      1,
		3.2:    3,
		3.8:    3,
		-1:     math.MaxUint64,
		-3.2:   math.MaxUint64 - 2,
		-3.8:   math.MaxUint64 - 2,
		"0":    0,
		"1":    1,
		"3.2":  3,
		"3.8":  3,
		"-0":   0,
		"-1":   math.MaxUint64,
		"-3.2": math.MaxUint64 - 2,
		"-3.8": math.MaxUint64 - 2,
	} {
		val, err := ToUint64(obj)
		if !reflect2.IsNil(err) {
			t.Error(fmt.Sprintf("error occured: err=%v, obj=%s", err, mustMarshalToString(obj)))
		} else if val != expect {
			t.Error(fmt.Sprintf("accert faild: expect %v, but %v, obj=%s", expect, val, mustMarshalToString(obj)))
		}
	}

	for _, obj := range []interface{}{
		"-",
		time.Now(),
		[]int{1}, []interface{}{"a"},
		map[string]interface{}{}, map[string]interface{}{"name": "abc"},
		Person{}, &Person{},
		Person{Name: "ZhangSan"}, &Person{Name: "ZhangSan"},
		Person{Name: "ZhangSan", Age: 24}, &Person{Name: "ZhangSan", Age: 24},
	} {
		val, err := ToUint64(obj)
		if reflect2.IsNil(err) {
			t.Error(fmt.Sprintf("accert faild: val=%v, item=%s", val, mustMarshalToString(obj)))
		}
	}
}

func TestToFloat(t *testing.T) {
	for obj, expect := range map[interface{}]float64{
		nil:    0,
		true:   1,
		false:  0,
		"":     0,
		0:      0,
		1:      1,
		3.2:    3.2,
		3.8:    3.8,
		-1:     -1,
		-3.2:   -3.2,
		-3.8:   -3.8,
		"0":    0,
		"1":    1,
		"3.2":  3.2,
		"3.8":  3.8,
		"-0":   0,
		"-1":   -1,
		"-3.2": -3.2,
		"-3.8": -3.8,
	} {
		val, err := ToFloat(obj)
		if !reflect2.IsNil(err) {
			t.Error(fmt.Sprintf("error occured: err=%v, obj=%s", err, mustMarshalToString(obj)))
		} else if val != expect {
			t.Error(fmt.Sprintf("accert faild: expect %v, but %v, obj=%s", expect, val, mustMarshalToString(obj)))
		}
	}

	for _, obj := range []interface{}{
		"-",
		time.Now(),
		[]int{1}, []interface{}{"a"},
		map[string]interface{}{}, map[string]interface{}{"name": "abc"},
		Person{}, &Person{},
		Person{Name: "ZhangSan"}, &Person{Name: "ZhangSan"},
		Person{Name: "ZhangSan", Age: 24}, &Person{Name: "ZhangSan", Age: 24},
	} {
		val, err := ToUint64(obj)
		if reflect2.IsNil(err) {
			t.Error(fmt.Sprintf("accert faild: val=%v, item=%s", val, mustMarshalToString(obj)))
		}
	}
}
