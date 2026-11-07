package git_commands

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/spf13/afero"
	"github.com/stefanhaller/git-todo-parser/todo"
)

type StatusCommands struct {
	*GitCommon
}

func NewStatusCommands(
	gitCommon *GitCommon,
) *StatusCommands {
	return &StatusCommands{
		GitCommon: gitCommon,
	}
}

func (self *StatusCommands) WorkingTreeState() models.WorkingTreeState {
	result := models.WorkingTreeState{}
	result.Rebasing, _ = self.IsInRebase()
	result.Merging, _ = self.IsInMergeState()
	result.CherryPicking, _ = self.IsInCherryPick()
	result.Reverting, _ = self.IsInRevert()
	result.ApplyingPatches, _ = self.IsApplyingPatches()
	result.SquashMerging, _ = self.IsInSquashMerge()
	return result
}

func (self *StatusCommands) IsBareRepo() bool {
	return self.repoPaths.isBareRepo
}

func (self *StatusCommands) IsInRebase() (bool, error) {
	exists, err := self.os.FileExists(filepath.Join(self.repoPaths.WorktreeGitDirPath(), "rebase-merge"))
	if err == nil && exists {
		return true, nil
	}

	// "git am" uses the rebase-apply directory too
	exists, err = self.os.FileExists(filepath.Join(self.repoPaths.WorktreeGitDirPath(), "rebase-apply"))
	if err != nil || !exists {
		return false, err
	}
	isApplyingPatches, err := self.IsApplyingPatches()
	return !isApplyingPatches, err
}

// IsApplyingPatches tells whether a "git am" is in progress
func (self *StatusCommands) IsApplyingPatches() (bool, error) {
	return self.os.FileExists(filepath.Join(self.repoPaths.WorktreeGitDirPath(), "rebase-apply", "applying"))
}

// IsInMergeState states whether we are still mid-merge
func (self *StatusCommands) IsInMergeState() (bool, error) {
	return self.os.FileExists(filepath.Join(self.repoPaths.WorktreeGitDirPath(), "MERGE_HEAD"))
}

// IsInSquashMerge tells whether "git merge --squash" has staged changes that
// haven't been committed yet. "git merge --squash" writes SQUASH_MSG, a message
// for committing them, but no MERGE_HEAD. "git commit" uses the message and
// removes the file; so do "git reset" without paths, switching branches, and
// "git stash".
func (self *StatusCommands) IsInSquashMerge() (bool, error) {
	return self.os.FileExists(filepath.Join(self.repoPaths.WorktreeGitDirPath(), "SQUASH_MSG"))
}

func (self *StatusCommands) IsInCherryPick() (bool, error) {
	hasCherryPickHead, err := self.hasCherryPickHead()
	if err != nil || hasCherryPickHead {
		return hasCherryPickHead, err
	}
	return self.sequencerCommand() == todo.Pick, nil
}

func (self *StatusCommands) hasCherryPickHead() (bool, error) {
	gitDir := self.repoPaths.WorktreeGitDirPath()
	exists, err := self.os.FileExists(filepath.Join(gitDir, "CHERRY_PICK_HEAD"))
	if err != nil || !exists {
		return exists, err
	}
	cherryPickHead, err := os.ReadFile(filepath.Join(gitDir, "CHERRY_PICK_HEAD"))
	if err != nil {
		return false, err
	}
	return !cherryPickHeadBelongsToRebase(gitDir, strings.TrimSpace(string(cherryPickHead)), self.config.GetCoreCommentChar()), nil
}

// cherryPickHeadBelongsToRebase tells whether CHERRY_PICK_HEAD, whose value is
// cherryPickHead, belongs to a pick of the rebase that is in progress in
// gitDir, rather than to a cherry-pick. A rebase writes CHERRY_PICK_HEAD when it
// has applied the changes of a pick, and removes it when it has committed them.
// So it is still there if the rebase stopped because the pick conflicted; git
// records the pick in rebase-merge/stopped-sha then. It is also still there if
// git couldn't commit the changes, for example because it couldn't take the
// index lock. Git exits without writing stopped-sha in that case, but the pick
// is the last line of rebase-merge/done.
func cherryPickHeadBelongsToRebase(gitDir string, cherryPickHead string, commentChar byte) bool {
	// The hashes in these files may be abbreviated
	isPrefixOfCherryPickHead := func(hash string) bool {
		return hash != "" && strings.HasPrefix(cherryPickHead, hash)
	}

	stoppedSha, err := os.ReadFile(filepath.Join(gitDir, "rebase-merge", "stopped-sha"))
	if err == nil && isPrefixOfCherryPickHead(strings.TrimSpace(string(stoppedSha))) {
		return true
	}

	doneTodos, err := utils.ReadRebaseTodoFile(filepath.Join(gitDir, "rebase-merge", "done"), commentChar)
	return err == nil && len(doneTodos) > 0 && isPrefixOfCherryPickHead(doneTodos[len(doneTodos)-1].Commit)
}

// rebaseStoppedBeforeCommittingPick tells whether the rebase that is in
// progress in gitDir has applied the changes of a pick, but exited before it
// could commit them (see cherryPickHeadBelongsToRebase). "git rebase
// --continue" refuses to commit the changes in that case, because git didn't
// write the message file that it writes when a pick conflicts.
func rebaseStoppedBeforeCommittingPick(gitDir string, commentChar byte) bool {
	cherryPickHead, err := os.ReadFile(filepath.Join(gitDir, "CHERRY_PICK_HEAD"))
	if err != nil {
		return false
	}
	if _, err := os.Stat(filepath.Join(gitDir, "rebase-merge", "message")); err == nil {
		return false
	}
	return cherryPickHeadBelongsToRebase(gitDir, strings.TrimSpace(string(cherryPickHead)), commentChar)
}

func (self *StatusCommands) IsInRevert() (bool, error) {
	hasRevertHead, err := self.os.FileExists(filepath.Join(self.repoPaths.WorktreeGitDirPath(), "REVERT_HEAD"))
	if err != nil || hasRevertHead {
		return hasRevertHead, err
	}
	return self.sequencerCommand() == todo.Revert, nil
}

// sequencerCommand returns the command of the first todo of a cherry-pick or
// revert of several commits that is in progress, or 0 if there is none. Git
// keeps the todos of such an operation in .git/sequencer. If it stops before it
// starts to apply a commit (for example because the commit would overwrite an
// untracked file), it doesn't write CHERRY_PICK_HEAD or REVERT_HEAD. Like git
// status, we still consider the cherry-pick or revert to be in progress then.
func (self *StatusCommands) sequencerCommand() todo.TodoCommand {
	todos, err := utils.ReadRebaseTodoFile(
		filepath.Join(self.repoPaths.WorktreeGitDirPath(), "sequencer", "todo"), self.config.GetCoreCommentChar())
	if err != nil || len(todos) == 0 {
		return 0
	}
	return todos[0].Command
}

// RefsSnapshot returns a string fingerprint of the current state of local
// branches and HEAD. Comparing two snapshots byte-for-byte tells us whether
// any local ref or HEAD has moved since the last snapshot.
func (self *StatusCommands) RefsSnapshot() (string, error) {
	t := time.Now()
	defer func() { self.Log.Infof("RefsSnapshot took %s", time.Since(t)) }()

	refsArgs := NewGitCmd("for-each-ref").
		Arg("--format=%(objectname) %(refname)").
		Arg("refs/heads").
		ToArgv()
	refs, err := self.cmd.New(refsArgs).DontLog().RunWithOutput()
	if err != nil {
		return "", err
	}

	head, err := self.headSnapshot()
	if err != nil {
		return "", err
	}

	return refs + head, nil
}

// headSnapshot returns a fingerprint of HEAD that distinguishes "detached at
// commit X" from "on a branch that points at X". The commit hash alone can't
// tell those apart, which matters at the end of a rebase: HEAD reattaches to
// the branch without the hash changing, and we'd otherwise miss that refresh.
//
// We read .git/HEAD directly rather than shelling out: it's faster (no child
// process) and its content is exactly the symref-or-hash distinction we want
// ("ref: refs/heads/foo" when attached, the raw hash when detached). The
// reftable backend, however, doesn't keep a real .git/HEAD — it writes a fixed
// stub ("ref: refs/heads/.invalid") that never reflects the actual HEAD. When
// we see that stub (or the file is missing/unreadable) we fall back to
// porcelain commands, which are backend-agnostic.
func (self *StatusCommands) headSnapshot() (string, error) {
	headPath := filepath.Join(self.repoPaths.WorktreeGitDirPath(), "HEAD")
	if content, err := afero.ReadFile(self.Fs, headPath); err == nil {
		head := strings.TrimSpace(string(content))
		if head != "" && head != "ref: refs/heads/.invalid" {
			return head, nil
		}
	}

	// symbolic-ref gives the branch when HEAD is attached and fails when it's
	// detached, in which case rev-parse gives the commit HEAD points at.
	symbolicRefArgs := NewGitCmd("symbolic-ref").Arg("HEAD").ToArgv()
	if symref, err := self.cmd.New(symbolicRefArgs).DontLog().RunWithOutput(); err == nil {
		return strings.TrimSpace(symref), nil
	}

	revParseArgs := NewGitCmd("rev-parse").Arg("HEAD").ToArgv()
	head, err := self.cmd.New(revParseArgs).DontLog().RunWithOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(head), nil
}

// Full ref (e.g. "refs/heads/mybranch") of the branch that is currently
// being rebased, or empty string when we're not in a rebase
func (self *StatusCommands) BranchBeingRebased() string {
	for _, dir := range []string{"rebase-merge", "rebase-apply"} {
		if bytesContent, err := os.ReadFile(filepath.Join(self.repoPaths.WorktreeGitDirPath(), dir, "head-name")); err == nil {
			return strings.TrimSpace(string(bytesContent))
		}
	}
	return ""
}
