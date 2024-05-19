package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type TextboxController struct {
	baseController
	c *ControllerCommon
}

var _ types.IController = &TextboxController{}

func NewTextboxController(
	c *ControllerCommon,
) *TextboxController {
	return &TextboxController{
		baseController: baseController{},
		c:              c,
	}
}

func (self *TextboxController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:             opts.GetKey(opts.Config.Universal.ConfirmInEditor),
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
	}
	return bindings
}

func (self *TextboxController) GetOnFocusLost() func(types.OnFocusLostOpts) error {
	return func(types.OnFocusLostOpts) error {
		self.c.Helpers().Textbox.DeactivateTextboxPrompt()
		return nil
	}
}

func (self *TextboxController) Context() types.Context {
	return self.context()
}

func (self *TextboxController) context() *context.TextboxContext {
	return self.c.Contexts().Textbox
}
