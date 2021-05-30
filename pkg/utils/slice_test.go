package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			"two elements, giving second one",
			[]int{1, 2},
			2,
			1,
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
			0,
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
