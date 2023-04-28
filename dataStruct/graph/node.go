package graph

import "github.com/3th1nk/easygo/util/mapUtil"

// Node 节点
type Node struct {
	Id      string                  `json:"id,omitempty"`       // 节点 ID，在一个图中唯一
	Level   int                     `json:"level,omitempty"`    // 节点层级（可选，可以用于树状布局时指定 节点 的层级）
	Type    string                  `json:"type,omitempty"`     // 节点类型（可选，可以用于指示 Data 的结构）
	ComboId string                  `json:"combo_id,omitempty"` // 节点所属的 Combo 的 ID
	Data    mapUtil.StringObjectMap `json:"data,omitempty"`     // 节点数据
}

type nodeWrap struct {
	*Node
	lines map[string]*lineWrap
}

func (this *nodeWrap) GetLines() []*Line {
	if n := len(this.lines); n == 0 {
		return nil
	} else {
		arr := make([]*Line, 0, n)
		for _, v := range this.lines {
			arr = append(arr, v.Line)
		}
		return arr
	}
}
