package hy

import (
	"encoding/json"
	"os"
	"testing"
)

func TestCodec_Read(t *testing.T) {
	jsonWriter := JSONWriter
	jsonWriter.MarshalFunc = func(v interface{}) ([]byte, error) {
		return json.MarshalIndent(v, "", "  ")
	}
	c := NewCodec(func(c *Codec) {
		c.TreeReader = NewFileTreeReader("json", "_")
		c.Reader = jsonWriter
		c.Writer = jsonWriter
	})

	v := TestStruct{}

	if err := c.Read("testdata/in", &v); err != nil {
		t.Fatal(err)
	}

	if err := os.RemoveAll("testdata/roundtripped"); err != nil {
		t.Fatal(err)
	}
	if err := c.Write("testdata/roundtripped", v); err != nil {
		t.Fatal(err)
	}

	v2 := TestStruct{}
	if err := c.Read("testdata/roundtripped", &v2); err != nil {
		t.Fatal(err)
	}

	if err := os.RemoveAll("testdata/roundtripped2"); err != nil {
		t.Fatal(err)
	}
	if err := c.Write("testdata/roundtripped2", &v2); err != nil {
		t.Fatal(err)
	}
}
