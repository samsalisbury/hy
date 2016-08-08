package hy

var expectedWriteFileTargets = map[string]FileTarget{
	"": FileTarget{
		Data: map[string]interface{}{
			"Name":        "Test struct writing",
			"Int":         1,
			"InlineSlice": []string{"a", "string", "slice"},
			"InlineMap":   map[string]int{"one": 1, "two": 2, "three": 3},
			"StructB":     StructB{},
			"StructBPtr":  nil,
		},
	},
	"a-file": FileTarget{
		Data: map[string]interface{}{
			"Name": "A file",
		},
	},
	"a-string-file": FileTarget{
		Data: "A string in a file.",
	},
	"nested": FileTarget{
		Data: map[string]interface{}{
			"Name":        "A nested struct pointer.",
			"Int":         2,
			"InlineMap":   nil,
			"InlineSlice": nil,
			"StructB":     StructB{},
			"StructBPtr":  nil,
		},
	},
	"nested/a-file": FileTarget{
		Data: map[string]interface{}{"Name": "Struct B file"},
	},
	"nested/slice/0": FileTarget{
		Data: map[string]interface{}{"Name": "Nested One"},
	},
	"nested/slice/1": FileTarget{
		Data: map[string]interface{}{"Name": "Nested Two"},
	},
	"nested/nested": FileTarget{
		Data: map[string]interface{}{
			"Name":        "",
			"Int":         0,
			"InlineSlice": nil,
			"InlineMap":   nil,
			"StructB":     StructB{},
			"StructBPtr":  nil,
		},
	},
	"nested/nested/a-slice-file": FileTarget{
		Data: []string{"this", "is", "a", "slice", "in", "a", "file"},
	},
	"nested/nested/a-map-file": FileTarget{
		Data: map[string]string{"deeply-nested": "map", "in a file": "yes"},
	},
	"nested/struct": FileTarget{
		Data: map[string]interface{}{
			"Name": "Struct B file",
		},
	},
	"nested/map-of-ptr/a-nil-file":       FileTarget{Data: nil},
	"nested/map-of-ptr/another-nil-file": FileTarget{Data: nil},
	"nested/map-of-ptr/this-one-has-a-value": FileTarget{
		Data: map[string]interface{}{
			"Name": "",
		},
	},
	"nested/map/a-zero-file": FileTarget{
		Data: map[string]interface{}{
			"Name": "",
		},
	},
	"nested/map/another-zero-file": FileTarget{
		Data: map[string]interface{}{
			"Name": "",
		},
	},
	"nested/map/nonzero-file": FileTarget{
		Data: map[string]interface{}{
			"Name": "I am not zero.",
		},
	},
	"slice/0": FileTarget{
		Data: map[string]interface{}{
			"Name": "One",
		},
	},
	"slice/1": FileTarget{
		Data: map[string]interface{}{
			"Name": "Two",
		},
	},
	"map/First": FileTarget{
		Data: map[string]interface{}{
			"Name": "First",
		},
	},
	"map/Second": FileTarget{
		Data: map[string]interface{}{
			"Name": "Second",
		},
	},
}
