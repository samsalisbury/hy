package hy

import "reflect"

// Node represents a generic node in the structure.
type Node interface {
	// ID returns this node's ID.
	ID() NodeID
	// Tag returns the parsed tag of this node.
	GetTag() Tag
	// FixedPathName returns the indubitable path segment name of this node.
	FixedPathName() (string, bool)
	// ChildPathName returns the path segment for children of this node.
	// If the node's parent is a map or slice, both key and val will have
	// valid values, with val having the same type as this node.
	// If the node's parent is a map, then the key will be a value of the
	// parent's key type.
	// If the node's parent is a slice, then key will be an int value
	// representing the index of this element.
	// If the node's parent is a struct, then key will be an invalid value,
	// and val will be the value of that struct field.
	ChildPathName(child Node, key, val reflect.Value) string

	// PathName returns the path name of this node. Implemented in NodeBase.
	PathName(key, val reflect.Value) string
	// WriteTargets generates the write targets from this node.
	WriteTargets(c WriteContext, key, val reflect.Value) (FileTargets, error)
	//Read(NodeContext)
}
