package mapUtil

import (
	"github.com/3th1nk/easygo/util"
	"github.com/3th1nk/easygo/util/arrUtil"
	"github.com/3th1nk/easygo/util/convertor"
	"github.com/3th1nk/easygo/util/jsonUtil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFlatMap_0(t *testing.T) {
	in := map[string]interface{}{
		"code":    0,
		"message": "hello",
		"data": map[string]interface{}{
			"id":   1,
			"name": "a",
			"list": []map[string]interface{}{
				{"id": 1, "name": "a"},
				{"id": 2, "name": "b"},
			},
			"extra": map[string]string{
				"abc": "123",
			},
		},
	}
	out, err := FlatMap(in)
	util.Println(jsonUtil.MustMarshalToStringIndent(out))
	assert.NoError(t, err)
	assert.Equal(t, 9, len(out))
	assert.Equal(t, 0, out["code"])
	assert.Equal(t, "hello", out["message"])
	assert.Equal(t, 1, out["data.id"])
	assert.Equal(t, "a", out["data.name"])
	assert.Equal(t, 1, out["data.list.0.id"])
	assert.Equal(t, "a", out["data.list.0.name"])
	assert.Equal(t, 2, out["data.list.1.id"])
	assert.Equal(t, "b", out["data.list.1.name"])
	assert.Equal(t, "123", out["data.extra.abc"])
}

func TestFlatMap_1(t *testing.T) {
	in := map[string]interface{}{
		"code":    0,
		"message": "hello",
		"data": map[string]interface{}{
			"id":   1,
			"name": "a",
			"list": []map[string]interface{}{
				{"id": 1, "name": "a"},
				{"id": 2, "name": "b"},
			},
			"extra": map[string]string{
				"abc": "123",
			},
		},
	}
	out, err := FlatMap(in, "_")
	util.Println(jsonUtil.MustMarshalToStringIndent(out))
	assert.NoError(t, err)
	assert.Equal(t, 9, len(out))
	assert.Equal(t, 0, out["code"])
	assert.Equal(t, "hello", out["message"])
	assert.Equal(t, 1, out["data_id"])
	assert.Equal(t, "a", out["data_name"])
	assert.Equal(t, 1, out["data_list_0_id"])
	assert.Equal(t, "a", out["data_list_0_name"])
	assert.Equal(t, 2, out["data_list_1_id"])
	assert.Equal(t, "b", out["data_list_1_name"])
	assert.Equal(t, "123", out["data_extra_abc"])
}

func TestFlatStringMap_1(t *testing.T) {
	in := map[string]interface{}{
		"code":    0,
		"message": "hello",
		"data": map[string]interface{}{
			"id":   1,
			"name": "a",
			"list": []map[string]interface{}{
				{"id": 1, "name": "a"},
				{"id": 2, "name": "b"},
			},
			"extra": map[string]string{
				"abc": "123",
			},
		},
	}
	out, err := FlatStringMap(in, "_")
	assert.NoError(t, err)
	assert.Equal(t, 9, len(out))
	assert.Equal(t, "0", out["code"])
	assert.Equal(t, "hello", out["message"])
	assert.Equal(t, "1", out["data_id"])
	assert.Equal(t, "a", out["data_name"])
	assert.Equal(t, "1", out["data_list_0_id"])
	assert.Equal(t, "a", out["data_list_0_name"])
	assert.Equal(t, "2", out["data_list_1_id"])
	assert.Equal(t, "b", out["data_list_1_name"])
	assert.Equal(t, "123", out["data_extra_abc"])
	util.Println(jsonUtil.MustMarshalToStringIndent(out))
}

func TestFlatMap_2(t *testing.T) {
	in := []interface{}{
		map[string]interface{}{"id": 1, "name": "a"},
		map[string]interface{}{"id": 2, "name": "b"},
		map[string]interface{}{"id": 3, "name": "c"},
	}
	out, err := FlatMap(in)
	util.Println(jsonUtil.MustMarshalToStringIndent(out))
	assert.NoError(t, err)
	assert.Equal(t, 1, out.MustGetInt("0.id"))
	assert.Equal(t, "a", out.MustGetString("0.name"))
	assert.Equal(t, 2, out.MustGetInt("1.id"))
	assert.Equal(t, "b", out.MustGetString("1.name"))
	assert.Equal(t, 3, out.MustGetInt("2.id"))
	assert.Equal(t, "c", out.MustGetString("2.name"))
}

func TestFlatMap_3(t *testing.T) {
	in := []interface{}{
		map[string]interface{}{"id": 1, "name": "a"},
		map[string]interface{}{"id": 2, "name": "b"},
		map[string]interface{}{"id": 3, "name": "c"},
	}
	out, err := FlatMap(in, "_")
	util.Println(jsonUtil.MustMarshalToStringIndent(out))
	assert.NoError(t, err)
	assert.Equal(t, 1, out.MustGetInt("0_id"))
	assert.Equal(t, "a", out.MustGetString("0_name"))
	assert.Equal(t, 2, out.MustGetInt("1_id"))
	assert.Equal(t, "b", out.MustGetString("1_name"))
	assert.Equal(t, 3, out.MustGetInt("2_id"))
	assert.Equal(t, "c", out.MustGetString("2_name"))
}

func TestUnFlatMap_0(t *testing.T) {
	out, err := convertor.ToStringObjectMap(`
{
    "code": 1,
    "data.extra.abc": "123123",
    "data.id": 1,
    "data.list.0.id": 1,
    "data.list.0.name": "a",
    "data.list.1.id": 2,
    "data.list.1.name": "b",
    "data.name": "a",
    "message": "hello"
}`)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(out))
	obj, err := UnFlatMap(out)
	assert.NoError(t, err)

	dict := StringObjectMap(obj.(map[string]interface{}))
	assert.NotEqual(t, 0, len(dict))
	assert.Equal(t, 1, dict.MustGetInt("code"))
	assert.Equal(t, 1, dict.MustGetStringObjectMap("data").MustGetInt("id"))
	assert.Equal(t, "a", dict.MustGetStringObjectMap("data").MustGetString("name"))

	util.Println(jsonUtil.SortMapKeysApi(false).MustMarshalToStringIndent(dict))
}

func TestUnFlatMap_1(t *testing.T) {
	out, err := convertor.ToStringObjectMap(`
{
    "code": 1,
    "data_extra_abc": "123123",
    "data_id": 1,
    "data_list_0_id": 1,
    "data_list_0_name": "a",
    "data_list_1_id": 2,
    "data_list_1_name": "b",
    "data_name": "a",
    "message": "hello"
}`)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(out))
	obj, err := UnFlatMap(out, "_")
	assert.NoError(t, err)

	dict := StringObjectMap(obj.(map[string]interface{}))
	assert.NotEqual(t, 0, len(dict))
	assert.Equal(t, 1, dict.MustGetInt("code"))
	assert.Equal(t, 1, dict.MustGetStringObjectMap("data").MustGetInt("id"))
	assert.Equal(t, "a", dict.MustGetStringObjectMap("data").MustGetString("name"))

	util.Println(jsonUtil.SortMapKeysApi(false).MustMarshalToStringIndent(dict))
}

func TestUnFlatMap_2(t *testing.T) {
	out, err := convertor.ToStringObjectMap(`
{
    "0.id": 1,
    "0.name": "a",
    "1.id": 2,
    "1.name": "b",
    "2.id": 3,
    "2.name": "c"
}
`)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(out))
	obj, err := UnFlatMap(out)
	assert.NoError(t, err)
	util.Println(jsonUtil.SortMapKeysApi(false).MustMarshalToStringIndent(obj))

	arr := arrUtil.ToInterface(obj)
	assert.Equal(t, 3, len(arr))
	assert.Equal(t, 1, StringObjectMap(arr[0].(map[string]interface{})).MustGetInt("id"))
	assert.Equal(t, "a", StringObjectMap(arr[0].(map[string]interface{})).MustGetString("name"))
	assert.Equal(t, 2, StringObjectMap(arr[1].(map[string]interface{})).MustGetInt("id"))
	assert.Equal(t, "b", StringObjectMap(arr[1].(map[string]interface{})).MustGetString("name"))
	assert.Equal(t, 3, StringObjectMap(arr[2].(map[string]interface{})).MustGetInt("id"))
	assert.Equal(t, "c", StringObjectMap(arr[2].(map[string]interface{})).MustGetString("name"))
}

func TestUnFlatMap_3(t *testing.T) {
	out, err := convertor.ToStringObjectMap(`
{
    "0_id": 1,
    "0_name": "a",
    "1_id": 2,
    "1_name": "b",
    "2_id": 3,
    "2_name": "c"
}
`)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(out))
	obj, err := UnFlatMap(out, "_")
	assert.NoError(t, err)
	util.Println(jsonUtil.SortMapKeysApi(false).MustMarshalToStringIndent(obj))

	arr := arrUtil.ToInterface(obj)
	assert.Equal(t, 3, len(arr))
	assert.Equal(t, 1, StringObjectMap(arr[0].(map[string]interface{})).MustGetInt("id"))
	assert.Equal(t, "a", StringObjectMap(arr[0].(map[string]interface{})).MustGetString("name"))
	assert.Equal(t, 2, StringObjectMap(arr[1].(map[string]interface{})).MustGetInt("id"))
	assert.Equal(t, "b", StringObjectMap(arr[1].(map[string]interface{})).MustGetString("name"))
	assert.Equal(t, 3, StringObjectMap(arr[2].(map[string]interface{})).MustGetInt("id"))
	assert.Equal(t, "c", StringObjectMap(arr[2].(map[string]interface{})).MustGetString("name"))
}
