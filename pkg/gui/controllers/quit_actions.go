package controllers

import (
	"github.com/jesseduffield/gocui"
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

	if self.c.UserConfig.ConfirmOnQuit {
		return self.c.Confirm(types.ConfirmOpts{
			Title:  "",
			Prompt: self.c.Tr.ConfirmQuit,
			HandleConfirm: func() error {
				return gocui.ErrQuit
			},
		})
	}

	return gocui.ErrQuit
}

func (self *QuitActions) confirmQuitDuringUpdate() error {
	return self.c.Confirm(types.ConfirmOpts{
		Title:  self.c.Tr.ConfirmQuitDuringUpdateTitle,
		Prompt: self.c.Tr.ConfirmQuitDuringUpdate,
		HandleConfirm: func() error {
			return gocui.ErrQuit
		},
	})
}

func (self *QuitActions) Escape() error {
	currentContext := self.c.CurrentContext()

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

	parentContext, hasParent := currentContext.GetParentContext()
	if hasParent && currentContext != nil && parentContext != nil {
		// TODO: think about whether this should be marked as a return rather than adding to the stack
		return self.c.PushContext(parentContext)
	}

	for _, mode := range self.c.Helpers().Mode.Statuses() {
		if mode.IsActive() {
			return mode.Reset()
		}
	}

	repoPathStack := self.c.State().GetRepoPathStack()
	if !repoPathStack.IsEmpty() {
		return self.c.Helpers().Repos.DispatchSwitchToRepo(repoPathStack.Pop(), true)
	}

	if self.c.UserConfig.QuitOnTopLevelReturn {
		return self.Quit()
	}

	return nil
}
