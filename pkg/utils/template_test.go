package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		{
			"{{a}}",
			map[string]string{
				"a": "X{{.a}}X",
			},
			"X{{.a}}X",
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, s.expected, ResolvePlaceholderString(s.templateString, s.arguments))
	}
}
