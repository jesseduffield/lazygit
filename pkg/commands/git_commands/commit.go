package git_commands

import (
	"fmt"
	"strings"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

type CommitCommands struct {
	*GitCommon
}

func NewCommitCommands(gitCommon *GitCommon) *CommitCommands {
	return &CommitCommands{
		GitCommon: gitCommon,
	}
}

// RewordLastCommit rewords the topmost commit with the given message
func (self *CommitCommands) RewordLastCommit(message string) error {
	return self.cmd.New("git commit --allow-empty --amend --only -m " + self.cmd.Quote(message)).Run()
}

// ResetAuthor resets the author of the topmost commit
func (self *CommitCommands) ResetAuthor() error {
	return self.cmd.New("git commit --allow-empty --only --no-edit --amend --reset-author").Run()
}

// Sets the commit's author to the supplied value. Value is expected to be of the form 'Name <Email>'
func (self *CommitCommands) SetAuthor(value string) error {
	commandStr := fmt.Sprintf("git commit --allow-empty --only --no-edit --amend --author=%s", self.cmd.Quote(value))
	return self.cmd.New(commandStr).Run()
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

func (self *CommitCommands) CommitCmdObj(message string) oscommands.ICmdObj {
	splitMessage := strings.Split(message, "\n")
	lineArgs := ""
	for _, line := range splitMessage {
		lineArgs += fmt.Sprintf(" -m %s", self.cmd.Quote(line))
	}

	skipHookPrefix := self.UserConfig.Git.SkipHookPrefix
	noVerifyFlag := ""
	if skipHookPrefix != "" && strings.HasPrefix(message, skipHookPrefix) {
		noVerifyFlag = " --no-verify"
	}

	return self.cmd.New(fmt.Sprintf("git commit%s%s%s", noVerifyFlag, self.signoffFlag(), lineArgs))
}

// runs git commit without the -m argument meaning it will invoke the user's editor
func (self *CommitCommands) CommitEditorCmdObj() oscommands.ICmdObj {
	return self.cmd.New(fmt.Sprintf("git commit%s%s", self.signoffFlag(), self.verboseFlag()))
}

func (self *CommitCommands) signoffFlag() string {
	if self.UserConfig.Git.Commit.SignOff {
		return " --signoff"
	} else {
		return ""
	}
}

func (self *CommitCommands) verboseFlag() string {
	switch self.config.UserConfig.Git.Commit.Verbose {
	case "always":
		return " --verbose"
	case "never":
		return " --no-verbose"
	default:
		return ""
	}
}

// Get the subject of the HEAD commit
func (self *CommitCommands) GetHeadCommitMessage() (string, error) {
	message, err := self.cmd.New("git log -1 --pretty=%s").DontLog().RunWithOutput()
	return strings.TrimSpace(message), err
}

func (self *CommitCommands) GetCommitMessage(commitSha string) (string, error) {
	cmdStr := "git rev-list --format=%B --max-count=1 " + commitSha
	messageWithHeader, err := self.cmd.New(cmdStr).DontLog().RunWithOutput()
	message := strings.Join(strings.SplitAfter(messageWithHeader, "\n")[1:], "")
	return strings.TrimSpace(message), err
}

func (self *CommitCommands) GetCommitDiff(commitSha string) (string, error) {
	cmdStr := "git show --no-color " + commitSha
	diff, err := self.cmd.New(cmdStr).DontLog().RunWithOutput()
	return diff, err
}

type Author struct {
	Name  string
	Email string
}

func (self *CommitCommands) GetCommitAuthor(commitSha string) (Author, error) {
	cmdStr := "git show --no-patch --pretty=format:'%an%x00%ae' " + commitSha
	output, err := self.cmd.New(cmdStr).DontLog().RunWithOutput()
	if err != nil {
		return Author{}, err
	}

	split := strings.SplitN(strings.TrimSpace(output), "\x00", 2)
	if len(split) < 2 {
		return Author{}, errors.New("unexpected git output")
	}

	author := Author{Name: split[0], Email: split[1]}
	return author, err
}

func (self *CommitCommands) GetCommitMessageFirstLine(sha string) (string, error) {
	return self.GetCommitMessagesFirstLine([]string{sha})
}

func (self *CommitCommands) GetCommitMessagesFirstLine(shas []string) (string, error) {
	return self.cmd.New(
		fmt.Sprintf("git show --no-patch --pretty=format:%%s %s", strings.Join(shas, " ")),
	).DontLog().RunWithOutput()
}

func (self *CommitCommands) GetCommitsOneline(shas []string) (string, error) {
	return self.cmd.New(
		fmt.Sprintf("git show --no-patch --oneline %s", strings.Join(shas, " ")),
	).DontLog().RunWithOutput()
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
