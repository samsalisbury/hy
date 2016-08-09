package hy

import (
	"fmt"
	"reflect"
	"testing"
)

type A struct {
	Name string
}
type M map[string]A
type MP map[string]*A

type FieldInfoTestStruct struct {
	// no tag
	NoTag1 *A ``      // IsField
	NoTag2 *A `hy:""` // IsField

	// ignore
	Ignore1 *A `hy:"-"`   // Ignore
	Ignore2 *A `json:"-"` // Ignore

	// hy path tags
	FileTag      *A `hy:"some-path1"`  // PathName = "some-path1"
	DirTag       *A `hy:"some-path2/"` // IsDir + PathName  = "some-path2"
	AutoFileTag1 *A `hy:"."`           // AutoPathName
	AutoFileTag2 *A `hy:","`           // AutoPathName
	AutoDirTag1  *A `hy:"/"`           // AutoPathName + IsDir

	// json tags
	JSONName          *A  `json:"jsonName1"`           // IsField +             FieldName = jsonName1
	JSONOmitEmpty     *A  `json:",omitempty"`          // IsField + OmitEmpty + AutoFieldName
	JSONNameOmitEmpty *A  `json:"jsonName2,omitempty"` // IsField + OmitEmpty + FieldName = jsonName2
	JSONIntString     int `json:",string"`             // IsField + IsString  + AutoFieldName
	JSONIntStringOmit int `json:",string,omitempty"`   // IsField + IsString  + AutoFieldName

	// hy path tags + json tags
	FileTagWithJSON *A `hy:"some-path3" json:"x"`  // PathName  = "some-path3"
	DirTagWithJSON  *A `hy:"some-path4/" json:"x"` // PathName  = "some-path4" + IsDir

	// hy key field tags
	KeyFieldTag1 M  `hy:"/,Name"` // AutoPathName + KeyField = "Name" + IsDir
	KeyFieldTag2 M  `hy:",Name"`  // AutoPathName + KeyField = "Name" + IsDir
	KeyFieldTag3 MP `hy:"/,Name"` // AutoPathName + KeyField = "Name" + IsDir
	KeyFieldTag4 MP `hy:",Name"`  // AutoPathName + KeyField = "Name" + IsDir

	// hy key get/set tags
	KeyGetSet1 M  `hy:"/,GetName(),SetName()"` // AutoPathName + IsDir + GetKey = "GetName" + SetKey = "SetName"
	KeyGetSet2 M  `hy:",GetName(),SetName()"`  // AutoPathName + IsDir + GetKey = "GetName" + SetKey = "SetName"
	KeyGetSet3 MP `hy:"/,GetName(),SetName()"` // AutoPathName + IsDir + GetKey = "GetName" + SetKey = "SetName"
	KeyGetSet4 MP `hy:",GetName(),SetName()"`  // AutoPathName + IsDir + GetKey = "GetName" + SetKey = "SetName"
}

var fitsType = reflect.TypeOf(FieldInfoTestStruct{})

var fieldInfoGoodCalls = map[string]FieldInfo{
	"NoTag1": {IsField: true, AutoFieldName: true},
	"NoTag2": {IsField: true, AutoFieldName: true},

	"FileTag":      {PathName: "some-path1"},
	"DirTag":       {PathName: "some-path2", IsDir: true},
	"AutoFileTag1": {AutoPathName: true},
	"AutoFileTag2": {AutoPathName: true},
	"AutoDirTag1":  {AutoPathName: true, IsDir: true},

	"JSONName":          {IsField: true, FieldName: "jsonName1"},
	"JSONOmitEmpty":     {IsField: true, AutoFieldName: true, OmitEmpty: true},
	"JSONNameOmitEmpty": {IsField: true, FieldName: "jsonName2", OmitEmpty: true},

	"FileTagWithJSON": {PathName: "some-path3"},
	"DirTagWithJSON":  {PathName: "some-path4", IsDir: true},

	"KeyFieldTag1": {AutoPathName: true, KeyField: "Name", IsDir: true},
	"KeyFieldTag2": {AutoPathName: true, KeyField: "Name"},
	"KeyFieldTag3": {AutoPathName: true, KeyField: "Name", IsDir: true},
	"KeyFieldTag4": {AutoPathName: true, KeyField: "Name"},
}

func TestNewFieldInfo(t *testing.T) {
	numFailed := 0
	for fieldName, expected := range fieldInfoGoodCalls {
		field, ok := fitsType.FieldByName(fieldName)
		if !ok {
			t.Errorf("no field named %q", fieldName)
			continue
		}
		actual, err := NewFieldInfo(field)
		if err != nil {
			t.Error(err)
			continue
		}
		actualVal := reflect.ValueOf(actual).Elem()
		expectedVal := reflect.ValueOf(expected)
		// check bool fields
		for _, f := range []string{
			"Ignore", "IsField", "IsDir", "AutoPathName", "OmitEmpty", "AutoFieldName"} {
			actualBool := actualVal.FieldByName(f).Bool()
			expectedBool := expectedVal.FieldByName(f).Bool()
			if actualBool != expectedBool {
				numFailed++
				issue := f
				if expectedBool {
					issue = "!" + issue
				}
				t.Errorf("%-14s for %s %s %# q", issue, field.Name, field.Type, field.Tag)
			}
		}
		// check string fields
		for _, f := range []string{
			"FieldName", "PathName", "KeyField", "GetKeyName", "SetKeyName"} {
			actualString := actualVal.FieldByName(f).Interface().(string)
			expectedString := expectedVal.FieldByName(f).Interface().(string)
			if actualString != expectedString {
				numFailed++
				issue := fmt.Sprintf("%s == %q; want %q", f, actualString, expectedString)
				t.Errorf("%s for %s %s %# q", issue, field.Name, field.Type, field.Tag)
			}
		}
	}
	t.Logf("%d assertions failed", numFailed)
}
