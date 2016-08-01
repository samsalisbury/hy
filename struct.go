package hy

import (
	"reflect"

	"github.com/pkg/errors"
)

// Struct represents a struct to be stored in a file.
type Struct struct {
	File
	// Fields is a map of simple struct field names to their types.
	Fields map[string]reflect.Type
	// Children is a map of field named to node pointers.
	Children map[string]*Node
}

func (n *Struct) Write(c NodeContext, v reflect.Value) error {
	return nil
}

func (c *Codec) analyseStruct(base NodeBase) (Node, error) {
	// Children need a pointer to this node, so create it first.
	n := &Struct{
		File: File{
			NodeBase: base,
		},
		Fields:   map[string]reflect.Type{},
		Children: map[string]*Node{},
	}
	for i := 0; i < n.Type.NumField(); i++ {
		field := n.Type.Field(i)
		tagStr := field.Tag.Get("hy")
		tag, err := parseTag(tagStr)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid tag %q", tagStr)
		}
		if tag.None {
			n.Fields[field.Name] = field.Type
			continue
		}
		if tag.Ignore {
			continue
		}
		child, err := c.analyse(n, field.Type, field.Name)
		if err != nil {
			return nil, errors.Wrapf(err, "analysing %T.%s", n.Type, field.Name)
		}
		n.Children[field.Name] = child
	}
	return n, nil
}
