package hy

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
)

// ReaadContext is context collected during a read operation.
type ReadContext struct {
	// Targets is the collected targets in this write context.
	Targets FileTargets
	// Parent is the parent write context.
	Parent *ReadContext
	// PathName is the name of this section of the path.
	PathName string
	// FileExtension is the extension of files to be considered.
	FileExtension string
	// Prefix is the base directory.
	Prefix string
	// UnmarshalFunc is a function to unmarshal a file into a value.
	UnmarshalFunc func(interface{}, []byte) error
}

// NewReadContext returns a new read context.
func NewReadContext(prefix, ext string) ReadContext {
	return ReadContext{
		Targets:       MakeFileTargets(0),
		FileExtension: ext,
		Prefix:        prefix,
	}
}

// Push creates a derivative read context.
func (c ReadContext) Push(pathName string) ReadContext {
	return ReadContext{
		Targets:       c.Targets,
		Parent:        &c,
		PathName:      pathName,
		FileExtension: c.FileExtension,
	}
}

// Path returns the path of this context.
func (c ReadContext) Path() string {
	if c.Parent == nil {
		return path.Join(c.Prefix, c.PathName)
	}
	return path.Join(c.Parent.Path(), c.PathName)
}

// FilePath returns path with file extension.
func (c ReadContext) FilePath() string {
	return c.Path() + "." + c.FileExtension
}

func (c *ReadContext) ListFiles() ([]string, error) {
	fs, err := ioutil.ReadDir(c.Path())
	if err != nil {
		return nil, errors.Wrapf(err, "read dir failed")
	}
	files := make([]string, 0, len(fs))
	for _, f := range fs {
		if f.IsDir() || !strings.HasSuffix(f.Name(), "."+c.FileExtension) {
			continue
		}
		files = append(files, strings.TrimSuffix(f.Name(), "."+c.FileExtension))
	}
	return files, nil
}

func (c *ReadContext) ReadFile(name string) ([]byte, error) {
	b, err := ioutil.ReadFile(c.FilePath())
	if os.IsNotExist(err) {
		err = nil // not existing is fine, means zero value
	}
	return b, errors.Wrapf(err, "reading %q", c.FilePath())
}
