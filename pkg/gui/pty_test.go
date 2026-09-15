package gui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithPtyGitConfig(t *testing.T) {
	args := []string{"git", "-C", "/repo", "diff", "--color=always"}

	assert.Equal(t,
		[]string{"git", "-c", "diff.autoRefreshIndex=false", "-C", "/repo", "diff", "--color=always"},
		withPtyGitConfig(args, "windows"))

	assert.Equal(t, args, withPtyGitConfig(args, "linux"))
	assert.Equal(t, args, withPtyGitConfig(args, "darwin"))

	// A user-configured command that wraps git in a shell must not have git
	// flags injected into it.
	shellArgs := []string{"sh", "-c", "git log --graph {{branchName}} -- | sed -e s/x/y/"}
	assert.Equal(t, shellArgs, withPtyGitConfig(shellArgs, "windows"))

	// The guard recognizes git regardless of case and extension.
	exeArgs := []string{"GIT.EXE", "diff"}
	assert.Equal(t,
		[]string{"GIT.EXE", "-c", "diff.autoRefreshIndex=false", "diff"},
		withPtyGitConfig(exeArgs, "windows"))
}
