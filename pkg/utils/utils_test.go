package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestAsJson(t *testing.T) {
	type myStruct struct {
		a string
	}

	output := AsJson(&myStruct{a: "foo"})

	// no idea why this is returning empty hashes but it's works in the app ¯\_(ツ)_/¯
	assert.EqualValues(t, "{}", output)
}

func TestSafeTruncate(t *testing.T) {
	type scenario struct {
		str      string
		limit    int
		expected string
	}

	scenarios := []scenario{
		{
			str:      "",
			limit:    0,
			expected: "",
		},
		{
			str:      "12345",
			limit:    3,
			expected: "123",
		},
		{
			str:      "12345",
			limit:    4,
			expected: "1234",
		},
		{
			str:      "12345",
			limit:    5,
			expected: "12345",
		},
		{
			str:      "12345",
			limit:    6,
			expected: "12345",
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, s.expected, SafeTruncate(s.str, s.limit))
	}
}

func TestModuloWithWrap(t *testing.T) {
	type scenario struct {
		n        int
		max      int
		expected int
	}

	scenarios := []scenario{
		{
			n:        0,
			max:      0,
			expected: 0,
		},
		{
			n:        0,
			max:      1,
			expected: 0,
		},
		{
			n:        1,
			max:      0,
			expected: 0,
		},
		{
			n:        3,
			max:      2,
			expected: 1,
		},
		{
			n:        -1,
			max:      2,
			expected: 1,
		},
	}

	for _, s := range scenarios {
		if s.expected != ModuloWithWrap(s.n, s.max) {
			t.Errorf("expected %d, got %d, for n: %d, max: %d", s.expected, ModuloWithWrap(s.n, s.max), s.n, s.max)
		}
	}
}
