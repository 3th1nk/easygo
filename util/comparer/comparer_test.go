package comparer

import (
	"fmt"
	"github.com/3th1nk/easygo/util/convertor"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompare(t *testing.T) {
	expectFunc := func(a, b interface{}, opr Operator, expectVal, expectErr bool) {
		val, err := Compare(a, b, opr, Option_CaseInsensitive)
		if (err != nil) != expectErr {
			if expectErr {
				t.Error(fmt.Sprintf("should has error: a=%v, b=%v, opr=%v", a, b, opr))
			} else {
				t.Error(fmt.Sprintf("error occured: a=%v, b=%v, opr=%v, err=%v", a, b, opr, err))
			}
		} else if expectVal != val {
			t.Error(fmt.Sprintf("assert faild: expect %v, but %v, a=%v, b=%v, opr=%v", expectVal, val, a, b, opr))
		}
	}

	for _, arr := range [][]interface{}{
		{true, true, "=", true},
		{true, true, "!=", false},

		{true, false, Operator_Eq, false},
		{true, false, Operator_Ueq, true},

		{false, true, Operator_Eq, false},
		{false, true, Operator_Ueq, true},

		{false, false, Operator_Eq, true},
		{false, false, Operator_Ueq, false},

		{true, true, ">", false, true},
		{true, true, ">=", false, true},
		{true, true, Operator_Lt, false, true},
		{true, true, Operator_Elt, false, true},
		{true, true, Operator_In, false, true},
		{true, true, "!in", false, true},
		{true, true, Operator_Contains, false, true},
		{true, true, Operator_NotContains, false, true},

		{1, 1, Operator_In, false, true},
		{1, 1, Operator_NotIn, false, true},
		{1, 1, Operator_Contains, false, true},
		{1, 1, Operator_NotContains, false, true},
		{3, "abc", Operator_Eq, false, true},
		{3, "abc", Operator_Ueq, false, true},
		{1, "abc", Operator_Gt, false, true},
		{1, "abc", Operator_Egt, false, true},
		{1, "abc", Operator_Lt, false, true},
		{1, "abc", Operator_Elt, false, true},

		{1, 1, Operator_Eq, true},
		{1, true, Operator_Eq, true},
		{1, "true", Operator_Eq, true},
		{1, "t", Operator_Eq, true},
		{1, "1.2", Operator_Eq, true},
		{1, "1.8", Operator_Eq, true},
		{0, "", Operator_Eq, true},
		{0, false, Operator_Eq, true},
		{0, nil, Operator_Eq, true},

		{1, 2, Operator_Ueq, true},
		{1, "", Operator_Ueq, true},
		{1, nil, Operator_Ueq, true},

		{3, 1, Operator_Gt, true},
		{3, "", Operator_Gt, true},
		{3, false, Operator_Gt, true},
		{3, "false", Operator_Gt, true},
	} {
		a, b, opr, expectVal, expectErr := arr[0], arr[1], Operator(fmt.Sprintf("%v", arr[2])), true, false
		if n := len(arr); n > 3 {
			expectVal = arr[3].(bool)
			if n > 4 {
				expectErr = arr[4].(bool)
			}
		}
		expectFunc(a, b, opr, expectVal, expectErr)
	}
}

func TestCompare2(t *testing.T) {
	val, err := Compare(3, "false", Operator_Gt, 0)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(val)
	}
}

func TestCompareBool(t *testing.T) {
	for _, arr := range [][]interface{}{
		{true, true, Operator_Eq, true},
		{true, true, Operator_Ueq, false},

		{true, false, Operator_Eq, false},
		{true, false, Operator_Ueq, true},

		{false, true, Operator_Eq, false},
		{false, true, Operator_Ueq, true},

		{false, false, Operator_Eq, true},
		{false, false, Operator_Ueq, false},

		{true, true, Operator_Gt},
		{true, true, Operator_Egt},
		{true, true, Operator_Lt},
		{true, true, Operator_Elt},
		{true, true, Operator_In},
		{true, true, Operator_NotIn},
		{true, true, Operator_Contains},
		{true, true, Operator_NotContains},
	} {
		a, b, opr, expectErr, expectVal := arr[0].(bool), arr[1].(bool), arr[2].(Operator), len(arr) == 3, false
		if !expectErr {
			expectVal = arr[3].(bool)
		}

		expectFunc := func(a, b bool, opr Operator, expectVal bool) {
			val, err := CompareBool(a, b, opr)
			if (err != nil) != expectErr {
				if expectErr {
					t.Error(fmt.Sprintf("should has error: a=%v, b=%v, opr=%v", a, b, opr))
				} else {
					t.Error(fmt.Sprintf("error occured: a=%v, b=%v, opr=%v, err=%v", a, b, opr, err))
				}
			}
			if !expectErr && expectVal != val {
				t.Error(fmt.Sprintf("assert faild: expect %v, but %v, a=%v, b=%v, opr=%v", expectVal, val, a, b, opr))
			}
		}
		expectFunc(a, b, opr, expectVal)

		if expectErr {
			continue
		}

		if opr == Operator_Eq || opr == Operator_Ueq {
			expectFunc(b, a, opr, expectVal)
			if opr == Operator_Eq {
				expectFunc(a, b, Operator_Ueq, !expectVal)
				expectFunc(b, a, Operator_Ueq, !expectVal)
			} else {
				expectFunc(a, b, Operator_Eq, !expectVal)
				expectFunc(b, a, Operator_Eq, !expectVal)
			}
		}
	}
}

func TestCompareInt(t *testing.T) {
	for _, arr := range [][]interface{}{
		{1, 1, Operator_In},
		{1, 1, Operator_NotIn},
		{1, 1, Operator_Contains},
		{1, 1, Operator_NotContains},
		{3, "abc", Operator_Eq},
		{3, "abc", Operator_Ueq},
		{1, "abc", Operator_Gt},
		{1, "abc", Operator_Egt},
		{1, "abc", Operator_Lt},
		{1, "abc", Operator_Elt},

		{1, 1, Operator_Eq, true},
		{1, true, Operator_Eq, true},
		{1, "true", Operator_Eq, true},
		{1, "t", Operator_Eq, true},
		{1, "1.2", Operator_Eq, true},
		{1, "1.8", Operator_Eq, true},
		{0, "", Operator_Eq, true},
		{0, false, Operator_Eq, true},
		{0, nil, Operator_Eq, true},

		{1, 2, Operator_Ueq, true},
		{1, "", Operator_Ueq, true},
		{1, nil, Operator_Ueq, true},

		{3, 1, Operator_Gt, true},
		{3, "", Operator_Gt, true},
		{3, false, Operator_Gt, true},
		{3, "false", Operator_Gt, true},
	} {
		opr, expectErr, expectVal, convertErr := arr[2].(Operator), len(arr) == 3, false, error(nil)
		if !expectErr {
			expectVal = arr[3].(bool)
		}
		a, err := convertor.ToInt64(arr[0])
		if err != nil {
			a, expectErr, convertErr = 0, true, err
		}
		b, err := convertor.ToInt64(arr[1])
		if err != nil {
			b, expectErr, convertErr = 0, true, err
		}
		if convertErr != nil {
			if !expectErr {
				t.Error(fmt.Sprintf("convert error: %v, a=%v, b=%v", err, a, b))
			}
			continue
		}

		expectFunc := func(a, b int64, opr Operator, expectVal bool) {
			val, err := CompareInt(a, b, opr)
			if (err != nil) != expectErr {
				if expectErr {
					t.Error(fmt.Sprintf("should has error: a=%v, b=%v, opr=%v", a, b, opr))
				} else {
					t.Error(fmt.Sprintf("error occured: a=%v, b=%v, opr=%v, err=%v", a, b, opr, err))
				}
			} else if expectVal != val {
				t.Error(fmt.Sprintf("assert faild: expect %v, but %v, a=%v, b=%v, opr=%v", expectVal, val, a, b, opr))
			}
		}
		expectFunc(a, b, opr, expectVal)

		if expectErr {
			continue
		}

		if opr == Operator_Eq || opr == Operator_Ueq {
			expectFunc(b, a, opr, expectVal)
			if opr == Operator_Eq {
				expectFunc(a, b, Operator_Ueq, !expectVal)
				expectFunc(b, a, Operator_Ueq, !expectVal)
			} else {
				expectFunc(a, b, Operator_Eq, !expectVal)
				expectFunc(b, a, Operator_Eq, !expectVal)
			}
		}
		if opr == Operator_Eq {
			expectFunc(a, b, Operator_Egt, expectVal)
			expectFunc(a, b, Operator_Elt, expectVal)
			expectFunc(b, a, Operator_Egt, expectVal)
			expectFunc(b, a, Operator_Elt, expectVal)
		}
		if opr == Operator_Gt {
			expectFunc(b, a, Operator_Lt, expectVal)

			expectFunc(a, b, Operator_Egt, expectVal)
			expectFunc(a, b, Operator_Lt, !expectVal)
			expectFunc(a, b, Operator_Elt, !expectVal)

			expectFunc(b, a, Operator_Elt, expectVal)
			expectFunc(b, a, Operator_Gt, !expectVal)
			expectFunc(b, a, Operator_Egt, !expectVal)
		}
		if opr == Operator_Lt {
			expectFunc(b, a, Operator_Gt, expectVal)

			expectFunc(a, b, Operator_Elt, expectVal)
			expectFunc(a, b, Operator_Gt, !expectVal)
			expectFunc(a, b, Operator_Egt, !expectVal)

			expectFunc(b, a, Operator_Egt, expectVal)
			expectFunc(b, a, Operator_Lt, !expectVal)
			expectFunc(b, a, Operator_Elt, !expectVal)
		}
	}
}

func TestCompareString(t *testing.T) {
	ok, err := false, error(nil)

	ok, err = Compare("abc", "a", Operator_Gt)
	if assert.NoError(t, err) {
		assert.Equal(t, true, ok)
	}

	ok, err = Compare("abc", "a", Operator_Contains)
	if assert.NoError(t, err) {
		assert.Equal(t, true, ok)
	}

	ok, err = Compare("abc", "ab", Operator_Like)
	if assert.NoError(t, err) {
		assert.Equal(t, true, ok)
	}

	ok, err = Compare("abc", "AB", Operator_Like, 2)
	if assert.NoError(t, err) {
		assert.Equal(t, true, ok)
	}

	ok, err = Compare("abc", "abcabc", Operator_In)
	if assert.NoError(t, err) {
		assert.Equal(t, true, ok)
	}
}

func TestCompareString_Regex(t *testing.T) {
	ok, err := Compare("aaa,bbb,ccc,name=tom and jerry", `name\s*=\s*.*?`, Operator_Regex)
	if assert.NoError(t, err) {
		assert.Equal(t, true, ok)
	}

	ok, err = Compare("aaa,bbb,ccc,name=tom and jerry", `names\s*=\s*.*?`, Operator_Regex)
	if assert.NoError(t, err) {
		assert.Equal(t, false, ok)
	}
}
