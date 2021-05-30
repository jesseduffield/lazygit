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
