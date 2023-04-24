package dynamicExpr

import (
	"github.com/3th1nk/easygo/expr/jsonValidator"
	"github.com/3th1nk/easygo/util/convertor"
	"github.com/3th1nk/easygo/util/jsonUtil"
	"github.com/3th1nk/easygo/util/mapUtil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Format(t *testing.T) {
	provider := jsonValidator.New(map[string]interface{}{
		"int": 1,
		"str": "abc",
		"obj": map[string]interface{}{
			"name":  "abc",
			"value": 123,
		},
	})

	formatter := _default
	var str string
	var err error
	var dict mapUtil.StringObjectMap

	str, err = formatter.Format(`{"id":{$obj.value},"name":{$obj.name},"str":{$str}}`, provider)
	if assert.NoError(t, err) {
		assert.Equal(t, `{"id":123,"name":"abc","str":"abc"}`, str)
	}
	str, err = formatter.Format(`{"id":$obj.value,"name":$obj.name,"str":$str}`, provider)
	if assert.NoError(t, err) {
		assert.Equal(t, `{"id":123,"name":"abc","str":"abc"}`, str)
	}

	str, err = formatter.Format(`{"{$obj.name}":{$obj.value}}`, provider)
	if assert.NoError(t, err) {
		assert.Equal(t, `{"abc":123}`, str)
	}
	str, err = formatter.Format(`{"$obj.name":$obj.value}`, provider)
	if assert.NoError(t, err) {
		assert.Equal(t, `{"abc":123}`, str)
	}

	str, err = formatter.Format(`{"$str":{$obj.value}}`, provider)
	if assert.NoError(t, err) {
		assert.Equal(t, `{"abc":123}`, str)
	}
	str, err = formatter.Format(`{"$str":$obj.value}`, provider)
	if assert.NoError(t, err) {
		assert.Equal(t, `{"abc":123}`, str)
	}

	str, err = formatter.Format("/api/device-info/{$obj.value}", provider)
	if err != nil || str != "/api/device-info/123" {
		t.Errorf("assert faild: %v, %v", str, err)
	}

	str, err = formatter.Format(`{
		"id": {$obj.value },
		"name": {$obj.name },
		"obj": {$obj},
		"objStr": "{$obj}",
		"id2": "{$obj.value}",
		"name2": "{$obj.name}",
		"id_3": "{$id-not-exist}"
	}`, provider)
	if !assert.NoError(t, err) {
		return
	}
	err = jsonUtil.UnmarshalFromString(str, &dict)
	if !assert.NoError(t, err, str) {
		return
	}

	assert.Equal(t, 123, dict.MustGetInt("id"))
	assert.Equal(t, "123", dict.MustGet("id2"))
	assert.Equal(t, "abc", dict.MustGet("name"))
	assert.Equal(t, "abc", dict.MustGet("name2"))
	assert.NotEqual(t, "", dict.MustGet("objStr"))

	var obj mapUtil.StringObjectMap
	_, err = dict.GetToObject("obj", &obj)
	if assert.NoError(t, err) && assert.NotNil(t, obj) {
		assert.Equal(t, 123, obj.MustGetInt("value"))
		assert.Equal(t, "abc", obj.MustGetString("name"))
	}

	str, err = formatter.Format("$int", provider)
	assert.NoError(t, err)
	assert.Equal(t, "1", str)

	str, err = formatter.Format("$str", provider)
	assert.NoError(t, err)
	assert.Equal(t, "abc", str)

	str, err = formatter.Format("$admin", provider)
	assert.Error(t, err)

	str, err = formatter.Format("test-{$int}", provider)
	assert.NoError(t, err)
	assert.Equal(t, "test-1", str)

	str, err = formatter.Format("test-{$str}", provider)
	assert.NoError(t, err)
	assert.Equal(t, "test-abc", str)

	str, err = formatter.Format("test-{$abc}", provider)
	assert.NoError(t, err)
	assert.Equal(t, "test-{$abc}", str)

	str, err = formatter.Format("test-{$int},{$str},test", provider)
	assert.NoError(t, err)
	assert.Equal(t, "test-1,abc,test", str)

	str, err = formatter.Format("test-{$int},{$str},test", provider)
	assert.NoError(t, err)
	assert.Equal(t, "test-1,abc,test", str)

	str, err = formatter.Format("$obj", provider)
	assert.NoError(t, err)
	obj = mapUtil.StringObjectMap{}
	err = jsonUtil.UnmarshalFromString(str, &obj)
	if assert.NoError(t, err) && assert.NotNil(t, obj) {
		assert.Equal(t, 123, obj.MustGetInt("value"))
		assert.Equal(t, "abc", obj.MustGetString("name"))
	}

	str, err = formatter.Format(`{"objStr": "$obj"}`, provider)
	assert.NoError(t, err)
	dict, err = convertor.ToStringObjectMap(str)
	assert.NoError(t, err)
	assert.NotEmpty(t, dict["objStr"])
}

func Test_FormatList(t *testing.T) {
	provider := jsonValidator.New(map[string]interface{}{
		"int": 1,
		"str": "abc",
		"list": []map[string]interface{}{
			{"id": 1, "name": "a"},
			{"id": 2, "name": "b"},
			{"id": 3, "name": "c"},
		},
	})

	formatter := _default
	var str string
	var err error

	str, err = formatter.Format(`{"id":{$list.*.id},"name":{$list.*.name}}`, provider)
	if assert.NoError(t, err) {
		assert.Equal(t, `{"id":[1,2,3],"name":["a","b","c"]}`, str)
	}
}

func Test_FormatUnknownExpr(t *testing.T) {
	provider := jsonValidator.New(map[string]interface{}{
		"int": 1,
		"str": "abc",
		"obj": map[string]interface{}{
			"name":  "abc",
			"value": 123,
		},
	})

	formatter := _default
	var str string
	var err error

	str, err = formatter.Format(`{"$strstr":{$obj.value}}`, provider)
	if assert.NoError(t, err) {
		assert.Equal(t, `{"$strstr":123}`, str)
	}

	str, err = formatter.Format(`$(($used * 100 / $total))`, provider)
	if assert.NoError(t, err) {
		assert.Equal(t, `$(($used * 100 / $total))`, str)
	}
}

func Test_BraceStrictFormatter(t *testing.T) {
	provider := jsonValidator.New(map[string]interface{}{
		"int": 1,
		"str": "abc",
		"obj": map[string]interface{}{
			"name":  "abc",
			"value": 123,
		},
	})

	formatter := _braceStrict
	var str string
	var err error

	str, err = formatter.Format(`{"id":{$obj.value},"name":{$obj.name},"str":{$str}}`, provider)
	if assert.NoError(t, err) {
		assert.Equal(t, `{"id":123,"name":"abc","str":"abc"}`, str)
	}
	str, err = formatter.Format(`{"id":$obj.value,"name":$obj.name,"str":$str}`, provider)
	if assert.NoError(t, err) {
		assert.Equal(t, `{"id":$obj.value,"name":$obj.name,"str":$str}`, str)
	}

	str, err = formatter.Format(`{"{$obj.name}":{$obj.value}}`, provider)
	if assert.NoError(t, err) {
		assert.Equal(t, `{"abc":123}`, str)
	}
	str, err = formatter.Format(`{"$obj.name":$obj.value}`, provider)
	if assert.NoError(t, err) {
		assert.Equal(t, `{"$obj.name":$obj.value}`, str)
	}

	str, err = formatter.Format(`{"$str":{$obj.value}}`, provider)
	if assert.NoError(t, err) {
		assert.Equal(t, `{"$str":123}`, str)
	}
	if assert.NoError(t, err) {
		str, err = formatter.Format(`{"$str":$obj.value}`, provider)
		assert.Equal(t, `{"$str":$obj.value}`, str)
	}

	str, err = formatter.Format(`{"$strstr":{$obj.value}}`, provider)
	if assert.NoError(t, err) {
		assert.Equal(t, `{"$strstr":123}`, str)
	}

	str, err = formatter.Format(`$(($used * 100 / $total))`, provider)
	if assert.NoError(t, err) {
		assert.Equal(t, `$(($used * 100 / $total))`, str)
	}
}
