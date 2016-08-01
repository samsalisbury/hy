package hy

// NodeBase is a node in an analysis.
type NodeBase struct {
	NodeID
	// Parent is the parent of this node. It is nil only for the root node.
	Parent Node
}

func (base NodeBase) ID() NodeID {
	return base.NodeID
}
