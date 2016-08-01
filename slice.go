package hy

import "reflect"

// A SliceNode represents a slice to be stored in a directory.
type SliceNode struct {
	DirNodeBase
}

func (n *SliceNode) Write(c *NodeContext, v reflect.Value) error {
	return nil
}

func (c *Codec) analyseSlice(base NodeBase) (Node, error) {
	return nil, nil
}
