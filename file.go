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

func (n *FileNode) ReadTargets(c ReadContext, key reflect.Value) (reflect.Value, error) {
	val := reflect.New(n.Type)
	err := c.Read(val.Interface())
	return val.Elem(), errors.Wrapf(err, "reading file")
	//return val.Elem(), errors.Wrapf(c.Read(val.Interface()), "reading file")
}

// WriteTargets returns the write target for this file.
func (n *FileNode) WriteTargets(c WriteContext, key, val reflect.Value) error {
	return errors.Wrap(c.SetValue(val.Interface()), "writing file target")
}
