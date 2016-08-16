package hy

var expectedFileTargets FileTargets
var expectedFileTargetsSnapshot map[string]*FileTarget

func init() {
	var err error
	expectedFiles, err := NewFileTargets([]*FileTarget{
		{FilePath: "",
			Value: map[string]interface{}{
				"Name":        "Test struct writing",
				"Int":         1,
				"InlineSlice": []string{"a", "string", "slice"},
				"InlineMap":   map[string]int{"one": 1, "two": 2, "three": 3},
				"StructB":     StructB{},
				"StructBPtr":  nil,
			},
		},
		{FilePath: "a-file",
			Value: map[string]interface{}{
				"Name": "A file",
			},
		},
		{FilePath: "a-string-file",
			Value: "A string in a file.",
		},
		{FilePath: "nested",
			Value: map[string]interface{}{
				"Name":        "A nested struct pointer.",
				"Int":         2,
				"InlineMap":   nil,
				"InlineSlice": nil,
				"StructB":     StructB{},
				"StructBPtr":  nil,
			},
		},
		{FilePath: "nested/a-file",
			Value: map[string]interface{}{"Name": "Struct B file"},
		},
		{FilePath: "nested/slice/0",
			Value: map[string]interface{}{"Name": "Nested One"},
		},
		{FilePath: "nested/slice/1",
			Value: map[string]interface{}{"Name": "Nested Two"},
		},
		{FilePath: "nested/nested",
			Value: map[string]interface{}{
				"Name":        "",
				"Int":         0,
				"InlineSlice": nil,
				"InlineMap":   nil,
				"StructB":     StructB{},
				"StructBPtr":  nil,
			},
		},
		{FilePath: "nested/nested/a-slice-file",
			Value: []string{"this", "is", "a", "slice", "in", "a", "file"},
		},
		{FilePath: "nested/nested/a-map-file",
			Value: map[string]string{"deeply-nested": "map", "in a file": "yes"},
		},
		{FilePath: "nested/map-of-ptr/a-nil-file",
			Value: nil},
		{FilePath: "nested/map-of-ptr/another-nil-file",
			Value: nil},
		{FilePath: "nested/map-of-ptr/this-one-has-a-value",
			Value: map[string]interface{}{
				// set automatically
				"Name": "this-one-has-a-value",
			},
		},
		{FilePath: "nested/map/a-zero-file",
			Value: map[string]interface{}{
				// set automatically
				"Name": "a-zero-file",
			},
		},
		{FilePath: "nested/map/another-zero-file",
			Value: map[string]interface{}{
				"Name": "another-zero-file",
			},
		},
		{FilePath: "slice/0",
			Value: map[string]interface{}{
				"Name": "One",
			},
		},
		{FilePath: "slice/1",
			Value: map[string]interface{}{
				"Name": "Two",
			},
		},
		{FilePath: "map/First",
			Value: map[string]interface{}{
				"Name": "First",
			},
		},
		{FilePath: "map/Second",
			Value: map[string]interface{}{
				"Name": "Second",
			},
		},
		{FilePath: "textmarshaler/Test-1",
			Value: nil,
		},
		{FilePath: "textmarshaler/Another-13",
			Value: nil,
		},
	}...)
	if err != nil {
		panic(err)
	}
	expectedFileTargetsSnapshot = expectedFiles.Snapshot()
}
