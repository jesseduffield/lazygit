package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		{
			"Güçlü",
			7,
			"Güçlü  ",
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, s.expected, WithPadding(s.str, s.padding))
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
		if !assert.EqualValues(t, output, test.expected) {
			t.Errorf("getPadWidths(%v) = %v, want %v", test.input, output, test.expected)
		}
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
		input    [][]string
		expected string
	}

	tests := []scenario{
		{
			[][]string{{""}, {""}},
			"",
		},
		{
			[][]string{{"a"}, {""}},
			"a\n",
		},
		{
			[][]string{{"a"}, {"b"}},
			"a\nb",
		},
		{
			[][]string{{"a", "b"}, {"c", "d"}},
			"a b\nc d",
		},
		{
			[][]string{{"a", "", "c"}, {"d", "", "f"}},
			"a c\nd f",
		},
		{
			[][]string{{"a", "", "c", ""}, {"d", "", "f", ""}},
			"a c\nd f",
		},
	}

	for _, test := range tests {
		output := RenderDisplayStrings(test.input)
		if !assert.EqualValues(t, output, test.expected) {
			t.Errorf("RenderDisplayStrings(%v) = %v, want %v", test.input, output, test.expected)
		}
	}
}
