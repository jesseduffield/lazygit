package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSplitLines is a function
func TestSplitLines(t *testing.T) {
	type scenario struct {
		multilineString string
		expected        []string
	}

	scenarios := []scenario{
		{
			"",
			[]string{},
		},
		{
			"\n",
			[]string{},
		},
		{
			"hello world !\nhello universe !\n",
			[]string{
				"hello world !",
				"hello universe !",
			},
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, s.expected, SplitLines(s.multilineString))
	}
}

// TestWithPadding is a function
func TestWithPadding(t *testing.T) {
	type scenario struct {
		str      string
		padding  int
		expected string
	}

	scenarios := []scenario{
		{
			"hello world !",
			1,
			"hello world !",
		},
		{
			"hello world !",
			14,
			"hello world ! ",
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, s.expected, WithPadding(s.str, s.padding))
	}
}

// TestTrimTrailingNewline is a function
func TestTrimTrailingNewline(t *testing.T) {
	type scenario struct {
		str      string
		expected string
	}

	scenarios := []scenario{
		{
			"hello world !\n",
			"hello world !",
		},
		{
			"hello world !",
			"hello world !",
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, s.expected, TrimTrailingNewline(s.str))
	}
}

// TestNormalizeLinefeeds is a function
func TestNormalizeLinefeeds(t *testing.T) {
	type scenario struct {
		byteArray []byte
		expected  []byte
	}
	var scenarios = []scenario{
		{
			// \r\n
			[]byte{97, 115, 100, 102, 13, 10},
			[]byte{97, 115, 100, 102, 10},
		},
		{
			// bash\r\nblah
			[]byte{97, 115, 100, 102, 13, 10, 97, 115, 100, 102},
			[]byte{97, 115, 100, 102, 10, 97, 115, 100, 102},
		},
		{
			// \r
			[]byte{97, 115, 100, 102, 13},
			[]byte{97, 115, 100, 102},
		},
		{
			// \n
			[]byte{97, 115, 100, 102, 10},
			[]byte{97, 115, 100, 102, 10},
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, string(s.expected), NormalizeLinefeeds(string(s.byteArray)))
	}
}

// TestResolvePlaceholderString is a function
func TestResolvePlaceholderString(t *testing.T) {
	type scenario struct {
		templateString string
		arguments      map[string]string
		expected       string
	}

	scenarios := []scenario{
		{
			"",
			map[string]string{},
			"",
		},
		{
			"hello",
			map[string]string{},
			"hello",
		},
		{
			"hello {{arg}}",
			map[string]string{},
			"hello {{arg}}",
		},
		{
			"hello {{arg}}",
			map[string]string{"arg": "there"},
			"hello there",
		},
		{
			"hello",
			map[string]string{"arg": "there"},
			"hello",
		},
		{
			"{{nothing}}",
			map[string]string{"nothing": ""},
			"",
		},
		{
			"{{}} {{ this }} { should not throw}} an {{{{}}}} error",
			map[string]string{
				"blah": "blah",
				"this": "won't match",
			},
			"{{}} {{ this }} { should not throw}} an {{{{}}}} error",
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, string(s.expected), ResolvePlaceholderString(s.templateString, s.arguments))
	}
}

// TestDisplayArraysAligned is a function
func TestDisplayArraysAligned(t *testing.T) {
	type scenario struct {
		input    [][]string
		expected bool
	}

	scenarios := []scenario{
		{
			[][]string{{"", ""}, {"", ""}},
			true,
		},
		{
			[][]string{{""}, {"", ""}},
			false,
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, s.expected, displayArraysAligned(s.input))
	}
}

type myDisplayable struct {
	strings []string
}

type myStruct struct{}

// GetDisplayStrings is a function
func (d *myDisplayable) GetDisplayStrings() []string {
	return d.strings
}

// TestGetDisplayStringArrays is a function
func TestGetDisplayStringArrays(t *testing.T) {
	type scenario struct {
		input    []Displayable
		expected [][]string
	}

	scenarios := []scenario{
		{
			[]Displayable{
				Displayable(&myDisplayable{[]string{"a", "b"}}),
				Displayable(&myDisplayable{[]string{"c", "d"}}),
			},
			[][]string{{"a", "b"}, {"c", "d"}},
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, s.expected, getDisplayStringArrays(s.input))
	}
}

// TestRenderDisplayableList is a function
func TestRenderDisplayableList(t *testing.T) {
	type scenario struct {
		input          []Displayable
		expectedString string
		expectedError  error
	}

	scenarios := []scenario{
		{
			[]Displayable{
				Displayable(&myDisplayable{[]string{}}),
				Displayable(&myDisplayable{[]string{}}),
			},
			"\n",
			nil,
		},
		{
			[]Displayable{
				Displayable(&myDisplayable{[]string{"aa", "b"}}),
				Displayable(&myDisplayable{[]string{"c", "d"}}),
			},
			"aa b\nc  d",
			nil,
		},
		{
			[]Displayable{
				Displayable(&myDisplayable{[]string{"a"}}),
				Displayable(&myDisplayable{[]string{"b", "c"}}),
			},
			"",
			errors.New("Each item must return the same number of strings to display"),
		},
	}

	for _, s := range scenarios {
		str, err := renderDisplayableList(s.input)
		assert.EqualValues(t, s.expectedString, str)
		assert.EqualValues(t, s.expectedError, err)
	}
}

// TestRenderList is a function
func TestRenderList(t *testing.T) {
	type scenario struct {
		input          interface{}
		expectedString string
		expectedError  error
	}

	scenarios := []scenario{
		{
			[]*myDisplayable{
				{[]string{"aa", "b"}},
				{[]string{"c", "d"}},
			},
			"aa b\nc  d",
			nil,
		},
		{
			[]*myStruct{
				{},
				{},
			},
			"",
			errors.New("item does not implement the Displayable interface"),
		},
		{
			&myStruct{},
			"",
			errors.New("RenderList given a non-slice type"),
		},
	}

	for _, s := range scenarios {
		str, err := RenderList(s.input)
		assert.EqualValues(t, s.expectedString, str)
		assert.EqualValues(t, s.expectedError, err)
	}
}

// TestGetPaddedDisplayStrings is a function
func TestGetPaddedDisplayStrings(t *testing.T) {
	type scenario struct {
		stringArrays [][]string
		padWidths    []int
		expected     []string
	}

	scenarios := []scenario{
		{
			[][]string{{"a", "b"}, {"c", "d"}},
			[]int{1},
			[]string{"a b", "c d"},
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, s.expected, getPaddedDisplayStrings(s.stringArrays, s.padWidths))
	}
}

// TestGetPadWidths is a function
func TestGetPadWidths(t *testing.T) {
	type scenario struct {
		stringArrays [][]string
		expected     []int
	}

	scenarios := []scenario{
		{
			[][]string{{""}, {""}},
			[]int{},
		},
		{
			[][]string{{"a"}, {""}},
			[]int{},
		},
		{
			[][]string{{"aa", "b", "ccc"}, {"c", "d", "e"}},
			[]int{2, 1},
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, s.expected, getPadWidths(s.stringArrays))
	}
}

// TestMin is a function
func TestMin(t *testing.T) {
	type scenario struct {
		a        int
		b        int
		expected int
	}

	scenarios := []scenario{
		{
			1,
			1,
			1,
		},
		{
			1,
			2,
			1,
		},
		{
			2,
			1,
			1,
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, s.expected, Min(s.a, s.b))
	}
}

// TestIncludesString is a function
func TestIncludesString(t *testing.T) {
	type scenario struct {
		list     []string
		element  string
		expected bool
	}

	scenarios := []scenario{
		{
			[]string{"a", "b"},
			"a",
			true,
		},
		{
			[]string{"a", "b"},
			"c",
			false,
		},
		{
			[]string{"a", "b"},
			"",
			false,
		},
		{
			[]string{""},
			"",
			true,
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, s.expected, IncludesString(s.list, s.element))
	}
}
