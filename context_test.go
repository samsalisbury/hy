package hy

import "testing"

func TestNodeContext_Path(t *testing.T) {
	c := NodeContext{}
	c = c.Push(Tag{}, "Foo").Push(Tag{}, "Bar").Push(Tag{}, "Bat")
	expected := "Foo/Bar/Bat"
	actual := c.Path()
	if actual != expected {
		t.Errorf("got %q; want %q", actual, expected)
	}
}
