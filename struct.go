package hy

import (
	"path"
	"reflect"

	"github.com/pkg/errors"
)

// StructNode represents a struct to be stored in a file.
type StructNode struct {
	FileNode
	// Fields is a map of simple struct field names to their types.
	Fields map[string]reflect.Type
	// Children is a map of field named to node pointers.
	Children map[string]*Node
}

// ChildPathName returns the path segment for this node's children.
func (n *StructNode) ChildPathName(child Node, key, val reflect.Value) string {
	name, ok := child.FixedPathName()
	if !ok {
		panic("child of struct must have fixed name")
	}
	return name
}

// WriteTargets generates file targets.
func (n *StructNode) WriteTargets(c WriteContext, key, val reflect.Value) (FileTargets, error) {
	writePath := path.Join(c.Path(), n.PathName(key, val))
	fts := MakeFileTargets(len(n.Children) + 1)
	if err := fts.Add(&FileTarget{
		Path: writePath,
		Data: n.prepareFileData(val),
	},
	); err != nil {
		return fts, errors.Wrapf(err, "failed to write self")
	}
	for name, childPtr := range n.Children {
		childNode := *childPtr
		childKey := reflect.ValueOf(name)
		childVal := val.FieldByName(name)
		childContext := c.Push(childNode.PathName(childKey, childVal))
		childTargets, err := childNode.WriteTargets(childContext, childKey, childVal)
		if err != nil {
			return fts, errors.Wrapf(err, "failed to write child %s", name)
		}
		if err := fts.AddAll(childTargets); err != nil {
			return fts, errors.Wrapf(err, "failed to add targets from child %s", name)
		}
	}
	return fts, nil
}

func (n *StructNode) prepareFileData(val reflect.Value) map[string]interface{} {
	out := make(map[string]interface{}, len(n.Fields))
	for name := range n.Fields {
		out[name] = val.FieldByName(name).Interface()
	}
	return out
}

func (c *Codec) analyseStruct(base NodeBase) (Node, error) {
	// Children need a pointer to this node, so create it first.
	n := &StructNode{
		FileNode: FileNode{
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
		fieldInfo := FieldInfo{Tag: tag, Name: field.Name}
		if tag.IsDir || field.Type.Kind() == reflect.Struct {
			child, err := c.analyse(n, field.Type, fieldInfo)
			if err != nil {
				return nil, errors.Wrapf(err, "analysing %T.%s", n.Type, field.Name)
			}
			n.Children[field.Name] = child
			continue
		}
		n.Children[field.Name] = NewFileNode(n.Type, field.Type, fieldInfo)
	}
	return n, nil
}
