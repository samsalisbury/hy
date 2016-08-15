package hy

import "reflect"

// Val wraps a reflect.Value which is a pointer.
type Val struct {
	Base *NodeBase
	// Ptr is the underlying pointer value.
	Ptr reflect.Value
	// Key is the associated key for this value. May be invalid.
	Key reflect.Value
	// IsPtr indicates whether the final version of this value should be a
	// pointer.
	IsPtr bool
}

// Final returns the final reflect.Value.
func (v Val) Final() reflect.Value {
	if v.IsPtr {
		return v.Ptr
	}
	return v.Ptr.Elem()
}

func (v Val) IsZero() bool {
	return v.Ptr.IsNil() ||
		!v.Ptr.Elem().IsValid() ||
		reflect.DeepEqual(v.Ptr.Elem().Interface(), v.Base.Zero)
}

func (v Val) ShouldWrite() bool {
	return v.Key.IsValid() || !v.IsZero()
}

func (v Val) SetField(name string, val Val) {
	v.Ptr.Elem().FieldByName(name).Set(val.Final())
}

func (v Val) GetField(name string) reflect.Value {
	return v.Ptr.Elem().FieldByName(name)
}

func (v Val) SetMapElement(val Val) {
	v.Ptr.Elem().SetMapIndex(val.Key, val.Final())
}

func (v Val) MapElements(elemNode Node) []Val {
	m := v.Ptr.Elem()
	vals := make([]Val, m.Len())
	for i, key := range m.MapKeys() {
		vals[i] = elemNode.NewKeyedValFrom(key, m.MapIndex(key))
	}
	return vals
}

func (v Val) Append(val Val) {
	reflect.Append(v.Ptr.Elem(), val.Final())
}

func (v Val) SliceElements(elemNode Node) []Val {
	s := v.Ptr.Elem()
	vals := make([]Val, s.Len())
	for i := 0; i < s.Len(); i++ {
		vals[i] = elemNode.NewKeyedVal(reflect.ValueOf(i))
	}
	return vals
}
