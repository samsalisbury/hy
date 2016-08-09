package hy

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

// FieldInfo is information about a field.
type FieldInfo struct {
	Name, FieldName, PathName, KeyField, GetKeyName, SetKeyName string
	// Type is the type of this fields.
	Type reflect.Type
	// Tag is the parsed hy tag.
	Tag Tag
	// Ignore indicates this field should not be written or read by hy.
	Ignore,
	// IsField indicates this is a regular field.
	IsField,
	// IsString indicates this field should be encoded as a string.
	IsString,
	// AutoFieldName indicates this field should use a field name derived from
	// the field's name.
	AutoFieldName,
	// IsDir indicates that this map or slice field should have its elements
	// stored in a directory.
	IsDir,
	// AutoPathName indicates the file or directory storing this field should
	// have its name derived from the field's name.
	AutoPathName,
	// OmitEmpty means this field should only be written if it is not empty,
	// according to the meaning of "not empty" defined by encoding/json.
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
	jsonTag := ParseJSONTag(f)

	var fieldName, pathName, keyField,
		getKeyName, setKeyName string
	var ignore, isField, isString, autoFieldName,
		isDir, autoPathName, omitEmpty bool

	if tag.Ignore || (tag.None && jsonTag.Ignore) {
		ignore = true
		goto done
	}

	if tag.None {
		isField = true
		if jsonTag.Name != "" {
			fieldName = jsonTag.Name
		} else {
			autoFieldName = true
		}
		omitEmpty = jsonTag.OmitEmpty
		isString = jsonTag.String
		goto done
	}

	isDir = tag.IsDir

	if tag.PathName == "." {
		autoPathName = true
	} else {
		pathName = tag.PathName
	}
	if strings.HasSuffix(tag.Key, "()") {
		getKeyName = strings.TrimSuffix(tag.Key, "()")
	} else {
		keyField = tag.Key
	}
	if strings.HasSuffix(tag.SetKey, "()") {
		setKeyName = strings.TrimSuffix(tag.SetKey, "()")
	}

done:
	return &FieldInfo{
		Name:          f.Name,
		FieldName:     fieldName,
		PathName:      pathName,
		KeyField:      keyField,
		GetKeyName:    getKeyName,
		SetKeyName:    setKeyName,
		Type:          f.Type,
		Tag:           tag,
		Ignore:        ignore,
		IsField:       isField,
		IsString:      isString,
		AutoFieldName: autoFieldName,
		IsDir:         isDir,
		AutoPathName:  autoPathName,
		OmitEmpty:     omitEmpty,
	}, nil
}

// Following code copied from https://golang.org/src/encoding/json

// tagOptions is the string following a comma in a struct field's "json"
// tag, or the empty string. It does not include the leading comma.
type tagOptions string

// parseTag splits a struct field's json tag into its name and
// comma-separated options.
func parseJSONTagOptions(tag string) (string, tagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], tagOptions(tag[idx+1:])
	}
	return tag, tagOptions("")
}

// Contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o tagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if s == optionName {
			return true
		}
		s = next
	}
	return false
}

// JSONTag represents a json field tag.
type JSONTag struct {
	Ignore, String, OmitEmpty bool
	Name                      string
}

// ParseJSONTag parses a json field tag from a struct field.
// Only strings, floats, integers, and booleans can be quoted.
func ParseJSONTag(field reflect.StructField) JSONTag {
	jsonTag := field.Tag.Get("json")
	var ignore, str, omitEmpty bool
	name, opts := parseJSONTagOptions(jsonTag)
	if opts.Contains("string") {
		switch field.Type.Kind() {
		case reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64,
			reflect.String:
			str = true
		}
	}
	if opts.Contains("omitempty") {
		omitEmpty = true
	}
	if jsonTag == "-" {
		ignore = true
	}
	return JSONTag{Name: name, Ignore: ignore, String: str, OmitEmpty: omitEmpty}
}
