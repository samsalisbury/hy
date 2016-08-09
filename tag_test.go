package hy

import "testing"

var goodTagTable = map[Tag][]string{
	Tag{None: true}: {
		"",
	},
	Tag{Ignore: true}: {
		"-", "-,", "-,,",
	},
	Tag{PathName: "mypath"}: {
		"mypath", "mypath,", "mypath,,",
	},
	Tag{PathName: "mypath", Key: "MyID"}: {
		"mypath,MyID", "mypath,MyID,",
	},
	Tag{PathName: "mypath", Key: "MyID", SetKey: "SetMyID()"}: {
		"mypath,MyID,SetMyID()",
	},
	Tag{PathName: ".", IsDir: true}: {
		"./", "/",
	},
	Tag{PathName: ".", IsDir: false}: {
		".", ",",
	},
	Tag{PathName: "mypath", IsDir: true}: {
		"mypath/", "mypath/,", "mypath/,,",
	},
	Tag{PathName: "mypath", IsDir: true, Key: "MyID"}: {
		"mypath/,MyID", "mypath/,MyID,",
	},
	Tag{PathName: "mypath", IsDir: true, Key: "MyID", SetKey: "SetMyID()"}: {
		"mypath/,MyID,SetMyID()",
	},
}

func TestParseTag_success(t *testing.T) {
	for expected, inputs := range goodTagTable {
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

var badTagTable = map[string][]string{
	"malformed tag, too many commas":                     {",,,", "mypath,key,setkey,"},
	`path name "/mypath" invalid: must not begin with /`: {"/mypath", "/mypath,", "/mypath,,"},
}

func TestParseTag_failure(t *testing.T) {
	for expected, inputs := range badTagTable {
		for _, input := range inputs {
			_, actualErr := parseTag(input)
			if actualErr == nil {
				t.Errorf("got nil; want error %q", expected)
				continue
			}
			actual := actualErr.Error()
			if actual != expected {
				t.Errorf("got error %q; want error %q", actual, expected)
			}
		}
	}
}
