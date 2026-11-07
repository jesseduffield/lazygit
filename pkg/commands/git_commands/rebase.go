package git_commands

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/app/daemon"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
	"github.com/stefanhaller/git-todo-parser/todo"
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
	// This check is currently unreachable (handled in LocalCommitsController.reword),
	// but kept as a safeguard in case this method is used elsewhere.
	if self.config.NeedsGpgSubprocessForCommit() {
		return errors.New(self.Tr.DisabledForGPG)
	}

	err := self.BeginInteractiveRebaseForCommit(commits, index, false)
	if err != nil {
		return err
	}

	// now the selected commit should be our head so we'll amend it with the new message
	err = self.commit.RewordLastCommit(summary, description).Run()
	if err != nil {
		return err
	}

	return self.ContinueRebase()
}

func (self *RebaseCommands) RewordCommitInEditor(commits []*models.Commit, index int) (*oscommands.CmdObj, error) {
	changes := []daemon.ChangeTodoAction{{
		Hash:      commits[index].Hash(),
		NewAction: todo.Reword,
	}}
	self.os.LogCommand(logTodoChanges(changes), false)

	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseHashOrRoot: getBaseHashOrRoot(commits, index+1),
		instruction:    daemon.NewChangeTodoActionsInstruction(changes),
	}), nil
}

func (self *RebaseCommands) ResetCommitAuthor(commits []*models.Commit, start, end int) error {
	return self.GenericAmend(commits, start, end, func(_ *models.Commit) error {
		return self.commit.ResetAuthor()
	})
}

func (self *RebaseCommands) SetCommitAuthor(commits []*models.Commit, start, end int, value string) error {
	return self.GenericAmend(commits, start, end, func(_ *models.Commit) error {
		return self.commit.SetAuthor(value)
	})
}

func (self *RebaseCommands) AddCommitCoAuthor(commits []*models.Commit, start, end int, value string) error {
	return self.GenericAmend(commits, start, end, func(commit *models.Commit) error {
		return self.commit.AddCoAuthor(commit.Hash(), value)
	})
}

func (self *RebaseCommands) GenericAmend(commits []*models.Commit, start, end int, f func(commit *models.Commit) error) error {
	if start == end && models.IsHeadCommit(commits, start) {
		// we've selected the top commit so no rebase is required
		return f(commits[start])
	}

	err := self.BeginInteractiveRebaseForCommitRange(commits, start, end, false)
	if err != nil {
		return err
	}

	for commitIndex := end; commitIndex >= start; commitIndex-- {
		err = f(commits[commitIndex])
		if err != nil {
			return err
		}

		if err := self.ContinueRebase(); err != nil {
			return err
		}
	}

	return nil
}

func (self *RebaseCommands) MoveCommitsDown(commits []*models.Commit, startIdx int, endIdx int) error {
	return self.MoveCommits(commits, startIdx, endIdx, 1)
}

func (self *RebaseCommands) MoveCommitsUp(commits []*models.Commit, startIdx int, endIdx int) error {
	return self.MoveCommits(commits, startIdx, endIdx, -1)
}

func (self *RebaseCommands) MoveCommits(commits []*models.Commit, startIdx int, endIdx int, offset int) error {
	baseHashOrRoot := getBaseHashOrRoot(commits, endIdx+max(offset, 0)+1)

	hashes := lo.Map(commits[startIdx:endIdx+1], func(commit *models.Commit, _ int) string {
		return commit.Hash()
	})

	var instruction daemon.Instruction
	if offset > 0 {
		instruction = daemon.NewMoveTodosDownInstruction(hashes, offset)
	} else {
		instruction = daemon.NewMoveTodosUpInstruction(hashes, -offset)
	}

	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseHashOrRoot: baseHashOrRoot,
		instruction:    instruction,
		overrideEditor: true,
	}).Run()
}

func (self *RebaseCommands) InteractiveRebase(commits []*models.Commit, startIdx int, endIdx int, action todo.TodoCommand, flag string) error {
	baseIndex := endIdx + 1
	if action == todo.Squash || action == todo.Fixup {
		baseIndex++
	}

	baseHashOrRoot := getBaseHashOrRoot(commits, baseIndex)

	changes := lo.FilterMap(commits[startIdx:endIdx+1], func(commit *models.Commit, _ int) (daemon.ChangeTodoAction, bool) {
		return daemon.ChangeTodoAction{
			Hash:      commit.Hash(),
			NewAction: action,
			Flag:      flag,
		}, !commit.IsMerge()
	})

	self.os.LogCommand(logTodoChanges(changes), false)

	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseHashOrRoot: baseHashOrRoot,
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
		baseHashOrRoot: branchRef,
		instruction:    daemon.NewInsertBreakInstruction(),
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
		baseHashOrRoot: baseCommit,
		onto:           targetBranchName,
		instruction:    daemon.NewInsertBreakInstruction(),
	}).Run()
}

func logTodoChanges(changes []daemon.ChangeTodoAction) string {
	changeTodoStr := strings.Join(lo.Map(changes, func(c daemon.ChangeTodoAction, _ int) string {
		return fmt.Sprintf("%s:%s", c.Hash, c.NewAction)
	}), "\n")
	return fmt.Sprintf("Changing TODO actions:\n%s", changeTodoStr)
}

type PrepareInteractiveRebaseCommandOpts struct {
	baseHashOrRoot             string
	onto                       string
	instruction                daemon.Instruction
	overrideEditor             bool
	keepCommitsThatBecomeEmpty bool
}

// PrepareInteractiveRebaseCommand returns the cmd for an interactive rebase
// we tell git to run lazygit to edit the todo list, and we pass the client
// lazygit instructions what to do with the todo file
func (self *RebaseCommands) PrepareInteractiveRebaseCommand(opts PrepareInteractiveRebaseCommandOpts) *oscommands.CmdObj {
	ex := oscommands.GetLazygitPath()

	cmdArgs := NewGitCmd("rebase").
		Arg("--interactive").
		Arg("--autostash").
		Arg("--keep-empty").
		ArgIf(opts.keepCommitsThatBecomeEmpty, "--empty=keep").
		Arg("--no-autosquash").
		Arg("--rebase-merges").
		ArgIf(opts.onto != "", "--onto", opts.onto).
		Arg(opts.baseHashOrRoot).
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
		cmdObj.AddEnvVars(daemon.ToEnvVars(daemon.NewRemoveUpdateRefsForCopiedBranchInstruction())...)
	}

	cmdObj.AddEnvVars(
		"DEBUG="+debug,
		"LANG=C",        // Force using English language
		"LC_ALL=C",      // Force using English language
		"LC_MESSAGES=C", // Force using English language
		"GIT_SEQUENCE_EDITOR="+gitSequenceEditor,
	)

	if opts.overrideEditor {
		cmdObj.AddEnvVars("GIT_EDITOR=" + ex)
	}

	return self.resumeAfterLockError(cmdObj)
}

// GitRebaseEditTodo runs "git rebase --edit-todo", saving the given todosFileContent to the file
func (self *RebaseCommands) GitRebaseEditTodo(todosFileContent []byte) error {
	ex := oscommands.GetLazygitPath()

	cmdArgs := NewGitCmd("rebase").
		Arg("--edit-todo").
		ToArgv()

	debug := "FALSE"
	if self.Debug {
		debug = "TRUE"
	}

	self.Log.WithField("command", cmdArgs).Debug("RunCommand")

	cmdObj := self.cmd.New(cmdArgs)

	cmdObj.AddEnvVars(daemon.ToEnvVars(daemon.NewWriteRebaseTodoInstruction(todosFileContent))...)

	cmdObj.AddEnvVars(
		"DEBUG="+debug,
		"LANG=C",        // Force using English language
		"LC_ALL=C",      // Force using English language
		"LC_MESSAGES=C", // Force using English language
		"GIT_EDITOR="+ex,
		"GIT_SEQUENCE_EDITOR="+ex,
	)

	return cmdObj.Run()
}

func (self *RebaseCommands) getHashOfLastCommitMade() (string, error) {
	cmdArgs := NewGitCmd("rev-parse").Arg("--verify", "HEAD").ToArgv()
	return self.cmd.New(cmdArgs).RunWithOutput()
}

// AmendTo amends the given commit with whatever files are staged
func (self *RebaseCommands) AmendTo(commits []*models.Commit, commitIndex int) error {
	commit := commits[commitIndex]

	if err := self.commit.CreateFixupCommit(commit.Hash()); err != nil {
		return err
	}

	fixupHash, err := self.getHashOfLastCommitMade()
	if err != nil {
		return err
	}

	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseHashOrRoot: getBaseHashOrRoot(commits, commitIndex+1),
		overrideEditor: true,
		instruction:    daemon.NewMoveFixupCommitDownInstruction(commit.Hash(), fixupHash, true),
	}).Run()
}

func (self *RebaseCommands) MoveFixupCommitDown(commits []*models.Commit, targetCommitIndex int) error {
	fixupHash, err := self.getHashOfLastCommitMade()
	if err != nil {
		return err
	}

	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseHashOrRoot: getBaseHashOrRoot(commits, targetCommitIndex+1),
		overrideEditor: true,
		instruction:    daemon.NewMoveFixupCommitDownInstruction(commits[targetCommitIndex].Hash(), fixupHash, false),
	}).Run()
}

func todoFromCommit(commit *models.Commit) utils.Todo {
	if commit.Action == todo.UpdateRef {
		return utils.Todo{Ref: commit.Name}
	}
	return utils.Todo{Hash: commit.Hash()}
}

// Sets the action for the given commits in the git-rebase-todo file
func (self *RebaseCommands) EditRebaseTodo(commits []*models.Commit, action todo.TodoCommand, flag string) error {
	commitsWithAction := lo.Map(commits, func(commit *models.Commit, _ int) utils.TodoChange {
		return utils.TodoChange{
			Hash:      commit.Hash(),
			NewAction: action,
			Flag:      flag,
		}
	})

	return utils.EditRebaseTodo(
		filepath.Join(self.repoPaths.WorktreeGitDirPath(), "rebase-merge/git-rebase-todo"),
		commitsWithAction,
		self.config.GetCoreCommentChar(),
	)
}

// Makes the todo for the given commit in the git-rebase-todo file use a
// different commit instead
func (self *RebaseCommands) ChangeTodoCommit(hash string, newHash string) error {
	return utils.ChangeTodoCommit(
		filepath.Join(self.repoPaths.WorktreeGitDirPath(), "rebase-merge/git-rebase-todo"),
		hash,
		newHash,
		self.config.GetCoreCommentChar(),
	)
}

func (self *RebaseCommands) DeleteUpdateRefTodos(commits []*models.Commit) error {
	todosToDelete := lo.Map(commits, func(commit *models.Commit, _ int) utils.Todo {
		return todoFromCommit(commit)
	})

	todosFileContent, err := utils.DeleteTodos(
		filepath.Join(self.repoPaths.WorktreeGitDirPath(), "rebase-merge/git-rebase-todo"),
		todosToDelete,
		self.config.GetCoreCommentChar(),
	)
	if err != nil {
		return err
	}

	return self.GitRebaseEditTodo(todosFileContent)
}

func (self *RebaseCommands) MoveTodosDown(commits []*models.Commit) error {
	return self.MoveTodos(commits, 1)
}

func (self *RebaseCommands) MoveTodosUp(commits []*models.Commit) error {
	return self.MoveTodos(commits, -1)
}

func (self *RebaseCommands) MoveTodos(commits []*models.Commit, offset int) error {
	fileName := filepath.Join(self.repoPaths.WorktreeGitDirPath(), "rebase-merge/git-rebase-todo")
	todosToMove := lo.Map(commits, func(commit *models.Commit, _ int) utils.Todo {
		return todoFromCommit(commit)
	})

	return utils.MoveTodos(fileName, todosToMove, true, offset, self.config.GetCoreCommentChar())
}

// SquashAllAboveFixupCommits squashes all fixup! commits above the given one
func (self *RebaseCommands) SquashAllAboveFixupCommits(commit *models.Commit) error {
	hashOrRoot := commit.Hash() + "^"
	if commit.IsFirstCommit() {
		hashOrRoot = "--root"
	}

	cmdArgs := NewGitCmd("rebase").
		Arg("--interactive", "--rebase-merges", "--autostash", "--autosquash", hashOrRoot).
		ToArgv()

	return self.resumeAfterLockError(skipEditor(self.cmd.New(cmdArgs))).Run()
}

// BeginInteractiveRebaseForCommit starts an interactive rebase to edit the current
// commit and pick all others. After this you'll want to call `self.ContinueRebase()
func (self *RebaseCommands) BeginInteractiveRebaseForCommit(
	commits []*models.Commit, commitIndex int, keepCommitsThatBecomeEmpty bool,
) error {
	if commitIndex < len(commits) && commits[commitIndex].IsMerge() {
		if self.config.NeedsGpgSubprocessForCommit() {
			return errors.New(self.Tr.DisabledForGPG)
		}

		return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
			baseHashOrRoot:             getBaseHashOrRoot(commits, commitIndex),
			instruction:                daemon.NewInsertBreakInstruction(),
			keepCommitsThatBecomeEmpty: keepCommitsThatBecomeEmpty,
		}).Run()
	}

	return self.BeginInteractiveRebaseForCommitRange(commits, commitIndex, commitIndex, keepCommitsThatBecomeEmpty)
}

func (self *RebaseCommands) BeginInteractiveRebaseForCommitRange(
	commits []*models.Commit, start, end int, keepCommitsThatBecomeEmpty bool,
) error {
	if len(commits)-1 < end {
		return errors.New("index outside of range of commits")
	}

	// we can make this GPG thing possible it just means we need to do this in two parts:
	// one where we handle the possibility of a credential request, and the other
	// where we continue the rebase
	if self.config.NeedsGpgSubprocessForCommit() {
		return errors.New(self.Tr.DisabledForGPG)
	}

	changes := make([]daemon.ChangeTodoAction, 0, end-start)
	for commitIndex := end; commitIndex >= start; commitIndex-- {
		changes = append(changes, daemon.ChangeTodoAction{
			Hash:      commits[commitIndex].Hash(),
			NewAction: todo.Edit,
		})
	}
	self.os.LogCommand(logTodoChanges(changes), false)

	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseHashOrRoot:             getBaseHashOrRoot(commits, end+1),
		overrideEditor:             true,
		keepCommitsThatBecomeEmpty: keepCommitsThatBecomeEmpty,
		instruction:                daemon.NewChangeTodoActionsInstruction(changes),
	}).Run()
}

// RebaseBranch interactive rebases onto a branch
func (self *RebaseCommands) RebaseBranch(branchName string) error {
	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{baseHashOrRoot: branchName}).Run()
}

func (self *RebaseCommands) RebaseBranchFromBaseCommit(targetBranchName string, baseCommit string) error {
	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseHashOrRoot: baseCommit,
		onto:           targetBranchName,
	}).Run()
}

// GenericMergeOrRebaseActionCmdObj returns the command that continues, skips or
// aborts the given operation; command is "continue", "skip" or "abort".
func (self *RebaseCommands) GenericMergeOrRebaseActionCmdObj(operation models.EffectiveWorkingTreeState, command string) *oscommands.CmdObj {
	cmdArgs := NewGitCmd(operation.CommandName()).Arg("--" + command).ToArgv()

	return self.cmd.New(cmdArgs)
}

func (self *RebaseCommands) ContinueRebase() error {
	return self.GenericMergeOrRebaseAction(models.WORKING_TREE_STATE_REBASING, "continue")
}

func (self *RebaseCommands) AbortRebase() error {
	return self.GenericMergeOrRebaseAction(models.WORKING_TREE_STATE_REBASING, "abort")
}

// GenericMergeOrRebaseAction continues, skips or aborts the given operation;
// command is "continue", "skip" or "abort". If this makes a commit, we skip the
// editor.
func (self *RebaseCommands) GenericMergeOrRebaseAction(operation models.EffectiveWorkingTreeState, command string) error {
	undoRepeatingTodo := func() {}
	if command == "continue" &&
		(operation == models.WORKING_TREE_STATE_CHERRY_PICKING || operation == models.WORKING_TREE_STATE_REVERTING) {
		var err error
		undoRepeatingTodo, err = repeatCurrentSequencerTodo(
			self.cmd, self.repoPaths.WorktreeGitDirPath(), self.config.GetCoreCommentChar())
		if err != nil {
			return err
		}
	}

	if operation == models.WORKING_TREE_STATE_REBASING && command == "continue" {
		if err := self.commitPickThatRebaseStoppedAt(); err != nil {
			return err
		}
	}

	cmdObj := skipEditor(self.GenericMergeOrRebaseActionCmdObj(operation, command))
	if command == "continue" || command == "skip" {
		switch operation {
		case models.WORKING_TREE_STATE_REBASING:
			cmdObj = self.resumeAfterLockError(cmdObj)
		case models.WORKING_TREE_STATE_CHERRY_PICKING, models.WORKING_TREE_STATE_REVERTING:
			cmdObj = retryContinueOrSkipAfterLockError(cmdObj,
				filepath.Join(self.repoPaths.WorktreeGitDirPath(), "sequencer", "todo"))
		case models.WORKING_TREE_STATE_APPLYING_PATCHES:
			cmdObj = retryContinueOrSkipAfterLockError(cmdObj,
				filepath.Join(self.repoPaths.WorktreeGitDirPath(), "rebase-apply", "next"))
		case models.WORKING_TREE_STATE_MERGING, models.WORKING_TREE_STATE_NONE:
		}
	}

	err := cmdObj.Run()
	if err != nil {
		undoRepeatingTodo()
		if !strings.Contains(err.Error(), "no rebase in progress") {
			return err
		}
		self.Log.Warn(err)
	}

	// sometimes we need to do a sequence of things in a rebase but the user needs to
	// fix merge conflicts along the way. When this happens we queue up the next step
	// so that after the next successful rebase continue we can continue from where we left off
	if operation == models.WORKING_TREE_STATE_REBASING && command == "continue" && self.onSuccessfulContinue != nil {
		f := self.onSuccessfulContinue
		self.onSuccessfulContinue = nil
		return f()
	}
	if command == "abort" {
		self.onSuccessfulContinue = nil
	}
	return nil
}

// commitPickThatRebaseStoppedAt commits the changes of a pick if the rebase
// stopped after it applied them, but before it committed them (see
// rebaseStoppedBeforeCommittingPick). "git rebase --continue" refuses to commit
// them in that case. "git cherry-pick --continue" commits them with the author
// and message of the pick, like the rebase would have.
func (self *RebaseCommands) commitPickThatRebaseStoppedAt() error {
	if !rebaseStoppedBeforeCommittingPick(self.repoPaths.WorktreeGitDirPath(), self.config.GetCoreCommentChar()) {
		return nil
	}

	return skipEditor(self.GenericMergeOrRebaseActionCmdObj(models.WORKING_TREE_STATE_CHERRY_PICKING, "continue")).Run()
}

// resumeAfterLockError makes cmdObj, a command that starts or continues a
// rebase, resume the rebase if a lock error stops it partway through; see
// resumeRebaseAfterLockError.
func (self *RebaseCommands) resumeAfterLockError(cmdObj *oscommands.CmdObj) *oscommands.CmdObj {
	return resumeRebaseAfterLockError(cmdObj, self.repoPaths.WorktreeGitDirPath(), self.config.GetCoreCommentChar(),
		func() *oscommands.CmdObj {
			return skipEditor(self.GenericMergeOrRebaseActionCmdObj(models.WORKING_TREE_STATE_REBASING, "continue"))
		})
}

// skipEditor makes git use lazygit as the editor for the given command; lazygit
// then exits right away, so git uses the file as it is.
func skipEditor(cmdObj *oscommands.CmdObj) *oscommands.CmdObj {
	instruction := daemon.NewExitImmediatelyInstruction()
	lazyGitPath := oscommands.GetLazygitPath()
	return cmdObj.
		AddEnvVars(
			"GIT_EDITOR="+lazyGitPath,
			"GIT_SEQUENCE_EDITOR="+lazyGitPath,
			"EDITOR="+lazyGitPath,
			"VISUAL="+lazyGitPath,
		).
		AddEnvVars(daemon.ToEnvVars(instruction)...)
}

// DiscardOldFileChanges discards changes to a file from an old commit
func (self *RebaseCommands) DiscardOldFileChanges(commits []*models.Commit, commitIndex int, filePaths []string) error {
	if err := self.BeginInteractiveRebaseForCommit(commits, commitIndex, false); err != nil {
		return err
	}

	for _, filePath := range filePaths {
		doesFileExistInPreviousCommit := false
		if commitIndex < len(commits)-1 {
			// check if file exists in previous commit (this command returns an empty string if the file doesn't exist)
			cmdArgs := NewGitCmd("ls-tree").Arg("--name-only", "HEAD^", "--", filePath).ToArgv()
			output, err := self.cmd.New(cmdArgs).DontLog().RunWithOutput()
			if err != nil {
				return err
			}
			doesFileExistInPreviousCommit = strings.TrimRight(output, "\n") == filePath
		}
		if !doesFileExistInPreviousCommit {
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

// CherryPickCommits begins an interactive rebase with the given hashes being cherry picked onto HEAD
func (self *RebaseCommands) CherryPickCommits(commits []*models.Commit) error {
	hasMergeCommit := lo.SomeBy(commits, func(c *models.Commit) bool { return c.IsMerge() })
	cmdArgs := NewGitCmd("cherry-pick").
		Arg("--allow-empty").
		ArgIf(self.version.IsAtLeast(2, 45, 0), "--empty=keep", "--keep-redundant-commits").
		ArgIf(hasMergeCommit, "-m1").
		Arg(lo.Reverse(lo.Map(commits, func(c *models.Commit, _ int) string { return c.Hash() }))...).
		ToArgv()

	return abortSequenceAfterLockError(self.cmd.New(cmdArgs), self.repoPaths.WorktreeGitDirPath(), len(commits),
		func() *oscommands.CmdObj {
			return self.cmd.New(NewGitCmd("cherry-pick").Arg("--abort").ToArgv())
		}).Run()
}

func (self *RebaseCommands) DropMergeCommit(commits []*models.Commit, commitIndex int) error {
	return self.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseHashOrRoot: getBaseHashOrRoot(commits, commitIndex+1),
		instruction:    daemon.NewDropMergeCommitInstruction(commits[commitIndex].Hash()),
	}).Run()
}

// we can't start an interactive rebase from the first commit without passing the
// '--root' arg
func getBaseHashOrRoot(commits []*models.Commit, index int) string {
	// We assume that the commits slice contains the initial commit of the repo.
	// Technically this assumption could prove false, but it's unlikely you'll
	// be starting a rebase from 300 commits ago (which is the original commit limit
	// at time of writing)
	if index < len(commits) {
		return commits[index].Hash()
	}

	return "--root"
}
