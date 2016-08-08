package hy

import (
	"reflect"

	"github.com/pkg/errors"
)

// NodeSet is a set of Node pointers indexed by ID.
type NodeSet struct {
	nodes map[NodeID]*Node
}

// NewNodeSet creates a new node set.
func NewNodeSet() NodeSet {
	return NodeSet{nodes: map[NodeID]*Node{}}
}

// Register tries to register a node ID. If the ID is not yet registered, it
// returns a new node pointer and true. Otherwise it returns the already
// registered node pointer and false.
func (ns NodeSet) Register(id NodeID) (*Node, bool) {
	n, ok := ns.nodes[id]
	if ok {
		return n, false
	}
	n = new(Node)
	ns.nodes[id] = n
	return n, true
}

// Codec provides the primary encoding and decoding facility of this package.
type Codec struct {
	Nodes NodeSet
}

// NewCodec creates a new codec.
func NewCodec() *Codec {
	return &Codec{Nodes: NewNodeSet()}
}

// Analyse analyses a tree starting at root.
func (c *Codec) Analyse(root interface{}) (Node, error) {
	if root == nil {
		return nil, errors.New("cannot analyse nil")
	}
	t := reflect.TypeOf(root)
	id, err := NewNodeID(nil, t, "")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to analyse %T", root)
	}
	if id.IsLeaf {
		return nil, errors.Errorf("failed to analyse %s: cannot analyse kind %s",
			id.Type, id.Type.Kind())
	}
	n, err := c.NewNode(nil, id, Tag{})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to analyse %T", root)
	}
	return *n, err
}

// NewNode creates a new node.
func (c *Codec) NewNode(parent Node, id NodeID, tag Tag) (*Node, error) {
	n, new := c.Nodes.Register(id)
	if !new {
		return n, nil
	}
	var err error
	k := id.Type.Kind()
	base := NewNodeBase(id, parent, tag, n)
	if k == reflect.Struct {
		*n, err = c.NewStructNode(base)
		return n, err
	}
	if id.IsLeaf || !tag.IsDir {
		*n = NewFileNode(base)
		return n, nil
	}
	switch k {
	default:
		*n = NewFileNode(base)
	case reflect.Map:
		*n, err = c.NewMapNode(base)
	case reflect.Slice:
		*n, err = c.NewSliceNode(base)
	}
	return n, errors.Wrapf(err, "analysing %s failed", id)
}
