package hy

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

// A MapNode represents a map node to be stored in a directory.
type MapNode struct {
	*DirNodeBase
	KeyType reflect.Type
}

// Detect returns nil if this base is a map.
func (MapNode) Detect(base NodeBase) error {
	if base.Kind == reflect.Map {
		return nil
	}
	return errors.Errorf("got kind %s; want map", base.Kind)
}

// New returns a new MapNode.
func (MapNode) New(base NodeBase, c *Codec) (Node, error) {
	n := &MapNode{
		DirNodeBase: &DirNodeBase{
			NodeBase: base,
		},
		KeyType: base.Type.Key(),
	}
	return n, errors.Wrap(n.AnalyseElemNode(n, c), "analysing map element node")
}

// ChildPathName returns the key as a string.
func (n *MapNode) ChildPathName(child Node, key, val reflect.Value) string {
	//if n.Field != nil && n.Field.KeyField != "" {
	//	n.Field.GetKeyFunc.Call([]reflect.Value{val})
	//}
	return fmt.Sprint(key)
}

// ReadTargets reads targets into map entries.
func (n *MapNode) ReadTargets(c ReadContext, val Val) error {
	list := c.List()
	for _, keyStr := range list {
		elemKey := reflect.ValueOf(keyStr)
		elem := *n.ElemNode
		elemContext := c.Push(keyStr)
		elemVal := elem.NewKeyedVal(elemKey)
		err := elem.Read(elemContext, elemVal)
		if err != nil {
			return errors.Wrapf(err, "reading child %s", keyStr)
		}
		val.SetMapElement(elemVal)
	}
	return nil
}

// WriteTargets writes all map elements.
func (n *MapNode) WriteTargets(c WriteContext, val Val) error {
	if !val.ShouldWrite() {
		return nil
	}
	elemNode := *n.ElemNode
	for _, elemVal := range val.MapElements(elemNode) {
		if n.Field != nil && n.Field.KeyField != "" {
			n.Field.SetKeyFunc.Call([]reflect.Value{elemVal.Ptr, elemVal.Key})
		}
		childContext := c.Push(elemNode.PathName(elemVal))
		if err := elemNode.Write(childContext, elemVal); err != nil {
			return errors.Wrapf(err, "writing map index %q failed",
				fmt.Sprint(elemVal.Key))
		}
	}
	return nil
}
