package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/context"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type SuggestionsController struct {
	baseController
	*controllerCommon
}

var _ types.IController = &SuggestionsController{}

func NewSuggestionsController(
	common *controllerCommon,
) *SuggestionsController {
	return &SuggestionsController{
		baseController:   baseController{},
		controllerCommon: common,
	}
}

func (self *SuggestionsController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{}

	return bindings
}

func (self *SuggestionsController) GetOnFocusLost() func(types.OnFocusLostOpts) error {
	return func(types.OnFocusLostOpts) error {
		deactivateConfirmationPrompt
		return nil
	}
}

func (self *SuggestionsController) Context() types.Context {
	return self.context()
}

func (self *SuggestionsController) context() *context.SuggestionsContext {
	return self.contexts.Suggestions
}
