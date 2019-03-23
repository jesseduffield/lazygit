package git

import (
	"os/exec"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/stretchr/testify/assert"
)

// NewDummyCommitListBuilder creates a new dummy CommitListBuilder for testing
func NewDummyCommitListBuilder() *CommitListBuilder {
	osCommand := commands.NewDummyOSCommand()

	return &CommitListBuilder{
		Log:                 commands.NewDummyLog(),
		GitCommand:          commands.NewDummyGitCommandWithOSCommand(osCommand),
		OSCommand:           osCommand,
		Tr:                  i18n.NewLocalizer(commands.NewDummyLog()),
		CherryPickedCommits: []*commands.Commit{},
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
			s.test(c.getUnpushedCommits())
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

// TestCommitListBuilderGetLog is a function.
func TestCommitListBuilderGetLog(t *testing.T) {
	type scenario struct {
		testName string
		command  func(string, ...string) *exec.Cmd
		test     func(string)
	}

	scenarios := []scenario{
		{
			"Retrieves logs",
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)
				assert.EqualValues(t, []string{"log", "--oneline", "-30"}, args)

				return exec.Command("echo", "6f0b32f commands/git : add GetCommits tests refactor\n9d9d775 circle : remove new line")
			},
			func(output string) {
				assert.EqualValues(t, "6f0b32f commands/git : add GetCommits tests refactor\n9d9d775 circle : remove new line\n", output)
			},
		},
		{
			"An error occurred when retrieving logs",
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)
				assert.EqualValues(t, []string{"log", "--oneline", "-30"}, args)
				return exec.Command("test")
			},
			func(output string) {
				assert.Empty(t, output)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			c := NewDummyCommitListBuilder()
			c.OSCommand.SetCommand(s.command)
			s.test(c.getLog())
		})
	}
}

// TestCommitListBuilderGetCommits is a function.
func TestCommitListBuilderGetCommits(t *testing.T) {
	type scenario struct {
		testName string
		command  func(string, ...string) *exec.Cmd
		test     func([]*commands.Commit, error)
	}

	scenarios := []scenario{
		{
			"No data found",
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)

				switch args[0] {
				case "rev-list":
					assert.EqualValues(t, []string{"rev-list", "@{u}..HEAD", "--abbrev-commit"}, args)
					return exec.Command("echo")
				case "log":
					assert.EqualValues(t, []string{"log", "--oneline", "-30"}, args)
					return exec.Command("echo")
				case "merge-base":
					assert.EqualValues(t, []string{"merge-base", "HEAD", "master"}, args)
					return exec.Command("test")
				case "symbolic-ref":
					assert.EqualValues(t, []string{"symbolic-ref", "--short", "HEAD"}, args)
					return exec.Command("echo", "master")
				}

				return nil
			},
			func(commits []*commands.Commit, err error) {
				assert.NoError(t, err)
				assert.Len(t, commits, 0)
			},
		},
		{
			"GetCommits returns 2 commits, 1 unpushed, the other merged",
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)

				switch args[0] {
				case "rev-list":
					assert.EqualValues(t, []string{"rev-list", "@{u}..HEAD", "--abbrev-commit"}, args)
					return exec.Command("echo", "8a2bb0e")
				case "log":
					assert.EqualValues(t, []string{"log", "--oneline", "-30"}, args)
					return exec.Command("echo", "8a2bb0e commit 1\n78976bc commit 2")
				case "merge-base":
					assert.EqualValues(t, []string{"merge-base", "HEAD", "master"}, args)
					return exec.Command("echo", "78976bc")
				case "symbolic-ref":
					assert.EqualValues(t, []string{"symbolic-ref", "--short", "HEAD"}, args)
					return exec.Command("echo", "master")
				}

				return nil
			},
			func(commits []*commands.Commit, err error) {
				assert.NoError(t, err)
				assert.Len(t, commits, 2)
				assert.EqualValues(t, []*commands.Commit{
					{
						Sha:           "8a2bb0e",
						Name:          "commit 1",
						Status:        "unpushed",
						DisplayString: "8a2bb0e commit 1",
					},
					{
						Sha:           "78976bc",
						Name:          "commit 2",
						Status:        "merged",
						DisplayString: "78976bc commit 2",
					},
				}, commits)
			},
		},
		{
			"GetCommits bubbles up an error from setCommitMergedStatuses",
			func(cmd string, args ...string) *exec.Cmd {
				assert.EqualValues(t, "git", cmd)

				switch args[0] {
				case "rev-list":
					assert.EqualValues(t, []string{"rev-list", "@{u}..HEAD", "--abbrev-commit"}, args)
					return exec.Command("echo", "8a2bb0e")
				case "log":
					assert.EqualValues(t, []string{"log", "--oneline", "-30"}, args)
					return exec.Command("echo", "8a2bb0e commit 1\n78976bc commit 2")
				case "merge-base":
					assert.EqualValues(t, []string{"merge-base", "HEAD", "master"}, args)
					return exec.Command("echo", "78976bc")
				case "symbolic-ref":
					assert.EqualValues(t, []string{"symbolic-ref", "--short", "HEAD"}, args)
					// here's where we are returning the error
					return exec.Command("test")
				case "rev-parse":
					assert.EqualValues(t, []string{"rev-parse", "--short", "HEAD"}, args)
					// here too
					return exec.Command("test")
				}

				return nil
			},
			func(commits []*commands.Commit, err error) {
				assert.Error(t, err)
				assert.Len(t, commits, 0)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			c := NewDummyCommitListBuilder()
			c.OSCommand.SetCommand(s.command)
			s.test(c.GetCommits())
		})
	}
}
