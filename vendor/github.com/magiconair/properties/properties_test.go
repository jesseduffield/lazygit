// Copyright 2018 Frank Schroeder. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package properties

import (
	"bytes"
	"flag"
	"fmt"
	"math"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
)

var verbose = flag.Bool("verbose", false, "Verbose output")

func init() {
	ErrorHandler = PanicHandler
}

// ----------------------------------------------------------------------------

// define test cases in the form of
// {"input", "key1", "value1", "key2", "value2", ...}
var complexTests = [][]string{
	// whitespace prefix
	{" key=value", "key", "value"},     // SPACE prefix
	{"\fkey=value", "key", "value"},    // FF prefix
	{"\tkey=value", "key", "value"},    // TAB prefix
	{" \f\tkey=value", "key", "value"}, // mix prefix

	// multiple keys
	{"key1=value1\nkey2=value2\n", "key1", "value1", "key2", "value2"},
	{"key1=value1\rkey2=value2\r", "key1", "value1", "key2", "value2"},
	{"key1=value1\r\nkey2=value2\r\n", "key1", "value1", "key2", "value2"},

	// blank lines
	{"\nkey=value\n", "key", "value"},
	{"\rkey=value\r", "key", "value"},
	{"\r\nkey=value\r\n", "key", "value"},
	{"\nkey=value\n \nkey2=value2", "key", "value", "key2", "value2"},
	{"\nkey=value\n\t\nkey2=value2", "key", "value", "key2", "value2"},

	// escaped chars in key
	{"k\\ ey = value", "k ey", "value"},
	{"k\\:ey = value", "k:ey", "value"},
	{"k\\=ey = value", "k=ey", "value"},
	{"k\\fey = value", "k\fey", "value"},
	{"k\\ney = value", "k\ney", "value"},
	{"k\\rey = value", "k\rey", "value"},
	{"k\\tey = value", "k\tey", "value"},

	// escaped chars in value
	{"key = v\\ alue", "key", "v alue"},
	{"key = v\\:alue", "key", "v:alue"},
	{"key = v\\=alue", "key", "v=alue"},
	{"key = v\\falue", "key", "v\falue"},
	{"key = v\\nalue", "key", "v\nalue"},
	{"key = v\\ralue", "key", "v\ralue"},
	{"key = v\\talue", "key", "v\talue"},

	// silently dropped escape character
	{"k\\zey = value", "kzey", "value"},
	{"key = v\\zalue", "key", "vzalue"},

	// unicode literals
	{"key\\u2318 = value", "key⌘", "value"},
	{"k\\u2318ey = value", "k⌘ey", "value"},
	{"key = value\\u2318", "key", "value⌘"},
	{"key = valu\\u2318e", "key", "valu⌘e"},

	// multiline values
	{"key = valueA,\\\n    valueB", "key", "valueA,valueB"},   // SPACE indent
	{"key = valueA,\\\n\f\f\fvalueB", "key", "valueA,valueB"}, // FF indent
	{"key = valueA,\\\n\t\t\tvalueB", "key", "valueA,valueB"}, // TAB indent
	{"key = valueA,\\\n \f\tvalueB", "key", "valueA,valueB"},  // mix indent

	// comments
	{"# this is a comment\n! and so is this\nkey1=value1\nkey#2=value#2\n\nkey!3=value!3\n# and another one\n! and the final one", "key1", "value1", "key#2", "value#2", "key!3", "value!3"},

	// expansion tests
	{"key=value\nkey2=${key}", "key", "value", "key2", "value"},
	{"key=value\nkey2=aa${key}", "key", "value", "key2", "aavalue"},
	{"key=value\nkey2=${key}bb", "key", "value", "key2", "valuebb"},
	{"key=value\nkey2=aa${key}bb", "key", "value", "key2", "aavaluebb"},
	{"key=value\nkey2=${key}\nkey3=${key2}", "key", "value", "key2", "value", "key3", "value"},
	{"key=value\nkey2=${key}${key}", "key", "value", "key2", "valuevalue"},
	{"key=value\nkey2=${key}${key}${key}${key}", "key", "value", "key2", "valuevaluevaluevalue"},
	{"key=value\nkey2=${key}${key3}\nkey3=${key}", "key", "value", "key2", "valuevalue", "key3", "value"},
	{"key=value\nkey2=${key3}${key}${key4}\nkey3=${key}\nkey4=${key}", "key", "value", "key2", "valuevaluevalue", "key3", "value", "key4", "value"},
	{"key=${USER}", "key", os.Getenv("USER")},
	{"key=${USER}\nUSER=value", "key", "value", "USER", "value"},
}

// ----------------------------------------------------------------------------

var commentTests = []struct {
	input, key, value string
	comments          []string
}{
	{"key=value", "key", "value", nil},
	{"#\nkey=value", "key", "value", []string{""}},
	{"#comment\nkey=value", "key", "value", []string{"comment"}},
	{"# comment\nkey=value", "key", "value", []string{"comment"}},
	{"#  comment\nkey=value", "key", "value", []string{"comment"}},
	{"# comment\n\nkey=value", "key", "value", []string{"comment"}},
	{"# comment1\n# comment2\nkey=value", "key", "value", []string{"comment1", "comment2"}},
	{"# comment1\n\n# comment2\n\nkey=value", "key", "value", []string{"comment1", "comment2"}},
	{"!comment\nkey=value", "key", "value", []string{"comment"}},
	{"! comment\nkey=value", "key", "value", []string{"comment"}},
	{"!  comment\nkey=value", "key", "value", []string{"comment"}},
	{"! comment\n\nkey=value", "key", "value", []string{"comment"}},
	{"! comment1\n! comment2\nkey=value", "key", "value", []string{"comment1", "comment2"}},
	{"! comment1\n\n! comment2\n\nkey=value", "key", "value", []string{"comment1", "comment2"}},
}

// ----------------------------------------------------------------------------

var errorTests = []struct {
	input, msg string
}{
	// unicode literals
	{"key\\u1 = value", "invalid unicode literal"},
	{"key\\u12 = value", "invalid unicode literal"},
	{"key\\u123 = value", "invalid unicode literal"},
	{"key\\u123g = value", "invalid unicode literal"},
	{"key\\u123", "invalid unicode literal"},

	// circular references
	{"key=${key}", "circular reference"},
	{"key1=${key2}\nkey2=${key1}", "circular reference"},

	// malformed expressions
	{"key=${ke", "malformed expression"},
	{"key=valu${ke", "malformed expression"},
}

// ----------------------------------------------------------------------------

var writeTests = []struct {
	input, output, encoding string
}{
	// ISO-8859-1 tests
	{"key = value", "key = value\n", "ISO-8859-1"},
	{"key = value \\\n   continued", "key = value continued\n", "ISO-8859-1"},
	{"key⌘ = value", "key\\u2318 = value\n", "ISO-8859-1"},
	{"ke\\ \\:y = value", "ke\\ \\:y = value\n", "ISO-8859-1"},

	// UTF-8 tests
	{"key = value", "key = value\n", "UTF-8"},
	{"key = value \\\n   continued", "key = value continued\n", "UTF-8"},
	{"key⌘ = value⌘", "key⌘ = value⌘\n", "UTF-8"},
	{"ke\\ \\:y = value", "ke\\ \\:y = value\n", "UTF-8"},
}

// ----------------------------------------------------------------------------

var writeCommentTests = []struct {
	input, output, encoding string
}{
	// ISO-8859-1 tests
	{"key = value", "key = value\n", "ISO-8859-1"},
	{"#\nkey = value", "key = value\n", "ISO-8859-1"},
	{"#\n#\n#\nkey = value", "key = value\n", "ISO-8859-1"},
	{"# comment\nkey = value", "# comment\nkey = value\n", "ISO-8859-1"},
	{"\n# comment\nkey = value", "# comment\nkey = value\n", "ISO-8859-1"},
	{"# comment\n\nkey = value", "# comment\nkey = value\n", "ISO-8859-1"},
	{"# comment1\n# comment2\nkey = value", "# comment1\n# comment2\nkey = value\n", "ISO-8859-1"},
	{"#comment1\nkey1 = value1\n#comment2\nkey2 = value2", "# comment1\nkey1 = value1\n\n# comment2\nkey2 = value2\n", "ISO-8859-1"},

	// UTF-8 tests
	{"key = value", "key = value\n", "UTF-8"},
	{"# comment⌘\nkey = value⌘", "# comment⌘\nkey = value⌘\n", "UTF-8"},
	{"\n# comment⌘\nkey = value⌘", "# comment⌘\nkey = value⌘\n", "UTF-8"},
	{"# comment⌘\n\nkey = value⌘", "# comment⌘\nkey = value⌘\n", "UTF-8"},
	{"# comment1⌘\n# comment2⌘\nkey = value⌘", "# comment1⌘\n# comment2⌘\nkey = value⌘\n", "UTF-8"},
	{"#comment1⌘\nkey1 = value1⌘\n#comment2⌘\nkey2 = value2⌘", "# comment1⌘\nkey1 = value1⌘\n\n# comment2⌘\nkey2 = value2⌘\n", "UTF-8"},
}

// ----------------------------------------------------------------------------

var boolTests = []struct {
	input, key string
	def, value bool
}{
	// valid values for TRUE
	{"key = 1", "key", false, true},
	{"key = on", "key", false, true},
	{"key = On", "key", false, true},
	{"key = ON", "key", false, true},
	{"key = true", "key", false, true},
	{"key = True", "key", false, true},
	{"key = TRUE", "key", false, true},
	{"key = yes", "key", false, true},
	{"key = Yes", "key", false, true},
	{"key = YES", "key", false, true},

	// valid values for FALSE (all other)
	{"key = 0", "key", true, false},
	{"key = off", "key", true, false},
	{"key = false", "key", true, false},
	{"key = no", "key", true, false},

	// non existent key
	{"key = true", "key2", false, false},
}

// ----------------------------------------------------------------------------

var durationTests = []struct {
	input, key string
	def, value time.Duration
}{
	// valid values
	{"key = 1", "key", 999, 1},
	{"key = 0", "key", 999, 0},
	{"key = -1", "key", 999, -1},
	{"key = 0123", "key", 999, 123},

	// invalid values
	{"key = 0xff", "key", 999, 999},
	{"key = 1.0", "key", 999, 999},
	{"key = a", "key", 999, 999},

	// non existent key
	{"key = 1", "key2", 999, 999},
}

// ----------------------------------------------------------------------------

var parsedDurationTests = []struct {
	input, key string
	def, value time.Duration
}{
	// valid values
	{"key = -1ns", "key", 999, -1 * time.Nanosecond},
	{"key = 300ms", "key", 999, 300 * time.Millisecond},
	{"key = 5s", "key", 999, 5 * time.Second},
	{"key = 3h", "key", 999, 3 * time.Hour},
	{"key = 2h45m", "key", 999, 2*time.Hour + 45*time.Minute},

	// invalid values
	{"key = 0xff", "key", 999, 999},
	{"key = 1.0", "key", 999, 999},
	{"key = a", "key", 999, 999},
	{"key = 1", "key", 999, 999},
	{"key = 0", "key", 999, 0},

	// non existent key
	{"key = 1", "key2", 999, 999},
}

// ----------------------------------------------------------------------------

var floatTests = []struct {
	input, key string
	def, value float64
}{
	// valid values
	{"key = 1.0", "key", 999, 1.0},
	{"key = 0.0", "key", 999, 0.0},
	{"key = -1.0", "key", 999, -1.0},
	{"key = 1", "key", 999, 1},
	{"key = 0", "key", 999, 0},
	{"key = -1", "key", 999, -1},
	{"key = 0123", "key", 999, 123},

	// invalid values
	{"key = 0xff", "key", 999, 999},
	{"key = a", "key", 999, 999},

	// non existent key
	{"key = 1", "key2", 999, 999},
}

// ----------------------------------------------------------------------------

var int64Tests = []struct {
	input, key string
	def, value int64
}{
	// valid values
	{"key = 1", "key", 999, 1},
	{"key = 0", "key", 999, 0},
	{"key = -1", "key", 999, -1},
	{"key = 0123", "key", 999, 123},

	// invalid values
	{"key = 0xff", "key", 999, 999},
	{"key = 1.0", "key", 999, 999},
	{"key = a", "key", 999, 999},

	// non existent key
	{"key = 1", "key2", 999, 999},
}

// ----------------------------------------------------------------------------

var uint64Tests = []struct {
	input, key string
	def, value uint64
}{
	// valid values
	{"key = 1", "key", 999, 1},
	{"key = 0", "key", 999, 0},
	{"key = 0123", "key", 999, 123},

	// invalid values
	{"key = -1", "key", 999, 999},
	{"key = 0xff", "key", 999, 999},
	{"key = 1.0", "key", 999, 999},
	{"key = a", "key", 999, 999},

	// non existent key
	{"key = 1", "key2", 999, 999},
}

// ----------------------------------------------------------------------------

var stringTests = []struct {
	input, key string
	def, value string
}{
	// valid values
	{"key = abc", "key", "def", "abc"},

	// non existent key
	{"key = abc", "key2", "def", "def"},
}

// ----------------------------------------------------------------------------

var keysTests = []struct {
	input string
	keys  []string
}{
	{"", []string{}},
	{"key = abc", []string{"key"}},
	{"key = abc\nkey2=def", []string{"key", "key2"}},
	{"key2 = abc\nkey=def", []string{"key2", "key"}},
	{"key = abc\nkey=def", []string{"key"}},
}

// ----------------------------------------------------------------------------

var filterTests = []struct {
	input   string
	pattern string
	keys    []string
	err     string
}{
	{"", "", []string{}, ""},
	{"", "abc", []string{}, ""},
	{"key=value", "", []string{"key"}, ""},
	{"key=value", "key=", []string{}, ""},
	{"key=value\nfoo=bar", "", []string{"foo", "key"}, ""},
	{"key=value\nfoo=bar", "f", []string{"foo"}, ""},
	{"key=value\nfoo=bar", "fo", []string{"foo"}, ""},
	{"key=value\nfoo=bar", "foo", []string{"foo"}, ""},
	{"key=value\nfoo=bar", "fooo", []string{}, ""},
	{"key=value\nkey2=value2\nfoo=bar", "ey", []string{"key", "key2"}, ""},
	{"key=value\nkey2=value2\nfoo=bar", "key", []string{"key", "key2"}, ""},
	{"key=value\nkey2=value2\nfoo=bar", "^key", []string{"key", "key2"}, ""},
	{"key=value\nkey2=value2\nfoo=bar", "^(key|foo)", []string{"foo", "key", "key2"}, ""},
	{"key=value\nkey2=value2\nfoo=bar", "[ abc", nil, "error parsing regexp.*"},
}

// ----------------------------------------------------------------------------

var filterPrefixTests = []struct {
	input  string
	prefix string
	keys   []string
}{
	{"", "", []string{}},
	{"", "abc", []string{}},
	{"key=value", "", []string{"key"}},
	{"key=value", "key=", []string{}},
	{"key=value\nfoo=bar", "", []string{"foo", "key"}},
	{"key=value\nfoo=bar", "f", []string{"foo"}},
	{"key=value\nfoo=bar", "fo", []string{"foo"}},
	{"key=value\nfoo=bar", "foo", []string{"foo"}},
	{"key=value\nfoo=bar", "fooo", []string{}},
	{"key=value\nkey2=value2\nfoo=bar", "key", []string{"key", "key2"}},
}

// ----------------------------------------------------------------------------

var filterStripPrefixTests = []struct {
	input  string
	prefix string
	keys   []string
}{
	{"", "", []string{}},
	{"", "abc", []string{}},
	{"key=value", "", []string{"key"}},
	{"key=value", "key=", []string{}},
	{"key=value\nfoo=bar", "", []string{"foo", "key"}},
	{"key=value\nfoo=bar", "f", []string{"foo"}},
	{"key=value\nfoo=bar", "fo", []string{"foo"}},
	{"key=value\nfoo=bar", "foo", []string{"foo"}},
	{"key=value\nfoo=bar", "fooo", []string{}},
	{"key=value\nkey2=value2\nfoo=bar", "key", []string{"key", "key2"}},
}

// ----------------------------------------------------------------------------

var setTests = []struct {
	input      string
	key, value string
	prev       string
	ok         bool
	err        string
	keys       []string
}{
	{"", "", "", "", false, "", []string{}},
	{"", "key", "value", "", false, "", []string{"key"}},
	{"key=value", "key2", "value2", "", false, "", []string{"key", "key2"}},
	{"key=value", "abc", "value3", "", false, "", []string{"key", "abc"}},
	{"key=value", "key", "value3", "value", true, "", []string{"key"}},
}

// ----------------------------------------------------------------------------

// TestBasic tests basic single key/value combinations with all possible
// whitespace, delimiter and newline permutations.
func TestBasic(t *testing.T) {
	testWhitespaceAndDelimiterCombinations(t, "key", "")
	testWhitespaceAndDelimiterCombinations(t, "key", "value")
	testWhitespaceAndDelimiterCombinations(t, "key", "value   ")
}

func TestComplex(t *testing.T) {
	for _, test := range complexTests {
		testKeyValue(t, test[0], test[1:]...)
	}
}

func TestErrors(t *testing.T) {
	for _, test := range errorTests {
		_, err := Load([]byte(test.input), ISO_8859_1)
		assert.Equal(t, err != nil, true, "want error")
		assert.Equal(t, strings.Contains(err.Error(), test.msg), true)
	}
}

func TestVeryDeep(t *testing.T) {
	input := "key0=value\n"
	prefix := "${"
	postfix := "}"
	i := 0
	for i = 0; i < maxExpansionDepth-1; i++ {
		input += fmt.Sprintf("key%d=%skey%d%s\n", i+1, prefix, i, postfix)
	}

	p, err := Load([]byte(input), ISO_8859_1)
	assert.Equal(t, err, nil)
	p.Prefix = prefix
	p.Postfix = postfix

	assert.Equal(t, p.MustGet(fmt.Sprintf("key%d", i)), "value")

	// Nudge input over the edge
	input += fmt.Sprintf("key%d=%skey%d%s\n", i+1, prefix, i, postfix)

	_, err = Load([]byte(input), ISO_8859_1)
	assert.Equal(t, err != nil, true, "want error")
	assert.Equal(t, strings.Contains(err.Error(), "expansion too deep"), true)
}

func TestDisableExpansion(t *testing.T) {
	input := "key=value\nkey2=${key}"
	p := mustParse(t, input)
	p.DisableExpansion = true
	assert.Equal(t, p.MustGet("key"), "value")
	assert.Equal(t, p.MustGet("key2"), "${key}")

	// with expansion disabled we can introduce circular references
	p.MustSet("keyA", "${keyB}")
	p.MustSet("keyB", "${keyA}")
	assert.Equal(t, p.MustGet("keyA"), "${keyB}")
	assert.Equal(t, p.MustGet("keyB"), "${keyA}")
}

func TestDisableExpansionStillUpdatesKeys(t *testing.T) {
	p := NewProperties()
	p.MustSet("p1", "a")
	assert.Equal(t, p.Keys(), []string{"p1"})
	assert.Equal(t, p.String(), "p1 = a\n")

	p.DisableExpansion = true
	p.MustSet("p2", "b")

	assert.Equal(t, p.Keys(), []string{"p1", "p2"})
	assert.Equal(t, p.String(), "p1 = a\np2 = b\n")
}

func TestMustGet(t *testing.T) {
	input := "key = value\nkey2 = ghi"
	p := mustParse(t, input)
	assert.Equal(t, p.MustGet("key"), "value")
	assert.Panic(t, func() { p.MustGet("invalid") }, "unknown property: invalid")
}

func TestGetBool(t *testing.T) {
	for _, test := range boolTests {
		p := mustParse(t, test.input)
		assert.Equal(t, p.Len(), 1)
		assert.Equal(t, p.GetBool(test.key, test.def), test.value)
	}
}

func TestMustGetBool(t *testing.T) {
	input := "key = true\nkey2 = ghi"
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetBool("key"), true)
	assert.Panic(t, func() { p.MustGetBool("invalid") }, "unknown property: invalid")
}

func TestGetDuration(t *testing.T) {
	for _, test := range durationTests {
		p := mustParse(t, test.input)
		assert.Equal(t, p.Len(), 1)
		assert.Equal(t, p.GetDuration(test.key, test.def), test.value)
	}
}

func TestMustGetDuration(t *testing.T) {
	input := "key = 123\nkey2 = ghi"
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetDuration("key"), time.Duration(123))
	assert.Panic(t, func() { p.MustGetDuration("key2") }, "strconv.ParseInt: parsing.*")
	assert.Panic(t, func() { p.MustGetDuration("invalid") }, "unknown property: invalid")
}

func TestGetParsedDuration(t *testing.T) {
	for _, test := range parsedDurationTests {
		p := mustParse(t, test.input)
		assert.Equal(t, p.Len(), 1)
		assert.Equal(t, p.GetParsedDuration(test.key, test.def), test.value)
	}
}

func TestMustGetParsedDuration(t *testing.T) {
	input := "key = 123ms\nkey2 = ghi"
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetParsedDuration("key"), 123*time.Millisecond)
	assert.Panic(t, func() { p.MustGetParsedDuration("key2") }, "time: invalid duration ghi")
	assert.Panic(t, func() { p.MustGetParsedDuration("invalid") }, "unknown property: invalid")
}

func TestGetFloat64(t *testing.T) {
	for _, test := range floatTests {
		p := mustParse(t, test.input)
		assert.Equal(t, p.Len(), 1)
		assert.Equal(t, p.GetFloat64(test.key, test.def), test.value)
	}
}

func TestMustGetFloat64(t *testing.T) {
	input := "key = 123\nkey2 = ghi"
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetFloat64("key"), float64(123))
	assert.Panic(t, func() { p.MustGetFloat64("key2") }, "strconv.ParseFloat: parsing.*")
	assert.Panic(t, func() { p.MustGetFloat64("invalid") }, "unknown property: invalid")
}

func TestGetInt(t *testing.T) {
	for _, test := range int64Tests {
		p := mustParse(t, test.input)
		assert.Equal(t, p.Len(), 1)
		assert.Equal(t, p.GetInt(test.key, int(test.def)), int(test.value))
	}
}

func TestMustGetInt(t *testing.T) {
	input := "key = 123\nkey2 = ghi"
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetInt("key"), int(123))
	assert.Panic(t, func() { p.MustGetInt("key2") }, "strconv.ParseInt: parsing.*")
	assert.Panic(t, func() { p.MustGetInt("invalid") }, "unknown property: invalid")
}

func TestGetInt64(t *testing.T) {
	for _, test := range int64Tests {
		p := mustParse(t, test.input)
		assert.Equal(t, p.Len(), 1)
		assert.Equal(t, p.GetInt64(test.key, test.def), test.value)
	}
}

func TestMustGetInt64(t *testing.T) {
	input := "key = 123\nkey2 = ghi"
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetInt64("key"), int64(123))
	assert.Panic(t, func() { p.MustGetInt64("key2") }, "strconv.ParseInt: parsing.*")
	assert.Panic(t, func() { p.MustGetInt64("invalid") }, "unknown property: invalid")
}

func TestGetUint(t *testing.T) {
	for _, test := range uint64Tests {
		p := mustParse(t, test.input)
		assert.Equal(t, p.Len(), 1)
		assert.Equal(t, p.GetUint(test.key, uint(test.def)), uint(test.value))
	}
}

func TestMustGetUint(t *testing.T) {
	input := "key = 123\nkey2 = ghi"
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetUint("key"), uint(123))
	assert.Panic(t, func() { p.MustGetUint64("key2") }, "strconv.ParseUint: parsing.*")
	assert.Panic(t, func() { p.MustGetUint64("invalid") }, "unknown property: invalid")
}

func TestGetUint64(t *testing.T) {
	for _, test := range uint64Tests {
		p := mustParse(t, test.input)
		assert.Equal(t, p.Len(), 1)
		assert.Equal(t, p.GetUint64(test.key, test.def), test.value)
	}
}

func TestMustGetUint64(t *testing.T) {
	input := "key = 123\nkey2 = ghi"
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetUint64("key"), uint64(123))
	assert.Panic(t, func() { p.MustGetUint64("key2") }, "strconv.ParseUint: parsing.*")
	assert.Panic(t, func() { p.MustGetUint64("invalid") }, "unknown property: invalid")
}

func TestGetString(t *testing.T) {
	for _, test := range stringTests {
		p := mustParse(t, test.input)
		assert.Equal(t, p.Len(), 1)
		assert.Equal(t, p.GetString(test.key, test.def), test.value)
	}
}

func TestMustGetString(t *testing.T) {
	input := `key = value`
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetString("key"), "value")
	assert.Panic(t, func() { p.MustGetString("invalid") }, "unknown property: invalid")
}

func TestComment(t *testing.T) {
	for _, test := range commentTests {
		p := mustParse(t, test.input)
		assert.Equal(t, p.MustGetString(test.key), test.value)
		assert.Equal(t, p.GetComments(test.key), test.comments)
		if test.comments != nil {
			assert.Equal(t, p.GetComment(test.key), test.comments[len(test.comments)-1])
		} else {
			assert.Equal(t, p.GetComment(test.key), "")
		}

		// test setting comments
		if len(test.comments) > 0 {
			// set single comment
			p.ClearComments()
			assert.Equal(t, len(p.c), 0)
			p.SetComment(test.key, test.comments[0])
			assert.Equal(t, p.GetComment(test.key), test.comments[0])

			// set multiple comments
			p.ClearComments()
			assert.Equal(t, len(p.c), 0)
			p.SetComments(test.key, test.comments)
			assert.Equal(t, p.GetComments(test.key), test.comments)

			// clear comments for a key
			p.SetComments(test.key, nil)
			assert.Equal(t, p.GetComment(test.key), "")
			assert.Equal(t, p.GetComments(test.key), ([]string)(nil))
		}
	}
}

func TestFilter(t *testing.T) {
	for _, test := range filterTests {
		p := mustParse(t, test.input)
		pp, err := p.Filter(test.pattern)
		if err != nil {
			assert.Matches(t, err.Error(), test.err)
			continue
		}
		assert.Equal(t, pp != nil, true, "want properties")
		assert.Equal(t, pp.Len(), len(test.keys))
		for _, key := range test.keys {
			v1, ok1 := p.Get(key)
			v2, ok2 := pp.Get(key)
			assert.Equal(t, ok1, true)
			assert.Equal(t, ok2, true)
			assert.Equal(t, v1, v2)
		}
	}
}

func TestFilterPrefix(t *testing.T) {
	for _, test := range filterPrefixTests {
		p := mustParse(t, test.input)
		pp := p.FilterPrefix(test.prefix)
		assert.Equal(t, pp != nil, true, "want properties")
		assert.Equal(t, pp.Len(), len(test.keys))
		for _, key := range test.keys {
			v1, ok1 := p.Get(key)
			v2, ok2 := pp.Get(key)
			assert.Equal(t, ok1, true)
			assert.Equal(t, ok2, true)
			assert.Equal(t, v1, v2)
		}
	}
}

func TestFilterStripPrefix(t *testing.T) {
	for _, test := range filterStripPrefixTests {
		p := mustParse(t, test.input)
		pp := p.FilterPrefix(test.prefix)
		assert.Equal(t, pp != nil, true, "want properties")
		assert.Equal(t, pp.Len(), len(test.keys))
		for _, key := range test.keys {
			v1, ok1 := p.Get(key)
			v2, ok2 := pp.Get(key)
			assert.Equal(t, ok1, true)
			assert.Equal(t, ok2, true)
			assert.Equal(t, v1, v2)
		}
	}
}

func TestKeys(t *testing.T) {
	for _, test := range keysTests {
		p := mustParse(t, test.input)
		assert.Equal(t, p.Len(), len(test.keys))
		assert.Equal(t, len(p.Keys()), len(test.keys))
		assert.Equal(t, p.Keys(), test.keys)
	}
}

func TestSet(t *testing.T) {
	for _, test := range setTests {
		p := mustParse(t, test.input)
		prev, ok, err := p.Set(test.key, test.value)
		if test.err != "" {
			assert.Matches(t, err.Error(), test.err)
			continue
		}

		assert.Equal(t, err, nil)
		assert.Equal(t, ok, test.ok)
		if ok {
			assert.Equal(t, prev, test.prev)
		}
		assert.Equal(t, p.Keys(), test.keys)
	}
}

func TestSetValue(t *testing.T) {
	tests := []interface{}{
		true, false,
		int8(123), int16(123), int32(123), int64(123), int(123),
		uint8(123), uint16(123), uint32(123), uint64(123), uint(123),
		float32(1.23), float64(1.23),
		"abc",
	}

	for _, v := range tests {
		p := NewProperties()
		err := p.SetValue("x", v)
		assert.Equal(t, err, nil)
		assert.Equal(t, p.GetString("x", ""), fmt.Sprintf("%v", v))
	}
}

func TestMustSet(t *testing.T) {
	input := "key=${key}"
	p := mustParse(t, input)
	assert.Panic(t, func() { p.MustSet("key", "${key}") }, "circular reference .*")
}

func TestWrite(t *testing.T) {
	for _, test := range writeTests {
		p, err := parse(test.input)

		buf := new(bytes.Buffer)
		var n int
		switch test.encoding {
		case "UTF-8":
			n, err = p.Write(buf, UTF8)
		case "ISO-8859-1":
			n, err = p.Write(buf, ISO_8859_1)
		}
		assert.Equal(t, err, nil)
		s := string(buf.Bytes())
		assert.Equal(t, n, len(test.output), fmt.Sprintf("input=%q expected=%q obtained=%q", test.input, test.output, s))
		assert.Equal(t, s, test.output, fmt.Sprintf("input=%q expected=%q obtained=%q", test.input, test.output, s))
	}
}

func TestWriteComment(t *testing.T) {
	for _, test := range writeCommentTests {
		p, err := parse(test.input)

		buf := new(bytes.Buffer)
		var n int
		switch test.encoding {
		case "UTF-8":
			n, err = p.WriteComment(buf, "# ", UTF8)
		case "ISO-8859-1":
			n, err = p.WriteComment(buf, "# ", ISO_8859_1)
		}
		assert.Equal(t, err, nil)
		s := string(buf.Bytes())
		assert.Equal(t, n, len(test.output), fmt.Sprintf("input=%q expected=%q obtained=%q", test.input, test.output, s))
		assert.Equal(t, s, test.output, fmt.Sprintf("input=%q expected=%q obtained=%q", test.input, test.output, s))
	}
}

func TestCustomExpansionExpression(t *testing.T) {
	testKeyValuePrePostfix(t, "*[", "]*", "key=value\nkey2=*[key]*", "key", "value", "key2", "value")
}

func TestPanicOn32BitIntOverflow(t *testing.T) {
	is32Bit = true
	var min, max int64 = math.MinInt32 - 1, math.MaxInt32 + 1
	input := fmt.Sprintf("min=%d\nmax=%d", min, max)
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetInt64("min"), min)
	assert.Equal(t, p.MustGetInt64("max"), max)
	assert.Panic(t, func() { p.MustGetInt("min") }, ".* out of range")
	assert.Panic(t, func() { p.MustGetInt("max") }, ".* out of range")
}

func TestPanicOn32BitUintOverflow(t *testing.T) {
	is32Bit = true
	var max uint64 = math.MaxUint32 + 1
	input := fmt.Sprintf("max=%d", max)
	p := mustParse(t, input)
	assert.Equal(t, p.MustGetUint64("max"), max)
	assert.Panic(t, func() { p.MustGetUint("max") }, ".* out of range")
}

func TestDeleteKey(t *testing.T) {
	input := "#comments should also be gone\nkey=to-be-deleted\nsecond=key"
	p := mustParse(t, input)
	assert.Equal(t, len(p.m), 2)
	assert.Equal(t, len(p.c), 1)
	assert.Equal(t, len(p.k), 2)
	p.Delete("key")
	assert.Equal(t, len(p.m), 1)
	assert.Equal(t, len(p.c), 0)
	assert.Equal(t, len(p.k), 1)
	assert.Equal(t, p.k[0], "second")
	assert.Equal(t, p.m["second"], "key")
}

func TestDeleteUnknownKey(t *testing.T) {
	input := "#comments should also be gone\nkey=to-be-deleted"
	p := mustParse(t, input)
	assert.Equal(t, len(p.m), 1)
	assert.Equal(t, len(p.c), 1)
	assert.Equal(t, len(p.k), 1)
	p.Delete("wrong-key")
	assert.Equal(t, len(p.m), 1)
	assert.Equal(t, len(p.c), 1)
	assert.Equal(t, len(p.k), 1)
}

func TestMerge(t *testing.T) {
	input1 := "#comment\nkey=value\nkey2=value2"
	input2 := "#another comment\nkey=another value\nkey3=value3"
	p1 := mustParse(t, input1)
	p2 := mustParse(t, input2)
	p1.Merge(p2)
	assert.Equal(t, len(p1.m), 3)
	assert.Equal(t, len(p1.c), 1)
	assert.Equal(t, len(p1.k), 3)
	assert.Equal(t, p1.MustGet("key"), "another value")
	assert.Equal(t, p1.GetComment("key"), "another comment")
}

func TestMap(t *testing.T) {
	input := "key=value\nabc=def"
	p := mustParse(t, input)
	m := map[string]string{"key": "value", "abc": "def"}
	assert.Equal(t, p.Map(), m)
}

func TestFilterFunc(t *testing.T) {
	input := "key=value\nabc=def"
	p := mustParse(t, input)
	pp := p.FilterFunc(func(k, v string) bool {
		return k != "abc"
	})
	m := map[string]string{"key": "value"}
	assert.Equal(t, pp.Map(), m)
}

func TestLoad(t *testing.T) {
	x := "key=${value}\nvalue=${key}"
	p := NewProperties()
	p.DisableExpansion = true
	err := p.Load([]byte(x), UTF8)
	assert.Equal(t, err, nil)
}

// ----------------------------------------------------------------------------

// tests all combinations of delimiters, leading and/or trailing whitespace and newlines.
func testWhitespaceAndDelimiterCombinations(t *testing.T, key, value string) {
	whitespace := []string{"", " ", "\f", "\t"}
	delimiters := []string{"", " ", "=", ":"}
	newlines := []string{"", "\r", "\n", "\r\n"}
	for _, dl := range delimiters {
		for _, ws1 := range whitespace {
			for _, ws2 := range whitespace {
				for _, nl := range newlines {
					// skip the one case where there is nothing between a key and a value
					if ws1 == "" && dl == "" && ws2 == "" && value != "" {
						continue
					}

					input := fmt.Sprintf("%s%s%s%s%s%s", key, ws1, dl, ws2, value, nl)
					testKeyValue(t, input, key, value)
				}
			}
		}
	}
}

// tests whether key/value pairs exist for a given input.
// keyvalues is expected to be an even number of strings of "key", "value", ...
func testKeyValue(t *testing.T, input string, keyvalues ...string) {
	testKeyValuePrePostfix(t, "${", "}", input, keyvalues...)
}

// tests whether key/value pairs exist for a given input.
// keyvalues is expected to be an even number of strings of "key", "value", ...
func testKeyValuePrePostfix(t *testing.T, prefix, postfix, input string, keyvalues ...string) {
	p, err := Load([]byte(input), ISO_8859_1)
	assert.Equal(t, err, nil)
	p.Prefix = prefix
	p.Postfix = postfix
	assertKeyValues(t, input, p, keyvalues...)
}

// tests whether key/value pairs exist for a given input.
// keyvalues is expected to be an even number of strings of "key", "value", ...
func assertKeyValues(t *testing.T, input string, p *Properties, keyvalues ...string) {
	assert.Equal(t, p != nil, true, "want properties")
	assert.Equal(t, 2*p.Len(), len(keyvalues), "Odd number of key/value pairs.")

	for i := 0; i < len(keyvalues); i += 2 {
		key, value := keyvalues[i], keyvalues[i+1]
		v, ok := p.Get(key)
		if !ok {
			t.Errorf("No key %q found (input=%q)", key, input)
		}
		if got, want := v, value; !reflect.DeepEqual(got, want) {
			t.Errorf("Value %q does not match %q (input=%q)", v, value, input)
		}
	}
}

func mustParse(t *testing.T, s string) *Properties {
	p, err := parse(s)
	if err != nil {
		t.Fatalf("parse failed with %s", err)
	}
	return p
}

// prints to stderr if the -verbose flag was given.
func printf(format string, args ...interface{}) {
	if *verbose {
		fmt.Fprintf(os.Stderr, format, args...)
	}
}
