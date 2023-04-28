package graph

import "github.com/3th1nk/easygo/util/mapUtil"

// Line 节点间的连线
type Line struct {
	Id        string                  `json:"id,omitempty"`        // 连线 ID，在一个图中唯一
	Left      string                  `json:"left,omitempty"`      // 左侧节点 ID
	Right     string                  `json:"right,omitempty"`     // 右侧节点 ID
	Direction Direction               `json:"direction,omitempty"` // 连线方向
	Type      string                  `json:"type,omitempty"`      // 连线类型（可选，可以用于指示 Data 的结构）
	Data      mapUtil.StringObjectMap `json:"data,omitempty"`      // 节点数据
}

type lineWrap struct {
	*Line
	leftNode  *nodeWrap
	rightNode *nodeWrap
}

func (this *lineWrap) safeLeft() *Node {
	if this.leftNode != nil {
		return this.leftNode.Node
	}
	return &Node{}
}

func (this *lineWrap) safeRight() *Node {
	if this.rightNode != nil {
		return this.rightNode.Node
	}
	return &Node{}
}
