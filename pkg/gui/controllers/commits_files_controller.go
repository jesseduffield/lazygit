package controllers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/git_commands"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/filetree"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
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
			Description:       self.c.Tr.CheckoutCommitFile,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.Remove),
			Handler:           self.withItem(self.discard),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.DiscardOldFileChange,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.OpenFile),
			Handler:           self.withItem(self.open),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.OpenFile,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.Edit),
			Handler:           self.withItem(self.edit),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.EditFile,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.OpenDiffTool),
			Handler:           self.withItem(self.openDiffTool),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.OpenDiffTool,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.Select),
			Handler:           self.withItem(self.toggleForPatch),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.ToggleAddToPatch,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.ToggleStagedAll),
			Handler:     self.withItem(self.toggleAllForPatch),
			Description: self.c.Tr.ToggleAllInPatch,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.GoInto),
			Handler:           self.withItem(self.enter),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.EnterFile,
		},
		{
			Key:         opts.GetKey(opts.Config.Files.ToggleTreeView),
			Handler:     self.toggleTreeView,
			Description: self.c.Tr.ToggleTreeView,
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

func (self *CommitFilesController) discard(node *filetree.CommitFileNode) error {
	parentContext, ok := self.c.CurrentContext().GetParentContext()
	if !ok || parentContext.GetKey() != context.LOCAL_COMMITS_CONTEXT_KEY {
		return self.c.ErrorMsg(self.c.Tr.CanOnlyDiscardFromLocalCommits)
	}

	if node.File == nil {
		return self.c.ErrorMsg(self.c.Tr.DiscardNotSupportedForDirectory)
	}

	if ok, err := self.c.Helpers().PatchBuilding.ValidateNormalWorkingTreeState(); !ok {
		return err
	}

	prompt := self.c.Tr.DiscardFileChangesPrompt
	if node.File.Added() {
		prompt = self.c.Tr.DiscardAddedFileChangesPrompt
	} else if node.File.Deleted() {
		prompt = self.c.Tr.DiscardDeletedFileChangesPrompt
	}

	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.DiscardFileChangesTitle,
		Prompt: prompt,
		HandleConfirm: func() error {
			return self.c.WithWaitingStatus(self.c.Tr.RebasingStatus, func(gocui.Task) error {
				self.c.LogAction(self.c.Tr.Actions.DiscardOldFileChange)
				if err := self.c.Git().Rebase.DiscardOldFileChanges(self.c.Model().Commits, self.c.Contexts().LocalCommits.GetSelectedLineIdx(), node.GetPath()); err != nil {
					if err := self.c.Helpers().MergeAndRebase.CheckMergeOrRebase(err); err != nil {
						return err
					}
				}

				return self.c.Refresh(types.RefreshOptions{Mode: types.BLOCK_UI})
			})
		},
	})
}

func (self *CommitFilesController) open(node *filetree.CommitFileNode) error {
	return self.c.Helpers().Files.OpenFile(node.GetPath())
}

func (self *CommitFilesController) edit(node *filetree.CommitFileNode) error {
	if node.File == nil {
		return self.c.ErrorMsg(self.c.Tr.ErrCannotEditDirectory)
	}

	return self.c.Helpers().Files.EditFile(node.GetPath())
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

func (self *CommitFilesController) toggleForPatch(node *filetree.CommitFileNode) error {
	toggle := func() error {
		return self.c.WithWaitingStatus(self.c.Tr.UpdatingPatch, func(gocui.Task) error {
			if !self.c.Git().Patch.PatchBuilder.Active() {
				if err := self.startPatchBuilder(); err != nil {
					return err
				}
			}

			// if there is any file that hasn't been fully added we'll fully add everything,
			// otherwise we'll remove everything
			adding := node.SomeFile(func(file *models.CommitFile) bool {
				return self.c.Git().Patch.PatchBuilder.GetFileStatus(file.Name, self.context().GetRef().RefName()) != patch.WHOLE
			})

			err := node.ForEachFile(func(file *models.CommitFile) error {
				if adding {
					return self.c.Git().Patch.PatchBuilder.AddFileWhole(file.Name)
				} else {
					return self.c.Git().Patch.PatchBuilder.RemoveFile(file.Name)
				}
			})
			if err != nil {
				return self.c.Error(err)
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
	return self.toggleForPatch(root)
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
