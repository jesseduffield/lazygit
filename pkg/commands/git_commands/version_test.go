package git_commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGitVersion(t *testing.T) {
	scenarios := []struct {
		input    string
		expected GitVersion
	}{
		{
			input:    "git version 2.39.0",
			expected: GitVersion{Major: 2, Minor: 39, Patch: 0, Additional: ""},
		},
		{
			input:    "git version 2.37.1 (Apple Git-137.1)",
			expected: GitVersion{Major: 2, Minor: 37, Patch: 1, Additional: "(Apple Git-137.1)"},
		},
		{
			input:    "git version 2.37 (Apple Git-137.1)",
			expected: GitVersion{Major: 2, Minor: 37, Patch: 0, Additional: "(Apple Git-137.1)"},
		},
	}

	for _, s := range scenarios {
		actual, err := ParseGitVersion(s.input)

		assert.NoError(t, err)
		assert.NotNil(t, actual)
		assert.Equal(t, s.expected.Major, actual.Major)
		assert.Equal(t, s.expected.Minor, actual.Minor)
		assert.Equal(t, s.expected.Patch, actual.Patch)
		assert.Equal(t, s.expected.Additional, actual.Additional)
	}
}

func TestGitVersionIsOlderThan(t *testing.T) {
	assert.False(t, (&GitVersion{2, 0, 0, ""}).IsOlderThan(1, 99, 99))
	assert.False(t, (&GitVersion{2, 0, 0, ""}).IsOlderThan(2, 0, 0))
	assert.False(t, (&GitVersion{2, 1, 0, ""}).IsOlderThan(2, 0, 9))

	assert.True(t, (&GitVersion{2, 0, 1, ""}).IsOlderThan(2, 1, 0))
	assert.True(t, (&GitVersion{2, 0, 1, ""}).IsOlderThan(3, 0, 0))
}

func TestGitVersionIsAtLeast(t *testing.T) {
	assert.True(t, (&GitVersion{2, 0, 0, ""}).IsAtLeast(1, 99, 99))
	assert.True(t, (&GitVersion{2, 0, 0, ""}).IsAtLeast(2, 0, 0))
	assert.True(t, (&GitVersion{2, 1, 0, ""}).IsAtLeast(2, 0, 9))

	assert.False(t, (&GitVersion{2, 0, 1, ""}).IsAtLeast(2, 1, 0))
	assert.False(t, (&GitVersion{2, 0, 1, ""}).IsAtLeast(3, 0, 0))
}
