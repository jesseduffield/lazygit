package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type ConfirmationController struct {
	baseController
	c *ControllerCommon
}

var _ types.IController = &ConfirmationController{}

func NewConfirmationController(
	common *ControllerCommon,
) *ConfirmationController {
	return &ConfirmationController{
		baseController: baseController{},
		c:              common,
	}
}

func (self *ConfirmationController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:         opts.GetKey(opts.Config.Universal.Confirm),
			Handler:     func() error { return self.context().State.OnConfirm() },
			Description: self.c.Tr.Confirm,
			Display:     true,
		},
		{
			Key:         opts.GetKey(opts.Config.Universal.Return),
			Handler:     func() error { return self.context().State.OnClose() },
			Description: self.c.Tr.CloseCancel,
			Display:     true,
		},
		{
			Key: opts.GetKey(opts.Config.Universal.TogglePanel),
			Handler: func() error {
				if len(self.c.Contexts().Suggestions.State.Suggestions) > 0 {
					return self.c.ReplaceContext(self.c.Contexts().Suggestions)
				}
				return nil
			},
		},
	}

	return bindings
}

func (self *ConfirmationController) GetOnFocusLost() func(types.OnFocusLostOpts) error {
	return func(types.OnFocusLostOpts) error {
		self.c.Helpers().Confirmation.DeactivateConfirmationPrompt()
		return nil
	}
}

func (self *ConfirmationController) Context() types.Context {
	return self.context()
}

func (self *ConfirmationController) context() *context.ConfirmationContext {
	return self.c.Contexts().Confirmation
}
