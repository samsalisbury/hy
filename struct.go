package hy

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

// StructNode represents a struct to be stored in a file.
type StructNode struct {
	FileNode
	// Fields is a map of simple struct field names to their types.
	Fields map[string]reflect.Type
	// Children is a map of field name to node pointer.
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

// ReadTargets reads files into values.
func (n *StructNode) ReadTargets(c ReadContext, key reflect.Value) (reflect.Value, error) {
	val := reflect.New(n.Type)
	valInterface := val.Interface()
	fieldData, err := c.ReadFile(fmt.Sprint(key))
	if err != nil {
		return val, errors.Wrapf(err, "readding own fields")
	}
	if len(fieldData) != 0 {
		if err := c.UnmarshalFunc(fieldData, valInterface); err != nil {
			return val, errors.Wrapf(err, "unmarshaling %q", c.FilePath())
		}
	}
	val = val.Elem()
	for fieldName, child := range n.Children {
		childPathName, _ := (*child).FixedPathName()
		childContext := c.Push(childPathName)
		// TODO: Not this...
		d, err := c.ReadFile(childContext)
		childKey := reflect.ValueOf(childPathName)
		childVal, err := (*child).Read(childContext, childKey)
		if err != nil {
			return val, errors.Wrapf(err, "reading child %q", childPathName)
		}
		val.FieldByName(fieldName).Set(childVal)
	}
	return val, nil
}

// WriteTargets generates file targets.
func (n *StructNode) WriteTargets(c WriteContext, key, val reflect.Value) error {
	if err := n.WriteSelfTarget(c, key, val); err != nil {
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

// WriteSelfTarget writes the struct fields that are not stored in other files.
func (n *StructNode) WriteSelfTarget(c WriteContext, key, val reflect.Value) error {
	data := n.prepareFileData(val)
	t := &FileTarget{FilePath: c.Path(), Value: data}
	return errors.Wrap(c.Targets.Add(t), "failed to write self")
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
