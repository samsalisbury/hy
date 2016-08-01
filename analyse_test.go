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

var badAnalysesTable = map[string]interface{}{
	"cannot analyse nil": nil,
}

func TestAnalyse_failure(t *testing.T) {
	for expected, input := range badAnalysesTable {
		node, actualErr := Analyse(input)
		if node != nil {
			t.Errorf("got node %v; want nil", node)
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
	for expected, input := range goodAnalysesTable {
		actual, err := Analyse(input)
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
	root, err := Analyse(StructA{})
	if err != nil {
		return nil, err
	}
	structNode, ok := root.(*Struct)
	if !ok {
		return nil, fmt.Errorf("got %T; want *StructNode", root)
	}
	n, ok := structNode.Children[name]
	if !ok {
		return nil, fmt.Errorf("%T does not have a child %s", s, name)
	}
	return n, nil
}

func TestAnalyse_mapDir(t *testing.T) {
	child, err := getStructChildNode(StructA{}, "Map")
	if err != nil {
		t.Fatal(err)
	}
	mapDir, ok := child.(*MapNode)
	if !ok {
		t.Fatalf("got a %T; want *Map", child)
	}
	stringType := reflect.TypeOf("")
	if mapDir.KeyType != stringType {
		t.Fatalf("got key type %s; want %s", mapDir.KeyType, stringType)
	}
}

func (expected ExpectedStructAnalysis) Matches(n Node) error {
	actual, ok := n.(*Struct)
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
