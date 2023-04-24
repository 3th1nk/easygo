package mapUtil

import (
	"github.com/3th1nk/easygo/util/convertor"
	"testing"
)

func TestRemoveMapObjectDuplicate(t *testing.T) {
	arr := []StringObjectMap{
		{"name": "aaa", "val": "aaa"},
		{"name": "bbb", "val": "bbb"},
		{"name": "aaa", "val": "aaa"},
		{"name": "ccc", "val": "ccc"},
		{"name": "aaa", "val": "bbb"},
	}
	t.Log(convertor.ToStringNoError(RemoveMapObjectDuplicate(arr, "name")))
	t.Log(convertor.ToStringNoError(RemoveMapObjectDuplicate(arr, "name", "val")))

}
