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
	return n, errors.Wrap(n.AnalyseElemNode(c), "analysing slice element node")
}

// ChildPathName returns the slice index as a string.
func (n *SliceNode) ChildPathName(child Node, key, val reflect.Value) string {
	return fmt.Sprint(key)
}

// WriteTargets writes all the elements of the slice.
func (n *SliceNode) WriteTargets(c WriteContext, key, val reflect.Value) (FileTargets, error) {
	fts := MakeFileTargets(val.Len())
	elemNode := *n.ElemNode
	for i := 0; i < val.Len(); i++ {
		v := val.Index(i)
		k := reflect.ValueOf(i)
		childContext := c.Push(elemNode.PathName(k, v))
		childTargets, err := elemNode.WriteTargets(childContext, k, v)
		if err != nil {
			return fts, errors.Wrapf(err, "writing slice index %d failed", i)
		}
		if err := fts.AddAll(childTargets); err != nil {
			return fts, errors.Wrapf(err, "adding children of index %d failed", i)
		}
	}
	return fts, nil
}
