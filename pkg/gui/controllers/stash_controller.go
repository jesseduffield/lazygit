package controllers

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
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
		ListControllerTrait: NewListControllerTrait(
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
			Keys:              opts.GetKeys(opts.Config.Universal.Select),
			Handler:           self.withItem(self.handleStashApply),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Apply,
			Tooltip:           self.c.Tr.StashApplyTooltip,
			DisplayOnScreen:   true,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Stash.PopStash),
			Handler:           self.withItem(self.handleStashPop),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.Pop,
			Tooltip:           self.c.Tr.StashPopTooltip,
			DisplayOnScreen:   true,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Universal.Remove),
			Handler:           self.withItems(self.handleStashDrop),
			GetDisabledReason: self.require(self.itemRangeSelected()),
			Description:       self.c.Tr.Drop,
			Tooltip:           self.c.Tr.StashDropTooltip,
			DisplayOnScreen:   true,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Universal.New),
			Handler:           self.withItem(self.handleNewBranchOffStashEntry),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.NewBranch,
			Tooltip:           self.c.Tr.NewBranchFromStashTooltip,
		},
		{
			Keys:        opts.GetKeys(opts.Config.Universal.NewWorktree),
			Handler:     self.withItem(self.c.Helpers().Worktree.NewWorktreeMenuForStash),
			Description: self.c.Tr.NewWorktree,
			OpensMenu:   true,
		},
		{
			Keys:              opts.GetKeys(opts.Config.Stash.RenameStash),
			Handler:           self.withItem(self.handleRenameStashEntry),
			GetDisabledReason: self.require(self.singleItemSelected()),
			Description:       self.c.Tr.RenameStash,
		},
	}

	return bindings
}

func (self *StashController) GetOnRenderToMain() func() {
	return func() {
		self.c.Helpers().Diff.WithDiffModeCheck(func() {
			var task types.UpdateTask
			stashEntry := self.context().GetSelected()
			if stashEntry == nil {
				task = types.NewRenderStringTask(self.c.Tr.NoStashEntries)
			} else {
				prefix := style.FgYellow.Sprintf("%s\n\n", stashEntry.Description())
				task = types.NewRunPtyTaskWithPrefix(
					self.c.Git().Stash.ShowStashEntryCmdObj(stashEntry.Index).GetCmd(),
					prefix,
				)
			}

			self.c.RenderToMainViews(types.RefreshMainOpts{
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
	return self.c.ConfirmIf(!self.c.UserConfig().Gui.SkipStashWarning,
		types.ConfirmOpts{
			Title:  self.c.Tr.StashApply,
			Prompt: self.c.Tr.SureApplyStashEntry,
			HandleConfirm: func() error {
				return self.c.WithWaitingStatusBlockingInput(
					types.WaitingStatusOpts{Message: self.c.Tr.ApplyingStashStatus},
					func(gocui.Task) error {
						self.c.LogAction(self.c.Tr.Actions.ApplyStash)
						err := self.c.Git().Stash.Apply(stashEntry.Index)
						self.postStashRefresh(err == nil && self.c.UserConfig().Gui.SwitchToFilesAfterStashApply)
						return err
					})
			},
		})
}

func (self *StashController) handleStashPop(stashEntry *models.StashEntry) error {
	pop := func() error {
		return self.c.WithWaitingStatusBlockingInput(
			types.WaitingStatusOpts{Message: self.c.Tr.PoppingStashStatus},
			func(gocui.Task) error {
				self.c.LogAction(self.c.Tr.Actions.PopStash)
				self.c.LogCommand(fmt.Sprintf(self.c.Tr.Log.PoppingStash, stashEntry.Hash), false)
				err := self.c.Git().Stash.Pop(stashEntry.Index)
				self.postStashRefresh(err == nil && self.c.UserConfig().Gui.SwitchToFilesAfterStashPop)
				return err
			})
	}

	if self.c.UserConfig().Gui.SkipStashWarning {
		return pop()
	}

	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.StashPop,
		Prompt: self.c.Tr.SurePopStashEntry,
		HandleConfirm: func() error {
			return pop()
		},
	})

	return nil
}

func (self *StashController) handleStashDrop(stashEntries []*models.StashEntry) error {
	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.StashDrop,
		Prompt: self.c.Tr.SureDropStashEntry,
		HandleConfirm: func() error {
			self.c.LogAction(self.c.Tr.Actions.DropStash)
			// Refresh once at the end rather than after each drop: a refresh
			// from the UI thread finishes in the background, so firing one per
			// iteration lets the workers race and an earlier, stale result can
			// land last. The indices are captured up front and we drop
			// highest-first, so the remaining lower indices stay valid without
			// an intervening refresh.
			var dropErr error
			for i := len(stashEntries) - 1; i >= 0; i-- {
				self.c.LogCommand(fmt.Sprintf(self.c.Tr.Log.DroppingStash, stashEntries[i].Hash), false)
				if dropErr = self.c.Git().Stash.Drop(stashEntries[i].Index); dropErr != nil {
					break
				}
			}
			// Block input until the refresh has landed, so that dropping the
			// next entry in quick succession (confirming and pressing the key
			// again right away) sees the refreshed list and not the stale,
			// pre-drop indices.
			self.c.RefreshBlockingInput(types.RefreshOptions{
				Scope: []types.RefreshableView{types.STASH},
				Then: func() error {
					// Collapse the range selection from here, so that it lands
					// in the same frame as the shortened list. The refresh has
					// painted the list by the time Then runs, so the new
					// selection needs a focus update of its own.
					if dropErr == nil {
						self.context().CollapseRangeSelectionToTop()
						self.context().HandleFocus(types.OnFocusOpts{})
					}
					return nil
				},
			})
			return dropErr
		},
	})

	return nil
}

// postStashRefresh refreshes the panels that applying or popping a stash
// affects, moving the focus to the files panel if switchToFiles is set.
//
// Call it from the worker that ran the stash command, from inside a
// WithWaitingStatusBlockingInput: popping shifts the indices of the remaining
// stash entries, so acting on the next entry in quick succession (confirming
// the popup and pressing the key again right away) has to be held off until
// the refreshed list is in place, or it would target the wrong stash.
func (self *StashController) postStashRefresh(switchToFiles bool) {
	self.c.RefreshFromWorker(types.RefreshOptions{
		BatchUIUpdates: true,
		Scope:          []types.RefreshableView{types.STASH, types.FILES},
		Then: func() error {
			// Switch panels from here, so that the focus change lands in the
			// same frame as the refreshed panel contents.
			if switchToFiles {
				self.c.Context().Push(self.c.Contexts().Files, types.OnFocusOpts{})
			}
			return nil
		},
	})
}

func (self *StashController) handleNewBranchOffStashEntry(stashEntry *models.StashEntry) error {
	return self.c.Helpers().Refs.NewBranch(stashEntry.FullRefName(), stashEntry.Description(), "")
}

func (self *StashController) handleRenameStashEntry(stashEntry *models.StashEntry) error {
	message := utils.ResolvePlaceholderString(
		self.c.Tr.RenameStashPrompt,
		map[string]string{
			"stashName": stashEntry.RefName(),
		},
	)

	self.c.Prompt(types.PromptOpts{
		Title:          message,
		InitialContent: stashEntry.Name,
		HandleConfirm: func(response string) error {
			self.c.LogAction(self.c.Tr.Actions.RenameStash)
			err := self.c.Git().Stash.Rename(stashEntry.Index, response)
			if err != nil {
				self.c.RefreshBlockingInput(types.RefreshOptions{Scope: []types.RefreshableView{types.STASH}})
				return err
			}
			self.context().SetSelection(0) // Select the renamed stash
			self.context().FocusLine(true)
			// Renaming re-creates the stash at the top, shifting the other
			// entries' indices; block input so that a quick next action sees
			// the refreshed list rather than the stale indices.
			self.c.RefreshBlockingInput(types.RefreshOptions{Scope: []types.RefreshableView{types.STASH}})
			return nil
		},
		AllowEmptyInput: true,
	})

	return nil
}
