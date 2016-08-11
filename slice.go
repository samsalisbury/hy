package hy

import (
	"fmt"
	"reflect"
	"strconv"

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

func (n *SliceNode) ReadTargets(c ReadContext, key reflect.Value) (reflect.Value, error) {
	files, err := c.ListFiles()
	if err != nil {
		return reflect.Value{}, errors.Wrapf(err, "listing files")
	}
	val := reflect.MakeSlice(n.Type, len(files), len(files))
	for _, keyStr := range files {
		childContext := c.Push(keyStr)
		childIndex, err := strconv.Atoi(keyStr)
		if err != nil {
			return val, errors.Wrapf(err, "parsing slice file %d", childIndex)
		}
		childKey := reflect.ValueOf(childIndex)
		childVal, err := (*n.ElemNode).Read(childContext, childKey)
		if err != nil {
			return reflect.Value{}, errors.Wrapf(err, "reading child %q", keyStr)
		}
		val.Index(childIndex).Set(childVal)
	}
	return val, nil
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
