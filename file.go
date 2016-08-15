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

// ReadTargets reads a single file target.
func (n *FileNode) ReadTargets(c ReadContext, val Val) error {
	err := c.Read(val.Ptr.Interface())
	return errors.Wrapf(err, "reading file")
}

// WriteTargets returns the write target for this file.
func (n *FileNode) WriteTargets(c WriteContext, val Val) error {
	return errors.Wrap(c.SetValue(val), "writing file target")
}
