package hy

import (
	"encoding/json"
	"testing"
)

func TestCodec_Read(t *testing.T) {

	c := NewCodec(func(c *Codec) {
		c.UnmarshalFunc = json.Unmarshal
	})

	v := TestWriteStruct{}
	if err := c.Read("testdata/in", "json", &v); err != nil {
		t.Fatal(err)
	}

}
