package hy

import "testing"

func TestWriteContext_Path(t *testing.T) {
	c := WriteContext{}
	c = c.Push("foo").Push("bar").Push("bat")
	expected := "foo/bar/bat"
	actual := c.Path()
	if actual != expected {
		t.Errorf("got %q; want %q", actual, expected)
	}
}
