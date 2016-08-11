package hy

import (
	"encoding/json"
	"sync/atomic"
	"testing"
)

func counter() *int64    { var c int64; return &c }
func increment(c *int64) { atomic.AddInt64(c, 1) }

func TestCodec_Write(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	w := JSONWriter
	w.RootDir = "testdata/out"
	numCalls := counter()
	w.MarshalFunc = func(v interface{}) ([]byte, error) {
		increment(numCalls)
		return json.MarshalIndent(v, "", "  ")
	}
	c := NewCodec(func(c *Codec) {
		c.Writer = w
	})

	if err := c.Write(testWriteStructData); err != nil {
		t.Fatal(err)
	}

	expectedNumCalls := int64(19)
	if *numCalls != expectedNumCalls {
		t.Errorf("MarshalFunc called %d times; want %d", *numCalls, expectedNumCalls)
	}
}
