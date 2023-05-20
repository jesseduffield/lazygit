package git_commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitCommandBuilder(t *testing.T) {
	scenarios := []struct {
		input    string
		expected string
	}{
		{
			input: NewGitCmd("push").
				Arg("--force-with-lease").
				Arg("--set-upstream").
				Arg("origin").
				Arg("master").
				ToString(),
			expected: "git push --force-with-lease --set-upstream origin master",
		},
		{
			input:    NewGitCmd("push").ArgIf(true, "--test").ToString(),
			expected: "git push --test",
		},
		{
			input:    NewGitCmd("push").ArgIf(false, "--test").ToString(),
			expected: "git push",
		},
		{
			input:    NewGitCmd("push").ArgIfElse(true, "-b", "-a").ToString(),
			expected: "git push -b",
		},
		{
			input:    NewGitCmd("push").ArgIfElse(false, "-a", "-b").ToString(),
			expected: "git push -b",
		},
		{
			input:    NewGitCmd("push").Arg("-a", "-b").ToString(),
			expected: "git push -a -b",
		},
		{
			input:    NewGitCmd("push").Config("user.name=foo").Config("user.email=bar").ToString(),
			expected: "git -c user.email=bar -c user.name=foo push",
		},
		{
			input:    NewGitCmd("push").RepoPath("a/b/c").ToString(),
			expected: "git -C a/b/c push",
		},
	}

	for _, s := range scenarios {
		assert.Equal(t, s.input, s.expected)
	}
}
