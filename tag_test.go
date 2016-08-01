package hy

import "testing"

var tagTable = map[Tag][]string{
	Tag{}:                                                     {""},
	Tag{Ignore: true}:                                         {"-", "-,", "-,,"},
	Tag{PathName: "mypath"}:                                   {"mypath", "mypath,", "mypath,,"},
	Tag{PathName: "mypath", Key: "MyID"}:                      {"mypath,MyID", "mypath,MyID,"},
	Tag{PathName: "mypath", Key: "MyID", SetKey: "SetMyID()"}: {"mypath,MyID,SetMyID()"},
}

func TestParseTag_success(t *testing.T) {
	for expected, inputs := range tagTable {
		for _, input := range inputs {
			actual, err := parseTag(input)
			if err != nil {
				t.Error(err)
			}
			if actual != expected {
				t.Errorf("got %+v from %q; want %+v", actual, input, expected)
			}
		}
	}
}
