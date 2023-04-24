package arrUtil

import (
	jsonIter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Unmarshal_CommaSepInt_1(t *testing.T) {
	var a CommaSepInt
	err := jsonIter.UnmarshalFromString(`null`, &a)
	assert.NoError(t, err)
	assert.Nil(t, a)
}

func Test_Unmarshal_CommaSepInt_2(t *testing.T) {
	var a CommaSepInt
	err := jsonIter.UnmarshalFromString(``, &a)
	assert.NoError(t, err)
	assert.Nil(t, a)
}

func Test_Unmarshal_CommaSepInt_3(t *testing.T) {
	var a CommaSepInt
	err := jsonIter.UnmarshalFromString(`[1,2,3]`, &a)
	assert.NoError(t, err)
	if assert.Equal(t, 3, len(a)) {
		assert.Equal(t, 1, a[0])
		assert.Equal(t, 2, a[1])
		assert.Equal(t, 3, a[2])
	}
	str, _ := jsonIter.MarshalToString(a)
	assert.Equal(t, `"1,2,3"`, str)
}

func Test_Unmarshal_CommaSepInt_4(t *testing.T) {
	var a CommaSepInt
	err := jsonIter.UnmarshalFromString(`"1,2,3"`, &a)
	assert.NoError(t, err)
	if assert.Equal(t, 3, len(a)) {
		assert.Equal(t, 1, a[0])
		assert.Equal(t, 2, a[1])
		assert.Equal(t, 3, a[2])
	}
	str, _ := jsonIter.MarshalToString(a)
	assert.Equal(t, `"1,2,3"`, str)
}

func Test_Unmarshal_CommaSepInt_5(t *testing.T) {
	type tmp struct {
		Ids CommaSepInt `json:"ids,omitempty"`
	}

	{
		var a *tmp
		err := jsonIter.UnmarshalFromString(`{}`, &a)
		if assert.NoError(t, err) && !assert.NotNil(t, a) {
			assert.Equal(t, 0, len(a.Ids))
		}
		str, _ := jsonIter.MarshalToString(a)
		assert.Equal(t, `{}`, str)
	}

	{
		var a *tmp
		err := jsonIter.UnmarshalFromString(`{"ids":[1,2,3]}`, &a)
		if assert.NoError(t, err) && !assert.NotNil(t, a) {
			if assert.Equal(t, 3, len(a.Ids)) {
				assert.Equal(t, 1, a.Ids[0])
				assert.Equal(t, 2, a.Ids[1])
				assert.Equal(t, 3, a.Ids[2])
			}
		}
		str, _ := jsonIter.MarshalToString(a)
		assert.Equal(t, `{"ids":"1,2,3"}`, str)
	}

	{
		var a *tmp
		err := jsonIter.UnmarshalFromString(`{"ids":null}`, &a)
		if assert.NoError(t, err) && !assert.NotNil(t, a) {
			assert.Equal(t, 0, len(a.Ids))
		}
		str, _ := jsonIter.MarshalToString(a)
		assert.Equal(t, `{}`, str)
	}

	{
		var a *tmp
		err := jsonIter.UnmarshalFromString(`{"ids":[1,2,3]}`, &a)
		if assert.NoError(t, err) && !assert.NotNil(t, a) {
			if assert.Equal(t, 3, len(a.Ids)) {
				assert.Equal(t, 1, a.Ids[0])
				assert.Equal(t, 2, a.Ids[1])
				assert.Equal(t, 3, a.Ids[2])
			}
		}
		str, _ := jsonIter.MarshalToString(a)
		assert.Equal(t, `{"ids":"1,2,3"}`, str)
	}

	{
		var a *tmp
		err := jsonIter.UnmarshalFromString(`{"ids":"1,2,3"}`, &a)
		if assert.NoError(t, err) && !assert.NotNil(t, a) {
			if assert.Equal(t, 3, len(a.Ids)) {
				assert.Equal(t, 1, a.Ids[0])
				assert.Equal(t, 2, a.Ids[1])
				assert.Equal(t, 3, a.Ids[2])
			}
		}
		str, _ := jsonIter.MarshalToString(a)
		assert.Equal(t, `{"ids":"1,2,3"}`, str)
	}
}

func Test_Unmarshal_CommaSepString(t *testing.T) {
	type tmp struct {
		Ids CommaSepString `json:"ids,omitempty"`
	}

	{
		var a *tmp
		err := jsonIter.UnmarshalFromString(`{}`, &a)
		if assert.NoError(t, err) && !assert.NotNil(t, a) {
			assert.Equal(t, 0, len(a.Ids))
		}
		str, _ := jsonIter.MarshalToString(a)
		assert.Equal(t, `{}`, str)
	}

	{
		var a *tmp
		err := jsonIter.UnmarshalFromString(`{"ids":null}`, &a)
		if assert.NoError(t, err) && !assert.NotNil(t, a) {
			assert.Equal(t, 0, len(a.Ids))
		}
		str, _ := jsonIter.MarshalToString(a)
		assert.Equal(t, `{}`, str)
	}

	{
		var a *tmp
		err := jsonIter.UnmarshalFromString(`{"ids":[]}`, &a)
		if assert.NoError(t, err) && !assert.NotNil(t, a) {
			assert.Equal(t, 0, len(a.Ids))
		}
		str, _ := jsonIter.MarshalToString(a)
		assert.Equal(t, `{}`, str)
	}

	{
		var a *tmp
		err := jsonIter.UnmarshalFromString(`{"ids":["1","2","3"]}`, &a)
		if assert.NoError(t, err) && !assert.NotNil(t, a) {
			if assert.Equal(t, 3, len(a.Ids)) {
				assert.Equal(t, "1", a.Ids[0])
				assert.Equal(t, "2", a.Ids[1])
				assert.Equal(t, "3", a.Ids[2])
			}
		}
		str, _ := jsonIter.MarshalToString(a)
		assert.Equal(t, `{"ids":"1,2,3"}`, str)
	}

	{
		var a *tmp
		err := jsonIter.UnmarshalFromString(`{"ids":["1","2","3"]}`, &a)
		if assert.NoError(t, err) && !assert.NotNil(t, a) {
			if assert.Equal(t, 3, len(a.Ids)) {
				assert.Equal(t, "1", a.Ids[0])
				assert.Equal(t, "2", a.Ids[1])
				assert.Equal(t, "3", a.Ids[2])
			}
		}
		str, _ := jsonIter.MarshalToString(a)
		assert.Equal(t, `{"ids":"1,2,3"}`, str)
	}

	{
		var a *tmp
		err := jsonIter.UnmarshalFromString(`{"ids":"1,2,3"}`, &a)
		if assert.NoError(t, err) && !assert.NotNil(t, a) {
			if assert.Equal(t, 3, len(a.Ids)) {
				assert.Equal(t, "1", a.Ids[0])
				assert.Equal(t, "2", a.Ids[1])
				assert.Equal(t, "3", a.Ids[2])
			}
		}
		str, _ := jsonIter.MarshalToString(a)
		assert.Equal(t, `{"ids":"1,2,3"}`, str)
	}

	{
		var a *tmp
		err := jsonIter.UnmarshalFromString(`{"ids":"1"}`, &a)
		if assert.NoError(t, err) && !assert.NotNil(t, a) {
			if assert.Equal(t, 3, len(a.Ids)) {
				assert.Equal(t, "1", a.Ids[0])
				assert.Equal(t, "2", a.Ids[1])
				assert.Equal(t, "3", a.Ids[2])
			}
		}
		str, _ := jsonIter.MarshalToString(a)
		assert.Equal(t, `{"ids":"1"}`, str)
	}
}

func Test(t *testing.T) {
	type tmp struct {
		Ids CommaSepInt `json:"ids,omitempty"`
	}
	var a *tmp
	err := jsonIter.UnmarshalFromString(`{"ids":"1"}`, &a)
	if assert.NoError(t, err) && !assert.NotNil(t, a) {
		if assert.Equal(t, 3, len(a.Ids)) {
			assert.Equal(t, "1", a.Ids[0])
			assert.Equal(t, "2", a.Ids[1])
			assert.Equal(t, "3", a.Ids[2])
		}
	}
	str, _ := jsonIter.MarshalToString(a)
	assert.Equal(t, `{"ids":"1"}`, str)
}
