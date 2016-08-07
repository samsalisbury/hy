package hy

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

// A SliceNode represents a slice to be stored in a directory.
type SliceNode struct {
	*DirNodeBase
}

// NewSliceNode makes a new slice node.
func (c *Codec) NewSliceNode(base NodeBase) (Node, error) {
	n := &SliceNode{&DirNodeBase{NodeBase: base}}
	return n, errors.Wrap(n.AnalyseElemNode(n, c), "analysing slice element node")
}

// ChildPathName returns the slice index as a string.
func (n *SliceNode) ChildPathName(child Node, key, val reflect.Value) string {
	return fmt.Sprint(key)
}

// WriteTargets writes all the elements of the slice.
func (n *SliceNode) WriteTargets(c WriteContext, key, val reflect.Value) error {
	elemNode := *n.ElemNode
	for i := 0; i < val.Len(); i++ {
		v := val.Index(i)
		k := reflect.ValueOf(i)
		childContext := c.Push(elemNode.PathName(k, v))
		if err := elemNode.Write(childContext, k, v); err != nil {
			return errors.Wrapf(err, "writing slice index %d failed", i)
		}
	}
	return nil
}
