package hy

import "reflect"

// A SliceDirNode represents a slice to be stored in a directory.
type SliceDirNode struct {
	DirNodeBase
}

func (n *SliceDirNode) Write(c *NodeContext, v reflect.Value) error {
	return nil
}

func analyseSlice(base NodeBase, t reflect.Type, isPtr bool) (Node, error) {
	return nil, nil
}
