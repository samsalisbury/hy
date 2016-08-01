package hy

import "reflect"

// Node represents a generic node in the structure.
type Node interface {
	Type() reflect.Type
	Write(NodeContext, reflect.Value) error
	//Read(NodeContext)
}
