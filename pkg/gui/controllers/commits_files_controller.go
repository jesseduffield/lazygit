package controllers

import (
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/constants"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
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
		ListControllerTrait: NewListControllerTrait[*filetree.CommitFileNode](
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
	}

	return bindings
}

func (self *CommitFilesController) GetMouseKeybindings(opts types.KeybindingsOpts) []*gocui.ViewMouseBinding {
	return []*gocui.ViewMouseBinding{
		{
			ViewName:    "patchBuilding",
			Key:         gocui.MouseLeft,
			Handler:     self.onClickMain,
			FocusedView: self.context().GetViewName(),
		},
	}
}

func (self *CommitFilesController) context() *context.CommitFilesContext {
	return self.c.Contexts().CommitFiles
}

func (self *CommitFilesController) GetOnRenderToMain() func() error {
	return func() error {
		node := self.context().GetSelected()
		if node == nil {
			return nil
		}

		ref := self.context().GetRef()
		to := ref.RefName()
		from, reverse := self.c.Modes().Diffing.GetFromAndReverseArgsForDiff(ref.ParentRefName())

		cmdObj := self.c.Git().WorkingTree.ShowFileDiffCmdObj(from, to, reverse, node.GetPath(), false)
		task := types.NewRunPtyTask(cmdObj.GetCmd())

		pair := self.c.MainViewPairs().Normal
		if node.File != nil {
			pair = self.c.MainViewPairs().PatchBuilding
		}

		return self.c.RenderToMainViews(types.RefreshMainOpts{
			Pair: pair,
			Main: &types.ViewUpdateOpts{
				Title:    self.c.Tr.Patch,
				SubTitle: self.c.Helpers().Diff.IgnoringWhitespaceSubTitle(),
				Task:     task,
			},
			Secondary: secondaryPatchPanelUpdateOpts(self.c),
		})
	}
}

func (self *CommitFilesController) onClickMain(opts gocui.ViewMouseBindingOpts) error {
	node := self.context().GetSelected()
	if node == nil {
		return nil
	}
	return self.enterCommitFile(node, types.OnFocusOpts{ClickedWindowName: "main", ClickedViewLineIdx: opts.Y})
}

func (self *CommitFilesController) checkout(node *filetree.CommitFileNode) error {
	self.c.LogAction(self.c.Tr.Actions.CheckoutFile)
	if err := self.c.Git().WorkingTree.CheckoutFile(self.context().GetRef().RefName(), node.GetPath()); err != nil {
		return self.c.Error(err)
	}

	return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
}

func (self *CommitFilesController) discard(selectedNodes []*filetree.CommitFileNode) error {
	parentContext, ok := self.c.CurrentContext().GetParentContext()
	if !ok || parentContext.GetKey() != context.LOCAL_COMMITS_CONTEXT_KEY {
		return self.c.ErrorMsg(self.c.Tr.CanOnlyDiscardFromLocalCommits)
	}

	if ok, err := self.c.Helpers().PatchBuilding.ValidateNormalWorkingTreeState(); !ok {
		return err
	}

	return self.c.Confirm(types.ConfirmOpts{
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
						return self.c.Error(err)
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
	ref := self.context().GetRef()
	to := ref.RefName()
	from, reverse := self.c.Modes().Diffing.GetFromAndReverseArgsForDiff(ref.ParentRefName())
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
					fileStatus := self.c.Git().Patch.PatchBuilder.GetFileStatus(file.Name, self.context().GetRef().RefName())
					return fileStatus == patch.PART || fileStatus == patch.UNSELECTED
				})
			})

			patchOperationFunction := self.c.Git().Patch.PatchBuilder.RemoveFile

			if adding {
				patchOperationFunction = self.c.Git().Patch.PatchBuilder.AddFileWhole
			}

			for _, node := range selectedNodes {
				err := node.ForEachFile(func(file *models.CommitFile) error {
					return patchOperationFunction(file.Name)
				})
				if err != nil {
					return self.c.Error(err)
				}
			}

			if self.c.Git().Patch.PatchBuilder.IsEmpty() {
				self.c.Git().Patch.PatchBuilder.Reset()
			}

			return self.c.PostRefreshUpdate(self.context())
		})
	}

	if self.c.Git().Patch.PatchBuilder.Active() && self.c.Git().Patch.PatchBuilder.To != self.context().GetRef().RefName() {
		return self.c.Confirm(types.ConfirmOpts{
			Title:  self.c.Tr.DiscardPatch,
			Prompt: self.c.Tr.DiscardPatchConfirm,
			HandleConfirm: func() error {
				self.c.Git().Patch.PatchBuilder.Reset()
				return toggle()
			},
		})
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
	ref := commitFilesContext.GetRef()
	to := ref.RefName()
	from, reverse := self.c.Modes().Diffing.GetFromAndReverseArgsForDiff(ref.ParentRefName())

	self.c.Git().Patch.PatchBuilder.Start(from, to, reverse, canRebase)
	return nil
}

func (self *CommitFilesController) enter(node *filetree.CommitFileNode) error {
	return self.enterCommitFile(node, types.OnFocusOpts{ClickedWindowName: "", ClickedViewLineIdx: -1})
}

func (self *CommitFilesController) enterCommitFile(node *filetree.CommitFileNode, opts types.OnFocusOpts) error {
	if node.File == nil {
		return self.handleToggleCommitFileDirCollapsed(node)
	}

	enterTheFile := func() error {
		if !self.c.Git().Patch.PatchBuilder.Active() {
			if err := self.startPatchBuilder(); err != nil {
				return err
			}
		}

		return self.c.PushContext(self.c.Contexts().CustomPatchBuilder, opts)
	}

	if self.c.Git().Patch.PatchBuilder.Active() && self.c.Git().Patch.PatchBuilder.To != self.context().GetRef().RefName() {
		return self.c.Confirm(types.ConfirmOpts{
			Title:  self.c.Tr.DiscardPatch,
			Prompt: self.c.Tr.DiscardPatchConfirm,
			HandleConfirm: func() error {
				self.c.Git().Patch.PatchBuilder.Reset()
				return enterTheFile()
			},
		})
	}

	return enterTheFile()
}

func (self *CommitFilesController) handleToggleCommitFileDirCollapsed(node *filetree.CommitFileNode) error {
	self.context().CommitFileTreeViewModel.ToggleCollapsed(node.GetPath())

	if err := self.c.PostRefreshUpdate(self.context()); err != nil {
		self.c.Log.Error(err)
	}

	return nil
}

// NOTE: this is very similar to handleToggleFileTreeView, could be DRY'd with generics
func (self *CommitFilesController) toggleTreeView() error {
	self.context().CommitFileTreeViewModel.ToggleShowTree()

	return self.c.PostRefreshUpdate(self.context())
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
