package influxdb

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"net/url"
	"sort"
	"strings"
	"sync/atomic"
)

func (this *Client) buildWriteUrl(db, rp string) string {
	uri := this.addr + "/write?precision=" + this.writePrecision
	if db != "" {
		uri += "&db=" + url.QueryEscape(db)
	}
	if rp != "" {
		uri += "&rp=" + url.QueryEscape(rp)
	}
	if this.username != "" {
		uri += "&u=" + url.QueryEscape(this.username)
	}
	if this.password != "" {
		uri += "&p=" + url.QueryEscape(this.password)
	}
	return uri
}

// RawWrite 执行写入操作
//	- db: 数据库名，必须指定
//	- rp: 保留策略名，为空时使用默认
//	- lines: 符合influxdb行协议的数据，需自行处理转义
func (this *Client) RawWrite(db, rp string, lines []string) error {
	if db == "" {
		return fmt.Errorf("missing database")
	}

	writeUrl := this.buildWriteUrl(db, rp)
	return this.doBatchWrite(writeUrl, lines)
}

// Point 数据点
//	注意：写入精度默认为秒，可以通过 WithWritePrecision 设置
type Point struct {
	Measurement string                 // 表名，必须指定
	Tags        map[string]interface{} // tag字段
	Fields      map[string]interface{} // field字段
	Time        int64                  // 时间戳，不传时数据库会自动写入当前时间
}

// ToLineData 转换为行协议数据
// 	行协议：<measurement>[,<tag_key>=<tag_value>[,<tag_key>=<tag_value>]] <field_key>=<field_value>[,<field_key>=<field_value>] [<timestamp>]
func (p *Point) ToLineData(sortTagKey bool) string {

	tagArr := make([]string, 0, len(p.Tags))
	if sortTagKey {
		// 标签按照key排序，提升写入性能
		// https://docs.influxdata.com/influxdb/v1/write_protocols/line_protocol_reference/#performance-tips
		sortedTagKey := make([]string, 0, len(p.Tags))
		for k := range p.Tags {
			sortedTagKey = append(sortedTagKey, k)
		}
		sort.Slice(sortedTagKey, func(i, j int) bool {
			return bytes.Compare([]byte(sortedTagKey[i]), []byte(sortedTagKey[j])) < 0
		})

		for _, k := range sortedTagKey {
			v := p.Tags[k]
			if t, ok := v.(string); ok {
				v = EscapeTagValue(t)
			}
			tagArr = append(tagArr, fmt.Sprintf("%s=%v", k, v))
		}

	} else {
		for k, v := range p.Tags {
			if t, ok := v.(string); ok {
				v = EscapeTagValue(t)
			}
			tagArr = append(tagArr, fmt.Sprintf("%s=%v", k, v))
		}
	}

	fieldArr := make([]string, 0, len(p.Fields))
	for k, v := range p.Fields {
		if t, ok := v.(string); ok {
			v = Quote(EscapeFieldValue(t))
		}
		fieldArr = append(fieldArr, fmt.Sprintf(`%s=%v`, k, v))
	}

	data := p.Measurement
	if len(tagArr) > 0 {
		data += "," + strings.Join(tagArr, ",")
	}
	if len(fieldArr) > 0 {
		data += " " + strings.Join(fieldArr, ",")
	}
	if p.Time > 0 {
		data += fmt.Sprintf(" %d", p.Time)
	}
	return data
}

// Write 写入数据
//	- db: 数据库名，必须指定
//	- rp: 保留策略名，为空时使用默认
//	- points: 数据点
// 	- immediate: 是否立即写入
func (this *Client) Write(db, rp string, points []*Point, immediate bool) error {
	if db == "" {
		return fmt.Errorf("missing database")
	}

	if len(points) == 0 {
		return nil
	}

	writeUrl := this.buildWriteUrl(db, rp)

	if immediate || atomic.LoadInt32(&this.state) != stateRunning {
		if !immediate {
			this.logger.Warn("[InfluxDB] 异步写入已停止，立即写入")
		}

		lines := make([]string, 0, len(points))
		for _, pt := range points {
			lines = append(lines, pt.ToLineData(this.writeSortTagKey))
		}
		return this.doBatchWrite(writeUrl, lines)
	}

	groupName := makeBucketGroupName(db, rp)
	this.mu.Lock()
	group, ok := this.bucketGroups[groupName]
	if !ok {
		group = newBucketGroup(db, rp, this.groupSize, this.flushSize)
		this.bucketGroups[groupName] = group
	}
	this.mu.Unlock()

	m := make(map[string][]string)
	for _, pt := range points {
		m[pt.Measurement] = append(m[pt.Measurement], pt.ToLineData(this.writeSortTagKey))
	}
	for measurement, lines := range m {
		idx := int(crc32.ChecksumIEEE([]byte(measurement))) % group.Size()
		if err := group.Push(idx, lines, func(lines []string) error {
			return this.writePool.Submit(func() {
				_ = this.doBatchWrite(writeUrl, lines)
			})
		}); err != nil {
			this.logger.Warn("[InfluxDB] 异步写入异常: %v", err)
		}
	}

	return nil
}
