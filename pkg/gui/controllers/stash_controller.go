package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type StashController struct {
	baseController
	*ListControllerTrait[*models.StashEntry]
	c *ControllerCommon
}

var _ types.IController = &StashController{}

func NewStashController(
	c *ControllerCommon,
) *StashController {
	return &StashController{
		baseController: baseController{},
		ListControllerTrait: NewListControllerTrait[*models.StashEntry](
			c,
			c.Contexts().Stash,
			c.Contexts().Stash.GetSelected,
			c.Contexts().Stash.GetSelectedItems,
		),
		c: c,
	}
}

func (self *StashController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:               opts.GetKey(opts.Config.Universal.Select),
			Handler:           self.withItem(self.handleStashApply),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Apply,
			Tooltip:           self.c.Tr.StashApplyTooltip,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Stash.PopStash),
			Handler:           self.withItem(self.handleStashPop),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Pop,
			Tooltip:           self.c.Tr.StashPopTooltip,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.Remove),
			Handler:           self.withItem(self.handleStashDrop),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Drop,
			Tooltip:           self.c.Tr.StashDropTooltip,
			DisplayOnScreen:   true,
		},
		{
			Key:               opts.GetKey(opts.Config.Universal.New),
			Handler:           self.withItem(self.handleNewBranchOffStashEntry),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.NewBranch,
			Tooltip:           self.c.Tr.NewBranchFromStashTooltip,
		},
		{
			Key:               opts.GetKey(opts.Config.Stash.RenameStash),
			Handler:           self.withItem(self.handleRenameStashEntry),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.RenameStash,
		},
	}

	return bindings
}

func (self *StashController) GetOnRenderToMain() func() error {
	return func() error {
		return self.c.Helpers().Diff.WithDiffModeCheck(func() error {
			var task types.UpdateTask
			stashEntry := self.context().GetSelected()
			if stashEntry == nil {
				task = types.NewRenderStringTask(self.c.Tr.NoStashEntries)
			} else {
				task = types.NewRunPtyTask(
					self.c.Git().Stash.ShowStashEntryCmdObj(stashEntry.Index).GetCmd(),
				)
			}

			return self.c.RenderToMainViews(types.RefreshMainOpts{
				Pair: self.c.MainViewPairs().Normal,
				Main: &types.ViewUpdateOpts{
					Title:    "Stash",
					SubTitle: self.c.Helpers().Diff.IgnoringWhitespaceSubTitle(),
					Task:     task,
				},
			})
		})
	}
}

func (self *StashController) context() *context.StashContext {
	return self.c.Contexts().Stash
}

func (self *StashController) handleStashApply(stashEntry *models.StashEntry) error {
	apply := func() error {
		self.c.LogAction(self.c.Tr.Actions.Stash)
		err := self.c.Git().Stash.Apply(stashEntry.Index)
		_ = self.postStashRefresh()
		if err != nil {
			return err
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
		err := self.c.Git().Stash.Pop(stashEntry.Index)
		_ = self.postStashRefresh()
		if err != nil {
			return err
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
			err := self.c.Git().Stash.Drop(stashEntry.Index)
			_ = self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.STASH}})
			if err != nil {
				return err
			}
			return nil
		},
	})
}

func (self *StashController) postStashRefresh() error {
	return self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.STASH, types.FILES}})
}

func (self *StashController) handleNewBranchOffStashEntry(stashEntry *models.StashEntry) error {
	return self.c.Helpers().Refs.NewBranch(stashEntry.RefName(), stashEntry.Description(), "")
}

func (self *StashController) handleRenameStashEntry(stashEntry *models.StashEntry) error {
	message := utils.ResolvePlaceholderString(
		self.c.Tr.RenameStashPrompt,
		map[string]string{
			"stashName": stashEntry.RefName(),
		},
	)

	return self.c.Prompt(types.PromptOpts{
		Title:          message,
		InitialContent: stashEntry.Name,
		HandleConfirm: func(response string) error {
			self.c.LogAction(self.c.Tr.Actions.RenameStash)
			err := self.c.Git().Stash.Rename(stashEntry.Index, response)
			_ = self.c.Refresh(types.RefreshOptions{Scope: []types.RefreshableView{types.STASH}})
			if err != nil {
				return err
			}
			self.context().SetSelection(0) // Select the renamed stash
			self.context().FocusLine()
			return nil
		},
	})
}
