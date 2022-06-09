package git_commands

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/app/daemon"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

type RebaseCommands struct {
	*GitCommon
	commit      *CommitCommands
	workingTree *WorkingTreeCommands

	onSuccessfulContinue func() error
}

func NewRebaseCommands(
	gitCommon *GitCommon,
	commitCommands *CommitCommands,
	workingTreeCommands *WorkingTreeCommands,
) *RebaseCommands {
	return &RebaseCommands{
		GitCommon:   gitCommon,
		commit:      commitCommands,
		workingTree: workingTreeCommands,
	}
}

func (self *RebaseCommands) RewordCommit(commits []*models.Commit, index int, message string) error {
	if index == 0 {
		// we've selected the top commit so no rebase is required
		return self.commit.RewordLastCommit(message)
	}

	err := self.BeginInteractiveRebaseForCommit(commits, index)
	if err != nil {
		return err
	}

	// now the selected commit should be our head so we'll amend it with the new message
	err = self.commit.RewordLastCommit(message)
	if err != nil {
		return err
	}

	return self.ContinueRebase()
}

func (self *RebaseCommands) RewordCommitInEditor(commits []*models.Commit, index int) (oscommands.ICmdObj, error) {
	todo, sha, err := self.BuildSingleActionTodo(commits, index, "reword")
	if err != nil {
		return nil, err
	}

	return self.PrepareInteractiveRebaseCommand(sha, todo, false), nil
}

func (self *RebaseCommands) ResetCommitAuthor(commits []*models.Commit, index int) error {
	return self.GenericAmend(commits, index, func() error {
		return self.commit.ResetAuthor()
	})
}

func (self *RebaseCommands) SetCommitAuthor(commits []*models.Commit, index int, value string) error {
	return self.GenericAmend(commits, index, func() error {
		return self.commit.SetAuthor(value)
	})
}

func (self *RebaseCommands) GenericAmend(commits []*models.Commit, index int, f func() error) error {
	if index == 0 {
		// we've selected the top commit so no rebase is required
		return f()
	}

	err := self.BeginInteractiveRebaseForCommit(commits, index)
	if err != nil {
		return err
	}

	// now the selected commit should be our head so we'll amend it
	err = f()
	if err != nil {
		return err
	}

	return self.ContinueRebase()
}

func (self *RebaseCommands) MoveCommitDown(commits []*models.Commit, index int) error {
	// we must ensure that we have at least two commits after the selected one
	if len(commits) <= index+2 {
		// assuming they aren't picking the bottom commit
		return errors.New(self.Tr.NoRoom)
	}

	orderedCommits := append(commits[0:index], commits[index+1], commits[index])

	todoLines := self.BuildTodoLinesSingleAction(orderedCommits, "pick")

	return self.PrepareInteractiveRebaseCommand(commits[index+2].Sha, todoLines, true).Run()
}

func (self *RebaseCommands) InteractiveRebase(commits []*models.Commit, index int, action string) error {
	todo, sha, err := self.BuildSingleActionTodo(commits, index, action)
	if err != nil {
		return err
	}

	return self.PrepareInteractiveRebaseCommand(sha, todo, true).Run()
}

// PrepareInteractiveRebaseCommand returns the cmd for an interactive rebase
// we tell git to run lazygit to edit the todo list, and we pass the client
// lazygit a todo string to write to the todo file
func (self *RebaseCommands) PrepareInteractiveRebaseCommand(baseSha string, todoLines []TodoLine, overrideEditor bool) oscommands.ICmdObj {
	todo := self.buildTodo(todoLines)
	ex := oscommands.GetLazygitPath()

	debug := "FALSE"
	if self.Debug {
		debug = "TRUE"
	}

	cmdStr := fmt.Sprintf("git rebase --interactive --autostash --keep-empty %s", baseSha)
	self.Log.WithField("command", cmdStr).Info("RunCommand")

	cmdObj := self.cmd.New(cmdStr)

	gitSequenceEditor := ex
	if todo == "" {
		gitSequenceEditor = "true"
	} else {
		self.os.LogCommand(fmt.Sprintf("Creating TODO file for interactive rebase: \n\n%s", todo), false)
	}

	cmdObj.AddEnvVars(
		daemon.DaemonKindEnvKey+"="+string(daemon.InteractiveRebase),
		daemon.RebaseTODOEnvKey+"="+todo,
		"DEBUG="+debug,
		"LANG=en_US.UTF-8",   // Force using EN as language
		"LC_ALL=en_US.UTF-8", // Force using EN as language
		"GIT_SEQUENCE_EDITOR="+gitSequenceEditor,
	)

	if overrideEditor {
		cmdObj.AddEnvVars("GIT_EDITOR=" + ex)
	}

	return cmdObj
}

// produces TodoLines where every commit is picked (or dropped for merge commits) except for the commit at the given index, which
// will have the given action applied to it.
func (self *RebaseCommands) BuildSingleActionTodo(commits []*models.Commit, actionIndex int, action string) ([]TodoLine, string, error) {
	baseIndex := actionIndex + 1

	if len(commits) <= baseIndex {
		return nil, "", errors.New(self.Tr.CannotRebaseOntoFirstCommit)
	}

	if action == "squash" || action == "fixup" {
		baseIndex++

		if len(commits) <= baseIndex {
			return nil, "", errors.New(self.Tr.CannotSquashOntoSecondCommit)
		}
	}

	todoLines := self.BuildTodoLines(commits[0:baseIndex], func(commit *models.Commit, i int) string {
		if i == actionIndex {
			return action
		} else if commit.IsMerge() {
			// your typical interactive rebase will actually drop merge commits by default. Damn git CLI, you scary!
			// doing this means we don't need to worry about rebasing over merges which always causes problems.
			// you typically shouldn't be doing rebases that pass over merge commits anyway.
			return "drop"
		} else {
			return "pick"
		}
	})

	return todoLines, commits[baseIndex].Sha, nil
}

// AmendTo amends the given commit with whatever files are staged
func (self *RebaseCommands) AmendTo(sha string) error {
	if err := self.commit.CreateFixupCommit(sha); err != nil {
		return err
	}

	return self.SquashAllAboveFixupCommits(sha)
}

// EditRebaseTodo sets the action at a given index in the git-rebase-todo file
func (self *RebaseCommands) EditRebaseTodo(index int, action string) error {
	fileName := filepath.Join(self.dotGitDir, "rebase-merge/git-rebase-todo")
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	content := strings.Split(string(bytes), "\n")
	commitCount := self.getTodoCommitCount(content)

	// we have the most recent commit at the bottom whereas the todo file has
	// it at the bottom, so we need to subtract our index from the commit count
	contentIndex := commitCount - 1 - index
	splitLine := strings.Split(content[contentIndex], " ")
	content[contentIndex] = action + " " + strings.Join(splitLine[1:], " ")
	result := strings.Join(content, "\n")

	return ioutil.WriteFile(fileName, []byte(result), 0o644)
}

func (self *RebaseCommands) getTodoCommitCount(content []string) int {
	// count lines that are not blank and are not comments
	commitCount := 0
	for _, line := range content {
		if line != "" && !strings.HasPrefix(line, "#") {
			commitCount++
		}
	}
	return commitCount
}

// MoveTodoDown moves a rebase todo item down by one position
func (self *RebaseCommands) MoveTodoDown(index int) error {
	fileName := filepath.Join(self.dotGitDir, "rebase-merge/git-rebase-todo")
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	content := strings.Split(string(bytes), "\n")
	commitCount := self.getTodoCommitCount(content)
	contentIndex := commitCount - 1 - index

	rearrangedContent := append(content[0:contentIndex-1], content[contentIndex], content[contentIndex-1])
	rearrangedContent = append(rearrangedContent, content[contentIndex+1:]...)
	result := strings.Join(rearrangedContent, "\n")

	return ioutil.WriteFile(fileName, []byte(result), 0o644)
}

// SquashAllAboveFixupCommits squashes all fixup! commits above the given one
func (self *RebaseCommands) SquashAllAboveFixupCommits(sha string) error {
	return self.runSkipEditorCommand(
		self.cmd.New(
			fmt.Sprintf(
				"git rebase --interactive --autostash --autosquash %s^",
				sha,
			),
		),
	)
}

// BeginInteractiveRebaseForCommit starts an interactive rebase to edit the current
// commit and pick all others. After this you'll want to call `self.ContinueRebase()
func (self *RebaseCommands) BeginInteractiveRebaseForCommit(commits []*models.Commit, commitIndex int) error {
	if len(commits)-1 < commitIndex {
		return errors.New("index outside of range of commits")
	}

	// we can make this GPG thing possible it just means we need to do this in two parts:
	// one where we handle the possibility of a credential request, and the other
	// where we continue the rebase
	if self.config.UsingGpg() {
		return errors.New(self.Tr.DisabledForGPG)
	}

	todo, sha, err := self.BuildSingleActionTodo(commits, commitIndex, "edit")
	if err != nil {
		return err
	}

	return self.PrepareInteractiveRebaseCommand(sha, todo, true).Run()
}

// RebaseBranch interactive rebases onto a branch
func (self *RebaseCommands) RebaseBranch(branchName string) error {
	return self.PrepareInteractiveRebaseCommand(branchName, nil, false).Run()
}

func (self *RebaseCommands) GenericMergeOrRebaseActionCmdObj(commandType string, command string) oscommands.ICmdObj {
	return self.cmd.New("git " + commandType + " --" + command)
}

func (self *RebaseCommands) ContinueRebase() error {
	return self.GenericMergeOrRebaseAction("rebase", "continue")
}

func (self *RebaseCommands) AbortRebase() error {
	return self.GenericMergeOrRebaseAction("rebase", "abort")
}

// GenericMerge takes a commandType of "merge" or "rebase" and a command of "abort", "skip" or "continue"
// By default we skip the editor in the case where a commit will be made
func (self *RebaseCommands) GenericMergeOrRebaseAction(commandType string, command string) error {
	err := self.runSkipEditorCommand(self.GenericMergeOrRebaseActionCmdObj(commandType, command))
	if err != nil {
		if !strings.Contains(err.Error(), "no rebase in progress") {
			return err
		}
		self.Log.Warn(err)
	}

	// sometimes we need to do a sequence of things in a rebase but the user needs to
	// fix merge conflicts along the way. When this happens we queue up the next step
	// so that after the next successful rebase continue we can continue from where we left off
	if commandType == "rebase" && command == "continue" && self.onSuccessfulContinue != nil {
		f := self.onSuccessfulContinue
		self.onSuccessfulContinue = nil
		return f()
	}
	if command == "abort" {
		self.onSuccessfulContinue = nil
	}
	return nil
}

func (self *RebaseCommands) runSkipEditorCommand(cmdObj oscommands.ICmdObj) error {
	lazyGitPath := oscommands.GetLazygitPath()
	return cmdObj.
		AddEnvVars(
			daemon.DaemonKindEnvKey+"="+string(daemon.ExitImmediately),
			"GIT_EDITOR="+lazyGitPath,
			"EDITOR="+lazyGitPath,
			"VISUAL="+lazyGitPath,
		).
		Run()
}

// DiscardOldFileChanges discards changes to a file from an old commit
func (self *RebaseCommands) DiscardOldFileChanges(commits []*models.Commit, commitIndex int, fileName string) error {
	if err := self.BeginInteractiveRebaseForCommit(commits, commitIndex); err != nil {
		return err
	}

	// check if file exists in previous commit (this command returns an error if the file doesn't exist)
	if err := self.cmd.New("git cat-file -e HEAD^:" + self.cmd.Quote(fileName)).Run(); err != nil {
		if err := self.os.Remove(fileName); err != nil {
			return err
		}
		if err := self.workingTree.StageFile(fileName); err != nil {
			return err
		}
	} else if err := self.workingTree.CheckoutFile("HEAD^", fileName); err != nil {
		return err
	}

	// amend the commit
	err := self.commit.AmendHead()
	if err != nil {
		return err
	}

	// continue
	return self.ContinueRebase()
}

// CherryPickCommits begins an interactive rebase with the given shas being cherry picked onto HEAD
func (self *RebaseCommands) CherryPickCommits(commits []*models.Commit) error {
	todoLines := self.BuildTodoLinesSingleAction(commits, "pick")

	return self.PrepareInteractiveRebaseCommand("HEAD", todoLines, false).Run()
}

func (self *RebaseCommands) buildTodo(todoLines []TodoLine) string {
	lines := slices.Map(todoLines, func(todoLine TodoLine) string {
		return todoLine.ToString()
	})

	return strings.Join(slices.Reverse(lines), "")
}

func (self *RebaseCommands) BuildTodoLines(commits []*models.Commit, f func(*models.Commit, int) string) []TodoLine {
	return slices.MapWithIndex(commits, func(commit *models.Commit, i int) TodoLine {
		return TodoLine{Action: f(commit, i), Commit: commit}
	})
}

func (self *RebaseCommands) BuildTodoLinesSingleAction(commits []*models.Commit, action string) []TodoLine {
	return self.BuildTodoLines(commits, func(commit *models.Commit, i int) string {
		return action
	})
}

type TodoLine struct {
	Action string
	Commit *models.Commit
}

func (self *TodoLine) ToString() string {
	return self.Action + " " + self.Commit.Sha + " " + self.Commit.Name + "\n"
}
