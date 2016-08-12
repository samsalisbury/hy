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

// ReadTargets reads targets into slice indicies.
func (n *SliceNode) ReadTargets(c ReadContext, key reflect.Value) (reflect.Value, error) {
	val := reflect.New(n.Type).Elem() // TODO: Maybe use MakeSlice
	list := c.List()
	for _, indexStr := range list {
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			return val, errors.Wrapf(err, "converting %q to int", indexStr)
		}
		elemKey := reflect.ValueOf(index)
		elem := *n.ElemNode
		elemContext := c.Push(indexStr)
		elemVal, err := elem.Read(elemContext, elemKey)
		if err != nil {
			return val, errors.Wrapf(err, "reading index %d", index)
		}
		reflect.Append(val, elemVal)
		val.SetMapIndex(elemKey, elemVal)
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
