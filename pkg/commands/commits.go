package commands

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
)

type CommitCommands struct {
	*common.Common

	cmd oscommands.ICmdObjBuilder
}

func NewCommitCommands(
	common *common.Common,
	cmd oscommands.ICmdObjBuilder,
) *CommitCommands {
	return &CommitCommands{
		Common: common,
		cmd:    cmd,
	}
}

// RewordLastCommit renames the topmost commit with the given name
func (self *CommitCommands) RewordLastCommit(name string) error {
	return self.cmd.New("git commit --allow-empty --amend --only -m " + self.cmd.Quote(name)).Run()
}

// ResetToCommit reset to commit
func (self *CommitCommands) ResetToCommit(sha string, strength string, envVars []string) error {
	return self.cmd.New(fmt.Sprintf("git reset --%s %s", strength, sha)).
		// prevents git from prompting us for input which would freeze the program
		// TODO: see if this is actually needed here
		AddEnvVars("GIT_TERMINAL_PROMPT=0").
		AddEnvVars(envVars...).
		Run()
}

func (self *CommitCommands) CommitCmdObj(message string, flags string) oscommands.ICmdObj {
	splitMessage := strings.Split(message, "\n")
	lineArgs := ""
	for _, line := range splitMessage {
		lineArgs += fmt.Sprintf(" -m %s", self.cmd.Quote(line))
	}

	flagsStr := ""
	if flags != "" {
		flagsStr = fmt.Sprintf(" %s", flags)
	}

	return self.cmd.New(fmt.Sprintf("git commit%s%s", flagsStr, lineArgs))
}

// Get the subject of the HEAD commit
func (self *CommitCommands) GetHeadCommitMessage() (string, error) {
	message, err := self.cmd.New("git log -1 --pretty=%s").DontLog().RunWithOutput()
	return strings.TrimSpace(message), err
}

func (self *CommitCommands) GetCommitMessage(commitSha string) (string, error) {
	cmdStr := "git rev-list --format=%B --max-count=1 " + commitSha
	messageWithHeader, err := self.cmd.New(cmdStr).DontLog().RunWithOutput()
	message := strings.Join(strings.SplitAfter(messageWithHeader, "\n")[1:], "\n")
	return strings.TrimSpace(message), err
}

func (self *CommitCommands) GetCommitMessageFirstLine(sha string) (string, error) {
	return self.cmd.New(fmt.Sprintf("git show --no-patch --pretty=format:%%s %s", sha)).DontLog().RunWithOutput()
}

// AmendHead amends HEAD with whatever is staged in your working tree
func (self *CommitCommands) AmendHead() error {
	return self.AmendHeadCmdObj().Run()
}

func (self *CommitCommands) AmendHeadCmdObj() oscommands.ICmdObj {
	return self.cmd.New("git commit --amend --no-edit --allow-empty")
}

func (self *CommitCommands) ShowCmdObj(sha string, filterPath string) oscommands.ICmdObj {
	contextSize := self.UserConfig.Git.DiffContextSize
	filterPathArg := ""
	if filterPath != "" {
		filterPathArg = fmt.Sprintf(" -- %s", self.cmd.Quote(filterPath))
	}

	cmdStr := fmt.Sprintf("git show --submodule --color=%s --unified=%d --no-renames --stat -p %s %s", self.UserConfig.Git.Paging.ColorArg, contextSize, sha, filterPathArg)
	return self.cmd.New(cmdStr).DontLog()
}

// Revert reverts the selected commit by sha
func (self *CommitCommands) Revert(sha string) error {
	return self.cmd.New(fmt.Sprintf("git revert %s", sha)).Run()
}

func (self *CommitCommands) RevertMerge(sha string, parentNumber int) error {
	return self.cmd.New(fmt.Sprintf("git revert %s -m %d", sha, parentNumber)).Run()
}

// CreateFixupCommit creates a commit that fixes up a previous commit
func (self *CommitCommands) CreateFixupCommit(sha string) error {
	return self.cmd.New(fmt.Sprintf("git commit --fixup=%s", sha)).Run()
}
