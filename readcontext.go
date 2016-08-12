package hy

import (
	"log"
	"os"
	"path"
	"path/filepath"
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
func NewReadContext(prefix string, targets FileTargets, reader FileReader) ReadContext {
	return ReadContext{Prefix: prefix, targets: targets, Reader: reader}
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
		p, err := filepath.Rel(path, c.Path())
		if err != nil {
			panic(err)
		}
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
	if !c.Exists() {
		return nil
	}
	filePath := filepath.Join(c.Prefix, c.Path())
	return errors.Wrapf(c.Reader.ReadFile(filePath, v), "reading %q", c.Path())
}

// Exists checks that a file exists at the current path.
func (c ReadContext) Exists() bool {
	log.Printf("CHECKING EXISTENCE %q\n", c.Path())
	_, ok := c.targets.Snapshot()[c.Path()]
	log.Println(c.targets.Paths())
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
