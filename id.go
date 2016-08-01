package hy

import "reflect"

// NodeID identifies a node in the tree.
type NodeID struct {
	// ParentType is the type of this node's parent.
	ParentType,
	// Type is the type of this node.
	Type reflect.Type
	// IsPtr indicates if OwnType is a pointer really.
	IsPtr bool
	// FieldName is the name of the parent field containing this node. FieldName
	// will be empty unless ParentType is a struct.
	FieldName string
}
