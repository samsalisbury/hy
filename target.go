package hy

type Target interface {
	Path() string
	Data() interface{}
}
