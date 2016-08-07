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

// WriteTargets returns the write target for this file.
func (n *FileNode) WriteTargets(c WriteContext, key, val reflect.Value) (FileTargets, error) {
	fts, err := NewFileTargets(&FileTarget{
		Path: c.Path(), // path.Join(c.Path(), n.PathName(key, val)),
		Data: val.Interface(),
	})
	return fts, errors.Wrapf(err, "failed making write targets")
}
