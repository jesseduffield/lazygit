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
	"github.com/samber/lo"
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
	changes := []daemon.ChangeTodoAction{{
		Sha:       commits[index].Sha,
		NewAction: todo.Reword,
	}}
	self.os.LogCommand(logTodoChanges(changes), false)

	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseShaOrRoot: getBaseShaOrRoot(commits, index+1),
		instruction:   daemon.NewChangeTodoActionsInstruction(changes),
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
	if models.IsHeadCommit(commits, index) {
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
	baseShaOrRoot := getBaseShaOrRoot(commits, index+2)

	sha := commits[index].Sha

	self.os.LogCommand(fmt.Sprintf("Moving TODO down: %s", utils.ShortSha(sha)), false)

	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseShaOrRoot:  baseShaOrRoot,
		instruction:    daemon.NewMoveTodoDownInstruction(sha),
		overrideEditor: true,
	}).Run()
}

func (self *RebaseCommands) MoveCommitUp(commits []*models.Commit, index int) error {
	baseShaOrRoot := getBaseShaOrRoot(commits, index+1)

	sha := commits[index].Sha

	self.os.LogCommand(fmt.Sprintf("Moving TODO up: %s", utils.ShortSha(sha)), false)

	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseShaOrRoot:  baseShaOrRoot,
		instruction:    daemon.NewMoveTodoUpInstruction(sha),
		overrideEditor: true,
	}).Run()
}

func (self *RebaseCommands) InteractiveRebase(commits []*models.Commit, index int, action todo.TodoCommand) error {
	baseIndex := index + 1
	if action == todo.Squash || action == todo.Fixup {
		baseIndex++
	}

	baseShaOrRoot := getBaseShaOrRoot(commits, baseIndex)

	changes := []daemon.ChangeTodoAction{{
		Sha:       commits[index].Sha,
		NewAction: action,
	}}
	self.os.LogCommand(logTodoChanges(changes), false)

	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseShaOrRoot:  baseShaOrRoot,
		overrideEditor: true,
		instruction:    daemon.NewChangeTodoActionsInstruction(changes),
	}).Run()
}

func (self *RebaseCommands) EditRebase(branchRef string) error {
	self.os.LogCommand(fmt.Sprintf("Beginning interactive rebase at '%s'", branchRef), false)
	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseShaOrRoot: branchRef,
		instruction:   daemon.NewInsertBreakInstruction(),
	}).Run()
}

func logTodoChanges(changes []daemon.ChangeTodoAction) string {
	changeTodoStr := strings.Join(slices.Map(changes, func(c daemon.ChangeTodoAction) string {
		return fmt.Sprintf("%s:%s", c.Sha, c.NewAction)
	}), "\n")
	return fmt.Sprintf("Changing TODO actions: %s", changeTodoStr)
}

type PrepareInteractiveRebaseCommandOpts struct {
	baseShaOrRoot  string
	instruction    daemon.Instruction
	overrideEditor bool
}

// PrepareInteractiveRebaseCommand returns the cmd for an interactive rebase
// we tell git to run lazygit to edit the todo list, and we pass the client
// lazygit a todo string to write to the todo file
func (self *RebaseCommands) PrepareInteractiveRebaseCommand(opts PrepareInteractiveRebaseCommandOpts) oscommands.ICmdObj {
	ex := oscommands.GetLazygitPath()

	debug := "FALSE"
	if self.Debug {
		debug = "TRUE"
	}

	emptyArg := " --empty=keep"
	if self.version.IsOlderThan(2, 26, 0) {
		emptyArg = ""
	}

	rebaseMergesArg := " --rebase-merges"
	if self.version.IsOlderThan(2, 22, 0) {
		rebaseMergesArg = ""
	}

	cmdStr := fmt.Sprintf("git rebase --interactive --autostash --keep-empty%s --no-autosquash%s %s",
		emptyArg, rebaseMergesArg, opts.baseShaOrRoot)
	self.Log.WithField("command", cmdStr).Debug("RunCommand")

	cmdObj := self.cmd.New(cmdStr)

	gitSequenceEditor := ex

	if opts.instruction != nil {
		cmdObj.AddEnvVars(daemon.ToEnvVars(opts.instruction)...)
	} else {
		gitSequenceEditor = "true"
	}

	cmdObj.AddEnvVars(
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

// AmendTo amends the given commit with whatever files are staged
func (self *RebaseCommands) AmendTo(commits []*models.Commit, commitIndex int) error {
	commit := commits[commitIndex]

	if err := self.commit.CreateFixupCommit(commit.Sha); err != nil {
		return err
	}

	// Get the sha of the commit we just created
	fixupSha, err := self.cmd.New("git rev-parse --verify HEAD").RunWithOutput()
	if err != nil {
		return err
	}

	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseShaOrRoot:  getBaseShaOrRoot(commits, commitIndex+1),
		overrideEditor: true,
		instruction:    daemon.NewMoveFixupCommitDownInstruction(commit.Sha, fixupSha),
	}).Run()
}

// EditRebaseTodo sets the action for a given rebase commit in the git-rebase-todo file
func (self *RebaseCommands) EditRebaseTodo(commit *models.Commit, action todo.TodoCommand) error {
	return utils.EditRebaseTodo(
		filepath.Join(self.dotGitDir, "rebase-merge/git-rebase-todo"), commit.Sha, commit.Action, action)
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

	changes := []daemon.ChangeTodoAction{{
		Sha:       commits[commitIndex].Sha,
		NewAction: todo.Edit,
	}}
	self.os.LogCommand(logTodoChanges(changes), false)

	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseShaOrRoot:  getBaseShaOrRoot(commits, commitIndex+1),
		overrideEditor: true,
		instruction:    daemon.NewChangeTodoActionsInstruction(changes),
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
	instruction := daemon.NewExitImmediatelyInstruction()
	lazyGitPath := oscommands.GetLazygitPath()
	return cmdObj.
		AddEnvVars(
			"GIT_EDITOR="+lazyGitPath,
			"GIT_SEQUENCE_EDITOR="+lazyGitPath,
			"EDITOR="+lazyGitPath,
			"VISUAL="+lazyGitPath,
		).
		AddEnvVars(daemon.ToEnvVars(instruction)...).
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
	commitLines := lo.Map(commits, func(commit *models.Commit, _ int) string {
		return fmt.Sprintf("%s %s", utils.ShortSha(commit.Sha), commit.Name)
	})
	self.os.LogCommand(fmt.Sprintf("Cherry-picking commits:\n%s", strings.Join(commitLines, "\n")), false)

	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseShaOrRoot: "HEAD",
		instruction:   daemon.NewCherryPickCommitsInstruction(commits),
	}).Run()
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
