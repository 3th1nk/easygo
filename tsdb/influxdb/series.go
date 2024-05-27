package influxdb

import (
	"github.com/3th1nk/easygo/util/convertor"
	"github.com/3th1nk/easygo/util/mapUtil"
)

type Series struct {
	Name    string            `json:"name,omitempty"` // 表名，多表查询时有用
	Tags    map[string]string `json:"tags,omitempty"`
	Columns []string          `json:"columns"`
	Values  [][]interface{}   `json:"values"`
}

// ToStringObjectMap converts the Series to a slice of mapUtil.StringObjectMap.
// 	[
//		{column1: value1, column2: value2, ...}, // row 1
//		{column1: value1, column2: value2, ...}, // row 2
//		...
//	]
func (this *Series) ToStringObjectMap() []mapUtil.StringObjectMap {
	m := make([]mapUtil.StringObjectMap, len(this.Values))
	for i, arr := range this.Values {
		row := make(mapUtil.StringObjectMap, len(this.Columns))
		for j, v := range arr {
			row[this.Columns[j]] = v
		}
		m[i] = row
	}
	return m
}

type RetentionPolicy struct {
	Name               string `json:"name"`
	Duration           string `json:"duration"`           // "8760h0m0s"
	ShardGroupDuration string `json:"shardGroupDuration"` // "168h0m0s"
	Replication        int    `json:"replication"`        // 副本数量
	Default            bool   `json:"default"`            // 是否默认策略
}

func (this *Series) toRetentionPolicies() []*RetentionPolicy {
	if len(this.Values) == 0 {
		return nil
	}

	rps := make([]*RetentionPolicy, 0, len(this.Values))
	for _, arr := range this.Values {
		rps = append(rps, &RetentionPolicy{
			Name:               convertor.ToStringNoError(arr[0]),
			Duration:           convertor.ToStringNoError(arr[1]),
			ShardGroupDuration: convertor.ToStringNoError(arr[2]),
			Replication:        convertor.ToIntNoError(arr[3]),
			Default:            convertor.ToBoolNoError(arr[4]),
		})
	}
	return rps
}

// toStringSlice 将指定列的值转换为字符串切片
//	- idx: 指定列的索引，不传时返回第一列
func (this *Series) toStringSlice(idx ...int) []string {
	if len(this.Values) == 0 {
		return nil
	}

	var colIdx int
	if len(idx) > 0 {
		colIdx = idx[0]
	}

	arr := make([]string, 0, len(this.Values))
	for _, values := range this.Values {
		if colIdx >= len(values) {
			continue
		}
		arr = append(arr, convertor.ToStringNoError(values[colIdx]))
	}
	return arr
}
