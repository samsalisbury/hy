package hy

import "testing"

func TestWriteContext_Path(t *testing.T) {
	c := WriteContext{}
	c = c.Push(Tag{}, "foo").Push(Tag{}, "bar").Push(Tag{}, "bat")
	expected := "foo/bar/bat"
	actual := c.Path()
	if actual != expected {
		t.Errorf("got %q; want %q", actual, expected)
	}
}
