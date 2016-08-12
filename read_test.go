package hy

import (
	"encoding/json"
	"os"
	"testing"
)

func TestCodec_Read(t *testing.T) {
	c := NewCodec(func(c *Codec) {
		c.TreeReader = NewFileTreeReader("json")
		c.Reader = JSONWriter
		c.Writer = JSONWriter
	})

	v := TestWriteStruct{}

	if err := c.Read("testdata/in", &v); err != nil {
		t.Fatal(err)
	}

	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	if err := os.RemoveAll("testdata/roundtripped"); err != nil {
		t.Fatal(err)
	}
	if err := c.Write("testdata/roundtripped", v); err != nil {
		t.Fatal(err)
	}

	t.Fatal(string(b))
}
