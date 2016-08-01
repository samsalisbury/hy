package hy

import "reflect"

// Node represents a generic node in the structure.
type Node interface {
	ID() NodeID
	Write(NodeContext, reflect.Value) error
	//Read(NodeContext)
}
