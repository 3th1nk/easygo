package jsonValidator

import (
	"fmt"
	"github.com/3th1nk/easygo/util/comparer"
	"github.com/3th1nk/easygo/util/convertor"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidator_GetValue_0(t *testing.T) {
	validator := newValidator(`
[
	{"id": 1, "name": "a", "data": {"str": "aa"}},
	{"id": 2, "name": "b", "data": {"str": "bb"}},
	{"id": 3, "name": "c", "data": {"str": "cc"}},
	{"id": 4, "name": "d", "data": {"str": "dd"}},
	{"id": 5, "name": "e", "data": {"str": "ee"}}
]
`)

	val, err := validator.StringSlice("*.data.str")
	assert.NoError(t, err)
	assert.Equal(t, convertor.BasicType_Slice, convertor.GetBasicType0(val))
	assert.Equal(t, 5, len(val))
	assert.Equal(t, "aa", convertor.ToStringNoError(val[0]))
	assert.Equal(t, "bb", convertor.ToStringNoError(val[1]))
	assert.Equal(t, "cc", convertor.ToStringNoError(val[2]))
	assert.Equal(t, "dd", convertor.ToStringNoError(val[3]))
	assert.Equal(t, "ee", convertor.ToStringNoError(val[4]))
}

func TestValidator_GetValue_1(t *testing.T) {
}

func TestValidator_GetValue_2(t *testing.T) {
	{
		val, err := newValidator().Value("code")
		assert.NoError(t, err)
		assert.Equal(t, 1, convertor.ToIntNoError(val))
	}

	{
		val, err := newValidator().String("message")
		assert.NoError(t, err)
		assert.Equal(t, "test", val)
	}

	{
		val, err := newValidator().Value("data")
		assert.NoError(t, err)
		assert.Equal(t, convertor.BasicType_Map, convertor.GetBasicType0(val))
	}

	{
		val, err := newValidator().StringObjectMap("data.user")
		assert.NoError(t, err)
		assert.Equal(t, 1, val.MustGetInt("id"))
		assert.Equal(t, "test", val.MustGetString("name"))
	}

	{
		val, err := newValidator().IntSlice("data.id")
		assert.NoError(t, err)
		if assert.Equal(t, 3, len(val)) {
			assert.Equal(t, 1, convertor.ToIntNoError(val[0]))
			assert.Equal(t, 2, convertor.ToIntNoError(val[1]))
			assert.Equal(t, 3, convertor.ToIntNoError(val[2]))
		}
	}

	{
		val, err := newValidator().Value("data.id.0")
		assert.NoError(t, err)
		assert.Equal(t, convertor.BasicType_Int, convertor.GetBasicType0(val))
		assert.Equal(t, 1, convertor.ToIntNoError(val))
	}

	{
		val, err := newValidator().Value("data.id.1")
		assert.NoError(t, err)
		assert.Equal(t, convertor.BasicType_Int, convertor.GetBasicType0(val))
		assert.Equal(t, 2, convertor.ToIntNoError(val))
	}

	{
		val, err := newValidator().Value("data.id.2")
		assert.NoError(t, err)
		assert.Equal(t, convertor.BasicType_Int, convertor.GetBasicType0(val))
		assert.Equal(t, 3, convertor.ToIntNoError(val))
	}

	{
		val, err := newValidator().IntSlice("data.items.*.id")
		assert.NoError(t, err)
		if assert.Equal(t, 5, len(val)) {
			assert.Equal(t, 1, convertor.ToIntNoError(val[0]))
			assert.Equal(t, 2, convertor.ToIntNoError(val[1]))
			assert.Equal(t, 3, convertor.ToIntNoError(val[2]))
			assert.Equal(t, 4, convertor.ToIntNoError(val[3]))
			assert.Equal(t, 5, convertor.ToIntNoError(val[4]))
		}
	}

	{
		val, err := newValidator().StringSlice("data.items.*.name")
		assert.NoError(t, err)
		if assert.Equal(t, 5, len(val)) {
			assert.Equal(t, "a", convertor.ToStringNoError(val[0]))
			assert.Equal(t, "b", convertor.ToStringNoError(val[1]))
			assert.Equal(t, "c", convertor.ToStringNoError(val[2]))
			assert.Equal(t, "d", convertor.ToStringNoError(val[3]))
			assert.Equal(t, "e", convertor.ToStringNoError(val[4]))
		}
	}
}

func TestValidator_Validate_1(t *testing.T) {
	validator := newValidator()

	for _, arr := range [][]interface{}{
		{"code", "eq", 1},
		{"code", "eq", "1"},
		{"code", "eq", "1.0"},
		{"code", "ueq", 1, false},
		{"code", "egt", 1},
		{"code", "elt", 1},
		{"code", "gt", 1, false},
		{"code", "lt", 1, false},
		{"code", "in", []int{1, 2}},
		{"code", "not-in", []int{2, 3}},
		{"code", "in", []int{2, 3}, false},
		{"code", "not-in", []int{1, 2}, false},
	} {
		path, opr, val, expectOk := arr[0].(string), comparer.Operator(arr[1].(string)), arr[2], true
		if len(arr) > 3 {
			expectOk = arr[3].(bool)
		}

		ok, err := validator.Validate(path, opr, val, 0)
		if err != nil {
			t.Error(fmt.Errorf("error occured: %v，path=%v, opr=%v, val=%v", err, path, opr, convertor.ToStringNoError(val)))
		} else if ok != expectOk {
			t.Error(fmt.Errorf("assert faild: expect %v, but %v，path=%v, opr=%v, val=%v", expectOk, ok, path, opr, convertor.ToStringNoError(val)))
		}
	}
}

func TestValidator_Validate_2(t *testing.T) {
	validator := newValidator()

	for _, arr := range [][]interface{}{
		{"data.userId", "ueq", 0},
		{"data.userId", "gt", 0},
		{"data.userId", "in", []int{1, 2, 3}},
	} {
		path, opr, val, expectOk := arr[0].(string), comparer.Operator(arr[1].(string)), arr[2], true
		if len(arr) > 3 {
			expectOk = arr[3].(bool)
		}

		ok, err := validator.Validate(path, opr, val, 0)
		if err != nil {
			t.Error(fmt.Errorf("error occured: %v，path=%v, opr=%v, val=%v", err, path, opr, convertor.ToStringNoError(val)))
		} else if ok != expectOk {
			t.Error(fmt.Errorf("assert faild: expect %v, but %v，path=%v, opr=%v, val=%v", expectOk, ok, path, opr, convertor.ToStringNoError(val)))
		}
	}
}

func TestValidator_Validate_3(t *testing.T) {
	validator := newValidator()

	for _, arr := range [][]interface{}{
		{"data.id.0", "exist", nil},
		{"data.id.3", "exist", nil, false},
	} {
		path, opr, val, expectOk := arr[0].(string), comparer.Operator(arr[1].(string)), arr[2], true
		if len(arr) > 3 {
			expectOk = arr[3].(bool)
		}

		ok, _ := validator.Validate(path, opr, val, 0)
		if ok != expectOk {
			t.Error(fmt.Errorf("assert faild: expect %v, but %v，path=%v, opr=%v, val=%v", expectOk, ok, path, opr, convertor.ToStringNoError(val)))
		}
	}
}

func TestValidator_Validate_4(t *testing.T) {
	validator := newValidator()

	for _, arr := range [][]interface{}{
		{"data.id.0", "eq", 1},
		{"data.id.1", "eq", 2},
	} {
		path, opr, val, expectOk := arr[0].(string), comparer.Operator(arr[1].(string)), arr[2], true
		if len(arr) > 3 {
			expectOk = arr[3].(bool)
		}

		ok, err := validator.Validate(path, opr, val, 0)
		if err != nil {
			t.Error(fmt.Errorf("error occured: %v，path=%v, opr=%v, val=%v", err, path, opr, convertor.ToStringNoError(val)))
		} else if ok != expectOk {
			t.Error(fmt.Errorf("assert faild: expect %v, but %v，path=%v, opr=%v, val=%v", expectOk, ok, path, opr, convertor.ToStringNoError(val)))
		}
	}
}

func TestValidator_Validate_5(t *testing.T) {
	validator := newValidator()

	for _, arr := range [][]interface{}{
		{"data.id", "contains", 1},
		{"data.id", "contains", 2},
		{"data.id", "contains", 3},
	} {
		path, opr, val, expectOk := arr[0].(string), comparer.Operator(arr[1].(string)), arr[2], true
		if len(arr) > 3 {
			expectOk = arr[3].(bool)
		}

		ok, err := validator.Validate(path, opr, val, 0)
		if err != nil {
			t.Error(fmt.Errorf("error occured: %v，path=%v, opr=%v, val=%v", err, path, opr, convertor.ToStringNoError(val)))
		} else if ok != expectOk {
			t.Error(fmt.Errorf("assert faild: expect %v, but %v，path=%v, opr=%v, val=%v", expectOk, ok, path, opr, convertor.ToStringNoError(val)))
		}
	}
}

func TestValidator_Validate_Error(t *testing.T) {
	validator := newValidator()

	for _, arr := range [][]interface{}{
		{"code", "eq", []int{0}},
		{"code", "asfdasdf", 0},
		{"code", "contains", 0},
		{"code", "contains", []int{0, 304}},
		{"code", "not-in", 0},
	} {
		path, opr, val := arr[0].(string), comparer.Operator(arr[1].(string)), arr[2]

		ok, err := validator.Validate(path, opr, val, 0)
		if err == nil {
			t.Error(fmt.Errorf("should have error: path=%v, opr=%v, val=%v, ok=%v", path, opr, convertor.ToStringNoError(val), ok))
		}
	}
}

func TestValidator_GetMapValue(t *testing.T) {
	validator := New(map[string]interface{}{
		"device": map[string]interface{}{
			"brand":  "cisco",
			"module": "WS-C3750X-48U-E",
			"abc":    nil,
		},
	})

	if v, err := validator.Value("device"); err != nil {
		t.Error(err)
	} else {
		t.Logf("device: %+v", v)
	}

	if v, err := validator.Value("device.brand"); err != nil {
		t.Error(err)
	} else {
		t.Logf("device.brand: %+v", v)
	}

	if v, err := validator.Value("device.abc"); err != nil {
		t.Error(err)
	} else {
		t.Logf("device.abc: %+v", v)
	}
}

func newValidator(str ...string) *Validator {
	var theStr string
	if len(str) != 0 {
		theStr = str[0]
	}
	if theStr == "" {
		theStr = `
{
    "code": 1,
    "message": "test",
    "data": {
        "userId": 1,
		"user": {"id": 1, "name": "test"},
        "id": [1,2,3],
        "items": [
            {"id": 1, "name": "a"},
            {"id": 2, "name": "b"},
            {"id": 3, "name": "c"},
            {"id": 4, "name": "d"},
            {"id": 5, "name": "e"}
        ]
    }
}
`
	}
	return New(theStr)
}

func ExampleNew() {
	obj := New(`
	 {
	     "code": 1,
	     "message": "test",
	     "data": {
	         "userId": 1,
	 		 "user": {"id": 1, "name": "test"},
	         "id": [1,2,3],
	         "items": [
	             {"id": 1, "name": "a"},
	             {"id": 2, "name": "b"},
	             {"id": 3, "name": "c"},
	             {"id": 4, "name": "d"},
	             {"id": 5, "name": "e"}
	         ]
	     }
	}`)

	// float64(1)
	val, _ := obj.Value("code")

	// int(1)
	nVal, _ := obj.Int("code")

	// string("test")
	val, _ = obj.Value("message")

	// map[string]interface{}{...}
	val, _ = obj.Value("data")

	// map[string]interface{}{"id": 1, "name": "test"}
	val, _ = obj.Value("data.user")

	// []int64{1,2,3}
	val, _ = obj.Value("data.id")

	// []int{1,2,3}
	arrVal, _ := obj.IntSlice("data.id")

	// int64(2)
	val, _ = obj.Value("data.id.1")

	// []int64{1,2,3,4,5}
	val, _ = obj.Value("data.items.*.id")

	// []string{"a","b","c","d","e"}
	val, _ = obj.Value("data.items.*.name")

	print(val, nVal, arrVal)
}
