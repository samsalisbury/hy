package hy

import "reflect"

// A FileNode represents a node to be stored in a file.
type FileNode struct {
	NodeBase
}

// NewFileNode creates a new file node.
func NewFileNode(parentType, t reflect.Type, field FieldInfo) *Node {
	var n Node
	n = &FileNode{
		NodeBase{
			NodeID: NodeID{
				ParentType: parentType,
				Type:       t,
				IsPtr:      t.Kind() == reflect.Ptr,
				FieldName:  field.Name,
			},
			Tag: field.Tag,
		},
	}
	return &n
}

func (n *FileNode) Write(c NodeContext, v reflect.Value) error {
	return nil
}
