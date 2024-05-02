package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SuggestionsController struct {
	baseController
	*ListControllerTrait[*types.Suggestion]
	c *ControllerCommon
}

var _ types.IController = &SuggestionsController{}

func NewSuggestionsController(
	c *ControllerCommon,
) *SuggestionsController {
	return &SuggestionsController{
		baseController: baseController{},
		ListControllerTrait: NewListControllerTrait[*types.Suggestion](
			c,
			c.Contexts().Suggestions,
			c.Contexts().Suggestions.GetSelected,
			c.Contexts().Suggestions.GetSelectedItems,
		),
		c: c,
	}
}

func (self *SuggestionsController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:               opts.GetKey(opts.Config.Universal.Confirm),
			Handler:           func() error { return self.context().State.OnConfirm() },
			GetDisabledReason: self.require(self.singleItemSelected()),
		},
		{
			Key:     opts.GetKey(opts.Config.Universal.Return),
			Handler: func() error { return self.context().State.OnClose() },
		},
		{
			Key:     opts.GetKey(opts.Config.Universal.TogglePanel),
			Handler: func() error { return self.c.ReplaceContext(self.c.Contexts().Confirmation) },
		},
	}

	return bindings
}

func (self *SuggestionsController) GetOnFocusLost() func(types.OnFocusLostOpts) error {
	return func(types.OnFocusLostOpts) error {
		self.c.Helpers().Confirmation.DeactivateConfirmationPrompt()
		return nil
	}
}

func (self *SuggestionsController) context() *context.SuggestionsContext {
	return self.c.Contexts().Suggestions
}
