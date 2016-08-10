package hy

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

type TestWriteStruct struct {
	Name         string              // regular field
	Int          int                 // regular field
	InlineSlice  []string            // regular field
	InlineMap    map[string]int      // regular field
	StructB      StructB             // regular field
	StructBPtr   *StructB            // regular field
	IgnoredField string              `hy:"-"`                // not output anywhere
	StructFile   StructB             `hy:"a-file"`           // a single file
	StringFile   string              `hy:"a-string-file"`    // a single file
	SliceFile    []string            `hy:"a-slice-file"`     // a single file
	MapFile      map[string]string   `hy:"a-map-file"`       // a single file
	Nested       *TestWriteStruct    `hy:"nested"`           // like a new root
	Slice        []StructB           `hy:"slice/"`           // file per element
	Map          map[string]StructB  `hy:"map/,Name"`        // file per element
	MapOfPtr     map[string]*StructB `hy:"map-of-ptr/,Name"` // file per element
}

var testWriteStructData = TestWriteStruct{
	Name:        "Test struct writing",
	Int:         1,
	InlineSlice: []string{"a", "string", "slice"},
	InlineMap:   map[string]int{"one": 1, "two": 2, "three": 3},
	StructFile:  StructB{Name: "A file"},
	StringFile:  "A string in a file.",
	Nested: &TestWriteStruct{
		Name: "A nested struct pointer.",
		Int:  2,
		Slice: []StructB{
			{Name: "Nested One"}, {Name: "Nested Two"},
		},
		Nested: &TestWriteStruct{
			SliceFile: []string{"this", "is", "a", "slice", "in", "a", "file"},
			MapFile:   map[string]string{"deeply-nested": "map", "in a file": "yes"},
		},
		StructFile: StructB{
			Name: "Struct B file",
		},
		MapOfPtr: map[string]*StructB{
			"a-nil-file":           nil,
			"another-nil-file":     nil,
			"this-one-has-a-value": &StructB{},
		},
		Map: map[string]StructB{
			// Notice how we don't set the Name field here. Hy sets it in the write
			// data because of the ",Name" tag.
			"a-zero-file":       StructB{},
			"another-zero-file": StructB{},
		},
	},
	Slice: []StructB{{Name: "One"}, {Name: "Two"}},
	Map: map[string]StructB{
		// Notice how we don't set the Name field here. Hy sets it in the write
		// data because of the ",Name" tag.
		"First":  StructB{},
		"Second": StructB{},
	},
}

// TODO:
//   - Use default path names for "." and "./" tags.
//   - Add options for default path names:
//     - lowerCamelCase
//     - CamelCase
//     - snake-case
//     - underscores_only
//     - lowercase
//     - UPPERCASE
//   - Respect JSON tags for field names.
//   - Respect YAML tags for field names?
//   - Add support for reading FileTargets.
//   - Add support for auto-filling ID fields in map/slice elements on read.
//     - Default field:  ID string
//     - Default getter: ID() string
//     - Default setter: SetID(string)
//   - On write, need to pick:
//     - Fail if ID field not matching key or index?
//     - Overwrite ID with current key or index?
//     - Elide ID field from output altogether? (This should be the default, so
//       it only matters in memory.)
//     - Other?
//   - Add support for writing special maps with default fields/methods:
//   - Add support for writing actual files with a marshaller.
//   - Add support for reading actual files with a marshaller.

func TestNode_Write_struct(t *testing.T) {
	c := NewCodec()
	n, err := c.Analyse(TestWriteStruct{})
	if err != nil {
		t.Fatal(err)
	}
	wc := NewWriteContext()
	v := reflect.ValueOf(testWriteStructData)
	if err := n.Write(wc, reflect.Value{}, v); err != nil {
		t.Fatal(err)
	}
	targets := wc.Targets
	expectedLen := len(expectedFileTargetsSnapshot)
	if targets.Len() != expectedLen {
		t.Errorf("got len %d; want %d", targets.Len(), expectedLen)
	}
	actualTargets := targets.Snapshot()
	for fileName, actual := range actualTargets {
		expected, ok := expectedFileTargetsSnapshot[fileName]
		if !ok {
			t.Errorf("extra file generated at %s:\n%s", fileName, actual.TestDump())
			continue
		}
		if actual.Data == nil && expected.Data == nil {
			continue
		}
		var actualType, expectedType reflect.Type
		if actual.Data != nil {
			actualType = reflect.ValueOf(actual.Data).Type()
			if expected.Data == nil {
				t.Errorf("at %q got: %v; want nil", fileName, actual.Data)
			}
		}
		if expected.Data != nil {
			expectedType = reflect.ValueOf(expected.Data).Type()
			if actual.Data == nil {
				t.Errorf("at %q got: nil; want:\n%v", fileName, expected.Data)
			}
		}

		if actualType != expectedType {
			t.Errorf("got type %s; want %s at %q", actualType, expectedType, fileName)
			t.Errorf("values: got:\n%# v\nwant:\n%# v", actual.Data, expected.Data)
		}
		if actual.TestDataDump() != expected.TestDataDump() {
			t.Errorf("\ngot rendered data:\n%s\nwant:\n%s\n",
				actual.TestDump(), expected.TestDump())
		}
	}
	for fileName := range expectedFileTargetsSnapshot {
		if _, ok := actualTargets[fileName]; !ok {
			t.Errorf("missing file %q", fileName)
		}
	}
}

func (ft FileTarget) TestDump() string {
	return fmt.Sprintf("file: %q\n%s\n", ft.Path, ft.TestDataDump())
}

func (ft FileTarget) TestDataDump() string {
	data, err := json.MarshalIndent(ft.Data, "  ", "  ")
	if err != nil {
		panic(err)
	}
	return string(data)
}
