package hy

import (
	"reflect"

	"github.com/pkg/errors"
)

type analysisInContext struct {
	Type reflect.Type
	Path string
}

type nodeInContext struct {
	Node    *Node
	Context NodeContext
}

// analyses is global state and should probably be removed
// using Node pointer so we can all refer to the same Node
var analyses = map[analysisInContext]nodeInContext{}

func getAnalysisForType(t reflect.Type) (*Node, bool) {
	for aic, nic := range analyses {
		if aic.Type == t {
			return nic.Node, true
		}
	}
	return nil, false
}

func registerAnalysis(t reflect.Type, n *Node, c NodeContext) {
	analyses[analysisInContext{t, c.Path()}] = nodeInContext{n, c}
}

var nodesByType = map[reflect.Type]*Node{}

func registerNodeForType(t reflect.Type, n *Node) {
	nodesByType[t] = n
}
func getNodeForType(t reflect.Type) (*Node, bool) {
	n, ok := nodesByType[t]
	return n, ok
}

var paths = map[string]Node{}

// Analyse analyses a tree starting at root.
func Analyse(root interface{}) (Node, error) {
	if root == nil {
		return nil, errors.New("cannot analyse nil")
	}
	return analyse(reflect.TypeOf(root))
}

func analyse(t reflect.Type) (Node, error) {
	if analysis, ok := getNodeForType(t); ok {
		return *analysis, nil
	}
	// allocate the node now, child nodes can refer to it
	n := new(Node)
	registerNodeForType(t, n)
	var isPtr bool
	k := t.Kind()
	if k == reflect.Ptr {
		isPtr = true
		t = t.Elem()
		k = t.Kind()
	}
	var err error
	base := NodeBase{
		OwnType: t,
		IsPtr:   isPtr,
	}
	switch k {
	default:
		return nil, errors.Errorf("cannot analyse %s (%T)", k, t)
	case reflect.Struct:
		*n, err = analyseStruct(base, t, isPtr)
	case reflect.Map:
		*n, err = analyseMap(base, t, isPtr)
	case reflect.Slice:
		*n, err = analyseSlice(base, t, isPtr)
	}
	return *n, errors.Wrapf(err, "analysing %s failed", t)
}
