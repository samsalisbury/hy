package hy

import "testing"

type TestWriteStruct struct {
	Name        string
	Int         int
	InlineSlice []string
	InlineMap   map[string]int
	Slice       []StructB          `hy:"./"`
	Map         map[string]StructB `hy:"./,Name"`
}

var testWriteStructData = TestWriteStruct{
	Name:        "Test struct writing",
	Int:         1,
	InlineSlice: []string{"a", "string", "slice"},
	InlineMap:   map[string]int{"one": 1, "two": 2, "three": 3},
	Slice:       []StructB{StructB{Name: "One"}, StructB{Name: "Two"}},
	Map: map[string]StructB{
		"First":  StructB{},
		"Second": StructB{},
	},
}

var testWriteFS = `
file: TestWriteStruct.yaml
Name: Test struct writing
Int: 1
InlineSlice:
	- a
	- string
	- slice
InlineMap:
	one: 1
	two: 2
	three: 3

file: Slice/1.yaml
Name: One

file: Slice/2.yaml
Name: Two

file: Map/First.yaml
Name: First

file: Map/Second.yaml
Name: Second
`

func TestWrite_struct(t *testing.T) {

}
