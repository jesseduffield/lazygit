package utils

import (
	"strings"
	"testing"

	"github.com/mattn/go-runewidth"
	"github.com/stretchr/testify/assert"
)

func TestWithPadding(t *testing.T) {
	type scenario struct {
		str       string
		padding   int
		alignment Alignment
		expected  string
	}

	scenarios := []scenario{
		{
			str:       "hello world !",
			padding:   1,
			alignment: AlignLeft,
			expected:  "hello world !",
		},
		{
			str:       "hello world !",
			padding:   14,
			alignment: AlignLeft,
			expected:  "hello world ! ",
		},
		{
			str:       "hello world !",
			padding:   14,
			alignment: AlignRight,
			expected:  " hello world !",
		},
		{
			str:       "Güçlü",
			padding:   7,
			alignment: AlignLeft,
			expected:  "Güçlü  ",
		},
		{
			str:       "Güçlü",
			padding:   7,
			alignment: AlignRight,
			expected:  "  Güçlü",
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, s.expected, WithPadding(s.str, s.padding, s.alignment))
	}
}

func TestGetPadWidths(t *testing.T) {
	type scenario struct {
		input    [][]string
		expected []int
	}

	tests := []scenario{
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
		{
			[][]string{{"AŁ", "b", "ccc"}, {"c", "d", "e"}},
			[]int{2, 1},
		},
	}

	for _, test := range tests {
		output := getPadWidths(test.input)
		assert.EqualValues(t, test.expected, output)
	}
}

func TestTruncateWithEllipsis(t *testing.T) {
	// will need to check chinese characters as well
	// important that we have a three dot ellipsis within the limit
	type scenario struct {
		str      string
		limit    int
		expected string
	}

	scenarios := []scenario{
		{
			"hello world !",
			1,
			".",
		},
		{
			"hello world !",
			2,
			"..",
		},
		{
			"hello world !",
			3,
			"he…",
		},
		{
			"hello world !",
			4,
			"hel…",
		},
		{
			"hello world !",
			5,
			"hell…",
		},
		{
			"hello world !",
			12,
			"hello world…",
		},
		{
			"hello world !",
			13,
			"hello world !",
		},
		{
			"hello world !",
			14,
			"hello world !",
		},
		{
			"大大大大",
			5,
			"大大…",
		},
		{
			"大大大大",
			2,
			"..",
		},
		{
			"大大大大",
			1,
			".",
		},
		{
			"大大大大",
			0,
			"",
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, s.expected, TruncateWithEllipsis(s.str, s.limit))
	}
}

func TestRenderDisplayStrings(t *testing.T) {
	type scenario struct {
		input                   [][]string
		columnAlignments        []Alignment
		expectedOutput          string
		expectedColumnPositions []int
	}

	tests := []scenario{
		{
			input:                   [][]string{{""}, {""}},
			columnAlignments:        nil,
			expectedOutput:          "",
			expectedColumnPositions: []int{0, 0},
		},
		{
			input:                   [][]string{{"a"}, {""}},
			columnAlignments:        nil,
			expectedOutput:          "a\n",
			expectedColumnPositions: []int{0},
		},
		{
			input:                   [][]string{{"a"}, {"b"}},
			columnAlignments:        nil,
			expectedOutput:          "a\nb",
			expectedColumnPositions: []int{0},
		},
		{
			input:                   [][]string{{"a", "b"}, {"c", "d"}},
			columnAlignments:        nil,
			expectedOutput:          "a b\nc d",
			expectedColumnPositions: []int{0, 2},
		},
		{
			input:                   [][]string{{"a", "", "c"}, {"d", "", "f"}},
			columnAlignments:        nil,
			expectedOutput:          "a c\nd f",
			expectedColumnPositions: []int{0, 2, 2},
		},
		{
			input:                   [][]string{{"a", "", "c", ""}, {"d", "", "f", ""}},
			columnAlignments:        nil,
			expectedOutput:          "a c\nd f",
			expectedColumnPositions: []int{0, 2, 2},
		},
		{
			input:                   [][]string{{"abc", "", "d", ""}, {"e", "", "f", ""}},
			columnAlignments:        nil,
			expectedOutput:          "abc d\ne   f",
			expectedColumnPositions: []int{0, 4, 4},
		},
		{
			input:                   [][]string{{"", "abc", "", "", "d", "e"}, {"", "f", "", "", "g", "h"}},
			columnAlignments:        nil,
			expectedOutput:          "abc d e\nf   g h",
			expectedColumnPositions: []int{0, 0, 4, 4, 4, 6},
		},
		{
			input:                   [][]string{{"abc", "", "d", ""}, {"e", "", "f", ""}},
			columnAlignments:        []Alignment{AlignLeft, AlignLeft}, // same as nil (default)
			expectedOutput:          "abc d\ne   f",
			expectedColumnPositions: []int{0, 4, 4},
		},
		{
			input:                   [][]string{{"abc", "", "d", ""}, {"e", "", "f", ""}},
			columnAlignments:        []Alignment{AlignRight, AlignLeft},
			expectedOutput:          "abc d\n  e f",
			expectedColumnPositions: []int{0, 4, 4},
		},
		{
			input:                   [][]string{{"a", "", "bcd", "efg", "h"}, {"i", "", "j", "k", "l"}},
			columnAlignments:        []Alignment{AlignLeft, AlignLeft, AlignRight, AlignLeft},
			expectedOutput:          "a bcd efg h\ni   j k   l",
			expectedColumnPositions: []int{0, 2, 2, 6, 10},
		},
		{
			input:                   [][]string{{"abc", "", "d", ""}, {"e", "", "f", ""}},
			columnAlignments:        []Alignment{AlignRight}, // gracefully defaults unspecified columns to left-align
			expectedOutput:          "abc d\n  e f",
			expectedColumnPositions: []int{0, 4, 4},
		},
	}

	for _, test := range tests {
		output, columnPositions := RenderDisplayStrings(test.input, test.columnAlignments)
		assert.EqualValues(t, test.expectedOutput, strings.Join(output, "\n"))
		assert.EqualValues(t, test.expectedColumnPositions, columnPositions)
	}
}

func BenchmarkStringWidthAsciiOriginal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runewidth.StringWidth("some ASCII string")
	}
}

func BenchmarkStringWidthAsciiOptimized(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StringWidth("some ASCII string")
	}
}

func BenchmarkStringWidthNonAsciiOriginal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runewidth.StringWidth("some non-ASCII string 🍉")
	}
}

func BenchmarkStringWidthNonAsciiOptimized(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StringWidth("some non-ASCII string 🍉")
	}
}
