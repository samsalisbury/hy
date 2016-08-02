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

func (c *Codec) analyseMap(base NodeBase) (Node, error) {
	n := &MapNode{
		DirNodeBase: DirNodeBase{
			NodeBase: base,
		},
	}
	elemType := n.Type.Elem()
	elemNode, err := c.analyse(n, elemType, FieldInfo{})
	if err != nil {
		return nil, errors.Wrapf(err, "analysing type %T failed", elemType)
	}
	n.KeyType = n.Type.Key()
	n.ElemNode = elemNode
	return n, nil
}

func (n *MapNode) Write(c NodeContext, v reflect.Value) error {
	//for _, k := range v.MapKeys() {
	//	indexString := "name"
	//	elemContext := c.Push(Tag{}, indexString)
	//	if err := n.ElemNode.Write(elemContext, v.MapIndex(k)); err != nil {
	//		return errors.Wrapf(err, "writing index %q", indexString)
	//	}
	//}
	//return nil
	return nil
}
