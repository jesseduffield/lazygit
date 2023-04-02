package git_commands

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fsmiamoto/git-todo-parser/todo"
	"github.com/go-errors/errors"
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/app/daemon"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/utils"
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
	if models.IsHeadCommit(commits, index) {
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

	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseShaOrRoot: sha,
		todoLines:     todo,
	}), nil
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
	// not appending to original slice so that we don't mutate it
	orderedCommits := append([]*models.Commit{}, commits[0:index]...)
	orderedCommits = append(orderedCommits, commits[index+1], commits[index])

	todoLines := self.BuildTodoLinesSingleAction(orderedCommits, "pick")

	baseShaOrRoot := getBaseShaOrRoot(commits, index+2)

	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseShaOrRoot:  baseShaOrRoot,
		todoLines:      todoLines,
		overrideEditor: true,
	}).Run()
}

func (self *RebaseCommands) InteractiveRebase(commits []*models.Commit, index int, action string) error {
	todo, sha, err := self.BuildSingleActionTodo(commits, index, action)
	if err != nil {
		return err
	}

	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseShaOrRoot:  sha,
		todoLines:      todo,
		overrideEditor: true,
	}).Run()
}

func (self *RebaseCommands) InteractiveRebaseBreakAfter(commits []*models.Commit, index int) error {
	todo, sha, err := self.BuildSingleActionTodo(commits, index-1, "pick")
	if err != nil {
		return err
	}

	todo = append(todo, TodoLine{Action: "break", Commit: nil})
	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseShaOrRoot:  sha,
		todoLines:      todo,
		overrideEditor: true,
	}).Run()
}

func (self *RebaseCommands) EditRebase(branchRef string) error {
	commands := []TodoLine{{Action: "break"}}
	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseShaOrRoot: branchRef,
		todoLines:     commands,
		prepend:       true,
	}).Run()
}

type PrepareInteractiveRebaseCommandOpts struct {
	baseShaOrRoot  string
	todoLines      []TodoLine
	overrideEditor bool
	prepend        bool
}

// PrepareInteractiveRebaseCommand returns the cmd for an interactive rebase
// we tell git to run lazygit to edit the todo list, and we pass the client
// lazygit a todo string to write to the todo file
func (self *RebaseCommands) PrepareInteractiveRebaseCommand(opts PrepareInteractiveRebaseCommandOpts) oscommands.ICmdObj {
	todo := self.buildTodo(opts.todoLines)
	ex := oscommands.GetLazygitPath()
	prependLines := ""
	if opts.prepend {
		prependLines = "TRUE"
	}

	debug := "FALSE"
	if self.Debug {
		debug = "TRUE"
	}

	cmdStr := fmt.Sprintf("git rebase --interactive --autostash --keep-empty --empty=keep --no-autosquash %s", opts.baseShaOrRoot)
	self.Log.WithField("command", cmdStr).Debug("RunCommand")

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
		daemon.PrependLinesEnvKey+"="+prependLines,
		"DEBUG="+debug,
		"LANG=en_US.UTF-8",   // Force using EN as language
		"LC_ALL=en_US.UTF-8", // Force using EN as language
		"GIT_SEQUENCE_EDITOR="+gitSequenceEditor,
	)

	if opts.overrideEditor {
		cmdObj.AddEnvVars("GIT_EDITOR=" + ex)
	}

	return cmdObj
}

// produces TodoLines where every commit is picked (or dropped for merge commits) except for the commit at the given index, which
// will have the given action applied to it.
func (self *RebaseCommands) BuildSingleActionTodo(commits []*models.Commit, actionIndex int, action string) ([]TodoLine, string, error) {
	baseIndex := actionIndex + 1

	if action == "squash" || action == "fixup" {
		baseIndex++
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

	baseShaOrRoot := getBaseShaOrRoot(commits, baseIndex)

	return todoLines, baseShaOrRoot, nil
}

// AmendTo amends the given commit with whatever files are staged
func (self *RebaseCommands) AmendTo(commit *models.Commit) error {
	if err := self.commit.CreateFixupCommit(commit.Sha); err != nil {
		return err
	}

	return self.SquashAllAboveFixupCommits(commit)
}

// EditRebaseTodo sets the action for a given rebase commit in the git-rebase-todo file
func (self *RebaseCommands) EditRebaseTodo(commit *models.Commit, action todo.TodoCommand) error {
	fileName := filepath.Join(self.dotGitDir, "rebase-merge/git-rebase-todo")
	todos, err := utils.ReadRebaseTodoFile(fileName)
	if err != nil {
		return err
	}

	for i := range todos {
		t := &todos[i]
		// Comparing just the sha is not enough; we need to compare both the
		// action and the sha, as the sha could appear multiple times (e.g. in a
		// pick and later in a merge)
		if t.Command == commit.Action && t.Commit == commit.Sha {
			t.Command = action
			return utils.WriteRebaseTodoFile(fileName, todos)
		}
	}

	// Should never get here
	return fmt.Errorf("Todo %s not found in git-rebase-todo", commit.Sha)
}

// MoveTodoDown moves a rebase todo item down by one position
func (self *RebaseCommands) MoveTodoDown(commit *models.Commit) error {
	fileName := filepath.Join(self.dotGitDir, "rebase-merge/git-rebase-todo")
	return utils.MoveTodoDown(fileName, commit.Sha, commit.Action)
}

// MoveTodoDown moves a rebase todo item down by one position
func (self *RebaseCommands) MoveTodoUp(commit *models.Commit) error {
	fileName := filepath.Join(self.dotGitDir, "rebase-merge/git-rebase-todo")
	return utils.MoveTodoUp(fileName, commit.Sha, commit.Action)
}

// SquashAllAboveFixupCommits squashes all fixup! commits above the given one
func (self *RebaseCommands) SquashAllAboveFixupCommits(commit *models.Commit) error {
	shaOrRoot := commit.Sha + "^"
	if commit.IsFirstCommit() {
		shaOrRoot = "--root"
	}

	return self.runSkipEditorCommand(
		self.cmd.New(
			fmt.Sprintf(
				"git rebase --interactive --rebase-merges --autostash --autosquash %s",
				shaOrRoot,
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

	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseShaOrRoot:  sha,
		todoLines:      todo,
		overrideEditor: true,
	}).Run()
}

// RebaseBranch interactive rebases onto a branch
func (self *RebaseCommands) RebaseBranch(branchName string) error {
	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{baseShaOrRoot: branchName}).Run()
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
			"GIT_SEQUENCE_EDITOR="+lazyGitPath,
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

	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseShaOrRoot: "HEAD",
		todoLines:     todoLines,
	}).Run()
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
	if self.Action == "break" {
		return self.Action + "\n"
	} else {
		return self.Action + " " + self.Commit.Sha + " " + self.Commit.Name + "\n"
	}
}

// we can't start an interactive rebase from the first commit without passing the
// '--root' arg
func getBaseShaOrRoot(commits []*models.Commit, index int) string {
	// We assume that the commits slice contains the initial commit of the repo.
	// Technically this assumption could prove false, but it's unlikely you'll
	// be starting a rebase from 300 commits ago (which is the original commit limit
	// at time of writing)
	if index < len(commits) {
		return commits[index].Sha
	} else {
		return "--root"
	}
}
