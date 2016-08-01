package hy

import (
	"reflect"

	"github.com/pkg/errors"
)

// NodeBase is a node in an analysis.
type NodeBase struct {
	Context NodeContext
	// OwnType is the type of the underlying value of this Node. If the node
	// represents a non-pointer type, ValueType will be the same as the type it
	// represents. Otherwise, it will be the element type of that pointer type.
	OwnType reflect.Type
	// IsPtr is true when the real type this node represents is a pointer type.
	// It is false when this node represents a non-pointer type.
	IsPtr bool
	// Parent is the parent of this node. It is nil only for the root node.
	Parent Node
	// FieldName is the name of the struct field which contains this node. It is
	// empty for the child nodes of maps and slices.
	FieldName string // empty for children of maps and slices
	// PathNameFunc is a function that returns the file path segment this Node
	// should be written to, based on its key and value. For struct fields, the
	// key and value do not matter. For slice fields, the keyVal will be sey to
	// the value of the int index, for map fields keyVal will be the value of
	// the map key. For both slice and map fields, val will be the value of the
	// element.
	PathNameFunc func(keyVal, val reflect.Value) string
}

// SetContext sets the context of this node.
func (n *NodeBase) SetContext(c NodeContext) {
	n.Context = c
}

// Type returns the type this node represents.
func (n *NodeBase) Type() reflect.Type {
	return n.OwnType
}

// StructNode represents a struct to be stored in a file.
type StructNode struct {
	FileNode
	// Fields is a map of simple struct field names to their types.
	Fields map[string]reflect.Type
	// Children is a map of
	Children map[string]Node
}

// A FileNode represents a node to be stored in a file.
type FileNode struct {
	NodeBase
}

// A DirNodeBase is the base type for a node stored in a directory.
type DirNodeBase struct {
	NodeBase
	ElemNode Node
}

// A MapDirNode represents a map node to be stored in a directory.
type MapDirNode struct {
	DirNodeBase
	KeyType reflect.Type
}

// A SliceDirNode represents a slice to be stored in a directory.
type SliceDirNode struct {
	DirNodeBase
}

// Clone returns a clone.
func (n StructNode) Clone() Node { return &n }

// Clone returns a clone.
func (n FileNode) Clone() Node { return &n }

// Clone returns a clone.
func (n MapDirNode) Clone() Node { return &n }

// Clone returns a clone.
func (n SliceDirNode) Clone() Node { return &n }

// Node represents a generic node in the structure.
type Node interface {
	SetContext(NodeContext)
	Type() reflect.Type
	Clone() Node
}

type analysisInContext struct {
	Type reflect.Type
	Path string
}

type nodeInContext struct {
	Node    *Node
	Context NodeContext
}

// analyses is global state and should probably be removed
// using Node pointer so we can all refer to the same Node
var analyses = map[analysisInContext]nodeInContext{}

func getAnalysisForType(t reflect.Type) (*Node, bool) {
	for aic, nic := range analyses {
		if aic.Type == t {
			return nic.Node, true
		}
	}
	return nil, false
}
func registerAnalysis(t reflect.Type, n *Node, c NodeContext) {
	analyses[analysisInContext{t, c.Path()}] = nodeInContext{n, c}
}

var paths = map[string]Node{}

// Analyse analyses a tree starting at root.
func Analyse(root interface{}) (Node, error) {
	if root == nil {
		return nil, errors.New("cannot analyse nil")
	}
	return analyse(NodeContext{}, reflect.TypeOf(root))
}

func analyse(c NodeContext, t reflect.Type) (Node, error) {
	if analysis, ok := getAnalysisForType(t); ok {
		registerAnalysis(t, analysis, c)
		return *analysis, nil
	}
	// allocate the node now, child nodes can refer to it
	n := new(Node)
	registerAnalysis(t, n, c)
	var isPtr bool
	k := t.Kind()
	if k == reflect.Ptr {
		isPtr = true
		t = t.Elem()
		k = t.Kind()
	}
	var err error
	base := NodeBase{
		Context: c,
		OwnType: t,
	}
	switch k {
	default:
		return nil, errors.Errorf("cannot analyse %s fields", t)
	case reflect.Struct:
		*n, err = analyseStruct(base, t, isPtr)
	case reflect.Map:
		*n, err = analyseMap(base, t, isPtr)
	case reflect.Slice:
		*n, err = analyseSlice(base, t, isPtr)
	}
	return *n, errors.Wrapf(err, "analysing %s at %q", t, c.Path())
}

func analyseStruct(base NodeBase, t reflect.Type, isPtr bool) (Node, error) {
	fields := map[string]reflect.Type{}
	children := map[string]Node{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagStr := field.Tag.Get("hy")
		tag, err := parseTag(tagStr)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid tag %q", tagStr)
		}
		if tag.None {
			fields[field.Name] = field.Type
			continue
		}
		if tag.Ignore {
			continue
		}
		childContext := base.Context.Push(tag, field.Name)
		child, err := analyse(childContext, field.Type)
		if err != nil {
			return nil, errors.Wrapf(err, "analysing file field %T.%s", t, field.Name)
		}
		children[field.Name] = child
	}
	return &StructNode{
		FileNode: FileNode{
			NodeBase: base,
		},
		Fields:   fields,
		Children: children,
	}, nil
}

func analyseMap(base NodeBase, t reflect.Type, isPtr bool) (Node, error) {
	return nil, nil
}

func analyseSlice(base NodeBase, t reflect.Type, isPtr bool) (Node, error) {
	return nil, nil
}
