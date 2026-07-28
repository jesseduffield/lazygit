package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestExpandTilde(t *testing.T) {
	home, err := os.UserHomeDir()
	assert.NoError(t, err)

	scenarios := []struct {
		name     string
		path     string
		expected string
	}{
		{"bare tilde", "~", home},
		{"tilde with subpath", "~/worktrees", filepath.Join(home, "worktrees")},
		{"absolute path is untouched", "/absolute/path", "/absolute/path"},
		{"relative path is untouched", "relative/path", "relative/path"},
		{"tilde not at the start is untouched", "/foo/~/bar", "/foo/~/bar"},
		{"tilde followed by a username is untouched", "~other/worktrees", "~other/worktrees"},
		{"empty string is untouched", "", ""},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.Equal(t, s.expected, ExpandTilde(s.path))
		})
	}
}
