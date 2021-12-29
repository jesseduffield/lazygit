package commands

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

// RenameCommit renames the topmost commit with the given name
func (c *GitCommand) RenameCommit(name string) error {
	return c.NewCmdObj("git commit --allow-empty --amend --only -m " + c.OSCommand.Quote(name)).Run()
}

// ResetToCommit reset to commit
func (c *GitCommand) ResetToCommit(sha string, strength string, envVars []string) error {
	return c.NewCmdObj(fmt.Sprintf("git reset --%s %s", strength, sha)).
		// prevents git from prompting us for input which would freeze the program
		// TODO: see if this is actually needed here
		AddEnvVars("GIT_TERMINAL_PROMPT=0").
		AddEnvVars(envVars...).
		Run()
}

func (c *GitCommand) CommitCmdObj(message string, flags string) oscommands.ICmdObj {
	splitMessage := strings.Split(message, "\n")
	lineArgs := ""
	for _, line := range splitMessage {
		lineArgs += fmt.Sprintf(" -m %s", c.OSCommand.Quote(line))
	}

	flagsStr := ""
	if flags != "" {
		flagsStr = fmt.Sprintf(" %s", flags)
	}

	return c.NewCmdObj(fmt.Sprintf("git commit%s%s", flagsStr, lineArgs))
}

// Get the subject of the HEAD commit
func (c *GitCommand) GetHeadCommitMessage() (string, error) {
	message, err := c.NewCmdObj("git log -1 --pretty=%s").RunWithOutput()
	return strings.TrimSpace(message), err
}

func (c *GitCommand) GetCommitMessage(commitSha string) (string, error) {
	cmdStr := "git rev-list --format=%B --max-count=1 " + commitSha
	messageWithHeader, err := c.NewCmdObj(cmdStr).RunWithOutput()
	message := strings.Join(strings.SplitAfter(messageWithHeader, "\n")[1:], "\n")
	return strings.TrimSpace(message), err
}

func (c *GitCommand) GetCommitMessageFirstLine(sha string) (string, error) {
	return c.NewCmdObj(fmt.Sprintf("git show --no-patch --pretty=format:%%s %s", sha)).RunWithOutput()
}

// AmendHead amends HEAD with whatever is staged in your working tree
func (c *GitCommand) AmendHead() error {
	return c.AmendHeadCmdObj().Run()
}

func (c *GitCommand) AmendHeadCmdObj() oscommands.ICmdObj {
	return c.NewCmdObj("git commit --amend --no-edit --allow-empty")
}

func (c *GitCommand) ShowCmdObj(sha string, filterPath string) oscommands.ICmdObj {
	contextSize := c.UserConfig.Git.DiffContextSize
	filterPathArg := ""
	if filterPath != "" {
		filterPathArg = fmt.Sprintf(" -- %s", c.OSCommand.Quote(filterPath))
	}

	cmdStr := fmt.Sprintf("git show --submodule --color=%s --unified=%d --no-renames --stat -p %s %s", c.colorArg(), contextSize, sha, filterPathArg)
	return c.NewCmdObj(cmdStr)
}

// Revert reverts the selected commit by sha
func (c *GitCommand) Revert(sha string) error {
	return c.NewCmdObj(fmt.Sprintf("git revert %s", sha)).Run()
}

func (c *GitCommand) RevertMerge(sha string, parentNumber int) error {
	return c.NewCmdObj(fmt.Sprintf("git revert %s -m %d", sha, parentNumber)).Run()
}

// CherryPickCommits begins an interactive rebase with the given shas being cherry picked onto HEAD
func (c *GitCommand) CherryPickCommits(commits []*models.Commit) error {
	todo := ""
	for _, commit := range commits {
		todo = "pick " + commit.Sha + " " + commit.Name + "\n" + todo
	}

	cmdObj, err := c.PrepareInteractiveRebaseCommand("HEAD", todo, false)
	if err != nil {
		return err
	}

	return cmdObj.Run()
}

// CreateFixupCommit creates a commit that fixes up a previous commit
func (c *GitCommand) CreateFixupCommit(sha string) error {
	return c.NewCmdObj(fmt.Sprintf("git commit --fixup=%s", sha)).Run()
}
