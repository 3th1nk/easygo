package nilType

import (
	"github.com/3th1nk/easygo/util/convertor"
	"github.com/3th1nk/easygo/util/jsonUtil"
	"github.com/3th1nk/easygo/util/mapUtil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNullType_1(t *testing.T) {
	type structA struct {
		A Int    `json:"a,omitempty"`
		B Bool   `json:"b,omitempty"`
		C Time   `json:"c,omitempty"`
		D String `json:"d,omitempty"`
	}

	{
		a := &structA{}
		jsonUtil.UnmarshalFromString(`{"a":null,"b":"","d":null}`, a)
		assert.Equal(t, false, a.A.IsNotNil())
		assert.Equal(t, false, a.B.IsNotNil())
		assert.Equal(t, false, a.C.IsNotNil())
		assert.Equal(t, false, a.D.IsNotNil())
		d, _ := convertor.ToStringObjectMap(a)
		assert.Equal(t, 0, len(d))
	}

	{
		a := &structA{}
		jsonUtil.UnmarshalFromString(`{"a":1,"b":"t","c":"2021-01-01","d":""}`, a)
		assert.Equal(t, true, a.A.IsNotNil())
		assert.Equal(t, true, a.B.IsNotNil())
		assert.Equal(t, true, a.C.IsNotNil())
		assert.Equal(t, true, a.D.IsNotNil())
		var d mapUtil.StringObjectMap = convertor.ToStringObjectMapNoError(a)
		assert.Equal(t, 4, len(d))
		assert.Equal(t, 1, d.MustGetInt("a"))
		assert.Equal(t, true, d.MustGetBool("b"))
		assert.Equal(t, "", d.MustGetString("d"))
	}
}
