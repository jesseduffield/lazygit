package controllers

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type ConfirmationController struct {
	baseController
	c *ControllerCommon
}

var _ types.IController = &ConfirmationController{}

func NewConfirmationController(
	c *ControllerCommon,
) *ConfirmationController {
	return &ConfirmationController{
		baseController: baseController{},
		c:              c,
	}
}

func (self *ConfirmationController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:             opts.GetKey(opts.Config.Universal.Confirm),
			Handler:         func() error { return self.context().State.OnConfirm() },
			Description:     self.c.Tr.Confirm,
			DisplayOnScreen: true,
		},
		{
			Key:             opts.GetKey(opts.Config.Universal.Return),
			Handler:         func() error { return self.context().State.OnClose() },
			Description:     self.c.Tr.CloseCancel,
			DisplayOnScreen: true,
		},
		{
			Key: opts.GetKey(opts.Config.Universal.TogglePanel),
			Handler: func() error {
				if len(self.c.Contexts().Suggestions.State.Suggestions) > 0 {
					subtitle := ""
					if self.c.State().GetRepoState().GetCurrentPopupOpts().HandleDeleteSuggestion != nil {
						// We assume that whenever things are deletable, they
						// are also editable, so we show both keybindings
						subtitle = fmt.Sprintf(self.c.Tr.SuggestionsSubtitle,
							self.c.UserConfig().Keybinding.Universal.Remove, self.c.UserConfig().Keybinding.Universal.Edit)
					}
					self.c.Views().Suggestions.Subtitle = subtitle
					self.c.Context().Replace(self.c.Contexts().Suggestions)
				}
				return nil
			},
		},
	}

	return bindings
}

func (self *ConfirmationController) GetOnFocusLost() func(types.OnFocusLostOpts) {
	return func(types.OnFocusLostOpts) {
		self.c.Helpers().Confirmation.DeactivateConfirmationPrompt()
	}
}

func (self *ConfirmationController) Context() types.Context {
	return self.context()
}

func (self *ConfirmationController) context() *context.ConfirmationContext {
	return self.c.Contexts().Confirmation
}
