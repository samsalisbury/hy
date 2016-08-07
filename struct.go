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
func (n *StructNode) WriteTargets(c WriteContext, key, val reflect.Value) error {
	if err := n.WriteSelfTarget(c, key, val); err != nil {
		return errors.Wrap(err, "writing self")
	}
	for name, childPtr := range n.Children {
		childNode := *childPtr
		childKey := reflect.ValueOf(name)
		childVal := val.FieldByName(name)
		childContext := c.Push(childNode.PathName(childKey, childVal))
		if err := childNode.Write(childContext, childKey, childVal); err != nil {
			return errors.Wrapf(err, "failed to write child %s", name)
		}
	}
	return nil
}

// WriteSelfTarget writes the struct fields that are not stored in other files.
func (n *StructNode) WriteSelfTarget(c WriteContext, key, val reflect.Value) error {
	t := &FileTarget{Path: c.Path(), Data: n.prepareFileData(val)}
	if t == nil {
		panic("NO CONTEXT")
	}
	err := c.Targets.Add(t)

	return errors.Wrap(err, "failed to write self")
}

func (n *StructNode) prepareFileData(val reflect.Value) map[string]interface{} {
	out := make(map[string]interface{}, len(n.Fields))
	for name := range n.Fields {
		out[name] = val.FieldByName(name).Interface()
	}
	return out
}
