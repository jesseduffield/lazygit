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

// WorkingTreeDiffActions implements what the files panel offers on the diff it renders
// into the focused main view: the diff itself, for the commands that need to read lines
// out of it rather than off the screen.
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
		WorktreeFileDiffCmdObj(node, true, self.showsStagedSide(pane), paths).
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
	return self.applyDiffLineSelection(infos, onStagedSide,
		git_commands.ApplyPatchOpts{Reverse: onStagedSide, Cached: true})
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
func (self *WorkingTreeDiffActions) applyDiffLineSelection(
	infos []types.DiffLineInfo, onStagedSide bool, opts git_commands.ApplyPatchOpts,
) error {
	self.c.LogAction(self.c.Tr.Actions.ApplyPatch)

	// A directory's diff spans several files, and a patch is of one file, so the
	// selected lines are grouped by the file they belong to and applied file by file.
	infosByFile := lo.GroupBy(infos, func(info types.DiffLineInfo) string { return info.Path })
	for path, fileInfos := range infosByFile {
		file := self.fileForDiffLinePath(path)
		if file == nil {
			continue
		}
		if err := self.applyDiffLines(file, fileInfos, onStagedSide, opts); err != nil {
			return err
		}
	}

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
// Each selected line is looked for by where it sits in the file. This tells the two
// halves of a modified line apart: the deletion and the addition replacing it share a
// position in the new file and differ only in being a deletion. Context lines are not
// selected: a patch of the lines you picked keeps whatever context it needs around
// them by itself.
func (self *WorkingTreeDiffActions) applyDiffLines(
	file *models.File, infos []types.DiffLineInfo, sourceCached bool, opts git_commands.ApplyPatchOpts,
) error {
	parsedPatch := patch.Parse(self.c.Git().WorkingTree.WorktreeFileDiff(file, true, sourceCached))

	type changeLine struct {
		lineNumber int
		isDeletion bool
	}
	selected := set.New[changeLine]()
	for _, info := range infos {
		if info.Type == types.DiffLineDeleted {
			selected.Add(changeLine{info.OldLine, true})
		} else {
			selected.Add(changeLine{info.NewLine, false})
		}
	}

	var patchLineIndices []int
	for idx, line := range parsedPatch.Lines() {
		var key changeLine
		switch {
		case line.IsAddition():
			key = changeLine{parsedPatch.LineNumberOfLine(idx), false}
		case line.IsDeletion():
			key = changeLine{parsedPatch.OldLineNumberOfLine(idx), true}
		default:
			continue
		}
		if selected.Includes(key) {
			patchLineIndices = append(patchLineIndices, idx)
		}
	}

	patchToApply := parsedPatch.
		Transform(patch.TransformOpts{
			Reverse:             opts.Reverse,
			IncludedLineIndices: patchLineIndices,
			FileNameOverride:    file.GetPath(),
		}).
		FormatPlain()
	if patchToApply == "" {
		return nil
	}

	return self.c.Git().Patch.ApplyPatch(patchToApply, opts)
}
