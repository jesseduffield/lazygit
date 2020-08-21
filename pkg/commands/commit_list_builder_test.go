package commands

import (
	"os/exec"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/stretchr/testify/assert"
)

// NewDummyCommitListBuilder creates a new dummy CommitListBuilder for testing
func NewDummyCommitListBuilder() *CommitListBuilder {
	osCommand := NewDummyOSCommand()

	return &CommitListBuilder{
		Log:                 NewDummyLog(),
		GitCommand:          NewDummyGitCommandWithOSCommand(osCommand),
		OSCommand:           osCommand,
		Tr:                  i18n.NewLocalizer(NewDummyLog()),
		CherryPickedCommits: []*Commit{},
	}
}

// TestCommitListBuilderGetUnpushedCommits is a function.
func TestCommitListBuilderGetUnpushedCommits(t *testing.T) {
	type scenario struct {
		testName string
		command  func(string, ...string) *exec.Cmd
		test     func(map[string]bool)
	}

	scenarios := []scenario{
		{
			"Can't retrieve pushable commits",
			func(string, ...string) *exec.Cmd {
				return exec.Command("test")
			},
			func(pushables map[string]bool) {
				assert.EqualValues(t, map[string]bool{}, pushables)
			},
		},
		{
			"Retrieve pushable commits",
			func(cmd string, args ...string) *exec.Cmd {
				return exec.Command("echo", "8a2bb0e\n78976bc")
			},
			func(pushables map[string]bool) {
				assert.Len(t, pushables, 2)
				assert.EqualValues(t, map[string]bool{"8a2bb0e": true, "78976bc": true}, pushables)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			c := NewDummyCommitListBuilder()
			c.OSCommand.SetCommand(s.command)
			s.test(c.getUnpushedCommits("HEAD"))
		})
	}
}

// TestCommitListBuilderGetMergeBase is a function.
func TestCommitListBuilderGetMergeBase(t *testing.T) {
	type scenario struct {
		testName string
		command  func(string, ...string) *exec.Cmd
		test     func(string, error)
	}

	scenarios := []scenario{
		{
			"swallows an error if the call to merge-base returns an error",
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)

				switch args[0] {
				case "symbolic-ref":
					assert.EqualValues(t, []string{"symbolic-ref", "--short", "HEAD"}, args)
					return exec.Command("echo", "master")
				case "merge-base":
					assert.EqualValues(t, []string{"merge-base", "HEAD", "master"}, args)
					return exec.Command("test")
				}
				return nil
			},
			func(output string, err error) {
				assert.NoError(t, err)
				assert.EqualValues(t, "", output)
			},
		},
		{
			"returns the commit when master",
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)

				switch args[0] {
				case "symbolic-ref":
					assert.EqualValues(t, []string{"symbolic-ref", "--short", "HEAD"}, args)
					return exec.Command("echo", "master")
				case "merge-base":
					assert.EqualValues(t, []string{"merge-base", "HEAD", "master"}, args)
					return exec.Command("echo", "blah")
				}
				return nil
			},
			func(output string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "blah\n", output)
			},
		},
		{
			"checks against develop when a feature branch",
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)

				switch args[0] {
				case "symbolic-ref":
					assert.EqualValues(t, []string{"symbolic-ref", "--short", "HEAD"}, args)
					return exec.Command("echo", "feature/test")
				case "merge-base":
					assert.EqualValues(t, []string{"merge-base", "HEAD", "develop"}, args)
					return exec.Command("echo", "blah")
				}
				return nil
			},
			func(output string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "blah\n", output)
			},
		},
		{
			"bubbles up error if there is one",
			func(cmd string, args ...string) *exec.Cmd {
				return exec.Command("test")
			},
			func(output string, err error) {
				assert.Error(t, err)
				assert.Equal(t, "", output)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			c := NewDummyCommitListBuilder()
			c.OSCommand.SetCommand(s.command)
			s.test(c.getMergeBase())
		})
	}
}
