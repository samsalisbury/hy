package hy

import (
	"path"
	"reflect"
)

// NodeContext provides data about a node passed down by its parent.
type NodeContext struct {
	// ParentType is the type containing this node.
	ParentType reflect.Type
	// Tag is the tag this context is based on. If ParentType is not a struct,
	// this will be a zero tag.
	Tag Tag
	// FieldName is the name of the field holding this node. If ParentType is
	// not a struct, FieldName will be empty.
	FieldName string
	// KeyField is the name of a field in a map element type that should match
	// the index in the map.
	KeyField,
	// GetKeyFuncName is the name of a function on a map element type that
	// returns the key of that element in its containing map.
	GetKeyFuncName,
	// SetKeyFuncName is the name of a function on a map element type that
	// is called when that element is added to its containing map.
	SetKeyFuncName string
}

// WriteContext is context collected during a write opration.
type WriteContext struct {
	// Parent is the parent write context.
	Parent *WriteContext
	// PathName is the name of this section of the path.
	PathName string
}

// Push creates a derivative node context.
func (c WriteContext) Push(pathName string) WriteContext {
	return WriteContext{
		Parent:   &c,
		PathName: pathName,
	}
}

// Path returns the path of this context.
func (c WriteContext) Path() string {
	if c.Parent == nil {
		return c.PathName
	}
	return path.Join(c.Parent.Path(), c.PathName)
}
