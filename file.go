package hy

import "reflect"

// A FileNode represents a node to be stored in a file.
type File struct {
	NodeBase
}

func (n *File) Write(c NodeContext, v reflect.Value) error {
	return nil
}
