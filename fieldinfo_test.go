package hy

import (
	"fmt"
	"reflect"
	"testing"
)

type A struct {
	Name string
}

func (a *A) GetName() string     { return a.Name }
func (a *A) SetName(name string) { a.Name = name }

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

func TestNewFieldInfo_success(t *testing.T) {
	fitsType := reflect.TypeOf(FieldInfoTestStruct{})
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
		if actual.Name != field.Name {
			// No point repeating this in each test row.
			issue := fmt.Sprintf("%s == %q; want %q", "Name", actual.Name, field.Name)
			t.Errorf("%s for %s %s %# q", issue, field.Name, field.Type, field.Tag)
		}
	}
	t.Logf("%d assertions failed", numFailed)
}

type FieldInfoErrors struct {
	IllegalGet1  M `hy:",/"`
	IllegalGet2  M `hy:",_"`
	IllegalGet3  M `hy:",1"`
	IllegalGet4  M `hy:",."`
	IllegalGet5  M `hy:",1abc"`
	IllegalGet6  M `hy:",ab.c"`
	IllegalGet7  M `hy:",ab-c"`
	IllegalGet8  M `hy:",GetName"`   // no field named GetName (did you mean "GetName()"?
	IllegalGet9  M `hy:",GetName("`  // illegal token "GetName("
	IllegalGet10 M `hy:",GetName)"`  // illegal token "GetName)"
	IllegalGet11 M `hy:",Name()"`    // no method "Name"
	IllegalGet12 M `hy:",SetName()"` // wrong signature

	IllegalSet1 M `hy:",,Name()"`   // No method called "Name"
	IllegalSet2 M `hy:",,SetName"`  // setter must end with ()
	IllegalSet3 M `hy:",,SetName("` // illegal token "SetName("
	IllegalSet4 M `hy:",,SetName)"` // illegal token "SetName)"
	IllegalSet5 M `hy:",,/()"`      // illegal token /
	IllegalSet6 M `hy:",,_()"`      // illegal token _
	IllegalSet7 M `hy:",,1()"`      // illegal token 1
	IllegalSet8 M `hy:",,.()"`      // illegal token .
}

func quoteTag(tag string) string { return fmt.Sprintf("%# q", tag) }

var newFieldInfoBadCalls = map[string]string{
	"IllegalGet1":  `reading key field name: illegal token "/"`,
	"IllegalGet2":  `reading key field name: illegal token "_"`,
	"IllegalGet3":  `reading key field name: illegal token "1"`,
	"IllegalGet4":  `reading key field name: illegal token "."`,
	"IllegalGet5":  `reading key field name: illegal token "1abc"`,
	"IllegalGet6":  `reading key field name: illegal token "ab.c"`,
	"IllegalGet7":  `reading key field name: illegal token "ab-c"`,
	"IllegalGet8":  `reading key field name: *hy.A has no field "GetName""`,
	"IllegalGet9":  `reading key field name: illegal token "GetName("`,
	"IllegalGet10": `reading key field name: illegal token "GetName)"`,

	"IllegalGet11": `reading get key method name: *hy.A has no method "Name"`,
	"IllegalGet12": `reading get key method name: *hy.Ai.SetName() has wrong signature`,

	"IllegalSet1": `*hy.A has no method called "Name" for set key func name in`,
	"IllegalSet2": `reading set key method name: setter should end with "()"`,
	"IllegalSet3": `reading set key method name: illegal token "SetName("`,
	"IllegalSet4": `reading set key method name: illegal token "SetName)"`,
	"IllegalSet5": `reading set key method name: illegal token "/"`,
	"IllegalSet6": `reading set key method name: illegal token "_"`,
	"IllegalSet7": `reading set key method name: illegal token "1"`,
	"IllegalSet8": `reading set key method name: illegal token "."`,
}

func TestNewFieldInfo_failure(t *testing.T) {
	fieType := reflect.TypeOf(FieldInfoErrors{})
	for fieldName, expected := range newFieldInfoBadCalls {
		field, ok := fieType.FieldByName(fieldName)
		if !ok {
			t.Errorf("no field named %q", fieldName)
			continue
		}
		// complete the expectation
		expected = fmt.Sprintf("analysing field %s %s %# q: %s",
			field.Name, field.Type, field.Tag, expected)
		_, actualErr := NewFieldInfo(field)
		if actualErr == nil {
			t.Errorf("got nil; want error:\n\t%s", expected)
			continue
		}
		actual := actualErr.Error()
		if actual != expected {
			t.Errorf("got error:\n\t%s'\nwant:\n\t%s", actual, expected)
		}
	}
}
