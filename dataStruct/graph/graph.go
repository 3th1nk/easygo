package graph

import (
	"fmt"
	"github.com/3th1nk/easygo/util"
	"github.com/3th1nk/easygo/util/strUtil"
)

// Graph 图结构
type Graph struct {
	nodeMap  map[string]*nodeWrap  // id -> Node
	lineMap  map[string]*lineWrap  // id -> Line
	comboMap map[string]*comboWrap // id -> Combo
}

// Path 路径
type Path struct {
	Node []*Node `json:"node,omitempty"` // 该路径上的节点列表
	Line []*Line `json:"line,omitempty"` // 该路径上的连线列表
}

// New 创建一个图
func New() *Graph {
	return &Graph{}
}

func (this *Graph) init(renew bool) {
	if this.nodeMap == nil || renew {
		this.nodeMap = make(map[string]*nodeWrap, 32)
		this.lineMap = make(map[string]*lineWrap, 32)
		this.comboMap = make(map[string]*comboWrap, 8)
	}
}

// NodeCount 获取节点数量
func (this *Graph) NodeCount() int {
	return len(this.nodeMap)
}

// LineCount 获取连线数量
func (this *Graph) LineCount() int {
	return len(this.lineMap)
}

// ComboCount 获取组合数量
func (this *Graph) ComboCount() int {
	return len(this.comboMap)
}

// AddNode 添加节点
//
//   注意：在变更数据之后，需要调用 Update 更新图。
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
		this.nodeMap[v.Id] = &nodeWrap{Node: v, lines: make(map[string]*lineWrap)}
	}
	return nil
}

// AddLine 添加连线
//
//   注意：在变更数据之后，需要调用 Update 更新图。
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
		this.lineMap[v.Id] = &lineWrap{Line: v}
	}
	return nil
}

// AddCombo 添加组合框
//
//   注意：在变更数据之后，需要调用 Update 更新图。
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
		this.comboMap[v.Id] = &comboWrap{
			Combo:  v,
			nodes:  make(map[string]*nodeWrap),
			combos: make(map[string]*comboWrap),
		}
	}
	return nil
}

// DeleteNode 删除指定的节点
//   注意：在变更数据之后，需要调用 Update 更新图。
func (this *Graph) DeleteNode(ids ...string) {
	for _, id := range ids {
		delete(this.nodeMap, id)
	}
}

// DeleteLine 删除指定的连线
//   注意：在变更数据之后，需要调用 Update 更新图。
func (this *Graph) DeleteLine(ids ...string) {
	for _, id := range ids {
		delete(this.lineMap, id)
	}
}

// DeleteCombo 删除指定的组合
//   注意：在变更数据之后，需要调用 Update 更新图。
func (this *Graph) DeleteCombo(ids ...string) {
	for _, id := range ids {
		delete(this.comboMap, id)
	}
}

func (this *Graph) updateNode(node *nodeWrap) error {
	if node.ComboId != "" {
		if v, ok := this.comboMap[node.ComboId]; !ok {
			return fmt.Errorf("node combo '%v.%v' not exists", node.Id, node.ComboId)
		} else {
			v.nodes[node.Id] = node
		}
	}
	for _, line := range node.lines {
		if _, ok := this.lineMap[line.Id]; !ok {
			delete(node.lines, line.Id)
		}
	}
	return nil
}

func (this *Graph) updateLine(line *lineWrap) error {
	if line.Left == "" || line.Right == "" {
		return fmt.Errorf("line '%v(%v %v %v)' missing node", line.Id, line.Left, line.Direction, line.Right)
	}

	if v, ok := this.nodeMap[line.Left]; !ok {
		return fmt.Errorf("line left '%v.%v' not exists", line.Id, line.Left)
	} else {
		line.leftNode = v
		v.lines[line.Id] = line
	}

	if v, ok := this.nodeMap[line.Right]; !ok {
		return fmt.Errorf("line right '%v.%v' not exists", line.Id, line.Right)
	} else {
		line.rightNode = v
		v.lines[line.Id] = line
	}

	return nil
}

func (this *Graph) updateCombo(combo *comboWrap) error {
	if combo.ComboId != "" {
		if v, ok := this.comboMap[combo.ComboId]; !ok {
			return fmt.Errorf("combo parent '%v.%v' not exists", combo.Id, combo.ComboId)
		} else {
			v.combos[combo.Id] = combo
		}
	}
	for _, c := range combo.combos {
		if v, ok := this.comboMap[c.Id]; !ok || v.ComboId != combo.Id {
			delete(combo.combos, c.Id)
		}
	}
	for _, n := range combo.nodes {
		if v, ok := this.nodeMap[n.Id]; !ok || v.ComboId != combo.Id {
			delete(combo.nodes, n.Id)
		}
	}
	return nil
}

// Update 更新图结构，并校验数据，确认数据格式正确。
//
// 	当变更了节点、连线、组合时调用。
//
// 返回：
//   err: 当校验失败时返回对应的错误：
//     * 节点的 ComboId 非空、但不存在对应的 Combo;
//     * 节点的 ComboId 非空、但不存在对应的 Combo;
func (this *Graph) Update() (err error) {
	return this.doUpdate()
}

func (this *Graph) doUpdate() (err error) {
	for _, line := range this.lineMap {
		if err = this.updateLine(line); err != nil {
			return err
		}
	}

	for _, node := range this.nodeMap {
		if err = this.updateNode(node); err != nil {
			return err
		}
	}

	for _, combo := range this.comboMap {
		if err = this.updateCombo(combo); err != nil {
			return err
		}
	}

	return nil
}

// 根据 ‘节点 ID’ 获取节点信息
func (this *Graph) GetNode(id string) *Node {
	if v := this.nodeMap[id]; v != nil {
		return v.Node
	}
	return nil
}

// 根据 ‘连线 ID’ 获取连线信息
func (this *Graph) GetLine(id string) *Line {
	if v := this.lineMap[id]; v != nil {
		return v.Line
	}
	return nil
}

// 根据 ‘组合 ID’ 获取组合信息
func (this *Graph) GetCombo(id string) *Combo {
	if v := this.comboMap[id]; v != nil {
		return v.Combo
	}
	return nil
}

func (this *Graph) FindNode(f ...func(node *Node) bool) []*Node {
	arr := make([]*Node, 0, len(this.nodeMap))
	for _, v := range this.nodeMap {
		if len(f) == 0 || f[0] == nil || f[0](v.Node) {
			arr = append(arr, v.Node)
		}
	}
	if len(arr) == 0 {
		return nil
	}
	return arr
}

func (this *Graph) FindLine(f ...func(line *Line, left, right *Node) bool) []*Line {
	arr := make([]*Line, 0, len(this.lineMap))
	for _, v := range this.lineMap {
		if len(f) == 0 || f[0] == nil || f[0](v.Line, v.safeLeft(), v.safeRight()) {
			arr = append(arr, v.Line)
		}
	}
	if len(arr) == 0 {
		return nil
	}
	return arr
}

func (this *Graph) FindCombo(f ...func(combo *Combo) bool) []*Combo {
	arr := make([]*Combo, 0, len(this.comboMap))
	for _, v := range this.comboMap {
		if len(f) == 0 || f[0] == nil || f[0](v.Combo) {
			arr = append(arr, v.Combo)
		}
	}
	if len(arr) == 0 {
		return nil
	}
	return arr
}

// 根据 ‘节点 ID’，获取与该节点直接关联的连线
func (this *Graph) GetNodeLines(id string) []*Line {
	if node := this.nodeMap[id]; node == nil {
		return nil
	} else {
		return node.GetLines()
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
