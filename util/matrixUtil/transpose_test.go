package matrixUtil

import (
	"github.com/3th1nk/easygo/util"
	"github.com/3th1nk/easygo/util/jsonUtil"
	"github.com/3th1nk/easygo/util/mapUtil"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestTransposeInt(t *testing.T) {
	a := [][]int{
		{11, 12, 13, 14, 15},
		{21, 22, 23, 24, 25},
		{31, 32, 33, 34, 35},
	}
	b := TransposeInt(a)
	assert.Equal(t, 5, len(b))
	assert.Equal(t, 3, len(b[0]))
	assert.Equal(t, 11, b[0][0])
	assert.Equal(t, 21, b[0][1])
	assert.Equal(t, 31, b[0][2])
	assert.Equal(t, 13, b[2][0])
	assert.Equal(t, 23, b[2][1])
	assert.Equal(t, 33, b[2][2])
	assert.Equal(t, 15, b[4][0])
	assert.Equal(t, 25, b[4][1])
	assert.Equal(t, 35, b[4][2])
}

func TestTransposeStringObjectTable_1(t *testing.T) {
	rows, err := TransposeStringObjectTable("name", []mapUtil.StringObjectMap{
		{"name": "auto_increment_increment", "value": "1"},
		{"name": "auto_increment_offset", "value": "1"},
		{"name": "big_tables", "value": "OFF"},
		{"name": "binlog_direct_non_transactional_updates", "value": "OFF"},
		{"name": "binlog_order_commits", "value": "ON"},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(rows))
	for _, key := range mapUtil.SortedStringKeys(rows[0]) {
		util.Println("%v = %v", key, rows[0][key])
	}
	assert.Equal(t, "1", rows[0]["auto_increment_increment"])
	assert.Equal(t, "OFF", rows[0]["big_tables"])
	assert.Equal(t, "OFF", rows[0]["binlog_direct_non_transactional_updates"])
	assert.Equal(t, "ON", rows[0]["binlog_order_commits"])
}

func TestTransposeStringObjectTable_2(t *testing.T) {
	rows, err := TransposeStringObjectTable("name", []mapUtil.StringObjectMap{
		{"id": 1, "name": "aaa", "count": 10},
		{"id": 2, "name": "bbb", "count": 20},
		{"id": 3, "name": "ccc", "count": 30},
		{"id": 4, "name": "ddd", "count": 40},
		{"id": 5, "name": "eee", "count": 50},
	})
	assert.NoError(t, err)
	sort.Slice(rows, func(i, j int) bool { return rows[i].MustGetInt("aaa") < rows[j].MustGetInt("aaa") })
	for _, row := range rows {
		util.Println(jsonUtil.MustMarshalToString(row))
	}
	assert.Equal(t, 2, len(rows))
	assert.Equal(t, 1, rows[0]["aaa"])
	assert.Equal(t, 2, rows[0]["bbb"])
	assert.Equal(t, 3, rows[0]["ccc"])
	assert.Equal(t, 4, rows[0]["ddd"])
	assert.Equal(t, 5, rows[0]["eee"])
	assert.Equal(t, 10, rows[1]["aaa"])
	assert.Equal(t, 20, rows[1]["bbb"])
	assert.Equal(t, 30, rows[1]["ccc"])
	assert.Equal(t, 40, rows[1]["ddd"])
	assert.Equal(t, 50, rows[1]["eee"])
}
