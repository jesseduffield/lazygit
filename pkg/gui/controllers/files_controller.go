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
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type FilesController struct {
	baseController
	*ListControllerTrait[*filetree.FileNode]
	c *ControllerCommon
}

var _ types.IController = &FilesController{}

func NewFilesController(
	c *ControllerCommon,
) *FilesController {
	return &FilesController{
		c: c,
		ListControllerTrait: NewListControllerTrait(
			c,
			c.Contexts().Files,
			c.Contexts().Files.GetSelected,
			c.Contexts().Files.GetSelectedItems,
		),
	}
}

func (self *FilesController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	return []*types.Binding{
		{
			Keys:              opts.GetKeys(opts.Config.Universal.Select),
			Handler:           self.withItems(self.press),
			GetDisabledReason: self.require(self.withFileTreeViewModelMutex(self.itemsSelected(self.canStageSelection))),
			Description:       self.c.Tr.Stage,
			Tooltip:           self.c.Tr.StageTooltip,
			DisplayOnScreen:   true,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Files.OpenStatusFilter),
			Handler:     self.handleStatusFilterPressed,
			Description: self.c.Tr.FileFilter,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Files.CopyFileInfoToClipboard),
			Handler:     self.openCopyMenu,
			Description: self.c.Tr.CopyToClipboardMenu,
			OpensMenu:   true,
		},
		{
			Keys:            opts.GetKeys(opts.Config.Files.CommitChanges),
			Handler:         self.c.Helpers().WorkingTree.HandleCommitPress,
			Description:     self.c.Tr.Commit,
			Tooltip:         self.c.Tr.CommitTooltip,
			DisplayOnScreen: true,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Files.CommitChangesWithoutHook),
			Handler:     self.c.Helpers().WorkingTree.HandleWIPCommitPress,
			Description: self.c.Tr.CommitChangesWithoutHook,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Files.AmendLastCommit),
			Handler:     self.handleAmendCommitPress,
			Description: self.c.Tr.AmendLastCommit,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Files.CommitChangesWithEditor),
			Handler:     self.c.Helpers().WorkingTree.HandleCommitEditorPress,
			Description: self.c.Tr.CommitChangesWithEditor,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Files.FindBaseCommitForFixup),
			Handler:     self.c.Helpers().FixupHelper.HandleFindBaseCommitForFixupPress,
			Description: self.c.Tr.FindBaseCommitForFixup,
			Tooltip:     self.c.Tr.FindBaseCommitForFixupTooltip,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Universal.Edit),
			Handler:           self.withItems(self.edit),
			GetDisabledReason: self.require(self.withFileTreeViewModelMutex(self.itemsSelected(self.canEditFiles))),
			Description:       self.c.Tr.Edit,
			Tooltip:           self.c.Tr.EditFileTooltip,
			DisplayOnScreen:   true,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Universal.OpenFile),
			Handler:           self.Open,
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.OpenFile,
			Tooltip:           self.c.Tr.OpenFileTooltip,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Files.IgnoreFile),
			Handler:           self.withItem(self.ignoreOrExcludeMenu),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Actions.IgnoreExcludeFile,
			OpensMenu:         true,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Files.RefreshFiles),
			Handler:     self.refresh,
			Description: self.c.Tr.RefreshFiles,
		},
		{
			Keys:            opts.GetKeys(opts.Config.Files.StashAllChanges),
			Handler:         self.stash,
			Description:     self.c.Tr.Stash,
			Tooltip:         self.c.Tr.StashTooltip,
			DisplayOnScreen: true,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Files.ViewStashOptions),
			Handler:     self.createStashMenu,
			Description: self.c.Tr.ViewStashOptions,
			Tooltip:     self.c.Tr.ViewStashOptionsTooltip,
			OpensMenu:   true,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Files.ToggleStagedAll),
			Handler:     self.toggleStagedAll,
			Description: self.c.Tr.ToggleStagedAll,
			Tooltip:     self.c.Tr.ToggleStagedAllTooltip,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Universal.GoInto),
			Handler:           self.enter,
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.FileEnter,
			Tooltip:           self.c.Tr.FileEnterTooltip,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Universal.Remove),
			Handler:           self.withItems(self.remove),
			GetDisabledReason: self.withFileTreeViewModelMutex(self.require(self.itemsSelected(self.canRemove))),
			Description:       self.c.Tr.Discard,
			Tooltip:           self.c.Tr.DiscardFileChangesTooltip,
			OpensMenu:         true,
			DisplayOnScreen:   true,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Commits.ViewResetOptions),
			Handler:     self.createResetToUpstreamMenu,
			Description: self.c.Tr.ViewResetToUpstreamOptions,
			OpensMenu:   true,
		},
		{
			Keys:            opts.GetKeys(opts.Config.Files.ViewResetOptions),
			Handler:         self.createResetMenu,
			Description:     self.c.Tr.Reset,
			Tooltip:         self.c.Tr.FileResetOptionsTooltip,
			OpensMenu:       true,
			DisplayOnScreen: true,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Files.ToggleTreeView),
			Handler:     self.toggleTreeView,
			Description: self.c.Tr.ToggleTreeView,
			Tooltip:     self.c.Tr.ToggleTreeViewTooltip,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Universal.OpenDiffTool),
			Handler:           self.withItem(self.openDiffTool),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.OpenDiffTool,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Files.OpenMergeOptions),
			Handler:           self.withItems(self.openMergeConflictMenu),
			Description:       self.c.Tr.ViewMergeConflictOptions,
			Tooltip:           self.c.Tr.ViewMergeConflictOptionsTooltip,
			GetDisabledReason: self.require(self.withFileTreeViewModelMutex(self.itemsSelected(self.canOpenMergeConflictMenu))),
			OpensMenu:         true,
			DisplayOnScreen:   true,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Files.Fetch),
			Handler:     self.fetch,
			Description: self.c.Tr.Fetch,
			Tooltip:     self.c.Tr.FetchTooltip,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Files.CollapseAll),
			Handler:           self.collapseAll,
			Description:       self.c.Tr.CollapseAll,
			Tooltip:           self.c.Tr.CollapseAllTooltip,
			GetDisabledReason: self.require(self.isInTreeMode),
		},
		{
			Keys:              opts.GetKeys(opts.Config.Files.ExpandAll),
			Handler:           self.expandAll,
			Description:       self.c.Tr.ExpandAll,
			Tooltip:           self.c.Tr.ExpandAllTooltip,
			GetDisabledReason: self.require(self.isInTreeMode),
		},
	}
}

func (self *FilesController) withFileTreeViewModelMutex(callback func() *types.DisabledReason) func() *types.DisabledReason {
	return func() *types.DisabledReason {
		self.c.Contexts().Files.FileTreeViewModel.RWMutex.RLock()
		defer self.c.Contexts().Files.FileTreeViewModel.RWMutex.RUnlock()

		return callback()
	}
}

func (self *FilesController) GetMouseKeybindings(opts types.KeybindingsOpts) []*gocui.ViewMouseBinding {
	return []*gocui.ViewMouseBinding{
		{
			ViewName:    "mergeConflicts",
			Key:         gocui.MouseLeft,
			Handler:     self.onClickMain,
			FocusedView: self.context().GetViewName(),
		},
	}
}

func (self *FilesController) GetOnClick() func(opts gocui.ViewMouseBindingOpts) error {
	return func(opts gocui.ViewMouseBindingOpts) error {
		clickedIdx := self.context().GetSelectedLineIdx()
		node := self.context().FileTreeViewModel.Get(clickedIdx)
		if node == nil || node.File != nil {
			return nil
		}

		// The arrow is at column visualDepth*2 (after indentation of 2 spaces per level).
		// Only treat clicks on the arrow and the trailing space as arrow clicks.
		visualDepth := self.context().FileTreeViewModel.GetVisualDepth(clickedIdx)
		arrowStartCol := visualDepth * 2
		arrowEndCol := arrowStartCol + 1
		if opts.X < arrowStartCol || opts.X > arrowEndCol {
			return nil
		}

		self.context().FileTreeViewModel.ToggleCollapsed(node.GetInternalPath())
		self.c.PostRefreshUpdate(self.context())

		return nil
	}
}

func (self *FilesController) GetOnRenderToMain() func() {
	return func() {
		self.c.Helpers().Diff.WithDiffModeCheck(func() {
			node := self.context().GetSelected()

			if node == nil {
				self.renderToMainWithTask(types.NewRenderStringTask(self.c.Tr.NoChangedFiles))
				return
			}

			if self.isSubmoduleCommitConflict(node.File) {
				self.renderSubmoduleConflict(node)
				return
			}

			if node.File != nil && node.File.HasInlineMergeConflicts {
				if self.renderInlineMergeConflict(node) {
					return
				}
				// The file is marked as conflicted but has no conflict markers (it
				// was resolved in an editor), so fall through to show its diff.
			} else if node.File != nil && node.File.HasMergeConflicts {
				self.renderNonTextualConflict(node)
				return
			}

			self.renderWorkingTreeDiff(node)
		})
	}
}

// renderToMainWithTask renders the given task to the main view with the standard
// diff title and subtitle.
func (self *FilesController) renderToMainWithTask(task types.UpdateTask) {
	self.c.RenderToMainViews(types.RefreshMainOpts{
		Pair: self.c.MainViewPairs().Normal,
		Main: &types.ViewUpdateOpts{
			Title:    self.c.Tr.DiffTitle,
			SubTitle: self.c.Helpers().Diff.IgnoringWhitespaceSubTitle(),
			Task:     task,
		},
	})
}

// renderSubmoduleConflict shows, for a conflicted submodule, an explanation plus
// the commits each side added relative to their common ancestor as two separate,
// indented logs. If a side added nothing of its own (e.g. it was rewound to an
// ancestor of the other), the commit it points at is shown instead.
func (self *FilesController) renderSubmoduleConflict(node *filetree.FileNode) {
	self.c.Helpers().MergeConflicts.ResetMergeState()

	path := node.GetPath()
	_, ours, theirs, err := self.c.Git().Submodule.GetConflictCommits(path)
	if err != nil {
		return
	}

	sideBlock := func(header string, side string, otherSide string) string {
		log, err := self.c.Git().Submodule.ConflictSideLog(path, side, otherSide)
		if err != nil {
			return header
		}
		if log = strings.TrimRight(log, "\n"); log == "" {
			if log, err = self.c.Git().Submodule.GetCommitSummary(path, side); err != nil {
				return header
			}
		}
		return header + "\n\n  " + strings.ReplaceAll(log, "\n", "\n  ")
	}

	message := strings.Join([]string{
		self.conflictResolutionHint(utils.ResolvePlaceholderString(self.c.Tr.SubmoduleMergeConflictDescription, map[string]string{"path": path})),
		sideBlock(self.c.Tr.MergeConflictCurrentDiff, ours, theirs),
		sideBlock(self.c.Tr.MergeConflictIncomingDiff, theirs, ours),
	}, "\n\n")

	self.renderToMainWithTask(types.NewRenderStringTask(message))
}

// renderInlineMergeConflict renders the merge-conflict view for a file with
// inline conflict markers. It returns false if the file has no actual markers
// (it was resolved in an editor), in which case the caller should fall back to
// showing the file's diff.
func (self *FilesController) renderInlineMergeConflict(node *filetree.FileNode) bool {
	hasConflicts, err := self.c.Helpers().MergeConflicts.SetMergeState(node.GetPath())
	if err != nil {
		return true
	}

	if !hasConflicts {
		return false
	}

	self.c.Helpers().MergeConflicts.Render()
	return true
}

// renderNonTextualConflict shows the resolution hint for a non-textual text-file
// conflict (DD/AU/UA/UD/DU), plus the base diff for the modify/delete cases.
func (self *FilesController) renderNonTextualConflict(node *filetree.FileNode) {
	message := self.conflictResolutionHint(node.File.GetMergeStateDescription(self.c.Tr))

	if node.File.ShortStatus == "DU" || node.File.ShortStatus == "UD" {
		cmdObj := self.c.Git().Diff.DiffCmdObj([]string{"--base", "--", node.GetPath()})
		prefix := message + "\n\n"
		if node.File.ShortStatus == "DU" {
			prefix += self.c.Tr.MergeConflictIncomingDiff
		} else {
			prefix += self.c.Tr.MergeConflictCurrentDiff
		}
		prefix += "\n\n"
		self.renderToMainWithTask(types.NewRunPtyTaskWithPrefix(cmdObj.GetCmd(), prefix))
		return
	}

	self.renderToMainWithTask(types.NewRenderStringTask(message))
}

func (self *FilesController) renderWorkingTreeDiff(node *filetree.FileNode) {
	self.c.Helpers().MergeConflicts.ResetMergeState()

	split, mainShowsStaged := self.diffSplitState(node)

	pathOverrides := self.pathOverridesForDiff(node)
	cmdObj := self.c.Git().WorkingTree.WorktreeFileDiffCmdObj(node, false, mainShowsStaged, pathOverrides)
	title := self.c.Tr.UnstagedChanges
	if mainShowsStaged {
		title = self.c.Tr.StagedChanges
	}
	refreshOpts := types.RefreshMainOpts{
		Pair: self.c.MainViewPairs().Normal,
		Main: &types.ViewUpdateOpts{
			Task:     types.NewRunPtyTask(cmdObj.GetCmd()),
			SubTitle: self.c.Helpers().Diff.IgnoringWhitespaceSubTitle(),
			Title:    title,
		},
	}

	if split {
		cmdObj := self.c.Git().WorkingTree.WorktreeFileDiffCmdObj(node, false, true, pathOverrides)

		title := self.c.Tr.StagedChanges
		if mainShowsStaged {
			title = self.c.Tr.UnstagedChanges
		}

		refreshOpts.Secondary = &types.ViewUpdateOpts{
			Title:    title,
			SubTitle: self.c.Helpers().Diff.IgnoringWhitespaceSubTitle(),
			Task:     types.NewRunPtyTask(cmdObj.GetCmd()),
		}
	}

	self.c.RenderToMainViews(refreshOpts)
}

func (self *FilesController) GetOnDoubleClick() func() error {
	return self.withItemGraceful(func(node *filetree.FileNode) error {
		return self.press([]*filetree.FileNode{node})
	})
}

func (self *FilesController) GetFocusedMainViewActions() types.FocusedMainViewActions {
	return self
}

func (self *FilesController) OnClick(mainViewName string, clickedLineIdx int) error {
	// Capture before any mutation below that might re-render the main view.
	snapshot := focusedMainViewSnapshot(self.c, mainViewName, self.context())

	info, ok := self.c.Helpers().Staging.GetDiffLineInfo(mainViewName, clickedLineIdx)
	line := -1
	isDeletion := false
	if ok {
		line, isDeletion = info.PatchSelectLine()
	}

	node := self.context().GetSelected()
	if node == nil {
		return nil
	}

	if !node.IsFile() && ok {
		relativePath, err := filepath.Rel(self.c.Git().RepoPaths.WorktreePath(), info.Path)
		if err != nil {
			return err
		}
		relativePath = "./" + relativePath
		self.context().FileTreeViewModel.ExpandToPath(relativePath)
		self.c.PostRefreshUpdate(self.context())

		idx, ok := self.context().FileTreeViewModel.GetIndexForPath(relativePath)
		if ok {
			self.context().SetSelectedLineIdx(idx)
			self.context().GetViewTrait().FocusPoint(
				self.context().ModelIndexToViewIndex(idx), false)
		}
	}

	return self.EnterFile(snapshot, types.OnFocusOpts{ClickedWindowName: mainViewName, ClickedViewLineIdx: line, ClickedViewRealLineIdx: line, ClickedViewRealLineIsDeletion: isDeletion, SelectLineInDefaultMode: true})
}

// diffSplitState reports, for the given file node, how the focused main view lays
// out its diff: whether it's split into unstaged (Normal) and staged
// (NormalSecondary) halves, and — when not split — whether the single Normal view
// shows the staged diff (which happens when the file has only staged changes).
// GetOnRenderToMain and PrimaryAction share this so the staging direction
// can't drift from what's on screen.
func (self *FilesController) diffSplitState(node *filetree.FileNode) (split bool, mainShowsStaged bool) {
	split = self.c.UserConfig().Gui.SplitDiff == "always" || (node.GetHasUnstagedChanges() && node.GetHasStagedChanges())
	mainShowsStaged = !split && node.GetHasStagedChanges()
	return split, mainShowsStaged
}

// PrimaryAction stages (or unstages) the selected diff line(s) when space is pressed in
// the focused main view of the working-tree files panel.
func (self *FilesController) PrimaryAction(mainViewName string, firstLineIdx int, lastLineIdx int) error {
	if self.c.UserConfig().Git.DiffContextSize == 0 {
		return fmt.Errorf(self.c.Tr.Actions.NotEnoughContextToStage,
			self.c.UserConfig().Keybinding.Universal.IncreaseContextInDiffView)
	}

	node := self.context().GetSelected()
	if node == nil {
		return nil
	}

	infos := self.c.Helpers().Staging.ChangeLinesInViewRange(mainViewName, firstLineIdx, lastLineIdx)
	if len(infos) == 0 {
		return nil
	}

	// The whole diff shown in the main view is on one side — the staged diff in
	// the secondary half of a split, and in the main half when there are only
	// staged changes; in those cases space unstages, otherwise it stages. The
	// direction is the same for every file in a multi-file (directory) diff.
	_, mainShowsStaged := self.diffSplitState(node)
	reverse := mainShowsStaged || mainViewName == self.c.Contexts().NormalSecondary.GetViewName()

	// A directory diff spans several files; group the selected change lines by
	// file and apply one patch per file.
	infosByFile := lo.GroupBy(infos, func(info types.DiffLineInfo) string { return info.Path })

	self.c.LogAction(self.c.Tr.Actions.ApplyPatch)
	for path, fileInfos := range infosByFile {
		file := self.fileForDiffLinePath(path)
		if file == nil {
			continue
		}
		// Staging reads the side being acted on (unstaged when staging, staged when
		// unstaging — both = reverse here) and applies to the index either way.
		if err := self.applyDiffLines(file, fileInfos, reverse, git_commands.ApplyPatchOpts{Reverse: reverse, Cached: true}); err != nil {
			return err
		}
	}

	self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES, types.STAGING}})

	// Focus follows the side that was acted on. Staging keeps it in the main half
	// (which always holds the unstaged side, or the staged side once the file has
	// only staged changes). Unstaging keeps it on the staged side, which lives in
	// the secondary half once the file is split into staged + unstaged, and moves
	// back to the main half when the staged side empties and the split collapses.
	// The model is up to date now (Refresh above is synchronous), so the post-op
	// split is read from the freshly selected node.
	focusViewName := self.c.Contexts().Normal.GetViewName()
	if reverse {
		if node := self.context().GetSelected(); node != nil {
			if split, _ := self.diffSplitState(node); split {
				focusViewName = self.c.Contexts().NormalSecondary.GetViewName()
			}
		}
	}

	// The staging Refresh above queued the main-view re-render; re-establish the
	// selection in whichever pane now holds the acted-on side once that render lands,
	// and focus that pane if staging moved it there.
	revealSelectionAfterPrimaryAction(self.c, mainViewName, focusViewName, firstLineIdx)
	if focusViewName != mainViewName {
		self.c.Context().Push(mainContextForViewName(self.c, focusViewName), types.OnFocusOpts{})
	}
	return nil
}

// fileForDiffLinePath maps a diff line's absolute file path (as carried by the
// diff-line metadata) to the working-tree file it belongs to, or nil if it isn't a
// tracked working-tree file.
func (self *FilesController) fileForDiffLinePath(path string) *models.File {
	relativePath, err := filepath.Rel(self.c.Git().RepoPaths.WorktreePath(), path)
	if err != nil {
		return nil
	}
	return self.context().FileTreeViewModel.GetFile(filepath.ToSlash(relativePath))
}

// applyDiffLines applies, to the diff lines identified by infos (a single line, a
// range, or a hunk) — all belonging to file — a patch built from the file's staged or
// unstaged diff (sourceCached: the side the selection was made on) and applied per opts:
//   - stage:   read unstaged, apply forward to the index    (sourceCached=false, Reverse=false, Cached=true)
//   - unstage: read staged,   apply reverse to the index    (sourceCached=true,  Reverse=true,  Cached=true)
//   - discard: read the shown side, apply reverse — not cached on the unstaged side
//     removes the change from the working tree, cached on the staged side just unstages it.
//
// The read side and the apply direction are independent (they only coincide for
// staging/unstaging), so they're passed separately: the patch lines are matched in the
// diff actually shown, then Transform reverses them to match opts. A selection covering
// no change lines yields an empty patch and is a no-op. The caller logs the action and
// refreshes, once, around the (possibly several) files it touches.
//
// Each selected change line is keyed by its (file line number, deletion?) identity,
// and the freshly parsed patch is scanned for the body lines matching those
// identities. Scanning the patch and matching identities — rather than looking each
// line number up with PatchLineForLineNumber — is what makes a modified line work: a
// deletion and the addition replacing it share a position but have distinct
// identities, and the line-number lookup can't tell an addition at the start of a
// hunk from the deletion above it.
func (self *FilesController) applyDiffLines(file *models.File, infos []types.DiffLineInfo, sourceCached bool, opts git_commands.ApplyPatchOpts) error {
	parsedPatch := patch.Parse(self.c.Git().WorkingTree.WorktreeFileDiff(file, true, sourceCached))

	type changeLineKey struct {
		lineNumber int
		isDeletion bool
	}
	selected := make(map[changeLineKey]bool, len(infos))
	for _, info := range infos {
		lineNumber, isDeletion := info.PatchSelectLine()
		selected[changeLineKey{lineNumber, isDeletion}] = true
	}

	var patchLineIndices []int
	for idx, line := range parsedPatch.Lines() {
		var key changeLineKey
		switch {
		case line.IsAddition():
			key = changeLineKey{parsedPatch.LineNumberOfLine(idx), false}
		case line.IsDeletion():
			key = changeLineKey{parsedPatch.OldLineNumberOfLine(idx), true}
		default:
			continue
		}
		if selected[key] {
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

// if we are dealing with a status for which there is no key in this map,
// then we won't optimistically render: we'll just let `git status` tell
// us what the new status is.
// There are no doubt more entries that could be added to these two maps.
var stageStatusMap = map[string]string{
	"??": "A ",
	" M": "M ",
	"MM": "M ",
	" D": "D ",
	" A": "A ",
	"AM": "A ",
	"MD": "D ",
}

var unstageStatusMap = map[string]string{
	"A ": "??",
	"M ": " M",
	"D ": " D",
	// A submodule with both a staged commit and unstageable dirty content; the
	// staged commit gets unstaged, the dirty content stays.
	"MM": " M",
}

func (self *FilesController) optimisticStage(file *models.File) bool {
	newShortStatus, ok := stageStatusMap[file.ShortStatus]
	if !ok {
		return false
	}

	models.SetStatusFields(file, newShortStatus)
	return true
}

func (self *FilesController) optimisticUnstage(file *models.File) bool {
	newShortStatus, ok := unstageStatusMap[file.ShortStatus]
	if !ok {
		return false
	}

	models.SetStatusFields(file, newShortStatus)
	return true
}

// Running a git add command followed by a git status command can take some time (e.g. 200ms).
// Given how often users stage/unstage files in Lazygit, we're adding some
// optimistic rendering to make things feel faster. When we go to stage
// a file, we'll first update that file's status in-memory, then re-render
// the files panel. Then we'll immediately do a proper git status call
// so that if the optimistic rendering got something wrong, it's quickly
// corrected.
func (self *FilesController) optimisticChange(nodes []*filetree.FileNode, optimisticChangeFn func(*models.File) bool) error {
	rerender := false

	for _, node := range nodes {
		err := node.ForEachFile(func(f *models.File) error {
			// can't act on the file itself: we need to update the original model file
			for _, modelFile := range self.c.Model().Files {
				if modelFile.Path == f.Path {
					if optimisticChangeFn(modelFile) {
						rerender = true
					}
					break
				}
			}

			return nil
		})
		if err != nil {
			return err
		}
	}

	if rerender {
		self.c.PostRefreshUpdate(self.c.Contexts().Files)
	}

	return nil
}

// toggleStaged decides whether to stage or unstage the given nodes, updates the
// model optimistically, and then runs the matching git command via the supplied
// callbacks. press() (acting on the selection) and toggleStagedAll() (acting on
// the whole tree) share this; they differ only in the git commands they run,
// which is why those are passed in.
//
// If any node has unstaged changes we stage the nodes that have them (staging
// already-staged deleted files/folders would fail); otherwise we unstage all
// the nodes.
func (self *FilesController) toggleStaged(
	nodes []*filetree.FileNode,
	stageAction string,
	unstageAction string,
	stage func(unstagedNodes []*filetree.FileNode) error,
	unstage func(nodes []*filetree.FileNode) error,
) error {
	for _, node := range nodes {
		// if any files within have inline merge conflicts we can't stage or unstage,
		// or it'll end up with those >>>>>> lines actually staged
		if node.GetHasInlineMergeConflicts() {
			return errors.New(self.c.Tr.ErrStageDirWithInlineMergeConflicts)
		}
	}

	nodes = normalisedSelectedNodes(nodes)

	unstagedNodes := filterNodesHaveUnstagedChanges(nodes, self.c.Model().Submodules)

	// Staging a submodule that only has dirty or untracked content (no new
	// commit) is a no-op: the parent repo can't stage that content. When that's
	// the only thing that looks stageable, don't stage; fall through to
	// unstaging instead. That keeps the toggle symmetric (e.g. a fully-staged
	// tree that also contains a dirty submodule still unstages on the next
	// press) rather than getting stuck trying to stage the unstageable content.
	shouldStage := len(unstagedNodes) > 0
	if shouldStage {
		noOp, err := self.stagingWouldBeNoOp(unstagedNodes)
		if err != nil {
			return err
		}
		shouldStage = !noOp
	}

	if shouldStage {
		self.c.LogAction(stageAction)

		if err := self.optimisticChange(unstagedNodes, self.optimisticStage); err != nil {
			return err
		}

		return stage(unstagedNodes)
	}

	// If there's nothing staged to unstage either, then the only thing we acted
	// on was an unstageable submodule and nothing happened, so say why.
	if !someNodesHaveStagedChanges(nodes) {
		return errors.New(self.c.Tr.NothingToStageForSubmodule)
	}

	self.c.LogAction(unstageAction)

	if err := self.optimisticChange(nodes, self.optimisticUnstage); err != nil {
		return err
	}

	return unstage(nodes)
}

func (self *FilesController) pressWithLock(selectedNodes []*filetree.FileNode) error {
	// Obtaining this lock because optimistic rendering requires us to mutate
	// the files in our model.
	self.c.Mutexes().RefreshingFilesMutex.Lock()
	defer self.c.Mutexes().RefreshingFilesMutex.Unlock()

	// When filtering, expand directory nodes to individual visible file paths
	// so that only filtered files are staged/unstaged.
	toPaths := func(nodes []*filetree.FileNode) []string {
		if self.context().IsFiltering() {
			var paths []string
			for _, node := range nodes {
				_ = node.ForEachFile(func(file *models.File) error {
					paths = append(paths, file.Path)
					return nil
				})
			}
			return paths
		}
		return lo.Map(nodes, func(node *filetree.FileNode, _ int) string {
			return node.GetPath()
		})
	}

	stage := func(unstagedNodes []*filetree.FileNode) error {
		var extraArgs []string
		if self.context().GetStatusFilter() == filetree.DisplayTracked {
			extraArgs = []string{"-u"}
		}

		return self.c.Git().WorkingTree.StageFiles(toPaths(unstagedNodes), extraArgs)
	}

	unstage := func(nodes []*filetree.FileNode) error {
		if self.context().IsFiltering() {
			// When filtering, only unstage visible files
			return self.unstageFilteredFiles(nodes)
		}

		// need to partition the paths into tracked and untracked (where we assume directories are tracked). Then we'll run the commands separately.
		trackedNodes, untrackedNodes := utils.Partition(nodes, func(node *filetree.FileNode) bool {
			// We treat all directories as tracked. I'm not actually sure why we do this but
			// it's been the existing behaviour for a while and nobody has complained
			return !node.IsFile() || node.GetIsTracked()
		})

		if len(untrackedNodes) > 0 {
			if err := self.c.Git().WorkingTree.UnstageUntrackedFiles(toPaths(untrackedNodes)); err != nil {
				return err
			}
		}

		if len(trackedNodes) > 0 {
			if err := self.c.Git().WorkingTree.UnstageTrackedFiles(toPaths(trackedNodes)); err != nil {
				return err
			}
		}

		return nil
	}

	return self.toggleStaged(selectedNodes,
		self.c.Tr.Actions.StageFile, self.c.Tr.Actions.UnstageFile,
		stage, unstage)
}

func (self *FilesController) press(nodes []*filetree.FileNode) error {
	// A single file with a conflict that can only be resolved through a dialog
	// can't be staged; route it to the same picker that `enter` uses instead.
	if len(nodes) == 1 && self.conflictNeedsResolutionDialog(nodes[0].File) {
		return self.openConflictResolutionMenu(nodes[0].File)
	}

	if err := self.pressWithLock(nodes); err != nil {
		return err
	}

	self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}, Mode: types.ASYNC})

	self.context().HandleFocus(types.OnFocusOpts{})
	return nil
}

// pathOverridesForDiff returns file paths to override the node's path in diff
// commands when a text filter is active and the node is a directory. This
// ensures the diff only shows filtered/visible files.
func (self *FilesController) pathOverridesForDiff(node *filetree.FileNode) []string {
	if !node.IsFile() && self.context().IsFiltering() {
		var paths []string
		_ = node.ForEachFile(func(file *models.File) error {
			paths = append(paths, file.Path)
			return nil
		})
		return paths
	}
	return nil
}

// unstageFilteredFiles unstages only the visible (filtered) files from the
// given nodes, correctly partitioning by tracked/untracked.
func (self *FilesController) unstageFilteredFiles(nodes []*filetree.FileNode) error {
	var trackedPaths, untrackedPaths []string
	for _, node := range nodes {
		_ = node.ForEachFile(func(file *models.File) error {
			if file.Tracked || file.HasStagedChanges {
				trackedPaths = append(trackedPaths, file.Path)
			} else {
				untrackedPaths = append(untrackedPaths, file.Path)
			}
			return nil
		})
	}
	if len(untrackedPaths) > 0 {
		if err := self.c.Git().WorkingTree.UnstageUntrackedFiles(untrackedPaths); err != nil {
			return err
		}
	}
	if len(trackedPaths) > 0 {
		if err := self.c.Git().WorkingTree.UnstageTrackedFiles(trackedPaths); err != nil {
			return err
		}
	}
	return nil
}

func (self *FilesController) Context() types.Context {
	return self.context()
}

func (self *FilesController) context() *context.WorkingTreeContext {
	return self.c.Contexts().Files
}

func (self *FilesController) getSelectedFile() *models.File {
	node := self.context().GetSelected()
	if node == nil {
		return nil
	}
	return node.File
}

func (self *FilesController) enter() error {
	return self.EnterFile(nil, types.OnFocusOpts{ClickedWindowName: "", ClickedViewLineIdx: -1, ClickedViewRealLineIdx: -1})
}

func (self *FilesController) collapseAll() error {
	self.context().FileTreeViewModel.CollapseAll()

	self.c.PostRefreshUpdate(self.context())

	return nil
}

func (self *FilesController) expandAll() error {
	self.context().FileTreeViewModel.ExpandAll()

	self.c.PostRefreshUpdate(self.context())

	return nil
}

// focusedMainViewSnapshot records the focused main view to return to when
// escaping the staging view, for the case where we're entering it straight from
// there; it's nil for the normal flow that goes through the files panel. See
// types.FocusedMainViewSnapshot.
func (self *FilesController) EnterFile(focusedMainViewSnapshot *types.FocusedMainViewSnapshot, opts types.OnFocusOpts) error {
	node := self.context().GetSelected()
	if node == nil {
		return nil
	}

	if node.File == nil {
		return self.handleToggleDirCollapsed()
	}

	file := node.File

	if self.conflictNeedsResolutionDialog(file) {
		return self.openConflictResolutionMenu(file)
	}

	submoduleConfigs := self.c.Model().Submodules
	if file.IsSubmodule(submoduleConfigs) {
		submoduleConfig := file.SubmoduleConfig(submoduleConfigs)
		return self.c.Helpers().Repos.EnterSubmodule(submoduleConfig)
	}

	if file.HasInlineMergeConflicts {
		return self.switchToMerge()
	}

	context := lo.Ternary(opts.ClickedWindowName == "secondary", self.c.Contexts().StagingSecondary, self.c.Contexts().Staging)
	// Record the focused-main-view return on *both* staging halves, not just the
	// one we're entering: staging the last unstaged hunk (or tabbing) moves the
	// selection to the other half, and escaping from there must still know to
	// return to the focused main view rather than falling back to the files panel.
	// Set on every entry (nil for a normal entry through the files panel) so a
	// snapshot can't leak from a previous main-view entry into a subsequent normal one.
	self.c.Contexts().Staging.SetFocusedMainViewSnapshot(focusedMainViewSnapshot)
	self.c.Contexts().StagingSecondary.SetFocusedMainViewSnapshot(focusedMainViewSnapshot)
	self.c.Context().Push(context, opts)
	self.c.Helpers().PatchBuilding.ShowHunkStagingHint()

	return nil
}

// conflictResolutionHint formats a conflict description for the main view,
// appending the "press <enter> to resolve" hint and wrapping it when the view is
// wide enough that long lines would otherwise hurt readability.
func (self *FilesController) conflictResolutionHint(description string) string {
	message := description + "\n\n" + fmt.Sprintf(self.c.Tr.MergeConflictPressEnterToResolve,
		self.c.UserConfig().Keybinding.Universal.GoInto)
	if self.c.Views().Main.InnerWidth() > 70 {
		lines, _, _ := utils.WrapViewLinesToWidth(true, false, message, 70, 4)
		message = strings.Join(lines, "\n")
	}
	return message
}

// conflictNeedsResolutionDialog reports whether a file's merge conflict can only
// be resolved through a dialog that picks one side, as opposed to editing
// conflict markers in the merge view. These are the "non-textual" conflicts:
// text files where one side modified and the other deleted/renamed the file
// (DD/AU/UA/UD/DU), and submodules where both sides moved the gitlink (UU).
func (self *FilesController) conflictNeedsResolutionDialog(file *models.File) bool {
	if file == nil || !file.HasMergeConflicts {
		return false
	}

	// A conflicted submodule has no conflict markers to edit; it's resolved by
	// picking which commit to point at.
	if file.IsSubmodule(self.c.Model().Submodules) {
		return true
	}

	return !file.HasInlineMergeConflicts
}

// canStageSelection disables staging when a multiple selection includes a file
// with a conflict that must be resolved through a dialog; those have to be
// resolved one at a time.
func (self *FilesController) canStageSelection(nodes []*filetree.FileNode) *types.DisabledReason {
	if len(nodes) > 1 {
		for _, node := range nodes {
			if node.SomeFile(self.conflictNeedsResolutionDialog) {
				return &types.DisabledReason{
					Text: utils.ResolvePlaceholderString(
						self.c.Tr.StageConflictsRangeDisabled, map[string]string{
							"goIntoKey": self.c.UserConfig().Keybinding.Universal.GoInto.String(),
						},
					),
				}
			}
		}
	}

	return nil
}

// isSubmoduleCommitConflict reports whether the file is a submodule whose commit
// pointer conflicts (status UU or AA): both sides recorded a different commit,
// with no base content to merge. These are resolved by picking one side's
// commit. Other submodule conflicts (e.g. modify/delete) are handled like
// ordinary non-textual conflicts, with the keep/delete picker.
func (self *FilesController) isSubmoduleCommitConflict(file *models.File) bool {
	return file != nil && file.HasInlineMergeConflicts && file.IsSubmodule(self.c.Model().Submodules)
}

func (self *FilesController) openConflictResolutionMenu(file *models.File) error {
	if self.isSubmoduleCommitConflict(file) {
		return self.openSubmoduleConflictMenu(file)
	}

	return self.openFileConflictMenu(file)
}

func (self *FilesController) openFileConflictMenu(file *models.File) error {
	handle := func(command func(command string) error, logText string) error {
		self.c.LogAction(logText)
		if err := command(file.GetPath()); err != nil {
			return err
		}
		self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}})
		return nil
	}
	keepItem := &types.MenuItem{
		Label: self.c.Tr.MergeConflictKeepFile,
		OnPress: func() error {
			return handle(self.c.Git().WorkingTree.StageFile, self.c.Tr.Actions.ResolveConflictByKeepingFile)
		},
		Keys: menuKey('k'),
	}
	deleteItem := &types.MenuItem{
		Label: self.c.Tr.MergeConflictDeleteFile,
		OnPress: func() error {
			return handle(self.c.Git().WorkingTree.RemoveConflictedFile, self.c.Tr.Actions.ResolveConflictByDeletingFile)
		},
		Keys: menuKey('d'),
	}
	items := []*types.MenuItem{}
	switch file.ShortStatus {
	case "DD":
		// For "both deleted" conflicts, deleting the file is the only reasonable thing you can do.
		// Restoring to the state before deletion is not the responsibility of a conflict resolution tool.
		items = append(items, deleteItem)
	case "DU", "UD":
		// For these, we put the delete option first because it's the most common one,
		// even if it's more destructive.
		items = append(items, deleteItem, keepItem)
	case "AU", "UA":
		// For these, we put the keep option first because it's less destructive,
		// and the chances between keep and delete are 50/50.
		items = append(items, keepItem, deleteItem)
	default:
		panic("should only be called if there's a merge conflict")
	}
	return self.c.Menu(types.CreateMenuOptions{
		Title:  self.c.Tr.MergeConflictsTitle,
		Prompt: file.GetMergeStateDescription(self.c.Tr),
		Items:  items,
	})
}

func (self *FilesController) openSubmoduleConflictMenu(file *models.File) error {
	path := file.GetPath()
	_, ours, theirs, err := self.c.Git().Submodule.GetConflictCommits(path)
	if err != nil {
		return err
	}

	resolve := func(sha string, logAction string) error {
		self.c.LogAction(logAction)
		if err := self.c.Git().Submodule.CheckoutConflictCommit(path, sha); err != nil {
			return err
		}
		if err := self.c.Git().WorkingTree.StageFile(path); err != nil {
			return err
		}
		self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}})
		return nil
	}

	// Append the commit summary to the label so the user can tell the two
	// candidates apart, falling back to the bare label if we can't read it.
	label := func(text string, sha string) string {
		if summary, err := self.c.Git().Submodule.GetCommitSummary(path, sha); err == nil && summary != "" {
			return fmt.Sprintf("%s (%s)", text, summary)
		}
		return text
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title:  self.c.Tr.MergeConflictsTitle,
		Prompt: utils.ResolvePlaceholderString(self.c.Tr.SubmoduleMergeConflictDescription, map[string]string{"path": path}),
		Items: []*types.MenuItem{
			{
				Label:   label(self.c.Tr.MergeConflictTakeCurrentCommit, ours),
				OnPress: func() error { return resolve(ours, self.c.Tr.Actions.TakeCurrentSubmoduleCommit) },
				Keys:    menuKey('c'),
			},
			{
				Label:   label(self.c.Tr.MergeConflictTakeIncomingCommit, theirs),
				OnPress: func() error { return resolve(theirs, self.c.Tr.Actions.TakeIncomingSubmoduleCommit) },
				Keys:    menuKey('i'),
			},
		},
	})
}

func (self *FilesController) toggleStagedAll() error {
	if err := self.toggleStagedAllWithLock(); err != nil {
		return err
	}

	self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}, Mode: types.ASYNC})

	self.context().HandleFocus(types.OnFocusOpts{})
	return nil
}

func (self *FilesController) toggleStagedAllWithLock() error {
	self.c.Mutexes().RefreshingFilesMutex.Lock()
	defer self.c.Mutexes().RefreshingFilesMutex.Unlock()

	root := self.context().FileTreeViewModel.GetRoot()

	stage := func(unstagedNodes []*filetree.FileNode) error {
		if self.context().IsFiltering() {
			// When filtering, only stage visible files
			var paths []string
			_ = root.ForEachFile(func(file *models.File) error {
				paths = append(paths, file.Path)
				return nil
			})
			return self.c.Git().WorkingTree.StageFiles(paths, nil)
		}

		onlyTrackedFiles := self.context().GetStatusFilter() == filetree.DisplayTracked
		return self.c.Git().WorkingTree.StageAll(onlyTrackedFiles)
	}

	unstage := func(nodes []*filetree.FileNode) error {
		if self.context().IsFiltering() {
			// When filtering, only unstage visible files
			return self.unstageFilteredFiles(nodes)
		}

		return self.c.Git().WorkingTree.UnstageAll()
	}

	return self.toggleStaged([]*filetree.FileNode{root},
		self.c.Tr.Actions.StageAllFiles, self.c.Tr.Actions.UnstageAllFiles,
		stage, unstage)
}

func (self *FilesController) unstageFiles(node *filetree.FileNode) error {
	return node.ForEachFile(func(file *models.File) error {
		if file.HasStagedChanges {
			if err := self.c.Git().WorkingTree.UnStageFile(file.Names(), file.Tracked); err != nil {
				return err
			}
		}

		return nil
	})
}

func (self *FilesController) ignoreOrExcludeTracked(node *filetree.FileNode, trAction string, f func(string) error) error {
	self.c.LogAction(trAction)
	// not 100% sure if this is necessary but I'll assume it is
	if err := self.unstageFiles(node); err != nil {
		return err
	}

	if err := self.c.Git().WorkingTree.RemoveTrackedFiles(node.GetPath()); err != nil {
		return err
	}

	if err := f(node.GetPath()); err != nil {
		return err
	}

	self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}})
	return nil
}

func (self *FilesController) ignoreOrExcludeUntracked(node *filetree.FileNode, trAction string, f func(string) error) error {
	self.c.LogAction(trAction)

	if err := f(node.GetPath()); err != nil {
		return err
	}

	self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}})
	return nil
}

func (self *FilesController) ignoreOrExcludeFile(node *filetree.FileNode, trText string, trPrompt string, trAction string, f func(string) error) error {
	if node.GetIsTracked() {
		self.c.Confirm(types.ConfirmOpts{
			Title:  trText,
			Prompt: trPrompt,
			HandleConfirm: func() error {
				return self.ignoreOrExcludeTracked(node, trAction, f)
			},
		})

		return nil
	}
	return self.ignoreOrExcludeUntracked(node, trAction, f)
}

func (self *FilesController) ignore(node *filetree.FileNode) error {
	if node.GetPath() == ".gitignore" {
		return errors.New(self.c.Tr.Actions.IgnoreFileErr)
	}
	return self.ignoreOrExcludeFile(node, self.c.Tr.IgnoreTracked, self.c.Tr.IgnoreTrackedPrompt, self.c.Tr.Actions.IgnoreExcludeFile, self.c.Git().WorkingTree.Ignore)
}

func (self *FilesController) exclude(node *filetree.FileNode) error {
	if node.GetPath() == ".gitignore" {
		return errors.New(self.c.Tr.Actions.ExcludeGitIgnoreErr)
	}

	return self.ignoreOrExcludeFile(node, self.c.Tr.ExcludeTracked, self.c.Tr.ExcludeTrackedPrompt, self.c.Tr.Actions.ExcludeFile, self.c.Git().WorkingTree.Exclude)
}

func (self *FilesController) ignoreOrExcludeMenu(node *filetree.FileNode) error {
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.Actions.IgnoreExcludeFile,
		Items: []*types.MenuItem{
			{
				LabelColumns: []string{self.c.Tr.IgnoreFile},
				OnPress: func() error {
					if err := self.ignore(node); err != nil {
						return err
					}
					return nil
				},
				Keys: menuKey('i'),
			},
			{
				LabelColumns: []string{self.c.Tr.ExcludeFile},
				OnPress: func() error {
					if err := self.exclude(node); err != nil {
						return err
					}
					return nil
				},
				Keys: menuKey('e'),
			},
		},
	})
}

func (self *FilesController) refresh() error {
	self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}})
	return nil
}

func (self *FilesController) handleAmendCommitPress() error {
	doAmend := func() error {
		return self.c.Helpers().WorkingTree.WithEnsureCommittableFiles(func() error {
			if len(self.c.Model().Commits) == 0 {
				return errors.New(self.c.Tr.NoCommitToAmend)
			}

			return self.c.Helpers().AmendHelper.AmendHead()
		})
	}

	if self.isResolvingConflicts() {
		return self.c.Menu(types.CreateMenuOptions{
			Title:      self.c.Tr.AmendCommitTitle,
			Prompt:     self.c.Tr.AmendCommitWithConflictsMenuPrompt,
			HideCancel: true, // We want the cancel item first, so we add one manually
			Items: []*types.MenuItem{
				{
					Label: self.c.Tr.Cancel,
					OnPress: func() error {
						return nil
					},
				},
				{
					Label: self.c.Tr.AmendCommitWithConflictsContinue,
					OnPress: func() error {
						return self.c.Helpers().MergeAndRebase.ContinueRebase()
					},
				},
				{
					Label: self.c.Tr.AmendCommitWithConflictsAmend,
					OnPress: func() error {
						return doAmend()
					},
				},
			},
		})
	}

	return self.c.ConfirmIf(!self.c.UserConfig().Gui.SkipAmendWarning,
		types.ConfirmOpts{
			Title:  self.c.Tr.AmendLastCommitTitle,
			Prompt: self.c.Tr.SureToAmend,
			HandleConfirm: func() error {
				return doAmend()
			},
		},
	)
}

func (self *FilesController) isResolvingConflicts() bool {
	commits := self.c.Model().Commits
	for _, c := range commits {
		if c.Status == models.StatusConflicted {
			return true
		}
		if !c.IsTODO() {
			break
		}
	}
	return false
}

func (self *FilesController) handleStatusFilterPressed() error {
	currentFilter := self.context().GetStatusFilter()
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.FilteringMenuTitle,
		Items: []*types.MenuItem{
			{
				Label: self.c.Tr.FilterStagedFiles,
				OnPress: func() error {
					return self.setStatusFiltering(filetree.DisplayStaged)
				},
				Keys:   menuKey('s'),
				Widget: types.MakeMenuRadioButton(currentFilter == filetree.DisplayStaged),
			},
			{
				Label: self.c.Tr.FilterUnstagedFiles,
				OnPress: func() error {
					return self.setStatusFiltering(filetree.DisplayUnstaged)
				},
				Keys:   menuKey('u'),
				Widget: types.MakeMenuRadioButton(currentFilter == filetree.DisplayUnstaged),
			},
			{
				Label: self.c.Tr.FilterTrackedFiles,
				OnPress: func() error {
					return self.setStatusFiltering(filetree.DisplayTracked)
				},
				Keys:   menuKey('t'),
				Widget: types.MakeMenuRadioButton(currentFilter == filetree.DisplayTracked),
			},
			{
				Label: self.c.Tr.FilterUntrackedFiles,
				OnPress: func() error {
					return self.setStatusFiltering(filetree.DisplayUntracked)
				},
				Keys:   menuKey('T'),
				Widget: types.MakeMenuRadioButton(currentFilter == filetree.DisplayUntracked),
			},
			{
				Label: self.c.Tr.NoFilter,
				OnPress: func() error {
					return self.setStatusFiltering(filetree.DisplayAll)
				},
				Keys:   menuKey('r'),
				Widget: types.MakeMenuRadioButton(currentFilter == filetree.DisplayAll),
			},
		},
	})
}

func (self *FilesController) filteringLabel(filter filetree.FileTreeDisplayFilter) string {
	switch filter {
	case filetree.DisplayAll:
		return ""
	case filetree.DisplayStaged:
		return self.c.Tr.FilterLabelStagedFiles
	case filetree.DisplayUnstaged:
		return self.c.Tr.FilterLabelUnstagedFiles
	case filetree.DisplayTracked:
		return self.c.Tr.FilterLabelTrackedFiles
	case filetree.DisplayUntracked:
		return self.c.Tr.FilterLabelUntrackedFiles
	case filetree.DisplayConflicted:
		return self.c.Tr.FilterLabelConflictingFiles
	}

	panic(fmt.Sprintf("Unexpected files display filter: %d", filter))
}

func (self *FilesController) setStatusFiltering(filter filetree.FileTreeDisplayFilter) error {
	previousFilter := self.context().GetStatusFilter()

	self.context().FileTreeViewModel.SetStatusFilter(filter)
	self.c.Contexts().Files.GetView().Subtitle = self.filteringLabel(filter)

	// Whenever we switch between untracked and other filters, we need to refresh the files view
	// because the untracked files filter applies when running `git status`.
	if previousFilter != filter && (previousFilter == filetree.DisplayUntracked || filter == filetree.DisplayUntracked) {
		self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.FILES}, Mode: types.ASYNC})
	} else {
		self.c.PostRefreshUpdate(self.context())
	}
	return nil
}

func (self *FilesController) edit(nodes []*filetree.FileNode) error {
	return self.c.Helpers().Files.EditFiles(lo.FilterMap(nodes,
		func(node *filetree.FileNode, _ int) (string, bool) {
			return node.GetPath(), node.IsFile()
		}))
}

func (self *FilesController) canEditFiles(nodes []*filetree.FileNode) *types.DisabledReason {
	if lo.NoneBy(nodes, func(node *filetree.FileNode) bool { return node.IsFile() }) {
		return &types.DisabledReason{
			Text:             self.c.Tr.ErrCannotEditDirectory,
			ShowErrorInPanel: true,
		}
	}

	return nil
}

func (self *FilesController) Open() error {
	node := self.context().GetSelected()
	if node == nil {
		return nil
	}

	return self.c.Helpers().Files.OpenFile(node.GetPath())
}

func (self *FilesController) openDiffTool(node *filetree.FileNode) error {
	fromCommit := ""
	reverse := false
	if self.c.Modes().Diffing.Active() {
		fromCommit = self.c.Modes().Diffing.Ref
		reverse = self.c.Modes().Diffing.Reverse
	}
	return self.c.RunSubprocessAndRefresh(
		self.c.Git().Diff.OpenDiffToolCmdObj(
			git_commands.DiffToolCmdOptions{
				Filepath:    node.GetPath(),
				FromCommit:  fromCommit,
				ToCommit:    "",
				Reverse:     reverse,
				IsDirectory: !node.IsFile(),
				Staged:      !node.GetHasUnstagedChanges(),
			}),
	)
}

func (self *FilesController) switchToMerge() error {
	file := self.getSelectedFile()
	if file == nil {
		return nil
	}

	return self.c.Helpers().MergeConflicts.SwitchToMerge(file.Path)
}

func (self *FilesController) createStashMenu() error {
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.StashOptions,
		Items: []*types.MenuItem{
			{
				Label: self.c.Tr.StashAllChanges,
				OnPress: func() error {
					if !self.c.Helpers().WorkingTree.IsWorkingTreeDirtyExceptSubmodules() {
						return errors.New(self.c.Tr.NoFilesToStash)
					}
					return self.handleStashSave(self.c.Git().Stash.Push, self.c.Tr.Actions.StashAllChanges)
				},
				Keys: menuKey('a'),
			},
			{
				Label: self.c.Tr.StashAllChangesKeepIndex,
				OnPress: func() error {
					if !self.c.Helpers().WorkingTree.IsWorkingTreeDirtyExceptSubmodules() {
						return errors.New(self.c.Tr.NoFilesToStash)
					}
					// if there are no staged files it behaves the same as Stash.Save
					return self.handleStashSave(self.c.Git().Stash.StashAndKeepIndex, self.c.Tr.Actions.StashAllChangesKeepIndex)
				},
				Keys: menuKey('i'),
			},
			{
				Label: self.c.Tr.StashIncludeUntrackedChanges,
				OnPress: func() error {
					return self.handleStashSave(self.c.Git().Stash.StashIncludeUntrackedChanges, self.c.Tr.Actions.StashIncludeUntrackedChanges)
				},
				Keys: menuKey('U'),
			},
			{
				Label: self.c.Tr.StashStagedChanges,
				OnPress: func() error {
					// there must be something in staging otherwise the current implementation mucks the stash up
					if !self.c.Helpers().WorkingTree.AnyStagedFilesExceptSubmodules() {
						return errors.New(self.c.Tr.NoTrackedStagedFilesStash)
					}
					return self.handleStashSave(self.c.Git().Stash.SaveStagedChanges, self.c.Tr.Actions.StashStagedChanges)
				},
				Keys: menuKey('s'),
			},
			{
				Label: self.c.Tr.StashUnstagedChanges,
				OnPress: func() error {
					if !self.c.Helpers().WorkingTree.IsWorkingTreeDirtyExceptSubmodules() {
						return errors.New(self.c.Tr.NoFilesToStash)
					}
					if self.c.Helpers().WorkingTree.AnyStagedFilesExceptSubmodules() {
						return self.handleStashSave(self.c.Git().Stash.StashUnstagedChanges, self.c.Tr.Actions.StashUnstagedChanges)
					}
					// ordinary stash
					return self.handleStashSave(self.c.Git().Stash.Push, self.c.Tr.Actions.StashUnstagedChanges)
				},
				Keys: menuKey('u'),
			},
		},
	})
}

func (self *FilesController) openMergeConflictMenu(nodes []*filetree.FileNode) error {
	normalizedNodes := flattenSelectedNodesToFiles(nodes)

	fileNodesWithConflicts := lo.Filter(normalizedNodes, func(node *filetree.FileNode, _ int) bool {
		return node.File != nil && node.File.HasInlineMergeConflicts
	})

	filepaths := lo.Map(fileNodesWithConflicts, func(node *filetree.FileNode, _ int) string {
		return node.GetPath()
	})

	return self.c.Helpers().WorkingTree.CreateMergeConflictMenu(filepaths)
}

func (self *FilesController) canOpenMergeConflictMenu(nodes []*filetree.FileNode) *types.DisabledReason {
	normalizedNodes := flattenSelectedNodesToFiles(nodes)

	hasFileNodesWithConflicts := lo.SomeBy(normalizedNodes, func(node *filetree.FileNode) bool {
		return node.File != nil && node.File.HasInlineMergeConflicts
	})

	if !hasFileNodesWithConflicts {
		return &types.DisabledReason{Text: self.c.Tr.NoFilesWithMergeConflicts}
	}

	return nil
}

func (self *FilesController) openCopyMenu() error {
	node := self.context().GetSelected()

	copyNameItem := &types.MenuItem{
		Label: self.c.Tr.CopyFileName,
		OnPress: func() error {
			if err := self.c.OS().CopyToClipboard(node.Name()); err != nil {
				return err
			}
			self.c.Toast(self.c.Tr.FileNameCopiedToast)
			return nil
		},
		DisabledReason: self.require(self.singleItemSelected())(),
		Keys:           menuKey('n'),
	}
	copyRelativePathItem := &types.MenuItem{
		Label: self.c.Tr.CopyRelativeFilePath,
		OnPress: func() error {
			if err := self.c.OS().CopyToClipboard(node.GetPath()); err != nil {
				return err
			}
			self.c.Toast(self.c.Tr.FilePathCopiedToast)
			return nil
		},
		DisabledReason: self.require(self.singleItemSelected())(),
		Keys:           menuKey('p'),
	}
	copyAbsolutePathItem := &types.MenuItem{
		Label: self.c.Tr.CopyAbsoluteFilePath,
		OnPress: func() error {
			absPath, err := filepath.Abs(node.GetPath())
			if err != nil {
				return err
			}
			if err := self.c.OS().CopyToClipboard(absPath); err != nil {
				return err
			}
			self.c.Toast(self.c.Tr.FilePathCopiedToast)
			return nil
		},
		DisabledReason: self.require(self.singleItemSelected())(),
		Keys:           menuKey('P'),
	}
	copyFileDiffItem := &types.MenuItem{
		Label:   self.c.Tr.CopySelectedDiff,
		Tooltip: self.c.Tr.CopyFileDiffTooltip,
		OnPress: func() error {
			path := self.context().GetSelectedPath()
			hasStaged := self.hasPathStagedChanges(node)
			diff, err := self.c.Git().Diff.GetDiff(hasStaged, "--", path)
			if err != nil {
				return err
			}
			if err := self.c.OS().CopyToClipboard(diff); err != nil {
				return err
			}
			self.c.Toast(self.c.Tr.FileDiffCopiedToast)
			return nil
		},
		DisabledReason: self.require(self.singleItemSelected(
			func(file *filetree.FileNode) *types.DisabledReason {
				if !node.GetHasStagedOrTrackedChanges() {
					return &types.DisabledReason{Text: self.c.Tr.NoContentToCopyError}
				}
				return nil
			},
		))(),
		Keys: menuKey('s'),
	}
	copyAllDiff := &types.MenuItem{
		Label:   self.c.Tr.CopyAllFilesDiff,
		Tooltip: self.c.Tr.CopyFileDiffTooltip,
		OnPress: func() error {
			hasStaged := self.c.Helpers().WorkingTree.AnyStagedFiles()
			diff, err := self.c.Git().Diff.GetDiff(hasStaged, "--")
			if err != nil {
				return err
			}
			if err := self.c.OS().CopyToClipboard(diff); err != nil {
				return err
			}
			self.c.Toast(self.c.Tr.AllFilesDiffCopiedToast)
			return nil
		},
		DisabledReason: self.require(
			func() *types.DisabledReason {
				if !self.anyStagedOrTrackedFile() {
					return &types.DisabledReason{Text: self.c.Tr.NoContentToCopyError}
				}
				return nil
			},
		)(),
		Keys: menuKey('a'),
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.CopyToClipboardMenu,
		Items: []*types.MenuItem{
			copyNameItem,
			copyRelativePathItem,
			copyAbsolutePathItem,
			copyFileDiffItem,
			copyAllDiff,
		},
	})
}

func (self *FilesController) anyStagedOrTrackedFile() bool {
	if !self.c.Helpers().WorkingTree.AnyStagedFiles() {
		return self.c.Helpers().WorkingTree.AnyTrackedFiles()
	}
	return true
}

func (self *FilesController) hasPathStagedChanges(node *filetree.FileNode) bool {
	return node.SomeFile(func(t *models.File) bool {
		return t.HasStagedChanges
	})
}

func (self *FilesController) stash() error {
	return self.handleStashSave(self.c.Git().Stash.Push, self.c.Tr.Actions.StashAllChanges)
}

func (self *FilesController) createResetToUpstreamMenu() error {
	return self.c.Helpers().Refs.CreateGitResetMenu("@{upstream}", "@{upstream}")
}

func (self *FilesController) handleToggleDirCollapsed() error {
	node := self.context().GetSelected()
	if node == nil {
		return nil
	}

	self.context().FileTreeViewModel.ToggleCollapsed(node.GetInternalPath())

	self.c.PostRefreshUpdate(self.c.Contexts().Files)

	return nil
}

func (self *FilesController) toggleTreeView() error {
	self.context().FileTreeViewModel.ToggleShowTree()

	self.c.PostRefreshUpdate(self.context())
	return nil
}

func (self *FilesController) handleStashSave(stashFunc func(message string) error, action string) error {
	self.c.Prompt(types.PromptOpts{
		Title: self.c.Tr.StashChanges,
		HandleConfirm: func(stashComment string) error {
			self.c.LogAction(action)

			if err := stashFunc(stashComment); err != nil {
				return err
			}
			self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.STASH, types.FILES}})
			return nil
		},
		AllowEmptyInput: true,
	})

	return nil
}

func (self *FilesController) onClickMain(opts gocui.ViewMouseBindingOpts) error {
	return self.EnterFile(nil, types.OnFocusOpts{ClickedWindowName: "main", ClickedViewLineIdx: opts.Y})
}

func (self *FilesController) fetch() error {
	return self.c.WithWaitingStatus(self.c.Tr.FetchingStatus, func(task gocui.Task) error {
		self.c.LogAction("Fetch")
		err := self.c.Git().Sync.Fetch(task)

		if err != nil && strings.Contains(err.Error(), "exit status 128") {
			return errors.New(self.c.Tr.PassUnameWrong)
		}

		return self.c.Helpers().BranchesHelper.PostFetchRefresh(err, false)
	})
}

// Couldn't think of a better term than 'normalised'. Alas.
// The idea is that when you select a range of nodes, you will often have both
// a node and its parent node selected. If we are trying to discard changes to the
// selected nodes, we'll get an error if we try to discard the child after the parent.
// So we just need to filter out any nodes from the selection that are descendants
// of other nodes
func normalisedSelectedNodes(selectedNodes []*filetree.FileNode) []*filetree.FileNode {
	return lo.Filter(selectedNodes, func(node *filetree.FileNode, _ int) bool {
		return !isDescendentOfSelectedNodes(node, selectedNodes)
	})
}

func isDescendentOfSelectedNodes(node *filetree.FileNode, selectedNodes []*filetree.FileNode) bool {
	nodePath := node.GetInternalPath()

	for _, selectedNode := range selectedNodes {
		if selectedNode.IsFile() {
			continue
		}

		selectedNodePath := selectedNode.GetInternalPath()

		if strings.HasPrefix(nodePath, selectedNodePath+"/") {
			return true
		}
	}
	return false
}

// BFS algorithm for expanding directories into their children,
// and for collecting the unique file nodes
func flattenSelectedNodesToFiles(selectedNodes []*filetree.FileNode) []*filetree.FileNode {
	queue := append(make([]*filetree.FileNode, 0, len(selectedNodes)), selectedNodes...)
	visited := set.New[string]()
	var files []*filetree.FileNode

	for len(queue) > 0 {
		// pop node from queue
		node := queue[0]
		queue = queue[1:]

		nodeID := node.ID()
		if visited.Includes(nodeID) {
			continue
		}
		visited.Add(nodeID)

		if node.File != nil {
			// unique file node -> collect it
			files = append(files, node)
			continue
		}

		// directory node -> enqueue children
		for _, ch := range node.Children {
			queue = append(queue, &filetree.FileNode{Node: ch})
		}
	}
	return files
}

func someNodesHaveUnstagedChanges(nodes []*filetree.FileNode) bool {
	return lo.SomeBy(nodes, (*filetree.FileNode).GetHasUnstagedChanges)
}

func someNodesHaveStagedChanges(nodes []*filetree.FileNode) bool {
	return lo.SomeBy(nodes, (*filetree.FileNode).GetHasStagedChanges)
}

func filterNodesHaveUnstagedChanges(nodes []*filetree.FileNode, submodules []*models.SubmoduleConfig) []*filetree.FileNode {
	return lo.Filter(nodes, func(node *filetree.FileNode, _ int) bool {
		return node.SomeFile(func(file *models.File) bool {
			return fileHasStageableUnstagedChanges(file, submodules)
		})
	})
}

// For a submodule, the only thing the parent repo can stage is the
// commit-pointer change; dirty or untracked content within the submodule
// shows up as an unstaged change but can never be staged from the parent. So
// once the submodule's commit is staged (leaving it at e.g. "MM"), we mustn't
// treat the leftover unstaged change as stageable, or pressing space would
// keep trying to stage it instead of unstaging it.
func fileHasStageableUnstagedChanges(file *models.File, submodules []*models.SubmoduleConfig) bool {
	if !file.HasUnstagedChanges {
		return false
	}

	if file.IsSubmodule(submodules) {
		return !file.HasStagedChanges
	}

	return true
}

// stagingWouldBeNoOp reports whether staging the given nodes would have no
// visible effect, which happens when the only things being staged are
// submodules that have dirty or untracked content but no new commit: the
// parent repo can't stage that content. If a regular file (or a submodule with
// a stageable new commit) is among them, staging does something, so this
// returns false.
func (self *FilesController) stagingWouldBeNoOp(nodes []*filetree.FileNode) (bool, error) {
	submodules := self.c.Model().Submodules

	var submodulePaths []string
	hasOtherStageableChanges := false
	for _, node := range nodes {
		_ = node.ForEachFile(func(file *models.File) error {
			if file.IsSubmodule(submodules) {
				submodulePaths = append(submodulePaths, file.Path)
			} else if file.HasUnstagedChanges {
				hasOtherStageableChanges = true
			}
			return nil
		})
	}

	if hasOtherStageableChanges || len(submodulePaths) == 0 {
		return false, nil
	}

	anyStageable, err := self.c.Git().Submodule.AnyHaveStageableChanges(submodulePaths)
	if err != nil {
		return false, err
	}

	return !anyStageable, nil
}

func findSubmoduleNode(nodes []*filetree.FileNode, submodules []*models.SubmoduleConfig) *models.File {
	for _, node := range nodes {
		submoduleNode := node.FindFirstFileBy(func(f *models.File) bool {
			return f.IsSubmodule(submodules)
		})
		if submoduleNode != nil {
			return submoduleNode
		}
	}
	return nil
}

func (self *FilesController) canRemove(selectedNodes []*filetree.FileNode) *types.DisabledReason {
	// Return disabled if the selection contains multiple changed items and includes a submodule change.
	submodules := self.c.Model().Submodules
	hasFiles := false
	uniqueSelectedSubmodules := set.New[*models.SubmoduleConfig]()

	for _, node := range selectedNodes {
		_ = node.ForEachFile(func(f *models.File) error {
			if submodule := f.SubmoduleConfig(submodules); submodule != nil {
				uniqueSelectedSubmodules.Add(submodule)
			} else {
				hasFiles = true
			}
			return nil
		})
		if uniqueSelectedSubmodules.Len() > 0 && (hasFiles || uniqueSelectedSubmodules.Len() > 1) {
			return &types.DisabledReason{Text: self.c.Tr.MultiSelectNotSupportedForSubmodules}
		}
	}

	return nil
}

func (self *FilesController) remove(selectedNodes []*filetree.FileNode) error {
	submodules := self.c.Model().Submodules

	selectedNodes = normalisedSelectedNodes(selectedNodes)

	// If we have one submodule then we must only have one submodule or `canRemove` would have
	// returned an error
	submoduleNode := findSubmoduleNode(selectedNodes, submodules)
	if submoduleNode != nil {
		submodule := submoduleNode.SubmoduleConfig(submodules)

		menuItems := []*types.MenuItem{
			{
				Label: self.c.Tr.SubmoduleStashAndReset,
				OnPress: func() error {
					return self.ResetSubmodule(submodule)
				},
			},
		}

		return self.c.Menu(types.CreateMenuOptions{Title: submoduleNode.GetPath(), Items: menuItems})
	}

	discardAllChangesItem := types.MenuItem{
		Label: self.c.Tr.DiscardAllChanges,
		OnPress: func() error {
			self.c.LogAction(self.c.Tr.Actions.DiscardAllChangesInFile)

			if self.context().IsSelectingRange() {
				defer self.context().CancelRangeSelect()
			}

			nodes := lo.Map(selectedNodes, func(n *filetree.FileNode, _ int) git_commands.IFileNode { return n })
			if err := self.c.Git().WorkingTree.DiscardAllDirChanges(nodes); err != nil {
				return err
			}

			self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES, types.WORKTREES}})
			return nil
		},
		Keys: self.c.KeybindingsOpts().GetKeys(self.c.UserConfig().Keybinding.Files.ConfirmDiscard),
		Tooltip: utils.ResolvePlaceholderString(
			self.c.Tr.DiscardAllTooltip,
			map[string]string{
				"path": self.formattedPaths(selectedNodes),
			},
		),
	}

	discardUnstagedChangesItem := types.MenuItem{
		Label: self.c.Tr.DiscardUnstagedChanges,
		OnPress: func() error {
			self.c.LogAction(self.c.Tr.Actions.DiscardAllUnstagedChangesInFile)

			if self.context().IsSelectingRange() {
				defer self.context().CancelRangeSelect()
			}

			nodes := lo.Map(selectedNodes, func(n *filetree.FileNode, _ int) git_commands.IFileNode { return n })
			if err := self.c.Git().WorkingTree.DiscardUnstagedDirChanges(nodes); err != nil {
				return err
			}

			self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES, types.WORKTREES}})
			return nil
		},
		Keys: menuKey('u'),
		Tooltip: utils.ResolvePlaceholderString(
			self.c.Tr.DiscardUnstagedTooltip,
			map[string]string{
				"path": self.formattedPaths(selectedNodes),
			},
		),
	}

	if !someNodesHaveStagedChanges(selectedNodes) || !someNodesHaveUnstagedChanges(selectedNodes) {
		discardUnstagedChangesItem.DisabledReason = &types.DisabledReason{Text: self.c.Tr.DiscardUnstagedDisabled}
	}

	menuItems := []*types.MenuItem{
		&discardAllChangesItem,
		&discardUnstagedChangesItem,
	}

	return self.c.Menu(types.CreateMenuOptions{Title: self.c.Tr.DiscardChangesTitle, Items: menuItems})
}

func (self *FilesController) ResetSubmodule(submodule *models.SubmoduleConfig) error {
	return self.c.WithWaitingStatus(self.c.Tr.ResettingSubmoduleStatus, func(gocui.Task) error {
		self.c.LogAction(self.c.Tr.Actions.ResetSubmodule)

		file := self.c.Helpers().WorkingTree.FileForSubmodule(submodule)
		if file != nil {
			if err := self.c.Git().WorkingTree.UnStageFile(file.Names(), file.Tracked); err != nil {
				return err
			}
		}

		if err := self.c.Git().Submodule.Stash(submodule); err != nil {
			return err
		}
		if err := self.c.Git().Submodule.Reset(submodule); err != nil {
			return err
		}

		self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC, Scope: []types.RefreshableView{types.FILES, types.SUBMODULES}})
		return nil
	})
}

func (self *FilesController) formattedPaths(nodes []*filetree.FileNode) string {
	return utils.FormatPaths(lo.Map(nodes, func(node *filetree.FileNode, _ int) string {
		return node.GetPath()
	}))
}

func (self *FilesController) isInTreeMode() *types.DisabledReason {
	if !self.context().FileTreeViewModel.InTreeMode() {
		return &types.DisabledReason{Text: self.c.Tr.DisabledInFlatView}
	}

	return nil
}
