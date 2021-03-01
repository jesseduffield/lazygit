package commands

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

// RenameCommit renames the topmost commit with the given name
func (c *GitCommand) RenameCommit(name string) error {
	return c.OSCommand.RunCommand("git commit --allow-empty --amend -m %s", c.OSCommand.Quote(name))
}

// ResetToCommit reset to commit
func (c *GitCommand) ResetToCommit(sha string, strength string, options oscommands.RunCommandOptions) error {
	return c.OSCommand.RunCommandWithOptions(fmt.Sprintf("git reset --%s %s", strength, sha), options)
}

// Commit commits to git
func (c *GitCommand) Commit(message string, flags string) (*exec.Cmd, error) {
	splitMessage := strings.Split(message, "\n")
	lineArgs := ""
	for _, line := range splitMessage {
		lineArgs += fmt.Sprintf(" -m %s", c.OSCommand.Quote(line))
	}

	command := fmt.Sprintf("git commit %s%s", flags, lineArgs)
	if c.usingGpg() {
		return c.OSCommand.ShellCommandFromString(command), nil
	}

	return nil, c.OSCommand.RunCommand(command)
}

// Get the subject of the HEAD commit
func (c *GitCommand) GetHeadCommitMessage() (string, error) {
	cmdStr := "git log -1 --pretty=%s"
	message, err := c.OSCommand.RunCommandWithOutput(cmdStr)
	return strings.TrimSpace(message), err
}

func (c *GitCommand) GetCommitMessage(commitSha string) (string, error) {
	cmdStr := "git rev-list --format=%B --max-count=1 " + commitSha
	messageWithHeader, err := c.OSCommand.RunCommandWithOutput(cmdStr)
	message := strings.Join(strings.SplitAfter(messageWithHeader, "\n")[1:], "\n")
	return strings.TrimSpace(message), err
}

// AmendHead amends HEAD with whatever is staged in your working tree
func (c *GitCommand) AmendHead() (*exec.Cmd, error) {
	command := "git commit --amend --no-edit --allow-empty"
	if c.usingGpg() {
		return c.OSCommand.ShellCommandFromString(command), nil
	}

	return nil, c.OSCommand.RunCommand(command)
}

// PrepareCommitAmendSubProcess prepares a subprocess for `git commit --amend --allow-empty`
func (c *GitCommand) PrepareCommitAmendSubProcess() *exec.Cmd {
	return c.OSCommand.PrepareSubProcess("git", "commit", "--amend", "--allow-empty")
}

func (c *GitCommand) ShowCmdStr(sha string, filterPath string) string {
	filterPathArg := ""
	if filterPath != "" {
		filterPathArg = fmt.Sprintf(" -- %s", c.OSCommand.Quote(filterPath))
	}
	return fmt.Sprintf("git show --submodule --color=%s --no-renames --stat -p %s %s", c.colorArg(), sha, filterPathArg)
}

// Revert reverts the selected commit by sha
func (c *GitCommand) Revert(sha string) error {
	return c.OSCommand.RunCommand("git revert %s", sha)
}

// CherryPickCommits begins an interactive rebase with the given shas being cherry picked onto HEAD
func (c *GitCommand) CherryPickCommits(commits []*models.Commit) error {
	todo := ""
	for _, commit := range commits {
		todo = "pick " + commit.Sha + " " + commit.Name + "\n" + todo
	}

	cmd, err := c.PrepareInteractiveRebaseCommand("HEAD", todo, false)
	if err != nil {
		return err
	}

	return c.OSCommand.RunPreparedCommand(cmd)
}

// CreateFixupCommit creates a commit that fixes up a previous commit
func (c *GitCommand) CreateFixupCommit(sha string) error {
	return c.OSCommand.RunCommand("git commit --fixup=%s", sha)
}
