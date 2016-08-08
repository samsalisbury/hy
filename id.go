package hy

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

// NodeID identifies a node in the tree.
type NodeID struct {
	// ParentType is the type of this node's parent.
	ParentType,
	// Type is the type of this node.
	Type reflect.Type
	// IsPtr indicates if OwnType is a pointer really.
	IsPtr bool
	// IsLeaf indicates if this node can only be a leaf.
	IsLeaf bool
	// FieldName is the name of the parent field containing this node. FieldName
	// will be empty unless ParentType is a struct.
	FieldName string
}

// NewNodeID creates a new node ID.
func NewNodeID(parentType, typ reflect.Type, fieldName string) (NodeID, error) {
	t := typ
	var isPtr bool
	k := t.Kind()
	if k == reflect.Ptr {
		isPtr = true
		t = t.Elem()
		k = t.Kind()
		if k == reflect.Ptr {
			return NodeID{}, errors.New("cannot analyse pointer to pointer")
		}
	}
	if k == reflect.Interface {
		// TODO: We should allow interfaces that implement map.Get/Set and elem.GetID/SetID.
		return NodeID{}, errors.New("cannot analyse kind interface")
	}
	isLeaf := (k != reflect.Struct && k != reflect.Map && k != reflect.Slice)
	return NodeID{
		ParentType: parentType,
		Type:       t,
		IsPtr:      isPtr,
		IsLeaf:     isLeaf,
		FieldName:  fieldName,
	}, nil
}

func (id NodeID) String() string {
	ptr := ""
	if id.IsPtr {
		ptr = "*"
	}
	parent := "nil"
	if id.ParentType != nil {
		parent = id.ParentType.String()
	}
	return fmt.Sprintf("%s%s.%s(%s)", ptr, parent, id.FieldName, id.Type)
}
