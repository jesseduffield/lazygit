package controllers

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/constants"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type CommitFilesController struct {
	baseController
	*ListControllerTrait[*filetree.CommitFileNode]
	c *ControllerCommon
}

var _ types.IController = &CommitFilesController{}

func NewCommitFilesController(
	c *ControllerCommon,
) *CommitFilesController {
	return &CommitFilesController{
		baseController: baseController{},
		c:              c,
		ListControllerTrait: NewListControllerTrait(
			c,
			c.Contexts().CommitFiles,
			c.Contexts().CommitFiles.GetSelected,
			c.Contexts().CommitFiles.GetSelectedItems,
		),
	}
}

func (self *CommitFilesController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Keys:        opts.GetKeys(opts.Config.Files.CopyFileInfoToClipboard),
			Handler:     self.openCopyMenu,
			Description: self.c.Tr.CopyToClipboardMenu,
			OpensMenu:   true,
		},
		{
			Keys:              opts.GetKeys(opts.Config.CommitFiles.CheckoutCommitFile),
			Handler:           self.withItem(self.checkout),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Checkout,
			Tooltip:           self.c.Tr.CheckoutCommitFileTooltip,
			DisplayOnScreen:   true,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Universal.Remove),
			Handler:           self.withItems(self.discard),
			GetDisabledReason: self.require(self.itemsSelected(self.canDiscardFileChanges)),
			Description:       self.c.Tr.Discard,
			Tooltip:           self.c.Tr.DiscardOldFileChangeTooltip,
			DisplayOnScreen:   true,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Universal.OpenFile),
			Handler:           self.withItem(self.open),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.OpenFile,
			Tooltip:           self.c.Tr.OpenFileTooltip,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Universal.Edit),
			Handler:           self.withItems(self.edit),
			GetDisabledReason: self.require(self.itemsSelected(self.canEditFiles)),
			Description:       self.c.Tr.Edit,
			Tooltip:           self.c.Tr.EditFileTooltip,
			DisplayOnScreen:   true,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Universal.OpenDiffTool),
			Handler:           self.withItem(self.openDiffTool),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.OpenDiffTool,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Universal.Select),
			Handler:           self.withItems(self.toggleForPatch),
			GetDisabledReason: self.require(self.itemsSelected()),
			Description:       self.c.Tr.ToggleAddToPatch,
			Tooltip: utils.ResolvePlaceholderString(self.c.Tr.ToggleAddToPatchTooltip,
				map[string]string{"doc": constants.Links.Docs.CustomPatchDemo},
			),
			DisplayOnScreen: true,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Files.ToggleStagedAll),
			Handler:     self.withItem(self.toggleAllForPatch),
			Description: self.c.Tr.ToggleAllInPatch,
			Tooltip: utils.ResolvePlaceholderString(self.c.Tr.ToggleAllInPatchTooltip,
				map[string]string{"doc": constants.Links.Docs.CustomPatchDemo},
			),
		},
		{
			Keys:              opts.GetKeys(opts.Config.Universal.GoInto),
			Handler:           self.withItem(self.enter),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.EnterCommitFile,
			Tooltip:           self.c.Tr.EnterCommitFileTooltip,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Files.ToggleTreeView),
			Handler:     self.toggleTreeView,
			Description: self.c.Tr.ToggleTreeView,
			Tooltip:     self.c.Tr.ToggleTreeViewTooltip,
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

	return bindings
}

func (self *CommitFilesController) context() *context.CommitFilesContext {
	return self.c.Contexts().CommitFiles
}

func (self *CommitFilesController) GetOnClick() func(opts gocui.ViewMouseBindingOpts) error {
	return func(opts gocui.ViewMouseBindingOpts) error {
		clickedIdx := self.context().GetSelectedLineIdx()
		node := self.context().CommitFileTreeViewModel.Get(clickedIdx)
		if node == nil || node.File != nil {
			return nil
		}

		// The arrow is at column visualDepth*2 (after indentation of 2 spaces per level).
		// Only treat clicks on the arrow and the trailing space as arrow clicks.
		visualDepth := self.context().CommitFileTreeViewModel.GetVisualDepth(clickedIdx)
		arrowStartCol := visualDepth * 2
		arrowEndCol := arrowStartCol + 1
		if opts.X < arrowStartCol || opts.X > arrowEndCol {
			return nil
		}

		self.context().CommitFileTreeViewModel.ToggleCollapsed(node.GetInternalPath())
		self.c.PostRefreshUpdate(self.context())

		return nil
	}
}

func (self *CommitFilesController) GetOnRenderToMain() func() {
	return func() {
		node := self.context().GetSelected()
		if node == nil {
			return
		}

		from, to := self.context().GetFromAndToForDiff()
		from, reverse := self.c.Modes().Diffing.GetFromAndReverseArgsForDiff(from)

		paths := self.pathsForDiff(node)
		cmdObj := self.c.Git().WorkingTree.ShowFileDiffCmdObj(from, to, reverse, paths, false)
		task := types.NewRunPtyTask(cmdObj.GetCmd())

		// Keep the inclusion gutter in step with the content as this diff (re-)renders.
		// It's a no-op unless the main view is focused and a patch is being built (see
		// RefreshInclusionGutter); a patch toggle re-renders this same diff, so the marks
		// recomputed over the current content stay valid through the swap.
		self.c.Helpers().Staging.RefreshInclusionGutter()

		self.c.RenderToMainViews(types.RefreshMainOpts{
			Pair: self.c.MainViewPairs().Normal,
			Main: &types.ViewUpdateOpts{
				Title:    self.c.Tr.Patch,
				SubTitle: self.c.Helpers().Diff.IgnoringWhitespaceSubTitle(),
				Task:     task,
			},
			Secondary: secondaryPatchPanelUpdateOpts(self.c),
		})
	}
}

func (self *CommitFilesController) copyDiffToClipboard(path string, toastMessage string) error {
	from, to := self.context().GetFromAndToForDiff()
	from, reverse := self.c.Modes().Diffing.GetFromAndReverseArgsForDiff(from)

	cmdObj := self.c.Git().WorkingTree.ShowFileDiffCmdObj(from, to, reverse, []string{path}, true)
	diff, err := cmdObj.RunWithOutput()
	if err != nil {
		return err
	}
	if err := self.c.OS().CopyToClipboard(diff); err != nil {
		return err
	}
	self.c.Toast(toastMessage)
	return nil
}

func (self *CommitFilesController) copyFileContentToClipboard(path string) error {
	_, to := self.context().GetFromAndToForDiff()
	cmdObj := self.c.Git().Commit.ShowFileContentCmdObj(to, path)
	diff, err := cmdObj.RunWithOutput()
	if err != nil {
		return err
	}
	return self.c.OS().CopyToClipboard(diff)
}

func (self *CommitFilesController) openCopyMenu() error {
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
		Label: self.c.Tr.CopySelectedDiff,
		OnPress: func() error {
			return self.copyDiffToClipboard(node.GetPath(), self.c.Tr.FileDiffCopiedToast)
		},
		DisabledReason: self.require(self.singleItemSelected())(),
		Keys:           menuKey('s'),
	}
	copyAllDiff := &types.MenuItem{
		Label: self.c.Tr.CopyAllFilesDiff,
		OnPress: func() error {
			return self.copyDiffToClipboard(".", self.c.Tr.AllFilesDiffCopiedToast)
		},
		DisabledReason: self.require(self.itemsSelected())(),
		Keys:           menuKey('a'),
	}
	copyFileContentItem := &types.MenuItem{
		Label: self.c.Tr.CopyFileContent,
		OnPress: func() error {
			if err := self.copyFileContentToClipboard(node.GetPath()); err != nil {
				return err
			}
			self.c.Toast(self.c.Tr.FileContentCopiedToast)
			return nil
		},
		DisabledReason: self.require(self.singleItemSelected(
			func(node *filetree.CommitFileNode) *types.DisabledReason {
				if !node.IsFile() {
					return &types.DisabledReason{
						Text:             self.c.Tr.ErrCannotCopyContentOfDirectory,
						ShowErrorInPanel: true,
					}
				}
				return nil
			}))(),
		Keys: menuKey('c'),
	}

	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.CopyToClipboardMenu,
		Items: []*types.MenuItem{
			copyNameItem,
			copyRelativePathItem,
			copyAbsolutePathItem,
			copyFileDiffItem,
			copyAllDiff,
			copyFileContentItem,
		},
	})
}

func (self *CommitFilesController) checkout(node *filetree.CommitFileNode) error {
	hasModifiedFiles := helpers.AnyTrackedFilesInPathExceptSubmodules(node.GetPath(),
		self.c.Model().Files, self.c.Model().Submodules)
	if hasModifiedFiles {
		return errors.New(self.c.Tr.CannotCheckoutWithModifiedFilesErr)
	}

	self.c.LogAction(self.c.Tr.Actions.CheckoutFile)
	_, to := self.context().GetFromAndToForDiff()
	if err := self.c.Git().WorkingTree.CheckoutFile(to, node.GetPath()); err != nil {
		return err
	}

	self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
	return nil
}

func (self *CommitFilesController) discard(selectedNodes []*filetree.CommitFileNode) error {
	prompt := lo.Ternary(self.c.Git().Patch.PatchBuilder.Active(),
		self.c.Tr.DiscardFileChangesPromptResetPatch,
		self.c.Tr.DiscardFileChangesPrompt)

	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.DiscardFileChangesTitle,
		Prompt: prompt,
		HandleConfirm: func() error {
			return self.c.WithWaitingStatus(self.c.Tr.RebasingStatus, func(gocui.Task) error {
				var filePaths []string
				selectedNodes = normalisedSelectedCommitFileNodes(selectedNodes)

				// Reset the current patch if there is one.
				if self.c.Git().Patch.PatchBuilder.Active() {
					self.c.Git().Patch.PatchBuilder.Reset()
				}

				for _, node := range selectedNodes {
					_ = node.ForEachFile(func(file *models.CommitFile) error {
						filePaths = append(filePaths, file.GetPath())
						return nil
					})
				}

				err := self.c.Git().Rebase.DiscardOldFileChanges(self.c.Model().Commits, self.c.Contexts().LocalCommits.GetSelectedLineIdx(), filePaths)
				if err := self.c.Helpers().MergeAndRebase.CheckMergeOrRebase(err); err != nil {
					return err
				}

				if self.context().RangeSelectEnabled() {
					self.context().GetList().CancelRangeSelect()
				}

				return nil
			})
		},
	})

	return nil
}

func (self *CommitFilesController) canDiscardFileChanges(nodes []*filetree.CommitFileNode) *types.DisabledReason {
	parentContext := self.c.Context().Current().GetParentContext()
	if parentContext == nil || parentContext.GetKey() != context.LOCAL_COMMITS_CONTEXT_KEY {
		return &types.DisabledReason{
			Text:             self.c.Tr.CanOnlyDiscardFromLocalCommits,
			ShowErrorInPanel: true,
		}
	}

	if self.c.Contexts().LocalCommits.AreMultipleItemsSelected() {
		return &types.DisabledReason{
			Text:             self.c.Tr.CannotDiscardFromMultipleCommits,
			ShowErrorInPanel: true,
		}
	}

	if self.c.Git().Status.WorkingTreeState().Any() {
		return &types.DisabledReason{
			Text:             self.c.Tr.CantPatchWhileRebasingError,
			ShowErrorInPanel: true,
		}
	}

	return nil
}

func (self *CommitFilesController) open(node *filetree.CommitFileNode) error {
	return self.c.Helpers().Files.OpenFile(node.GetPath())
}

func (self *CommitFilesController) edit(nodes []*filetree.CommitFileNode) error {
	return self.c.Helpers().Files.EditFiles(lo.FilterMap(nodes,
		func(node *filetree.CommitFileNode, _ int) (string, bool) {
			return node.GetPath(), node.IsFile()
		}))
}

func (self *CommitFilesController) canEditFiles(nodes []*filetree.CommitFileNode) *types.DisabledReason {
	if lo.NoneBy(nodes, func(node *filetree.CommitFileNode) bool { return node.IsFile() }) {
		return &types.DisabledReason{
			Text:             self.c.Tr.ErrCannotEditDirectory,
			ShowErrorInPanel: true,
		}
	}

	return nil
}

func (self *CommitFilesController) openDiffTool(node *filetree.CommitFileNode) error {
	from, to := self.context().GetFromAndToForDiff()
	from, reverse := self.c.Modes().Diffing.GetFromAndReverseArgsForDiff(from)
	_, err := self.c.RunSubprocess(self.c.Git().Diff.OpenDiffToolCmdObj(
		git_commands.DiffToolCmdOptions{
			Filepath:    node.GetPath(),
			FromCommit:  from,
			ToCommit:    to,
			Reverse:     reverse,
			IsDirectory: !node.IsFile(),
			Staged:      false,
		}))
	return err
}

func (self *CommitFilesController) toggleForPatch(selectedNodes []*filetree.CommitFileNode) error {
	if self.c.UserConfig().Git.DiffContextSize == 0 {
		return fmt.Errorf(self.c.Tr.Actions.NotEnoughContextForCustomPatch,
			self.c.UserConfig().Keybinding.Universal.IncreaseContextInDiffView)
	}

	toggle := func() error {
		return self.c.WithWaitingStatus(self.c.Tr.UpdatingPatch, func(gocui.Task) error {
			if !self.c.Git().Patch.PatchBuilder.Active() {
				if err := self.c.Helpers().CommitFiles.StartPatchBuilder(); err != nil {
					return err
				}
			}

			selectedNodes = normalisedSelectedCommitFileNodes(selectedNodes)

			// Find if any file in the selection is unselected or partially added
			adding := lo.SomeBy(selectedNodes, func(node *filetree.CommitFileNode) bool {
				return node.SomeFile(func(file *models.CommitFile) bool {
					fileStatus := self.c.Git().Patch.PatchBuilder.GetFileStatus(file.Path, self.context().GetRef().RefName())
					return fileStatus == patch.PART || fileStatus == patch.UNSELECTED
				})
			})

			patchOperationFunction := self.c.Git().Patch.PatchBuilder.RemoveFile

			if adding {
				patchOperationFunction = self.c.Git().Patch.PatchBuilder.AddFileWhole
			}

			for _, node := range selectedNodes {
				err := node.ForEachFile(func(file *models.CommitFile) error {
					return patchOperationFunction(file.Path)
				})
				if err != nil {
					return err
				}
			}

			if self.c.Git().Patch.PatchBuilder.IsEmpty() {
				self.c.Git().Patch.PatchBuilder.Reset()
			}

			self.c.OnUIThread(func() error {
				self.c.PostRefreshUpdate(self.context())
				return nil
			})

			return nil
		})
	}

	from, to, reverse := self.c.Helpers().CommitFiles.CurrentFromToReverseForPatchBuilding()
	mustDiscardPatch := self.c.Git().Patch.PatchBuilder.Active() && self.c.Git().Patch.PatchBuilder.NewPatchRequired(from, to, reverse)
	return self.c.ConfirmIf(mustDiscardPatch, types.ConfirmOpts{
		Title:  self.c.Tr.DiscardPatch,
		Prompt: self.c.Tr.DiscardPatchConfirm,
		HandleConfirm: func() error {
			if mustDiscardPatch {
				self.c.Git().Patch.PatchBuilder.Reset()
			}

			return toggle()
		},
	})
}

func (self *CommitFilesController) toggleAllForPatch(_ *filetree.CommitFileNode) error {
	root := self.context().CommitFileTreeViewModel.GetRoot()
	return self.toggleForPatch([]*filetree.CommitFileNode{root})
}

func (self *CommitFilesController) enter(node *filetree.CommitFileNode) error {
	return self.c.Helpers().CommitFiles.EnterCommitFile(node, nil, types.OnFocusOpts{ClickedWindowName: "", ClickedViewLineIdx: -1, ClickedViewRealLineIdx: -1})
}

// NOTE: this is very similar to handleToggleFileTreeView, could be DRY'd with generics
func (self *CommitFilesController) toggleTreeView() error {
	self.context().CommitFileTreeViewModel.ToggleShowTree()

	self.c.PostRefreshUpdate(self.context())
	return nil
}

func (self *CommitFilesController) collapseAll() error {
	self.context().CommitFileTreeViewModel.CollapseAll()

	self.c.PostRefreshUpdate(self.context())

	return nil
}

func (self *CommitFilesController) expandAll() error {
	self.context().CommitFileTreeViewModel.ExpandAll()

	self.c.PostRefreshUpdate(self.context())

	return nil
}

func (self *CommitFilesController) GetFocusedMainViewActions() types.FocusedMainViewActions {
	return self
}

func (self *CommitFilesController) OnClick(mainViewName string, clickedLineIdx int) error {
	// Capture before any mutation below that might re-render the main view.
	snapshot := focusedMainViewSnapshot(self.c, mainViewName, self.context())

	info, ok := self.c.Helpers().Staging.GetDiffLineInfo(mainViewName, clickedLineIdx)
	line := -1
	isDeletion := false
	if ok {
		line, isDeletion = info.PatchSelectLine()
	}

	node := self.getSelectedItem()
	if node == nil {
		return nil
	}

	if !node.IsFile() && ok {
		relativePath, err := filepath.Rel(self.c.Git().RepoPaths.WorktreePath(), info.Path)
		if err != nil {
			return err
		}
		relativePath = "./" + relativePath
		self.context().CommitFileTreeViewModel.ExpandToPath(relativePath)
		self.c.PostRefreshUpdate(self.context())

		idx, ok := self.context().CommitFileTreeViewModel.GetIndexForPath(relativePath)
		if ok {
			self.context().SetSelectedLineIdx(idx)
			self.context().GetViewTrait().FocusPoint(
				self.context().ModelIndexToViewIndex(idx), false)
			node = self.context().GetSelected()
		}
	}

	// Entered from the focused main view, so escaping returns there.
	return self.c.Helpers().CommitFiles.EnterCommitFile(node, snapshot, types.OnFocusOpts{ClickedWindowName: "main", ClickedViewLineIdx: line, ClickedViewRealLineIdx: line, ClickedViewRealLineIsDeletion: isDeletion, SelectLineInDefaultMode: true})
}

// PrimaryAction toggles the selected diff line(s) into or out of the custom patch when
// space is pressed in the focused main view of a commit's files. The per-file diff's
// patch target comes from the commit files context. It refreshes normally afterwards so
// the file's patch-status indicator in the browser updates along with the secondary patch
// view (the commits / sub-commits / stash panels, which build from the whole-commit diff,
// have no such indicator and refresh more cheaply).
func (self *CommitFilesController) PrimaryAction(mainViewName string, firstLineIdx int, lastLineIdx int) error {
	from, to, reverse := self.c.Helpers().CommitFiles.CurrentFromToReverseForPatchBuilding()
	canRebase := self.context().GetCanRebase()
	return togglePatchFromFocusedMainView(self.c, mainViewName, firstLineIdx, lastLineIdx,
		from, to, reverse, canRebase,
		func() {
			self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.COMMIT_FILES}})
		})
}

// pathsForDiff returns the file paths to use for a diff command. When a text
// filter is active and the node is a directory, only the visible (filtered)
// file paths are returned so the diff reflects what the user sees.
func (self *CommitFilesController) pathsForDiff(node *filetree.CommitFileNode) []string {
	if !node.IsFile() && self.context().IsFiltering() {
		var paths []string
		_ = node.ForEachFile(func(file *models.CommitFile) error {
			paths = append(paths, file.Path)
			return nil
		})
		return paths
	}
	return []string{node.GetPath()}
}

// NOTE: these functions are identical to those in files_controller.go (except for types) and
// could also be cleaned up with some generics
func normalisedSelectedCommitFileNodes(selectedNodes []*filetree.CommitFileNode) []*filetree.CommitFileNode {
	return lo.Filter(selectedNodes, func(node *filetree.CommitFileNode, _ int) bool {
		return !isDescendentOfSelectedCommitFileNodes(node, selectedNodes)
	})
}

func isDescendentOfSelectedCommitFileNodes(node *filetree.CommitFileNode, selectedNodes []*filetree.CommitFileNode) bool {
	for _, selectedNode := range selectedNodes {
		selectedNodePath := selectedNode.GetPath()
		nodePath := node.GetPath()

		if strings.HasPrefix(nodePath, selectedNodePath) && nodePath != selectedNodePath {
			return true
		}
	}
	return false
}

func (self *CommitFilesController) isInTreeMode() *types.DisabledReason {
	if !self.context().CommitFileTreeViewModel.InTreeMode() {
		return &types.DisabledReason{Text: self.c.Tr.DisabledInFlatView}
	}

	return nil
}
