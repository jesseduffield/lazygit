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
