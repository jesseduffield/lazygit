package commands

import (
	"fmt"
	"strings"

	. "github.com/jesseduffield/lazygit/pkg/commands/types"
)

type ICommitsMgr interface {
	RewordHead(name string) error
	CommitCmdObj(message string, flags string) ICmdObj
	GetHeadMessage() (string, error)
	GetMessage(commitSha string) (string, error)
	GetMessageFirstLine(sha string) (string, error)
	AmendHead() error
	AmendHeadCmdObj() ICmdObj
	ShowCmdObj(sha string, filterPath string) ICmdObj
	Revert(sha string) error
	RevertMerge(sha string, parentNumber int) error
	CreateFixupCommit(sha string) error
}

type CommitsMgr struct {
	commander ICommander
	config    *GitConfig
	quote     func(string) string
}

func NewCommitsMgr(commander *Commander, config *GitConfig, quote func(string) string) *CommitsMgr {
	return &CommitsMgr{
		commander: commander,
		config:    config,
		quote:     quote,
	}
}

// RenameCommit renames the topmost commit with the given name
func (c *CommitsMgr) RewordHead(name string) error {
	return c.commander.RunGitCmdFromStr(fmt.Sprintf("commit --allow-empty --amend --only -m %s", c.quote(name)))
}

type ResetToCommitOptions struct {
	EnvVars []string
}

func (c *CommitsMgr) CommitCmdObj(message string, flags string) ICmdObj {
	splitMessage := strings.Split(message, "\n")
	lineArgs := ""
	for _, line := range splitMessage {
		lineArgs += fmt.Sprintf(" -m %s", c.quote(line))
	}

	flagsStr := ""
	if flags != "" {
		flagsStr = fmt.Sprintf(" %s", flags)
	}

	cmdStr := fmt.Sprintf("commit%s%s", flagsStr, lineArgs)

	return BuildGitCmdObjFromStr(cmdStr)
}

// Get the subject of the HEAD commit
func (c *CommitsMgr) GetHeadMessage() (string, error) {
	cmdObj := BuildGitCmdObjFromStr("log -1 --pretty=%s")
	message, err := c.commander.RunWithOutput(cmdObj)
	return strings.TrimSpace(message), err
}

func (c *CommitsMgr) GetMessage(commitSha string) (string, error) {
	messageWithHeader, err := c.commander.RunWithOutput(
		BuildGitCmdObjFromStr("rev-list --format=%B --max-count=1 " + commitSha),
	)
	message := strings.Join(strings.SplitAfter(messageWithHeader, "\n")[1:], "\n")
	return strings.TrimSpace(message), err
}

func (c *CommitsMgr) GetMessageFirstLine(sha string) (string, error) {
	return c.commander.RunWithOutput(
		BuildGitCmdObjFromStr(fmt.Sprintf("show --no-patch --pretty=format:%%s %s", sha)),
	)
}

// AmendHead amends HEAD with whatever is staged in your working tree
func (c *CommitsMgr) AmendHead() error {
	return c.commander.Run(c.AmendHeadCmdObj())
}

func (c *CommitsMgr) AmendHeadCmdObj() ICmdObj {
	return BuildGitCmdObjFromStr("commit --amend --no-edit --allow-empty")
}

func (c *CommitsMgr) ShowCmdObj(sha string, filterPath string) ICmdObj {
	filterPathArg := ""
	if filterPath != "" {
		filterPathArg = fmt.Sprintf(" -- %s", c.quote(filterPath))
	}
	return BuildGitCmdObjFromStr(
		fmt.Sprintf("show --submodule --color=%s --no-renames --stat -p %s %s", c.config.colorArg(), sha, filterPathArg),
	)
}

// Revert reverts the selected commit by sha
func (c *CommitsMgr) Revert(sha string) error {
	return c.commander.RunGitCmdFromStr(fmt.Sprintf("revert %s", sha))
}

func (c *CommitsMgr) RevertMerge(sha string, parentNumber int) error {
	return c.commander.RunGitCmdFromStr(fmt.Sprintf("revert %s -m %d", sha, parentNumber))
}

// CreateFixupCommit creates a commit that fixes up a previous commit
func (c *CommitsMgr) CreateFixupCommit(sha string) error {
	return c.commander.RunGitCmdFromStr(fmt.Sprintf("commit --fixup=%s", sha))
}
