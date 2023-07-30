package git_commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitCommandBuilder(t *testing.T) {
	scenarios := []struct {
		input    []string
		expected []string
	}{
		{
			input: NewGitCmd("push").
				Arg("--force-with-lease").
				Arg("--set-upstream").
				Arg("origin").
				Arg("master").
				ToArgv(),
			expected: []string{"git", "push", "--force-with-lease", "--set-upstream", "origin", "master"},
		},
		{
			input:    NewGitCmd("push").ArgIf(true, "--test").ToArgv(),
			expected: []string{"git", "push", "--test"},
		},
		{
			input:    NewGitCmd("push").ArgIf(false, "--test").ToArgv(),
			expected: []string{"git", "push"},
		},
		{
			input:    NewGitCmd("push").ArgIfElse(true, "-b", "-a").ToArgv(),
			expected: []string{"git", "push", "-b"},
		},
		{
			input:    NewGitCmd("push").ArgIfElse(false, "-a", "-b").ToArgv(),
			expected: []string{"git", "push", "-b"},
		},
		{
			input:    NewGitCmd("push").Arg("-a", "-b").ToArgv(),
			expected: []string{"git", "push", "-a", "-b"},
		},
		{
			input:    NewGitCmd("push").Config("user.name=foo").Config("user.email=bar").ToArgv(),
			expected: []string{"git", "-c", "user.email=bar", "-c", "user.name=foo", "push"},
		},
		{
			input:    NewGitCmd("push").Dir("a/b/c").ToArgv(),
			expected: []string{"git", "-C", "a/b/c", "push"},
		},
	}

	for _, s := range scenarios {
		assert.Equal(t, s.input, s.expected)
	}
}
