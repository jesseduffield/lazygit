package git_commands

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fsmiamoto/git-todo-parser/todo"
	"github.com/go-errors/errors"
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

func (self *RebaseCommands) RewordCommit(commits []*models.Commit, index int, summary string, description string) error {
	if models.IsHeadCommit(commits, index) {
		// we've selected the top commit so no rebase is required
		return self.commit.RewordLastCommit(summary, description)
	}

	err := self.BeginInteractiveRebaseForCommit(commits, index, false)
	if err != nil {
		return err
	}

	// now the selected commit should be our head so we'll amend it with the new message
	err = self.commit.RewordLastCommit(summary, description)
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

func (self *RebaseCommands) AddCommitCoAuthor(commits []*models.Commit, index int, value string) error {
	return self.GenericAmend(commits, index, func() error {
		return self.commit.AddCoAuthor(commits[index].Sha, value)
	})
}

func (self *RebaseCommands) GenericAmend(commits []*models.Commit, index int, f func() error) error {
	if models.IsHeadCommit(commits, index) {
		// we've selected the top commit so no rebase is required
		return f()
	}

	err := self.BeginInteractiveRebaseForCommit(commits, index, false)
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

func (self *RebaseCommands) MoveCommitsDown(commits []*models.Commit, startIdx int, endIdx int) error {
	baseShaOrRoot := getBaseShaOrRoot(commits, endIdx+2)

	shas := lo.Map(commits[startIdx:endIdx+1], func(commit *models.Commit, _ int) string {
		return commit.Sha
	})

	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseShaOrRoot:  baseShaOrRoot,
		instruction:    daemon.NewMoveTodosDownInstruction(shas),
		overrideEditor: true,
	}).Run()
}

func (self *RebaseCommands) MoveCommitsUp(commits []*models.Commit, startIdx int, endIdx int) error {
	baseShaOrRoot := getBaseShaOrRoot(commits, endIdx+1)

	shas := lo.Map(commits[startIdx:endIdx+1], func(commit *models.Commit, _ int) string {
		return commit.Sha
	})

	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseShaOrRoot:  baseShaOrRoot,
		instruction:    daemon.NewMoveTodosUpInstruction(shas),
		overrideEditor: true,
	}).Run()
}

func (self *RebaseCommands) InteractiveRebase(commits []*models.Commit, startIdx int, endIdx int, action todo.TodoCommand) error {
	baseIndex := endIdx + 1
	if action == todo.Squash || action == todo.Fixup {
		baseIndex++
	}

	baseShaOrRoot := getBaseShaOrRoot(commits, baseIndex)

	changes := lo.Map(commits[startIdx:endIdx+1], func(commit *models.Commit, _ int) daemon.ChangeTodoAction {
		return daemon.ChangeTodoAction{
			Sha:       commit.Sha,
			NewAction: action,
		}
	})

	self.os.LogCommand(logTodoChanges(changes), false)

	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseShaOrRoot:  baseShaOrRoot,
		overrideEditor: true,
		instruction:    daemon.NewChangeTodoActionsInstruction(changes),
	}).Run()
}

func (self *RebaseCommands) EditRebase(branchRef string) error {
	msg := utils.ResolvePlaceholderString(
		self.Tr.Log.EditRebase,
		map[string]string{
			"ref": branchRef,
		},
	)
	self.os.LogCommand(msg, false)
	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseShaOrRoot: branchRef,
		instruction:   daemon.NewInsertBreakInstruction(),
	}).Run()
}

func (self *RebaseCommands) EditRebaseFromBaseCommit(targetBranchName string, baseCommit string) error {
	msg := utils.ResolvePlaceholderString(
		self.Tr.Log.EditRebaseFromBaseCommit,
		map[string]string{
			"baseCommit":       baseCommit,
			"targetBranchName": targetBranchName,
		},
	)
	self.os.LogCommand(msg, false)
	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseShaOrRoot: baseCommit,
		onto:          targetBranchName,
		instruction:   daemon.NewInsertBreakInstruction(),
	}).Run()
}

func logTodoChanges(changes []daemon.ChangeTodoAction) string {
	changeTodoStr := strings.Join(lo.Map(changes, func(c daemon.ChangeTodoAction, _ int) string {
		return fmt.Sprintf("%s:%s", c.Sha, c.NewAction)
	}), "\n")
	return fmt.Sprintf("Changing TODO actions:\n%s", changeTodoStr)
}

type PrepareInteractiveRebaseCommandOpts struct {
	baseShaOrRoot              string
	onto                       string
	instruction                daemon.Instruction
	overrideEditor             bool
	keepCommitsThatBecomeEmpty bool
}

// PrepareInteractiveRebaseCommand returns the cmd for an interactive rebase
// we tell git to run lazygit to edit the todo list, and we pass the client
// lazygit a todo string to write to the todo file
func (self *RebaseCommands) PrepareInteractiveRebaseCommand(opts PrepareInteractiveRebaseCommandOpts) oscommands.ICmdObj {
	ex := oscommands.GetLazygitPath()

	cmdArgs := NewGitCmd("rebase").
		Arg("--interactive").
		Arg("--autostash").
		Arg("--keep-empty").
		ArgIf(opts.keepCommitsThatBecomeEmpty && self.version.IsAtLeast(2, 26, 0), "--empty=keep").
		Arg("--no-autosquash").
		ArgIf(self.version.IsAtLeast(2, 22, 0), "--rebase-merges").
		ArgIf(opts.onto != "", "--onto", opts.onto).
		Arg(opts.baseShaOrRoot).
		ToArgv()

	debug := "FALSE"
	if self.Debug {
		debug = "TRUE"
	}

	self.Log.WithField("command", cmdArgs).Debug("RunCommand")

	cmdObj := self.cmd.New(cmdArgs)

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
	cmdArgs := NewGitCmd("rev-parse").Arg("--verify", "HEAD").ToArgv()
	fixupSha, err := self.cmd.New(cmdArgs).RunWithOutput()
	if err != nil {
		return err
	}

	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseShaOrRoot:  getBaseShaOrRoot(commits, commitIndex+1),
		overrideEditor: true,
		instruction:    daemon.NewMoveFixupCommitDownInstruction(commit.Sha, fixupSha),
	}).Run()
}

func todoFromCommit(commit *models.Commit) utils.Todo {
	if commit.Action == todo.UpdateRef {
		return utils.Todo{Ref: commit.Name, Action: commit.Action}
	} else {
		return utils.Todo{Sha: commit.Sha, Action: commit.Action}
	}
}

// Sets the action for the given commits in the git-rebase-todo file
func (self *RebaseCommands) EditRebaseTodo(commits []*models.Commit, action todo.TodoCommand) error {
	commitsWithAction := lo.Map(commits, func(commit *models.Commit, _ int) utils.TodoChange {
		return utils.TodoChange{
			Sha:       commit.Sha,
			OldAction: commit.Action,
			NewAction: action,
		}
	})

	return utils.EditRebaseTodo(
		filepath.Join(self.repoPaths.WorktreeGitDirPath(), "rebase-merge/git-rebase-todo"),
		commitsWithAction,
		self.config.GetCoreCommentChar(),
	)
}

func (self *RebaseCommands) DeleteUpdateRefTodos(commits []*models.Commit) error {
	todosToDelete := lo.Map(commits, func(commit *models.Commit, _ int) utils.Todo {
		return todoFromCommit(commit)
	})

	return utils.DeleteTodos(
		filepath.Join(self.repoPaths.WorktreeGitDirPath(), "rebase-merge/git-rebase-todo"),
		todosToDelete,
		self.config.GetCoreCommentChar(),
	)
}

func (self *RebaseCommands) MoveTodosDown(commits []*models.Commit) error {
	fileName := filepath.Join(self.repoPaths.WorktreeGitDirPath(), "rebase-merge/git-rebase-todo")
	todosToMove := lo.Map(commits, func(commit *models.Commit, _ int) utils.Todo {
		return todoFromCommit(commit)
	})

	return utils.MoveTodosDown(fileName, todosToMove, self.config.GetCoreCommentChar())
}

func (self *RebaseCommands) MoveTodosUp(commits []*models.Commit) error {
	fileName := filepath.Join(self.repoPaths.WorktreeGitDirPath(), "rebase-merge/git-rebase-todo")
	todosToMove := lo.Map(commits, func(commit *models.Commit, _ int) utils.Todo {
		return todoFromCommit(commit)
	})

	return utils.MoveTodosUp(fileName, todosToMove, self.config.GetCoreCommentChar())
}

// SquashAllAboveFixupCommits squashes all fixup! commits above the given one
func (self *RebaseCommands) SquashAllAboveFixupCommits(commit *models.Commit) error {
	shaOrRoot := commit.Sha + "^"
	if commit.IsFirstCommit() {
		shaOrRoot = "--root"
	}

	cmdArgs := NewGitCmd("rebase").
		Arg("--interactive", "--rebase-merges", "--autostash", "--autosquash", shaOrRoot).
		ToArgv()

	return self.runSkipEditorCommand(self.cmd.New(cmdArgs))
}

// BeginInteractiveRebaseForCommit starts an interactive rebase to edit the current
// commit and pick all others. After this you'll want to call `self.ContinueRebase()
func (self *RebaseCommands) BeginInteractiveRebaseForCommit(
	commits []*models.Commit, commitIndex int, keepCommitsThatBecomeEmpty bool,
) error {
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
		baseShaOrRoot:              getBaseShaOrRoot(commits, commitIndex+1),
		overrideEditor:             true,
		keepCommitsThatBecomeEmpty: keepCommitsThatBecomeEmpty,
		instruction:                daemon.NewChangeTodoActionsInstruction(changes),
	}).Run()
}

// RebaseBranch interactive rebases onto a branch
func (self *RebaseCommands) RebaseBranch(branchName string) error {
	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{baseShaOrRoot: branchName}).Run()
}

func (self *RebaseCommands) RebaseBranchFromBaseCommit(targetBranchName string, baseCommit string) error {
	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseShaOrRoot: baseCommit,
		onto:          targetBranchName,
	}).Run()
}

func (self *RebaseCommands) GenericMergeOrRebaseActionCmdObj(commandType string, command string) oscommands.ICmdObj {
	cmdArgs := NewGitCmd(commandType).Arg("--" + command).ToArgv()

	return self.cmd.New(cmdArgs)
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
func (self *RebaseCommands) DiscardOldFileChanges(commits []*models.Commit, commitIndex int, filePaths []string) error {
	if err := self.BeginInteractiveRebaseForCommit(commits, commitIndex, false); err != nil {
		return err
	}

	for _, filePath := range filePaths {
		// check if file exists in previous commit (this command returns an error if the file doesn't exist)
		cmdArgs := NewGitCmd("cat-file").Arg("-e", "HEAD^:"+filePath).ToArgv()

		if err := self.cmd.New(cmdArgs).Run(); err != nil {
			if err := self.os.Remove(filePath); err != nil {
				return err
			}
			if err := self.workingTree.StageFile(filePath); err != nil {
				return err
			}
		} else if err := self.workingTree.CheckoutFile("HEAD^", filePath); err != nil {
			return err
		}
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
	msg := utils.ResolvePlaceholderString(
		self.Tr.Log.CherryPickCommits,
		map[string]string{
			"commitLines": strings.Join(commitLines, "\n"),
		},
	)
	self.os.LogCommand(msg, false)

	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseShaOrRoot: "HEAD",
		instruction:   daemon.NewCherryPickCommitsInstruction(commits),
	}).Run()
}

// CherryPickCommitsDuringRebase simply prepends the given commits to the existing git-rebase-todo file
func (self *RebaseCommands) CherryPickCommitsDuringRebase(commits []*models.Commit) error {
	todoLines := lo.Map(commits, func(commit *models.Commit, _ int) daemon.TodoLine {
		return daemon.TodoLine{
			Action: "pick",
			Commit: commit,
		}
	})

	todo := daemon.TodoLinesToString(todoLines)
	filePath := filepath.Join(self.repoPaths.worktreeGitDirPath, "rebase-merge/git-rebase-todo")
	return utils.PrependStrToTodoFile(filePath, []byte(todo))
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
