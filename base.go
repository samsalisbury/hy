package hy

import (
	"reflect"

	"github.com/pkg/errors"
)

// NodeBase is a node in an analysis.
type NodeBase struct {
	NodeID
	// Parent is the parent of this node. It is nil only for the root node.
	Parent Node
	// FieldInfo is the field info for this node.
	Field *FieldInfo
	// Zero is a zero value of this node's Type.
	Zero interface{}
	// HasKey indicates if this type has a key (e.g. maps and slices)
	HasKey bool
	// self is a pointer to the node based on this node base. This means more
	// common functionality can be handled by NodeBase, by allowing it to call
	// methods on it's differentiated self.
	//
	// self is only safe to use after analysis is complete.
	self *Node
}

// ID returns the ID of this node base.
func (base NodeBase) ID() NodeID {
	return base.NodeID
}

// NewNodeBase returns a new NodeBase.
func NewNodeBase(id NodeID, parent Node, field *FieldInfo, self *Node) NodeBase {
	var k reflect.Kind
	if parent != nil {
		k = parent.ID().Type.Kind()
	}
	var zero interface{}
	if !id.IsPtr {
		zero = reflect.Zero(id.Type).Interface()
	}
	return NodeBase{
		NodeID: id,
		Parent: parent,
		Field:  field,
		Zero:   zero,
		HasKey: k == reflect.Map || k == reflect.Slice,
		self:   self,
	}
}

func (base NodeBase) Read(c ReadContext, key reflect.Value) (reflect.Value, error) {
	v, err := (*base.self).ReadTargets(c, key)
	if err != nil {
		return v, errors.Wrapf(err, "reading node")
	}
	if base.IsPtr {
		v = v.Addr()
	}
	return v, nil
}

func (base NodeBase) Write(c WriteContext, key, val reflect.Value) error {
	if base.IsPtr {
		val = val.Elem()
	}
	if !base.HasKey &&
		(!val.IsValid() || reflect.DeepEqual(val.Interface(), base.Zero)) {
		return nil
	}
	return (*base.self).WriteTargets(c, key, val)
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
	return base.Field.Tag
}

// FixedPathName returns the fixed path name of this node.
// If there is no fixed path name, returns empty string and false.
// Otherwise returns the fixed path name and true.
func (base NodeBase) FixedPathName() (string, bool) {
	if base.Field == nil {
		return "", false
	}
	if base.Field.PathName != "" {
		return base.Field.PathName, true
	}
	if base.FieldName != "" {
		return base.FieldName, true
	}
	return "", false
}
