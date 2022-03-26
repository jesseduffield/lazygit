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
			Handler:     self.checkSelected(self.checkout),
			Description: self.c.Tr.LcCheckoutCommitFile,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Remove),
			Handler:     self.checkSelected(self.discard),
			Description: self.c.Tr.LcDiscardOldFileChange,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.OpenFile),
			Handler:     self.checkSelected(self.open),
			Description: self.c.Tr.LcOpenFile,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Edit),
			Handler:     self.checkSelected(self.edit),
			Description: self.c.Tr.LcEditFile,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Select),
			Handler:     self.checkSelected(self.toggleForPatch),
			Description: self.c.Tr.LcToggleAddToPatch,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.ToggleStagedAll),
			Handler:     self.checkSelected(self.toggleAllForPatch),
			Description: self.c.Tr.LcToggleAllInPatch,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.GoInto),
			Handler:     self.checkSelected(self.enter),
			Description: self.c.Tr.LcEnterFile,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.ToggleTreeView),
			Handler:     self.toggleTreeView,
			Description: self.c.Tr.LcToggleTreeView,
		},
	}

	return bindings
}

func (self *CommitFilesController) GetMouseKeybindings(opts types.KeybindingsOpts) []*gocui.ViewMouseBinding {
	return []*gocui.ViewMouseBinding{
		{
			ViewName:    "main",
			Key:         gocui.MouseLeft,
			Handler:     self.onClickMain,
			FromContext: string(self.context().GetKey()),
		},
	}
}

func (self *CommitFilesController) checkSelected(callback func(*filetree.CommitFileNode) error) func() error {
	return func() error {
		selected := self.context().GetSelected()
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
	node := self.context().GetSelected()
	if node == nil {
		return nil
	}
	return self.enterCommitFile(node, types.OnFocusOpts{ClickedViewName: "main", ClickedViewLineIdx: opts.Y})
}

func (self *CommitFilesController) checkout(node *filetree.CommitFileNode) error {
	self.c.LogAction(self.c.Tr.Actions.CheckoutFile)
	if err := self.git.WorkingTree.CheckoutFile(self.context().GetRef().RefName(), node.GetPath()); err != nil {
		return self.c.Error(err)
	}

	return self.c.Refresh(types.RefreshOptions{Mode: types.ASYNC})
}

func (self *CommitFilesController) discard(node *filetree.CommitFileNode) error {
	if ok, err := self.helpers.PatchBuilding.ValidateNormalWorkingTreeState(); !ok {
		return err
	}

	return self.c.Confirm(types.ConfirmOpts{
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

func (self *CommitFilesController) open(node *filetree.CommitFileNode) error {
	return self.helpers.Files.OpenFile(node.GetPath())
}

func (self *CommitFilesController) edit(node *filetree.CommitFileNode) error {
	if node.File == nil {
		return self.c.ErrorMsg(self.c.Tr.ErrCannotEditDirectory)
	}

	return self.helpers.Files.EditFile(node.GetPath())
}

func (self *CommitFilesController) toggleForPatch(node *filetree.CommitFileNode) error {
	toggle := func() error {
		return self.c.WithWaitingStatus(self.c.Tr.LcUpdatingPatch, func() error {
			if !self.git.Patch.PatchManager.Active() {
				if err := self.startPatchManager(); err != nil {
					return err
				}
			}

			// if there is any file that hasn't been fully added we'll fully add everything,
			// otherwise we'll remove everything
			adding := node.AnyFile(func(file *models.CommitFile) bool {
				return self.git.Patch.PatchManager.GetFileStatus(file.Name, self.context().GetRef().RefName()) != patch.WHOLE
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
		})
	}

	if self.git.Patch.PatchManager.Active() && self.git.Patch.PatchManager.To != self.context().GetRef().RefName() {
		return self.c.Confirm(types.ConfirmOpts{
			Title:  self.c.Tr.DiscardPatch,
			Prompt: self.c.Tr.DiscardPatchConfirm,
			HandleConfirm: func() error {
				self.git.Patch.PatchManager.Reset()
				return toggle()
			},
		})
	}

	return toggle()
}

func (self *CommitFilesController) toggleAllForPatch(_ *filetree.CommitFileNode) error {
	// not a fan of type assertions but this will be fixed very soon thanks to generics
	root := self.context().CommitFileTreeViewModel.Tree().(*filetree.CommitFileNode)
	return self.toggleForPatch(root)
}

func (self *CommitFilesController) startPatchManager() error {
	commitFilesContext := self.context()

	canRebase := commitFilesContext.GetCanRebase()
	ref := commitFilesContext.GetRef()
	to := ref.RefName()
	from, reverse := self.modes.Diffing.GetFromAndReverseArgsForDiff(ref.ParentRefName())

	self.git.Patch.PatchManager.Start(from, to, reverse, canRebase)
	return nil
}

func (self *CommitFilesController) enter(node *filetree.CommitFileNode) error {
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

	if self.git.Patch.PatchManager.Active() && self.git.Patch.PatchManager.To != self.context().GetRef().RefName() {
		return self.c.Confirm(types.ConfirmOpts{
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
func (self *CommitFilesController) toggleTreeView() error {
	self.context().CommitFileTreeViewModel.ToggleShowTree()

	return self.c.PostRefreshUpdate(self.context())
}
