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
	"github.com/samber/lo"
	"github.com/stefanhaller/git-todo-parser/todo"
)

type PatchCommands struct {
	*GitCommon
	rebase *RebaseCommands
	commit *CommitCommands
	status *StatusCommands
	stash  *StashCommands

	PatchBuilder *patch.PatchBuilder

	// The version of the patch the diff trees were last written for, so that they are
	// written again when, and only when, the patch has changed since.
	treesWrittenForGeneration int
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

// EnsureCustomPatchDiffTrees writes the custom patch's diff trees if what is there no
// longer describes the patch. Call it before rendering the patch, which is often — every
// time the panel showing it re-renders — while the patch itself changes rarely.
func (self *PatchCommands) EnsureCustomPatchDiffTrees() error {
	if self.PatchBuilder.Generation() == self.treesWrittenForGeneration {
		return nil
	}
	if err := self.WriteCustomPatchDiffTrees(); err != nil {
		return err
	}
	self.treesWrittenForGeneration = self.PatchBuilder.Generation()
	return nil
}

// WriteCustomPatchDiffTrees materializes the custom patch as two file trees under the
// directory the patch builder keeps for it: `a` holds each of the patch's files as it is
// before the patch, `b` as it is after. Diffing those two trees against each other
// (DiffCommands.CustomPatchDiffCmdObj) turns the patch into a diff of real files, which
// can then be rendered exactly as any other diff is — through a diff renderer of any
// kind, and with git's own idea of how much context to show.
//
// The trees are named a and b so that the diff's paths, with git's own prefixes
// suppressed, come out reading like the a/ and b/ of an ordinary diff, over the real
// repo-relative paths.
func (self *PatchCommands) WriteCustomPatchDiffTrees() error {
	dir := self.PatchBuilder.TempDir()
	if dir == "" {
		return nil
	}

	before := filepath.Join(dir, "a")
	after := filepath.Join(dir, "b")
	for _, tree := range []string{before, after} {
		if err := os.RemoveAll(tree); err != nil {
			return err
		}
		if err := os.MkdirAll(tree, 0o700); err != nil {
			return err
		}
	}

	for _, file := range self.PatchBuilder.FilesInPatch() {
		content, err := self.commit.ShowFileContentCmdObj(self.PatchBuilder.From, file.ContentPath).RunWithOutput()
		// A file the patch adds has no content on the before side, so git has nothing to
		// show for it.
		added := err != nil

		// The before side holds an added file as an empty file rather than not at all, so
		// that the diff pairs the two sides up and states the file's real path, instead of
		// reporting a file that only one of the trees has.
		if err := self.os.CreateFileWithContent(filepath.Join(before, file.Path),
			lo.Ternary(added, "", content)); err != nil {
			return err
		}
		// The after side is seeded with the same content, for the patch to change; a file
		// the patch adds is left absent, for the patch to create.
		if !added {
			if err := self.os.CreateFileWithContent(filepath.Join(after, file.Path), content); err != nil {
				return err
			}
		}
	}

	// Write added files as creations rather than as diffs against an empty file: the
	// patch is applied in one go, so a file it expects to be there already would make
	// the whole of it fail.
	patchText := self.PatchBuilder.PatchToApply(false, false)
	if strings.TrimSpace(patchText) == "" {
		// Nothing in the patch, so the two trees are alike and the diff is empty.
		return nil
	}
	patchFilePath, err := self.SaveTemporaryPatch(patchText)
	if err != nil {
		return err
	}
	return self.cmd.New(NewGitCmd("apply").
		NoLineEndingConversion().
		Arg(patchFilePath).
		Dir(after).
		ToArgv()).Run()
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

		if err := self.pickCopyWithPatchInParent(commits[sourceCommitIdx]); err != nil {
			_ = self.rebase.AbortRebase()
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

// When moving the patch to an earlier commit, we apply it to that commit, and
// the rebase then picks the commit that the patch was built from. Picking that
// commit unchanged would rely on git noticing that the patch is already
// applied, and git doesn't always notice this. If the patch contains only some
// of several adjacent added lines, the pick conflicts; if git can place the
// patch's added lines in more than one position, it may add them a second
// time. So we make the rebase pick a copy of the commit instead, whose parent
// already contains the patch.
func (self *PatchCommands) pickCopyWithPatchInParent(commit *models.Commit) error {
	copyHash, err := self.copyCommitWithPatchInParent(commit)
	if err != nil {
		return err
	}

	return self.rebase.ChangeTodoCommit(commit.Hash(), copyHash)
}

// Creates a copy of the given commit, with the same tree, author, and message,
// whose parent is the commit's original parent with the custom patch applied.
// It does this without touching the working tree or the index.
func (self *PatchCommands) copyCommitWithPatchInParent(commit *models.Commit) (string, error) {
	indexDir, err := os.MkdirTemp(self.os.GetTempDir(), "index")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(indexDir)
	indexFileEnvVar := "GIT_INDEX_FILE=" + filepath.Join(indexDir, "index")

	cmd := self.cmd.New(NewGitCmd("read-tree").Arg(commit.Parents()[0]).ToArgv())
	err = cmd.AddEnvVars(indexFileEnvVar).DontLog().Run()
	if err != nil {
		return "", err
	}

	patchFilePath, err := self.SaveTemporaryPatch(self.PatchBuilder.PatchToApply(false, false))
	if err != nil {
		return "", err
	}
	cmd = self.cmd.New(NewGitCmd("apply").Arg("--cached", patchFilePath).ToArgv())
	err = cmd.AddEnvVars(indexFileEnvVar).DontLog().Run()
	if err != nil {
		return "", err
	}

	cmd = self.cmd.New(NewGitCmd("write-tree").ToArgv())
	tree, _, err := cmd.AddEnvVars(indexFileEnvVar).DontLog().RunWithOutputs()
	if err != nil {
		return "", err
	}

	cmd = self.cmd.New(NewGitCmd("cat-file").Arg("commit", commit.Hash()).ToArgv())
	commitObject, _, err := cmd.DontLog().RunWithOutputs()
	if err != nil {
		return "", err
	}

	parentHash, err := self.writeCommitObject(
		commitObjectWithHeader(commitObject, "tree", strings.TrimSpace(tree)))
	if err != nil {
		return "", err
	}

	return self.writeCommitObject(commitObjectWithHeader(commitObject, "parent", parentHash))
}

func (self *PatchCommands) writeCommitObject(commitObject string) (string, error) {
	cmd := self.cmd.New(
		NewGitCmd("hash-object").
			Arg("-t", "commit", "-w", "--stdin").ToArgv())
	hash, _, err := cmd.
		SetStdin(commitObject).
		DontLog().
		RunWithOutputs()
	return strings.TrimSpace(hash), err
}

// Returns the given raw commit object with the value of the given header
// replaced. The commit's signature is dropped, since it no longer matches.
func commitObjectWithHeader(commitObject string, name string, value string) string {
	headers, message, _ := strings.Cut(commitObject, "\n\n")

	newHeaders := []string{}
	isInSignature := false
	for line := range strings.SplitSeq(headers, "\n") {
		// Continuation lines of a multi-line header start with a space
		if isInSignature && strings.HasPrefix(line, " ") {
			continue
		}

		isInSignature = strings.HasPrefix(line, "gpgsig ") || strings.HasPrefix(line, "gpgsig-sha256 ")
		if isInSignature {
			continue
		}

		if strings.HasPrefix(line, name+" ") {
			line = name + " " + value
		}
		newHeaders = append(newHeaders, line)
	}

	return strings.Join(newHeaders, "\n") + "\n\n" + message
}

// We have just applied a patch in reverse to discard it from a commit; if we
// now try to apply the patch again to move it to a later commit, or to the
// index, then this would conflict "with itself" in case the patch contained
// only some lines of a range of adjacent added lines. To solve this, we
// get the diff of HEAD and the original commit and then apply that.
func (self *PatchCommands) diffHeadAgainstCommit(commit *models.Commit) (string, error) {
	cmdArgs := NewGitCmd("diff").
		ParseableDiffPrefixes().
		Arg("--no-ext-diff", "--no-color").
		Arg("HEAD.." + commit.Hash()).
		ToArgv()

	return self.cmd.New(cmdArgs).RunWithOutput()
}
