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
