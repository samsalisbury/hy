package hy

import (
	"reflect"

	"github.com/pkg/errors"
)

// Codec provides the primary encoding and decoding facility of this package.
type Codec struct {
	nodes      NodeSet
	Writer     FileWriter
	TreeReader *FileTreeReader
}

// NewCodec creates a new codec.
func NewCodec(configure ...func(*Codec)) *Codec {
	c := &Codec{nodes: NewNodeSet()}
	for _, cfg := range configure {
		cfg(c)
	}
	if c.Writer == nil {
		c.Writer = JSONWriter
	}
	return c
}

func (c *Codec) Write(prefix string, root interface{}) error {
	rootNode, err := c.Analyse(root)
	if err != nil {
		return errors.Wrapf(err, "analysing structure")
	}
	wc := NewWriteContext()
	v := reflect.ValueOf(root)
	if err := rootNode.Write(wc, reflect.Value{}, v); err != nil {
		return errors.Wrapf(err, "generating write targets")
	}
	for _, t := range wc.targets.Snapshot() {
		if err := c.Writer.WriteFile(prefix, t); err != nil {
			return errors.Wrapf(err, "writing target %q", t.Path())
		}
	}
	return nil
}

func (c *Codec) Read(prefix string, root interface{}) error {
	rootNode, err := c.Analyse(root)
	if err != nil {
		return errors.Wrapf(err, "analysing structure")
	}
	targets, err := c.TreeReader.ReadTree(prefix)
	if err != nil {
		return errors.Wrapf(err, "reading tree at %q", prefix)
	}
	rc := NewReadContext(targets)
	val, err := rootNode.Read(rc, reflect.Value{})
	if err != nil {
		return errors.Wrapf(err, "reading root")
	}
	reflect.ValueOf(root).Elem().Set(val.Elem())
	return nil
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
	n, err := c.NewNode(nil, id, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to analyse %T", root)
	}
	return *n, err
}

// NewNode creates a new node.
func (c *Codec) NewNode(parent Node, id NodeID, field *FieldInfo) (*Node, error) {
	n, new := c.nodes.Register(id)
	if !new {
		return n, nil
	}
	var err error
	k := id.Type.Kind()
	base := NewNodeBase(id, parent, field, n)
	if k == reflect.Struct {
		*n, err = c.NewStructNode(base)
		return n, err
	}
	if id.IsLeaf || !field.Tag.IsDir {
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
