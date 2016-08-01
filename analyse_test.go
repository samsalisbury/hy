package hy

import (
	"fmt"
	"reflect"
	"testing"
)

type (
	StructA struct {
		EmbeddedString string
		EmbeddedStruct StructB
		FileStruct     StructB `hy:"."`
		EmbeddedMap    map[string]string
		FileMap        map[string]StructB `hy:"."`
		Map            map[string]StructB `hy:"./"`
		EmbeddedSlice  []string
		FileSlice      []string `hy:"."`
		DirSlice       []string `hy:"./"`
	}
	StructB struct {
		Name          string
		FileSubStruct *StructB            `hy:"."`
		DirSubMap     map[string]*StructB `hy:"./"`
	}
	ExpectedStructAnalysis struct {
		IsPtr                  bool
		NumChildren, NumFields int
	}
)

var nodePtr = new(Node)

var badAnalysesTable = map[string]interface{}{
	"cannot analyse nil":                                        nil,
	"failed to analyse int: cannot analyse kind int":            1,
	"failed to analyse string: cannot analyse kind string":      "",
	"failed to analyse *hy.Node: cannot analyse kind interface": new(Node),
	"failed to analyse **hy.Node: cannot analyse kind ptr":      &nodePtr,
}

func TestAnalyse_failure(t *testing.T) {
	c := NewCodec()
	for expected, input := range badAnalysesTable {
		node, actualErr := c.Analyse(input)
		if actualErr == nil || node != nil {
			t.Errorf("got (%v, %q); want (nil, %q)", node, actualErr, expected)
			continue
		}
		actual := actualErr.Error()
		if actual != expected {
			t.Errorf("got error %q; want %q", actual, expected)
		}
	}
}

var goodAnalysesTable = map[ExpectedStructAnalysis]interface{}{
	{NumChildren: 5, NumFields: 4}:              StructA{},
	{NumChildren: 5, NumFields: 4, IsPtr: true}: &StructA{},
}

func TestAnalyse_success(t *testing.T) {
	c := NewCodec()
	for expected, input := range goodAnalysesTable {
		actual, err := c.Analyse(input)
		if err != nil {
			t.Error(err)
			continue
		}
		if err := expected.Matches(actual); err != nil {
			t.Error(err)
		}
	}
}

func getStructChildNode(s interface{}, name string) (Node, error) {
	c := NewCodec()
	root, err := c.Analyse(StructA{})
	if err != nil {
		return nil, err
	}
	structNode, ok := root.(*StructNode)
	if !ok {
		return nil, fmt.Errorf("got %T; want *StructNode", root)
	}
	n, ok := structNode.Children[name]
	if !ok {
		return nil, fmt.Errorf("%T does not have a child %s", s, name)
	}
	return *n, nil
}

func TestCodec_Analyse_mapNode(t *testing.T) {
	child, err := getStructChildNode(StructA{}, "Map")
	if err != nil {
		t.Fatal(err)
	}
	mapNode, ok := child.(*MapNode)
	if !ok {
		t.Fatalf("got a %T; want *MapNode", child)
	}
	stringType := reflect.TypeOf("")
	if mapNode.KeyType != stringType {
		t.Fatalf("got key type %s; want %s", mapNode.KeyType, stringType)
	}
}

func (expected ExpectedStructAnalysis) Matches(n Node) error {
	actual, ok := n.(*StructNode)
	if !ok {
		return fmt.Errorf("got a %T; want *StructNode", n)
	}
	if actual.Parent != nil {
		return fmt.Errorf("Parent was %v; want nil", actual.Parent)
	}
	if len(actual.Fields) != expected.NumFields {
		return fmt.Errorf("len(Fields) == %d; want %d", len(actual.Fields), expected.NumFields)
	}
	if len(actual.Children) != expected.NumChildren {
		return fmt.Errorf("len(Children) == %d; want %d", len(actual.Children), expected.NumChildren)
	}
	return nil
}
