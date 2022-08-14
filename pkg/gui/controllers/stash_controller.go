package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type StashController struct {
	baseController
	*controllerCommon
}

var _ types.IController = &StashController{}

func NewStashController(
	common *controllerCommon,
) *StashController {
	return &StashController{
		baseController:   baseController{},
		controllerCommon: common,
	}
}

func (self *StashController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.Select),
			Handler:     self.checkSelected(self.handleStashApply),
			Description: self.c.Tr.LcApply,
		},
		{
			Key:         opts.GetKey(opts.Config.Stash.PopStash),
			Handler:     self.checkSelected(self.handleStashPop),
			Description: self.c.Tr.LcPop,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Remove),
			Handler:     self.checkSelected(self.handleStashDrop),
			Description: self.c.Tr.LcDrop,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.New),
			Handler:     self.checkSelected(self.handleNewBranchOffStashEntry),
			Description: self.c.Tr.LcNewBranch,
		},
	}

	return bindings
}

func (self *StashController) checkSelected(callback func(*models.StashEntry) error) func() error {
	return func() error {
		item := self.context().GetSelected()
		if item == nil {
			return nil
		}

		return callback(item)
	}
}

func (self *StashController) Context() types.Context {
	return self.context()
}

func (self *StashController) context() *context.StashContext {
	return self.contexts.Stash
}

func (self *StashController) handleStashApply(stashEntry *models.StashEntry) error {
	apply := func() error {
		self.c.LogAction(self.c.Tr.Actions.Stash)
		err := self.git.Stash.Apply(stashEntry.Index)
		_ = self.postStashRefresh()
		if err != nil {
			return self.c.Error(err)
		}
		return nil
	}

	if self.c.UserConfig.Gui.SkipStashWarning {
		return apply()
	}

	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.StashApply,
		Prompt: self.c.Tr.SureApplyStashEntry,
		HandleConfirm: func() error {
			return apply()
		},
	})
}

func (self *StashController) handleStashPop(stashEntry *models.StashEntry) error {
	pop := func() error {
		self.c.LogAction(self.c.Tr.Actions.Stash)
		err := self.git.Stash.Pop(stashEntry.Index)
		_ = self.postStashRefresh()
		if err != nil {
			return self.c.Error(err)
		}
		return nil
	}

	if self.c.UserConfig.Gui.SkipStashWarning {
		return pop()
	}

	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.StashPop,
		Prompt: self.c.Tr.SurePopStashEntry,
		HandleConfirm: func() error {
			return pop()
		},
	})
}

func (self *StashController) handleStashDrop(stashEntry *models.StashEntry) error {
	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.StashDrop,
		Prompt: self.c.Tr.SureDropStashEntry,
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.Stash)
			err := self.git.Stash.Drop(stashEntry.Index)
			_ = self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.STASH}})
			if err != nil {
				return self.c.Error(err)
			}
			return nil
		},
	})
}

func (self *StashController) postStashRefresh() error {
	return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.STASH, types.FILES}})
}

func (self *StashController) handleNewBranchOffStashEntry(stashEntry *models.StashEntry) error {
	return self.helpers.Refs.NewBranch(stashEntry.RefName(), stashEntry.Description(), "")
}
