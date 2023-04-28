package graph

import "github.com/3th1nk/easygo/util/mapUtil"

// Combo 组合
type Combo struct {
	Id      string                  `json:"id,omitempty"`       // 组合 ID，在一个图中唯一
	Type    string                  `json:"type,omitempty"`     // 组合类型（可选，可以用于指示 Data 的结构）
	ComboId string                  `json:"combo_id,omitempty"` // 父级 Combo 的 ID
	Data    mapUtil.StringObjectMap `json:"data,omitempty"`     // 节点数据
}

type comboWrap struct {
	*Combo
	nodes  map[string]*nodeWrap
	combos map[string]*comboWrap
}

// GetNodes 获取组合下的节点（不含子组合下的节点）
func (this *comboWrap) GetNodes() []*Node {
	if n := len(this.nodes); n == 0 {
		return nil
	} else {
		arr := make([]*Node, 0, n)
		for _, v := range this.nodes {
			arr = append(arr, v.Node)
		}
		return arr
	}
}

// GetAllNodes 获取 组合 及 子组合 下所有的节点
func (this *comboWrap) GetAllNodes() []*Node {
	if n := len(this.nodes); n == 0 {
		return nil
	} else {
		arr := make([]*Node, 0, n)
		for _, v := range this.nodes {
			arr = append(arr, v.Node)
		}
		for _, v := range this.combos {
			arr = append(arr, v.GetAllNodes()...)
		}
		return arr
	}
}

// GetCombos 获取组合下直接嵌套的组合
func (this *comboWrap) GetCombos() []*Combo {
	if n := len(this.combos); n == 0 {
		return nil
	} else {
		arr := make([]*Combo, 0, n)
		for _, v := range this.combos {
			arr = append(arr, v.Combo)
		}
		return arr
	}
}

// GetAllCombos 获取组合 及 子组合 下所有的组合
func (this *comboWrap) GetAllCombos() []*Combo {
	if n := len(this.combos); n == 0 {
		return nil
	} else {
		arr := make([]*Combo, 0, n)
		for _, v := range this.combos {
			arr = append(arr, v.Combo)
			arr = append(arr, v.GetAllCombos()...)
		}
		return arr
	}
}
