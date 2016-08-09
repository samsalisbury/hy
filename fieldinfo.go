package hy

import (
	"reflect"

	"github.com/pkg/errors"
)

// FieldInfo is information about a field.
type FieldInfo struct {
	Name, FieldName, PathName, KeyField string
	Type                                reflect.Type
	Tag                                 Tag
	IsField,
	AutoFieldName,
	IsDir,
	AutoPathName,
	OmitEmpty bool
}

// NewFieldInfo creates a new FieldInfo, analysing the tag and checking the
// tag's named ID field or ID get/set methods for consistency.
func NewFieldInfo(f reflect.StructField) (*FieldInfo, error) {
	tagStr := f.Tag.Get("hy")
	tag, err := parseTag(tagStr)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid tag %q", tagStr)
	}
	//if tag.Key == ""
	return &FieldInfo{
		Name: f.Name,
		Type: f.Type,
		Tag:  tag,
	}, nil
}
