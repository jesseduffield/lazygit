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

	if self.c.UserConfig().ConfirmOnQuit {
		self.c.Confirm(types.ConfirmOpts{
			Title:  "",
			Prompt: self.c.Tr.ConfirmQuit,
			HandleConfirm: func() error {
				return gocui.ErrQuit
			},
		})

		return nil
	}

	return gocui.ErrQuit
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
	currentContext := self.c.Context().Current()

	if listContext, ok := currentContext.(types.IListContext); ok {
		if listContext.GetList().IsSelectingRange() {
			listContext.GetList().CancelRangeSelect()
			return self.c.PostRefreshUpdate(listContext)
		}
	}

	switch ctx := currentContext.(type) {
	case types.IFilterableContext:
		if ctx.IsFiltering() {
			self.c.Helpers().Search.Cancel()
			return nil
		}
	case types.ISearchableContext:
		if ctx.IsSearching() {
			self.c.Helpers().Search.Cancel()
			return nil
		}
	}

	parentContext := currentContext.GetParentContext()
	if parentContext != nil {
		// TODO: think about whether this should be marked as a return rather than adding to the stack
		self.c.Context().Push(parentContext)
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
