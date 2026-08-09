package git_commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/app/daemon"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/stefanhaller/git-todo-parser/todo"
)

type PatchCommands struct {
	*GitCommon
	rebase *RebaseCommands
	commit *CommitCommands
	status *StatusCommands
	stash  *StashCommands

	PatchBuilder *patch.PatchBuilder

	// lastBuiltTreeGeneration is the patch builder generation the diff trees were last
	// materialized for, so EnsureCustomPatchDiffTrees rebuilds them only when the patch
	// actually changed (not on every render / navigation).
	lastBuiltTreeGeneration int
}

func NewPatchCommands(
	gitCommon *GitCommon,
	rebase *RebaseCommands,
	commit *CommitCommands,
	status *StatusCommands,
	stash *StashCommands,
	patchBuilder *patch.PatchBuilder,
) *PatchCommands {
	return &PatchCommands{
		GitCommon:    gitCommon,
		rebase:       rebase,
		commit:       commit,
		status:       status,
		stash:        stash,
		PatchBuilder: patchBuilder,
	}
}

// EnsureCustomPatchDiffTrees materializes the custom patch's diff trees if they're stale —
// i.e. the patch changed since they were last built. Called before rendering the secondary
// pane, so a mere re-render (e.g. navigating commits, which doesn't change the patch) reuses
// the existing trees, while a toggle/removal/whole-file change rebuilds them. This is what
// keeps the trees current across every path that mutates the patch.
func (self *PatchCommands) EnsureCustomPatchDiffTrees() error {
	if self.PatchBuilder.Generation() == self.lastBuiltTreeGeneration {
		return nil
	}
	if err := self.WriteCustomPatchDiffTrees(); err != nil {
		return err
	}
	self.lastBuiltTreeGeneration = self.PatchBuilder.Generation()
	return nil
}

// WriteCustomPatchDiffTrees materializes the custom patch under the patch builder's temp
// dir as two file trees — a/ holds each patched file's "from"-side content, b/ holds that
// content with the patch applied — so the patch can be re-diffed with `git diff --no-index`
// and rendered through any pager (see DiffCommands.CustomPatchDiffCmdObj). This is what lets
// a partial, in-memory custom patch be shown the same way as any other diff; the in-memory
// aggregated patch on its own could only be fed to a stdin pager, never an external diff
// tool. The dirs are named a/b so that, with `--no-prefix`, the diff's paths come out as the
// real repo-relative paths (git's conventional a//b/ prefixes).
//
// Called whenever the patch's contents change. The temp dir's lifetime is owned by the
// patch builder (created on Start, removed on Reset), so there's nothing to clean up here
// beyond wiping the trees before rebuilding them.
func (self *PatchCommands) WriteCustomPatchDiffTrees() error {
	dir := self.PatchBuilder.TempDir()
	if dir == "" {
		return nil
	}

	aDir := filepath.Join(dir, "a")
	bDir := filepath.Join(dir, "b")
	for _, d := range []string{aDir, bDir} {
		if err := os.RemoveAll(d); err != nil {
			return err
		}
		if err := os.MkdirAll(d, 0o700); err != nil {
			return err
		}
	}

	for _, filename := range self.PatchBuilder.ActiveFilenames() {
		// The "from"-side content; the patch applied below turns b/ into the "after" side
		// while a/ stays the "before".
		content, err := self.commit.ShowFileContentCmdObj(self.PatchBuilder.From, filename).RunWithOutput()
		added := err != nil // absent on the "from" side — a file the patch adds

		// a/ always holds the "before" content — empty for an added file, so the rendered
		// diff still pairs it with b/ and shows the real a//b/ paths (rather than git's
		// directory-comparison "added in b" form, which mangles the header).
		beforeContent := content
		if added {
			beforeContent = ""
		}
		if err := self.os.CreateFileWithContent(filepath.Join(aDir, filename), beforeContent); err != nil {
			return err
		}
		// b/ is seeded only for existing files (which the patch modifies in place); an added
		// file is left absent so the patch — rendered as a /dev/null creation — creates it.
		if !added {
			if err := self.os.CreateFileWithContent(filepath.Join(bDir, filename), content); err != nil {
				return err
			}
		}
	}

	// Render added files as /dev/null creations (the natural form), not as diffs against an
	// empty file: the latter (--- a/file) would make the atomic `git apply` expect them to
	// already exist in b/, where they don't.
	patchText := self.PatchBuilder.PatchToApply(false, false)
	if strings.TrimSpace(patchText) == "" {
		// An empty patch leaves a/ and b/ identical (or empty), so the diff is empty.
		return nil
	}
	patchFilePath, err := self.SaveTemporaryPatch(patchText)
	if err != nil {
		return err
	}
	return self.cmd.New(NewGitCmd("apply").Arg(patchFilePath).Dir(bDir).ToArgv()).Run()
}

type ApplyPatchOpts struct {
	ThreeWay bool
	Cached   bool
	Index    bool
	Reverse  bool
}

func (self *PatchCommands) ApplyCustomPatch(reverse bool, turnAddedFilesIntoDiffAgainstEmptyFile bool) error {
	patch := self.PatchBuilder.PatchToApply(reverse, turnAddedFilesIntoDiffAgainstEmptyFile)

	return self.ApplyPatch(patch, ApplyPatchOpts{
		Index:    true,
		ThreeWay: true,
		Reverse:  reverse,
	})
}

func (self *PatchCommands) ApplyPatch(patch string, opts ApplyPatchOpts) error {
	filepath, err := self.SaveTemporaryPatch(patch)
	if err != nil {
		return err
	}

	return self.applyPatchFile(filepath, opts)
}

func (self *PatchCommands) applyPatchFile(filepath string, opts ApplyPatchOpts) error {
	cmdArgs := NewGitCmd("apply").
		ArgIf(opts.ThreeWay, "--3way").
		ArgIf(opts.Cached, "--cached").
		ArgIf(opts.Index, "--index").
		ArgIf(opts.Reverse, "--reverse").
		Arg(filepath).
		ToArgv()

	return self.cmd.New(cmdArgs).Run()
}

func (self *PatchCommands) SaveTemporaryPatch(patch string) (string, error) {
	filepath := filepath.Join(self.os.GetTempDir(), self.repoPaths.RepoName(), time.Now().Format("Jan _2 15.04.05.000000000")+".patch")
	self.Log.Infof("saving temporary patch to %s", filepath)
	if err := self.os.CreateFileWithContent(filepath, patch); err != nil {
		return "", err
	}
	return filepath, nil
}

// DeletePatchesFromCommit applies a patch in reverse for a commit
func (self *PatchCommands) DeletePatchesFromCommit(commits []*models.Commit, commitIndex int) error {
	if err := self.rebase.BeginInteractiveRebaseForCommit(commits, commitIndex, false); err != nil {
		return err
	}

	// apply each patch in reverse
	if err := self.ApplyCustomPatch(true, true); err != nil {
		_ = self.rebase.AbortRebase()
		return err
	}

	// time to amend the selected commit
	if err := self.commit.AmendHead(); err != nil {
		return err
	}

	self.rebase.onSuccessfulContinue = func() error {
		self.PatchBuilder.Reset()
		return nil
	}

	// continue
	return self.rebase.ContinueRebase()
}

func (self *PatchCommands) MovePatchToSelectedCommit(commits []*models.Commit, sourceCommitIdx int, destinationCommitIdx int) error {
	if sourceCommitIdx < destinationCommitIdx {
		// Passing true for keepCommitsThatBecomeEmpty: if the moved-from
		// commit becomes empty, we want to keep it, mainly for consistency with
		// moving the patch to a *later* commit, which behaves the same.
		if err := self.rebase.BeginInteractiveRebaseForCommit(commits, destinationCommitIdx, true); err != nil {
			return err
		}

		// apply each patch forward
		if err := self.ApplyCustomPatch(false, false); err != nil {
			// Don't abort the rebase here; this might cause conflicts, so give
			// the user a chance to resolve them
			return err
		}

		// amend the destination commit
		if err := self.commit.AmendHead(); err != nil {
			return err
		}

		self.rebase.onSuccessfulContinue = func() error {
			self.PatchBuilder.Reset()
			return nil
		}

		// continue
		return self.rebase.ContinueRebase()
	}

	if len(commits)-1 < sourceCommitIdx {
		return errors.New("index outside of range of commits")
	}

	// we can make this GPG thing possible it just means we need to do this in two parts:
	// one where we handle the possibility of a credential request, and the other
	// where we continue the rebase
	if self.config.NeedsGpgSubprocessForCommit() {
		return errors.New(self.Tr.DisabledForGPG)
	}

	baseIndex := sourceCommitIdx + 1

	changes := []daemon.ChangeTodoAction{
		{Hash: commits[sourceCommitIdx].Hash(), NewAction: todo.Edit},
		{Hash: commits[destinationCommitIdx].Hash(), NewAction: todo.Edit},
	}
	self.os.LogCommand(logTodoChanges(changes), false)

	err := self.rebase.PrepareInteractiveRebaseCommand(PrepareInteractiveRebaseCommandOpts{
		baseHashOrRoot: getBaseHashOrRoot(commits, baseIndex),
		overrideEditor: true,
		instruction:    daemon.NewChangeTodoActionsInstruction(changes),
	}).Run()
	if err != nil {
		return err
	}

	// apply each patch in reverse
	if err := self.ApplyCustomPatch(true, true); err != nil {
		_ = self.rebase.AbortRebase()
		return err
	}

	// amend the source commit
	if err := self.commit.AmendHead(); err != nil {
		return err
	}

	patch, err := self.diffHeadAgainstCommit(commits[sourceCommitIdx])
	if err != nil {
		_ = self.rebase.AbortRebase()
		return err
	}

	if self.rebase.onSuccessfulContinue != nil {
		return errors.New("You are midway through another rebase operation. Please abort to start again")
	}

	self.rebase.onSuccessfulContinue = func() error {
		// now we should be up to the destination, so let's apply forward these patches to that.
		// ideally we would ensure we're on the right commit but I'm not sure if that check is necessary
		if err := self.ApplyPatch(patch, ApplyPatchOpts{Index: true, ThreeWay: true}); err != nil {
			// Don't abort the rebase here; this might cause conflicts, so give
			// the user a chance to resolve them
			return err
		}

		// amend the destination commit
		if err := self.commit.AmendHead(); err != nil {
			return err
		}

		self.rebase.onSuccessfulContinue = func() error {
			self.PatchBuilder.Reset()
			return nil
		}

		return self.rebase.ContinueRebase()
	}

	return self.rebase.ContinueRebase()
}

func (self *PatchCommands) MovePatchIntoIndex(commits []*models.Commit, commitIdx int, stash bool) error {
	if stash {
		if err := self.stash.Push(fmt.Sprintf(self.Tr.AutoStashForMovingPatchToIndex, commits[commitIdx].ShortHash())); err != nil {
			return err
		}
	}

	if err := self.rebase.BeginInteractiveRebaseForCommit(commits, commitIdx, false); err != nil {
		return err
	}

	if err := self.ApplyCustomPatch(true, true); err != nil {
		if self.status.WorkingTreeState().Rebasing {
			_ = self.rebase.AbortRebase()
		}
		return err
	}

	// amend the commit
	if err := self.commit.AmendHead(); err != nil {
		return err
	}

	patch, err := self.diffHeadAgainstCommit(commits[commitIdx])
	if err != nil {
		_ = self.rebase.AbortRebase()
		return err
	}

	if self.rebase.onSuccessfulContinue != nil {
		return errors.New("You are midway through another rebase operation. Please abort to start again")
	}

	self.rebase.onSuccessfulContinue = func() error {
		// add patches to index
		if err := self.ApplyPatch(patch, ApplyPatchOpts{Index: true, ThreeWay: true}); err != nil {
			if self.status.WorkingTreeState().Rebasing {
				_ = self.rebase.AbortRebase()
			}
			return err
		}

		if stash {
			if err := self.stash.Pop(0); err != nil {
				return err
			}
		}

		self.PatchBuilder.Reset()
		return nil
	}

	return self.rebase.ContinueRebase()
}

func (self *PatchCommands) PullPatchIntoNewCommit(
	commits []*models.Commit,
	commitIdx int,
	commitSummary string,
	commitDescription string,
) error {
	if err := self.rebase.BeginInteractiveRebaseForCommit(commits, commitIdx, false); err != nil {
		return err
	}

	if err := self.ApplyCustomPatch(true, true); err != nil {
		_ = self.rebase.AbortRebase()
		return err
	}

	// amend the commit
	if err := self.commit.AmendHead(); err != nil {
		return err
	}

	patch, err := self.diffHeadAgainstCommit(commits[commitIdx])
	if err != nil {
		_ = self.rebase.AbortRebase()
		return err
	}

	if err := self.ApplyPatch(patch, ApplyPatchOpts{Index: true, ThreeWay: true}); err != nil {
		_ = self.rebase.AbortRebase()
		return err
	}

	if err := self.commit.CommitCmdObj(commitSummary, commitDescription, false).Run(); err != nil {
		return err
	}

	if self.rebase.onSuccessfulContinue != nil {
		return errors.New("You are midway through another rebase operation. Please abort to start again")
	}

	self.PatchBuilder.Reset()
	return self.rebase.ContinueRebase()
}

func (self *PatchCommands) PullPatchIntoNewCommitBefore(
	commits []*models.Commit,
	commitIdx int,
	commitSummary string,
	commitDescription string,
) error {
	if err := self.rebase.BeginInteractiveRebaseForCommit(commits, commitIdx+1, true); err != nil {
		return err
	}

	if err := self.ApplyCustomPatch(false, false); err != nil {
		_ = self.rebase.AbortRebase()
		return err
	}

	if err := self.commit.CommitCmdObj(commitSummary, commitDescription, false).Run(); err != nil {
		return err
	}

	self.PatchBuilder.Reset()
	return self.rebase.ContinueRebase()
}

// We have just applied a patch in reverse to discard it from a commit; if we
// now try to apply the patch again to move it to a later commit, or to the
// index, then this would conflict "with itself" in case the patch contained
// only some lines of a range of adjacent added lines. To solve this, we
// get the diff of HEAD and the original commit and then apply that.
func (self *PatchCommands) diffHeadAgainstCommit(commit *models.Commit) (string, error) {
	cmdArgs := NewGitCmd("diff").
		Config("diff.noprefix=false").
		Arg("--no-ext-diff", "--no-color").
		Arg("HEAD.." + commit.Hash()).
		ToArgv()

	return self.cmd.New(cmdArgs).RunWithOutput()
}
