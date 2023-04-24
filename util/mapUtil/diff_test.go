package mapUtil

import (
	"github.com/3th1nk/easygo/internal/_test"
	"github.com/3th1nk/easygo/util/convertor"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

type strIntMap map[string]int

func TestDiffEmptyOfError(t *testing.T) {
	var diff interface{}
	var err error

	//
	diff = Diff(nil, nil)
	if diff != nil {
		t.Errorf("assert faild 1")
	}
	diff = Diff(nil, strIntMap{})
	if diff == nil {
		t.Errorf("assert faild 2")
	}
	diff = Diff(map[string]string{}, nil)
	if diff == nil {
		t.Errorf("assert faild 3")
	}

	// 类型不同，需要报错
	_, err = testDiffError(map[string]string{}, strIntMap{})
	if err == nil {
		t.Errorf("类型不同，应该报错而未报")
	}

	// 虽然实际上是同一个类型，但仍然需要报错
	_, err = testDiffError(map[string]int{}, strIntMap{})
	if err == nil {
		t.Errorf("类型不同，应该报错而未报")
	}

	// 同一个类型，一个是结构体一个是指针
	_, err = testDiffError(&strIntMap{}, strIntMap{})
	if err == nil {
		t.Errorf("类型不同，应该报错而未报")
	}
}

func testDiffError(a, b interface{}) (v interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = convertor.ToError(e)
		}
	}()
	return Diff(a, b), nil
}

func TestDiffType1(t *testing.T) {
	diff := Diff(strIntMap{}, strIntMap{})
	_, ok := diff.(strIntMap)
	if !ok {
		t.Errorf("类型转换出错: %v", diff)
	}
}

func TestDiffType2(t *testing.T) {
	diff := Diff(&strIntMap{}, &strIntMap{})
	_, ok := diff.(*strIntMap)
	if !ok {
		t.Errorf("类型转换出错: %v", diff)
	}
}

func TestDiffMap(t *testing.T) {
	diff := Diff(strIntMap{"a": 1, "b": 2, "c": 3}, strIntMap{"a": 1, "b": 222})
	dict, ok := diff.(strIntMap)
	if !ok {
		t.Errorf("类型转换出错: %v", diff)
		return
	}
	if len(dict) != 2 {
		t.Errorf("assert faild: %v", len(dict))
	}
	if dict["b"] != 222 {
		t.Errorf("assert faild: %v", dict["b"])
	}
	if dict["c"] != 0 {
		t.Errorf("assert faild: %v", dict["c"])
	}
}

func TestUnionMap(t *testing.T) {
	diff := Diff(nil, strIntMap{"a": 1, "b": 2, "c": 3}, strIntMap{"a": 1, "b": 222})
	dict, ok := diff.(strIntMap)
	if !ok {
		t.Errorf("类型转换出错: %v", diff)
		return
	}
	if len(dict) != 3 {
		t.Errorf("assert faild: %v", len(dict))
	}
	if dict["a"] != 1 {
		t.Errorf("assert faild: %v", dict["b"])
	}
	if dict["b"] != 222 {
		t.Errorf("assert faild: %v", dict["b"])
	}
	if dict["c"] != 3 {
		t.Errorf("assert faild: %v", dict["c"])
	}
}

func TestDiffMapPtr(t *testing.T) {
	diff := Diff(&strIntMap{"a": 1, "b": 2, "c": 3}, &strIntMap{"a": 1, "b": 222})
	dict, ok := diff.(*strIntMap)
	if !ok {
		t.Errorf("类型转换出错: %v", diff)
		return
	}
	if len(*dict) != 2 {
		t.Errorf("assert faild: %v", len(*dict))
	}
	if (*dict)["b"] != 222 {
		t.Errorf("assert faild: %v", (*dict)["b"])
	}
	if (*dict)["c"] != 0 {
		t.Errorf("assert faild: %v", (*dict)["c"])
	}
}

func TestUnionMapPtr(t *testing.T) {
	diff := Diff(nil, &strIntMap{"a": 1, "b": 2, "c": 3}, &strIntMap{"a": 1, "b": 222})
	dict, ok := diff.(*strIntMap)
	if !ok {
		t.Errorf("类型转换出错: %v", diff)
		return
	}
	if len(*dict) != 3 {
		t.Errorf("assert faild: %v", len(*dict))
	}
	if (*dict)["a"] != 1 {
		t.Errorf("assert faild: %v", (*dict)["b"])
	}
	if (*dict)["b"] != 222 {
		t.Errorf("assert faild: %v", (*dict)["b"])
	}
	if (*dict)["c"] != 3 {
		t.Errorf("assert faild: %v", (*dict)["c"])
	}
}

func TestDiff2(t *testing.T) {
	src := map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	dest := map[string]interface{}{
		"a": 1,
		"b": 3,
		"d": 3,
	}
	a, b, c, d := Diff2(src, dest, nil)
	matches, _ := a.(map[string]interface{})
	changed, _ := b.(map[string]interface{})
	added, _ := c.(map[string]interface{})
	removed, _ := d.(map[string]interface{})

	assert.Equal(t, 1, matches["a"])
	assert.Equal(t, 3, changed["b"])
	assert.Equal(t, 3, added["d"])
	assert.Equal(t, 3, removed["c"])
}

func TestDiff2Str(t *testing.T) {
	src := map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}
	dest := map[string]string{
		"a": "1",
		"b": "3",
		"d": "3",
	}
	matches, changed, added, removed := DiffStr2(src, dest)

	assert.Equal(t, "1", matches["a"])
	assert.Equal(t, "3", changed["b"])
	assert.Equal(t, "3", added["d"])
	assert.Equal(t, "3", removed["c"])
}

func TestDiffStr_Perf(t *testing.T) {
	src := map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	}
	dest := map[string]string{
		"a": "1",
		"b": "3",
		"d": "3",
	}

	_test.Perf(func(i int) {
		DiffStr(src, dest)
	})
	_test.Perf(func(i int) {
		DiffStr2(src, dest)
	})
}

func TestLeftDiffKeyInt64(t *testing.T) {

	rand.Seed(time.Now().UnixNano())
	left, right := map[int64]bool{}, map[int64]bool{}
	for i := 0; i < 100; i++ {
		left[int64(rand.Intn(100))] = true
		right[int64(rand.Intn(100))] = true
	}

	t.Log("left", SortedInt64Keys(left))
	t.Log("right", SortedInt64Keys(right))
	t.Log("leftDiff", "\n\t", LeftDiffKeyInt64(left, right, true))

}
