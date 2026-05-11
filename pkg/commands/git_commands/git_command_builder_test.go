package git_commands

import (
	"strings"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
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

func TestRunGitCmdOnPaths(t *testing.T) {
	// Each path is 9000 bytes. Three fit within the 30 KB limit (27001 bytes
	// including spaces), four do not (36002 bytes), so a four-path slice must
	// be split into two calls of three and one.
	longPath := func(ch string) string { return strings.Repeat(ch, 9_000) }
	p1, p2, p3, p4 := longPath("a"), longPath("b"), longPath("c"), longPath("d")

	scenarios := []struct {
		name   string
		paths  []string
		runner *oscommands.FakeCmdObjRunner
	}{
		{
			name:   "empty list makes no calls",
			paths:  []string{},
			runner: oscommands.NewFakeRunner(t),
		},
		{
			name:  "paths that fit in one batch make a single call",
			paths: []string{p1, p2, p3},
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs(append([]string{"checkout", "--"}, p1, p2, p3), "", nil),
		},
		{
			name:  "paths that exceed the limit are split across multiple calls",
			paths: []string{p1, p2, p3, p4},
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs(append([]string{"checkout", "--"}, p1, p2, p3), "", nil).
				ExpectGitArgs(append([]string{"checkout", "--"}, p4), "", nil),
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			cmd := oscommands.NewDummyCmdObjBuilder(s.runner)
			assert.NoError(t, runGitCmdOnPaths("checkout", s.paths, cmd))
			s.runner.CheckForMissingCalls()
		})
	}
}
