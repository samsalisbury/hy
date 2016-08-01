package hy

import "path"

// NodeContext provides data about a node passed down by its parent.
type NodeContext struct {
	// Parent is the parent node context.
	Parent *NodeContext
	// Tag is the tag this context is based on.
	Tag Tag
	// FieldName is the name of the field holding this node.
	FieldName string
	// FixedPathName is the fixed path name for this node (only nonempty
	// for struct fields.
	FixedPathName,
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

// Push creates a derivative node context.
func (c NodeContext) Push(tag Tag, fieldName string) NodeContext {
	return NodeContext{
		Parent:    &c,
		Tag:       tag,
		FieldName: fieldName,
	}
}

// Path returns the path of this context.
func (c NodeContext) Path() string {
	if c.Parent == nil {
		return c.FieldName
	}
	return path.Join(c.Parent.Path(), c.FieldName)
}
