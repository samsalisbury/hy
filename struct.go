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
	// Children is a map of
	Children map[string]Node
}

func (n *Struct) Write(c NodeContext, v reflect.Value) error {
	return nil
}

func analyseStruct(base NodeBase, t reflect.Type, isPtr bool) (Node, error) {
	fields := map[string]reflect.Type{}
	children := map[string]Node{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagStr := field.Tag.Get("hy")
		tag, err := parseTag(tagStr)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid tag %q", tagStr)
		}
		if tag.None {
			fields[field.Name] = field.Type
			continue
		}
		if tag.Ignore {
			continue
		}
		//childContext := base.Context.Push(tag, field.Name)
		child, err := analyse(field.Type)
		if err != nil {
			return nil, errors.Wrapf(err, "analysing file field %T.%s", t, field.Name)
		}
		children[field.Name] = child
	}
	return &Struct{
		File: File{
			NodeBase: base,
		},
		Fields:   fields,
		Children: children,
	}, nil
}
