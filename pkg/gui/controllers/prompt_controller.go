package controllers

import (
	"fmt"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type PromptController struct {
	baseController
	c *ControllerCommon
}

var _ types.IController = &PromptController{}

func NewPromptController(
	c *ControllerCommon,
) *PromptController {
	return &PromptController{
		baseController: baseController{},
		c:              c,
	}
}

func (self *PromptController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:             gocui.KeyEnter,
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
					self.switchToSuggestions()
				}
				return nil
			},
		},
	}

	return bindings
}

func (self *PromptController) GetMouseKeybindings(opts types.KeybindingsOpts) []*gocui.ViewMouseBinding {
	return []*gocui.ViewMouseBinding{
		{
			ViewName:    self.c.Contexts().Suggestions.GetViewName(),
			FocusedView: self.c.Contexts().Prompt.GetViewName(),
			Key:         gocui.MouseLeft,
			Handler: func(gocui.ViewMouseBindingOpts) error {
				self.switchToSuggestions()
				// Let it fall through to the ListController's click handler so that
				// the clicked line gets selected:
				return gocui.ErrKeybindingNotHandled
			},
		},
	}
}

func (self *PromptController) GetOnFocusLost() func(types.OnFocusLostOpts) {
	return func(types.OnFocusLostOpts) {
		self.c.Helpers().Confirmation.DeactivatePrompt()
	}
}

func (self *PromptController) Context() types.Context {
	return self.context()
}

func (self *PromptController) context() *context.PromptContext {
	return self.c.Contexts().Prompt
}

func (self *PromptController) switchToSuggestions() {
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
