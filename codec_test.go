package hy

import (
	"encoding/json"
	"os"
	"sync/atomic"
	"testing"
)

func counter() *int64    { var c int64; return &c }
func increment(c *int64) { atomic.AddInt64(c, 1) }

const prefix = "testdata/out"

func TestCodec_Write(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	if err := os.RemoveAll(prefix); err != nil {
		t.Fatalf("failed to remove output dir: %s", err)
	}

	w := JSONWriter
	numCalls := counter()
	w.MarshalFunc = func(v interface{}) ([]byte, error) {
		increment(numCalls)
		return json.MarshalIndent(v, "", "  ")
	}
	c := NewCodec(func(c *Codec) {
		c.Writer = w
	})

	if err := c.Write(prefix, testWriteStructData); err != nil {
		t.Fatal(err)
	}

	expectedNumCalls := int64(23)
	if *numCalls != expectedNumCalls {
		t.Errorf("MarshalFunc called %d times; want %d", *numCalls, expectedNumCalls)
	}
}
