package hy

import "testing"

type FileTargetsAssertion struct {
	Len int
	Err string
}

var goodFileTargets = map[FileTargetsAssertion]func() (FileTargets, error){
	{2, ""}: func() (FileTargets, error) {
		ts := MakeFileTargets(4)
		err := ts.Add(
			&FileTarget{FilePath: "a"},
			&FileTarget{FilePath: "b"},
		)
		return ts, err
	},
	{3, ""}: func() (FileTargets, error) {
		return NewFileTargets(
			&FileTarget{FilePath: "c"},
			&FileTarget{FilePath: "d"},
			&FileTarget{FilePath: "e"},
		)
	},
	{5, ""}: func() (FileTargets, error) {

		ts := MakeFileTargets(4)
		err := ts.Add(
			&FileTarget{FilePath: "a"},
			&FileTarget{FilePath: "b"},
		)
		if err != nil {
			return ts, err
		}
		ts2, err := NewFileTargets(
			&FileTarget{FilePath: "c"},
			&FileTarget{FilePath: "d"},
			&FileTarget{FilePath: "e"},
		)
		if err != nil {
			return ts, err
		}
		err = ts.AddAll(ts2)
		return ts, err
	},
	{-1, `duplicate file target "a"`}: func() (FileTargets, error) {
		return NewFileTargets(
			&FileTarget{FilePath: "a"},
			&FileTarget{FilePath: "a"},
		)
	},
}

func TestFileTargets(t *testing.T) {
	for expected, f := range goodFileTargets {
		fts, err := f()
		if err == nil && expected.Err != "" {
			t.Errorf("got nil; want error %q", expected.Err)
		}
		if err != nil && expected.Err == "" {
			t.Errorf("got error %q; want nil", err)
		}
		if expected.Err != "" {
			actual := err.Error()
			if actual != expected.Err {
				t.Errorf("got error %q; want error %q", actual, expected.Err)
			}
		}
		if expected.Len == -1 {
			continue
		}
		if fts.Len() != expected.Len {
			t.Errorf("got len %d; want %d", fts.Len(), expected.Len)
		}
	}
}
