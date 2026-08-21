package controllers

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

// CommitDiffActions is what a panel showing a commit's diff offers on that diff in the
// focused main view. Five panels do: the commit files panel shows the diff of one file
// of a commit, and the commits, sub-commits, stash and reflog panels the whole diff of
// whatever they have selected. What they offer is the same for all of them, differing
// only in which diff it is — so they share this, each saying what it is showing.
type CommitDiffActions struct {
	c *ControllerCommon

	// The panel this belongs to, and what it is showing the diff of — nil when it has
	// nothing selected, and so no diff.
	panel  types.Context
	target func() *commitDiffTarget
}

// commitDiffTarget is the diff a panel is showing: the two ends of it, and whether it
// belongs to a commit lazygit may rewrite.
type commitDiffTarget struct {
	from      string
	to        string
	canRebase bool
}

var _ types.FocusedMainViewActions = &CommitDiffActions{}

func NewCommitDiffActions(
	c *ControllerCommon, panel types.Context, target func() *commitDiffTarget,
) *CommitDiffActions {
	return &CommitDiffActions{c: c, panel: panel, target: target}
}

// PlainDiff hands out the diff the asking pane is showing, for the given files — the
// commit's diff as in the main view, only without the commit's message and stat above it,
// or the diff the custom patch is previewed as, whose lines are the patch's own rather
// than the commit's.
//
// The patch's own diff is handed out whole: it is only ever as big as the patch, and it
// names its files under the trees the patch was materialized into rather than under the
// paths asked for.
func (self *CommitDiffActions) PlainDiff(pane types.DiffPaneContext, paths []string) string {
	if self.showsCustomPatch(pane) {
		return self.customPatchDiff()
	}

	target := self.target()
	if target == nil {
		return ""
	}
	return self.c.Helpers().Diff.PlainDiffBetweenRefs(target.from, target.to, paths)
}

// customPatchDiff is the diff the custom patch is previewed as, as git writes it — the
// diff behind what the pane previewing the patch shows, in which the lines shown there
// can be found again.
func (self *CommitDiffActions) customPatchDiff() string {
	treesDir := self.c.Git().Patch.PatchBuilder.TempDir()
	if treesDir == "" {
		return ""
	}
	// An error means the two trees differ, which is what a patch with anything in it looks
	// like; the diff itself is what we are after either way.
	diff, _ := self.c.Git().Diff.
		CustomPatchDiffCmdObj(treesDir, git_commands.DiffModePlain).
		RunWithOutput()
	return diff
}

// PrimaryAction takes the selected lines into the custom patch being built from this
// diff, or back out of it when the first of them is already in — the same toggling the
// commit files panel does to a whole file at a time.
//
// The commit is not touched, so the diff stays as it is: what changes is the patch
// beside it, and which of its lines are marked as being in that patch.
func (self *CommitDiffActions) PrimaryAction(pane types.DiffPaneContext, firstLineIdx int, lastLineIdx int) error {
	// In the pane showing the patch, the lines are the patch's own, so there they only
	// come back out of it.
	if self.showsCustomPatch(pane) {
		return self.removePatchLines(pane, firstLineIdx, lastLineIdx)
	}

	if self.c.UserConfig().Git.DiffContextSize == 0 {
		return fmt.Errorf(self.c.Tr.Actions.NotEnoughContextForCustomPatch,
			self.c.UserConfig().Keybinding.Universal.IncreaseContextInDiffView)
	}

	target := self.target()
	if target == nil {
		return nil
	}
	lines := self.c.Helpers().DiffLine.ChangeLinesInViewRange(pane.GetView(), firstLineIdx, lastLineIdx)
	if len(lines) == 0 {
		return nil
	}

	patchBuilder := self.c.Git().Patch.PatchBuilder
	from, reverse := self.patchEndpoints(target)
	// A patch is built from one diff, so building from another one means giving up the
	// patch there is — which the user is asked about, as entering the patch builder asks.
	mustDiscardPatch := patchBuilder.Active() && patchBuilder.NewPatchRequired(from, target.to, reverse)
	return self.c.ConfirmIf(mustDiscardPatch, types.ConfirmOpts{
		Title:  self.c.Tr.DiscardPatch,
		Prompt: self.c.Tr.DiscardPatchConfirm,
		HandleConfirm: func() error {
			if mustDiscardPatch {
				patchBuilder.Reset()
			}
			if !patchBuilder.Active() {
				patchBuilder.Start(from, target.to, reverse, target.canRebase)
			}

			if err := self.togglePatchLines(lines); err != nil {
				return err
			}
			// Taking the last line back out ends the patch rather than leaving an empty
			// one, so that the pane previewing it and the marks over the diff go with it.
			if patchBuilder.IsEmpty() {
				patchBuilder.Reset()
			}

			// The diff on screen is the one the marks belong to, so they can be brought up
			// to date at once rather than waiting for the render below.
			self.c.Helpers().DiffLine.RefreshInclusionGutter()

			// The selection moves on past the lines just toggled, to the next change of
			// the diff — which is still there, a toggle leaving the diff as it was, so
			// hold input back until it has moved: a second press meanwhile would toggle
			// the same lines straight back.
			self.c.GocuiGui().BeginBlockingEvents()
			self.c.Helpers().DiffLine.RevealSelectionAfterAction(pane, pane, firstLineIdx, len(lines),
				func() { _ = self.c.GocuiGui().EndBlockingEvents() })

			// The panel's own render, which is all that is needed: the marks over the diff
			// and the patch previewed beside it have changed, while the commit has not.
			self.c.PostRefreshUpdate(self.panel)
			return nil
		},
	})
}

// removePatchLines takes the selected lines of the custom patch out of it, which is what
// the primary action does in the pane showing the patch: everything shown there is in the
// patch already, so there is nothing else it could mean.
//
// A line of the patch is named by which of its file's changes it is, counted in the diff
// the patch is shown as, which is the same place that line has among the file's changes
// the patch holds. Line numbers would not do: a patch that leaves an earlier addition out
// numbers everything after it differently from the commit's diff.
func (self *CommitDiffActions) removePatchLines(
	pane types.DiffPaneContext, firstLineIdx int, lastLineIdx int,
) error {
	lines := self.c.Helpers().DiffLine.ChangeLinesInViewRange(pane.GetView(), firstLineIdx, lastLineIdx)
	if len(lines) == 0 {
		return nil
	}

	patchBuilder := self.c.Git().Patch.PatchBuilder
	previousPaths := self.previousPaths()
	for path, ordinals := range self.c.Helpers().DiffLine.ChangeLineOrdinals(self.customPatchDiff(), lines) {
		filename := self.patchBuilderPath(path)
		if filename == "" {
			continue
		}
		included := patchBuilder.IncludedChangeLineIndices(filename)
		indices := []int{}
		for _, ordinal := range ordinals {
			if ordinal < len(included) {
				indices = append(indices, included[ordinal])
			}
		}
		if len(indices) == 0 {
			continue
		}
		if err := patchBuilder.RemoveFileLineRange(filename, previousPaths[filename], indices); err != nil {
			return err
		}
	}
	// Taking the last line out ends the patch rather than leaving an empty one, as it does
	// in the diff beside this pane.
	if patchBuilder.IsEmpty() {
		patchBuilder.Reset()
	}

	self.c.Helpers().DiffLine.RefreshInclusionGutter()

	// The lines are gone from the patch, so the selection carries on from where they were,
	// as unstaging leaves it. Input is held until it has moved, so that a second press acts
	// on the patch as it now is.
	self.c.GocuiGui().BeginBlockingEvents()
	self.c.Helpers().DiffLine.RevealSelectionAfterAction(pane, pane, firstLineIdx, 0,
		func() { _ = self.c.GocuiGui().EndBlockingEvents() })

	self.c.PostRefreshUpdate(self.panel)
	return nil
}

// DiscardSelection takes the selected lines out of the commit they are part of, by
// building a patch of exactly those lines and removing that patch from the commit. It is
// a rebase, so a later commit that touches the same lines can conflict with it.
//
// The patch it needs is its own, so a patch being built is given up first — which the
// prompt says, there being no way to get it back.
func (self *CommitDiffActions) DiscardSelection(pane types.DiffPaneContext, firstLineIdx int, lastLineIdx int) error {
	target := self.target()
	if target == nil {
		return nil
	}
	lines := self.c.Helpers().DiffLine.ChangeLinesInViewRange(pane.GetView(), firstLineIdx, lastLineIdx)
	if len(lines) == 0 {
		return nil
	}
	commitIndex := self.indexOfTargetCommit(target)
	if commitIndex == -1 {
		return nil
	}

	patchBuilder := self.c.Git().Patch.PatchBuilder
	prompt := lo.Ternary(patchBuilder.IsEmpty(),
		self.c.Tr.DiscardLinesFromCommitPrompt,
		self.c.Tr.DiscardLinesFromCommitPromptWithReset)

	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.DiscardLinesFromCommitTitle,
		Prompt: prompt,
		HandleConfirm: func() error {
			from, reverse := self.patchEndpoints(target)
			patchBuilder.Reset()
			patchBuilder.Start(from, target.to, reverse, target.canRebase)
			if err := self.togglePatchLines(lines); err != nil {
				return err
			}
			if patchBuilder.IsEmpty() {
				return nil
			}

			// The rebase runs on a worker, which may not read the model, so the commits
			// it rewrites are taken here.
			commits := self.c.Model().Commits
			return self.c.WithWaitingStatusBlockingInput(types.WaitingStatusOpts{
				Message:              self.c.Tr.RebasingStatus,
				HideWorkingTreeState: true,
			}, func(gocui.Task) error {
				self.c.LogAction(self.c.Tr.Actions.RemovePatchFromCommit)
				err := self.c.Git().Patch.DeletePatchesFromCommit(commits, commitIndex)
				return self.c.Helpers().MergeAndRebase.CheckMergeOrRebase(err)
			})
		},
	})
	return nil
}

// DiscardSelectionDisabledReason says why the selected lines can't be taken out of the
// commit: doing so rewrites it, which is only ours to do for a commit of the branch we
// are on, and not while a rebase is already under way. In the pane previewing the custom
// patch there is nothing to discard from — the lines there are the patch's, and space
// takes them back out of it.
func (self *CommitDiffActions) DiscardSelectionDisabledReason(pane types.DiffPaneContext) *types.DisabledReason {
	if self.showsCustomPatch(pane) {
		return &types.DisabledReason{Text: self.c.Tr.CannotDiscardFromCustomPatchView, ShowErrorInPanel: true}
	}
	target := self.target()
	if target == nil || !target.canRebase {
		return &types.DisabledReason{Text: self.c.Tr.CanOnlyDiscardFromLocalCommits, ShowErrorInPanel: true}
	}
	if self.c.Git().Status.WorkingTreeState().Any() {
		return &types.DisabledReason{Text: self.c.Tr.CantPatchWhileRebasingError, ShowErrorInPanel: true}
	}
	if self.c.UserConfig().Git.DiffContextSize == 0 {
		return &types.DisabledReason{
			Text: fmt.Sprintf(self.c.Tr.Actions.NotEnoughContextToRemoveLines,
				self.c.UserConfig().Keybinding.Universal.IncreaseContextInDiffView),
			ShowErrorInPanel: true,
		}
	}
	return nil
}

// PatchInclusion says which lines of the commit's diff are in the custom patch being
// built from it. nil when there is no such patch: none is being built at all, or the one
// being built is of another diff, whose lines are not these however alike they look.
func (self *CommitDiffActions) PatchInclusion() func(types.DiffLineInfo) bool {
	patchBuilder := self.c.Git().Patch.PatchBuilder
	target := self.target()
	if !patchBuilder.Active() || target == nil {
		return nil
	}
	from, reverse := self.patchEndpoints(target)
	if patchBuilder.NewPatchRequired(from, target.to, reverse) {
		return nil
	}

	// Which lines of a file are in the patch is asked of the patch builder per file, and
	// a diff can span many, so each is asked about when a line of it first comes up.
	includedByPath := map[string]*set.Set[patch.LineIdentity]{}
	return func(info types.DiffLineInfo) bool {
		path := self.patchBuilderPath(info.Path)
		if path == "" {
			return false
		}
		included, asked := includedByPath[path]
		if !asked {
			included = set.NewFromSlice(patchBuilder.IncludedLineIdentities(path))
			includedByPath[path] = included
		}
		return included.Includes(info.PatchLineIdentity())
	}
}

// togglePatchLines takes the given lines of the commit's diff into the custom patch, or
// out of it. Which of the two it is is decided once, by the first line of the selection:
// pointing at a line that is already in the patch takes the whole selection out of it,
// as toggling a selection of files in the commit files panel does.
func (self *CommitDiffActions) togglePatchLines(lines []types.DiffLineInfo) error {
	patchBuilder := self.c.Git().Patch.PatchBuilder

	// The files the selection covers, in the order the diff shows them, and per file the
	// lines of it that are selected: a patch is built a file at a time, while a selection
	// can span several of them.
	paths := []string{}
	linesByPath := map[string][]patch.LineIdentity{}
	for _, line := range lines {
		path := self.patchBuilderPath(line.Path)
		if path == "" {
			continue
		}
		if _, seen := linesByPath[path]; !seen {
			paths = append(paths, path)
		}
		linesByPath[path] = append(linesByPath[path], line.PatchLineIdentity())
	}
	if len(paths) == 0 {
		return nil
	}

	previousPaths := self.previousPaths()
	indicesByPath := map[string][]int{}
	wholeFileByPath := map[string]bool{}
	for _, path := range paths {
		indices, err := patchBuilder.PatchLineIndicesForLines(path, previousPaths[path], linesByPath[path])
		if err != nil {
			return err
		}
		indicesByPath[path] = indices

		wholeFile, err := patchBuilder.SelectionRepresentsWholeFile(path, previousPaths[path], linesByPath[path])
		if err != nil {
			return err
		}
		wholeFileByPath[path] = wholeFile
	}

	included, err := patchBuilder.GetFileIncLineIndices(paths[0], previousPaths[paths[0]])
	if err != nil {
		return err
	}
	removing := len(indicesByPath[paths[0]]) > 0 && lo.Contains(included, indicesByPath[paths[0]][0])

	for _, path := range paths {
		if len(indicesByPath[path]) == 0 {
			continue
		}
		var err error
		switch {
		case wholeFileByPath[path] && removing:
			err = patchBuilder.RemoveFile(path, previousPaths[path])
		case wholeFileByPath[path]:
			err = patchBuilder.AddFileWhole(path, previousPaths[path])
		case removing:
			err = patchBuilder.RemoveFileLineRange(path, previousPaths[path], indicesByPath[path])
		default:
			err = patchBuilder.AddFileLineRange(path, previousPaths[path], indicesByPath[path])
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// previousPaths says which files of the diff were renamed, and what they were called
// before. A renamed file's diff only comes out as a rename when git is asked about both
// of its paths, and its lines are numbered in the file under its old name, so the patch
// builder has to be told the old path along with them.
func (self *CommitDiffActions) previousPaths() map[string]string {
	target := self.target()
	if target == nil {
		return nil
	}
	from, reverse := self.patchEndpoints(target)
	files, err := self.c.Git().Loaders.CommitFileLoader.GetFilesInDiff(from, target.to, reverse)
	if err != nil {
		return nil
	}

	previousPaths := map[string]string{}
	for _, file := range files {
		if file.PreviousPath != "" {
			previousPaths[file.Path] = file.PreviousPath
		}
	}
	return previousPaths
}

// patchEndpoints gives the two ends of the diff a patch is built from. They are the ends
// of the diff shown, except in diffing mode, where what is shown is a diff against
// another ref, possibly the other way around.
func (self *CommitDiffActions) patchEndpoints(target *commitDiffTarget) (string, bool) {
	return self.c.Modes().Diffing.GetFromAndReverseArgsForDiff(target.from)
}

// patchBuilderPath turns the absolute path a diff line carries into the repo-relative
// one the patch builder keys a file by, and "" for a path that is no file of this repo.
func (self *CommitDiffActions) patchBuilderPath(path string) string {
	relativePath, err := filepath.Rel(self.c.Git().RepoPaths.WorktreePath(), path)
	if err != nil || strings.HasPrefix(relativePath, "..") {
		return ""
	}
	return filepath.ToSlash(relativePath)
}

// indexOfTargetCommit finds the commit the diff belongs to among the commits of the
// branch we are on, which is how a rebase is told which commit to rewrite. -1 when it
// isn't one of them, in which case there is nothing we can rewrite.
func (self *CommitDiffActions) indexOfTargetCommit(target *commitDiffTarget) int {
	return lo.IndexOf(
		lo.Map(self.c.Model().Commits, func(commit *models.Commit, _ int) string { return commit.Hash() }),
		target.to)
}

// showsCustomPatch reports whether the given main pane is the one previewing the custom
// patch being built, rather than the commit's diff — which for a commit's diff is always
// the lower one.
func (self *CommitDiffActions) showsCustomPatch(pane types.DiffPaneContext) bool {
	return pane.GetKey() == self.c.Contexts().NormalSecondary.GetKey()
}
