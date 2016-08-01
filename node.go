package hy

// Node represents a generic node in the structure.
type Node interface {
	ID() NodeID
	//Write(WriteContext, reflect.Value) error
	//Read(NodeContext)
}
