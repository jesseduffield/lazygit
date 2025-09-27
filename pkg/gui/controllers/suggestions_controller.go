package controllers

import (
	"github.com/jesseduffield/gocui"
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
		ListControllerTrait: NewListControllerTrait(
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
			Key:               opts.GetKey(opts.Config.Universal.ConfirmSuggestion),
			Handler:           func() error { return self.context().State.OnConfirm() },
			GetDisabledReason: self.require(self.singleItemSelected()),
		},
		{
			Key:     opts.GetKey(opts.Config.Universal.Return),
			Handler: func() error { return self.context().State.OnClose() },
		},
		{
			Key:     opts.GetKey(opts.Config.Universal.TogglePanel),
			Handler: self.switchToPrompt,
		},
		{
			Key: opts.GetKey(opts.Config.Universal.Remove),
			Handler: func() error {
				return self.context().State.OnDeleteSuggestion()
			},
		},
		{
			Key: opts.GetKey(opts.Config.Universal.Edit),
			Handler: func() error {
				if self.context().State.AllowEditSuggestion {
					if selectedItem := self.c.Contexts().Suggestions.GetSelected(); selectedItem != nil {
						self.c.Contexts().Prompt.GetView().TextArea.Clear()
						self.c.Contexts().Prompt.GetView().TextArea.TypeString(selectedItem.Value)
						self.c.Contexts().Prompt.GetView().RenderTextArea()
						self.c.Contexts().Suggestions.RefreshSuggestions()
						return self.switchToPrompt()
					}
				}
				return nil
			},
		},
	}

	return bindings
}

func (self *SuggestionsController) GetMouseKeybindings(opts types.KeybindingsOpts) []*gocui.ViewMouseBinding {
	return []*gocui.ViewMouseBinding{
		{
			ViewName:    self.c.Contexts().Prompt.GetViewName(),
			FocusedView: self.c.Contexts().Suggestions.GetViewName(),
			Key:         gocui.MouseLeft,
			Handler: func(gocui.ViewMouseBindingOpts) error {
				return self.switchToPrompt()
			},
		},
	}
}

func (self *SuggestionsController) switchToPrompt() error {
	self.c.Views().Suggestions.Subtitle = ""
	self.c.Views().Suggestions.Highlight = false
	self.c.Context().Replace(self.c.Contexts().Prompt)
	return nil
}

func (self *SuggestionsController) GetOnFocusLost() func(types.OnFocusLostOpts) {
	return func(types.OnFocusLostOpts) {
		self.c.Helpers().Confirmation.DeactivatePrompt()
	}
}

func (self *SuggestionsController) context() *context.SuggestionsContext {
	return self.c.Contexts().Suggestions
}
