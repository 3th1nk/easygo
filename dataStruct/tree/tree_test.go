package tree

import (
	"github.com/3th1nk/easygo/util"
	"github.com/3th1nk/easygo/util/jsonUtil"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func Test_MarshalJsonNode(t *testing.T) {
	assertNode0 := func(node *Node) {
		assert.Equal(t, node.id, "1")
		obj, _ := node.data.(*nodeA)
		if assert.NotNil(t, obj) {
			assert.Equal(t, 1, obj.A)
		}

		if assert.Equal(t, 2, len(node.nodes)) {
			assert.Equal(t, "1.1", node.nodes[0].id)
			assert.Equal(t, 11, node.nodes[0].data.(*nodeA).A)

			assert.Equal(t, "1.2", node.nodes[1].id)
			assert.Equal(t, "aaa", node.nodes[1].data.(*nodeB).B)
		}
	}

	tree := newTestTree()
	var node *Node
	str := jsonUtil.MustMarshalToString(tree.root.nodes[0])
	if assert.NoError(t, jsonUtil.UnmarshalFromString(str, &node)) {
		assertNode0(node)
	}

	var tree2 *Tree
	str = jsonUtil.MustMarshalToString(tree)
	if !assert.NoError(t, jsonUtil.UnmarshalFromString(str, &tree2)) || !assert.NotNil(t, tree2) {
		return
	}
	if !assert.NotNil(t, tree2.root) || !assert.Equal(t, 3, len(tree2.root.nodes)) {
		return
	}
	node = nil
	str = jsonUtil.MustMarshalToString(tree2.root.nodes[0])
	if assert.NoError(t, jsonUtil.UnmarshalFromString(str, &node)) {
		assertNode0(node)
	}
}

func Test_MarshalJsonTree(t *testing.T) {
	tree := newTestTree()
	var tree2 *Tree
	if !assert.NoError(t, jsonUtil.UnmarshalFromString(jsonUtil.MustMarshalToString(tree), &tree2)) || !assert.NotNil(t, tree2) {
		return
	}

	nodeId1, nodeId2 := make([]string, 0, 8), make([]string, 0, 8)
	tree.Each(OrderPre, func(node *Node, depth int) bool {
		nodeId1 = append(nodeId1, node.id)
		return true
	})
	tree2.Each(OrderPre, func(node *Node, depth int) bool {
		nodeId2 = append(nodeId2, node.id)
		return true
	})
	if assert.NotEqual(t, 0, len(nodeId1)) && assert.Equal(t, len(nodeId1), len(nodeId2)) {
		for i, v := range nodeId1 {
			assert.Equal(t, v, nodeId2[i])
		}
	}
}

func Test_Each(t *testing.T) {
	nodeIds := make([]string, 0, 16)
	newTestTree().Each(OrderPre, func(node *Node, depth int) bool {
		nodeIds = append(nodeIds, node.id)
		return true
	})

	expected := []string{
		"",
		"1",
		"1.1",
		"1.2",
		"2",
		"2.1",
		"2.2",
		"2.2.1",
		"2.2.2",
		"2.2.3",
		"3",
		"3.1",
	}
	if assert.Equal(t, len(expected), len(nodeIds)) {
		for i, s := range expected {
			assert.Equal(t, s, nodeIds[i])
		}
	}
}

func newTestTree() *Tree {
	tree := &Tree{}
	var node *Node

	node = tree.Root().MustAddChild("1", "a", &nodeA{A: 1})
	node.MustAddChild("1.1", "a", &nodeA{A: 11})
	node.MustAddChild("1.2", "b", &nodeB{B: "aaa"})
	node = tree.Root().MustAddChild("2", "a", &nodeA{A: 2})
	node.MustAddChild("2.1", "a", &nodeA{A: 21})
	node = node.MustAddChild("2.2", "b", &nodeB{B: "bbb"})
	node.MustAddChild("2.2.1", "b", &nodeB{B: "bbb.1"})
	node.MustAddChild("2.2.2", "b", &nodeB{B: "bbb.2"})
	node.MustAddChild("2.2.3", "b", &nodeB{B: "bbb.3"})
	node = tree.Root().MustAddChild("3", "a", &nodeA{A: 3})
	node.MustAddChild("3.1", "a", &nodeA{A: 31})

	return tree
}

func init() {
	RegisterNodeDataDecoder("a", func(bytes []byte) (interface{}, error) {
		var obj *nodeA
		err := jsonUtil.Unmarshal(bytes, &obj)
		return obj, err
	})
	RegisterNodeDataDecoder("b", func(bytes []byte) (interface{}, error) {
		var obj *nodeB
		err := jsonUtil.Unmarshal(bytes, &obj)
		return obj, err
	})
}

type nodeA struct {
	A int `json:"a,omitempty"`
}

func (this *nodeA) GetNodeType() string {
	return "a"
}

func (this *nodeA) GetNodeId() string {
	return strconv.Itoa(this.A)
}

type nodeB struct {
	B string `json:"b,omitempty"`
}

func (this *nodeB) GetNodeId() string {
	return this.B
}

func (this *nodeB) GetNodeType() string {
	return "b"
}

func TestEachOrder_String(t *testing.T) {
	util.Println("%v", OrderPre)
}
