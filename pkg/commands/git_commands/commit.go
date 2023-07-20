package git_commands

import (
	"fmt"
	"strings"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

var ErrInvalidCommitIndex = errors.New("invalid commit index")

type CommitCommands struct {
	*GitCommon
}

func NewCommitCommands(gitCommon *GitCommon) *CommitCommands {
	return &CommitCommands{
		GitCommon: gitCommon,
	}
}

// ResetAuthor resets the author of the topmost commit
func (self *CommitCommands) ResetAuthor() error {
	cmdArgs := NewGitCmd("commit").
		Arg("--allow-empty", "--only", "--no-edit", "--amend", "--reset-author").
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

// Sets the commit's author to the supplied value. Value is expected to be of the form 'Name <Email>'
func (self *CommitCommands) SetAuthor(value string) error {
	cmdArgs := NewGitCmd("commit").
		Arg("--allow-empty", "--only", "--no-edit", "--amend", "--author="+value).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

// ResetToCommit reset to commit
func (self *CommitCommands) ResetToCommit(sha string, strength string, envVars []string) error {
	cmdArgs := NewGitCmd("reset").Arg("--"+strength, sha).ToArgv()

	return self.cmd.New(cmdArgs).
		// prevents git from prompting us for input which would freeze the program
		// TODO: see if this is actually needed here
		AddEnvVars("GIT_TERMINAL_PROMPT=0").
		AddEnvVars(envVars...).
		Run()
}

func (self *CommitCommands) CommitCmdObj(summary string, description string) oscommands.ICmdObj {
	messageArgs := self.commitMessageArgs(summary, description)

	skipHookPrefix := self.UserConfig.Git.SkipHookPrefix

	cmdArgs := NewGitCmd("commit").
		ArgIf(skipHookPrefix != "" && strings.HasPrefix(summary, skipHookPrefix), "--no-verify").
		ArgIf(self.signoffFlag() != "", self.signoffFlag()).
		Arg(messageArgs...).
		ToArgv()

	return self.cmd.New(cmdArgs)
}

func (self *CommitCommands) RewordLastCommitInEditorCmdObj() oscommands.ICmdObj {
	return self.cmd.New(NewGitCmd("commit").Arg("--allow-empty", "--amend", "--only").ToArgv())
}

// RewordLastCommit rewords the topmost commit with the given message
func (self *CommitCommands) RewordLastCommit(summary string, description string) error {
	messageArgs := self.commitMessageArgs(summary, description)

	cmdArgs := NewGitCmd("commit").
		Arg("--allow-empty", "--amend", "--only").
		Arg(messageArgs...).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *CommitCommands) commitMessageArgs(summary string, description string) []string {
	args := []string{"-m", summary}

	if description != "" {
		args = append(args, "-m", description)
	}

	return args
}

// runs git commit without the -m argument meaning it will invoke the user's editor
func (self *CommitCommands) CommitEditorCmdObj() oscommands.ICmdObj {
	cmdArgs := NewGitCmd("commit").
		ArgIf(self.signoffFlag() != "", self.signoffFlag()).
		ToArgv()

	return self.cmd.New(cmdArgs)
}

func (self *CommitCommands) signoffFlag() string {
	if self.UserConfig.Git.Commit.SignOff {
		return "--signoff"
	} else {
		return ""
	}
}

func (self *CommitCommands) GetCommitMessage(commitSha string) (string, error) {
	cmdArgs := NewGitCmd("rev-list").
		Arg("--format=%B", "--max-count=1", commitSha).
		ToArgv()

	messageWithHeader, err := self.cmd.New(cmdArgs).DontLog().RunWithOutput()
	message := strings.Join(strings.SplitAfter(messageWithHeader, "\n")[1:], "")
	return strings.TrimSpace(message), err
}

func (self *CommitCommands) GetCommitDiff(commitSha string) (string, error) {
	cmdArgs := NewGitCmd("show").Arg("--no-color", commitSha).ToArgv()

	diff, err := self.cmd.New(cmdArgs).DontLog().RunWithOutput()
	return diff, err
}

type Author struct {
	Name  string
	Email string
}

func (self *CommitCommands) GetCommitAuthor(commitSha string) (Author, error) {
	cmdArgs := NewGitCmd("show").
		Arg("--no-patch", "--pretty=format:'%an%x00%ae'", commitSha).
		ToArgv()

	output, err := self.cmd.New(cmdArgs).DontLog().RunWithOutput()
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
	cmdArgs := NewGitCmd("show").
		Arg("--no-patch", "--pretty=format:%s").
		Arg(shas...).
		ToArgv()

	return self.cmd.New(cmdArgs).DontLog().RunWithOutput()
}

func (self *CommitCommands) GetCommitsOneline(shas []string) (string, error) {
	cmdArgs := NewGitCmd("show").
		Arg("--no-patch", "--oneline").
		Arg(shas...).
		ToArgv()

	return self.cmd.New(cmdArgs).DontLog().RunWithOutput()
}

// AmendHead amends HEAD with whatever is staged in your working tree
func (self *CommitCommands) AmendHead() error {
	return self.AmendHeadCmdObj().Run()
}

func (self *CommitCommands) AmendHeadCmdObj() oscommands.ICmdObj {
	cmdArgs := NewGitCmd("commit").
		Arg("--amend", "--no-edit", "--allow-empty").
		ToArgv()

	return self.cmd.New(cmdArgs)
}

func (self *CommitCommands) ShowCmdObj(sha string, filterPath string, ignoreWhitespace bool) oscommands.ICmdObj {
	contextSize := self.UserConfig.Git.DiffContextSize

	cmdArgs := NewGitCmd("show").
		Arg("--submodule").
		Arg("--color="+self.UserConfig.Git.Paging.ColorArg).
		Arg(fmt.Sprintf("--unified=%d", contextSize)).
		Arg("--stat").
		Arg("--decorate").
		Arg("-p").
		Arg(sha).
		ArgIf(ignoreWhitespace, "--ignore-all-space").
		ArgIf(filterPath != "", "--", filterPath).
		ToArgv()

	return self.cmd.New(cmdArgs).DontLog()
}

// Revert reverts the selected commit by sha
func (self *CommitCommands) Revert(sha string) error {
	cmdArgs := NewGitCmd("revert").Arg(sha).ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *CommitCommands) RevertMerge(sha string, parentNumber int) error {
	cmdArgs := NewGitCmd("revert").Arg(sha, "-m", fmt.Sprintf("%d", parentNumber)).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

// CreateFixupCommit creates a commit that fixes up a previous commit
func (self *CommitCommands) CreateFixupCommit(sha string) error {
	cmdArgs := NewGitCmd("commit").Arg("--fixup=" + sha).ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

// a value of 0 means the head commit, 1 is the parent commit, etc
func (self *CommitCommands) GetCommitMessageFromHistory(value int) (string, error) {
	cmdArgs := NewGitCmd("log").Arg("-1", fmt.Sprintf("--skip=%d", value), "--pretty=%H").
		ToArgv()

	hash, _ := self.cmd.New(cmdArgs).DontLog().RunWithOutput()
	formattedHash := strings.TrimSpace(hash)
	if len(formattedHash) == 0 {
		return "", ErrInvalidCommitIndex
	}
	return self.GetCommitMessage(formattedHash)
}
