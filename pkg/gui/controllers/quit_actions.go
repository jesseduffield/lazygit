package controllers

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type QuitActions struct {
	c *ControllerCommon
}

func (self *QuitActions) Quit() error {
	self.c.State().SetRetainOriginalDir(false)
	return self.quitAux()
}

func (self *QuitActions) QuitWithoutChangingDirectory() error {
	self.c.State().SetRetainOriginalDir(true)
	return self.quitAux()
}

func (self *QuitActions) quitAux() error {
	if self.c.State().GetUpdating() {
		return self.confirmQuitDuringUpdate()
	}

	return self.c.ConfirmIf(self.c.UserConfig().ConfirmOnQuit,
		types.ConfirmOpts{
			Title:  "",
			Prompt: self.c.Tr.ConfirmQuit,
			HandleConfirm: func() error {
				return gocui.ErrQuit
			},
		})
}

func (self *QuitActions) confirmQuitDuringUpdate() error {
	self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.ConfirmQuitDuringUpdateTitle,
		Prompt: self.c.Tr.ConfirmQuitDuringUpdate,
		HandleConfirm: func() error {
			return gocui.ErrQuit
		},
	})

	return nil
}

func (self *QuitActions) Escape() error {
	// If you make changes to this function, be sure to update EscapeEnabled and EscapeDescription accordingly.

	currentContext := self.c.Context().Current()

	if listContext, ok := currentContext.(types.IListContext); ok {
		if listContext.GetList().IsSelectingRange() {
			listContext.GetList().CancelRangeSelect()
			self.c.PostRefreshUpdate(listContext)
			return nil
		}
	}

	// Cancelling searching (as opposed to filtering) is handled by gocui
	if ctx, ok := currentContext.(types.IFilterableContext); ok {
		if ctx.IsFiltering() {
			self.c.Helpers().Search.Cancel()
			return nil
		}
	}

	parentContext := currentContext.GetParentContext()
	if parentContext != nil {
		// TODO: think about whether this should be marked as a return rather than adding to the stack
		self.c.Context().Push(parentContext, types.OnFocusOpts{})
		return nil
	}

	for _, mode := range self.c.Helpers().Mode.Statuses() {
		if mode.IsActive() {
			return mode.Reset()
		}
	}

	repoPathStack := self.c.State().GetRepoPathStack()
	if !repoPathStack.IsEmpty() {
		return self.c.Helpers().Repos.DispatchSwitchToRepo(repoPathStack.Pop(), context.NO_CONTEXT)
	}

	if self.c.UserConfig().QuitOnTopLevelReturn {
		return self.Quit()
	}

	return nil
}

func (self *QuitActions) EscapeEnabled() bool {
	currentContext := self.c.Context().Current()

	if listContext, ok := currentContext.(types.IListContext); ok {
		if listContext.GetList().IsSelectingRange() {
			return true
		}
	}

	if ctx, ok := currentContext.(types.IFilterableContext); ok {
		if ctx.IsFiltering() {
			return true
		}
	}

	parentContext := currentContext.GetParentContext()
	if parentContext != nil {
		return true
	}

	for _, mode := range self.c.Helpers().Mode.Statuses() {
		if mode.IsActive() {
			return true
		}
	}

	repoPathStack := self.c.State().GetRepoPathStack()
	if !repoPathStack.IsEmpty() {
		return true
	}

	if self.c.UserConfig().QuitOnTopLevelReturn {
		return true
	}

	return false
}

func (self *QuitActions) EscapeDescription() string {
	currentContext := self.c.Context().Current()

	if listContext, ok := currentContext.(types.IListContext); ok {
		if listContext.GetList().IsSelectingRange() {
			return self.c.Tr.DismissRangeSelect
		}
	}

	if ctx, ok := currentContext.(types.IFilterableContext); ok {
		if ctx.IsFiltering() {
			return self.c.Tr.ExitFilterMode
		}
	}

	parentContext := currentContext.GetParentContext()
	if parentContext != nil {
		return self.c.Tr.ExitSubview
	}

	for _, mode := range self.c.Helpers().Mode.Statuses() {
		if mode.IsActive() {
			return mode.CancelLabel()
		}
	}

	repoPathStack := self.c.State().GetRepoPathStack()
	if !repoPathStack.IsEmpty() {
		return self.c.Tr.BackToParentRepo
	}

	if self.c.UserConfig().QuitOnTopLevelReturn {
		return self.c.Tr.Quit
	}

	return self.c.Tr.Cancel
}
