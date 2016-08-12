package hy

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// FileTreeReader gets targets from the filesystem.
type FileTreeReader struct {
	// FileExtension is the extension of files to consider.
	FileExtension string
	// Prefix is the path prefix.
	Prefix string
	// Targets is the collection of discovered targets.
	Targets FileTargets
}

// NewFileTreeReader returns a new FileTreeReader configured to consider files
// with extension ext.
func NewFileTreeReader(ext string) *FileTreeReader {
	return &FileTreeReader{
		FileExtension: ext,
		Targets:       MakeFileTargets(0),
	}
}

// ReadTree reads a tree rooted at prefix and generates a target from each file
// with extension FileExtension found in the tree.
func (ftr *FileTreeReader) ReadTree(prefix string) (FileTargets, error) {
	ftr.Prefix = prefix
	if err := filepath.Walk(prefix, ftr.WalkFunc); err != nil {
		return ftr.Targets, errors.Wrapf(err, "walking tree")
	}
	return ftr.Targets, nil
}

// WalkFunc processes a single filesystem object.
func (ftr *FileTreeReader) WalkFunc(p string, fi os.FileInfo, err error) error {
	if err != nil || fi.IsDir() || filepath.Ext(p) != "."+ftr.FileExtension {
		return err
	}
	path, err := filepath.Rel(ftr.Prefix, p)
	if err != nil {
		return errors.Wrapf(err, "getting relative path")
	}
	path = strings.TrimSuffix(path, "."+ftr.FileExtension)
	t := &FileTarget{
		FilePath: path,
	}
	return errors.Wrapf(ftr.Targets.Add(t), "adding file target %q", p)
}
