package hy

import (
	"reflect"

	"github.com/pkg/errors"
)

// A SliceNode represents a slice to be stored in a directory.
type SliceNode struct {
	DirNodeBase
}

func (n *SliceNode) Write(c WriteContext, v reflect.Value) error {
	return nil
}

func (c *Codec) analyseSlice(base NodeBase) (Node, error) {
	n := &SliceNode{
		DirNodeBase{
			NodeBase: base,
		},
	}
	elemType := n.Type.Elem()
	elemNode, err := c.analyse(n, elemType, "")
	if err != nil {
		return nil, errors.Wrapf(err, "analysing %T failed", elemType)
	}
	n.ElemNode = elemNode
	return n, nil
}
