package controllers

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/constants"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
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
			Key:         opts.GetKey(opts.Config.Files.CopyFileInfoToClipboard),
			Handler:     self.openCopyMenu,
			Description: self.c.Tr.CopyToClipboardMenu,
			OpensMenu:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.CommitFiles.CheckoutCommitFile),
			Handler:           self.withItem(self.checkout),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Checkout,
			Tooltip:           self.c.Tr.CheckoutCommitFileTooltip,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.Remove),
			Handler:           self.withItems(self.discard),
			GetDisabledReason: self.require(self.itemsSelected()),
			Description:       self.c.Tr.Remove,
			Tooltip:           self.c.Tr.DiscardOldFileChangeTooltip,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.OpenFile),
			Handler:           self.withItem(self.open),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.OpenFile,
			Tooltip:           self.c.Tr.OpenFileTooltip,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.Edit),
			Handler:           self.withItems(self.edit),
			GetDisabledReason: self.require(self.itemsSelected(self.canEditFiles)),
			Description:       self.c.Tr.Edit,
			Tooltip:           self.c.Tr.EditFileTooltip,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.OpenDiffTool),
			Handler:           self.withItem(self.openDiffTool),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.OpenDiffTool,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.Select),
			Handler:           self.withItems(self.toggleForPatch),
			GetDisabledReason: self.require(self.itemsSelected()),
			Description:       self.c.Tr.ToggleAddToPatch,
			Tooltip: utils.ResolvePlaceholderString(self.c.Tr.ToggleAddToPatchTooltip,
				map[string]string{"doc": constants.Links.Docs.CustomPatchDemo},
			),
			DisplayOnScreen: true,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.ToggleStagedAll),
			Handler:     self.withItem(self.toggleAllForPatch),
			Description: self.c.Tr.ToggleAllInPatch,
			Tooltip: utils.ResolvePlaceholderString(self.c.Tr.ToggleAllInPatchTooltip,
				map[string]string{"doc": constants.Links.Docs.CustomPatchDemo},
			),
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.GoInto),
			Handler:           self.withItem(self.enter),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.EnterCommitFile,
			Tooltip:           self.c.Tr.EnterCommitFileTooltip,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.ToggleTreeView),
			Handler:     self.toggleTreeView,
			Description: self.c.Tr.ToggleTreeView,
			Tooltip:     self.c.Tr.ToggleTreeViewTooltip,
		},
		{
			Key:               opts.GetKey(opts.Config.Files.CollapseAll),
			Handler:           self.collapseAll,
			Description:       self.c.Tr.CollapseAll,
			Tooltip:           self.c.Tr.CollapseAllTooltip,
			GetDisabledReason: self.require(self.isInTreeMode),
		},
		{
			Key:               opts.GetKey(opts.Config.Files.ExpandAll),
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

func (self *CommitFilesController) GetOnRenderToMain() func() {
	return func() {
		node := self.context().GetSelected()
		if node == nil {
			return
		}

		from, to := self.context().GetFromAndToForDiff()
		from, reverse := self.c.Modes().Diffing.GetFromAndReverseArgsForDiff(from)

		cmdObj := self.c.Git().WorkingTree.ShowFileDiffCmdObj(from, to, reverse, node.GetPath(), false)
		task := types.NewRunPtyTask(cmdObj.GetCmd())

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

	cmdObj := self.c.Git().WorkingTree.ShowFileDiffCmdObj(from, to, reverse, path, true)
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
		Key:            'n',
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
		Key:            'p',
	}
	copyAbsolutePathItem := &types.MenuItem{
		Label: self.c.Tr.CopyAbsoluteFilePath,
		OnPress: func() error {
			if err := self.c.OS().CopyToClipboard(filepath.Join(self.c.Git().RepoPaths.RepoPath(), node.GetPath())); err != nil {
				return err
			}
			self.c.Toast(self.c.Tr.FilePathCopiedToast)
			return nil
		},
		DisabledReason: self.require(self.singleItemSelected())(),
		Key:            'P',
	}
	copyFileDiffItem := &types.MenuItem{
		Label: self.c.Tr.CopySelectedDiff,
		OnPress: func() error {
			return self.copyDiffToClipboard(node.GetPath(), self.c.Tr.FileDiffCopiedToast)
		},
		DisabledReason: self.require(self.singleItemSelected())(),
		Key:            's',
	}
	copyAllDiff := &types.MenuItem{
		Label: self.c.Tr.CopyAllFilesDiff,
		OnPress: func() error {
			return self.copyDiffToClipboard(".", self.c.Tr.AllFilesDiffCopiedToast)
		},
		DisabledReason: self.require(self.itemsSelected())(),
		Key:            'a',
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
		Key: 'c',
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
	self.c.LogAction(self.c.Tr.Actions.CheckoutFile)
	_, to := self.context().GetFromAndToForDiff()
	if err := self.c.Git().WorkingTree.CheckoutFile(to, node.GetPath()); err != nil {
		return err
	}

	return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
}

func (self *CommitFilesController) discard(selectedNodes []*filetree.CommitFileNode) error {
	parentContext := self.c.Context().Current().GetParentContext()
	if parentContext == nil || parentContext.GetKey() != context.LOCAL_COMMITS_CONTEXT_KEY {
		return errors.New(self.c.Tr.CanOnlyDiscardFromLocalCommits)
	}

	if ok, err := self.c.Helpers().PatchBuilding.ValidateNormalWorkingTreeState(); !ok {
		return err
	}

	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.DiscardFileChangesTitle,
		Prompt: self.c.Tr.DiscardFileChangesPrompt,
		HandleConfirm: func() error {
			return self.c.WithWaitingStatus(self.c.Tr.RebasingStatus, func(gocui.Task) error {
				var filePaths []string
				selectedNodes = normalisedSelectedCommitFileNodes(selectedNodes)

				// Reset the current patch if there is one.
				if self.c.Git().Patch.PatchBuilder.Active() {
					self.c.Git().Patch.PatchBuilder.Reset()
					if err := self.c.Refresh(types.RefreshOptions{Mode: types.BLOCK_UI}); err != nil {
						return err
					}
				}

				for _, node := range selectedNodes {
					err := node.ForEachFile(func(file *models.CommitFile) error {
						filePaths = append(filePaths, file.GetPath())
						return nil
					})
					if err != nil {
						return err
					}
				}

				err := self.c.Git().Rebase.DiscardOldFileChanges(self.c.Model().Commits, self.c.Contexts().LocalCommits.GetSelectedLineIdx(), filePaths)
				if err := self.c.Helpers().MergeAndRebase.CheckMergeOrRebase(err); err != nil {
					return err
				}

				if self.context().RangeSelectEnabled() {
					self.context().GetList().CancelRangeSelect()
				}
				return self.c.Refresh(types.RefreshOptions{Mode: types.SYNC})
			})
		},
	})

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
	if self.c.AppState.DiffContextSize == 0 {
		return fmt.Errorf(self.c.Tr.Actions.NotEnoughContextToStage,
			keybindings.Label(self.c.UserConfig().Keybinding.Universal.IncreaseContextInDiffView))
	}

	toggle := func() error {
		return self.c.WithWaitingStatus(self.c.Tr.UpdatingPatch, func(gocui.Task) error {
			if !self.c.Git().Patch.PatchBuilder.Active() {
				if err := self.startPatchBuilder(); err != nil {
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

			self.c.PostRefreshUpdate(self.context())
			return nil
		})
	}

	from, to, reverse := self.currentFromToReverseForPatchBuilding()
	if self.c.Git().Patch.PatchBuilder.Active() && self.c.Git().Patch.PatchBuilder.NewPatchRequired(from, to, reverse) {
		self.c.Confirm(types.ConfirmOpts{
			Title:  self.c.Tr.DiscardPatch,
			Prompt: self.c.Tr.DiscardPatchConfirm,
			HandleConfirm: func() error {
				self.c.Git().Patch.PatchBuilder.Reset()
				return toggle()
			},
		})

		return nil
	}

	return toggle()
}

func (self *CommitFilesController) toggleAllForPatch(_ *filetree.CommitFileNode) error {
	root := self.context().CommitFileTreeViewModel.GetRoot()
	return self.toggleForPatch([]*filetree.CommitFileNode{root})
}

func (self *CommitFilesController) startPatchBuilder() error {
	commitFilesContext := self.context()

	canRebase := commitFilesContext.GetCanRebase()
	from, to, reverse := self.currentFromToReverseForPatchBuilding()

	self.c.Git().Patch.PatchBuilder.Start(from, to, reverse, canRebase)
	return nil
}

func (self *CommitFilesController) currentFromToReverseForPatchBuilding() (string, string, bool) {
	commitFilesContext := self.context()

	from, to := commitFilesContext.GetFromAndToForDiff()
	from, reverse := self.c.Modes().Diffing.GetFromAndReverseArgsForDiff(from)
	return from, to, reverse
}

func (self *CommitFilesController) enter(node *filetree.CommitFileNode) error {
	return self.enterCommitFile(node, types.OnFocusOpts{ClickedWindowName: "", ClickedViewLineIdx: -1, ClickedViewRealLineIdx: -1})
}

func (self *CommitFilesController) enterCommitFile(node *filetree.CommitFileNode, opts types.OnFocusOpts) error {
	if node.File == nil {
		return self.handleToggleCommitFileDirCollapsed(node)
	}

	if self.c.AppState.DiffContextSize == 0 {
		return fmt.Errorf(self.c.Tr.Actions.NotEnoughContextToStage,
			keybindings.Label(self.c.UserConfig().Keybinding.Universal.IncreaseContextInDiffView))
	}

	enterTheFile := func() error {
		if !self.c.Git().Patch.PatchBuilder.Active() {
			if err := self.startPatchBuilder(); err != nil {
				return err
			}
		}

		self.c.Context().Push(self.c.Contexts().CustomPatchBuilder, opts)
		return nil
	}

	from, to, reverse := self.currentFromToReverseForPatchBuilding()
	if self.c.Git().Patch.PatchBuilder.Active() && self.c.Git().Patch.PatchBuilder.NewPatchRequired(from, to, reverse) {
		self.c.Confirm(types.ConfirmOpts{
			Title:  self.c.Tr.DiscardPatch,
			Prompt: self.c.Tr.DiscardPatchConfirm,
			HandleConfirm: func() error {
				self.c.Git().Patch.PatchBuilder.Reset()
				return enterTheFile()
			},
		})

		return nil
	}

	return enterTheFile()
}

func (self *CommitFilesController) handleToggleCommitFileDirCollapsed(node *filetree.CommitFileNode) error {
	self.context().CommitFileTreeViewModel.ToggleCollapsed(node.GetInternalPath())

	self.c.PostRefreshUpdate(self.context())

	return nil
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

func (self *CommitFilesController) GetOnClickFocusedMainView() func(mainViewName string, clickedLineIdx int) error {
	return func(mainViewName string, clickedLineIdx int) error {
		clickedFile, line, ok := self.c.Helpers().Staging.GetFileAndLineForClickedDiffLine(mainViewName, clickedLineIdx)
		if !ok {
			line = -1
		}

		node := self.getSelectedItem()
		if node == nil {
			return nil
		}

		if !node.IsFile() && ok {
			relativePath, err := filepath.Rel(self.c.Git().RepoPaths.RepoPath(), clickedFile)
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
					self.context().ModelIndexToViewIndex(idx))
				node = self.context().GetSelected()
			}
		}

		return self.enterCommitFile(node, types.OnFocusOpts{ClickedWindowName: "main", ClickedViewLineIdx: line, ClickedViewRealLineIdx: line})
	}
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
