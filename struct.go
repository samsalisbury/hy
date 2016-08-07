package hy

import (
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

// NewStructNode makes a new struct node.
func (c *Codec) NewStructNode(base NodeBase) (Node, error) {
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
		childNodeID, _ := NewNodeID(n.Type, field.Type, field.Name)
		child, err := c.NewNode(n, childNodeID, fieldInfo.Tag)
		if err != nil {
			return nil, errors.Wrapf(err, "analysing %T.%s", n.Type, field.Name)
		}
		if child != nil {
			n.Children[field.Name] = child
		}
	}
	return n, nil
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
	if val.Type().Kind() == reflect.Ptr {
		if val.IsNil() {
			return MakeFileTargets(0), nil
		}
		val = val.Elem()
		if !val.IsValid() {
			return MakeFileTargets(0), nil
		}
	}
	fts := MakeFileTargets(len(n.Children) + 1)
	if err := fts.Add(&FileTarget{
		Path: c.Path(),
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
