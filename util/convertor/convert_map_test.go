package convertor

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToStringObjectMap(t *testing.T) {
	{
		src := "123"
		map1, err1 := ToStringObjectMap(src)
		map2 := make(map[string]interface{})
		data, _ := json.Marshal(src)
		err2 := json.Unmarshal(data, &map2)
		assert.Error(t, err1)
		assert.Error(t, err2)
		assert.Equal(t, 0, len(map1))
	}
	{
		src := "abc"
		map1, err1 := ToStringObjectMap(src)
		map2 := make(map[string]interface{})
		data, _ := json.Marshal(src)
		err2 := json.Unmarshal(data, &map2)
		assert.Error(t, err1)
		assert.Error(t, err2)
		assert.Equal(t, 0, len(map1))
	}
	{
		src := []interface{}{1, 2, 3}
		map1, err1 := ToStringObjectMap(src)
		map2 := make(map[string]interface{})
		data, _ := json.Marshal(src)
		err2 := json.Unmarshal(data, &map2)
		assert.Error(t, err1)
		assert.Error(t, err2)
		assert.Equal(t, 0, len(map1))
	}
}
