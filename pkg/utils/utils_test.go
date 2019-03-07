package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSplitLines is a function.
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

// TestWithPadding is a function.
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

// TestTrimTrailingNewline is a function.
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

// TestNormalizeLinefeeds is a function.
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

// TestResolvePlaceholderString is a function.
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

// TestDisplayArraysAligned is a function.
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

// GetDisplayStrings is a function.
func (d *myDisplayable) GetDisplayStrings(isFocused bool) []string {
	if isFocused {
		return append(d.strings, "blah")
	}
	return d.strings
}

// TestGetDisplayStringArrays is a function.
func TestGetDisplayStringArrays(t *testing.T) {
	type scenario struct {
		input     []Displayable
		isFocused bool
		expected  [][]string
	}

	scenarios := []scenario{
		{
			[]Displayable{
				Displayable(&myDisplayable{[]string{"a", "b"}}),
				Displayable(&myDisplayable{[]string{"c", "d"}}),
			},
			false,
			[][]string{{"a", "b"}, {"c", "d"}},
		},
		{
			[]Displayable{
				Displayable(&myDisplayable{[]string{"a", "b"}}),
				Displayable(&myDisplayable{[]string{"c", "d"}}),
			},
			true,
			[][]string{{"a", "b", "blah"}, {"c", "d", "blah"}},
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, s.expected, getDisplayStringArrays(s.input, s.isFocused))
	}
}

// TestRenderDisplayableList is a function.
func TestRenderDisplayableList(t *testing.T) {
	type scenario struct {
		input                []Displayable
		isFocused            bool
		expectedString       string
		expectedErrorMessage string
	}

	scenarios := []scenario{
		{
			[]Displayable{
				Displayable(&myDisplayable{[]string{}}),
				Displayable(&myDisplayable{[]string{}}),
			},
			false,
			"\n",
			"",
		},
		{
			[]Displayable{
				Displayable(&myDisplayable{[]string{"aa", "b"}}),
				Displayable(&myDisplayable{[]string{"c", "d"}}),
			},
			false,
			"aa b\nc  d",
			"",
		},
		{
			[]Displayable{
				Displayable(&myDisplayable{[]string{"a"}}),
				Displayable(&myDisplayable{[]string{"b", "c"}}),
			},
			false,
			"",
			"Each item must return the same number of strings to display",
		},
		{
			[]Displayable{
				Displayable(&myDisplayable{[]string{"a"}}),
				Displayable(&myDisplayable{[]string{"b"}}),
			},
			true,
			"a blah\nb blah",
			"",
		},
	}

	for _, s := range scenarios {
		str, err := renderDisplayableList(s.input, s.isFocused)
		assert.EqualValues(t, s.expectedString, str)
		if s.expectedErrorMessage != "" {
			assert.EqualError(t, err, s.expectedErrorMessage)
		} else {
			assert.NoError(t, err)
		}
	}
}

// TestRenderList is a function.
func TestRenderList(t *testing.T) {
	type scenario struct {
		input                interface{}
		isFocused            bool
		expectedString       string
		expectedErrorMessage string
	}

	scenarios := []scenario{
		{
			[]*myDisplayable{
				{[]string{"aa", "b"}},
				{[]string{"c", "d"}},
			},
			false,
			"aa b\nc  d",
			"",
		},
		{
			[]*myStruct{
				{},
				{},
			},
			false,
			"",
			"item does not implement the Displayable interface",
		},
		{
			&myStruct{},
			false,
			"",
			"RenderList given a non-slice type",
		},
		{
			[]*myDisplayable{
				{[]string{"a"}},
			},
			true,
			"a blah",
			"",
		},
	}

	for _, s := range scenarios {
		str, err := RenderList(s.input, s.isFocused)
		assert.EqualValues(t, s.expectedString, str)
		if s.expectedErrorMessage != "" {
			assert.EqualError(t, err, s.expectedErrorMessage)
		} else {
			assert.NoError(t, err)
		}
	}
}

// TestGetPaddedDisplayStrings is a function.
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

// TestGetPadWidths is a function.
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

// TestMin is a function.
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

// TestIncludesString is a function.
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

func TestNextIndex(t *testing.T) {
	type scenario struct {
		testName string
		list     []int
		element  int
		expected int
	}

	scenarios := []scenario{
		{
			// I'm not really fussed about how it behaves here
			"no elements",
			[]int{},
			1,
			0,
		},
		{
			"one element",
			[]int{1},
			1,
			0,
		},
		{
			"two elements",
			[]int{1, 2},
			1,
			1,
		},
		{
			"two elements, giving second one",
			[]int{1, 2},
			2,
			0,
		},
		{
			"three elements, giving second one",
			[]int{1, 2, 3},
			2,
			2,
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			assert.EqualValues(t, s.expected, NextIndex(s.list, s.element))
		})
	}
}

func TestPrevIndex(t *testing.T) {
	type scenario struct {
		testName string
		list     []int
		element  int
		expected int
	}

	scenarios := []scenario{
		{
			// I'm not really fussed about how it behaves here
			"no elements",
			[]int{},
			1,
			-1,
		},
		{
			"one element",
			[]int{1},
			1,
			0,
		},
		{
			"two elements",
			[]int{1, 2},
			1,
			1,
		},
		{
			"three elements, giving second one",
			[]int{1, 2, 3},
			2,
			0,
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			assert.EqualValues(t, s.expected, PrevIndex(s.list, s.element))
		})
	}
}

func TestAsJson(t *testing.T) {
	type myStruct struct {
		a string
	}

	output := AsJson(&myStruct{a: "foo"})

	// no idea why this is returning empty hashes but it's works in the app ¯\_(ツ)_/¯
	assert.EqualValues(t, "{}", output)
}
