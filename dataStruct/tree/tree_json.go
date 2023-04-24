package tree

import (
	"fmt"
	"github.com/3th1nk/easygo/util/jsonUtil"
	jsonIter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"reflect"
	"unsafe"
)

var (
	nodeDataDecodeFunc = make(map[string]func(bytes []byte) (interface{}, error))
)

func RegisterNodeDataDecoder(typ string, f func(bytes []byte) (interface{}, error)) {
	nodeDataDecodeFunc[typ] = f
}

func init() {
	treeType := reflect.TypeOf(Tree{}).String()
	jsonIter.RegisterTypeEncoderFunc(treeType, func(ptr unsafe.Pointer, stream *jsonIter.Stream) {
		if data, err := getTreeJsonData((*Tree)(ptr)); err != nil {
			stream.Error = err
		} else {
			stream.WriteVal(data)
		}
	}, func(ptr unsafe.Pointer) bool {
		return reflect2.IsNil((*Tree)(ptr))
	})
	jsonIter.RegisterTypeDecoderFunc(treeType, func(ptr unsafe.Pointer, iter *jsonIter.Iterator) {
		var data *jsonTree
		if iter.ReadVal(&data); data != nil {
			if err := data.toTree((*Tree)(ptr)); err != nil {
				iter.ReportError("decode tree", err.Error())
			}
		}
	})

	nodeType := reflect.TypeOf(Node{}).String()
	jsonIter.RegisterTypeEncoderFunc(nodeType, func(ptr unsafe.Pointer, stream *jsonIter.Stream) {
		if data, err := getJsonNode((*Node)(ptr)); err != nil {
			stream.Error = err
		} else {
			stream.WriteVal(data)
		}
	}, func(ptr unsafe.Pointer) bool {
		return reflect2.IsNil((*Node)(ptr))
	})
	jsonIter.RegisterTypeDecoderFunc(nodeType, func(ptr unsafe.Pointer, iter *jsonIter.Iterator) {
		var data *jsonNode
		if iter.ReadVal(&data); data != nil {
			if err := data.toNode((*Node)(ptr)); err != nil {
				iter.ReportError("decode node", err.Error())
			}
		}
	})
}

type jsonTree struct {
	Name string `json:"name,omitempty"`
	*jsonNode
}

type jsonNode struct {
	Id    string              `json:"id,omitempty"`
	Type  string              `json:"type,omitempty"`
	Data  jsonIter.RawMessage `json:"data,omitempty"`
	Nodes []*jsonNode         `json:"nodes,omitempty"`
}

func getTreeJsonData(tree *Tree) (*jsonTree, error) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()

	root, err := getJsonNode(tree.root)
	if err != nil {
		return nil, err
	}
	return &jsonTree{Name: tree.Name, jsonNode: root}, nil
}

func (this *jsonTree) toTree(p *Tree) error {
	p.ensureInit()
	p.Name = this.Name
	if err := this.jsonNode.toNode(p.root); err != nil {
		return err
	}
	return nil
}

func getJsonNode(node *Node) (*jsonNode, error) {
	if node == nil {
		return nil, nil
	}

	obj := &jsonNode{
		Id:    node.id,
		Type:  node.typ,
		Nodes: make([]*jsonNode, len(node.nodes)),
	}
	if !reflect2.IsNil(node.data) {
		var err error
		if obj.Data, err = jsonUtil.Marshal(node.data); err != nil {
			return nil, fmt.Errorf("marshal node '%v' data error: %v", node.id, err)
		}
	}
	for i, v := range node.nodes {
		if val, err := getJsonNode(v); err != nil {
			return nil, err
		} else {
			obj.Nodes[i] = val
		}
	}
	return obj, nil
}

func (this *jsonNode) toNode(p *Node) error {
	p.id = this.Id
	p.typ = this.Type
	// Data
	if len(this.Data) == 0 {
		p.data = nil
	} else if f := nodeDataDecodeFunc[this.Type]; f != nil {
		if val, err := f(this.Data); err != nil {
			return err
		} else {
			p.data = val
		}
	} else {
		p.data = jsonUtil.Get(this.Data).GetInterface()
	}
	// Nodes
	if n := len(this.Nodes); n == 0 {
		p.nodes = nil
	} else {
		p.nodes = make([]*Node, n)
		for i, v := range this.Nodes {
			obj := &Node{tree: p.tree, pid: p.id}
			if err := v.toNode(obj); err != nil {
				return err
			} else {
				p.nodes[i] = obj
				if p.tree != nil {
					p.tree.nodeMap[obj.id] = obj
				}
			}
		}
	}
	//
	return nil
}
