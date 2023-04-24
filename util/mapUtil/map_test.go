package mapUtil

import (
	"github.com/3th1nk/easygo/util/jsonUtil"
	"github.com/3th1nk/easygo/util/strUtil"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

// 验证 map 是引用类型：
// 1、函数内部对 map 的改动在函数外面是可以拿到新值的
// 2、map 变量赋值，两个变量实际上指向的是同一个 map
// 3、map 是引用，但不是指针，在函数内部 new 之后赋值，外层并不会变化
func TestMapIsRef(t *testing.T) {
	setMap := func(a map[string]interface{}) {
		if a == nil {
			a = make(map[string]interface{})
		}
		a["a"] = 1
		a["b"] = 2
		a["c"] = 3
	}

	// 验证 1
	{
		dict := make(map[string]interface{})
		setMap(dict)
		assert.Equal(t, 1, dict["a"])
		assert.Equal(t, 2, dict["b"])
		assert.Equal(t, 3, dict["c"])
	}

	// 验证 2
	{
		dict1 := make(map[string]interface{})
		dict2 := dict1
		setMap(dict1)
		assert.Equal(t, 1, dict2["a"])
		assert.Equal(t, 2, dict2["b"])
		assert.Equal(t, 3, dict2["c"])
	}

	// 验证 2
	{
		dict1 := make(map[string]interface{})
		var dict2 StringObjectMap = dict1
		setMap(dict1)
		assert.Equal(t, 1, dict2["a"])
		assert.Equal(t, 2, dict2["b"])
		assert.Equal(t, 3, dict2["c"])
	}

	// 验证 3
	{
		var dict map[string]interface{}
		setMap(dict)
		assert.Nil(t, dict)
	}
}

// 验证从 nil map 中获取值不会产生 panic
func TestGetFromNilMapIsSafe(t *testing.T) {
	var a map[string]int
	assert.Nil(t, a)
	assert.NotPanics(t, func() {
		assert.Equal(t, 0, a["id"])
	})
}

// 验证：在 for 循环遍历过程中进行 delete，不会引起程序异常（其他语言在迭代期间修改容器可能会引起异常）
func TestDeleteInEachIsSafe(t *testing.T) {
	const loop = 100000
	dict := make(map[string]int, loop)
	for i := 0; i < loop; i++ {
		dict[strUtil.Rand(12)] = rand.Int()
	}

	const mod = 10
	deletedKeys := make(map[string]bool, loop/mod)
	for key, val := range dict {
		// 在遍历过程中，对其中 1/10 的数据进行删除
		if val%mod == 0 {
			deletedKeys[key] = true
			delete(dict, key)
		}
	}

	assert.Equal(t, loop-len(deletedKeys), len(dict))
	for key := range dict {
		assert.Equal(t, false, deletedKeys[key])
	}
	for key := range deletedKeys {
		_, ok := dict[key]
		assert.Equal(t, false, ok)
	}
}

func TestUnmarshal_1(t *testing.T) {
	var a StringMap
	var b CaseInsensitiveStringMap
	var c StringObjectMap
	var d CaseInsensitiveStringObjectMap

	assert.NoError(t, jsonUtil.UnmarshalFromString(`{"id": "1", "name": "abc"}`, &a))
	assert.NoError(t, jsonUtil.UnmarshalFromString(`{"id": "1", "name": "abc"}`, &b))
	assert.NoError(t, jsonUtil.UnmarshalFromString(`{"id": "1", "name": "abc"}`, &c))
	assert.NoError(t, jsonUtil.UnmarshalFromString(`{"id": "1", "name": "abc"}`, &d))
	assert.Equal(t, 2, len(a))
	assert.Equal(t, 2, len(b))
	assert.Equal(t, 2, len(c))
	assert.Equal(t, 2, len(d))
	assert.Equal(t, 1, a.MustGetInt("id"))
	assert.Equal(t, 1, b.MustGetInt("id"))
	assert.Equal(t, 1, c.MustGetInt("id"))
	assert.Equal(t, 1, d.MustGetInt("id"))
	assert.Equal(t, "abc", a.MustGet("name"))
	assert.Equal(t, "abc", b.MustGet("name"))
	assert.Equal(t, "abc", c.MustGet("name"))
	assert.Equal(t, "abc", d.MustGet("name"))

	type tmp struct {
		Map interface{}
	}
	for _, v := range []interface{}{StringMap{}, StringObjectMap{}, a, b, c, d} {
		str, err := jsonUtil.MarshalToString(&tmp{Map: v})
		t.Log(str)
		if assert.NoError(t, err) {
			b := &tmp{}
			if assert.NoError(t, jsonUtil.UnmarshalFromString(str, b)) {
				assert.Equal(t, jsonUtil.SortMapKeysApi().MustMarshalToString(v), jsonUtil.SortMapKeysApi().MustMarshalToString(b.Map))
			}
		}
	}
}
