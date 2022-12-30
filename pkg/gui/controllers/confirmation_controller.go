package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type ConfirmationController struct {
	baseController
	*controllerCommon
}

var _ types.IController = &ConfirmationController{}

func NewConfirmationController(
	common *controllerCommon,
) *ConfirmationController {
	return &ConfirmationController{
		baseController:   baseController{},
		controllerCommon: common,
	}
}

func (self *ConfirmationController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{}

	return bindings
}

func (self *ConfirmationController) GetOnFocusLost() func(types.OnFocusLostOpts) error {
	return func(types.OnFocusLostOpts) error {
		deactivateConfirmationPrompt(self.controllerCommon)
		return nil
	}
}

func (self *ConfirmationController) Context() types.Context {
	return self.context()
}

func (self *ConfirmationController) context() types.Context {
	return self.contexts.Confirmation
}

func deactivateConfirmationPrompt(c *controllerCommon) {
	c.mutexes.PopupMutex.Lock()
	c.c.State().GetRepoState().SetCurrentPopupOpts(nil)
	c.mutexes.PopupMutex.Unlock()

	c.c.Views().Confirmation.Visible = false
	c.c.Views().Suggestions.Visible = false

	gui.clearConfirmationViewKeyBindings()
}
