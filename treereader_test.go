package hy

import "testing"

func TestFileTreeReader_ReadTree(t *testing.T) {

	tr := NewFileTreeReader("json", "_")

	targets, err := tr.ReadTree("testdata/in")
	if err != nil {
		t.Fatal(err)
	}

	actualLen := targets.Len()
	expectedLen := 21
	if actualLen != expectedLen {
		t.Errorf("got %d targets; want %d", actualLen, expectedLen)
	}

}
