package controllers

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/generics/set"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
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
// side of the index that the asking pane shows. It fails if the diff of one of the
// files is no longer the one the pane shows (see shownDiff).
func (self *WorkingTreeDiffActions) PlainDiff(pane types.DiffPaneContext, paths []string) (string, error) {
	node := self.context().GetSelected()
	if node == nil {
		return "", nil
	}
	shown := self.shownDiffs(pane)
	onStagedSide := self.showsStagedSide(pane)
	// An error means there is no diff to be had, which for our purposes is the same as
	// an empty one.
	diff, _ := self.c.Git().WorkingTree.
		WorktreeFileDiffCmdObj(node, git_commands.DiffModePlain, onStagedSide, paths).
		RunWithOutput()

	versions := self.c.Helpers().DiffLine.FileVersionsInDiff(diff)
	worktreePath := self.c.Git().RepoPaths.WorktreePath()
	for _, path := range paths {
		absolutePath := filepath.Join(worktreePath, path)
		filePaths := []string{path}
		if file := self.fileForDiffLinePath(absolutePath); file != nil {
			filePaths = file.Names()
		}
		if err := self.checkShown(shown[absolutePath], filePaths, onStagedSide, versions[absolutePath]); err != nil {
			return "", err
		}
	}
	return diff, nil
}

// showsStagedSide reports whether the given main pane is the one showing the staged
// side of a file's diff, which is always the lower one.
func (self *WorkingTreeDiffActions) showsStagedSide(pane types.DiffPaneContext) bool {
	return pane.GetKey() == self.c.Contexts().NormalSecondary.GetKey()
}

// PrimaryAction stages the selected diff lines, or takes them back out of the index
// when what is selected is the staged side of the diff.
func (self *WorkingTreeDiffActions) PrimaryAction(pane types.DiffPaneContext, firstBufferLine int, lastBufferLine int) error {
	if self.c.UserConfig().Git.DiffContextSize == 0 {
		return fmt.Errorf(self.c.Tr.Actions.NotEnoughContextToStage,
			self.c.UserConfig().Keybinding.Universal.IncreaseContextInDiffView)
	}

	linesByFile, onStagedSide, err := self.diffLineSelection(pane, firstBufferLine, lastBufferLine)
	if err != nil {
		return err
	}
	if len(linesByFile) == 0 {
		return nil
	}

	// Either way the patch goes to the index: forwards from the unstaged side to stage
	// it, backwards from the staged side to take it back out.
	return self.applyDiffLineSelection(pane, firstBufferLine, linesByFile, onStagedSide,
		git_commands.ApplyPatchOpts{Reverse: onStagedSide, Cached: true})
}

// DiscardSelection takes the selected diff lines out of the working tree — or, on the
// staged side, out of the index, which is where "discard this" means "I don't want it
// staged".
func (self *WorkingTreeDiffActions) DiscardSelection(pane types.DiffPaneContext, firstBufferLine int, lastBufferLine int) error {
	if self.c.UserConfig().Git.DiffContextSize == 0 {
		return fmt.Errorf(self.c.Tr.Actions.NotEnoughContextToDiscard,
			self.c.UserConfig().Keybinding.Universal.IncreaseContextInDiffView)
	}

	linesByFile, onStagedSide, err := self.diffLineSelection(pane, firstBufferLine, lastBufferLine)
	if err != nil {
		return err
	}
	if len(linesByFile) == 0 {
		return nil
	}

	// Either way the change is applied backwards; which side it is applied to decides
	// how destructive that is. On the staged side it goes to the index; this is the
	// same as unstaging, so nothing is lost. On the unstaged side it goes to the
	// working tree, where the change is gone for good, so we ask first.
	askFirst := !onStagedSide && !self.c.UserConfig().Gui.SkipDiscardChangeWarning

	// Building the patch reads the diff again and fails if it changed since it was
	// shown. That happens once the user has confirmed, since the diff can change while
	// they decide; checking before asking as well spares them confirming a discard
	// that won't happen.
	if askFirst {
		for _, fileLines := range linesByFile {
			if _, err := self.readFileDiff(fileLines, onStagedSide); err != nil {
				return err
			}
		}
	}

	return self.c.ConfirmIf(askFirst,
		types.ConfirmOpts{
			Title:  self.c.Tr.DiscardChangeTitle,
			Prompt: self.c.Tr.DiscardChangePrompt,
			HandleConfirm: func() error {
				return self.applyDiffLineSelection(pane, firstBufferLine, linesByFile, onStagedSide,
					git_commands.ApplyPatchOpts{Reverse: true, Cached: onStagedSide})
			},
		})
}

// DiscardSelectionDisabledReason is nil: a change of the working tree can always be
// thrown away, and one in the index always taken back out of it.
func (self *WorkingTreeDiffActions) DiscardSelectionDisabledReason(types.DiffPaneContext) *types.DisabledReason {
	return nil
}

// EditHunk opens the git hunk holding the selection in an editor, as a patch against
// the index, and applies whatever comes back. It is how you stage something the diff
// can't express — half of a changed line, or a change written differently from either
// side — since what the editor hands back is applied rather than matched against the
// file's own diff.
//
// The hunk is the git one, context and all, rather than lazygit's block of adjacent
// changes: an editable patch is one that still applies, and the context lines are what
// let git place it.
func (self *WorkingTreeDiffActions) EditHunk(
	pane types.DiffPaneContext, firstBufferLine int, lastBufferLine int,
) error {
	linesByFile, onStagedSide, err := self.diffLineSelection(pane, firstBufferLine, lastBufferLine)
	if err != nil {
		return err
	}
	if len(linesByFile) == 0 {
		return nil
	}
	file := linesByFile[0].file

	parsedPatch, err := self.readFileDiff(linesByFile[0], onStagedSide)
	if err != nil {
		return err
	}
	lineIndices := patch.ChangeLineIndicesForLines(parsedPatch,
		[]patch.LineIdentity{linesByFile[0].lines[0].PatchLineIdentity()})
	if len(lineIndices) == 0 {
		return diffLinesNotFoundError(self.c, file.GetPath())
	}

	hunkIdx := parsedPatch.HunkContainingLine(lineIndices[0])
	hunkStartIdx := parsedPatch.HunkStartIdx(hunkIdx)
	patchText := parsedPatch.
		Transform(patch.TransformOpts{
			Reverse:             onStagedSide,
			IncludedLineIndices: patch.ExpandRange(hunkStartIdx, parsedPatch.HunkEndIdx(hunkIdx)),
			FileNameOverride:    file.GetPath(),
		}).
		FormatPlain()

	patchFilepath, err := self.c.Git().Patch.SaveTemporaryPatch(patchText)
	if err != nil {
		return err
	}

	// The patch is written with a two-line header before its hunk, so the line the
	// user was on sits that much further down the file they are about to edit.
	const headerLineCount = 2
	if err := self.c.Helpers().Files.EditFileAtLineAndWait(patchFilepath,
		lineIndices[0]-hunkStartIdx+headerLineCount+1); err != nil {
		return err
	}

	editedPatchText, err := self.c.Git().File.Cat(patchFilepath)
	if err != nil {
		return err
	}

	self.c.LogAction(self.c.Tr.Actions.ApplyPatch)

	// Everything the editor left behind is taken, this being a patch the user wrote
	// rather than a selection out of one of ours.
	lineCount := strings.Count(editedPatchText, "\n") + 1
	newPatchText := patch.
		Parse(editedPatchText).
		Transform(patch.TransformOpts{
			IncludedLineIndices: patch.ExpandRange(0, lineCount),
			FileNameOverride:    file.GetPath(),
		}).
		FormatPlain()

	if err := self.c.Git().Patch.ApplyPatch(newPatchText, git_commands.ApplyPatchOpts{
		Reverse: onStagedSide,
		Cached:  true,
	}); err != nil {
		return err
	}

	// Block input until the refresh has landed, as the staging commands do: the diff is
	// about to be rebuilt from a file that no longer looks the way it did.
	self.c.RefreshBlockingInput(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}})
	return nil
}

// PatchInclusion is nil: a custom patch is built from a commit's diff, never from the
// working tree's, so no line of this diff is ever in one.
func (self *WorkingTreeDiffActions) PatchInclusion() func(types.DiffLineInfo) bool {
	return nil
}

// fileDiffLines is the selected change lines of one file of the diff.
type fileDiffLines struct {
	file  *models.File
	lines []types.DiffLineInfo

	// What the pane showed of the file's diff when the lines were selected.
	shown shownDiff
}

// shownDiff is what a pane showed of a file's diff at one moment, to tell whether a
// diff of the file read later is still that diff. The selected lines are looked up in
// the diff read later by where they sit in the file, and in a different diff that
// place can hold a different line. The diff changes whenever the file, the index or
// HEAD is changed outside of lazygit and lazygit hasn't refreshed since.
type shownDiff struct {
	// The versions of the file the diff is between, as the header of the diff on screen
	// states them; nil if the screen shows no header to read them from (see
	// DiffLineHelper.FileVersionsOnScreen).
	versions *patch.FileVersions

	// What the diff on screen was made from, as its render took it before running git
	// (see FilesController.diffInputsTag); nil if it took nothing. A diff renderer that
	// states its rows shows no header, and this is what is left to tell by.
	inputs *git_commands.DiffInputs
}

// shownDiffs reads what the pane shows of the diff of each of its files, by the file's
// absolute path.
func (self *WorkingTreeDiffActions) shownDiffs(pane types.DiffPaneContext) map[string]shownDiff {
	shown := map[string]shownDiff{}
	view := pane.GetView()

	if inputs, ok := view.ContentTag().(*git_commands.DiffInputs); ok {
		worktreePath := self.c.Git().RepoPaths.WorktreePath()
		for _, path := range inputs.Paths() {
			shown[filepath.Join(worktreePath, path)] = shownDiff{inputs: inputs}
		}
	}

	// The versions on screen are the diff itself, so they win where there are both.
	versions, _ := self.c.Helpers().DiffLine.FileVersionsOnScreen(view)
	for path, fileVersions := range versions {
		shown[path] = shownDiff{versions: &fileVersions}
	}
	return shown
}

// checkShown fails if a diff of a file, just read from the given side of the index and
// between the given versions, isn't the diff that was shown. paths are the file's
// paths in the repo: the one it has, and the one it was renamed from, if any (see
// models.File.Names). With nothing to tell by, any diff passes.
//
// The inputs of the diff are taken after the diff was read, so that a change in
// between counts as a change.
func (self *WorkingTreeDiffActions) checkShown(
	shown shownDiff, paths []string, onStagedSide bool, versions patch.FileVersions,
) error {
	same := true
	switch {
	case shown.versions != nil:
		same = shown.versions.Same(versions)
	case shown.inputs != nil:
		current, err := self.c.Git().WorkingTree.GetDiffInputs(paths, onStagedSide, !onStagedSide)
		if err != nil {
			return err
		}
		same = shown.inputs.SameFor(current, paths, onStagedSide)
	}
	if !same {
		return self.diffChangedError(paths[0])
	}
	return nil
}

// diffChangedError reports that the diff of the file at the given path in the repo is
// no longer the one on screen. It has the files panel refreshed, so that the diff on
// screen is the current one when the user looks at it again.
func (self *WorkingTreeDiffActions) diffChangedError(path string) error {
	self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}})
	return errors.New(utils.ResolvePlaceholderString(self.c.Tr.DiffChangedSinceShown,
		map[string]string{"path": path}))
}

// diffLineSelection resolves what the user has selected in a pane of the focused main
// view to the change lines to act on, and reports whether they are the staged side of
// the diff — which is a question about the pane, so it is the same for every file of a
// directory's diff. When the selection holds no change line, it returns no file, and
// there is nothing to act on.
//
// A directory's diff spans several files, and a patch is of one file, so the lines are
// grouped by the file they belong to, in the order the diff shows the files. If a line
// belongs to no file of the files panel, the diff and the panel disagree. Acting on the
// other lines would do only part of what the user asked for, so this fails instead.
//
// With each file's lines goes what the pane shows of the file's diff, read at the same
// moment as the lines.
func (self *WorkingTreeDiffActions) diffLineSelection(
	pane types.DiffPaneContext, firstBufferLine int, lastBufferLine int,
) (linesByFile []fileDiffLines, onStagedSide bool, err error) {
	infos := self.c.Helpers().DiffLine.ChangeLinesInBufferRange(pane.GetView(), firstBufferLine, lastBufferLine)
	infosByPath := lo.GroupBy(infos, func(info types.DiffLineInfo) string { return info.Path })
	paths := lo.Uniq(lo.Map(infos, func(info types.DiffLineInfo, _ int) string { return info.Path }))
	shown := self.shownDiffs(pane)
	for _, path := range paths {
		file := self.fileForDiffLinePath(path)
		if file == nil {
			return nil, false, diffLinesFileNotFoundError(self.c, path)
		}
		linesByFile = append(linesByFile, fileDiffLines{file: file, lines: infosByPath[path], shown: shown[path]})
	}
	return linesByFile, self.showsStagedSide(pane), nil
}

// applyDiffLineSelection applies the selected change lines, a patch per file, and
// re-renders what that changed. onStagedSide says which of the file's two diffs the
// lines were selected in and so are to be found in; opts says how to apply them.
// firstBufferLine is where the selection started, which is where the work carries on
// from once the diff has changed under it.
func (self *WorkingTreeDiffActions) applyDiffLineSelection(
	pane types.DiffPaneContext, firstBufferLine int,
	linesByFile []fileDiffLines, onStagedSide bool, opts git_commands.ApplyPatchOpts,
) error {
	// If the lines of one file can't be found, no file is acted on, so every file's
	// patch is built before any is applied.
	patches := make([]diffLinesPatch, 0, len(linesByFile))
	for _, fileLines := range linesByFile {
		p, err := self.buildDiffLinesPatch(fileLines, onStagedSide, opts)
		if err != nil {
			return err
		}
		patches = append(patches, p)
	}

	self.c.LogAction(self.c.Tr.Actions.ApplyPatch)

	acted := set.New[string]()
	actedSideRemains := false
	for _, p := range patches {
		if err := self.applyDiffLinesPatch(p, opts); err != nil {
			return err
		}
		acted.Add(p.file.GetPath())
		actedSideRemains = actedSideRemains || p.changesLeft
	}
	if !actedSideRemains {
		actedSideRemains = self.anyFileHasChangesOnSide(acted, onStagedSide)
	}

	// Whether the other side has anything decides which pane the work carries on in.
	// If the lines were staged, they are in the index now, so that side has them. If
	// they were discarded, the other side was not touched, so the model still
	// describes it correctly.
	otherSideHasChanges := opts.Cached || self.anyFileHasChangesOnSide(set.New[string](), !onStagedSide)

	// The refresh below queues the re-render of the diff we just changed; this rides it,
	// so that the selection ends up on the change that took the place of the one acted
	// on rather than at a position that means nothing any more — and in the pane the
	// work carries on in, which is not always the one it was in.
	self.revealSelectionInPaneItLandsIn(pane, firstBufferLine, actedSideRemains, otherSideHasChanges)

	// Block input until the refresh has landed, so that a quick second keypress acts on
	// the diff as it now is rather than on the one we just changed.
	self.c.RefreshBlockingInput(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}})
	return nil
}

// fileForDiffLinePath maps the absolute path a diff line carries to the working tree
// file it belongs to, or nil for a path that is no file of this repo's working tree.
func (self *WorkingTreeDiffActions) fileForDiffLinePath(path string) *models.File {
	relativePath := repoRelativePath(self.c.Git().RepoPaths.WorktreePath(), path)
	if relativePath == "" {
		return nil
	}
	return self.context().FileTreeViewModel.GetFile(relativePath)
}

// diffLinesPatch is how the selected change lines of one file are applied.
type diffLinesPatch struct {
	file *models.File

	// Whether the lines are every change of the file's diff. Acting on all of them is
	// acting on the file itself, which is staged or unstaged as a whole.
	wholeFile bool

	// The patch of the lines, where they are not the whole file.
	patch string

	// Whether the diff the lines were found in holds changes they don't cover.
	changesLeft bool
}

// buildDiffLinesPatch builds what applies the given change lines of one file — a line,
// a hunk, a range — as a patch built from that file's own diff:
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
//
// It records whether the diff it read holds changes the selection didn't cover. The
// caller uses this to tell whether the side acted on still has anything of this file
// in it once we are done.
//
// It fails if the diff it reads is no longer the one on screen (see readFileDiff), and
// if none of the lines are in it.
func (self *WorkingTreeDiffActions) buildDiffLinesPatch(
	fileLines fileDiffLines, sourceCached bool, opts git_commands.ApplyPatchOpts,
) (diffLinesPatch, error) {
	file := fileLines.file
	parsedPatch, err := self.readFileDiff(fileLines, sourceCached)
	if err != nil {
		return diffLinesPatch{}, err
	}

	patchLineIndices := patch.ChangeLineIndicesForLines(parsedPatch,
		lo.Map(fileLines.lines, func(info types.DiffLineInfo, _ int) patch.LineIdentity {
			return info.PatchLineIdentity()
		}))
	if len(patchLineIndices) == 0 {
		return diffLinesPatch{}, diffLinesNotFoundError(self.c, file.GetPath())
	}

	changesLeft := len(patchLineIndices) < changeLineCount(parsedPatch)

	// Acting on every change of a file is acting on the file itself, and saying so is
	// not the same as applying its diff. The diff of a deleted file is its content
	// going away, and putting that into the index line by line leaves an empty file
	// there rather than the deletion; the diff of an added one is its whole content,
	// and taking that back out leaves an empty file in the index rather than an
	// untracked one.
	if !changesLeft && opts.Cached {
		return diffLinesPatch{file: file, wholeFile: true}, nil
	}

	patchToApply := parsedPatch.
		Transform(patch.TransformOpts{
			Reverse:             opts.Reverse,
			IncludedLineIndices: patchLineIndices,
			FileNameOverride:    file.GetPath(),
		}).
		FormatPlain()
	return diffLinesPatch{file: file, patch: patchToApply, changesLeft: changesLeft}, nil
}

// readFileDiff reads the diff of the file the given lines belong to, from the side of
// the index they were selected on: the diff to look them up in. It fails if this is no
// longer the diff that was on screen when they were selected (see shownDiff).
func (self *WorkingTreeDiffActions) readFileDiff(fileLines fileDiffLines, onStagedSide bool) (*patch.Patch, error) {
	parsedPatch := patch.Parse(
		self.c.Git().WorkingTree.WorktreeFileDiff(fileLines.file, git_commands.DiffModePlain, onStagedSide))
	if err := self.checkShown(fileLines.shown, fileLines.file.Names(), onStagedSide, parsedPatch.FileVersions()); err != nil {
		return nil, err
	}
	return parsedPatch, nil
}

// applyDiffLinesPatch applies what buildDiffLinesPatch built, as opts says.
func (self *WorkingTreeDiffActions) applyDiffLinesPatch(p diffLinesPatch, opts git_commands.ApplyPatchOpts) error {
	if p.wholeFile {
		if opts.Reverse {
			return self.c.Git().WorkingTree.UnStageFile(p.file.Names(), p.file.Tracked)
		}
		return self.c.Git().WorkingTree.StageFile(p.file.GetPath())
	}
	return self.c.Git().Patch.ApplyPatch(p.patch, opts)
}

// changeLineCount returns how many of a patch's lines are changes rather than context
// or header. A selection of the whole diff covers exactly that many.
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
	pane types.DiffPaneContext, firstBufferLine int, actedSideRemains bool, otherSideHasChanges bool,
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
	self.c.Helpers().DiffLine.RevealSelectionAfterAction(pane, target, firstBufferLine, 0,
		self.c.GocuiGui().EndBlockingEvents)
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
