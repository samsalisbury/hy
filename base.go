package hy

import "reflect"

// NodeBase is a node in an analysis.
type NodeBase struct {
	// OwnType is the type of the underlying value of this Node. If the node
	// represents a non-pointer type, ValueType will be the same as the type it
	// represents. Otherwise, it will be the element type of that pointer type.
	OwnType reflect.Type
	// IsPtr is true when the real type this node represents is a pointer type.
	// It is false when this node represents a non-pointer type.
	IsPtr bool
	// Parent is the parent of this node. It is nil only for the root node.
	Parent Node
}

// Type returns the type this node represents.
func (n *NodeBase) Type() reflect.Type {
	return n.OwnType
}
