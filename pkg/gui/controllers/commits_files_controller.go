package controllers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type CommitFilesController struct {
	baseController
	*controllerCommon
}

var _ types.IController = &CommitFilesController{}

func NewCommitFilesController(
	common *controllerCommon,
) *CommitFilesController {
	return &CommitFilesController{
		baseController:   baseController{},
		controllerCommon: common,
	}
}

func (self *CommitFilesController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.CommitFiles.CheckoutCommitFile),
			Handler:     self.checkSelected(self.handleCheckoutCommitFile),
			Description: self.c.Tr.LcCheckoutCommitFile,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Remove),
			Handler:     self.checkSelected(self.handleDiscardOldFileChange),
			Description: self.c.Tr.LcDiscardOldFileChange,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.OpenFile),
			Handler:     self.checkSelected(self.handleOpenOldCommitFile),
			Description: self.c.Tr.LcOpenFile,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Edit),
			Handler:     self.checkSelected(self.handleEditCommitFile),
			Description: self.c.Tr.LcEditFile,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Select),
			Handler:     self.checkSelected(self.handleToggleFileForPatch),
			Description: self.c.Tr.LcToggleAddToPatch,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.GoInto),
			Handler:     self.checkSelected(self.handleEnterCommitFile),
			Description: self.c.Tr.LcEnterFile,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.ToggleTreeView),
			Handler:     self.handleToggleCommitFileTreeView,
			Description: self.c.Tr.LcToggleTreeView,
		},
	}

	return bindings
}

func (self *CommitFilesController) GetMouseKeybindings(opts types.KeybindingsOpts) []*gocui.ViewMouseBinding {
	return []*gocui.ViewMouseBinding{
		{
			ViewName: "main",
			Key:      gocui.MouseLeft,
			Handler:  self.onClickMain,
		},
	}
}

func (self *CommitFilesController) checkSelected(callback func(*filetree.CommitFileNode) error) func() error {
	return func() error {
		selected := self.context().GetSelectedFileNode()
		if selected == nil {
			return nil
		}

		return callback(selected)
	}
}

func (self *CommitFilesController) Context() types.Context {
	return self.context()
}

func (self *CommitFilesController) context() *context.CommitFilesContext {
	return self.contexts.CommitFiles
}

func (self *CommitFilesController) onClickMain(opts gocui.ViewMouseBindingOpts) error {
	clickedViewLineIdx := opts.Cy + opts.Oy
	node := self.context().GetSelectedFileNode()
	if node == nil {
		return nil
	}
	return self.enterCommitFile(node, types.OnFocusOpts{ClickedViewName: "main", ClickedViewLineIdx: clickedViewLineIdx})
}

func (self *CommitFilesController) handleCheckoutCommitFile(node *filetree.CommitFileNode) error {
	self.c.LogAction(self.c.Tr.Actions.CheckoutFile)
	if err := self.git.WorkingTree.CheckoutFile(self.context().GetRefName(), node.GetPath()); err != nil {
		return self.c.Error(err)
	}

	return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
}

func (self *CommitFilesController) handleDiscardOldFileChange(node *filetree.CommitFileNode) error {
	if ok, err := self.helpers.PatchBuilding.ValidateNormalWorkingTreeState(); !ok {
		return err
	}

	return self.c.Ask(types.AskOpts{
		Title:  self.c.Tr.DiscardFileChangesTitle,
		Prompt: self.c.Tr.DiscardFileChangesPrompt,
		HandleConfirm: func() error {
			return self.c.WithWaitingStatus(self.c.Tr.RebasingStatus, func() error {
				self.c.LogAction(self.c.Tr.Actions.DiscardOldFileChange)
				if err := self.git.Rebase.DiscardOldFileChanges(self.model.Commits, self.contexts.LocalCommits.GetSelectedLineIdx(), node.GetPath()); err != nil {
					if err := self.helpers.MergeAndRebase.CheckMergeOrRebase(err); err != nil {
						return err
					}
				}

				return self.c.Refresh(types.RefreshOptions{Mode: types.BLOCK_UI})
			})
		},
	})
}

func (self *CommitFilesController) handleOpenOldCommitFile(node *filetree.CommitFileNode) error {
	return self.helpers.Files.OpenFile(node.GetPath())
}

func (self *CommitFilesController) handleEditCommitFile(node *filetree.CommitFileNode) error {
	if node.File == nil {
		return self.c.ErrorMsg(self.c.Tr.ErrCannotEditDirectory)
	}

	return self.helpers.Files.EditFile(node.GetPath())
}

func (self *CommitFilesController) handleToggleFileForPatch(node *filetree.CommitFileNode) error {
	toggleTheFile := func() error {
		if !self.git.Patch.PatchManager.Active() {
			if err := self.startPatchManager(); err != nil {
				return err
			}
		}

		// if there is any file that hasn't been fully added we'll fully add everything,
		// otherwise we'll remove everything
		adding := node.AnyFile(func(file *models.CommitFile) bool {
			return self.git.Patch.PatchManager.GetFileStatus(file.Name, self.context().GetRefName()) != patch.WHOLE
		})

		err := node.ForEachFile(func(file *models.CommitFile) error {
			if adding {
				return self.git.Patch.PatchManager.AddFileWhole(file.Name)
			} else {
				return self.git.Patch.PatchManager.RemoveFile(file.Name)
			}
		})

		if err != nil {
			return self.c.Error(err)
		}

		if self.git.Patch.PatchManager.IsEmpty() {
			self.git.Patch.PatchManager.Reset()
		}

		return self.c.PostRefreshUpdate(self.context())
	}

	if self.git.Patch.PatchManager.Active() && self.git.Patch.PatchManager.To != self.context().GetRefName() {
		return self.c.Ask(types.AskOpts{
			Title:  self.c.Tr.DiscardPatch,
			Prompt: self.c.Tr.DiscardPatchConfirm,
			HandleConfirm: func() error {
				self.git.Patch.PatchManager.Reset()
				return toggleTheFile()
			},
		})
	}

	return toggleTheFile()
}

func (self *CommitFilesController) startPatchManager() error {
	commitFilesContext := self.context()

	canRebase := commitFilesContext.GetCanRebase()
	to := commitFilesContext.GetRefName()

	from, reverse := self.modes.Diffing.GetFromAndReverseArgsForDiff(to)

	self.git.Patch.PatchManager.Start(from, to, reverse, canRebase)
	return nil
}

func (self *CommitFilesController) handleEnterCommitFile(node *filetree.CommitFileNode) error {
	return self.enterCommitFile(node, types.OnFocusOpts{ClickedViewName: "", ClickedViewLineIdx: -1})
}

func (self *CommitFilesController) enterCommitFile(node *filetree.CommitFileNode, opts types.OnFocusOpts) error {
	if node.File == nil {
		return self.handleToggleCommitFileDirCollapsed(node)
	}

	enterTheFile := func() error {
		if !self.git.Patch.PatchManager.Active() {
			if err := self.startPatchManager(); err != nil {
				return err
			}
		}

		return self.c.PushContext(self.contexts.PatchBuilding, opts)
	}

	if self.git.Patch.PatchManager.Active() && self.git.Patch.PatchManager.To != self.context().GetRefName() {
		return self.c.Ask(types.AskOpts{
			Title:  self.c.Tr.DiscardPatch,
			Prompt: self.c.Tr.DiscardPatchConfirm,
			HandleConfirm: func() error {
				self.git.Patch.PatchManager.Reset()
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
func (self *CommitFilesController) handleToggleCommitFileTreeView() error {
	self.context().CommitFileTreeViewModel.ToggleShowTree()

	return self.c.PostRefreshUpdate(self.context())
}
