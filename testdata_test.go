package hy

var expectedFileTargets FileTargets
var expectedFileTargetsSnapshot map[string]*FileTarget

func init() {
	var err error
	expectedFiles, err := NewFileTargets([]*FileTarget{
		{Path: "",
			Data: map[string]interface{}{
				"Name":        "Test struct writing",
				"Int":         1,
				"InlineSlice": []string{"a", "string", "slice"},
				"InlineMap":   map[string]int{"one": 1, "two": 2, "three": 3},
				"StructB":     StructB{},
				"StructBPtr":  nil,
			},
		},
		{Path: "a-file",
			Data: map[string]interface{}{
				"Name": "A file",
			},
		},
		{Path: "a-string-file",
			Data: "A string in a file.",
		},
		{Path: "nested",
			Data: map[string]interface{}{
				"Name":        "A nested struct pointer.",
				"Int":         2,
				"InlineMap":   nil,
				"InlineSlice": nil,
				"StructB":     StructB{},
				"StructBPtr":  nil,
			},
		},
		{Path: "nested/a-file",
			Data: map[string]interface{}{"Name": "Struct B file"},
		},
		{Path: "nested/slice/0",
			Data: map[string]interface{}{"Name": "Nested One"},
		},
		{Path: "nested/slice/1",
			Data: map[string]interface{}{"Name": "Nested Two"},
		},
		{Path: "nested/nested",
			Data: map[string]interface{}{
				"Name":        "",
				"Int":         0,
				"InlineSlice": nil,
				"InlineMap":   nil,
				"StructB":     StructB{},
				"StructBPtr":  nil,
			},
		},
		{Path: "nested/nested/a-slice-file",
			Data: []string{"this", "is", "a", "slice", "in", "a", "file"},
		},
		{Path: "nested/nested/a-map-file",
			Data: map[string]string{"deeply-nested": "map", "in a file": "yes"},
		},
		{Path: "nested/map-of-ptr/a-nil-file",
			Data: nil},
		{Path: "nested/map-of-ptr/another-nil-file",
			Data: nil},
		{Path: "nested/map-of-ptr/this-one-has-a-value",
			Data: map[string]interface{}{
				// set automatically
				"Name": "this-one-has-a-value",
			},
		},
		{Path: "nested/map/a-zero-file",
			Data: map[string]interface{}{
				// set automatically
				"Name": "a-zero-file",
			},
		},
		{Path: "nested/map/another-zero-file",
			Data: map[string]interface{}{
				"Name": "another-zero-file",
			},
		},
		{Path: "slice/0",
			Data: map[string]interface{}{
				"Name": "One",
			},
		},
		{Path: "slice/1",
			Data: map[string]interface{}{
				"Name": "Two",
			},
		},
		{Path: "map/First",
			Data: map[string]interface{}{
				"Name": "First",
			},
		},
		{Path: "map/Second",
			Data: map[string]interface{}{
				"Name": "Second",
			},
		},
	}...)
	if err != nil {
		panic(err)
	}
	expectedFileTargetsSnapshot = expectedFiles.Snapshot()
}
