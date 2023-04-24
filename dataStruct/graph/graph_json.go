package graph

import (
	jsonIter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"reflect"
	"sort"
	"strings"
	"unsafe"
)

func init() {
	graphType := reflect.TypeOf(Graph{}).String()
	jsonIter.RegisterTypeEncoderFunc(graphType, func(ptr unsafe.Pointer, stream *jsonIter.Stream) {
		stream.WriteVal(toJsonData((*Graph)(ptr)))
	}, func(ptr unsafe.Pointer) bool {
		return reflect2.IsNil((*Graph)(ptr))
	})
	jsonIter.RegisterTypeDecoderFunc(graphType, func(ptr unsafe.Pointer, iter *jsonIter.Iterator) {
		var data *graphJsonData
		if iter.ReadVal(&data); data != nil {
			if err := fromJsonData((*Graph)(ptr), data); err != nil {
				iter.ReportError("decode graph", err.Error())
			}
		}
	})
}

type graphJsonData struct {
	Node  []*Node  `json:"node,omitempty"`
	Line  []*Line  `json:"line,omitempty"`
	Combo []*Combo `json:"combo,omitempty"`
}

func toJsonData(graph *Graph) *graphJsonData {
	data, idx := &graphJsonData{}, 0

	//
	data.Node, idx = make([]*Node, len(graph.nodeMap)), 0
	for _, item := range graph.nodeMap {
		data.Node[idx], idx = item.Node, idx+1
	}
	sort.Slice(data.Node, func(i, j int) bool {
		a, b := data.Node[i], data.Node[j]
		if n := a.Level - b.Level; n != 0 {
			return n < 0
		}
		return a.Id < b.Id
	})

	//
	data.Line, idx = make([]*Line, len(graph.lineMap)), 0
	for _, item := range graph.lineMap {
		data.Line[idx], idx = item.Line, idx+1
	}
	sort.Slice(data.Line, func(i, j int) bool {
		a, b := data.Line[i], data.Line[j]
		if n := strings.Compare(a.Left, b.Left); n != 0 {
			return n < 0
		} else if n := strings.Compare(a.Right, b.Right); n != 0 {
			return n < 0
		}
		return a.Id < b.Id
	})

	//
	data.Combo, idx = make([]*Combo, len(graph.comboMap)), 0
	for _, item := range graph.comboMap {
		data.Combo[idx], idx = item, idx+1
	}
	sort.Slice(data.Combo, func(i, j int) bool {
		return data.Combo[i].Id < data.Combo[j].Id
	})

	return data
}

func fromJsonData(g *Graph, data *graphJsonData) (err error) {
	g.init(true)
	if err = g.AddNode(data.Node...); err == nil {
		if err = g.AddLine(data.Line...); err == nil {
			if err = g.AddCombo(data.Combo...); err == nil {
				err = g.Update()
			}
		}
	}
	return
}
