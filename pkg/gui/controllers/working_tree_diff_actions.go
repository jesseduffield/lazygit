package controllers

import (
	"fmt"
	"path/filepath"

	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/samber/lo"
)

// WorkingTreeDiffActions is what the files panel offers on the diff it renders into
// the focused main view: the diff itself, for the commands that need to read lines out
// of it rather than off the screen.
type WorkingTreeDiffActions struct {
	c *ControllerCommon
}

var _ types.FocusedMainViewActions = &WorkingTreeDiffActions{}

func NewWorkingTreeDiffActions(c *ControllerCommon) *WorkingTreeDiffActions {
	return &WorkingTreeDiffActions{c: c}
}

func (self *WorkingTreeDiffActions) context() *context.WorkingTreeContext {
	return self.c.Contexts().Files
}

// PlainDiff hands out the working tree's diff for the given files, taken from the
// side of the index that the asking pane shows.
func (self *WorkingTreeDiffActions) PlainDiff(pane types.DiffPaneContext, paths []string) string {
	node := self.context().GetSelected()
	if node == nil {
		return ""
	}
	// An error means there is no diff to be had, which for our purposes is the same as
	// an empty one.
	diff, _ := self.c.Git().WorkingTree.
		WorktreeFileDiffCmdObj(node, git_commands.DiffModePlain, self.showsStagedSide(pane), paths).
		RunWithOutput()
	return diff
}

// showsStagedSide reports whether the given main pane is the one showing the staged
// side of a file's diff, which is always the lower one.
func (self *WorkingTreeDiffActions) showsStagedSide(pane types.DiffPaneContext) bool {
	return pane.GetKey() == self.c.Contexts().NormalSecondary.GetKey()
}

// PrimaryAction stages the selected diff lines, or takes them back out of the index
// when what is selected is the staged side of the diff.
func (self *WorkingTreeDiffActions) PrimaryAction(pane types.DiffPaneContext, firstLineIdx int, lastLineIdx int) error {
	if self.c.UserConfig().Git.DiffContextSize == 0 {
		return fmt.Errorf(self.c.Tr.Actions.NotEnoughContextToStage,
			self.c.UserConfig().Keybinding.Universal.IncreaseContextInDiffView)
	}

	infos, onStagedSide, ok := self.diffLineSelection(pane, firstLineIdx, lastLineIdx)
	if !ok {
		return nil
	}

	// Either way the patch goes to the index: forwards from the unstaged side to stage
	// it, backwards from the staged side to take it back out.
	return self.applyDiffLineSelection(pane, firstLineIdx, infos, onStagedSide,
		git_commands.ApplyPatchOpts{Reverse: onStagedSide, Cached: true})
}

// DiscardSelection takes the selected diff lines out of the working tree — or, on the
// staged side, out of the index, which is where "discard this" means "I don't want it
// staged".
func (self *WorkingTreeDiffActions) DiscardSelection(pane types.DiffPaneContext, firstLineIdx int, lastLineIdx int) error {
	if self.c.UserConfig().Git.DiffContextSize == 0 {
		return fmt.Errorf(self.c.Tr.Actions.NotEnoughContextToDiscard,
			self.c.UserConfig().Keybinding.Universal.IncreaseContextInDiffView)
	}

	infos, onStagedSide, ok := self.diffLineSelection(pane, firstLineIdx, lastLineIdx)
	if !ok {
		return nil
	}

	// Either way the change is applied backwards; the side it is applied to is what
	// makes the difference. On the staged side that is the index, which is the same
	// thing as unstaging and so is not destructive. On the unstaged side it is the
	// working tree, where the change is gone for good — hence the confirmation.
	return self.c.ConfirmIf(!onStagedSide && !self.c.UserConfig().Gui.SkipDiscardChangeWarning,
		types.ConfirmOpts{
			Title:  self.c.Tr.DiscardChangeTitle,
			Prompt: self.c.Tr.DiscardChangePrompt,
			HandleConfirm: func() error {
				return self.applyDiffLineSelection(pane, firstLineIdx, infos, onStagedSide,
					git_commands.ApplyPatchOpts{Reverse: true, Cached: onStagedSide})
			},
		})
}

// DiscardSelectionDisabledReason is nil: a change of the working tree can always be
// thrown away, and one in the index always taken back out of it.
func (self *WorkingTreeDiffActions) DiscardSelectionDisabledReason(types.DiffPaneContext) *types.DisabledReason {
	return nil
}

// PatchInclusion is nil: a custom patch is built from a commit's diff, never from the
// working tree's, so no line of this diff is ever in one.
func (self *WorkingTreeDiffActions) PatchInclusion() func(types.DiffLineInfo) bool {
	return nil
}

// diffLineSelection resolves what the user has selected in a pane of the focused main
// view to the change lines to act on, and reports whether they are the staged side of
// the diff — which is a question about the pane, so it is the same for every file of a
// directory's diff. ok is false when the selection holds no change line, in which case
// there is nothing to act on.
func (self *WorkingTreeDiffActions) diffLineSelection(
	pane types.DiffPaneContext, firstLineIdx int, lastLineIdx int,
) (infos []types.DiffLineInfo, onStagedSide bool, ok bool) {
	infos = self.c.Helpers().DiffLine.ChangeLinesInViewRange(pane.GetView(), firstLineIdx, lastLineIdx)
	if len(infos) == 0 {
		return nil, false, false
	}
	return infos, self.showsStagedSide(pane), true
}

// applyDiffLineSelection applies the selected change lines, a patch per file, and
// re-renders what that changed. onStagedSide says which of the file's two diffs the
// lines were selected in and so are to be found in; opts says how to apply them.
// firstLineIdx is where the selection started, which is where the work carries on from
// once the diff has changed under it.
func (self *WorkingTreeDiffActions) applyDiffLineSelection(
	pane types.DiffPaneContext, firstLineIdx int,
	infos []types.DiffLineInfo, onStagedSide bool, opts git_commands.ApplyPatchOpts,
) error {
	self.c.LogAction(self.c.Tr.Actions.ApplyPatch)

	// A directory's diff spans several files, and a patch is of one file, so the
	// selected lines are grouped by the file they belong to and applied file by file.
	infosByFile := lo.GroupBy(infos, func(info types.DiffLineInfo) string { return info.Path })
	acted := set.New[string]()
	actedSideRemains := false
	for path, fileInfos := range infosByFile {
		file := self.fileForDiffLinePath(path)
		if file == nil {
			continue
		}
		changesLeft, err := self.applyDiffLines(file, fileInfos, onStagedSide, opts)
		if err != nil {
			return err
		}
		acted.Add(file.GetPath())
		actedSideRemains = actedSideRemains || changesLeft
	}
	if !actedSideRemains {
		actedSideRemains = self.anyFileHasChangesOnSide(acted, onStagedSide)
	}

	// Whether the other side has anything is a question about the pane the work would
	// carry on in. Where the lines went into the index they are there now; where they
	// were thrown away that side is as the model describes it, having not been touched.
	otherSideHasChanges := opts.Cached || self.anyFileHasChangesOnSide(set.New[string](), !onStagedSide)

	// The refresh below queues the re-render of the diff we just changed; this rides it,
	// so that the selection ends up on the change that took the place of the one acted
	// on rather than at a position that means nothing any more — and in the pane the
	// work carries on in, which is not always the one it was in.
	self.revealSelectionInPaneItLandsIn(pane, firstLineIdx, actedSideRemains, otherSideHasChanges)

	// Block input until the refresh has landed, so that a quick second keypress acts on
	// the diff as it now is rather than on the one we just changed.
	self.c.RefreshBlockingInput(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}})
	return nil
}

// fileForDiffLinePath maps the absolute path a diff line carries to the working tree
// file it belongs to, or nil for a path that is no file of this repo's working tree.
func (self *WorkingTreeDiffActions) fileForDiffLinePath(path string) *models.File {
	relativePath, err := filepath.Rel(self.c.Git().RepoPaths.WorktreePath(), path)
	if err != nil {
		return nil
	}
	return self.context().FileTreeViewModel.GetFile(filepath.ToSlash(relativePath))
}

// applyDiffLines applies the given change lines of one file — a line, a hunk, a range —
// as a patch built from that file's own diff:
//
//   - stage:   read the unstaged diff, apply it to the index
//   - unstage: read the staged diff, apply it to the index backwards
//
// sourceCached names the diff the lines were selected in, which is where they are found
// again; opts says how to apply what is built from them. The two are independent — a
// discard reads one side and reverses it — so they are passed separately.
//
// Each selected line is looked for by where it sits in the file, which is what tells
// the two halves of a modified line apart: the deletion and the addition replacing it
// share a position in the new file and differ only in being a deletion. Context lines
// are not selected: a patch of the lines you picked keeps whatever context it needs
// around them by itself.
//
// It reports whether the diff it read holds changes the selection didn't cover, which
// is how the caller knows whether the side acted on still has anything of this file in
// it once we are done.
func (self *WorkingTreeDiffActions) applyDiffLines(
	file *models.File, infos []types.DiffLineInfo, sourceCached bool, opts git_commands.ApplyPatchOpts,
) (bool, error) {
	parsedPatch := patch.Parse(self.c.Git().WorkingTree.WorktreeFileDiff(file, git_commands.DiffModePlain, sourceCached))

	patchLineIndices := patch.ChangeLineIndicesForLines(parsedPatch,
		lo.Map(infos, func(info types.DiffLineInfo, _ int) patch.LineIdentity {
			return info.PatchLineIdentity()
		}))

	changesLeft := len(patchLineIndices) < changeLineCount(parsedPatch)

	// Acting on every change of a file is acting on the file itself, and saying so is
	// not the same as applying its diff. The diff of a deleted file is its content
	// going away, and putting that into the index line by line leaves an empty file
	// there rather than the deletion; the diff of an added one is its whole content,
	// and taking that back out leaves an empty file in the index rather than an
	// untracked one.
	if !changesLeft && opts.Cached {
		if opts.Reverse {
			return false, self.c.Git().WorkingTree.UnStageFile(file.Names(), file.Tracked)
		}
		return false, self.c.Git().WorkingTree.StageFile(file.GetPath())
	}

	patchToApply := parsedPatch.
		Transform(patch.TransformOpts{
			Reverse:             opts.Reverse,
			IncludedLineIndices: patchLineIndices,
			FileNameOverride:    file.GetPath(),
		}).
		FormatPlain()
	if patchToApply == "" {
		return changesLeft, nil
	}

	return changesLeft, self.c.Git().Patch.ApplyPatch(patchToApply, opts)
}

// changeLineCount returns how many of a patch's lines are changes rather than context
// or header, which is how many of them a selection of the whole diff covers.
func changeLineCount(p *patch.Patch) int {
	return lo.CountBy(p.Lines(), func(line *patch.PatchLine) bool {
		return line.IsAddition() || line.IsDeletion()
	})
}

// revealSelectionInPaneItLandsIn arranges for the selection to carry on where the work
// does, which is not always the pane it was in.
//
// Each side of the diff has a pane of its own, so acting on one usually leaves
// everything where it is. But a pane is only shown while its side has something in it:
// staging the last unstaged change takes the upper pane away, and unstaging the last
// staged one takes the lower one away. The refresh moves the focus into whichever pane
// is left, and this puts the selection there to meet it — on the lines just acted on,
// which are in that pane now, unless they were discarded rather than moved, in which
// case on what is left of the file.
func (self *WorkingTreeDiffActions) revealSelectionInPaneItLandsIn(
	pane types.DiffPaneContext, firstLineIdx int, actedSideRemains bool, otherSideHasChanges bool,
) {
	target := pane
	if !actedSideRemains && otherSideHasChanges {
		target = self.otherPane(pane)
	}

	// Hold input back until the selection is on the change the work carries on from. The
	// refresh holds it until the model is up to date, but the diff is re-rendered after
	// that, and until it has been the selection is still on lines that aren't there any
	// more — so a key pressed meanwhile would act on nothing.
	self.c.GocuiGui().BeginBlockingEvents()
	self.c.Helpers().DiffLine.RevealSelectionAfterAction(pane, target, firstLineIdx, 0,
		func() { _ = self.c.GocuiGui().EndBlockingEvents() })
}

// otherPane returns the main pane that isn't the given one.
func (self *WorkingTreeDiffActions) otherPane(pane types.DiffPaneContext) types.DiffPaneContext {
	if pane.GetKey() == self.c.Contexts().Normal.GetKey() {
		return self.c.Contexts().NormalSecondary
	}
	return self.c.Contexts().Normal
}

// anyFileHasChangesOnSide reports whether any file under the selected node, other than
// the ones named by except, has changes on the given side of the index, as the model
// has them. The model is right about any file the action didn't touch; the ones it did
// touch report for themselves, their entry not being right until the refresh lands.
func (self *WorkingTreeDiffActions) anyFileHasChangesOnSide(except *set.Set[string], staged bool) bool {
	node := self.context().GetSelected()
	if node == nil {
		return false
	}

	found := false
	_ = node.ForEachFile(func(file *models.File) error {
		if except.Includes(file.GetPath()) {
			return nil
		}
		if (staged && file.HasStagedChanges) || (!staged && file.HasUnstagedChanges) {
			found = true
		}
		return nil
	})
	return found
}
