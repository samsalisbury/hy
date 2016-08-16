package hy

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

type TestStruct struct {
	Name                string                      // regular field
	Int                 int                         // regular field
	InlineSlice         []string                    // regular field
	InlineMap           map[string]int              // regular field
	StructB             StructB                     // regular field
	StructBPtr          *StructB                    // regular field
	IgnoredField        string                      `hy:"-"`                 // not output anywhere
	StructFile          StructB                     `hy:"a-file"`            // a single file
	StringFile          string                      `hy:"a-string-file"`     // a single file
	SliceFile           []string                    `hy:"a-slice-file"`      // a single file
	MapFile             map[string]string           `hy:"a-map-file"`        // a single file
	Nested              *TestStruct                 `hy:"nested"`            // like a new root
	Slice               []StructB                   `hy:"slice/"`            // file per element
	Map                 map[string]StructB          `hy:"map/,Name"`         // file per element
	MapOfPtr            map[string]*StructB         `hy:"map-of-ptr/,Name"`  // file per element
	TextMarshalerKey    map[TextMarshaler]*StructB  `hy:"textmarshaler/"`    // file per element
	TextMarshalerPtrKey map[*TextMarshaler]*StructB `hy:"textmarshalerptr/"` // file per element
	SpecialMap          SpecialMap                  `hy:"specialmap/"`       // file per element
}

type SpecialMap struct {
	m map[TextMarshaler]*StructB
}

func (s SpecialMap) SetAll(m map[TextMarshaler]*StructB) { s.m = m }
func (s SpecialMap) GetAll() map[TextMarshaler]*StructB  { return s.m }

type TextMarshaler struct {
	String string
	Int    int
}

func (tm TextMarshaler) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%s-%d", tm.String, tm.Int)), nil
}

func (tm *TextMarshaler) UnmarshalText(text []byte) error {
	s := string(text)
	s = strings.Replace(s, "-", " ", -1)
	n, err := fmt.Sscanf(s, "%s %d", &tm.String, &tm.Int)
	if err != nil && err != io.EOF {
		return errors.Wrapf(err, "unmarshaling %s", s)
	}
	if n != 2 {
		return errors.Errorf("%s has %d missing fields", s, 2-n)
	}
	return nil
}

var testWriteStructData = TestStruct{
	Name:        "Test struct writing",
	Int:         1,
	InlineSlice: []string{"a", "string", "slice"},
	InlineMap:   map[string]int{"one": 1, "two": 2, "three": 3},
	StructFile:  StructB{Name: "A file"},
	StringFile:  "A string in a file.",
	Nested: &TestStruct{
		Name: "A nested struct pointer.",
		Int:  2,
		Slice: []StructB{
			{Name: "Nested One"}, {Name: "Nested Two"},
		},
		Nested: &TestStruct{
			SliceFile: []string{"this", "is", "a", "slice", "in", "a", "file"},
			MapFile:   map[string]string{"deeply-nested": "map", "in a file": "yes"},
		},
		StructFile: StructB{
			Name: "Struct B file",
		},
		MapOfPtr: map[string]*StructB{
			"a-nil-file":           nil,
			"another-nil-file":     nil,
			"this-one-has-a-value": {},
		},
		Map: map[string]StructB{
			// Notice how we don't set the Name field here. Hy sets it in the write
			// data because of the ",Name" tag.
			"a-zero-file":       {},
			"another-zero-file": {},
		},
	},
	Slice: []StructB{{Name: "One"}, {Name: "Two"}},
	Map: map[string]StructB{
		// Notice how we don't set the Name field here. Hy sets it in the write
		// data because of the ",Name" tag.
		"First":  {},
		"Second": {},
	},
	TextMarshalerKey: map[TextMarshaler]*StructB{
		{"Test", 1}:     nil,
		{"Another", 13}: nil,
	},
	TextMarshalerPtrKey: map[*TextMarshaler]*StructB{
		{"Test", 2}:     nil,
		{"Another", 14}: nil,
	},
	SpecialMap: SpecialMap{m: map[TextMarshaler]*StructB{
		{"Special", 3}:  {Name: "Special"},
		{"Another", 15}: nil,
	}},
}

func TestNode_Write_struct(t *testing.T) {
	c := NewCodec()
	n, err := c.Analyse(TestStruct{})
	if err != nil {
		t.Fatal(err)
	}
	wc := NewWriteContext()
	v := reflect.ValueOf(testWriteStructData)
	val := n.NewValFrom(v)
	if err := n.Write(wc, val); err != nil {
		t.Fatal(err)
	}
	targets := wc.targets
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
		if actual.Value == nil && expected.Value == nil {
			continue
		}
		var actualType, expectedType reflect.Type
		if actual.Value != nil {
			actualType = reflect.ValueOf(actual.Value).Type()
			if expected.Value == nil {
				t.Errorf("at %q got: %v; want nil", fileName, actual.Value)
			}
		}
		if expected.Value != nil {
			expectedType = reflect.ValueOf(expected.Value).Type()
			if actual.Value == nil {
				t.Errorf("at %q got: nil; want:\n%v", fileName, expected.Value)
			}
		}

		if actualType != expectedType {
			t.Errorf("got type %s; want %s at %q", actualType, expectedType, fileName)
			t.Errorf("values: got:\n%# v\nwant:\n%# v", actual.Value, expected.Value)
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
	return fmt.Sprintf("file: %q\n%s\n", ft.FilePath, ft.TestDataDump())
}

func (ft FileTarget) TestDataDump() string {
	data, err := json.MarshalIndent(ft.Value, "  ", "  ")
	if err != nil {
		panic(err)
	}
	return string(data)
}
