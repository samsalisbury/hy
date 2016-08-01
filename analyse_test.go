package hy

import (
	"fmt"
	"testing"
)

type (
	StructA struct {
		EmbeddedString string
		EmbeddedStruct StructB
		FileStruct     StructB `hy:"."`
		EmbeddedMap    map[string]string
		FileMap        map[string]StructB `hy:"."`
		DirMap         map[string]StructB `hy:"./"`
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
