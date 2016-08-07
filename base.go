package hy

import "reflect"

// NodeBase is a node in an analysis.
type NodeBase struct {
	NodeID
	// Parent is the parent of this node. It is nil only for the root node.
	Parent Node
	// Tag is the struct tag applying to this node.
	Tag Tag
	// self is a pointer to the node based on this node base. This means more
	// common functionality can be handled by NodeBase, by allowing it to call
	// methods on it's differentiated self.
	self *Node
}

// ID returns the ID of this node base.
func (base NodeBase) ID() NodeID {
	return base.NodeID
}

// PathName returns the path name segment of this node by querying its tag,
// field name, and parent's ChildPathName func.
func (base NodeBase) PathName(key, val reflect.Value) string {
	if fixedName, ok := base.FixedPathName(); ok {
		return fixedName
	}
	if base.Parent == nil {
		return ""
	}
	return base.Parent.ChildPathName(*base.self, key, val)
}

// GetTag returns the tag associated with this node.
func (base NodeBase) GetTag() Tag {
	return base.Tag
}

// FixedPathName returns the fixed path name of this node.
// If there is no fixed path name, returns empty string and false.
// Otherwise returns the fixed path name and true.
func (base NodeBase) FixedPathName() (string, bool) {
	if base.Tag.PathName != "" {
		return base.Tag.PathName, true
	}
	if base.FieldName != "" {
		return base.FieldName, true
	}
	return "", false
}
