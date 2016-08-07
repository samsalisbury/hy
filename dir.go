package hy

import "github.com/pkg/errors"

// A DirNodeBase is the base type for a node stored in a directory.
type DirNodeBase struct {
	NodeBase
	ElemNode *Node
}

// AnalyseElemNode sets the element node for this directory-bound node.
func (n *DirNodeBase) AnalyseElemNode(parent Node, c *Codec) error {
	elemType := n.Type.Elem()
	elemID, err := NewNodeID(n.Type, elemType, "")
	n.ElemNode, err = c.NewNode(parent, elemID, Tag{})
	return errors.Wrapf(err, "analysing type %T failed", elemType)
}
