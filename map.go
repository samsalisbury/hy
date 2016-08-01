package hy

import (
	"reflect"

	"github.com/pkg/errors"
)

// A MapNode represents a map node to be stored in a directory.
type MapNode struct {
	DirNodeBase
	KeyType reflect.Type
}

func analyseMap(base NodeBase, t reflect.Type, isPtr bool) (Node, error) {
	child, err := analyse(t.Elem())
	if err != nil {
		return nil, errors.Wrapf(err, "analysing map element type %T failed", t.Elem())
	}
	return &MapNode{
		DirNodeBase: DirNodeBase{ElemNode: child},
		KeyType:     t.Key(),
	}, nil
}

func (n *MapNode) Write(c NodeContext, v reflect.Value) error {
	for _, k := range v.MapKeys() {
		indexString := "name"
		elemContext := c.Push(Tag{}, indexString)
		if err := n.ElemNode.Write(elemContext, v.MapIndex(k)); err != nil {
			return errors.Wrapf(err, "writing index %q", indexString)
		}
	}
	return nil
}
