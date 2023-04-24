package graph

import (
	"fmt"
	"github.com/3th1nk/easygo/internal/_test"
	"github.com/3th1nk/easygo/util"
	"github.com/3th1nk/easygo/util/jsonUtil"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestGraph_DecodeData(t *testing.T) {
	g := New()
	g.AddNode(&Node{
		Id:   "n1",
		Type: "a",
		Data: map[string]interface{}{"id": 123},
	}, &Node{
		Id:    "n2",
		Level: 0,
		Type:  "b",
		Data:  map[string]interface{}{"name": "abc"},
	})
	g.AddLine(&Line{
		Id:        "l1",
		Left:      "n1",
		Right:     "n2",
		Direction: LeftToRight,
		Type:      "a",
		Data:      map[string]interface{}{"a": "AAA"},
	}, &Line{
		Id:        "l2",
		Left:      "n1",
		Right:     "n2",
		Direction: RightToLeft,
		Type:      "b",
		Data:      map[string]interface{}{"b": "BBB"},
	})
	str := jsonUtil.MustMarshalToStringIndent(g)

	var g2 *Graph
	if assert.NoError(t, jsonUtil.UnmarshalFromString(str, &g2)) && assert.NotNil(t, g2) {
		if node := g2.GetNode("n1"); assert.NotNil(t, node) {
			assert.Equal(t, 123, node.Data.MustGetInt("id"))
		}
		if node := g2.GetNode("n2"); assert.NotNil(t, node) {
			assert.Equal(t, "abc", node.Data.MustGetString("name"))
		}
		if line := g2.GetLine("l1"); assert.NotNil(t, line) {
			assert.Equal(t, "AAA", line.Data.MustGetString("a"))
		}
		if line := g2.GetLine("l2"); assert.NotNil(t, line) {
			assert.Equal(t, "BBB", line.Data.MustGetString("b"))
		}
	}

	util.Println(str)
}

func TestGraph_JsonMarshal(t *testing.T) {
	a := &struct {
		Graph *Graph `json:"graph,omitempty"`
	}{}
	var str string

	str = jsonUtil.MustMarshalToString(a)
	assert.Equal(t, `{}`, str)

	str = jsonUtil.NoOmitemptyApi().MustMarshalToString(a)
	assert.Equal(t, `{"graph":null}`, str)

	a.Graph = newTestGraph()
	str = jsonUtil.MustMarshalToStringIndent(a)
	util.Println(str)
}

func TestGraph_JsonUnmarshal(t *testing.T) {
	a := newTestGraph()

	str := jsonUtil.MustMarshalToString(a)

	var b *Graph
	assert.NoError(t, jsonUtil.UnmarshalFromString(str, &b))

	assert.Equal(t, len(a.nodeMap), len(b.nodeMap))
	for k, v1 := range a.nodeMap {
		v2 := b.nodeMap[k]
		assert.Equal(t, jsonUtil.MustMarshalToString(v1.Node), jsonUtil.MustMarshalToString(v2.Node))
	}
	assert.Equal(t, len(a.lineMap), len(b.lineMap))
	for k, v1 := range a.lineMap {
		v2 := b.lineMap[k]
		assert.Equal(t, jsonUtil.MustMarshalToString(v1), jsonUtil.MustMarshalToString(v2))
	}
	assert.Equal(t, len(a.comboMap), len(b.comboMap))
	for k, v1 := range a.comboMap {
		v2 := b.comboMap[k]
		assert.Equal(t, jsonUtil.MustMarshalToString(v1), jsonUtil.MustMarshalToString(v2))
	}
}

func TestGraph_Traversal_Perf(t *testing.T) {
	g := newTestGraph()
	_test.Perf(func(i int) {
		g.Traversal("1", "5")
	})
}

func TestGraph_Traversal_1(t *testing.T) {
	path, _, lastNode := testTraversal("1", "5")

	assert.Equal(t, 4, len(path))
	assert.Contains(t, path, "1,2,3,5")
	assert.Contains(t, path, "1,2,3,4,5")
	assert.Contains(t, path, "1,2,4,5")
	assert.Contains(t, path, "1,2,4,3,5")

	for _, v := range strings.Split("6,7,8,9", ",") {
		n := lastNode[v]
		assert.LessOrEqual(t, n, 1)
	}
}

func TestGraph_Traversal_2(t *testing.T) {
	path, _, _ := testTraversal("6", "7")
	assert.Equal(t, 3, len(path))
	assert.Contains(t, path, "6,7")
	assert.Contains(t, path, "6,8,7")
	assert.Contains(t, path, "6,8,9,7")
}

func testTraversal(from, to string) (path, walked []string, lastNode map[string]int) {
	path, walked, lastNode = make([]string, 0, 4), make([]string, 0, 16), make(map[string]int, 10)
	newTestGraph().Traversal(from, to, func(a *Path, ended bool) bool {
		str := a.ToString()
		if ended {
			path = append(path, str)
		}
		walked = append(walked, fmt.Sprintf("%v, %v", str, ended))
		lastId := a.Node[len(a.Node)-1].Id
		lastNode[lastId] = lastNode[lastId] + 1
		return true
	})
	util.Println("============= paths: [%v]", len(path))
	for _, s := range path {
		util.Println(s)
	}
	util.Println("============= walked: [%v]", len(walked))
	for _, s := range walked {
		util.Println(s)
	}
	return
}

//         1
//     <--
//   <--
// 2      <->      3      -->      6      -->      7
//   -->       <->   -->             <->       <->   <--
//     -->   <->       -->             <->   <->       <--
//         4     -->       5               8      -->      9
func newTestGraph() *Graph {
	g := New()
	g.AddNode(
		&Node{Id: "1", Level: 1, Type: "a", Data: map[string]interface{}{"a": 1}},
		&Node{Id: "2", Level: 2, Type: "a", ComboId: "aaa", Data: map[string]interface{}{"a": 2}},
		&Node{Id: "3", Level: 2, Type: "a", ComboId: "aaa", Data: map[string]interface{}{"b": 3}},
		&Node{Id: "4", Level: 3, Type: "b", ComboId: "aaa"},
		&Node{Id: "5", Level: 3, Type: "b", ComboId: "aaa"},
		&Node{Id: "6", Level: 2, Type: "b", ComboId: "bbb"},
		&Node{Id: "7", Level: 2, Type: "b", ComboId: "bbb"},
		&Node{Id: "8", Level: 3, Type: "b", ComboId: "bbb"},
		&Node{Id: "9", Level: 3, Type: "b", ComboId: "bbb"},
	)
	g.AddLine(
		&Line{Left: "1", Right: "2", Direction: LeftToRight},
		&Line{Left: "2", Right: "3", Direction: BothWay},
		&Line{Left: "2", Right: "4", Direction: LeftToRight},
		&Line{Left: "3", Right: "4", Direction: BothWay},
		&Line{Left: "3", Right: "5", Direction: LeftToRight},
		&Line{Left: "3", Right: "6", Direction: LeftToRight},
		&Line{Left: "4", Right: "5", Direction: LeftToRight},
		&Line{Left: "6", Right: "7", Direction: LeftToRight},
		&Line{Left: "6", Right: "8", Direction: BothWay},
		&Line{Left: "7", Right: "8", Direction: BothWay},
		&Line{Left: "7", Right: "9", Direction: RightToLeft},
		&Line{Left: "8", Right: "9", Direction: LeftToRight},
	)
	g.AddCombo(
		&Combo{Id: "aaa", Type: "a"},
		&Combo{Id: "bbb", Type: "b"},
	)
	return g
}
