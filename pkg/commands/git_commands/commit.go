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

// Add a commit's coauthor using Github/Gitlab Co-authored-by metadata. Value is expected to be of the form 'Name <Email>'
func (self *CommitCommands) AddCoAuthor(hash string, author string) error {
	message, err := self.GetCommitMessage(hash)
	if err != nil {
		return err
	}

	message = AddCoAuthorToMessage(message, author)

	cmdArgs := NewGitCmd("commit").
		Arg("--allow-empty", "--amend", "--only", "-m", message).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func AddCoAuthorToMessage(message string, author string) string {
	subject, body, _ := strings.Cut(message, "\n")

	return strings.TrimSpace(subject) + "\n\n" + AddCoAuthorToDescription(strings.TrimSpace(body), author)
}

func AddCoAuthorToDescription(description string, author string) string {
	if description != "" {
		lines := strings.Split(description, "\n")
		if strings.HasPrefix(lines[len(lines)-1], "Co-authored-by:") {
			description += "\n"
		} else {
			description += "\n\n"
		}
	}

	return description + fmt.Sprintf("Co-authored-by: %s", author)
}

// ResetToCommit reset to commit
func (self *CommitCommands) ResetToCommit(hash string, strength string, envVars []string) error {
	cmdArgs := NewGitCmd("reset").Arg("--"+strength, hash).ToArgv()

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

func (self *CommitCommands) RewordLastCommitInEditorWithMessageFileCmdObj(tmpMessageFile string) oscommands.ICmdObj {
	return self.cmd.New(NewGitCmd("commit").
		Arg("--allow-empty", "--amend", "--only", "--edit", "--file="+tmpMessageFile).ToArgv())
}

func (self *CommitCommands) CommitInEditorWithMessageFileCmdObj(tmpMessageFile string) oscommands.ICmdObj {
	return self.cmd.New(NewGitCmd("commit").
		Arg("--edit").
		Arg("--file="+tmpMessageFile).
		ArgIf(self.signoffFlag() != "", self.signoffFlag()).
		ToArgv())
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

func (self *CommitCommands) GetCommitMessage(commitHash string) (string, error) {
	cmdArgs := NewGitCmd("log").
		Arg("--format=%B", "--max-count=1", commitHash).
		Config("log.showsignature=false").
		ToArgv()

	message, err := self.cmd.New(cmdArgs).DontLog().RunWithOutput()
	return strings.ReplaceAll(strings.TrimSpace(message), "\r\n", "\n"), err
}

func (self *CommitCommands) GetCommitSubject(commitHash string) (string, error) {
	cmdArgs := NewGitCmd("log").
		Arg("--format=%s", "--max-count=1", commitHash).
		Config("log.showsignature=false").
		ToArgv()

	subject, err := self.cmd.New(cmdArgs).DontLog().RunWithOutput()
	return strings.TrimSpace(subject), err
}

func (self *CommitCommands) GetCommitDiff(commitHash string) (string, error) {
	cmdArgs := NewGitCmd("show").Arg("--no-color", commitHash).ToArgv()

	diff, err := self.cmd.New(cmdArgs).DontLog().RunWithOutput()
	return diff, err
}

type Author struct {
	Name  string
	Email string
}

func (self *CommitCommands) GetCommitAuthor(commitHash string) (Author, error) {
	cmdArgs := NewGitCmd("show").
		Arg("--no-patch", "--pretty=format:'%an%x00%ae'", commitHash).
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

func (self *CommitCommands) GetCommitMessageFirstLine(hash string) (string, error) {
	return self.GetCommitMessagesFirstLine([]string{hash})
}

func (self *CommitCommands) GetCommitMessagesFirstLine(hashes []string) (string, error) {
	cmdArgs := NewGitCmd("show").
		Arg("--no-patch", "--pretty=format:%s").
		Arg(hashes...).
		ToArgv()

	return self.cmd.New(cmdArgs).DontLog().RunWithOutput()
}

// Example output:
//
//	cd50c79ae Preserve the commit message correctly even if the description has blank lines
//	3ebba5f32 Add test demonstrating a bug with preserving the commit message
//	9a423c388 Remove unused function
func (self *CommitCommands) GetHashesAndCommitMessagesFirstLine(hashes []string) (string, error) {
	cmdArgs := NewGitCmd("show").
		Arg("--no-patch", "--pretty=format:%h %s").
		Arg(hashes...).
		ToArgv()

	return self.cmd.New(cmdArgs).DontLog().RunWithOutput()
}

func (self *CommitCommands) GetCommitsOneline(hashes []string) (string, error) {
	cmdArgs := NewGitCmd("show").
		Arg("--no-patch", "--oneline").
		Arg(hashes...).
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

func (self *CommitCommands) ShowCmdObj(hash string, filterPath string) oscommands.ICmdObj {
	contextSize := self.AppState.DiffContextSize

	extDiffCmd := self.UserConfig.Git.Paging.ExternalDiffCommand
	cmdArgs := NewGitCmd("show").
		Config("diff.noprefix=false").
		ConfigIf(extDiffCmd != "", "diff.external="+extDiffCmd).
		ArgIfElse(extDiffCmd != "", "--ext-diff", "--no-ext-diff").
		Arg("--submodule").
		Arg("--color="+self.UserConfig.Git.Paging.ColorArg).
		Arg(fmt.Sprintf("--unified=%d", contextSize)).
		Arg("--stat").
		Arg("--decorate").
		Arg("-p").
		Arg(hash).
		ArgIf(self.AppState.IgnoreWhitespaceInDiffView, "--ignore-all-space").
		ArgIf(filterPath != "", "--", filterPath).
		Dir(self.repoPaths.worktreePath).
		ToArgv()

	return self.cmd.New(cmdArgs).DontLog()
}

// Revert reverts the selected commit by hash
func (self *CommitCommands) Revert(hash string) error {
	cmdArgs := NewGitCmd("revert").Arg(hash).ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *CommitCommands) RevertMerge(hash string, parentNumber int) error {
	cmdArgs := NewGitCmd("revert").Arg(hash, "-m", fmt.Sprintf("%d", parentNumber)).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

// CreateFixupCommit creates a commit that fixes up a previous commit
func (self *CommitCommands) CreateFixupCommit(hash string) error {
	cmdArgs := NewGitCmd("commit").Arg("--fixup=" + hash).ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

// CreateAmendCommit creates a commit that changes the commit message of a previous commit
func (self *CommitCommands) CreateAmendCommit(originalSubject, newSubject, newDescription string, includeFileChanges bool) error {
	description := newSubject
	if newDescription != "" {
		description += "\n\n" + newDescription
	}
	cmdArgs := NewGitCmd("commit").
		Arg("-m", "amend! "+originalSubject).
		Arg("-m", description).
		ArgIf(!includeFileChanges, "--only", "--allow-empty").
		ToArgv()

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
