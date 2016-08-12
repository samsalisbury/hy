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
		field, err := NewFieldInfo(n.Type.Field(i)) //tag, Name: field.Name)
		if err != nil {
			return nil, errors.Wrapf(err, "reading field %s.%s", n.Type, n.Type.Field(i).Name)
		}
		if field.Tag.None {
			n.Fields[field.Name] = field.Type
			continue
		}
		if field.Tag.Ignore {
			continue
		}
		childNodeID, err := NewNodeID(n.Type, field.Type, field.Name)
		if err != nil {
			return nil, errors.Wrapf(err, "getting ID for %T.%s", n.Type, field.Name)
		}
		child, err := c.NewNode(n, childNodeID, field)
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
	name, _ := child.FixedPathName()
	return name
}

// ReadTargets reads targets into struct fields.
func (n *StructNode) ReadTargets(c ReadContext, key reflect.Value) (reflect.Value, error) {
	val := reflect.New(n.Type)
	if err := c.Read(val.Interface()); err != nil {
		return val, errors.Wrapf(err, "reading struct fields")
	}
	val = val.Elem()
	for fieldName, childPtr := range n.Children {
		child := *childPtr
		childPathName, _ := child.FixedPathName()
		childContext := c.Push(childPathName)
		if !childContext.Exists() {
			continue
		}
		childVal, err := child.Read(childContext, reflect.Value{})
		if err != nil {
			return val, errors.Wrapf(err, "reading child %s", fieldName)
		}
		val.FieldByName(fieldName).Set(childVal)
	}
	return val, nil
}

// WriteTargets generates file targets.
func (n *StructNode) WriteTargets(c WriteContext, key, val reflect.Value) error {
	if err := c.SetValue(n.prepareFileData(val)); err != nil {
		return errors.Wrap(err, "writing self")
	}
	if !val.IsValid() {
		return nil
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

func (n *StructNode) prepareFileData(val reflect.Value) interface{} {
	if !val.IsValid() {
		return nil
	}
	out := make(map[string]interface{}, len(n.Fields))
	for name := range n.Fields {
		out[name] = val.FieldByName(name).Interface()
	}
	return out
}
