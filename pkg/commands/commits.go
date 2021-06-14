package commands

import (
	"fmt"
	"strings"

	. "github.com/jesseduffield/lazygit/pkg/commands/types"
)

// type IGitCommon interface {

// }

// type Commits struct {
// 	IGitCommon
// }

// RenameCommit renames the topmost commit with the given name
func (c *Git) RenameHeadCommit(name string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("commit --allow-empty --amend --only -m %s", c.GetOS().Quote(name)))
}

type ResetToCommitOptions struct {
	EnvVars []string
}

// ResetToCommit reset to commit
func (c *Git) ResetToRef(ref string, strength string, options ResetToCommitOptions) error {
	cmdObj := BuildGitCmdObjFromStr(fmt.Sprintf("reset --%s %s", strength, ref))
	cmdObj.AddEnvVars(options.EnvVars...)

	return c.Run(cmdObj)
}

func (c *Git) CommitCmdObj(message string, flags string) ICmdObj {
	splitMessage := strings.Split(message, "\n")
	lineArgs := ""
	for _, line := range splitMessage {
		lineArgs += fmt.Sprintf(" -m %s", c.GetOS().Quote(line))
	}

	flagsStr := ""
	if flags != "" {
		flagsStr = fmt.Sprintf(" %s", flags)
	}

	cmdStr := fmt.Sprintf("commit%s%s", flagsStr, lineArgs)

	return BuildGitCmdObjFromStr(cmdStr)
}

// Get the subject of the HEAD commit
func (c *Git) GetHeadCommitMessage() (string, error) {
	cmdObj := BuildGitCmdObjFromStr("log -1 --pretty=%s")
	message, err := c.RunWithOutput(cmdObj)
	return strings.TrimSpace(message), err
}

func (c *Git) GetCommitMessage(commitSha string) (string, error) {
	messageWithHeader, err := c.RunWithOutput(
		BuildGitCmdObjFromStr("rev-list --format=%B --max-count=1 " + commitSha),
	)
	message := strings.Join(strings.SplitAfter(messageWithHeader, "\n")[1:], "\n")
	return strings.TrimSpace(message), err
}

func (c *Git) GetCommitMessageFirstLine(sha string) (string, error) {
	return c.RunWithOutput(
		BuildGitCmdObjFromStr(fmt.Sprintf("show --no-patch --pretty=format:%%s %s", sha)),
	)
}

// AmendHead amends HEAD with whatever is staged in your working tree
func (c *Git) AmendHead() error {
	return c.Run(c.AmendHeadCmdObj())
}

func (c *Git) AmendHeadCmdObj() ICmdObj {
	return BuildGitCmdObjFromStr("commit --amend --no-edit --allow-empty")
}

func (c *Git) ShowCmdObj(sha string, filterPath string) ICmdObj {
	filterPathArg := ""
	if filterPath != "" {
		filterPathArg = fmt.Sprintf(" -- %s", c.GetOS().Quote(filterPath))
	}
	return BuildGitCmdObjFromStr(
		fmt.Sprintf("show --submodule --color=%s --no-renames --stat -p %s %s", c.colorArg(), sha, filterPathArg),
	)
}

// Revert reverts the selected commit by sha
func (c *Git) Revert(sha string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("revert %s", sha))
}

func (c *Git) RevertMerge(sha string, parentNumber int) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("revert %s -m %d", sha, parentNumber))
}

// CreateFixupCommit creates a commit that fixes up a previous commit
func (c *Git) CreateFixupCommit(sha string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("commit --fixup=%s", sha))
}
