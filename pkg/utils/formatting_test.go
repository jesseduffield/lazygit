package utils

import (
	"testing"

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
		assert.EqualValues(t, output, test.expected)
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
			"...",
		},
		{
			"hello world !",
			4,
			"h...",
		},
		{
			"hello world !",
			5,
			"he...",
		},
		{
			"hello world !",
			12,
			"hello wor...",
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
			"大...",
		},
		{
			"大大大大",
			2,
			"..",
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
		input            [][]string
		columnAlignments []Alignment
		expected         string
	}

	tests := []scenario{
		{
			input:            [][]string{{""}, {""}},
			columnAlignments: nil,
			expected:         "",
		},
		{
			input:            [][]string{{"a"}, {""}},
			columnAlignments: nil,
			expected:         "a\n",
		},
		{
			input:            [][]string{{"a"}, {"b"}},
			columnAlignments: nil,
			expected:         "a\nb",
		},
		{
			input:            [][]string{{"a", "b"}, {"c", "d"}},
			columnAlignments: nil,
			expected:         "a b\nc d",
		},
		{
			input:            [][]string{{"a", "", "c"}, {"d", "", "f"}},
			columnAlignments: nil,
			expected:         "a c\nd f",
		},
		{
			input:            [][]string{{"a", "", "c", ""}, {"d", "", "f", ""}},
			columnAlignments: nil,
			expected:         "a c\nd f",
		},
		{
			input:            [][]string{{"abc", "", "d", ""}, {"e", "", "f", ""}},
			columnAlignments: nil,
			expected:         "abc d\ne   f",
		},
		{
			input:            [][]string{{"abc", "", "d", ""}, {"e", "", "f", ""}},
			columnAlignments: []Alignment{AlignLeft, AlignLeft}, // same as nil (default)
			expected:         "abc d\ne   f",
		},
		{
			input:            [][]string{{"abc", "", "d", ""}, {"e", "", "f", ""}},
			columnAlignments: []Alignment{AlignRight, AlignLeft},
			expected:         "abc d\n  e f",
		},
		{
			input:            [][]string{{"abc", "", "d", ""}, {"e", "", "f", ""}},
			columnAlignments: []Alignment{AlignRight}, // gracefully defaults unspecified columns to left-align
			expected:         "abc d\n  e f",
		},
	}

	for _, test := range tests {
		output := RenderDisplayStrings(test.input, test.columnAlignments)
		assert.EqualValues(t, output, test.expected)
	}
}
