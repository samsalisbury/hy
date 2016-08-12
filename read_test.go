package hy

import "testing"

func TestCodec_Read(t *testing.T) {

	c := NewCodec(func(c *Codec) {
		c.TreeReader = NewFileTreeReader("json")
	})

	v := TestWriteStruct{}

	if err := c.Read("testdata/in", &v); err != nil {
		t.Fatal(err)
	}

}
