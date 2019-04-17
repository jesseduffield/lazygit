package toml

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
	"time"
)

type basicMarshalTestStruct struct {
	String     string                      `toml:"Zstring"`
	StringList []string                    `toml:"Ystrlist"`
	Sub        basicMarshalTestSubStruct   `toml:"Xsubdoc"`
	SubList    []basicMarshalTestSubStruct `toml:"Wsublist"`
}

type basicMarshalTestSubStruct struct {
	String2 string
}

var basicTestData = basicMarshalTestStruct{
	String:     "Hello",
	StringList: []string{"Howdy", "Hey There"},
	Sub:        basicMarshalTestSubStruct{"One"},
	SubList:    []basicMarshalTestSubStruct{{"Two"}, {"Three"}},
}

var basicTestToml = []byte(`Ystrlist = ["Howdy","Hey There"]
Zstring = "Hello"

[[Wsublist]]
  String2 = "Two"

[[Wsublist]]
  String2 = "Three"

[Xsubdoc]
  String2 = "One"
`)

var basicTestTomlOrdered = []byte(`Zstring = "Hello"
Ystrlist = ["Howdy","Hey There"]

[Xsubdoc]
  String2 = "One"

[[Wsublist]]
  String2 = "Two"

[[Wsublist]]
  String2 = "Three"
`)

func TestBasicMarshal(t *testing.T) {
	result, err := Marshal(basicTestData)
	if err != nil {
		t.Fatal(err)
	}
	expected := basicTestToml
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

func TestBasicMarshalOrdered(t *testing.T) {
	var result bytes.Buffer
	err := NewEncoder(&result).Order(OrderPreserve).Encode(basicTestData)
	if err != nil {
		t.Fatal(err)
	}
	expected := basicTestTomlOrdered
	if !bytes.Equal(result.Bytes(), expected) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result.Bytes())
	}
}

func TestBasicMarshalWithPointer(t *testing.T) {
	result, err := Marshal(&basicTestData)
	if err != nil {
		t.Fatal(err)
	}
	expected := basicTestToml
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

func TestBasicMarshalOrderedWithPointer(t *testing.T) {
	var result bytes.Buffer
	err := NewEncoder(&result).Order(OrderPreserve).Encode(&basicTestData)
	if err != nil {
		t.Fatal(err)
	}
	expected := basicTestTomlOrdered
	if !bytes.Equal(result.Bytes(), expected) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result.Bytes())
	}
}

func TestBasicUnmarshal(t *testing.T) {
	result := basicMarshalTestStruct{}
	err := Unmarshal(basicTestToml, &result)
	expected := basicTestData
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Bad unmarshal: expected %v, got %v", expected, result)
	}
}

type testDoc struct {
	Title       string            `toml:"title"`
	BasicLists  testDocBasicLists `toml:"basic_lists"`
	SubDocPtrs  []*testSubDoc     `toml:"subdocptrs"`
	BasicMap    map[string]string `toml:"basic_map"`
	Subdocs     testDocSubs       `toml:"subdoc"`
	Basics      testDocBasics     `toml:"basic"`
	SubDocList  []testSubDoc      `toml:"subdoclist"`
	err         int               `toml:"shouldntBeHere"`
	unexported  int               `toml:"shouldntBeHere"`
	Unexported2 int               `toml:"-"`
}

type testDocBasics struct {
	Uint       uint      `toml:"uint"`
	Bool       bool      `toml:"bool"`
	Float      float32   `toml:"float"`
	Int        int       `toml:"int"`
	String     *string   `toml:"string"`
	Date       time.Time `toml:"date"`
	unexported int       `toml:"shouldntBeHere"`
}

type testDocBasicLists struct {
	Floats  []*float32  `toml:"floats"`
	Bools   []bool      `toml:"bools"`
	Dates   []time.Time `toml:"dates"`
	Ints    []int       `toml:"ints"`
	UInts   []uint      `toml:"uints"`
	Strings []string    `toml:"strings"`
}

type testDocSubs struct {
	Second *testSubDoc `toml:"second"`
	First  testSubDoc  `toml:"first"`
}

type testSubDoc struct {
	Name       string `toml:"name"`
	unexported int    `toml:"shouldntBeHere"`
}

var biteMe = "Bite me"
var float1 float32 = 12.3
var float2 float32 = 45.6
var float3 float32 = 78.9
var subdoc = testSubDoc{"Second", 0}

var docData = testDoc{
	Title:       "TOML Marshal Testing",
	unexported:  0,
	Unexported2: 0,
	Basics: testDocBasics{
		Bool:       true,
		Date:       time.Date(1979, 5, 27, 7, 32, 0, 0, time.UTC),
		Float:      123.4,
		Int:        5000,
		Uint:       5001,
		String:     &biteMe,
		unexported: 0,
	},
	BasicLists: testDocBasicLists{
		Bools: []bool{true, false, true},
		Dates: []time.Time{
			time.Date(1979, 5, 27, 7, 32, 0, 0, time.UTC),
			time.Date(1980, 5, 27, 7, 32, 0, 0, time.UTC),
		},
		Floats:  []*float32{&float1, &float2, &float3},
		Ints:    []int{8001, 8001, 8002},
		Strings: []string{"One", "Two", "Three"},
		UInts:   []uint{5002, 5003},
	},
	BasicMap: map[string]string{
		"one": "one",
		"two": "two",
	},
	Subdocs: testDocSubs{
		First:  testSubDoc{"First", 0},
		Second: &subdoc,
	},
	SubDocList: []testSubDoc{
		{"List.First", 0},
		{"List.Second", 0},
	},
	SubDocPtrs: []*testSubDoc{&subdoc},
}

func TestDocMarshal(t *testing.T) {
	result, err := Marshal(docData)
	if err != nil {
		t.Fatal(err)
	}
	expected, _ := ioutil.ReadFile("marshal_test.toml")
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

func TestDocMarshalOrdered(t *testing.T) {
	var result bytes.Buffer
	err := NewEncoder(&result).Order(OrderPreserve).Encode(docData)
	if err != nil {
		t.Fatal(err)
	}
	expected, _ := ioutil.ReadFile("marshal_OrderPreserve_test.toml")
	if !bytes.Equal(result.Bytes(), expected) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result.Bytes())
	}
}

func TestDocMarshalPointer(t *testing.T) {
	result, err := Marshal(&docData)
	if err != nil {
		t.Fatal(err)
	}
	expected, _ := ioutil.ReadFile("marshal_test.toml")
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

func TestDocUnmarshal(t *testing.T) {
	result := testDoc{}
	tomlData, _ := ioutil.ReadFile("marshal_test.toml")
	err := Unmarshal(tomlData, &result)
	expected := docData
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, expected) {
		resStr, _ := json.MarshalIndent(result, "", "  ")
		expStr, _ := json.MarshalIndent(expected, "", "  ")
		t.Errorf("Bad unmarshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expStr, resStr)
	}
}

func TestDocPartialUnmarshal(t *testing.T) {
	result := testDocSubs{}

	tree, _ := LoadFile("marshal_test.toml")
	subTree := tree.Get("subdoc").(*Tree)
	err := subTree.Unmarshal(&result)
	expected := docData.Subdocs
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, expected) {
		resStr, _ := json.MarshalIndent(result, "", "  ")
		expStr, _ := json.MarshalIndent(expected, "", "  ")
		t.Errorf("Bad partial unmartial: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expStr, resStr)
	}
}

type tomlTypeCheckTest struct {
	name string
	item interface{}
	typ  int //0=primitive, 1=otherslice, 2=treeslice, 3=tree
}

func TestTypeChecks(t *testing.T) {
	tests := []tomlTypeCheckTest{
		{"integer", 2, 0},
		{"time", time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC), 0},
		{"stringlist", []string{"hello", "hi"}, 1},
		{"timelist", []time.Time{time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)}, 1},
		{"objectlist", []tomlTypeCheckTest{}, 2},
		{"object", tomlTypeCheckTest{}, 3},
	}

	for _, test := range tests {
		expected := []bool{false, false, false, false}
		expected[test.typ] = true
		result := []bool{
			isPrimitive(reflect.TypeOf(test.item)),
			isOtherSlice(reflect.TypeOf(test.item)),
			isTreeSlice(reflect.TypeOf(test.item)),
			isTree(reflect.TypeOf(test.item)),
		}
		if !reflect.DeepEqual(expected, result) {
			t.Errorf("Bad type check on %q: expected %v, got %v", test.name, expected, result)
		}
	}
}

type unexportedMarshalTestStruct struct {
	String      string                      `toml:"string"`
	StringList  []string                    `toml:"strlist"`
	Sub         basicMarshalTestSubStruct   `toml:"subdoc"`
	SubList     []basicMarshalTestSubStruct `toml:"sublist"`
	unexported  int                         `toml:"shouldntBeHere"`
	Unexported2 int                         `toml:"-"`
}

var unexportedTestData = unexportedMarshalTestStruct{
	String:      "Hello",
	StringList:  []string{"Howdy", "Hey There"},
	Sub:         basicMarshalTestSubStruct{"One"},
	SubList:     []basicMarshalTestSubStruct{{"Two"}, {"Three"}},
	unexported:  0,
	Unexported2: 0,
}

var unexportedTestToml = []byte(`string = "Hello"
strlist = ["Howdy","Hey There"]
unexported = 1
shouldntBeHere = 2

[subdoc]
  String2 = "One"

[[sublist]]
  String2 = "Two"

[[sublist]]
  String2 = "Three"
`)

func TestUnexportedUnmarshal(t *testing.T) {
	result := unexportedMarshalTestStruct{}
	err := Unmarshal(unexportedTestToml, &result)
	expected := unexportedTestData
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Bad unexported unmarshal: expected %v, got %v", expected, result)
	}
}

type errStruct struct {
	Bool   bool      `toml:"bool"`
	Date   time.Time `toml:"date"`
	Float  float64   `toml:"float"`
	Int    int16     `toml:"int"`
	String *string   `toml:"string"`
}

var errTomls = []string{
	"bool = truly\ndate = 1979-05-27T07:32:00Z\nfloat = 123.4\nint = 5000\nstring = \"Bite me\"",
	"bool = true\ndate = 1979-05-27T07:3200Z\nfloat = 123.4\nint = 5000\nstring = \"Bite me\"",
	"bool = true\ndate = 1979-05-27T07:32:00Z\nfloat = 123a4\nint = 5000\nstring = \"Bite me\"",
	"bool = true\ndate = 1979-05-27T07:32:00Z\nfloat = 123.4\nint = j000\nstring = \"Bite me\"",
	"bool = true\ndate = 1979-05-27T07:32:00Z\nfloat = 123.4\nint = 5000\nstring = Bite me",
	"bool = true\ndate = 1979-05-27T07:32:00Z\nfloat = 123.4\nint = 5000\nstring = Bite me",
	"bool = 1\ndate = 1979-05-27T07:32:00Z\nfloat = 123.4\nint = 5000\nstring = \"Bite me\"",
	"bool = true\ndate = 1\nfloat = 123.4\nint = 5000\nstring = \"Bite me\"",
	"bool = true\ndate = 1979-05-27T07:32:00Z\n\"sorry\"\nint = 5000\nstring = \"Bite me\"",
	"bool = true\ndate = 1979-05-27T07:32:00Z\nfloat = 123.4\nint = \"sorry\"\nstring = \"Bite me\"",
	"bool = true\ndate = 1979-05-27T07:32:00Z\nfloat = 123.4\nint = 5000\nstring = 1",
}

type mapErr struct {
	Vals map[string]float64
}

type intErr struct {
	Int1  int
	Int2  int8
	Int3  int16
	Int4  int32
	Int5  int64
	UInt1 uint
	UInt2 uint8
	UInt3 uint16
	UInt4 uint32
	UInt5 uint64
	Flt1  float32
	Flt2  float64
}

var intErrTomls = []string{
	"Int1 = []\nInt2 = 2\nInt3 = 3\nInt4 = 4\nInt5 = 5\nUInt1 = 1\nUInt2 = 2\nUInt3 = 3\nUInt4 = 4\nUInt5 = 5\nFlt1 = 1.0\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = []\nInt3 = 3\nInt4 = 4\nInt5 = 5\nUInt1 = 1\nUInt2 = 2\nUInt3 = 3\nUInt4 = 4\nUInt5 = 5\nFlt1 = 1.0\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = 2\nInt3 = []\nInt4 = 4\nInt5 = 5\nUInt1 = 1\nUInt2 = 2\nUInt3 = 3\nUInt4 = 4\nUInt5 = 5\nFlt1 = 1.0\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = 2\nInt3 = 3\nInt4 = []\nInt5 = 5\nUInt1 = 1\nUInt2 = 2\nUInt3 = 3\nUInt4 = 4\nUInt5 = 5\nFlt1 = 1.0\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = 2\nInt3 = 3\nInt4 = 4\nInt5 = []\nUInt1 = 1\nUInt2 = 2\nUInt3 = 3\nUInt4 = 4\nUInt5 = 5\nFlt1 = 1.0\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = 2\nInt3 = 3\nInt4 = 4\nInt5 = 5\nUInt1 = []\nUInt2 = 2\nUInt3 = 3\nUInt4 = 4\nUInt5 = 5\nFlt1 = 1.0\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = 2\nInt3 = 3\nInt4 = 4\nInt5 = 5\nUInt1 = 1\nUInt2 = []\nUInt3 = 3\nUInt4 = 4\nUInt5 = 5\nFlt1 = 1.0\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = 2\nInt3 = 3\nInt4 = 4\nInt5 = 5\nUInt1 = 1\nUInt2 = 2\nUInt3 = []\nUInt4 = 4\nUInt5 = 5\nFlt1 = 1.0\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = 2\nInt3 = 3\nInt4 = 4\nInt5 = 5\nUInt1 = 1\nUInt2 = 2\nUInt3 = 3\nUInt4 = []\nUInt5 = 5\nFlt1 = 1.0\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = 2\nInt3 = 3\nInt4 = 4\nInt5 = 5\nUInt1 = 1\nUInt2 = 2\nUInt3 = 3\nUInt4 = 4\nUInt5 = []\nFlt1 = 1.0\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = 2\nInt3 = 3\nInt4 = 4\nInt5 = 5\nUInt1 = 1\nUInt2 = 2\nUInt3 = 3\nUInt4 = 4\nUInt5 = 5\nFlt1 = []\nFlt2 = 2.0",
	"Int1 = 1\nInt2 = 2\nInt3 = 3\nInt4 = 4\nInt5 = 5\nUInt1 = 1\nUInt2 = 2\nUInt3 = 3\nUInt4 = 4\nUInt5 = 5\nFlt1 = 1.0\nFlt2 = []",
}

func TestErrUnmarshal(t *testing.T) {
	for ind, toml := range errTomls {
		result := errStruct{}
		err := Unmarshal([]byte(toml), &result)
		if err == nil {
			t.Errorf("Expected err from case %d\n", ind)
		}
	}
	result2 := mapErr{}
	err := Unmarshal([]byte("[Vals]\nfred=\"1.2\""), &result2)
	if err == nil {
		t.Errorf("Expected err from map")
	}
	for ind, toml := range intErrTomls {
		result3 := intErr{}
		err := Unmarshal([]byte(toml), &result3)
		if err == nil {
			t.Errorf("Expected int err from case %d\n", ind)
		}
	}
}

type emptyMarshalTestStruct struct {
	Title      string                  `toml:"title"`
	Bool       bool                    `toml:"bool"`
	Int        int                     `toml:"int"`
	String     string                  `toml:"string"`
	StringList []string                `toml:"stringlist"`
	Ptr        *basicMarshalTestStruct `toml:"ptr"`
	Map        map[string]string       `toml:"map"`
}

var emptyTestData = emptyMarshalTestStruct{
	Title:      "Placeholder",
	Bool:       false,
	Int:        0,
	String:     "",
	StringList: []string{},
	Ptr:        nil,
	Map:        map[string]string{},
}

var emptyTestToml = []byte(`bool = false
int = 0
string = ""
stringlist = []
title = "Placeholder"

[map]
`)

type emptyMarshalTestStruct2 struct {
	Title      string                  `toml:"title"`
	Bool       bool                    `toml:"bool,omitempty"`
	Int        int                     `toml:"int, omitempty"`
	String     string                  `toml:"string,omitempty "`
	StringList []string                `toml:"stringlist,omitempty"`
	Ptr        *basicMarshalTestStruct `toml:"ptr,omitempty"`
	Map        map[string]string       `toml:"map,omitempty"`
}

var emptyTestData2 = emptyMarshalTestStruct2{
	Title:      "Placeholder",
	Bool:       false,
	Int:        0,
	String:     "",
	StringList: []string{},
	Ptr:        nil,
	Map:        map[string]string{},
}

var emptyTestToml2 = []byte(`title = "Placeholder"
`)

func TestEmptyMarshal(t *testing.T) {
	result, err := Marshal(emptyTestData)
	if err != nil {
		t.Fatal(err)
	}
	expected := emptyTestToml
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad empty marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

func TestEmptyMarshalOmit(t *testing.T) {
	result, err := Marshal(emptyTestData2)
	if err != nil {
		t.Fatal(err)
	}
	expected := emptyTestToml2
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad empty omit marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

func TestEmptyUnmarshal(t *testing.T) {
	result := emptyMarshalTestStruct{}
	err := Unmarshal(emptyTestToml, &result)
	expected := emptyTestData
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Bad empty unmarshal: expected %v, got %v", expected, result)
	}
}

func TestEmptyUnmarshalOmit(t *testing.T) {
	result := emptyMarshalTestStruct2{}
	err := Unmarshal(emptyTestToml, &result)
	expected := emptyTestData2
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Bad empty omit unmarshal: expected %v, got %v", expected, result)
	}
}

type pointerMarshalTestStruct struct {
	Str       *string
	List      *[]string
	ListPtr   *[]*string
	Map       *map[string]string
	MapPtr    *map[string]*string
	EmptyStr  *string
	EmptyList *[]string
	EmptyMap  *map[string]string
	DblPtr    *[]*[]*string
}

var pointerStr = "Hello"
var pointerList = []string{"Hello back"}
var pointerListPtr = []*string{&pointerStr}
var pointerMap = map[string]string{"response": "Goodbye"}
var pointerMapPtr = map[string]*string{"alternate": &pointerStr}
var pointerTestData = pointerMarshalTestStruct{
	Str:       &pointerStr,
	List:      &pointerList,
	ListPtr:   &pointerListPtr,
	Map:       &pointerMap,
	MapPtr:    &pointerMapPtr,
	EmptyStr:  nil,
	EmptyList: nil,
	EmptyMap:  nil,
}

var pointerTestToml = []byte(`List = ["Hello back"]
ListPtr = ["Hello"]
Str = "Hello"

[Map]
  response = "Goodbye"

[MapPtr]
  alternate = "Hello"
`)

func TestPointerMarshal(t *testing.T) {
	result, err := Marshal(pointerTestData)
	if err != nil {
		t.Fatal(err)
	}
	expected := pointerTestToml
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad pointer marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

func TestPointerUnmarshal(t *testing.T) {
	result := pointerMarshalTestStruct{}
	err := Unmarshal(pointerTestToml, &result)
	expected := pointerTestData
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Bad pointer unmarshal: expected %v, got %v", expected, result)
	}
}

func TestUnmarshalTypeMismatch(t *testing.T) {
	result := pointerMarshalTestStruct{}
	err := Unmarshal([]byte("List = 123"), &result)
	if !strings.HasPrefix(err.Error(), "(1, 1): Can't convert 123(int64) to []string(slice)") {
		t.Errorf("Type mismatch must be reported: got %v", err.Error())
	}
}

type nestedMarshalTestStruct struct {
	String [][]string
	//Struct [][]basicMarshalTestSubStruct
	StringPtr *[]*[]*string
	// StructPtr *[]*[]*basicMarshalTestSubStruct
}

var str1 = "Three"
var str2 = "Four"
var strPtr = []*string{&str1, &str2}
var strPtr2 = []*[]*string{&strPtr}

var nestedTestData = nestedMarshalTestStruct{
	String:    [][]string{{"Five", "Six"}, {"One", "Two"}},
	StringPtr: &strPtr2,
}

var nestedTestToml = []byte(`String = [["Five","Six"],["One","Two"]]
StringPtr = [["Three","Four"]]
`)

func TestNestedMarshal(t *testing.T) {
	result, err := Marshal(nestedTestData)
	if err != nil {
		t.Fatal(err)
	}
	expected := nestedTestToml
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad nested marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

func TestNestedUnmarshal(t *testing.T) {
	result := nestedMarshalTestStruct{}
	err := Unmarshal(nestedTestToml, &result)
	expected := nestedTestData
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Bad nested unmarshal: expected %v, got %v", expected, result)
	}
}

type customMarshalerParent struct {
	Self    customMarshaler   `toml:"me"`
	Friends []customMarshaler `toml:"friends"`
}

type customMarshaler struct {
	FirsName string
	LastName string
}

func (c customMarshaler) MarshalTOML() ([]byte, error) {
	fullName := fmt.Sprintf("%s %s", c.FirsName, c.LastName)
	return []byte(fullName), nil
}

var customMarshalerData = customMarshaler{FirsName: "Sally", LastName: "Fields"}
var customMarshalerToml = []byte(`Sally Fields`)
var nestedCustomMarshalerData = customMarshalerParent{
	Self:    customMarshaler{FirsName: "Maiku", LastName: "Suteda"},
	Friends: []customMarshaler{customMarshalerData},
}
var nestedCustomMarshalerToml = []byte(`friends = ["Sally Fields"]
me = "Maiku Suteda"
`)

func TestCustomMarshaler(t *testing.T) {
	result, err := Marshal(customMarshalerData)
	if err != nil {
		t.Fatal(err)
	}
	expected := customMarshalerToml
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad custom marshaler: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

func TestNestedCustomMarshaler(t *testing.T) {
	result, err := Marshal(nestedCustomMarshalerData)
	if err != nil {
		t.Fatal(err)
	}
	expected := nestedCustomMarshalerToml
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad nested custom marshaler: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

var commentTestToml = []byte(`
# it's a comment on type
[postgres]
  # isCommented = "dvalue"
  noComment = "cvalue"

  # A comment on AttrB with a
  # break line
  password = "bvalue"

  # A comment on AttrA
  user = "avalue"

  [[postgres.My]]

    # a comment on my on typeC
    My = "Foo"

  [[postgres.My]]

    # a comment on my on typeC
    My = "Baar"
`)

func TestMarshalComment(t *testing.T) {
	type TypeC struct {
		My string `comment:"a comment on my on typeC"`
	}
	type TypeB struct {
		AttrA string `toml:"user" comment:"A comment on AttrA"`
		AttrB string `toml:"password" comment:"A comment on AttrB with a\n break line"`
		AttrC string `toml:"noComment"`
		AttrD string `toml:"isCommented" commented:"true"`
		My    []TypeC
	}
	type TypeA struct {
		TypeB TypeB `toml:"postgres" comment:"it's a comment on type"`
	}

	ta := []TypeC{{My: "Foo"}, {My: "Baar"}}
	config := TypeA{TypeB{AttrA: "avalue", AttrB: "bvalue", AttrC: "cvalue", AttrD: "dvalue", My: ta}}
	result, err := Marshal(config)
	if err != nil {
		t.Fatal(err)
	}
	expected := commentTestToml
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

type mapsTestStruct struct {
	Simple map[string]string
	Paths  map[string]string
	Other  map[string]float64
	X      struct {
		Y struct {
			Z map[string]bool
		}
	}
}

var mapsTestData = mapsTestStruct{
	Simple: map[string]string{
		"one plus one": "two",
		"next":         "three",
	},
	Paths: map[string]string{
		"/this/is/a/path": "/this/is/also/a/path",
		"/heloo.txt":      "/tmp/lololo.txt",
	},
	Other: map[string]float64{
		"testing": 3.9999,
	},
	X: struct{ Y struct{ Z map[string]bool } }{
		Y: struct{ Z map[string]bool }{
			Z: map[string]bool{
				"is.Nested": true,
			},
		},
	},
}
var mapsTestToml = []byte(`
[Other]
  "testing" = 3.9999

[Paths]
  "/heloo.txt" = "/tmp/lololo.txt"
  "/this/is/a/path" = "/this/is/also/a/path"

[Simple]
  "next" = "three"
  "one plus one" = "two"

[X]

  [X.Y]

    [X.Y.Z]
      "is.Nested" = true
`)

func TestEncodeQuotedMapKeys(t *testing.T) {
	var buf bytes.Buffer
	if err := NewEncoder(&buf).QuoteMapKeys(true).Encode(mapsTestData); err != nil {
		t.Fatal(err)
	}
	result := buf.Bytes()
	expected := mapsTestToml
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad maps marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

func TestDecodeQuotedMapKeys(t *testing.T) {
	result := mapsTestStruct{}
	err := NewDecoder(bytes.NewBuffer(mapsTestToml)).Decode(&result)
	expected := mapsTestData
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Bad maps unmarshal: expected %v, got %v", expected, result)
	}
}

type structArrayNoTag struct {
	A struct {
		B []int64
		C []int64
	}
}

func TestMarshalArray(t *testing.T) {
	expected := []byte(`
[A]
  B = [1,2,3]
  C = [1]
`)

	m := structArrayNoTag{
		A: struct {
			B []int64
			C []int64
		}{
			B: []int64{1, 2, 3},
			C: []int64{1},
		},
	}

	b, err := Marshal(m)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(b, expected) {
		t.Errorf("Bad arrays marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, b)
	}
}

func TestMarshalArrayOnePerLine(t *testing.T) {
	expected := []byte(`
[A]
  B = [
    1,
    2,
    3,
  ]
  C = [1]
`)

	m := structArrayNoTag{
		A: struct {
			B []int64
			C []int64
		}{
			B: []int64{1, 2, 3},
			C: []int64{1},
		},
	}

	var buf bytes.Buffer
	encoder := NewEncoder(&buf).ArraysWithOneElementPerLine(true)
	err := encoder.Encode(m)

	if err != nil {
		t.Fatal(err)
	}

	b := buf.Bytes()

	if !bytes.Equal(b, expected) {
		t.Errorf("Bad arrays marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, b)
	}
}

var customTagTestToml = []byte(`
[postgres]
  password = "bvalue"
  user = "avalue"

  [[postgres.My]]
    My = "Foo"

  [[postgres.My]]
    My = "Baar"
`)

func TestMarshalCustomTag(t *testing.T) {
	type TypeC struct {
		My string
	}
	type TypeB struct {
		AttrA string `file:"user"`
		AttrB string `file:"password"`
		My    []TypeC
	}
	type TypeA struct {
		TypeB TypeB `file:"postgres"`
	}

	ta := []TypeC{{My: "Foo"}, {My: "Baar"}}
	config := TypeA{TypeB{AttrA: "avalue", AttrB: "bvalue", My: ta}}
	var buf bytes.Buffer
	err := NewEncoder(&buf).SetTagName("file").Encode(config)
	if err != nil {
		t.Fatal(err)
	}
	expected := customTagTestToml
	result := buf.Bytes()
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

var customCommentTagTestToml = []byte(`
# db connection
[postgres]

  # db pass
  password = "bvalue"

  # db user
  user = "avalue"
`)

func TestMarshalCustomComment(t *testing.T) {
	type TypeB struct {
		AttrA string `toml:"user" descr:"db user"`
		AttrB string `toml:"password" descr:"db pass"`
	}
	type TypeA struct {
		TypeB TypeB `toml:"postgres" descr:"db connection"`
	}

	config := TypeA{TypeB{AttrA: "avalue", AttrB: "bvalue"}}
	var buf bytes.Buffer
	err := NewEncoder(&buf).SetTagComment("descr").Encode(config)
	if err != nil {
		t.Fatal(err)
	}
	expected := customCommentTagTestToml
	result := buf.Bytes()
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

var customCommentedTagTestToml = []byte(`
[postgres]
  # password = "bvalue"
  # user = "avalue"
`)

func TestMarshalCustomCommented(t *testing.T) {
	type TypeB struct {
		AttrA string `toml:"user" disable:"true"`
		AttrB string `toml:"password" disable:"true"`
	}
	type TypeA struct {
		TypeB TypeB `toml:"postgres"`
	}

	config := TypeA{TypeB{AttrA: "avalue", AttrB: "bvalue"}}
	var buf bytes.Buffer
	err := NewEncoder(&buf).SetTagCommented("disable").Encode(config)
	if err != nil {
		t.Fatal(err)
	}
	expected := customCommentedTagTestToml
	result := buf.Bytes()
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

var customMultilineTagTestToml = []byte(`int_slice = [
  1,
  2,
  3,
]
`)

func TestMarshalCustomMultiline(t *testing.T) {
	type TypeA struct {
		AttrA []int `toml:"int_slice" mltln:"true"`
	}

	config := TypeA{AttrA: []int{1, 2, 3}}
	var buf bytes.Buffer
	err := NewEncoder(&buf).ArraysWithOneElementPerLine(true).SetTagMultiline("mltln").Encode(config)
	if err != nil {
		t.Fatal(err)
	}
	expected := customMultilineTagTestToml
	result := buf.Bytes()
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

var testDocBasicToml = []byte(`
[document]
  bool_val = true
  date_val = 1979-05-27T07:32:00Z
  float_val = 123.4
  int_val = 5000
  string_val = "Bite me"
  uint_val = 5001
`)

type testDocCustomTag struct {
	Doc testDocBasicsCustomTag `file:"document"`
}
type testDocBasicsCustomTag struct {
	Bool       bool      `file:"bool_val"`
	Date       time.Time `file:"date_val"`
	Float      float32   `file:"float_val"`
	Int        int       `file:"int_val"`
	Uint       uint      `file:"uint_val"`
	String     *string   `file:"string_val"`
	unexported int       `file:"shouldntBeHere"`
}

var testDocCustomTagData = testDocCustomTag{
	Doc: testDocBasicsCustomTag{
		Bool:       true,
		Date:       time.Date(1979, 5, 27, 7, 32, 0, 0, time.UTC),
		Float:      123.4,
		Int:        5000,
		Uint:       5001,
		String:     &biteMe,
		unexported: 0,
	},
}

func TestUnmarshalCustomTag(t *testing.T) {
	buf := bytes.NewBuffer(testDocBasicToml)

	result := testDocCustomTag{}
	err := NewDecoder(buf).SetTagName("file").Decode(&result)
	if err != nil {
		t.Fatal(err)
	}
	expected := testDocCustomTagData
	if !reflect.DeepEqual(result, expected) {
		resStr, _ := json.MarshalIndent(result, "", "  ")
		expStr, _ := json.MarshalIndent(expected, "", "  ")
		t.Errorf("Bad unmarshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expStr, resStr)

	}
}

func TestUnmarshalMap(t *testing.T) {
	m := make(map[string]int)
	m["a"] = 1

	err := Unmarshal(basicTestToml, m)
	if err.Error() != "Only a pointer to struct can be unmarshaled from TOML" {
		t.Fail()
	}
}

func TestMarshalSlice(t *testing.T) {
	m := make([]int, 1)
	m[0] = 1

	var buf bytes.Buffer
	err := NewEncoder(&buf).Encode(&m)
	if err == nil {
		t.Error("expected error, got nil")
		return
	}
	if err.Error() != "Only pointer to struct can be marshaled to TOML" {
		t.Fail()
	}
}

func TestMarshalSlicePointer(t *testing.T) {
	m := make([]int, 1)
	m[0] = 1

	var buf bytes.Buffer
	err := NewEncoder(&buf).Encode(m)
	if err == nil {
		t.Error("expected error, got nil")
		return
	}
	if err.Error() != "Only a struct or map can be marshaled to TOML" {
		t.Fail()
	}
}

type testDuration struct {
	Nanosec   time.Duration  `toml:"nanosec"`
	Microsec1 time.Duration  `toml:"microsec1"`
	Microsec2 *time.Duration `toml:"microsec2"`
	Millisec  time.Duration  `toml:"millisec"`
	Sec       time.Duration  `toml:"sec"`
	Min       time.Duration  `toml:"min"`
	Hour      time.Duration  `toml:"hour"`
	Mixed     time.Duration  `toml:"mixed"`
	AString   string         `toml:"a_string"`
}

var testDurationToml = []byte(`
nanosec = "1ns"
microsec1 = "1us"
microsec2 = "1µs"
millisec = "1ms"
sec = "1s"
min = "1m"
hour = "1h"
mixed = "1h1m1s1ms1µs1ns"
a_string = "15s"
`)

func TestUnmarshalDuration(t *testing.T) {
	buf := bytes.NewBuffer(testDurationToml)

	result := testDuration{}
	err := NewDecoder(buf).Decode(&result)
	if err != nil {
		t.Fatal(err)
	}
	ms := time.Duration(1) * time.Microsecond
	expected := testDuration{
		Nanosec:   1,
		Microsec1: time.Microsecond,
		Microsec2: &ms,
		Millisec:  time.Millisecond,
		Sec:       time.Second,
		Min:       time.Minute,
		Hour:      time.Hour,
		Mixed: time.Hour +
			time.Minute +
			time.Second +
			time.Millisecond +
			time.Microsecond +
			time.Nanosecond,
		AString: "15s",
	}
	if !reflect.DeepEqual(result, expected) {
		resStr, _ := json.MarshalIndent(result, "", "  ")
		expStr, _ := json.MarshalIndent(expected, "", "  ")
		t.Errorf("Bad unmarshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expStr, resStr)

	}
}

var testDurationToml2 = []byte(`a_string = "15s"
hour = "1h0m0s"
microsec1 = "1µs"
microsec2 = "1µs"
millisec = "1ms"
min = "1m0s"
mixed = "1h1m1.001001001s"
nanosec = "1ns"
sec = "1s"
`)

func TestMarshalDuration(t *testing.T) {
	ms := time.Duration(1) * time.Microsecond
	data := testDuration{
		Nanosec:   1,
		Microsec1: time.Microsecond,
		Microsec2: &ms,
		Millisec:  time.Millisecond,
		Sec:       time.Second,
		Min:       time.Minute,
		Hour:      time.Hour,
		Mixed: time.Hour +
			time.Minute +
			time.Second +
			time.Millisecond +
			time.Microsecond +
			time.Nanosecond,
		AString: "15s",
	}

	var buf bytes.Buffer
	err := NewEncoder(&buf).Encode(data)
	if err != nil {
		t.Fatal(err)
	}
	expected := testDurationToml2
	result := buf.Bytes()
	if !bytes.Equal(result, expected) {
		t.Errorf("Bad marshal: expected\n-----\n%s\n-----\ngot\n-----\n%s\n-----\n", expected, result)
	}
}

type testBadDuration struct {
	Val time.Duration `toml:"val"`
}

var testBadDurationToml = []byte(`val = "1z"`)

func TestUnmarshalBadDuration(t *testing.T) {
	buf := bytes.NewBuffer(testBadDurationToml)

	result := testBadDuration{}
	err := NewDecoder(buf).Decode(&result)
	if err == nil {
		t.Fatal()
	}
	if err.Error() != "(1, 1): Can't convert 1z(string) to time.Duration. time: unknown unit z in duration 1z" {
		t.Fatalf("unexpected error: %s", err)
	}
}

var testCamelCaseKeyToml = []byte(`fooBar = 10`)

func TestUnmarshalCamelCaseKey(t *testing.T) {
	var x struct {
		FooBar int
		B      int
	}

	if err := Unmarshal(testCamelCaseKeyToml, &x); err != nil {
		t.Fatal(err)
	}

	if x.FooBar != 10 {
		t.Fatal("Did not set camelCase'd key")
	}
}

func TestUnmarshalDefault(t *testing.T) {
	var doc struct {
		StringField  string  `default:"a"`
		BoolField    bool    `default:"true"`
		IntField     int     `default:"1"`
		Int64Field   int64   `default:"2"`
		Float64Field float64 `default:"3.1"`
	}

	err := Unmarshal([]byte(``), &doc)
	if err != nil {
		t.Fatal(err)
	}
	if doc.BoolField != true {
		t.Errorf("BoolField should be true, not %t", doc.BoolField)
	}
	if doc.StringField != "a" {
		t.Errorf("StringField should be \"a\", not %s", doc.StringField)
	}
	if doc.IntField != 1 {
		t.Errorf("IntField should be 1, not %d", doc.IntField)
	}
	if doc.Int64Field != 2 {
		t.Errorf("Int64Field should be 2, not %d", doc.Int64Field)
	}
	if doc.Float64Field != 3.1 {
		t.Errorf("Float64Field should be 3.1, not %f", doc.Float64Field)
	}
}

func TestUnmarshalDefaultFailureBool(t *testing.T) {
	var doc struct {
		Field bool `default:"blah"`
	}

	err := Unmarshal([]byte(``), &doc)
	if err == nil {
		t.Fatal("should error")
	}
}

func TestUnmarshalDefaultFailureInt(t *testing.T) {
	var doc struct {
		Field int `default:"blah"`
	}

	err := Unmarshal([]byte(``), &doc)
	if err == nil {
		t.Fatal("should error")
	}
}

func TestUnmarshalDefaultFailureInt64(t *testing.T) {
	var doc struct {
		Field int64 `default:"blah"`
	}

	err := Unmarshal([]byte(``), &doc)
	if err == nil {
		t.Fatal("should error")
	}
}

func TestUnmarshalDefaultFailureFloat64(t *testing.T) {
	var doc struct {
		Field float64 `default:"blah"`
	}

	err := Unmarshal([]byte(``), &doc)
	if err == nil {
		t.Fatal("should error")
	}
}

func TestUnmarshalDefaultFailureUnsupported(t *testing.T) {
	var doc struct {
		Field struct{} `default:"blah"`
	}

	err := Unmarshal([]byte(``), &doc)
	if err == nil {
		t.Fatal("should error")
	}
}
