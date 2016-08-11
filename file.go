package hy

import (
	"reflect"

	"github.com/pkg/errors"
)

// A FileNode represents a node to be stored in a file.
type FileNode struct {
	NodeBase
}

// NewFileNode creates a new file node.
func NewFileNode(base NodeBase) Node {
	return &FileNode{NodeBase: base}
}

// ChildPathName returns an empty string (file targets don't have children).
func (n *FileNode) ChildPathName(child Node, key, val reflect.Value) string {
	return ""
}

// ReadTargets reads a single file into a value.
func (n *FileNode) ReadTargets(c ReadContext, key reflect.Value) (reflect.Value, error) {
	val := reflect.New(n.Type)
	valInterface := val.Interface()
	b, err := c.ReadFile(c.PathName)
	if err != nil {
		return reflect.Value{}, errors.Wrapf(err, "reading file node")
	}
	if len(b) == 0 {
		return val, nil
	}
	if err := c.UnmarshalFunc(b, valInterface); err != nil {
		return reflect.Value{}, errors.Wrapf(err, "unmarshalling file %q", c.FilePath())
	}
	return val.Elem(), nil
}

// WriteTargets returns the write target for this file.
func (n *FileNode) WriteTargets(c WriteContext, key, val reflect.Value) error {
	t := &FileTarget{FilePath: c.Path(), Value: val.Interface()}
	return errors.Wrap(c.Targets.Add(t), "writing file target")
}
