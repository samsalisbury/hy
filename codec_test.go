package hy

import (
	"encoding/json"
	"testing"
)

func TestCodec_Write(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	w := JSONWriter
	w.RootDir = "testdata/out"
	w.MarshalFunc = func(v interface{}) ([]byte, error) {
		return json.MarshalIndent(v, "  ", "  ")
	}
	c := NewCodec(func(c *Codec) {
		c.Writer = w
	})

	if err := c.Write(testWriteStructData); err != nil {
		t.Fatal(err)
	}
}
