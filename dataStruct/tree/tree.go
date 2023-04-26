package tree

import (
	"fmt"
	"github.com/3th1nk/easygo/util"
	"github.com/3th1nk/easygo/util/mathUtil"
	"github.com/3th1nk/easygo/util/strUtil"
	"sync"
)

type Tree struct {
	Name string `json:"name,omitempty"`

	mu      sync.RWMutex     //
	root    *Node            //
	nodeMap map[string]*Node //
}

type Node struct {
	id    string
	pid   string
	typ   string
	data  interface{}
	tree  *Tree
	nodes []*Node
}

type NodeData struct {
	Id   string
	Type string
	Data interface{}
}

//go:generate stringer -type EachOrder -trimprefix Order
type EachOrder int

const (
	OrderPre  EachOrder = 1 // 前序遍历，先遍历根节点，然后再遍历子节点
	OrderPost EachOrder = 2 // 后序遍历，先遍历子节点，然后再遍历根节点
)

func (this *Tree) ensureInit() *Tree {
	if this.root == nil {
		this.mu.Lock()
		defer this.mu.Unlock()
		if this.root == nil {
			this.root = &Node{tree: this}
			this.nodeMap = make(map[string]*Node, 16)
		}
	}
	return this
}

func (this *Tree) Root() *Node {
	return this.ensureInit().root
}

func (this *Tree) GetNode(id string) (node *Node) {
	if id == "" {
		return this.ensureInit().root
	} else if this.nodeMap != nil {
		return this.nodeMap[id]
	}
	return nil
}

func (this *Tree) MustAddNode(pid, id string, typ string, data interface{}) (added *Node) {
	v, _ := this.AddNode(pid, id, typ, data)
	return v
}

func (this *Tree) AddNode(pid, id string, typ string, data interface{}) (added *Node, err error) {
	parent := this.GetNode(pid)
	if parent == nil {
		return nil, fmt.Errorf("pid '%v' not exists", pid)
	}
	return parent.AddChild(id, typ, data)
}

// 遍历节点及其所有子孙节点。
//   order: 排序方式。 PreOrder=前序遍历; PostOrder=后序遍历。
//   f: 回调函数，如果 f 返回 false，则会中止遍历。
//      depth: 当前遍历深度，从 0 开始累加。
//   maxDepth: 最大遍历深度。
func (this *Tree) Each(order EachOrder, f func(node *Node, depth int) bool, maxDepth ...int) {
	this.root.eachCallback(order, util.IfEmptyIntSlice(maxDepth, -1), 0, f)
}

func (this *Tree) Update() {
	this.mu.Lock()
	defer this.mu.Unlock()

	dict := make(map[string]*Node, mathUtil.MaxInt(len(this.nodeMap), 16))
	this.root.eachCallback(OrderPre, -1, 0, func(node *Node, depth int) (stop bool) {
		dict[node.id] = node
		return true
	})
	this.nodeMap = dict
}

func (this *Node) Id() string {
	return this.id
}

func (this *Node) Type() string {
	return this.typ
}

func (this *Node) Data() interface{} {
	return this.data
}

func (this *Node) SetData(typ string, data interface{}) *Node {
	this.typ, this.data = typ, data
	return this
}

func (this *Node) Parent() *Node {
	return this.tree.GetNode(this.id)
}

func (this *Node) Children() []*Node {
	return this.nodes
}

// 遍历节点及其所有子孙节点。
//   order: 排序方式。 PreOrder=前序遍历; PostOrder=后序遍历。
//   f: 回调函数，如果 f 返回 false，则会中止遍历。
//      depth: 当前遍历深度，从 0 开始累加。
//   maxDepth: 最大遍历深度。
func (this *Node) Each(order EachOrder, f func(node *Node, depth int) bool, maxDepth ...int) {
	this.eachCallback(order, util.IfEmptyIntSlice(maxDepth, -1), 0, f)
}

func (this *Node) eachCallback(order EachOrder, maxDepth, depth int, f func(node *Node, depth int) bool) bool {
	nextDepth := depth + 1
	if order == OrderPost {
		// 后序遍历
		if maxDepth < 0 || nextDepth < maxDepth {
			for _, v := range this.nodes {
				if !v.eachCallback(OrderPost, maxDepth, nextDepth, f) {
					return false
				}
			}
		}
		//
		if !f(this, depth) {
			return false
		}
	} else {
		// 前序遍历
		if !f(this, depth) {
			return false
		}
		//
		if maxDepth < 0 || nextDepth < maxDepth {
			for _, v := range this.nodes {
				if !v.eachCallback(OrderPre, maxDepth, nextDepth, f) {
					return false
				}
			}
		}
	}
	return true
}

func (this *Node) MustAddChild(id string, typ string, data interface{}) (added *Node) {
	added, _ = this.AddChild(id, typ, data)
	return
}

func (this *Node) AddChild(id string, typ string, data interface{}) (added *Node, err error) {
	if arr, err := this.AddChildData(NodeData{Id: id, Type: typ, Data: data}); err != nil {
		return nil, err
	} else {
		return arr[0], nil
	}
}

func (this *Node) AddChildData(data ...NodeData) (added []*Node, err error) {
	if len(data) == 0 {
		return nil, nil
	}

	this.tree.ensureInit()

	this.tree.mu.Lock()
	defer this.tree.mu.Unlock()

	if this.nodes == nil {
		this.nodes = make([]*Node, 0, 4)
	}

	added = make([]*Node, len(data))
	for i, v := range data {
		if v.Id == "" {
			v.Id = strUtil.Rand(8)
		} else if obj, _ := this.tree.nodeMap[v.Id]; obj != nil {
			return nil, fmt.Errorf("node '%v' already exists", v.Id)
		}

		node := &Node{
			id:   v.Id,
			pid:  this.id,
			typ:  v.Type,
			data: v.Data,
			tree: this.tree,
		}
		this.nodes = append(this.nodes, node)
		this.tree.nodeMap[node.id] = node
		added[i] = node
	}
	return added, nil
}
