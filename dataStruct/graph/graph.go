package graph

import (
	"fmt"
	"github.com/3th1nk/easygo/util"
	"github.com/3th1nk/easygo/util/mapUtil"
	"github.com/3th1nk/easygo/util/strUtil"
)

// 图结构
type Graph struct {
	nodeMap  map[string]*nodeWrap // id -> Node
	lineMap  map[string]*lineWrap // id -> Line
	comboMap map[string]*Combo    // id -> Combo
}

// 创建一个图
func New() *Graph {
	g := &Graph{}
	return g
}

// 节点
type Node struct {
	Id      string                  `json:"id,omitempty"`       // 节点 ID，在一个图中唯一
	Level   int                     `json:"level,omitempty"`    // 节点等级（可选）
	Type    string                  `json:"type,omitempty"`     // 节点类型（可选，可以用于指示 Data 的结构）
	ComboId string                  `json:"combo_id,omitempty"` // 节点所属的 Combo 的 ID
	Data    mapUtil.StringObjectMap `json:"data,omitempty"`     // 节点数据
}

// 节点间的连线
type Line struct {
	Id        string                  `json:"id,omitempty"`        // 连线 ID，在一个图中唯一
	Left      string                  `json:"left,omitempty"`      // 左侧节点 ID
	Right     string                  `json:"right,omitempty"`     // 右侧节点 ID
	Direction Direction               `json:"direction,omitempty"` // 连线方向
	Type      string                  `json:"type,omitempty"`      // 连线类型（可选，可以用于指示 Data 的结构）
	Data      mapUtil.StringObjectMap `json:"data,omitempty"`      // 节点数据
}

// 组合框
type Combo struct {
	Id   string                  `json:"id,omitempty"`   // 组合框 ID，在一个图中唯一
	Type string                  `json:"type,omitempty"` // 组合框类型（可选，可以用于指示 Data 的结构）
	Data mapUtil.StringObjectMap `json:"data,omitempty"` // 节点数据
}

// 路径
type Path struct {
	Node []*Node `json:"node,omitempty"` // 该路径上的节点列表
	Line []*Line `json:"line,omitempty"` // 该路径上的连线列表
}

// 连线方向
type Direction int

const (
	LeftToRight Direction = 0 // 单向，Left -> Right
	RightToLeft Direction = 1 // 单向，Left <- Right
	BothWay     Direction = 2 // 双向，Left <-> Right
)

type nodeWrap struct {
	*Node
	lines []*lineWrap
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
	return nil
}

func (this *lineWrap) safeRight() *Node {
	if this.rightNode != nil {
		return this.rightNode.Node
	}
	return nil
}

func (this *Graph) init(renew bool) {
	if this.nodeMap == nil || renew {
		this.nodeMap = make(map[string]*nodeWrap, 32)
		this.lineMap = make(map[string]*lineWrap, 32)
		this.comboMap = make(map[string]*Combo, 8)
	}
}

// 获取节点数量
func (this *Graph) NodeCount() int {
	return len(this.nodeMap)
}

// 获取连线数量
func (this *Graph) LineCount() int {
	return len(this.lineMap)
}

// 获取所有节点
func (this *Graph) GetNodes() []*Node {
	arr, idx := make([]*Node, len(this.nodeMap)), 0
	for _, v := range this.nodeMap {
		arr[idx] = v.Node
	}
	return arr
}

// 获取所有连线
func (this *Graph) GetLines() []*Line {
	arr, idx := make([]*Line, len(this.lineMap)), 0
	for _, v := range this.lineMap {
		arr[idx] = v.Line
	}
	return arr
}

// 获取所有组合框
func (this *Graph) GetCombos() []*Combo {
	arr, idx := make([]*Combo, len(this.comboMap)), 0
	for _, v := range this.comboMap {
		arr[idx] = v
	}
	return arr
}

// 添加节点
//   注意：在变更节点数据之后，需要调用 Update 更新图。
//
// 返回值：
//   error: 如果要添加的节点已经存在，则会报 “node 'xxx' already exists” 错误。
func (this *Graph) AddNode(node ...*Node) error {
	if len(node) == 0 {
		return nil
	}

	this.init(false)

	for _, v := range node {
		if v.Id == "" {
			v.Id = strUtil.Rand(8)
		}

		if _, ok := this.nodeMap[v.Id]; ok {
			return fmt.Errorf("node '%v' already exists", v.Id)
		}
		this.nodeMap[v.Id] = &nodeWrap{Node: v, lines: make([]*lineWrap, 0, 2)}
	}
	return nil
}

// 添加连线
//   注意：在变更连线数据之后，需要调用 Update 更新图。
//
// 返回值：
//   error: 如果要添加的连线已经存在，则会报 “line 'xxx' already exists” 错误。
func (this *Graph) AddLine(line ...*Line) error {
	if len(line) == 0 {
		return nil
	}

	this.init(false)

	for _, v := range line {
		if v.Id == "" {
			v.Id = strUtil.Rand(8)
		}

		if _, ok := this.lineMap[v.Id]; ok {
			return fmt.Errorf("line '%v' already exists", v.Id)
		}
		wrap := &lineWrap{Line: v}
		_ = this.updateLine(wrap)
		this.lineMap[v.Id] = wrap
	}
	return nil
}

// 添加组合框
//   注意：在变更组合框数据之后，需要调用 Update 更新图。
//
// 返回值：
//   error: 如果要添加的组合框已经存在，则会报 “combo 'xxx' already exists” 错误。
func (this *Graph) AddCombo(combo ...*Combo) error {
	if len(combo) == 0 {
		return nil
	}

	this.init(false)

	for _, v := range combo {
		if v.Id == "" {
			v.Id = strUtil.Rand(8)
		}

		if _, ok := this.comboMap[v.Id]; ok {
			return fmt.Errorf("combo '%v' already exists", v.Id)
		}
		this.comboMap[v.Id] = v
	}
	return nil
}

// 删除指定的节点
//   注意：在变更数据之后，需要调用 Update 更新图。
func (this *Graph) DeleteNode(id ...string) {
	for _, v := range id {
		delete(this.nodeMap, v)
	}
}

// 删除指定的连线
//   注意：在变更数据之后，需要调用 Update 更新图。
func (this *Graph) DeleteLine(id ...string) {
	for _, v := range id {
		delete(this.lineMap, v)
	}
}

// 删除指定的组合框
//   注意：在变更数据之后，需要调用 Update 更新图。
func (this *Graph) DeleteCombo(id ...string) {
	for _, v := range id {
		delete(this.comboMap, v)
	}
}

// 更新图结构，并校验数据，确认数据格式正确。
//
// 当变更了节点、连线、组合框时调用。
//
// 返回：
//   err: 当校验失败时返回对应的错误：
//     * 节点的 ComboId 非空、但不存在对应的 Combo;
//     * 节点的 ComboId 非空、但不存在对应的 Combo;
func (this *Graph) Update() (err error) {
	return this.doUpdate()
}

func (this *Graph) doUpdate() (err error) {
	for _, v := range this.nodeMap {
		if v.ComboId != "" && nil == this.comboMap[v.ComboId] {
			return fmt.Errorf("node combo '%v.%v' not exists", v.Id, v.ComboId)
		}
	}

	for _, line := range this.lineMap {
		if err := this.updateLine(line); err != nil {
			return err
		}
	}

	return nil
}

func (this *Graph) updateLine(line *lineWrap) error {
	updateLineNode := func(line *lineWrap, nodeId string) *nodeWrap {
		node := this.nodeMap[nodeId]
		if node != nil {
			exists := false
			for _, v := range node.lines {
				if v == line {
					exists = true
					break
				}
			}
			if !exists {
				node.lines = append(node.lines, line)
			}
		}
		return node
	}

	if line.Left != "" {
		if line.leftNode = updateLineNode(line, line.Left); line.leftNode == nil {
			return fmt.Errorf("line left '%v.%v' not exists", line.Id, line.Left)
		}
	}
	if line.Right != "" {
		if line.rightNode = updateLineNode(line, line.Right); line.rightNode == nil {
			return fmt.Errorf("line right '%v.%v' not exists", line.Id, line.Right)
		}
	}

	return nil
}

// 根据 ‘节点 ID’ 获取节点信息
func (this *Graph) GetNode(id string) (node *Node) {
	if v := this.nodeMap[id]; v != nil {
		return v.Node
	}
	return nil
}

// 根据 ‘连线 ID’ 获取连线信息
func (this *Graph) GetLine(id string) (line *Line) {
	if v := this.lineMap[id]; v != nil {
		return v.Line
	}
	return nil
}

func (this *Graph) FindNode(f func(node *Node) bool) []*Node {
	cnt := len(this.nodeMap)
	arr, idx := make([]*Node, cnt), 0
	for _, v := range this.nodeMap {
		if f == nil || f(v.Node) {
			arr[idx], idx = v.Node, idx+1
		}
	}
	if idx == 0 {
		return nil
	} else if idx < cnt {
		arr = arr[:idx]
	}
	return arr
}

func (this *Graph) FindLine(f func(line *Line, left, right *Node) bool) []*Line {
	cnt := len(this.lineMap)
	arr, idx := make([]*Line, cnt), 0
	for _, v := range this.lineMap {
		if f == nil || f(v.Line, v.safeLeft(), v.safeRight()) {
			arr[idx], idx = v.Line, idx+1
		}
	}
	if idx == 0 {
		return nil
	} else if idx < cnt {
		arr = arr[:idx]
	}
	return arr
}

// 根据 ‘节点 ID’，获取与该节点直接关联的连线
func (this *Graph) GetNodeLines(id string) []*Line {
	if node := this.nodeMap[id]; node == nil {
		return nil
	} else if n := len(node.lines); n == 0 {
		return nil
	} else {
		arr := make([]*Line, n)
		for i, v := range node.lines {
			arr[i] = v.Line
		}
		return arr
	}
}

// 根据 ‘节点 ID’，获取与该节点直接关联的邻居节点。
func (this *Graph) GetNodeNeighbours(id string) []*Node {
	if node := this.nodeMap[id]; node == nil {
		return nil
	} else if n := len(node.lines); n == 0 {
		return nil
	} else {
		arr, dict := make([]*Node, 0, n), make(map[string]bool, n)
		for _, line := range node.lines {
			var another *nodeWrap
			switch line.Direction {
			case LeftToRight:
				another = line.rightNode
			case RightToLeft:
				another = line.leftNode
			case BothWay:
				if line.Left == node.Id {
					another = line.rightNode
				} else if line.Right == node.Id {
					another = line.leftNode
				}
			}
			if another != nil && !dict[another.Id] {
				dict[another.Id] = true
				arr = append(arr, another.Node)
			}
		}
		return arr
	}
}

// 路径搜索：
//   搜索能够联通起止的节点的所有路径。
//
// 参数：
//   from: 起始节点 ID
//   to:   截至节点 ID
//   f:    回调函数，如果该函数返回 false 则会停止在当前路径上的搜索。
//         参数：
//           path: 当前路径，包含路径上的所有节点和连线（节点数等于连线数+1，因为包含一个起始节点）。
//                 比如：可以通过 return len(path.Line) < maxDepth 限制搜索路径的最大深度。
func (this *Graph) Traversal(from, to string, f ...func(path *Path, ended bool) bool) []*Path {
	fromNode := this.nodeMap[from]
	if fromNode == nil || this.nodeMap[to] == nil {
		return nil
	}

	var theF func(path *Path, ended bool) bool
	if len(f) != 0 {
		theF = f[0]
	}

	const defaultPathCap = 8
	// path 用来存储当前走过的路径
	path := &Path{Node: append(make([]*Node, 0, defaultPathCap), fromNode.Node), Line: make([]*Line, 0, defaultPathCap)}
	// deadNode 用来存储已经被证明无法到达 endNode 的节点。
	deadNode := make(map[string]bool, defaultPathCap)

	reached := make([]*Path, 0, 8)
	this.traversalCallback(path, deadNode, fromNode, to, &reached, theF)
	if len(reached) != 0 {
		return reached
	}
	return nil
}

// 路径搜索回调函数。
//
// 参数：
//   path:        当前已经走过的路径
//   deadNode:    已经被证明无法到达 endNode 的节点
//   node:        当前要搜索的节点
//   endNode:     结束节点
//   reachable:   能够到达目标的路径列表
//   f:           回调函数
//
// 返回值
//   reached: 是否曾经到达过 endNode
//   stopped: 由于 f 函数返回了 false，搜索被中止。
func (this *Graph) traversalCallback(path *Path, deadNode map[string]bool, node *nodeWrap, endNode string, reachable *[]*Path, f func(path *Path, ended bool) bool) (reached, stopped bool) {
	reachableLine := 0
	// 遍历当前节点的连线，判断从该连线是否可以走到 endNode
	for _, line := range node.lines {
		var another *nodeWrap
		switch line.Direction {
		case LeftToRight:
			another = line.rightNode
		case RightToLeft:
			another = line.leftNode
		case BothWay:
			if line.Left == node.Id {
				another = line.rightNode
			} else if line.Right == node.Id {
				another = line.leftNode
			}
		}
		if another == nil || another.Id == node.Id {
			continue
		} else if deadNode[another.Id] {
			// 如果节点在之前的遍历中已经被证明无法到达，则忽略
			continue
		} else if -1 != path.NodeIndex(another.Id) {
			// 如果 another 已经在当前 path 中已经走过，则忽略
			continue
		}

		// 将当前节点添加到 path 中
		path.Node = append(path.Node, another.Node)
		path.Line = append(path.Line, line.Line)

		// 判断 another 是否就是 endNode，如果是就加入到 reached 中
		isEnd := another.Id == endNode
		if isEnd {
			pathCopy := &Path{Node: make([]*Node, len(path.Node)), Line: make([]*Line, len(path.Line))}
			copy(pathCopy.Node, path.Node)
			copy(pathCopy.Line, path.Line)
			*reachable = append(*reachable, pathCopy)
			reachableLine++
		}

		if f != nil && !f(path, another.Id == endNode) {
			return reachableLine != 0, true
		}

		// 如果还没有到达 endPoint，则深度遍历
		if !isEnd {
			deepReached, stopped := this.traversalCallback(path, deadNode, another, endNode, reachable, f)
			if deepReached {
				reachableLine++
			}
			if stopped {
				return reachableLine != 0, true
			}
		}

		// 将当前节点从 path 中移除、使其恢复深度遍历前的状态。
		path.Node = path.Node[:len(path.Node)-1]
		path.Line = path.Line[:len(path.Line)-1]
	}
	if reachableLine == 0 {
		deadNode[node.Id] = true
		return false, false
	}
	return true, false
}

func (this *Path) NodeIndex(id string) int {
	for i, v := range this.Node {
		if v.Id == id {
			return i
		}
	}
	return -1
}

// 拼接路径上所有的 ‘节点 ID’，参数可以指定分隔符。
func (this *Path) ToString(sep ...string) string {
	return strUtil.Join(this.Node, util.IfEmptyStringSlice(sep, ","), func(i int) string { return this.Node[i].Id })
}

func (this *Path) String() string {
	return this.ToString()
}

func (this Direction) ToString() string {
	switch this {
	default:
		return "-->"
	case BothWay:
		return "<->"
	case RightToLeft:
		return "<--"
	}
}
