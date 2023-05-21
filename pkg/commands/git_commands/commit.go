package git_commands

import (
	"fmt"
	"strings"
	"time"

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
	cmdStr := NewGitCmd("commit").
		Arg("--allow-empty", "--only", "--no-edit", "--amend", "--reset-author").
		ToString()

	return self.cmd.New(cmdStr).Run()
}

// Sets the commit's author to the supplied value. Value is expected to be of the form 'Name <Email>'
func (self *CommitCommands) SetAuthor(value string) error {
	cmdStr := NewGitCmd("commit").
		Arg("--allow-empty", "--only", "--no-edit", "--amend", "--author="+self.cmd.Quote(value)).
		ToString()

	return self.cmd.New(cmdStr).Run()
}

// ResetToCommit reset to commit
func (self *CommitCommands) ResetToCommit(sha string, strength string, envVars []string) error {
	cmdStr := NewGitCmd("reset").Arg("--"+strength, sha).ToString()

	return self.cmd.New(cmdStr).
		// prevents git from prompting us for input which would freeze the program
		// TODO: see if this is actually needed here
		AddEnvVars("GIT_TERMINAL_PROMPT=0").
		AddEnvVars(envVars...).
		Run()
}

func (self *CommitCommands) CommitCmdObj(message string) oscommands.ICmdObj {
	messageArgs := self.commitMessageArgs(message)

	skipHookPrefix := self.UserConfig.Git.SkipHookPrefix

	cmdStr := NewGitCmd("commit").
		ArgIf(skipHookPrefix != "" && strings.HasPrefix(message, skipHookPrefix), "--no-verify").
		ArgIf(self.signoffFlag() != "", self.signoffFlag()).
		Arg(messageArgs...).
		ToString()

	return self.cmd.New(cmdStr)
}

// RewordLastCommit rewords the topmost commit with the given message
func (self *CommitCommands) RewordLastCommit(message string) error {
	messageArgs := self.commitMessageArgs(message)

	cmdStr := NewGitCmd("commit").
		Arg("--allow-empty", "--amend", "--only").
		Arg(messageArgs...).
		ToString()

	return self.cmd.New(cmdStr).Run()
}

func (self *CommitCommands) commitMessageArgs(message string) []string {
	msg, description, _ := strings.Cut(message, "\n")
	args := []string{"-m", self.cmd.Quote(msg)}

	if description != "" {
		args = append(args, "-m", self.cmd.Quote(description))
	}

	return args
}

// runs git commit without the -m argument meaning it will invoke the user's editor
func (self *CommitCommands) CommitEditorCmdObj() oscommands.ICmdObj {
	cmdStr := NewGitCmd("commit").
		ArgIf(self.signoffFlag() != "", self.signoffFlag()).
		ArgIf(self.verboseFlag() != "", self.verboseFlag()).
		ToString()

	return self.cmd.New(cmdStr)
}

func (self *CommitCommands) signoffFlag() string {
	if self.UserConfig.Git.Commit.SignOff {
		return "--signoff"
	} else {
		return ""
	}
}

func (self *CommitCommands) verboseFlag() string {
	switch self.config.UserConfig.Git.Commit.Verbose {
	case "always":
		return "--verbose"
	case "never":
		return "--no-verbose"
	default:
		return ""
	}
}

// Get the subject of the HEAD commit
func (self *CommitCommands) GetHeadCommitMessage() (string, error) {
	cmdStr := NewGitCmd("log").Arg("-1", "--pretty=%s").ToString()

	message, err := self.cmd.New(cmdStr).DontLog().RunWithOutput()
	return strings.TrimSpace(message), err
}

func (self *CommitCommands) GetCommitMessage(commitSha string) (string, error) {
	cmdStr := NewGitCmd("rev-list").
		Arg("--format=%B", "--max-count=1", commitSha).
		ToString()

	messageWithHeader, err := self.cmd.New(cmdStr).DontLog().RunWithOutput()
	message := strings.Join(strings.SplitAfter(messageWithHeader, "\n")[1:], "")
	return strings.TrimSpace(message), err
}

func (self *CommitCommands) GetCommitDiff(commitSha string) (string, error) {
	cmdStr := NewGitCmd("show").Arg("--no-color", commitSha).ToString()

	diff, err := self.cmd.New(cmdStr).DontLog().RunWithOutput()
	return diff, err
}

type Author struct {
	Name  string
	Email string
}

func (self *CommitCommands) GetCommitAuthor(commitSha string) (Author, error) {
	cmdStr := NewGitCmd("show").
		Arg("--no-patch", "--pretty=format:'%an%x00%ae'", commitSha).
		ToString()

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
	cmdStr := NewGitCmd("show").
		Arg("--no-patch", "--pretty=format:%s").
		Arg(shas...).
		ToString()

	return self.cmd.New(cmdStr).DontLog().RunWithOutput()
}

func (self *CommitCommands) GetCommitsOneline(shas []string) (string, error) {
	cmdStr := NewGitCmd("show").
		Arg("--no-patch", "--oneline").
		Arg(shas...).
		ToString()

	return self.cmd.New(cmdStr).DontLog().RunWithOutput()
}

// AmendHead amends HEAD with whatever is staged in your working tree
func (self *CommitCommands) AmendHead() error {
	return self.AmendHeadCmdObj().Run()
}

func (self *CommitCommands) AmendHeadCmdObj() oscommands.ICmdObj {
	cmdStr := NewGitCmd("commit").
		Arg("--amend", "--no-edit", "--allow-empty").
		ToString()

	return self.cmd.New(cmdStr)
}

func (self *CommitCommands) ShowCmdObj(sha string, filterPath string, ignoreWhitespace bool) oscommands.ICmdObj {
	contextSize := self.UserConfig.Git.DiffContextSize

	cmdStr := NewGitCmd("show").
		Arg("--submodule").
		Arg("--color="+self.UserConfig.Git.Paging.ColorArg).
		Arg(fmt.Sprintf("--unified=%d", contextSize)).
		Arg("--stat").
		Arg("-p").
		Arg(sha).
		ArgIf(ignoreWhitespace, "--ignore-all-space").
		ArgIf(filterPath != "", "--", self.cmd.Quote(filterPath)).
		ToString()

	return self.cmd.New(cmdStr).DontLog()
}

// Revert reverts the selected commit by sha
func (self *CommitCommands) Revert(sha string) error {
	cmdStr := NewGitCmd("revert").Arg(sha).ToString()

	return self.cmd.New(cmdStr).Run()
}

func (self *CommitCommands) RevertMerge(sha string, parentNumber int) error {
	cmdStr := NewGitCmd("revert").Arg(sha, "-m", fmt.Sprintf("%d", parentNumber)).
		ToString()

	return self.cmd.New(cmdStr).Run()
}

// CreateFixupCommit creates a commit that fixes up a previous commit
func (self *CommitCommands) CreateFixupCommit(sha string) error {
	cmdStr := NewGitCmd("commit").Arg("--fixup=" + sha).ToString()

	return self.cmd.New(cmdStr).Run()
}

// a value of 0 means the head commit, 1 is the parent commit, etc
func (self *CommitCommands) GetCommitMessageFromHistory(value int) (string, error) {
	cmdStr := NewGitCmd("log").Arg("-1", fmt.Sprintf("--skip=%d", value), "--pretty=%H").
		ToString()

	hash, _ := self.cmd.New(cmdStr).DontLog().RunWithOutput()
	formattedHash := strings.TrimSpace(hash)
	if len(formattedHash) == 0 {
		return "", ErrInvalidCommitIndex
	}
	return self.GetCommitMessage(formattedHash)
}

// Returns hashes of recent commits which changed the given file
// Note: This does not look for the last X commits to change a file, instead
// it looks among the last X commits and see which of them happened to have changed the file.
// This is more efficient.
func (self *CommitCommands) GetRecentCommitsWhichChangedFile(path string) []string {
	t := time.Now()
	// Checking last X commits. Funnily this seems to actually consider more than the last
	// X, perhaps because of topological sorting.
	cmdStr := NewGitCmd("log").Arg("HEAD~50..HEAD", "--pretty=%H", "--", self.cmd.Quote(path)).
		ToString()

	hashes, _ := self.cmd.New(cmdStr).DontLog().RunWithOutput()
	self.Log.Warn(fmt.Sprintf("GetRecentCommitsWhichChangedFile took %s", time.Since(t)))
	return strings.Split(strings.TrimSpace(hashes), "\n")
}
