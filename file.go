package hy

import "reflect"

// A FileNode represents a node to be stored in a file.
type FileNode struct {
	NodeBase
}

// NewFileNode creates a new file node.
func NewFileNode(parentType reflect.Type, field reflect.StructField) *Node {
	t := field.Type
	var n Node
	n = &FileNode{
		NodeBase{
			NodeID: NodeID{
				ParentType: parentType,
				Type:       t,
				IsPtr:      t.Kind() == reflect.Ptr,
				FieldName:  field.Name,
			},
		},
	}
	return &n
}

func (n *FileNode) Write(c NodeContext, v reflect.Value) error {
	return nil
}
