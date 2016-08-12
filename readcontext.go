package hy

import (
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
)

// ReadContext is context collected during a read opration.
type ReadContext struct {
	// Targets is the collected targets for this read context.
	targets FileTargets
	// Reader reads data from path.
	Reader FileReader
	// Parent is the parent read context.
	Parent *ReadContext
	// PathName is the name of this section of the path.
	PathName string
	// Prefix is the path prefix.
	Prefix string
}

// NewReadContext returns a new read context.
func NewReadContext(targets FileTargets) ReadContext {
	return ReadContext{targets: targets}
}

// Push creates a derivative node context.
func (c ReadContext) Push(pathName string) ReadContext {
	return ReadContext{
		targets:  c.targets,
		Reader:   c.Reader,
		Parent:   &c,
		PathName: pathName,
		Prefix:   c.Prefix,
	}
}

// List lists files in the current directory.
func (c ReadContext) List() []string {
	l := []string{}
	for path := range c.targets.Snapshot() {
		if !strings.HasPrefix(path, c.Path()) {
			continue
		}
		p := strings.TrimPrefix(path, c.Path())
		if p == "" || strings.ContainsRune(p, os.PathSeparator) {
			continue
		}
		if p == "_" {
			p = ""
		}
		l = append(l, p)
	}
	return l
}

func (c ReadContext) Read(v interface{}) error {
	t, ok := c.targets.Snapshot()[c.Path()]
	if !ok {
		return nil
	}
	return errors.Wrapf(c.Reader.ReadFile(c.Prefix, t, v), "reading %q", c.Path())
}

// Exists checks that a file exists at the current path.
func (c ReadContext) Exists() bool {
	_, ok := c.targets.Snapshot()[c.Path()]
	return ok
}

// Path returns the path of this context.
func (c ReadContext) Path() string {
	if c.Parent == nil {
		return c.PathName
	}
	return path.Join(c.Parent.Path(), c.PathName)
}

// SetValue sets the value of the current path.
func (c ReadContext) SetValue(v interface{}) error {
	t := &FileTarget{FilePath: c.Path(), Value: v}
	return errors.Wrapf(c.targets.Add(t), "setting value at %q", c.Path())
}
