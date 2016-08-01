package hy

import "reflect"

// A SliceDirNode represents a slice to be stored in a directory.
type SliceDirNode struct {
	DirNodeBase
}

func (n *SliceDirNode) Write(c *NodeContext, v reflect.Value) error {
	return nil
}

func (c *Codec) analyseSlice(base NodeBase) (Node, error) {
	return nil, nil
}
