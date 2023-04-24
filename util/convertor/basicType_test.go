package convertor

import (
	"fmt"
	"github.com/3th1nk/easygo/util"
	"github.com/stretchr/testify/assert"
	"math"
	"strconv"
	"testing"
)

func Test(t *testing.T) {
	{
		a := uint64(math.MaxUint64)
		str := strconv.FormatUint(a, 10)
		b, err := strconv.ParseInt(str, 10, 64)
		assert.Error(t, err)
		assert.NotEqual(t, int64(a), b, str)
		util.Println("%v\n%v", a, b)
	}
	{
		a := int64(-1)
		str := strconv.FormatInt(a, 10)
		b, err := strconv.ParseUint(str, 10, 64)
		assert.Error(t, err)
		assert.NotEqual(t, int64(a), b, str)
		util.Println("%v\n%v", a, b)
	}
}

func TestGetBasicType(t *testing.T) {
	for _, arr := range [][]interface{}{
		{nil},
		{1, BasicType_Int},
		{1.1, BasicType_Float},
		{[]int{1}, BasicType_Slice},
		{[]string{"a", "b", "c"}, BasicType_Slice},
	} {
		n := len(arr)
		val := arr[0]
		basicType, refType, refVal := GetBasicType(val)
		valStr, _ := jsonApi.MarshalToString(val)
		t.Logf("GetBasicType(%v): basicType=%v(%v), refType=%v, refVal=%v]", valStr, basicType, int(basicType), refType, refVal)
		if n != 1 {
			expectType := arr[1].(BasicType)
			if basicType != expectType {
				t.Error(fmt.Errorf("assert faild: expect %v, but %v, val=%v", expectType, basicType, valStr))
			}
		}
	}
}

func TestGetBasicType2(t *testing.T) {
	for _, arr := range [][]interface{}{
		{"1", BasicType_Int},
		{strconv.FormatUint(uint64(math.MaxInt64+10), 10), BasicType_Uint},
		{"1.1", BasicType_Float},
		{"[]", BasicType_Slice},
		{"[1]", BasicType_Slice},
		{"[1,2,3]", BasicType_Slice},
		{`["a"]`, BasicType_Slice},
		{`["a", "b", "c"]`, BasicType_Slice},
		{`1,2,3`, BasicType_String},
		{`{}`, BasicType_Map},
		{`{"id":123}`, BasicType_Map},
	} {
		n := len(arr)
		val := arr[0]
		basicType, refType, refVal := GetBasicType(val, true)
		basicValue := basicType.GetValue(refVal)
		basicValueStr, _ := jsonApi.MarshalToString(basicValue)
		fmt.Println(fmt.Sprintf("GetBasicType(%v): basicType=%v(%v), basicValue=%v: %v", val, basicType, int(basicType), refType, basicValueStr))
		if n != 1 {
			expectType := arr[1].(BasicType)
			if basicType != expectType {
				t.Errorf("basicType assert faild: expect %v, but %v, val=%v", expectType, basicType, val)
			}
		}
	}
}

func TestGetStrValueType(t *testing.T) {
	for _, arr := range [][]interface{}{
		{"1", BasicType_Int},
		{strconv.FormatUint(uint64(math.MaxInt64+10), 10), BasicType_Uint},
		{"1.1", BasicType_Float},
		{"[]", BasicType_Slice},
		{"[1]", BasicType_Slice},
		{"[1,2,3]", BasicType_Slice},
		{`["a"]`, BasicType_Slice},
		{`["a", "b", "c"]`, BasicType_Slice},
		{`1,2,3`, BasicType_String},
		{`{}`, BasicType_Map},
		{`{"id":123}`, BasicType_Map},
		{`["id":123`, BasicType_String},
		{`{"id":123`, BasicType_String},
		{`number`, BasicType_String},
		{`-ddd`, BasicType_String},
		{`false`, BasicType_Bool},
		{`true`, BasicType_Bool},
	} {
		basicType, refType, _, basicValue := GetStrValueType(arr[0].(string))
		basicValueStr, _ := jsonApi.MarshalToString(basicValue)
		fmt.Println(fmt.Sprintf("GetBasicType(%v): basicType=%v(%v), basicValue=%v: %v", arr[0], basicType, int(basicType), refType, basicValueStr))
		expectType := arr[1].(BasicType)
		if basicType != expectType {
			t.Errorf("basicType assert faild: expect %v, but %v, val=%v", expectType, basicType, arr[0])
		}
	}
}
