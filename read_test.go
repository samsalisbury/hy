package hy

import (
	"encoding/json"
	"testing"
)

func TestCodec_Read(t *testing.T) {

	c := NewCodec(func(c *Codec) {
		c.TreeReader = NewFileTreeReader("json")
		c.Reader = JSONWriter
	})

	v := TestWriteStruct{}

	if err := c.Read("testdata/in", &v); err != nil {
		t.Fatal(err)
	}

	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	t.Fatal(string(b))
}
